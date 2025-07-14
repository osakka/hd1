# ADR-002: Specification-Driven Development

**Status**: Accepted  
**Date**: 2025-07-14  
**Authors**: HD1 Development Team  
**Supersedes**: None  
**Related**: ADR-001 (API-First Architecture)

## Context

Building on ADR-001's API-First Architecture, we need to establish a development workflow that maintains the specification as the authoritative source while enabling rapid development and deployment cycles. Manual synchronization between specification and implementation creates maintenance overhead and introduces inconsistencies.

## Decision

We implement **Specification-Driven Development** where all code generation, validation, and deployment processes are automated from the OpenAPI specification.

### Development Workflow

1. **Specification First**: All changes begin with `api.yaml` modifications
2. **Automatic Generation**: Build system regenerates all derived code
3. **Implementation**: Developers implement only business logic handlers
4. **Validation**: Automatic testing against specification contracts

### Code Generation Strategy

```go
// Automatic generation from api.yaml
func generateFromSpecification() {
    // Router generation
    generateRouting()
    
    // Client library generation  
    generateJavaScriptClient()
    
    // Validation middleware
    generateValidation()
    
    // Documentation
    generateAPIDocs()
}
```

## Implementation

### Single Source of Truth
```yaml
# src/api.yaml - Complete system definition
openapi: 3.0.0
info:
  title: "HD1 (Holodeck One) Three.js API"
  version: "6.0.0"
x-code-generation:
  strict-validation: true
  auto-routing: true
  fail-on-missing-handlers: true
```

### Build System Integration
```makefile
# Makefile automation
generate:
    @echo "GENERATING THREE.JS CODE FROM SPECIFICATION..."
    go run codegen/generator.go

build: generate
    go build -o ../build/bin/hd1 .
```

### Handler Template System
```go
// api/threejs/create.go - Developer implements business logic only
func CreateEntity(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
    // Business logic here - routing/validation auto-generated
}
```

## Consequences

### Positive
- **Consistency**: Specification and implementation never drift apart
- **Velocity**: Developers focus on business logic, not boilerplate
- **Quality**: Automatic validation prevents common errors
- **Documentation**: Always accurate, executable documentation
- **Client Support**: Multiple client libraries from single specification

### Negative
- **Complexity**: Build process more complex than manual approaches
- **Dependencies**: Requires functioning code generation at all times
- **Learning Curve**: Team must understand OpenAPI specifications deeply
- **Debugging**: Generated code adds layer between specification and runtime

### Neutral
- **Tool Investment**: Requires building and maintaining generation toolchain
- **Process Change**: Different development workflow than traditional approaches

## Technical Implementation

### Code Generator Architecture
```
api.yaml → Parser → Validation → Templates → Generated Code
```

### Generated Artifacts
- **`auto_router.go`**: Complete HTTP routing with parameter extraction
- **`hd1lib.js`**: JavaScript client library with typed methods  
- **Validation**: Request/response schema validation middleware
- **Documentation**: Interactive API explorer

### Template System
```go
// External templates for maintainable generation
templateFS embed.FS
templates/
├── go/router.tmpl           # Go HTTP routing
├── javascript/client.tmpl   # JavaScript API client
└── validation/schemas.tmpl  # JSON schema validation
```

## Validation Strategy

### Build-Time Validation
- **Specification Syntax**: OpenAPI 3.0 compliance
- **Handler Existence**: All endpoints have corresponding handler files
- **Template Validity**: Code generation templates compile successfully

### Runtime Validation  
- **Request Validation**: Automatic schema validation for all requests
- **Response Validation**: Ensure responses match specification schemas
- **Type Safety**: Generated code provides compile-time type checking

### Development Validation
- **Contract Testing**: Automated tests verify implementation matches specification
- **Live Documentation**: Interactive API explorer for manual testing
- **Client Validation**: Generated clients tested against live API

## Development Process

### Adding New Endpoints
1. **Define in Specification**: Add endpoint to `api.yaml`
2. **Generate Code**: Run `make generate` to create routing
3. **Implement Handler**: Write business logic in designated handler file
4. **Test**: Validate against specification contracts
5. **Deploy**: Generated client libraries automatically available

### Modifying Existing Endpoints
1. **Update Specification**: Modify endpoint definition in `api.yaml`
2. **Regenerate**: Build system updates all derived code
3. **Update Handler**: Modify business logic if needed
4. **Validate**: Ensure backward compatibility if required
5. **Deploy**: Updated clients automatically generated

## Success Metrics

- **Development Velocity**: Time to add new endpoints (target: <30 minutes)
- **API Consistency**: Zero drift between specification and implementation
- **Code Coverage**: Generated code coverage (target: >95%)
- **Client Adoption**: Number of different client implementations using API

## Future Enhancements

### Versioning Strategy
- **Specification Versioning**: Semantic versioning for API contracts
- **Backward Compatibility**: Automated compatibility checking
- **Migration Tools**: Generated migration guides for breaking changes

### Advanced Generation
- **Performance Optimization**: Generated code optimized for high throughput
- **Custom Extensions**: HD1-specific OpenAPI extensions
- **Integration Testing**: Generated integration test suites

### Tooling Improvements
- **IDE Integration**: Real-time specification validation in editors
- **Visual Design**: Graphical API design tools
- **Monitoring**: Runtime monitoring of specification compliance

---

*ADR-002 establishes HD1's specification-driven development workflow*