# ADR-004: 3D Coordinate System Design

## Status
✅ **ACCEPTED** - Implemented and Operational

## Context

VWS (Virtual World Synthesizer) requires a consistent, predictable 3D coordinate system for object placement, camera positioning, and world boundaries. The system must be intuitive for developers, efficient for rendering, and provide consistent behavior across different client screen sizes and orientations.

## Decision

**We implement a fixed-grid coordinate system with universal boundaries of [-12, +12] on all axes (X, Y, Z), providing a 25×25×25 coordinate space with automatic boundary validation.**

### Core Coordinate System
```
Coordinate Space: [-12, +12] × [-12, +12] × [-12, +12]
Grid Resolution: 25 × 25 × 25 = 15,625 possible positions
Origin: (0, 0, 0) at world center
Axis Orientation: Right-handed coordinate system
```

## Implementation Details

### Boundary Validation
```go
// Automatic coordinate validation at storage layer
if x < -12 || x > 12 || y < -12 || y > 12 || z < -12 || z > 12 {
    return nil, &CoordinateError{
        Message: "Coordinates must be within [-12, +12] bounds"
    }
}
```

### Object Positioning
```go
type Object struct {
    Name     string  `json:"name"`      // Unique identifier
    Type     string  `json:"type"`      // cube, sphere, etc.
    X        float64 `json:"x"`         // X coordinate [-12, +12]
    Y        float64 `json:"y"`         // Y coordinate [-12, +12]
    Z        float64 `json:"z"`         // Z coordinate [-12, +12]
    Scale    float64 `json:"scale"`     // Size multiplier (default: 1.0)
    Rotation string  `json:"rotation"`  // Rotation specification
    Color    string  `json:"color"`     // Color specification
}
```

### World Configuration
```go
type World struct {
    Size         int     `json:"size"`          // Grid size (default: 25)
    Transparency float64 `json:"transparency"`  // World transparency (0.0-1.0)
    CameraX      float64 `json:"camera_x"`      // Camera X position
    CameraY      float64 `json:"camera_y"`      // Camera Y position  
    CameraZ      float64 `json:"camera_z"`      // Camera Z position
}
```

## Coordinate System Properties

### 1. **Fixed Boundaries**
- **Universal Limits**: All sessions use identical coordinate boundaries
- **Predictable Behavior**: Objects cannot be placed outside valid ranges
- **Error Prevention**: Invalid coordinates rejected at API level

### 2. **Grid-Aligned Positioning**
- **Integer Coordinates**: Encourages grid-aligned object placement
- **Float Support**: Allows sub-grid positioning for precise placement
- **Scaling Independence**: Coordinate system independent of rendering scale

### 3. **Origin-Centered Design**
- **Center Origin**: (0, 0, 0) at world center for intuitive positioning
- **Symmetric Bounds**: Equal positive and negative ranges
- **Balanced Distribution**: Objects naturally distributed around center

## Axis Conventions

### Coordinate Orientation
```
X-Axis: Left (-12) to Right (+12)
Y-Axis: Down (-12) to Up (+12)  
Z-Axis: Back (-12) to Front (+12)
```

### Right-Handed System
- **Positive X**: Right direction
- **Positive Y**: Up direction
- **Positive Z**: Forward/toward viewer
- **Rotation**: Counter-clockwise around positive axes

## Camera System Integration

### Default Camera Positioning
```go
// Standard camera placement for optimal viewing
CameraX: 15   // Slightly outside world bounds
CameraY: 15   // Elevated view
CameraZ: 15   // Angled perspective
```

### Orbital Camera Mathematics
- **Focus Point**: World center (0, 0, 0)
- **Orbit Radius**: Configurable (default: ~21 units)
- **Elevation Range**: Prevents camera from going below world floor

## Rendering Considerations

### WebGL Transformation
```javascript
// Convert VWS coordinates to WebGL normalized coordinates
function vwsToWebGL(vwsCoord) {
    return vwsCoord / 12.0;  // Maps [-12, +12] to [-1, +1]
}

// Scaling factor for consistent rendering
const RENDER_SCALE = 0.65;  // Discovered optimal scale factor
```

### Viewport Adaptation
- **Aspect Ratio Independence**: Coordinate system unaffected by screen size
- **Consistent Scale**: Objects maintain relative sizes across devices
- **Responsive Rendering**: Rendering adapts while coordinates remain fixed

## Validation and Error Handling

