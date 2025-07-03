/**
 * HD1 Input Manager - Keyboard and Mouse Input Handling
 * 
 * Manages all input events, keyboard shortcuts, mouse controls, and pointer lock.
 * Provides mouse look mode with ESC notification and input state management.
 */
class HD1InputManager {
    constructor(domManager) {
        this.dom = domManager;
        this.ready = false;
        
        // Input state
        this.isMouseLookActive = false;
        this.keys = new Set();
        this.mouseButtons = new Set();
        
        // Mouse look settings
        this.mouseSensitivity = 0.002;
        this.invertY = false;
        
        // Notification state
        this.escNotificationVisible = false;
        this.escNotificationTimeout = null;
        
        // Event handlers
        this.keyDownHandler = null;
        this.keyUpHandler = null;
        this.mouseDownHandler = null;
        this.mouseUpHandler = null;
        this.mouseMoveHandler = null;
        this.pointerLockChangeHandler = null;
        this.contextMenuHandler = null;
        
        // Canvas reference
        this.canvas = null;
    }

    /**
     * Initialize input manager
     */
    async initialize() {
        console.log('[HD1-Input] Initializing Input Manager');
        
        try {
            // Get canvas element
            this.canvas = this.dom.get('hd1-playcanvas-canvas');
            
            // Setup input handlers
            this.setupKeyboardHandlers();
            this.setupMouseHandlers();
            this.setupPointerLockHandlers();
            
            // Create ESC notification element
            this.createEscNotification();
            
            this.ready = true;
            console.log('[HD1-Input] Input Manager ready');
        } catch (error) {
            console.error('[HD1-Input] Initialization failed:', error);
        }
    }

    /**
     * Setup keyboard event handlers
     */
    setupKeyboardHandlers() {
        this.keyDownHandler = (event) => {
            this.handleKeyDown(event);
        };
        
        this.keyUpHandler = (event) => {
            this.handleKeyUp(event);
        };
        
        document.addEventListener('keydown', this.keyDownHandler);
        document.addEventListener('keyup', this.keyUpHandler);
        
        console.log('[HD1-Input] Keyboard handlers setup');
    }

    /**
     * Setup mouse event handlers
     */
    setupMouseHandlers() {
        if (!this.canvas) {
            console.warn('[HD1-Input] Canvas not found, skipping mouse handlers');
            return;
        }
        
        this.mouseDownHandler = (event) => {
            this.handleMouseDown(event);
        };
        
        this.mouseUpHandler = (event) => {
            this.handleMouseUp(event);
        };
        
        this.mouseMoveHandler = (event) => {
            this.handleMouseMove(event);
        };
        
        this.contextMenuHandler = (event) => {
            // Prevent context menu when in mouse look mode
            if (this.isMouseLookActive) {
                event.preventDefault();
            }
        };
        
        this.canvas.addEventListener('mousedown', this.mouseDownHandler);
        this.canvas.addEventListener('mouseup', this.mouseUpHandler);
        document.addEventListener('mousemove', this.mouseMoveHandler);
        this.canvas.addEventListener('contextmenu', this.contextMenuHandler);
        
        console.log('[HD1-Input] Mouse handlers setup');
    }

    /**
     * Setup pointer lock handlers
     */
    setupPointerLockHandlers() {
        this.pointerLockChangeHandler = () => {
            this.handlePointerLockChange();
        };
        
        document.addEventListener('pointerlockchange', this.pointerLockChangeHandler);
        document.addEventListener('pointerlockerror', () => {
            console.error('[HD1-Input] Pointer lock error');
            this.exitMouseLook();
        });
        
        console.log('[HD1-Input] Pointer lock handlers setup');
    }

    /**
     * Handle key down events
     */
    handleKeyDown(event) {
        this.keys.add(event.code);
        
        // Handle ESC key
        if (event.code === 'Escape') {
            if (this.isMouseLookActive) {
                this.exitMouseLook();
                event.preventDefault();
                return;
            }
        }
        
        // Handle other keyboard shortcuts
        this.handleKeyboardShortcuts(event);
        
        // Log debug info
        this.dom.addDebugEntry('KEY_DOWN', { code: event.code, key: event.key });
    }

    /**
     * Handle key up events
     */
    handleKeyUp(event) {
        this.keys.delete(event.code);
        
        this.dom.addDebugEntry('KEY_UP', { code: event.code, key: event.key });
    }

    /**
     * Handle keyboard shortcuts
     */
    handleKeyboardShortcuts(event) {
        // Console toggle (backtick/grave accent)
        if (event.code === 'Backquote' && !this.isMouseLookActive) {
            this.toggleConsole();
            event.preventDefault();
            return;
        }
        
        // F11 for fullscreen
        if (event.code === 'F11') {
            this.toggleFullscreen();
            event.preventDefault();
            return;
        }
        
        // WASD movement (when in mouse look mode)
        if (this.isMouseLookActive) {
            this.handleMovementKeys(event);
        }
    }

