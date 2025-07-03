# HD1 (Holodeck One) - API-First Game Engine Platform

**The world's first Game Engine as a Service - Professional 3D game development via REST APIs**

## Overview

HD1 v5.0.0 is a production-ready **API-first game engine platform** that exposes complete game engine functionality through REST endpoints. Built with PlayCanvas professional 3D rendering and real-time WebSocket synchronization.

## Key Features

### 🎮 API-First Game Engine Architecture
- **77 REST Endpoints**: Complete game engine control via HTTP APIs
- **Professional 3D Rendering**: PlayCanvas integration with WebGL/WebXR support
- **Entity-Component-System**: Full ECS architecture with lifecycle management
- **Real-Time Synchronization**: <10ms WebSocket state sync across all clients
- **Single Source of Truth**: All functionality auto-generated from api.yaml specification

### 🏗️ Complete Game Engine APIs
- **Entity Management**: Create, update, delete entities with full component systems
- **Physics Engine**: Rigidbodies, force application, collision detection
- **Animation System**: Timeline-based animations with play/stop controls
- **Audio Engine**: 3D positional audio sources with spatial audio
- **Scene Graph**: Hierarchical transforms, parent-child relationships
- **Camera Controls**: Position, orbit, and movement APIs

### 🌐 Channel-Based Architecture
- **YAML Configuration**: Declarative scene definition via channel files
- **Multi-Channel Support**: Isolated collaborative environments
- **Real-Time Collaboration**: Multiple users per channel with live synchronization
- **Session Management**: Per-user session isolation with state restoration

### 🛠️ Development Features
- **Auto-Generated Clients**: JavaScript, Go CLI, and shell functions from specification
- **Template Architecture**: 8 externalized templates for maintainable code generation
- **Clean Build System**: Make-based with automatic code generation
- **Professional Console**: Real-time monitoring with smooth animations
- **Performance Optimized**: <50ms API response, <10ms WebSocket latency

## Quick Start

```bash
# Build the system
cd src && make clean && make

# Start the daemon
make start

# Access the console
open http://localhost:8080
```

## API Usage Examples

### Entity Management
```bash
# Create a session
SESSION_ID=$(./build/bin/hd1-client create-session | jq -r '.session_id')

# Create an entity
./build/bin/hd1-client create-entity "$SESSION_ID" '{
  "name": "my-cube",
  "components": {
    "transform": {"position": {"x": 0, "y": 1, "z": 0}},
    "render": {"geometry": "box", "material": {"color": "#ff0000"}}
  }
}'

# List all entities
./build/bin/hd1-client list-entities "$SESSION_ID"
```

### Channel-Based Scene Management
```bash
# List available channels
./build/bin/hd1-client list-channels

# Join a session to a channel
./build/bin/hd1-client join-session-channel "$SESSION_ID" '{
  "channel_id": "channel_one",
  "client_id": "player1"
}'

# Get channel status
./build/bin/hd1-client get-session-channel-status "$SESSION_ID"
```

### Physics and Animation
```bash
# Apply force to an entity
./build/bin/hd1-client apply-force "$SESSION_ID" "entity-id" '{
  "force": {"x": 100, "y": 0, "z": 0},
  "point": {"x": 0, "y": 0, "z": 0}
}'

# Create and play animation
./build/bin/hd1-client create-animation "$SESSION_ID" '{
  "name": "rotate-cube",
  "target": "entity-id",
  "duration": 2.0,
  "properties": {"rotation": {"y": 360}}
}'
```

## Architecture

### System Flow
```
HTTP APIs → Game Commands → Server State → WebSocket Events → PlayCanvas Rendering
```

### Key Components
- **API Specification**: `src/api.yaml` - Single source of truth (77 endpoints)
- **Auto-Generated Router**: `src/auto_router.go` - HTTP routing from specification
- **PlayCanvas Engine**: Professional 3D rendering with ECS architecture
- **WebSocket Hub**: Real-time state synchronization across clients
- **Channel System**: YAML-based scene configuration in `share/channels/`

### Template Architecture
```
src/codegen/templates/
├── go/router.tmpl              # Auto-router generation
├── javascript/api-client.tmpl  # JS API wrapper
├── javascript/playcanvas-bridge.tmpl # PlayCanvas integration
└── shell/playcanvas-functions.tmpl   # Shell function library
```

## API Endpoints (77 Total)

### Core APIs
- **Sessions**: Create, list, get, delete sessions
- **Entities**: Full CRUD with component management
- **Channels**: Multi-user collaboration spaces
- **Physics**: World simulation and rigidbody control
- **Animation**: Timeline-based animation system
- **Audio**: 3D positional audio management

### Advanced APIs
- **Scene Graph**: Hierarchical entity relationships
- **Lifecycle**: Entity activation, deactivation, destruction
- **Recording**: Session capture and playback
- **Hierarchy**: Parent-child entity relationships
- **Components**: Dynamic component attachment/detachment

## Performance Metrics

- **API Response**: <50ms average
- **WebSocket Latency**: <10ms real-time sync
- **Entity Creation**: ~5ms per entity
- **Scene Loading**: <200ms for complex scenes
- **Memory Usage**: ~100MB baseline, scales with entities

## Development

### Build Requirements
- Go 1.21+
- Make
- jq (for JSON processing)

### Development Commands
```bash
# Clean build
cd src && make clean && make

# Start development server
make start

# Generate code from specification
make generate

# Run tests
make test
```

### Architecture Principles
- **Specification-Driven**: All code generated from api.yaml
- **Single Source of Truth**: No manual synchronization needed
- **Zero Regressions**: Surgical precision in all changes
- **Production Ready**: Clean builds, comprehensive testing

## License

MIT License - See LICENSE file for details.

## Support

- **Documentation**: `/docs/` directory with comprehensive guides
- **API Reference**: Auto-generated from OpenAPI specification
- **Issues**: Report via git repository
- **Architecture Decisions**: See `/docs/adr/` for detailed design rationale