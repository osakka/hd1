# HD1 Downstream Integrations

**Standard A-Frame WebXR Integration Libraries**

This directory contains **downstream integration libraries** that bridge HD1's upstream API capabilities with external frameworks and technologies. Currently focused on A-Frame WebXR integration with identical function signatures across shell and JavaScript environments.

## 🎯 Downstream Architecture

### Integration Philosophy
**Downstream libraries extend upstream capabilities** by:
- **Importing upstream core** - All downstream libraries source `../hd1lib.sh`
- **Adding framework-specific features** - A-Frame schema validation, WebXR capabilities
- **Maintaining identical signatures** - Shell and JavaScript functions have identical parameters
- **Standard validation** - Complete parameter and schema validation

### Advanced Bridge System
```
api.yaml (upstream spec) → hd1lib.* (core API)
         ↓
aframe-schemas (downstream spec) → aframelib.* (A-Frame integration)
```

## 📋 A-Frame Integration Libraries

### [aframelib.sh](aframelib.sh)
**A-Frame Shell Integration** - Enhanced shell functions with complete A-Frame WebXR capabilities

- **Imports upstream core** - Sources `../hd1lib.sh` for API access
- **A-Frame schema validation** - Complete validation using A-Frame component schemas
- **Enhanced object creation** - Materials, lighting, physics, particles
- **Standard error handling** - Actionable error messages with A-Frame context

**Usage:**
```bash
source /opt/holodeck-one/lib/downstream/aframelib.sh

# Enhanced A-Frame functions (identical to JavaScript signatures)
hd1::create_enhanced_object "crystal" "cone" 0 3 0 --color "#ff0000" --metalness 0.8
hd1::create_enhanced_light "key_light" "directional" 10 10 5 --intensity 1.5
```

### [aframelib.js](aframelib.js)
**A-Frame JavaScript Bridge** - Enhanced JavaScript functions with identical shell signatures

- **Advanced bridge** - Identical function signatures to shell functions
- **A-Frame integration** - Direct A-Frame entity manipulation
- **Schema validation** - Browser-side A-Frame schema validation
- **Standard API** - Async/await patterns with comprehensive error handling

**Usage:**
```javascript
// Identical signatures to shell functions
await hd1.createEnhancedObject('crystal', 'cone', 0, 3, 0, {
    color: '#ff0000', 
    metalness: 0.8
});

await hd1.createEnhancedLight('key_light', 'directional', 10, 10, 5, {
    intensity: 1.5
});
```

## 🏆 Advanced Features

### Identical Function Signatures
**Perfect symmetry** between shell and JavaScript:

```bash
# Shell
hd1::create_enhanced_object "name" "type" x y z --option value
```

```javascript
// JavaScript (identical parameters)
hd1.createEnhancedObject('name', 'type', x, y, z, {option: value})
```

### A-Frame Schema Validation
- **Complete schema coverage** - All A-Frame component schemas available
- **Real-time validation** - Parameters validated against A-Frame specifications
- **Standard error messages** - Clear guidance when validation fails
- **Browser and shell** - Validation works in both environments

### Enhanced Capabilities
Beyond basic API functions, downstream libraries provide:
- **Materials & Shaders** - PBR materials, custom shaders, metalness/roughness
- **Advanced Lighting** - Directional, point, ambient, spot lights with shadows
- **Physics Integration** - Dynamic, static, kinematic physics bodies
- **Particle Systems** - Fire, smoke, sparkle effects
- **3D Text Rendering** - Standard 3D text with material properties

## 🔗 Upstream Integration

### Dependency Chain
```
aframelib.sh sources ../hd1lib.sh
    ↓
Enhanced functions use core API functions
    ↓
Standard layered architecture maintained
```

### Architectural Benefits
- **Single source of truth** - Core API changes automatically affect enhanced functions
- **Perfect compatibility** - Enhanced functions build on stable core foundation
- **Standard separation** - Clear distinction between core API and framework integration
- **Maintainability** - Changes to A-Frame integration don't affect core API

## 🎨 A-Frame Specific Features

### Enhanced Object Creation
```bash
hd1::create_enhanced_object "name" "geometry" x y z \
    --color "#rrggbb" \
    --metalness 0.0-1.0 \
    --roughness 0.0-1.0 \
    --emission true/false
```

### Standard Lighting
```bash
hd1::create_enhanced_light "name" "type" x y z \
    --intensity 0.0-5.0 \
    --color "#rrggbb" \
    --cast-shadow true/false
```

### Physics Integration
```bash
hd1::create_physics_object "name" "geometry" x y z \
    --mass 0.0-100.0 \
    --physics-type dynamic/static/kinematic
```

## 🚀 Development Workflow

### Adding New A-Frame Features
1. **Study A-Frame schemas** - Understand component parameters
2. **Add to enhanced functions** - Implement in both shell and JavaScript
3. **Maintain identical signatures** - Ensure shell and JS functions match
4. **Test validation** - Verify A-Frame schema validation works
5. **Document usage** - Update examples and documentation

### Integration Testing
```bash
# Test enhanced functions
source lib/downstream/aframelib.sh
hd1::create_enhanced_object "test" "box" 0 1 0 --color "#ff0000"

# Test scenes using enhanced functions
HD1_SESSION=test bash share/scenes/complete-demo.sh
```

## 📊 A-Frame Integration Coverage

| A-Frame Feature | Shell Support | JS Support | Schema Validation |
|-----------------|---------------|------------|-------------------|
| Basic Geometry | ✅ Complete | ✅ Complete | ✅ Full |
| Materials/PBR | ✅ Complete | ✅ Complete | ✅ Full |
| Lighting | ✅ Complete | ✅ Complete | ✅ Full |
| Physics | ✅ Complete | ✅ Complete | ✅ Full |
| Particles | ✅ Complete | ✅ Complete | ✅ Full |
| Text Rendering | ✅ Complete | ✅ Complete | ✅ Full |

## 🎯 Future Downstream Integrations

The downstream architecture supports future integrations:
- **Three.js integration** - `lib/downstream/threelib.*`
- **Babylon.js integration** - `lib/downstream/babylonlib.*`
- **Unity integration** - `lib/downstream/unitylib.*`
- **Custom frameworks** - Easy to add new downstream integrations

---

**The downstream libraries represent HD1's advanced capability to provide consistent, standard integration with any 3D framework while maintaining perfect single source of truth architecture.**