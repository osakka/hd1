# HD1 Lighting Library

## Overview
The HD1 Lighting Library provides reusable lighting components, utilities, and calculations for creating realistic illumination in holodeck environments.

## Architecture
- **Library Components**: Reusable lighting utilities and calculations
- **Props Integration**: Used by lighting props (light-bulb, etc.)
- **Scene Integration**: Supports dynamic lighting in scenes
- **Environment Adaptation**: Lighting adapts to different environment contexts

## Components

### Core Lighting Types
- `point-light.sh` - Point light sources (bulbs, LEDs)
- `directional-light.sh` - Directional lighting (sun, laser)
- `ambient-light.sh` - Ambient lighting calculations
- `spot-light.sh` - Focused cone lighting

### Utilities
- `color-temperature.sh` - Color temperature calculations (1000K-10000K)
- `light-physics.sh` - Light physics and falloff calculations
- `lighting-presets.sh` - Common lighting scenarios

### Standards
- **Color Temperature**: Kelvin-based color calculations
- **Intensity**: Real-world lumens scaling
- **Distance**: Realistic light falloff
- **Materials**: PBR-compatible light interactions

## Usage
```bash
source "${HD1_ROOT}/share/lighting/point-light.sh"
create_point_light "desk_lamp" 2 1.5 0 "warm_white" 1.0
```

## Integration
- Props call lighting library functions
- Scenes use lighting presets
- Environments modify lighting behavior
- API creates lighting objects using these utilities