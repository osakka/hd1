# THD Holodeck Scenes

**Professional Scene Collection for THD (The Holo-Deck)**

This directory contains the complete collection of holodeck scenes available through THD's professional scene management system.

## 🎭 Scene Collection

### [empty-grid.sh](empty-grid.sh)
Clean holodeck grid foundation - perfect starting point for custom content creation.

### [basic-shapes.sh](basic-shapes.sh) 
Fundamental geometric shapes demonstration showcasing THD's primitive object capabilities.

### [anime-ui.sh](anime-ui.sh)
Advanced UI demonstration with anime-inspired visual elements and interactive components.

### [ultimate-demo.sh](ultimate-demo.sh)
Complete THD capabilities showcase featuring lighting, materials, physics, and A-Frame integration.

### [my-custom-scene.sh](my-custom-scene.sh)
Template for custom scene development with professional patterns and best practices.

## 🎯 Scene Management

### API Integration
All scenes are accessible through THD's professional scene management API:

```bash
# List available scenes
curl http://localhost:8080/api/scenes

# Load a scene
curl -X POST http://localhost:8080/api/scenes/ultimate-demo
```

### Web Interface
Professional scene dropdown available in the holodeck console:
- 30-day cookie persistence
- Automatic scene restoration 
- Session-specific scene management

## 📋 Scene Development Standards

### Scene Structure
Each scene file follows the professional THD pattern:

```bash
#!/bin/bash
# Scene: [NAME] - [DESCRIPTION]
# Professional THD scene with enhanced capabilities

# Load THD enhanced functions
source "/opt/holo-deck/lib/thd-enhanced-functions.sh"

# Scene implementation
thd::create_enhanced_object "example" "cube" 0 1 0 --color "#ff0000"
```

### Professional Standards
- **Enhanced Functions**: Use `/opt/holo-deck/lib/thd-enhanced-functions.sh`
- **A-Frame Integration**: Leverage complete A-Frame schema validation
- **Error Handling**: Professional validation and error reporting
- **Documentation**: Clear scene purpose and feature descriptions

## 🔧 Technical Integration

### Server Integration
Scenes are served by THD's API endpoint handlers in:
- `src/api/scenes/list.go` - Scene discovery and listing
- `src/api/scenes/load.go` - Scene execution and loading

### File Path Resolution
THD automatically resolves scene paths from this directory:
```
/opt/holo-deck/share/scenes/[scene-name].sh
```

## 🎨 Scene Categories

### Foundation Scenes
- **empty-grid.sh**: Clean starting point
- **basic-shapes.sh**: Primitive demonstrations

### Advanced Demonstrations
- **anime-ui.sh**: UI and interface showcase
- **ultimate-demo.sh**: Complete capability demonstration

### Development Templates
- **my-custom-scene.sh**: Custom scene development template

---

**Professional holodeck scenes powered by THD's revolutionary VR/AR architecture with A-Frame WebXR integration.**