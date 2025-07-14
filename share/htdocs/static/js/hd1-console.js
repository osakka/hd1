/**
 * HD1 Three.js Console System
 * Comprehensive console interface for HD1 Three.js engine with robot trapping
 */

class HD1Console {
    constructor() {
        console.log('🚀 Initializing HD1 Three.js Console...');
        
        // Core components
        this.apiClient = new HD1ThreeJSAPIClient();
        this.threeJS = null;
        this.sync = null;
        this.uiComponents = null;
        this.formSystem = new HD1FormSystem();
        
        // Console state
        this.isVisible = false;
        this.isDragging = false;
        this.robotTrapActive = true;
        
        // WebSocket for real-time updates
        this.websocket = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 10;
        
        // Initialize console
        this.initializeConsole();
        this.setupEventListeners();
        this.initializeWebSocket();
        this.initializeThreeJS();
        
        console.log('✅ HD1 Console initialized successfully');
    }
    
    initializeConsole() {
        // Create console container
        this.container = document.createElement('div');
        this.container.id = 'hd1-console';
        this.container.innerHTML = `
            <div class="console-header">
                <div class="console-title">HD1 Three.js Console</div>
                <div class="console-controls">
                    <button id="console-toggle-robot-trap" class="control-btn ${this.robotTrapActive ? 'active' : ''}">
                        🤖 Robot Trap
                    </button>
                    <button id="console-minimize" class="control-btn">−</button>
                    <button id="console-close" class="control-btn">×</button>
                </div>
            </div>
            <div class="console-content">
                <div class="console-tabs">
                    <div class="tab active" data-tab="scene">Scene</div>
                    <div class="tab" data-tab="sync">Sync</div>
                    <div class="tab" data-tab="entities">Entities</div>
                    <div class="tab" data-tab="console">Console</div>
                </div>
                <div class="console-panels">
                    <div class="panel active" id="panel-scene">
                        <div class="panel-section">
                            <h4>Scene Control</h4>
                            <button class="action-btn" id="get-scene">Get Scene</button>
                            <button class="action-btn" id="update-scene">Update Scene</button>
                        </div>
                        <div class="panel-section">
                            <h4>Camera Control</h4>
                            <div class="control-group">
                                <label>Position</label>
                                <input type="number" id="cam-x" placeholder="X" value="5">
                                <input type="number" id="cam-y" placeholder="Y" value="5">
                                <input type="number" id="cam-z" placeholder="Z" value="5">
                                <button class="action-btn" id="set-camera">Set</button>
                            </div>
                        </div>
                    </div>
                    <div class="panel" id="panel-sync">
                        <div class="panel-section">
                            <h4>Synchronization</h4>
                            <div class="status-display">
                                <div class="status-item">
                                    <span>Status:</span> 
                                    <span id="sync-status">Connecting...</span>
                                </div>
                                <div class="status-item">
                                    <span>Last Seq:</span> 
                                    <span id="sync-last-seq">0</span>
                                </div>
                                <div class="status-item">
                                    <span>Buffered:</span> 
                                    <span id="sync-buffered">0</span>
                                </div>
                            </div>
                            <button class="action-btn" id="sync-stats">Get Stats</button>
                            <button class="action-btn" id="sync-full">Full Sync</button>
                        </div>
                    </div>
                    <div class="panel" id="panel-entities">
                        <div class="panel-section">
                            <h4>Entity Management</h4>
                            <div class="control-group">
                                <select id="entity-type">
                                    <option value="box">Box</option>
                                    <option value="sphere">Sphere</option>
                                    <option value="plane">Plane</option>
                                    <option value="cylinder">Cylinder</option>
                                </select>
                                <input type="color" id="entity-color" value="#777777">
                                <button class="action-btn" id="create-entity">Create</button>
                            </div>
                            <div class="entity-list" id="entity-list">
                                <div class="list-header">Active Entities</div>
                            </div>
                        </div>
                    </div>
                    <div class="panel" id="panel-console">
                        <div class="panel-section">
                            <h4>Console Output</h4>
                            <div class="console-output" id="console-output"></div>
                            <div class="console-input-group">
                                <input type="text" id="console-input" placeholder="Enter command...">
                                <button class="action-btn" id="console-execute">Execute</button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        `;
        
        // Add styles
        this.addConsoleStyles();
        
        // Add to DOM
        document.body.appendChild(this.container);
        
        // Initial position (top-right)
        this.container.style.right = '20px';
        this.container.style.top = '20px';
    }
    
