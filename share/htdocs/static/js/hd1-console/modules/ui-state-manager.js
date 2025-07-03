/**
 * HD1 UI State Manager - Panel State and User Interface Management
 * 
 * Manages console panel states, user preferences, and UI interactions.
 * Handles collapse/expand, persistence, and responsive behavior.
 */
class HD1UIStateManager {
    constructor(domManager) {
        this.dom = domManager;
        this.ready = false;
        
        // UI state
        this.isCollapsed = false;
        this.refreshTimer = null;
        this.refreshStartTime = Date.now();
        
        // Version info
        this.apiVersion = 'loading...';
        this.jsVersion = 'loading...';
        
        // Timer state
        this.timerInterval = null;
    }

    /**
     * Initialize UI state manager
     */
    async initialize() {
        console.log('[HD1-UI] Initializing UI State Manager');
        
        try {
            // Load saved UI state
            this.loadUIState();
            
            // Setup UI event handlers
            this.setupUIHandlers();
            
            // Initialize version display
            await this.loadVersionInfo();
            
            // Start refresh timer
            this.startRefreshTimer();
            
            // Setup responsive behavior
            this.setupResponsiveBehavior();
            
            // Session management is now handled by dedicated SessionManager module
            
            this.ready = true;
            console.log('[HD1-UI] UI State Manager ready');
        } catch (error) {
            console.error('[HD1-UI] Initialization failed:', error);
        }
    }

    /**
     * Load saved UI state from localStorage
     */
    loadUIState() {
        // Load collapse state
        const savedCollapsed = localStorage.getItem('hd1_console_collapsed');
        this.isCollapsed = savedCollapsed === 'true';
        
        // Apply initial state
        this.setCollapseState(this.isCollapsed, false); // Don't save on initial load
    }

    /**
     * Setup UI event handlers
     */
    setupUIHandlers() {
        // Debug header click handler for collapse/expand
        this.dom.addEventListener('debug-header', 'click', (e) => {
            this.toggleCollapse();
        });

        // Status bar click handler (secondary collapse method)
        this.dom.addEventListener('debug-status-bar', 'click', (e) => {
            // Don't toggle if clicking on session ID tag
            const sessionTag = this.dom.get('session-id-tag-status');
            if (e.target !== sessionTag) {
                this.toggleCollapse();
            }
        });

        // Session ID tag click handler for copying (handled by session manager)
        this.dom.addEventListener('session-id-tag-status', 'click', (e) => {
            const sessionManager = window.hd1ConsoleManager?.getModule('session');
            if (sessionManager && sessionManager.copySessionId) {
                sessionManager.copySessionId();
            }
            e.stopPropagation(); // Prevent triggering collapse
        });

        // Refresh timer click handler for rebootstrap
        this.dom.addEventListener('refresh-timer', 'click', (e) => {
            this.triggerRebootstrap();
        });

        // Window resize handler
        window.addEventListener('resize', () => {
            this.handleResize();
        });
    }

    /**
     * Toggle collapse state
     */
    toggleCollapse() {
        this.setCollapseState(!this.isCollapsed, true);
    }

