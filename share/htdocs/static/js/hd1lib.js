/**
 * ===================================================================
 * WARNING: AUTO-GENERATED CODE - DO NOT MODIFY THIS FILE
 * ===================================================================
 *
 * This file is automatically generated from api_new.yaml specification.
 * 
 * ⚠️  CRITICAL WARNING: ALL MANUAL CHANGES WILL BE LOST ⚠️
 *
 * • This file is regenerated on every build
 * • Changes made here are NON-PERSISTENT
 * • Manual modifications will be OVERWRITTEN
 * • To modify API client: Update api_new.yaml specification
 *
 * Generation Command: make generate
 * Source File: /opt/hd1/src/api.yaml
 * Generated: Auto-generated by HD1 Three.js specification-driven development
 *
 * ===================================================================
 * SINGLE SOURCE OF TRUTH: api.yaml drives ALL client methods
 * ===================================================================
 */

class HD1ThreeJSAPIClient {
    constructor(baseURL = '/api', hd1Id = null) {
        this.baseURL = baseURL;
        this.hd1Id = hd1Id; // Server-provided hd1_id only
    }

    // Set hd1_id from server (called when WebSocket receives client_init)
    setHd1Id(hd1Id) {
        this.hd1Id = hd1Id;
    }

    async request(method, path, data = null) {
        const url = this.baseURL + path;
        const headers = {
            'Content-Type': 'application/json',
            'X-Client-ID': this.hd1Id
        };

        const options = {
            method: method,
            headers: headers
        };

        if (data && (method === 'POST' || method === 'PUT')) {
            options.body = JSON.stringify(data);
        }

        const response = await fetch(url, options);

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }

        return response.json();
    }

    extractPathParams(pathTemplate, params) {
        let path = pathTemplate;
        let paramIndex = 0;
        
        // Replace path parameters like {sessionId} with actual values
        path = path.replace(/\{[^}]+\}/g, () => {
            if (paramIndex < params.length) {
                return params[paramIndex++];
            }
            throw new Error(`Missing parameter for path template: ${pathTemplate}`);
        });
        
        return path;
    }

    // ========================================
    // SYNC OPERATIONS (Generated from spec)
    // ========================================


    /**
     * GET /sync/missing/{from}/{to} - getMissingOperations
     */
    async getMissingOperations(param1, param2) {
        const path = this.extractPathParams('/sync/missing/{from}/{to}', [param1, param2]);
        return this.request('GET', path);
    }

    /**
     * GET /sync/stats - getSyncStats
     */
    async getSyncStats() {
        return this.request('GET', '/sync/stats');
    }

    /**
     * POST /sync/operations - submitOperation
     */
    async submitOperation(data = null) {
        return this.request('POST', '/sync/operations', data);
    }

    /**
     * GET /sync/full - getFullSync
     */
    async getFullSync() {
        return this.request('GET', '/sync/full');
    }


    // ========================================
    // ENTITIES (Generated from spec)
    // ========================================


    /**
     * GET /entities - getEntities
     */
    async getEntities() {
        return this.request('GET', '/entities');
    }

    /**
     * PUT /entities/{entityId} - updateEntity
     */
    async updateEntity(param1, data = null) {
        const path = this.extractPathParams('/entities/{entityId}', [param1]);
        return this.request('PUT', path, data);
    }

    /**
     * DELETE /entities/{entityId} - deleteEntity
     */
    async deleteEntity(param1) {
        const path = this.extractPathParams('/entities/{entityId}', [param1]);
        return this.request('DELETE', path);
    }


    // ========================================
    // AVATARS (Generated from spec)
    // ========================================


    /**
     * POST /avatars/{sessionId}/move - moveAvatar
     */
    async moveAvatar(param1, data = null) {
        const path = this.extractPathParams('/avatars/{sessionId}/move', [param1]);
        return this.request('POST', path, data);
    }

    /**
     * PUT /avatars/{avatarId} - updateAvatar
     */
    async updateAvatar(param1, data = null) {
        const path = this.extractPathParams('/avatars/{avatarId}', [param1]);
        return this.request('PUT', path, data);
    }

    /**
     * DELETE /avatars/{avatarId} - removeAvatar
     */
    async removeAvatar(param1) {
        const path = this.extractPathParams('/avatars/{avatarId}', [param1]);
        return this.request('DELETE', path);
    }

    /**
     * GET /avatars - getAvatars
     */
    async getAvatars() {
        return this.request('GET', '/avatars');
    }

    /**
     * POST /avatars - createAvatar
     */
    async createAvatar(data = null) {
        return this.request('POST', '/avatars', data);
    }


    // ========================================
    // SCENE MANAGEMENT (Generated from spec)
    // ========================================


    /**
     * GET /scene - getScene
     */
    async getScene() {
        return this.request('GET', '/scene');
    }

    /**
     * PUT /scene - updateScene
     */
    async updateScene(data = null) {
        return this.request('PUT', '/scene', data);
    }


    // ========================================
    // MATERIALS (Generated from spec)
    // ========================================


    /**
     * POST /materials/basic - createBasicMaterial
     */
    async createBasicMaterial(data = null) {
        return this.request('POST', '/materials/basic', data);
    }

    /**
     * POST /materials/phong - createPhongMaterial
     */
    async createPhongMaterial(data = null) {
        return this.request('POST', '/materials/phong', data);
    }

    /**
     * POST /materials/standard - createStandardMaterial
     */
    async createStandardMaterial(data = null) {
        return this.request('POST', '/materials/standard', data);
    }

    /**
     * POST /materials/physical - createPhysicalMaterial
     */
    async createPhysicalMaterial(data = null) {
        return this.request('POST', '/materials/physical', data);
    }


    // ========================================
    // SYSTEM (Generated from spec)
    // ========================================


    /**
     * GET /system/version - getVersion
     */
    async getVersion() {
        return this.request('GET', '/system/version');
    }


    // ========================================
    // CONVENIENCE METHODS
    // ========================================

    /**
     * Create a box entity with default material
     */
    async createBox(width = 1, height = 1, depth = 1, color = '#777777', position = {x: 0, y: 0, z: 0}) {
        return this.createEntityWithGeometry({
            geometry: {
                type: 'box',
                width: width,
                height: height,
                depth: depth
            },
            material: {
                type: 'phong',
                color: color
            },
            position: position
        });
    }

    /**
     * Create a sphere entity with default material
     */
    async createSphere(radius = 0.5, color = '#777777', position = {x: 0, y: 0, z: 0}) {
        return this.createEntityWithGeometry({
            geometry: {
                type: 'sphere',
                radius: radius
            },
            material: {
                type: 'phong',
                color: color
            },
            position: position
        });
    }

    /**
     * Move entity to new position
     */
    async moveEntity(entityId, position) {
        return this.updateEntity(entityId, { position: position });
    }

    /**
     * Change entity color
     */
    async changeEntityColor(entityId, color) {
        return this.updateEntity(entityId, {
            material: {
                type: 'phong',
                color: color
            }
        });
    }

    /**
     * Set scene background color
     */
    async setBackground(color) {
        return this.updateScene({ background: color });
    }

    /**
     * Add fog to scene
     */
    async addFog(color = '#ffffff', near = 1, far = 100) {
        return this.updateScene({
            fog: {
                color: color,
                near: near,
                far: far
            }
        });
    }
}

// Export for module systems
if (typeof module !== 'undefined' && module.exports) {
    module.exports = HD1ThreeJSAPIClient;
}

// Global export
if (typeof window !== 'undefined') {
    window.HD1ThreeJSAPIClient = HD1ThreeJSAPIClient;
}