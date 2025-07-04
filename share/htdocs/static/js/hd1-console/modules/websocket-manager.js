/**
 * HD1 WebSocket Manager - Bulletproof Connection Handling
 * 
 * Manages WebSocket connections with infinite resilience and graceful error handling.
 * Includes automatic reconnection, message queuing, and connection state management.
 */
class HD1WebSocketManager {
    constructor(domManager) {
        this.dom = domManager;
        this.ws = null;
        this.ready = false;
        
        // Connection state
        this.isConnected = false;
        this.isConnecting = false;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 99;
        this.reconnectTimeout = null;
        this.reconnectDelay = 1000; // Start with 1 second
        this.maxReconnectDelay = 30000; // Max 30 seconds
        
        // Message handling
        this.messageQueue = [];
        this.lastMessageTime = 0;
        this.heartbeatInterval = null;
        this.heartbeatFrequency = 5000; // 5 seconds for better latency tracking
        
        // Latency tracking
        this.pingTimestamps = new Map();
        this.lastLatency = 0;
        
        // Event handlers
        this.eventHandlers = new Map();
        
        // 🔥 HIGH-FREQUENCY UPDATE PROTECTION
        this.avatarRegistry = new Map(); // session_id -> {entity, lastUpdate, creationStatus}
        this.positionUpdateQueue = new Map(); // session_id -> latest position
        this.creationLocks = new Set(); // session_ids currently being created
        this.positionUpdateThrottle = 16; // ~60fps max updates
        this.lastPositionUpdate = 0;
        this.positionUpdateTimer = null;
        
        // Entity lifecycle protection
        this.entityProtectionMode = true;
        this.maxConcurrentCreations = 2;
        this.activeCreations = 0;
        
        // 🛡️ ROBUSTNESS IMPROVEMENTS
        this.avatarHealthChecker = null; // Periodic avatar health monitoring
        this.avatarRecoveryAttempts = new Map(); // session_id -> attempt_count
        this.maxRecoveryAttempts = 3;
        this.healthCheckInterval = 5000; // Check every 5 seconds
    }

    /**
     * Initialize WebSocket manager
     */
    async initialize() {
        console.log('[HD1-WebSocket] Initializing WebSocket Manager');
        
        // Setup UI connection indicators
        this.setupConnectionUI();
        
        // Start connection
        await this.connect();
        
        this.ready = true;
        console.log('[HD1-WebSocket] WebSocket Manager ready');
    }

    /**
     * Setup connection status UI elements
     */
    setupConnectionUI() {
        // Update connection indicator
        this.updateConnectionStatus('connecting');
        
        // Setup heartbeat monitoring
        this.startHeartbeat();
    }

    /**
     * Establish WebSocket connection
     */
    async connect() {
        if (this.isConnecting || this.isConnected) {
            console.log('[HD1-WebSocket] Connection already in progress or established');
            return;
        }

        this.isConnecting = true;
        this.updateConnectionStatus('connecting');

        try {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = `${protocol}//${window.location.host}/ws`;
            
            console.log(`[HD1-WebSocket] Connecting to ${wsUrl}`);
            
            this.ws = new WebSocket(wsUrl);
            this.setupWebSocketHandlers();
            
        } catch (error) {
            console.error('[HD1-WebSocket] Connection failed:', error);
            this.handleConnectionError(error);
        }
    }

    /**
     * Setup WebSocket event handlers
     */
    setupWebSocketHandlers() {
        this.ws.onopen = (event) => {
            console.log('[HD1-WebSocket] Connection established');
            this.isConnected = true;
            this.isConnecting = false;
            this.reconnectAttempts = 0;
            this.reconnectDelay = 1000;
            
            this.updateConnectionStatus('connected');
            this.processMessageQueue();
            this.startAvatarHealthMonitoring();
            
            this.dom.addDebugEntry('WS_CONNECT', 'WebSocket connected');
        };

        this.ws.onmessage = (event) => {
            this.lastMessageTime = Date.now();
            this.handleMessage(event.data);
        };

        this.ws.onclose = (event) => {
            console.log('[HD1-WebSocket] Connection closed:', event.code, event.reason);
            this.isConnected = false;
            this.isConnecting = false;
            
            this.updateConnectionStatus('disconnected');
            this.stopAvatarHealthMonitoring();
            this.dom.addDebugEntry('WS_DISCONNECT', { code: event.code, reason: event.reason });
            
            // Attempt reconnection
            this.scheduleReconnect();
        };

        this.ws.onerror = (error) => {
            console.error('[HD1-WebSocket] WebSocket error:', error);
            this.handleConnectionError(error);
        };
    }

    /**
     * Handle incoming WebSocket messages
     */
    handleMessage(data) {
        try {
            // Track incoming traffic
            this.trackTraffic(data.length, 0);
            
            const message = JSON.parse(data);
            
            // Route message to appropriate handler
            const messageType = message.type || message.action || 'unknown';
            
            if (this.eventHandlers.has(messageType)) {
                this.eventHandlers.get(messageType)(message);
            } else {
                // Default message handling
                this.handleDefaultMessage(message);
            }
            
        } catch (error) {
            console.error('[HD1-WebSocket] Failed to parse message:', error, data);
        }
    }

    /**
     * Handle default messages (session updates, etc.)
     */
    handleDefaultMessage(message) {
        // Handle common message types
        switch (message.type) {
            case 'session_update':
                this.handleSessionUpdate(message);
                break;
            case 'object_update':
                this.handleObjectUpdate(message);
                break;
            case 'entity_created':
                this.handleEntityCreated(message);
                break;
            case 'entity_deleted':
                this.handleEntityDeleted(message);
                break;
            case 'entity_updated':
                this.handleEntityUpdated(message);
                break;
            case 'component_added':
                this.handleComponentAdded(message);
                break;
            case 'component_removed':
                this.handleComponentRemoved(message);
                break;
            case 'client_count':
                this.handleClientCount(message);
                break;
            case 'pong':
                this.handlePongMessage(message);
                break;
            case 'force_refresh':
                this.handleForceRefresh(message);
                break;
            case 'avatar_position_update':
                this.handleAvatarPositionUpdate(message);
                break;
            case 'avatar_asset_response':
                this.handleAvatarAssetResponse(message);
                break;
            case 'avatar_asset_error':
                this.handleAvatarAssetError(message);
                break;
            default:
                this.dom.addDebugEntry('WS_MESSAGE', message);
        }
    }

