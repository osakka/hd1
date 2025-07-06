/**
 * HD1 Session Manager - Session Lifecycle and State Management
 * 
 * Manages HD1 sessions, automatic creation, persistence, and cleanup.
 * Integrates with UI for session display and WebSocket for real-time updates.
 */
class HD1SessionManager {
    constructor(apiClient, domManager) {
        this.api = apiClient;
        this.dom = domManager;
        this.ready = false;
        
        // Session state
        this.currentSession = null;
        this.sessionId = null;
        this.sessionCreated = false;
        
        // Session monitoring
        this.sessionCheckInterval = null;
        this.sessionCheckFrequency = 30000; // 30 seconds
    }

    /**
     * Initialize session manager
     */
    async initialize() {
        console.log('[HD1-Session] Initializing Session Manager');
        
        try {
            // Load existing session or create new one
            await this.loadOrCreateSession();
            
            // Start session monitoring
            this.startSessionMonitoring();
            
            this.ready = true;
            console.log('[HD1-Session] Session Manager ready');
        } catch (error) {
            console.error('[HD1-Session] Initialization failed:', error);
        }
    }

    /**
     * Load existing session from localStorage or create new one
     */
    async loadOrCreateSession() {
        const existingSessionId = localStorage.getItem('hd1_session_id');
        
        if (existingSessionId) {
            // Verify the session still exists on server
            if (await this.validateSession(existingSessionId)) {
                this.sessionId = existingSessionId;
                await this.loadSessionDetails();
                this.updateSessionDisplay();
                
                // CRITICAL: Associate existing session with WebSocket for proper broadcasting
                this.associateWebSocketSession();
                
                console.log(`[HD1-Session] Loaded existing session: ${existingSessionId}`);
                return;
            } else {
                // Session no longer exists, clear it
                localStorage.removeItem('hd1_session_id');
                console.log('[HD1-Session] Existing session invalid, will create new one');
            }
        }
        
        // Create new session
        await this.createNewSession();
    }

    /**
     * Validate that a session still exists on the server
     */
    async validateSession(sessionId) {
        try {
            const response = await this.api.getSession(sessionId);
            // API response doesn't have 'success' field, just check if we got valid session data
            return response && response.id === sessionId && response.status === 'active';
        } catch (error) {
            console.warn('[HD1-Session] Session validation failed:', error);
            return false;
        }
    }

    /**
     * Create a new session
     */
    async createNewSession() {
        try {
            console.log('[HD1-Session] Creating new session...');
            const response = await this.api.createSession();
            
            if (response.success && response.session_id) {
                this.sessionId = response.session_id;
                this.currentSession = response;
                this.sessionCreated = true;
                
                // Persist session
                localStorage.setItem('hd1_session_id', this.sessionId);
                
                // Update UI
                this.updateSessionDisplay();
                
                console.log(`[HD1-Session] New session created: ${this.sessionId}`);
                
                // Notify other modules
                this.notifySessionCreated();
                
                return this.sessionId;
            } else {
                throw new Error('Session creation failed: ' + JSON.stringify(response));
            }
        } catch (error) {
            console.error('[HD1-Session] Session creation failed:', error);
            this.handleSessionError('Failed to create session');
            return null;
        }
    }

    /**
     * Load detailed session information
     */
    async loadSessionDetails() {
        if (!this.sessionId) return;
        
        try {
            const response = await this.api.getSession(this.sessionId);
            if (response.success) {
                this.currentSession = response;
                console.log('[HD1-Session] Session details loaded');
            }
        } catch (error) {
            console.error('[HD1-Session] Failed to load session details:', error);
        }
    }

    /**
     * Update session display in UI
     */
    updateSessionDisplay() {
        if (this.sessionId) {
            this.dom.setText('session-id-tag-status', this.sessionId);
            
            // Show session tag (it's hidden by default)
            const sessionTag = this.dom.get('session-id-tag-status');
            if (sessionTag) {
                sessionTag.style.display = 'inline';
            }
            
            console.log(`[HD1-Session] Session display updated: ${this.sessionId}`);
        } else {
            this.dom.setText('session-id-tag-status', '---');
        }
    }

    /**
     * Handle session updates from WebSocket
     */
    handleSessionUpdate(message) {
        console.log('[HD1-Session] Session update received:', message);
        
        if (message.session_id) {
            this.sessionId = message.session_id;
            localStorage.setItem('hd1_session_id', this.sessionId);
            this.updateSessionDisplay();
        }
        
        // Update session details if provided
        if (message.session) {
            this.currentSession = { ...this.currentSession, ...message.session };
        }
    }

