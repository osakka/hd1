```
    ██╗  ██╗██████╗  ██╗
    ██║  ██║██╔══██╗███║
    ███████║██║  ██║╚██║
    ██╔══██║██║  ██║ ██║
    ██║  ██║██████╔╝ ██║
    ╚═╝  ╚═╝╚═════╝  ╚═╝
```

# HD1 (Holodeck One)
## The API-First 3D Engine for Distributed Services

> **⚠️ DEVELOPMENT & EXPERIMENTAL SOFTWARE**  
> **HD1 is active development code and serves as an experiment in the efficacy of [XVC (Extreme Vibe Coding)](https://github.com/osakka/xvc) - a methodology for human-LLM collaboration that treats LLMs as "pattern reflection engines" for systematic development. This is NOT production-ready software. Use for development, research, and experimentation only.**

---

> **Where REST APIs meet immersive 3D worlds**  
> Holodeck One transforms any service into a visual, interactive experience through 59 REST endpoints and real-time synchronization.

**🎯 XVC Experiment**: HD1 demonstrates XVC principles through **Pattern Consistency** (API-first specification-driven development), **Surgical Precision** (template-based code generation), and **Bar-Raising Solutions** (zero regressions, single source of truth). The project validates XVC's effectiveness for complex system architecture while building toward the vision of any distributed service integrating visual representation through HD1's unified API surface.

### 🌟 What Makes HD1 Unique

**API-First Architecture** → Every 3D operation accessible via REST (XVC Pattern Consistency)  
**Real-Time Multi-User** → <10ms WebSocket synchronization across clients  
**Service Integration Ready** → Built for distributed services to present, animate, and interact  
**WebGL/3D Native** → Professional PlayCanvas rendering engine  
**Specification-Driven** → 100% auto-generated from OpenAPI specification (XVC Single Source of Truth)

## 🚀 Current State: HD1 v5.0.1 (In Development)

HD1 has a **complete API surface** for 3D engine control with **59 REST endpoints**, **advanced multiplayer avatar synchronization**, and **complete configuration management standardization**. The foundation is built for the distributed services revolution.

### 🎮 Complete 3D Engine via REST APIs
```
59 REST Endpoints → Entity-Component-System → PlayCanvas Rendering → Real-Time Avatar Sync
```

**Core Engine APIs**
- **Entity Lifecycle**: Create, update, delete with full component systems
- **Avatar Management**: Real-time multiplayer avatar tracking and synchronization
- **Advanced Camera**: Free/orbital modes, smooth movement with momentum, TAB toggle
- **Physics Simulation**: Rigidbodies, forces, collision detection
- **Animation Control**: Timeline-based animations with precise control
- **3D Audio Engine**: Spatial audio sources with positional audio
- **Scene Hierarchy**: Parent-child relationships and transform inheritance

### 🌐 Multi-User Collaborative Environment
- **Channel Architecture**: YAML-based scene configuration
- **Real-Time Sync**: <10ms WebSocket state synchronization
- **Avatar Persistence**: High-frequency position updates (100+ per second) without entity deletion
- **Multiplayer Camera**: Orbital mode for multi-session viewing with automatic centering
- **Session Isolation**: Per-user state with clean restoration
- **Concurrent Users**: Multiple operators per channel with conflict resolution

### 🔗 Service Integration Platform
> **The Holodeck Vision**: Any service can **present**, **animate**, **demonstrate**, **write**, **talk**, and **play** through HD1's visual interface.

**Coming Integration Points**:
- **Authentication Layer**: Operator profiles and secure service access
- **Service Keys**: Distributed service registration and authorization  
- **XVC Integration**: Extended Version Control for collaborative development
- **API Gateways**: Service discovery and visual service representation

**Use Cases**:
- **DevOps Dashboards**: Visualize infrastructure as interactive 3D environments
- **Data Analytics**: Transform datasets into explorable 3D visualizations  
- **Service Monitoring**: Real-time system health in immersive interfaces
- **Collaborative Tools**: Multi-user spaces for distributed team interaction
- **Educational Platforms**: Interactive learning environments with 3D content

## ⚡ Quick Start

### Get HD1 Running in 30 Seconds
```bash
cd src && make clean && make && make start
```

**🎯 Access Points**:
- **Visual Console**: http://localhost:8080 (Professional monitoring interface)
- **API Explorer**: Interactive testing of all 59 endpoints
- **WebSocket Feed**: Real-time entity updates and state synchronization

### Your First 3D Scene via API
```bash
# Create a session
SESSION_ID=$(./build/bin/hd1-client create-session | jq -r '.session_id')

# Create a spinning cube
./build/bin/hd1-client create-entity "$SESSION_ID" '{
  "name": "hello-cube",
  "components": {
    "transform": {"position": {"x": 0, "y": 1, "z": 0}},
    "render": {"geometry": "box", "material": {"color": "#00ffff"}}
  }
}'

# Add physics and animation
./build/bin/hd1-client create-animation "$SESSION_ID" '{
  "name": "spin",
  "target": "hello-cube",
  "duration": 3.0,
  "properties": {"rotation": {"y": 360}},
  "loop": true
}'
```

**Result**: A teal cube spinning in 3D space, controlled entirely via REST APIs.

## 🛠️ Architecture & Technical Excellence

### API-First Development Flow
```
OpenAPI Specification (api.yaml) → Templates → Generated Code → Running System
```

**Core Principle**: Every feature starts as an API specification, ensuring consistency and maintainability.

### Template-Driven Code Generation
```
src/codegen/templates/
├── go/router.tmpl              # Auto-router from specification  
├── javascript/api-client.tmpl  # JS API wrapper
├── javascript/playcanvas-bridge.tmpl # 3D engine integration
└── shell/api-functions.tmpl    # Command-line interface
```

**8 External Templates** → **100% Generated Code** → **Zero Manual Synchronization**

### System Performance
- **API Response**: Efficient response handling for all 59 endpoints
- **WebSocket**: Real-time synchronization  
- **Entity Operations**: Efficient create/update/delete operations
- **Memory**: Scales with entity count
- **Concurrent Users**: Multi-user collaborative support

## 🌌 The Holodeck Vision

### Inspired by Star Trek's Immersive Environments
HD1 draws inspiration from Star Trek's **Holodeck** - a space where any scenario can be created, experienced, and interacted with through technology. In our case, **distributed services** become the scenarios.

### From Concept to Reality
**Today**: HD1 provides 59 REST endpoints for 3D engine control (in development)  
**Tomorrow**: Any distributed service integrates visual representation through HD1's API surface  
**Future**: The **visual interface layer** for the entire distributed service ecosystem

### Logo Design

**Design Concept**: 
- **HD Base**: Concrete foundation, compact and cracked under computational weight
- **Wireframe "1"**: Teal/blue wireframe floating above, representing the first platform

*Symbolizing the solid foundation (HD) supporting the innovative interface layer (1)*

---

## 📖 Documentation & Development

### Complete Documentation Suite
- **`/docs/README.md`**: Master documentation index with professional taxonomy
- **`/docs/decisions/`**: 25 Architectural Decision Records (ADRs) 
- **`/src/README.md`**: Developer-focused with API → Templates → Code workflow
- **`/src/codegen/templates/README.md`**: Complete template architecture guide

### Development Commands
```bash
# Clean build and start
cd src && make clean && make && make start

# Code generation from specification  
make generate

# Development tools
make test        # Run test suite
make stop        # Stop daemon
```

### Quality Standards (XVC Principles)
- **Specification-Driven**: All code auto-generated from `api.yaml` (Single Source of Truth)
- **Zero Regressions**: Surgical precision in all changes (XVC Surgical Precision)
- **Template Architecture**: 8 externalized templates for maintainability (Pattern Consistency)
- **Performance Optimized**: Efficient API response and WebSocket synchronization (Bar-Raising Solutions)

---

## 🚦 Getting Started

### System Requirements
- **Go 1.21+**: Core runtime and build system
- **Make**: Build automation and task runner
- **jq**: JSON processing for API interactions

### Explore the Platform
1. **Start HD1**: `cd src && make clean && make && make start`
2. **Open Console**: http://localhost:8080 for visual monitoring
3. **Try APIs**: Use the 59 REST endpoints to create your first 3D scene
4. **Join Community**: Explore `/docs/` for comprehensive guides

---

## 🤝 Contributing & Community

### Project Status
**HD1 v5.0.0**: API-first 3D engine platform in active development  
**Next Phase**: Service integration layer and distributed service ecosystem

### Architecture Philosophy (XVC Methodology)
- **API-First**: Every capability exposed via REST endpoints (Pattern Consistency)
- **Template-Driven**: All code generated from external templates (LLM as Pattern Reflection Engine)
- **Single Source of Truth**: OpenAPI specification drives everything (XVC Core Principle)
- **Quality Focus**: Zero regressions, surgical precision in changes (XVC Principles in Practice)

### Support Resources
- **Complete Documentation**: `/docs/` with professional taxonomy
- **Architectural Decisions**: `/docs/decisions/` with 23 ADRs
- **Developer Guides**: `/src/README.md` for implementation details
- **Template Architecture**: `/src/codegen/templates/README.md`

---

**HD1 v5.0.0** | MIT License | API-First 3D Engine for Distributed Services  
*Where distributed services meet immersive 3D experiences*