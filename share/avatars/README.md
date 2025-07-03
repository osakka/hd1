# HD1 Avatar Asset System

## Overview
Specification-driven avatar management for HD1 collaborative sessions. Each avatar is defined by a YAML specification that drives entity creation, 3D asset loading, and real-time synchronization.

## Architecture Principles
- **API-First**: All avatar operations via unified API surface
- **Specification-Driven**: 100% YAML-defined avatar configurations
- **Single Source of Truth**: Avatar specs drive all behavior
- **Real-Time Sync**: WebSocket synchronization with <10ms latency

## Directory Structure
```
avatars/
├── config.yaml              # System configuration
├── README.md                # This file
├── default/                 # Standard avatar
│   ├── avatar.yaml          # Avatar specification
│   ├── model.glb           # 3D model
│   ├── textures/           # Character textures
│   └── animations/         # Animation assets
├── business/               # Professional avatar
│   ├── avatar.yaml
│   ├── business_model.glb
│   ├── textures/
│   └── animations/
└── casual/                 # Relaxed avatar
    ├── avatar.yaml
    ├── casual_model.glb
    ├── textures/
    └── animations/
```

## Avatar Specification Schema

Each `avatar.yaml` defines:
- **Metadata**: Name, description, version
- **Assets**: 3D models, textures, animations
- **Physical**: Transforms, camera positioning, physics
- **Entity**: HD1 entity definition with PlayCanvas components
- **Movement**: Speed, acceleration, animation states
- **Network**: Sync frequency, interpolation settings

## Usage

### API Integration (Planned)
```bash
# List available avatars
GET /avatars

# Get avatar specification
GET /avatars/default

# Set session avatar
POST /sessions/{id}/avatar {"type": "business"}

# Get current session avatar
GET /sessions/{id}/avatar
```

### Channel Integration
Avatars auto-select based on channel context:
- `channel_one`: Default avatar (general purpose)
- `channel_two`: Business avatar (meetings)
- `channel_three`: Casual avatar (creative work)

### Session Lifecycle
1. **Join Session** → Auto-create avatar entity
2. **WASD Movement** → Update camera position via API
3. **Position Change** → Broadcast via WebSocket
4. **Other Clients** → Receive update, render movement
5. **Leave Session** → Clean up avatar entity

## Data Flow

```
Avatar Selection → YAML Spec → Entity Creation → PlayCanvas Rendering
       ↓              ↓              ↓               ↓
   API Control → State Update → WebSocket Sync → Real-time Display
```

## Integration Points

### HD1 Core Systems
- **API Layer**: Avatar management endpoints
- **Entity System**: Avatar as tagged entities
- **Channel System**: Context-aware avatar selection
- **Session Management**: Avatar lifecycle
- **WebSocket Sync**: Real-time position updates

### PlayCanvas Integration
- **Model Loading**: GLB asset pipeline
- **Animation System**: State-driven animations
- **Physics**: Collision detection and movement
- **Rendering**: Real-time 3D display with shadows/lighting

## Performance Optimizations
- **Asset Compression**: Draco models, Basis textures
- **LOD System**: Distance-based quality adjustment
- **Network Culling**: Limit visible avatars and sync frequency
- **Caching**: 1-hour asset cache with preloading

## Configuration
See `config.yaml` for system-wide settings including:
- Default avatar types per channel
- Network sync optimization
- Asset loading performance
- Integration parameters

---

**HD1 v5.0**: Where avatar specifications become immersive collaborative experiences.