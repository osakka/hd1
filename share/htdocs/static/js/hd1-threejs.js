/**
 * HD1 Three.js Scene Manager
 * Pure Three.js implementation with zero abstraction layers
 */

import * as THREE from 'https://cdn.jsdelivr.net/npm/three@0.160.0/build/three.module.js';

class HD1ThreeJS {
    constructor(canvasId = 'hd1-canvas') {
        this.canvas = document.getElementById(canvasId);
        if (!this.canvas) {
            this.canvas = document.createElement('canvas');
            this.canvas.id = canvasId;
            document.body.appendChild(this.canvas);
        }
        
        // Core Three.js components
        this.scene = new THREE.Scene();
        this.camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
        this.renderer = new THREE.WebGLRenderer({ canvas: this.canvas, antialias: true });
        
        // Object tracking
        this.objects = new Map();      // entity_id -> THREE.Object3D
        this.avatars = new Map();      // session_id -> THREE.Object3D
        this.materials = new Map();    // material_id -> THREE.Material
        this.geometries = new Map();   // geometry_id -> THREE.Geometry
        
        // Font loading
        this.fontLoader = null;
        this.loadedFont = null;
        this.pendingTextEntities = new Map(); // Store entities waiting for font
        this.loadFontModules();
        
        // Camera controls
        this.controls = null;
        this.cameraTarget = new THREE.Vector3(0, 0, 0);
        this.lastTime = 0;
        
        // Initialize scene
        this.setupRenderer();
        this.setupLighting();
        this.setupCamera();
        this.setupEventListeners();
        
        // Start render loop
        this.animate();
        
        // Make globally accessible for sync operations
        window.hd1ThreeJS = this;
        
        console.log('[HD1-ThreeJS] Scene manager initialized');
        console.log('[HD1-ThreeJS] Three.js version:', THREE.REVISION);
        
        // Request full sync once Three.js is ready
        this.requestFullSyncWhenReady();
    }
    
    async loadFontModules() {
        try {
            // Use simplified approach - create geometric text using boxes
            console.log('[HD1-ThreeJS] Using simplified text rendering');
            
            // Mark font as "loaded" to enable text processing
            this.loadedFont = true;
            this.TextGeometry = null; // Will use fallback rendering
            
            // Process any pending text geometries
            this.processPendingTextGeometries();
        } catch (error) {
            console.error('[HD1-ThreeJS] Failed to load font modules:', error);
            console.log('[HD1-ThreeJS] Text rendering will use placeholder boxes');
        }
    }

    loadFont() {
        if (!this.FontLoader) {
            console.log('[HD1-ThreeJS] FontLoader not available, using placeholder');
            return;
        }
        
        this.fontLoader = new this.FontLoader();
        
        // Load Helvetiker font (commonly used Three.js font)
        this.fontLoader.load(
            'https://threejs.org/examples/fonts/helvetiker_regular.typeface.json',
            (font) => {
                this.loadedFont = font;
                console.log('[HD1-ThreeJS] Font loaded successfully');
                
                // Process any pending text geometries
                this.processPendingTextGeometries();
            },
            (progress) => {
                console.log('[HD1-ThreeJS] Font loading progress:', progress);
            },
            (error) => {
                console.error('[HD1-ThreeJS] Font loading failed:', error);
                console.log('[HD1-ThreeJS] Falling back to placeholder text rendering');
            }
        );
    }
    
    processPendingTextGeometries() {
        // Process any text entities that were created before the font loaded
        console.log('[HD1-ThreeJS] Processing pending text geometries:', this.pendingTextEntities.size);
        
        this.pendingTextEntities.forEach((entityData, entityId) => {
            // Recreate the entity with proper text geometry
            this.recreateTextEntity(entityId, entityData);
        });
        
        this.pendingTextEntities.clear();
    }
    
