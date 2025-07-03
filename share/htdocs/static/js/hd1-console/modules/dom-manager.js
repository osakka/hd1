/**
 * HD1 DOM Manager - Bulletproof DOM Operations with Null Safety
 * 
 * Handles all DOM interactions with comprehensive error boundaries.
 * No operation can break if elements don't exist.
 */
class HD1DOMManager {
    constructor() {
        this.elements = new Map();
        this.observers = new Map();
        this.ready = false;
    }

    /**
     * Initialize DOM manager - called when DOM is ready
     */
    async initialize() {
        console.log('[HD1-DOM] Initializing DOM Manager');
        
        // Cache all known elements with null safety
        this.cacheElements();
        
        // Setup DOM observers
        this.setupObservers();
        
        this.ready = true;
        console.log('[HD1-DOM] DOM Manager ready');
    }

    /**
     * Cache all DOM elements with null safety
     */
    cacheElements() {
        const elementIds = [
            'hd1-playcanvas-canvas',
            'debug-panel',
            'debug-log',
            'debug-header', 
            'debug-collapse-icon',
            'status-connection-indicator',
            'status-connection-text',
            'status-lock-indicator',
            'session-id-tag-status',
            'debug-status-bar',
            'status-collapse-arrow',
            'channel-selector',
            'debug-content',
            'debug-persistent-content',
            'debug-info-panels',
            'playcanvas-controls',
            'playcanvas-stats-panel',
            'pc-scene-name',
            'pc-entities-count',
            'pc-active-entities',
            'pc-avatar-entities',
            'pc-sync-events',
            'pc-camera-position',
            'pc-last-update',
            'refresh-timer',
            'api-version',
            'js-version',
            'fps-graph-canvas',
            'memory-graph-canvas',
            'latency-graph-canvas',
            'websocket-graph-canvas',
            'debug-collapsed-stats',
            'collapsed-fps-value',
            'collapsed-memory-value',
            'collapsed-latency-value',
            'collapsed-rx-value',
            'collapsed-tx-value'
        ];

        elementIds.forEach(id => {
            const element = document.getElementById(id);
            if (element) {
                this.elements.set(id, element);
                console.log(`[HD1-DOM] Cached element: ${id}`);
            } else {
                console.warn(`[HD1-DOM] Element not found: ${id}`);
            }
        });
    }

    /**
     * Get element with null safety
     */
    get(elementId) {
        return this.elements.get(elementId) || null;
    }

    /**
     * Safely set text content
     */
    setText(elementId, text) {
        const element = this.get(elementId);
        if (element) {
            element.textContent = text;
            return true;
        }
        console.warn(`[HD1-DOM] Cannot set text on missing element: ${elementId}`);
        return false;
    }

    /**
     * Safely set HTML content
     */
    setHTML(elementId, html) {
        const element = this.get(elementId);
        if (element) {
            element.innerHTML = html;
            return true;
        }
        console.warn(`[HD1-DOM] Cannot set HTML on missing element: ${elementId}`);
        return false;
    }

    /**
     * Safely add event listener
     */
    addEventListener(elementId, event, handler, options = {}) {
        const element = this.get(elementId);
        if (element) {
            element.addEventListener(event, handler, options);
            return true;
        }
        console.warn(`[HD1-DOM] Cannot add event listener to missing element: ${elementId}`);
        return false;
    }

    /**
     * Safely add CSS class
     */
    addClass(elementId, className) {
        const element = this.get(elementId);
        if (element) {
            element.classList.add(className);
            return true;
        }
        return false;
    }

    /**
     * Safely remove CSS class
     */
    removeClass(elementId, className) {
        const element = this.get(elementId);
        if (element) {
            element.classList.remove(className);
            return true;
        }
        return false;
    }

    /**
     * Safely check if element exists
     */
    exists(elementId) {
        return this.elements.has(elementId);
    }

    /**
     * Create and append debug log entry with fallback
     */
    addDebugEntry(command, data = null) {
        const debugLog = this.get('debug-log');
        
        if (!debugLog) {
            // Fallback to console if debug log not available
            console.log(`[HD1-Debug] ${command}${data ? ' ' + JSON.stringify(data, null, 0) : ''}`);
            return;
        }

        const time = new Date().toLocaleTimeString();
        const entry = document.createElement('div');
        entry.className = 'debug-entry';
        
        const timeSpan = document.createElement('span');
        timeSpan.className = 'debug-time';
        timeSpan.textContent = time + ' ';
        
        const commandSpan = document.createElement('span');
        commandSpan.className = 'debug-command';
        commandSpan.textContent = command;
        
        entry.appendChild(timeSpan);
        entry.appendChild(commandSpan);
        
        if (data) {
            const dataSpan = document.createElement('span');
            dataSpan.className = 'debug-data';
            dataSpan.textContent = ' ' + JSON.stringify(data, null, 0);
            entry.appendChild(dataSpan);
        }
        
        debugLog.appendChild(entry);
        debugLog.scrollTop = debugLog.scrollHeight;
        
        // Keep only last 50 entries
        while (debugLog.children.length > 50) {
            debugLog.removeChild(debugLog.firstChild);
        }
    }

    /**
     * Setup DOM mutation observers for dynamic content
     */
    setupObservers() {
        // Observer for new elements being added
        const observer = new MutationObserver((mutations) => {
            mutations.forEach((mutation) => {
                if (mutation.type === 'childList') {
                    mutation.addedNodes.forEach((node) => {
                        if (node.nodeType === 1 && node.id) { // Element node with ID
                            if (!this.elements.has(node.id)) {
                                this.elements.set(node.id, node);
                                console.log(`[HD1-DOM] Dynamically cached element: ${node.id}`);
                            }
                        }
                    });
                }
            });
        });

        observer.observe(document.body, {
            childList: true,
            subtree: true
        });

        this.observers.set('mutation', observer);
    }

    /**
     * Cleanup all observers
     */
    cleanup() {
        this.observers.forEach((observer) => {
            observer.disconnect();
        });
        this.observers.clear();
        this.elements.clear();
        this.ready = false;
    }
}

// Export for use in console manager
window.HD1DOMManager = HD1DOMManager;