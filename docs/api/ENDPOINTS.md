# HD1 API Endpoints Reference

**Complete reference for HD1 v0.7.2 API endpoints**

## 📋 Endpoint Summary

**Total Endpoints**: 15 REST endpoints  
**API Version**: v0.7.2  
**Base URL**: `http://localhost:8080/api`  
**Specification**: `/src/api.yaml`

## 🔄 Sync Operations (4 endpoints)

### 1. Submit Operation
- **Endpoint**: `POST /sync/operations`
- **Purpose**: Submit synchronization operation
- **Handler**: `sync.SubmitOperation`

### 2. Get Missing Operations
- **Endpoint**: `GET /sync/missing/{from}/{to}`
- **Purpose**: Retrieve missing operations in range
- **Handler**: `sync.GetMissingOperations`
- **Parameters**: `from` (start sequence), `to` (end sequence)

### 3. Get Full Sync
- **Endpoint**: `GET /sync/full`
- **Purpose**: Retrieve all operations for full synchronization
- **Handler**: `sync.GetFullSync`

### 4. Get Sync Stats
- **Endpoint**: `GET /sync/stats`
- **Purpose**: Retrieve synchronization statistics
- **Handler**: `sync.GetSyncStats`

## 🎯 Entity Operations (3 endpoints)

### 1. Create Entity
- **Endpoint**: `POST /entities`
- **Purpose**: Create new 3D entity (box, sphere, etc.)
- **Handler**: `entities.CreateEntity`

### 2. Update Entity
- **Endpoint**: `PUT /entities/{entityId}`
- **Purpose**: Update existing entity properties
- **Handler**: `entities.UpdateEntity`
- **Parameters**: `entityId` (entity identifier)

### 3. Delete Entity
- **Endpoint**: `DELETE /entities/{entityId}`
- **Purpose**: Remove entity from scene
- **Handler**: `entities.DeleteEntity`
- **Parameters**: `entityId` (entity identifier)

## 👥 Avatar Operations (5 endpoints)

### 1. Get Avatars
- **Endpoint**: `GET /avatars`
- **Purpose**: Retrieve all active avatars
- **Handler**: `avatars.GetAvatars`

### 2. Create Avatar
- **Endpoint**: `POST /avatars`
- **Purpose**: Create new avatar
- **Handler**: `avatars.CreateAvatar`

### 3. Update Avatar
- **Endpoint**: `PUT /avatars/{avatarId}`
- **Purpose**: Update avatar properties
- **Handler**: `avatars.UpdateAvatar`
- **Parameters**: `avatarId` (avatar identifier)

### 4. Remove Avatar
- **Endpoint**: `DELETE /avatars/{avatarId}`
- **Purpose**: Remove avatar from scene
- **Handler**: `avatars.RemoveAvatar`
- **Parameters**: `avatarId` (avatar identifier)

### 5. Move Avatar
- **Endpoint**: `POST /avatars/{sessionId}/move`
- **Purpose**: Update avatar position and rotation
- **Handler**: `avatars.MoveAvatar`
- **Parameters**: `sessionId` (session identifier)

## 🌍 Scene Operations (2 endpoints)

### 1. Get Scene
- **Endpoint**: `GET /scene`
- **Purpose**: Retrieve scene configuration
- **Handler**: `scene.GetScene`

### 2. Update Scene
- **Endpoint**: `PUT /scene`
- **Purpose**: Update scene properties (background, lighting, etc.)
- **Handler**: `scene.UpdateScene`

## 🔧 System Operations (1 endpoint)

### 1. Get Version
- **Endpoint**: `GET /system/version`
- **Purpose**: Retrieve system version information
- **Handler**: `system.GetVersionHandler`

## 🔗 API Architecture

### Request Flow
```
HTTP Request → API Router → Handler → Business Logic → Response
```

### Auto-Generation
- **Source**: `/src/api.yaml` (OpenAPI 3.0.3)
- **Generated**: `/src/auto_router.go` (DO NOT EDIT)
- **Client**: `/share/htdocs/static/js/hd1lib.js` (auto-generated)

### Authentication
- **Method**: X-Client-ID header
- **Source**: Server-provided client ID via WebSocket

### Response Format
- **Content-Type**: `application/json`
- **CORS**: Enabled for all origins
- **Methods**: GET, POST, PUT, DELETE, OPTIONS

## 📊 Endpoint Categories

| Category | Count | Purpose |
|----------|--------|---------|
| Sync | 4 | Real-time synchronization |
| Entities | 3 | 3D object management |
| Avatars | 5 | Avatar lifecycle management |
| Scene | 2 | Scene configuration |
| System | 1 | System information |
| **Total** | **15** | **Complete API** |

## 🎯 Key Features

### Avatar Lifecycle Management
- Automatic avatar creation on WebSocket connection
- Real-time position updates via `/avatars/{sessionId}/move`
- Automatic cleanup via session inactivity timeout
- Manual removal via DELETE endpoint

### Mobile Support
- Touch controls integrated with avatar movement
- Left screen: movement control
- Right screen: camera control

### Single Source of Truth
- All operations flow through API endpoints
- WebSocket delivers operation results
- No parallel implementation paths

---

**Accuracy Verified**: ✅ All endpoints confirmed against `/src/api.yaml` and `/src/auto_router.go`  
**Last Updated**: 2025-07-18  
**Version**: v0.7.2