    /**
     * Handle movement keys in mouse look mode
     */
    handleMovementKeys(event) {
        const movementKeys = ['KeyW', 'KeyA', 'KeyS', 'KeyD', 'Space', 'ShiftLeft', 'ShiftRight'];
        
        if (movementKeys.includes(event.code)) {
            // Send movement to PlayCanvas if available
            if (window.hd1GameEngine) {
                this.sendMovementToPlayCanvas(event.code, true);
            }
        }
    }

    /**
     * Handle mouse down events
     */
    handleMouseDown(event) {
        this.mouseButtons.add(event.button);
        
        // Left click on canvas - enter mouse look mode
        if (event.button === 0 && this.canvas && event.target === this.canvas) {
            this.enterMouseLook();
            event.preventDefault();
        }
        
        this.dom.addDebugEntry('MOUSE_DOWN', { button: event.button });
    }

    /**
     * Handle mouse up events
     */
    handleMouseUp(event) {
        this.mouseButtons.delete(event.button);
        
        this.dom.addDebugEntry('MOUSE_UP', { button: event.button });
    }

    /**
     * Handle mouse move events
     */
    handleMouseMove(event) {
        if (this.isMouseLookActive) {
            this.handleMouseLookMovement(event);
        }
    }

    /**
     * Handle mouse movement in mouse look mode
     */
    handleMouseLookMovement(event) {
        if (!this.isMouseLookActive) return;
        
        const deltaX = event.movementX * this.mouseSensitivity;
        const deltaY = event.movementY * this.mouseSensitivity * (this.invertY ? -1 : 1);
        
        // Send mouse look data to PlayCanvas
        if (window.hd1GameEngine) {
            this.sendMouseLookToPlayCanvas(deltaX, deltaY);
        }
    }

    /**
     * Enter mouse look mode
     */
    enterMouseLook() {
        if (!this.canvas) {
            console.warn('[HD1-Input] Cannot enter mouse look - canvas not available');
            return;
        }
        
        // Request pointer lock
        this.canvas.requestPointerLock();
    }

    /**
     * Exit mouse look mode
     */
    exitMouseLook() {
        if (document.pointerLockElement) {
            document.exitPointerLock();
        }
        
        this.isMouseLookActive = false;
        this.hideEscNotification();
        
        console.log('[HD1-Input] Exited mouse look mode');
    }

    /**
     * Handle pointer lock state changes
     */
    handlePointerLockChange() {
        if (document.pointerLockElement === this.canvas) {
            this.isMouseLookActive = true;
            this.showEscNotification();
            console.log('[HD1-Input] Entered mouse look mode');
        } else {
            this.isMouseLookActive = false;
            this.hideEscNotification();
            console.log('[HD1-Input] Exited mouse look mode');
        }
        
        // Notify other modules
        this.notifyMouseLookStateChange();
    }

    /**
     * Create ESC notification element
     */
    createEscNotification() {
        // Check if notification already exists
        if (document.getElementById('hd1-esc-notification')) {
            return;
        }
        
        const notification = document.createElement('div');
        notification.id = 'hd1-esc-notification';
        notification.style.cssText = `
            position: fixed;
            top: 20px;
            left: 50%;
            transform: translateX(-50%);
            background: rgba(0, 0, 0, 0.8);
            color: #00ffff;
            padding: 10px 20px;
            border: 1px solid #00ffff;
            border-radius: 4px;
            font-family: monospace;
            font-size: 14px;
            z-index: 10000;
            display: none;
            animation: fadeIn 0.3s ease-in;
        `;
        notification.innerHTML = 'Press <strong>ESC</strong> to exit mouse look mode';
        
        // Add CSS animation
        if (!document.getElementById('hd1-input-styles')) {
            const style = document.createElement('style');
            style.id = 'hd1-input-styles';
            style.textContent = `
                @keyframes fadeIn {
                    from { opacity: 0; transform: translateX(-50%) translateY(-10px); }
                    to { opacity: 1; transform: translateX(-50%) translateY(0); }
                }
                @keyframes fadeOut {
                    from { opacity: 1; transform: translateX(-50%) translateY(0); }
                    to { opacity: 0; transform: translateX(-50%) translateY(-10px); }
                }
            `;
            document.head.appendChild(style);
        }
        
        document.body.appendChild(notification);
        console.log('[HD1-Input] ESC notification created');
    }

    /**
     * Show ESC notification
     */
    showEscNotification() {
        const notification = document.getElementById('hd1-esc-notification');
        if (notification) {
            notification.style.display = 'block';
            this.escNotificationVisible = true;
            
            // Auto-hide after 3 seconds
            if (this.escNotificationTimeout) {
                clearTimeout(this.escNotificationTimeout);
            }
            
            this.escNotificationTimeout = setTimeout(() => {
                this.hideEscNotification();
            }, 3000);
        }
    }

