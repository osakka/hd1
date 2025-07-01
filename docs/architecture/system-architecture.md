# HD1 (Holodeck One) - System Architecture v2.0

> **API-first 3D visualization platform with three-layer game engine architecture**

## 🏗️ **CURRENT SYSTEM OVERVIEW (2025-07-01)**

HD1 implements a **specification-driven architecture** with **three-layer game engine design** where the OpenAPI 3.0.3 specification (`api.yaml`) serves as the single source of truth for all HTTP APIs, while WebSockets provide real-time broadcasting only.

```
┌─────────────────────────────────────────────────────────────────┐
│                     HD1 SYSTEM ARCHITECTURE v2.0              │
└─────────────────────────────────────────────────────────────────┘

┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Browser   │    │  A-Frame    │    │  WebSocket  │    │ API Client  │
│   WebUI     │    │   WebXR     │    │   Client    │    │ (hd1-client)│
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │                   │
       │ HTTP              │ HTTP              │ WS (broadcast)    │ HTTP
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
│                  THREE-LAYER GAME ENGINE                       │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐│
│  │ENVIRONMENTS │ │    PROPS    │ │   SCENES    │ │  SESSIONS   ││
│  │(4 Physics   │ │(6 Categories│ │(Compositions│ │(User State) ││
│  │ Contexts)   │ │ w/ Physics) │ │ & Scripts)  │ │ Management) ││
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

---

## 🎮 **THREE-LAYER GAME ENGINE ARCHITECTURE**

### **Current Implementation (Revolutionary Achievement)**

HD1 now implements a **complete game engine architecture** matching Unity/Unreal patterns:

```
┌─────────────────────────────────────────────────────────────────┐
│                    THREE-LAYER SYSTEM                          │
├─────────────────────────────────────────────────────────────────┤
│  ENVIRONMENT LAYER                                              │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐│
│  │Earth Surface│ │Molecular    │ │Space Vacuum │ │ Underwater  ││
│  │9.81m/s²     │ │Scale 10⁻⁹m  │ │0.0m/s²      │ │20bar press. ││
│  │1.0 atm      │ │Van der Waals│ │2.7K temp    │ │Dense medium ││
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘│
│                               │                                 │
│  PROPS LAYER                  │                                 │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐│
│  │   Vehicles  │ │ Electronics │ │  Furniture  │ │   Lighting  ││
│  │  Buildings  │ │   Tools     │ │(Auto-adapt │ │(Physics     ││
│  │(YAML-based) │ │(Components) │ │ to environ.)│ │ aware)      ││
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘│
│                               │                                 │
│  SCENE LAYER                  │                                 │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐│
│  │Basic Shapes │ │Ultimate Demo│ │ Custom      │ │ Procedural  ││
│  │(Geometric)  │ │(Full System)│ │ Compositions│ │ Generated   ││
│  │             │ │             │ │             │ │             ││
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

---

## 🔄 **HTTP vs WEBSOCKET ARCHITECTURE (CLARIFIED)**

### **Current Protocol Split**

