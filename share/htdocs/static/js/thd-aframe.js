/**
 * THD A-Frame Holodeck Manager
 * Replaces custom WebGL renderer with A-Frame ECS system
 * Maintains 100% compatibility with existing THD WebSocket protocol
 */

// Register THD Holodeck component for scene management
AFRAME.registerComponent('thd-holodeck', {
    schema: {
        sessionId: {type: 'string', default: ''},
        gridSize: {type: 'number', default: 25},
        gridTransparency: {type: 'number', default: 0.01}
    },

    init: function () {
        console.log('[THD-AFrame] Initializing holodeck scene');
        this.objects = new Map();
        this.sceneEl = this.el;
        this.objectsContainer = document.getElementById('holodeck-objects');
        this.gridContainer = document.getElementById('holodeck-grid');
        
        // Initialize coordinate system
        this.setupCoordinateSystem();
        
        // Ready for WebSocket integration
        this.isReady = true;
        console.log('[THD-AFrame] Scene ready for holodeck integration');
    },

    setupCoordinateSystem: function() {
        // Create holodeck coordinate grid (25x25x25, Y=0 floor, Y=1.7 eye level)
        const gridMaterial = `color: #00ffff; transparent: true; opacity: ${this.data.gridTransparency}`;
        
        // Floor grid lines
        for (let x = -12; x <= 12; x++) {
            for (let z = -12; z <= 12; z++) {
                if (x % 2 === 0 || z % 2 === 0) {
                    const gridLine = document.createElement('a-box');
                    gridLine.setAttribute('position', `${x} 0 ${z}`);
                    gridLine.setAttribute('scale', '0.1 0.02 0.1');
                    gridLine.setAttribute('material', gridMaterial);
                    this.gridContainer.appendChild(gridLine);
                }
            }
        }
        
        console.log('[THD-AFrame] Holodeck coordinate system initialized');
    }
});

// THD Manager class - compatible with existing WebSocket interface
class THDAFrameManager {
    constructor(scene) {
        this.scene = scene;
        this.objects = new Map();
        this.camera = this.setupCamera();
        this.objectsContainer = document.getElementById('holodeck-objects');
        
        console.log('[THD-AFrame] Manager initialized');
    }

    setupCamera() {
        const cameraEl = document.getElementById('holodeck-camera');
        return {
            position: [0, 1.7, 5], // Compatible with existing interface
            target: [0, 0, 0],
            element: cameraEl
        };
    }

    // Main interface method - processes WebSocket messages
    processMessage(message) {
        console.log('[THD-AFrame] Processing message:', message.type, message);
        
        switch (message.type) {
            case 'create':
                this.createObjects(message.objects);
                break;
            case 'clear':
                this.clearObjects();
                break;
            case 'camera':
                this.updateCamera(message.camera);
                break;
            case 'update':
                this.updateObjects(message.objects);
                break;
            default:
                console.warn('[THD-AFrame] Unknown message type:', message.type);
        }
    }

    createObjects(objects) {
        if (!objects || !Array.isArray(objects)) return;
        
        objects.forEach(obj => {
            this.createObject(obj);
        });
        
        console.log('[THD-AFrame] Created', objects.length, 'objects');
    }

    createObject(obj) {
        const entity = document.createElement('a-entity');
        const id = obj.id || obj.name;
        
        // Set position using holodeck coordinates
        const pos = obj.transform?.position || {x: obj.x || 0, y: obj.y || 0, z: obj.z || 0};
        entity.setAttribute('position', `${pos.x} ${pos.y} ${pos.z}`);
        
        // Set scale
        const scale = obj.transform?.scale || {x: obj.scale || 1, y: obj.scale || 1, z: obj.scale || 1};
        entity.setAttribute('scale', `${scale.x} ${scale.y} ${scale.z}`);
        
        // Add geometry based on type
        this.setGeometry(entity, obj);
        
        // Add material with enhanced A-Frame capabilities
        this.setMaterial(entity, obj);
        
        // Store reference
        entity.setAttribute('id', `thd-${id}`);
        this.objects.set(id, entity);
        this.objectsContainer.appendChild(entity);
        
        console.log('[THD-AFrame] Created object:', id, 'at', pos);
    }

    setGeometry(entity, obj) {
        switch (obj.type) {
            case 'cube':
                entity.setAttribute('geometry', {
                    primitive: 'box',
                    width: 1,
                    height: 1,
                    depth: 1
                });
                break;
            case 'sphere':
                entity.setAttribute('geometry', {
                    primitive: 'sphere',
                    radius: 0.5
                });
                break;
            case 'plane':
                entity.setAttribute('geometry', {
                    primitive: 'plane',
                    width: 1,
                    height: 1
                });
                break;
            case 'light':
                this.createLight(entity, obj);
                return; // Lights don't need geometry
            default:
                // Default to cube for unknown types
                entity.setAttribute('geometry', {
                    primitive: 'box',
                    width: 1,
                    height: 1,
                    depth: 1
                });
        }
    }

    setMaterial(entity, obj) {
        const color = obj.color || {r: 0.2, g: 0.8, b: 0.2, a: 1.0};
        const hexColor = `#${Math.round(color.r * 255).toString(16).padStart(2, '0')}${Math.round(color.g * 255).toString(16).padStart(2, '0')}${Math.round(color.b * 255).toString(16).padStart(2, '0')}`;
        
        const materialProps = {
            color: hexColor,
            transparent: color.a < 1.0,
            opacity: color.a
        };
        
        if (obj.wireframe) {
            materialProps.wireframe = true;
        }
        
        entity.setAttribute('material', materialProps);
    }

    createLight(entity, obj) {
        // Enhanced lighting system
        const lightProps = {
            type: obj.lightType || 'point',
            color: obj.color ? this.colorToHex(obj.color) : '#ffffff',
            intensity: obj.intensity || 1.0
        };
        
        entity.setAttribute('light', lightProps);
        console.log('[THD-AFrame] Created light:', lightProps);
    }

    clearObjects() {
        // Clear all THD objects but preserve grid and environment
        this.objects.forEach((entity, id) => {
            if (entity.parentNode) {
                entity.parentNode.removeChild(entity);
            }
        });
        this.objects.clear();
        console.log('[THD-AFrame] Cleared all objects');
    }

    updateCamera(cameraData) {
        if (cameraData.position) {
            this.camera.position = cameraData.position;
            this.camera.element.setAttribute('position', 
                `${cameraData.position[0]} ${cameraData.position[1]} ${cameraData.position[2]}`);
        }
    }

    initializeWorld(worldData) {
        console.log('[THD-AFrame] Initializing world:', worldData);
        
        if (worldData && worldData.grid_size) {
            // Update grid if needed
            this.scene.setAttribute('thd-holodeck', {
                gridSize: worldData.grid_size,
                gridTransparency: worldData.transparency || 0.01
            });
        }
        
        console.log('[THD-AFrame] World initialized');
    }

    // Utility methods
    colorToHex(color) {
        const r = Math.round(color.r * 255).toString(16).padStart(2, '0');
        const g = Math.round(color.g * 255).toString(16).padStart(2, '0');
        const b = Math.round(color.b * 255).toString(16).padStart(2, '0');
        return `#${r}${g}${b}`;
    }
}

console.log('[THD-AFrame] THD A-Frame manager loaded');