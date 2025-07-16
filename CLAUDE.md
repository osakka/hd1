# HD1 (Holodeck One) - Development Context

## Core Principles
- **Universal Platform**: HD1 enables any service to render 3D interfaces
- **Multi-Tenant Architecture**: Unlimited concurrent sessions with isolation
- **Real-Time Collaboration**: WebRTC P2P with operational transforms
- **AI-Native Integration**: LLM avatars with content generation
- **Cross-Platform**: Web, mobile, desktop clients with plugin architecture
- **Single Source of Truth**: Centralized database with incremental schemas

## Current State (2025-07-15)
HD1 v0.7.0 is a **universal 3D interface platform** with complete **multi-tenant architecture**, **real-time collaboration**, **AI integration**, **cross-platform support**, and **enterprise features**.

### ✅ Current Features
- **Three.js Integration**: 3D rendering with native Three.js r170
- **TCP-Simple Sync**: Sequence-based synchronization for bulletproof reliability
- **Real-Time WebSocket**: Entity lifecycle synchronization with operation ordering
- **Pure WebGL**: Direct Three.js Scene/Mesh/Material/Geometry operations
- **Avatar System**: Real-time multiplayer avatars with position synchronization
- **Entity Management**: Create/update/delete boxes, spheres, and custom geometries
- **Scene Control**: Background colors, lighting, fog, and camera management
- **Console UI**: Debug monitoring with Three.js statistics
- **HTTP Asset Delivery**: Direct GLB model serving with proper headers
- **Configuration Management**: Complete system with priority order: Flags > Environment Variables > .env File > Defaults
- **Multi-Tenant Sessions**: Unlimited concurrent sessions with isolation
- **WebRTC Collaboration**: Peer-to-peer real-time communication
- **Operational Transforms**: Conflict-free collaborative editing
- **Asset Management**: File upload, versioning, and usage tracking
- **LLM Integration**: Multi-provider AI avatars (OpenAI, Claude)
- **Content Generation**: AI-powered template-based content creation
- **Cross-Platform Clients**: Web, mobile, desktop with platform adapters
- **Plugin Architecture**: Extensible hook-based system with sandboxing
- **Enterprise RBAC**: Role-based access control with granular permissions
- **Analytics Platform**: Event tracking, aggregation, and reporting
- **Security & Compliance**: Audit logging, API keys, compliance records

### Architecture
```
HTTP APIs → Sync Operations → WebSocket Events → Three.js Rendering
```

**Key Files:**
- `src/api.yaml` - Single source of truth for all Three.js operations
- `src/auto_router.go` - Auto-generated routing from specification
- `src/sync/reliable.go` - TCP-simple synchronization system
- `src/server/hub.go` - WebSocket hub with sequence ordering
- `src/threejs/bridge.go` - Three.js JSON bridge
- `share/htdocs/static/js/hd1-threejs.js` - Three.js scene manager
- `share/htdocs/static/js/hd1-sync.js` - Client-side sync system
- `share/htdocs/static/js/hd1-console-new.js` - Three.js console interface
- `src/codegen/templates/` - External template files for code generation
- `src/webrtc/` - WebRTC infrastructure for P2P communication
- `src/ot/` - Operational transform system for collaboration
- `src/assets/` - Asset management and CDN integration
- `src/llm/` - LLM avatar and content generation system
- `src/clients/` - Cross-platform client management
- `src/plugins/` - Plugin architecture and management
- `src/enterprise/` - Enterprise features (orgs, RBAC, analytics, security)

## Development Commands
```bash
cd src && make clean && make && make start  # Build & start
make stop                                   # Stop server
make generate                              # Auto-generate from api.yaml
```

## Three.js Sync System
**TCP-Simple Reliability**: HD1 uses sequence numbers for guaranteed operation ordering

### Core Synchronization
```
Client Operations → Sequence Assignment → WebSocket Broadcast → Client Application
```

**Sync Operations**:
- **Sequence Tracking**: Each operation gets incremental sequence number
- **Missing Detection**: Clients detect gaps and request missing operations
- **Ordered Application**: Operations applied in sequence order
- **Guaranteed Delivery**: TCP-like reliability without complexity

