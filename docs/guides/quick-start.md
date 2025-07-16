# HD1 Quick Start Guide

Get HD1 (Holodeck One) running in under 5 minutes.

## Prerequisites

- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **Git** - For cloning the repository
- **Modern Browser** - Chrome, Firefox, Safari, or Edge
- **Port 8080** - Available for HD1 server

## Installation

### 1. Clone Repository
```bash
git clone <repository-url>
cd hd1
```

### 2. Build HD1
```bash
cd src
make clean && make
```

### 3. Start Server
```bash
make start
```

### 4. Access Console
Open your browser to: **http://localhost:8080**

## Verification

### Check Server Status
```bash
# Verify server is running
curl http://localhost:8080/api/system/version

# Expected response:
{
  "api_version": "0.7.0",
  "js_version": "...",
  "build_timestamp": "...",
  "title": "HD1 (Holodeck One) Three.js API"
}
```

### Test WebSocket Connection
1. Open browser console (F12)
2. Navigate to **http://localhost:8080**
3. Look for WebSocket connection messages
4. Debug panel should show "Connected" status

## Basic Usage

### Using the Console
- **Debug Panel**: Click to expand/collapse system information
- **WebSocket Status**: Real-time connection monitoring
- **Rebootstrap**: Automatic recovery on connection failures
- **Version Sync**: Automatic page refresh on server updates

### API Access
```bash
# List available endpoints
curl http://localhost:8080/api/

# Create Three.js entity
curl -X POST http://localhost:8080/api/threejs/entities \
  -H "Content-Type: application/json" \
  -d '{"type": "cube", "position": {"x": 0, "y": 0, "z": 0}}'

# Get scene information
curl http://localhost:8080/api/threejs/scene
```

### WebSocket Communication
```javascript
// Connect to WebSocket
const ws = new WebSocket('ws://localhost:8080/ws');

// Send message
ws.send(JSON.stringify({
    type: 'session_associate',
    session_id: 'my_session'
}));

// Receive messages
ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('Received:', data);
};
```

## Configuration

### Environment Variables
```bash
# Server configuration
export HD1_HOST=0.0.0.0
export HD1_PORT=8080
export HD1_LOG_LEVEL=INFO

# Start with custom config
make start
```

### .env File
```bash
# Create .env file in project root
cat > .env << EOF
HD1_HOST=127.0.0.1
HD1_PORT=9090
HD1_LOG_LEVEL=DEBUG
HD1_DAEMON=false
EOF
```

## Development Mode

### Live Reload Development
```bash
# Terminal 1: Build with file watching
make dev

# Terminal 2: Test changes
curl http://localhost:8080/api/system/version
```

### View Logs
```bash
# Real-time logs
make logs

# Specific log level
HD1_LOG_LEVEL=DEBUG make start
```

## Troubleshooting

### Server Won't Start
```bash
# Check if port is in use
lsof -i :8080

# Stop existing HD1 instance
make stop

# Start with debug logging
HD1_LOG_LEVEL=DEBUG make start
```

### WebSocket Connection Issues
1. **Check browser console** for connection errors
2. **Verify firewall** allows port 8080
3. **Try different browser** to isolate issues
4. **Check server logs** for WebSocket errors

### Build Failures
```bash
# Clean and rebuild
make clean
make generate
make build

# Check Go version
go version  # Should be 1.21+

# Verify dependencies
go mod tidy
```

## Next Steps

### Learn the API
- **[API Reference](../architecture/api-design.md)** - Complete endpoint documentation
- **[WebSocket Protocol](../architecture/websocket.md)** - Real-time communication guide

### Development
- **[Development Guide](development.md)** - Extending HD1
- **[Configuration Guide](configuration.md)** - Advanced configuration

### Deployment
- **[Deployment Guide](deployment.md)** - Production deployment

## Examples

### Basic Three.js Scene
```javascript
// Create a simple scene with entities
fetch('/api/threejs/entities', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({
        type: 'cube',
        position: {x: 0, y: 1, z: 0},
        color: '#ff0000'
    })
});
```

### Real-time Updates
```javascript
// Listen for real-time entity updates
ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    
    switch(message.type) {
        case 'entity_created':
            console.log('New entity:', message.entity);
            break;
        case 'entity_updated':
            console.log('Entity changed:', message.entity);
            break;
    }
};
```

## Support

### Getting Help
- **Documentation**: Check [docs/](../) for comprehensive guides
- **Logs**: Use `make logs` for detailed error information
- **Issues**: Review troubleshooting section above

### Common Issues
- **Port conflicts**: Use different port with `HD1_PORT=9090`
- **Permission errors**: Check file permissions in build/ directory
- **WebSocket failures**: Verify no proxy/firewall blocking connections

---

*Quick Start Guide for HD1 v0.7.0 - Three.js Game Engine Platform*