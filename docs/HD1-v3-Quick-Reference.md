# HD1 v3.0 - Quick Reference Guide

**Version**: 3.0  
**Updated**: 2025-07-03  

## Quick Start

### Start the Server
```bash
cd /opt/hd1/src
make clean && make && make start
```

### Access the Interface
- **Web UI**: http://localhost:8080
- **API Documentation**: http://localhost:8080/api
- **Console**: Click header to expand/collapse (smooth animations!)

## Essential API Endpoints

### Session Management
```bash
# Create a new session
curl -X POST http://localhost:8080/api/sessions

# Join a channel
curl -X POST http://localhost:8080/api/sessions/{sessionId}/channel/join \
  -H "Content-Type: application/json" \
  -d '{"channel_id": "channel_one", "client_id": "client_test"}'

# List sessions
curl http://localhost:8080/api/sessions
```

### Entity Management
```bash
# Create entity
curl -X POST http://localhost:8080/api/sessions/{sessionId}/entities \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my_box",
    "components": {
      "model": {"type": "box"},
      "transform": {"position": [0, 1, 0], "scale": [1, 1, 1]},
      "material": {"diffuse": "#ff0000"}
    }
  }'

# List entities
curl http://localhost:8080/api/sessions/{sessionId}/entities

# Delete entity
curl -X DELETE http://localhost:8080/api/sessions/{sessionId}/entities/{entityId}
```

### Channel Information
```bash
# List available channels
curl http://localhost:8080/api/channels

# Get channel details
curl http://localhost:8080/api/channels/channel_one
```

## Available Channels

### Channel One - Red Box Scene
- **ID**: `channel_one`
- **Theme**: Primary collaborative environment
- **Objects**: Red box, floor, overhead lighting
- **Floor**: 10×1×10 dimensions, flat horizontal

### Channel Two - Blue Sphere Scene  
- **ID**: `channel_two`
- **Theme**: Physics simulation environment
- **Objects**: Blue sphere, floor, overhead lighting
- **Floor**: 10×1×10 dimensions, flat horizontal

### Channel Three - Green Pyramid Scene
- **ID**: `channel_three`
- **Theme**: Exploration environment
- **Objects**: Green pyramid (cone), floor, overhead lighting
- **Floor**: 10×1×10 dimensions, flat horizontal

## Console Features

### Performance Monitoring
- **CPU Usage**: Real-time system CPU monitoring
- **Memory Usage**: RAM consumption tracking
- **WebSocket Status**: Connection health and message counts
- **Entity Count**: Live entity tracking per session

### Interactive Controls
- **Expand/Collapse**: Click header or status bar (smooth 0.4s animation)
- **Channel Switching**: Dropdown selector with live channel list
- **Session Management**: Auto-session creation and management
- **Version Display**: Real-time API and JS version tracking

### Stats Panels (Expanded Mode)
- **Top Left**: CPU usage graph with live updates
- **Top Right**: Memory usage visualization
- **Bottom Left**: WebSocket connection status and metrics
- **Bottom Right**: Entity count and session information

### Collapsed Mode
- **Compact Tiles**: 5 key metrics in minimal space
- **Smooth Transition**: Professional easing animations
- **Status Indicators**: Critical information always visible

## WebSocket Events

### Entity Lifecycle
```javascript
// Entity created
{
  "type": "entity_created",
  "data": {
    "session_id": "session-abc123",
    "entity_id": "entity-xyz789", 
    "name": "red_box",
    "components": { ... }
  }
}

// Entity deleted
{
  "type": "entity_deleted",
  "data": {"entity_id": "entity-xyz789"}
}

// Entity updated
{
  "type": "entity_updated", 
  "data": {"entity_id": "entity-xyz789", "components": { ... }}
}
```

## Entity Components

### Transform Component
```json
{
  "transform": {
    "position": [x, y, z],      // World position
    "rotation": [x, y, z],      // Euler angles in degrees
    "scale": [x, y, z]          // Local scale factors
  }
}
```

### Model Component
```json
{
  "model": {
    "type": "box|sphere|plane|cone|cylinder"
  }
}
```

### Material Component
```json
{
  "material": {
    "diffuse": "#ff0000",       // Hex color
    "metalness": 0.0,           // 0.0 to 1.0
    "roughness": 0.7,           // 0.0 to 1.0  
    "shininess": 50             // Specular intensity
  }
}
```

