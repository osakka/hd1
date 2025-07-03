# HD1 v3.0 PlayCanvas API Design Specification

> **Revolutionary Game Engine as a Service - API-First PlayCanvas Integration**

## 🎯 **API DESIGN PHILOSOPHY**

**Core Principle:** Every PlayCanvas game engine feature accessible via clean HTTP REST endpoints.

**Session-Centric Architecture:** All game engine operations scoped within HD1 sessions for multi-tenancy.

---

## 📋 **PHASE 1: ENTITY MANAGEMENT APIS**

### **1. Entity Lifecycle Management**

#### **Create Entity**
```yaml
POST /api/sessions/{sessionId}/entities
Content-Type: application/json

{
  "name": "MyEntity",
  "tags": ["player", "dynamic"],
  "enabled": true,
  "position": {"x": 0, "y": 0, "z": 0},
  "rotation": {"x": 0, "y": 0, "z": 0, "w": 1},
  "scale": {"x": 1, "y": 1, "z": 1}
}

Response: 201 Created
{
  "entity_id": "entity-abc123",
  "name": "MyEntity",
  "guid": "playcanvas-guid-xyz789",
  "created_at": "2025-07-01T16:10:00Z"
}
```

#### **List Entities**
```yaml
GET /api/sessions/{sessionId}/entities
Query Parameters:
  - tag: Filter by tag
  - enabled: Filter by enabled state
  - limit: Maximum results (default: 100)

Response: 200 OK
{
  "entities": [
    {
      "entity_id": "entity-abc123",
      "name": "MyEntity",
      "enabled": true,
      "tags": ["player", "dynamic"],
      "transform": {
        "position": {"x": 0, "y": 0, "z": 0},
        "rotation": {"x": 0, "y": 0, "z": 0, "w": 1},
        "scale": {"x": 1, "y": 1, "z": 1}
      },
      "components": ["model", "rigidbody"]
    }
  ],
  "total": 1
}
```

#### **Get Entity Details**
```yaml
GET /api/sessions/{sessionId}/entities/{entityId}

Response: 200 OK
{
  "entity_id": "entity-abc123",
  "name": "MyEntity",
  "enabled": true,
  "tags": ["player", "dynamic"],
  "transform": {
    "position": {"x": 0, "y": 0, "z": 0},
    "rotation": {"x": 0, "y": 0, "z": 0, "w": 1},
    "scale": {"x": 1, "y": 1, "z": 1}
  },
  "components": {
    "model": {
      "asset": "model-asset-id",
      "type": "asset"
    },
    "rigidbody": {
      "type": "dynamic",
      "mass": 1.0
    }
  },
  "children": ["entity-child-1", "entity-child-2"],
  "parent": "entity-parent-1"
}
```

#### **Update Entity**
```yaml
PUT /api/sessions/{sessionId}/entities/{entityId}
Content-Type: application/json

{
  "name": "UpdatedEntity",
  "enabled": false,
  "tags": ["player", "static"],
  "position": {"x": 5, "y": 2, "z": -3},
  "rotation": {"x": 0, "y": 45, "z": 0, "w": 1}
}

Response: 200 OK
{
  "entity_id": "entity-abc123",
  "updated_at": "2025-07-01T16:15:00Z",
  "changes": ["name", "enabled", "tags", "position", "rotation"]
}
```

#### **Delete Entity**
```yaml
DELETE /api/sessions/{sessionId}/entities/{entityId}
Query Parameters:
  - cascade: Delete children (default: false)

Response: 204 No Content
```

---

## 📋 **PHASE 2: COMPONENT SYSTEM APIS**

### **2. Component Management**

#### **Add Component**
```yaml
POST /api/sessions/{sessionId}/entities/{entityId}/components
Content-Type: application/json

{
  "type": "model",
  "properties": {
    "asset": "model-asset-id",
    "type": "asset",
    "castShadows": true,
    "receiveShadows": true
  }
}

Response: 201 Created
{
  "component_type": "model",
  "entity_id": "entity-abc123",
  "created_at": "2025-07-01T16:20:00Z"
}
```

#### **Update Component**
```yaml
PUT /api/sessions/{sessionId}/entities/{entityId}/components/{componentType}
Content-Type: application/json

{
  "properties": {
    "asset": "new-model-asset-id",
    "castShadows": false
  }
}

Response: 200 OK
{
  "component_type": "model",
  "entity_id": "entity-abc123",
  "updated_at": "2025-07-01T16:25:00Z"
}
```

#### **Remove Component**
```yaml
DELETE /api/sessions/{sessionId}/entities/{entityId}/components/{componentType}

Response: 204 No Content
```

### **3. Transform Operations**

#### **Set Position**
```yaml
PUT /api/sessions/{sessionId}/entities/{entityId}/transform/position
Content-Type: application/json

{
  "x": 10,
  "y": 5,
  "z": -2
}

Response: 200 OK
{
  "entity_id": "entity-abc123",
  "position": {"x": 10, "y": 5, "z": -2},
  "updated_at": "2025-07-01T16:30:00Z"
}
```

#### **Set Rotation**
```yaml
PUT /api/sessions/{sessionId}/entities/{entityId}/transform/rotation
Content-Type: application/json

{
  "x": 0,
  "y": 90,
  "z": 0
}

Response: 200 OK
{
  "entity_id": "entity-abc123", 
  "rotation": {"x": 0, "y": 90, "z": 0, "w": 1},
  "updated_at": "2025-07-01T16:35:00Z"
}
```

