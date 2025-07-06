/**
 * HD1 v3.0 PlayCanvas Game Engine Integration
 * Complete API-first 3D game development interface
 */

// Global HD1 PlayCanvas Game Engine
let hd1GameEngine = null;
let currentSession = null; // Will be set dynamically from active session

// Get current active session ID
function getCurrentSession() {
    if (currentSession) return currentSession;
    
    // Try to get from local storage first (primary source)
    const storedSession = localStorage.getItem('hd1_session_id');
    if (storedSession) {
        currentSession = storedSession;
        console.log('[HD1] Found session in localStorage:', currentSession);
        return currentSession;
    }
    
    // Try to get from global console manager
    if (window.hd1Console && window.hd1Console.sessionManager && window.hd1Console.sessionManager.sessionId) {
        currentSession = window.hd1Console.sessionManager.sessionId;
        console.log('[HD1] Found session from console manager:', currentSession);
        return currentSession;
    }
    
    // Try to get from global variable if set by console
    if (window.activeSessionId) {
        currentSession = window.activeSessionId;
        console.log('[HD1] Found session from global variable:', currentSession);
        return currentSession;
    }
    
    // Try to get from console display elements
    const sessionElement = document.getElementById('session-id-tag-status');
    if (sessionElement && sessionElement.textContent && sessionElement.textContent.startsWith('session-')) {
        currentSession = sessionElement.textContent;
        console.log('[HD1] Found session from UI element:', currentSession);
        return currentSession;
    }
    
    console.warn('[HD1] No active session found, API calls may fail');
    return null;
}

// Initialize PlayCanvas engine on page load
document.addEventListener('DOMContentLoaded', function() {
    initializePlayCanvasEngine();
    setupGameEngineControls();
    updateGameStats();
});

/**
 * Initialize PlayCanvas Engine with HD1 configuration
 */
function initializePlayCanvasEngine() {
    console.log('[HD1] Starting PlayCanvas initialization...');
    
    // Check if PlayCanvas is available
    if (typeof pc === 'undefined') {
        console.error('[HD1] PlayCanvas library not loaded! Check if playcanvas.min.js is included.');
        return;
    }
    
    const canvas = document.getElementById('hd1-playcanvas-canvas');
    if (!canvas) {
        console.error('[HD1] PlayCanvas canvas not found');
        return;
    }
    
    console.log('[HD1] Canvas found, creating PlayCanvas application...');

    try {
        // Detect WebGL renderer backend using WEBGL_debug_renderer_info
        function detectWebGLRenderer() {
            const testCanvas = document.createElement('canvas');
            const gl = testCanvas.getContext('webgl') || testCanvas.getContext('experimental-webgl');
            
            if (!gl) return { renderer: 'unknown', isMetal: false };
            
            const debugInfo = gl.getExtension('WEBGL_debug_renderer_info');
            if (debugInfo) {
                const renderer = gl.getParameter(debugInfo.UNMASKED_RENDERER_WEBGL);
                const isMetal = renderer && renderer.includes('Metal');
                console.log('[HD1] WebGL Renderer detected:', renderer);
                return { renderer, isMetal };
            }
            
            return { renderer: 'unknown', isMetal: false };
        }
        
        const rendererInfo = detectWebGLRenderer();
        const isChromeMac = /Chrome/.test(navigator.userAgent) && /Mac/.test(navigator.platform);
        
        // Create PlayCanvas application with proper backend configuration
        const graphicsOptions = {
            antialias: true,
            alpha: false,
            depth: true,
            stencil: true
        };
        
        // Force OpenGL backend instead of Metal to avoid shader compilation issues
        if (isChromeMac || rendererInfo.isMetal) {
            console.log('[HD1] Metal backend detected - implementing OpenGL fallback');
            
            // Multiple strategies to force OpenGL backend
            graphicsOptions.powerPreference = 'default'; // Avoid discrete GPU
            graphicsOptions.failIfMajorPerformanceCaveat = false; // Accept performance tradeoffs
            graphicsOptions.antialias = false; // Reduce Metal backend triggers
            
            // Force WebGL 1.0 context to avoid Metal backend entirely
            graphicsOptions.preferWebGl1 = true;
            graphicsOptions.webgl1 = true;
            
            // Additional Metal backend avoidance
            graphicsOptions.preserveDrawingBuffer = false;
            graphicsOptions.premultipliedAlpha = false;
            
            console.log('[HD1] OpenGL fallback configured - WebGL 1.0 forced');
        }
        
        const app = new pc.Application(canvas, {
            mouse: new pc.Mouse(canvas),
            keyboard: new pc.Keyboard(window),
            touch: new pc.TouchDevice(canvas),
            elementInput: new pc.ElementInput(canvas),
            graphicsDeviceOptions: graphicsOptions
        });

        console.log('[HD1] PlayCanvas application created, configuring...');

        // Configure rendering settings
        app.setCanvasFillMode(pc.FILLMODE_FILL_WINDOW);
        app.setCanvasResolution(pc.RESOLUTION_AUTO);

        // Add error handling for shader compilation issues
        app.graphicsDevice.on('error', function(message) {
            console.warn('[HD1] Graphics device error (possibly Metal backend):', message);
            // Continue running - don't crash the application
        });
        
        // Start the application
        app.start();
        console.log('[HD1] PlayCanvas application started');
        
        // Check final WebGL renderer after PlayCanvas initialization
        const finalCanvas = app.graphicsDevice.canvas;
        const finalGl = finalCanvas.getContext('webgl') || finalCanvas.getContext('webgl2');
        if (finalGl) {
            const finalDebugInfo = finalGl.getExtension('WEBGL_debug_renderer_info');
            if (finalDebugInfo) {
                const finalRenderer = finalGl.getParameter(finalDebugInfo.UNMASKED_RENDERER_WEBGL);
                console.log('[HD1] Final WebGL Renderer:', finalRenderer);
                
                if (finalRenderer && finalRenderer.includes('Metal')) {
                    console.warn('[HD1] Warning: Still using Metal backend - shader compilation issues may occur');
                    console.warn('[HD1] Chrome Metal backend detected - this causes shader compilation errors');
                    console.warn('[HD1] To fix: Go to chrome://flags/#use-angle-gl and set to "OpenGL"');
                    console.warn('[HD1] Or try: chrome://flags/#use-angle-metal and set to "Disabled"');
                    
                    // Display user-friendly alert
                    setTimeout(() => {
                        alert('Chrome Metal Backend Issue:\n\nTo fix the blank screen:\n1. Go to chrome://flags/#use-angle-gl\n2. Set to "OpenGL"\n3. Restart Chrome\n\nOr use Firefox which works correctly.');
                    }, 2000);
                } else {
                    console.log('[HD1] Success: Using OpenGL backend - shader compilation should work correctly');
                }
            }
        }

        // Initialize physics world with realistic settings - PHYSICS EXCELLENCE
        app.systems.rigidbody.setGravity(0, -9.81, 0); // Standard Earth gravity
        console.log('[HD1] 🎮 Physics world initialized: gravity=9.81m/s², physics engine ready');

        // Create empty scene - content loaded from channels/scenes
        createEmptyScene(app);

        // Store globally
        hd1GameEngine = app;
        window.hd1GameEngine = app;

        console.log('[HD1] PlayCanvas engine fully initialized and ready');
    } catch (error) {
        console.error('[HD1] Failed to initialize PlayCanvas:', error);
        return;
    }
    
    // Update status if element exists
    const statusElement = document.getElementById('playcanvas-status');
    if (statusElement) {
        statusElement.textContent = 'PlayCanvas Active';
    }
    
    // Load any existing objects for the current session
    setTimeout(() => {
        loadExistingSessionObjects();
    }, 500); // Small delay to ensure other systems are ready
}