### Coordinate Validation
```go
func validateCoordinates(x, y, z float64) error {
    if x < -12 || x > 12 {
        return &CoordinateError{Message: "X coordinate must be within [-12, +12] bounds"}
    }
    if y < -12 || y > 12 {
        return &CoordinateError{Message: "Y coordinate must be within [-12, +12] bounds"}
    }
    if z < -12 || z > 12 {
        return &CoordinateError{Message: "Z coordinate must be within [-12, +12] bounds"}
    }
    return nil
}
```

### Error Messages
- **Specific Axis**: Error messages identify which coordinate is invalid
- **Clear Boundaries**: Always specify valid range in error messages
- **User-Friendly**: Non-technical language for API consumers

## Alternative Approaches Considered

### 1. **Normalized Coordinates [-1, +1]**
**Rejected**: Too abstract for users, difficult to reason about object positions

### 2. **Large Integer Grid [0, 1000]**
**Rejected**: Arbitrary scale, no intuitive meaning, waste of precision

### 3. **Dynamic Boundaries Based on Screen Size**
**Rejected**: Inconsistent behavior across clients, complex synchronization

### 4. **Infinite Coordinate Space**
**Rejected**: No natural boundaries, difficult to implement camera limits

## Consequences

### ✅ Benefits
- **Predictable Behavior**: Consistent coordinate system across all sessions
- **Intuitive Scale**: [-12, +12] range easy to understand and remember
- **Error Prevention**: Automatic boundary validation prevents invalid states
- **Rendering Efficiency**: Fixed bounds enable optimization in 3D renderer
- **Grid Alignment**: Natural grid encourages organized object placement

### ⚠️ Trade-offs
- **Limited Space**: 25×25×25 grid may be insufficient for complex scenes
- **Fixed Scale**: Cannot dynamically adjust world size based on content
- **Integer Bias**: Decimal coordinates less intuitive than integers

### 🔧 Mitigation
- **Future Expansion**: Design allows changing boundaries in future versions
- **Sub-grid Precision**: Float coordinates enable precise positioning
- **Multiple Sessions**: Create multiple worlds for larger scenes

## Real-world Usage Patterns

### Common Object Placements
```javascript
// Center objects
{x: 0, y: 0, z: 0}      // World center
{x: 0, y: -10, z: 0}    // Floor level
{x: 0, y: 10, z: 0}     // Ceiling level

// Corner placements
{x: -10, y: -10, z: -10}  // Back-left-bottom corner
{x: 10, y: 10, z: 10}     // Front-right-top corner

// Grid-aligned positioning
{x: -6, y: 0, z: 3}       // Grid intersection
{x: 2.5, y: 4.5, z: -1.5} // Sub-grid precision
```

### Camera Positioning Patterns
```javascript
// Standard views
{x: 15, y: 15, z: 15}     // Isometric view
{x: 0, y: 0, z: 20}       // Front view
{x: 20, y: 0, z: 0}       // Side view
{x: 0, y: 20, z: 0}       // Top view
```

## Performance Impact

### Validation Overhead
- **Constant Time**: O(1) coordinate validation per object
- **Minimal CPU**: Simple comparison operations
- **Early Rejection**: Invalid coordinates rejected before storage

### Memory Efficiency  
- **Compact Storage**: Float64 coordinates (8 bytes × 3 = 24 bytes per position)
- **No Index Overhead**: Direct coordinate storage without spatial indexing
- **Cache Friendly**: Sequential coordinate access patterns

## Future Enhancements

### 1. **Configurable Boundaries**
```go
type WorldConfig struct {
    MinX, MaxX float64  // Customizable X boundaries
    MinY, MaxY float64  // Customizable Y boundaries  
    MinZ, MaxZ float64  // Customizable Z boundaries
}
```

### 2. **Coordinate Transformations**
- World-to-world coordinate mapping
- Scale transformations for different world sizes
- Coordinate system conversions (spherical, cylindrical)

### 3. **Spatial Indexing**
- Quadtree/Octree for efficient spatial queries
- Collision detection optimization
- Proximity-based object clustering

## Related Decisions
- [ADR-001: Specification-Driven Development](001-specification-driven-development.md)
- [ADR-002: Thread-Safe SessionStore](002-thread-safe-session-store.md)
- [ADR-003: WebSocket Real-time Architecture](003-websocket-realtime-architecture.md)

## References
- [WebGL Coordinate Systems](https://webglfundamentals.org/webgl/lessons/webgl-3d-orthographic.html)
- [3D Mathematics for Computer Graphics](https://mathworld.wolfram.com/CoordinateSystem.html)

---

**Decision Date**: 2025-06-28  
**Decision Makers**: VWS Architecture Team  
**Review Date**: 2025-12-28  

*This ADR establishes the spatial foundation that enables precise object placement in VWS virtual worlds.*