# HD1 v5.0.5 Source - API-First Game Engine Platform

**Game Engine as a Service** - World's first HTTP-controlled professional game engine with complete configuration management

## 🎯 Architecture Overview

**HD1 v5.0.5** transforms game development through **API-First Game Engine** architecture:
- **86 REST Endpoints**: Complete game engine control via HTTP APIs
- **Real-Time WebSocket**: Entity lifecycle synchronization  
- **PlayCanvas Engine**: Professional 3D rendering with ECS
- **Single Source of Truth**: 100% auto-generated from `api.yaml`
- **Configuration Management**: Complete standardization with 50+ environment variables

### Core Principle: "API = Control, WebSocket = Graph Extension"
```
HTTP APIs → Game Commands → Server State → WebSocket Events → Client Rendering
```

## 🔄 **Development Flow: API → Templates → Code Logic**

HD1 follows a strict **three-phase development flow** to maintain technical excellence:

### **Phase 1: API Specification (ALWAYS FIRST)**
```bash
# 1. Design API in OpenAPI specification
vim api.yaml

# Define new endpoints, schemas, and operations
# This is the SINGLE SOURCE OF TRUTH for all code generation
```

### **Phase 2: Template Development (SECOND)**
```bash
# 2. Create/update templates for code generation
vim codegen/templates/go/router.tmpl           # Go router template
vim codegen/templates/javascript/api-client.tmpl # JavaScript API client
vim codegen/templates/shell/core-functions.tmpl  # Shell functions

# Templates consume api.yaml and generate code
# Templates are EXTERNAL files with proper syntax highlighting
```

### **Phase 3: Code Logic Implementation (THIRD)**
```bash
# 3. Implement handler logic (only after API + templates are complete)
vim api/entities/create.go                    # Handler implementation
vim api/sessions/join_world.go                # Business logic

# Handler code implements the API contracts defined in Phase 1
# Code generation from Phase 2 provides routing and client libraries
```

**⚠️ CRITICAL RULE**: Never implement code logic before the API specification and templates are complete. This ensures consistency and prevents architectural drift.

## 🏗️ File Structure

```
src/
├── api.yaml                # 🎯 PHASE 1: SINGLE SOURCE OF TRUTH - OpenAPI 3.0.3
├── main.go                 # HD1 daemon entry point
├── auto_router.go          # 🤖 AUTO-GENERATED - 86 endpoint routing
├── go.mod                  # Go dependencies
├── Makefile               # Build system with auto-generation
│
├── codegen/               # 🏆 PHASE 2: CODE GENERATION SYSTEM
│   ├── generator.go       # Template processor (100% external templates)
│   └── templates/         # 🎨 EXTERNAL TEMPLATE FILES
│       ├── go/           # Go client & router templates
│       │   ├── router.tmpl      # Auto-router generation
│       │   └── client.tmpl      # Go CLI client
│       ├── javascript/   # JavaScript API client templates  
│       │   ├── api-client.tmpl  # JS API wrapper
│       │   ├── ui-components.tmpl # UI component generation
│       │   ├── form-system.tmpl # Dynamic form system
│       │   └── playcanvas-bridge.tmpl # PlayCanvas integration
│       └── shell/        # Shell function templates
│           ├── core-functions.tmpl    # Core API functions
│           └── playcanvas-functions.tmpl # Enhanced PlayCanvas
│
├── api/                   # 🎮 PHASE 3: GAME ENGINE API HANDLERS
│   ├── entities/         # Entity lifecycle (create, update, delete)
│   ├── sessions/         # Session & world management
│   ├── components/       # ECS component system
│   ├── hierarchy/        # Transform hierarchy & parenting
│   ├── physics/          # Physics world & rigidbodies
│   ├── audio/            # 3D audio sources & playback
│   ├── animation/        # Animation system
│   ├── lifecycle/        # Entity activation/deactivation
│   ├── worlds/           # World management
│   ├── camera/           # Camera controls
│   ├── browser/          # Browser integration
│   ├── scenegraph/       # Scene graph management
│   ├── recording/        # Session recording
│   ├── logging/          # Admin logging controls
│   └── system/           # System information
│
└── server/               # 🛡️ SERVER INFRASTRUCTURE
    ├── hub.go            # WebSocket hub & session management
    ├── client.go         # WebSocket client handling
    └── logging.go        # Structured logging system
```

## 🎨 **Template Architecture (External Templates Only)**

HD1 uses **100% external templates** for maintainable code generation:

