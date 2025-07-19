# ADR-012: Three.js API Expansion Strategy for Universal 3D Interface

**Date**: 2025-07-19  
**Status**: Accepted  
**Deciders**: HD1 Architecture Team  
**Technical Story**: Expand from 11 to 200+ Three.js endpoints for comprehensive 3D API coverage

## Context

HD1 currently exposes only 11 Three.js endpoints (sync, entities, avatars, scene, system) but claims to be a "universal 3D interface platform." Analysis reveals massive API gaps:

### Current Limited Coverage
- **Geometries**: Only box, sphere, cylinder (3 of 50+ available)
- **Materials**: Only basic phong materials
- **Lights**: No lighting control APIs
- **Cameras**: No camera manipulation
- **Textures**: No texture loading or manipulation
- **Animation**: No animation timeline control
- **Effects**: No post-processing or shaders
- **Physics**: No collision detection or raycasting

### Three.js API Potential
- **50+ Geometry Types**: box, sphere, cylinder, torus, text, extrude, lathe, etc.
- **Complete Material System**: basic, phong, standard, physical, shader materials
- **Full Lighting**: directional, point, spot, ambient, hemisphere lights
- **Camera Control**: perspective, orthographic, stereo cameras
- **Texture Management**: loading, creation, manipulation, procedural
- **Animation System**: keyframes, mixing, clips, timeline control
- **Post-Processing**: effects, filters, shaders, render passes
- **Physics Integration**: collision, raycasting, constraints

### Vision Gap
**Current**: Basic 3D object creation  
**Target**: Complete Three.js functionality as REST APIs  
**Goal**: "GraphQL for 3D Graphics" - every WebGL feature via HTTP

## Decision

**Implement comprehensive Three.js API expansion** from 11 to 200+ endpoints covering every Three.js class and method.

### API Architecture Strategy
```yaml
# Proposed API Structure
/api/geometries/{type}      # All 50+ geometry types
/api/materials/{type}       # All material types with properties
/api/lights/{type}          # All light types with controls
/api/cameras/{type}         # Camera creation and manipulation
/api/textures/{operation}   # Texture loading and processing
/api/animations/{operation} # Animation system control
/api/effects/{type}         # Post-processing effects
/api/physics/{operation}    # Physics and collision detection
/api/shaders/{type}         # Custom shader management
/api/rendering/{operation}  # Render controls and optimization
```

### Implementation Approach
1. **Auto-Generation from TypeScript**: Parse Three.js definitions
2. **OpenAPI Specification**: Generate comprehensive spec from TS types
3. **Handler Generation**: Auto-generate Go handlers for all endpoints
4. **Real-time Sync**: WebSocket broadcast all state changes
5. **Documentation**: Auto-generate API docs with examples

### Coverage Expansion Plan
- **Phase 1**: Complete geometry system (50+ types)
- **Phase 2**: Full material and lighting system
- **Phase 3**: Camera controls and texture management
- **Phase 4**: Animation and timeline controls
- **Phase 5**: Post-processing and effects
- **Phase 6**: Physics and advanced rendering

## Consequences

### Positive
- **Universal 3D Interface**: Every Three.js feature accessible via REST
- **Developer Experience**: Intuitive HTTP APIs for complex 3D operations
- **Language Agnostic**: Any language can create 3D interfaces
- **Microservice Friendly**: Services can add 3D UIs via HTTP calls
- **Real-time Collaboration**: WebSocket sync for multiplayer experiences
- **Complete Documentation**: Auto-generated docs for all endpoints

### Negative
- **Large API Surface**: 200+ endpoints to maintain
- **Complex State Management**: Synchronizing complex 3D state
- **Performance Considerations**: HTTP overhead for real-time operations

### Mitigation
- **Auto-Generation**: Reduces maintenance overhead
- **WebSocket Optimization**: Real-time operations use WebSocket
- **Intelligent Caching**: Client-side caching for performance
- **Progressive Disclosure**: Basic operations simple, advanced features available

## Implementation Timeline

**Target**: HD1 v0.9.0  
**Effort**: 4 weeks comprehensive implementation  
**Risk**: Medium (large scope, requires careful design)

### Development Phases
1. **Week 1**: Geometry expansion (50+ types)
2. **Week 2**: Material and lighting systems  
3. **Week 3**: Camera, texture, animation APIs
4. **Week 4**: Effects, physics, advanced features

### Success Criteria
- [ ] All Three.js geometry types accessible via REST
- [ ] Complete material system with property control
- [ ] Full lighting API (directional, point, spot, ambient)
- [ ] Camera manipulation and control
- [ ] Texture loading and processing
- [ ] Animation timeline and keyframe control
- [ ] Post-processing effects and shaders
- [ ] Physics and collision detection
- [ ] Real-time synchronization for all operations
- [ ] Comprehensive API documentation
- [ ] Zero performance regressions

## API Design Principles

### RESTful Design
```bash
# Create geometry
POST /api/geometries/torus {"radiusOuter": 2, "radiusInner": 1}

# Update material  
PUT /api/materials/standard/material123 {"metalness": 0.8, "roughness": 0.2}

# Control lighting
POST /api/lights/directional {"intensity": 1, "position": [10, 10, 10]}

# Animate properties
POST /api/animations/keyframe {"target": "mesh123", "property": "rotation.y", "keyframes": [...]}
```

### Real-time Synchronization
- All REST operations broadcast via WebSocket
- Clients receive real-time updates
- Conflict resolution for simultaneous edits
- Optimistic updates for performance

## Alignment with HD1 Principles

✅ **One Source of Truth**: api.yaml specification drives all endpoints  
✅ **No Regressions**: Existing functionality preserved and enhanced  
✅ **No Parallel Implementations**: Single REST API for each Three.js feature  
✅ **No Hacks**: Clean, systematic API expansion  
✅ **Bar Raising**: Comprehensive 3D API coverage  
✅ **Zero Compile Warnings**: Auto-generated code quality  

**Result**: HD1 becomes the definitive "Stripe for 3D Graphics" - where any service can add rich 3D interfaces with simple HTTP calls.