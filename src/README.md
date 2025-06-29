# 🏆 THD Source Code - Revolutionary Specification-Driven Architecture

**Professional VR/AR Holodeck Platform with Revolutionary Upstream/Downstream Integration**

This directory contains the core implementation of **THD (The Holo-Deck)** - the world's first **professional holodeck platform** with revolutionary specification-driven architecture and complete upstream/downstream API integration.

## 🎯 Revolutionary Architecture Overview

### **Single Source of Truth Pipeline**
```
api.yaml (specification) → generator.go → {
    auto_router.go (Go routing)
    thdlib.sh (shell API client)
    thdlib.js (JavaScript API client)
    aframelib.sh (A-Frame shell integration)
    aframelib.js (A-Frame JavaScript bridge)
}
```

**Revolutionary Achievement**: Change the API specification = automatically regenerate ALL client libraries and routing across ALL environments.

## 🏗️ Architectural Principles

### **1. Specification-Driven Development**
- **OpenAPI 3.0.3** as absolute single source of truth
- **Zero manual routing** - everything auto-generated
- **Perfect synchronization** - clients never fall out of sync
- **Build-time validation** - prevents deployment of incomplete implementations

### **2. Revolutionary Upstream/Downstream Integration**
- **Upstream**: Core THD API wrappers (`thdlib.*`)
- **Downstream**: A-Frame WebXR integration (`aframelib.*`)
- **Identical signatures** - shell and JavaScript functions have identical parameters
- **Perfect layering** - downstream imports upstream maintaining single source of truth

### **3. Professional Engineering Standards**
- **Thread-safe concurrency** - mutex-protected session management
- **Enterprise-grade logging** - structured, timestamped, trace modules
- **Quality assurance** - comprehensive validation preventing regressions
- **Professional build system** - Make-based with daemon control

## 📁 Revolutionary Directory Structure

```
src/
├── api.yaml                # 🎯 SINGLE SOURCE OF TRUTH - OpenAPI 3.0.3 Specification
├── main.go                 # Professional THD daemon with holodeck integration
├── auto_router.go          # 🤖 AUTO-GENERATED - Revolutionary routing from spec
├── go.mod & go.sum         # Go module dependencies
├── Makefile               # Professional build system with validation
│
├── codegen/               # 🏆 REVOLUTIONARY CODE GENERATION SYSTEM
│   ├── generator.go       # Unified generator - upstream + downstream
│   ├── enhanced_generator.go  # A-Frame integration generator
│   └── aframe_schema_reader.go  # A-Frame schema validation
│
├── api/                   # 🎪 PROFESSIONAL API HANDLER IMPLEMENTATIONS
│   ├── sessions/          # Session lifecycle management
│   │   ├── create.go      # POST /sessions - Create holodeck session
│   │   ├── list.go        # GET /sessions - List active sessions
│   │   ├── get.go         # GET /sessions/{id} - Session details
│   │   └── delete.go      # DELETE /sessions/{id} - Terminate session
│   │
│   ├── objects/           # 3D object management with A-Frame integration
│   │   ├── create.go      # POST /objects - Create 3D objects
│   │   ├── list.go        # GET /objects - List session objects
│   │   ├── get.go         # GET /objects/{name} - Object details
│   │   ├── update.go      # PUT /objects/{name} - Update properties
│   │   └── delete.go      # DELETE /objects/{name} - Remove objects
│   │
│   ├── scenes/            # 🎭 PROFESSIONAL SCENE MANAGEMENT
│   │   ├── list.go        # GET /scenes - Available holodeck scenes
│   │   ├── load.go        # POST /scenes/{id} - Load scene into session
│   │   ├── save.go        # POST /scenes/save - Save session as scene
│   │   └── fork.go        # POST /scenes/{id}/fork - Scene forking
│   │
│   ├── camera/            # Professional camera control
│   │   ├── position.go    # PUT /camera/position - Set coordinates
│   │   └── orbit.go       # POST /camera/orbit - Orbital motion
│   │
│   ├── browser/           # Canvas control and rendering
│   │   └── control.go     # POST /browser/canvas - Canvas manipulation
│   │
│   ├── recording/         # 🎥 TEMPORAL RECORDING SYSTEM
│   │   ├── start.go       # POST /recording/start - Begin recording
│   │   ├── stop.go        # POST /recording/stop - End recording
│   │   ├── play.go        # POST /recording/play - Playback session
│   │   └── status.go      # GET /recording/status - Recording state
│   │
│   └── logging/           # 📊 PROFESSIONAL LOGGING CONTROL
│       └── handlers.go    # Logging configuration and trace modules
│
├── server/                # 🛡️ ENTERPRISE-GRADE SERVER INFRASTRUCTURE
│   ├── hub.go             # WebSocket hub with thread-safe SessionStore
│   ├── client.go          # WebSocket client association and management
│   ├── handlers.go        # Professional static file serving + A-Frame
│   ├── logging.go         # Structured logging with enterprise features
│   ├── semantic.go        # Professional UI component generation
│   └── version.go         # Version management and build info
│
└── logging/               # 🔍 PROFESSIONAL LOGGING SYSTEM
    ├── config.go          # Logging configuration management
    └── logger.go          # Thread-safe structured logging
```

