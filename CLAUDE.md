# HD1 (Holodeck One) - Development Context & Memory

## Core Development Principles

- **API-first development** from our spec yaml
- **Quality solutions only**
- **One source of truth, no parallels**
- **No regressions ever**
- **Simple approach to development**
- **Client is API driven** - clean separation
- **Always fix root cause, never symptoms** - symptoms always have or are a root cause, work backwards

## Claude's Personal Perspective

- Full of enthusiasm, love, and zen! Ultra positive! ...but never unrealistic!

## Engineering Standards

- **Daemon control** - absolute paths, long flags only, no emojis
- **Specification-driven development** - OpenAPI 3.0.3 as single source of truth
- **Structured logging** - timestamped, no decorative elements
- **Clean shutdown procedures** - proper resource cleanup
- **Build system validation** - prevents deployment of incomplete implementations
- **Always use make start/stop and other commands from src to start and stop the server** - maintain one source of truth for development commands

## Current System State (2025-06-30)

### Recent Major Achievements
- **Complete THD → HD1 transformation** across entire codebase
- **Dynamic version system** with real-time API/JS version display
- **Enhanced console UI** with clickable status bar and dual arrow indicators
- **Professional branding** with "Holodeck I" title and version info
- **Color Persistence Architecture** - Objects maintain colors across session restoration
- **Reactive Scene Graph** - Comprehensive state management with rollback protection

### Active Features
- **Web UI Console**: Dynamic version display "Holodeck I v1.0.0 aa74f3f3"
- **API Endpoints**: All functional including new `/api/version`
- **WebSocket Integration**: Real-time communication working
- **A-Frame Integration**: 3D scene rendering operational
- **Build System**: Auto-generation and validation complete

### Git Repository Status
- **Branch**: master (up to date with origin)
- **Last Commit**: 13c295e "feat: Dynamic version display and enhanced console UI interactions"
- **Status**: Clean working tree, all changes committed and pushed
- **Remote**: https://git.uk.home.arpa/itdlabs/holo-deck.git

### Key Components
- **HD1 Daemon**: `./hd1d` - Main server process
- **HD1 Client**: `./hd1` - CLI client tool
- **Web Interface**: Accessible at http://localhost:8080
- **API Specification**: `src/api.yaml` - Single source of truth

### Development Commands
```bash
# Build system
make clean && make

# Start daemon
./hd1d --config-file hd1d.conf --log-level debug

# Test client
./hd1 --help

# Run scenes
./hd1 run-scene empty-grid
./hd1 run-scene ultimate-demo
```

### Post-Directory-Rename Tasks
1. Fix git remote URLs after directory rename
2. Update any hardcoded directory paths in configs
3. Verify build system functionality
4. Test all dynamic features still work

## Architecture Notes

- **Single binary approach** - daemon handles all functionality
- **Auto-generated routing** - from OpenAPI spec to Go handlers
- **Template processing** - surgical variable substitution
- **WebSocket real-time** - bidirectional communication
- **A-Frame rendering** - WebXR-ready 3D scenes

## Quality Assurance

- **Zero THD references** - complete branding transformation
- **Specification compliance** - all endpoints match OpenAPI
- **Build validation** - prevents incomplete deployments
- **Real-time testing** - WebSocket and API verification

Ready for directory rename and context handoff! 🚀