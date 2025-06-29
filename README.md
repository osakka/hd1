# THD (The Holo-Deck) - Professional VR/AR Holodeck Platform

**Revolutionary API-first 3D/VR visualization engine powered by A-Frame WebXR**

Where professional engineering meets immersive holodeck technology.

## 🌟 Revolutionary Features

### 🏆 Upstream/Downstream API Integration (v3.5.0)
- **🔗 Revolutionary Bridge System**: Complete upstream/downstream API integration
- **⚙️ Enhanced Shell Functions**: A-Frame capabilities exposed through professional shell interface
- **🌐 JavaScript Function Bridge**: Identical signatures to shell functions with A-Frame validation
- **📐 Single Source of Truth**: API specification + A-Frame schemas drive all generation
- **🎯 Professional Parameter Validation**: Enterprise-grade error handling with actionable messages
- **🚀 Zero Manual Synchronization**: Shell and JavaScript functions auto-generated from specs

### Core Holodeck Technology
- **🥽 Full VR/AR Support**: Complete WebXR integration with headset compatibility
- **🎨 A-Frame WebXR Engine**: Built on Mozilla's powerful A-Frame framework (MIT License)
- **⚡ API-First Architecture**: Everything controlled via REST API
- **🌍 Universal Coordinate System**: Professional [-12, +12] holodeck boundaries
- **🔄 Real-time WebSocket**: Instant 3D object synchronization
- **📊 Specification-Driven**: OpenAPI 3.0.3 single source of truth

### 🎬 Scene Forking & Recording System (v3.4.0)
- **📷 FREEZE-FRAME Mode**: Save current session state as new scenes
- **🎥 TEMPORAL SEQUENCE Mode**: Full temporal recording with playback capabilities
- **🔀 Scene Forking**: Load scenes for non-destructive editing
- **🏷️ Object Tracking**: Complete provenance system (base/modified/new)
- **📝 Script Generation**: Dynamic scene file creation from session state
- **🎛️ Professional Console**: Integrated FREEZE-FRAME/TEMPORAL SEQUENCE controls

### 👑 Crown Jewel - Auto-Generated Web UI Client
- **🌐 JavaScript API Client**: Complete API wrapper auto-generated from specification
- **🎨 UI Component Library**: Each API endpoint becomes interactive component
- **📝 Dynamic Form System**: Forms automatically generated from request schemas
- **🔄 100% Single Source of Truth**: All clients generated from same OpenAPI spec
- **⚡ Zero Manual Sync**: API changes automatically update all clients

### Professional Engineering
- **🏗️ Auto-Generated APIs**: Complete routing from OpenAPI specification
- **🛡️ Session Isolation**: Perfect multi-user separation
- **📱 Professional UI**: Holodeck-themed interface with scene management
- **⚡ Sprint Controls**: Shift-key acceleration for enhanced navigation
- **🎭 Multiple Rendering Backends**: Extensible architecture for future engines

## 🚀 Quick Start

```bash
# Build the holodeck system
cd src && make all

# Start the holodeck daemon
make start

# Navigate to http://localhost:8080
# Use WASD to move, mouse to look around
# Click VR button for full immersive experience
```

## 🎬 Scene Forking & Recording Workflow

### FREEZE-FRAME Mode (Scene Snapshots)
```bash
# Create a session
SESSION_ID=$(thd-client create-session | jq -r '.session_id')

# Fork an existing scene for editing
thd-client fork-scene basic-shapes "{\"session_id\": \"$SESSION_ID\", \"clear_existing\": true}"

# Make modifications in the UI (add/move/delete objects)
# Then save as new scene using the FREEZE-FRAME button in console
# Or via CLI:
thd-client save-scene-from-session "$SESSION_ID" '{"scene_id": "my-scene", "name": "My Custom Scene"}'
```

### TEMPORAL SEQUENCE Mode (Recording)
```bash
# Start recording session interactions
thd-client start-recording "$SESSION_ID" '{"name": "Demo Recording", "description": "User interaction demo"}'

# Perform actions (create/modify objects, move camera, etc.)
# Recording captures ALL changes with timestamps

# Stop recording
thd-client stop-recording "$SESSION_ID"

# Play back recorded session in another session
NEW_SESSION=$(thd-client create-session | jq -r '.session_id')
thd-client play-recording "$NEW_SESSION" '{"recording_id": "demo-recording"}'
```

### Object Tracking System
- **Base Objects**: Loaded from forked scenes (original state preserved)
- **Modified Objects**: Base objects that have been changed (tracks source)
- **New Objects**: Created directly in session (marked for inclusion in saved scenes)

## 🏗️ Architecture

