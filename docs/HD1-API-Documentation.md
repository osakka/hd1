# HD1 (Holodeck One) API Documentation

## Overview

HD1 is a revolutionary API-first 3D coordinate system with universal world boundaries. It provides a single source of truth for sessions, worlds, objects, and real-time control through a comprehensive REST API.

**Base URL**: `http://localhost:8080/api`  
**API Version**: 1.0.0  
**OpenAPI Specification**: 3.0.3

## Core Principles

- **Everything API-driven** - Zero shell commands required
- **Universal coordinate system** - [-12, +12] on all axes
- **Objects honor world boundaries absolutely**
- **Named object persistence** with full lifecycle management
- **Real-time WebSocket + REST hybrid** architecture

## Architecture

- **Specification drives code generation**
- **Build fails if handlers missing** - ensures complete implementation
- **Zero manual routing** - auto-generated from OpenAPI spec

## Complete API Endpoints

HD1 provides **77 endpoints** across **5 HTTP methods**:

- **GET**: 27 endpoints (data retrieval)
- **POST**: 33 endpoints (creation and actions)
- **PUT**: 12 endpoints (updates)
- **DELETE**: 5 endpoints (removal)

---

## 1. Channel Management

Channels provide collaboration environments for multiple users.

### List Channels
```http
GET /channels
```

**Response Example**:
```json
{
  "success": true,
  "message": "Available channels retrieved successfully",
  "channels": [
    {
      "id": "channel_one",
      "name": "Channel One - Main Collaboration",
      "description": "Primary collaborative environment for general holodeck activities",
      "environment": "channel_one",
      "max_clients": 100
    }
  ]
}
```

### Create Channel
```http
POST /channels
```

**Request Body**:
```json
{
  "name": "Custom Collaboration Channel",
  "description": "Custom environment for specialized workflows",
  "environment": "channel_one",
  "max_clients": 50,
  "enabled": true,
  "priority": 100
}
```

**Response (201)**:
```json
{
  "success": true,
  "message": "Channel created successfully",
  "channel": {
    "id": "channel_custom",
    "name": "Custom Collaboration Channel",
    "description": "Custom environment for specialized workflows",
    "environment": "channel_one",
    "max_clients": 50,
    "enabled": true,
    "created_at": "2025-07-03T10:30:00Z"
  }
}
```

### Get Channel Details
```http
GET /channels/{channelId}
```

**Parameters**:
- `channelId` (path, required): Channel identifier

### Update Channel
```http
PUT /channels/{channelId}
```

### Delete Channel
```http
DELETE /channels/{channelId}
```

---

## 2. Session Management

Sessions represent individual 3D worlds with complete state management.

### List Sessions
```http
GET /sessions
```

**Response**:
```json
{
  "sessions": [
    {
      "id": "session_a1b2c3d4",
      "created": "2025-07-03T10:00:00Z",
      "world_initialized": true,
      "object_count": 15,
      "status": "active"
    }
  ],
  "total": 1,
  "timestamp": "2025-07-03T10:30:00Z"
}
```

### Create Session
```http
POST /sessions
```

**Response (201)**:
```json
{
  "success": true,
  "session_id": "session_a1b2c3d4",
  "created": "2025-07-03T10:00:00Z",
  "message": "Session created - initialize world next"
}
```

### Get Session Details
```http
GET /sessions/{sessionId}
```

**Session ID Pattern**: `^session_[a-f0-9]{8}$`

### Delete Session
```http
DELETE /sessions/{sessionId}
```

---

## 3. Channel Collaboration

Real-time collaboration features within channels.

### Join Channel
```http
POST /sessions/{sessionId}/channel/join
```

### Leave Channel
```http
POST /sessions/{sessionId}/channel/leave
```

### Get Channel Graph
```http
GET /sessions/{sessionId}/channel/graph
```

### Update Channel Graph
```http
PUT /sessions/{sessionId}/channel/graph
```

### Sync with Channel
```http
POST /sessions/{sessionId}/channel/sync
```

### Get Channel Status
```http
GET /sessions/{sessionId}/channel/status
```

---

## 4. Entity Management

Complete entity lifecycle management within sessions.

### List Entities
```http
GET /sessions/{sessionId}/entities
```

### Create Entity
```http
POST /sessions/{sessionId}/entities
```

### Get Entity Details
```http
GET /sessions/{sessionId}/entities/{entityId}
```

### Update Entity
```http
PUT /sessions/{sessionId}/entities/{entityId}
```

### Delete Entity
```http
DELETE /sessions/{sessionId}/entities/{entityId}
```

---

## 5. Component Management

Entity components provide functionality and behavior.

### List Components
```http
GET /sessions/{sessionId}/entities/{entityId}/components
```

### Add Component
```http
POST /sessions/{sessionId}/entities/{entityId}/components
```

### Get Component
```http
GET /sessions/{sessionId}/entities/{entityId}/components/{componentType}
```

### Update Component
```http
PUT /sessions/{sessionId}/entities/{entityId}/components/{componentType}
```

### Delete Component
```http
DELETE /sessions/{sessionId}/entities/{entityId}/components/{componentType}
```