    /**
     * Handle session update messages
     */
    handleSessionUpdate(message) {
        console.log('[HD1-WebSocket] Session update received:', message);
        
        // Route to session manager
        if (window.hd1ConsoleManager) {
            const sessionManager = window.hd1ConsoleManager.getModule('session');
            if (sessionManager && sessionManager.handleSessionUpdate) {
                sessionManager.handleSessionUpdate(message);
            }
        }
        
        this.dom.addDebugEntry('SESSION_UPDATE', message);
    }

    /**
     * Handle object update messages
     */
    handleObjectUpdate(message) {
        console.log('[HD1-WebSocket] Object update received:', message);
        
        // Trigger PlayCanvas updates if engine is available
        if (window.hd1GameEngine && message.object) {
            // Handle object updates in PlayCanvas scene
            this.updatePlayCanvasObject(message.object);
        }
        
        this.dom.addDebugEntry('OBJECT_UPDATE', message);
    }

    /**
     * Handle client count updates
     */
    handleClientCount(message) {
        console.log('[HD1-WebSocket] Client count update:', message.count);
        
        // Update any client count displays
        this.dom.addDebugEntry('CLIENT_COUNT', { count: message.count });
    }

    /**
     * Handle pong messages for latency calculation
     */
    handlePongMessage(message) {
        const now = Date.now();
        const pingId = message.ping_id || message.timestamp;
        
        if (pingId && this.pingTimestamps.has(pingId)) {
            const pingTime = this.pingTimestamps.get(pingId);
            const latency = now - pingTime;
            
            this.lastLatency = latency;
            this.pingTimestamps.delete(pingId);
            
            // Send latency to stats manager
            if (window.hd1ConsoleManager) {
                const statsManager = window.hd1ConsoleManager.getModule('stats');
                if (statsManager && statsManager.trackLatency) {
                    statsManager.trackLatency(latency);
                }
            }
            
            console.log(`[HD1-WebSocket] Latency: ${latency}ms`);
        }
    }

    /**
     * Handle force refresh messages - triggers browser reload
     */
    handleForceRefresh(message) {
        console.log('[HD1-WebSocket] Force refresh received:', message);
        
        // Clear storage if requested
        if (message.clear_storage) {
            localStorage.clear();
            sessionStorage.clear();
            console.log('[HD1-WebSocket] Browser storage cleared');
        }
        
        // Log the refresh event
        this.dom.addDebugEntry('FORCE_REFRESH', { 
            session_id: message.session_id,
            clear_storage: message.clear_storage,
            timestamp: message.timestamp 
        });
        
        // Short delay to allow logging, then reload
        setTimeout(() => {
            console.log('[HD1-WebSocket] Reloading browser...');
            window.location.reload();
        }, 100);
    }

    /**
     * Handle avatar position updates - HIGH-FREQUENCY PROTECTED real-time movement
     */
    handleAvatarPositionUpdate(message) {
        // Extract position data
        const { session_id, avatar_name, position, camera_position, channel_id } = message.data || message;
        
        if (!session_id || !avatar_name || !position) {
            console.warn('[HD1-WebSocket] Invalid avatar position update - missing required fields');
            return;
        }
        
        // 🔍 TRACE: Avatar disappearance debugging
        const currentSessionId = this.getCurrentSessionId();
        const isOwnAvatar = (session_id === currentSessionId);
        
        // Log every 10th update for debugging disappearance
        if (Math.random() < 0.1) {
            console.log(`[HD1-WebSocket] 🔍 TRACE: Position update for ${isOwnAvatar ? 'OWN' : 'OTHER'} avatar ${avatar_name} (${session_id})`);
        }
        
        // 🔥 HIGH-FREQUENCY PROTECTION: Queue position updates and throttle processing
        this.queuePositionUpdate(session_id, {
            avatar_name,
            position,
            camera_position,
            channel_id,
            timestamp: Date.now()
        });
        
        // Start throttled processing if not already running
        if (!this.positionUpdateTimer) {
            this.positionUpdateTimer = setTimeout(() => {
                this.processQueuedPositionUpdates();
                this.positionUpdateTimer = null;
            }, this.positionUpdateThrottle);
        }
    }
    
    /**
     * Queue position update for throttled processing
     */
    queuePositionUpdate(session_id, updateData) {
        // Always keep the latest position update per session
        this.positionUpdateQueue.set(session_id, updateData);
    }
    
    /**
     * Process all queued position updates in a batch
     */
    processQueuedPositionUpdates() {
        if (!window.hd1GameEngine) {
            console.warn('[HD1-WebSocket] PlayCanvas not available for position updates');
            return;
        }
        
        const currentSessionId = this.getCurrentSessionId();
        const now = Date.now();
        
        // 🔍 TRACE: Enhanced logging for disappearance debugging
        const queueSize = this.positionUpdateQueue.size;
        const sceneEntityCount = window.hd1GameEngine.root.children.length;
        
        // Process each session's latest position update
        for (const [session_id, updateData] of this.positionUpdateQueue.entries()) {
            try {
                const isOwnAvatar = (session_id === currentSessionId);
                
                // 🛡️ PROTECTED avatar entity management
                let avatarEntity = this.getOrCreateAvatar(session_id, updateData.avatar_name, updateData.position, isOwnAvatar);
                
                if (avatarEntity) {
                    // 🔍 TRACE: Check if entity is about to be destroyed
                    const entityExists = window.hd1GameEngine.root.children.includes(avatarEntity);
                    
                    if (!entityExists) {
                        console.error(`[HD1-WebSocket] 🚨 CRITICAL: Avatar entity ${updateData.avatar_name} (${session_id}) missing from scene during update!`);
                        // Remove from registry and continue
                        this.avatarRegistry.delete(session_id);
                        continue;
                    }
                    
                    // ⚡ OPTIMIZED position update - direct and fast
                    avatarEntity.setPosition(updateData.position.x, updateData.position.y, updateData.position.z);
                    
                    // Update registry
                    this.updateAvatarRegistry(session_id, avatarEntity, now);
                    
                    // Enhanced logging for own avatar
                    if (isOwnAvatar && Math.random() < 0.2) {
                        console.log(`[HD1-WebSocket] 🔍 TRACE: Updated OWN avatar ${updateData.avatar_name} at (${updateData.position.x}, ${updateData.position.y}, ${updateData.position.z})`);
                    }
                    
                    // Minimal logging for high-frequency updates
                    if (Math.random() < 0.05) { // Log 5% of updates to trace disappearance
                        console.log(`[HD1-WebSocket] 📍 Batch updated ${isOwnAvatar ? 'own' : 'other'} avatar ${updateData.avatar_name} - scene entities: ${sceneEntityCount}`);
                    }
                } else if (!this.creationLocks.has(session_id) && this.activeCreations < this.maxConcurrentCreations) {
                    // 🚀 PROTECTED avatar creation
                    console.log(`[HD1-WebSocket] 🔍 TRACE: Creating missing avatar for ${isOwnAvatar ? 'OWN' : 'OTHER'} session ${session_id}`);
                    this.createAvatarWithProtection(session_id, updateData.avatar_name, updateData.position, isOwnAvatar);
                } else {
                    // 🔍 TRACE: Avatar not found and creation blocked
                    console.warn(`[HD1-WebSocket] 🔍 TRACE: Avatar missing for ${session_id}, creation blocked - locks: ${this.creationLocks.size}, active: ${this.activeCreations}`);
                }
            } catch (error) {
                console.error(`[HD1-WebSocket] Failed to process position update for ${session_id}:`, error);
            }
        }
        
        // Clear processed updates
        this.positionUpdateQueue.clear();
        
        // 📊 Enhanced logging for disappearance monitoring
        if (Math.random() < 0.1) {
            console.log(`[HD1-WebSocket] 🔍 BATCH STATS: Queue: ${queueSize}, Registry: ${this.avatarRegistry.size}, Scene entities: ${sceneEntityCount}, Active creations: ${this.activeCreations}`);
            
            this.dom.addDebugEntry('AVATAR_BATCH', { 
                processed_sessions: queueSize,
                total_avatars: this.avatarRegistry.size,
                scene_entities: sceneEntityCount,
                active_creations: this.activeCreations,
                timestamp: Date.now()
            });
        }
    }

