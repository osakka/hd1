# ADR-004: WebSocket Synchronization Protocol

**Status**: Accepted  
**Date**: 2025-07-14  
**Authors**: HD1 Development Team  
**Related**: ADR-001 (API-First Architecture), ADR-003 (Three.js Minimal Console)

## Context

HD1 requires real-time synchronization between server state and multiple connected clients for multiplayer 3D environments. Traditional polling approaches introduce latency and scaling issues. The system needs to support high-frequency updates (avatar positions, entity changes) while maintaining connection reliability and efficient resource utilization.

## Decision

We implement a **TCP-Simple WebSocket Synchronization Protocol** that provides real-time bidirectional communication with intelligent reconnection, session isolation, and high-frequency update optimization.

### Protocol Characteristics

1. **Bidirectional**: Server-to-client and client-to-server messaging
2. **Session-Isolated**: Multi-tenant support with session-based message routing
3. **High-Frequency**: Optimized for rapid avatar position updates
4. **Reliable**: Automatic reconnection with exponential backoff
5. **Efficient**: Object pooling and zero-allocation hot paths

### Message Format

```javascript
// Standard message envelope
{
    "type": "message_type",
    "session_id": "session_identifier", 
    "timestamp": 1642518000000,
    "data": { /* payload */ }
}
```

## Implementation

### Server-Side Hub Architecture

```go
// WebSocket hub manages all client connections
type Hub struct {
    clients    map[*Client]bool           // Active client connections
    broadcast  chan []byte               // Broadcast channel
    register   chan *Client              // Client registration
    unregister chan *Client              // Client unregistration
    sessions   map[string][]*Client      // Session-based client grouping
}

// High-performance message broadcasting
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.registerClient(client)
            
        case client := <-h.unregister:
            h.unregisterClient(client)
            
        case message := <-h.broadcast:
            h.broadcastMessage(message)
        }
    }
}
```

### Client-Side Connection Management

```javascript
// Intelligent WebSocket client with automatic recovery
class WebSocketManager {
    constructor() {
        this.ws = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 99;
        this.connected = false;
        this.sessionId = null;
    }
    
    // Exponential backoff reconnection strategy
    reconnect() {
        const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
        
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            this.triggerRebootstrap(); // Nuclear option
            return;
        }
        
        setTimeout(() => this.connect(), delay);
        this.reconnectAttempts++;
    }
}
```

### Session Isolation

```go
// Session-based message routing
func (h *Hub) broadcastToSession(sessionID string, message []byte) {
    if clients, exists := h.sessions[sessionID]; exists {
        for _, client := range clients {
            select {
            case client.send <- message:
                // Message sent successfully
            default:
                // Client channel blocked - remove client
                h.unregisterClient(client)
            }
        }
    }
}
```

## Message Types

### Connection Management
```javascript
// Session association
{
    "type": "session_associate",
    "session_id": "world_one_session_123"
}

// Ping/pong for latency measurement
{
    "type": "ping",
    "ping_id": "ping_001",
    "timestamp": 1642518000000
}
```

### High-Frequency Updates
```javascript
// Avatar position updates (high frequency)
{
    "type": "avatar_position_update", 
    "session_id": "session_123",
    "avatar_id": "avatar_456",
    "position": {"x": 10.5, "y": 2.0, "z": -5.3},
    "rotation": {"x": 0, "y": 1.57, "z": 0}
}

// Entity lifecycle events
{
    "type": "entity_created",
    "session_id": "session_123", 
    "entity": {
        "id": "entity_789",
        "type": "cube",
        "position": {"x": 0, "y": 0, "z": 0}
    }
}
```

### System Messages
```javascript
// Version synchronization
{
    "type": "version_check",
    "js_version": "d63795c8-7ea10484-98025332-53f4cdfd"
}

// Client information
{
    "type": "client_info",
    "screen": {"width": 1920, "height": 1080},
    "capabilities": {"webgl": true, "touch": false}
}
```

## Performance Optimizations

### Memory Pooling
```go
// Object pooling for high-frequency messages
var messagePool = sync.Pool{
    New: func() interface{} {
        return &WebSocketMessage{}
    },
}

// Zero-allocation message handling
func handleHighFrequencyUpdate(data []byte) {
    msg := messagePool.Get().(*WebSocketMessage)
    defer messagePool.Put(msg)
    
    // Process message without allocation
}
```

