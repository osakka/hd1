# HD1 Changelog

All notable changes to HD1 (Holodeck One) are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.7.2] - 2025-07-18

### Major Enhancements - Complete Avatar Lifecycle Management
This release achieves surgical precision avatar lifecycle management with automatic cleanup, mobile touch controls, and comprehensive single source of truth architecture.

### Added
- **Avatar Lifecycle Management**: Automatic cleanup of inactive avatars via session inactivity timeout
- **Mobile Touch Controls**: Left side for movement, right side for camera look with seamless desktop/mobile experience
- **Session-Avatar Integration**: Avatar registry connected to session cleanup worker for surgical precision cleanup
- **Pointer Lock Security**: Fixed user gesture requirements for browser security compliance
- **Enhanced Error Handling**: Comprehensive error handling for avatar operations and cleanup processes

### Changed
- **Avatar Cleanup Architecture**: Session manager now drives avatar cleanup through interface pattern
- **Mobile UX**: Touch controls enable full 3D navigation without external keyboards
- **Pointer Lock Behavior**: Escape key now only exits (no entry) to comply with browser security
- **Code Generation**: Auto-generated API client method ordering improved for consistency

### Fixed
- **Browser Security Compliance**: Removed pointer lock request from keyboard events (user gesture required)
- **Avatar Persistence**: Avatars now properly cleaned up when sessions become inactive
- **Single Source of Truth**: Avatar registry cleanup integrated with session database lifecycle
- **Mobile Touch Navigation**: Proper touch event handling for both movement and camera controls

### Technical Details
- **Interface Pattern**: `AvatarRegistryInterface` enables clean session-avatar coupling
- **Cleanup Query**: Pre-query inactive participants before database update for avatar cleanup
- **Touch Controls**: Screen-split approach (left=move, right=look) for intuitive mobile UX
- **Error Prevention**: Eliminated "NotAllowedError" and "SecurityError" from pointer lock operations

### Architecture Validation
- ✅ **Single Source of Truth**: Session database drives avatar lifecycle
- ✅ **No Parallel Implementations**: All avatar cleanup flows through session manager
- ✅ **Mobile-First**: Touch controls work seamlessly alongside desktop controls
- ✅ **Security Compliant**: Proper browser security model compliance
- ✅ **Zero Regressions**: All existing functionality preserved and enhanced

## [0.7.1] - 2025-07-17

### Critical Fixes - Surgical Precision Avatar Movement
This release completes the surgical precision single source of truth avatar movement system with full API-first architecture compliance.

### Fixed
- **Avatar Movement API Client**: Added missing `extractPathParams` method to JavaScript API client template
- **Avatar Operation Data Structure**: Fixed `session_id` vs `client_id` handling in Three.js avatar operations
- **Global API Client Access**: Made API client globally available for Three.js integration
- **Avatar Movement Logging**: Restored comprehensive console logging and error handling
- **Avatar Remove Operation**: Fixed avatar removal handling in Three.js scene manager

### Technical Details
- **Single Source of Truth**: Validated 100% compliance with no parallel implementations
- **API-First Architecture**: All avatar movement flows through `/avatars/{sessionId}/move` endpoint
- **Surgical Precision**: Zero regressions, zero hacks, zero architectural compromises
- **Code Generation**: Fixed template to include proper path parameter extraction
- **Build System**: Maintained clean build with zero warnings

### Architecture Validation
- ✅ **One Source of Truth**: `/avatars/{sessionId}/move` API endpoint only
- ✅ **No Parallel Implementations**: All movement flows through sync system
- ✅ **No WebSocket Bypasses**: All avatar updates come from sync operations
- ✅ **No Manual Updates**: No code bypasses API/sync architecture
- ✅ **Zero Ambiguity**: Clean, consistent, and foolproof implementation

## [0.7.0] - 2025-07-16

### Major Changes - Clean HD1-Native API Architecture
This release completes the transformation to a clean, intuitive HD1-native API architecture with WebSocket-driven avatar lifecycle management.

