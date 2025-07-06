/**
 * HD1 World Manager - 100% API-Driven World Operations
 * 
 * Manages all world-related functionality with bulletproof error handling.
 * All operations go through the HD1 API specification.
 */
class HD1WorldManager {
    constructor(apiClient, domManager) {
        this.api = apiClient;
        this.dom = domManager;
        this.currentWorld = null;
        this.worlds = [];
        this.ready = false;
    }

    /**
     * Initialize world manager
     */
    async initialize() {
        console.log('[HD1-World] Initializing World Manager');
        
        try {
            // Load worlds from API
            await this.loadWorlds();
            
            // Setup world selector
            this.setupWorldSelector();
            
            // Load saved world or default
            await this.loadInitialWorld();
            
            this.ready = true;
            console.log('[HD1-World] World Manager ready');
        } catch (error) {
            console.error('[HD1-World] Initialization failed:', error);
            this.handleWorldError('Initialization failed');
        }
    }

    /**
     * Load worlds from API with error handling
     */
    async loadWorlds() {
        try {
            console.log('[HD1-World] Loading worlds from API...');
            
            const response = await this.api.listWorlds();
            
            if (response.success && response.worlds) {
                this.worlds = response.worlds.filter(world => world.id); // Remove empty worlds
                console.log(`[HD1-World] Loaded ${this.worlds.length} worlds`);
                return this.worlds;
            } else {
                throw new Error('Invalid API response format');
            }
        } catch (error) {
            console.error('[HD1-World] Failed to load worlds:', error);
            this.worlds = [];
            throw error;
        }
    }

    /**
     * Setup world selector dropdown
     */
    setupWorldSelector() {
        if (!this.dom.exists('channel-selector')) {
            console.warn('[HD1-World] Channel selector element not found');
            return;
        }

        // Clear existing options
        this.dom.setHTML('channel-selector', '');
        
        // Add default option
        const selector = this.dom.get('channel-selector');
        const defaultOption = document.createElement('option');
        defaultOption.value = '';
        defaultOption.textContent = this.worlds.length > 0 ? 'Select World...' : 'No worlds available';
        selector.appendChild(defaultOption);
        
        // Add worlds
        this.worlds.forEach(world => {
            const option = document.createElement('option');
            option.value = world.id;
            option.textContent = world.name;
            option.title = world.description;
            selector.appendChild(option);
        });

        // Setup change handler
        this.dom.addEventListener('channel-selector', 'change', (e) => {
            const worldId = e.target.value;
            if (worldId) {
                this.switchWorld(worldId);
            }
        });

        console.log('[HD1-World] World selector setup complete');
    }

    /**
     * Load initial world (saved or default)
     */
    async loadInitialWorld() {
        // Try to load saved world
        const savedWorldId = localStorage.getItem('hd1_current_world');
        
        if (savedWorldId && this.worlds.find(c => c.id === savedWorldId)) {
            await this.switchWorld(savedWorldId);
            return;
        }

        // Load default world from API response
        const defaultWorldId = this.getDefaultWorldId();
        if (defaultWorldId) {
            await this.switchWorld(defaultWorldId);
            return;
        }

        console.log('[HD1-World] No initial world to load');
    }

    /**
     * Get default world ID from server config or first available
     */
    getDefaultWorldId() {
        // If API returned default_world, use that
        if (this.worlds.length > 0) {
            return this.worlds[0].id;
        }
        return null;
    }

    /**
     * Switch to specified world
     */
    async switchWorld(worldId) {
        if (!worldId) {
            console.warn('[HD1-World] No world ID provided');
            return false;
        }

        const world = this.worlds.find(c => c.id === worldId);
        if (!world) {
            console.error('[HD1-World] World not found:', worldId);
            return false;
        }

        try {
            console.log(`[HD1-World] Switching to world: ${worldId}`);
            
            // Get current session ID for world join
            const sessionId = this.getCurrentSessionId();
            if (!sessionId) {
                throw new Error('No active session found to join world');
            }
            
            // PHASE 2 FIX: Pure API call - server handles clearing + loading + broadcasting
            // No client-side business logic - just call API and let WebSocket handle rendering
            await this.joinSessionToWorld(sessionId, worldId);
            
            // Get detailed world info from API for UI updates
            const worldResponse = await this.api.getWorld(worldId);
            
            if (!worldResponse.success) {
                throw new Error('Failed to get world details');
            }

            const worldData = worldResponse.world;
            
            // Update current world state
            this.currentWorld = worldData;
            
            // Save to localStorage
            localStorage.setItem('hd1_current_world', worldId);
            
            // Update UI only - no entity management (handled by WebSocket broadcasts)
            this.updateWorldUI(worldData);
            
            // Update scene name in PlayCanvas panel
            this.updateSceneName(world.name);
            
            // CRITICAL FIX: Notify session manager about world switch for avatar control recovery
            if (window.hd1ConsoleManager) {
                const sessionManager = window.hd1ConsoleManager.getModule('session');
                if (sessionManager && sessionManager.handleEvent) {
                    sessionManager.handleEvent('world_switched', {
                        worldId: worldId,
                        worldName: world.name,
                        sessionId: sessionId
                    });
                }
            }
            
            console.log(`[HD1-World] Successfully switched to world: ${worldId}`);
            return true;
            
        } catch (error) {
            console.error('[HD1-World] Failed to switch world:', error);
            this.handleWorldError(`Failed to switch to ${world.name}`);
            return false;
        }
    }

