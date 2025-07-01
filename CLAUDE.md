# HD1 (Holodeck One) - Development Context & Memory

## Core Development Principles

- **Tied API Architecture** - HD1 exposes everything via unified API surface (platform + downstream APIs)
- **API-first development** from our spec yaml - 100% API-driven service
- **Scripts are bidirectional bridges** - API targets AND API consumers
- **WebSocket for state sync only** - all commands flow through API
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

## Logging Standards & Troubleshooting

### Single Source of Truth Logging Architecture
HD1 uses a unified logging system (`holodeck1/logging`) across the entire codebase. **No standard library `log` package usage anywhere.**

### Log Format Standard
**Format**: `timestamp [processid:threadid] [level] functionname.filename(without extension) line_num: short message goes here`

**Example**: `2025-07-01T10:30:45.123456789Z [1234:g1] [INFO] main.main:15 server starting`

### Log Levels & Audience Targeting
- **TRACE**: Development only - function-level debugging with module filtering (`logging.Trace("module", "message", data)`)
- **DEBUG**: Development environment - detailed operational information
- **INFO**: Production SRE - normal operational events worth noting
- **WARN**: Production SRE - potential issues that need attention  
- **ERROR**: Production SRE - error conditions requiring immediate action
- **FATAL**: All environments - fatal errors that terminate the process

### Dynamic Log Level Control
Supports real-time log level adjustment via:
- **API**: `POST /admin/logging/level {"level": "DEBUG"}`
- **Environment**: `HD1_LOG_LEVEL=DEBUG`
- **Flag**: `--log-level debug`

### Trace Module System
Enable/disable tracing per functionality:
```bash
# API control
POST /admin/logging/trace {"modules": ["sessions", "objects"]}

# Environment variable  
HD1_TRACE_MODULES=sessions,objects,websocket
```

### Thread Safety & Performance
- **Thread-safe**: All logging operations use sync.RWMutex for concurrent access
- **Zero-overhead when disabled**: Log levels below threshold consume no CPU cycles
- **Structured JSON**: Machine-readable with optional human-readable console output
- **Log rotation**: 10MB max size, 3 rotated files retained

### Message Quality Standards
- **Concise**: Remove redundant context (function/file already captured)
- **Actionable**: Messages should point to specific issues
- **No decorations**: No emojis, "SUCCESS:", "[TAGS]", etc.
- **Error separation**: System logging separate from user console output

### Troubleshooting Workflow
1. **Increase verbosity**: `curl -X POST /api/admin/logging/level -d '{"level":"TRACE"}'`
2. **Enable module tracing**: `curl -X POST /api/admin/logging/trace -d '{"modules":["target_module"]}'`
3. **Check structured logs**: `tail -f /opt/hd1/build/logs/hd1.log | jq .`
4. **Restore normal levels**: Return to INFO/WARN for production

### Implementation Examples
```go
// Good - structured logging with relevant context
logging.Error("session creation failed", map[string]interface{}{
    "user_id": userID,
    "error": err.Error(),
})

// Bad - redundant context and decorative elements  
logging.Error("ERROR: failed to create session in CreateSessionHandler()", map[string]interface{}{
    "function": "CreateSessionHandler", // redundant - already captured
    "message": "Session creation failed", // redundant - already in message
})
```

### Pre-logging Initialization
Use `fmt.Fprintf(os.Stderr, "FATAL: message")` **only** for errors before logging system initialization. All post-initialization fatal errors use `logging.Fatal()`.

## Current System State (2025-07-01)

### Revolutionary Three-Layer Architecture Achievement
- **Complete Game Engine Architecture** - Environment + Props + Scene layers matching Unity/Unreal patterns
- **Environment System** - 4 physics contexts with realistic material adaptation
- **Props System** - 6 categories of reusable objects with YAML-based definitions  
- **Physics Cohesion Engine** - Props automatically adapt to environment physics
- **API Extension** - 31 endpoints including comprehensive three-layer APIs
- **Single Source of Truth** - All three-layer functionality auto-generated from specification

### Recent Major Achievements  
- **Complete Lighting System Implementation** - API-first lighting with WebSocket reactivity and A-Frame integration
- **Session Restoration Fix** - Per-client object synchronization replacing global session locks
- **Method Context Binding Resolution** - Fixed JavaScript prototype chain issues in reactive system
- **WebSocket Message Handler Completion** - Added missing `prop_instantiated` handler for full reactivity
- **API Specification Alignment** - Unified object creation format across client/server boundaries
- **Object Visibility System** - Clean, API-first visibility management with client validation and toggle
- **Infinite WebSocket Resilience** - Graceful error handling with intelligent rebootstrapping after 99 attempts
- **Complete THD → HD1 transformation** across entire codebase
- **Template Architecture Revolution** - Surgically externalized 8 hardcoded templates (2,000+ lines) to maintainable external files
- **Zero Regression Refactoring** - Complete template system overhaul with identical output validation
- **Dynamic version system** with real-time API/JS version display
- **Enhanced console UI** with clickable status bar and dual arrow indicators
- **Professional branding** with "Holodeck I" title and version info
- **Color Persistence Architecture** - Objects maintain colors across session restoration
- **Reactive Scene Graph** - Comprehensive state management with rollback protection