    recreateTextEntity(entityId, entityData) {
        // Remove the placeholder entity if it exists
        const existingEntity = this.objects.get(entityId);
        if (existingEntity) {
            this.scene.remove(existingEntity);
            
            // Dispose of old geometry and material
            if (existingEntity.geometry) existingEntity.geometry.dispose();
            if (existingEntity.material) existingEntity.material.dispose();
        }
        
        // Create new entity with proper text geometry
        const geometry = this.createGeometry(entityData.geometry);
        const material = this.createMaterial(entityData.material);
        
        const entity = new THREE.Mesh(geometry, material);
        entity.castShadow = true;
        entity.receiveShadow = true;
        
        // Apply transform
        if (entityData.position) {
            entity.position.set(entityData.position.x, entityData.position.y, entityData.position.z);
        }
        if (entityData.rotation) {
            entity.rotation.set(entityData.rotation.x, entityData.rotation.y, entityData.rotation.z);
        }
        if (entityData.scale) {
            entity.scale.set(entityData.scale.x, entityData.scale.y, entityData.scale.z);
        }
        
        // Add to scene and track
        this.scene.add(entity);
        this.objects.set(entityId, entity);
        
        console.log('[HD1-ThreeJS] Text entity recreated with proper geometry:', entityId);
    }
    
    setupRenderer() {
        this.renderer.setSize(window.innerWidth, window.innerHeight);
        this.renderer.setPixelRatio(window.devicePixelRatio);
        this.renderer.shadowMap.enabled = true;
        this.renderer.shadowMap.type = THREE.PCFSoftShadowMap;
        this.renderer.outputColorSpace = THREE.SRGBColorSpace;
        this.renderer.toneMapping = THREE.ACESFilmicToneMapping;
    }
    
    setupLighting() {
        // Ambient light
        const ambientLight = new THREE.AmbientLight(0x404040, 0.4);
        this.scene.add(ambientLight);
        
        // Directional light with shadows
        const directionalLight = new THREE.DirectionalLight(0xffffff, 0.8);
        directionalLight.position.set(10, 10, 5);
        directionalLight.castShadow = true;
        directionalLight.shadow.camera.near = 0.1;
        directionalLight.shadow.camera.far = 50;
        directionalLight.shadow.camera.left = -10;
        directionalLight.shadow.camera.right = 10;
        directionalLight.shadow.camera.top = 10;
        directionalLight.shadow.camera.bottom = -10;
        directionalLight.shadow.mapSize.width = 2048;
        directionalLight.shadow.mapSize.height = 2048;
        this.scene.add(directionalLight);
        
        // Grid helper
        const gridHelper = new THREE.GridHelper(20, 20, 0x444444, 0x444444);
        this.scene.add(gridHelper);
    }
    
    setupCamera() {
        // FPS camera setup
        this.camera.position.set(0, 1.6, 0); // Human eye height
        this.camera.rotation.order = 'YXZ';   // Yaw-Pitch-Roll order
        
        // FPS controls state
        this.keys = {};
        this.mouseLook = false;
        this.mouseX = 0;
        this.mouseY = 0;
        this.pitch = 0;
        this.yaw = 0;
        this.moveSpeed = 5.0;
        this.sprintMultiplier = 2.0;
        this.mouseSensitivity = 0.002;
        
        // Movement direction vectors
        this.moveForward = false;
        this.moveBackward = false;
        this.moveLeft = false;
        this.moveRight = false;
        this.sprint = false;
        
        // FPS controls
        this.enableFPSControls();
    }
    