/**
 * Create empty 3D scene with camera only - content loaded from channels
 */
function createEmptyScene(app) {
    // Create camera entity
    const camera = new pc.Entity('camera');
    camera.addComponent('camera', {
        clearColor: new pc.Color(0.1, 0.1, 0.1),
        farClip: 1000
    });
    camera.setPosition(0, 5, 15);
    camera.lookAt(0, 0, 0);
    
    app.root.addChild(camera);
    
    // Store camera reference globally
    app.camera = camera;
    
    // Create camera controller
    app.cameraController = new HD1CameraController(app, camera);
    
    // Camera controller will be updated in setupCameraControls
    
    // Setup manual camera controls
    setupCameraControls(app, camera);
    
    // Setup automatic avatar-driven camera system
    setupAvatarDrivenCamera(app.cameraController);
    
    console.log('[HD1] Empty PlayCanvas scene created - ready for channel content');
}


/**
 * 🎮 ORBITAL CAMERA: Enhanced camera system with multiple modes
 */
class HD1CameraController {
    constructor(app, camera) {
        this.app = app;
        this.camera = camera;
        this.mode = 'avatar-driven'; // Only avatar-driven mode
        
        // Avatar-driven camera properties
        this.boundAvatar = null;
        this.avatarCameraOffset = new pc.Vec3(0, 1.7, 0); // Default head level
        this.availableAvatars = [];
        
        console.log('[HD1] 🎮 Camera Controller initialized - Avatar-driven mode only');
    }
    
    
    // Switch to avatar-driven camera mode based on avatar type
    setAvatarDrivenMode(avatar = null) {
        this.mode = 'avatar-driven';
        this.updateAvailableAvatars();
        
        if (avatar) {
            this.boundAvatar = avatar;
        } else if (this.availableAvatars.length > 0) {
            // Find the current session's avatar
            const sessionId = getCurrentSession();
            this.boundAvatar = this.availableAvatars.find(av => 
                av.name && av.name.includes(sessionId)
            ) || this.availableAvatars[0];
        } else {
            console.warn('[HD1] 🎮 No avatars available for avatar-driven camera');
            return;
        }
        
        // Determine avatar type and set appropriate vision system
        if (!this.configureAvatarVision(this.boundAvatar)) {
            console.error('[HD1] 🎮 Failed to configure camera for avatar, unbinding');
            this.boundAvatar = null;
            return;
        }
        
        console.log('[HD1] 🎮 Avatar-driven camera activated:', this.boundAvatar.name);
        // NOTE: No updateAvatarPosition() - camera moves independently, avatar follows via API sync
    }
    
    // Configure camera and physics based on avatar configuration - SINGLE SOURCE OF TRUTH
    configureAvatarVision(avatar) {
        if (!avatar) {
            console.error('[HD1] 🎮 Avatar configuration failed: No avatar provided');
            return false;
        }
        
        // Camera configuration from avatar data
        if (avatar.hd1CameraConfig) {
            const config = avatar.hd1CameraConfig;
            this.avatarCameraOffset.set(
                config.position.x || 0,
                config.position.y || 1.7,
                config.position.z || 0
            );
            if (this.camera.camera && config.field_of_view) {
                this.camera.camera.fov = config.field_of_view;
            }
        } else {
            // Fallback defaults
            this.avatarCameraOffset.set(0, 1.7, 0);
            if (this.camera.camera) {
                this.camera.camera.fov = 75;
            }
        }
        
        // Simple avatar configuration - RESTORE WORKING STATE
        console.log('[HD1] 🎮 Basic avatar configuration (physics disabled)');
        
        console.log(`[HD1] 🎮 Avatar configured successfully:`, {
            avatar: avatar.name,
            camera_offset: this.avatarCameraOffset,
            physics_enabled: !!avatar.hd1PhysicsConfig
        });
        
        return true;
    }
    
    
    
    
    
    // Update available avatars list
    updateAvailableAvatars() {
        this.availableAvatars = this.app.root.children.filter(entity => 
            entity.hd1Tags && entity.hd1Tags.includes('session-avatar')
        );
    }
    
    // Update avatar-driven camera position based on avatar type
    updateAvatarPosition() {
        if (!this.boundAvatar || this.mode !== 'avatar-driven') return;
        
        const avatarPos = this.boundAvatar.getPosition();
        const avatarRot = this.boundAvatar.getRotation();
        
        // Calculate camera position based on avatar vision type
        const offset = this.avatarCameraOffset.clone();
        const rotatedOffset = new pc.Vec3();
        avatarRot.transformVector(offset, rotatedOffset);
        
        this.targetPosition.copy(avatarPos).add(rotatedOffset);
        
        // Configure look target based on perspective type
        let lookTarget;
        
        switch(this.avatarPerspective) {
            case 'first-person':
                // First-person: look in the direction the avatar is facing
                lookTarget = avatarPos.clone();
                const forwardDir = new pc.Vec3();
                avatarRot.transformVector(pc.Vec3.FORWARD, forwardDir);
                lookTarget.add(forwardDir.scale(10)); // Look ahead
                break;
                
            case 'third-person-elevated':
                // Elevated third-person: look down at avatar from above
                lookTarget = avatarPos.clone();
                lookTarget.y += 0.5; // Look at avatar center
                break;
                
            default: // third-person
                // Standard third-person: look at avatar
                lookTarget = avatarPos.clone();
                lookTarget.y += 1; // Look at avatar's upper body
                break;
        }
        
        const lookDirection = new pc.Vec3();
        lookDirection.sub2(lookTarget, this.targetPosition).normalize();
        const lookAtMatrix = new pc.Mat4();
        lookAtMatrix.setLookAt(this.targetPosition, lookTarget, pc.Vec3.UP);
        this.targetRotation.setFromMat4(lookAtMatrix);
    }
    
    
    // Update camera (called every frame)
    update(dt) {
        // Update avatar position to match camera - SIMPLE WORKING APPROACH
        if (this.boundAvatar) {
            const avatar = this.boundAvatar;
            const cameraPos = this.camera.getPosition();
            const cameraRot = this.camera.getRotation();
            
            // Avatar follows camera (restore working behavior)
            const avatarPos = new pc.Vec3();
            avatarPos.copy(cameraPos).sub(this.avatarCameraOffset);
            
            // SIMPLE GROUND COLLISION: Prevent avatar from going below Y=0
            if (avatarPos.y < 0) {
                avatarPos.y = 0;
                // Also adjust camera to maintain offset
                const newCameraPos = avatarPos.clone().add(this.avatarCameraOffset);
                this.camera.setPosition(newCameraPos);
            }
            
            avatar.setPosition(avatarPos);
            avatar.setRotation(cameraRot);
        }
    }
    
    
    
}

/**
 * Setup camera controls for mouse look and WASD movement
 */
