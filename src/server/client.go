package server

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"holodeck1/config"
	"holodeck1/logging"
	"holodeck1/sync"
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
	hub            *Hub
	conn           *websocket.Conn
	send           chan []byte
	info           *ClientInfo
	lastSeen       time.Time
	hd1ID          string  // Single unified identifier - SINGLE SOURCE OF TRUTH
	avatarCreated  bool    // Track if avatar has been created for this client
	syncChan       chan *sync.Operation  // Sync system channel - SINGLE SOURCE OF TRUTH
}

// generateHD1ID generates a unified HD1 identifier
func generateHD1ID() string {
	return fmt.Sprintf("hd1-%d-%d", time.Now().Unix(), rand.Intn(100000))
}

// GetHD1ID returns the unified HD1 identifier
func (c *Client) GetHD1ID() string {
	if c.hd1ID == "" {
		c.hd1ID = generateHD1ID()
	}
	return c.hd1ID
}

// Legacy compatibility methods - maintain avatar creation tracking
func (c *Client) GetClientID() string { return c.GetHD1ID() }
func (c *Client) GetSessionID() string { return c.GetHD1ID() }
func (c *Client) GetAvatarID() string { 
	if c.avatarCreated {
		return c.GetHD1ID()
	}
	return ""
}
func (c *Client) SetSessionID(id string) { c.hd1ID = id }
func (c *Client) SetAvatarID(id string) { 
	c.hd1ID = id
	c.avatarCreated = true
}

// ensureRegistered ensures the client is registered with the hub (lazy registration)
func (c *Client) ensureRegistered() {
	// Check if client is already registered by checking if it's in the hub's clients map
	c.hub.mutex.RLock()
	_, isRegistered := c.hub.clients[c]
	c.hub.mutex.RUnlock()
	
	if !isRegistered {
		// Register client and send client_init message
		c.hub.register <- c
		
		// Send client ID to browser
		clientID := c.GetClientID()
		initMessage := map[string]interface{}{
			"type":    "client_init",
			"hd1_id":  clientID,
			"message": "HD1 ID assigned by server",
		}
		
		if initData, err := json.Marshal(initMessage); err == nil {
			select {
			case c.send <- initData:
				logging.Info("late client ID sent to browser", map[string]interface{}{
					"hd1_id": clientID,
				})
			default:
				logging.Error("failed to send late client ID to browser", map[string]interface{}{
					"hd1_id": clientID,
					"error":   "send channel blocked",
				})
			}
		}
	}
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
		c.lastSeen = time.Now()
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
		
		// Update last seen time for any message activity
		c.lastSeen = time.Now()
		
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
		// Not JSON, skip invalid message
		return
	}
	
	msgType, ok := msg["type"].(string)
	if !ok {
		// No type field, skip invalid message
		return
	}
	
	switch msgType {
	case "client_reconnect":
		// Handle client reconnection with existing client ID
		if existingClientID, ok := msg["hd1_id"].(string); ok {
			// Try to reconnect to existing avatar
			if avatar := c.hub.avatarRegistry.ReconnectClient(existingClientID, c); avatar != nil {
				// Set client ID to the existing one
				c.hd1ID = existingClientID
				
				// Register client with hub (since we skipped it in ServeWS)
				c.hub.register <- c
				
				// Pure in-memory architecture - no session persistence needed
				
				// Send confirmation back to client
				confirmMsg := map[string]interface{}{
					"type":      "client_reconnect_success",
					"hd1_id":    existingClientID,
					"avatar_id": avatar.ID,
					"message":   "Reconnected to existing avatar",
				}
				if jsonData, err := json.Marshal(confirmMsg); err == nil {
					select {
					case c.send <- jsonData:
						logging.Info("client reconnection confirmed", map[string]interface{}{
							"hd1_id":    existingClientID,
							"avatar_id": avatar.ID,
						})
					default:
						// Client Go channel blocked, don't wait
					}
				}
				return // Don't broadcast this message
			} else {
				logging.Info("client reconnection failed, creating new identity", map[string]interface{}{
					"requested_hd1_id": existingClientID,
				})
				// Avatar not found, client will get new client_init message
			}
		}
		
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
		// Client logging disabled for minimal build
		
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
		// Legacy session association - eliminated for unified HD1 ID system
		// Client already has proper hd1_id, no need to override it
		logging.Info("legacy session_associate ignored - using unified HD1 ID", map[string]interface{}{
			"hd1_id": c.GetHD1ID(),
		})
		
	case "interaction":
		c.lastSeen = time.Now()
		var interaction map[string]interface{}
		if err := json.Unmarshal(message, &interaction); err == nil {
			logging.Debug("user interaction", interaction)
		}
		// Interaction messages - handled locally, no sync needed
		
	case "avatar_asset_request":
		// Avatar asset requests not used in minimal build
		
	default:
		// Ensure client is registered if not already (for first non-reconnect message)
		c.ensureRegistered()
		
		// Regular 3D visualization message - REMOVED: Using sync system directly
	}
}