    /**
     * Start session monitoring
     */
    startSessionMonitoring() {
        this.sessionCheckInterval = setInterval(async () => {
            if (this.sessionId) {
                const isValid = await this.validateSession(this.sessionId);
                if (!isValid) {
                    console.warn('[HD1-Session] Session became invalid, creating new one');
                    await this.createNewSession();
                }
            }
        }, this.sessionCheckFrequency);
    }

    /**
     * Stop session monitoring
     */
    stopSessionMonitoring() {
        if (this.sessionCheckInterval) {
            clearInterval(this.sessionCheckInterval);
            this.sessionCheckInterval = null;
        }
    }

    /**
     * Notify other modules of session creation
     */
    notifySessionCreated() {
        // CRITICAL: Associate WebSocket with session for proper broadcasting
        this.associateWebSocketSession();
        
        // Notify via console manager if available
        if (window.hd1ConsoleManager) {
            window.hd1ConsoleManager.notifyModules('session_created', {
                sessionId: this.sessionId,
                session: this.currentSession
            });
        }
    }
    
    /**
     * Re-associate WebSocket session after world switching
     * CRITICAL FIX: Handle avatar control loss during world transitions
     */
    reAssociateAfterWorldSwitch() {
        console.log('[HD1-Session] Re-associating WebSocket after world switch');
        
        // Reset retry counter for fresh attempts
        this.wsAssociationRetries = 0;
        
        // Force immediate re-association
        this.associateWebSocketSession();
        
        // CRITICAL FIX: Additional avatar control recovery
        setTimeout(() => {
            this.ensureAvatarControlRecovery();
        }, 2000); // Give time for avatar creation
    }
    
    /**
     * Ensure avatar control is recovered after world switching
     */
    ensureAvatarControlRecovery() {
        console.log('[HD1-Session] Ensuring avatar control recovery');
        
        // Check if PlayCanvas has loaded avatars but camera isn't bound
        if (window.hd1GameEngine && window.cameraController) {
            const avatarEntities = window.hd1GameEngine.root.children.filter(entity => 
                entity.hd1Tags && entity.hd1Tags.includes('session-avatar')
            );
            
            console.log(`[HD1-Session] Found ${avatarEntities.length} avatar entities after world switch`);
            
            if (avatarEntities.length > 0 && !window.cameraController.boundAvatar) {
                console.log('[HD1-Session] Attempting to restore avatar camera binding');
                
                // Force camera controller to re-detect avatars
                if (window.cameraController.updateAvailableAvatars) {
                    window.cameraController.updateAvailableAvatars();
                }
                
                // Set avatar-driven mode if avatars are available
                if (window.cameraController.setAvatarDrivenMode) {
                    window.cameraController.setAvatarDrivenMode();
                    console.log('[HD1-Session] ✅ Avatar control restored after world switch');
                }
            }
        }
    }
    
    /**
     * Associate the WebSocket connection with this session for proper message broadcasting
     */
    associateWebSocketSession() {
        if (!this.sessionId) {
            console.warn('[HD1-Session] Cannot associate WebSocket: no session ID');
            return;
        }
        
        // Get WebSocket manager and send session association message
        if (window.hd1ConsoleManager) {
            const wsManager = window.hd1ConsoleManager.getModule('websocket');
            if (wsManager) {
                // Check if WebSocket is connected
                if (wsManager.isConnected && wsManager.send) {
                    const associationMessage = {
                        type: 'session_associate',
                        session_id: this.sessionId
                    };
                    
                    if (wsManager.send(associationMessage)) {
                        console.log(`[HD1-Session] ✅ WebSocket associated with session: ${this.sessionId}`);
                        
                        // CRITICAL FIX: Reset association retry counter on successful association
                        this.wsAssociationRetries = 0;
                        
                        // CRITICAL FIX: Track association success for world switching recovery
                        this.lastAssociationTime = Date.now();
                        
                        return;
                    }
                }
                
                // WebSocket not connected or send failed - retry with exponential backoff
                console.warn(`[HD1-Session] ❌ WebSocket not ready for session association: ${this.sessionId}`);
                this.retryWebSocketAssociation();
            }
        }
    }

    /**
     * Retry WebSocket association with exponential backoff
     */
    retryWebSocketAssociation() {
        if (!this.wsAssociationRetries) {
            this.wsAssociationRetries = 0;
        }
        
        this.wsAssociationRetries++;
        const maxRetries = 5;
        const baseDelay = 1000; // 1 second
        
        if (this.wsAssociationRetries > maxRetries) {
            console.error(`[HD1-Session] Failed to associate WebSocket after ${maxRetries} attempts`);
            return;
        }
        
        const delay = baseDelay * Math.pow(2, this.wsAssociationRetries - 1); // Exponential backoff
        console.log(`[HD1-Session] Retrying WebSocket association in ${delay}ms (attempt ${this.wsAssociationRetries}/${maxRetries})`);
        
        setTimeout(() => {
            this.associateWebSocketSession();
        }, delay);
    }

