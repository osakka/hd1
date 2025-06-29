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

### Current State (v3.6.0 - Ultra-Simple Scene Updates & Code Audit Perfection)
- **Binary**: `thd` (professional VR holodeck daemon)
- **Module**: `holodeck` (Go module name)
- **Rendering Engine**: A-Frame WebXR 1.4.0 (100% Local - 2.5MB ecosystem)
- **🎯 Ultra-Simple Architecture**: API-based scene updates replacing complex fsnotify monitoring
- **🔍 Code Audit Excellence**: Surgical precision audit achieving zero ambiguity and duplicates
- **🏆 Single Source of Truth**: Perfect specification-driven architecture validated
- **✅ Clean Build System**: Zero warnings, all auto-generation verified current
- **🌐 Scene System**: Ultra-simple API-driven scene loading with WebSocket connection trigger
- **📋 Professional Standards**: Complete elimination of duplicate ADR files and ambiguity
- **🏆 Revolutionary Integration**: Complete upstream/downstream API bridge system maintained
- **👑 Crown Jewel**: Auto-generated web UI client achieving 100% single source of truth
- **🎭 Scene Forking System**: Revolutionary "photo vs video" content creation paradigm
- **⚡ Bulletproof Reliability**: Simple API calls immune to filesystem mount option constraints
- **Holodeck Containment**: Escape-proof boundary enforcement with 60fps monitoring
- **Session Architecture**: Single source of truth with perfect isolation
- **A-Frame Components**: Complete local library with zero CDN dependencies
- **Git Repository**: `https://git.uk.home.arpa/itdlabs/holo-deck.git`
- **Credentials**: `claude-3/claude-password`
- **Directory Structure**: Professional organization with clean .gitignore

### Revolutionary Technical Architecture
- **🏆 Upstream/Downstream Integration**: Single source of truth bridge between THD API and A-Frame capabilities
- **Enhanced Shell Functions**: Complete A-Frame schema validation with professional parameter checking
- **JavaScript Function Bridge**: Identical signatures to shell functions enabling seamless API usage
- **Revolutionary Code Generation**: Unified generator producing both standard and enhanced clients
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
- **👑 Crown Jewel Web UI**: Complete JavaScript client auto-generated from OpenAPI spec
- **Three-Tier Generation**: Go router + CLI client + Web UI client all from single spec
- **Scene Forking Architecture**: Photo mode (snapshots) and video mode (temporal recording)
- **Object Lifecycle Tracking**: Complete provenance with base/modified/new state transitions
- **Professional Build System**: Make-based with daemon control targets
- **A-Frame Integration**: Seamless WebXR with 100% API compatibility

## Key File Locations

### Source Code
- `/opt/holo-deck/src/api.yaml` - Single source of truth specification
- `/opt/holo-deck/src/main.go` - THD daemon with professional standards
- `/opt/holo-deck/src/auto_router.go` - Auto-generated routing (DO NOT EDIT)
- `/opt/holo-deck/src/server/handlers.go` - A-Frame WebXR interface with FREEZE-FRAME/TEMPORAL SEQUENCE controls
- `/opt/holo-deck/src/codegen/generator.go` - 👑 Crown jewel code generator with web UI generation
- `/opt/holo-deck/share/htdocs/static/js/thd-aframe.js` - A-Frame holodeck manager
- `/opt/holo-deck/share/htdocs/static/js/thd-api-client.js` - 👑 Auto-generated JavaScript API client
- `/opt/holo-deck/share/htdocs/static/js/thd-ui-components.js` - 👑 Auto-generated UI components
- `/opt/holo-deck/share/htdocs/static/js/thd-form-system.js` - 👑 Auto-generated dynamic form system
- `/opt/holo-deck/src/Makefile` - Professional build system

### Runtime
- `/opt/holo-deck/build/bin/thd` - Professional daemon binary
- `/opt/holo-deck/build/runtime/thd.pid` - Process ID file
- `/opt/holo-deck/build/logs/` - Professional logging directory

