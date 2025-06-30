# A-Frame Component Library

This directory contains the complete A-Frame ecosystem for the HD1 Holodeck platform.

## Core A-Frame
- `aframe.min.js` - A-Frame 1.4.0 core framework
- `aframe-extras.min.js` - Extra components (controls, loaders, primitives)

## Physics & Animation
- `aframe-physics-system.min.js` - Physics simulation (Cannon.js/Ammo.js)
- `aframe-animation-component.min.js` - Smooth property animations

## Visual Effects & Environment
- `aframe-environment-component.min.js` - Procedural environments
- `aframe-particle-system.js` - Particle effects (fire, smoke, sparkles)

## Interaction & Controls
- `aframe-teleport-controls.min.js` - VR teleportation system
- `aframe-orbit-controls.min.js` - Orbit camera controls
- `aframe-controller-cursor-component.min.js` - VR controller interactions

## Utilities & Components
- `aframe-event-set-component.min.js` - Event-driven property changes
- `aframe-look-at-component.min.js` - Object orientation utilities
- `aframe-text-geometry-component.min.js` - 3D text rendering
- `aframe-state-component.min.js` - Application state management

## Data Visualization
- `aframe-forcegraph-component.min.js` - Interactive force-directed graphs

## Usage
Include required components in HTML:
```html
<script src="/static/js/vendor/aframe.min.js"></script>
<script src="/static/js/vendor/aframe-physics-system.min.js"></script>
<!-- Add other components as needed -->
```