# 🔧 VWS Source Code - Technical Implementation

## Overview

This directory contains the core implementation of the **VWS (Virtual World Synthesizer)** engine. The architecture follows strict specification-driven development principles where `api.yaml` serves as the single source of truth.

## 🏗️ Architecture Principles

### **Specification-Driven Development**
```
api.yaml → codegen/generator.go → auto_router.go → api/ handlers → Virtual World
```

1. **Single Source of Truth**: All API routing is generated from `api.yaml`
2. **Zero Manual Routing**: No hand-written route configurations
3. **Build-time Validation**: Missing handlers fail the build
4. **Automatic Code Generation**: Router regenerated on every build

## 📁 Directory Structure

```
src/
├── api.yaml                # 🎯 SINGLE SOURCE OF TRUTH - OpenAPI 3.0.3 Specification
├── main.go                 # Server entry point and WebSocket setup
├── auto_router.go          # 🤖 AUTO-GENERATED - Never edit manually
├── go.mod & go.sum         # Go module dependencies
├── Makefile               # Build automation and validation
│
├── codegen/               # Code generation system
│   └── generator.go       # Reads api.yaml, generates auto_router.go
│
├── api/                   # Modular API handler implementations
│   ├── sessions/          # Session lifecycle management
│   │   ├── create.go      # POST /sessions - Create virtual world
│   │   ├── list.go        # GET /sessions - List all worlds
│   │   ├── get.go         # GET /sessions/{id} - Get world details
│   │   └── delete.go      # DELETE /sessions/{id} - Terminate world
│   │
│   ├── world/             # Virtual world configuration
│   │   ├── init.go        # POST /world - Initialize world parameters
│   │   └── spec.go        # GET /world - Get world specifications
│   │
│   ├── objects/           # 3D object management
│   │   ├── create.go      # POST /objects - Create 3D objects
│   │   ├── list.go        # GET /objects - List all objects
│   │   ├── get.go         # GET /objects/{name} - Get object details
│   │   ├── update.go      # PUT /objects/{name} - Update properties
│   │   └── delete.go      # DELETE /objects/{name} - Remove objects
│   │
│   └── camera/            # Camera control system
│       ├── position.go    # PUT /camera/position - Set coordinates
│       └── orbit.go       # POST /camera/orbit - Start orbital motion
│
├── server/                # Core server infrastructure
│   ├── hub.go             # WebSocket hub and SessionStore
│   ├── client.go          # WebSocket client management
│   ├── handlers.go        # Static file serving
│   ├── logging.go         # Structured logging system
│   ├── semantic.go        # Semantic UI components
│   └── version.go         # Version management
│
└── renderer/              # 3D visualization engine
    └── static/            # WebGL renderer and UI
        └── js/            # JavaScript 3D engine
```

## 🎯 Core Files Explained

### **api.yaml** - The Heart of VWS
```yaml
# This file IS the system architecture
# Changes here automatically update the entire system
paths:
  /sessions:
    post:
      operationId: createSession
      x-handler: "api/sessions/create.go"
      x-function: "CreateSession"
```

### **auto_router.go** - Generated Routing
```go
// AUTO-GENERATED FROM api.yaml - DO NOT EDIT MANUALLY
// This file is the SINGLE SOURCE OF TRUTH for routing

func (r *APIRouter) generateRoutes() {
    r.routes = []Route{
        {
            Path:       "/sessions",
            Method:     "POST", 
            Handler:    r.CreateSession,
            OperationID: "createSession",
        },
        // ... all routes generated from specification
    }
}
```

### **main.go** - Server Bootstrap
```go
func main() {
    hub := server.NewHub()  // Initialize WebSocket hub with SessionStore
    go hub.Run()            // Start real-time communication

    // REVOLUTIONARY: Auto-generated API router from specification
    apiRouter := NewAPIRouter(hub)
    http.Handle("/api/", apiRouter)
    
    // Static file serving for 3D renderer
    http.Handle("/static/", ...)
    
    log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
```

## 🔄 Development Workflow

### **1. Define New Functionality**
Edit `api.yaml` to add new endpoints:
```yaml
/sessions/{sessionId}/physics:
  post:
    operationId: enablePhysics
    x-handler: "api/physics/enable.go"
    x-function: "EnablePhysics"
```

### **2. Auto-Generate Routing**
```bash
make generate  # Reads api.yaml, regenerates auto_router.go
```

### **3. Implement Handler**
Create `api/physics/enable.go`:
```go
func EnablePhysicsHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
    // Implementation here
}
```

### **4. Build & Test**
```bash
make build  # Validates all handlers exist, builds binary
make test   # Tests all endpoints
```

## 🧠 SessionStore Architecture

The SessionStore provides thread-safe persistence for all virtual world state:

```go
type SessionStore struct {
    mutex    sync.RWMutex
    sessions map[string]*Session              // Session metadata
    objects  map[string]map[string]*Object    // sessionId -> objectName -> Object
    worlds   map[string]*World                // sessionId -> World config
}
```

### **Key Features**
- **Thread-Safe**: Concurrent access with proper locking
- **Coordinate Validation**: Enforces [-12, +12] world boundaries
- **Real-time Updates**: Broadcasts changes via WebSocket hub
- **Session Isolation**: Complete data separation between sessions