    /**
     * Handle session errors
     */
    handleSessionError(message) {
        console.error(`[HD1-Session] ${message}`);
        
        // Update UI to show error state
        this.dom.setText('session-id-tag-status', 'ERROR');
        
        // Hide session tag on error
        const sessionTag = this.dom.get('session-id-tag-status');
        if (sessionTag) {
            sessionTag.style.display = 'none';
        }
    }

    /**
     * Copy session ID to clipboard
     */
    async copySessionId() {
        if (!this.sessionId) {
            console.warn('[HD1-Session] No session ID to copy');
            return false;
        }

        console.log('[HD1-Session] Attempting to copy session ID:', this.sessionId);

        // Try modern clipboard API first (only available in secure contexts)
        if (navigator.clipboard && navigator.clipboard.writeText && window.isSecureContext) {
            try {
                await navigator.clipboard.writeText(this.sessionId);
                console.log('[HD1-Session] Session ID copied to clipboard via modern API');
                this.showCopyFeedback();
                return true;
            } catch (error) {
                console.warn('[HD1-Session] Modern clipboard API failed:', error);
            }
        } else {
            console.log('[HD1-Session] Modern clipboard API not available, using fallback');
        }

        // Fallback to execCommand method
        try {
            const textArea = document.createElement('textarea');
            textArea.value = this.sessionId;
            textArea.style.position = 'fixed';
            textArea.style.opacity = '0';
            document.body.appendChild(textArea);
            textArea.select();
            textArea.setSelectionRange(0, 99999);
            
            const successful = document.execCommand('copy');
            document.body.removeChild(textArea);
            
            if (successful) {
                console.log('[HD1-Session] Session ID copied to clipboard via execCommand');
                this.showCopyFeedback();
                return true;
            } else {
                throw new Error('execCommand copy failed');
            }
        } catch (error) {
            console.error('[HD1-Session] All copy methods failed:', error);
            this.selectSessionIdText();
            return false;
        }
    }

    /**
     * Show visual feedback for copy action
     */
    showCopyFeedback() {
        // Get notification manager for glow effect
        const notificationManager = window.hd1ConsoleManager?.getModule('notification');
        if (notificationManager && notificationManager.showSessionCopied) {
            notificationManager.showSessionCopied(this.sessionId);
        } else {
            // Fallback to text change
            const sessionTag = this.dom.get('session-id-tag-status');
            if (!sessionTag) return;
            
            const originalText = sessionTag.textContent;
            sessionTag.textContent = 'Copied!';
            
            setTimeout(() => {
                sessionTag.textContent = originalText;
            }, 1000);
        }
    }

    /**
     * Select session ID text as fallback for copy
     */
    selectSessionIdText() {
        const sessionTag = this.dom.get('session-id-tag-status');
        if (!sessionTag) return;
        
        const range = document.createRange();
        range.selectNode(sessionTag);
        window.getSelection().removeAllRanges();
        window.getSelection().addRange(range);
    }

    /**
     * Clear current session
     */
    clearSession() {
        this.sessionId = null;
        this.currentSession = null;
        this.sessionCreated = false;
        
        localStorage.removeItem('hd1_session_id');
        this.updateSessionDisplay();
        
        console.log('[HD1-Session] Session cleared');
    }

    /**
     * Refresh session (create new one)
     */
    async refreshSession() {
        console.log('[HD1-Session] Refreshing session...');
        this.clearSession();
        await this.createNewSession();
    }

    /**
     * Get current session information
     */
    getCurrentSession() {
        return {
            sessionId: this.sessionId,
            session: this.currentSession,
            created: this.sessionCreated
        };
    }

    /**
     * Get session ID
     */
    getSessionId() {
        return this.sessionId;
    }

    /**
     * Check if session is ready
     */
    isSessionReady() {
        return this.sessionId !== null;
    }

    /**
     * Handle events from other modules
     */
    handleEvent(eventType, data) {
        switch (eventType) {
            case 'websocket_session_update':
                this.handleSessionUpdate(data);
                break;
            case 'rebootstrap_requested':
                this.refreshSession();
                break;
            case 'world_switched':
                // CRITICAL FIX: Handle world switching avatar control recovery
                console.log('[HD1-Session] World switch detected, re-associating WebSocket');
                this.reAssociateAfterWorldSwitch();
                break;
        }
    }

    /**
     * Cleanup session manager
     */
    cleanup() {
        this.stopSessionMonitoring();
        this.ready = false;
    }
}

// Export for use in console manager
window.HD1SessionManager = HD1SessionManager;