# HD1 Development Guide

Complete guide for developing with and extending HD1.

## Development Workflow

### 1. Specification-Driven Development

HD1 follows specification-first development where all changes begin with the OpenAPI specification:

```bash
# 1. Edit the specification
vim src/api.yaml

# 2. Generate code from specification  
make generate

# 3. Implement business logic in handlers
vim src/api/threejs/your_handler.go

# 4. Test the implementation
make build && make start
```

### 2. Adding New API Endpoints

#### Step 1: Define in Specification
```yaml
# src/api.yaml
paths:
  /threejs/meshes:
    post:
      operationId: createMesh
      tags: [threejs]
      summary: Create a new Three.js mesh
      x-handler: api/threejs/mesh.go
      x-function: CreateMesh
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                geometry:
                  type: string
                material:
                  type: string
      responses:
        201:
          description: Mesh created successfully
```

#### Step 2: Generate Code
```bash
make generate
# This creates routing in auto_router.go and updates client libraries
```

#### Step 3: Implement Handler
```go
// src/api/threejs/mesh.go
package threejs

import (
    "encoding/json"
    "net/http"
    "holodeck1/logging"
    "holodeck1/server"
)

// CreateMesh creates a new Three.js mesh entity
func CreateMesh(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
    var request struct {
        Geometry string `json:"geometry"`
        Material string `json:"material"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // Business logic here
    mesh := createMeshEntity(request.Geometry, request.Material)
    
    // Broadcast to WebSocket clients
    broadcastMessage := map[string]interface{}{
        "type": "mesh_created",
        "mesh": mesh,
    }
    
    if jsonData, err := json.Marshal(broadcastMessage); err == nil {
        hub.Broadcast(jsonData)
    }
    
    // Return response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(mesh)
    
    logging.Info("mesh created", map[string]interface{}{
        "geometry": request.Geometry,
        "material": request.Material,
    })
}
```

### 3. Development Commands

```bash
# Full development cycle
make clean && make && make start

# Individual steps
make validate     # Validate api.yaml specification
make generate     # Generate code from specification
make build        # Build server binary
make start        # Start development server
make stop         # Stop running server
make logs         # View server logs

# Development with auto-reload
make dev          # Watch for changes and rebuild
```

## Project Structure

### Core Directories
```
src/
├── api.yaml              # OpenAPI specification (source of truth)
├── main.go               # Server entry point
├── auto_router.go        # Generated routing (DO NOT EDIT)
├── Makefile              # Build automation
├── api/                  # API handler implementations
│   ├── threejs/          # Three.js-specific endpoints
│   ├── sync/             # Synchronization endpoints
│   └── system/           # System information endpoints
├── codegen/              # Code generation from specification
│   ├── generator.go      # Main generator logic
│   └── templates/        # Code generation templates
├── config/               # Configuration management
├── logging/              # Structured logging system
└── server/               # WebSocket hub and client management
```

### Generated Files (Never Edit)
- `auto_router.go` - HTTP routing generated from api.yaml
- `share/htdocs/static/js/hd1lib.js` - JavaScript client library

### Configuration Files
- `src/api.yaml` - OpenAPI specification
- `.env` - Environment variable overrides
- `Makefile` - Build and development commands

## Code Generation

### How It Works
HD1 uses a custom code generator that reads `src/api.yaml` and generates:

1. **HTTP Routing** (`auto_router.go`)
2. **JavaScript Client** (`hd1lib.js`) 
3. **Validation Middleware** (embedded in router)

### Generator Configuration
```yaml
# src/api.yaml
x-code-generation:
  strict-validation: true        # Fail build on validation errors
  auto-routing: true            # Generate HTTP routing
  handler-validation: true      # Check handler files exist
  fail-on-missing-handlers: true # Fail if handlers missing
```

### Template System
```
src/codegen/templates/
├── go/
│   └── router.tmpl           # Go HTTP router template
└── javascript/
    └── threejs-client.tmpl   # JavaScript API client template
```

### Custom Templates
```go
// Add custom generation logic
func generateCustomCode(spec OpenAPISpec, routes []RouteInfo) {
    tmpl, err := loadTemplate("templates/custom/my-template.tmpl")
    if err != nil {
        return err
    }
    
    // Execute template with data
    return tmpl.Execute(outputFile, templateData)
}
```

## WebSocket Development

### Adding WebSocket Message Types
```javascript
// Client-side message handling
switch(message.type) {
    case 'my_custom_message':
        handleCustomMessage(message.data);
        break;
}

function handleCustomMessage(data) {
    console.log('Custom message received:', data);
}
```

```go
// Server-side message handling in server/client.go
case "my_custom_message":
    // Handle custom message type
    handleCustomMessage(msg)
```

### Broadcasting Messages
```go
// From any API handler
func broadcastToClients(hub *server.Hub, messageType string, data interface{}) {
    message := map[string]interface{}{
        "type": messageType,
        "data": data,
        "timestamp": time.Now().Unix(),
    }
    
    if jsonData, err := json.Marshal(message); err == nil {
        hub.Broadcast(jsonData)
    }
}
```

## Configuration Management

### Environment Variables
```bash
# All HD1 configuration uses HD1_ prefix
export HD1_HOST=0.0.0.0
export HD1_PORT=8080
export HD1_LOG_LEVEL=DEBUG
export HD1_STATIC_DIR=/opt/hd1/share/htdocs/static
```

### Configuration Priority
1. **Command-line flags** (highest priority)
2. **Environment variables** 
3. **.env file**
4. **Default values** (lowest priority)

### Adding New Configuration
```go
// 1. Add to config/config.go
type Config struct {
    MyNewSetting string `yaml:"my_new_setting"`
}