    /**
     * Handle entity created messages - PHASE 2: Pure graph extension
     */
    handleEntityCreated(message) {
        console.log('[HD1-WebSocket] Entity created received:', message);
        
        // Pure graph extension - automatically render new entity in PlayCanvas
        if (window.hd1GameEngine && message.data) {
            this.createPlayCanvasEntityFromBroadcast(message.data);
        }
        
        this.dom.addDebugEntry('ENTITY_CREATED', message.data);
    }

    /**
     * Handle entity deleted messages - PHASE 2: Pure graph extension  
     */
    handleEntityDeleted(message) {
        console.log('[HD1-WebSocket] Entity deleted received:', message);
        
        // 🔍 TRACE: Check if this is an avatar being deleted
        if (message.data && message.data.entity_id) {
            // 🛡️  AVATAR PROTECTION: Check if this is an avatar entity
            const isAvatarEntity = message.data.entity_id.includes('session_client_') || message.data.entity_id.includes('session_');
            
            if (isAvatarEntity) {
                console.error(`[HD1-WebSocket] 🚨 AVATAR DELETION BLOCKED: Avatar entity ${message.data.entity_id} deletion blocked for protection!`);
                this.dom.addDebugEntry('ENTITY_DELETED', { ...message.data, blocked: true, reason: 'avatar_protection' });
                return;
            }
            
            // Check if this entity is in our avatar registry
            for (const [session_id, registryEntry] of this.avatarRegistry.entries()) {
                if (registryEntry.entity && (
                    registryEntry.entity.name === message.data.entity_id ||
                    registryEntry.entity.hd1Id === message.data.entity_id
                )) {
                    console.error(`[HD1-WebSocket] 🚨 AVATAR DELETION: Avatar entity ${message.data.entity_id} for session ${session_id} is being deleted via WebSocket!`);
                    // Remove from registry
                    this.avatarRegistry.delete(session_id);
                    break;
                }
            }
        }
        
        // Pure graph extension - automatically remove entity from PlayCanvas
        if (window.hd1GameEngine && message.data && message.data.entity_id) {
            this.destroyPlayCanvasEntityFromBroadcast(message.data.entity_id);
        }
        
        this.dom.addDebugEntry('ENTITY_DELETED', message.data);
    }

    /**
     * Handle entity updated messages - PHASE 2: Pure graph extension
     */
    handleEntityUpdated(message) {
        console.log('[HD1-WebSocket] Entity updated received:', message);
        
        // Pure graph extension - automatically update entity in PlayCanvas
        if (window.hd1GameEngine && message.data) {
            this.updatePlayCanvasEntityFromBroadcast(message.data);
        }
        
        this.dom.addDebugEntry('ENTITY_UPDATED', message.data);
    }

    /**
     * Handle component added messages - PHASE 2: Pure graph extension
     */
    handleComponentAdded(message) {
        console.log('[HD1-WebSocket] Component added received:', message);
        
        // Pure graph extension - automatically add component to PlayCanvas entity
        if (window.hd1GameEngine && message.data) {
            this.addPlayCanvasComponentFromBroadcast(message.data);
        }
        
        this.dom.addDebugEntry('COMPONENT_ADDED', message.data);
    }

    /**
     * Handle component removed messages - PHASE 2: Pure graph extension
     */
    handleComponentRemoved(message) {
        console.log('[HD1-WebSocket] Component removed received:', message);
        
        // Pure graph extension - automatically remove component from PlayCanvas entity
        if (window.hd1GameEngine && message.data) {
            this.removePlayCanvasComponentFromBroadcast(message.data);
        }
        
        this.dom.addDebugEntry('COMPONENT_REMOVED', message.data);
    }

    /**
     * Create PlayCanvas entity from WebSocket broadcast (pure graph extension)
     */
    createPlayCanvasEntityFromBroadcast(entityData) {
        try {
            if (!window.createObjectFromData) {
                console.warn('[HD1-WebSocket] createObjectFromData not available, deferring entity creation');
                return;
            }
            
            // Use existing PlayCanvas entity creation function
            window.createObjectFromData(entityData);
            console.log('[HD1-WebSocket] Entity rendered in PlayCanvas:', entityData.name);
        } catch (error) {
            console.error('[HD1-WebSocket] Failed to create PlayCanvas entity from broadcast:', error);
        }
    }

    /**
     * Destroy PlayCanvas entity from WebSocket broadcast (pure graph extension)
     */
    destroyPlayCanvasEntityFromBroadcast(entityId) {
        try {
            if (!window.deleteObjectByName) {
                console.warn('[HD1-WebSocket] deleteObjectByName not available, cannot destroy entity');
                return;
            }
            
            // 🔥 CRITICAL FIX: Don't allow avatar entity deletion via WebSocket!
            // Check if this entity is an avatar (has session-avatar tag or session_ name)
            const isAvatarEntity = entityId.includes('session_client_') || entityId.includes('session_');
            
            if (isAvatarEntity) {
                console.warn(`[HD1-WebSocket] 🛡️  AVATAR PROTECTION: Blocked deletion of avatar entity ${entityId} via WebSocket`);
                return;
            }
            
            // Use existing PlayCanvas entity deletion function for non-avatar entities
            window.deleteObjectByName(entityId);
            console.log('[HD1-WebSocket] Entity removed from PlayCanvas:', entityId);
        } catch (error) {
            console.error('[HD1-WebSocket] Failed to destroy PlayCanvas entity from broadcast:', error);
        }
    }

