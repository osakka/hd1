# HD1 Source Code - Three-Layer Game Engine Architecture

**Game Engine Architecture VR/AR Platform with Environment Physics and Props System**

This directory contains the core implementation of **HD1 (Holodeck One)** with revolutionary three-layer architecture (Environment + Props + Scenes), advanced specification-driven development, and industry-standard game engine patterns matching Unity and Unreal Engine.

## 🎯 Advanced Architecture Overview

### **Three-Layer Architecture Pipeline**
```
api.yaml (31 endpoints) → generator.go → {
    auto_router.go (Go routing with 3-layer APIs)
    hd1lib.sh (shell API client with environment/props)
    hd1lib.js (JavaScript API client)
    aframelib.sh (A-Frame shell integration)
    aframelib.js (A-Frame JavaScript bridge)
}
```

**Revolutionary Achievement**: Complete game engine architecture with Environment + Props + Scene layers, all auto-generated from single specification source.

## 🏗️ Architectural Principles

### **1. Three-Layer Game Engine Architecture**
- **Environment Layer**: Physics contexts (gravity, atmosphere, scale) - 4 distinct environments
- **Props Layer**: Reusable objects with realistic physics properties - 6 categories
- **Scene Layer**: Orchestration combining environments + prop placement (future Phase 3)
- **Game Engine Parity**: Matches Unity (World Settings + Prefabs + Scenes) and Unreal patterns

### **2. Specification-Driven Development**
- **OpenAPI 3.0.3** as absolute single source of truth (31 endpoints)
- **Zero manual routing** - everything auto-generated including three-layer APIs
- **Perfect synchronization** - clients never fall out of sync
- **Build-time validation** - prevents deployment of incomplete implementations

### **3. Advanced Physics Cohesion Engine**
- **Environment-Aware Props**: Props automatically adapt physics based on session environment
- **Realistic Materials**: Material-accurate properties (wood: 600 kg/m³, metal: 7800 kg/m³)
- **Real-time Adaptation**: Physics recalculated instantly on environment changes
- **Scale Compatibility**: Props work across multiple scale units (nm, mm, cm, m, km)

### **4. Standard Engineering Standards**
- **Thread-safe concurrency** - mutex-protected session management
- **High-quality logging** - structured, timestamped, trace modules
- **Quality assurance** - comprehensive validation preventing regressions
- **Standard build system** - Make-based with daemon control

## 📁 Advanced Directory Structure

```
src/
├── api.yaml                # 🎯 SINGLE SOURCE OF TRUTH - OpenAPI 3.0.3 Specification
├── main.go                 # Standard HD1 daemon with holodeck integration
├── auto_router.go          # 🤖 AUTO-GENERATED - Advanced routing from spec
├── go.mod & go.sum         # Go module dependencies
├── Makefile               # Standard build system with validation
│
├── codegen/               # 🏆 REVOLUTIONARY CODE GENERATION SYSTEM
│   ├── generator.go       # Unified generator - upstream + downstream
│   ├── enhanced_generator.go  # A-Frame integration generator
│   └── aframe_schema_reader.go  # A-Frame schema validation
│
├── api/                   # 🎪 THREE-LAYER API HANDLER IMPLEMENTATIONS
│   ├── environments/      # 🌍 Environment System (Layer 1)
│   │   ├── list.go        # GET /environments - List physics contexts
│   │   └── apply.go       # POST /environments/{id} - Apply environment to session
│   │
│   ├── props/             # 🏗️ Props System (Layer 2)
│   │   ├── list.go        # GET /props - List available props
│   │   └── instantiate.go # POST /props/{id} - Instantiate prop in session
│   │
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
│   ├── camera/            # Standard camera control
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
│   ├── handlers.go        # Standard static file serving + A-Frame
│   ├── logging.go         # Structured logging with enterprise features
│   ├── semantic.go        # Standard UI component generation
│   └── version.go         # Version management and build info
│
└── logging/               # 🔍 PROFESSIONAL LOGGING SYSTEM
    ├── config.go          # Logging configuration management
    └── logger.go          # Thread-safe structured logging
```

## 🎯 Advanced Code Generation