function setupCameraControls(app, camera) {
    let pitch = 0;
    let yaw = 0;
    let moveSpeed = 10;
    let lookSpeed = 0.2;
    let isMouseDown = false;
    
    // 🎮 SMOOTH MOVEMENT: Add interpolation and momentum variables
    let targetPosition = camera.getPosition().clone();
    let currentVelocity = new pc.Vec3(0, 0, 0);
    let targetVelocity = new pc.Vec3(0, 0, 0);
    const smoothingFactor = 0.15; // Higher = snappier, Lower = smoother
    const momentumDecay = 0.85; // How quickly momentum fades
    const accelerationRate = 0.2; // How quickly we reach target speed
    
    // 🎮 ORBITAL CAMERA: Use existing camera controller
    const cameraController = app.cameraController;
    window.hd1CameraController = cameraController; // Make it globally accessible
    
    // Mouse look controls
    app.mouse.on(pc.EVENT_MOUSEDOWN, function(event) {
        if (event.button === pc.MOUSEBUTTON_LEFT) {
            // Check if mouse is over console panel
            const debugPanel = document.getElementById('debug-panel');
            const rect = debugPanel.getBoundingClientRect();
            const mouseX = event.x;
            const mouseY = event.y;
            
            // Don't enable camera controls if clicking on console
            if (mouseX >= rect.left && mouseX <= rect.right && 
                mouseY >= rect.top && mouseY <= rect.bottom) {
                return;
            }
            
            isMouseDown = true;
            app.mouse.enablePointerLock();
        }
    });
    
    app.mouse.on(pc.EVENT_MOUSEUP, function(event) {
        if (event.button === pc.MOUSEBUTTON_LEFT) {
            isMouseDown = false;
            app.mouse.disablePointerLock();
        }
    });
    
    app.mouse.on(pc.EVENT_MOUSEMOVE, function(event) {
        if (isMouseDown && document.pointerLockElement) {
            // Avatar-driven camera mode
            yaw -= event.dx * lookSpeed;
            pitch -= event.dy * lookSpeed;
            pitch = pc.math.clamp(pitch, -90, 90);
            
            camera.setEulerAngles(pitch, yaw, 0);
        }
    });
    
    
    // WASD movement controls with HD1 API synchronization
    let lastCameraUpdate = 0;
    const cameraUpdateThrottle = 100; // Update camera position via API every 100ms
    
    // 🎮 PHYSICS-BASED CONTROLS: Avatar-specific abilities
    let lastKeyTime = 0;
    app.keyboard.on(pc.EVENT_KEYDOWN, function(event) {
        const now = Date.now();
        
        // Simple space key handling (no physics)
        if (event.key === pc.KEY_SPACE) {
            // Basic space key - can add simple jump later
        }
        
        // Prevent rapid key presses for other keys
        if (now - lastKeyTime < 200) return;
        lastKeyTime = now;
    });
    
    app.on('update', function(dt) {
        // Update camera controller
        cameraController.update(dt);
        
        // Process WASD movement in avatar-driven mode
        if (cameraController.mode === 'avatar-driven') {
            const forward = camera.forward;
            const right = camera.right;
            const up = pc.Vec3.UP;
            
            // 🎮 SMOOTH MOVEMENT: Calculate target velocity based on input
            targetVelocity.set(0, 0, 0);
            let hasInput = false;
            
            if (app.keyboard.isPressed(pc.KEY_W)) {
                targetVelocity.add(forward.clone().scale(moveSpeed));
                hasInput = true;
            }
            if (app.keyboard.isPressed(pc.KEY_S)) {
                targetVelocity.add(forward.clone().scale(-moveSpeed));
                hasInput = true;
            }
            if (app.keyboard.isPressed(pc.KEY_A)) {
                targetVelocity.add(right.clone().scale(-moveSpeed));
                hasInput = true;
            }
            if (app.keyboard.isPressed(pc.KEY_D)) {
                targetVelocity.add(right.clone().scale(moveSpeed));
                hasInput = true;
            }
            if (app.keyboard.isPressed(pc.KEY_Q)) {
                targetVelocity.add(up.clone().scale(-moveSpeed));
                hasInput = true;
            }
            if (app.keyboard.isPressed(pc.KEY_E)) {
                targetVelocity.add(up.clone().scale(moveSpeed));
                hasInput = true;
            }
            
            // 🎮 MOMENTUM: Apply acceleration/deceleration
            if (hasInput) {
                // Accelerate towards target velocity
                currentVelocity.lerp(currentVelocity, targetVelocity, accelerationRate);
            } else {
                // Apply momentum decay when no input
                currentVelocity.scale(momentumDecay);
            }
            
            // 🎮 SIMPLE MOVEMENT: Apply direct camera movement (restore working behavior)
            if (currentVelocity.length() > 0.01) {
                targetPosition.add(currentVelocity.clone().scale(dt));
                const currentPos = camera.getPosition();
                const newPos = new pc.Vec3();
                newPos.lerp(currentPos, targetPosition, smoothingFactor);
                camera.setPosition(newPos);
                
                // Sync avatar position to API for multiplayer
                const now = Date.now();
                if (now - lastCameraUpdate > cameraUpdateThrottle) {
                    lastCameraUpdate = now;
                    const avatarPos = camera.getPosition();
                    syncCameraPositionToAPI(avatarPos);
                }
            }
            
            // Update target position to current position (for smooth catch-up)
            targetPosition.copy(camera.getPosition());
        }
    });
    
    // Object interaction - right click to select (only if not over console)
    document.addEventListener('mousedown', function(event) {
        if (event.button === 2) { // Right mouse button
            // Check if mouse is over console panel
            const debugPanel = document.getElementById('debug-panel');
            const rect = debugPanel.getBoundingClientRect();
            
            // Don't handle object selection if clicking on console
            if (event.clientX >= rect.left && event.clientX <= rect.right && 
                event.clientY >= rect.top && event.clientY <= rect.bottom) {
                return;
            }
            
            // Simple object highlighting (since proper ray casting would be complex)
            if (hd1GameEngine) {
                hd1GameEngine.root.children.forEach(entity => {
                    if (entity.hd1Id && entity.model && entity.model.meshInstances[0]) {
                        // Toggle highlight on all objects (simplified interaction)
                        const material = entity.model.meshInstances[0].material;
                        if (material.emissive && material.emissive.r > 0) {
                            material.emissive = new pc.Color(0, 0, 0); // Remove highlight
                        } else {
                            material.emissive = new pc.Color(0.3, 0.3, 0.3); // Add highlight
                        }
                        material.update();
                        
                        console.log('[HD1] Toggled highlight on object:', entity.name);
                    }
                });
            }
        }
    });
    
    console.log('[HD1] Camera controls enabled: Mouse to look, WASD to move, QE for up/down, Right-click to select objects');
}

/**
 * Setup camera control buttons and keyboard shortcuts
 */