    /**
     * Update PlayCanvas entity from WebSocket broadcast (pure graph extension) 
     */
    updatePlayCanvasEntityFromBroadcast(entityData) {
        try {
            // 🔥 CRITICAL FIX: Don't delete/recreate avatars - just update position!
            // Check if this is an avatar entity (has session-avatar tag or session_ name)
            const isAvatar = (entityData.tags && entityData.tags.includes('session-avatar')) ||
                            (entityData.name && entityData.name.includes('session_client_'));
            
            if (isAvatar && window.hd1GameEngine) {
                // For avatar updates, just update position without deleting entity
                const entity = window.hd1GameEngine.root.findByName(entityData.name);
                if (entity && entityData.components && entityData.components.transform && entityData.components.transform.position) {
                    const pos = entityData.components.transform.position;
                    entity.setPosition(pos[0], pos[1], pos[2]);
                    console.log(`[HD1-WebSocket] 🔄 Avatar position updated directly: ${entityData.name} -> (${pos[0]}, ${pos[1]}, ${pos[2]})`);
                    return;
                }
            }
            
            // For non-avatar updates, use the standard delete/recreate approach
            if (entityData.entity_id && window.deleteObjectByName) {
                window.deleteObjectByName(entityData.entity_id);
            }
            
            if (window.createObjectFromData) {
                window.createObjectFromData(entityData);
            }
            
            console.log('[HD1-WebSocket] Entity updated in PlayCanvas:', entityData.name);
        } catch (error) {
            console.error('[HD1-WebSocket] Failed to update PlayCanvas entity from broadcast:', error);
        }
    }

    /**
     * Add component to PlayCanvas entity from WebSocket broadcast
     */
    addPlayCanvasComponentFromBroadcast(componentData) {
        try {
            // Component additions are handled by entity updates
            console.log('[HD1-WebSocket] Component added to entity:', componentData.entity_id);
        } catch (error) {
            console.error('[HD1-WebSocket] Failed to add component from broadcast:', error);
        }
    }

    /**
     * Remove component from PlayCanvas entity from WebSocket broadcast  
     */
    removePlayCanvasComponentFromBroadcast(componentData) {
        try {
            // Component removals are handled by entity updates
            console.log('[HD1-WebSocket] Component removed from entity:', componentData.entity_id);
        } catch (error) {
            console.error('[HD1-WebSocket] Failed to remove component from broadcast:', error);
        }
    }

    /**
     * Update PlayCanvas object from WebSocket message
     */
    updatePlayCanvasObject(objectData) {
        try {
            if (!window.hd1GameEngine) return;
            
            const app = window.hd1GameEngine;
            const entity = app.root.findByName(objectData.name);
            
            if (entity && objectData.transform) {
                // Update transform
                if (objectData.transform.position) {
                    entity.setPosition(
                        objectData.transform.position.x,
                        objectData.transform.position.y,
                        objectData.transform.position.z
                    );
                }
                
                if (objectData.transform.rotation) {
                    entity.setEulerAngles(
                        objectData.transform.rotation.x,
                        objectData.transform.rotation.y,
                        objectData.transform.rotation.z
                    );
                }
                
                if (objectData.transform.scale) {
                    entity.setLocalScale(
                        objectData.transform.scale.x,
                        objectData.transform.scale.y,
                        objectData.transform.scale.z
                    );
                }
            }
        } catch (error) {
            console.error('[HD1-WebSocket] Failed to update PlayCanvas object:', error);
        }
    }

    /**
     * Send message through WebSocket
     */
    send(message) {
        if (!this.isConnected) {
            console.warn('[HD1-WebSocket] Not connected, queuing message');
            this.messageQueue.push(message);
            return false;
        }

        try {
            const messageStr = typeof message === 'string' ? message : JSON.stringify(message);
            this.ws.send(messageStr);
            
            // Track outgoing traffic
            this.trackTraffic(0, messageStr.length);
            
            return true;
        } catch (error) {
            console.error('[HD1-WebSocket] Failed to send message:', error);
            this.messageQueue.push(message);
            return false;
        }
    }

    /**
     * Process queued messages when connection is restored
     */
    processMessageQueue() {
        while (this.messageQueue.length > 0) {
            const message = this.messageQueue.shift();
            this.send(message);
        }
    }

    /**
     * Schedule reconnection attempt
     */
    scheduleReconnect() {
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            console.log('[HD1-WebSocket] Max reconnect attempts reached, triggering rebootstrap');
            this.triggerRebootstrap();
            return;
        }

        this.reconnectAttempts++;
        
        console.log(`[HD1-WebSocket] Scheduling reconnect attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts} in ${this.reconnectDelay}ms`);
        
        this.reconnectTimeout = setTimeout(() => {
            this.connect();
        }, this.reconnectDelay);
        