```
┌─────────────────────────────────────────────────────────────────┐
│                    PROTOCOL SEPARATION                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  HTTP APIs (State Management) - 31 Endpoints                   │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  • Sessions: CREATE, READ, UPDATE, DELETE              │    │
│  │  • Objects: Full CRUD operations                       │    │
│  │  • Environments: Apply physics contexts                │    │
│  │  • Props: Instantiate with physics adaptation          │    │
│  │  • Scenes: Load, fork, save compositions               │    │
│  │  • Admin: Logging, system control                      │    │
│  └─────────────────────────────────────────────────────────┘    │
│                               │                                 │
│  WebSocket (Real-time Broadcasting Only)                       │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  • Session restoration (existing objects on connect)   │    │
│  │  • Object lifecycle events (create, update, delete)    │    │
│  │  • Prop instantiation notifications                    │    │
│  │  • Canvas control commands (rendering updates)         │    │
│  │  • Version synchronization & client info exchange      │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### **Key Principle: HTTP Commands, WebSocket Events**

- **ALL state changes initiated via HTTP APIs**
- **WebSocket provides real-time notification of state changes**
- **Clean separation of concerns**

---

## 🗄️ **SESSION-CENTRIC ARCHITECTURE**

### **Current Data Model**

```
┌─────────────────────────────────────────────────────────────────┐
│                      SESSION STORE                             │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                 Thread-Safe RWMutex                     │    │
│  └─────────────────────────────────────────────────────────┘    │
│                               │                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐          │
│  │  sessions   │    │   objects   │    │   worlds    │          │
│  │[string]*    │    │[sessionId]  │    │[sessionId]  │          │
│  │ Session     │    │[objName]*   │    │*World       │          │
│  │             │    │ Object      │    │(optional)   │          │
│  └─────────────┘    └─────────────┘    └─────────────┘          │
│         │                   │                   │               │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐          │
│  │ ID, Status  │    │Type, Pos,   │    │Size, Camera │          │
│  │ CreatedAt   │    │Color, Scale │    │Transparency │          │
│  │ Metadata    │    │Environment  │    │25×25×25 Grid│          │
│  └─────────────┘    └─────────────┘    └─────────────┘          │
└─────────────────────────────────────────────────────────────────┘
```

### **World Concept Evolution**

**Current Implementation:**
- **World = 3D coordinate system configuration for a session**
- **Optional component** - sessions can exist without worlds
- **Simplified from complex world management to coordinate system config**

```go
type World struct {
    Size         int     `json:"size"`           // Default: 25 (25×25×25 grid)
    Transparency float64 `json:"transparency"`   // Grid visibility
    CameraX      float64 `json:"camera_x"`      // Initial camera position
    CameraY      float64 `json:"camera_y"`
    CameraZ      float64 `json:"camera_z"`
}
```

---

## 📡 **COMPLETE API ENDPOINT MAPPING (31 Endpoints)**

### **Current API Structure**

| **Category** | **Method** | **Endpoint** | **Purpose** |
|--------------|------------|--------------|-------------|
| **Sessions** | `POST` | `/api/sessions` | Create new session |
| | `GET` | `/api/sessions` | List all sessions |
| | `GET` | `/api/sessions/{id}` | Get session details |
| | `DELETE` | `/api/sessions/{id}` | Delete session |
| **Objects** | `POST` | `/api/sessions/{id}/objects` | Create object in session |
| | `GET` | `/api/sessions/{id}/objects` | List session objects |
| | `GET` | `/api/sessions/{id}/objects/{name}` | Get specific object |
| | `PUT` | `/api/sessions/{id}/objects/{name}` | Update object |
| | `DELETE` | `/api/sessions/{id}/objects/{name}` | Delete object |
| **Environments** | `GET` | `/api/environments` | List available physics contexts |
| | `POST` | `/api/environments/{id}` | Apply environment to session |
| **Props** | `GET` | `/api/props` | List available prop categories |
| | `POST` | `/api/sessions/{id}/props/{propId}` | Instantiate prop in session |
| **Scenes** | `GET` | `/api/scenes` | List available scene compositions |
| | `POST` | `/api/scenes/{id}` | Load scene into session |
| | `POST` | `/api/scenes/{id}/fork` | Create scene variation |
| | `POST` | `/api/sessions/{id}/scenes/save` | Save session as scene |
| **Camera** | `PUT` | `/api/sessions/{id}/camera/position` | Set camera position |
| | `POST` | `/api/sessions/{id}/camera/orbit` | Start camera animation |
| **Recording** | `POST` | `/api/sessions/{id}/recording/start` | Begin session recording |
| | `POST` | `/api/sessions/{id}/recording/stop` | End session recording |
| | `POST` | `/api/sessions/{id}/recording/play` | Playback recording |
| | `GET` | `/api/sessions/{id}/recording/status` | Get recording status |
| **Browser** | `POST` | `/api/browser/refresh` | Force browser refresh |
| | `POST` | `/api/browser/canvas` | Canvas control commands |
| **Admin** | `GET` | `/api/admin/logging/config` | Get logging configuration |
| | `POST` | `/api/admin/logging/config` | Update logging settings |
| | `POST` | `/api/admin/logging/level` | Change log level |
| | `POST` | `/api/admin/logging/trace` | Enable module tracing |
| | `GET` | `/api/admin/logging/logs` | Retrieve log entries |
| **System** | `GET` | `/api/version` | Get system version info |

---

## 🔐 **COORDINATE SYSTEM & VALIDATION**

### **Universal Boundaries (ENFORCED)**

```
┌─────────────────────────────────────────────────────────────────┐
│                    COORDINATE VALIDATION                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Universal Bounds: [-12, +12] on all axes                     │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                                                         │    │
│  │     Z                                                   │    │
│  │     ↑                                                   │    │
│  │     │                                                   │    │
│  │     │         25×25×25 Grid                             │    │
│  │     │       (center at origin)                          │    │
│  │     │                                                   │    │
│  │     └────────────────→ X                                │    │
│  │    /                                                    │    │
│  │   /                                                     │    │
│  │  ↙                                                      │    │
│  │ Y                                                       │    │
│  │                                                         │    │
│  │ Validation enforced at:                                 │    │
│  │ • API boundary (api.yaml schema)                        │    │
│  │ • SessionStore level (Go validation)                    │    │
│  │ • Client-side (JavaScript validation)                   │    │
│  │                                                         │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🚀 **WEBSOCKET REACTIVE SYSTEM**