## 🎯 Revolutionary Code Generation

### **Core Generator (`codegen/generator.go`)**
**Revolutionary unified generator** producing:
- **Go routing** from OpenAPI specification
- **Shell API client** (`thdlib.sh`) from API endpoints
- **JavaScript API client** (`thdlib.js`) from API endpoints
- **Web UI components** auto-generated from schemas
- **A-Frame integration** (`aframelib.*`) with schema validation

### **Generation Command**
```bash
make generate
```

**Produces:**
```
✅ auto_router.go - Go routing (28 routes)
✅ ../lib/thdlib.sh - Shell API client (upstream)
✅ ../share/htdocs/static/js/thdlib.js - JavaScript client (upstream)
✅ ../lib/downstream/aframelib.sh - A-Frame shell integration
✅ ../lib/downstream/aframelib.js - A-Frame JavaScript bridge
✅ Web UI components with dynamic forms
```

## 🏆 Revolutionary Features

### **1. Perfect Upstream/Downstream Integration**

**Upstream Libraries** (Auto-generated from `api.yaml`):
```bash
# Shell
source /opt/holo-deck/lib/thdlib.sh
thd::create_object "cube1" "box" 0 1 0
thd::camera 5 5 5
```

```javascript
// JavaScript (identical API coverage)
await thdAPI.createObject('session-id', {name: 'cube1', type: 'box', x: 0, y: 1, z: 0});
await thdAPI.setCameraPosition('session-id', {x: 5, y: 5, z: 5});
```

**Downstream A-Frame Integration** (Identical signatures):
```bash
# Shell A-Frame integration
source /opt/holo-deck/lib/downstream/aframelib.sh
thd::create_enhanced_object "crystal" "cone" 0 3 0 --color "#ff0000" --metalness 0.8
```

```javascript
// JavaScript A-Frame bridge (identical signature)
await thd.createEnhancedObject('crystal', 'cone', 0, 3, 0, {color: '#ff0000', metalness: 0.8});
```

### **2. Professional Scene Management**
- **Scene Collection**: Professional scenes in `share/scenes/`
- **API Integration**: Scenes accessible via `/api/scenes` endpoints
- **Web UI**: Scene dropdown with 30-day cookie persistence
- **Session Isolation**: Perfect scene separation across sessions

### **3. Holodeck Containment System**
- **Universal Coordinates**: Professional [-12, +12] holodeck boundaries
- **60fps Monitoring**: Real-time position checking with visual feedback
- **Escape-proof Design**: Dual boundary enforcement system
- **Professional Standards**: Enterprise-grade spatial control

## 🔄 Professional Development Workflow

### **1. API Specification Changes**
```bash
# Edit the single source of truth
vim api.yaml

# Regenerate everything automatically
make generate

# All clients now updated automatically:
# - Go routing (auto_router.go)
# - Shell functions (lib/thdlib.sh)  
# - JavaScript client (share/htdocs/static/js/thdlib.js)
# - A-Frame integration (lib/downstream/aframelib.*)
```

### **2. Build & Deployment**
```bash
make all        # Complete build pipeline with validation
make start      # Start THD daemon professionally  
make status     # Professional status reporting
make stop       # Clean shutdown with resource cleanup
```

### **3. Professional Testing**
```bash
# Test scene functionality
THD_SESSION=test-session bash share/scenes/basic-shapes.sh

# Test API endpoints
curl -X POST http://localhost:8080/api/sessions
curl -X GET http://localhost:8080/api/scenes
```

## 🌐 Thread-Safe Session Architecture

### **SessionStore - Enterprise Concurrency**
```go
type SessionStore struct {
    mutex    sync.RWMutex                    // Thread-safe operations
    sessions map[string]*Session             // Session metadata
    objects  map[string]map[string]*Object   // sessionId -> objectName -> Object
    worlds   map[string]*World               // sessionId -> World config
}
```

**Revolutionary Features:**
- **Perfect Isolation**: Sessions cannot access each other's data
- **Thread Safety**: Concurrent operations with proper mutex protection
- **Real-time Updates**: WebSocket broadcasting with session association
- **Enterprise Quality**: Production-ready concurrency patterns

### **WebSocket Hub - Real-time Communication**
```go
type Hub struct {
    clients       map[*Client]bool    // Connected WebSocket clients
    sessionClients map[string][]*Client // Session-specific client mapping
    broadcast     chan []byte         // Message broadcasting
    store         *SessionStore       // Thread-safe persistence
}
```

