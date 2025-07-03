# ADR-023: Legacy Code Elimination for v5.0.0 Clean Architecture

## Status
Accepted - Implemented 2025-07-03

## Context
HD1 v5.0.0 migration to PlayCanvas revealed extensive legacy code from previous A-Frame iterations that was incompatible with the new architecture. A comprehensive code audit identified multiple legacy systems that needed complete removal to achieve surgical precision and maintain single source of truth principles.

### Legacy Systems Identified
```
Legacy Components (Incompatible with v5.0.0):
├── A-Frame Shell Scripts
│   ├── Environment shell scripts (4 files)
│   ├── Scene generation scripts (8+ files) 
│   └── Props with embedded shell scripts
├── Legacy API Endpoints
│   ├── Shell script-based scene management
│   ├── Environment execution endpoints
│   └── Props instantiation via shell scripts
├── Legacy File Watchers
│   ├── Scenes directory monitoring (.sh files)
│   └── Props file system watching
└── Legacy Configuration
    ├── Hardcoded environment types
    ├── Shell script generation in channel creation
    └── Mixed YAML/shell script props
```

### Problems with Legacy Systems
- **Security vulnerabilities**: Shell script execution in web context
- **Performance overhead**: Process spawning and shell interpretation
- **Maintenance burden**: Mixed shell/YAML/Go codebase
- **PlayCanvas incompatibility**: A-Frame specific implementations
- **Build failures**: Missing handler files for legacy endpoints
- **Architectural inconsistency**: Multiple sources of truth

## Decision
Perform comprehensive surgical removal of all legacy code incompatible with PlayCanvas v5.0.0 architecture, maintaining zero regressions for supported functionality.

### Elimination Scope
1. **Complete A-Frame removal**: All shell script dependencies
2. **Legacy API cleanup**: Remove endpoints with missing handlers
3. **File system cleanup**: Remove obsolete directories and watchers
4. **Code audit**: Systematic legacy reference removal
5. **Documentation accuracy**: Update all references to reflect reality

### Surgical Precision Approach
- **Build-driven validation**: Use build failures to identify legacy references
- **Systematic code review**: File-by-file analysis and cleanup
- **Zero regression policy**: Maintain all v5.0.0 functionality
- **Single source of truth**: Ensure api.yaml drives all routing

## Implementation Evidence

### File System Cleanup
```
DIRECTORIES REMOVED:
- share/environments/ (4 shell scripts)
- share/props/ (7 YAML files with embedded shell scripts)
- share/scenes/ (entire shell script system)
- share/lighting/ (legacy lighting shell scripts)
- share/configs/ (obsolete configuration)
- share/templates/ (legacy template system)

FILES REMOVED:
- server/scenes_watcher.go (legacy .sh file monitoring)
- api/environments/apply.go (shell script execution)
- api/environments/list.go 
- api/props/instantiate.go (shell script-based props)
- api/props/list.go
- api/scenes/fork.go (shell script generation)
- api/scenes/list.go
- api/scenes/load.go  
- api/scenes/save.go
```

### API Specification Cleanup
```yaml
# REMOVED LEGACY ENDPOINTS:
- GET /scenes
- POST /scenes/{sceneId}
- POST /sessions/{sessionId}/scenes/save
- POST /scenes/{sceneId}/fork
- GET /environments
- POST /environments/{environmentId}
- GET /props
- POST /sessions/{sessionId}/props/{propId}

# RESULT: Clean 77 endpoints, all with handlers
```

### Code Reference Cleanup
```go
// REMOVED: Legacy shell script generation
func createChannelScript() {
    scriptContent := fmt.Sprintf(`#!/bin/bash
    source "/opt/hd1/lib/hd1lib.sh"
    hd1::api_call "POST" "/environments/%s"
    `, environment)
    // Entire function removed
}

// REMOVED: Legacy scenes watcher
type ScenesWatcher struct {
    watcher    *fsnotify.Watcher
    scenesPath string  // Monitored deleted directory
}

// REMOVED: Legacy props with shell scripts
// Example: wooden-chair.yaml contained:
# Lines 31-118: Embedded shell script calling hd1::create_object
```

### Build System Validation
```
BEFORE CLEANUP (Build Failures):
✗ Missing handler: api/environments/apply.go
✗ Missing handler: api/props/instantiate.go  
✗ Missing handler: api/scenes/load.go
✗ Orphaned endpoints referencing deleted files

AFTER CLEANUP (Clean Build):
✅ 77 API endpoints processed
✅ 77 routes generated  
✅ 77 handler stubs validated
✅ Auto-router generated successfully
✅ All endpoints have proper handlers
```

## Systematic Cleanup Process