### Documentation
- `/opt/holo-deck/README.md` - Professional project overview
- `/opt/holo-deck/CHANGELOG.md` - Complete project transformation history
- `/opt/holo-deck/docs/adr/ADR-001-aframe-webxr-integration.md` - A-Frame integration decision
- `/opt/holo-deck/docs/adr/ADR-002-Specification-Driven-Development.md` - OpenAPI specification-driven architecture
- `/opt/holo-deck/docs/adr/ADR-003-Professional-UI-Enhancement.md` - UI excellence & scene management
- `/opt/holo-deck/docs/adr/ADR-004-Scene-Forking-System.md` - Scene forking and photo/video paradigm
- `/opt/holo-deck/docs/adr/ADR-005-ultra-simple-scene-updates.md` - Ultra-simple scene update architecture
- `/opt/holo-deck/docs/adr/ADR-006-Auto-Generated-Web-UI-Client.md` - 👑 Crown jewel implementation
- `/opt/holo-deck/docs/adr/ADR-007-Revolutionary-Upstream-Downstream-Integration.md` - 🏆 Revolutionary API integration
- `/opt/holo-deck/docs/adr/ADR-008-Thread-Safe-Session-Store.md` - Thread-safe session management
- `/opt/holo-deck/docs/adr/ADR-009-WebSocket-Realtime-Architecture.md` - WebSocket real-time architecture
- `/opt/holo-deck/docs/adr/ADR-010-3D-Coordinate-System.md` - 3D coordinate system design
- `/opt/holo-deck/docs/adr/ADR-011-Build-System-Validation.md` - Build system validation architecture
- `/opt/holo-deck/docs/api/README.md` - THD API documentation

### Holodeck Libraries
- `/opt/holo-deck/lib/thd-functions.sh` - Comprehensive shell function library
- `/opt/holo-deck/lib/thd-enhanced-functions.sh` - 🏆 Enhanced shell functions with A-Frame integration
- `/opt/holo-deck/lib/thd-enhanced-bridge.js` - 🏆 JavaScript function bridge with identical signatures
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
- **Phase 8 (RESOLVED)**: Session restoration loop causing object flickering
- **Phase 9 (RESOLVED)**: Hardcoded scene metadata replaced with dynamic discovery
- **Phase 10 (IMPLEMENTED)**: 👑 Crown jewel - Auto-generated web UI client system

### Professional Standards Implementation
- **🎯 Ultra-Simple Solutions** - elegant API-based designs over complex infrastructure
- **🔍 Single Source of Truth** - zero duplicate files or ambiguous documentation
- **🚫 No emojis** in any system output or logging
- **📁 Absolute paths only** - all configured from THD_* constants
- **🏷️ Long flags only** - no short flags to eliminate confusion
- **📝 Professional error messages** - clear, actionable, no decorative elements
- **🌐 Clean URLs** - no query string versioning, standards-compliant caching
- **⚡ Surgical Precision** - bar-raising solutions with zero regression tolerance

### Breaking Changes from VWS to THD
- **Module imports**: `visualstream/*` → `holodeck/*`
- **Binary names**: `vws` → `thd`, `vws-client` → `thd-client`
- **Constants**: `VWS_*` → `THD_*`
- **PID files**: `vws.pid` → `thd.pid`

## Ultra-Simple Scene Update Context (v3.6.0)

### Revolutionary Scene Update Architecture
- **🎯 Ultra-Simple Design**: API-based scene loading instead of complex fsnotify monitoring
- **🔍 Filesystem Discovery**: Root cause analysis revealed `noatime,lazytime` mount options interfering with fsnotify
- **🌐 WebSocket Trigger**: `setTimeout(refreshSceneDropdown, 1000)` on connection for automatic scene discovery
- **📊 Complete Scene Detection**: All 11+ scenes automatically discovered from filesystem via API
- **👤 Natural Workflow**: Page refresh pattern users expect instead of complex file monitoring
- **⚡ Bulletproof Reliability**: Simple API calls immune to filesystem mount option constraints

### Professional Code Audit Excellence
- **🔧 Surgical Precision**: Complete elimination of duplicate ADR-005 files
- **🎯 Zero Ambiguity**: 100% single source of truth achieved throughout codebase  
- **✅ Clean Build Verification**: `make all` and `make test` pass with zero warnings
- **🤖 Auto-Generation Validation**: All clients perfectly generated from api.yaml specification
- **📋 Documentation Audit**: Complete elimination of conflicting or redundant documentation
- **🏆 Bar-Raising Standards**: Ultra-simple solutions preferred over complex infrastructure

### Ultra-Simple Implementation Details
```javascript
// Ultra-simple scene loading on WebSocket connection
setTimeout(refreshSceneDropdown, 1000);

// Perfect scene discovery from API
async function refreshSceneDropdown() {
    const response = await fetch('/api/scenes');
    const data = await response.json();
    
    // Update dropdown with all scenes from filesystem
    // Preserve saved scene selection from cookies
    // Zero complexity, 100% reliability
}
```