// 2. Add environment variable support
func loadFromEnvironment() {
    if val := os.Getenv("HD1_MY_NEW_SETTING"); val != "" {
        Config.MyNewSetting = val
    }
}

// 3. Add command-line flag
var myFlag = flag.String("my-setting", "", "Description of setting")

// 4. Use in application
settingValue := config.GetMyNewSetting()
```

## Logging and Debugging

### Structured Logging
```go
// Use structured logging throughout
logging.Info("operation completed", map[string]interface{}{
    "operation": "create_entity",
    "entity_id": entityID,
    "duration_ms": duration.Milliseconds(),
})

logging.Error("operation failed", map[string]interface{}{
    "operation": "create_entity",
    "error": err.Error(),
    "request_id": requestID,
})
```

### Log Levels
```bash
# Development
HD1_LOG_LEVEL=DEBUG make start

# Production  
HD1_LOG_LEVEL=INFO make start

# Troubleshooting
HD1_LOG_LEVEL=TRACE HD1_TRACE_MODULES=websocket,entities make start
```

### Debug Console
```javascript
// Client-side debugging
console.log('HD1 Debug:', {
    websocketStatus: this.connected,
    reconnectAttempts: this.reconnectAttempts,
    sessionId: this.sessionId
});
```

## Testing

### API Testing
```bash
# Test new endpoint
curl -X POST http://localhost:8080/api/threejs/meshes \
  -H "Content-Type: application/json" \
  -d '{"geometry": "box", "material": "basic"}'

# Test WebSocket
wscat -c ws://localhost:8080/ws
> {"type": "session_associate", "session_id": "test"}
```

### Automated Testing
```go
// Example handler test
func TestCreateMesh(t *testing.T) {
    hub := server.NewHub()
    req := httptest.NewRequest("POST", "/api/threejs/meshes", strings.NewReader(`{"geometry":"box"}`))
    w := httptest.NewRecorder()
    
    CreateMesh(w, req, hub)
    
    assert.Equal(t, http.StatusCreated, w.Code)
}
```

## Performance Optimization

### Memory Management
```go
// Use object pooling for high-frequency operations
var messagePool = sync.Pool{
    New: func() interface{} {
        return &Message{}
    },
}

func handleMessage() {
    msg := messagePool.Get().(*Message)
    defer messagePool.Put(msg)
    // Use msg...
}
```

### WebSocket Optimization
```go
// Efficient broadcasting
func (h *Hub) broadcastToSession(sessionID string, message []byte) {
    clients := h.sessions[sessionID]
    for _, client := range clients {
        select {
        case client.send <- message:
            // Success
        default:
            // Channel full - remove client
            h.unregisterClient(client)
        }
    }
}
```

## Best Practices

### Code Organization
- **Single Responsibility**: Each handler does one thing well
- **Error Handling**: Always handle and log errors appropriately
- **Resource Cleanup**: Use defer for cleanup (file handles, connections)
- **Context Awareness**: Pass context through request chains

### API Design
- **RESTful Conventions**: Use appropriate HTTP methods and status codes
- **Consistent Responses**: Standardize response formats
- **Validation**: Validate all inputs at API boundaries
- **Documentation**: Update OpenAPI spec before implementation

### WebSocket Guidelines
- **Message Types**: Use clear, specific message type names
- **Data Validation**: Validate all incoming WebSocket messages
- **Error Recovery**: Handle connection failures gracefully
- **Resource Limits**: Implement appropriate rate limiting

## Common Patterns

### Handler Pattern
```go
func HandlerName(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
    // 1. Parse and validate input
    var request RequestType
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    // 2. Business logic
    result, err := processRequest(request)
    if err != nil {
        logging.Error("processing failed", map[string]interface{}{
            "error": err.Error(),
        })
        http.Error(w, "Processing failed", http.StatusInternalServerError)
        return
    }
    
    // 3. WebSocket broadcast (if needed)
    broadcastMessage(hub, "event_type", result)
    
    // 4. HTTP response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(result)
}
```

### WebSocket Message Pattern
```javascript
// Client-side message sending
sendMessage(type, data) {
    const message = {
        type: type,
        data: data,
        timestamp: Date.now(),
        session_id: this.sessionId
    };
    
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.ws.send(JSON.stringify(message));
    }
}
```

## Troubleshooting Development Issues

### Build Failures
```bash
# Clean and regenerate
make clean
make generate
go mod tidy

# Check for syntax errors in api.yaml
make validate
```

### WebSocket Issues
```javascript
// Debug WebSocket in browser
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onopen = () => console.log('Connected');
ws.onerror = (e) => console.error('WebSocket error:', e);
ws.onclose = (e) => console.log('Disconnected:', e.code, e.reason);
```

### Performance Issues
```bash
# Profile server performance
go tool pprof http://localhost:8080/debug/pprof/heap
go tool pprof http://localhost:8080/debug/pprof/profile
```

---

*Development Guide for HD1 v0.7.0 - Three.js Game Engine Platform*