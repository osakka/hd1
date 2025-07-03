# HD1 API Reference

**HD1 v5.0.0 - Complete API-First Game Engine Platform**

HD1 provides a comprehensive REST API for 3D game engine control with **82 endpoints** covering all aspects of entity management, physics simulation, audio systems, and real-time multiplayer capabilities.

## API Overview

- **Total Endpoints**: 82 REST endpoints
- **API Specification**: OpenAPI 3.0.3
- **Base URL**: `http://localhost:8080/api`
- **Response Format**: JSON
- **Real-Time Updates**: WebSocket at `ws://localhost:8080/ws`

## Endpoint Categories

### Sessions (67 endpoints)
Core game engine functionality including entity management, physics, audio, and multiplayer features.

#### Entity Management
- `GET /sessions/{sessionId}/entities` - List all entities in session
- `POST /sessions/{sessionId}/entities` - Create new entity
- `GET /sessions/{sessionId}/entities/{entityId}` - Get entity details
- `PUT /sessions/{sessionId}/entities/{entityId}` - Update entity
- `DELETE /sessions/{sessionId}/entities/{entityId}` - Delete entity

#### Avatar System
- `GET /sessions/{sessionId}/avatar` - Get session avatar
- `POST /sessions/{sessionId}/avatar` - Set session avatar

#### Component System
- `GET /sessions/{sessionId}/entities/{entityId}/components` - List entity components
- `POST /sessions/{sessionId}/entities/{entityId}/components` - Add component to entity
- `GET /sessions/{sessionId}/entities/{entityId}/components/{componentType}` - Get specific component
- `PUT /sessions/{sessionId}/entities/{entityId}/components/{componentType}` - Update component
- `DELETE /sessions/{sessionId}/entities/{entityId}/components/{componentType}` - Remove component
- `POST /sessions/{sessionId}/entities/{entityId}/components/bulk` - Bulk component operations

#### Entity Lifecycle
- `POST /sessions/{sessionId}/entities/{entityId}/lifecycle/activate` - Activate entity
- `PUT /sessions/{sessionId}/entities/{entityId}/lifecycle/enable` - Enable entity
- `PUT /sessions/{sessionId}/entities/{entityId}/lifecycle/disable` - Disable entity
- `POST /sessions/{sessionId}/entities/{entityId}/lifecycle/deactivate` - Deactivate entity
- `DELETE /sessions/{sessionId}/entities/{entityId}/lifecycle/destroy` - Destroy entity
- `GET /sessions/{sessionId}/entities/{entityId}/lifecycle/status` - Get lifecycle status
- `POST /sessions/{sessionId}/entities/lifecycle/bulk` - Bulk lifecycle operations

#### Hierarchy & Transforms
- `GET /sessions/{sessionId}/entities/{entityId}/hierarchy/children` - Get entity children
- `GET /sessions/{sessionId}/entities/{entityId}/hierarchy/parent` - Get entity parent
- `PUT /sessions/{sessionId}/entities/{entityId}/hierarchy/parent` - Set entity parent
- `GET /sessions/{sessionId}/entities/{entityId}/hierarchy/transforms` - Get transforms
- `PUT /sessions/{sessionId}/entities/{entityId}/hierarchy/transforms` - Set transforms
- `GET /sessions/{sessionId}/entities/hierarchy/tree` - Get complete hierarchy tree

#### Physics System
- `GET /sessions/{sessionId}/physics/world` - Get physics world state
- `PUT /sessions/{sessionId}/physics/world` - Update physics world
- `GET /sessions/{sessionId}/physics/rigidbodies` - List rigidbodies
- `POST /sessions/{sessionId}/physics/rigidbodies/{entityId}/force` - Apply force to rigidbody

#### Animation System
- `GET /sessions/{sessionId}/animations` - List animations
- `POST /sessions/{sessionId}/animations` - Create animation
- `POST /sessions/{sessionId}/animations/{animationId}/play` - Play animation
- `POST /sessions/{sessionId}/animations/{animationId}/stop` - Stop animation

#### Audio System
- `GET /sessions/{sessionId}/audio/sources` - List audio sources
- `POST /sessions/{sessionId}/audio/sources` - Create audio source
- `POST /sessions/{sessionId}/audio/sources/{audioId}/play` - Play audio
- `POST /sessions/{sessionId}/audio/sources/{audioId}/stop` - Stop audio

#### Camera System
- `GET /sessions/{sessionId}/camera/position` - Get camera position
- `PUT /sessions/{sessionId}/camera/position` - Set camera position
- `POST /sessions/{sessionId}/camera/orbit` - Start camera orbit

