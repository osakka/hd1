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
        
        // Camera controls
        this.controls = null;
        this.cameraTarget = new THREE.Vector3(0, 0, 0);
        
        // Initialize scene
        this.setupRenderer();
        this.setupLighting();
        this.setupCamera();
        this.setupEventListeners();
        
        // Start render loop
        this.animate();
        
        console.log('[HD1-ThreeJS] Scene manager initialized');
        console.log('[HD1-ThreeJS] Three.js version:', THREE.REVISION);
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
        this.camera.position.set(5, 5, 5);
        this.camera.lookAt(0, 0, 0);
        
        // Simple orbit controls
        this.enableOrbitControls();
    }
    
    enableOrbitControls() {
        // Basic orbit controls without external dependencies
        let isMouseDown = false;
        let mouseX = 0;
        let mouseY = 0;
        let cameraDistance = 10;
        let cameraAngleX = 0;
        let cameraAngleY = 0;
        
        this.canvas.addEventListener('mousedown', (e) => {
            isMouseDown = true;
            mouseX = e.clientX;
            mouseY = e.clientY;
        });
        
        this.canvas.addEventListener('mouseup', () => {
            isMouseDown = false;
        });
        
        this.canvas.addEventListener('mousemove', (e) => {
            if (!isMouseDown) return;
            
            const deltaX = e.clientX - mouseX;
            const deltaY = e.clientY - mouseY;
            
            cameraAngleX -= deltaY * 0.01;
            cameraAngleY -= deltaX * 0.01;
            
            // Limit vertical rotation
            cameraAngleX = Math.max(-Math.PI / 2, Math.min(Math.PI / 2, cameraAngleX));
            
            this.updateCameraPosition(cameraDistance, cameraAngleX, cameraAngleY);
            
            mouseX = e.clientX;
            mouseY = e.clientY;
        });
        
        this.canvas.addEventListener('wheel', (e) => {
            cameraDistance += e.deltaY * 0.01;
            cameraDistance = Math.max(2, Math.min(50, cameraDistance));
            this.updateCameraPosition(cameraDistance, cameraAngleX, cameraAngleY);
        });
    }
    
    updateCameraPosition(distance, angleX, angleY) {
        const x = distance * Math.cos(angleX) * Math.cos(angleY);
        const y = distance * Math.sin(angleX);
        const z = distance * Math.cos(angleX) * Math.sin(angleY);
        
        this.camera.position.set(x, y, z);
        this.camera.lookAt(this.cameraTarget);
    }
    
    setupEventListeners() {
        window.addEventListener('resize', () => {
            this.camera.aspect = window.innerWidth / window.innerHeight;
            this.camera.updateProjectionMatrix();
            this.renderer.setSize(window.innerWidth, window.innerHeight);
        });
    }
    
    animate() {
        requestAnimationFrame(() => this.animate());
        this.renderer.render(this.scene, this.camera);
    }
    
    // Operation handlers - called by sync system
    applyOperation(operation) {
        switch (operation.type) {
            case 'avatar_move':
                this.updateAvatar(operation.data.session_id, operation.data);
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