    addConsoleStyles() {
        const styles = `
            #hd1-console {
                position: fixed;
                width: 400px;
                height: 500px;
                background: rgba(0, 20, 40, 0.95);
                border: 1px solid #00ff00;
                border-radius: 8px;
                font-family: 'Courier New', monospace;
                font-size: 12px;
                color: #00ff00;
                z-index: 10000;
                box-shadow: 0 4px 20px rgba(0, 255, 0, 0.3);
                backdrop-filter: blur(10px);
                display: none;
            }
            
            #hd1-console.visible {
                display: block;
            }
            
            .console-header {
                display: flex;
                justify-content: space-between;
                align-items: center;
                padding: 8px 12px;
                background: rgba(0, 40, 80, 0.8);
                border-bottom: 1px solid #00ff00;
                cursor: move;
                user-select: none;
            }
            
            .console-title {
                font-weight: bold;
                color: #00ffff;
            }
            
            .console-controls {
                display: flex;
                gap: 4px;
            }
            
            .control-btn {
                background: transparent;
                border: 1px solid #00ff00;
                color: #00ff00;
                padding: 2px 8px;
                cursor: pointer;
                border-radius: 3px;
                font-family: inherit;
                font-size: 11px;
            }
            
            .control-btn:hover {
                background: rgba(0, 255, 0, 0.2);
            }
            
            .control-btn.active {
                background: rgba(0, 255, 0, 0.3);
                color: #ffffff;
            }
            
            .console-content {
                height: calc(100% - 40px);
                display: flex;
                flex-direction: column;
            }
            
            .console-tabs {
                display: flex;
                background: rgba(0, 20, 40, 0.5);
                border-bottom: 1px solid #00ff00;
            }
            
            .tab {
                padding: 8px 16px;
                cursor: pointer;
                border-right: 1px solid #00ff00;
                user-select: none;
            }
            
            .tab:hover {
                background: rgba(0, 255, 0, 0.1);
            }
            
            .tab.active {
                background: rgba(0, 255, 0, 0.2);
                color: #ffffff;
            }
            
            .console-panels {
                flex: 1;
                overflow-y: auto;
                padding: 12px;
            }
            
            .panel {
                display: none;
            }
            
            .panel.active {
                display: block;
            }
            
            .panel-section {
                margin-bottom: 16px;
                padding: 8px;
                border: 1px solid rgba(0, 255, 0, 0.3);
                border-radius: 4px;
            }
            
            .panel-section h4 {
                margin: 0 0 8px 0;
                color: #00ffff;
                font-size: 13px;
            }
            
            .action-btn {
                background: transparent;
                border: 1px solid #00ff00;
                color: #00ff00;
                padding: 4px 12px;
                cursor: pointer;
                border-radius: 3px;
                font-family: inherit;
                font-size: 11px;
                margin: 2px;
            }
            
            .action-btn:hover {
                background: rgba(0, 255, 0, 0.2);
            }
            
            .control-group {
                display: flex;
                gap: 4px;
                align-items: center;
                margin: 4px 0;
            }
            
            .control-group input, .control-group select {
                background: rgba(0, 0, 0, 0.5);
                border: 1px solid #00ff00;
                color: #00ff00;
                padding: 2px 4px;
                font-family: inherit;
                font-size: 11px;
                border-radius: 2px;
            }
            
            .control-group input[type="number"] {
                width: 50px;
            }
            
            .control-group input[type="color"] {
                width: 30px;
                height: 24px;
            }
            
            .status-display {
                font-family: 'Courier New', monospace;
                font-size: 11px;
            }
            
            .status-item {
                display: flex;
                justify-content: space-between;
                margin: 2px 0;
            }
            
            .console-output {
                background: rgba(0, 0, 0, 0.8);
                border: 1px solid #00ff00;
                height: 200px;
                overflow-y: auto;
                padding: 8px;
                font-family: 'Courier New', monospace;
                font-size: 10px;
                margin-bottom: 8px;
            }
            
            .console-input-group {
                display: flex;
                gap: 4px;
            }
            
            .console-input-group input {
                flex: 1;
                background: rgba(0, 0, 0, 0.5);
                border: 1px solid #00ff00;
                color: #00ff00;
                padding: 4px;
                font-family: inherit;
                font-size: 11px;
            }
            
            .entity-list {
                max-height: 150px;
                overflow-y: auto;
                border: 1px solid rgba(0, 255, 0, 0.3);
                background: rgba(0, 0, 0, 0.3);
            }
            
            .list-header {
                padding: 4px 8px;
                background: rgba(0, 255, 0, 0.2);
                font-weight: bold;
                border-bottom: 1px solid rgba(0, 255, 0, 0.3);
            }
            
            /* Robot trap indicator */
            .robot-trap-indicator {
                position: fixed;
                top: 10px;
                left: 10px;
                background: rgba(255, 0, 0, 0.8);
                color: white;
                padding: 4px 8px;
                border-radius: 4px;
                font-family: 'Courier New', monospace;
                font-size: 10px;
                z-index: 10001;
                display: none;
            }
            
            .robot-trap-indicator.active {
                display: block;
            }
        `;
        
        const styleElement = document.createElement('style');
        styleElement.textContent = styles;
        document.head.appendChild(styleElement);
    }
    
