# HD1 Architecture Overview

## Three-Layer Game Engine Architecture

HD1 implements a comprehensive three-layer architecture matching industry-standard game engine patterns:

### Layer 1: Environment System
- **Physics Contexts**: 4 distinct environments (Earth Surface, Molecular Scale, Space Vacuum, Underwater)
- **Dynamic Properties**: Gravity, atmosphere, scale units, density, temperature
- **API Endpoints**: `/environments` (GET/POST)
- **Storage**: `/opt/hd1/share/environments/` shell scripts

### Layer 2: Props System  
- **Reusable Objects**: 6 categories of physics-accurate objects
- **YAML Definitions**: Structured specifications with material properties
- **API Endpoints**: `/props` (GET), `/sessions/{sessionId}/props/{propId}` (POST)
- **Storage**: `/opt/hd1/share/props/` YAML files

### Layer 3: Scene Orchestration
- **Composition**: Environment + Props placement
- **Scene Management**: Loading, forking, saving
- **API Integration**: Full scene lifecycle through API

## Core Systems

### Specification-Driven Development
- **OpenAPI 3.0.3**: Single source of truth for all 31 endpoints
- **Auto-Generation**: Routing, clients, documentation from specification
- **Build Validation**: Prevents deployment of incomplete implementations

### Physics Cohesion Engine
- **Environment-Aware**: Props adapt physics based on session environment
- **Real-Time**: Physics recalculated instantly on environment changes
- **Material Accuracy**: Realistic properties (wood: 600 kg/m³, metal: 7800 kg/m³)

### Concurrent Session Management
- **Thread-Safe**: Mutex-protected session isolation
- **Multi-User**: Independent 3D worlds per session
- **WebSocket**: Real-time object synchronization

## Technology Stack
- **Backend**: Go with OpenAPI-generated routing
- **Frontend**: A-Frame WebXR for VR/AR rendering
- **API**: REST/WebSocket hybrid
- **Storage**: File-based with YAML configurations
- **Build**: Make-based with validation pipeline