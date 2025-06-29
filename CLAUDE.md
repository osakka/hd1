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

### Current State (v3.3.0 - Professional UI Excellence & Scene Management)
- **Binary**: `thd` (professional VR holodeck daemon)
- **Module**: `holodeck` (Go module name)
- **Rendering Engine**: A-Frame WebXR 1.4.0 (100% Local - 2.5MB ecosystem)
- **UI Standards**: Complete professional interface with scene management system
- **Cache Control**: Standards-compliant HTTP headers replacing hacky query strings
- **Scene System**: API-driven scene selection with cookie persistence
- **Sprint Controls**: Shift key acceleration for enhanced holodeck traversal
- **Holodeck Containment**: Escape-proof boundary enforcement with 60fps monitoring
- **Session Architecture**: Single source of truth with perfect isolation
- **A-Frame Components**: Complete local library with zero CDN dependencies
- **Git Repository**: `https://git.uk.home.arpa/itdlabs/holo-deck.git`
- **Credentials**: `claude-3/claude-password`
- **Directory Structure**: Professional organization with clean .gitignore

### Revolutionary Technical Architecture
- **VR/AR Holodeck**: Full WebXR support with A-Frame Entity-Component-System
- **Session Isolation**: Single source of truth with perfect multi-user session separation
- **Professional Error Handling**: All API responses return consistent JSON (zero parsing errors)
- **Enhanced Text Rendering**: Complete A-Frame text field transmission and display
- **WebSocket Session Association**: Automatic client-session binding for perfect isolation
- **3D Coordinate System**: 25×25×25 grid with [-12, +12] boundaries (holodeck-grade)
- **Holodeck Containment**: 100% escape-proof dual boundary enforcement system
- **60fps Position Monitoring**: Real-time boundary checking with visual feedback
- **100% Local A-Frame**: Complete 2.5MB ecosystem with zero CDN dependencies  
- **Multi-Backend Ready**: Framework-agnostic API layer supporting future engines
- **Real-time Communication**: WebSocket Hub with thread-safe SessionStore
- **Auto-generated Routing**: Complete routing layer generated from api.yaml
- **Professional Build System**: Make-based with daemon control targets
- **A-Frame Integration**: Seamless WebXR with 100% API compatibility

## Key File Locations

### Source Code
- `/opt/holo-deck/src/api.yaml` - Single source of truth specification
- `/opt/holo-deck/src/main.go` - THD daemon with professional standards
- `/opt/holo-deck/src/auto_router.go` - Auto-generated routing (DO NOT EDIT)
- `/opt/holo-deck/src/server/handlers.go` - A-Frame WebXR interface
- `/opt/holo-deck/share/htdocs/static/js/thd-aframe.js` - A-Frame holodeck manager
- `/opt/holo-deck/src/Makefile` - Professional build system

### Runtime
- `/opt/holo-deck/build/bin/thd` - Professional daemon binary
- `/opt/holo-deck/build/runtime/thd.pid` - Process ID file
- `/opt/holo-deck/build/logs/` - Professional logging directory

### Documentation
- `/opt/holo-deck/README.md` - Professional project overview
- `/opt/holo-deck/CHANGELOG.md` - Complete project transformation history
- `/opt/holo-deck/docs/adr/ADR-001-aframe-webxr-integration.md` - A-Frame integration decision
- `/opt/holo-deck/docs/adr/ADR-003-Professional-UI-Enhancement.md` - UI excellence & scene management
- `/opt/holo-deck/docs/api/README.md` - THD API documentation

### Holodeck Libraries
- `/opt/holo-deck/lib/thd-functions.sh` - Comprehensive shell function library
- `/opt/holo-deck/scenarios/ultimate-holodeck.thd` - Ultimate demonstration scenario

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

### Session Management Evolution (Historical Context)
- **Phase 1 (RESOLVED)**: Empty object lists due to hardcoded returns instead of store calls
- **Phase 2 (RESOLVED)**: Wireframe property missing from canvas control objects
- **Phase 3 (RESOLVED)**: PID file writing issues in daemon mode
- **Phase 4 (RESOLVED)**: Emoji removal for professional standards
- **Phase 5 (RESOLVED)**: Cross-session object visibility (session isolation implemented)
- **Phase 6 (RESOLVED)**: Text field transmission ("Holodeck Text" fallback eliminated)
- **Phase 7 (RESOLVED)**: jq parsing errors from non-JSON API responses

### Professional Standards Implementation
- **Zero Hacky Implementations** - proper HTTP cache control, professional status indicators
- **No emojis** in any system output or logging
- **Absolute paths only** - all configured from THD_* constants
- **Long flags only** - no short flags to eliminate confusion
- **Professional error messages** - clear, actionable, no decorative elements
- **Clean URLs** - no query string versioning, standards-compliant caching
- **Semantic UI** - professional status indicators replacing unicode characters

