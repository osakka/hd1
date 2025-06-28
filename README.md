# THD (The Holo-Deck) - Professional 3D Visualization Engine

Revolutionary API-first 3D coordinate system with universal world boundaries.

## Core Principles

- **API-Driven Architecture**: Everything controlled via REST API, zero shell commands
- **Universal Coordinate System**: Fixed [-12, +12] bounds on all axes  
- **Specification-Driven Development**: OpenAPI 3.0.3 as single source of truth
- **Professional Daemon Control**: Robust process management with PID files
- **Real-time Communication**: WebSocket + REST hybrid architecture

## Quick Start

```bash
# Build the system
cd src && make all

# Start daemon
make start

# Check status  
make status

# Stop daemon
make stop
```

## Architecture

- **Go Backend**: Professional daemon with absolute path configuration
- **WebGL Frontend**: Real-time 3D rendering with Three.js
- **Auto-Generated Routing**: API specification drives code generation
- **Session Management**: Isolated 3D worlds with persistence
- **Object Lifecycle**: Named objects with full CRUD operations

## API Endpoints

- `POST /api/sessions` - Create new 3D session
- `GET /api/sessions` - List all active sessions  
- `POST /api/sessions/{id}/world` - Initialize world coordinate system
- `POST /api/sessions/{id}/objects` - Create named 3D objects
- `PUT /api/sessions/{id}/camera/position` - Control camera position

## Development

Built with professional standards:
- Long flags only (no short flags)
- 100% absolute paths in configuration
- Professional logging with timestamps
- Comprehensive error handling
- No emojis in system output

## World Constraints

- **Grid Size**: 25×25×25 coordinate system
- **Boundaries**: [-12, +12] on X, Y, Z axes
- **Max Objects**: 1000 per session
- **Max Sessions**: 100 concurrent

## Building

```bash
cd src
make validate  # Validate API specification
make generate  # Generate routing code
make build     # Build THD binary
make test      # Test API endpoints
```

## Daemon Control

```bash
make start     # Start THD daemon
make stop      # Stop THD daemon  
make restart   # Restart THD daemon
make status    # Show daemon status
```

## File Structure

```
/home/claude-3/3dv/
├── src/           # Go source code
├── share/         # Web assets (htdocs, static files)
├── build/         # Build artifacts (excluded from git)
├── api.yaml       # OpenAPI specification (single source of truth)
└── README.md      # This file
```

## Professional Standards

- Absolute path configuration
- Professional error handling
- Comprehensive logging
- PID file management
- Clean shutdown procedures
- No runtime artifacts in git

---

**THD (The Holo-Deck)** - Where 3D visualization meets professional engineering.