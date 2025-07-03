/**
 * HD1 Channel Manager - 100% API-Driven Channel Operations
 * 
 * Manages all channel-related functionality with bulletproof error handling.
 * All operations go through the HD1 API specification.
 */
class HD1ChannelManager {
    constructor(apiClient, domManager) {
        this.api = apiClient;
        this.dom = domManager;
        this.currentChannel = null;
        this.channels = [];
        this.ready = false;
    }

    /**
     * Initialize channel manager
     */
    async initialize() {
        console.log('[HD1-Channel] Initializing Channel Manager');
        
        try {
            // Load channels from API
            await this.loadChannels();
            
            // Setup channel selector
            this.setupChannelSelector();
            
            // Load saved channel or default
            await this.loadInitialChannel();
            
            this.ready = true;
            console.log('[HD1-Channel] Channel Manager ready');
        } catch (error) {
            console.error('[HD1-Channel] Initialization failed:', error);
            this.handleChannelError('Initialization failed');
        }
    }

    /**
     * Load channels from API with error handling
     */
    async loadChannels() {
        try {
            console.log('[HD1-Channel] Loading channels from API...');
            
            const response = await this.api.listChannels();
            
            if (response.success && response.channels) {
                this.channels = response.channels.filter(channel => channel.id); // Remove empty channels
                console.log(`[HD1-Channel] Loaded ${this.channels.length} channels`);
                return this.channels;
            } else {
                throw new Error('Invalid API response format');
            }
        } catch (error) {
            console.error('[HD1-Channel] Failed to load channels:', error);
            this.channels = [];
            throw error;
        }
    }

    /**
     * Setup channel selector dropdown
     */
    setupChannelSelector() {
        if (!this.dom.exists('channel-selector')) {
            console.warn('[HD1-Channel] Channel selector element not found');
            return;
        }

        // Clear existing options
        this.dom.setHTML('channel-selector', '');
        
        // Add default option
        const selector = this.dom.get('channel-selector');
        const defaultOption = document.createElement('option');
        defaultOption.value = '';
        defaultOption.textContent = this.channels.length > 0 ? 'Select Channel...' : 'No channels available';
        selector.appendChild(defaultOption);
        
        // Add channels
        this.channels.forEach(channel => {
            const option = document.createElement('option');
            option.value = channel.id;
            option.textContent = channel.name;
            option.title = channel.description;
            selector.appendChild(option);
        });

        // Setup change handler
        this.dom.addEventListener('channel-selector', 'change', (e) => {
            const channelId = e.target.value;
            if (channelId) {
                this.switchChannel(channelId);
            }
        });

        console.log('[HD1-Channel] Channel selector setup complete');
    }

    /**
     * Load initial channel (saved or default)
     */
    async loadInitialChannel() {
        // Try to load saved channel
        const savedChannelId = localStorage.getItem('hd1_current_channel');
        
        if (savedChannelId && this.channels.find(c => c.id === savedChannelId)) {
            await this.switchChannel(savedChannelId);
            return;
        }

        // Load default channel from API response
        const defaultChannelId = this.getDefaultChannelId();
        if (defaultChannelId) {
            await this.switchChannel(defaultChannelId);
            return;
        }

        console.log('[HD1-Channel] No initial channel to load');
    }

    /**
     * Get default channel ID from server config or first available
     */
    getDefaultChannelId() {
        // If API returned default_channel, use that
        if (this.channels.length > 0) {
            return this.channels[0].id;
        }
        return null;
    }

    /**
     * Switch to specified channel
     */
    async switchChannel(channelId) {
        if (!channelId) {
            console.warn('[HD1-Channel] No channel ID provided');
            return false;
        }

        const channel = this.channels.find(c => c.id === channelId);
        if (!channel) {
            console.error('[HD1-Channel] Channel not found:', channelId);
            return false;
        }

        try {
            console.log(`[HD1-Channel] Switching to channel: ${channelId}`);
            
            // Get current session ID for channel join
            const sessionId = this.getCurrentSessionId();
            if (!sessionId) {
                throw new Error('No active session found to join channel');
            }
            
            // PHASE 2 FIX: Pure API call - server handles clearing + loading + broadcasting
            // No client-side business logic - just call API and let WebSocket handle rendering
            await this.joinSessionToChannel(sessionId, channelId);
            
            // Get detailed channel info from API for UI updates
            const channelResponse = await this.api.getChannel(channelId);
            
            if (!channelResponse.success) {
                throw new Error('Failed to get channel details');
            }

            const channelData = channelResponse.channel;
            
            // Update current channel state
            this.currentChannel = channelData;
            
            // Save to localStorage
            localStorage.setItem('hd1_current_channel', channelId);
            
            // Update UI only - no entity management (handled by WebSocket broadcasts)
            this.updateChannelUI(channelData);
            
            // Update scene name in PlayCanvas panel
            this.updateSceneName(channel.name);
            
            console.log(`[HD1-Channel] Successfully switched to channel: ${channelId}`);
            return true;
            
        } catch (error) {
            console.error('[HD1-Channel] Failed to switch channel:', error);
            this.handleChannelError(`Failed to switch to ${channel.name}`);
            return false;
        }
    }

