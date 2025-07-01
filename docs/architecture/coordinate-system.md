# HOLODECK COORDINATE SYSTEM STANDARD
## Single Source of Truth - Production Holodeck Technology

### Physical Holodeck Mapping
```
Y-Axis (Vertical):
  Y = 0.0   : Physical floor (where humans stand)
  Y = 1.7   : Average human eye level (standing)
  Y = 3.0   : Holodeck ceiling
  Y = -0.5  : Below floor (subsurface)

X/Z-Axis (Horizontal Plane):
  Range: [-12.0, +12.0] on both axes
  Grid: 25x25 unit cells
  Center: (0, 0) at holodeck center
  
Camera Default Position:
  Position: (0, 1.7, 0) - Human standing at center
  Looking: Forward along Z-axis
  
Grid Visualization:
  Floor: Y = 0.0 (solid surface)
  Grid lines: Transparent wireframe at Y = 0.0
```

### Coordinate Conventions
- **Right-handed coordinate system**
- **Y-up convention** (Y=0 is floor, positive Y is up)
- **Z-forward convention** (positive Z is forward/north)
- **X-right convention** (positive X is right/east)

### Object Placement Examples
```javascript
// Floor tile at center
{ x: 0, y: 0, z: 0 }

// Object at human eye level
{ x: 0, y: 1.7, z: 0 }

// Object floating near ceiling  
{ x: 0, y: 2.5, z: 0 }

// Floor corner positions
{ x: -12, y: 0, z: -12 }  // Southwest corner
{ x: +12, y: 0, z: +12 }  // Northeast corner
```

### Implementation Rules
1. **All coordinate values MUST use this standard**
2. **No coordinate transformations between server/client**
3. **Grid boundaries are enforced at [-12, +12]**
4. **Floor is always Y=0, never negative**
5. **Camera defaults to human perspective**

---
*This standard defines the coordinate system for production holodeck technology.*