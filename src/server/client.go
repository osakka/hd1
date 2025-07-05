package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
	"holodeck1/config"
	"holodeck1/logging"
)

// WebSocket configuration functions (using config system)
func getWriteWait() time.Duration {
	return config.GetWebSocketWriteTimeout()
}

func getPongWait() time.Duration {
	return config.GetWebSocketPongTimeout()
}

func getPingPeriod() time.Duration {
	return config.GetWebSocketPingPeriod()
}

func getMaxMessageSize() int64 {
	return config.GetWebSocketMaxMessageSize()
}

func getUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  config.GetWebSocketReadBufferSize(),
		WriteBufferSize: config.GetWebSocketWriteBufferSize(),
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for now
		},
	}
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

// readPump handles incoming WebSocket messages from the client.
// It runs in a separate goroutine and manages the client's read lifecycle:
// - Sets connection limits and deadlines for message size and pong timeouts
// - Processes incoming messages through handleClientMessage
// - Automatically unregisters client and closes connection on errors or disconnection
// - Handles graceful and unexpected connection closures
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	
	c.conn.SetReadLimit(getMaxMessageSize())
	c.conn.SetReadDeadline(time.Now().Add(getPongWait()))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(getPongWait()))
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

// handleClientMessage processes WebSocket messages received from clients.
// It handles three types of messages:
// 1. avatar_position_update: High-frequency avatar movement with direct position updates
// 2. session_change: Client requests to switch between HD1 worlds
// 3. Regular 3D visualization messages: Standard scene graph operations
// 
// Parameters:
//   message: Raw JSON message bytes from the WebSocket connection
//
// The function ensures avatar persistence during rapid updates and manages
// bidirectional session isolation for multiplayer synchronization.
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
					// Client Go channel blocked, don't wait
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
		
	case "ping":
		// Handle client ping for latency measurement
		pongMsg := map[string]interface{}{
			"type": "pong",
		}
		
		// Copy ping_id and timestamp for round-trip calculation
		if pingID, ok := msg["ping_id"]; ok {
			pongMsg["ping_id"] = pingID
		}
		if timestamp, ok := msg["timestamp"]; ok {
			pongMsg["timestamp"] = timestamp
		}
		
		// Send pong response immediately
		if jsonData, err := json.Marshal(pongMsg); err == nil {
			select {
			case c.send <- jsonData:
			default:
				// Client Go channel blocked, don't wait
			}
		}
		
		logging.Trace("websocket", "ping pong latency", map[string]interface{}{
			"ping_id": msg["ping_id"],
		})

	case "session_associate":
		// Associate this client with a specific HD1 session
		if sessionID, ok := msg["session_id"].(string); ok {
			c.sessionID = sessionID
			logging.Info("client session associated", map[string]interface{}{
				"session_id": sessionID,
			})
			
			if c.hub.store != nil {
				// Join the session room (this handles duplicate prevention)
				_, _, _ = c.hub.JoinSessionWorld(sessionID, fmt.Sprintf("%p", c), false)
				
				// Legacy object loading removed - entities now managed via worlds/PlayCanvas
				// Session restoration handled by world manager when switching worlds
				logging.Info("session connected, entities managed via worlds", map[string]interface{}{
					"session_id": sessionID,
				})
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
		
	case "avatar_asset_request":
		// Handle GLB asset request for avatars
		if avatarType, ok := msg["avatar_type"].(string); ok {
			c.handleAvatarAssetRequest(avatarType)
		}
		
	default:
		// Regular 3D visualization message
		c.hub.broadcast <- message
	}
}

// handleAvatarAssetRequest loads and sends GLB avatar assets via WebSocket.
// Reads the GLB file from the avatars directory and sends it as base64-encoded data.
// This maintains HD1's real-time architecture by delivering assets through WebSocket
// rather than HTTP, ensuring coordinated asset loading with the session state.
func (c *Client) handleAvatarAssetRequest(avatarType string) {
	// Construct path to GLB asset file
	avatarsDir := config.GetAvatarsDir()
	glbPath := filepath.Join(avatarsDir, avatarType, "model.glb")
	
	// Read GLB file data
	glbData, err := ioutil.ReadFile(glbPath)
	if err != nil {
		logging.Error("failed to read GLB asset", map[string]interface{}{
			"avatar_type": avatarType,
			"path": glbPath,
			"error": err.Error(),
		})
		
		// Send error response
		errorResponse := map[string]interface{}{
			"type": "avatar_asset_error",
			"avatar_type": avatarType,
			"error": "Failed to load avatar asset",
		}
		if jsonData, err := json.Marshal(errorResponse); err == nil {
			select {
			case c.send <- jsonData:
			default:
				// Client Go channel blocked, don't wait
			}
		}
		return
	}
	
	// Encode GLB data as base64
	base64Data := base64.StdEncoding.EncodeToString(glbData)
	
	// Send asset response via WebSocket
	assetResponse := map[string]interface{}{
		"type": "avatar_asset_response",
		"avatar_type": avatarType,
		"data": base64Data,
		"size": len(glbData),
	}
	
	if jsonData, err := json.Marshal(assetResponse); err == nil {
		select {
		case c.send <- jsonData:
			logging.Info("GLB asset sent via WebSocket", map[string]interface{}{
				"avatar_type": avatarType,
				"size_bytes": len(glbData),
			})
		default:
			logging.Warn("failed to send GLB asset - client Go channel blocked", map[string]interface{}{
				"avatar_type": avatarType,
			})
		}
	} else {
		logging.Error("failed to marshal GLB asset response", map[string]interface{}{
			"avatar_type": avatarType,
			"error": err.Error(),
		})
	}
}

// writePump handles outgoing WebSocket messages to the client.
// It runs in a separate goroutine and manages the client's write lifecycle:
// - Sends queued messages from the client's send Go channel
// - Implements ping/pong keepalive mechanism with configurable intervals
// - Manages write deadlines to prevent connection hangs
// - Gracefully handles Go channel closure and connection errors
// - Automatically closes connection when send Go channel is closed
func (c *Client) writePump() {
	ticker := time.NewTicker(getPingPeriod())
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(getWriteWait()))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
			
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(getWriteWait()))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	upgrader := getUpgrader()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.Error("websocket upgrade failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, config.GetWebSocketClientWorldBuffer())}
	client.hub.register <- client
	
	go client.writePump()
	go client.readPump()
}