function setupAvatarDrivenCamera(cameraController) {
    // Update vision indicator
    function updateVisionIndicator() {
        const indicator = document.getElementById('vision-type');
        if (!indicator) return;
        
        let displayText = '👁️ Avatar';
        
        if (cameraController.boundAvatar) {
            displayText = '👁️ Avatar';
        } else {
            displayText = '👁️ No Avatar';
        }
        
        indicator.textContent = displayText;
    }
    
    // Automatically detect and switch to avatar-driven camera when avatars are available
    function checkForAvatars() {
        cameraController.updateAvailableAvatars();
        
        if (cameraController.availableAvatars.length > 0 && !cameraController.boundAvatar) {
            // Automatically bind to avatar when avatars become available
            cameraController.setAvatarDrivenMode();
            updateVisionIndicator();
            console.log('[HD1] 🎮 Auto-bound to avatar - avatars detected');
        } else if (cameraController.availableAvatars.length === 0 && cameraController.boundAvatar) {
            // Unbind if no avatars available
            cameraController.boundAvatar = null;
            updateVisionIndicator();
            console.log('[HD1] 🎮 Unbound from avatar - no avatars detected');
        }
    }
    
    // Check for avatars periodically
    setInterval(checkForAvatars, 2000); // Check every 2 seconds
    
    // Manual override keys (for debugging/admin purposes)
    document.addEventListener('keydown', (e) => {
        // Only handle if not typing in input fields
        if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') return;
        
        // Hold Shift for manual camera overrides
        if (e.shiftKey) {
            switch(e.key.toLowerCase()) {
                case 'f':
                    e.preventDefault();
                    cameraController.setFreeMode();
                    updateVisionIndicator();
                    console.log('[HD1] 🎮 Manual override: Free camera mode');
                    break;
                case 'o':
                    e.preventDefault();
                    const avatarCenter = cameraController.calculateAvatarCenter();
                    cameraController.setOrbitMode(avatarCenter, 20, 8);
                    updateVisionIndicator();
                    console.log('[HD1] 🎮 Manual override: Orbit camera mode');
                    break;
                case 'a':
                    e.preventDefault();
                    cameraController.setAvatarDrivenMode();
                    updateVisionIndicator();
                    console.log('[HD1] 🎮 Manual override: Avatar-driven camera mode');
                    break;
            }
        }
    });
    
    // Initial check and indicator update
    updateVisionIndicator();
    setTimeout(checkForAvatars, 1000); // Initial check after 1 second
    
    console.log('[HD1] 🎮 Avatar-driven camera system enabled - automatic avatar detection active');
}

/**
 * Sync camera position to HD1 API for avatar updates and participant synchronization
 */
async function syncCameraPositionToAPI(position) {
    try {
        const sessionId = getCurrentSession();
        if (!sessionId) {
            console.warn('[HD1] ❌ No session ID found - cannot sync camera position!');
            console.log('[HD1] localStorage hd1_session_id:', localStorage.getItem('hd1_session_id'));
            return;
        }
        
        // DEBUG: Log API call to verify it's happening
        console.log('[HD1] 🎯 Syncing camera position to API:', { sessionId, position: { x: position.x, y: position.y, z: position.z } });
        
        // Call HD1 API to update camera position (which updates avatar and broadcasts to participants)
        const result = await window.hd1API.setCameraPosition(sessionId, {
            x: position.x,
            y: position.y, 
            z: position.z
        });
        
        console.log('[HD1] ✅ Camera position sync successful:', result);
        
        // Track camera position in stats manager
        if (window.hd1ConsoleManager) {
            const statsManager = window.hd1ConsoleManager.getModule('stats');
            if (statsManager && statsManager.trackCameraPosition) {
                statsManager.trackCameraPosition(position);
            }
        }
        
    } catch (error) {
        // LOG ALL ERRORS to debug the issue
        console.error('[HD1] ❌ Camera position sync failed:', error);
        console.log('[HD1] Error details:', {
            message: error.message,
            sessionId: getCurrentSession(),
            position: position
        });
    }
}

/**
 * Setup game engine control buttons
 */
function setupGameEngineControls() {
    // Bind control functions to global scope
    window.createEntity = createEntity;
    window.startAnimation = startAnimation;
    window.playAudio = playAudio;
    window.applyPhysics = applyPhysics;
    window.loadDemoScene = loadDemoScene;

    // Setup dropdown event listeners
    setupDropdownControls();
}

/**
 * Create holodeck entity via HD1 API
 */
async function createEntity() {
    try {
        const sessionId = getCurrentSession();
        if (!sessionId) {
            console.error('[HD1] Cannot create entity: no active session');
            return;
        }
        
        const response = await hd1API.createEntity(sessionId, {
            name: `entity_${Date.now()}`,
            components: {
                transform: {
                    position: { x: Math.random() * 10 - 5, y: 0, z: Math.random() * 10 - 5 },
                    rotation: { x: 0, y: 0, z: 0 }
                }
            }
        });

        if (response.success) {
            // Add model component
            await hd1API.addComponent(sessionId, response.entity_id, {
                type: 'model',
                properties: {
                    type: 'box',
                    material: { 
                        color: `#${Math.floor(Math.random()*16777215).toString(16)}`,
                        metalness: 0.8,
                        roughness: 0.2
                    }
                }
            });

            // Create visual representation in PlayCanvas
            createPlayCanvasEntity(response.entity_id, response.name);
            updateGameStats();
            console.log('[HD1] Entity created:', response.entity_id);
        }
    } catch (error) {
        console.error('[HD1] Failed to create entity:', error);
    }
}

/**
 * Create PlayCanvas visual entity
 */
function createPlayCanvasEntity(entityId, name) {
    if (!hd1GameEngine) return;

    const entity = new pc.Entity(name);
    entity.hd1Id = entityId;

    // Add model component
    entity.addComponent('model', {
        type: 'box'
    });

    // Random position
    entity.setPosition(
        Math.random() * 10 - 5,
        0,
        Math.random() * 10 - 5
    );

    // Random color material
    const material = new pc.StandardMaterial();
    material.diffuse = new pc.Color(Math.random(), Math.random(), Math.random());
    material.metalness = 0.8;
    material.shininess = 80;
    material.update();

    entity.model.meshInstances[0].material = material;

    hd1GameEngine.root.addChild(entity);
}

/**
 * Start animation via HD1 API
 */
async function startAnimation() {
    try {
        const sessionId = getCurrentSession();
        if (!sessionId) {
            console.error('[HD1] Cannot start animation: no active session');
            return;
        }
        
        const response = await hd1API.createAnimation(sessionId, {
            name: `rotation_${Date.now()}`,
            targets: [{
                entity_id: 'all',
                property: 'rotation.y',
                from: 0,
                to: 360,
                duration: 3000,
                loop: true
            }]
        });

        if (response.success) {
            await hd1API.playAnimation(sessionId, response.animation_id, {});
            
            // Start PlayCanvas rotation for all entities
            startPlayCanvasRotation();
            updateGameStats();
            console.log('[HD1] Animation started:', response.animation_id);
        }
    } catch (error) {
        console.error('[HD1] Failed to start animation:', error);
    }
}

/**
 * Start PlayCanvas rotation animation
 */
function startPlayCanvasRotation() {
    if (!hd1GameEngine) return;

    hd1GameEngine.root.children.forEach(entity => {
        if (entity.hd1Id) {
            entity.rotationSpeed = 30; // degrees per second
        }
    });

    // Add rotation update to app loop
    if (!hd1GameEngine.rotationHandler) {
        hd1GameEngine.rotationHandler = hd1GameEngine.on('update', function(dt) {
            hd1GameEngine.root.children.forEach(entity => {
                if (entity.rotationSpeed) {
                    entity.rotateLocal(0, entity.rotationSpeed * dt, 0);
                }
            });
        });
    }
}

/**
 * Play audio via HD1 API
 */