        // Exponential backoff with jitter
        this.reconnectDelay = Math.min(this.reconnectDelay * 1.5 + Math.random() * 1000, this.maxReconnectDelay);
    }

    /**
     * Trigger system rebootstrap after max reconnect attempts
     */
    triggerRebootstrap() {
        console.log('[HD1-WebSocket] Triggering system rebootstrap');
        this.updateConnectionStatus('rebootstrap');
        
        // Clear all state
        localStorage.removeItem('hd1_session_id');
        localStorage.removeItem('hd1_current_channel');
        
        // Reload page after short delay
        setTimeout(() => {
            window.location.reload();
        }, 2000);
    }

    /**
     * Update connection status in UI
     */
    updateConnectionStatus(status) {
        const statusMessages = {
            connecting: 'Connecting',
            connected: 'Connected',
            disconnected: 'Disconnected',
            error: 'Connection Error',
            rebootstrap: 'Rebootstrapping'
        };

        const statusClasses = {
            connecting: 'connecting',
            connected: 'connected',
            disconnected: 'disconnected',
            error: 'error',
            rebootstrap: 'rebootstrap'
        };

        this.dom.setText('status-connection-text', statusMessages[status] || 'Unknown');
        
        // Update indicator classes
        const indicator = this.dom.get('status-connection-indicator');
        if (indicator) {
            indicator.className = statusClasses[status] || 'unknown';
        }
    }

    /**
     * Handle connection errors
     */
    handleConnectionError(error) {
        console.error('[HD1-WebSocket] Connection error:', error);
        this.updateConnectionStatus('error');
        this.isConnecting = false;
        this.isConnected = false;
        
        this.dom.addDebugEntry('WS_ERROR', { error: error.message });
        
        // Schedule reconnect
        this.scheduleReconnect();
    }

    /**
     * Start heartbeat monitoring with latency tracking
     */
    startHeartbeat() {
        this.heartbeatInterval = setInterval(() => {
            if (this.isConnected) {
                const pingId = Date.now();
                this.pingTimestamps.set(pingId, pingId);
                this.send({ type: 'ping', timestamp: pingId, ping_id: pingId });
                
                // Clean up old ping timestamps (older than 30 seconds)
                const cutoff = pingId - 30000;
                for (const [id, timestamp] of this.pingTimestamps) {
                    if (timestamp < cutoff) {
                        this.pingTimestamps.delete(id);
                    }
                }
            }
        }, this.heartbeatFrequency);
    }

    /**
     * Register event handler for specific message types
     */
    on(messageType, handler) {
        this.eventHandlers.set(messageType, handler);
    }

    /**
     * Unregister event handler
     */
    off(messageType) {
        this.eventHandlers.delete(messageType);
    }

    /**
     * ROBUST avatar entity finder - multiple search strategies for bulletproof sync
     */
    findAvatarEntity(avatar_name, session_id) {
        if (!window.hd1GameEngine) return null;
        
        const entities = window.hd1GameEngine.root.children;
        
        // Strategy 1: Exact name match
        let entity = entities.find(e => e.name === avatar_name);
        if (entity) return entity;
        
        // Strategy 2: HD1 ID match
        entity = entities.find(e => e.hd1Id === avatar_name);
        if (entity) return entity;
        
        // Strategy 3: Session-based name matching
        entity = entities.find(e => 
            e.name && (
                e.name.includes(session_id) || 
                e.name.includes('session') ||
                e.name.includes('avatar')
            )
        );
        if (entity) return entity;
        
        // Strategy 4: Tag-based matching (session-avatar tag)
        entity = entities.find(e => 
            e.hd1Tags && e.hd1Tags.includes && (
                e.hd1Tags.includes('session-avatar') ||
                e.hd1Tags.includes('avatar')
            )
        );
        if (entity) return entity;
        
        // Strategy 5: Any entity with matching session in metadata
        entity = entities.find(e => 
            e.hd1SessionId === session_id ||
            (e.hd1Metadata && e.hd1Metadata.sessionId === session_id)
        );
        
        return entity;
    }
    
    /**
     * Get current session ID from session manager
     */
    getCurrentSessionId() {
        if (window.hd1ConsoleManager) {
            const sessionManager = window.hd1ConsoleManager.getModule('session');
            if (sessionManager && sessionManager.sessionId) {
                return sessionManager.sessionId;
            }
        }
        
        // Fallback to localStorage
        return localStorage.getItem('hd1_session_id');
    }
    
    /**
     * 🛡️ PROTECTED avatar entity retrieval with registry management
     */
    getOrCreateAvatar(session_id, avatar_name, position, isOwnAvatar) {
        // Check registry first
        const registryEntry = this.avatarRegistry.get(session_id);
        
        if (registryEntry && registryEntry.entity) {
            // Verify entity still exists in PlayCanvas scene
            if (window.hd1GameEngine.root.children.includes(registryEntry.entity)) {
                return registryEntry.entity;
            } else {
                // 🔍 TRACE: Entity was destroyed - critical debugging info
                console.error(`[HD1-WebSocket] 🚨 AVATAR DISAPPEARED: Entity for ${isOwnAvatar ? 'OWN' : 'OTHER'} session ${session_id} (${avatar_name}) was destroyed from scene!`);
                console.error(`[HD1-WebSocket] 🚨 Scene entity count: ${window.hd1GameEngine.root.children.length}`);
                console.error(`[HD1-WebSocket] 🚨 Registry size before cleanup: ${this.avatarRegistry.size}`);
                
                // Log remaining entities in scene for debugging
                const remainingEntities = window.hd1GameEngine.root.children.map(e => e.name || 'unnamed');
                console.error(`[HD1-WebSocket] 🚨 Remaining entities: ${remainingEntities.join(', ')}`);
                
                this.avatarRegistry.delete(session_id);
            }
        }
        
        // Try to find existing entity in scene
        const existingEntity = this.findAvatarEntity(avatar_name, session_id);
        if (existingEntity) {
            // 🔍 TRACE: Found existing entity
            if (isOwnAvatar && Math.random() < 0.1) {
                console.log(`[HD1-WebSocket] 🔍 TRACE: Found existing OWN avatar entity ${avatar_name} in scene`);
            }
            
            // Register the found entity
            this.updateAvatarRegistry(session_id, existingEntity, Date.now());
            return existingEntity;
        }
        
        // 🔍 TRACE: No entity found
        if (isOwnAvatar && Math.random() < 0.2) {
            console.warn(`[HD1-WebSocket] 🔍 TRACE: No avatar entity found for OWN session ${session_id} (${avatar_name})`);
        }
        
        return null; // No entity found, creation needed
    }
    
    /**
     * 🚀 PROTECTED avatar creation with concurrency limits
     */
    createAvatarWithProtection(session_id, avatar_name, position, isOwnAvatar) {
        // Check if already being created
        if (this.creationLocks.has(session_id)) {
            return;
        }
        
        // Check concurrency limit
        if (this.activeCreations >= this.maxConcurrentCreations) {
            console.warn(`[HD1-WebSocket] Avatar creation rate limited for ${session_id}`);
            return;
        }
        
        // Lock creation for this session
        this.creationLocks.add(session_id);
        this.activeCreations++;
        
        console.log(`[HD1-WebSocket] 🔒 Creating avatar for ${session_id} (${this.activeCreations}/${this.maxConcurrentCreations})`);
        
        try {
            if (isOwnAvatar) {
                this.createOwnSessionAvatarProtected(session_id, avatar_name, position);
            } else {
                this.createOtherSessionAvatarProtected(session_id, avatar_name, position);
            }
        } catch (error) {
            console.error(`[HD1-WebSocket] Protected avatar creation failed for ${session_id}:`, error);
            this.unlockAvatarCreation(session_id);
        }
    }
    
    /**
     * Update avatar registry with entity and timestamp
     */
    updateAvatarRegistry(session_id, entity, timestamp) {
        this.avatarRegistry.set(session_id, {
            entity: entity,
            lastUpdate: timestamp,
            creationStatus: 'active'
        });
    }
    
    /**
     * Release avatar creation lock
     */
    unlockAvatarCreation(session_id) {
        this.creationLocks.delete(session_id);
        this.activeCreations = Math.max(0, this.activeCreations - 1);
        console.log(`[HD1-WebSocket] 🔓 Released creation lock for ${session_id} (${this.activeCreations}/${this.maxConcurrentCreations})`);
    }
    
    /**
     * AUTO-CREATE missing avatar for bulletproof bidirectional visibility
     */
    createMissingAvatar(session_id, avatar_name, position, isOwnAvatar = false) {
        if (!window.hd1API) {
            console.warn('[HD1-WebSocket] Cannot auto-create avatar: hd1API not available');
            return;
        }
        
        console.log(`[HD1-WebSocket] 🚀 Auto-creating ${isOwnAvatar ? 'OWN' : 'OTHER'} avatar for session ${session_id}`);
        
        if (isOwnAvatar) {
            // For our own avatar, try to fetch from our session's API
            this.createOwnSessionAvatar(session_id, avatar_name, position);
        } else {
            // For other sessions' avatars, create a visual representation directly
            this.createOtherSessionAvatarVisual(session_id, avatar_name, position);
        }
    }
    
    /**
     * Create our own session avatar (uses API)
     */
    createOwnSessionAvatar(session_id, avatar_name, position) {
        // Strategy 1: Try to fetch existing avatar entity from API
        window.hd1API.listEntities(session_id)
            .then(response => {
                if (response.entities) {
                    const avatarEntity = response.entities.find(entity => 
                        entity.tags && entity.tags.includes('session-avatar') &&
                        (entity.name === avatar_name || entity.name.includes(session_id))
                    );
                    
                    if (avatarEntity) {
                        // Create the avatar in PlayCanvas
                        if (window.createObjectFromData) {
                            window.createObjectFromData(avatarEntity);
                            
                            // Position it after creation
                            setTimeout(() => {
                                const newEntity = this.findAvatarEntity(avatar_name, session_id);
                                if (newEntity) {
                                    newEntity.setPosition(position.x, position.y, position.z);
                                    console.log(`[HD1-WebSocket] ✅ Auto-created and positioned OWN avatar ${avatar_name}`);
                                    this.trackAvatarSync(newEntity, 'own_auto_created');
                                }
                            }, 100);
                        }
                    } else {
                        // Strategy 2: Create basic avatar entity via API
                        this.createBasicAvatarEntity(session_id, avatar_name, position);
                    }
                }
            })
            .catch(error => {
                console.warn(`[HD1-WebSocket] Failed to fetch own session entities, creating basic avatar:`, error);
                this.createBasicAvatarEntity(session_id, avatar_name, position);
            });
    }
    
    /**
     * 🛡️ PROTECTED: Create our own session avatar with lifecycle management
     */
    createOwnSessionAvatarProtected(session_id, avatar_name, position) {
        setTimeout(async () => {
            try {
                // Try API-based approach first
                if (window.hd1API) {
                    const response = await window.hd1API.listEntities(session_id);
                    if (response.entities) {
                        const avatarEntity = response.entities.find(entity => 
                            entity.tags && entity.tags.includes('session-avatar')
                        );
                        
                        if (avatarEntity && window.createObjectFromData) {
                            window.createObjectFromData(avatarEntity);
                            
                            // Wait for entity creation and register
                            setTimeout(() => {
                                const newEntity = this.findAvatarEntity(avatar_name, session_id);
                                if (newEntity) {
                                    newEntity.setPosition(position.x, position.y, position.z);
                                    this.updateAvatarRegistry(session_id, newEntity, Date.now());
                                    console.log(`[HD1-WebSocket] ✅ Protected own avatar created: ${avatar_name}`);
                                }
                                this.unlockAvatarCreation(session_id);
                            }, 100);
                            return;
                        }
                    }
                }
                
                // Fallback to basic creation
                this.createBasicAvatarEntityProtected(session_id, avatar_name, position);
                
            } catch (error) {
                console.error(`[HD1-WebSocket] Protected own avatar creation failed:`, error);
                this.unlockAvatarCreation(session_id);
            }
        }, 50); // Small delay to avoid immediate collisions
    }
    
    /**
     * 🛡️ PROTECTED: Create visual representation of another session's avatar
     */
    createOtherSessionAvatarProtected(session_id, avatar_name, position) {
        if (!window.hd1GameEngine) {
            this.unlockAvatarCreation(session_id);
            return;
        }
        
        setTimeout(() => {
            try {
                const app = window.hd1GameEngine;
                
                // Create entity with protection against naming conflicts
                const safeEntityName = `${avatar_name}_${session_id.substring(0, 8)}_${Date.now()}`;
                const entity = new pc.Entity(safeEntityName, app);
                
                // Add transform component
                entity.addComponent('transform', {
                    position: [position.x, position.y, position.z],
                    rotation: [0, 0, 0, 1],
                    scale: [1, 1, 1]
                });
                
                // Add render component with distinct appearance
                entity.addComponent('render', {
                    type: 'capsule',
                    material: {
                        ambient: [0.0, 1.0, 0.5], // Green-teal for other sessions
                        diffuse: [0.0, 1.0, 0.5],
                        emissive: [0.0, 0.2, 0.1]
                    }
                });
                
                // Protected metadata assignment
                entity.hd1Id = avatar_name;
                entity.hd1SessionId = session_id;
                entity.hd1Tags = ['session-avatar', 'other-session', 'visual-only'];
                entity.hd1IsOtherSession = true;
                entity.hd1Protected = true; // Mark as protected entity
                
                // Add to scene
                app.root.addChild(entity);
                
                // Register in protected registry
                this.updateAvatarRegistry(session_id, entity, Date.now());
                
                console.log(`[HD1-WebSocket] ✅ Protected other avatar created: ${avatar_name} for ${session_id}`);
                this.unlockAvatarCreation(session_id);
                
            } catch (error) {
                console.error(`[HD1-WebSocket] Protected other avatar creation failed:`, error);
                this.unlockAvatarCreation(session_id);
            }
        }, Math.random() * 100); // Random delay to prevent collision storms
    }
    
    /**
     * 🛡️ PROTECTED: Create basic avatar entity with lifecycle management
     */
    createBasicAvatarEntityProtected(session_id, avatar_name, position) {
        const avatarData = {
            name: avatar_name,
            tags: ['session-avatar', 'avatar'],
            enabled: true,
            components: {
                transform: {
                    position: [position.x, position.y, position.z],
                    rotation: [0, 0, 0, 1],
                    scale: [1, 1, 1]
                },
                model: {
                    type: 'box',
                    material: {
                        ambient: [0.3, 0.6, 1.0],
                        diffuse: [0.3, 0.6, 1.0]
                    }
                }
            }
        };
        
        // Create via API with protection
        window.hd1API.createEntity(session_id, avatarData)
            .then(response => {
                if (response.success && window.createObjectFromData) {
                    avatarData.entity_id = response.entity_id;
                    window.createObjectFromData(avatarData);
                    
                    // Register after creation
                    setTimeout(() => {
                        const newEntity = this.findAvatarEntity(avatar_name, session_id);
                        if (newEntity) {
                            this.updateAvatarRegistry(session_id, newEntity, Date.now());
                        }
                        this.unlockAvatarCreation(session_id);
                    }, 100);
                } else {
                    this.unlockAvatarCreation(session_id);
                }
            })
            .catch(error => {
                console.error('[HD1-WebSocket] Protected basic avatar creation failed:', error);
                this.unlockAvatarCreation(session_id);
            });
    }

    /**
     * Create visual representation of another session's avatar (client-side only)
     */
    createOtherSessionAvatarVisual(session_id, avatar_name, position) {
        if (!window.hd1GameEngine) {
            console.warn('[HD1-WebSocket] Cannot create other session avatar: PlayCanvas not available');
            return;
        }
        
        try {
            console.log(`[HD1-WebSocket] 🎨 Creating visual representation for OTHER session ${session_id} avatar: ${avatar_name}`);
            
            const app = window.hd1GameEngine;
            
            // Create a distinct visual representation for other session's avatar
            const entity = new pc.Entity(avatar_name, app);
            
            // Add transform component
            entity.addComponent('transform', {
                position: [position.x, position.y, position.z],
                rotation: [0, 0, 0, 1],
                scale: [1, 1, 1]
            });
            
            // Add render component with distinct appearance
            entity.addComponent('render', {
                type: 'capsule',
                material: {
                    ambient: [0.0, 1.0, 0.5], // Green-teal for other sessions
                    diffuse: [0.0, 1.0, 0.5],
                    emissive: [0.0, 0.2, 0.1]
                }
            });
            
            // Add text label to identify the session
            if (pc.Entity.prototype.addComponent) {
                entity.addComponent('element', {
                    type: 'text',
                    text: `Session: ${session_id.substring(0, 8)}...`,
                    fontSize: 0.3,
                    color: [1, 1, 1],
                    anchor: [0.5, 0.5, 0.5, 0.5],
                    pivot: [0.5, 0.5]
                });
            }
            
            // Mark as other session entity
            entity.hd1Id = avatar_name;
            entity.hd1SessionId = session_id;
            entity.hd1Tags = ['session-avatar', 'other-session', 'visual-only'];
            entity.hd1IsOtherSession = true;
            
            // Add to scene
            app.root.addChild(entity);
            
            console.log(`[HD1-WebSocket] ✅ Created visual representation for OTHER session avatar: ${avatar_name}`);
            this.trackAvatarSync(entity, 'other_visual_created');
            
        } catch (error) {
            console.error('[HD1-WebSocket] Failed to create other session avatar visual:', error);
        }
    }
    
    /**
     * Create basic avatar entity when none exists
     */
    createBasicAvatarEntity(session_id, avatar_name, position) {
        const avatarData = {
            name: avatar_name,
            tags: ['session-avatar', 'avatar'],
            enabled: true,
            components: {
                transform: {
                    position: [position.x, position.y, position.z],
                    rotation: [0, 0, 0, 1],
                    scale: [1, 1, 1]
                },
                model: {
                    type: 'box',  // Basic placeholder
                    material: {
                        ambient: [0.3, 0.6, 1.0],
                        diffuse: [0.3, 0.6, 1.0]
                    }
                }
            }
        };
        
        // Create entity via API first
        window.hd1API.createEntity(session_id, avatarData)
            .then(response => {
                if (response.success) {
                    console.log(`[HD1-WebSocket] ✅ Created basic avatar entity via API`);
                    
                    // Then render in PlayCanvas
                    if (window.createObjectFromData) {
                        avatarData.entity_id = response.entity_id;
                        window.createObjectFromData(avatarData);
                        console.log(`[HD1-WebSocket] ✅ Rendered basic avatar in PlayCanvas`);
                    }
                }
            })
            .catch(error => {
                console.error('[HD1-WebSocket] Failed to create basic avatar entity:', error);
            });
    }
    
    /**
     * Track avatar synchronization events for stats and debugging
     */
    trackAvatarSync(avatarEntity, eventType) {
        if (window.hd1ConsoleManager) {
            const statsManager = window.hd1ConsoleManager.getModule('stats');
            if (statsManager && statsManager.trackEntitySync) {
                statsManager.trackEntitySync(avatarEntity.hd1Id || avatarEntity.name, eventType);
            }
        }
        
        // Enhanced debug logging
        console.log(`[HD1-WebSocket] 📊 Avatar sync event: ${eventType}`, {
            entity_name: avatarEntity.name,
            entity_id: avatarEntity.hd1Id,
            position: avatarEntity.getPosition(),
            tags: avatarEntity.hd1Tags
        });
    }

    /**
     * Track WebSocket traffic for stats
     */
    trackTraffic(inBytes, outBytes) {
        // Get stats manager and track traffic
        if (window.hd1ConsoleManager) {
            const statsManager = window.hd1ConsoleManager.getModule('stats');
            if (statsManager && statsManager.trackWebSocketTraffic) {
                // Update cumulative totals
                this.totalInBytes = (this.totalInBytes || 0) + inBytes;
                this.totalOutBytes = (this.totalOutBytes || 0) + outBytes;
                
                // Pass to stats manager
                statsManager.trackWebSocketTraffic(this.totalInBytes, this.totalOutBytes);
            }
        }
    }

    /**
     * Get connection status
     */
    getConnectionStatus() {
        return {
            connected: this.isConnected,
            connecting: this.isConnecting,
            reconnectAttempts: this.reconnectAttempts,
            lastMessageTime: this.lastMessageTime,
            totalInBytes: this.totalInBytes || 0,
            totalOutBytes: this.totalOutBytes || 0
        };
    }

    /**
     * 🔧 SCENE PROTECTION: Preserve avatars during channel clearing
     */
    preserveAvatarsBeforeSceneClear() {
        if (!window.hd1GameEngine || !this.entityProtectionMode) {
            return;
        }
        
        const protectedEntities = [];
        
        // Find all protected avatar entities
        window.hd1GameEngine.root.children.forEach(entity => {
            if (entity.hd1Protected || 
                (entity.hd1Tags && entity.hd1Tags.includes('session-avatar')) ||
                this.avatarRegistry.has(entity.hd1SessionId)) {
                protectedEntities.push({
                    entity: entity,
                    sessionId: entity.hd1SessionId,
                    position: entity.getPosition().clone(),
                    name: entity.name
                });
            }
        });
        
        console.log(`[HD1-WebSocket] 🛡️ Preserving ${protectedEntities.length} avatar entities before scene clear`);
        
        // Store for restoration
        this.preservedAvatars = protectedEntities;
        
        return protectedEntities;
    }
    
    /**
     * 🔧 SCENE PROTECTION: Restore avatars after channel clearing
     */
    restoreAvatarsAfterSceneClear() {
        if (!this.preservedAvatars || this.preservedAvatars.length === 0) {
            return;
        }
        
        setTimeout(() => {
            console.log(`[HD1-WebSocket] 🔄 Restoring ${this.preservedAvatars.length} preserved avatars`);
            
            this.preservedAvatars.forEach(preserved => {
                // Clear old registry entry
                this.avatarRegistry.delete(preserved.sessionId);
                this.creationLocks.delete(preserved.sessionId);
                
                // Force recreation with preserved position
                if (preserved.sessionId !== this.getCurrentSessionId()) {
                    this.createOtherSessionAvatarProtected(
                        preserved.sessionId, 
                        preserved.name, 
                        preserved.position
                    );
                }
            });
            
            this.preservedAvatars = [];
        }, 200); // Allow scene clearing to complete
    }

    /**
     * Handle avatar asset response - receives GLB data via WebSocket
     */
    handleAvatarAssetResponse(message) {
        console.log('[HD1-WebSocket] Avatar asset response received:', message.avatar_type, `${message.size} bytes`);
        
        try {
            // Forward to PlayCanvas integration for GLB processing
            if (window.handleAvatarAssetResponse) {
                window.handleAvatarAssetResponse(message.avatar_type, message.data);
            } else {
                console.warn('[HD1-WebSocket] handleAvatarAssetResponse not available in PlayCanvas integration');
            }
        } catch (error) {
            console.error('[HD1-WebSocket] Failed to handle avatar asset response:', error);
        }
        
        this.dom.addDebugEntry('AVATAR_ASSET_RESPONSE', {
            avatar_type: message.avatar_type,
            size: message.size
        });
    }

    /**
     * Handle avatar asset error - GLB asset loading failed
     */
    handleAvatarAssetError(message) {
        console.error('[HD1-WebSocket] Avatar asset error:', message.avatar_type, message.error);
        
        this.dom.addDebugEntry('AVATAR_ASSET_ERROR', {
            avatar_type: message.avatar_type,
            error: message.error
        });
    }

    /**
     * Cleanup WebSocket manager
     */
    cleanup() {
        if (this.reconnectTimeout) {
            clearTimeout(this.reconnectTimeout);
        }
        
        if (this.heartbeatInterval) {
            clearInterval(this.heartbeatInterval);
        }
        
        if (this.positionUpdateTimer) {
            clearTimeout(this.positionUpdateTimer);
        }
        
        if (this.ws) {
            this.ws.close();
        }
        
        this.eventHandlers.clear();
        this.messageQueue = [];
        this.pingTimestamps.clear();
        
        // Clear avatar protection state
        this.avatarRegistry.clear();
        this.positionUpdateQueue.clear();
        this.creationLocks.clear();
        this.activeCreations = 0;
        
        this.ready = false;
    }
    
    /**
     * 🛡️ ROBUSTNESS: Start avatar health monitoring
     */
    startAvatarHealthMonitoring() {
        if (this.avatarHealthChecker) {
            clearInterval(this.avatarHealthChecker);
        }
        
        this.avatarHealthChecker = setInterval(() => {
            this.performAvatarHealthCheck();
        }, this.healthCheckInterval);
        
        console.log('[HD1-WebSocket] 🛡️ Avatar health monitoring started');
    }
    
    /**
     * 🛡️ ROBUSTNESS: Stop avatar health monitoring
     */
    stopAvatarHealthMonitoring() {
        if (this.avatarHealthChecker) {
            clearInterval(this.avatarHealthChecker);
            this.avatarHealthChecker = null;
        }
    }
    
    /**
     * 🛡️ ROBUSTNESS: Perform avatar health check
     */
    performAvatarHealthCheck() {
        if (!window.hd1GameEngine || !this.isConnected) return;
        
        try {
            const app = window.hd1GameEngine;
            let recoveredAvatars = 0;
            
            // Check each registered avatar
            this.avatarRegistry.forEach((avatarInfo, sessionId) => {
                const entity = avatarInfo.entity;
                
                // Check if entity still exists in scene
                if (entity && entity.parent === null) {
                    console.warn(`[HD1-WebSocket] 🛡️ Avatar health issue detected for session ${sessionId} - entity orphaned`);
                    this.recoverAvatar(sessionId, avatarInfo);
                    recoveredAvatars++;
                }
                
                // Check if entity is properly tagged
                if (entity && (!entity.hd1Tags || !entity.hd1Tags.includes('session-avatar'))) {
                    console.warn(`[HD1-WebSocket] 🛡️ Avatar health issue detected for session ${sessionId} - missing tags`);
                    entity.hd1Tags = entity.hd1Tags || [];
                    if (!entity.hd1Tags.includes('session-avatar')) {
                        entity.hd1Tags.push('session-avatar');
                    }
                }
            });
            
            if (recoveredAvatars > 0) {
                console.log(`[HD1-WebSocket] 🛡️ Avatar health check completed - recovered ${recoveredAvatars} avatars`);
            }
            
        } catch (error) {
            console.error('[HD1-WebSocket] 🛡️ Avatar health check failed:', error);
        }
    }
    
    /**
     * 🛡️ ROBUSTNESS: Recover a damaged avatar
     */
    recoverAvatar(sessionId, avatarInfo) {
        const attempts = this.avatarRecoveryAttempts.get(sessionId) || 0;
        
        if (attempts >= this.maxRecoveryAttempts) {
            console.warn(`[HD1-WebSocket] 🛡️ Max recovery attempts reached for session ${sessionId}`);
            this.avatarRecoveryAttempts.delete(sessionId);
            return;
        }
        
        this.avatarRecoveryAttempts.set(sessionId, attempts + 1);
        
        try {
            const app = window.hd1GameEngine;
            
            // Re-add to scene if orphaned
            if (avatarInfo.entity && avatarInfo.entity.parent === null) {
                app.root.addChild(avatarInfo.entity);
                console.log(`[HD1-WebSocket] 🛡️ Avatar recovered for session ${sessionId} (attempt ${attempts + 1})`);
            }
            
        } catch (error) {
            console.error(`[HD1-WebSocket] 🛡️ Avatar recovery failed for session ${sessionId}:`, error);
        }
    }
}

// Export for use in console manager
window.HD1WebSocketManager = HD1WebSocketManager;