    enableFPSControls() {
        // Keyboard controls
        document.addEventListener('keydown', (event) => {
            if (event.target.tagName === 'INPUT' || event.target.tagName === 'TEXTAREA') return;
            
            switch (event.code) {
                case 'KeyW':
                    this.moveForward = true;
                    break;
                case 'KeyS':
                    this.moveBackward = true;
                    break;
                case 'KeyA':
                    this.moveLeft = true;
                    break;
                case 'KeyD':
                    this.moveRight = true;
                    break;
                case 'ShiftLeft':
                    this.sprint = true;
                    break;
                case 'Escape':
                    // Only exit pointer lock on Escape - can't enter without user gesture
                    if (this.mouseLook) {
                        this.exitPointerLock();
                    }
                    break;
            }
        });
        
        document.addEventListener('keyup', (event) => {
            switch (event.code) {
                case 'KeyW':
                    this.moveForward = false;
                    break;
                case 'KeyS':
                    this.moveBackward = false;
                    break;
                case 'KeyA':
                    this.moveLeft = false;
                    break;
                case 'KeyD':
                    this.moveRight = false;
                    break;
                case 'ShiftLeft':
                    this.sprint = false;
                    break;
            }
        });
        
        // Mouse controls
        this.canvas.addEventListener('click', () => {
            this.requestPointerLock();
        });
        
        document.addEventListener('pointerlockchange', () => {
            this.mouseLook = document.pointerLockElement === this.canvas;
            
            // Update mouse look indicator
            if (window.setMouselookStatus) {
                window.setMouselookStatus(this.mouseLook);
            }
        });
        
        document.addEventListener('mousemove', (event) => {
            if (!this.mouseLook) return;
            
            const movementX = event.movementX || 0;
            const movementY = event.movementY || 0;
            
            this.yaw -= movementX * this.mouseSensitivity;
            this.pitch -= movementY * this.mouseSensitivity;
            
            // Limit pitch to prevent flipping
            this.pitch = Math.max(-Math.PI / 2, Math.min(Math.PI / 2, this.pitch));
            
            // Update camera rotation
            this.camera.rotation.set(this.pitch, this.yaw, 0);
        });
        
        // Touch controls for mobile
        this.canvas.addEventListener('touchstart', (event) => {
            event.preventDefault();
            this.touchStartX = event.touches[0].clientX;
            this.touchStartY = event.touches[0].clientY;
        });
        
        this.canvas.addEventListener('touchmove', (event) => {
            event.preventDefault();
            if (!this.touchStartX || !this.touchStartY) return;
            
            const touchX = event.touches[0].clientX;
            const touchY = event.touches[0].clientY;
            
            const deltaX = touchX - this.touchStartX;
            const deltaY = touchY - this.touchStartY;
            
            // Right side of screen: camera look
            if (this.touchStartX > window.innerWidth / 2) {
                this.yaw -= deltaX * this.mouseSensitivity;
                this.pitch -= deltaY * this.mouseSensitivity;
                
                // Limit pitch to prevent flipping
                this.pitch = Math.max(-Math.PI / 2, Math.min(Math.PI / 2, this.pitch));
                
                // Update camera rotation
                this.camera.rotation.set(this.pitch, this.yaw, 0);
            }
            // Left side of screen: movement
            else {
                const moveThreshold = 30;
                this.moveForward = deltaY < -moveThreshold;
                this.moveBackward = deltaY > moveThreshold;
                this.moveLeft = deltaX < -moveThreshold;
                this.moveRight = deltaX > moveThreshold;
            }
            
            this.touchStartX = touchX;
            this.touchStartY = touchY;
        });
        
        this.canvas.addEventListener('touchend', () => {
            this.touchStartX = null;
            this.touchStartY = null;
            // Stop movement on touch end
            this.moveForward = false;
            this.moveBackward = false;
            this.moveLeft = false;
            this.moveRight = false;
        });
    }
    
    requestPointerLock() {
        if (this.canvas.requestPointerLock) {
            this.canvas.requestPointerLock();
        }
    }
    
    exitPointerLock() {
        if (document.exitPointerLock) {
            document.exitPointerLock();
        }
    }
    