    setupEventListeners() {
        // Console visibility toggle (Ctrl+~)
        document.addEventListener('keydown', (e) => {
            if (e.ctrlKey && e.key === '`') {
                e.preventDefault();
                this.toggleConsole();
            }
        });
        
        // Robot trap toggle
        document.getElementById('console-toggle-robot-trap').addEventListener('click', () => {
            this.toggleRobotTrap();
        });
        
        // Console controls
        document.getElementById('console-minimize').addEventListener('click', () => {
            this.minimizeConsole();
        });
        
        document.getElementById('console-close').addEventListener('click', () => {
            this.hideConsole();
        });
        
        // Tab switching
        document.querySelectorAll('.tab').forEach(tab => {
            tab.addEventListener('click', () => {
                this.switchTab(tab.dataset.tab);
            });
        });
        
        // Dragging
        const header = document.querySelector('.console-header');
        header.addEventListener('mousedown', (e) => {
            this.startDragging(e);
        });
        
        document.addEventListener('mousemove', (e) => {
            this.handleDrag(e);
        });
        
        document.addEventListener('mouseup', () => {
            this.stopDragging();
        });
        
        // Action buttons
        this.setupActionButtons();
        
        // Auto-show console on startup
        setTimeout(() => {
            this.showConsole();
        }, 1000);
    }
    