### **Three-Layer Generator (`codegen/generator.go`)**
**Revolutionary unified generator** producing:
- **Go routing** from OpenAPI specification (31 endpoints including environment/props)
- **Shell API client** (`hd1lib.sh`) with three-layer architecture support
- **JavaScript API client** (`hd1lib.js`) from API endpoints
- **Web UI components** auto-generated from schemas
- **A-Frame integration** (`aframelib.*`) with schema validation

### **Generation Command**
```bash
make generate
```

**Produces:**
```
✅ auto_router.go - Go routing (31 routes including environment/props APIs)
✅ ../lib/hd1lib.sh - Shell API client with three-layer support
✅ ../share/htdocs/static/js/hd1lib.js - JavaScript client (upstream)
✅ ../lib/downstream/aframelib.sh - A-Frame shell integration
✅ ../lib/downstream/aframelib.js - A-Frame JavaScript bridge
✅ Web UI components with dynamic forms
```

## 🏆 Advanced Features

### **1. Three-Layer Game Engine Architecture**

**Environment System (Layer 1):**
```bash
# Shell API
hd1::apply_environment "session-id" "underwater"

# JavaScript API  
await hd1API.applyEnvironment('session-id', 'underwater');
```

**Props System (Layer 2):**
```bash
# Shell API
hd1::instantiate_prop "session-id" "wooden-chair" '{"x": 0, "y": 2, "z": 0}'

# JavaScript API
await hd1API.instantiateProp('session-id', 'wooden-chair', {x: 0, y: 2, z: 0});
```

**Physics Cohesion Engine:**
- Props automatically adapt to environment physics
- Real-time physics recalculation on environment changes
- Material-accurate properties (wood, metal, plastic densities)

### **2. Perfect Upstream/Downstream Integration**

**Upstream Libraries** (Auto-generated from `api.yaml`):
```bash
# Shell
source /opt/holodeck-one/lib/hd1lib.sh
hd1::create_object "cube1" "box" 0 1 0
hd1::camera 5 5 5
```

```javascript
// JavaScript (identical API coverage)
await hd1API.createObject('session-id', {name: 'cube1', type: 'box', x: 0, y: 1, z: 0});
await hd1API.setCameraPosition('session-id', {x: 5, y: 5, z: 5});
```

**Downstream A-Frame Integration** (Identical signatures):
```bash
# Shell A-Frame integration
source /opt/holodeck-one/lib/downstream/aframelib.sh
hd1::create_enhanced_object "crystal" "cone" 0 3 0 --color "#ff0000" --metalness 0.8
```

```javascript
// JavaScript A-Frame bridge (identical signature)
await hd1.createEnhancedObject('crystal', 'cone', 0, 3, 0, {color: '#ff0000', metalness: 0.8});
```

### **2. Standard Scene Management**
- **Scene Collection**: Standard scenes in `share/scenes/`
- **API Integration**: Scenes accessible via `/api/scenes` endpoints
- **Web UI**: Scene dropdown with 30-day cookie persistence
- **Session Isolation**: Perfect scene separation across sessions

### **3. Holodeck Containment System**
- **Universal Coordinates**: Standard [-12, +12] holodeck boundaries
- **60fps Monitoring**: Real-time position checking with visual feedback
- **Escape-proof Design**: Dual boundary enforcement system
- **Standard Standards**: High-quality spatial control

## 🔄 Standard Development Workflow

### **1. API Specification Changes**
```bash
# Edit the single source of truth
vim api.yaml

# Regenerate everything automatically
make generate

# All clients now updated automatically:
# - Go routing (auto_router.go)
# - Shell functions (lib/hd1lib.sh)  
# - JavaScript client (share/htdocs/static/js/hd1lib.js)
# - A-Frame integration (lib/downstream/aframelib.*)
```

### **2. Build & Deployment**
```bash
make all        # Complete build pipeline with validation
make start      # Start HD1 daemon standardly  
make status     # Standard status reporting
make stop       # Clean shutdown with resource cleanup
```

