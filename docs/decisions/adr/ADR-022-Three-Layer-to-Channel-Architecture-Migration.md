# ADR-022: Three-Layer to Channel Architecture Migration (v5.0.0)

## Status
Accepted - Implemented 2025-07-03

## Context
HD1 v4.0.0 and earlier used a complex "three-layer architecture" (Environment + Props + Scene) with shell script-based scene generation. This system proved incompatible with the PlayCanvas migration and modern game engine requirements.

### Legacy Three-Layer Architecture Problems
```
Problems with v4.0.0 Architecture:
├── Environment System
│   ├── Shell script execution (security risk)
│   ├── Complex physics mapping
│   └── Hardcoded environment types
├── Props System  
│   ├── YAML files with embedded shell scripts
│   ├── A-Frame specific implementations
│   └── Manual physics adaptation
└── Scene System
    ├── Shell script generation
    ├── Legacy file monitoring
    └── Complex state management
```

### Key Issues
- **Shell script execution**: Security and maintenance concerns
- **A-Frame dependency**: Incompatible with PlayCanvas migration
- **Complex state management**: Multiple sources of truth
- **Performance**: Shell script interpretation overhead
- **Maintainability**: Mixed YAML and shell script files
- **Scalability**: Rigid environment-props coupling

## Decision
Replace the three-layer architecture with a modern channel-based YAML configuration system that directly configures PlayCanvas scenes.

### New Channel Architecture
```
HD1 v5.0.0: Channel-Based Configuration
├── Channel Definition (YAML)
│   ├── Channel metadata
│   ├── PlayCanvas scene settings
│   └── Entity definitions
├── Direct PlayCanvas Integration
│   ├── Native component system
│   ├── Direct physics configuration
│   └── Entity-component-system
└── API-First Management
    ├── Channel CRUD operations
    ├── Real-time configuration
    └── WebSocket synchronization
```

### Channel Configuration Format
```yaml
# Example: channel_one.yaml
channel:
  id: "channel_one" 
  name: "Scene 1 - Red Box"
  description: "Minimalist scene with red box and overhead lighting"
  max_clients: 100

playcanvas:
  scene:
    ambientLight: [0.3, 0.3, 0.3]  # Direct PlayCanvas settings
    gravity: [0, -9.81, 0]          # No environment abstraction
    
  entities:
    - name: "floor"
      components:
        model: {type: "plane"}
        transform:
          position: [0, -0.1, 0]
          scale: [10, 1, 10]
        material: {diffuse: "#cccccc"}
        rigidbody: {type: "static"}
        collision: {type: "box"}
        
    - name: "red_box"
      components:
        model: {type: "box"}
        transform: {position: [0, 1, 0]}
        material: {diffuse: "#ff0000"}
        rigidbody: {type: "dynamic", mass: 1.0}
        collision: {type: "box"}
```

## Implementation Evidence

### File System Migration
```
REMOVED (Legacy Three-Layer):
- share/environments/earth-surface.sh
- share/environments/molecular-scale.sh  
- share/environments/space-vacuum.sh
- share/environments/underwater.sh
- share/props/furniture/wooden-chair.yaml
- share/props/electronic/lightbulb.yaml
- share/props/decorative/flower-pot.yaml
- share/scenes/README.md (entire directory)
- server/scenes_watcher.go (legacy file monitoring)

ADDED (Channel System):
+ share/channels/channel_one.yaml
+ share/channels/channel_two.yaml
+ share/channels/channel_three.yaml
+ share/channels/config.yaml
+ api/channels/ (complete CRUD API)
```

### API Transformation
```
Legacy APIs (REMOVED):
- POST /environments/{environmentId}
- GET /environments
- POST /props/{propId}
- GET /props
- POST /scenes/{sceneId}
- GET /scenes

Modern APIs (ADDED):
- GET /channels
- POST /channels
- GET /channels/{channelId}
- PUT /channels/{channelId}
- DELETE /channels/{channelId}
- POST /sessions/{sessionId}/channel/join
- POST /sessions/{sessionId}/channel/leave
```

