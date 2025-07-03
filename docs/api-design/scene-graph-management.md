# Scene Graph Management API Design

## Overview
The Scene Graph Management system provides comprehensive control over PlayCanvas scene hierarchy, composition, and state management through REST API endpoints. This system enables developers to orchestrate complex 3D scenes through HTTP calls.

## Core Concepts

### Scene Graph Structure
- **Scene Root**: Top-level container for all scene entities
- **Scene Hierarchy**: Parent-child relationships between scene nodes
- **Scene Composition**: Multiple scenes within a single session
- **Scene State**: Serializable scene configuration and data

### API Endpoints Design

#### Scene Hierarchy Management
```
GET    /api/sessions/{sessionId}/scene/hierarchy          # Get complete scene hierarchy
PUT    /api/sessions/{sessionId}/scene/hierarchy          # Update scene hierarchy structure
POST   /api/sessions/{sessionId}/scene/hierarchy/nodes    # Add new hierarchy node
DELETE /api/sessions/{sessionId}/scene/hierarchy/nodes/{nodeId} # Remove hierarchy node
```

#### Scene State Management
```
GET    /api/sessions/{sessionId}/scene/state              # Get current scene state
PUT    /api/sessions/{sessionId}/scene/state              # Update scene state
POST   /api/sessions/{sessionId}/scene/state/save         # Save scene state snapshot
POST   /api/sessions/{sessionId}/scene/state/load         # Load scene state snapshot
POST   /api/sessions/{sessionId}/scene/state/reset        # Reset scene to initial state
```

#### Scene Composition
```
GET    /api/sessions/{sessionId}/scenes                   # List all scenes in session
POST   /api/sessions/{sessionId}/scenes                   # Create new scene in session
GET    /api/sessions/{sessionId}/scenes/{sceneId}         # Get specific scene details
PUT    /api/sessions/{sessionId}/scenes/{sceneId}         # Update scene properties
DELETE /api/sessions/{sessionId}/scenes/{sceneId}         # Remove scene from session
POST   /api/sessions/{sessionId}/scenes/{sceneId}/activate # Set scene as active
```

#### Scene Serialization
```
GET    /api/sessions/{sessionId}/scene/export             # Export scene definition
POST   /api/sessions/{sessionId}/scene/import             # Import scene definition
POST   /api/sessions/{sessionId}/scene/fork               # Fork scene to new session
POST   /api/sessions/{sessionId}/scene/merge              # Merge scenes together
```

## PlayCanvas Integration Points

### Scene Graph Mapping
- **PlayCanvas Scene**: Maps to HD1 session scene root
- **PlayCanvas Entity Hierarchy**: Exposed through hierarchy endpoints
- **PlayCanvas Scene Settings**: Accessible via state management
- **PlayCanvas Asset References**: Tracked in scene serialization

### Real-time Synchronization
- WebSocket events for hierarchy changes
- Scene state change notifications
- Multi-scene composition updates
- Asset loading progress tracking

## Request/Response Schemas

### Scene Hierarchy Response
```json
{
  "hierarchy": {
    "root": {
      "id": "scene-root",
      "name": "Scene Root",
      "children": [
        {
          "id": "entity-123",
          "name": "Main Camera",
          "type": "camera",
          "children": []
        }
      ]
    }
  },
  "metadata": {
    "totalNodes": 25,
    "maxDepth": 4,
    "lastModified": "2025-07-01T10:30:45Z"
  }
}
```

### Scene State Response
```json
{
  "state": {
    "lighting": {
      "ambientColor": "#404040",
      "skybox": "urban-environment"
    },
    "physics": {
      "gravity": {"x": 0, "y": -9.8, "z": 0},
      "enabled": true
    },
    "rendering": {
      "shadows": true,
      "fog": {"enabled": false}
    }
  },
  "entities": 45,
  "assets": ["model-1", "texture-2", "audio-3"]
}
```

## Success Criteria
- Complete scene hierarchy management via HTTP
- Scene state persistence and restoration
- Multi-scene composition within sessions
- Scene serialization for reuse and sharing
- Zero regression with existing entity/component APIs