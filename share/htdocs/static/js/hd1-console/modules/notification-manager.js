/**
 * HD1 Notification Manager - Status Notifications and Visual Feedback
 * 
 * Manages status notifications, glow effects, and user feedback animations.
 * Provides centralized notification system for all console modules.
 */
class HD1NotificationManager {
    constructor(domManager) {
        this.dom = domManager;
        this.ready = false;
        
        // Notification state
        this.activeNotifications = new Map();
        this.notificationQueue = [];
        this.maxNotifications = 3;
        
        // Glow effects
        this.activeGlows = new Map();
        
        // Reboot countdown state
        this.rebootCountdown = null;
        this.rebootInterval = null;
        this.rebootCancelled = false;
        
        // Elements
        this.notificationPanel = null;
        this.notificationContainer = null;
        
        // Auto-hide timeouts
        this.notificationTimeouts = new Map();
    }

    /**
     * Initialize notification manager
     */
    async initialize() {
        console.log('[HD1-Notification] Initializing Notification Manager');
        
        try {
            // Create notification panel if it doesn't exist
            this.createNotificationPanel();
            
            // Setup CSS styles
            this.createNotificationStyles();
            
            this.ready = true;
            console.log('[HD1-Notification] Notification Manager ready');
        } catch (error) {
            console.error('[HD1-Notification] Initialization failed:', error);
        }
    }

    /**
     * Create notification panel
     */
    createNotificationPanel() {
        // Check if panel already exists
        this.notificationPanel = document.getElementById('status-notification-panel');
        
        if (!this.notificationPanel) {
            // Create new panel
            this.notificationPanel = document.createElement('div');
            this.notificationPanel.id = 'status-notification-panel';
            this.notificationPanel.className = 'status-panel-hidden';
            
            this.notificationContainer = document.createElement('div');
            this.notificationContainer.className = 'status-panel-content';
            this.notificationContainer.innerHTML = '<div id="status-message-text"></div>';
            
            this.notificationPanel.appendChild(this.notificationContainer);
            document.body.appendChild(this.notificationPanel);
        } else {
            // Use existing panel
            this.notificationContainer = this.notificationPanel.querySelector('.status-panel-content');
        }
        
        console.log('[HD1-Notification] Notification panel ready');
    }