// forwardSyncOperations listens to sync channel and forwards operations to WebSocket
func (c *Client) forwardSyncOperations() {
	for operation := range c.syncChan {
		// Convert sync operation to WebSocket message
		message := map[string]interface{}{
			"type":      "sync_operation",
			"operation": operation,
		}
		
		if messageData, err := json.Marshal(message); err == nil {
			select {
			case c.send <- messageData:
				logging.Trace("websocket", "sync operation forwarded to client", map[string]interface{}{
					"hd1_id":  c.GetClientID(),
					"seq_num": operation.SeqNum,
					"op_type": operation.Type,
				})
			default:
				logging.Error("sync operation dropped - client send channel blocked", map[string]interface{}{
					"hd1_id":  c.GetClientID(),
					"seq_num": operation.SeqNum,
					"op_type": operation.Type,
				})
			}
		}
	}
}

// sendInitialSync sends existing operations to newly connected client
func (c *Client) sendInitialSync() {
	// Get all operations from sequence 1 to current
	currentSeq := c.hub.sync.GetCurrentSequence()
	if currentSeq > 0 {
		missingOps := c.hub.sync.GetMissingOperations(1, currentSeq)
		
		logging.Info("sending initial sync to client", map[string]interface{}{
			"hd1_id":     c.GetClientID(),
			"operations": len(missingOps),
			"from_seq":   1,
			"to_seq":     currentSeq,
		})
		
		for _, op := range missingOps {
			// Send each operation via sync channel (will be forwarded by forwardSyncOperations)
			select {
			case c.syncChan <- op:
				// Operation sent successfully
			default:
				logging.Error("initial sync operation dropped - sync channel blocked", map[string]interface{}{
					"hd1_id":  c.GetClientID(),
					"seq_num": op.SeqNum,
					"op_type": op.Type,
				})
			}
		}
	}
}

// Avatar asset handling removed for minimal build

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
	
	client := &Client{
		hub:  hub, 
		conn: conn, 
		send: make(chan []byte, config.GetWebSocketClientWorldBuffer()),
	}
	
	// Generate client ID immediately
	clientID := client.GetClientID()
	
	// Send client ID to browser for unified identification
	initMessage := map[string]interface{}{
		"type":    "client_init",
		"hd1_id":  clientID,
		"message": "HD1 ID assigned by server",
	}
	
	if initData, err := json.Marshal(initMessage); err == nil {
		select {
		case client.send <- initData:
			logging.Info("client ID sent to browser", map[string]interface{}{
				"hd1_id": clientID,
			})
		default:
			logging.Error("failed to send client ID to browser", map[string]interface{}{
				"hd1_id": clientID,
				"error":   "send channel blocked",
			})
		}
	} else {
		logging.Error("failed to marshal client init message", map[string]interface{}{
			"hd1_id": clientID,
			"error":   err.Error(),
		})
	}
	
	// Register client immediately - SINGLE SOURCE OF TRUTH
	hub.register <- client
	
	go client.writePump()
	go client.readPump()
}