### Active Features
- **Complete Lighting System**: Point, directional, ambient, and spot lights with API-first architecture
- **Per-Client Session Restoration**: Every client receives full session state on connection
- **WebSocket Reactive System**: Full message handling including prop instantiation events
- **Method Context Binding**: Resolved JavaScript prototype chain issues across reactive system
- **Three-Layer Architecture**: Environment + Props + Scene system fully operational
- **Environment APIs**: `/environments` (GET/POST) with 4 physics contexts
- **Props APIs**: `/props` (GET) and `/sessions/{id}/props/{propId}` (POST) with 6 categories
- **Physics Cohesion**: Real-time environment-aware physics adaptation
- **Web UI Console**: Dynamic version display "Holodeck I v1.0.0 aa74f3f3"
- **API Endpoints**: All 31 endpoints functional including three-layer architecture
- **WebSocket Integration**: Real-time communication working
- **A-Frame Integration**: 3D scene rendering operational
- **Build System**: Auto-generation and validation complete

### Git Repository Status
- **Branch**: master (up to date with origin)
- **Last Commit**: 73f7169 "refactor: Complete surgical code audit with comprehensive documentation update"
- **Status**: Working tree has technical documentation accuracy updates in progress
- **Remote**: https://git.uk.home.arpa/itdlabs/holo-deck.git

### Key Components
- **HD1 Daemon**: `./hd1d` - Main server process
- **HD1 Client**: `./hd1` - CLI client tool
- **Web Interface**: Accessible at http://localhost:8080
- **API Specification**: `src/api.yaml` - Single source of truth

### Development Commands
```bash
# Build system
cd src && make clean && make

# Start daemon
cd src && make start

# Test client
./build/bin/hd1-client --help

# Test three-layer architecture
./build/bin/hd1-client list-environments
./build/bin/hd1-client list-props

# Run scenes
./hd1 run-scene empty-grid
./hd1 run-scene ultimate-demo
```

### Three-Layer Architecture Achievements
1. ✅ Environment system with 4 physics contexts (Earth Surface, Molecular Scale, Space Vacuum, Underwater)
2. ✅ Props system with 5 categories and YAML-based definitions
3. ✅ Physics cohesion engine with real-time environment adaptation
4. ✅ Complete API extension with 31 endpoints including three-layer APIs
5. ✅ Comprehensive ADR-014 documentation and testing validation
6. ✅ Single source of truth maintained with auto-generated clients

## Architecture Notes

- **Three-layer game engine** - Environment + Props + Scene architecture matching Unity/Unreal patterns
- **Physics cohesion engine** - Props automatically adapt to environment physics with real-time recalculation
- **Single binary approach** - daemon handles all functionality including three-layer APIs
- **Auto-generated routing** - from OpenAPI spec to Go handlers (31 endpoints)
- **External template system** - 8 templates organized by language with Go embed filesystem
- **Template processing** - surgical variable substitution with caching optimization
- **WebSocket real-time** - bidirectional communication
- **A-Frame rendering** - WebXR-ready 3D scenes

### Template Architecture (2025-07-01)
Revolutionary externalization of code generation templates:

```
src/codegen/templates/
├── go/
│   ├── router.tmpl       # Auto-router generation
│   └── client.tmpl       # Go CLI client
├── javascript/
│   ├── api-client.tmpl   # JS API wrapper
│   ├── ui-components.tmpl # UI components
│   ├── form-system.tmpl  # Dynamic forms
│   └── aframe-bridge.tmpl # A-Frame integration
└── shell/
    ├── core-functions.tmpl    # Core API functions
    └── aframe-functions.tmpl  # Enhanced A-Frame
```

**Benefits:**
- **Maintainability**: Proper syntax highlighting and IDE support
- **Developer Experience**: Frontend developers can directly edit templates
- **Single Binary**: Go embed maintains deployment simplicity
- **Performance**: Template caching with `loadTemplate()` system
- **Zero Regression**: Identical output validated, surgical refactoring achieved

## Quality Assurance

- **Three-layer architecture** - complete game engine parity with Unity/Unreal patterns
- **Physics validation** - comprehensive testing of environment-props cohesion
- **API completeness** - all 31 endpoints including three-layer APIs functional
- **Zero THD references** - complete branding transformation
- **Specification compliance** - all endpoints match OpenAPI
- **Build validation** - prevents incomplete deployments
- **Real-time testing** - WebSocket and API verification

Revolutionary three-layer architecture complete with comprehensive documentation! 🏗️