### Breaking Changes from VWS to THD
- **Module imports**: `visualstream/*` → `holodeck/*`
- **Binary names**: `vws` → `thd`, `vws-client` → `thd-client`
- **Constants**: `VWS_*` → `THD_*`
- **PID files**: `vws.pid` → `thd.pid`

## Professional UI Enhancement Context

### Scene Management System (v3.3.0)
- **API Endpoints**: `/api/scenes` (list) and `/api/scenes/{sceneId}` (load)
- **Predefined Scenes**: Empty Grid, Anime UI Demo, Ultimate Demo, Basic Shapes
- **Cookie Persistence**: 30-day scene preference storage with automatic restoration
- **Session Integration**: Scene dropdown appears under session ID in holodeck console
- **Auto-Bootstrap**: Saved scenes automatically load on session restore/creation

### Professional Interface Standards
- **Console Status**: `THD Console [ACTIVE]` / `THD Console [MINIMIZED]` (no unicode)
- **Status LED**: 6px professional indicator with hover tooltips (50% size reduction)
- **VR Button**: Removed empty rectangle (`vr-mode-ui="enabled: false"`)
- **Sprint Controls**: Shift key modifier for 3x speed boost (20 → 60 acceleration)
- **Scrollbar Theming**: Cross-browser holodeck aesthetic with cyan accents
- **Cache Control**: Standards-compliant HTTP headers for development and production

### UI Component Architecture
- **A-Frame Sprint Component**: `thd-sprint-controls` with dynamic acceleration switching
- **Scene Selection Flow**: Dropdown → Cookie Save → API Call → WebSocket Broadcast
- **Professional Scrollbars**: WebKit and Firefox theming for consistent holodeck look
- **HTTP Cache Strategy**: Development no-cache for JS/CSS, production caching for assets

## Debugging Context

### A-Frame VR Holodeck Rendering
- **Coordinate system**: Fixed boundaries [-12, +12] on all axes (holodeck-grade)
- **Entity-Component-System**: A-Frame ECS architecture with THDAFrameManager
- **WebXR Integration**: Full VR/AR headset support via A-Frame 1.4.0
- **Object creation**: Enhanced with color, materials, physics, lighting support
- **Session isolation**: Perfect single source of truth with WebSocket client association
- **Text rendering**: Complete text field transmission from API to A-Frame display
- **Error handling**: All API responses return consistent JSON (eliminates parsing errors)
- **Color rendering**: Proper hex conversion from RGBA to A-Frame materials
- **Advanced features**: PBR materials, particle effects, 3D text, physics simulation

### API Control
- **Browser control**: Force refresh and canvas manipulation APIs available
- **Session bootstrap**: Complete world initialization with grid system
- **Scene management**: API-driven scene loading with session isolation
- **Real-time updates**: WebSocket hub broadcasts all changes

### UI Enhancement Debugging
- **Cache Issues**: Check HTTP headers in Network tab (should show no-cache for dev)
- **Scene Loading**: Monitor console for `SCENE_LOADED` and `AUTO_SCENE` debug entries
- **Sprint Controls**: Verify `thd-sprint-controls` component attachment on camera entity
- **Status Indicators**: Professional `[ACTIVE]`/`[MINIMIZED]` instead of unicode arrows
- **Scrollbar Theme**: CSS custom properties for cross-browser holodeck styling

## Recovery Context for New Sessions

When resuming development after session restart:

1. **Check daemon status**: `cd /opt/holo-deck/src && make status`
2. **Verify build**: `make all` to ensure clean build
3. **Git status**: Confirm all changes committed to `https://git.uk.home.arpa/itdlabs/holo-deck.git`
4. **API validation**: `make test` to verify all endpoints functional
5. **Professional standards**: Confirm no emojis, absolute paths maintained

## Project Philosophy

> **"Where immersive holodeck technology meets professional engineering"**

THD represents the revolutionary evolution from basic 3D visualization to **professional VR/AR holodeck platform**, powered by A-Frame WebXR while maintaining all professional engineering standards and 100% API compatibility.

**A-Frame Integration Philosophy:**
- **Framework Flexibility**: Clean API layer enables multi-backend architecture
- **Community Leverage**: Utilize best-in-class open-source WebXR framework
- **Professional Standards**: Maintain engineering quality across all integrations
- **Future-Ready**: WebXR standard compliance for long-term viability

**Never compromise on professional standards while achieving revolutionary capabilities.**