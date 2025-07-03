# ADR-018: API-First Game Engine Architecture

## Status
**ACCEPTED** - 2025-07-01

## Context

Following ADR-017's decision to adopt PlayCanvas as HD1's sole rendering backend, we need to define the technical architecture for creating the world's first API-driven game engine platform.

### Challenge
PlayCanvas is a JavaScript game engine library, not an HTTP API service. We must create an API layer that exposes PlayCanvas capabilities through REST endpoints while maintaining HD1's core principles.

### Technical Requirements
- Vendor PlayCanvas source locally (no CDN dependency)
- Generate API endpoints from PlayCanvas source analysis
- Maintain source synchronization capability with upstream PlayCanvas
- Preserve HD1's specification-driven development approach
- Enable complete game engine control via HTTP APIs

## Decision

**We will implement a dual-source API generation system that combines HD1 platform APIs with PlayCanvas engine APIs into a unified specification.**

### Architecture Components

#### 1. PlayCanvas Source Integration
```
/opt/hd1/
├── vendor/playcanvas/              # Source sync location (build-time)
└── share/htdocs/vendor/playcanvas/ # Web-accessible engine library
```

**Strategy:**
- Vendor PlayCanvas source locally for control and offline capability
- Copy to web-accessible location during build process
- Maintain sync capability with upstream releases
- Serve via HTTP static file serving (standard web pattern)

#### 2. Dual-Source API Generation
```
API Generation Pipeline:
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│ PlayCanvas      │    │ HD1 Platform     │    │ Unified         │
│ Source Analysis │ +  │ api.yaml         │ =  │ API Spec        │
│ (Auto-detect)   │    │ (Manual Design)  │    │ (Generated)     │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

**Sources:**
- **PlayCanvas Engine**: Auto-analyzed for entity, component, scene management APIs
- **HD1 Platform**: Existing session, object, environment management APIs
- **Unified Output**: Single `api.yaml` specification driving all code generation

#### 3. File Serving Strategy
```
HTTP Static Serving:
├── /vendor/playcanvas/playcanvas.js    # PlayCanvas engine library
├── /static/js/hd1-api.js               # Generated HD1 API client  
├── /static/js/hd1-playcanvas.js        # Generated PlayCanvas wrapper
└── /static/js/hd1-engine.js            # Unified game engine interface