### Bulk Component Operations
```http
POST /sessions/{sessionId}/entities/{entityId}/components/bulk
```

---

## 6. Hierarchy Management

Entity parent-child relationships and transformations.

### Get Parent
```http
GET /sessions/{sessionId}/entities/{entityId}/hierarchy/parent
```

### Set Parent
```http
PUT /sessions/{sessionId}/entities/{entityId}/hierarchy/parent
```

### Get Children
```http
GET /sessions/{sessionId}/entities/{entityId}/hierarchy/children
```

### Get Transforms
```http
GET /sessions/{sessionId}/entities/{entityId}/hierarchy/transforms
```

### Update Transforms
```http
PUT /sessions/{sessionId}/entities/{entityId}/hierarchy/transforms
```

### Get Hierarchy Tree
```http
GET /sessions/{sessionId}/entities/hierarchy/tree
```

---

## 7. Lifecycle Management

Entity lifecycle state management.

### Enable Entity
```http
PUT /sessions/{sessionId}/entities/{entityId}/lifecycle/enable
```

### Disable Entity
```http
PUT /sessions/{sessionId}/entities/{entityId}/lifecycle/disable
```

### Activate Entity
```http
POST /sessions/{sessionId}/entities/{entityId}/lifecycle/activate
```

### Deactivate Entity
```http
POST /sessions/{sessionId}/entities/{entityId}/lifecycle/deactivate
```

### Destroy Entity
```http
DELETE /sessions/{sessionId}/entities/{entityId}/lifecycle/destroy
```

### Get Lifecycle Status
```http
GET /sessions/{sessionId}/entities/{entityId}/lifecycle/status
```

### Bulk Lifecycle Operations
```http
POST /sessions/{sessionId}/entities/lifecycle/bulk
```

---

## 8. Scene Management

Scene state and hierarchy management.

### Get Scene Hierarchy
```http
GET /sessions/{sessionId}/scene/hierarchy
```

### Update Scene Hierarchy
```http
PUT /sessions/{sessionId}/scene/hierarchy
```

### Get Scene State
```http
GET /sessions/{sessionId}/scene/state
```

### Update Scene State
```http
PUT /sessions/{sessionId}/scene/state
```

### Save Scene State
```http
POST /sessions/{sessionId}/scene/state/save
```

### Load Scene State
```http
POST /sessions/{sessionId}/scene/state/load
```

### Reset Scene State
```http
POST /sessions/{sessionId}/scene/state/reset
```

### List Scenes
```http
GET /sessions/{sessionId}/scenes
```

### Create Scene
```http
POST /sessions/{sessionId}/scenes
```

### Activate Scene
```http
POST /sessions/{sessionId}/scenes/{sceneId}/activate
```

### Export Scene
```http
GET /sessions/{sessionId}/scene/export
```

### Import Scene
```http
POST /sessions/{sessionId}/scene/import
```

---

## 9. Animation System

Real-time animation control and management.

### List Animations
```http
GET /sessions/{sessionId}/animations
```

### Create Animation
```http
POST /sessions/{sessionId}/animations
```

### Play Animation
```http
POST /sessions/{sessionId}/animations/{animationId}/play
```

### Stop Animation
```http
POST /sessions/{sessionId}/animations/{animationId}/stop
```

---

## 10. Physics System

Physics world configuration and rigid body management.

### Get Physics World
```http
GET /sessions/{sessionId}/physics/world
```

### Update Physics World
```http
PUT /sessions/{sessionId}/physics/world
```

### List Rigid Bodies
```http
GET /sessions/{sessionId}/physics/rigidbodies
```

### Apply Force to Rigid Body
```http
POST /sessions/{sessionId}/physics/rigidbodies/{entityId}/force
```

---

## 11. Audio System

3D audio source management and control.

### List Audio Sources
```http
GET /sessions/{sessionId}/audio/sources
```

### Create Audio Source
```http
POST /sessions/{sessionId}/audio/sources
```

### Play Audio
```http
POST /sessions/{sessionId}/audio/sources/{audioId}/play
```

### Stop Audio
```http
POST /sessions/{sessionId}/audio/sources/{audioId}/stop
```

---

## 12. Recording System

Session recording and playback functionality.

### Start Recording
```http
POST /sessions/{sessionId}/recording/start
```

### Stop Recording
```http
POST /sessions/{sessionId}/recording/stop
```

### Play Recording
```http
POST /sessions/{sessionId}/recording/play
```

### Get Recording Status
```http
GET /sessions/{sessionId}/recording/status
```

---

## 13. Camera System

Camera positioning and control.

### Update Camera Position
```http
PUT /sessions/{sessionId}/camera/position
```

### Set Camera Orbit
```http
POST /sessions/{sessionId}/camera/orbit
```

---

## 14. Admin & Logging

System administration and logging configuration.

### Get Logging Config
```http
GET /admin/logging/config
```

### Update Logging Config
```http
POST /admin/logging/config
```

### Set Log Level
```http
POST /admin/logging/level
```