    updateMovement(deltaTime) {
        if (!this.mouseLook && !this.touchStartX) {
            // Debug: Show why movement isn't working
            if (this.moveForward || this.moveBackward || this.moveLeft || this.moveRight) {
                console.log('[HD1-ThreeJS] Movement blocked - click canvas or touch to enable!');
            }
            return; // Only move when mouse look is active or touch is active
        }
        
        // Calculate movement speed with sprint modifier
        const speed = this.moveSpeed * (this.sprint ? this.sprintMultiplier : 1.0) * deltaTime;
        
        // Get camera direction vectors
        const forward = new THREE.Vector3(0, 0, -1).applyQuaternion(this.camera.quaternion);
        const right = new THREE.Vector3(1, 0, 0).applyQuaternion(this.camera.quaternion);
        
        // Calculate movement vector
        const moveVector = new THREE.Vector3();
        
        if (this.moveForward) moveVector.add(forward.clone().multiplyScalar(speed));
        if (this.moveBackward) moveVector.add(forward.clone().multiplyScalar(-speed));
        if (this.moveLeft) moveVector.add(right.clone().multiplyScalar(-speed));
        if (this.moveRight) moveVector.add(right.clone().multiplyScalar(speed));
        
        // Apply movement to camera
        this.camera.position.add(moveVector);
        
        // Send position update to server if movement occurred
        if (moveVector.length() > 0) {
            console.log('[HD1-ThreeJS] Movement detected:', {
                position: {x: this.camera.position.x, y: this.camera.position.y, z: this.camera.position.z},
                apiClient: !!window.apiClient,
                clientId: window.clientId
            });
            this.sendAvatarPosition();
        }
    }
    
    sendAvatarPosition() {
        // Send avatar position update via API endpoint - SINGLE SOURCE OF TRUTH
        if (!window.apiClient) {
            // Client not ready yet - silently skip without warning
            return;
        }
        
        if (!window.clientId) {
            // Client ID not available yet - silently skip without warning
            return;
        }
        
        if (window.apiClient && window.clientId) {
            const positionData = {
                position: {
                    x: this.camera.position.x,
                    y: this.camera.position.y,
                    z: this.camera.position.z
                },
                rotation: {
                    x: this.camera.rotation.x,
                    y: this.camera.rotation.y,
                    z: this.camera.rotation.z
                }
            };
            
            // Use auto-generated API client to call /avatars/{sessionId}/move
            window.apiClient.moveAvatar(window.clientId, positionData)
                .then(response => {
                    console.log('[HD1-ThreeJS] Avatar moved:', {
                        position: positionData.position,
                        seq_num: response.seq_num,
                        client_id: window.clientId
                    });
                })
                .catch(error => {
                    console.warn('[HD1-ThreeJS] Avatar move failed:', error);
                });
        }
    }
    
    setupEventListeners() {
        window.addEventListener('resize', () => {
            this.camera.aspect = window.innerWidth / window.innerHeight;
            this.camera.updateProjectionMatrix();
            this.renderer.setSize(window.innerWidth, window.innerHeight);
        });
    }
    
    animate(currentTime = 0) {
        requestAnimationFrame((time) => this.animate(time));
        
        // Calculate delta time
        const deltaTime = this.lastTime ? (currentTime - this.lastTime) / 1000 : 0;
        this.lastTime = currentTime;
        
        // Update movement
        this.updateMovement(deltaTime);
        
        this.renderer.render(this.scene, this.camera);
    }
    
    // Operation handlers - called by sync system
    applyOperation(operation) {
        switch (operation.type) {
            case 'avatar_move':
                this.updateAvatar(operation.data.hd1_id || operation.client_id, operation.data);
                break;
            case 'avatar_create':
                this.createAvatar(operation.data.hd1_id || operation.client_id, operation.data);
                break;
            case 'avatar_remove':
                this.removeAvatar(operation.data.hd1_id || operation.client_id);
                break;
            case 'entity_create':
                this.createEntity(operation.data.id, operation.data);
                break;
            case 'entity_update':
                this.updateEntity(operation.data.id, operation.data);
                break;
            case 'entity_delete':
                this.deleteEntity(operation.data.id);
                break;
            case 'scene_update':
                this.updateScene(operation.data);
                break;
            default:
                console.warn('[HD1-ThreeJS] Unknown operation type:', operation.type);
        }
    }
    
    updateAvatar(sessionId, data) {
        let avatar = this.avatars.get(sessionId);
        
        if (!avatar) {
            // Avatar doesn't exist - this shouldn't happen with proper operation ordering
            console.log('[HD1-ThreeJS] Avatar not found for movement, creating:', sessionId);
            avatar = this.createAvatar(sessionId, data);
            this.avatars.set(sessionId, avatar);
        }
        
        // Update position
        if (data.position) {
            avatar.position.set(data.position.x, data.position.y, data.position.z);
        }
        
        // Update rotation
        if (data.rotation) {
            avatar.rotation.set(data.rotation.x, data.rotation.y, data.rotation.z);
        }
        
        // Update animation
        if (data.animation) {
            this.updateAvatarAnimation(avatar, data.animation);
        }
    }
    