### **3. Standard Testing**
```bash
# Test scene functionality
HD1_SESSION=test-session bash share/scenes/basic-shapes.sh

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

**Advanced Features:**
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

## 🎮 Standard Holodeck Coordinates

**Coordinate System** (Specification: `docs/design/holodeck-coordinates.md`):
- **Boundaries**: [-12, +12] on all axes (holodeck-grade)
- **Grid System**: 25×25×25 standard grid
- **Human Standards**: Y=0 floor, Y=1.7 eye level, Y=3.0 ceiling
- **Validation**: Automatic boundary enforcement throughout

## 🔨 Standard Build System

### **Makefile Targets**
```bash
make all        # Complete build pipeline with validation
make generate   # Advanced code generation from spec
make build      # Compile HD1 daemon 
make validate   # Validate API specification
make start      # Standard daemon startup
make stop       # Clean shutdown with resource cleanup
make status     # Standard status reporting
make restart    # Restart with validation
make test       # API endpoint testing
make clean      # Remove build artifacts
```

### **Quality Assurance Pipeline**
- ✅ **Specification Validation**: OpenAPI 3.0.3 schema checking
- ✅ **Handler Validation**: All referenced handlers must exist
- ✅ **Generation Verification**: All clients successfully generated
- ✅ **Build Validation**: Go compilation must succeed
- ✅ **Standard Standards**: No regressions allowed

## 🧪 API Handler Excellence

**Standard Handler Pattern:**
```go
func HandlerName(w http.ResponseWriter, r *http.Request, hub interface{}) {
    // 1. Standard type casting with error handling
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

    // 2. Standard parameter extraction and validation
    sessionID := extractSessionID(r.URL.Path)
    if sessionID == "" {
        // Standard error response...
    }

    // 3. Thread-safe business logic
    result, err := h.GetStore().SomeOperation(sessionID, request)

    // 4. Real-time WebSocket broadcasting
    h.BroadcastToSession(sessionID, "event_type", result)

    // 5. Standard JSON response
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
- `lib/hd1lib.sh` - Generated from API endpoints
- `share/htdocs/static/js/hd1lib.js` - Generated JavaScript client
- `lib/downstream/aframelib.*` - Generated A-Frame integration

### **ALWAYS EDIT (Source of Truth)**
- `api.yaml` - The specification that drives everything
- Handler implementations in `api/` directories
- `Makefile` build configuration
- Scene scripts in `share/scenes/`

### **DEVELOPMENT PRINCIPLES**
1. **Specification First**: Design in `api.yaml` before implementation
2. **Standard Standards**: High-quality logging, error handling
3. **Single Source of Truth**: Zero manual synchronization
4. **Thread Safety**: All concurrent operations properly protected
5. **Real-time Updates**: Broadcast all state changes via WebSocket
6. **Quality Assurance**: Comprehensive validation preventing regressions

## 🎯 Extension Capabilities

### **Adding New API Endpoints**
1. **Define in specification**: Add to `api.yaml` with handler reference
2. **Generate routing**: `make generate` creates routing automatically
3. **Implement handler**: Create Go handler following standard pattern
4. **Auto-generated clients**: Shell and JavaScript functions created automatically

### **Adding Downstream Integrations**
1. **Create integration directory**: `lib/downstream/newframework/`
2. **Implement shell functions**: Following `aframelib.sh` pattern
3. **Implement JavaScript bridge**: Following `aframelib.js` pattern
4. **Update generator**: Add generation logic for new framework

### **Adding Standard Scenes**
1. **Create scene script**: In `share/scenes/` using HD1 functions
2. **Update scene handler**: Add to scene mapping in `api/scenes/load.go`
3. **Test integration**: Verify scene loads via API and web interface

## 🏆 Architectural Achievements

- **🏗️ Three-Layer Game Engine**: Environment + Props + Scene architecture matching Unity/Unreal patterns
- **🌍 Physics Cohesion Engine**: Real-time environment-aware physics adaptation
- **🎯 100% Single Source of Truth**: API specification drives everything (31 endpoints)
- **🔄 Zero Manual Synchronization**: All clients auto-generated including three-layer APIs
- **🛡️ Enterprise Concurrency**: Thread-safe session management
- **⚡ Real-time Communication**: WebSocket hub with session association
- **🎮 Standard VR/AR**: Complete A-Frame WebXR integration
- **🔧 Quality Assurance**: Build-time validation preventing regressions
- **📊 Standard Logging**: Structured, timestamped, trace modules
- **🎭 Scene Management**: Standard scene collection with API integration

---

**HD1 Source Code represents the pinnacle of specification-driven development with advanced upstream/downstream integration architecture.**

*Where API specifications become immersive holodeck experiences through surgical engineering precision.*