async function playAudio() {
    try {
        const sessionId = getCurrentSession();
        if (!sessionId) {
            console.error('[HD1] Cannot play audio: no active session');
            return;
        }
        
        const response = await hd1API.createAudioSource(sessionId, {
            name: `audio_${Date.now()}`,
            type: 'positional',
            url: 'game_sound.ogg',
            loop: false,
            volume: 0.7
        });

        if (response.success) {
            await hd1API.playAudio(sessionId, response.audio_id, {});
            updateGameStats();
            console.log('[HD1] Audio playing:', response.audio_id);
        }
    } catch (error) {
        console.error('[HD1] Failed to play audio:', error);
    }
}

/**
 * Apply physics via HD1 API
 */
async function applyPhysics() {
    try {
        const sessionId = getCurrentSession();
        if (!sessionId) {
            console.error('[HD1] Cannot apply physics: no active session');
            return;
        }
        
        // Configure physics world
        await hd1API.updatePhysicsWorld(sessionId, {
            gravity: { x: 0, y: -9.8, z: 0 },
            timeStep: 0.016
        });

        // Get all entities and add physics
        const entities = await hd1API.listEntities(sessionId);
        
        for (const entity of entities.entities || []) {
            await hd1API.addComponent(sessionId, entity.entity_id, {
                type: 'rigidbody',
                properties: {
                    type: 'dynamic',
                    mass: 1.0
                }
            });
        }

        updateGameStats();
        console.log('[HD1] Physics applied to all entities');
    } catch (error) {
        console.error('[HD1] Failed to apply physics:', error);
    }
}

/**
 * Load demo holodeck experience
 */
async function loadDemoScene() {
    try {
        const sessionId = getCurrentSession();
        if (!sessionId) {
            console.error('[HD1] Cannot load demo scene: no active session');
            return;
        }
        
        // Create presentation screen
        const screenResponse = await hd1API.createEntity(sessionId, {
            name: 'presentation_screen',
            components: {
                transform: { position: { x: 0, y: 2, z: -5 } }
            }
        });

        if (screenResponse.success) {
            await hd1API.addComponent(sessionId, screenResponse.entity_id, {
                type: 'model',
                properties: { type: 'plane', material: { color: '#ffffff' } }
            });

            createPlayCanvasEntity(screenResponse.entity_id, 'presentation_screen');
        }

        // Create interactive objects
        for (let i = 0; i < 3; i++) {
            const objectResponse = await hd1API.createEntity(sessionId, {
                name: `interactive_${i}`,
                components: {
                    transform: { 
                        position: { 
                            x: (i - 1) * 3, 
                            y: 0, 
                            z: 0 
                        } 
                    }
                }
            });

            if (objectResponse.success) {
                await hd1API.addComponent(sessionId, objectResponse.entity_id, {
                    type: 'model',
                    properties: { type: 'box', material: { color: `#${Math.floor(Math.random()*16777215).toString(16)}` } }
                });

                createPlayCanvasEntity(objectResponse.entity_id, `interactive_${i}`);
            }
        }

        // Create ambient audio
        await hd1API.createAudioSource(sessionId, {
            name: 'ambient_sound',
            type: 'background',
            loop: true,
            volume: 0.2
        });

        updateGameStats();
        document.getElementById('content-experience').textContent = 'Interactive Demo';
        console.log('[HD1] Demo holodeck experience loaded');
    } catch (error) {
        console.error('[HD1] Failed to load demo:', error);
    }
}

/**
 * Setup dropdown control event listeners
 */
function setupDropdownControls() {
    // Scene selection (only if element exists)
    const sceneSelect = document.getElementById('debug-scene-select');
    if (sceneSelect) {
        sceneSelect.addEventListener('change', async function(e) {
            const scene = e.target.value;
            if (scene === 'interactive-demo') {
                await loadDemoScene();
            }
            e.target.value = '';
        });
    }

    // Object creation (only if element exists)
    const objectSelect = document.getElementById('debug-object-select');
    if (objectSelect) {
        objectSelect.addEventListener('change', async function(e) {
            const objectType = e.target.value;
            if (objectType) {
                await createEntity();
            }
            e.target.value = '';
        });
    }

    // Interaction control (removed during UI consolidation)
    // Previously debug-interaction-select - functionality moved to other controls
}

/**
 * Update holodeck statistics display
 */
function updateGameStats() {
    if (hd1GameEngine) {
        const entities = hd1GameEngine.root.children.filter(e => e.hd1Id);
        // Update consolidated panel (performance-entities was removed during UI consolidation)
        const objectsDisplay = document.getElementById('content-objects');
        if (objectsDisplay) {
            objectsDisplay.textContent = entities.length;
            // Check for overflow if checkTextOverflow function exists
            if (typeof checkTextOverflow === 'function') {
                checkTextOverflow(objectsDisplay);
            }
        }
    }

    // Update FPS in consolidated panel
    const fpsDisplay = document.getElementById('performance-fps');
    if (fpsDisplay) {
        fpsDisplay.textContent = hd1GameEngine ? Math.round(1000 / hd1GameEngine.stats.frame.ms) : 60;
    }
}

// Update stats every second
setInterval(updateGameStats, 1000);

/**
 * Create PlayCanvas object directly from server data
 * This method is called by the console when objects are received from the server
 */