    /**
     * Hide ESC notification
     */
    hideEscNotification() {
        const notification = document.getElementById('hd1-esc-notification');
        if (notification) {
            notification.style.animation = 'fadeOut 0.3s ease-out';
            setTimeout(() => {
                notification.style.display = 'none';
                notification.style.animation = 'fadeIn 0.3s ease-in';
            }, 300);
            
            this.escNotificationVisible = false;
            
            if (this.escNotificationTimeout) {
                clearTimeout(this.escNotificationTimeout);
                this.escNotificationTimeout = null;
            }
        }
    }

    /**
     * Toggle console visibility
     */
    toggleConsole() {
        const uiManager = window.hd1ConsoleManager?.getModule('ui-state');
        if (uiManager && uiManager.toggleCollapse) {
            uiManager.toggleCollapse();
        }
    }

    /**
     * Toggle fullscreen mode
     */
    toggleFullscreen() {
        if (!document.fullscreenElement) {
            document.documentElement.requestFullscreen().catch(err => {
                console.error('[HD1-Input] Fullscreen request failed:', err);
            });
        } else {
            document.exitFullscreen().catch(err => {
                console.error('[HD1-Input] Exit fullscreen failed:', err);
            });
        }
    }

    /**
     * Send movement input to PlayCanvas
     */
    sendMovementToPlayCanvas(keyCode, pressed) {
        if (!window.hd1GameEngine) return;
        
        // This would integrate with PlayCanvas character controller
        // For now, just log the input
        console.log(`[HD1-Input] Movement: ${keyCode} ${pressed ? 'pressed' : 'released'}`);
    }

    /**
     * Send mouse look input to PlayCanvas
     */
    sendMouseLookToPlayCanvas(deltaX, deltaY) {
        if (!window.hd1GameEngine) return;
        
        // This would integrate with PlayCanvas camera controller
        // For now, just log the input
        console.log(`[HD1-Input] Mouse look: deltaX=${deltaX.toFixed(3)}, deltaY=${deltaY.toFixed(3)}`);
    }

    /**
     * Notify other modules of mouse look state change
     */
    notifyMouseLookStateChange() {
        if (window.hd1ConsoleManager) {
            window.hd1ConsoleManager.notifyModules('mouse_look_changed', {
                active: this.isMouseLookActive
            });
        }
    }

    /**
     * Check if key is currently pressed
     */
    isKeyPressed(keyCode) {
        return this.keys.has(keyCode);
    }

    /**
     * Check if mouse button is currently pressed
     */
    isMouseButtonPressed(button) {
        return this.mouseButtons.has(button);
    }

    /**
     * Get current input state
     */
    getInputState() {
        return {
            mouseLookActive: this.isMouseLookActive,
            pressedKeys: Array.from(this.keys),
            pressedMouseButtons: Array.from(this.mouseButtons),
            escNotificationVisible: this.escNotificationVisible
        };
    }

    /**
     * Set mouse sensitivity
     */
    setMouseSensitivity(sensitivity) {
        this.mouseSensitivity = Math.max(0.0001, Math.min(0.01, sensitivity));
        console.log(`[HD1-Input] Mouse sensitivity set to ${this.mouseSensitivity}`);
    }

    /**
     * Toggle Y-axis inversion
     */
    toggleInvertY() {
        this.invertY = !this.invertY;
        console.log(`[HD1-Input] Y-axis inversion ${this.invertY ? 'enabled' : 'disabled'}`);
    }

    /**
     * Handle events from other modules
     */
    handleEvent(eventType, data) {
        switch (eventType) {
            case 'playcanvas_ready':
                console.log('[HD1-Input] PlayCanvas ready, input integration active');
                break;
        }
    }

    /**
     * Cleanup input manager
     */
    cleanup() {
        // Remove event listeners
        if (this.keyDownHandler) {
            document.removeEventListener('keydown', this.keyDownHandler);
        }
        if (this.keyUpHandler) {
            document.removeEventListener('keyup', this.keyUpHandler);
        }
        if (this.mouseDownHandler && this.canvas) {
            this.canvas.removeEventListener('mousedown', this.mouseDownHandler);
        }
        if (this.mouseUpHandler && this.canvas) {
            this.canvas.removeEventListener('mouseup', this.mouseUpHandler);
        }
        if (this.mouseMoveHandler) {
            document.removeEventListener('mousemove', this.mouseMoveHandler);
        }
        if (this.pointerLockChangeHandler) {
            document.removeEventListener('pointerlockchange', this.pointerLockChangeHandler);
        }
        if (this.contextMenuHandler && this.canvas) {
            this.canvas.removeEventListener('contextmenu', this.contextMenuHandler);
        }
        
        // Clear timeouts
        if (this.escNotificationTimeout) {
            clearTimeout(this.escNotificationTimeout);
        }
        
        // Exit mouse look if active
        if (this.isMouseLookActive) {
            this.exitMouseLook();
        }
        
        // Remove notification element
        const notification = document.getElementById('hd1-esc-notification');
        if (notification) {
            notification.remove();
        }
        
        // Clear state
        this.keys.clear();
        this.mouseButtons.clear();
        this.ready = false;
        
        console.log('[HD1-Input] Input Manager cleaned up');
    }
}

// Export for use in console manager
window.HD1InputManager = HD1InputManager;