    /**
     * Update channel-related UI elements
     */
    updateChannelUI(channelData) {
        // Update channel selector
        const selector = this.dom.get('channel-selector');
        if (selector) {
            selector.value = channelData.id;
        }
        
        // Update any channel info displays
        this.dom.addDebugEntry('CHANNEL_SWITCH', {
            channel: channelData.id,
            name: channelData.name,
            clients: channelData.current_clients || 0
        });
    }

    /**
     * Load PlayCanvas content from channel configuration
     */
    async loadPlayCanvasContent(playcanvasConfig) {
        try {
            console.log('[HD1-Channel] Loading PlayCanvas content');
            
            // PHASE 2 FIX: No client-side entity clearing - handled by WebSocket broadcasts
            
            // Apply scene settings (handle both capitalized and lowercase property names)
            const sceneConfig = playcanvasConfig.Scene || playcanvasConfig.scene;
            if (sceneConfig && window.hd1GameEngine) {
                const app = window.hd1GameEngine;
                
                // Apply ambient light (handle both string and array formats)
                const ambientLight = sceneConfig.AmbientLight || sceneConfig.ambientLight;
                if (ambientLight) {
                    let color;
                    if (Array.isArray(ambientLight)) {
                        // RGB array format [r, g, b]
                        color = new pc.Color(ambientLight[0] || 0, ambientLight[1] || 0, ambientLight[2] || 0);
                    } else {
                        // Hex string format
                        color = new pc.Color().fromString(ambientLight);
                    }
                    app.scene.ambientLight = color;
                    console.log('[HD1-Channel] Applied ambient light:', ambientLight);
                }
                
                // Apply gravity
                const gravity = sceneConfig.Gravity || sceneConfig.gravity;
                if (gravity && app.systems.rigidbody) {
                    app.systems.rigidbody.gravity.set(gravity[0], gravity[1], gravity[2]);
                    console.log('[HD1-Channel] Applied gravity:', gravity);
                }
            }
            
            // PHASE 2 FIX: No manual entity loading - entities automatically rendered via WebSocket broadcasts
            
        } catch (error) {
            console.error('[HD1-Channel] Failed to load PlayCanvas content:', error);
        }
    }

    /**
     * Create PlayCanvas entities from configuration
     */
    async createPlayCanvasEntities(entitiesConfig) {
        if (!window.hd1GameEngine) {
            console.warn('[HD1-Channel] PlayCanvas engine not ready');
            return;
        }
        
        const app = window.hd1GameEngine;
        
        // 🛡️ PRESERVE avatars before clearing scene
        let preservedAvatars = [];
        if (window.hd1ConsoleManager) {
            const wsManager = window.hd1ConsoleManager.getModule('websocket');
            if (wsManager && wsManager.preserveAvatarsBeforeSceneClear) {
                preservedAvatars = wsManager.preserveAvatarsBeforeSceneClear() || [];
            }
        }
        
        // Clear existing entities (except camera and protected avatars)
        const entitiesToRemove = app.root.children.filter(entity => 
            entity.name !== 'camera' && 
            entity.hd1Entity !== false &&
            !entity.hd1Protected &&
            !(entity.hd1Tags && entity.hd1Tags.includes('session-avatar'))
        );
        
        console.log(`[HD1-Channel] Clearing ${entitiesToRemove.length} entities, preserving ${preservedAvatars.length} avatars`);
        
        entitiesToRemove.forEach(entity => {
            app.root.removeChild(entity);
            entity.destroy();
        });
        
        // Reset ambient light to black - channel will set proper lighting
        app.scene.ambientLight = new pc.Color(0, 0, 0);
        console.log('[HD1-Channel] Scene cleared, ready for new channel content');
        
        // Create each entity
        for (const entityConfig of entitiesConfig) {
            await this.createNativePlayCanvasEntity(entityConfig, app);
        }
        
        // 🔄 RESTORE preserved avatars after scene creation
        if (preservedAvatars.length > 0 && window.hd1ConsoleManager) {
            const wsManager = window.hd1ConsoleManager.getModule('websocket');
            if (wsManager && wsManager.restoreAvatarsAfterSceneClear) {
                wsManager.restoreAvatarsAfterSceneClear();
            }
        }
    }