### Added
- **Clean HD1-Native API Paths**: Resource-focused endpoints without implementation details
  - `/api/avatars` (previously `/api/threejs/avatars`)
  - `/api/entities` (previously `/api/threejs/entities`)
  - `/api/scene` (previously `/api/threejs/scene`)
- **WebSocket-Driven Avatar Lifecycle**: Automatic avatar creation/cleanup on connection/disconnection
- **Avatar Registry System**: In-memory avatar state management tied to WebSocket connections
- **Resource-Based File Organization**: Clean separation of concerns with dedicated handler packages
- **Shared API Utilities**: Common helper functions for consistent request handling

### Changed
- **BREAKING**: API paths no longer expose Three.js implementation details
- **Architecture**: Resource-based handler organization (`api/avatars/`, `api/entities/`, `api/scene/`)
- **Avatar Management**: Automatic lifecycle tied to WebSocket connections instead of manual creation
- **Code Generation**: Updated templates for clean package prefixes and method names
- **JavaScript Client**: Intuitive method names (`api.getAvatars()` vs `api.getThreeJSAvatars()`)
- **Version Numbering**: Standardized to semantic versioning (0.x.x format)

### Improved
- **Developer Experience**: Clean, intuitive API paths that feel HD1-native
- **Real-time Updates**: Avatar position/animation updates through registry system
- **Code Organization**: Logical separation of API concerns by resource type
- **Maintainability**: Shared utilities reduce code duplication across handlers

### Fixed
- **Consistent Versioning**: Aligned all version references to semantic versioning scheme
- **API Response Consistency**: Standardized response formats across all endpoints
- **WebSocket Integration**: Seamless avatar lifecycle management without manual intervention

## [0.6.0] - 2025-07-14

### Major Changes - Ultra-Minimal Build Strategy
This release represents a complete architectural transformation focusing on ultra-minimal Three.js console implementation with massive codebase optimization.

### Added
- **Ultra-Minimal Three.js Console**: Zero-framework debugging interface with essential features only
- **Rebootstrap Recovery System**: Intelligent storage clearing and page reload on connection failures  
- **Specification-Driven Development**: Complete code generation from OpenAPI specification
- **Comprehensive Documentation**: Full documentation suite with ADRs, guides, and architecture docs
- **Configuration Management**: Hierarchical configuration with environment variables, flags, and .env support
- **WebSocket Synchronization Protocol**: High-performance real-time communication with session isolation
- **Object Pooling**: Memory optimization for high-frequency operations
- **Intelligent Reconnection**: Exponential backoff with automatic recovery mechanisms

### Changed
- **BREAKING**: Removed A-Frame framework dependencies (replaced with Three.js)
- **BREAKING**: Eliminated CLI client generation system
- **BREAKING**: Simplified avatar asset handling (removed WebSocket-based GLB streaming)
- **Architecture**: Complete migration to API-first, specification-driven development
- **Build System**: Streamlined Makefile with focused development workflow
- **Console Interface**: Replaced complex dashboard with minimal debug panel
- **WebSocket Protocol**: Optimized message handling with session-based routing
- **Logging System**: Enhanced structured logging with module-based tracing
- **Configuration**: Unified configuration system with clear priority hierarchy

### Removed
- **A-Frame WebVR Framework**: Complete removal of A-Frame dependencies and VR support
- **CLI Client System**: Removed 350+ lines of CLI generation code
- **UI Component Generation**: Removed dynamic UI component generation system (200+ lines)
- **Form System Generation**: Removed automatic form generation (150+ lines)  
- **Complex Avatar Asset Pipeline**: Removed base64 GLB streaming via WebSocket
- **Redundant Dependencies**: Removed unused imports and dependencies throughout codebase
- **Legacy Template System**: Simplified template architecture with external templates
- **Semantic UI Integration**: Removed unused semantic UI components
- **Shell Function Generation**: Removed shell integration code generation