    createAvatar(sessionId, data) {
        // Create avatar geometry
        const geometry = new THREE.CapsuleGeometry(0.3, 1.8, 4, 8);
        
        // Create avatar material with session-based color
        const color = this.getAvatarColor(sessionId);
        const material = new THREE.MeshPhongMaterial({ color: color });
        
        // Create avatar mesh
        const avatar = new THREE.Mesh(geometry, material);
        avatar.castShadow = true;
        avatar.receiveShadow = true;
        
        // Set initial position
        if (data.position) {
            avatar.position.set(data.position.x, data.position.y, data.position.z);
        }
        
        // Add to scene
        this.scene.add(avatar);
        
        console.log('[HD1-ThreeJS] Avatar created:', sessionId);
        return avatar;
    }
    
    getAvatarColor(sessionId) {
        // Generate consistent color from session ID
        let hash = 0;
        for (let i = 0; i < sessionId.length; i++) {
            hash = sessionId.charCodeAt(i) + ((hash << 5) - hash);
        }
        
        const hue = (hash % 360) / 360;
        return new THREE.Color().setHSL(hue, 0.7, 0.5);
    }
    
    removeAvatar(sessionId) {
        const avatar = this.avatars.get(sessionId);
        if (avatar) {
            this.scene.remove(avatar);
            this.avatars.delete(sessionId);
            
            // Clean up geometry and material
            if (avatar.geometry) avatar.geometry.dispose();
            if (avatar.material) avatar.material.dispose();
            
            console.log('[HD1-ThreeJS] Avatar removed:', sessionId);
        }
    }
    
    updateAvatarAnimation(avatar, animation) {
        // Simple animation system
        switch (animation) {
            case 'walk':
                // Simple bob animation
                avatar.userData.walkAnimation = true;
                break;
            case 'idle':
                avatar.userData.walkAnimation = false;
                break;
        }
    }
    
    createEntity(id, data) {
        // If this is a text entity and font is not loaded, store for later
        if (data.geometry && data.geometry.type === 'text' && !this.loadedFont) {
            console.log('[HD1-ThreeJS] Text entity created before font loaded, storing for later:', id);
            this.pendingTextEntities.set(id, data);
        }
        
        // Create geometry
        const geometry = this.createGeometry(data.geometry);
        
        // Create material
        const material = this.createMaterial(data.material);
        
        // Create mesh
        const entity = new THREE.Mesh(geometry, material);
        entity.castShadow = true;
        entity.receiveShadow = true;
        
        // Apply transform
        if (data.position) {
            entity.position.set(data.position.x, data.position.y, data.position.z);
        }
        if (data.rotation) {
            entity.rotation.set(data.rotation.x, data.rotation.y, data.rotation.z);
        }
        if (data.scale) {
            entity.scale.set(data.scale.x, data.scale.y, data.scale.z);
        }
        
        // Add to scene and track
        this.scene.add(entity);
        this.objects.set(id, entity);
        
        console.log('[HD1-ThreeJS] Entity created:', id);
    }
    
    updateEntity(id, data) {
        const entity = this.objects.get(id);
        if (!entity) return;
        
        // Update transform
        if (data.position) {
            entity.position.set(data.position.x, data.position.y, data.position.z);
        }
        if (data.rotation) {
            entity.rotation.set(data.rotation.x, data.rotation.y, data.rotation.z);
        }
        if (data.scale) {
            entity.scale.set(data.scale.x, data.scale.y, data.scale.z);
        }
        
        // Update visibility
        if (data.visible !== undefined) {
            entity.visible = data.visible;
        }
    }
    
    deleteEntity(id) {
        const entity = this.objects.get(id);
        if (!entity) return;
        
        this.scene.remove(entity);
        this.objects.delete(id);
        
        // Clean up geometry and material
        if (entity.geometry) entity.geometry.dispose();
        if (entity.material) entity.material.dispose();
        
        console.log('[HD1-ThreeJS] Entity deleted:', id);
    }
    