### Connection Optimization
```go
// WebSocket configuration for performance
upgrader := websocket.Upgrader{
    ReadBufferSize:  4096,  // Optimized buffer sizes
    WriteBufferSize: 4096,
    CheckOrigin: func(r *http.Request) bool {
        return true // CORS handling
    },
}

// Connection limits and timeouts
const (
    writeWait      = 10 * time.Second  // Write timeout
    pongWait       = 60 * time.Second  // Pong wait timeout  
    pingPeriod     = 54 * time.Second  // Ping interval
    maxMessageSize = 512               // Maximum message size
)
```

### Broadcast Efficiency
```go
// Efficient message broadcasting with channel selection
func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer ticker.Stop()
    
    for {
        select {
        case message, ok := <-c.send:
            if !ok {
                c.ws.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            
            // Non-blocking write with timeout
            c.ws.SetWriteDeadline(time.Now().Add(writeWait))
            c.ws.WriteMessage(websocket.TextMessage, message)
            
        case <-ticker.C:
            // Keepalive ping
            c.ws.SetWriteDeadline(time.Now().Add(writeWait))
            c.ws.WriteMessage(websocket.PingMessage, nil)
        }
    }
}
```

## Reliability Features

### Automatic Reconnection
```javascript
// Intelligent reconnection with backoff
handleDisconnection() {
    this.connected = false;
    this.addDebug('WEBSOCKET', `Connection lost. Attempting reconnection ${this.reconnectAttempts + 1}/${this.maxReconnectAttempts}`);
    
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
        this.scheduleReconnection();
    } else {
        this.triggerRebootstrap(); // Ultimate recovery
    }
}
```

### Connection Health Monitoring
```javascript
// Ping/pong latency measurement
sendPing() {
    const pingId = `ping_${Date.now()}`;
    const message = {
        type: 'ping',
        ping_id: pingId,
        timestamp: Date.now()
    };
    
    this.ws.send(JSON.stringify(message));
    this.pendingPings.set(pingId, Date.now());
}

handlePong(data) {
    const latency = Date.now() - this.pendingPings.get(data.ping_id);
    this.updateLatencyStats(latency);
}
```

### Error Recovery
```javascript
// Progressive recovery strategies
handleError(error) {
    console.error('WebSocket error:', error);
    
    switch(error.type) {
        case 'connection_timeout':
            this.reconnect();
            break;
        case 'protocol_error':
            this.resetConnection();
            break;
        case 'persistent_failure':
            this.triggerRebootstrap();
            break;
    }
}
```

## Consequences

### Positive
- **Real-Time Performance**: Sub-100ms message latency for local networks
- **High Throughput**: Thousands of messages per second capability
- **Reliable Recovery**: Automatic reconnection handles network instability
- **Session Isolation**: Multi-tenant support with clean message routing
- **Resource Efficiency**: Object pooling minimizes garbage collection

### Negative
- **Complexity**: More complex than simple HTTP polling
- **Connection Management**: Requires careful WebSocket lifecycle handling
- **State Synchronization**: Client/server state consistency challenges
- **Network Dependency**: Real-time features degrade with poor connectivity

### Neutral
- **Protocol Overhead**: WebSocket framing adds minimal overhead vs raw TCP
- **Browser Support**: Modern WebSocket APIs required

## Monitoring and Metrics

### Connection Statistics
```javascript
// Real-time connection monitoring
const stats = {
    connections: activeConnectionCount,
    messagesPerSecond: messageRate,
    averageLatency: latencyStats.average,
    reconnectionRate: reconnectionEvents.rate
};
```

### Performance Metrics
- **Message Throughput**: Messages processed per second
- **Connection Stability**: Average connection duration
- **Latency Distribution**: P50, P90, P99 message latency
- **Memory Usage**: WebSocket buffer and object pool efficiency

## Future Enhancements

### Protocol Extensions
- **Message Compression**: Gzip compression for large messages
- **Priority Queuing**: Critical message prioritization
- **Batch Operations**: Bulk message processing
- **Delta Compression**: State difference encoding

### Advanced Features
- **Selective Sync**: Client-specified interest areas
- **Quality of Service**: Guaranteed delivery for critical messages
- **Load Balancing**: Multi-server WebSocket distribution
- **Offline Support**: Message queuing for disconnected clients

---

*ADR-004 establishes HD1's real-time WebSocket synchronization protocol*