### Performance Improvements
- **Code Reduction**: Generator reduced from 803 to 538 lines (-33%)
- **Total Codebase**: Reduced from 8000+ to 4989 lines of Go code (-38%)
- **File Count**: Reduced from 50+ to 27 core source files (-46%)
- **Build Time**: Improved build performance with streamlined generation process
- **Memory Usage**: Reduced allocation overhead through object pooling
- **WebSocket Efficiency**: Optimized message broadcasting and connection handling

### Technical Details

#### Ultra-Minimal Console Implementation
```javascript
// Zero-framework Three.js console with essential features
class HD1Console {
    constructor() {
        this.setupWebSocket();      // Real-time connection monitoring
        this.setupDebugPanel();     // Collapsible system information
        this.setupRebootstrap();    // Nuclear recovery option
    }
}
```

#### Specification-Driven Architecture
```yaml
# Complete system definition in src/api.yaml
paths:
  /threejs/entities:
    post:
      operationId: createEntity
      x-handler: api/threejs/create.go
      x-function: CreateEntity
```

#### WebSocket Protocol Optimization
```go
// High-performance message broadcasting with session isolation
func (h *Hub) broadcastToSession(sessionID string, message []byte) {
    for _, client := range h.sessions[sessionID] {
        select {
        case client.send <- message: // Non-blocking send
        default: h.unregisterClient(client) // Handle blocked clients
        }
    }
}
```

### Migration Guide
For users migrating from v5.x:

1. **A-Frame Removal**: Update any A-Frame-specific code to use Three.js directly
2. **CLI Client**: Remove any dependencies on generated CLI clients
3. **Avatar Assets**: Switch to HTTP-based asset loading instead of WebSocket streaming
4. **Configuration**: Update environment variables to use HD1_ prefix consistently
5. **WebSocket Messages**: Review custom message types for compatibility

### Documentation
- **Complete Documentation Suite**: Added comprehensive docs/ directory
- **Architecture Decision Records**: 5 ADRs documenting key design decisions
- **User Guides**: Quick start, development, configuration, and troubleshooting guides
- **README Files**: Project, docs, and src-specific documentation
- **API Reference**: Enhanced OpenAPI specification documentation

---

## [5.0.6] - Previous Release

### Added
- Avatar control system with seamless transitions
- Seamless world transition recovery  
- GLB avatar asset loading with proper resource handling
- Advanced camera system with orbital mode
- Console UI with smooth animations
- Vendor cleanup and optimized structure

### Changed
- Enhanced avatar management with real GLB models
- Improved WebSocket avatar synchronization
- Streamlined configuration management

---

## Version History Summary

- **v6.0.0**: Ultra-minimal Three.js console with specification-driven architecture
- **v5.0.6**: Production avatar control system with seamless world transitions  
- **v5.0.5**: Configuration management standardization
- **v5.0.x**: Avatar system development and WebSocket optimization
- **v4.x.x**: Initial Three.js integration and API development
- **v3.x.x**: WebSocket implementation and real-time synchronization
- **v2.x.x**: API-first architecture establishment
- **v1.x.x**: Initial holodeck platform development

---

## Development Metrics

### v6.0.0 Optimization Results
- **Lines of Code**: 8000+ → 4989 (-38% reduction)
- **Source Files**: 50+ → 27 (-46% reduction)  
- **Code Generator**: 803 → 538 lines (-33% reduction)
- **Build Time**: <10 seconds (improved)
- **Binary Size**: Optimized with dependency removal
- **Memory Usage**: Reduced through object pooling and cleanup

### Quality Improvements
- **Zero Manual Routing**: 100% specification-driven route generation
- **Test Coverage**: Enhanced with simplified codebase
- **Documentation**: Complete documentation suite added
- **Maintainability**: Significant improvement through code reduction
- **Performance**: Optimized hot paths and memory management

---

*HD1 v6.0.0: Where OpenAPI specifications become ultra-minimal Three.js game engine platforms.*