    updateScene(data) {
        // Update scene-level properties
        if (data.background) {
            this.scene.background = new THREE.Color(data.background);
        }
        
        if (data.fog) {
            this.scene.fog = new THREE.Fog(data.fog.color, data.fog.near, data.fog.far);
        }
    }
    
    createHDLetterGeometry(size = 1, depth = 0.1) {
        // Create HD letters using box geometry shapes
        console.log('[HD1-ThreeJS] Creating HD letter geometry');
        
        const group = new THREE.Group();
        
        // H letter (left side)
        const hLeft = new THREE.BoxGeometry(size * 0.2, size, depth);
        const hMiddle = new THREE.BoxGeometry(size * 0.6, size * 0.2, depth);
        const hRight = new THREE.BoxGeometry(size * 0.2, size, depth);
        
        const hLeftMesh = new THREE.Mesh(hLeft);
        const hMiddleMesh = new THREE.Mesh(hMiddle);
        const hRightMesh = new THREE.Mesh(hRight);
        
        hLeftMesh.position.set(-size * 0.5, 0, 0);
        hMiddleMesh.position.set(-size * 0.1, 0, 0);
        hRightMesh.position.set(size * 0.3, 0, 0);
        
        // D letter (right side)
        const dLeft = new THREE.BoxGeometry(size * 0.2, size, depth);
        const dTop = new THREE.BoxGeometry(size * 0.6, size * 0.2, depth);
        const dBottom = new THREE.BoxGeometry(size * 0.6, size * 0.2, depth);
        const dRight = new THREE.BoxGeometry(size * 0.2, size * 0.6, depth);
        
        const dLeftMesh = new THREE.Mesh(dLeft);
        const dTopMesh = new THREE.Mesh(dTop);
        const dBottomMesh = new THREE.Mesh(dBottom);
        const dRightMesh = new THREE.Mesh(dRight);
        
        dLeftMesh.position.set(size * 0.8, 0, 0);
        dTopMesh.position.set(size * 1.1, size * 0.4, 0);
        dBottomMesh.position.set(size * 1.1, -size * 0.4, 0);
        dRightMesh.position.set(size * 1.4, 0, 0);
        
        // Add all parts to group
        group.add(hLeftMesh, hMiddleMesh, hRightMesh);
        group.add(dLeftMesh, dTopMesh, dBottomMesh, dRightMesh);
        
        // Convert group to single geometry
        const geometry = new THREE.BufferGeometry();
        group.updateMatrixWorld();
        
        // For simplicity, just return a box for now - proper HD letters need more complex merging
        return new THREE.BoxGeometry(size * 2, size, depth);
    }
    
    createGeometry(geometryData) {
        switch (geometryData.type) {
            case 'box':
                return new THREE.BoxGeometry(
                    geometryData.width || 1,
                    geometryData.height || 1,
                    geometryData.depth || 1
                );
            case 'sphere':
                return new THREE.SphereGeometry(
                    geometryData.radius || 0.5,
                    geometryData.widthSegments || 32,
                    geometryData.heightSegments || 16
                );
            case 'plane':
                return new THREE.PlaneGeometry(
                    geometryData.width || 1,
                    geometryData.height || 1
                );
            case 'cylinder':
                return new THREE.CylinderGeometry(
                    geometryData.radiusTop || 0.5,
                    geometryData.radiusBottom || 0.5,
                    geometryData.height || 1,
                    geometryData.radialSegments || 8
                );
            case 'text':
                const text = geometryData.text || 'TEXT';
                const size = geometryData.size || 1;
                const depth = geometryData.depth || 0.1;
                
                console.log('[HD1-ThreeJS] Creating text geometry:', text);
                
                // Since font loading is problematic, create geometric text using letter shapes
                if (text.toUpperCase() === 'HD') {
                    // Create composite geometry for HD letters
                    const hdGeometry = this.createHDLetterGeometry(size, depth);
                    console.log('[HD1-ThreeJS] Created geometric HD letters');
                    return hdGeometry;
                } else {
                    // For other text, create placeholder
                    console.log('[HD1-ThreeJS] Creating placeholder for:', text);
                    return new THREE.BoxGeometry(
                        size * text.length * 0.6,  // Approximate text width
                        size,                       // Text height
                        depth                       // Text depth
                    );
                }
            default:
                return new THREE.BoxGeometry(1, 1, 1);
        }
    }
    