function createObjectFromData(objectData) {
    console.log('[HD1] createObjectFromData called with:', objectData);
    
    if (!hd1GameEngine) {
        console.warn('[HD1] PlayCanvas engine not ready, storing object for later');
        return;
    }

    // Check if object already exists to prevent duplicates
    const entityId = objectData.entity_id || objectData.id;
    const existingEntity = hd1GameEngine.root.children.find(entity => 
        entity.name === objectData.name || entity.hd1Id === entityId
    );
    
    if (existingEntity) {
        console.log('[HD1] Object already exists, skipping:', objectData.name);
        return;
    }

    // Get model type from components
    const modelType = objectData.components?.model?.type || objectData.type || 'box';
    console.log('[HD1] Creating PlayCanvas object:', objectData.name, 'type:', modelType);

    const entity = new pc.Entity(objectData.name);
    entity.hd1Id = entityId;
    entity.hd1Tags = objectData.tags || []; // Store HD1 tags for WebSocket handler lookup

    // Add model component (only if not a light)
    if (!objectData.components?.light) {
        // Check if this entity has avatar tags to determine if it should load GLB assets
        let avatarTag = objectData.tags?.find(tag => tag.startsWith('avatar-'));
        let avatarType = avatarTag ? avatarTag.replace('avatar-', '') : null;
        
        // Special handling for session client avatars - check session avatar configuration
        if (!avatarType && objectData.tags?.includes('session-avatar') && objectData.name?.includes('session_client_')) {
            // This is a session client avatar, try to get the session's configured avatar type
            const sessionId = objectData.session_id || getCurrentSession();
            if (sessionId) {
                console.log('[HD1] Session client avatar detected, checking session avatar configuration for:', sessionId);
                
                // CRITICAL FIX: Select avatar type based on current world
                avatarType = getAvatarTypeForCurrentWorld();
                console.log('[HD1] Using world-based avatar type for session client:', avatarType);
            }
        }
        
        // Handle GLB asset loading for avatars via HTTP
        if (avatarType) {
            console.log('[HD1] Avatar entity detected, loading GLB asset via HTTP:', avatarType);
            
            // Store entity reference
            entity.hd1AvatarType = avatarType;
            entity.hd1WaitingForAsset = true;
            
            // Add entity to scene first
            hd1GameEngine.root.addChild(entity);
            
            // Load GLB asset directly via HTTP endpoint
            loadAvatarAssetHTTP(avatarType, entity);
            
            // Don't add model component yet - wait for GLB asset
            console.log('[HD1] Entity added to scene, waiting for GLB asset:', entity.name);
            return; // Skip the normal model component addition
        } else if (modelType === 'asset' && objectData.components?.model?.asset_path) {
            // Legacy asset loading (for non-avatar assets)
            console.warn('[HD1] Non-avatar asset entity, using primitive model fallback');
            entity.addComponent('model', { type: 'box' });
        } else {
            // Standard primitive model
            entity.addComponent('model', {
                type: modelType
            });
        }
    }
    
    // Add light component if entity has light
    if (objectData.components?.light) {
        const lightConfig = objectData.components.light;
        entity.addComponent('light', {
            type: lightConfig.type || 'directional',
            color: Array.isArray(lightConfig.color) ? 
                new pc.Color(lightConfig.color[0] || 1, lightConfig.color[1] || 1, lightConfig.color[2] || 1) :
                new pc.Color(1, 1, 1),
            intensity: lightConfig.intensity || 1,
            castShadows: lightConfig.castShadows || false
        });
        console.log('[HD1] Added light component:', lightConfig.type, 'intensity:', lightConfig.intensity);
    }

    // Set position from HD1 v3.0 API format
    const transform = objectData.components?.transform;
    if (transform?.position && Array.isArray(transform.position)) {
        entity.setPosition(
            transform.position[0] || 0,
            transform.position[1] || 0, 
            transform.position[2] || 0
        );
        console.log('[HD1] Setting position from v3.0 API:', transform.position);
    } else if (objectData.transform && objectData.transform.position) {
        // Legacy format support
        entity.setPosition(
            objectData.transform.position.x || 0,
            objectData.transform.position.y || 0,
            objectData.transform.position.z || 0
        );
    } else if (objectData.x !== undefined || objectData.y !== undefined || objectData.z !== undefined) {
        entity.setPosition(
            objectData.x || 0,
            objectData.y || 0,
            objectData.z || 0
        );
        console.log('[HD1] Setting position from legacy API format:', objectData.x, objectData.y, objectData.z);
    }

    // Set scale from HD1 v3.0 API format
    if (transform?.scale && Array.isArray(transform.scale)) {
        entity.setLocalScale(
            transform.scale[0] || 1,
            transform.scale[1] || 1,
            transform.scale[2] || 1
        );
        console.log('[HD1] Setting scale from v3.0 API:', transform.scale);
    } else if (objectData.transform && objectData.transform.scale) {
        // Legacy format support
        entity.setLocalScale(
            objectData.transform.scale.x || 1,
            objectData.transform.scale.y || 1,
            objectData.transform.scale.z || 1
        );
    } else {
        // Default scale for entities without scale data
        entity.setLocalScale(1, 1, 1);
    }

    // Set rotation from HD1 v3.0 API format
    if (transform?.rotation && Array.isArray(transform.rotation)) {
        entity.setEulerAngles(
            transform.rotation[0] || 0,
            transform.rotation[1] || 0,
            transform.rotation[2] || 0
        );
        console.log('[HD1] Setting rotation from v3.0 API:', transform.rotation);
    } else if (objectData.transform && objectData.transform.rotation) {
        // Legacy format support
        entity.setEulerAngles(
            objectData.transform.rotation.x || 0,
            objectData.transform.rotation.y || 0,
            objectData.transform.rotation.z || 0
        );
    }

    // Create material with color (only for model entities, not lights)
    if (!objectData.components?.light && entity.model) {
        const material = new pc.StandardMaterial();
        const materialConfig = objectData.components?.material;
        
        if (materialConfig?.diffuse) {
            // Parse hex color to RGB
            const hexColor = materialConfig.diffuse;
            if (hexColor.startsWith('#')) {
                const r = parseInt(hexColor.substr(1, 2), 16) / 255;
                const g = parseInt(hexColor.substr(3, 2), 16) / 255;
                const b = parseInt(hexColor.substr(5, 2), 16) / 255;
                material.diffuse = new pc.Color(r, g, b);
            }
        } else if (objectData.color) {
            // Legacy format support
            material.diffuse = new pc.Color(
                objectData.color.r || 0.5,
                objectData.color.g || 0.5,
                objectData.color.b || 0.5
            );
        } else {
            material.diffuse = new pc.Color(0.5, 0.5, 0.5);
        }

        // Set material properties from HD1 v3.0 API
        material.metalness = materialConfig?.metalness || 0.3;
        material.shininess = materialConfig?.shininess || 40;
        material.wireframe = objectData.wireframe || false;
        material.update();

        // Apply material to mesh instances
        entity.model.meshInstances.forEach(meshInstance => {
            meshInstance.material = material;
        });
    }

    hd1GameEngine.root.addChild(entity);
    updateGameStats();
    
    console.log('[HD1] PlayCanvas object created:', objectData.name);
}

/**
 * Delete PlayCanvas object by name
 */
function deleteObjectByName(objectName) {
    if (!hd1GameEngine) {
        console.warn('[HD1] PlayCanvas engine not ready, cannot delete object:', objectName);
        return false;
    }

    // 🛡️ AVATAR PROTECTION: Don't allow direct avatar deletion
    if (objectName && (objectName.includes('session_client_') || objectName.includes('session_'))) {
        console.warn(`[HD1] 🛡️ AVATAR PROTECTION: Blocked direct deletion of avatar object: ${objectName}`);
        return true; // Return true to indicate "handled" (blocked but processed)
    }

    console.log('[HD1] Deleting PlayCanvas object:', objectName);
    
    // Find and remove the entity
    const entityToRemove = hd1GameEngine.root.children.find(entity => 
        entity.name === objectName || entity.hd1Id === objectName
    );
    
    if (entityToRemove) {
        hd1GameEngine.root.removeChild(entityToRemove);
        entityToRemove.destroy();
        updateGameStats();
        console.log('[HD1] PlayCanvas object deleted:', objectName);
        return true;
    } else {
        // CRITICAL FIX: Limit repeated warnings using throttling
        if (!window.deletionWarnings) {
            window.deletionWarnings = new Map();
        }
        
        const warningKey = `missing_${objectName}`;
        const lastWarning = window.deletionWarnings.get(warningKey) || 0;
        const now = Date.now();
        
        if (now - lastWarning > 5000) { // Only warn every 5 seconds per entity
            console.warn('[HD1] PlayCanvas object not found for deletion:', objectName);
            window.deletionWarnings.set(warningKey, now);
        }
        return false;
    }
}

// Global storage for pending avatar asset requests
const pendingAvatarAssets = new Map(); // avatarType -> entity

/**
 * Load GLB avatar asset via HTTP endpoint
 * Uses HD1's API-first architecture with proper HTTP asset delivery
 */
