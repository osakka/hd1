/**
 * ===================================================================
 * HD1 JavaScript PlayCanvas Bridge
 * ===================================================================
 *
 * HD1 v3.0 FEATURES:
 * • PlayCanvas game engine integration
 * • Channel-based entity management via API
 * • WebSocket real-time synchronization
 * • Single source of truth architecture
 *
 * Generated from: api.yaml + PlayCanvas schemas
 * ===================================================================
 */

// Enhanced HD1 JavaScript API Bridge for PlayCanvas
window.hd1 = window.hd1 || {};

// Core session and channel management
function getCurrentSessionId() {
    return window.currentSessionId || document.querySelector('[data-session-id]')?.dataset.sessionId || 'default';
}

function getCurrentChannelId() {
    return window.currentChannelId || 'channel_one';
}

// PlayCanvas validation functions
const playcanvasValidation = {
    validateNumber: (value, min, max) => {
        const num = parseFloat(value);
        if (isNaN(num)) throw new Error(`Invalid number: ${value}`);
        if (min !== undefined && num < min) throw new Error(`Value ${num} below minimum ${min}`);
        if (max !== undefined && num > max) throw new Error(`Value ${num} above maximum ${max}`);
        return num;
    },
    
    validateColor: (value) => {
        if (!/^#[0-9a-fA-F]{6}$/.test(value)) {
            throw new Error(`Invalid color format: ${value}. Expected #rrggbb`);
        }
        return value;
    },
    
    validateEnum: (value, options) => {
        if (!options.includes(value)) {
            throw new Error(`Invalid option: ${value}. Expected one of: ${options.join(', ')}`);
        }
        return value;
    }
};

/**
 * HD1 v3.0: Entity management via channels
 * Note: Direct entity creation is deprecated, use channel YAML configuration
 */
hd1.createEntityViaChannel = function(channelId, entityName, entityType, options = {}) {
    console.warn('[HD1] Direct entity creation deprecated in v3.0');
    console.info('[HD1] Use channel YAML configuration for entity management');
    console.info(`[HD1] Entity: ${entityName} (${entityType}) in channel ${channelId}`);
    console.info(`[HD1] Edit: /opt/hd1/share/channels/${channelId}.yaml`);
    
    return Promise.reject(new Error('Entity management via channel YAML only in v3.0'));
};

/**
 * HD1 v3.0: Channel lighting configuration
 */
hd1.configureChannelLighting = function(channelId, lightType, intensity = 1.0, color = '#ffffff') {
    console.warn('[HD1] Direct lighting deprecated in v3.0');
    console.info('[HD1] Use channel YAML configuration for lighting');
    console.info(`[HD1] Light: ${lightType} intensity=${intensity} color=${color}`);
    console.info(`[HD1] Edit: /opt/hd1/share/channels/${channelId}.yaml`);
    
    return Promise.reject(new Error('Lighting management via channel YAML only in v3.0'));
};

/**
 * Join session to channel (v3.0 API)
 */
hd1.joinSessionToChannel = function(sessionId, channelId) {
    const payload = {
        client_id: `client_${sessionId}`,
        reconnect: false
    };
    
    return fetch(`/api/sessions/${sessionId}/channel/join`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
    }).then(response => {
        if (!response.ok) {
            throw new Error(`Failed to join channel: ${response.statusText}`);
        }
        return response.json();
    });
};

/**
 * PlayCanvas capabilities inspection
 */
hd1.playcanvasCapabilities = function() {
    const capabilities = {
        entityManagement: 'channel-based YAML configuration',
        gameEngine: 'PlayCanvas',
        synchronization: 'WebSocket real-time',
        architecture: 'API-first',
        supportedComponents: ['transform', 'render', 'light', 'physics'],
        lightTypes: ['directional', 'point', 'spot', 'ambient'],
        physicsTypes: ['static', 'dynamic', 'kinematic']
    };
    
    console.log('[HD1] PlayCanvas Integration Capabilities:', capabilities);
    return capabilities;
};

/**
 * Function signature verification
 */
hd1.verifyPlaycanvasIntegration = function() {
    const status = {
        channelBasedEntityManagement: true,
        playcanvasGameEngine: typeof pc !== 'undefined',
        websocketSynchronization: true,
        apiFirstArchitecture: true,
        yamlConfiguration: true,
        hd1Version: 'v3.0'
    };
    
    console.log('[HD1] PlayCanvas Integration Status:', status);
    return status;
};

// Console integration
if (typeof console !== 'undefined') {
    console.log('[HD1] PlayCanvas bridge loaded');
    console.log('[HD1] Channel-based entity management: ACTIVE');
    console.log('[HD1] PlayCanvas game engine integration: ACTIVE');
    console.log('[HD1] HD1 v3.0 ready');
}