### Core Engine: A-Frame WebXR
THD leverages **[A-Frame](https://aframe.io)** (MIT License) as its primary rendering backend:

- **A-Frame Version**: 1.4.0 WebXR
- **License**: MIT License - [https://github.com/aframevr/aframe/blob/master/LICENSE](https://github.com/aframevr/aframe/blob/master/LICENSE)
- **Entity-Component-System**: Professional ECS architecture
- **WebXR Standard**: Full VR/AR headset compatibility
- **Cross-Platform**: Desktop, mobile, and VR devices

### Multi-Backend Architecture
THD is designed for **framework flexibility**:

```
┌─────────────────────┐
│   THD API Layer     │  ← Universal REST/WebSocket interface
├─────────────────────┤
│  Rendering Backends │  ← Pluggable engine architecture
│                     │
│  🔹 A-Frame WebXR   │  ← Current: VR/AR holodeck
│  🔸 Three.js WebGL  │  ← Future: Direct WebGL
│  🔸 Babylon.js      │  ← Future: Alternative engine
│  🔸 Unity WebGL     │  ← Future: Game engine
│  🔸 Custom Engines  │  ← Future: Specialized renderers
└─────────────────────┘
```

**Future Vision**: Session-based backend selection allowing users to choose their preferred rendering engine based on needs:
- **A-Frame**: For VR/AR immersive experiences
- **Three.js**: For lightweight WebGL performance
- **Babylon.js**: For advanced graphics features
- **Custom**: For specialized industrial applications

### Professional Backend
- **Go Daemon**: High-performance concurrent server
- **Session Management**: Isolated 3D worlds with persistence
- **WebSocket Hub**: Real-time object synchronization
- **Auto-Generated API**: Specification-driven development

## 🎮 Holodeck Capabilities

### Basic 3D Objects
```bash
# Create objects via API
curl -X POST localhost:8080/api/sessions/{id}/objects \
  -d '{"name": "red_cube", "type": "cube", "x": 0, "y": 1, "z": 0, "color": {"r": 1.0, "g": 0.2, "b": 0.2, "a": 1.0}}'
```

### Advanced A-Frame Features
- **🏗️ Physically-Based Rendering**: Metalness, roughness, emissive materials
- **⚡ Physics Simulation**: Dynamic, static, and kinematic bodies  
- **💡 Advanced Lighting**: Directional, point, ambient, and spot lights
- **✨ Particle Effects**: Fire, smoke, sparkles, and custom systems
- **📝 3D Text Rendering**: Holographic text displays in 3D space
- **🌌 Environment Systems**: Sky domes, fog, and atmospheric effects
- **🎭 Animation Support**: Object movement and transformation
- **🔊 Spatial Audio**: 3D positional sound (future)

### VR/AR Interaction
- **🥽 Headset Support**: Oculus, HTC Vive, Magic Leap, etc.
- **🖱️ Desktop Controls**: WASD movement, mouse look
- **📱 Mobile Compatible**: Touch controls for mobile devices
- **🤲 Hand Tracking**: Future WebXR hand tracking support

## 📋 API Reference

### Session Management
```bash
POST /api/sessions              # Create new holodeck session
GET  /api/sessions              # List active sessions
GET  /api/sessions/{id}         # Get session details
```

### Object Creation
```bash
POST /api/sessions/{id}/objects              # Create 3D objects
GET  /api/sessions/{id}/objects              # List all objects
PUT  /api/sessions/{id}/objects/{name}       # Update object
DELETE /api/sessions/{id}/objects/{name}     # Delete object
```

### Advanced Controls
```bash
PUT  /api/sessions/{id}/camera/position      # Control camera
POST /api/sessions/{id}/camera/orbit         # Camera orbital motion
POST /api/browser/refresh                    # Force browser refresh
POST /api/browser/canvas                     # Direct canvas control
```

## 🛠️ Development

### Professional Standards
- **Absolute Paths Only**: No relative path confusion
- **Long Flags Only**: No short flags for clarity
- **API-First Design**: Zero shell command dependencies
- **Specification-Driven**: OpenAPI 3.0.3 generates all routing
- **Professional Logging**: Timestamped, structured output
- **Clean Architecture**: Separation of concerns

### Build System
```bash
cd src
make validate    # Validate OpenAPI specification
make generate    # Generate routing from spec
make build       # Build THD daemon binary
make test        # Run API endpoint tests
```

### Daemon Management
```bash
make start       # Start THD holodeck daemon
make stop        # Stop daemon with clean shutdown
make restart     # Restart with validation
make status      # Professional status reporting
```

## 🗂️ File Structure

```
/opt/holo-deck/
├── src/                          # Go source code
│   ├── main.go                   # THD daemon entry point
│   ├── auto_router.go            # Auto-generated API routing
│   ├── api.yaml                  # OpenAPI specification (single source of truth)
│   ├── api/                      # API handler packages
│   └── server/                   # Core server infrastructure
├── share/
│   └── htdocs/
│       └── static/js/
│           └── thd-aframe.js     # A-Frame holodeck integration
├── lib/
│   └── thd-functions.sh          # Holodeck shell function library
├── scenarios/
│   └── ultimate-holodeck.thd     # Example holodeck scenarios
├── docs/                         # Architecture Decision Records
└── build/                        # Build artifacts (excluded from git)
    ├── bin/thd                   # THD daemon binary
    ├── runtime/thd.pid           # Process management
    └── logs/                     # Professional logging
```

## 🎯 Holodeck Script Library

THD includes a comprehensive shell function library for rapid holodeck development:

```bash
# Load THD functions
source lib/thd-functions.sh

# Basic objects
thd::create_object "my_cube" cube 0 1 0

# Advanced A-Frame features
thd::create_light "sun" directional 10 10 5 1.2 "#ffffff"
thd::create_physics "bouncing_ball" sphere 0 5 0 2.0 "dynamic"
thd::create_material "metal_pillar" cylinder 2 1 2 "standard" 0.8 0.1
thd::create_particles "campfire" fire 0 0 0 1000
thd::create_text "welcome" "HOLODECK ACTIVE" 0 3 -5 1.0 1.0 0.0
thd::create_sky "environment" "#1a1a2e"
```

## 🌍 World Constraints

- **Coordinate System**: 25×25×25 professional grid
- **Boundaries**: [-12, +12] on all X, Y, Z axes
- **Floor Level**: Y=0 (world floor)
- **Eye Level**: Y=1.7 (human standing height)
- **Max Objects**: 1000 per session
- **Max Sessions**: 100 concurrent

## 📄 Licensing & Attribution

### THD (The Holo-Deck)
- **License**: MIT License
- **Copyright**: Professional engineering standards

### A-Frame WebXR Framework
- **Project**: [A-Frame](https://aframe.io) by Mozilla
- **Version**: 1.4.0
- **License**: MIT License
- **Repository**: [https://github.com/aframevr/aframe](https://github.com/aframevr/aframe)
- **Documentation**: [https://aframe.io/docs/](https://aframe.io/docs/)

THD gratefully acknowledges the A-Frame community for creating the world's most accessible WebXR framework. A-Frame's Entity-Component-System architecture and comprehensive WebXR support make THD's holodeck vision possible.

### Integration Philosophy
THD demonstrates how easy it is to integrate powerful open-source frameworks:
- **Clean API Layer**: Framework-agnostic REST/WebSocket interface
- **Pluggable Architecture**: Easy to swap rendering backends
- **Professional Standards**: Maintain engineering quality across integrations
- **Community-Driven**: Leverage the best open-source tools available

## 🔮 Future Roadmap

### Multi-Backend Support
- **Configuration-Based Selection**: Choose rendering engine per session
- **Performance Optimization**: Match engine to use case
- **Specialized Backends**: CAD, gaming, scientific visualization
- **Engine Comparison**: A/B testing different frameworks

### Advanced Features
- **Collaborative VR**: Multi-user shared holodeck spaces
- **Persistence Layer**: Save/load holodeck configurations
- **Asset Pipeline**: Import 3D models, textures, animations
- **Scripting Engine**: Lua/JavaScript holodeck programming
- **AI Integration**: Procedural content generation
- **Cloud Deployment**: Scalable holodeck infrastructure

## 🎬 Example Scenarios

Run the ultimate holodeck demonstration:
```bash
# Set your session ID
export THD_SESSION_ID="your-session-id"

# Run the ultimate scenario
./scenarios/ultimate-holodeck.thd
```

This creates a complete holodeck experience with:
- 🌌 Sky environment with atmospheric effects
- 💡 Cinematic lighting system (4 light sources)
- 🏗️ Circular metallic platform foundation
- 💎 Crystal formations with ultimate materials
- ✨ Particle effects (fire, smoke, sparkles)
- ⚡ Physics simulation with bouncing spheres
- 🏛️ Architectural elements (glass walls, metal beams)
- 📝 3D holographic text displays
- 🎛️ Interactive control panels
- 🎨 Floating artistic sculptures

## 🤝 Contributing

THD welcomes contributions to expand holodeck capabilities:
- **New Rendering Backends**: Integrate additional 3D engines
- **API Extensions**: Expand holodeck functionality
- **VR/AR Features**: Enhance immersive experiences
- **Performance Optimization**: Improve real-time rendering
- **Documentation**: Share holodeck knowledge

---

**THD (The Holo-Deck)** - Where professional engineering meets immersive holodeck technology.

*Powered by A-Frame WebXR • Engineered for the future • Ready for VR*