    createMaterial(materialData) {
        switch (materialData.type) {
            case 'basic':
                return new THREE.MeshBasicMaterial({
                    color: materialData.color || 0x777777,
                    transparent: materialData.transparent || false,
                    opacity: materialData.opacity || 1.0
                });
            case 'phong':
                return new THREE.MeshPhongMaterial({
                    color: materialData.color || 0x777777,
                    specular: materialData.specular || 0x111111,
                    shininess: materialData.shininess || 30,
                    transparent: materialData.transparent || false,
                    opacity: materialData.opacity || 1.0
                });
            case 'standard':
                return new THREE.MeshStandardMaterial({
                    color: materialData.color || 0x777777,
                    metalness: materialData.metalness || 0.0,
                    roughness: materialData.roughness || 0.5,
                    transparent: materialData.transparent || false,
                    opacity: materialData.opacity || 1.0
                });
            default:
                return new THREE.MeshPhongMaterial({ color: materialData.color || 0x777777 });
        }
    }
    
    // Public API methods
    getScene() {
        return this.scene;
    }
    
    getCamera() {
        return this.camera;
    }
    
    getRenderer() {
        return this.renderer;
    }
    
    getObject(id) {
        return this.objects.get(id);
    }
    
    getAvatar(sessionId) {
        return this.avatars.get(sessionId);
    }
    
    getAllObjects() {
        return Array.from(this.objects.values());
    }
    
    getAllAvatars() {
        return Array.from(this.avatars.values());
    }
    
    dispose() {
        // Clean up resources
        this.objects.forEach(obj => {
            if (obj.geometry) obj.geometry.dispose();
            if (obj.material) obj.material.dispose();
        });
        
        this.avatars.forEach(avatar => {
            if (avatar.geometry) avatar.geometry.dispose();
            if (avatar.material) avatar.material.dispose();
        });
        
        this.renderer.dispose();
        console.log('[HD1-ThreeJS] Scene manager disposed');
    }
    
    // Request full sync when Three.js client is ready
    async requestFullSyncWhenReady() {
        // Wait for API client and clientId to be available
        let attempts = 0;
        const maxAttempts = 50; // 5 seconds max wait
        
        const waitForReady = () => {
            attempts++;
            
            if (window.apiClient && window.clientId) {
                console.log('[HD1-ThreeJS] Ready to request full sync');
                this.requestFullSync();
                return;
            }
            
            if (attempts >= maxAttempts) {
                console.warn('[HD1-ThreeJS] Timeout waiting for API client and client ID');
                return;
            }
            
            // Wait 100ms and try again
            setTimeout(waitForReady, 100);
        };
        
        setTimeout(waitForReady, 100); // Start checking after 100ms
    }
    
    // Request full sync from server
    async requestFullSync() {
        if (!window.apiClient) {
            console.warn('[HD1-ThreeJS] API client not available for full sync');
            return;
        }
        
        try {
            console.log('[HD1-ThreeJS] Requesting full sync...');
            const response = await window.apiClient.getFullSync();
            
            if (response.success && response.operations) {
                console.log(`[HD1-ThreeJS] Received ${response.operations.length} operations for bootstrap`);
                
                // Apply all operations to scene
                for (const opWrapper of response.operations) {
                    this.handleSyncOperation(opWrapper.operation);
                }
                
                console.log(`[HD1-ThreeJS] Applied ${response.operations.length} operations to scene`);
            } else {
                console.error('[HD1-ThreeJS] Failed to get full sync:', response);
            }
        } catch (error) {
            console.error('[HD1-ThreeJS] Full sync request failed:', error);
        }
    }
    
