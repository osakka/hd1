/**
 * THD A-Frame Holodeck Manager
 * Replaces custom WebGL renderer with A-Frame ECS system
 * Maintains 100% compatibility with existing THD WebSocket protocol
 */

// Register boundary enforcement component
AFRAME.registerComponent('holodeck-boundaries', {
    schema: {
        xMin: {type: 'number', default: -11},
        xMax: {type: 'number', default: 11},
        zMin: {type: 'number', default: -11},
        zMax: {type: 'number', default: 11},
        yMin: {type: 'number', default: 0.5},
        yMax: {type: 'number', default: 7}
    },
    
    init: function() {
        this.lastValidPosition = this.el.getAttribute('position');
        this.boundaryCheckInterval = setInterval(this.checkBoundaries.bind(this), 16); // ~60fps checking
        console.log('[Holodeck-Boundaries] COMPONENT INITIALIZED with bounds:', this.data);
        console.log('[Holodeck-Boundaries] Element:', this.el.id);
        
        // Test immediate boundary check
        setTimeout(() => {
            console.log('[Holodeck-Boundaries] Test boundary check in 1 second...');
            this.checkBoundaries();
        }, 1000);
    },
    
    checkBoundaries: function() {
        const position = this.el.getAttribute('position');
        let corrected = false;
        
        // Check X boundaries
        if (position.x < this.data.xMin) {
            position.x = this.data.xMin;
            corrected = true;
        } else if (position.x > this.data.xMax) {
            position.x = this.data.xMax;
            corrected = true;
        }
        
        // Check Z boundaries  
        if (position.z < this.data.zMin) {
            position.z = this.data.zMin;
            corrected = true;
        } else if (position.z > this.data.zMax) {
            position.z = this.data.zMax;
            corrected = true;
        }
        
        // Check Y boundaries
        if (position.y < this.data.yMin) {
            position.y = this.data.yMin;
            corrected = true;
        } else if (position.y > this.data.yMax) {
            position.y = this.data.yMax;
            corrected = true;
        }
        
        if (corrected) {
            this.el.setAttribute('position', position);
            console.log('[Holodeck-Boundaries] Position corrected to:', position);
            
            // Flash red border as visual feedback
            const scene = document.querySelector('a-scene');
            if (scene) {
                scene.style.border = '5px solid red';
                setTimeout(() => {
                    scene.style.border = 'none';
                }, 200);
            }
        } else {
            this.lastValidPosition = {...position};
        }
    },
    
    remove: function() {
        if (this.boundaryCheckInterval) {
            clearInterval(this.boundaryCheckInterval);
        }
    }
});