    /**
     * Create native PlayCanvas entity from configuration
     */
    async createNativePlayCanvasEntity(config, app) {
        try {
            // Handle both capitalized and lowercase property names
            const entityName = config.Name || config.name;
            const entity = new pc.Entity(entityName);
            
            console.log(`[HD1-Channel] Creating entity: ${entityName}`);
            
            // Add components (handle both capitalized and lowercase)
            const components = config.Components || config.components;
            if (components) {
                for (const [componentType, componentData] of Object.entries(components)) {
                    console.log(`[HD1-Channel] Adding ${componentType} component to ${entityName}`);
                    entity.addComponent(componentType, componentData);
                }
            }
            
            app.root.addChild(entity);
            console.log(`[HD1-Channel] Successfully created entity: ${entityName}`);
            
        } catch (error) {
            const entityName = config.Name || config.name || 'unknown';
            console.error(`[HD1-Channel] Failed to create entity ${entityName}:`, error);
        }
    }

    /**
     * Update scene name in PlayCanvas panel
     */
    updateSceneName(channelName) {
        const sceneName = channelName.split(' - ')[0]; // Get first part before " - "
        this.dom.setText('pc-scene-name', sceneName);
    }

    /**
     * Handle channel-related errors
     */
    handleChannelError(message) {
        console.error(`[HD1-Channel] Error: ${message}`);
        
        // Update UI to show error state
        this.dom.setHTML('channel-selector', '<option value="">Error loading channels</option>');
        
        // Log error
        this.dom.addDebugEntry('CHANNEL_ERROR', { message });
    }

    /**
     * Get current channel info
     */
    getCurrentChannel() {
        return this.currentChannel;
    }

    /**
     * Get all available channels
     */
    getChannels() {
        return this.channels;
    }

    /**
     * Refresh channels from API
     */
    async refreshChannels() {
        try {
            await this.loadChannels();
            this.setupChannelSelector();
            console.log('[HD1-Channel] Channels refreshed');
            return true;
        } catch (error) {
            console.error('[HD1-Channel] Failed to refresh channels:', error);
            return false;
        }
    }

    /**
     * Get current session ID from session manager
     */
    getCurrentSessionId() {
        // Try multiple sources for session ID
        if (window.hd1Console && window.hd1Console.sessionManager) {
            return window.hd1Console.sessionManager.getSessionId();
        }
        
        // Fallback to localStorage
        const sessionId = localStorage.getItem('hd1_session_id');
        if (sessionId) {
            return sessionId;
        }
        
        console.warn('[HD1-Channel] No session ID found');
        return null;
    }

    // PHASE 2 FIX: Removed clearSessionEntities() - violates API-first principle
    // Entity clearing now handled server-side in channel join endpoint with proper WebSocket broadcasts

    /**
     * Join current session to specified channel
     */
    async joinSessionToChannel(sessionId, channelId) {
        try {
            console.log(`[HD1-Channel] Joining session ${sessionId} to channel ${channelId}`);
            
            const response = await fetch(`/api/sessions/${sessionId}/channel/join`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    client_id: `client_${sessionId}`,
                    channel_id: channelId,
                    reconnect: false
                })
            });

            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }

            const data = await response.json();
            
            if (data.success) {
                console.log(`[HD1-Channel] Successfully joined session to channel: ${channelId}`);
                console.log('[HD1-Channel] Channel join response:', data);
                
                // Store channel association
                localStorage.setItem('hd1_session_channel', channelId);
                
                return data;
            } else {
                throw new Error(data.message || 'Channel join failed');
            }
        } catch (error) {
            console.error('[HD1-Channel] Failed to join session to channel:', error);
            throw error;
        }
    }
}

// Export for use in console manager
window.HD1ChannelManager = HD1ChannelManager;