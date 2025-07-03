# ADR-001: Specification-Driven Development Architecture

## Status
✅ **ACCEPTED** - Implemented and Operational

## Context

Traditional API development suffers from inconsistency between API documentation and implementation. Manual routing configuration leads to errors, regressions, and maintenance overhead. For HD1 (Virtual World Synthesizer), we needed a advanced approach that ensures the API specification IS the system architecture.

## Decision

**We adopt 100% specification-driven development where the OpenAPI 3.0.3 specification (`api.yaml`) serves as the single source of truth for all routing and endpoint configuration.**

### Core Architecture
```
api.yaml → codegen/generator.go → auto_router.go → api/ handlers → Virtual World
```

### Key Components
1. **Single Source of Truth**: `api.yaml` defines all endpoints, routing, and handler mappings
2. **Automatic Code Generation**: `codegen/generator.go` reads specification and generates `auto_router.go`
3. **Build-time Validation**: Missing handlers cause build failures
4. **Zero Manual Routing**: No hand-written route configurations

## Implementation Details

### OpenAPI Extensions
```yaml
paths:
  /sessions:
    post:
      operationId: createSession
      x-handler: "api/sessions/create.go"    # Handler file location
      x-function: "CreateSession"            # Function name
```

### Generated Router Pattern
```go
// AUTO-GENERATED FROM api.yaml - DO NOT EDIT MANUALLY
func (r *APIRouter) generateRoutes() {
    r.routes = []Route{
        {
            Path:       "/sessions",
            Method:     "POST",
            Handler:    r.CreateSession,
            OperationID: "createSession",
        },
    }
}
```

### Handler Validation
- Build system validates that all `x-handler` files exist
- Missing handlers cause immediate build failure
- Prevents deployment of incomplete implementations

## Consequences

### ✅ Benefits
- **Zero Routing Bugs**: Impossible to have route/handler mismatches
- **Self-Documenting**: API specification IS the implementation
- **No Regressions**: Missing handlers fail build, preventing broken deployments
- **Rapid Development**: Add endpoint to spec, implement handler, done
- **Perfect Consistency**: Documentation and implementation always match

### ⚠️ Trade-offs
- **Learning Curve**: Developers must understand OpenAPI specification
- **Build Dependency**: Changes require regeneration step
- **Tool Complexity**: Code generator must be maintained

### 🔧 Mitigation
- Comprehensive documentation and examples
- Automated build pipeline handles regeneration
- Simple, focused code generator with clear error messages

## Validation

### Success Metrics
✅ **Zero Manual Routes**: No hand-written routing code  
✅ **Build Validation**: Missing handlers fail build  
✅ **Documentation Sync**: API docs always match implementation  
✅ **Rapid Development**: New endpoints added in minutes  

### Real-world Validation
- Successfully implemented 13 API endpoints with zero routing bugs
- Build system catches missing handlers before deployment
- Team velocity increased by eliminating manual routing tasks

## Examples

### Adding New Endpoint
1. **Define in Specification**:
```yaml
/sessions/{sessionId}/physics:
  post:
    operationId: enablePhysics
    x-handler: "api/physics/enable.go"
    x-function: "EnablePhysics"
```

2. **Auto-Generate Routing**:
```bash
make generate  # Reads api.yaml, updates auto_router.go
```

3. **Implement Handler**:
```go
// api/physics/enable.go
func EnablePhysicsHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
    // Implementation
}
```

4. **Build & Deploy**:
```bash
make build  # Validates handler exists, builds binary
```

## Related Decisions
- [ADR-002: Thread-Safe SessionStore](002-thread-safe-session-store.md)
- [ADR-003: WebSocket Real-time Architecture](003-websocket-realtime-architecture.md)

## References
- [OpenAPI 3.0.3 Specification](https://spec.openapis.org/oas/v3.0.3)
- [HD1 API Specification](../../src/api.yaml)
- [Code Generator Implementation](../../src/codegen/generator.go)

---

**Decision Date**: 2025-06-28  
**Decision Makers**: HD1 Architecture Team  
**Review Date**: 2025-12-28  

*This ADR represents the foundational architectural decision that makes HD1 possible.*