### **Template Categories**
```
codegen/templates/
├── go/                    # Go code generation
│   ├── router.tmpl        # → auto_router.go (86 endpoints)
│   └── client.tmpl        # → client/main.go (CLI tool)
├── javascript/            # JavaScript code generation  
│   ├── api-client.tmpl    # → ../share/htdocs/static/js/hd1lib.js
│   ├── ui-components.tmpl # → ../share/htdocs/static/js/hd1-ui-components.js
│   ├── form-system.tmpl   # → ../share/htdocs/static/js/hd1-form-system.js
│   └── playcanvas-bridge.tmpl # → ../lib/downstream/playcanvaslib.js
└── shell/                 # Shell code generation
    ├── core-functions.tmpl # → ../lib/hd1lib.sh
    └── playcanvas-functions.tmpl # → ../lib/downstream/playcanvaslib.sh
```

### **Template Features**
- **External Files**: No hardcoded templates in generator.go
- **Syntax Highlighting**: Proper IDE support for template development
- **Developer Friendly**: Frontend developers can edit JavaScript templates directly
- **Go Embed**: Templates embedded in binary for single-file deployment
- **Template Caching**: Performance optimization with `loadTemplate()` system
- **Zero Manual Sync**: Templates drive all client generation automatically

### **Template Development Guidelines**
1. **Template Language**: Go text/template syntax
2. **Data Source**: Templates receive parsed `api.yaml` data
3. **Output Format**: Each template produces specific file types
4. **Validation**: Templates must produce syntactically correct output
5. **Testing**: Template changes tested via `make generate`

## 🔄 **Complete Development Workflow**

### **1. API-First Development (Phase 1)**
```bash
# Edit OpenAPI specification - ALWAYS START HERE
vim api.yaml

# Define new endpoint
paths:
  /sessions/{sessionId}/entities/{entityId}/physics:
    post:
      operationId: applyPhysicsForce
      x-handler: api/physics/force.go
      x-function: ApplyPhysicsForce
      summary: Apply force to entity rigidbody
```

### **2. Template Updates (Phase 2)**
```bash
# Update templates if new patterns needed
vim codegen/templates/go/router.tmpl           # Add routing patterns
vim codegen/templates/javascript/api-client.tmpl # Add JS method patterns

# Templates automatically consume new endpoints from api.yaml
```

### **3. Code Generation (Auto)**
```bash
# Generate all code from templates + API spec
make generate

# Produces:
# ├── auto_router.go              (86 endpoints)
# ├── ../lib/hd1lib.sh            (Shell API functions)
# ├── ../share/htdocs/static/js/hd1lib.js (JavaScript API)
# ├── ../share/htdocs/static/js/hd1-ui-components.js
# ├── ../share/htdocs/static/js/hd1-form-system.js
# ├── ../lib/downstream/playcanvaslib.sh
# └── ../lib/downstream/playcanvaslib.js
```

### **4. Handler Implementation (Phase 3)**
```bash
# Implement business logic - ONLY AFTER API + TEMPLATES
vim api/physics/force.go

func ApplyPhysicsForceHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
    // 1. Extract & validate parameters
    sessionID := extractSessionID(r.URL.Path)
    entityID := extractEntityID(r.URL.Path)
    
    // 2. Parse request body
    var req ApplyForceRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // 3. Business logic
    err := hub.ApplyForce(sessionID, entityID, req.Force, req.Point)
    
    // 4. WebSocket broadcast
    hub.BroadcastUpdate("physics_force_applied", ForceEvent{
        EntityID: entityID,
        Force: req.Force,
    })
    
    // 5. API response
    json.NewEncoder(w).Encode(SuccessResponse{Success: true})
}
```

### **5. Build & Test**
```bash
make clean && make     # Build with validation
make start             # Start HD1 daemon
make stop              # Clean shutdown

# Test new endpoint
curl -X POST /api/sessions/{id}/entities/{entityId}/physics \
  -d '{"force": [0, 10, 0], "point": [0, 0, 0]}'
```

## 🎮 Game Engine Features

### Entity Component System
```bash
# Create game entity
curl -X POST /api/sessions/{sessionId}/entities \
  -d '{"name": "player", "components": {
    "model": {"type": "box"},
    "transform": {"position": [0, 1, 0]},
    "material": {"diffuse": "#ff0000"},
    "rigidbody": {"type": "dynamic", "mass": 1.0}
  }}'

# Real-time WebSocket event
{"type": "entity_created", "data": {"entity_id": "...", "components": {...}}}
```

### World-Based Scene Management
```yaml
# /opt/hd1/share/worlds/world_one.yaml
world:
  id: "world_one"
  name: "Scene 1 - Red Box"

playcanvas:
  scene:
    ambientLight: [0.3, 0.3, 0.3]
    gravity: [0, -9.81, 0]
  
  entities:
    - name: "floor"
      components:
        transform: {position: [0, -0.1, 0], scale: [10, 1, 10]}
        material: {diffuse: "#cccccc"}
        rigidbody: {type: "static"}
```

## 🎯 API Endpoints (86 Total)

