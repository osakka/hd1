# ADR-016: Tied API Architecture - Unified Platform and Downstream APIs

## Status
✅ **ACCEPTED** - Core Architecture Pattern (2025-07-01)

## Context

HD1 operates as a comprehensive holodeck platform that must expose both high-level framework abstractions (sessions, props, environments) and low-level capabilities (A-Frame, WebXR, physics). Traditional API architectures create silos between different capability layers, leading to inconsistent interfaces and complex client implementations.

## Decision

**We implement a Tied API Architecture where HD1 exposes everything through a unified API surface, with scripts serving as bidirectional bridges between API layers.**

### Core Architectural Principles

#### 1. **100% API-Driven Service**
- **All functionality exposed via API endpoints** - No direct client manipulation of underlying systems
- **Scripts are both API consumers AND API targets** - Called by API, call API internally
- **Single interface for all capabilities** - Clients only interact through API endpoints

#### 2. **Two-Layer API Exposure**
- **HD1 Platform APIs**: Sessions, props, environments, scenes (framework abstractions)
- **Downstream APIs**: Direct exposure of A-Frame, WebXR, physics capabilities
- **Unified client interface**: All capabilities accessible through single API surface

#### 3. **State Management Separation**
- **API handles all commands and mutations** - Every state change flows through API
- **WebSocket provides state synchronization only** - Real-time updates, no command execution
- **Shared memory consistency** - All clients maintain synchronized state

## Implementation Architecture

### **Tied API Flow Pattern**
```
Client Request → HD1 API → Script Execution → Downstream APIs → State Change → WebSocket Broadcast → All Clients
```

### **Bidirectional Script Integration**
```
API Triggers Scripts:
POST /scenes/{sceneId} → executes scene script
POST /props/{propId} → executes prop script  
POST /environments/{environmentId} → executes environment script

Scripts Call API:
hd1::create_object() → POST /sessions/{id}/objects
hd1::camera() → PUT /sessions/{id}/camera/position
```

### **Client Interaction Model**
- **API-Only Interface**: Clients never directly manipulate A-Frame/WebXR
- **Automatic Script Retrieval**: Scripts delivered and executed via API triggers
- **State Synchronization**: WebSocket keeps all clients synchronized with server state
- **Framework Abstraction**: HD1 APIs provide high-level constructs (sessions, props) built on downstream capabilities

## Technical Benefits

### **1. Unified Interface Consistency**
- Single API surface for all capabilities
- Consistent authentication, validation, error handling
- No client-side complexity for different API layers

### **2. Script-Mediated Flexibility**
- Scripts bridge between HD1 framework and downstream APIs
- Hot-swappable script implementations
- API-driven script execution with automatic client synchronization

### **3. State Management Clarity**
- Clear separation: API = commands, WebSocket = state sync
- Predictable state flow and consistency guarantees
- Real-time updates without command/state confusion

### **4. Development Efficiency**
- API-first development for all features
- Scripts provide reusable API orchestration
- Single source of truth through specification-driven development

## Examples

### **Scene Creation Flow**
```bash
# Client calls HD1 API
curl -X POST /api/scenes/basic-shapes

# API executes scene script
/opt/hd1/share/scenes/basic-shapes.sh

# Script calls HD1 APIs internally
hd1::create_object "cube1" "cube" 0 1 0
hd1::create_object "sphere1" "sphere" 2 1 0

# API processes object creation
POST /sessions/{id}/objects (called by script)

# State changes broadcast via WebSocket
{"type": "object_created", "name": "cube1", ...}
```

### **Prop Instantiation Flow**
```bash
# Client calls HD1 API
curl -X POST /api/sessions/{id}/props/lightbulb

# API executes prop script
/opt/hd1/share/props/electronic/lightbulb.yaml (script section)

# Script calls downstream APIs
hd1::create_object() → A-Frame entity creation
hd1::lighting() → A-Frame light component

# All clients receive state updates via WebSocket
```

## Architectural Consequences

### ✅ **Benefits**
- **Unified Developer Experience**: Single API interface for all capabilities
- **Clear Separation of Concerns**: Commands via API, state sync via WebSocket
- **Maximum Flexibility**: Scripts provide customizable API orchestration
- **State Consistency**: Guaranteed synchronization across all clients
- **Framework Abstraction**: High-level HD1 concepts built on solid downstream foundations

### ⚠️ **Considerations**
- **Script Execution Overhead**: API calls trigger script execution
- **Complexity in Script Design**: Scripts must handle API interactions correctly
- **Error Propagation**: Script failures must be handled gracefully by API layer

### 🔧 **Mitigations**
- **Performance**: Script caching and optimized execution paths
- **Error Handling**: Comprehensive script validation and error reporting
- **Documentation**: Clear patterns and examples for script development

## Validation Criteria

### **Success Metrics**
✅ **100% API Coverage**: All HD1 capabilities accessible via API  
✅ **Script Integration**: All props/scenes/environments execute via API  
✅ **State Consistency**: WebSocket maintains synchronized client state  
✅ **Downstream Exposure**: A-Frame/WebXR capabilities available through API  
✅ **Developer Experience**: Single interface for all client interactions  

### **Implementation Evidence**
- 31 API endpoints covering all HD1 functionality
- Props, scenes, environments all executable via API calls
- WebSocket real-time state synchronization operational
- Shell functions demonstrate API-driven script patterns
- A-Frame integration accessible through API layer

## Related ADRs

- [ADR-002: Specification-Driven Development](ADR-002-Specification-Driven-Development.md)
- [ADR-007: Upstream/Downstream API Integration](ADR-007-Revolutionary-Upstream-Downstream-Integration.md)
- [ADR-009: WebSocket Real-time Architecture](ADR-009-WebSocket-Realtime-Architecture.md)
- [ADR-014: Three-Layer Architecture](ADR-014-Three-Layer-Architecture-Environment-Props-System.md)

## Conclusion

The Tied API Architecture provides HD1 with a unified, consistent interface that exposes both high-level framework abstractions and low-level capabilities through a single API surface. Scripts serve as intelligent bridges between API layers, while WebSocket provides real-time state synchronization. This architecture ensures consistent developer experience, predictable state management, and maximum flexibility for complex holodeck operations.

**Status**: Core architectural pattern - implemented and operational across all HD1 systems.

---

**Decision Date**: 2025-07-01  
**Decision Makers**: HD1 Architecture Team  
**Review Date**: 2025-12-31  

*This ADR documents the foundational architectural pattern that makes HD1's unified platform capabilities possible.*