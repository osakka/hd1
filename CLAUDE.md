# HD1 (Holodeck One) - Development Context

## Core Principles
- **API-First**: HD1 exposes everything via unified API surface (82 endpoints)
- **Specification-Driven**: 100% auto-generated from api.yaml
- **Real-Time Sync**: WebSocket for state sync, API for commands
- **Single Source of Truth**: No parallel implementations
- **Quality Only**: Zero regressions, clean architecture
- **"API = Control, WebSocket = Graph Extension"**

## Current State (2025-07-03)
HD1 v5.0.1 is a **production-ready API-first game engine** platform with **complete configuration management standardization**.

### ✅ Completed Features
- **82 REST Endpoints**: Complete game engine control via HTTP + Avatar management
- **Real-Time WebSocket**: Entity lifecycle synchronization + avatar position updates
- **PlayCanvas Integration**: Professional 3D rendering with ECS + advanced camera system
- **Avatar Synchronization**: High-frequency multiplayer avatar tracking with persistence
- **Advanced Camera System**: Smooth movement, momentum, orbital mode with TAB toggle
- **Channel System**: YAML-based scene configuration with 3 channels
- **Console UI**: Professional monitoring with smooth animations
- **Vendor Cleanup**: Removed 1.1GB redundant directories, optimized structure
- **Template Architecture**: 8 externalized templates for maintainable code generation
- **Configuration Management**: Complete configuration system with priority order: Flags > Environment Variables > .env File > Defaults

### Architecture
```
HTTP APIs → Game Commands → Server State → WebSocket Events → Client Rendering
```

**Key Files:**
- `src/api.yaml` - Single source of truth (82 total route endpoints, 59 unique paths)
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

## Configuration Management System

**Priority Order**: Flags > Environment Variables > .env File > Defaults

### Environment Variables (HD1_* prefix)
```bash
# Server Configuration
HD1_HOST=0.0.0.0                          # Server bind host
HD1_PORT=8080                             # Server bind port  
HD1_API_BASE=http://0.0.0.0:8080/api     # External API base URL
HD1_INTERNAL_API_BASE=http://localhost:8080/api  # Internal API communications
HD1_VERSION=v5.0.1                       # HD1 version identifier
HD1_DAEMON=true                           # Run in daemon mode

# Directory Paths
HD1_ROOT_DIR=/opt/hd1                     # HD1 root directory
HD1_STATIC_DIR=/opt/hd1/share/htdocs/static  # Static files directory
HD1_LOG_DIR=/opt/hd1/build/logs          # Log directory
HD1_CHANNELS_DIR=/opt/hd1/share/channels # Channels configuration directory
HD1_AVATARS_DIR=/opt/hd1/share/avatars   # Avatars configuration directory

# Logging Configuration  
HD1_LOG_LEVEL=INFO                        # Logging level (TRACE, DEBUG, INFO, WARN, ERROR, FATAL)
HD1_TRACE_MODULES=websocket,entities     # Comma-separated trace modules

# Channel System
HD1_CHANNELS_DEFAULT_CHANNEL=channel_one # Default channel
HD1_CHANNELS_PROTECTED_LIST=channel_one,channel_two # Protected channels list
```

### Command-Line Flags (All with long options)
```bash
./hd1 --host=127.0.0.1 --port=8081 --daemon --log-level=DEBUG
./hd1 --root-dir=/custom/path --internal-api-base=http://api.internal:8080/api
./hd1 --protected-channels=secure_one,secure_two --version=v1.0.0
```

### .env File Support
Create `.env` file in project root:
```bash
HD1_HOST=127.0.0.1
HD1_PORT=9090
HD1_LOG_LEVEL=DEBUG
HD1_DAEMON=false
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
- **Performance Optimized**: Real-time WebSocket for multiplayer synchronization

**Camera System Architecture**:
- **Free Camera Mode**: WASD movement with smooth momentum and acceleration/deceleration
- **Orbital Camera Mode**: TAB to toggle, automatic avatar centering, mouse wheel zoom
- **PlayCanvas Integration**: Professional Vec3.lerp() interpolation and Quat.slerp() rotation
- **Error-Free Operation**: Fixed critical PlayCanvas constructor issues

## Configuration Management
**Priority**: Flags > Environment Variables > .env File > Defaults
**Environment Variables**: 50+ HD1_* variables for complete configurability
**Key Variables**:
- `HD1_HOST`, `HD1_PORT` - Server binding configuration
- `HD1_CHANNELS_DIR`, `HD1_AVATARS_DIR`, `HD1_RECORDINGS_DIR` - Path configuration
- `HD1_WEBSOCKET_*` - WebSocket timeout and buffer configuration
- `HD1_LOG_LEVEL`, `HD1_TRACE_MODULES` - Logging configuration

## Logging Standards
**Format**: `timestamp [pid:thread] [level] function.file:line message`  
**Levels**: TRACE (dev), DEBUG (dev), INFO (prod), WARN (prod), ERROR (prod), FATAL (all)  
**Control**: Real-time via `/api/admin/logging/level` and `/api/admin/logging/trace`  

### Log Level Usage
- **TRACE**: Module-specific debugging (websocket, entities, sessions)
- **DEBUG**: Detailed operations for development troubleshooting
- **INFO**: Production events, successful operations, startup/shutdown
- **WARN**: Non-critical issues, missing resources, degraded performance
- **ERROR**: Operation failures, request errors, system problems
- **FATAL**: Critical failures causing process termination

### Runtime Control
```bash
# API-based level changes
curl -X POST /api/admin/logging/level -d '{"level":"DEBUG"}'
curl -X POST /api/admin/logging/trace -d '{"enable":["websocket","entities"]}'

# Command-line flags
./hd1 --log-level=DEBUG --trace-modules=websocket,entities

# Environment variables
HD1_LOG_LEVEL=TRACE HD1_TRACE_MODULES=websocket,entities ./hd1
```

### Thread Safety & Performance
- **Zero-overhead**: Disabled log levels consume ~50ns CPU
- **Thread-safe**: RWMutex with single lock acquisition
- **Structured**: JSON data with consistent field names
- **Context-rich**: operation, session_id, entity_id, error details

## Quality Standards
- **Auto-Generated**: Never edit auto_router.go, hd1lib.js, hd1lib.sh
- **Source Files**: Always edit api.yaml, handler implementations
- **Zero Regressions**: All changes maintain compatibility
- **Clean Architecture**: Separation of concerns maintained
- **Single Source of Truth**: All configuration via config system, zero hardcoded values

---

**HD1 v5.0.1**: Where OpenAPI specifications become immersive multiplayer game worlds with complete configuration management standardization.