## 🌐 WebSocket Real-time System

```go
type Hub struct {
    clients    map[*Client]bool     // Connected WebSocket clients
    broadcast  chan []byte          // Message broadcasting channel
    register   chan *Client         // Client connection registration
    unregister chan *Client         // Client disconnection handling
    store      *SessionStore        // Persistent state management
}
```

### **Real-time Events**
- `session_created` - New virtual world spawned
- `world_initialized` - World configuration applied
- `object_created` - New 3D object added
- `object_updated` - Object properties changed
- `session_deleted` - Virtual world terminated

## 🎮 3D Coordinate System

### **World Boundaries**
- **Coordinate Space**: [-12, +12] on all axes (X, Y, Z)
- **Grid Size**: Configurable (default 25×25×25)
- **Validation**: Automatic boundary enforcement
- **Snapping**: Objects align to coordinate grid

### **Object Properties**
```go
type Object struct {
    Name     string  `json:"name"`      // Unique identifier
    Type     string  `json:"type"`      // cube, sphere, etc.
    X        float64 `json:"x"`         // X coordinate [-12, +12]
    Y        float64 `json:"y"`         // Y coordinate [-12, +12] 
    Z        float64 `json:"z"`         // Z coordinate [-12, +12]
    Color    string  `json:"color"`     // Color specification
    Scale    float64 `json:"scale"`     // Size multiplier
    Rotation string  `json:"rotation"`  // Rotation specification
}
```

## 🔨 Build System

### **Makefile Targets**
```bash
make all        # Complete build pipeline
make validate   # Check api.yaml exists
make generate   # Auto-generate router from spec
make build      # Compile Go binary
make test       # Test API endpoints
make run        # Start VWS server
make clean      # Remove build artifacts
```

### **Build Validation**
- ✅ **Specification Check**: `api.yaml` must exist
- ✅ **Handler Validation**: All referenced handlers must exist
- ✅ **Code Generation**: Router must be regeneratable
- ✅ **Compilation**: Go build must succeed

## 🚀 API Handler Pattern

All handlers follow this pattern:

```go
func HandlerName(w http.ResponseWriter, r *http.Request, hub interface{}) {
    // 1. Cast hub to proper type
    h, ok := hub.(*server.Hub)
    if !ok {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    
    // 2. Extract/validate parameters
    sessionID := extractSessionID(r.URL.Path)
    
    // 3. Parse request body (if needed)
    var req RequestStruct
    json.NewDecoder(r.Body).Decode(&req)
    
    // 4. Business logic using SessionStore
    result, err := h.GetStore().SomeOperation(sessionID, req)
    
    // 5. Broadcast real-time updates
    h.BroadcastUpdate("event_type", result)
    
    // 6. Return JSON response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}
```

## 🧪 Testing Strategy

### **API Testing**
```bash
# Test session creation
curl -X POST http://localhost:8080/api/sessions

# Test world initialization  
curl -X POST http://localhost:8080/api/sessions/{id}/world \
  -H "Content-Type: application/json" \
  -d '{"size":25,"transparency":0.2}'

# Test object creation with validation
curl -X POST http://localhost:8080/api/sessions/{id}/objects \
  -H "Content-Type: application/json" \
  -d '{"name":"cube1","type":"cube","x":0,"y":0,"z":0}'
```

### **Development Control**
```bash
./dev-control.sh status   # Show server status
./dev-control.sh cycle    # Full build/test cycle
./dev-control.sh logs     # View recent logs
```

## ⚠️ Critical Rules

### **DO NOT EDIT**
- `auto_router.go` - This file is AUTO-GENERATED
- Generated route definitions
- Auto-generated handler stubs

### **ALWAYS EDIT**
- `api.yaml` - The single source of truth
- Handler implementations in `api/` directories
- Build configuration in `Makefile`

### **DEVELOPMENT PRINCIPLES**
1. **Specification First**: Design API in `api.yaml` before coding
2. **Handler Implementation**: Focus on business logic only
3. **Real-time Updates**: Broadcast all state changes
4. **Coordinate Validation**: Enforce world boundaries
5. **Error Handling**: Proper HTTP status codes

## 🎯 Extension Points

### **Adding New Object Types**
1. Update `api.yaml` if new endpoints needed
2. Extend `Object` struct in `server/hub.go`
3. Add type validation in object handlers
4. Update 3D renderer for new geometry

### **Adding New Camera Modes**
1. Define new endpoints in `api.yaml`
2. Implement handlers in `api/camera/`
3. Extend camera controls in 3D renderer

### **Adding Physics**
1. Design physics API in `api.yaml`
2. Create `api/physics/` handler directory
3. Integrate physics engine with SessionStore
4. Broadcast physics updates via WebSocket

## 🏆 Architecture Achievements

- **Zero Manual Routing**: 100% specification-generated
- **Build-time Validation**: Missing components fail build
- **Real-time Synchronization**: WebSocket state broadcasting
- **Thread-safe Persistence**: Concurrent session management
- **Coordinate System**: Enforced 3D world boundaries
- **Modular Design**: Clean separation of concerns

---

**VWS Source Code represents the pinnacle of specification-driven development.**

*Every line of code serves the vision: transforming API specifications into virtual worlds.*