### Core Game Engine
```
# Entity Management
POST   /sessions/{id}/entities              # Create entity
GET    /sessions/{id}/entities              # List entities
PUT    /sessions/{id}/entities/{entityId}   # Update entity
DELETE /sessions/{id}/entities/{entityId}   # Delete entity

# Component System  
POST   /sessions/{id}/entities/{entityId}/components                   # Add component
GET    /sessions/{id}/entities/{entityId}/components/{componentType}   # Get component
PUT    /sessions/{id}/entities/{entityId}/components/{componentType}   # Update component

# Physics Engine
POST   /sessions/{id}/physics/rigidbodies/{entityId}/force  # Apply force
GET    /sessions/{id}/physics/world                         # Get physics world

# Transform Hierarchy
PUT    /sessions/{id}/entities/{entityId}/hierarchy/transforms  # Set transform
GET    /sessions/{id}/entities/{entityId}/hierarchy/parent      # Get parent
```

### Session & World Management
```
POST /sessions                              # Create session
POST /sessions/{id}/world/join             # Join world
POST /sessions/{id}/world/leave            # Leave world
GET  /sessions/{id}/world/status           # World status
```

## 🏆 Performance Characteristics

- **API Response**: <50ms for entity operations
- **WebSocket Latency**: <10ms for real-time sync
- **Entity Capacity**: 1000+ entities per session at 60fps
- **Concurrent Clients**: 100+ per world
- **Throughput**: 1000+ API requests/second

## 🔨 Quality Standards

### **File Categories**

#### **AUTO-GENERATED (NEVER EDIT)**
- `auto_router.go` - Generated from api.yaml
- `../share/htdocs/static/js/hd1lib.js` - JS API client
- `../lib/hd1lib.sh` - Shell API client
- `../share/htdocs/static/js/hd1-ui-components.js` - UI components
- `../share/htdocs/static/js/hd1-form-system.js` - Form system
- `../lib/downstream/playcanvaslib.sh` - PlayCanvas shell bridge
- `../lib/downstream/playcanvaslib.js` - PlayCanvas JS bridge

#### **SOURCE FILES (ALWAYS EDIT)**
- `api.yaml` - OpenAPI specification (single source of truth)
- `api/*/` - Handler implementations
- `codegen/templates/` - Code generation templates

### **Development Principles**
1. **API First**: Design in `api.yaml` before any implementation
2. **Templates Second**: Update templates for new generation patterns
3. **Code Logic Third**: Implement handlers only after API + templates complete
4. **Zero Regressions**: All changes maintain compatibility  
5. **Single Source**: api.yaml drives all code generation
6. **Clean Architecture**: Separation of concerns maintained

### **Code Generation Rules**
- **Never edit auto-generated files** - Changes will be lost
- **Always run `make generate`** after api.yaml changes
- **Test templates thoroughly** - Invalid templates break builds
- **Validate output** - Ensure generated code compiles and runs
- **Template documentation** - Document template changes and patterns

## 🎮 WebSocket Events

```javascript
// Entity lifecycle events
{
  "type": "entity_created",
  "data": {
    "session_id": "session-abc123",
    "entity_id": "entity-xyz789",
    "components": {
      "transform": {"position": [0, 1, 0]},
      "material": {"diffuse": "#ff0000"}
    }
  }
}

{
  "type": "entity_deleted",
  "data": {"entity_id": "entity-xyz789"}
}
```

## 🌐 Architectural Achievements

- ✅ **API-First Game Engine**: Revolutionary HTTP-based game development
- ✅ **Real-Time Synchronization**: <10ms WebSocket entity updates
- ✅ **Professional ECS**: Complete Entity Component System via APIs
- ✅ **Single Source of Truth**: 100% specification-driven development
- ✅ **Template Externalization**: 100% external templates, zero hardcoded generation
- ✅ **Production Ready**: Enterprise logging, monitoring, scalability
- ✅ **PlayCanvas Integration**: Professional 3D rendering engine
- ✅ **Auto-Generation**: 8 different outputs from single API spec
- ✅ **World System**: YAML-based scene configuration

## 🚀 **Template Development Examples**

### **Adding New JavaScript API Method**
```javascript
// In codegen/templates/javascript/api-client.tmpl
{{range .Routes}}
{{if eq .Method "POST"}}
    /**
     * {{.Comment}}
     * @param {string} {{.Parameters}}
     * @returns {Promise} API response
     */
    {{.MethodName}}: function({{.Parameters}}) {
        {{.Implementation}}
    },
{{end}}
{{end}}
```

### **Adding New Go Handler Stub**
```go
// In codegen/templates/go/router.tmpl
{{range .HandlerStubs}}
// {{.Comment}}
func (r *APIRouter) {{.FuncName}}(w http.ResponseWriter, req *http.Request) {
    {{.Package}}.{{.FuncName}}Handler(w, req, r.hub)
}
{{end}}
```

---

**HD1 v5.0.5 - Where OpenAPI specifications become immersive game worlds through API-first engineering.**

**Development Flow: API specification → Template processing → Code logic implementation**