### Filesystem Mount Constraints Resolution
- **Root Cause**: Filesystem mounted with `noatime,lazytime` options
- **Technical Impact**: Mount options interfere with fsnotify reliability
- **Ultra-Simple Solution**: API-based approach more reliable than filesystem events
- **Professional Decision**: User-driven workflow instead of complex monitoring
- **Future Ready**: Infrastructure preserved for when filesystem constraints resolved

## Scene Forking & Crown Jewel Context (v3.4.0)

### Revolutionary Scene Forking System
- **Photo Mode (FREEZE-FRAME)**: `POST /sessions/{sessionId}/scenes/save` - Save current session state as new scene
- **Video Mode (TEMPORAL SEQUENCE)**: Recording endpoints for complete session capture and playback
- **Scene Fork API**: `POST /scenes/{sceneId}/fork` - Load scenes into sessions for non-destructive editing
- **Object Tracking**: Complete provenance system (base/modified/new) with source scene references
- **Dynamic Scene Discovery**: Script-based metadata parsing replacing hardcoded scene lists
- **Script Generation**: Automatic creation of executable scene files from session state

### 👑 Crown Jewel Implementation - Auto-Generated Web UI Client
- **JavaScript API Client**: Complete API wrapper auto-generated from OpenAPI specification
- **UI Components**: Each API endpoint becomes an interactive UI component
- **Dynamic Form System**: Forms automatically generated from request schemas
- **100% Single Source of Truth**: All clients (Go router, CLI, Web UI) generated from same spec
- **Zero Manual Synchronization**: API changes automatically update all client systems

### Scene Management System (Legacy v3.3.0 Context)
- **API Endpoints**: `/api/scenes` (list) and `/api/scenes/{sceneId}` (load)
- **Dynamic Discovery**: Script-based metadata parsing replaces hardcoded scenes
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

### Scene Forking & Crown Jewel Debugging
- **Scene Forking**: Check object tracking status (base/modified/new) in session object listings
- **FREEZE-FRAME**: Monitor `POST /sessions/{sessionId}/scenes/save` for scene creation
- **TEMPORAL SEQUENCE**: Check recording status with `GET /sessions/{sessionId}/recording/status`
- **Crown Jewel Generation**: Verify auto-generated files in `/share/htdocs/static/js/thd-*.js`
- **UI Components**: Test auto-generated components via `window.THDUIComponents` in browser console
- **API Client**: Verify `window.thdAPI` provides all 23 auto-generated methods
- **Form System**: Check `window.THDFormSystem` for schema-driven form generation
- **Build Generation**: Ensure `make generate` produces all crown jewel files without errors

### UI Enhancement Debugging (Legacy)
- **Cache Issues**: Check HTTP headers in Network tab (should show no-cache for dev)
- **Scene Loading**: Monitor console for `SCENE_LOADED` and `AUTO_SCENE` debug entries
- **Sprint Controls**: Verify `thd-sprint-controls` component attachment on camera entity
- **Status Indicators**: Professional `[ACTIVE]`/`[MINIMIZED]` instead of unicode arrows
- **Scrollbar Theme**: CSS custom properties for cross-browser holodeck styling

## Recovery Context for New Sessions

When resuming development after session restart:

1. **Check daemon status**: `cd /opt/holo-deck/src && make status`
2. **Verify build**: `make all` to ensure clean build with zero warnings
3. **Git status**: Confirm all changes committed to `https://git.uk.home.arpa/itdlabs/holo-deck.git`
4. **API validation**: `make test` to verify all endpoints functional
5. **Professional standards**: Confirm no emojis, absolute paths maintained
6. **🎯 Ultra-simple architecture**: Verify scene loading works via API (not fsnotify)
7. **🔍 Single source validation**: Confirm no duplicate ADR files or ambiguous documentation

## Project Philosophy

> **"Where immersive holodeck technology meets professional engineering"**

THD represents the revolutionary evolution from basic 3D visualization to **professional VR/AR holodeck platform**, powered by A-Frame WebXR while maintaining all professional engineering standards and 100% API compatibility.

**A-Frame Integration Philosophy:**
- **Framework Flexibility**: Clean API layer enables multi-backend architecture
- **Community Leverage**: Utilize best-in-class open-source WebXR framework
- **Professional Standards**: Maintain engineering quality across all integrations
- **Future-Ready**: WebXR standard compliance for long-term viability

**Never compromise on professional standards while achieving revolutionary capabilities.**