### Implementation
- **Server**: `src/sync/reliable.go` - 200-line sequence-based sync
- **Client**: `src/static/js/hd1-sync.js` - Gap detection and ordering
- **Protocol**: JSON operations over WebSocket with sequence IDs

## Three.js Entity System
**Pure Three.js Operations**: Direct Scene/Mesh/Material/Geometry manipulation

### Entity Operations
- **Create**: `POST /api/threejs/entities` - Create boxes, spheres, custom geometries
- **Update**: `PUT /api/threejs/entities/{id}` - Modify position, rotation, scale, materials
- **Delete**: `DELETE /api/threejs/entities/{id}` - Remove from scene
- **List**: `GET /api/threejs/entities` - Get all entities

### Avatar System
- **Move**: `POST /api/threejs/avatars/{sessionId}/move` - Real-time position updates
- **Create**: `POST /api/threejs/avatars` - Spawn avatar in scene
- **List**: `GET /api/threejs/avatars` - Get all active avatars

### Scene Management
- **Background**: `PUT /api/threejs/scene` - Set background color
- **Lighting**: `PUT /api/threejs/scene/lighting` - Configure lights
- **Fog**: `PUT /api/threejs/scene/fog` - Add atmospheric effects
- **Camera**: `PUT /api/threejs/scene/camera` - Control camera position/rotation

## Configuration Management System

**Priority Order**: Flags > Environment Variables > .env File > Defaults

### Environment Variables (HD1_* prefix)
```bash
# Server Configuration
HD1_HOST=0.0.0.0                          # Server bind host
HD1_PORT=8080                             # Server bind port  
HD1_API_BASE=http://0.0.0.0:8080/api     # External API base URL
HD1_INTERNAL_API_BASE=http://localhost:8080/api  # Internal API communications
HD1_VERSION=v0.7.0                       # HD1 version identifier
HD1_DAEMON=true                           # Run in daemon mode

# Directory Paths
HD1_ROOT_DIR=/opt/hd1                     # HD1 root directory
HD1_STATIC_DIR=/opt/hd1/share/htdocs/static  # Static files directory
HD1_LOG_DIR=/opt/hd1/build/logs          # Log directory
HD1_WORLDS_DIR=/opt/hd1/share/worlds     # Worlds configuration directory
HD1_AVATARS_DIR=/opt/hd1/share/avatars   # Avatars configuration directory

# Logging Configuration  
HD1_LOG_LEVEL=INFO                        # Logging level (TRACE, DEBUG, INFO, WARN, ERROR, FATAL)
HD1_TRACE_MODULES=websocket,sync,threejs  # Comma-separated trace modules

# Sync System
HD1_SYNC_BUFFER_SIZE=1000                 # Operation buffer size
HD1_SYNC_TIMEOUT=30                       # Sync timeout in seconds
```

### Command-Line Flags (All with long options)
```bash
./hd1 --host=127.0.0.1 --port=8081 --daemon --log-level=DEBUG
./hd1 --root-dir=/custom/path --internal-api-base=http://api.internal:8080/api
./hd1 --sync-buffer-size=2000 --version=v1.0.0
```

### .env File Support
Create `.env` file in project root:
```bash
HD1_HOST=127.0.0.1
HD1_PORT=9090
HD1_LOG_LEVEL=DEBUG
HD1_DAEMON=false
```

## Console Features
- **Three.js Stats**: Object count, geometry count, material count
- **Sync Statistics**: Sequence numbers, operation counts, buffer status
- **Entity Controls**: Create boxes/spheres, modify materials, clear scene
- **Avatar Tracking**: Real-time position display and movement controls
- **Scene Management**: Background colors, lighting controls
- **Performance Monitoring**: WebGL renderer statistics

## Three.js Integration
**Development WebGL**: Direct Three.js r170 with zero abstraction layers

