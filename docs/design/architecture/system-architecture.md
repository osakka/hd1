# HD1 (Holodeck One) - System Architecture

> **Advanced API-first 3D visualization with specification-driven development**

## 🏗️ **SYSTEM OVERVIEW**

HD1 implements a **specification-driven architecture** where the OpenAPI 3.0.3 specification (`api.yaml`) serves as the single source of truth for all routing, validation, and API behavior.

```
┌─────────────────────────────────────────────────────────────────┐
│                     HD1 SYSTEM ARCHITECTURE                    │
└─────────────────────────────────────────────────────────────────┘

┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Client    │    │   Browser   │    │  WebSocket  │    │ API Client  │
│ WebGL/HTML  │    │   React     │    │   Client    │    │   (curl)    │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │                   │
       │ HTTP              │ HTTP              │ WS                │ HTTP
       │                   │                   │                   │
┌─────────────────────────────────────────────────────────────────┐
│                      HD1 DAEMON (main.go)                      │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐│
│  │   Static    │ │    Home     │ │  WebSocket  │ │ Auto-Router ││
│  │   Server    │ │   Handler   │ │    Hub      │ │ (Generated) ││
│  │ /static/*   │ │     /       │ │    /ws      │ │   /api/*    ││
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘│
└─────────────────────────────────────────────────────────────────┘
                                      │
                                      │
┌─────────────────────────────────────────────────────────────────┐
│                    CORE COMPONENTS                             │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐│
│  │   Session   │ │   WebSocket │ │    API      │ │   Build     ││
│  │    Store    │ │     Hub     │ │  Handlers   │ │   System    ││
│  │(Thread-Safe)│ │(Real-time)  │ │(Generated)  │ │(Makefile)   ││
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

---

## 🔄 **REQUEST LIFECYCLE & COMPONENT HANDOFFS**

### **1. API Request Flow**

```
[Client] → [HD1 Daemon] → [Auto-Router] → [Handler] → [SessionStore] → [Response]
    │           │              │             │            │              │
    │           │              │             │            │              │
    HTTP        │              │             │            │              JSON
    Request     │              │             │            │              Response
                │              │             │            │              │
            Port 8080      Pattern         Impl.      Thread-Safe     Status Code
                          Matching       Function      Operations     + Data
```

#### **Detailed Request Steps:**

1. **Client Sends Request**
   ```
   POST /api/sessions
   Content-Type: application/json
   ```

2. **HD1 Daemon Receives** (`main.go:line 83`)
   ```go
   http.Handle("/api/", apiRouter)
   ```

3. **Auto-Router Processes** (`auto_router.go:line 53`)
   ```go
   path := strings.TrimPrefix(req.URL.Path, "/api")
   // Matches against generated routes from api.yaml
   ```

4. **Handler Executes** (`api/sessions/create.go:line 10`)
   ```go
   func CreateSessionHandler(w http.ResponseWriter, r *http.Request, hub interface{})
   ```

5. **SessionStore Operation** (`server/hub.go:line 117`)
   ```go
   session := h.GetStore().CreateSession()
   ```

6. **WebSocket Broadcast** (`server/hub.go:line 317`)
   ```go
   h.BroadcastUpdate("session_created", data)
   ```

7. **JSON Response Sent**
   ```json
   {
     "success": true,
     "session_id": "session_a1b2c3d4",
     "status": "active"
   }
   ```

---

## 🤖 **AUTO-GENERATED ROUTING SYSTEM**

### **Specification-Driven Architecture**

HD1's advanced approach auto-generates ALL routing from the OpenAPI specification:

```
api.yaml (Single Source of Truth)
    │
    │ make generate
    ▼
auto_router.go (Generated Code)
    │
    │ imports
    ▼