function loadAvatarAssetHTTP(avatarType, entity) {
    console.log('[HD1] Loading avatar asset via HTTP with PlayCanvas v1.73.3:', avatarType);
    
    try {
        // Construct HTTP URL for avatar GLB asset
        const assetUrl = `/api/avatars/${avatarType}/asset`;
        
        // Use PlayCanvas loadFromUrlAndFilename for proper GLB loading
        hd1GameEngine.assets.loadFromUrlAndFilename(
            assetUrl,
            `${avatarType}.glb`,
            'container',
            function(err, asset) {
                if (err) {
                    console.error('[HD1] Avatar GLB asset loading failed:', err);
                    console.log('[HD1] Falling back to colored sphere for entity:', entity.name);
                    loadAvatarFallback(avatarType, entity);
                    return;
                }
                
                console.log('[HD1] Avatar GLB asset loaded successfully with v1.73.3:', avatarType);
                console.log('[HD1] Asset resource:', asset.resource);
                
                try {
                    // Remove existing model component if any
                    if (entity.model) {
                        entity.removeComponent('model');
                    }
                    
                    // Apply the loaded asset using the resource.model property
                    entity.addComponent('model', {
                        asset: asset.resource.model
                    });
                    
                    console.log('[HD1] Avatar GLB model applied to entity:', entity.name);
                    entity.hd1WaitingForAsset = false;
                    
                    // Simple avatar loading (no physics config)
                    
                } catch (componentError) {
                    console.error('[HD1] Error applying avatar GLB model component:', componentError);
                    // Fallback to colored sphere
                    loadAvatarFallback(avatarType, entity);
                }
            }
        );
        
        console.log('[HD1] Started loading avatar asset with loadFromUrlAndFilename:', assetUrl);
        
    } catch (error) {
        console.error('[HD1] Error in loadAvatarAssetHTTP:', error);
        loadAvatarFallback(avatarType, entity);
    }
}

// Fallback function for avatar loading
function loadAvatarFallback(avatarType, entity) {
    try {
        // Remove existing model component if any
        if (entity.model) {
            entity.removeComponent('model');
        }
        
        // Create a basic material first
        const material = new pc.StandardMaterial();
        
        // Set avatar-specific colors
        if (avatarType === 'claude_avatar') {
            material.diffuse = new pc.Color(0.2, 0.4, 0.8); // Blue for Claude
        } else if (avatarType === 'human_avatar') {
            material.diffuse = new pc.Color(0.8, 0.4, 0.2); // Orange for human
        } else {
            material.diffuse = new pc.Color(0.5, 0.5, 0.5); // Gray for unknown
        }
        
        material.update();
        
        // Use a sphere as avatar representation with proper material
        entity.addComponent('model', { 
            type: 'sphere'
        });
        
        // Apply the material after component is added
        if (entity.model && entity.model.model) {
            entity.model.model.meshInstances[0].material = material;
        }
        
        console.log('[HD1] Applied fallback avatar model:', avatarType);
        entity.hd1WaitingForAsset = false;
        
        // Simple avatar loading (no physics config)
        
    } catch (error) {
        console.error('[HD1] Error in fallback avatar loading:', error);
        
        // Ultimate fallback to box
        if (entity.model) {
            entity.removeComponent('model');
        }
        
        entity.addComponent('model', { type: 'box' });
        entity.hd1WaitingForAsset = false;
    }
}



/**
 * Handle incoming avatar asset data from WebSocket
 * Converts base64 data to ArrayBuffer and creates PlayCanvas GLB asset
 */
function handleAvatarAssetResponse(avatarType, base64Data) {
    console.log('[HD1] Processing avatar asset response:', avatarType, 'pending entities:', pendingAvatarAssets.size);
    
    const entity = pendingAvatarAssets.get(avatarType);
    if (!entity) {
        console.warn('[HD1] Received avatar asset but no pending entity found:', avatarType);
        console.log('[HD1] Available pending assets:', Array.from(pendingAvatarAssets.keys()));
        return;
    }
    
    try {
        // Convert base64 to ArrayBuffer
        const binaryString = atob(base64Data);
        const arrayBuffer = new ArrayBuffer(binaryString.length);
        const uint8Array = new Uint8Array(arrayBuffer);
        for (let i = 0; i < binaryString.length; i++) {
            uint8Array[i] = binaryString.charCodeAt(i);
        }
        
        console.log('[HD1] Converted GLB data to ArrayBuffer:', arrayBuffer.byteLength, 'bytes');
        
        // Create PlayCanvas container asset and set binary data directly
        const asset = new pc.Asset(entity.name + '_model', 'container');
        
        // Add asset to registry
        hd1GameEngine.assets.add(asset);
        
        // Set the GLB binary data directly on the asset
        asset.data = arrayBuffer;
        
        // Use PlayCanvas's built-in GLB loader through the asset system with proper binary handling
        try {
            // Create a temporary object URL for the GLB data
            const blob = new Blob([arrayBuffer], { type: 'application/octet-stream' });
            const tempUrl = URL.createObjectURL(blob);
            
            // Create a new asset with the temporary URL
            const tempAsset = new pc.Asset(entity.name + '_temp_model', 'container', {
                url: tempUrl
            });
            
            // Add to asset registry
            hd1GameEngine.assets.add(tempAsset);
            
            // Set up asset ready handler
            tempAsset.ready((loadedAsset) => {
                console.log('[HD1] GLB asset loaded successfully, applying to entity:', entity.name);
                
                try {
                    // Remove existing model component if any
                    if (entity.model) {
                        entity.removeComponent('model');
                    }
                    
                    // Apply the loaded asset to the entity
                    entity.addComponent('model', {
                        type: 'asset',
                        asset: loadedAsset
                    });
                    
                    console.log('[HD1] GLB asset applied to entity:', entity.name);
                    
                    // Clean up
                    URL.revokeObjectURL(tempUrl);
                    entity.hd1WaitingForAsset = false;
                    pendingAvatarAssets.delete(avatarType);
                    
                } catch (componentError) {
                    console.error('[HD1] Error applying GLB model component:', componentError);
                    // Fallback to primitive
                    entity.addComponent('model', { type: 'box' });
                    URL.revokeObjectURL(tempUrl);
                    entity.hd1WaitingForAsset = false;
                    pendingAvatarAssets.delete(avatarType);
                }
            });
            
            // Set up error handler
            tempAsset.on('error', (err, asset) => {
                console.error('[HD1] GLB asset loading failed:', err);
                console.log('[HD1] Falling back to primitive model for entity:', entity.name);
                
                // Fallback to primitive
                if (entity.model) {
                    entity.removeComponent('model');
                } else {
                    entity.addComponent('model', { type: 'box' });
                }
                
                // Clean up
                URL.revokeObjectURL(tempUrl);
                entity.hd1WaitingForAsset = false;
                pendingAvatarAssets.delete(avatarType);
            });
            
            // Load the asset
            hd1GameEngine.assets.load(tempAsset);
            
        } catch (parseError) {
            console.error('[HD1] GLB asset processing error:', parseError);
            console.log('[HD1] Falling back to primitive model for entity:', entity.name);
            
            // Fallback to primitive
            if (entity.model) {
                entity.removeComponent('model');
            }
            entity.addComponent('model', { type: 'box' });
            
            // Clean up after error
            entity.hd1WaitingForAsset = false;
            pendingAvatarAssets.delete(avatarType);
        }
        console.log('[HD1] GLB asset loading started via WebSocket:', avatarType);
        
    } catch (error) {
        console.error('[HD1] Failed to process GLB asset:', error);
        
        // Fallback to primitive
        if (entity.model) {
            entity.removeComponent('model');
        }
        entity.addComponent('model', { type: 'box' });
        entity.hd1WaitingForAsset = false;
        pendingAvatarAssets.delete(avatarType);
    }
}