    /**
     * Create notification and glow effect styles
     */
    createNotificationStyles() {
        if (document.getElementById('hd1-notification-styles')) {
            return; // Styles already exist
        }
        
        const style = document.createElement('style');
        style.id = 'hd1-notification-styles';
        style.textContent = `
            /* Notification Panel Styles */
            .status-panel-hidden {
                position: fixed;
                bottom: -100px;
                right: 20px;
                background: rgba(0, 0, 0, 0.9);
                color: #00ffff;
                border: 1px solid #00ffff;
                border-radius: 4px;
                padding: 12px 16px;
                font-family: monospace;
                font-size: 12px;
                z-index: 9999;
                min-width: 200px;
                max-width: 400px;
                transition: bottom 0.3s ease-in-out;
                box-shadow: 0 4px 12px rgba(0, 255, 255, 0.2);
            }
            
            .status-panel-visible {
                bottom: 20px;
            }
            
            /* Glow Effect Classes */
            .hd1-glow-cyan {
                animation: glowCyan 0.8s ease-in-out;
                transition: all 0.3s ease;
            }
            
            .hd1-glow-red {
                animation: glowRed 0.8s ease-in-out;
                transition: all 0.3s ease;
            }
            
            .hd1-glow-green {
                animation: glowGreen 0.8s ease-in-out;
                transition: all 0.3s ease;
            }
            
            .hd1-glow-yellow {
                animation: glowYellow 0.8s ease-in-out;
                transition: all 0.3s ease;
            }
            
            /* Glow Animations */
            @keyframes glowCyan {
                0% { box-shadow: 0 0 0 rgba(0, 255, 255, 0); }
                50% { box-shadow: 0 0 15px rgba(0, 255, 255, 0.8), inset 0 0 15px rgba(0, 255, 255, 0.2); }
                100% { box-shadow: 0 0 0 rgba(0, 255, 255, 0); }
            }
            
            @keyframes glowRed {
                0% { box-shadow: 0 0 0 rgba(255, 0, 0, 0); }
                50% { box-shadow: 0 0 15px rgba(255, 0, 0, 0.8), inset 0 0 15px rgba(255, 0, 0, 0.2); }
                100% { box-shadow: 0 0 0 rgba(255, 0, 0, 0); }
            }
            
            @keyframes glowGreen {
                0% { box-shadow: 0 0 0 rgba(0, 255, 0, 0); }
                50% { box-shadow: 0 0 15px rgba(0, 255, 0, 0.8), inset 0 0 15px rgba(0, 255, 0, 0.2); }
                100% { box-shadow: 0 0 0 rgba(0, 255, 0, 0); }
            }
            
            @keyframes glowYellow {
                0% { box-shadow: 0 0 0 rgba(255, 255, 0, 0); }
                50% { box-shadow: 0 0 15px rgba(255, 255, 0, 0.8), inset 0 0 15px rgba(255, 255, 0, 0.2); }
                100% { box-shadow: 0 0 0 rgba(255, 255, 0, 0); }
            }
            
            /* Multiple notification support */
            .notification-item {
                margin-bottom: 8px;
                padding: 8px;
                border-left: 3px solid #00ffff;
                background: rgba(0, 255, 255, 0.05);
            }
            
            .notification-item:last-child {
                margin-bottom: 0;
            }
            
            .notification-item.success {
                border-left-color: #00ff00;
                background: rgba(0, 255, 0, 0.05);
            }
            
            .notification-item.warning {
                border-left-color: #ffff00;
                background: rgba(255, 255, 0, 0.05);
            }
            
            .notification-item.error {
                border-left-color: #ff0000;
                background: rgba(255, 0, 0, 0.05);
            }
        `;
        
        document.head.appendChild(style);
        console.log('[HD1-Notification] Notification styles created');
    }

    /**
     * Show notification
     */
    showNotification(message, type = 'info', duration = 4000) {
        const id = Date.now().toString();
        
        // Add to active notifications
        this.activeNotifications.set(id, {
            message,
            type,
            duration,
            timestamp: Date.now()
        });
        
        // Update display
        this.updateNotificationDisplay();
        
        // Auto-hide
        if (duration > 0) {
            const timeout = setTimeout(() => {
                this.hideNotification(id);
            }, duration);
            
            this.notificationTimeouts.set(id, timeout);
        }
        
        console.log(`[HD1-Notification] Notification shown: ${message}`);
        return id;
    }

    /**
     * Hide specific notification
     */
    hideNotification(id) {
        if (this.activeNotifications.has(id)) {
            this.activeNotifications.delete(id);
            
            // Clear timeout
            if (this.notificationTimeouts.has(id)) {
                clearTimeout(this.notificationTimeouts.get(id));
                this.notificationTimeouts.delete(id);
            }
            
            // Update display
            this.updateNotificationDisplay();
        }
    }

    /**
     * Clear all notifications
     */
    clearAllNotifications() {
        // Clear all timeouts
        this.notificationTimeouts.forEach(timeout => clearTimeout(timeout));
        this.notificationTimeouts.clear();
        
        // Clear notifications
        this.activeNotifications.clear();
        
        // Update display
        this.updateNotificationDisplay();
        
        console.log('[HD1-Notification] All notifications cleared');
    }