### Core Three.js Components
- **Scene**: `new THREE.Scene()` - Root container for all objects
- **Camera**: `new THREE.PerspectiveCamera()` - ViewPort control
- **Renderer**: `new THREE.WebGLRenderer()` - Canvas WebGL context
- **Geometries**: `BoxGeometry`, `SphereGeometry`, `CylinderGeometry`
- **Materials**: `MeshPhongMaterial`, `MeshBasicMaterial`, `MeshStandardMaterial`
- **Lights**: `DirectionalLight`, `AmbientLight`, `PointLight`

### Asset Loading
- **GLB Models**: Native `GLTFLoader` for 3D model loading
- **Textures**: `TextureLoader` for material textures
- **HTTP Delivery**: Static file serving with proper MIME types

### Performance Optimization
- **Object Pooling**: Reuse geometries and materials
- **Frustum Culling**: Automatic by Three.js renderer
- **Level of Detail**: Distance-based model switching
- **Batch Operations**: Grouped sync operations for efficiency

## Logging Standards
**Format**: `timestamp [pid:thread] [level] function.file:line message`  
**Levels**: TRACE (dev), DEBUG (dev), INFO (prod), WARN (prod), ERROR (prod), FATAL (all)  
**Control**: Real-time via `/api/admin/logging/level` and `/api/admin/logging/trace`  

### Log Level Usage
- **TRACE**: Module-specific debugging (websocket, sync, threejs)
- **DEBUG**: Detailed operations for development troubleshooting
- **INFO**: Production events, successful operations, startup/shutdown
- **WARN**: Non-critical issues, missing resources, degraded performance
- **ERROR**: Operation failures, request errors, system problems
- **FATAL**: Critical failures causing process termination

### Runtime Control
```bash
# API-based level changes
curl -X POST /api/admin/logging/level -d '{"level":"DEBUG"}'
curl -X POST /api/admin/logging/trace -d '{"enable":["websocket","sync","threejs"]}'

# Command-line flags
./hd1 --log-level=DEBUG --trace-modules=websocket,sync,threejs

# Environment variables
HD1_LOG_LEVEL=TRACE HD1_TRACE_MODULES=websocket,sync,threejs ./hd1
```

### Thread Safety & Performance
- **Zero-overhead**: Disabled log levels consume ~50ns CPU
- **Thread-safe**: RWMutex with single lock acquisition
- **Structured**: JSON data with consistent field names
- **Context-rich**: operation, session_id, entity_id, sync_sequence details

## Quality Standards
- **Auto-Generated**: Never edit auto_router.go, hd1lib.js, hd1lib.sh
- **Source Files**: Always edit api.yaml, handler implementations
- **Zero Regressions**: All changes maintain compatibility
- **Clean Architecture**: Separation of concerns maintained
- **Single Source of Truth**: All configuration via config system, zero hardcoded values

## Naming Standards
**Atomic Naming Principles**: HD1 follows consistent naming standards with clear semantic meaning:
- **Lowercase Snake Case**: All functions and variables use lowercase_snake_case for consistency
- **Semantic Clarity**: Function names clearly indicate their purpose and behavior
- **Global Variable Visibility**: Package-level exports follow Go conventions
- **No Non-Indicative Patterns**: Clear, descriptive names without ambiguous abbreviations

## Implementation Phases (All Completed)

### Phase 1: Foundation ✅
- Multi-tenant session management with PostgreSQL
- Service registry for dynamic service integration
- JWT-based authentication with refresh tokens
- RESTful API endpoints with Gorilla Mux routing

### Phase 2: Collaboration ✅
- WebRTC infrastructure for P2P communication
- Operational transforms for conflict-free editing
- Asset management with versioning and CDN
- Real-time synchronization via WebSocket

### Phase 3: AI Integration ✅
- LLM avatar system with multi-provider support
- Content generation with template engine
- Natural language processing for commands
- Token usage tracking and rate limiting

### Phase 4: Universal Platform ✅
- Cross-platform client adapters (web, mobile, desktop)
- Plugin architecture with sandboxed execution
- Client registration and capability negotiation
- Enterprise features (RBAC, analytics, security)

---

**HD1 v0.7.0**: The universal 3D interface platform where any service becomes an immersive experience.