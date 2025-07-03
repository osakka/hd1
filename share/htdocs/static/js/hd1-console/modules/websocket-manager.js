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
            
            // Use existing PlayCanvas entity deletion function
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
            // For updates, remove and recreate entity to ensure consistency
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
     * Cleanup WebSocket manager
     */
    cleanup() {
        if (this.reconnectTimeout) {
            clearTimeout(this.reconnectTimeout);
        }
        
        if (this.heartbeatInterval) {
            clearInterval(this.heartbeatInterval);
        }
        
        if (this.ws) {
            this.ws.close();
        }
        
        this.eventHandlers.clear();
        this.messageQueue = [];
        this.pingTimestamps.clear();
        this.ready = false;
    }
}

// Export for use in console manager
window.HD1WebSocketManager = HD1WebSocketManager;