    /**
     * Set collapse state
     */
    setCollapseState(collapsed, save = true) {
        this.isCollapsed = collapsed;
        
        // Update UI elements
        const debugPanel = this.dom.get('debug-panel');
        const debugContent = this.dom.get('debug-content');
        const debugInfoPanels = this.dom.get('debug-info-panels');
        const debugCollapsedStats = this.dom.get('debug-collapsed-stats');
        const collapseIcon = this.dom.get('debug-collapse-icon');
        const statusArrow = this.dom.get('status-collapse-arrow');
        
        // Add/remove collapsed class from main panel (for CSS styling)
        if (debugPanel) {
            if (collapsed) {
                debugPanel.classList.add('collapsed');
            } else {
                debugPanel.classList.remove('collapsed');
            }
        }
        
        // Handle panel visibility with smooth transitions
        if (debugInfoPanels) {
            if (collapsed) {
                // Start animation to collapsed state
                debugInfoPanels.style.height = '0';
                debugInfoPanels.style.opacity = '0';
                debugInfoPanels.style.margin = '0 8px';
                debugInfoPanels.style.padding = '0 4px';
                // Hide after animation completes
                setTimeout(() => {
                    if (this.isCollapsed) { // Check if still collapsed
                        debugInfoPanels.style.display = 'none';
                    }
                }, 400);
            } else {
                // Show first, then animate to expanded state
                debugInfoPanels.style.display = 'grid';
                setTimeout(() => {
                    debugInfoPanels.style.height = '220px';
                    debugInfoPanels.style.opacity = '1';
                    debugInfoPanels.style.margin = '8px 8px 0 8px';
                    debugInfoPanels.style.padding = '4px';
                }, 10);
            }
        }
        
        if (debugCollapsedStats) {
            if (collapsed) {
                // Show collapsed stats with animation
                debugCollapsedStats.style.display = 'grid';
                setTimeout(() => {
                    debugCollapsedStats.style.height = '50px';
                    debugCollapsedStats.style.opacity = '1';
                }, 10);
            } else {
                // Hide collapsed stats with animation
                debugCollapsedStats.style.height = '0';
                debugCollapsedStats.style.opacity = '0';
                setTimeout(() => {
                    if (!this.isCollapsed) { // Check if still expanded
                        debugCollapsedStats.style.display = 'none';
                    }
                }, 400);
            }
        }
        
        // If uncollapsing, ensure graphs are initialized and render them
        if (!collapsed && window.hd1ConsoleManager) {
            setTimeout(() => {
                const statsManager = window.hd1ConsoleManager.getModule('stats');
                if (statsManager) {
                    // Ensure graphs are initialized first (handles deferred initialization)
                    if (statsManager.ensureGraphsInitialized) {
                        statsManager.ensureGraphsInitialized();
                    }
                    // Then render if they weren't already initialized
                    if (statsManager.renderAllGraphs) {
                        console.log('[HD1-UI] Triggering graph re-render after uncollapse');
                        statsManager.renderAllGraphs();
                    }
                }
            }, 100); // Small delay to ensure display is updated
        }
        
        // Handle debug content with smooth transitions
        if (debugContent) {
            if (collapsed) {
                debugContent.style.maxHeight = '0';
                debugContent.style.opacity = '0';
            } else {
                debugContent.style.maxHeight = '400px';
                debugContent.style.opacity = '1';
            }
        }
        
        // Update collapse icons
        if (collapseIcon) {
            collapseIcon.textContent = collapsed ? '↓' : '↑';
        }
        
        if (statusArrow) {
            statusArrow.textContent = collapsed ? '↑' : '↓';
        }
        
        // Save state if requested
        if (save) {
            localStorage.setItem('hd1_console_collapsed', collapsed.toString());
            
            // Log state change
            this.dom.addDebugEntry('UI_COLLAPSE', {
                collapsed: collapsed,
                saved: save
            });
        }
        
        console.log(`[HD1-UI] Console ${collapsed ? 'collapsed' : 'expanded'}`);
    }


    /**
     * Load version information
     */
    async loadVersionInfo() {
        try {
            // Load API version
            const versionResponse = await window.hd1API.getVersion();
            this.apiVersion = versionResponse.api_version || 'unknown';
            
            // Use JS version from API response (more reliable than URL parsing)
            this.jsVersion = versionResponse.js_version ? versionResponse.js_version.substring(0, 8) : this.extractJSVersion();
            
            // Update displays
            this.dom.setText('api-version', `v${this.apiVersion}`);
            this.dom.setText('js-version', this.jsVersion);
            
            console.log(`[HD1-UI] Version info loaded: API ${this.apiVersion}, JS ${this.jsVersion}`);
            
        } catch (error) {
            console.error('[HD1-UI] Failed to load version info:', error);
            this.dom.setText('api-version', 'v?.?.?');
            this.dom.setText('js-version', 'unknown');
        }
    }

    /**
     * Extract JS version from script URL dynamically
     */
    extractJSVersion() {
        // Try to extract from script URL parameter
        const scripts = document.getElementsByTagName('script');
        for (const script of scripts) {
            if (script.src && script.src.includes('hd1-console.js')) {
                const urlParts = script.src.split('?');
                if (urlParts.length > 1) {
                    const urlParams = new URLSearchParams(urlParts[1]);
                    const version = urlParams.get('v');
                    if (version && version !== '${JS_VERSION}') {
                        // Return first 8 characters of the actual version hash
                        return version.substring(0, 8);
                    }
                }
            }
        }
        
        // Fallback: generate from script modification time or current timestamp
        const scriptElement = document.querySelector('script[src*="hd1-console.js"]');
        if (scriptElement && scriptElement.src) {
            // Create hash from URL without parameters
            const baseUrl = scriptElement.src.split('?')[0];
            const hash = this.simpleHash(baseUrl + Date.now());
            return hash.substring(0, 8);
        }
        
        // Final fallback
        return 'loading';
    }