// Add the createObject and deleteObject methods to the engine
if (typeof window !== 'undefined') {
    window.addEventListener('DOMContentLoaded', function() {
        setTimeout(() => {
            if (window.hd1GameEngine) {
                window.hd1GameEngine.createObject = createObjectFromData;
                window.hd1GameEngine.deleteObject = deleteObjectByName;
                console.log('[HD1] Added createObject and deleteObject methods to PlayCanvas engine');
                
                // Add avatar asset handler to global scope for WebSocket manager
                window.handleAvatarAssetResponse = handleAvatarAssetResponse;
                console.log('[HD1] Added avatar asset handler to global scope');
                
                // Process any pending objects
                if (window.pendingObjects && window.pendingObjects.length > 0) {
                    console.log('[HD1] Processing', window.pendingObjects.length, 'pending objects');
                    window.pendingObjects.forEach(obj => {
                        createObjectFromData(obj);
                    });
                    window.pendingObjects = null;
                }
            }
        }, 100);
    });
}

/**
 * Load existing objects from the current session and all channel participants into PlayCanvas
 */
async function loadExistingSessionObjects() {
    console.log('[HD1] loadExistingSessionObjects() called');
    
    const sessionId = getCurrentSession();
    if (!sessionId) {
        console.warn('[HD1] No active session found - skipping object loading');
        console.log('[HD1] localStorage hd1_session_id:', localStorage.getItem('hd1_session_id'));
        return;
    }
    
    if (!hd1GameEngine) {
        console.warn('[HD1] PlayCanvas engine not ready - deferring object loading');
        return;
    }
    
    try {
        console.log(`[HD1] Loading existing entities from session: ${sessionId}`);
        
        // Load entities from current session
        const response = await fetch(`/api/sessions/${sessionId}/entities`);
        const data = await response.json();
        
        let totalEntitiesLoaded = 0;
        
        if (data.entities && Array.isArray(data.entities)) {
            console.log(`[HD1] Found ${data.entities.length} existing entities in current session`);
            
            data.entities.forEach(entity => {
                console.log('[HD1] Loading existing entity:', entity.name);
                createObjectFromData(entity);
            });
            
            totalEntitiesLoaded += data.entities.length;
        }
        
        // Load avatars from ALL sessions in the same channel for multi-user visibility
        try {
            // Get current session info to find its channel
            const sessionResponse = await fetch(`/api/sessions/${sessionId}`);
            const sessionData = await sessionResponse.json();
            
            if (sessionData.channel_id) {
                console.log(`[HD1] Loading avatars from all sessions in channel: ${sessionData.channel_id}`);
                
                // Get all sessions in the same channel
                const allSessionsResponse = await fetch('/api/sessions');
                const allSessionsData = await allSessionsResponse.json();
                
                if (allSessionsData.sessions && Array.isArray(allSessionsData.sessions)) {
                    const channelSessions = allSessionsData.sessions.filter(session => 
                        session.channel_id === sessionData.channel_id && session.id !== sessionId
                    );
                    
                    console.log(`[HD1] Found ${channelSessions.length} other sessions in channel ${sessionData.channel_id}`);
                    
                    // Load avatars from each session in the channel
                    for (const session of channelSessions) {
                        try {
                            const otherSessionResponse = await fetch(`/api/sessions/${session.id}/entities?tag=session-avatar`);
                            const otherSessionData = await otherSessionResponse.json();
                            
                            if (otherSessionData.entities && Array.isArray(otherSessionData.entities)) {
                                otherSessionData.entities.forEach(avatar => {
                                    console.log(`[HD1] Loading avatar from session ${session.id}:`, avatar.name);
                                    createObjectFromData(avatar);
                                    totalEntitiesLoaded++;
                                });
                            }
                        } catch (error) {
                            console.warn(`[HD1] Failed to load avatars from session ${session.id}:`, error);
                        }
                    }
                }
            } else {
                console.log('[HD1] Session not in a channel, skipping cross-session avatar loading');
            }
        } catch (error) {
            console.warn('[HD1] Failed to load cross-session avatars:', error);
        }
        
        if (totalEntitiesLoaded > 0) {
            updateGameStats();
            console.log(`[HD1] Successfully loaded ${totalEntitiesLoaded} total entities/avatars into PlayCanvas`);
        } else {
            console.log('[HD1] No existing entities found to load');
        }
    } catch (error) {
        console.error('[HD1] Failed to load existing session entities:', error);
    }
}

// Add rebootstrap and manual object loading functions to global scope
window.triggerRebootstrap = function() {
    console.log('[HD1] Triggering rebootstrap...');
    localStorage.removeItem('hd1_session_id');
    setTimeout(() => {
        window.location.reload(true);
    }, 1000);
};

/**
 * Get avatar type based on current world
 * Returns appropriate avatar type for the active world/channel
 */
function getAvatarTypeForCurrentWorld() {
    console.log('[HD1] Determining avatar type for current world');
    
    // Method 1: Get from localStorage (primary source) 
    const currentWorldId = localStorage.getItem('hd1_current_world');
    if (currentWorldId) {
        const avatarType = getAvatarTypeForWorld(currentWorldId);
        console.log(`[HD1] Avatar type from localStorage world ${currentWorldId}: ${avatarType}`);
        return avatarType;
    }
    
    // Method 2: Get from session storage (for world transitions)
    const sessionWorldId = localStorage.getItem('hd1_session_world');
    if (sessionWorldId) {
        const avatarType = getAvatarTypeForWorld(sessionWorldId);
        console.log(`[HD1] Avatar type from session world ${sessionWorldId}: ${avatarType}`);
        return avatarType;
    }
    
    // Method 3: Get from console manager world selection
    if (window.hd1ConsoleManager) {
        const worldManager = window.hd1ConsoleManager.getModule('world');
        if (worldManager && worldManager.currentWorld && worldManager.currentWorld.id) {
            const avatarType = getAvatarTypeForWorld(worldManager.currentWorld.id);
            console.log(`[HD1] Avatar type from world manager ${worldManager.currentWorld.id}: ${avatarType}`);
            return avatarType;
        }
    }
    
    // Method 4: Get from UI selector 
    const selector = document.getElementById('channel-selector');
    if (selector && selector.value) {
        const avatarType = getAvatarTypeForWorld(selector.value);
        console.log(`[HD1] Avatar type from UI selector ${selector.value}: ${avatarType}`);
        return avatarType;
    }
    
    // Last fallback: default avatar
    console.warn('[HD1] Could not determine current world, using default avatar');
    return 'claude_avatar';
}

/**
 * Map world IDs to avatar types
 */
function getAvatarTypeForWorld(worldId) {
    const worldAvatarMap = {
        'world_one': 'humanoid_avatar',   // World 1: Humanoid avatar
        'world_two': 'fox_avatar',        // World 2: Fox avatar  
        'world_three': 'human_avatar',    // World 3: Human avatar
        'channel_one': 'humanoid_avatar', // Legacy channel support
        'channel_two': 'fox_avatar',
        'channel_three': 'human_avatar'
    };
    
    const avatarType = worldAvatarMap[worldId] || 'claude_avatar';
    console.log(`[HD1] World ${worldId} mapped to avatar type: ${avatarType}`);
    return avatarType;
}

// Allow manual triggering of object loading
window.loadSessionObjects = loadExistingSessionObjects;

console.log('[HD1] PlayCanvas integration loaded');