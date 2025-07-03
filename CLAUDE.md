# HD1 (Holodeck One) - Development Context

## Core Principles
- **API-First**: HD1 exposes everything via unified API surface (85 endpoints)
- **Specification-Driven**: 100% auto-generated from api.yaml
- **Real-Time Sync**: WebSocket for state sync, API for commands
- **Single Source of Truth**: No parallel implementations
- **Quality Only**: Zero regressions, clean architecture
- **"API = Control, WebSocket = Graph Extension"**

## Current State (2025-07-03)
HD1 v3.0 is a **production-ready API-first game engine** platform.

### ✅ Completed Features
- **85 REST Endpoints**: Complete game engine control via HTTP
- **Real-Time WebSocket**: <10ms entity lifecycle synchronization
- **PlayCanvas Integration**: Professional 3D rendering with ECS
- **Channel System**: YAML-based scene configuration with 3 channels
- **Console UI**: Professional monitoring with smooth animations
- **Floor Fix**: All channels now have 10×1×10 flat floors
- **Template Architecture**: 8 externalized templates for maintainable code generation

### Architecture
```
HTTP APIs → Game Commands → Server State → WebSocket Events → Client Rendering
```

**Key Files:**
- `src/api.yaml` - Single source of truth (85 endpoints)
- `src/auto_router.go` - Auto-generated routing
- `src/codegen/templates/` - External template files
- `share/channels/*.yaml` - Scene configurations
- `share/htdocs/static/js/hd1-console/` - Modular console system

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
- **Entity Tracking**: Real-time entity counts

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

**HD1 v3.0**: Where OpenAPI specifications become immersive game worlds.