#### **Set Scale**
```yaml
PUT /api/sessions/{sessionId}/entities/{entityId}/transform/scale
Content-Type: application/json

{
  "x": 2,
  "y": 2,
  "z": 2
}

Response: 200 OK
{
  "entity_id": "entity-abc123",
  "scale": {"x": 2, "y": 2, "z": 2},
  "updated_at": "2025-07-01T16:40:00Z"
}
```

---

## 📋 **PHASE 3: SCENE GRAPH APIS**

### **4. Hierarchy Management**

#### **Set Parent**
```yaml
PUT /api/sessions/{sessionId}/entities/{entityId}/parent
Content-Type: application/json

{
  "parent_id": "entity-parent-123"
}

Response: 200 OK
{
  "entity_id": "entity-abc123",
  "parent_id": "entity-parent-123",
  "updated_at": "2025-07-01T16:45:00Z"
}
```

#### **Add Child**
```yaml
POST /api/sessions/{sessionId}/entities/{entityId}/children
Content-Type: application/json

{
  "child_id": "entity-child-456"
}

Response: 201 Created
{
  "parent_id": "entity-abc123",
  "child_id": "entity-child-456",
  "created_at": "2025-07-01T16:50:00Z"
}
```

#### **Remove Child**
```yaml
DELETE /api/sessions/{sessionId}/entities/{entityId}/children/{childId}

Response: 204 No Content
```

---

## 📋 **PHASE 4: ASSET MANAGEMENT APIS**

### **5. Asset Operations**

#### **Upload Asset**
```yaml
POST /api/sessions/{sessionId}/assets
Content-Type: multipart/form-data

Form Data:
  - file: [binary file data]
  - type: "model" | "texture" | "audio" | "script"
  - name: "MyAsset"

Response: 201 Created
{
  "asset_id": "asset-abc123",
  "name": "MyAsset",
  "type": "model",
  "size": 2048576,
  "url": "/static/assets/session-123/asset-abc123.glb",
  "created_at": "2025-07-01T16:55:00Z"
}
```

#### **List Assets**
```yaml
GET /api/sessions/{sessionId}/assets
Query Parameters:
  - type: Filter by asset type
  - limit: Maximum results

Response: 200 OK
{
  "assets": [
    {
      "asset_id": "asset-abc123",
      "name": "MyAsset",
      "type": "model",
      "size": 2048576,
      "url": "/static/assets/session-123/asset-abc123.glb"
    }
  ],
  "total": 1
}
```

---

## 📋 **PHASE 5: ANIMATION APIS**

### **6. Animation Control**

#### **Play Animation**
```yaml
POST /api/sessions/{sessionId}/entities/{entityId}/animations/{animationId}/play
Content-Type: application/json

{
  "loop": true,
  "speed": 1.0,
  "blend_time": 0.2
}

Response: 200 OK
{
  "animation_id": "anim-abc123",
  "entity_id": "entity-abc123",
  "playing": true,
  "started_at": "2025-07-01T17:00:00Z"
}
```

#### **Stop Animation**
```yaml
POST /api/sessions/{sessionId}/entities/{entityId}/animations/{animationId}/stop

Response: 200 OK
{
  "animation_id": "anim-abc123",
  "entity_id": "entity-abc123",
  "playing": false,
  "stopped_at": "2025-07-01T17:05:00Z"
}
```

---

## 🔄 **WEBSOCKET REAL-TIME EVENTS**

### **Entity Events**
```yaml
# Entity Created
{
  "type": "entity_created",
  "data": {
    "session_id": "session-123",
    "entity_id": "entity-abc123",
    "name": "MyEntity"
  },
  "timestamp": 1719849600
}

# Entity Updated
{
  "type": "entity_updated", 
  "data": {
    "session_id": "session-123",
    "entity_id": "entity-abc123",
    "changes": ["position", "rotation"],
    "transform": {
      "position": {"x": 10, "y": 5, "z": -2},
      "rotation": {"x": 0, "y": 90, "z": 0, "w": 1}
    }
  },
  "timestamp": 1719849660
}

# Component Added
{
  "type": "component_added",
  "data": {
    "session_id": "session-123", 
    "entity_id": "entity-abc123",
    "component_type": "model",
    "properties": {
      "asset": "model-asset-id"
    }
  },
  "timestamp": 1719849720
}
```

---

## 🎯 **IMPLEMENTATION STRATEGY**

### **Phase 1 Priority (This Session)**
1. ✅ **Entity CRUD APIs** - Core entity lifecycle management
2. 🔄 **Transform APIs** - Position, rotation, scale operations  
3. ⏳ **Component Add/Remove** - Basic component system

### **Phase 2 Priority (Next Session)**
1. **Scene Graph APIs** - Hierarchy and parenting
2. **Asset Management** - File upload and management
3. **Animation Control** - Basic play/stop functionality

### **Success Metrics**
- **API Coverage**: All core PlayCanvas features accessible via HTTP
- **Real-time Sync**: WebSocket events for all state changes
- **Session Isolation**: Multi-tenant game instances
- **Zero Manual Config**: All integration auto-generated from analysis

---

**🎮 HD1 v3.0: Transforming game development from complex SDKs to simple REST APIs**

*Last Updated: 2025-07-01*  
*Status: Phase 1 Entity Management Design Complete*