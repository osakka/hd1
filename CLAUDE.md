# HD1 (Holodeck One) - Development Context

## Core Principles
- **API-First**: HD1 exposes everything via unified API surface (77 endpoints)
- **Specification-Driven**: 100% auto-generated from api.yaml
- **Real-Time Sync**: WebSocket for state sync, API for commands
- **Single Source of Truth**: No parallel implementations
- **Quality Only**: Zero regressions, clean architecture
- **"API = Control, WebSocket = Graph Extension"**

## Current State (2025-07-03)
HD1 v5.0.0 is a **production-ready API-first game engine** platform with **advanced multiplayer avatar synchronization**.

### ✅ Completed Features
- **79 REST Endpoints**: Complete game engine control via HTTP + Avatar management
- **Real-Time WebSocket**: <10ms entity lifecycle synchronization + avatar position updates
- **PlayCanvas Integration**: Professional 3D rendering with ECS + advanced camera system
- **Avatar Synchronization**: High-frequency multiplayer avatar tracking with persistence
- **Advanced Camera System**: Smooth movement, momentum, orbital mode with TAB toggle
- **Channel System**: YAML-based scene configuration with 3 channels
- **Console UI**: Professional monitoring with smooth animations
- **Vendor Cleanup**: Removed 1.1GB redundant directories, optimized structure
- **Template Architecture**: 8 externalized templates for maintainable code generation

### Architecture
```
HTTP APIs → Game Commands → Server State → WebSocket Events → Client Rendering
```

**Key Files:**
- `src/api.yaml` - Single source of truth (79 endpoints)
- `src/auto_router.go` - Auto-generated routing
- `src/codegen/templates/` - External template files
- `share/channels/*.yaml` - Scene configurations
- `share/htdocs/static/js/hd1-console/` - Modular console system
- `share/htdocs/static/js/hd1-playcanvas.js` - Advanced camera system with orbital mode
- `src/api/sessions/avatar/` - Avatar management endpoints
- `src/api/camera/position.go` - Camera position API with avatar sync

## Development Commands
```bash
cd src && make clean && make && make start  # Build & start
make stop                                   # Stop server
make generate                              # Auto-generate from api.yaml
```

## Console Features
- **Smooth Animations**: 0.4s cubic-bezier expand/collapse
- **Performance Monitoring**: CPU, memory, WebSocket stats
- **Channel Switching**: Live channel selection
- **Entity Tracking**: Real-time entity counts + avatar count display
- **Advanced Camera**: TAB toggle between free/orbital modes, smooth movement with momentum
- **Avatar Sync**: Real-time multiplayer avatar position updates via WebSocket

## Avatar Synchronization System
**High-Frequency Multiplayer Support**: Handles 100+ movements/updates per second
- **Avatar Persistence**: Prevents avatar disappearing during rapid position updates
- **Dual Message Types**: `avatar_position_update` for movement, `entity_updated` for creation
- **Entity Lifecycle Protection**: Direct position updates avoid delete/recreate cycles
- **Channel Broadcasting**: Bidirectional avatar visibility across sessions
- **Performance Optimized**: <10ms WebSocket latency for real-time multiplayer

**Camera System Architecture**:
- **Free Camera Mode**: WASD movement with smooth momentum and acceleration/deceleration
- **Orbital Camera Mode**: TAB to toggle, automatic avatar centering, mouse wheel zoom
- **PlayCanvas Integration**: Professional Vec3.lerp() interpolation and Quat.slerp() rotation
- **Error-Free Operation**: Fixed critical PlayCanvas constructor issues

## Quality Standards
- **Auto-Generated**: Never edit auto_router.go, hd1lib.js, hd1lib.sh
- **Source Files**: Always edit api.yaml, handler implementations
- **Zero Regressions**: All changes maintain compatibility
- **Clean Architecture**: Separation of concerns maintained

## Logging Standards
**Format**: `timestamp [pid:thread] [level] function.file:line message`
**Levels**: TRACE (dev), DEBUG (dev), INFO (prod), WARN (prod), ERROR (prod), FATAL (all)
**Control**: Real-time via `/api/admin/logging/level` and `/api/admin/logging/trace`

---

**HD1 v5.0.0**: Where OpenAPI specifications become immersive multiplayer game worlds with real-time avatar synchronization.