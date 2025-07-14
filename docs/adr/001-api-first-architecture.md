# ADR-001: API-First Architecture

**Status**: Accepted  
**Date**: 2025-07-14  
**Authors**: HD1 Development Team  

## Context

HD1 requires a scalable, maintainable architecture that can support diverse client types (web, mobile, desktop) while maintaining consistency and reducing development overhead. Traditional approaches lead to API drift, inconsistent documentation, and manual synchronization between server and client implementations.

## Decision

We adopt an **API-First Architecture** where the OpenAPI specification (`src/api.yaml`) serves as the single source of truth for all system behavior.

### Core Principles

1. **Specification-Driven Development**: All routing, validation, and client libraries generated from `api.yaml`
2. **Single Source of Truth**: No manual route configuration or parallel implementations
3. **Auto-Generation**: Complete toolchain for routing, handlers, clients, and documentation
4. **Contract-First**: API contracts defined before implementation begins

### Implementation

```yaml
# src/api.yaml defines complete system behavior
paths:
  /threejs/entities:
    post:
      operationId: createEntity
      x-handler: api/threejs/create.go
      x-function: CreateEntity
```

Generated artifacts:
- **Go Router**: `auto_router.go` with complete routing table
- **JavaScript Client**: `hd1lib.js` with typed API methods
- **Validation**: Automatic request/response validation
- **Documentation**: Live API documentation

## Consequences

### Positive
- **Zero API Drift**: Specification and implementation always synchronized
- **Consistent Clients**: All client libraries generated from same specification
- **Reduced Boilerplate**: Automatic generation eliminates manual routing code
- **Type Safety**: Compile-time validation of API contracts
- **Documentation**: Always up-to-date, executable documentation

### Negative
- **Build Complexity**: Code generation adds build step dependency
- **Learning Curve**: Developers must understand OpenAPI specification format
- **Template Maintenance**: Code generation templates require periodic updates
- **Debug Complexity**: Generated code can be harder to debug

### Neutral
- **Tool Dependency**: Requires OpenAPI-compatible tooling throughout development
- **Version Control**: Generated files included in repository for transparency

## Implementation Details

### Directory Structure
```
src/
├── api.yaml           # Single source of truth
├── auto_router.go     # Generated (never edit)
├── codegen/           # Generation templates and logic
└── api/               # Handler implementations
```

### Build Integration
```bash
# Automatic regeneration on api.yaml changes
make generate  # Regenerate all code from specification
make build     # Build with generated code
```

### Validation Strategy
- **Build-time**: Specification syntax validation
- **Runtime**: Automatic request/response validation
- **Test-time**: Contract testing against specification

## Alternatives Considered

### Manual Routing
**Rejected**: High maintenance overhead, prone to inconsistencies, difficult to maintain documentation.

### Framework-Based (Gin, Echo)
**Rejected**: Adds dependency overhead, reduces control over generated code, limited OpenAPI integration.

### GraphQL
**Rejected**: Adds complexity for real-time use cases, WebSocket integration challenges, learning curve for team.

## Monitoring and Success Metrics

- **API Consistency**: Zero drift between specification and implementation
- **Development Velocity**: Reduced time to add new endpoints
- **Client Adoption**: Multiple client types using same API
- **Documentation Quality**: Up-to-date, executable API documentation

## Future Considerations

- **Versioning**: API versioning strategy for backward compatibility
- **Performance**: Optimize generated code for high-throughput scenarios
- **Extensions**: Custom OpenAPI extensions for HD1-specific requirements
- **Testing**: Automated contract testing against live endpoints

---

*ADR-001 establishes the foundation for HD1's specification-driven architecture*