### Physics Component
```json
{
  "rigidbody": {
    "type": "static|dynamic|kinematic",
    "mass": 1.0                 // For dynamic bodies
  },
  "collision": {
    "type": "box|sphere|plane|cone"
  }
}
```

### Light Component
```json
{
  "light": {
    "type": "directional|point|spot|ambient",
    "color": [r, g, b],         // RGB 0.0 to 1.0
    "intensity": 3.0,           // Light intensity
    "castShadows": true,        // Shadow casting
    "enabled": true             // Light state
  }
}
```

## Development Commands

### Build System
```bash
# Full rebuild
make clean && make

# Start server
make start

# Stop server  
make stop

# Generate code from API spec
make generate

# View logs
tail -f /opt/hd1/build/logs/hd1.log
```

### CLI Client
```bash
# List sessions
./build/bin/hd1-client list-sessions

# Create session
./build/bin/hd1-client create-session

# List entities in session
curl http://localhost:8080/api/sessions/{sessionId}/entities
```

## Configuration Files

### Channel Configuration
```yaml
# /opt/hd1/share/channels/channel_one.yaml
channel:
  id: "channel_one"
  name: "Scene 1 - Red Box"
  max_clients: 100

playcanvas:
  scene:
    ambientLight: [0.3, 0.3, 0.3]
    gravity: [0, -9.81, 0]
  
  entities:
    - name: "floor"
      components:
        model: { type: "plane" }
        transform:
          position: [0, -0.1, 0]
          rotation: [0, 0, 0]     # Flat horizontal
          scale: [10, 1, 10]      # 10×1×10 floor
```

### API Specification
- **Location**: `/opt/hd1/src/api.yaml`
- **Format**: OpenAPI 3.0.3
- **Auto-generates**: All routing, clients, documentation

## Troubleshooting

### Common Issues

#### Console Shows "Loading"
- **Cause**: Server not responding or WebSocket disconnected
- **Solution**: Check server status with `make status` or restart with `make stop && make start`

#### No Graph Data in Console
- **Cause**: Console manager not initialized
- **Solution**: Refresh browser or check browser console for JavaScript errors

#### Channel Selector Empty
- **Cause**: API not accessible or channels not loaded
- **Solution**: Verify server running and channels directory exists: `ls /opt/hd1/share/channels/`

#### WebSocket Connection Failed
- **Cause**: Port 8080 blocked or server not running
- **Solution**: Check port with `netstat -tlnp | grep 8080` and firewall settings

### Logging Commands
```bash
# Increase log verbosity
curl -X POST http://localhost:8080/api/admin/logging/level \
  -d '{"level": "DEBUG"}'

# Enable module tracing
curl -X POST http://localhost:8080/api/admin/logging/trace \
  -d '{"modules": ["sessions", "entities", "websocket"]}'

# View structured logs
tail -f /opt/hd1/build/logs/hd1.log | jq .

# Return to normal logging
curl -X POST http://localhost:8080/api/admin/logging/level \
  -d '{"level": "INFO"}'
```

### Browser Debug
```javascript
// Check console manager status
window.hd1ConsoleManager?.getStatus()

// Force WebSocket reconnection
window.hd1ConsoleManager?.getModule('websocket')?.forceReconnect()

// Check API client
window.hd1API.getVersion()

// Force browser refresh
window.hd1API.forceRefresh({force: true})
```

## Performance Tips

### Optimal Settings
- **Entity Limit**: Keep <100 entities per session for best performance
- **Update Frequency**: Avoid rapid entity modifications (>10/second)
- **WebSocket Messages**: Monitor message rate in console stats

### Monitoring
- **Console Stats**: Real-time performance monitoring in expanded console
- **Browser DevTools**: Network tab for API performance
- **Server Logs**: Structured logging with performance metrics

## Next Steps

1. **Read Full Documentation**: [HD1-v3-Current-State-Architecture.md](HD1-v3-Current-State-Architecture.md)
2. **Architecture Details**: [ADR-019-Production-Ready-API-First-Game-Engine.md](adr/ADR-019-Production-Ready-API-First-Game-Engine.md)
3. **API Reference**: View OpenAPI spec at `/opt/hd1/src/api.yaml`
4. **Community**: Join development discussions and contribute improvements

---

**Happy Building with HD1 v3.0! 🎮**