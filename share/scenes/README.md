# HD1 Holodeck Scenes

**Standard Scene Collection for HD1 (Holodeck One)**

This directory contains the complete collection of holodeck scenes available through HD1's standard scene management system.

## 🎭 Scene Collection

### [empty-grid.sh](empty-grid.sh)
Clean holodeck grid foundation - perfect starting point for custom content creation.

### [basic-shapes.sh](basic-shapes.sh) 
Fundamental geometric shapes demonstration showcasing HD1's primitive object capabilities.

### [anime-ui.sh](anime-ui.sh)
Advanced UI demonstration with anime-inspired visual elements and interactive components.

### [complete-demo.sh](complete-demo.sh)
Complete HD1 capabilities showcase featuring lighting, materials, physics, and A-Frame integration.

### [my-custom-scene.sh](my-custom-scene.sh)
Template for custom scene development with standard patterns and best practices.

## 🎯 Scene Management

### API Integration
All scenes are accessible through HD1's standard scene management API:

```bash
# List available scenes
curl http://localhost:8080/api/scenes

# Load a scene
curl -X POST http://localhost:8080/api/scenes/complete-demo
```

### Web Interface
Standard scene dropdown available in the holodeck console:
- 30-day cookie persistence
- Automatic scene restoration 
- Session-specific scene management

## 📋 Scene Development Standards

### Scene Structure
Each scene file follows the standard HD1 pattern:

```bash
#!/bin/bash
# Scene: [NAME] - [DESCRIPTION]
# Standard HD1 scene with enhanced capabilities

# Load HD1 enhanced functions
source "/opt/holodeck-one/lib/hd1-enhanced-functions.sh"

# Scene implementation
hd1::create_enhanced_object "example" "cube" 0 1 0 --color "#ff0000"
```

### Standard Standards
- **Enhanced Functions**: Use `/opt/holodeck-one/lib/hd1-enhanced-functions.sh`
- **A-Frame Integration**: Leverage complete A-Frame schema validation
- **Error Handling**: Standard validation and error reporting
- **Documentation**: Clear scene purpose and feature descriptions

## 🔧 Technical Integration

### Server Integration
Scenes are served by HD1's API endpoint handlers in:
- `src/api/scenes/list.go` - Scene discovery and listing
- `src/api/scenes/load.go` - Scene execution and loading

### File Path Resolution
HD1 automatically resolves scene paths from this directory:
```
/opt/holodeck-one/share/scenes/[scene-name].sh
```

## 🎨 Scene Categories

### Foundation Scenes
- **empty-grid.sh**: Clean starting point
- **basic-shapes.sh**: Primitive demonstrations

### Advanced Demonstrations
- **anime-ui.sh**: UI and interface showcase
- **complete-demo.sh**: Complete capability demonstration

### Development Templates
- **my-custom-scene.sh**: Custom scene development template

---

**Standard holodeck scenes powered by HD1's advanced VR/AR architecture with A-Frame WebXR integration.**