    /**
     * Update notification display
     */
    updateNotificationDisplay() {
        if (!this.notificationPanel || !this.notificationContainer) return;
        
        const messageElement = this.notificationContainer.querySelector('#status-message-text');
        if (!messageElement) return;
        
        if (this.activeNotifications.size === 0) {
            // Hide panel
            this.notificationPanel.className = 'status-panel-hidden';
            messageElement.innerHTML = '';
            return;
        }
        
        // Show panel
        this.notificationPanel.className = 'status-panel-visible';
        
        // Build notification content
        let content = '';
        const notifications = Array.from(this.activeNotifications.values());
        
        if (notifications.length === 1) {
            // Single notification
            content = notifications[0].message;
        } else {
            // Multiple notifications
            content = notifications.map(notif => 
                `<div class="notification-item ${notif.type}">${notif.message}</div>`
            ).join('');
        }
        
        messageElement.innerHTML = content;
    }

    /**
     * Apply glow effect to element
     */
    addGlowEffect(elementId, color = 'cyan', duration = 800) {
        const element = this.dom.get(elementId);
        if (!element) {
            console.warn(`[HD1-Notification] Element not found for glow: ${elementId}`);
            return;
        }
        
        const glowClass = `hd1-glow-${color}`;
        
        // Remove any existing glow
        this.removeGlowEffect(elementId);
        
        // Add new glow
        element.classList.add(glowClass);
        
        // Store active glow
        this.activeGlows.set(elementId, glowClass);
        
        // Remove after animation
        setTimeout(() => {
            this.removeGlowEffect(elementId);
        }, duration);
        
        console.log(`[HD1-Notification] Glow effect applied: ${elementId} (${color})`);
    }

    /**
     * Remove glow effect from element
     */
    removeGlowEffect(elementId) {
        const element = this.dom.get(elementId);
        if (!element) return;
        
        if (this.activeGlows.has(elementId)) {
            const glowClass = this.activeGlows.get(elementId);
            element.classList.remove(glowClass);
            this.activeGlows.delete(elementId);
        }
        
        // Remove all possible glow classes
        ['hd1-glow-cyan', 'hd1-glow-red', 'hd1-glow-green', 'hd1-glow-yellow'].forEach(cls => {
            element.classList.remove(cls);
        });
    }

    /**
     * Show success notification
     */
    showSuccess(message, duration = 3000) {
        return this.showNotification(message, 'success', duration);
    }

    /**
     * Show warning notification
     */
    showWarning(message, duration = 4000) {
        return this.showNotification(message, 'warning', duration);
    }

    /**
     * Show error notification
     */
    showError(message, duration = 5000) {
        return this.showNotification(message, 'error', duration);
    }

    /**
     * Show info notification
     */
    showInfo(message, duration = 4000) {
        return this.showNotification(message, 'info', duration);
    }

    /**
     * Show session copied notification with glow
     */
    showSessionCopied(sessionId) {
        this.addGlowEffect('session-id-tag-status', 'cyan');
        return this.showSuccess(`Session ID copied: ${sessionId}`, 2000);
    }

    /**
     * Show reboot notification with glow
     */
    showRebootNotification() {
        this.addGlowEffect('refresh-timer', 'red');
        return this.showWarning('System rebootstrap initiated...', 3000);
    }

    /**
     * Start reboot countdown with cancellation option
     */
    startRebootCountdown() {
        this.rebootCancelled = false;
        this.rebootCountdown = 5; // 5 seconds
        
        // Add red glow to refresh timer
        this.addGlowEffect('refresh-timer', 'red', 6000);
        
        // Show initial notification
        const notificationId = this.showWarning(`System rebooting in ${this.rebootCountdown} seconds... (ESC to cancel)`, 6000);
        
        // Start countdown interval
        this.rebootInterval = setInterval(() => {
            if (this.rebootCancelled) {
                this.cancelRebootCountdown();
                return;
            }
            
            this.rebootCountdown--;
            
            if (this.rebootCountdown <= 0) {
                // Execute reboot
                this.executeReboot();
            } else {
                // Update notification
                this.updateNotification(notificationId, `System rebooting in ${this.rebootCountdown} seconds... (ESC to cancel)`);
            }
        }, 1000);
        
        // Setup ESC key listener
        this.setupRebootCancellation();
        
        console.log('[HD1-Notification] Reboot countdown started');
        return notificationId;
    }