Handler Functions (api/*/*)
    │
    │ dependency injection
    ▼
SessionStore Operations
    │
    │ real-time updates
    ▼
WebSocket Broadcasts
```

#### **Code Generation Process:**

1. **Specification Parse** (`codegen/generator.go`)
   ```go
   // Loads api.yaml and extracts:
   // - Paths and operations
   // - Handler file locations (x-handler)
   // - Function names (x-function)
   ```

2. **Route Generation**
   ```go
   // Generated in auto_router.go:
   {
       Path: "/sessions",
       Method: "POST", 
       Handler: sessions.CreateSessionHandler,
       OperationID: "createSession"
   }
   ```

3. **Build-Time Validation**
   ```bash
   make generate  # Fails if handlers missing
   make build     # Validates all routes have implementations
   ```

#### **Handler Binding Pattern:**

```go
// From api.yaml:
// x-handler: "api/sessions/create.go"
// x-function: "CreateSession"

// Generated router code:
routes = append(routes, Route{
    Path: "/sessions",
    Method: "POST",
    Handler: func(w http.ResponseWriter, r *http.Request) {
        sessions.CreateSessionHandler(w, r, router.hub)
    },
})
```

---

## 🌐 **WEBSOCKET HUB ARCHITECTURE**

### **Real-Time Communication Flow**

```
┌─────────────────────────────────────────────────────────────────┐
│                      WEBSOCKET HUB                             │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐          │
│  │   Client    │    │   Client    │    │   Client    │          │
│  │  Connection │    │  Connection │    │  Connection │          │
│  └─────────────┘    └─────────────┘    └─────────────┘          │
│         │                   │                   │               │
│         │                   │                   │               │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                     HUB CHANNELS                       │    │
│  │  register   │  unregister  │  broadcast               │    │
│  │    chan     │     chan     │    chan                  │    │
│  └─────────────────────────────────────────────────────────┘    │
│                               │                                 │
│         ┌─────────────────────┼─────────────────────┐           │
│         │                     ▼                     │           │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐          │
│  │  Session    │    │    Real     │    │   Client    │          │
│  │   Store     │    │    Time     │    │ Management  │          │
│  │  Updates    │    │ Broadcast   │    │   (ping/    │          │
│  │             │    │             │    │   pong)     │          │
│  └─────────────┘    └─────────────┘    └─────────────┘          │
└─────────────────────────────────────────────────────────────────┘
```

#### **Hub Operations:**

1. **Client Registration** (`server/hub.go:line 33`)
   ```go
   case client := <-h.register:
       h.clients[client] = true
       log.Printf("Client connected. Total: %d", len(h.clients))
   ```

2. **Broadcast Distribution** (`server/hub.go:line 54`)
   ```go
   case message := <-h.broadcast:
       for client := range h.clients {
           select {
           case client.send <- message:
           default:
               close(client.send)
               delete(h.clients, client)
           }
       }
   ```

3. **Real-Time Updates** (`server/hub.go:line 317`)
   ```go
   func (h *Hub) BroadcastUpdate(updateType string, data interface{}) {
       update := map[string]interface{}{
           "type": updateType,
           "data": data,
           "timestamp": time.Now().Unix(),
       }
       jsonData, _ := json.Marshal(update)
       h.BroadcastMessage(jsonData)
   }
   ```

---

## 🗄️ **SESSION MANAGEMENT ARCHITECTURE**

### **Thread-Safe SessionStore**

```
┌─────────────────────────────────────────────────────────────────┐
│                      SESSION STORE                             │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    RWMutex Lock                         │    │
│  │           (Thread-Safe Operations)                      │    │
│  └─────────────────────────────────────────────────────────┘    │
│                               │                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐          │
│  │  sessions   │    │   objects   │    │   worlds    │          │
│  │    map      │    │     map     │    │     map     │          │
│  │[string]*    │    │[sessionId]  │    │[sessionId]  │          │
│  │ Session     │    │[objName]*   │    │  *World     │          │
│  │             │    │  Object     │    │             │          │
│  └─────────────┘    └─────────────┘    └─────────────┘          │
│         │                   │                   │               │
│         │                   │                   │               │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐          │
│  │  Session    │    │   Object    │    │   World     │          │
│  │ Lifecycle   │    │ Lifecycle   │    │Coordinate   │          │
│  │   CRUD      │    │   CRUD      │    │ System      │          │
│  │             │    │ Validation  │    │[-12,+12]    │          │
│  └─────────────┘    └─────────────┘    └─────────────┘          │
└─────────────────────────────────────────────────────────────────┘
```

#### **SessionStore Operations:**

1. **Session Creation** (`server/hub.go:line 117`)
   ```go
   func (s *SessionStore) CreateSession() *Session {
       s.mutex.Lock()
       defer s.mutex.Unlock()
       
       sessionID := generateSessionID()
       session := &Session{
           ID: sessionID,
           CreatedAt: time.Now(),
           Status: "active",
       }
       
       s.sessions[sessionID] = session
       s.objects[sessionID] = make(map[string]*Object)
       return session
   }
   ```

2. **Object Management** (`server/hub.go:line 172`)
   ```go
   func (s *SessionStore) CreateObject(sessionID, objectName, objectType string, x, y, z float64) (*Object, error) {
       // Coordinate validation: [-12, +12]
       if x < -12 || x > 12 || y < -12 || y > 12 || z < -12 || z > 12 {
           return nil, &CoordinateError{Message: "Coordinates must be within [-12, +12] bounds"}
       }
       // Thread-safe object creation
   }
   ```

3. **World Initialization** (`server/hub.go:line 282`)
   ```go
   func (s *SessionStore) InitializeWorld(sessionID string, size int, transparency float64, cameraX, cameraY, cameraZ float64) (*World, error) {
       world := &World{
           Size: size,           // 25x25x25 grid
           Transparency: transparency,
           CameraX: cameraX,
           CameraY: cameraY,
           CameraZ: cameraZ,
       }
       s.worlds[sessionID] = world
       return world, nil
   }
   ```

---

## 🔧 **BUILD SYSTEM ARCHITECTURE**

### **Standard Development Pipeline**

```
┌─────────────────────────────────────────────────────────────────┐
│                    BUILD SYSTEM FLOW                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  api.yaml (Specification)                                      │
│      │                                                         │
│      │ make validate                                           │
│      ▼                                                         │
│  ✅ Specification Valid                                         │
│      │                                                         │
│      │ make generate                                           │
│      ▼                                                         │
│  auto_router.go (Generated)                                    │
│      │                                                         │
│      │ Handler Validation                                      │
│      ▼                                                         │
│  ✅ All Handlers Exist                                          │
│      │                                                         │
│      │ go build                                                │
│      ▼                                                         │
│  thd (Binary)                                                  │
│      │                                                         │
│      │ make start                                              │
│      ▼                                                         │
│  🚀 HD1 Daemon Running                                          │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

#### **Build Validation Chain:**

1. **Specification Validation**
   ```bash
   make validate  # Checks api.yaml syntax
   ```

2. **Code Generation**
   ```bash
   make generate  # Creates auto_router.go from api.yaml
   ```

3. **Handler Validation**
   ```go
   // Build fails if handlers missing for any route
   // x-handler and x-function must exist
   ```

4. **Standard Binary Creation**
   ```bash
   make build  # Creates /opt/holodeck-one/build/bin/thd
   ```

5. **Daemon Management**
   ```bash
   make start   # Standard daemon with PID management
   make stop    # Clean shutdown with resource cleanup
   make status  # Complete health reporting
   ```

---

## 📡 **API ENDPOINT MAPPING**

### **Complete Route → Handler → Store Mapping**

| Method | Endpoint | Handler | Store Operation | WebSocket Event |
|--------|----------|---------|-----------------|-----------------|
| **Sessions Management** |
| `POST` | `/api/sessions` | `sessions.CreateSessionHandler` | `CreateSession()` | `session_created` |
| `GET` | `/api/sessions` | `sessions.ListSessionsHandler` | `ListSessions()` | - |
| `GET` | `/api/sessions/{id}` | `sessions.GetSessionHandler` | `GetSession()` | - |
| `DELETE` | `/api/sessions/{id}` | `sessions.DeleteSessionHandler` | `DeleteSession()` | `session_deleted` |
| **World Management** |
| `POST` | `/api/sessions/{id}/world` | `world.InitializeWorldHandler` | `InitializeWorld()` | `world_initialized` |
| `GET` | `/api/sessions/{id}/world` | `world.GetWorldSpecHandler` | `GetWorld()` | - |
| **Object Management** |
| `POST` | `/api/sessions/{id}/objects` | `objects.CreateObjectHandler` | `CreateObject()` | `object_created` |
| `GET` | `/api/sessions/{id}/objects` | `objects.ListObjectsHandler` | `ListObjects()` | - |
| `GET` | `/api/sessions/{id}/objects/{name}` | `objects.GetObjectHandler` | `GetObject()` | - |
| `PUT` | `/api/sessions/{id}/objects/{name}` | `objects.UpdateObjectHandler` | `UpdateObject()` | `object_updated` |
| `DELETE` | `/api/sessions/{id}/objects/{name}` | `objects.DeleteObjectHandler` | `DeleteObject()` | `object_deleted` |
| **Camera Control** |
| `PUT` | `/api/sessions/{id}/camera/position` | `camera.SetCameraPositionHandler` | - | `camera_moved` |
| `POST` | `/api/sessions/{id}/camera/orbit` | `camera.StartCameraOrbitHandler` | - | `camera_orbit` |
| **Browser Control** |
| `POST` | `/api/browser/refresh` | `browser.ForceRefreshHandler` | - | `force_refresh` |
| `POST` | `/api/browser/canvas` | `browser.SetCanvasHandler` | - | `canvas_update` |

---

## 🔐 **SECURITY ARCHITECTURE**

### **Input Validation & Boundary Enforcement**

```
┌─────────────────────────────────────────────────────────────────┐
│                      SECURITY LAYERS                           │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────────┐    │
│  │               API BOUNDARY VALIDATION                   │    │
│  │  • OpenAPI 3.0.3 Schema Validation                     │    │
│  │  • Request Format Checking                             │    │
│  │  • Parameter Type Validation                           │    │
│  └─────────────────────────────────────────────────────────┘    │
│                               │                                 │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │              COORDINATE VALIDATION                      │    │
│  │  • Universal Bounds: [-12, +12] on all axes           │    │
│  │  • Enforced at SessionStore Level                      │    │
│  │  • No Object Creation Outside Bounds                   │    │
│  └─────────────────────────────────────────────────────────┘    │
│                               │                                 │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                SESSION ISOLATION                       │    │
│  │  • Thread-Safe Operations (RWMutex)                    │    │
│  │  • Per-Session Object Stores                           │    │
│  │  • Clean Resource Management                           │    │
│  └─────────────────────────────────────────────────────────┘    │
│                               │                                 │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                 WEBSOCKET SECURITY                     │    │
│  │  • Origin Checking (configurable)                      │    │
│  │  • Connection Limits                                   │    │
│  │  • Proper Ping/Pong Handling                          │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🚀 **DEPLOYMENT ARCHITECTURE**

### **Standard Daemon Management**

```
┌─────────────────────────────────────────────────────────────────┐
│                    PRODUCTION DEPLOYMENT                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Process Management:                                            │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  HD1 Daemon (PID Management)                           │    │
│  │  • Absolute Path Configuration                         │    │
│  │  • Standard Logging                                │    │
│  │  • Clean Shutdown Procedures                           │    │
│  │  • Resource Cleanup on Exit                            │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                 │
│  Directory Structure:                                           │
│  /opt/holodeck-one/                                               │
│  ├── src/              # Source code & build system            │
│  ├── build/bin/thd     # Standard daemon binary            │
│  ├── build/logs/       # Timestamped, structured logs          │
│  ├── build/runtime/    # PID files, runtime data              │
│  └── share/htdocs/     # Static web assets                     │
│                                                                 │
│  Network Architecture:                                         │
│  Port 8080 (configurable):                                     │
│  ├── /                 # Web interface                         │
│  ├── /ws               # WebSocket connections                 │
│  ├── /api/*            # Auto-generated API routes            │
│  └── /static/*         # Static file serving                  │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📊 **MONITORING & OBSERVABILITY**

### **Health Check Architecture**

```
┌─────────────────────────────────────────────────────────────────┐
│                     MONITORING POINTS                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  System Health:                                                 │
│  • make status     → Daemon PID validation                     │
│  • Port 8080       → Network listener status                   │
│  • API Endpoints   → Core functionality testing                │
│  • WebSocket Hub   → Client connection monitoring              │
│                                                                 │
│  Application Metrics:                                           │
│  • Session Count   → Active visualization sessions             │
│  • Object Count    → Per-session 3D objects                    │
│  • Client Count    → WebSocket connections                     │
│  • Memory Usage    → SessionStore and Hub resources            │
│                                                                 │
│  Performance Monitoring:                                        │
│  • Request Latency → API response times                        │
│  • WebSocket RTT   → Real-time update delays                   │
│  • Coordinate Val. → Boundary enforcement performance          │
│  • Build Times     → Code generation efficiency                │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🔮 **INNOVATION HIGHLIGHTS**

### **Advanced Architectural Patterns**

1. **Specification-Driven Development**
   - OpenAPI 3.0.3 as executable specification
   - Auto-generated routing eliminates manual configuration
   - Build-time validation prevents incomplete deployments

2. **Hybrid Real-Time Architecture**
   - RESTful APIs for state management
   - WebSocket hub for real-time updates
   - Unified session store for both protocols

3. **Standard Engineering Standards**
   - Absolute path configuration throughout
   - Thread-safe operations with proper locking
   - Clean daemon lifecycle management

4. **Universal Coordinate System**
   - Fixed 25×25×25 grid with universal boundaries
   - Enforced validation at API boundary
   - Consistent 3D visualization standards

---

**HD1 represents the perfect fusion of innovative 3D visualization capabilities with standard software engineering practices, delivering a reliable, scalable, and maintainable system for real-time collaborative 3D environments.**

---

*Architecture Document Version: 1.0*  
*HD1 Version: 2.0.0*  
*Last Updated: 2025-06-28*