### **Per-Client Session Restoration**

```
┌─────────────────────────────────────────────────────────────────┐
│                    WEBSOCKET FLOW                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Client Connection Lifecycle:                                   │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  1. WebSocket Connect → /ws                             │    │
│  │  2. Session Associate → {"type": "session_associate"}   │    │
│  │  3. Full State Sync → {"type": "canvas_control"}        │    │
│  │  4. Real-time Updates → Event Broadcasting              │    │
│  │  5. Clean Disconnect → Resource cleanup                 │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                 │
│  Session Restoration (FIXED):                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  • EVERY client gets full session state on connect     │    │
│  │  • No global session locks (previous bug fixed)        │    │
│  │  • Browser refresh = complete object restoration       │    │
│  │  • Multiple clients per session supported              │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 💡 **LIGHTING SYSTEM IMPLEMENTATION**

### **Complete API-First Lighting Architecture**

```
┌─────────────────────────────────────────────────────────────────┐
│                      LIGHTING SYSTEM                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Light Types Supported:                                         │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐│
│  │    Point    │ │Directional  │ │   Ambient   │ │    Spot     ││
│  │(Omni-dir)   │ │(Parallel)   │ │ (Global)    │ │(Cone-shaped)││
│  │             │ │             │ │             │ │             ││
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘│
│                                                                 │
│  API Integration:                                               │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  • HTTP API: POST /api/sessions/{id}/objects           │    │
│  │  • Object type: "light"                                │    │
│  │  • Properties: lightType, intensity, color             │    │
│  │  • WebSocket broadcast: real-time light updates        │    │
│  │  • A-Frame rendering: <a-light> component mapping      │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🔧 **BUILD SYSTEM & DEPLOYMENT**

### **Current Build Pipeline**

```
┌─────────────────────────────────────────────────────────────────┐
│                    BUILD SYSTEM v2.0                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  api.yaml (Single Source of Truth)                             │
│      │                                                         │
│      │ make generate                                           │
│      ▼                                                         │
│  Auto-generation Pipeline:                                     │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  • auto_router.go (31 HTTP endpoints)                  │    │
│  │  • hd1-client binary (CLI tool)                        │    │
│  │  • JavaScript API client (hd1lib.js)                   │    │
│  │  • UI components (form system, A-Frame bridge)         │    │
│  │  • Shell functions (bash API wrappers)                 │    │
│  └─────────────────────────────────────────────────────────┘    │
│      │                                                         │
│      │ go build                                                │
│      ▼                                                         │
│  HD1 Binary: /opt/hd1/build/bin/hd1                           │
│      │                                                         │
│      │ make start                                              │
│      ▼                                                         │
│  🚀 HD1 Daemon (PID managed, port 8080)                       │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📊 **ARCHITECTURAL ACHIEVEMENTS**

### **Revolutionary Features Implemented**

1. **Three-Layer Game Engine** - Complete Environment/Props/Scene architecture
2. **Per-Client Session Restoration** - Fixed global session lock issue
3. **Complete Lighting System** - API-first with WebSocket reactivity
4. **31 Auto-Generated Endpoints** - 100% specification-driven
5. **Method Context Binding Resolution** - Fixed JavaScript prototype issues
6. **Clean Protocol Separation** - HTTP for commands, WebSocket for events

### **Evolution from v1.0**

- **From world-centric to session-centric architecture**
- **From basic object management to three-layer game engine**  
- **From incomplete lighting to full A-Frame light integration**
- **From manual routing to 100% auto-generated APIs**
- **From global session locks to per-client restoration**

---

## 🎯 **SYSTEM STATUS: PRODUCTION READY**

**HD1 v2.0 represents a mature, feature-complete 3D visualization platform with game engine capabilities, clean architectural patterns, and robust real-time communication.**

---

*Architecture Document Version: 2.0*  
*HD1 Version: 2.0.0*  
*Last Updated: 2025-07-01*  
*Revolutionary three-layer architecture complete with lighting system!*