WebSocket Real-time:
└── /ws → Game state synchronization, object lifecycle, session restoration
```

### API Layer Architecture

#### Game Engine API Categories

**1. Entity Management APIs**
```yaml
# Generated from PlayCanvas Entity system
POST /api/sessions/{id}/entities
GET  /api/sessions/{id}/entities
PUT  /api/sessions/{id}/entities/{entityId}
DELETE /api/sessions/{id}/entities/{entityId}
```

**2. Component System APIs**
```yaml
# Generated from PlayCanvas Component system
POST /api/sessions/{id}/entities/{entityId}/components
GET  /api/sessions/{id}/entities/{entityId}/components/{type}
PUT  /api/sessions/{id}/entities/{entityId}/components/{type}
DELETE /api/sessions/{id}/entities/{entityId}/components/{type}
```

**3. Asset Management APIs**
```yaml
# Generated from PlayCanvas Asset system
POST /api/sessions/{id}/assets
GET  /api/sessions/{id}/assets
PUT  /api/sessions/{id}/assets/{assetId}
DELETE /api/sessions/{id}/assets/{assetId}
```

**4. Scene Graph APIs**
```yaml
# Generated from PlayCanvas Scene system
POST /api/sessions/{id}/scene/hierarchy
GET  /api/sessions/{id}/scene/hierarchy
PUT  /api/sessions/{id}/scene/hierarchy/{nodeId}
```

**5. Animation System APIs**
```yaml
# Generated from PlayCanvas Animation system
POST /api/sessions/{id}/animations
GET  /api/sessions/{id}/animations
POST /api/sessions/{id}/animations/{animId}/play
POST /api/sessions/{id}/animations/{animId}/stop
```

### Implementation Strategy

#### Phase 1: Source Integration
1. **Vendor PlayCanvas**: Clone to `/opt/hd1/vendor/playcanvas/`
2. **Build Integration**: Copy to web-accessible location during `make generate`
3. **Version Management**: Track PlayCanvas version for sync capability
4. **Static Serving**: Configure HD1 daemon to serve `/vendor/playcanvas/*`

#### Phase 2: API Analysis & Generation
1. **PlayCanvas Analysis**: Parse PlayCanvas source to identify API-worthy methods
2. **Endpoint Generation**: Auto-generate REST endpoints from engine capabilities
3. **Specification Merge**: Combine with existing HD1 platform APIs
4. **Handler Generation**: Create HTTP handlers that invoke PlayCanvas methods

#### Phase 3: Client Integration
1. **JavaScript Bridge**: Generate HD1-PlayCanvas integration layer
2. **API Client**: Update generated JavaScript client with game engine endpoints
3. **WebSocket Integration**: Extend real-time synchronization for game state
4. **A-Frame Removal**: Clean removal of all A-Frame dependencies

#### Phase 4: Advanced Features
1. **Physics Integration**: API control of PlayCanvas physics engine
2. **Audio System**: 3D positional audio via REST endpoints
3. **Input Management**: Game input handling through APIs
4. **Asset Streaming**: Dynamic asset loading and management

### Code Generation Pipeline

#### Build Process Enhancement
```bash
make generate:
1. Analyze PlayCanvas source → extract API-worthy methods
2. Merge with existing api.yaml → unified specification
3. Generate HTTP handlers → PlayCanvas method bridges
4. Generate JavaScript clients → unified game engine interface
5. Copy PlayCanvas to web directory → static serving ready
```

#### Template System Extension
```
src/codegen/templates/
├── playcanvas/
│   ├── entity-handlers.tmpl     # Entity management endpoints
│   ├── component-handlers.tmpl  # Component system endpoints
│   ├── asset-handlers.tmpl      # Asset management endpoints
│   └── bridge.tmpl             # JavaScript PlayCanvas bridge
```

### Architectural Principles Preserved

#### Single Source of Truth
- **Unified api.yaml**: Combines platform + engine specifications
- **Generated Code**: All PlayCanvas integration auto-generated
- **No Manual Sync**: PlayCanvas capabilities auto-discovered and exposed

#### API-First Philosophy
- **Every Feature**: All PlayCanvas capabilities accessible via HTTP
- **Game Engine as a Service**: Professional game development via REST APIs
- **Clean Abstractions**: HTTP endpoints hide JavaScript complexity

#### HD1 Core Values
- **Specification-Driven**: PlayCanvas APIs generated from source analysis
- **Zero Regressions**: Existing platform APIs remain unchanged
- **Thread-Safe**: Game engine operations properly synchronized
- **Session-Centric**: Game instances mapped to HD1 sessions

## Expected Outcomes

### Revolutionary Capabilities
- **API-Controlled Game Engine**: Every PlayCanvas feature via HTTP endpoints
- **Game Development as a Service**: Create games through REST API calls
- **Real-time Collaboration**: Multi-user game development via WebSocket sync
- **Headless Game Server**: Game logic without mandatory visual client

### Technical Advantages
- **Professional Engine**: AAA game development capabilities via APIs
- **Modern Graphics**: WebGL2/WebGPU/WebXR through HTTP interface
- **Component Architecture**: Entity-component system via REST endpoints
- **Asset Pipeline**: Professional asset management through APIs

### Developer Experience
- **Familiar HTTP**: Game development via REST APIs instead of engine SDKs
- **Language Agnostic**: Any HTTP client can create games
- **Service Integration**: Games created by external services via API calls
- **Real-time Sync**: Collaborative development through WebSocket state sharing

## Success Metrics

- **Complete Engine Coverage**: All PlayCanvas capabilities accessible via HTTP APIs
- **Zero Manual Configuration**: All game engine integration auto-generated
- **Performance Parity**: No significant overhead from API abstraction layer
- **Developer Adoption**: External services successfully creating games via APIs
- **Market Validation**: Recognition as first API-driven game engine platform

---

**This ADR establishes the technical foundation for transforming HD1 from a visualization platform into the world's first API-driven game engine service.**