#### Recording System
- `POST /sessions/{sessionId}/recording/start` - Start recording
- `POST /sessions/{sessionId}/recording/stop` - Stop recording
- `POST /sessions/{sessionId}/recording/play` - Play recording
- `GET /sessions/{sessionId}/recording/status` - Get recording status

#### Scene Management
- `GET /sessions/{sessionId}/scene/state` - Get scene state
- `PUT /sessions/{sessionId}/scene/state` - Update scene state
- `POST /sessions/{sessionId}/scene/state/save` - Save scene state
- `POST /sessions/{sessionId}/scene/state/load` - Load scene state
- `POST /sessions/{sessionId}/scene/state/reset` - Reset scene state
- `GET /sessions/{sessionId}/scene/export` - Export scene definition
- `POST /sessions/{sessionId}/scene/import` - Import scene definition
- `GET /sessions/{sessionId}/scene/hierarchy` - Get scene hierarchy
- `PUT /sessions/{sessionId}/scene/hierarchy` - Update scene hierarchy

#### Session Scenes
- `GET /sessions/{sessionId}/scenes` - List session scenes
- `POST /sessions/{sessionId}/scenes` - Create session scene
- `POST /sessions/{sessionId}/scenes/{sceneId}/activate` - Activate session scene

#### Channel Management
- `POST /sessions/{sessionId}/channel/join` - Join session to channel
- `POST /sessions/{sessionId}/channel/leave` - Leave session channel
- `GET /sessions/{sessionId}/channel/status` - Get channel status
- `GET /sessions/{sessionId}/channel/graph` - Get session graph
- `PUT /sessions/{sessionId}/channel/graph` - Update session graph
- `POST /sessions/{sessionId}/channel/sync` - Sync session state

#### Session Management
- `GET /sessions` - List all sessions
- `POST /sessions` - Create new session
- `GET /sessions/{sessionId}` - Get session details
- `DELETE /sessions/{sessionId}` - Delete session

### Channels (5 endpoints)
Multi-user collaboration and scene management.

- `GET /channels` - List all channels
- `POST /channels` - Create new channel
- `GET /channels/{channelId}` - Get channel details
- `PUT /channels/{channelId}` - Update channel
- `DELETE /channels/{channelId}` - Delete channel

### Administration (5 endpoints)
System logging and configuration management.

- `GET /admin/logging/config` - Get logging configuration
- `POST /admin/logging/config` - Set logging configuration
- `POST /admin/logging/level` - Set log level
- `POST /admin/logging/trace` - Set trace modules
- `GET /admin/logging/logs` - Get log entries

### Avatars (2 endpoints)
Avatar system and specifications.

- `GET /avatars` - List available avatars
- `GET /avatars/{avatarType}` - Get avatar specification

### Browser (2 endpoints)
Browser integration and control.

- `POST /browser/canvas` - Set canvas configuration
- `POST /browser/refresh` - Force browser refresh

### System (1 endpoint)
API metadata and version information.

- `GET /version` - Get API version information

## WebSocket Events

HD1 provides real-time updates via WebSocket connection at `ws://localhost:8080/ws`:

### Avatar Synchronization Events
- `avatar_position_update` - High-frequency avatar position updates
- `entity_updated` - Entity creation and modification events
- `session_joined` - User joined session
- `session_left` - User left session

### Entity Events
- `entity_created` - New entity created
- `entity_modified` - Entity properties changed
- `entity_deleted` - Entity removed
- `component_added` - Component attached to entity
- `component_removed` - Component detached from entity

### Scene Events
- `scene_loaded` - Scene state loaded
- `scene_saved` - Scene state saved
- `hierarchy_changed` - Scene hierarchy modified

## Authentication

Currently, HD1 operates without authentication for development purposes. Production deployments should implement appropriate authentication and authorization mechanisms.

## Rate Limiting

HD1 supports high-frequency operations including:
- Avatar position updates: 100+ per second
- Entity updates: Real-time with <10ms latency
- WebSocket messaging: Optimized for multiplayer scenarios

## Error Handling

All API endpoints return structured JSON error responses:

```json
{
  "error": "Error description",
  "code": "ERROR_CODE",
  "details": "Additional error details"
}
```

## API Specification

The complete OpenAPI 3.0.3 specification is available in `src/api.yaml` and serves as the single source of truth for all API functionality.

---

**HD1 v5.0.0** - Complete API-First Game Engine Platform  
**Generated from**: `src/api.yaml` (Single Source of Truth)  
**Last Updated**: 2025-07-03