# HD1 v3.0 Documentation

**API-First Game Engine Platform** - Complete documentation suite for the world's first HTTP-controlled professional game engine.

## 📚 Documentation Index

### 🎯 **Quick Start**
- **[Quick Reference Guide](HD1-v3-Quick-Reference.md)** - Essential commands, API endpoints, and troubleshooting

### 🏗️ **Architecture Documentation**
- **[Current State Architecture](HD1-v3-Current-State-Architecture.md)** - Comprehensive system overview and technical details
- **[ADR-019: Production-Ready API-First Game Engine](adr/ADR-019-Production-Ready-API-First-Game-Engine.md)** - Complete architectural decisions and implementation details

### 📋 **Architecture Decision Records (ADRs)**
- **[ADR-017: PlayCanvas Engine Integration Strategy](adr/ADR-017-PlayCanvas-Engine-Integration-Strategy.md)** - PlayCanvas integration approach
- **[ADR-018: API-First Game Engine Architecture](adr/ADR-018-API-First-Game-Engine-Architecture.md)** - API-first design principles
- **[ADR-019: Production-Ready API-First Game Engine](adr/ADR-019-Production-Ready-API-First-Game-Engine.md)** - Current production architecture

### 🎮 **PlayCanvas Integration**
- **[PlayCanvas Migration Tracker](HD1-v3-PlayCanvas-Migration-Tracker.md)** - Migration progress and status
- **[PlayCanvas API Design](HD1-v3-PlayCanvas-API-Design.md)** - API integration patterns

### 🔧 **API Design Documentation**
- **[Complete API Interfaces](HD1-v3-Complete-API-Interfaces.md)** - Comprehensive API documentation
- **[OpenAPI Specification](../src/api.yaml)** - Single source of truth (85 endpoints)

## 🎯 **HD1 v3.0 Overview**

### **What is HD1?**
HD1 (Holodeck One) v3.0 is the world's first **"Game Engine as a Service"** platform - a production-ready, API-first game engine that provides professional game development capabilities through HTTP APIs with real-time WebSocket synchronization.

### **Core Innovation**
**API-First Game Engine Architecture**: "API = Control, WebSocket = Graph Extension"
```
HTTP APIs → Game Commands → Server State → WebSocket Events → Client Rendering
```

### **Key Features**
- ✅ **85 REST Endpoints**: Complete game engine control via HTTP
- ✅ **Real-Time WebSocket**: <10ms entity lifecycle synchronization
- ✅ **PlayCanvas Engine**: Professional 3D rendering with ECS
- ✅ **Channel System**: YAML-based scene configuration
- ✅ **Auto-Generation**: 100% specification-driven development
- ✅ **Production Ready**: Enterprise logging, monitoring, scalability

## 🚀 **Quick Start**

### **1. Start the Server**
```bash
cd /opt/hd1/src
make clean && make && make start
```

### **2. Access the Interface**
- **Web UI**: http://localhost:8080
- **Console**: Click header to expand/collapse (smooth animations!)
- **API Documentation**: Generated from OpenAPI specification

### **3. Create Your First Entity**
```bash
# Create session
curl -X POST http://localhost:8080/api/sessions

# Create entity
curl -X POST http://localhost:8080/api/sessions/{sessionId}/entities \
  -H "Content-Type: application/json" \
  -d '{"name": "my_box", "components": {
    "model": {"type": "box"},
    "transform": {"position": [0, 1, 0]},
    "material": {"diffuse": "#ff0000"}
  }}'
```

## 🎮 **Game Engine Capabilities**

### **Entity Component System**
Complete ECS implementation accessible via REST APIs:
- **Entities**: Game objects with unique IDs
- **Components**: Transform, Model, Material, Physics, Light, Audio
- **Systems**: Physics, Rendering, Animation (handled by PlayCanvas)

### **Channel-Based Scenes**
YAML-configured scenes with hot-swapping:
```yaml
# /opt/hd1/share/channels/channel_one.yaml
channel:
  id: "channel_one"
  name: "Scene 1 - Red Box"

playcanvas:
  entities:
    - name: "floor"
      components:
        transform: {scale: [10, 1, 10]}  # 10×1×10 floor
        material: {diffuse: "#cccccc"}
```