### Phase 1: Directory Structure Cleanup
1. **Remove legacy shell scripts**: Delete all .sh files from share/
2. **Remove mixed YAML/shell props**: Delete props with embedded scripts
3. **Remove obsolete directories**: Clean up empty/unused directories
4. **Update share/ README**: Reflect current v5.0.0 architecture

### Phase 2: API Specification Cleanup  
1. **Identify orphaned endpoints**: Use build failures as guide
2. **Remove legacy API definitions**: Clean up api.yaml
3. **Remove unused schemas**: Delete SceneInfo and other obsolete types
4. **Update API examples**: Use current channel IDs instead of legacy environments

### Phase 3: Code Reference Cleanup
1. **Remove legacy Go handlers**: Delete files for removed endpoints
2. **Remove legacy watchers**: Delete scenes_watcher.go
3. **Clean up channel creation**: Remove shell script generation
4. **Update error handling**: Remove references to deleted systems

### Phase 4: Build Validation
1. **Test clean build**: Verify all 77 endpoints have handlers
2. **Validate API generation**: Ensure auto-router completes
3. **Test functionality**: Verify channel system works
4. **Document changes**: Update all references

## Benefits Achieved

### Security Improvements
- **No shell execution**: Eliminated all shell script attack vectors
- **Pure API surface**: Only REST endpoints, no shell access
- **Clean separation**: No mixed script/configuration files

### Performance Gains
- **Direct PlayCanvas**: No shell script interpretation overhead
- **Streamlined build**: Faster builds without legacy processing
- **Memory efficiency**: Reduced complexity and resource usage

### Maintainability Excellence  
- **Single source of truth**: api.yaml drives all routing
- **Clean codebase**: No legacy code maintenance burden
- **Clear architecture**: Pure PlayCanvas entity-component-system
- **Modern patterns**: Industry-standard game engine architecture

### Build System Quality
- **Zero warnings**: Clean compilation across all components
- **100% validation**: All endpoints verified with handlers
- **Automated verification**: Build fails if handlers missing
- **Surgical precision**: Zero regressions in supported functionality

## Impact Assessment

### Repository Optimization
- **Size reduction**: Removed legacy shell scripts and props
- **Clean structure**: Focused on current v5.0.0 architecture
- **Dependency cleanup**: Eliminated A-Frame dependencies
- **Documentation accuracy**: All references reflect current reality

### API Surface Clarity
- **77 endpoints**: All modern PlayCanvas-based functionality
- **Consistent patterns**: RESTful entity-component management
- **Complete validation**: Build system prevents incomplete APIs
- **Documentation alignment**: API spec matches implementation

### Developer Experience
- **Clear boundaries**: No confusion between legacy and current
- **Modern tooling**: Pure PlayCanvas development
- **Predictable behavior**: No hidden shell script execution
- **Quality standards**: Enterprise-grade engineering practices

## Risk Mitigation

### Validation Strategies
- **Build-driven cleanup**: Use compile failures to find all references
- **Systematic approach**: File-by-file analysis and verification  
- **Functional testing**: Verify channel system operations
- **Documentation review**: Ensure all references updated

### Rollback Protection
- **Git history**: Complete change tracking for any needed rollbacks
- **Incremental approach**: Staged removal with validation at each step
- **Backup preservation**: Legacy code preserved in git history
- **Clear versioning**: v4.0.0 vs v5.0.0 architectural boundaries

## Future Maintenance

### Prevention Strategies
- **Build validation**: Automatic detection of missing handlers
- **Documentation standards**: Keep all references current
- **Code review**: Prevent reintroduction of legacy patterns
- **Architecture principles**: Maintain single source of truth

### Evolution Path
- **PlayCanvas focus**: All future development on modern engine
- **API-first**: Every feature accessible via REST endpoints
- **Channel expansion**: Additional channel types and configurations
- **Component library**: Standardized PlayCanvas components

## Alternatives Considered

1. **Gradual legacy removal**: Rejected for complexity and confusion
2. **Legacy compatibility layer**: Rejected for maintenance burden
3. **Hybrid shell/PlayCanvas**: Rejected for security and performance
4. **Manual legacy flagging**: Rejected for incomplete identification

## References
- ADR-021: PlayCanvas Migration Implementation
- ADR-022: Three-Layer to Channel Architecture Migration
- HD1 v5.0.0 Build System Validation
- PlayCanvas Entity-Component-System architecture
- Single source of truth principles

---
**Decision made by**: HD1 Development Team  
**Implementation date**: 2025-07-03  
**Version**: v5.0.0 - API-FIRST GAME ENGINE REVOLUTION  
**Completion evidence**: Clean build with 77 validated endpoints  
**Review date**: Ongoing maintenance and architecture evolution