## 🎮 Professional Holodeck Coordinates

**Coordinate System** (Specification: `docs/design/holodeck-coordinates.md`):
- **Boundaries**: [-12, +12] on all axes (holodeck-grade)
- **Grid System**: 25×25×25 professional grid
- **Human Standards**: Y=0 floor, Y=1.7 eye level, Y=3.0 ceiling
- **Validation**: Automatic boundary enforcement throughout

## 🔨 Professional Build System

### **Makefile Targets**
```bash
make all        # Complete build pipeline with validation
make generate   # Revolutionary code generation from spec
make build      # Compile THD daemon 
make validate   # Validate API specification
make start      # Professional daemon startup
make stop       # Clean shutdown with resource cleanup
make status     # Professional status reporting
make restart    # Restart with validation
make test       # API endpoint testing
make clean      # Remove build artifacts
```

### **Quality Assurance Pipeline**
- ✅ **Specification Validation**: OpenAPI 3.0.3 schema checking
- ✅ **Handler Validation**: All referenced handlers must exist
- ✅ **Generation Verification**: All clients successfully generated
- ✅ **Build Validation**: Go compilation must succeed
- ✅ **Professional Standards**: No regressions allowed

## 🧪 API Handler Excellence

**Professional Handler Pattern:**
```go
func HandlerName(w http.ResponseWriter, r *http.Request, hub interface{}) {
    // 1. Professional type casting with error handling
    h, ok := hub.(*server.Hub)
    if !ok {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(ErrorResponse{
            Success: false,
            Error:   "internal_error",
            Message: "Internal server error",
        })
        return
    }

    // 2. Professional parameter extraction and validation
    sessionID := extractSessionID(r.URL.Path)
    if sessionID == "" {
        // Professional error response...
    }

    // 3. Thread-safe business logic
    result, err := h.GetStore().SomeOperation(sessionID, request)

    // 4. Real-time WebSocket broadcasting
    h.BroadcastToSession(sessionID, "event_type", result)

    // 5. Professional JSON response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(SuccessResponse{
        Success: true,
        Data:    result,
    })
}
```

## ⚠️ Critical Development Rules

### **NEVER EDIT (Auto-Generated)**
- `auto_router.go` - Generated from specification
- `lib/thdlib.sh` - Generated from API endpoints
- `share/htdocs/static/js/thdlib.js` - Generated JavaScript client
- `lib/downstream/aframelib.*` - Generated A-Frame integration

### **ALWAYS EDIT (Source of Truth)**
- `api.yaml` - The specification that drives everything
- Handler implementations in `api/` directories
- `Makefile` build configuration
- Scene scripts in `share/scenes/`

### **DEVELOPMENT PRINCIPLES**
1. **Specification First**: Design in `api.yaml` before implementation
2. **Professional Standards**: Enterprise-grade logging, error handling
3. **Single Source of Truth**: Zero manual synchronization
4. **Thread Safety**: All concurrent operations properly protected
5. **Real-time Updates**: Broadcast all state changes via WebSocket
6. **Quality Assurance**: Comprehensive validation preventing regressions

## 🎯 Extension Capabilities

### **Adding New API Endpoints**
1. **Define in specification**: Add to `api.yaml` with handler reference
2. **Generate routing**: `make generate` creates routing automatically
3. **Implement handler**: Create Go handler following professional pattern
4. **Auto-generated clients**: Shell and JavaScript functions created automatically

### **Adding Downstream Integrations**
1. **Create integration directory**: `lib/downstream/newframework/`
2. **Implement shell functions**: Following `aframelib.sh` pattern
3. **Implement JavaScript bridge**: Following `aframelib.js` pattern
4. **Update generator**: Add generation logic for new framework

### **Adding Professional Scenes**
1. **Create scene script**: In `share/scenes/` using THD functions
2. **Update scene handler**: Add to scene mapping in `api/scenes/load.go`
3. **Test integration**: Verify scene loads via API and web interface

## 🏆 Architectural Achievements

- **🎯 100% Single Source of Truth**: API specification drives everything
- **🔄 Zero Manual Synchronization**: All clients auto-generated and consistent
- **🛡️ Enterprise Concurrency**: Thread-safe session management
- **⚡ Real-time Communication**: WebSocket hub with session association
- **🎮 Professional VR/AR**: Complete A-Frame WebXR integration
- **🔧 Quality Assurance**: Build-time validation preventing regressions
- **📊 Professional Logging**: Structured, timestamped, trace modules
- **🎭 Scene Management**: Professional scene collection with API integration

---

**THD Source Code represents the pinnacle of specification-driven development with revolutionary upstream/downstream integration architecture.**

*Where API specifications become immersive holodeck experiences through surgical engineering precision.*