    setupActionButtons() {
        // Scene controls
        document.getElementById('get-scene').addEventListener('click', () => {
            this.executeCommand('scene.get');
        });
        
        document.getElementById('update-scene').addEventListener('click', () => {
            this.executeCommand('scene.update', {background: '#000033'});
        });
        
        document.getElementById('set-camera').addEventListener('click', () => {
            const x = parseFloat(document.getElementById('cam-x').value) || 0;
            const y = parseFloat(document.getElementById('cam-y').value) || 0;
            const z = parseFloat(document.getElementById('cam-z').value) || 0;
            this.setCameraPosition(x, y, z);
        });
        
        // Sync controls
        document.getElementById('sync-stats').addEventListener('click', () => {
            this.getSyncStats();
        });
        
        document.getElementById('sync-full').addEventListener('click', () => {
            this.executeCommand('sync.full');
        });
        
        // Entity controls
        document.getElementById('create-entity').addEventListener('click', () => {
            this.createEntity();
        });
        
        // Console command execution
        document.getElementById('console-execute').addEventListener('click', () => {
            this.executeConsoleCommand();
        });
        
        document.getElementById('console-input').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                this.executeConsoleCommand();
            }
        });
    }
    
    toggleConsole() {
        if (this.isVisible) {
            this.hideConsole();
        } else {
            this.showConsole();
        }
    }
    
    showConsole() {
        this.container.classList.add('visible');
        this.isVisible = true;
        this.logToConsole('HD1 Console activated');
    }
    
    hideConsole() {
        this.container.classList.remove('visible');
        this.isVisible = false;
    }
    
    minimizeConsole() {
        // Toggle minimized state
        if (this.container.style.height === '40px') {
            this.container.style.height = '500px';
        } else {
            this.container.style.height = '40px';
        }
    }
    
    toggleRobotTrap() {
        this.robotTrapActive = !this.robotTrapActive;
        const btn = document.getElementById('console-toggle-robot-trap');
        btn.classList.toggle('active', this.robotTrapActive);
        
        // Add visual indicator
        this.updateRobotTrapIndicator();
        
        this.logToConsole(`Robot trap ${this.robotTrapActive ? 'ACTIVATED' : 'DEACTIVATED'}`);
        
        // Robot trap logic - detect automated behavior
        if (this.robotTrapActive) {
            this.activateRobotDetection();
        }
    }
    
    updateRobotTrapIndicator() {
        let indicator = document.querySelector('.robot-trap-indicator');
        if (!indicator) {
            indicator = document.createElement('div');
            indicator.className = 'robot-trap-indicator';
            indicator.textContent = '🤖 ROBOT TRAP ACTIVE';
            document.body.appendChild(indicator);
        }
        
        indicator.classList.toggle('active', this.robotTrapActive);
    }
    
    activateRobotDetection() {
        // Monitor for suspicious automated patterns
        let mouseMovements = [];
        let keystrokes = [];
        
        const detectRobot = () => {
            // Check for perfectly linear mouse movements
            if (mouseMovements.length > 10) {
                const isLinear = this.checkLinearMovement(mouseMovements);
                if (isLinear) {
                    this.logToConsole('⚠️ ROBOT DETECTED: Linear mouse movement pattern');
                    this.triggerRobotTrap('linear_movement');
                }
            }
            
            // Check for rapid, regular keystrokes
            if (keystrokes.length > 5) {
                const isRegular = this.checkRegularKeystrokes(keystrokes);
                if (isRegular) {
                    this.logToConsole('⚠️ ROBOT DETECTED: Regular keystroke pattern');
                    this.triggerRobotTrap('regular_keystrokes');
                }
            }
        };
        
        document.addEventListener('mousemove', (e) => {
            if (!this.robotTrapActive) return;
            
            mouseMovements.push({x: e.clientX, y: e.clientY, time: Date.now()});
            if (mouseMovements.length > 20) mouseMovements.shift();
            
            setTimeout(detectRobot, 100);
        });
        
        document.addEventListener('keydown', (e) => {
            if (!this.robotTrapActive) return;
            
            keystrokes.push({key: e.key, time: Date.now()});
            if (keystrokes.length > 10) keystrokes.shift();
            
            setTimeout(detectRobot, 100);
        });
    }
    
    checkLinearMovement(movements) {
        if (movements.length < 5) return false;
        
        // Check if movements are too perfectly linear
        let deltaX = movements[1].x - movements[0].x;
        let deltaY = movements[1].y - movements[0].y;
        
        for (let i = 2; i < movements.length; i++) {
            let currentDeltaX = movements[i].x - movements[i-1].x;
            let currentDeltaY = movements[i].y - movements[i-1].y;
            
            // Allow small variance for human movement
            if (Math.abs(currentDeltaX - deltaX) > 2 || Math.abs(currentDeltaY - deltaY) > 2) {
                return false;
            }
        }
        return true;
    }
    
    checkRegularKeystrokes(keystrokes) {
        if (keystrokes.length < 5) return false;
        
        // Check for suspiciously regular timing
        let intervals = [];
        for (let i = 1; i < keystrokes.length; i++) {
            intervals.push(keystrokes[i].time - keystrokes[i-1].time);
        }
        
        // Calculate variance in intervals
        let avg = intervals.reduce((a, b) => a + b) / intervals.length;
        let variance = intervals.reduce((a, b) => a + Math.pow(b - avg, 2), 0) / intervals.length;
        
        // If variance is too low, it's likely automated
        return variance < 100; // Very regular timing
    }
    
    triggerRobotTrap(reason) {
        this.logToConsole(`🚨 ROBOT TRAP TRIGGERED: ${reason}`);
        
        // Visual effect
        document.body.style.filter = 'hue-rotate(180deg) invert(1)';
        setTimeout(() => {
            document.body.style.filter = '';
        }, 1000);
        
        // Audio alert (if possible)
        try {
            const audio = new Audio('data:audio/wav;base64,UklGRnoGAABXQVZFZm10IBAAAAABAAEAQB8AAEAfAAABAAgAZGF0YQoGAACBhYqFbF1fdJivrJBhNjVgodDbq2EcBj+a2/LDciUFLIHO8tiJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsFJYHO8diJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsFJYHO8diJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsFJYHO8diJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsFJYHO8diJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsFJYHO8diJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsFJYHO8diJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsFJYHO8diJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsFJYHO8diJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsFJYHO8diJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsFJYHO8diJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsFJYHO8diJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsFJYHO8diJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsFJYHO8diJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsFJYHO8diJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsFJYHO8diJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsFJYHO8diJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsFJYHO8diJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsFJYHO8diJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsFJYHO8diJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsFJYHO8diJNwgZaLvt559NEAxQp+PwtmMcBjiR1/LMeSwFJHfH8N2QQAoUXrTp66hVFApGn+DyvmchBTaH0fLNeSsF');
            audio.play();
        } catch (e) {
            // Silent fail for audio
        }
    }
    
    initializeWebSocket() {
        const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${wsProtocol}//${window.location.host}/ws`;
        
        this.websocket = new WebSocket(wsUrl);
        
        this.websocket.onopen = () => {
            this.logToConsole('WebSocket connected');
            document.getElementById('sync-status').textContent = 'Connected';
            this.reconnectAttempts = 0;
        };
        
        this.websocket.onclose = () => {
            this.logToConsole('WebSocket disconnected');
            document.getElementById('sync-status').textContent = 'Disconnected';
            this.attemptReconnect();
        };
        
        this.websocket.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                this.handleWebSocketMessage(data);
            } catch (e) {
                console.error('Failed to parse WebSocket message:', e);
            }
        };
    }
    
    initializeThreeJS() {
        // Initialize Three.js scene
        this.threeJS = new HD1ThreeJS('hd1-canvas');
        
        // Initialize sync system
        this.sync = new HD1Sync(this.websocket, this.threeJS);
        
        // Initialize UI components
        this.uiComponents = new HD1UIComponents(this.apiClient);
        
        this.logToConsole('Three.js engine initialized');
    }
    
    handleWebSocketMessage(data) {
        // Update sync status display
        if (this.sync) {
            const stats = this.sync.getStats();
            document.getElementById('sync-last-seq').textContent = stats.lastSeenSeq;
            document.getElementById('sync-buffered').textContent = stats.bufferedOperations;
        }
    }
    
    attemptReconnect() {
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            this.logToConsole('Max reconnection attempts reached');
            return;
        }
        
        this.reconnectAttempts++;
        const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
        
        this.logToConsole(`Reconnecting in ${delay/1000}s (attempt ${this.reconnectAttempts})`);
        
        setTimeout(() => {
            this.initializeWebSocket();
        }, delay);
    }
    
    switchTab(tabName) {
        // Update tab appearance
        document.querySelectorAll('.tab').forEach(tab => {
            tab.classList.remove('active');
        });
        document.querySelector(`[data-tab="${tabName}"]`).classList.add('active');
        
        // Update panel visibility
        document.querySelectorAll('.panel').forEach(panel => {
            panel.classList.remove('active');
        });
        document.getElementById(`panel-${tabName}`).classList.add('active');
    }
    
    startDragging(e) {
        this.isDragging = true;
        this.dragOffset = {
            x: e.clientX - this.container.offsetLeft,
            y: e.clientY - this.container.offsetTop
        };
    }
    
    handleDrag(e) {
        if (!this.isDragging) return;
        
        const x = e.clientX - this.dragOffset.x;
        const y = e.clientY - this.dragOffset.y;
        
        // Keep console within viewport
        const maxX = window.innerWidth - this.container.offsetWidth;
        const maxY = window.innerHeight - this.container.offsetHeight;
        
        this.container.style.left = Math.max(0, Math.min(x, maxX)) + 'px';
        this.container.style.top = Math.max(0, Math.min(y, maxY)) + 'px';
        this.container.style.right = 'auto';
    }
    
    stopDragging() {
        this.isDragging = false;
    }
    
    setCameraPosition(x, y, z) {
        if (this.threeJS && this.threeJS.camera) {
            this.threeJS.camera.position.set(x, y, z);
            this.threeJS.camera.lookAt(0, 0, 0);
            this.logToConsole(`Camera position set to (${x}, ${y}, ${z})`);
        }
    }
    
    createEntity() {
        const type = document.getElementById('entity-type').value;
        const color = document.getElementById('entity-color').value;
        
        const entityData = {
            geometry: {
                type: type,
                width: 1,
                height: 1,
                depth: 1,
                radius: 0.5
            },
            material: {
                type: 'phong',
                color: color
            },
            position: {
                x: Math.random() * 4 - 2,
                y: Math.random() * 4 - 2,
                z: Math.random() * 4 - 2
            }
        };
        
        this.executeCommand('entity.create', entityData);
    }
    
    getSyncStats() {
        if (this.sync) {
            const stats = this.sync.getStats();
            this.logToConsole('Sync Stats: ' + JSON.stringify(stats, null, 2));
        }
    }
    
    executeCommand(command, data = null) {
        this.logToConsole(`> ${command} ${data ? JSON.stringify(data) : ''}`);
        
        const parts = command.split('.');
        const category = parts[0];
        const action = parts[1];
        
        switch (category) {
            case 'scene':
                this.handleSceneCommand(action, data);
                break;
            case 'entity':
                this.handleEntityCommand(action, data);
                break;
            case 'sync':
                this.handleSyncCommand(action, data);
                break;
            default:
                this.logToConsole('Unknown command category: ' + category);
        }
    }
    
    handleSceneCommand(action, data) {
        switch (action) {
            case 'get':
                this.apiClient.getScene().then(response => {
                    this.logToConsole('Scene: ' + JSON.stringify(response, null, 2));
                }).catch(error => {
                    this.logToConsole('Error: ' + error.message);
                });
                break;
            case 'update':
                this.apiClient.updateScene(data).then(response => {
                    this.logToConsole('Scene updated: ' + JSON.stringify(response, null, 2));
                }).catch(error => {
                    this.logToConsole('Error: ' + error.message);
                });
                break;
        }
    }
    
    handleEntityCommand(action, data) {
        switch (action) {
            case 'create':
                this.apiClient.createEntity(data).then(response => {
                    this.logToConsole('Entity created: ' + JSON.stringify(response, null, 2));
                }).catch(error => {
                    this.logToConsole('Error: ' + error.message);
                });
                break;
        }
    }
    
    handleSyncCommand(action, data) {
        switch (action) {
            case 'full':
                this.apiClient.getFullSync().then(response => {
                    this.logToConsole('Full sync: ' + JSON.stringify(response, null, 2));
                }).catch(error => {
                    this.logToConsole('Error: ' + error.message);
                });
                break;
        }
    }
    
    executeConsoleCommand() {
        const input = document.getElementById('console-input');
        const command = input.value.trim();
        
        if (!command) return;
        
        this.logToConsole('> ' + command);
        input.value = '';
        
        // Simple command parser
        try {
            if (command.startsWith('help')) {
                this.showHelp();
            } else if (command.startsWith('clear')) {
                this.clearConsole();
            } else if (command.startsWith('stats')) {
                this.getSyncStats();
            } else if (command.startsWith('trap')) {
                this.toggleRobotTrap();
            } else {
                this.executeCommand(command);
            }
        } catch (error) {
            this.logToConsole('Error: ' + error.message);
        }
    }
    
    showHelp() {
        const helpText = `
HD1 Console Commands:
  scene.get - Get current scene state
  scene.update - Update scene properties
  entity.create - Create new entity
  sync.full - Request full synchronization
  stats - Show sync statistics
  trap - Toggle robot trap
  clear - Clear console output
  help - Show this help
  
Keyboard Shortcuts:
  Ctrl+~ - Toggle console visibility
        `;
        this.logToConsole(helpText);
    }
    
    clearConsole() {
        document.getElementById('console-output').innerHTML = '';
    }
    
    logToConsole(message) {
        const output = document.getElementById('console-output');
        const timestamp = new Date().toLocaleTimeString();
        const logEntry = document.createElement('div');
        logEntry.innerHTML = `<span style="color: #666">[${timestamp}]</span> ${message}`;
        output.appendChild(logEntry);
        output.scrollTop = output.scrollHeight;
        
        // Also log to browser console
        console.log(`[HD1] ${message}`);
    }
}

// Auto-initialize when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
        window.hd1Console = new HD1Console();
    });
} else {
    window.hd1Console = new HD1Console();
}