// Register THD keyboard controls for Q/E turning  
AFRAME.registerComponent('thd-keyboard-controls', {
    init: function () {
        this.onKeyDown = this.onKeyDown.bind(this);
        this.onKeyUp = this.onKeyUp.bind(this);
        this.keys = {};
        this.rotationSpeed = 2.0; // degrees per frame
        this.lookControls = this.el.components['look-controls'];
        
        // Add event listeners with higher priority
        document.addEventListener('keydown', this.onKeyDown, true);
        document.addEventListener('keyup', this.onKeyUp, true);
        
        console.log('[THD-Controls] Q/E turning controls initialized with look-controls override');
    },

    onKeyDown: function (event) {
        console.log('[THD-Controls] ANY key pressed:', event.code, 'pointerLock:', !!document.pointerLockElement);
        
        // Handle ESC key to exit pointer lock
        if (event.code === 'Escape') {
            document.exitPointerLock();
            console.log('[THD-Controls] Pointer lock released with ESC');
            return;
        }
        
        // Handle Q/E rotation
        if (event.code === 'KeyQ' || event.code === 'KeyE') {
            event.preventDefault();
            event.stopPropagation();
            this.keys[event.code] = true;
            
            // Temporarily disable look-controls to prevent conflict
            if (this.lookControls) {
                this.lookControls.pause();
                console.log('[THD-Controls] Look-controls paused for Q/E rotation');
            }
            
            console.log('[THD-Controls] Q/E Key pressed:', event.code, 'keys state:', this.keys);
        }
    },

    onKeyUp: function (event) {
        if (event.code === 'KeyQ' || event.code === 'KeyE') {
            event.preventDefault();
            event.stopPropagation();
            this.keys[event.code] = false;
            
            // Re-enable look-controls when no Q/E keys are pressed
            if (!this.keys['KeyQ'] && !this.keys['KeyE'] && this.lookControls) {
                // Update look-controls internal rotation to match our current rotation
                const currentRotation = this.el.getAttribute('rotation');
                this.lookControls.yawObject.rotation.y = currentRotation.y * Math.PI / 180;
                this.lookControls.pitchObject.rotation.x = currentRotation.x * Math.PI / 180;
                
                this.lookControls.play();
                console.log('[THD-Controls] Look-controls resumed with synced rotation:', currentRotation.y);
            }
            
            console.log('[THD-Controls] Key released:', event.code);
        }
    },

    tick: function () {
        // Debug what's happening each tick
        if (this.keys['KeyQ'] || this.keys['KeyE']) {
            console.log('[THD-Controls] TICK - Keys active:', this.keys, 'pointerLock:', !!document.pointerLockElement);
        }
        
        // Rotate regardless of pointer lock for now (debugging)
        const el = this.el;
        const currentRotation = el.getAttribute('rotation');
        
        if (this.keys['KeyQ']) {
            // Turn left - increase Y rotation
            const newY = currentRotation.y + this.rotationSpeed;
            el.setAttribute('rotation', {
                x: currentRotation.x,
                y: newY,
                z: currentRotation.z
            });
            console.log('[THD-Controls] Turning left, new Y:', newY);
        }
        
        if (this.keys['KeyE']) {
            // Turn right - decrease Y rotation  
            const newY = currentRotation.y - this.rotationSpeed;
            el.setAttribute('rotation', {
                x: currentRotation.x,
                y: newY,
                z: currentRotation.z
            });
            console.log('[THD-Controls] Turning right, new Y:', newY);
        }
    },

    remove: function () {
        document.removeEventListener('keydown', this.onKeyDown, true);
        document.removeEventListener('keyup', this.onKeyUp, true);
    }
});

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
        // Set default grid transparency for holodeck environment
        this.data.gridTransparency = 0.15;
        this.createHolodeckGrid();
        console.log('[THD-AFrame] Holodeck coordinate system initialized with default grid');
    },
    
    createHolodeckGrid: function() {
        // Clear existing grid
        while (this.gridContainer.firstChild) {
            this.gridContainer.removeChild(this.gridContainer.firstChild);
        }
        
        if (this.data.gridTransparency === 0) {
            return; // Grid disabled
        }
        
        const gridColor = '#00ffff';
        const opacity = this.data.gridTransparency;
        const lineOpacity = opacity * 0.6;
        
        // Create floor grid points
        for (let x = -12; x <= 12; x += 4) {
            for (let z = -12; z <= 12; z += 4) {
                const gridPoint = document.createElement('a-sphere');
                gridPoint.setAttribute('position', `${x} 0.1 ${z}`);
                gridPoint.setAttribute('radius', '0.1');
                gridPoint.setAttribute('material', `color: ${gridColor}; transparent: true; opacity: ${opacity}`);
                this.gridContainer.appendChild(gridPoint);
            }
        }
        
        // Create horizontal connecting lines (X direction)
        for (let x = -12; x < 12; x += 4) {
            for (let z = -12; z <= 12; z += 4) {
                const lineX = document.createElement('a-box');
                lineX.setAttribute('position', `${x + 2} 0.1 ${z}`);
                lineX.setAttribute('scale', '4 0.02 0.02');
                lineX.setAttribute('material', `color: ${gridColor}; transparent: true; opacity: ${lineOpacity}`);
                this.gridContainer.appendChild(lineX);
            }
        }
        
        // Create horizontal connecting lines (Z direction)
        for (let x = -12; x <= 12; x += 4) {
            for (let z = -12; z < 12; z += 4) {
                const lineZ = document.createElement('a-box');
                lineZ.setAttribute('position', `${x} 0.1 ${z + 2}`);
                lineZ.setAttribute('scale', '0.02 0.02 4');
                lineZ.setAttribute('material', `color: ${gridColor}; transparent: true; opacity: ${lineOpacity}`);
                this.gridContainer.appendChild(lineZ);
            }
        }
        
        // Create ceiling grid points
        for (let x = -8; x <= 8; x += 8) {
            for (let z = -8; z <= 8; z += 8) {
                const ceilingPoint = document.createElement('a-sphere');
                ceilingPoint.setAttribute('position', `${x} 8 ${z}`);
                ceilingPoint.setAttribute('radius', '0.05');
                ceilingPoint.setAttribute('material', `color: ${gridColor}; transparent: true; opacity: ${opacity * 0.5}`);
                this.gridContainer.appendChild(ceilingPoint);
            }
        }
        
        // Create vertical connecting lines to ceiling
        for (let x = -8; x <= 8; x += 8) {
            for (let z = -8; z <= 8; z += 8) {
                const verticalLine = document.createElement('a-box');
                verticalLine.setAttribute('position', `${x} 4 ${z}`);
                verticalLine.setAttribute('scale', '0.01 8 0.01');
                verticalLine.setAttribute('material', `color: ${gridColor}; transparent: true; opacity: ${lineOpacity * 0.5}`);
                this.gridContainer.appendChild(verticalLine);
            }
        }
        
        // Create ceiling connecting lines
        const ceilingLineOpacity = lineOpacity * 0.3;
        
        // Ceiling X lines
        for (let z = -8; z <= 8; z += 8) {
            const ceilingLineX = document.createElement('a-box');
            ceilingLineX.setAttribute('position', `0 8 ${z}`);
            ceilingLineX.setAttribute('scale', '16 0.01 0.01');
            ceilingLineX.setAttribute('material', `color: ${gridColor}; transparent: true; opacity: ${ceilingLineOpacity}`);
            this.gridContainer.appendChild(ceilingLineX);
        }
        
        // Ceiling Z lines
        for (let x = -8; x <= 8; x += 8) {
            const ceilingLineZ = document.createElement('a-box');
            ceilingLineZ.setAttribute('position', `${x} 8 0`);
            ceilingLineZ.setAttribute('scale', '0.01 0.01 16');
            ceilingLineZ.setAttribute('material', `color: ${gridColor}; transparent: true; opacity: ${ceilingLineOpacity}`);
            this.gridContainer.appendChild(ceilingLineZ);
        }
        
        // Create invisible collision walls to contain player within holodeck
        this.createHolodeckBoundaries();
    },
    
    createHolodeckBoundaries: function() {
        const wallHeight = 8;
        const wallThickness = 0.5;
        const boundaryDistance = 11; // Closer to visible grid
        
        // North wall (negative Z) - RED with label
        const northWall = document.createElement('a-box');
        northWall.setAttribute('position', `0 ${wallHeight/2} ${-boundaryDistance}`);
        northWall.setAttribute('scale', `22 ${wallHeight} ${wallThickness}`);
        northWall.setAttribute('material', 'color: red');
        northWall.setAttribute('static-body', '');
        this.gridContainer.appendChild(northWall);
        
        const northText = document.createElement('a-text');
        northText.setAttribute('position', `0 ${wallHeight/2} ${-boundaryDistance + 0.3}`);
        northText.setAttribute('value', 'NORTH WALL - RED');
        northText.setAttribute('align', 'center');
        northText.setAttribute('color', 'white');
        northText.setAttribute('scale', '2 2 2');
        this.gridContainer.appendChild(northText);
        
        // South wall (positive Z) - GREEN with label
        const southWall = document.createElement('a-box');
        southWall.setAttribute('position', `0 ${wallHeight/2} ${boundaryDistance}`);
        southWall.setAttribute('scale', `22 ${wallHeight} ${wallThickness}`);
        southWall.setAttribute('material', 'color: green');
        southWall.setAttribute('static-body', '');
        this.gridContainer.appendChild(southWall);
        
        const southText = document.createElement('a-text');
        southText.setAttribute('position', `0 ${wallHeight/2} ${boundaryDistance - 0.3}`);
        southText.setAttribute('value', 'SOUTH WALL - GREEN');
        southText.setAttribute('align', 'center');
        southText.setAttribute('color', 'white');
        southText.setAttribute('scale', '2 2 2');
        southText.setAttribute('rotation', '0 180 0');
        this.gridContainer.appendChild(southText);
        
        // East wall (positive X) - BLUE with label
        const eastWall = document.createElement('a-box');
        eastWall.setAttribute('position', `${boundaryDistance} ${wallHeight/2} 0`);
        eastWall.setAttribute('scale', `${wallThickness} ${wallHeight} 22`);
        eastWall.setAttribute('material', 'color: blue');
        eastWall.setAttribute('static-body', '');
        this.gridContainer.appendChild(eastWall);
        
        const eastText = document.createElement('a-text');
        eastText.setAttribute('position', `${boundaryDistance - 0.3} ${wallHeight/2} 0`);
        eastText.setAttribute('value', 'EAST WALL - BLUE');
        eastText.setAttribute('align', 'center');
        eastText.setAttribute('color', 'white');
        eastText.setAttribute('scale', '2 2 2');
        eastText.setAttribute('rotation', '0 -90 0');
        this.gridContainer.appendChild(eastText);
        
        // West wall (negative X) - YELLOW with label
        const westWall = document.createElement('a-box');
        westWall.setAttribute('position', `${-boundaryDistance} ${wallHeight/2} 0`);
        westWall.setAttribute('scale', `${wallThickness} ${wallHeight} 22`);
        westWall.setAttribute('material', 'color: yellow');
        westWall.setAttribute('static-body', '');
        this.gridContainer.appendChild(westWall);
        
        const westText = document.createElement('a-text');
        westText.setAttribute('position', `${-boundaryDistance + 0.3} ${wallHeight/2} 0`);
        westText.setAttribute('value', 'WEST WALL - YELLOW');
        westText.setAttribute('align', 'center');
        westText.setAttribute('color', 'black');
        westText.setAttribute('scale', '2 2 2');
        westText.setAttribute('rotation', '0 90 0');
        this.gridContainer.appendChild(westText);
        
        // Add floor indicator (visual only, boundary checking handles containment)
        const floor = document.createElement('a-box');
        floor.setAttribute('position', '0 -0.5 0');
        floor.setAttribute('scale', '24 1 24');
        floor.setAttribute('material', 'color: #333; transparent: true; opacity: 0.1');
        floor.setAttribute('static-body', '');
        this.gridContainer.appendChild(floor);
        
        console.log('[THD-AFrame] Solid test barriers with labels and floor created');
    },
    
    createGridPlane: function(type, yPos, color, opacity) {
        // Create horizontal grid lines (floor/ceiling)
        for (let x = -12; x <= 12; x += 2) {
            for (let z = -12; z <= 12; z += 2) {
                // Grid intersection points
                const point = document.createElement('a-box');
                point.setAttribute('position', `${x} ${yPos} ${z}`);
                point.setAttribute('scale', '0.05 0.01 0.05');
                point.setAttribute('material', `color: ${color}; transparent: true; opacity: ${opacity}`);
                this.gridContainer.appendChild(point);
                
                // Horizontal lines
                if (x < 12) {
                    const lineX = document.createElement('a-box');
                    lineX.setAttribute('position', `${x + 1} ${yPos} ${z}`);
                    lineX.setAttribute('scale', '2 0.005 0.01');
                    lineX.setAttribute('material', `color: ${color}; transparent: true; opacity: ${opacity * 0.5}`);
                    this.gridContainer.appendChild(lineX);
                }
                
                // Vertical lines
                if (z < 12) {
                    const lineZ = document.createElement('a-box');
                    lineZ.setAttribute('position', `${x} ${yPos} ${z + 1}`);
                    lineZ.setAttribute('scale', '0.01 0.005 2');
                    lineZ.setAttribute('material', `color: ${color}; transparent: true; opacity: ${opacity * 0.5}`);
                    this.gridContainer.appendChild(lineZ);
                }
            }
        }
    },
    
    createGridWalls: function(color, opacity) {
        // Create vertical grid lines on walls
        const wallOpacity = opacity * 0.3;
        
        // North and South walls (Z = -12 and Z = 12)
        for (let wall of [-12, 12]) {
            for (let x = -12; x <= 12; x += 2) {
                for (let y = 0; y <= 12; y += 2) {
                    const point = document.createElement('a-box');
                    point.setAttribute('position', `${x} ${y} ${wall}`);
                    point.setAttribute('scale', '0.02 0.02 0.02');
                    point.setAttribute('material', `color: ${color}; transparent: true; opacity: ${wallOpacity}`);
                    this.gridContainer.appendChild(point);
                }
            }
        }
        
        // East and West walls (X = -12 and X = 12)
        for (let wall of [-12, 12]) {
            for (let z = -12; z <= 12; z += 2) {
                for (let y = 0; y <= 12; y += 2) {
                    const point = document.createElement('a-box');
                    point.setAttribute('position', `${wall} ${y} ${z}`);
                    point.setAttribute('scale', '0.02 0.02 0.02');
                    point.setAttribute('material', `color: ${color}; transparent: true; opacity: ${wallOpacity}`);
                    this.gridContainer.appendChild(point);
                }
            }
        }
    },
    
    updateGrid: function(transparency) {
        this.data.gridTransparency = transparency;
        this.createHolodeckGrid();
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
            case 'delete':
                this.deleteObject(message.object_name);
                break;
            case 'camera':
                this.updateCamera(message.camera);
                break;
            case 'update':
                this.updateObjects(message.objects);
                break;
            case 'session_created':
                console.log('[THD-AFrame] Session created:', message.data);
                break;
            case 'world_initialized':
                console.log('[THD-AFrame] World initialized:', message.data);
                if (message.data) {
                    this.initializeWorld(message.data);
                }
                break;
            case 'grid_control':
                console.log('[THD-AFrame] Grid control:', message.data);
                this.updateHolodeckGrid(message.data);
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
            case 'cylinder':
                entity.setAttribute('geometry', {
                    primitive: 'cylinder',
                    radius: 0.5,
                    height: 1
                });
                break;
            case 'cone':
                entity.setAttribute('geometry', {
                    primitive: 'cone',
                    radiusBottom: 0.5,
                    radiusTop: 0,
                    height: 1
                });
                break;
            case 'light':
                this.createLight(entity, obj);
                return; // Lights don't need geometry
            case 'sky':
                this.createSky(entity, obj);
                return; // Sky doesn't need geometry
            case 'text':
                this.createText(entity, obj);
                return; // Text uses special component
            case 'environment':
                this.createEnvironment(entity, obj);
                return; // Environment uses special component
            case 'particle':
                this.createParticles(entity, obj);
                return; // Particles use special component
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
        
        // Enhanced A-Frame material properties
        const materialProps = {
            shader: obj.material?.shader || 'standard',
            color: hexColor,
            metalness: obj.material?.metalness || 0.1,
            roughness: obj.material?.roughness || 0.7,
            transparent: obj.material?.transparent || color.a < 1.0,
            opacity: color.a
        };
        
        if (obj.wireframe) {
            materialProps.wireframe = true;
            materialProps.shader = 'flat'; // Wireframe works better with flat shader
        }
        
        entity.setAttribute('material', materialProps);
        
        // Add shadow properties if lighting data exists
        if (obj.lighting) {
            if (obj.lighting.castShadow) {
                entity.setAttribute('shadow', 'cast: true');
            }
            if (obj.lighting.receiveShadow) {
                entity.setAttribute('shadow', 'receive: true');
            }
        }
        
        // Add physics if enabled
        if (obj.physics && obj.physics.enabled) {
            if (obj.physics.type === 'static') {
                entity.setAttribute('static-body', {
                    shape: 'auto'
                });
            } else {
                entity.setAttribute('dynamic-body', {
                    shape: 'auto',
                    mass: obj.physics.mass || 1.0
                });
            }
        }
    }

    createLight(entity, obj) {
        // Enhanced lighting system
        const lightProps = {
            type: obj.lightType || 'point',
            color: obj.color ? this.colorToHex(obj.color) : '#ffffff',
            intensity: obj.intensity || 1.0
        };
        
        // Add directional light specific properties
        if (obj.lightType === 'directional') {
            lightProps.castShadow = true;
            lightProps.shadowMapWidth = 1024;
            lightProps.shadowMapHeight = 1024;
        }
        
        entity.setAttribute('light', lightProps);
        console.log('[THD-AFrame] Created light:', lightProps);
    }

    createSky(entity, obj) {
        // Create sky dome/environment
        const color = obj.color ? this.colorToHex(obj.color) : '#87CEEB';
        
        // Use A-Frame sky primitive
        entity.setAttribute('geometry', {primitive: 'sphere', radius: 5000});
        entity.setAttribute('material', {
            shader: 'flat',
            color: color,
            side: 'back'
        });
        entity.setAttribute('scale', '-1 1 1'); // Invert to see from inside
        
        console.log('[THD-AFrame] Created sky environment:', color);
    }

    createText(entity, obj) {
        // Create 3D text in holodeck space
        const textProps = {
            value: obj.text || 'Holodeck Text',
            color: obj.color ? this.colorToHex(obj.color) : '#ffffff',
            align: 'center',
            width: obj.width || 6,
            font: 'roboto'
        };
        
        entity.setAttribute('text', textProps);
        console.log('[THD-AFrame] Created 3D text:', obj.text);
    }

    createEnvironment(entity, obj) {
        // Create procedural environments using A-Frame environment component
        const envType = obj.envType || 'forest';
        const envProps = {
            preset: envType,
            groundColor: obj.groundColor || '#553e35',
            grid: '1x1',
            gridColor: obj.gridColor || '#ccc',
            groundTexture: 'none',
            skyType: 'atmosphere',
            lighting: 'distant'
        };
        
        // Override with custom environment settings if provided
        if (obj.environmentSettings) {
            Object.assign(envProps, obj.environmentSettings);
        }
        
        entity.setAttribute('environment', envProps);
        console.log('[THD-AFrame] Created environment:', envType, envProps);
    }

    createParticles(entity, obj) {
        // Create particle effects (fire, smoke, sparkles, etc.)
        const particleType = obj.particleType || 'sparkle';
        
        switch (particleType) {
            case 'fire':
                entity.setAttribute('particle-system', {
                    preset: 'fire',
                    particleCount: obj.count || 1000,
                    maxAge: 2,
                    color: '#ff6600,#ff0000,#ffaa00',
                    size: '0.5, 0'
                });
                break;
            case 'smoke':
                entity.setAttribute('particle-system', {
                    preset: 'smoke',
                    particleCount: obj.count || 500,
                    maxAge: 4,
                    color: '#888888,#aaaaaa,#cccccc',
                    size: '1, 3'
                });
                break;
            case 'sparkle':
            default:
                entity.setAttribute('particle-system', {
                    preset: 'sparkle',
                    particleCount: obj.count || 200,
                    maxAge: 3,
                    color: '#ffffff,#00ffff,#ffff00',
                    size: '0.1, 0.5'
                });
                break;
        }
        
        console.log('[THD-AFrame] Created particle system:', particleType);
    }

    deleteObject(objectName) {
        // Delete specific object by name
        const entity = this.objects.get(objectName);
        if (entity) {
            if (entity.parentNode) {
                entity.parentNode.removeChild(entity);
            }
            this.objects.delete(objectName);
            console.log('[THD-AFrame] Deleted object:', objectName);
        } else {
            console.warn('[THD-AFrame] Object not found for deletion:', objectName);
        }
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

    // Get current object count for status display
    getObjectCount() {
        return this.objects.size;
    }

    updateHolodeckGrid(gridData) {
        const sceneEl = document.getElementById('holodeck-scene');
        const holodeckComponent = sceneEl.components['thd-holodeck'];
        
        if (holodeckComponent && gridData.transparency !== undefined) {
            holodeckComponent.updateGrid(gridData.transparency);
        }
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
console.log('[THD-AFrame] THDAFrameManager class defined:', typeof THDAFrameManager);
console.log('[THD-AFrame] Registered components:', Object.keys(AFRAME.components));
console.log('[THD-AFrame] holodeck-boundaries registered:', !!AFRAME.components['holodeck-boundaries']);

// Check if camera element exists and has the component, if not add it manually
setTimeout(() => {
    const camera = document.getElementById('holodeck-camera');
    console.log('[THD-AFrame] Camera element found:', !!camera);
    if (camera) {
        console.log('[THD-AFrame] Camera components:', Object.keys(camera.components || {}));
        console.log('[THD-AFrame] Has holodeck-boundaries component:', !!camera.components['holodeck-boundaries']);
        
        // If component not attached, force attach it
        if (!camera.components['holodeck-boundaries']) {
            console.log('[THD-AFrame] Manually attaching holodeck-boundaries component...');
            camera.setAttribute('holodeck-boundaries', '');
            
            // Check again after 1 second
            setTimeout(() => {
                console.log('[THD-AFrame] After manual attach - Has holodeck-boundaries:', !!camera.components['holodeck-boundaries']);
            }, 1000);
        }
    }
}, 2000);

// Global export to ensure it's available
window.THDAFrameManager = THDAFrameManager;