**Request Body**:
```json
{
  "level": "DEBUG"
}
```

### Configure Trace Modules
```http
POST /admin/logging/trace
```

**Request Body**:
```json
{
  "modules": ["sessions", "objects", "websocket"]
}
```

### Get Logs
```http
GET /admin/logging/logs
```

---

## 15. Browser Integration

Browser refresh and canvas management.

### Refresh Browser
```http
POST /browser/refresh
```

### Update Canvas
```http
POST /browser/canvas
```

---

## 16. System Information

### Get Version
```http
GET /version
```

**Response**:
```json
{
  "version": "1.0.0",
  "build": "aa74f3f3",
  "timestamp": "2025-07-03T10:30:00Z"
}
```

---

## Error Handling

All endpoints return consistent error responses:

### 400 Bad Request
```json
{
  "success": false,
  "message": "Invalid request parameters",
  "errors": [
    "Field 'name' is required",
    "Field 'max_clients' must be positive"
  ]
}
```

### 404 Not Found
```json
{
  "success": false,
  "message": "Resource not found",
  "resource": "session_invalid123"
}
```

### 409 Conflict
```json
{
  "success": false,
  "message": "Resource already exists",
  "conflict": "Channel ID already exists"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "message": "Internal server error",
  "error_id": "error_abc123"
}
```

---

## Data Schemas

### Coordinate System
All coordinates use the universal boundary system:
- **X, Y, Z axes**: [-12, +12] range
- **Scale units**: nm, μm, mm, cm, m, km, Mm, Gm
- **Gravity**: 0.0 to 100.0 m/s²
- **Atmosphere**: air, vacuum, thin_air, liquid

### Session ID Format
Sessions use the pattern: `^session_[a-f0-9]{8}$`

Example: `session_a1b2c3d4`

### Entity Status
Entity lifecycle states:
- `active` - Entity is running and interactive
- `inactive` - Entity is paused
- `error` - Entity has encountered an error

### Environment Properties
```json
{
  "id": "channel_one",
  "name": "Earth Surface",
  "description": "Standard terrestrial environment",
  "scale_unit": "m",
  "gravity": 9.8,
  "atmosphere": "air",
  "boundaries": {
    "x": [-12, 12],
    "y": [-12, 12],
    "z": [-12, 12]
  }
}
```

---

## Authentication

Currently, HD1 API operates without authentication for development environments. All endpoints are accessible without tokens or credentials.

**Note**: Production deployments should implement proper authentication and authorization mechanisms.

---

## Rate Limiting

No rate limiting is currently implemented. All endpoints can be called without restrictions.

**Note**: Production deployments should implement appropriate rate limiting based on use case requirements.

---

## Real-time Features

HD1 provides WebSocket integration for real-time state synchronization:

- **WebSocket endpoint**: `ws://localhost:8080/ws`
- **Message format**: JSON
- **Automatic reconnection**: Implemented with exponential backoff
- **State synchronization**: Complete session state broadcast on connect

---

## Development Tools

### Build System
```bash
cd src && make clean && make
```

### Start Server
```bash
cd src && make start
```

### HD1 Client CLI
```bash
./build/bin/hd1-client --help
```

### Test Endpoints
```bash
# List environments
./build/bin/hd1-client list-environments

# List props
./build/bin/hd1-client list-props

# Test API endpoints
curl -X GET http://localhost:8080/api/version
curl -X POST http://localhost:8080/api/sessions
```

---

## Implementation Details

### Auto-generated Routing
All endpoints are auto-generated from the OpenAPI specification:
- **Go handlers**: Located in `src/api/` directory
- **Route generation**: Automatic from specification
- **Build validation**: Prevents deployment of incomplete implementations

### Handler Structure
Each endpoint maps to a specific Go handler file:
- **Pattern**: `x-handler: "api/path/handler_name.go"`
- **Function**: `x-function: "HandlerFunction"`
- **Validation**: Build system ensures all handlers exist

### Template System
HD1 uses an externalized template system:
- **Location**: `src/codegen/templates/`
- **Languages**: Go, JavaScript, Shell
- **Embedding**: Go embed for single binary deployment

---

## Getting Started

1. **Start the HD1 server**:
   ```bash
   cd /opt/hd1/src && make start
   ```

2. **Create a session**:
   ```bash
   curl -X POST http://localhost:8080/api/sessions
   ```

3. **Get session details**:
   ```bash
   curl -X GET http://localhost:8080/api/sessions/{sessionId}
   ```

4. **Create entities**:
   ```bash
   curl -X POST http://localhost:8080/api/sessions/{sessionId}/entities \
     -H "Content-Type: application/json" \
     -d '{"name": "test_entity", "type": "cube"}'
   ```

5. **Access web interface**:
   Open `http://localhost:8080` in your browser

---

## Support

For technical support and development questions:
- **Contact**: Holodeck One 3D Engine
- **Documentation**: Located in `/opt/hd1/docs/`
- **Source**: Available at current working directory
- **Build logs**: Check `/opt/hd1/build/logs/`

---

*Generated on 2025-07-03 | HD1 API v1.0.0*