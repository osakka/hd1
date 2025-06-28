# THD (The Holo-Deck) - Development Context & Memory

## Core Development Principles

- **100% API-first development** from our spec yaml
- **Bar raising solutions only**
- **One source of truth, no parallels**
- **No regressions ever**
- **Zen approach to development**
- **Client is 100% API driven** - clean separation
- **Always fix root cause, never symptoms** - symptoms always have or are a root cause, work backwards

## Professional Engineering Standards

- **Professional daemon control** - 100% absolute paths, long flags only, no emojis
- **Specification-driven development** - OpenAPI 3.0.3 as single source of truth
- **Professional logging** - timestamped, no decorative elements
- **Clean shutdown procedures** - proper resource cleanup
- **Build system validation** - prevents deployment of incomplete implementations

## Project Identity: THD (The Holo-Deck)

### Current State (v2.0.0)
- **Binary**: `thd` (professional daemon)
- **Module**: `holodeck` (Go module name)
- **Git Repository**: `https://git.uk.home.arpa/itdlabs/holo-deck.git`
- **Credentials**: `claude-3/claude-password`
- **Directory Structure**: Professional organization with clean .gitignore

### Technical Architecture
- **3D Coordinate System**: 25×25×25 grid with [-12, +12] boundaries
- **Real-time Communication**: WebSocket Hub with thread-safe SessionStore
- **Auto-generated Routing**: Complete routing layer generated from api.yaml
- **Professional Build System**: Make-based with daemon control targets

## Key File Locations

### Source Code
- `/home/claude-3/3dv/src/api.yaml` - Single source of truth specification
- `/home/claude-3/3dv/src/main.go` - THD daemon with professional standards
- `/home/claude-3/3dv/src/auto_router.go` - Auto-generated routing (DO NOT EDIT)
- `/home/claude-3/3dv/src/Makefile` - Professional build system

### Runtime
- `/home/claude-3/3dv/build/bin/thd` - Professional daemon binary
- `/home/claude-3/3dv/build/runtime/thd.pid` - Process ID file
- `/home/claude-3/3dv/build/logs/` - Professional logging directory

### Documentation
- `/home/claude-3/3dv/README.md` - Professional project overview
- `/home/claude-3/3dv/CHANGELOG.md` - Complete project transformation history
- `/home/claude-3/3dv/docs/api/README.md` - THD API documentation

## Development Workflow

### Professional Daemon Control
```bash
make start     # Start THD daemon
make stop      # Stop THD daemon with proper cleanup
make restart   # Restart with validation
make status    # Professional status reporting
```

### Build System
```bash
make all       # Complete build pipeline
make build     # Build THD binary
make validate  # Validate API specification
make generate  # Generate routing from spec
```

### Git Operations
```bash
git remote origin: https://git.uk.home.arpa/itdlabs/holo-deck.git
git config user.name "claude-3"
git config user.email "claude@anthropic.com"
```

## Critical Development Context

### Session Management Issues (Historical - RESOLVED)
- **Fixed**: Empty object lists due to hardcoded returns instead of store calls
- **Fixed**: Wireframe property missing from canvas control objects
- **Fixed**: PID file writing issues in daemon mode
- **Fixed**: Emoji removal for professional standards

### Professional Standards Implementation
- **No emojis** in any system output or logging
- **Absolute paths only** - all configured from THD_* constants
- **Long flags only** - no short flags to eliminate confusion
- **Professional error messages** - clear, actionable, no decorative elements

### Breaking Changes from VWS to THD
- **Module imports**: `visualstream/*` → `holodeck/*`
- **Binary names**: `vws` → `thd`, `vws-client` → `thd-client`
- **Constants**: `VWS_*` → `THD_*`
- **PID files**: `vws.pid` → `thd.pid`

## Debugging Context

### WebGL Rendering
- **Coordinate system**: Fixed boundaries [-12, +12] on all axes
- **Object creation**: Must validate coordinates at API level
- **Wireframe support**: Ensure `wireframe: obj.wireframe || false` in canvas control
- **Session isolation**: Each session maintains independent object store

### API Control
- **Browser control**: Force refresh and canvas manipulation APIs available
- **Session bootstrap**: Complete world initialization with grid system
- **Real-time updates**: WebSocket hub broadcasts all changes

## Recovery Context for New Sessions

When resuming development after session restart:

1. **Check daemon status**: `cd /home/claude-3/3dv/src && make status`
2. **Verify build**: `make all` to ensure clean build
3. **Git status**: Confirm all changes committed to `https://git.uk.home.arpa/itdlabs/holo-deck.git`
4. **API validation**: `make test` to verify all endpoints functional
5. **Professional standards**: Confirm no emojis, absolute paths maintained

## Project Philosophy

> **"Where 3D visualization meets professional engineering"**

THD represents the evolution from innovative VWS concept to professional-grade engineering solution, maintaining all revolutionary capabilities while implementing proper software engineering practices.

**Never compromise on professional standards while preserving innovation.**