    /**
     * Cancel reboot countdown
     */
    cancelRebootCountdown() {
        this.rebootCancelled = true;
        
        if (this.rebootInterval) {
            clearInterval(this.rebootInterval);
            this.rebootInterval = null;
        }
        
        // Remove ESC listener
        this.removeRebootCancellation();
        
        // Remove red glow
        this.removeGlowEffect('refresh-timer');
        
        // Clear reboot notifications
        this.clearAllNotifications();
        
        // Show cancellation message
        this.showInfo('System reboot cancelled', 2000);
        
        console.log('[HD1-Notification] Reboot countdown cancelled');
    }

    /**
     * Execute the actual reboot
     */
    executeReboot() {
        if (this.rebootInterval) {
            clearInterval(this.rebootInterval);
            this.rebootInterval = null;
        }
        
        // Remove ESC listener
        this.removeRebootCancellation();
        
        // Show final message
        this.showWarning('Rebooting system...', 1000);
        
        // Clear localStorage
        localStorage.removeItem('hd1_session_id');
        localStorage.removeItem('hd1_current_channel');
        
        // Reload page
        setTimeout(() => {
            window.location.reload();
        }, 500);
        
        console.log('[HD1-Notification] Executing system reboot');
    }

    /**
     * Setup ESC key cancellation
     */
    setupRebootCancellation() {
        this.escKeyHandler = (event) => {
            if (event.key === 'Escape') {
                event.preventDefault();
                this.cancelRebootCountdown();
            }
        };
        
        document.addEventListener('keydown', this.escKeyHandler);
    }

    /**
     * Remove ESC key cancellation
     */
    removeRebootCancellation() {
        if (this.escKeyHandler) {
            document.removeEventListener('keydown', this.escKeyHandler);
            this.escKeyHandler = null;
        }
    }

    /**
     * Update existing notification message
     */
    updateNotification(notificationId, newMessage) {
        const notification = this.activeNotifications.get(notificationId);
        if (notification) {
            notification.message = newMessage;
            this.updateNotificationDisplay();
        }
    }

    /**
     * Handle events from other modules
     */
    handleEvent(eventType, data) {
        switch (eventType) {
            case 'session_copied':
                this.showSessionCopied(data.sessionId);
                break;
            case 'reboot_requested':
                this.showRebootNotification();
                break;
            case 'websocket_connected':
                this.showSuccess('WebSocket connected');
                break;
            case 'websocket_disconnected':
                this.showWarning('WebSocket disconnected');
                break;
            case 'channel_changed':
                this.showInfo(`Switched to channel: ${data.channelId}`);
                break;
        }
    }

    /**
     * Get notification statistics
     */
    getStats() {
        return {
            active: this.activeNotifications.size,
            activeGlows: this.activeGlows.size,
            ready: this.ready
        };
    }

    /**
     * Cleanup notification manager
     */
    cleanup() {
        // Clear all notifications and timeouts
        this.clearAllNotifications();
        
        // Remove all glow effects
        this.activeGlows.forEach((glowClass, elementId) => {
            this.removeGlowEffect(elementId);
        });
        
        // Remove notification panel
        if (this.notificationPanel && this.notificationPanel.parentNode) {
            this.notificationPanel.parentNode.removeChild(this.notificationPanel);
        }
        
        // Remove styles
        const styles = document.getElementById('hd1-notification-styles');
        if (styles && styles.parentNode) {
            styles.parentNode.removeChild(styles);
        }
        
        this.ready = false;
        console.log('[HD1-Notification] Notification Manager cleaned up');
    }
}

// Export for use in console manager
window.HD1NotificationManager = HD1NotificationManager;