### **Real-Time Synchronization**
WebSocket events for instant updates:
```javascript
{
  "type": "entity_created",
  "data": {
    "entity_id": "entity-xyz789",
    "components": {...}
  }
}
```

## 🏗️ **Architecture Highlights**

### **Single Source of Truth**
Everything auto-generated from `src/api.yaml`:
- Go router with 85 endpoints
- JavaScript API client
- Shell function library
- UI components and forms

### **Professional Console**
Real-time monitoring with smooth animations:
- Performance graphs (CPU, memory)
- WebSocket connection status
- Entity count tracking
- 0.4s cubic-bezier transitions

### **Template Architecture**
External templates for maintainable code generation:
```
src/codegen/templates/
├── go/router.tmpl              # Auto-router generation
├── javascript/api-client.tmpl  # JS API wrapper
└── shell/core-functions.tmpl   # Shell API client
```

## 📊 **Performance Characteristics**

- **API Response Time**: <50ms for entity operations
- **WebSocket Latency**: <10ms for real-time synchronization
- **Entity Capacity**: 1000+ entities per session at 60fps
- **Concurrent Clients**: 100+ per channel
- **API Throughput**: 1000+ requests/second

## 🔧 **Development Workflow**

### **Specification-Driven Development**
```bash
# 1. Edit API specification
vim src/api.yaml

# 2. Auto-generate all code
make generate

# 3. Build and deploy
make clean && make && make start
```

### **Quality Standards**
- **Zero Regressions**: All changes maintain backward compatibility
- **API-First**: No functionality exists outside the specification
- **Clean Architecture**: Separation of concerns maintained
- **Comprehensive Testing**: API endpoints and WebSocket flows

## 🎯 **Use Cases**

### **Game Development**
- Rapid prototyping with HTTP APIs
- Real-time multiplayer games
- Educational game development
- VR/AR experiences

### **Enterprise Applications**
- 3D data visualization
- Training simulations
- Collaborative 3D environments
- IoT device visualization

### **Research & Education**
- Game engine architecture study
- API-first development patterns
- Real-time systems research
- Web-based 3D applications

## 📖 **Documentation Guide**

### **For Developers**
1. Start with **[Quick Reference](HD1-v3-Quick-Reference.md)** for essential commands
2. Read **[Current State Architecture](HD1-v3-Current-State-Architecture.md)** for system overview
3. Explore **[ADR-019](adr/ADR-019-Production-Ready-API-First-Game-Engine.md)** for technical details

### **For Architects**
1. Review **[ADR-017](adr/ADR-017-PlayCanvas-Engine-Integration-Strategy.md)** for integration strategy
2. Study **[ADR-018](adr/ADR-018-API-First-Game-Engine-Architecture.md)** for design principles
3. Examine **[Complete API Interfaces](HD1-v3-Complete-API-Interfaces.md)** for full API coverage

### **For Users**
1. Follow **[Quick Reference](HD1-v3-Quick-Reference.md)** for getting started
2. Check **[PlayCanvas Migration Tracker](HD1-v3-PlayCanvas-Migration-Tracker.md)** for current status
3. Use **[OpenAPI Specification](../src/api.yaml)** for API reference

## 🌟 **Strategic Vision**

HD1 v3.0 establishes the foundation for **Game Engine as a Service** - democratizing professional game development through standard web technologies and RESTful APIs. By making game engine capabilities accessible via HTTP, we're revolutionizing how games are developed, deployed, and scaled.

### **Future Roadmap**
- Advanced physics systems (constraints, joints)
- Asset pipeline and management
- Multi-server clustering
- AI/ML integration for procedural content
- Industry ecosystem integration (Unity/Unreal import/export)

---

**HD1 v3.0 Documentation** - Comprehensive guide to the API-first game engine revolution.

*Where OpenAPI specifications become immersive game worlds through specification-driven engineering.*