    /**
     * Simple hash function for version generation
     */
    simpleHash(str) {
        let hash = 0;
        for (let i = 0; i < str.length; i++) {
            const char = str.charCodeAt(i);
            hash = ((hash << 5) - hash) + char;
            hash = hash & hash; // Convert to 32-bit integer
        }
        return Math.abs(hash).toString(16);
    }

    /**
     * Start refresh timer
     */
    startRefreshTimer() {
        this.refreshStartTime = Date.now();
        
        this.timerInterval = setInterval(() => {
            this.updateRefreshTimer();
        }, 1000);
    }

    /**
     * Update refresh timer display
     */
    updateRefreshTimer() {
        const elapsed = Math.floor((Date.now() - this.refreshStartTime) / 1000);
        const minutes = Math.floor(elapsed / 60);
        const seconds = elapsed % 60;
        
        const timeStr = `${minutes}:${seconds.toString().padStart(2, '0')}`;
        this.dom.setText('refresh-timer', timeStr);
    }

    /**
     * Trigger system rebootstrap
     */
    triggerRebootstrap() {
        console.log('[HD1-UI] User triggered rebootstrap');
        
        // Get notification manager for reboot countdown
        const notificationManager = window.hd1ConsoleManager?.getModule('notification');
        if (notificationManager && notificationManager.startRebootCountdown) {
            notificationManager.startRebootCountdown();
        } else {
            // Fallback to immediate reboot if notification manager not available
            console.warn('[HD1-UI] Notification manager not available, immediate reboot');
            this.dom.setText('refresh-timer', 'Reboot...');
            localStorage.removeItem('hd1_session_id');
            localStorage.removeItem('hd1_current_channel');
            setTimeout(() => {
                window.location.reload();
            }, 500);
        }
    }

    /**
     * Handle window resize
     */
    handleResize() {
        // Update any responsive elements
        console.log('[HD1-UI] Window resized');
    }

    /**
     * Setup responsive behavior for different screen sizes
     */
    setupResponsiveBehavior() {
        // Auto-collapse on small screens
        if (window.innerWidth < 768) {
            this.setCollapseState(true, false);
        }
    }


    /**
     * Show status notification
     */
    showStatusNotification(message, duration = 3000) {
        // This would integrate with a notification panel if it exists
        console.log(`[HD1-UI] Status: ${message}`);
        
        // For now, use debug log
        this.dom.addDebugEntry('UI_STATUS', { message, duration });
    }


    /**
     * Get current UI state
     */
    getUIState() {
        return {
            collapsed: this.isCollapsed,
            apiVersion: this.apiVersion,
            jsVersion: this.jsVersion,
            refreshTime: Math.floor((Date.now() - this.refreshStartTime) / 1000)
        };
    }

    /**
     * Update PlayCanvas controls visibility
     */
    updatePlayCanvasControlsVisibility(show) {
        const controls = this.dom.get('playcanvas-controls');
        if (controls) {
            controls.style.display = show ? 'block' : 'none';
            console.log(`[HD1-UI] PlayCanvas controls ${show ? 'shown' : 'hidden'}`);
        }
    }

    /**
     * Show PlayCanvas controls when engine is ready
     */
    showPlayCanvasControls() {
        this.updatePlayCanvasControlsVisibility(true);
    }

    /**
     * Apply theme (if theme system is implemented)
     */
    applyTheme(theme = 'default') {
        document.body.className = `hd1-theme-${theme}`;
        console.log(`[HD1-UI] Applied theme: ${theme}`);
    }

    /**
     * Save user preference
     */
    savePreference(key, value) {
        try {
            localStorage.setItem(`hd1_pref_${key}`, JSON.stringify(value));
            console.log(`[HD1-UI] Saved preference: ${key} = ${value}`);
        } catch (error) {
            console.error(`[HD1-UI] Failed to save preference ${key}:`, error);
        }
    }

    /**
     * Load user preference
     */
    loadPreference(key, defaultValue = null) {
        try {
            const value = localStorage.getItem(`hd1_pref_${key}`);
            return value ? JSON.parse(value) : defaultValue;
        } catch (error) {
            console.error(`[HD1-UI] Failed to load preference ${key}:`, error);
            return defaultValue;
        }
    }

    /**
     * Cleanup UI state manager
     */
    cleanup() {
        if (this.timerInterval) {
            clearInterval(this.timerInterval);
        }
        
        this.ready = false;
    }
}

// Export for use in console manager
window.HD1UIStateManager = HD1UIStateManager;