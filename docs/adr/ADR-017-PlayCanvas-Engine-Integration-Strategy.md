# ADR-017: PlayCanvas Engine Integration Strategy

## Status
**ACCEPTED** - 2025-07-01

## Context

HD1 has successfully implemented a complete lighting system and session restoration architecture with A-Frame as the rendering backend. However, we've reached a strategic fork in the road regarding our graphics engine architecture.

### Current State Analysis
- HD1 v2.0 has 31 HTTP API endpoints with WebSocket real-time broadcasting
- Three-layer architecture (Environment/Props/Scene) successfully implemented
- A-Frame provides WebVR-focused 3D visualization capabilities
- System is production-ready with session-centric architecture

### Strategic Options Evaluated

**Option 1: Multi-Backend Extensible Architecture**
- Support multiple rendering backends (WebGL, A-Frame, PlayCanvas)
- Abstraction layer for cross-backend compatibility
- Maximum flexibility and rendering options

**Option 2: Streamlined PlayCanvas-Only Architecture**  
- Single professional game engine backend (PlayCanvas)
- Deep integration with game engine capabilities
- API-driven game development platform

### PlayCanvas Engine Analysis

**Technical Capabilities:**
- Professional WebGL/WebGPU game engine (10.4k GitHub stars)
- Complete game development toolkit with physics, animation, audio
- Component-based entity system
- Modern web graphics optimization (WebGL2, WebGPU, WebXR)
- MIT licensed with active development

**Strategic Advantages:**
- **Revolutionary Market Position**: World's first API-driven game engine
- **Game Engine as a Service**: Every PlayCanvas feature accessible via HTTP APIs
- **Professional Capabilities**: AAA game development features vs simple 3D visualization
- **Unique Value Proposition**: No competitors in API-controlled game engine space

## Decision

**We choose Option 2: Streamlined PlayCanvas-Only Architecture**

### Rationale

1. **Focus Enables Excellence**: One exceptional backend surpasses three mediocre ones
2. **Revolutionary Positioning**: API-driven game engine is unprecedented in the market
3. **Professional Capabilities**: Game engine features vs basic 3D visualization
4. **Bar-Raising Potential**: Unique "Game Engine as a Service" model
5. **Maintenance Efficiency**: Single backend eliminates feature parity complexity

### Strategic Vision

**HD1 v3.0 Transformation:**
```
From: "3D Visualization Platform" 
To:   "API-First Game Engine Platform"
```

**Market Positioning:**
- Game Engine as a Service (GEaaS)
- HTTP API interface to professional game development
- Real-time collaborative game creation
- Headless game server architecture

## Implementation Strategy

### Phase 1: Foundation Alignment
- Extend `api.yaml` with PlayCanvas-native game engine concepts
- Align HD1 object model with PlayCanvas entity-component system
- Create asset management APIs for PlayCanvas pipeline
- Map sessions to PlayCanvas application instances

### Phase 2: Engine Integration  
- Replace A-Frame bridge with PlayCanvas integration
- Implement component system mapping (objects → entities)
- Integrate environment system with PlayCanvas physics contexts
- Align lighting/materials with PBR rendering pipeline

### Phase 3: Game Engine APIs
- Animation system API endpoints
- Input management APIs  
- 3D positional audio integration
- Dynamic asset loading and streaming

### Zero-Regression Migration Strategy
1. **Parallel Development**: Build PlayCanvas bridge alongside existing A-Frame
2. **Feature Parity Validation**: Ensure all current capabilities translate
3. **Atomic Cutover**: Single switch from A-Frame to PlayCanvas backend
4. **Comprehensive Testing**: Validate all 31 API endpoints post-migration

## Architectural Principles Maintained

### One Source of Truth
- `api.yaml` remains single specification source
- PlayCanvas becomes sole rendering backend
- Zero parallel implementations
- Template system generates all integration code

### API-First Philosophy  
- Every PlayCanvas feature accessible via HTTP endpoints
- Game engine capabilities exposed through clean REST APIs
- WebSocket real-time synchronization for game state
- Session-centric architecture preserved

### HD1 Core Values
- Specification-driven development continues
- Thread-safe session management maintained  
- Clean separation: HTTP commands, WebSocket events
- Universal coordinate system validation preserved

## Expected Outcomes

### Revolutionary Capabilities
- **Industry First**: API-controlled professional game engine
- **Service Integration**: Any system can create games via HTTP calls
- **Collaborative Development**: Real-time multi-user game creation
- **Headless Architecture**: Game server without mandatory visual client

### Technical Advantages
- Professional game engine features (physics, animation, audio)
- Modern web graphics optimization (WebGL2, WebGPU)
- Component-based entity system
- Advanced asset management and streaming

### Market Differentiation
- **Unique Position**: No competitors in API-driven game engine space
- **Professional Grade**: AAA game development via REST APIs
- **Platform Approach**: Game Engine as a Service model
- **Developer Experience**: HTTP APIs vs complex game engine SDKs

## Migration Timeline

- **Documentation Phase**: ADRs, architecture updates, execution tracking
- **Foundation Phase**: API specification extensions, entity model alignment  
- **Integration Phase**: PlayCanvas bridge development, feature mapping
- **Cutover Phase**: Atomic backend replacement with validation
- **Enhancement Phase**: Advanced game engine features via APIs

## Success Metrics

- Zero functionality regression from A-Frame → PlayCanvas migration
- All 31 current API endpoints maintain compatibility
- PlayCanvas game engine features accessible via new API endpoints
- Performance improvements in rendering and 3D capabilities
- Successful transformation to "Game Engine as a Service" platform

---

**This ADR represents the most significant architectural evolution in HD1's history - from visualization platform to revolutionary API-driven game engine.**