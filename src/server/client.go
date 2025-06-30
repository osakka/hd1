package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"holodeck1/logging"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

type ClientInfo struct {
	Screen struct {
		Width            int     `json:"width"`
		Height           int     `json:"height"`
		DevicePixelRatio float64 `json:"devicePixelRatio"`
		Orientation      int     `json:"orientation"`
	} `json:"screen"`
	Canvas struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"canvas"`
	Capabilities struct {
		WebGL  bool `json:"webgl"`
		Touch  bool `json:"touch"`
		Mobile bool `json:"mobile"`
	} `json:"capabilities"`
}

type Client struct {
	hub       *Hub
	conn      *websocket.Conn
	send      chan []byte
	info      *ClientInfo
	lastSeen  time.Time
	sessionID string  // HD1 session isolation
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logging.Error("websocket connection error", map[string]interface{}{
					"error": err.Error(),
				})
			}
			break
		}
		
		// Handle special client messages
		c.handleClientMessage(message)
	}
}

func (c *Client) handleClientMessage(message []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		// Not JSON, broadcast as regular message
		c.hub.broadcast <- message
		return
	}
	
	msgType, ok := msg["type"].(string)
	if !ok {
		c.hub.broadcast <- message
		return
	}
	
	switch msgType {
	case "version_check":
		clientVersion, _ := msg["js_version"].(string)
		serverVersion := GetJSVersion()
		
		// Log version info and trigger reloads when versions don't match
		logging.Info("client version check", map[string]interface{}{
			"client": clientVersion,
			"server": serverVersion,
			"match": clientVersion == serverVersion,
		})
		
		// Send version mismatch response to trigger browser refresh
		if clientVersion != serverVersion {
			versionMismatchMsg := map[string]interface{}{
				"type": "version_mismatch",
				"server_version": serverVersion,
				"client_version": clientVersion,
			}
			if jsonData, err := json.Marshal(versionMismatchMsg); err == nil {
				select {
				case c.send <- jsonData:
				default:
					// Client channel blocked, don't wait
				}
			}
		}
		
	case "client_log":
		var logMsg LogMessage
		if err := json.Unmarshal(message, &logMsg); err == nil {
			logging.Debug("client log message", map[string]interface{}{
				"level": logMsg.Level,
				"message": logMsg.Message,
				"source": "browser",
			})
		}
		
	case "client_info":
		var info ClientInfo
		if err := json.Unmarshal(message, &info); err == nil {
			c.info = &info
			c.lastSeen = time.Now()
			
			logging.Info("client info updated", map[string]interface{}{
				"screen": info.Screen,
				"capabilities": info.Capabilities,
			})
		}

	case "session_associate":
		// Associate this client with a specific HD1 session
		if sessionID, ok := msg["session_id"].(string); ok {
			c.sessionID = sessionID
			logging.Info("client session associated", map[string]interface{}{
				"session_id": sessionID,
			})
			
			// SURGICAL FIX: Send existing session objects ONLY ONCE per session globally
			if c.hub.store != nil && !c.hub.IsSessionRestored(sessionID) {
				existingObjects := c.hub.store.ListObjects(sessionID)
				if len(existingObjects) > 0 {
					// Format objects for canvas control message (EXACT same format as normal object creation)
					var objectsData []map[string]interface{}
					for _, obj := range existingObjects {
						// Parse color string back to RGBA object (if stored as JSON string)
						var colorObj map[string]interface{}
						if obj.Color != "" {
							// Try to parse as JSON first, fallback to default red if parsing fails
							var parsedColor map[string]interface{}
							if err := json.Unmarshal([]byte(obj.Color), &parsedColor); err == nil {
								colorObj = parsedColor
							} else {
								// Default red color if parsing fails
								colorObj = map[string]interface{}{"r": 1.0, "g": 0.2, "b": 0.2, "a": 1.0}
							}
						} else {
							// Default red color if no color stored
							colorObj = map[string]interface{}{"r": 1.0, "g": 0.2, "b": 0.2, "a": 1.0}
						}
						
						// Use EXACT same format as normal object creation (/opt/hd1/src/api/objects/create.go)
						objectData := map[string]interface{}{
							"id":   obj.Name,
							"name": obj.Name,
							"type": obj.Type,
							"transform": map[string]interface{}{
								"position": map[string]interface{}{
									"x": obj.X,
									"y": obj.Y,
									"z": obj.Z,
								},
								"scale": map[string]interface{}{
									"x": obj.Scale,
									"y": obj.Scale,
									"z": obj.Scale,
								},
								"rotation": map[string]interface{}{
									"x": 0,
									"y": 0,
									"z": 0,
								},
							},
							"color": colorObj,
							"material": map[string]interface{}{
								"shader":      "standard",
								"metalness":   0.1,
								"roughness":   0.7,
								"transparent": false,
							},
							"physics": map[string]interface{}{
								"enabled": false,
								"mass":    1.0,
								"type":    "static",
							},
							"lighting": map[string]interface{}{
								"castShadow":    true,
								"receiveShadow": true,
							},
							"visible":   true,
							"wireframe": false,
							"text":      "",
							"lightType": "",
							"intensity": 1,
						}
						objectsData = append(objectsData, objectData)
					}
					
					// Send session restoration message to this client only
					message := map[string]interface{}{
						"type": "canvas_control",
						"data": map[string]interface{}{
							"command": "create",
							"objects": objectsData,
						},
						"timestamp": time.Now().Unix(),
					}
					
					if jsonData, err := json.Marshal(message); err == nil {
						select {
						case c.send <- jsonData:
						default:
							// Client channel blocked, don't wait
						}
					}
					
					logging.Info("session objects restored to client", map[string]interface{}{
						"session_id": sessionID,
						"object_count": len(existingObjects),
					})
					
					// Mark this session as globally restored to prevent future restoration loops
					c.hub.MarkSessionRestored(sessionID)
				}
			}
		}
		
	case "interaction":
		c.lastSeen = time.Now()
		var interaction map[string]interface{}
		if err := json.Unmarshal(message, &interaction); err == nil {
			logging.Debug("user interaction", interaction)
		}
		// Broadcast interaction to other systems that might be listening
		c.hub.broadcast <- message
		
	default:
		// Regular 3D visualization message
		c.hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
			
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.Error("websocket upgrade failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	
	go client.writePump()
	go client.readPump()
}