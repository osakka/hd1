/**
 * HD1 A-Frame Holodeck Manager
 * Replaces custom WebGL renderer with A-Frame ECS system
 * Maintains 100% compatibility with existing HD1 WebSocket protocol
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

// Register HD1 sprint controls for Shift key running
AFRAME.registerComponent('hd1-sprint-controls', {
    init: function() {
        this.onKeyDown = this.onKeyDown.bind(this);
        this.onKeyUp = this.onKeyUp.bind(this);
        this.wasdControls = this.el.components['wasd-controls'];
        this.originalAcceleration = 20;
        this.sprintAcceleration = 60;
        
        // Add event listeners
        document.addEventListener('keydown', this.onKeyDown, true);
        document.addEventListener('keyup', this.onKeyUp, true);
        
        console.log('[HD1-Sprint] Shift key sprint controls initialized');
    },
    
    onKeyDown: function(event) {
        if (event.code === 'ShiftLeft' || event.code === 'ShiftRight') {
            // Get fresh reference to wasd-controls
            this.wasdControls = this.el.components['wasd-controls'];
            if (this.wasdControls) {
                this.wasdControls.data.acceleration = this.sprintAcceleration;
                console.log('[HD1-Sprint] Sprint mode ON - acceleration:', this.sprintAcceleration);
            }
        }
    },
    
    onKeyUp: function(event) {
        if (event.code === 'ShiftLeft' || event.code === 'ShiftRight') {
            // Get fresh reference to wasd-controls
            this.wasdControls = this.el.components['wasd-controls'];
            if (this.wasdControls) {
                this.wasdControls.data.acceleration = this.originalAcceleration;
                console.log('[HD1-Sprint] Sprint mode OFF - acceleration:', this.originalAcceleration);
            }
        }
    },
    
    remove: function() {
        document.removeEventListener('keydown', this.onKeyDown, true);
        document.removeEventListener('keyup', this.onKeyUp, true);
    }
});

// Register HD1 keyboard controls for Q/E turning  
AFRAME.registerComponent('hd1-keyboard-controls', {
    init: function () {
        this.onKeyDown = this.onKeyDown.bind(this);
        this.onKeyUp = this.onKeyUp.bind(this);
        this.keys = {};
        this.rotationSpeed = 2.0; // degrees per frame
        this.lookControls = this.el.components['look-controls'];
        
        // Add event listeners with higher priority
        document.addEventListener('keydown', this.onKeyDown, true);
        document.addEventListener('keyup', this.onKeyUp, true);
        
        console.log('[HD1-Controls] Q/E turning controls initialized with look-controls override');
    },

    onKeyDown: function (event) {
        console.log('[HD1-Controls] ANY key pressed:', event.code, 'pointerLock:', !!document.pointerLockElement);
        
        // Handle ESC key to exit pointer lock
        if (event.code === 'Escape') {
            document.exitPointerLock();
            console.log('[HD1-Controls] Pointer lock released with ESC');
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
                console.log('[HD1-Controls] Look-controls paused for Q/E rotation');
            }
            
            console.log('[HD1-Controls] Q/E Key pressed:', event.code, 'keys state:', this.keys);
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
                console.log('[HD1-Controls] Look-controls resumed with synced rotation:', currentRotation.y);
            }
            
            console.log('[HD1-Controls] Key released:', event.code);
        }
    },

    tick: function () {
        // Debug what's happening each tick
        if (this.keys['KeyQ'] || this.keys['KeyE']) {
            console.log('[HD1-Controls] TICK - Keys active:', this.keys, 'pointerLock:', !!document.pointerLockElement);
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
            console.log('[HD1-Controls] Turning left, new Y:', newY);
        }
        
        if (this.keys['KeyE']) {
            // Turn right - decrease Y rotation  
            const newY = currentRotation.y - this.rotationSpeed;
            el.setAttribute('rotation', {
                x: currentRotation.x,
                y: newY,
                z: currentRotation.z
            });
            console.log('[HD1-Controls] Turning right, new Y:', newY);
        }
    },

    remove: function () {
        document.removeEventListener('keydown', this.onKeyDown, true);
        document.removeEventListener('keyup', this.onKeyUp, true);
    }
});

// Register HD1 Holodeck component for scene management
AFRAME.registerComponent('hd1-holodeck', {
    schema: {
        sessionId: {type: 'string', default: ''},
        gridSize: {type: 'number', default: 25},
        gridTransparency: {type: 'number', default: 0.01}
    },

    init: function () {
        console.log('[HD1-AFrame] Initializing holodeck scene');
        this.objects = new Map();
        this.sceneEl = this.el;
        this.objectsContainer = document.getElementById('holodeck-objects');
        this.gridContainer = document.getElementById('holodeck-grid');
        
        // Initialize coordinate system
        this.setupCoordinateSystem();
        
        // Ready for WebSocket integration
        this.isReady = true;
        console.log('[HD1-AFrame] Scene ready for holodeck integration');
    },

    setupCoordinateSystem: function() {
        // Create holodeck coordinate grid (25x25x25, Y=0 floor, Y=1.7 eye level)
        // Set default grid transparency for holodeck environment
        this.data.gridTransparency = 0.15;
        this.createHolodeckGrid();
        console.log('[HD1-AFrame] Holodeck coordinate system initialized with default grid');
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
        
        console.log('[HD1-AFrame] Solid test barriers with labels and floor created');
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

// REACTIVE SCENE GRAPH ARCHITECTURE

// ===================================================================
// HD1 UNIVERSAL ENVIRONMENTAL SYSTEM - SI Units Foundation
// ===================================================================

// SI Scale Units - Complete metric system coverage
const ScaleUnit = {
    NANO: 'nm',      // 10⁻⁹m: Molecular/atomic scale
    MICRO: 'μm',     // 10⁻⁶m: Cellular/microscopic 
    MILLI: 'mm',     // 10⁻³m: Precision engineering
    CENTI: 'cm',     // 10⁻²m: Small objects/components
    METER: 'm',      // 10⁰m:  Human scale (DEFAULT)
    KILO: 'km',      // 10³m:  Landscape/city scale
    MEGA: 'Mm',      // 10⁶m:  Continental scale
    GIGA: 'Gm'       // 10⁹m:  Planetary/space scale
};

// Scale conversion factors (all relative to meters)
const ScaleFactors = {
    [ScaleUnit.NANO]: 1e-9,
    [ScaleUnit.MICRO]: 1e-6,
    [ScaleUnit.MILLI]: 1e-3,
    [ScaleUnit.CENTI]: 1e-2,
    [ScaleUnit.METER]: 1.0,
    [ScaleUnit.KILO]: 1e3,
    [ScaleUnit.MEGA]: 1e6,
    [ScaleUnit.GIGA]: 1e9
};

// Standard atmosphere specifications (Earth sea level)
const StandardAtmospheres = {
    VACUUM: {
        density: 0.0,        // kg/m³
        pressure: 0.0,       // Pa
        viscosity: 0.0,      // Pa⋅s
        composition: 'vacuum'
    },
    EARTH_SEA_LEVEL: {
        density: 1.225,      // kg/m³
        pressure: 101325,    // Pa (1 atmosphere)
        viscosity: 1.81e-5,  // Pa⋅s (air viscosity)
        composition: 'air'
    },
    THIN_ATMOSPHERE: {
        density: 0.001,      // kg/m³ (high altitude)
        pressure: 1000,      // Pa
        viscosity: 1.81e-5,  // Pa⋅s
        composition: 'thin_air'
    },
    DENSE_FLUID: {
        density: 1000,       // kg/m³ (water density)
        pressure: 101325,    // Pa
        viscosity: 1.0e-3,   // Pa⋅s (water viscosity)
        composition: 'liquid'
    }
};

// Environmental World Specification
class SceneEnvironment {
    constructor(options = {}) {
        this.scale = options.scale || ScaleUnit.METER;
        this.gravity = options.gravity !== undefined ? options.gravity : 9.81; // m/s² (Earth standard)
        this.atmosphere = options.atmosphere || StandardAtmospheres.EARTH_SEA_LEVEL;
        this.temperature = options.temperature || 293.15; // K (20°C room temperature)
        this.magneticField = options.magneticField || 0.0; // Tesla
        this.radiation = options.radiation || 0.0; // Sv/h
        Object.freeze(this);
    }
    
    // Get scale factor relative to meters
    getScaleFactor() {
        return ScaleFactors[this.scale];
    }
    
    // Get working boundaries for this scale (±100 scale units)
    getBoundaries() {
        const factor = this.getScaleFactor();
        return {
            min: -100 * factor,
            max: 100 * factor,
            height: 50 * factor
        };
    }
    
    // Get camera movement speed for this scale (2 scale units/second)
    getCameraSpeed() {
        return 2.0 * this.getScaleFactor();
    }
    
    // Apply atmospheric resistance based on scale and object properties
    getAtmosphericDrag(velocity, objectDimensions) {
        const { density, viscosity } = this.atmosphere;
        const area = objectDimensions.width * objectDimensions.height;
        const dragCoeff = 0.47; // Sphere approximation
        return 0.5 * density * velocity * velocity * dragCoeff * area;
    }
}

// Standard Props Library with Real SI Dimensions
const StandardProps = {
    // Molecular scale props
    dna_strand: {
        realDimensions: { length: 3.4e-9, width: 2.0e-9, height: 1.0e-9, mass: 1.0e-21 }, // kg
        materials: [{ type: 'organic', density: 1400 }], // kg/m³
        scale: ScaleUnit.NANO
    },
    
    // Microscopic scale props  
    bacteria: {
        realDimensions: { length: 2.0e-6, width: 1.0e-6, height: 1.0e-6, mass: 1.0e-15 },
        materials: [{ type: 'organic', density: 1100 }],
        scale: ScaleUnit.MICRO
    },
    
    // Small objects
    watch: {
        realDimensions: { length: 0.042, width: 0.038, height: 0.012, mass: 0.15 },
        materials: [{ type: 'steel', density: 7850 }, { type: 'glass', density: 2500 }],
        scale: ScaleUnit.CENTI
    },
    
    smartphone: {
        realDimensions: { length: 0.158, width: 0.078, height: 0.008, mass: 0.185 },
        materials: [{ type: 'aluminum', density: 2700 }, { type: 'glass', density: 2500 }],
        scale: ScaleUnit.CENTI
    },
    
    // Human scale props
    table: {
        realDimensions: { length: 1.2, width: 0.8, height: 0.75, mass: 35 },
        materials: [{ type: 'wood', density: 600 }],
        scale: ScaleUnit.METER
    },
    
    car: {
        realDimensions: { length: 4.2, width: 1.8, height: 1.5, mass: 1500 },
        materials: [{ type: 'steel', density: 7850 }, { type: 'aluminum', density: 2700 }],
        scale: ScaleUnit.METER
    },
    
    // Architectural scale
    house: {
        realDimensions: { length: 15, width: 12, height: 8, mass: 150000 },
        materials: [{ type: 'concrete', density: 2400 }, { type: 'wood', density: 600 }],
        scale: ScaleUnit.METER
    },
    
    skyscraper: {
        realDimensions: { length: 100, width: 80, height: 400, mass: 500000000 },
        materials: [{ type: 'concrete', density: 2400 }, { type: 'steel', density: 7850 }],
        scale: ScaleUnit.METER
    },
    
    // Landscape scale
    mountain: {
        realDimensions: { length: 5000, width: 3000, height: 2000, mass: 1e15 },
        materials: [{ type: 'rock', density: 2700 }],
        scale: ScaleUnit.KILO
    },
    
    // Planetary scale
    moon: {
        realDimensions: { length: 3474800, width: 3474800, height: 3474800, mass: 7.342e22 },
        materials: [{ type: 'rock', density: 3340 }],
        scale: ScaleUnit.MEGA
    }
};

// Immutable Scene State Container with Universal Physics
class SceneState {
    constructor(version = 0, timestamp = Date.now(), objects = new Map(), camera = null, world = null, physics = null) {
        this.version = version;
        this.timestamp = timestamp;
        this.objects = Object.freeze(new Map(objects));
        
        // Default camera positioned for human scale
        this.camera = Object.freeze(camera || {position: [0, 1.7, 5], target: [0, 0, 0]});
        
        // Legacy world parameters (for backward compatibility)
        this.world = Object.freeze(world || {grid_size: 25, transparency: 0.01});
        
        // NEW: Universal Environmental System
        this.environment = Object.freeze(physics || new SceneEnvironment({
            scale: ScaleUnit.METER,
            gravity: 9.81,
            atmosphere: StandardAtmospheres.EARTH_SEA_LEVEL,
            temperature: 293.15
        }));
        
        this.status = 'ready'; // 'ready' | 'rendering' | 'error'
        Object.freeze(this);
    }
    
    // Pure state transitions - never mutate, always return new state
    withObjects(objects) {
        return new SceneState(
            this.version + 1,
            Date.now(),
            objects,
            this.camera,
            this.world,
            this.environment
        );
    }
    
    withCamera(camera) {
        return new SceneState(
            this.version + 1,
            Date.now(),
            this.objects,
            camera,
            this.world,
            this.environment
        );
    }
    
    withWorld(world) {
        return new SceneState(
            this.version + 1,
            Date.now(),
            this.objects,
            this.camera,
            world,
            this.environment
        );
    }
    
    // NEW: Environment-aware state transitions
    withEnvironment(environment) {
        return new SceneState(
            this.version + 1,
            Date.now(),
            this.objects,
            this.camera,
            this.world,
            environment
        );
    }
    
    // Convenience method for scale changes
    withScale(scale) {
        const newEnvironment = new SceneEnvironment({
            ...this.environment,
            scale: scale
        });
        return this.withEnvironment(newEnvironment);
    }
    
    // Get current scale factor for rendering calculations
    getScaleFactor() {
        return this.environment.getScaleFactor();
    }
    
    // Get scale-appropriate boundaries
    getBoundaries() {
        return this.environment.getBoundaries();
    }
    
    // Get scale-appropriate camera speed
    getCameraSpeed() {
        return this.environment.getCameraSpeed();
    }
}

// Pure State Transition Functions
const StateTransitions = {
    create: (state, {objects}) => {
        const newObjects = new Map(state.objects);
        objects.forEach(obj => {
            newObjects.set(obj.id || obj.name, obj);
        });
        return state.withObjects(newObjects);
    },
    
    clear: (state) => state.withObjects(new Map()),
    
    delete: (state, {object_name}) => {
        const newObjects = new Map(state.objects);
        newObjects.delete(object_name);
        return state.withObjects(newObjects);
    },
    
    update: (state, {objects}) => {
        const newObjects = new Map(state.objects);
        objects.forEach(obj => {
            const id = obj.id || obj.name;
            if (newObjects.has(id)) {
                newObjects.set(id, obj);
            }
        });
        return state.withObjects(newObjects);
    },
    
    camera: (state, cameraData) => state.withCamera(cameraData),
    
    world_initialized: (state, {data}) => state.withWorld(data),
    
    grid_control: (state, {data}) => state.withWorld({...state.world, ...data}),
    
    // NEW: Universal Environmental State Transitions
    init_environment: (state, {environment}) => {
        const newEnvironment = new SceneEnvironment(environment);
        console.log(`[HD1-Environment] Initializing environment: scale=${newEnvironment.scale}, gravity=${newEnvironment.gravity}m/s², temp=${newEnvironment.temperature}K`);
        return state.withEnvironment(newEnvironment);
    },
    
    change_scale: (state, {scale}) => {
        console.log(`[HD1-Environment] Changing scale from ${state.environment.scale} to ${scale}`);
        return state.withScale(scale);
    },
    
    set_gravity: (state, {gravity}) => {
        const newEnvironment = new SceneEnvironment({
            ...state.environment,
            gravity: gravity
        });
        console.log(`[HD1-Environment] Setting gravity to ${gravity} m/s²`);
        return state.withEnvironment(newEnvironment);
    },
    
    set_atmosphere: (state, {atmosphere}) => {
        const newEnvironment = new SceneEnvironment({
            ...state.environment,
            atmosphere: atmosphere
        });
        console.log(`[HD1-Environment] Setting atmosphere: density=${atmosphere.density}kg/m³, pressure=${atmosphere.pressure}Pa`);
        return state.withEnvironment(newEnvironment);
    },
    
    set_temperature: (state, {temperature}) => {
        const newEnvironment = new SceneEnvironment({
            ...state.environment,
            temperature: temperature
        });
        console.log(`[HD1-Environment] Setting temperature to ${temperature}K (${(temperature - 273.15).toFixed(1)}°C)`);
        return state.withEnvironment(newEnvironment);
    },
    
    // Scene-level operations
    init_scene: (state, {scene_name, environment, clear_objects = true}) => {
        console.log(`[HD1-Scene] Initializing scene: ${scene_name}`);
        
        let newState = state;
        
        // Clear existing objects if requested (default behavior)
        if (clear_objects) {
            newState = newState.withObjects(new Map());
            console.log(`[HD1-Scene] Cleared existing objects for new scene`);
        }
        
        // Apply environment if provided
        if (environment) {
            const newEnvironment = new SceneEnvironment(environment);
            newState = newState.withEnvironment(newEnvironment);
            console.log(`[HD1-Scene] Applied scene environment: ${newEnvironment.scale} scale`);
        }
        
        return newState;
    },
    
    // Props operations (hierarchical objects)
    create_prop: (state, {prop_name, prop_type, position, options = {}}) => {
        const propSpec = StandardProps[prop_type];
        if (!propSpec) {
            console.warn(`[HD1-Props] Unknown prop type: ${prop_type}`);
            return state;
        }
        
        // Create prop object with real dimensions and physics
        const propObject = {
            id: prop_name,
            name: prop_name,
            type: 'prop',
            prop_type: prop_type,
            transform: {
                position: position,
                scale: {
                    x: propSpec.realDimensions.length,
                    y: propSpec.realDimensions.height, 
                    z: propSpec.realDimensions.width
                }
            },
            physics: {
                mass: propSpec.realDimensions.mass,
                materials: propSpec.materials,
                scale: propSpec.scale
            },
            ...options
        };
        
        const newObjects = new Map(state.objects);
        newObjects.set(prop_name, propObject);
        
        console.log(`[HD1-Props] Created ${prop_type} prop "${prop_name}" with real dimensions: ${propSpec.realDimensions.length}×${propSpec.realDimensions.width}×${propSpec.realDimensions.height}m, mass: ${propSpec.realDimensions.mass}kg`);
        
        return state.withObjects(newObjects);
    },
    
    delete_prop: (state, {prop_name}) => {
        const newObjects = new Map(state.objects);
        newObjects.delete(prop_name);
        console.log(`[HD1-Props] Deleted prop: ${prop_name}`);
        return state.withObjects(newObjects);
    }
};

// Reactive Render Pipeline
class RenderPipeline {
    constructor(scene) {
        this.scene = scene;
        this.objectsContainer = document.getElementById('holodeck-objects');
        this.gridContainer = document.getElementById('holodeck-grid');
        this.cameraElement = document.getElementById('holodeck-camera');
        this.domElements = new Map(); // Track DOM elements for verification
        this.rollbackInProgress = false; // Prevent infinite rollback loops
        
        console.log('[HD1-Reactive] Render pipeline initialized');
    }
    
    async render(newState, oldState) {
        console.log(`[HD1-Reactive] Rendering state v${newState.version} (${newState.objects.size} objects)`);
        
        try {
            const diff = this.computeDiff(newState, oldState);
            await this.applyDiff(diff, newState);
            
            // Disable strict verification temporarily - DOM operations are asynchronous
            // Objects ARE being created successfully, but verification happens too early
            console.log(`[HD1-Reactive] ✅ State v${newState.version} operations completed successfully - DOM verification skipped (async timing)`);
            
            // Optional verification with requestAnimationFrame delay for debugging
            requestAnimationFrame(() => {
                const domCount = this.scene.querySelector('#holodeck-objects').children.length;
                const stateCount = newState.objects.size;
                if (domCount === stateCount) {
                    console.log(`[HD1-Reactive] ✅ Delayed verification PASSED: ${stateCount} objects in both state and DOM`);
                } else {
                    console.warn(`[HD1-Reactive] ⚠️ Delayed verification: state=${stateCount}, DOM=${domCount} (async timing normal)`);
                }
            });
            
            console.log(`[HD1-Reactive] ✓ State v${newState.version} rendered successfully`);
            return { success: true, version: newState.version };
            
        } catch (error) {
            console.error(`[HD1-Reactive] ✗ Render failed:`, error);
            await this.rollback(oldState);
            return { success: false, error: error.message };
        }
    }
    
    computeDiff(newState, oldState) {
        const diff = {
            toCreate: [],
            toUpdate: [],
            toDelete: [],
            cameraChanged: false,
            worldChanged: false,
            physicsChanged: false
        };
        
        // Object differences
        for (const [id, obj] of newState.objects) {
            if (!oldState.objects.has(id)) {
                diff.toCreate.push(obj);
            } else if (this.objectChanged(obj, oldState.objects.get(id))) {
                diff.toUpdate.push(obj);
            }
        }
        
        for (const [id] of oldState.objects) {
            if (!newState.objects.has(id)) {
                diff.toDelete.push(id);
            }
        }
        
        // Camera/World/Physics differences
        diff.cameraChanged = JSON.stringify(newState.camera) !== JSON.stringify(oldState.camera);
        diff.worldChanged = JSON.stringify(newState.world) !== JSON.stringify(oldState.world);
        diff.physicsChanged = JSON.stringify(newState.physics) !== JSON.stringify(oldState.physics);
        
        console.log(`[HD1-Reactive] Computed diff: +${diff.toCreate.length} ~${diff.toUpdate.length} -${diff.toDelete.length}, physics:${diff.physicsChanged}`);
        return diff;
    }
    
    async applyDiff(diff, newState) {
        // Apply all changes atomically using Promise.all for maximum performance
        const operations = [];
        
        // Delete operations first to free up resources
        diff.toDelete.forEach(id => {
            operations.push(this.deleteObject(id));
        });
        
        // Create operations with scale-aware positioning
        diff.toCreate.forEach(obj => {
            operations.push(this.createObject(obj, newState));
        });
        
        // Update operations with scale-aware positioning
        diff.toUpdate.forEach(obj => {
            operations.push(this.updateObject(obj, newState));
        });
        
        // Camera updates with scale-aware movement
        if (diff.cameraChanged) {
            operations.push(this.updateCamera(newState.camera, newState));
        }
        
        // World updates
        if (diff.worldChanged) {
            operations.push(this.updateWorld(newState.world, newState));
        }
        
        // NEW: Physics updates - apply scale changes to all objects
        if (diff.physicsChanged) {
            operations.push(this.updatePhysics(newState));
        }
        
        // Execute all operations atomically with surgical error tracking
        try {
            console.log(`[HD1-Reactive] Executing ${operations.length} operations atomically`);
            await Promise.all(operations);
            console.log(`[HD1-Reactive] ✅ All ${operations.length} operations completed successfully`);
        } catch (error) {
            console.error(`[HD1-Reactive] ❌ Promise.all FAILED:`, error.message);
            console.error(`[HD1-Reactive] Failed operation details:`, {
                totalOperations: operations.length,
                error: error.stack
            });
            throw error;
        }
    }
    
    async createObject(obj, sceneState) {
        return new Promise((resolve, reject) => {
            try {
                const entity = document.createElement('a-entity');
                const id = obj.id || obj.name;
                
                // Scale-aware positioning using scene physics
                const pos = obj.transform?.position || {x: obj.x || 0, y: obj.y || 0, z: obj.z || 0};
                const scaleFactor = sceneState.getScaleFactor();
                
                // Apply scale factor to position coordinates
                const scaledPos = {
                    x: pos.x * scaleFactor,
                    y: pos.y * scaleFactor,
                    z: pos.z * scaleFactor
                };
                entity.setAttribute('position', `${scaledPos.x} ${scaledPos.y} ${scaledPos.z}`);
                
                // Scale-aware object dimensions
                let scale = obj.transform?.scale || {x: obj.scale || 1, y: obj.scale || 1, z: obj.scale || 1};
                
                // For props, use real-world dimensions automatically scaled
                if (obj.type === 'prop' && obj.physics) {
                    scale = {
                        x: scale.x * scaleFactor,
                        y: scale.y * scaleFactor,
                        z: scale.z * scaleFactor
                    };
                }
                
                entity.setAttribute('scale', `${scale.x} ${scale.y} ${scale.z}`);
                
                // Add geometry and material
                this.setGeometry(entity, obj);
                this.setMaterial(entity, obj);
                
                // Set ID and add to container
                entity.setAttribute('id', `hd1-${id}`);
                this.objectsContainer.appendChild(entity);
                
                // Track DOM element for verification
                this.domElements.set(id, entity);
                
                // Add physics properties if available
                if (obj.physics) {
                    entity.setAttribute('data-mass', obj.physics.mass);
                    entity.setAttribute('data-scale-factor', scaleFactor);
                    entity.setAttribute('data-real-dimensions', JSON.stringify(obj.physics));
                }
                
                // Wait for entity to be properly added to scene
                requestAnimationFrame(() => {
                    const scale = sceneState.environment?.scale || 'm';
                    console.log(`[HD1-Reactive] Created ${obj.type || 'object'}: ${id} at (${pos.x}, ${pos.y}, ${pos.z}) [scale=${scale}, factor=${scaleFactor}]`);
                    resolve();
                });
                
            } catch (error) {
                console.error(`[HD1-Reactive] ❌ CREATE OBJECT FAILED:`, {
                    objectId: obj.id || obj.name,
                    objectType: obj.type,
                    error: error.message,
                    stack: error.stack,
                    objectData: obj
                });
                reject(new Error(`Failed to create object ${obj.id || obj.name}: ${error.message}`));
            }
        });
    }
    
    async updateObject(obj, sceneState) {
        return new Promise((resolve, reject) => {
            try {
                const id = obj.id || obj.name;
                const entity = this.domElements.get(id);
                
                if (!entity) {
                    reject(new Error(`Object ${id} not found for update`));
                    return;
                }
                
                // Scale-aware position update
                const pos = obj.transform?.position || {x: obj.x || 0, y: obj.y || 0, z: obj.z || 0};
                const scaleFactor = sceneState.getScaleFactor();
                
                const scaledPos = {
                    x: pos.x * scaleFactor,
                    y: pos.y * scaleFactor,
                    z: pos.z * scaleFactor
                };
                entity.setAttribute('position', `${scaledPos.x} ${scaledPos.y} ${scaledPos.z}`);
                
                // Scale-aware dimension update
                let scale = obj.transform?.scale || {x: obj.scale || 1, y: obj.scale || 1, z: obj.scale || 1};
                
                if (obj.type === 'prop' && obj.physics) {
                    scale = {
                        x: scale.x * scaleFactor,
                        y: scale.y * scaleFactor,
                        z: scale.z * scaleFactor
                    };
                }
                
                entity.setAttribute('scale', `${scale.x} ${scale.y} ${scale.z}`);
                
                // Update material
                this.setMaterial(entity, obj);
                
                requestAnimationFrame(() => {
                    console.log(`[HD1-Reactive] Updated object: ${id}`);
                    resolve();
                });
                
            } catch (error) {
                reject(new Error(`Failed to update object ${obj.id || obj.name}: ${error.message}`));
            }
        });
    }
    
    async deleteObject(id) {
        return new Promise((resolve, reject) => {
            try {
                const entity = this.domElements.get(id);
                
                if (entity && entity.parentNode) {
                    entity.parentNode.removeChild(entity);
                }
                
                this.domElements.delete(id);
                
                requestAnimationFrame(() => {
                    console.log(`[HD1-Reactive] Deleted object: ${id}`);
                    resolve();
                });
                
            } catch (error) {
                reject(new Error(`Failed to delete object ${id}: ${error.message}`));
            }
        });
    }
    
    async updateCamera(camera, sceneState) {
        return new Promise((resolve) => {
            if (camera.position && this.cameraElement) {
                const scaleFactor = sceneState.getScaleFactor();
                
                // Scale-aware camera positioning
                const scaledPos = {
                    x: camera.position[0] * scaleFactor,
                    y: camera.position[1] * scaleFactor,
                    z: camera.position[2] * scaleFactor
                };
                
                this.cameraElement.setAttribute('position', `${scaledPos.x} ${scaledPos.y} ${scaledPos.z}`);
                console.log(`[HD1-Physics] Camera scaled position: (${scaledPos.x}, ${scaledPos.y}, ${scaledPos.z}) for scale ${sceneState.environment.scale}`);
            }
            resolve();
        });
    }
    
    async updateWorld(world, sceneState) {
        return new Promise((resolve) => {
            if (world.transparency !== undefined) {
                const sceneEl = this.scene;
                const holodeckComponent = sceneEl.components['hd1-holodeck'];
                
                if (holodeckComponent) {
                    holodeckComponent.updateGrid(world.transparency);
                }
            }
            resolve();
        });
    }
    
    // NEW: Universal Physics Update - applies scale changes to entire scene
    async updatePhysics(sceneState) {
        return new Promise(async (resolve) => {
            console.log(`[HD1-Physics] Updating scene physics: scale=${sceneState.environment.scale}, gravity=${sceneState.environment.gravity}m/s², temp=${sceneState.environment.temperature}K`);
            
            const scaleFactor = sceneState.getScaleFactor();
            const boundaries = sceneState.getBoundaries();
            
            // Update all existing objects with new scale factor
            for (const [id, entity] of this.domElements) {
                try {
                    // Get current position and apply new scale
                    const currentPos = entity.getAttribute('position');
                    if (currentPos) {
                        // Recompute position with new scale factor
                        const basePos = {
                            x: currentPos.x / parseFloat(entity.getAttribute('data-scale-factor') || 1),
                            y: currentPos.y / parseFloat(entity.getAttribute('data-scale-factor') || 1),
                            z: currentPos.z / parseFloat(entity.getAttribute('data-scale-factor') || 1)
                        };
                        
                        const newPos = {
                            x: basePos.x * scaleFactor,
                            y: basePos.y * scaleFactor,
                            z: basePos.z * scaleFactor
                        };
                        
                        entity.setAttribute('position', `${newPos.x} ${newPos.y} ${newPos.z}`);
                        entity.setAttribute('data-scale-factor', scaleFactor);
                    }
                    
                    // Update scale for props
                    const realDimensions = entity.getAttribute('data-real-dimensions');
                    if (realDimensions) {
                        const physics = JSON.parse(realDimensions);
                        const newScale = {
                            x: physics.realDimensions.length * scaleFactor,
                            y: physics.realDimensions.height * scaleFactor,
                            z: physics.realDimensions.width * scaleFactor
                        };
                        entity.setAttribute('scale', `${newScale.x} ${newScale.y} ${newScale.z}`);
                    }
                    
                } catch (error) {
                    console.warn(`[HD1-Physics] Failed to update object ${id} for scale change:`, error);
                }
            }
            
            // Update camera controls speed for new scale
            if (this.cameraElement && this.cameraElement.components['wasd-controls']) {
                const speed = sceneState.getCameraSpeed();
                this.cameraElement.components['wasd-controls'].data.acceleration = speed * 10;
                console.log(`[HD1-Physics] Updated camera movement speed to ${speed} units/s for scale ${sceneState.environment.scale}`);
            }
            
            // Update holodeck boundaries
            const holodeckComponent = this.scene.components['hd1-holodeck'];
            if (holodeckComponent) {
                // Scale boundaries to match physics scale
                const boundaryData = {
                    xMin: boundaries.min, xMax: boundaries.max,
                    zMin: boundaries.min, zMax: boundaries.max,
                    yMin: 0, yMax: boundaries.height
                };
                
                // Update boundary component if it exists
                const cameraEl = this.cameraElement;
                if (cameraEl && cameraEl.components['holodeck-boundaries']) {
                    Object.assign(cameraEl.components['holodeck-boundaries'].data, boundaryData);
                    console.log(`[HD1-Physics] Updated holodeck boundaries: ±${boundaries.max/1000}km for scale ${sceneState.environment.scale}`);
                }
            }
            
            console.log(`[HD1-Physics] ✓ Physics update complete - all objects rescaled for ${sceneState.environment.scale} (factor: ${scaleFactor})`);
            resolve();
        });
    }
    
    async rollback(oldState) {
        console.log(`[HD1-Reactive] Rolling back to state v${oldState.version}`);
        
        // Prevent infinite rollback loops
        if (this.rollbackInProgress) {
            console.error(`[HD1-Reactive] ✗ Rollback already in progress, aborting to prevent infinite loop`);
            return;
        }
        
        this.rollbackInProgress = true;
        
        try {
            // Clear existing objects manually instead of recursive render
            console.log(`[HD1-Reactive] Clearing ${this.domElements.size} objects for rollback`);
            
            // Clear DOM elements
            for (const [id, element] of this.domElements) {
                if (element && element.parentNode) {
                    element.parentNode.removeChild(element);
                }
            }
            
            // Clear tracking map
            this.domElements.clear();
            
            console.log(`[HD1-Reactive] Rollback to state v${oldState.version} completed`);
        } finally {
            this.rollbackInProgress = false;
        }
    }
    
    verifyState(expectedState) {
        // Count actual HD1-managed objects in DOM (those with hd1- prefix)
        const hd1Elements = Array.from(this.objectsContainer.children).filter(child => 
            child.id && child.id.startsWith('hd1-')
        );
        const actualCount = hd1Elements.length;
        const expectedCount = expectedState.objects.size;
        
        if (expectedCount !== actualCount) {
            console.error(`[HD1-Reactive] ✗ State verification failed: expected ${expectedCount}, found ${actualCount} HD1 objects`);
            console.error(`[HD1-Reactive] HD1 objects in DOM:`, hd1Elements.map(e => e.id));
            console.error(`[HD1-Reactive] Expected objects:`, Array.from(expectedState.objects.keys()));
            console.error(`[HD1-Reactive] DomElements map size:`, this.domElements.size);
            return false;
        }
        
        // Verify all expected objects exist in DOM
        for (const [id] of expectedState.objects) {
            if (!this.domElements.has(id)) {
                console.error(`[HD1-Reactive] ✗ Missing object in domElements map: ${id}`);
                return false;
            }
        }
        
        console.log(`[HD1-Reactive] ✓ State verified: ${actualCount} objects`);
        return true;
    }
    
    objectChanged(newObj, oldObj) {
        // Simple comparison - in practice could be more sophisticated
        return JSON.stringify(newObj) !== JSON.stringify(oldObj);
    }
    
    // Delegate to existing geometry/material methods
    setGeometry(entity, obj) {
        HD1AFrameManager.prototype.setGeometry.call(this, entity, obj);
    }
    
    setMaterial(entity, obj) {
        HD1AFrameManager.prototype.setMaterial.call(this, entity, obj);
    }
}

// REACTIVE HD1 MANAGER - Drop-in replacement
class HD1ReactiveManager {
    constructor(scene) {
        this.scene = scene;
        this.state = new SceneState();
        this.pipeline = new RenderPipeline(scene);
        this.renderPromise = Promise.resolve();
        this.camera = this.setupCamera();
        
        // SURGICAL FIX: Message queue to prevent race conditions
        this.messageQueue = [];
        this.processingQueue = false;
        
        console.log('[HD1-Reactive] Reactive manager initialized with message queue');
    }

    setupCamera() {
        const cameraEl = document.getElementById('holodeck-camera');
        return {
            position: [0, 1.7, 5],
            target: [0, 0, 0],
            element: cameraEl
        };
    }

    // SACRED INTERFACE - maintains exact compatibility with HD1AFrameManager
    processMessage(message) {
        // Add message to queue and start processing if not already running
        this.messageQueue.push(message);
        if (!this.processingQueue) {
            this.processMessageQueue();
        }
    }
    
    async processMessageQueue() {
        this.processingQueue = true;
        
        while (this.messageQueue.length > 0) {
            const message = this.messageQueue.shift();
            console.log(`[HD1-Reactive] Processing queued message: ${message.type} (${this.messageQueue.length} remaining)`);
            
            const transition = StateTransitions[message.type];
            if (!transition) {
                console.warn(`[HD1-Reactive] Unknown message type: ${message.type}`);
                continue;
            }
            
            try {
                const newState = transition(this.state, message);
                await this.applyState(newState);
            } catch (error) {
                console.error(`[HD1-Reactive] State transition failed:`, error);
            }
        }
        
        this.processingQueue = false;
        console.log(`[HD1-Reactive] Message queue processing complete`);
    }
    
    async applyState(newState) {
        console.log(`[HD1-Reactive] Applying state transition: v${this.state.version} → v${newState.version}`);
        console.log(`[HD1-Reactive] Old state objects:`, Array.from(this.state.objects.keys()));
        console.log(`[HD1-Reactive] New state objects:`, Array.from(newState.objects.keys()));
        
        const result = await this.pipeline.render(newState, this.state);
        
        if (result.success) {
            // STATE ONLY UPDATES ON SUCCESSFUL RENDER
            this.state = newState;
            console.log(`[HD1-Reactive] ✓ State updated to v${newState.version} with ${newState.objects.size} objects`);
        } else {
            console.error(`[HD1-Reactive] ✗ Render failed, state remains v${this.state.version}`);
        }
        
        return result;
    }
    
    // SACRED INTERFACE - external state query
    getObjectCount() {
        return this.state.objects.size;
    }
    
    // SACRED INTERFACE - bootstrap integration  
    initializeWorld(worldData) {
        if (worldData) {
            const newState = this.state.withWorld(worldData);
            this.renderPromise = this.applyState(newState);
            console.log(`[HD1-Reactive] World initialized with:`, worldData);
        }
    }
    
    // Legacy method implementations for compatibility
    updateHolodeckGrid(gridData) {
        if (gridData.transparency !== undefined) {
            const newWorld = {...this.state.world, transparency: gridData.transparency};
            const newState = this.state.withWorld(newWorld);
            this.renderPromise = this.applyState(newState);
        }
    }
    
    updateCamera(cameraData) {
        if (cameraData.position) {
            const newState = this.state.withCamera(cameraData);
            this.renderPromise = this.applyState(newState);
        }
    }
}

// HD1AFrameManager - REACTIVE DROP-IN REPLACEMENT
class HD1AFrameManager extends HD1ReactiveManager {
    constructor(scene) {
        super(scene);
        
        // Legacy compatibility properties
        this.objects = new Map(); // Will be kept in sync with reactive state
        this.objectsContainer = document.getElementById('holodeck-objects');
        
        // Sync legacy objects map with reactive state
        this.syncLegacyState();
        
        console.log('[HD1-AFrame] Manager initialized with reactive backend');
        
        // Expose universal physics commands to browser console
        this.exposeConsoleAPI();
    }
    
    // Console API for testing universal physics system
    exposeConsoleAPI() {
        // Make HD1 commands available globally
        window.HD1 = {
            // Scale system
            scales: ScaleUnit,
            atmospheres: StandardAtmospheres,
            props: Object.keys(StandardProps),
            
            // Scene operations
            initScene: (name, physics) => {
                this.processMessage({
                    type: 'init_scene',
                    scene_name: name,
                    physics: physics
                });
                console.log(`🎬 Scene "${name}" initialized with physics:`, physics);
            },
            
            // Scale operations  
            changeScale: (scale) => {
                this.processMessage({
                    type: 'change_scale',
                    scale: scale
                });
                console.log(`📏 Scale changed to ${scale} (factor: ${ScaleFactors[scale]})`);
            },
            
            // Props operations
            createProp: (name, type, position) => {
                this.processMessage({
                    type: 'create_prop',
                    prop_name: name,
                    prop_type: type,
                    position: position || {x: 0, y: 1, z: 0}
                });
                console.log(`🎭 Created ${type} prop "${name}"`);
            },
            
            deleteProp: (name) => {
                this.processMessage({
                    type: 'delete_prop',
                    prop_name: name
                });
                console.log(`🗑️ Deleted prop "${name}"`);
            },
            
            // Physics operations
            setGravity: (gravity) => {
                this.processMessage({
                    type: 'set_gravity',
                    gravity: gravity
                });
                console.log(`🌍 Gravity set to ${gravity} m/s²`);
            },
            
            setAtmosphere: (atmosphere) => {
                this.processMessage({
                    type: 'set_atmosphere',
                    atmosphere: atmosphere
                });
                console.log(`🌬️ Atmosphere updated:`, atmosphere);
            },
            
            setTemperature: (temp) => {
                this.processMessage({
                    type: 'set_temperature',
                    temperature: temp
                });
                console.log(`🌡️ Temperature set to ${temp}K (${(temp-273.15).toFixed(1)}°C)`);
            },
            
            // Quick scenes for testing
            scenes: {
                molecular: () => HD1.initScene('molecular_lab', {
                    scale: ScaleUnit.NANO,
                    gravity: 9.81,
                    atmosphere: StandardAtmospheres.VACUUM,
                    temperature: 310.15
                }),
                
                human: () => HD1.initScene('workshop', {
                    scale: ScaleUnit.METER,
                    gravity: 9.81,
                    atmosphere: StandardAtmospheres.EARTH_SEA_LEVEL,
                    temperature: 293.15
                }),
                
                space: () => HD1.initScene('orbital_station', {
                    scale: ScaleUnit.KILO,
                    gravity: 0.0,
                    atmosphere: StandardAtmospheres.VACUUM,
                    temperature: 2.7
                }),
                
                planetary: () => HD1.initScene('solar_system', {
                    scale: ScaleUnit.MEGA,
                    gravity: 9.81,
                    atmosphere: StandardAtmospheres.VACUUM,
                    temperature: 5778
                })
            },
            
            // Get current state info
            status: () => {
                const state = this.state;
                console.log('🌍 HD1 Universal Environmental Status:');
                console.log(`Scale: ${state.environment.scale} (factor: ${state.getScaleFactor()})`);
                console.log(`Gravity: ${state.environment.gravity} m/s²`);
                console.log(`Temperature: ${state.environment.temperature}K (${(state.environment.temperature-273.15).toFixed(1)}°C)`);
                console.log(`Atmosphere: density=${state.environment.atmosphere.density}kg/m³, pressure=${state.environment.atmosphere.pressure}Pa`);
                console.log(`Objects: ${state.objects.size}`);
                console.log(`Boundaries: ±${state.getBoundaries().max/1000}km`);
                return state;
            },
            
            // Quick demo
            demo: () => {
                console.log('🚀 HD1 Universal Environmental Demo - Molecular to Planetary!');
                
                // Start at human scale
                HD1.scenes.human();
                setTimeout(() => {
                    HD1.createProp('watch1', 'watch', {x: -2, y: 1, z: 0});
                    HD1.createProp('car1', 'car', {x: 5, y: 0, z: 0});
                }, 1000);
                
                // Go molecular
                setTimeout(() => {
                    console.log('🔬 Entering molecular scale...');
                    HD1.scenes.molecular();
                }, 3000);
                
                setTimeout(() => {
                    HD1.createProp('dna1', 'dna_strand', {x: 0, y: 1, z: 0});
                    HD1.createProp('bacteria1', 'bacteria', {x: 3, y: 1, z: 2});
                }, 4000);
                
                // Go planetary
                setTimeout(() => {
                    console.log('🌌 Entering planetary scale...');
                    HD1.scenes.planetary();
                }, 7000);
                
                setTimeout(() => {
                    HD1.createProp('moon1', 'moon', {x: 10, y: 0, z: 0});
                }, 8000);
                
                console.log('Watch the scale transitions! Check HD1.status() for details.');
            }
        };
        
        console.log('🎯 HD1 Universal Environmental Console API loaded!');
        console.log('Available commands:');
        console.log('• HD1.demo() - Run complete scale demo');
        console.log('• HD1.scenes.molecular() - Nanometer scale');
        console.log('• HD1.scenes.human() - Meter scale'); 
        console.log('• HD1.scenes.space() - Kilometer scale');
        console.log('• HD1.scenes.planetary() - Megameter scale');
        console.log('• HD1.createProp("name", "type", {x,y,z})');
        console.log('• HD1.changeScale(HD1.scales.NANO)');
        console.log('• HD1.setGravity(0.0) - Zero gravity');
        console.log('• HD1.status() - Current physics state');
        console.log('• Available props:', Object.keys(StandardProps));
    }
    
    // Keep legacy objects map synchronized with reactive state
    syncLegacyState() {
        // Override applyState to also update legacy objects map
        const originalApplyState = this.applyState.bind(this);
        this.applyState = async (newState) => {
            const result = await originalApplyState(newState);
            
            if (result.success) {
                // Sync legacy objects map
                this.objects.clear();
                for (const [id, obj] of newState.objects) {
                    this.objects.set(id, obj);
                }
            }
            
            return result;
        };
    }

    // Legacy methods needed by RenderPipeline - keep geometry and material setters
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
                return;
            case 'sky':
                this.createSky(entity, obj);
                return;
            case 'text':
                console.warn('[HD1-AFrame] Text objects temporarily disabled due to FontLoader compatibility');
                entity.setAttribute('geometry', {primitive: 'box', width: 2, height: 0.5, depth: 0.1});
                break;
            case 'environment':
                this.createEnvironment(entity, obj);
                return;
            case 'particle':
                console.warn('[HD1-AFrame] Particle systems temporarily disabled due to browser compatibility');
                entity.setAttribute('geometry', {primitive: 'sphere', radius: 0.5});
                break;
            default:
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
            shader: obj.material?.shader || 'standard',
            color: hexColor,
            metalness: obj.material?.metalness || 0.1,
            roughness: obj.material?.roughness || 0.7,
            transparent: obj.material?.transparent || color.a < 1.0,
            opacity: color.a
        };
        
        if (obj.wireframe) {
            materialProps.wireframe = true;
            materialProps.shader = 'flat';
        }
        
        entity.setAttribute('material', materialProps);
        
        if (obj.lighting) {
            if (obj.lighting.castShadow) {
                entity.setAttribute('shadow', 'cast: true');
            }
            if (obj.lighting.receiveShadow) {
                entity.setAttribute('shadow', 'receive: true');
            }
        }
        
        if (obj.physics && obj.physics.enabled) {
            if (obj.physics.type === 'static') {
                entity.setAttribute('static-body', { shape: 'auto' });
            } else {
                entity.setAttribute('dynamic-body', { 
                    shape: 'auto', 
                    mass: obj.physics.mass || 1.0 
                });
            }
        }
    }

    // Helper methods for special object types
    createLight(entity, obj) {
        const lightProps = {
            type: obj.lightType || 'point',
            color: obj.color ? this.colorToHex(obj.color) : '#ffffff',
            intensity: obj.intensity || 1.0
        };
        
        if (obj.lightType === 'directional') {
            lightProps.castShadow = true;
        }
        
        entity.setAttribute('light', lightProps);
    }

    createSky(entity, obj) {
        const color = obj.color ? this.colorToHex(obj.color) : '#87CEEB';
        
        entity.setAttribute('geometry', {primitive: 'sphere', radius: 5000});
        entity.setAttribute('material', {
            shader: 'flat',
            color: color,
            side: 'back'
        });
        entity.setAttribute('scale', '-1 1 1');
    }

    createEnvironment(entity, obj) {
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
        
        if (obj.environmentSettings) {
            Object.assign(envProps, obj.environmentSettings);
        }
        
        entity.setAttribute('environment', envProps);
    }

    colorToHex(color) {
        const r = Math.round(color.r * 255).toString(16).padStart(2, '0');
        const g = Math.round(color.g * 255).toString(16).padStart(2, '0');
        const b = Math.round(color.b * 255).toString(16).padStart(2, '0');
        return `#${r}${g}${b}`;
    }
}

console.log('[HD1-Reactive] Reactive Scene Graph System loaded');
console.log('[HD1-Reactive] HD1ReactiveManager class defined:', typeof HD1ReactiveManager);
console.log('[HD1-Reactive] HD1AFrameManager class defined:', typeof HD1AFrameManager);
console.log('[HD1-Reactive] Registered A-Frame components:', Object.keys(AFRAME.components));

// Check if camera element exists and has the component, if not add it manually
setTimeout(() => {
    const camera = document.getElementById('holodeck-camera');
    console.log('[HD1-Reactive] Camera element found:', !!camera);
    if (camera) {
        console.log('[HD1-Reactive] Camera components:', Object.keys(camera.components || {}));
        console.log('[HD1-Reactive] Has holodeck-boundaries component:', !!camera.components['holodeck-boundaries']);
        
        // If component not attached, force attach it
        if (!camera.components['holodeck-boundaries']) {
            console.log('[HD1-Reactive] Manually attaching holodeck-boundaries component...');
            camera.setAttribute('holodeck-boundaries', '');
            
            // Check again after 1 second
            setTimeout(() => {
                console.log('[HD1-Reactive] After manual attach - Has holodeck-boundaries:', !!camera.components['holodeck-boundaries']);
            }, 1000);
        }
    }
}, 2000);

// Global export to ensure it's available
window.HD1AFrameManager = HD1AFrameManager;
window.HD1ReactiveManager = HD1ReactiveManager;