    /**
     * Update world-related UI elements
     */
    updateWorldUI(worldData) {
        // Update world selector
        const selector = this.dom.get('channel-selector');
        if (selector) {
            selector.value = worldData.id;
        }
        
        // Update any world info displays
        this.dom.addDebugEntry('WORLD_SWITCH', {
            world: worldData.id,
            name: worldData.name,
            clients: worldData.current_clients || 0
        });
    }

    /**
     * Load PlayCanvas content from world configuration
     */
    async loadPlayCanvasContent(playcanvasConfig) {
        try {
            console.log('[HD1-World] Loading PlayCanvas content');
            
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
                    console.log('[HD1-World] Applied ambient light:', ambientLight);
                }
                
                // Apply gravity
                const gravity = sceneConfig.Gravity || sceneConfig.gravity;
                if (gravity && app.systems.rigidbody) {
                    app.systems.rigidbody.gravity.set(gravity[0], gravity[1], gravity[2]);
                    console.log('[HD1-World] Applied gravity:', gravity);
                }
            }
            
            // PHASE 2 FIX: No manual entity loading - entities automatically rendered via WebSocket broadcasts
            
        } catch (error) {
            console.error('[HD1-World] Failed to load PlayCanvas content:', error);
        }
    }

    /**
     * Create PlayCanvas entities from configuration
     */
    async createPlayCanvasEntities(entitiesConfig) {
        if (!window.hd1GameEngine) {
            console.warn('[HD1-World] PlayCanvas engine not ready');
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
        
        console.log(`[HD1-World] Clearing ${entitiesToRemove.length} entities, preserving ${preservedAvatars.length} avatars`);
        
        entitiesToRemove.forEach(entity => {
            app.root.removeChild(entity);
            entity.destroy();
        });
        
        // Reset ambient light to black - world will set proper lighting
        app.scene.ambientLight = new pc.Color(0, 0, 0);
        console.log('[HD1-World] Scene cleared, ready for new world content');
        
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
            
            console.log(`[HD1-World] Creating entity: ${entityName}`);
            
            // Add components (handle both capitalized and lowercase)
            const components = config.Components || config.components;
            if (components) {
                for (const [componentType, componentData] of Object.entries(components)) {
                    console.log(`[HD1-World] Adding ${componentType} component to ${entityName}`);
                    entity.addComponent(componentType, componentData);
                }
            }
            
            app.root.addChild(entity);
            console.log(`[HD1-World] Successfully created entity: ${entityName}`);
            
        } catch (error) {
            const entityName = config.Name || config.name || 'unknown';
            console.error(`[HD1-World] Failed to create entity ${entityName}:`, error);
        }
    }

    /**
     * Update scene name in PlayCanvas panel
     */
    updateSceneName(worldName) {
        const sceneName = worldName.split(' - ')[0]; // Get first part before " - "
        this.dom.setText('pc-scene-name', sceneName);
    }

    /**
     * Handle world-related errors
     */
    handleWorldError(message) {
        console.error(`[HD1-World] Error: ${message}`);
        
        // Update UI to show error state
        this.dom.setHTML('channel-selector', '<option value="">Error loading worlds</option>');
        
        // Log error
        this.dom.addDebugEntry('WORLD_ERROR', { message });
    }

    /**
     * Get current world info
     */
    getCurrentWorld() {
        return this.currentWorld;
    }

    /**
     * Get all available worlds
     */
    getWorlds() {
        return this.worlds;
    }

    /**
     * Refresh worlds from API
     */
    async refreshWorlds() {
        try {
            await this.loadWorlds();
            this.setupWorldSelector();
            console.log('[HD1-World] Worlds refreshed');
            return true;
        } catch (error) {
            console.error('[HD1-World] Failed to refresh worlds:', error);
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
        
        console.warn('[HD1-World] No session ID found');
        return null;
    }

    // PHASE 2 FIX: Removed clearSessionEntities() - violates API-first principle
    // Entity clearing now handled server-side in world join endpoint with proper WebSocket broadcasts

    /**
     * Join current session to specified world
     */
    async joinSessionToWorld(sessionId, worldId) {
        try {
            console.log(`[HD1-World] Joining session ${sessionId} to world ${worldId}`);
            
            const response = await fetch(`/api/sessions/${sessionId}/world/join`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    client_id: `client_${sessionId}`,
                    world_id: worldId,
                    reconnect: false
                })
            });

            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }

            const data = await response.json();
            
            if (data.success) {
                console.log(`[HD1-World] Successfully joined session to world: ${worldId}`);
                console.log('[HD1-World] World join response:', data);
                
                // Store world association
                localStorage.setItem('hd1_session_world', worldId);
                
                return data;
            } else {
                throw new Error(data.message || 'World join failed');
            }
        } catch (error) {
            console.error('[HD1-World] Failed to join session to world:', error);
            throw error;
        }
    }
}

// Export for use in console manager
window.HD1WorldManager = HD1WorldManager;