    // Handle sync operations - SINGLE SOURCE OF TRUTH
    handleSyncOperation(operation) {
        console.log('[HD1-ThreeJS] Handling sync operation:', operation.type, operation);
        
        switch (operation.type) {
            case 'entity_create':
                this.handleEntityCreate(operation.data);
                break;
            case 'entity_update':
                this.handleEntityUpdate(operation.data);
                break;
            case 'entity_delete':
                this.handleEntityDelete(operation.data);
                break;
            case 'avatar_create':
                this.handleAvatarCreate(operation.data);
                break;
            case 'avatar_move':
                this.handleAvatarMove(operation.data);
                break;
            case 'avatar_remove':
                this.handleAvatarRemove(operation.data);
                break;
            case 'scene_update':
                this.handleSceneUpdate(operation.data);
                break;
            default:
                console.warn('[HD1-ThreeJS] Unknown operation type:', operation.type);
        }
    }
    
    handleEntityCreate(data) {
        const geometry = this.createGeometry(data.geometry);
        const material = this.createMaterial(data.material);
        const mesh = new THREE.Mesh(geometry, material);
        
        // Set position
        if (data.position) {
            mesh.position.set(data.position.x, data.position.y, data.position.z);
        }
        
        // Set rotation
        if (data.rotation) {
            mesh.rotation.set(data.rotation.x, data.rotation.y, data.rotation.z);
        }
        
        // Set scale
        if (data.scale) {
            mesh.scale.set(data.scale.x, data.scale.y, data.scale.z);
        }
        
        // Set visibility
        if (data.visible !== undefined) {
            mesh.visible = data.visible;
        }
        
        // Add to scene and track
        this.scene.add(mesh);
        this.objects.set(data.id, mesh);
        
        console.log('[HD1-ThreeJS] Entity created:', data.id);
    }
    
    handleEntityUpdate(data) {
        const mesh = this.objects.get(data.id);
        if (!mesh) return;
        
        // Update position
        if (data.position) {
            mesh.position.set(data.position.x, data.position.y, data.position.z);
        }
        
        // Update rotation
        if (data.rotation) {
            mesh.rotation.set(data.rotation.x, data.rotation.y, data.rotation.z);
        }
        
        // Update scale
        if (data.scale) {
            mesh.scale.set(data.scale.x, data.scale.y, data.scale.z);
        }
        
        // Update visibility
        if (data.visible !== undefined) {
            mesh.visible = data.visible;
        }
        
        console.log('[HD1-ThreeJS] Entity updated:', data.id);
    }
    
    handleEntityDelete(data) {
        const mesh = this.objects.get(data.id);
        if (mesh) {
            this.scene.remove(mesh);
            this.objects.delete(data.id);
            console.log('[HD1-ThreeJS] Entity deleted:', data.id);
        }
    }
    
    handleAvatarCreate(data) {
        const avatar = this.createAvatar(data.hd1_id, data);
        this.avatars.set(data.hd1_id, avatar);
        console.log('[HD1-ThreeJS] Avatar created:', data.hd1_id);
    }
    
    handleAvatarMove(data) {
        this.updateAvatar(data.hd1_id, data);
    }
    
    handleAvatarRemove(data) {
        this.removeAvatar(data.hd1_id);
    }
    
    handleSceneUpdate(data) {
        // Update scene properties
        if (data.background) {
            this.scene.background = new THREE.Color(data.background);
        }
        
        if (data.fog) {
            this.scene.fog = new THREE.Fog(
                data.fog.color,
                data.fog.near,
                data.fog.far
            );
        }
        
        console.log('[HD1-ThreeJS] Scene updated');
    }
}

// Initialize HD1ThreeJS when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    // Create global instance for compatibility
    window.HD1ThreeJS = HD1ThreeJS;
    
    // Auto-initialize with default canvas
    if (document.getElementById('holodeck-canvas')) {
        window.hd1ThreeJS = new HD1ThreeJS('holodeck-canvas');
        console.log('[HD1] Three.js scene manager initialized with ES modules');
    }
});

// Export for ES modules
export default HD1ThreeJS;