### Code Removal Evidence
```go
// REMOVED: Legacy shell script generation
func createChannelScript() {
    scriptContent := fmt.Sprintf(`#!/bin/bash
    # Apply channel environment
    hd1::api_call "POST" "/environments/%s"
    `, environment)
    // ... 30+ lines of shell script generation
}

// REMOVED: Legacy scenes watcher
type ScenesWatcher struct {
    // Monitored .sh files in deleted directory
}
```

## Migration Benefits

### Technical Improvements
- **Security**: No shell script execution
- **Performance**: Direct PlayCanvas configuration 
- **Maintainability**: Pure YAML configuration
- **Scalability**: Channel-based multi-tenancy
- **Simplicity**: Single source of truth per channel
- **Modern**: Entity-component-system paradigm

### Architectural Advantages
- **PlayCanvas native**: Direct engine configuration
- **API-first**: Everything configurable via REST
- **Real-time**: WebSocket channel synchronization
- **Flexible**: No rigid environment-props coupling
- **Extensible**: Easy addition of new components

### Developer Experience
- **Clear configuration**: Human-readable YAML
- **IDE support**: Syntax highlighting and validation
- **Version control**: Clean diffs for configuration changes
- **Testing**: Isolated channel configurations
- **Debugging**: Clear configuration vs runtime separation

## Migration Process

### Phase 1: Channel System Implementation
1. **Create channel YAML structure**: Define configuration schema
2. **Implement channel APIs**: CRUD operations for channels
3. **PlayCanvas integration**: Direct engine configuration
4. **WebSocket support**: Real-time channel management

### Phase 2: Legacy System Removal
1. **Remove shell scripts**: Delete all .sh files
2. **Remove YAML props**: Delete embedded shell script props
3. **Remove environment APIs**: Delete legacy endpoints
4. **Remove scenes watcher**: Delete file monitoring system

### Phase 3: Validation
1. **Build system**: Verify 82 endpoints with handlers
2. **API testing**: Validate all channel operations
3. **Scene rendering**: Verify PlayCanvas integration
4. **Documentation**: Update all references

## Risk Assessment

### Migration Risks
- **Breaking changes**: Complete API restructuring
- **Data loss**: Legacy scenes need manual conversion
- **Compatibility**: New clients required for v5.0.0

### Mitigation Strategies
- **Build validation**: Prevent incomplete migrations
- **API versioning**: Clear v4.0.0 vs v5.0.0 distinction
- **Documentation**: Complete migration guide
- **Testing**: Comprehensive channel system validation

## Performance Impact

### Before (Three-Layer System)
- **Shell execution**: Process spawning overhead
- **File monitoring**: Filesystem watcher complexity
- **State management**: Multiple configuration sources
- **Memory usage**: Shell script interpretation

### After (Channel System)
- **Direct configuration**: Native PlayCanvas setup
- **API-driven**: REST-based configuration
- **Single source**: Channel YAML as truth
- **Memory efficiency**: Direct engine integration

## Consequences

### Positive
- **Modern architecture**: Industry-standard configuration
- **Security**: No shell script execution vectors
- **Performance**: Direct PlayCanvas integration
- **Maintainability**: Pure YAML configuration
- **Flexibility**: Component-based entity system
- **API completeness**: Full channel management via REST

### Negative
- **Migration effort**: All legacy scenes need conversion
- **Learning curve**: New configuration paradigm
- **Breaking changes**: v4.0.0 scenes incompatible

### Long-term Benefits
- **Scalability**: Channel-based multi-tenancy
- **Extensibility**: Easy component additions
- **Standards compliance**: Modern game engine patterns
- **Market position**: Professional game development platform

## Future Evolution
- **Component library**: Standardized component definitions
- **Scene templates**: Reusable channel configurations
- **Visual editor**: GUI-based channel configuration
- **Asset management**: Integrated resource handling

## References
- ADR-021: PlayCanvas Migration Implementation
- PlayCanvas Entity-Component-System documentation
- HD1 v5.0.0 CHANGELOG: Legacy scene generation system removal
- Channel configuration examples: `/share/channels/`

---
**Decision made by**: HD1 Development Team  
**Implementation date**: 2025-07-03  
**Version**: v5.0.0 - API-FIRST GAME ENGINE REVOLUTION  
**Supersedes**: ADR-014 Three-Layer Architecture System  
**Review date**: Next major architecture evolution