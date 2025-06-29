package server

import (
	"net/http"
)

func ServeHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	html := `<!DOCTYPE html>
<html>
<head>
    <title>THD Holodeck - A-Frame VR</title>
    <script src="/static/js/vendor/aframe.min.js"></script>
    <script src="/static/js/vendor/aframe-environment-component.min.js"></script>
    <script src="/static/js/vendor/aframe-teleport-controls.min.js"></script>
    <script src="/static/js/vendor/aframe-event-set-component.min.js"></script>
    <script src="/static/js/vendor/aframe-look-at-component.min.js"></script>
    <script src="/static/js/vendor/aframe-orbit-controls.min.js"></script>
    <!-- Removed problematic components:
         - aframe-animation-component.min.js (conflicts with core A-Frame animation)
         - aframe-particle-system.js (uses Node.js require())
         - aframe-text-geometry-component.min.js (THREE.FontLoader constructor error)
         - aframe-state-component.min.js (not essential for core functionality)
         - aframe-controller-cursor-component.min.js (VR-specific, not needed for holodeck)
         - aframe-forcegraph-component.min.js (complex, not used yet)
    -->
    <script src="/static/js/thd-aframe.js"></script>
    <style>
        body { margin: 0; padding: 0; background: #000; overflow: hidden; font-family: monospace; }
        a-scene { display: block; }
        
        
        @keyframes pulse {
            0%, 100% { opacity: 1; transform: scale(1); }
            50% { opacity: 0.7; transform: scale(1.1); }
        }
        
        @keyframes flicker {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.3; }
        }
        
        
        #text-overlay {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            pointer-events: none;
            z-index: 10;
            font-family: 'Arial', sans-serif;
        }
        
        .text-element {
            position: absolute;
            white-space: nowrap;
            transform-origin: center center;
            transition: all 0.3s ease;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.5);
            font-weight: bold;
        }
        
        .text-element.animate {
            transition: transform 0.1s linear;
        }
        
        
        #debug-panel {
            position: absolute;
            top: 20px;
            right: 20px;
            width: 300px;
            min-height: 60px;
            max-height: 400px;
            background: rgba(0, 0, 0, 0.7);
            border: 1px solid rgba(0, 255, 255, 0.3);
            border-radius: 6px;
            font-family: 'Courier New', monospace;
            font-size: 10px;
            color: #00ffff;
            z-index: 100;
            overflow: hidden;
            transition: min-height 0.3s ease;
        }
        
        #debug-panel.collapsed {
            min-height: 40px;
        }
        
        #debug-header {
            background: rgba(0, 255, 255, 0.1);
            padding: 4px 8px;
            border-bottom: 1px solid rgba(0, 255, 255, 0.2);
            font-weight: bold;
            cursor: pointer;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        
        #debug-session-bar {
            background: rgba(0, 255, 255, 0.05);
            padding: 4px 8px;
            border-bottom: 1px solid rgba(0, 255, 255, 0.1);
            display: flex;
            justify-content: space-between;
            align-items: center;
            font-size: 9px;
        }
        
        #debug-scene-bar {
            background: rgba(0, 255, 255, 0.03);
            padding: 4px 8px;
            border-bottom: 1px solid rgba(0, 255, 255, 0.1);
            font-size: 9px;
        }
        
        #debug-scene-select {
            background: rgba(0, 0, 0, 0.5);
            border: 1px solid rgba(0, 255, 255, 0.3);
            color: #00ffff;
            font-family: 'Courier New', monospace;
            font-size: 9px;
            padding: 2px 4px;
            border-radius: 2px;
            width: 100%;
            margin-top: 2px;
        }
        
        #debug-scene-select:focus {
            outline: none;
            border-color: rgba(0, 255, 255, 0.6);
        }
        
        #debug-controls-bar {
            background: rgba(0, 255, 255, 0.03);
            padding: 4px 8px;
            border-bottom: 1px solid rgba(0, 255, 255, 0.1);
            display: flex;
            gap: 4px;
            align-items: center;
            font-size: 9px;
        }
        
        .control-btn {
            background: rgba(0, 255, 255, 0.1);
            border: 1px solid rgba(0, 255, 255, 0.3);
            color: #00ffff;
            font-family: 'Courier New', monospace;
            font-size: 8px;
            padding: 3px 6px;
            border-radius: 2px;
            cursor: pointer;
            transition: all 0.2s ease;
            flex: 1;
        }
        
        .control-btn:hover {
            background: rgba(0, 255, 255, 0.2);
            border-color: rgba(0, 255, 255, 0.5);
        }
        
        .control-btn:active {
            background: rgba(0, 255, 255, 0.3);
        }
        
        .control-btn.active {
            background: rgba(255, 0, 0, 0.2);
            border-color: rgba(255, 0, 0, 0.5);
            color: #ff6666;
        }
        
        .recording-status {
            color: rgba(255, 0, 0, 0.8);
            font-weight: bold;
            font-size: 8px;
            margin-left: 4px;
        }
        
        #debug-session-id {
            color: rgba(255, 255, 255, 0.8);
            font-weight: bold;
        }
        
        #debug-status-led {
            width: 6px;
            height: 6px;
            border-radius: 50%;
            background: #666;
            box-shadow: 0 0 3px rgba(102, 102, 102, 0.5);
            transition: all 0.3s ease;
            cursor: help;
            position: relative;
        }
        
        #debug-status-led:hover::after {
            content: attr(data-status);
            position: absolute;
            bottom: 120%;
            left: 50%;
            transform: translateX(-50%);
            background: rgba(0, 0, 0, 0.9);
            color: white;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 9px;
            white-space: nowrap;
            z-index: 1000;
            border: 1px solid rgba(255, 255, 255, 0.2);
        }
        
        #debug-lock-icon {
            font-size: 10px;
            color: #00ffff;
            transition: all 0.3s ease;
            cursor: help;
            position: relative;
            margin-right: 6px;
            opacity: 0.7;
            filter: grayscale(100%);
        }
        
        #debug-lock-icon:hover::after {
            content: attr(data-status);
            position: absolute;
            bottom: 120%;
            left: 50%;
            transform: translateX(-50%);
            background: rgba(0, 0, 0, 0.9);
            color: white;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 9px;
            white-space: nowrap;
            z-index: 1000;
            border: 1px solid rgba(255, 255, 255, 0.2);
        }
        
        #debug-lock-icon.unlocked {
            opacity: 0.5;
        }
        
        #debug-lock-icon.locked {
            opacity: 1.0;
        }
        
        .debug-indicators {
            display: flex;
            align-items: center;
            gap: 6px;
        }
        
        #debug-collapse-icon {
            font-size: 10px;
            transition: all 0.3s ease;
            cursor: pointer;
            opacity: 0.7;
            transform-origin: center;
        }
        
        #debug-collapse-icon.collapsed {
            transform: rotate(180deg);
        }
        
        #session-id-tag {
            background: rgba(0, 255, 255, 0.2);
            color: #00ffff;
            border: 1px solid rgba(0, 255, 255, 0.4);
            padding: 2px 6px;
            border-radius: 4px;
            font-size: 8px;
            font-weight: bold;
            cursor: pointer;
            margin-left: 6px;
            transition: all 0.2s ease;
            user-select: none;
        }
        
        #session-id-tag:hover {
            background: rgba(0, 255, 255, 0.3);
            border-color: rgba(0, 255, 255, 0.6);
            transform: scale(1.05);
        }
        
        #session-id-tag.copied {
            background: rgba(0, 255, 0, 0.3);
            border-color: rgba(0, 255, 0, 0.6);
            color: #00ff00;
            transform: scale(1.1);
        }
        
        #debug-status-led.connecting {
            background: #ff9500;
            box-shadow: 0 0 6px rgba(255, 149, 0, 0.8);
            animation: pulse 1.5s infinite;
        }
        
        #debug-status-led.connected {
            background: #00ff00;
            box-shadow: 0 0 6px rgba(0, 255, 0, 0.8);
        }
        
        #debug-status-led.receiving {
            background: #00ffff;
            box-shadow: 0 0 8px rgba(0, 255, 255, 1);
            animation: flicker 0.2s;
        }
        
        #debug-status-led.error {
            background: #ff0000;
            box-shadow: 0 0 6px rgba(255, 0, 0, 0.8);
            animation: pulse 0.8s infinite;
        }
        
        #debug-log {
            height: 200px;
            overflow-y: auto;
            padding: 4px 8px;
            line-height: 1.2;
            transition: height 0.3s ease;
        }
        
        #debug-log.collapsed {
            height: 0;
            padding: 0;
        }
        
        /* Standard holodeck-themed scrollbar */
        #debug-log::-webkit-scrollbar {
            width: 8px;
        }
        
        #debug-log::-webkit-scrollbar-track {
            background: rgba(0, 0, 0, 0.3);
            border-radius: 4px;
        }
        
        #debug-log::-webkit-scrollbar-thumb {
            background: rgba(0, 255, 255, 0.4);
            border-radius: 4px;
            border: 1px solid rgba(0, 255, 255, 0.2);
        }
        
        #debug-log::-webkit-scrollbar-thumb:hover {
            background: rgba(0, 255, 255, 0.6);
        }
        
        /* Firefox scrollbar theming */
        #debug-log {
            scrollbar-width: thin;
            scrollbar-color: rgba(0, 255, 255, 0.4) rgba(0, 0, 0, 0.3);
        }
        
        .debug-entry {
            margin-bottom: 2px;
            word-wrap: break-word;
        }
        
        .debug-time {
            color: rgba(0, 255, 255, 0.6);
        }
        
        .debug-command {
            color: #00ff00;
        }
        
        .debug-data {
            color: rgba(255, 255, 255, 0.8);
        }
        
    </style>
</head>
<body>
    <a-scene 
        id="holodeck-scene" 
        thd-holodeck
        embedded
        style="height: 100vh; width: 100vw;"
        vr-mode-ui="enabled: false"
        device-orientation-permission-ui="enabled: false">
        
        <!-- THD Holodeck Environment -->
        <a-assets>
            <!-- Asset preloading will go here -->
        </a-assets>
        
        <!-- Default holodeck setup -->
        <a-entity id="holodeck-camera" 
                  camera 
                  look-controls="enabled: true; pointerLockEnabled: true" 
                  wasd-controls="acceleration: 20; fly: false; enabled: true"
                  thd-keyboard-controls=""
                  thd-sprint-controls=""
                  holodeck-boundaries=""
                  position="0 1.7 5">
        </a-entity>
        
        <!-- Default lighting -->
        <a-light type="ambient" color="#404040" intensity="0.4"></a-light>
        <a-light type="directional" position="10 10 5" color="#ffffff" intensity="0.8"></a-light>
        
        <!-- Holodeck coordinate system grid -->
        <a-entity id="holodeck-grid"></a-entity>
        
        <!-- Dynamic objects container -->
        <a-entity id="holodeck-objects"></a-entity>
    </a-scene>
    
    <div id="text-overlay"></div>
    
    <div id="debug-panel">
        <div id="debug-header">
            <div style="display: flex; align-items: center; gap: 6px;">
                <div id="debug-status-led" class="connecting" data-status="Connecting..."></div>
                <span>THD Console</span>
                <span id="session-id-tag" style="display: none;" title="Click to copy session ID">---</span>
            </div>
            <div style="display: flex; align-items: center; gap: 6px;">
                <span id="debug-lock-icon" class="unlocked" data-status="Mouse look available">&#128274;</span>
                <span id="debug-collapse-icon">&#8679;</span>
            </div>
        </div>
        <div id="debug-scene-bar">
            <div style="display: flex; align-items: center; gap: 4px;">
                <span>Scene</span>
                <select id="debug-scene-select">
                    <option value="">Select Scene...</option>
                    <option value="empty">Empty Grid</option>
                    <option value="anime-ui">Anime UI Demo</option>
                    <option value="ultimate">Ultimate Demo</option>
                    <option value="basic-shapes">Basic Shapes</option>
                </select>
            </div>
        </div>
        <div id="debug-controls-bar">
            <button id="photo-btn" class="control-btn photo-btn" title="Freeze-Frame Mode: Save current session state as new scene">FREEZE-FRAME</button>
            <button id="video-btn" class="control-btn video-btn" title="Temporal Sequence Mode: Start/Stop temporal recording">TEMPORAL SEQUENCE</button>
            <span id="recording-status" class="recording-status"></span>
        </div>
        <div id="debug-log"></div>
    </div>
    
    <script>
        const scene = document.getElementById('holodeck-scene');
        const debugLog = document.getElementById('debug-log');
        const debugStatusLed = document.getElementById('debug-status-led');
        const debugLockIcon = document.getElementById('debug-lock-icon');
        const debugHeader = document.getElementById('debug-header');
        const debugCollapseIcon = document.getElementById('debug-collapse-icon');
        const sessionIdTag = document.getElementById('session-id-tag');
        
        let thdManager;
        let ws;
        let lastMessageTime = 0;
        let reconnectAttempts = 0;
        let maxReconnectAttempts = 99;
        let reconnectTimeout;
        let jsVersion = '${JS_VERSION}'; // Server will replace this
        
        // Status management
        function setStatus(status, message) {
            // Update debug panel status LED
            debugStatusLed.className = status;
            debugStatusLed.setAttribute('data-status', message || status);
        }
        
        // Lock status management
        function setLockStatus(status, message) {
            // Update debug panel lock icon
            debugLockIcon.className = status;
            debugLockIcon.setAttribute('data-status', message || status);
        }
        
        // Update debug session ID when session changes
        function updateDebugSession(sessionId) {
            // Update session ID tag
            if (sessionId && sessionId !== 'No Session') {
                const shortId = sessionId.replace('session-', '');
                sessionIdTag.textContent = shortId;
                sessionIdTag.style.display = 'inline';
            } else {
                sessionIdTag.style.display = 'none';
            }
        }
        
        // Debug logging function
        function addDebug(command, data = null) {
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
        
        // Persistent THD session management
        let currentSessionId = localStorage.getItem('thd_session_id');
        let sessionInitialized = localStorage.getItem('thd_session_initialized') === 'true'; // Persistent flag to prevent multiple scene loads
        
        async function ensureSession() {
            addDebug('SESSION_ENSURE', 'ensureSession() called - initialized: ' + sessionInitialized);
            try {
                // Check if we have a persistent session
                if (currentSessionId) {
                    // Verify session still exists
                    const checkResponse = await fetch('/api/sessions/' + currentSessionId);
                    if (checkResponse.ok) {
                        const sessionData = await checkResponse.json();
                        updateDebugSession(currentSessionId);
                        
                        // Initialize world grid if it exists (defer if manager not ready)
                        if (sessionData.world) {
                            if (thdManager && thdManager.initializeWorld) {
                                thdManager.initializeWorld(sessionData.world);
                            } else {
                                // Store world data for later initialization
                                window.pendingWorldData = sessionData.world;
                            }
                        }
                        
                        console.log('THD Session restored:', currentSessionId);
                        setStatus('connected', 'THD session: ' + currentSessionId.slice(-8));
                        
                        // Associate WebSocket client with this session
                        associateSession(currentSessionId);
                        
                        // Auto-load saved scene after session is restored (ONCE only)
                        if (!sessionInitialized) {
                            sessionInitialized = true;
                            localStorage.setItem('thd_session_initialized', 'true');
                            setTimeout(() => {
                                const savedScene = getCookie('thd_scene');
                                if (savedScene && debugSceneSelect) {
                                    debugSceneSelect.value = savedScene;
                                    if (savedScene !== '') {
                                        addDebug('AUTO_SCENE', {scene: savedScene, trigger: 'session_restore'});
                                        loadScene(savedScene);
                                    }
                                }
                            }, 1500); // Wait for session to fully restore
                        }
                        return;
                    } else {
                        // Session expired, clear it and reset initialization flag
                        localStorage.removeItem('thd_session_id');
                        localStorage.removeItem('thd_session_initialized');
                        sessionInitialized = false;
                        currentSessionId = null;
                    }
                }
                
                // Create new session only if needed
                const response = await fetch('/api/sessions', { method: 'POST' });
                const sessionData = await response.json();
                
                if (sessionData.success) {
                    currentSessionId = sessionData.session_id;
                    localStorage.setItem('thd_session_id', currentSessionId);
                    updateDebugSession(currentSessionId);
                    
                    // Initialize world grid in THD manager (defer if not ready)
                    if (sessionData.world) {
                        if (thdManager && thdManager.initializeWorld) {
                            thdManager.initializeWorld(sessionData.world);
                        } else {
                            // Store world data for later initialization
                            window.pendingWorldData = sessionData.world;
                        }
                    }
                    
                    console.log('THD Session created:', currentSessionId);
                    setStatus('connected', 'THD session: ' + currentSessionId.slice(-8));
                    
                    // Associate WebSocket client with this session
                    associateSession(currentSessionId);
                    
                    // Auto-load saved scene after session is established (ONCE only)
                    if (!sessionInitialized) {
                        sessionInitialized = true;
                        localStorage.setItem('thd_session_initialized', 'true');
                        setTimeout(() => {
                            const savedScene = getCookie('thd_scene');
                            if (savedScene && debugSceneSelect) {
                                debugSceneSelect.value = savedScene;
                                if (savedScene !== '') {
                                    addDebug('AUTO_SCENE', {scene: savedScene, trigger: 'session_create'});
                                    loadScene(savedScene);
                                }
                            }
                        }, 1000); // Wait for session to fully establish
                    }
                } else {
                    console.error('Failed to create session:', sessionData);
                    updateDebugSession('Session Failed');
                }
            } catch (error) {
                console.error('Error managing session:', error);
                updateDebugSession('Connection Error');
            }
        }
        
        // Auto-reconnect functionality
        function connectWebSocket() {
            if (reconnectAttempts >= maxReconnectAttempts) {
                setStatus('error', 'max reconnect attempts reached');
                return;
            }
            
            setStatus('connecting', 'connecting... attempt ' + (reconnectAttempts + 1));
            
            ws = new WebSocket('ws://' + window.location.host + '/ws');
            
            ws.onopen = function() {
                reconnectAttempts = 0;
                setStatus('connected', 'connected');
                
                // Send version check
                const versionMsg = {
                    type: 'version_check',
                    js_version: jsVersion
                };
                addDebug('WS_SEND', versionMsg);
                ws.send(JSON.stringify(versionMsg));
                
                // Send client capabilities
                setTimeout(sendClientInfo, 500);
                
                // Load scenes on initial connection
                setTimeout(refreshSceneDropdown, 1000);
                
                // Initialize session connection without object restoration
                setTimeout(ensureSession, 2000);
            };
            
            ws.onmessage = async function(event) {
                try {
                    const message = JSON.parse(event.data);
                    addDebug('WS_RECV', message);
                    
                    // Handle system messages
                    if (message.type === 'version_mismatch') {
                        setStatus('connecting', 'updating interface...');
                        setTimeout(() => window.location.reload(true), 1000);
                        return;
                    }
                    
                    if (message.type === 'log') {
                        console.log('[SERVER]', message.data);
                        return;
                    }
                    
                    if (message.type === 'reload') {
                        setStatus('connecting', 'forced reload...');
                        window.location.reload(true);
                        return;
                    }
                    
                    // Handle scene list changes
                    if (message.type === 'scene_list_changed') {
                        console.log('[THD] Scene list changed, refreshing dropdown');
                        addDebug('SCENE_LIST_CHANGED', 'Refreshing scene dropdown');
                        await refreshSceneDropdown();
                        return;
                    }
                    
                    // Handle browser control messages
                    if (message.type === 'force_refresh') {
                        console.log('[THD] Force refresh command received');
                        if (message.clear_storage) {
                            localStorage.clear();
                        }
                        if (message.session_id) {
                            localStorage.setItem('thd_session_id', message.session_id);
                        }
                        setStatus('connecting', 'THD forced refresh...');
                        window.location.reload(true);
                        return;
                    }
                    
                    // Handle direct canvas control
                    if (message.type === 'canvas_control') {
                        const controlData = message.data || message;
                        console.log('[THD] Canvas control command:', controlData.command, controlData.objects);
                        addDebug('CANVAS_CTRL', {cmd: controlData.command, objs: controlData.objects?.length || 0});
                        if (controlData.clear) {
                            thdManager.processMessage({type: 'clear'});
                        }
                        if (controlData.command === 'delete' && controlData.object_name) {
                            thdManager.processMessage({
                                type: 'delete', 
                                object_name: controlData.object_name
                            });
                            addDebug('DELETE', {obj: controlData.object_name});
                        }
                        if (controlData.objects) {
                            // Convert server objects to renderer format
                            const rendererObjects = controlData.objects.map(obj => {
                                const converted = {
                                    id: obj.id || obj.name,
                                    name: obj.name || obj.id,
                                    type: obj.type,
                                    transform: obj.transform || {
                                        position: { x: obj.x || 0, y: obj.y || 0, z: obj.z || 0 },
                                        scale: { x: obj.scale || 1, y: obj.scale || 1, z: obj.scale || 1 },
                                        rotation: { x: 0, y: 0, z: 0 }
                                    },
                                    color: obj.color || { r: 0.2, g: 0.8, b: 0.2, a: 1.0 },
                                    wireframe: obj.wireframe || false,
                                    visible: obj.visible !== undefined ? obj.visible : true,
                                    // A-Frame specific properties
                                    text: obj.text,
                                    lightType: obj.lightType,
                                    intensity: obj.intensity,
                                    particleType: obj.particleType,
                                    count: obj.count,
                                    material: obj.material,
                                    physics: obj.physics,
                                    lighting: obj.lighting
                                };
                                console.log('[THD] Converted object with text:', converted.text, converted);
                                return converted;
                            });
                            if (!thdManager) {
                                console.error('[THD] THD MANAGER NOT FOUND - This is the root cause!');
                                addDebug('ERROR', 'THD Manager not initialized');
                                return;
                            }
                            if (!thdManager.processMessage) {
                                console.error('[THD] THD_MANAGER.processMessage NOT FOUND - Missing method!');
                                addDebug('ERROR', 'THD Manager.processMessage missing');
                                return;
                            }
                            console.log('[THD] Calling thdManager.processMessage with:', controlData.command, rendererObjects);
                            addDebug('RENDER_CALL', {cmd: controlData.command, count: rendererObjects.length});
                            try {
                                thdManager.processMessage({
                                    type: controlData.command,
                                    objects: rendererObjects
                                });
                                console.log('[THD] THD Manager.processMessage SUCCESS');
                                addDebug('RENDER_OK', 'Objects sent to THD Manager');
                            } catch(e) {
                                console.error('[THD] THD Manager.processMessage FAILED:', e);
                                addDebug('RENDER_FAIL', {error: e.message});
                            }
                        }
                        if (message.camera) {
                            thdManager.processMessage({
                                type: 'camera',
                                ...message.camera
                            });
                        }
                        setStatus('receiving', 'THD canvas control');
                        setTimeout(() => {
                            setStatus('connected', 'THD active • objects: ' + (thdManager?.getObjectCount() || 0));
                        }, 500);
                        return;
                    }
                    
                    // Regular 3D messages
                    if (thdManager) {
                        thdManager.processMessage(message);
                        
                        setStatus('receiving', 'receiving data');
                        clearTimeout(window.receivingTimeout);
                        window.receivingTimeout = setTimeout(() => {
                            setStatus('connected', 'connected • objects: ' + (thdManager?.getObjectCount() || 0));
                        }, 200);
                    }
                    
                } catch (e) {
                    console.error('Failed to process message:', e, event.data);
                    setStatus('error', 'message parse error');
                }
            };
            
            ws.onclose = function() {
                if (reconnectAttempts < maxReconnectAttempts) {
                    reconnectAttempts++;
                    setStatus('error', 'disconnected • retrying in 2s');
                    reconnectTimeout = setTimeout(connectWebSocket, 2000);
                } else {
                    setStatus('error', 'connection failed permanently');
                }
            };
            
            ws.onerror = function(error) {
                console.error('WebSocket error:', error);
            };
        }
        
        // Send logs to server
        function sendLog(level, message, data = null) {
            if (ws && ws.readyState === WebSocket.OPEN) {
                const logMsg = {
                    type: 'client_log',
                    level: level,
                    message: message,
                    data: data,
                    timestamp: new Date().toISOString(),
                    url: window.location.href,
                    userAgent: navigator.userAgent
                };
                addDebug('LOG_' + level.toUpperCase(), {msg: message});
                ws.send(JSON.stringify(logMsg));
            }
        }
        
        // Override console methods to send logs to server
        const originalLog = console.log;
        const originalError = console.error;
        const originalWarn = console.warn;
        
        console.log = function(...args) {
            originalLog.apply(console, args);
            sendLog('info', args.join(' '), args);
        };
        
        console.error = function(...args) {
            originalError.apply(console, args);
            sendLog('error', args.join(' '), args);
        };
        
        console.warn = function(...args) {
            originalWarn.apply(console, args);
            sendLog('warn', args.join(' '), args);
        };
        
        // Global error handler
        window.addEventListener('error', function(e) {
            sendLog('error', 'JavaScript Error: ' + e.message, {
                filename: e.filename,
                lineno: e.lineno,
                colno: e.colno,
                stack: e.error ? e.error.stack : null
            });
        });
        
        try {
            // Initialize A-Frame manager when scene is ready
            scene.addEventListener('loaded', function() {
                thdManager = new THDAFrameManager(scene);
                sendLog('info', 'THD A-Frame Manager initialized successfully');
                console.log('[THD] A-Frame scene loaded and manager ready');
                
                // Initialize any pending world data
                if (window.pendingWorldData) {
                    thdManager.initializeWorld(window.pendingWorldData);
                    window.pendingWorldData = null;
                    console.log('[THD] Applied pending world data');
                }
            });
            
            // Fallback: initialize even if scene doesn't fire loaded event
            setTimeout(function() {
                if (!thdManager) {
                    thdManager = new THDAFrameManager(scene);
                    console.log('[THD] A-Frame manager initialized via fallback');
                    
                    // Initialize any pending world data
                    if (window.pendingWorldData) {
                        thdManager.initializeWorld(window.pendingWorldData);
                        window.pendingWorldData = null;
                        console.log('[THD] Applied pending world data via fallback');
                    }
                }
            }, 2000);
            
        } catch (e) {
            setStatus('error', 'A-Frame not supported');
            sendLog('error', 'A-Frame initialization failed', e.message);
        }
        
        // A-Frame handles canvas sizing automatically
        function resize() {
            // A-Frame scene automatically resizes
            console.log('Window resized - A-Frame handles automatically');
        }
        window.addEventListener('resize', resize);
        
        // A-Frame provides built-in WASD and mouse controls
        console.log('A-Frame controls: WASD to move, mouse to look, VR ready');
        
        
        // Associate WebSocket client with THD session for isolation
        function associateSession(sessionId) {
            if (ws && ws.readyState === WebSocket.OPEN) {
                const associateMsg = {
                    type: 'session_associate',
                    session_id: sessionId
                };
                addDebug('WS_SEND', {type: 'session_associate', session: sessionId});
                ws.send(JSON.stringify(associateMsg));
                console.log('[THD] WebSocket associated with session:', sessionId);
            }
        }
        
        // Send client info to server
        function sendClientInfo() {
            if (ws && ws.readyState === WebSocket.OPEN) {
                const clientInfo = {
                    type: 'client_info',
                    screen: {
                        width: window.innerWidth,
                        height: window.innerHeight,
                        devicePixelRatio: window.devicePixelRatio || 1,
                        orientation: screen.orientation ? screen.orientation.angle : 0
                    },
                    canvas: {
                        width: window.innerWidth,
                        height: window.innerHeight
                    },
                    capabilities: {
                        webgl: !!thdManager,
                        aframe: !!scene,
                        vr: AFRAME && AFRAME.utils.device.checkHeadsetConnected(),
                        touch: 'ontouchstart' in window,
                        mobile: /Mobi|Android/i.test(navigator.userAgent)
                    },
                    performance: {
                        memory: performance.memory ? {
                            used: performance.memory.usedJSHeapSize,
                            total: performance.memory.totalJSHeapSize,
                            limit: performance.memory.jsHeapSizeLimit
                        } : null
                    }
                };
                addDebug('WS_SEND', {type: 'client_info', summary: 'capabilities'});
                ws.send(JSON.stringify(clientInfo));
            }
        }
        
        // Send client info on connect and resize
        window.addEventListener('resize', function() {
            resize();
            setTimeout(sendClientInfo, 100); // Delay to get accurate dimensions
        });
        
        // Add click interaction tracking on A-Frame scene
        scene.addEventListener('click', function(e) {
            if (ws && ws.readyState === WebSocket.OPEN) {
                const interaction = {
                    type: 'interaction',
                    event: 'click',
                    position: {
                        x: e.clientX,
                        y: e.clientY,
                        normalized: {
                            x: e.clientX / window.innerWidth,
                            y: e.clientY / window.innerHeight
                        }
                    },
                    timestamp: Date.now()
                };
                addDebug('CLICK', {x: interaction.position.x, y: interaction.position.y});
                ws.send(JSON.stringify(interaction));
            }
        });
        
        // A-Frame provides built-in controls, no manual setup needed
        function setupMouseControls() {
            // A-Frame handles all mouse and keyboard controls automatically
            // WASD movement via wasd-controls component
            // Mouse look via look-controls component
            console.log('A-Frame controls initialized: WASD to move, mouse to look, VR ready');
            setStatus('connected', 'A-Frame controls ready - WASD to move, mouse to look');
        }
        
        // Start mouse controls after a delay
        setTimeout(setupMouseControls, 1000);
        
        // Debug panel collapsible functionality (initialized after cookie functions)
        let debugCollapsed = false;
        
        // Scene selection management
        const debugSceneSelect = document.getElementById('debug-scene-select');
        
        // Load saved scene from cookie
        function loadSavedScene() {
            const savedScene = getCookie('thd_scene');
            if (savedScene) {
                debugSceneSelect.value = savedScene;
            }
        }
        
        // Save scene to cookie
        function saveScene(sceneId) {
            setCookie('thd_scene', sceneId, 30); // 30 days
        }
        
        // Cookie utilities
        function setCookie(name, value, days) {
            const expires = new Date();
            expires.setTime(expires.getTime() + (days * 24 * 60 * 60 * 1000));
            document.cookie = name + '=' + value + ';expires=' + expires.toUTCString() + ';path=/';
        }
        
        function getCookie(name) {
            const nameEQ = name + '=';
            const ca = document.cookie.split(';');
            for (let i = 0; i < ca.length; i++) {
                let c = ca[i];
                while (c.charAt(0) === ' ') c = c.substring(1, c.length);
                if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
            }
            return null;
        }
        
        // Handle scene selection
        debugSceneSelect.addEventListener('change', function() {
            const selectedScene = this.value;
            if (selectedScene && currentSessionId) {
                addDebug('SCENE_SELECT', {scene: selectedScene, session: currentSessionId, manual: true});
                saveScene(selectedScene);
                loadScene(selectedScene);
            }
        });
        
        // Load scene via API call
        async function loadScene(sceneId) {
            try {
                const response = await fetch('/api/scenes/' + sceneId, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        session_id: currentSessionId
                    })
                });
                
                if (response.ok) {
                    const result = await response.json();
                    addDebug('SCENE_LOADED', {scene: sceneId, objects: result.objects_created || 0});
                    setStatus('receiving', 'Loading scene: ' + sceneId);
                } else {
                    addDebug('SCENE_ERROR', {scene: sceneId, status: response.status});
                    setStatus('error', 'Failed to load scene');
                }
            } catch (error) {
                console.error('Failed to load scene:', error);
                addDebug('SCENE_FAIL', {scene: sceneId, error: error.message});
                setStatus('error', 'Scene load failed');
            }
        }
        
        // Load saved scene on page load
        loadSavedScene();
        
        // Photo/Video Controls Implementation
        const photoBtn = document.getElementById('photo-btn');
        const videoBtn = document.getElementById('video-btn');
        const recordingStatus = document.getElementById('recording-status');
        let isRecording = false;
        let recordingStartTime = null;
        
        // Photo Mode: Save current session state as new scene
        photoBtn.addEventListener('click', async function() {
            if (!currentSessionId) {
                addDebug('PHOTO_ERROR', 'No active session');
                return;
            }
            
            // Prompt for scene name
            const sceneName = prompt('Enter scene name:');
            if (!sceneName) return;
            
            const sceneId = sceneName.toLowerCase().replace(/[^a-z0-9]/g, '-');
            
            try {
                addDebug('PHOTO_START', {session: currentSessionId, scene: sceneId});
                
                const response = await fetch('/api/sessions/' + currentSessionId + '/scenes/save', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        scene_id: sceneId,
                        name: sceneName,
                        description: 'Scene captured from session ' + currentSessionId
                    })
                });
                
                if (response.ok) {
                    const result = await response.json();
                    addDebug('PHOTO_SUCCESS', {scene: sceneId, objects: result.objects_count || 0});
                    setStatus('receiving', 'Freeze-frame saved: ' + sceneName);
                    
                    // Refresh scene dropdown
                    await refreshSceneDropdown();
                } else {
                    addDebug('PHOTO_ERROR', {scene: sceneId, status: response.status});
                    setStatus('error', 'Freeze-frame save failed');
                }
            } catch (error) {
                addDebug('PHOTO_FAIL', {error: error.message});
                setStatus('error', 'Freeze-frame operation failed');
            }
        });
        
        // Video Mode: Start/Stop recording
        videoBtn.addEventListener('click', async function() {
            if (!currentSessionId) {
                addDebug('VIDEO_ERROR', 'No active session');
                return;
            }
            
            try {
                if (!isRecording) {
                    // Start recording
                    const response = await fetch('/api/sessions/' + currentSessionId + '/recording/start', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({
                            name: 'Recording-' + Date.now(),
                            description: 'Temporal session recording'
                        })
                    });
                    
                    if (response.ok) {
                        isRecording = true;
                        recordingStartTime = Date.now();
                        videoBtn.classList.add('active');
                        videoBtn.textContent = 'STOP SEQUENCE';
                        recordingStatus.textContent = 'REC';
                        addDebug('VIDEO_START', {session: currentSessionId});
                        setStatus('receiving', 'Temporal sequence started');
                        
                        // Update recording timer
                        updateRecordingTimer();
                    } else {
                        addDebug('VIDEO_START_ERROR', {status: response.status});
                        setStatus('error', 'Temporal sequence start failed');
                    }
                } else {
                    // Stop recording
                    const response = await fetch('/api/sessions/' + currentSessionId + '/recording/stop', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' }
                    });
                    
                    if (response.ok) {
                        isRecording = false;
                        recordingStartTime = null;
                        videoBtn.classList.remove('active');
                        videoBtn.textContent = 'TEMPORAL SEQUENCE';
                        recordingStatus.textContent = '';
                        addDebug('VIDEO_STOP', {session: currentSessionId});
                        setStatus('receiving', 'Temporal sequence stopped');
                    } else {
                        addDebug('VIDEO_STOP_ERROR', {status: response.status});
                        setStatus('error', 'Temporal sequence stop failed');
                    }
                }
            } catch (error) {
                addDebug('VIDEO_FAIL', {error: error.message});
                setStatus('error', 'Temporal sequence operation failed');
            }
        });
        
        // Update recording timer
        function updateRecordingTimer() {
            if (isRecording && recordingStartTime) {
                const elapsed = Math.floor((Date.now() - recordingStartTime) / 1000);
                const minutes = Math.floor(elapsed / 60);
                const seconds = elapsed % 60;
                recordingStatus.textContent = 'REC ' + minutes + ':' + (seconds < 10 ? '0' : '') + seconds;
                setTimeout(updateRecordingTimer, 1000);
            }
        }
        
        // Refresh scene dropdown from API
        async function refreshSceneDropdown() {
            try {
                const response = await fetch('/api/scenes');
                if (response.ok) {
                    const data = await response.json();
                    const select = document.getElementById('debug-scene-select');
                    
                    // Save current selection
                    const currentValue = select.value;
                    
                    // Clear existing options except first
                    while (select.children.length > 1) {
                        select.removeChild(select.lastChild);
                    }
                    
                    // Add updated scenes
                    if (data.scenes) {
                        data.scenes.forEach(scene => {
                            const option = document.createElement('option');
                            option.value = scene.id;
                            option.textContent = scene.name;
                            select.appendChild(option);
                        });
                    }
                    
                    // Restore selection (saved scene or previous value)
                    const savedScene = getCookie('thd_scene');
                    if (savedScene) {
                        select.value = savedScene;
                    } else if (currentValue) {
                        select.value = currentValue;
                    }
                    
                    addDebug('SCENE_REFRESH', {count: data.scenes?.length || 0, restored: savedScene || currentValue});
                }
            } catch (error) {
                addDebug('SCENE_REFRESH_ERROR', {error: error.message});
            }
        }
        
        // Pointer lock status indicator
        function updatePointerLockStatus() {
            if (document.pointerLockElement) {
                setLockStatus('locked', 'Freelook engaged • Press ESC to exit');
                addDebug('FREELOOK_ACTIVE', 'Freelook engaged • Press ESC to exit');
            } else {
                setLockStatus('unlocked', 'Click holodeck for freelook navigation');
                addDebug('FREELOOK_READY', 'Click holodeck for freelook navigation');
            }
        }
        
        // Listen for pointer lock changes
        document.addEventListener('pointerlockchange', updatePointerLockStatus);
        document.addEventListener('pointerlockerror', function() {
            addDebug('POINTER_ERROR', 'Pointer lock failed');
        });
        
        // Initial status
        updatePointerLockStatus();
        
        // Initialize debug panel with cookie persistence (after cookie functions are defined)
        debugCollapsed = getCookie('thd_console_collapsed') === 'true';
        addDebug('CONSOLE_COOKIE', 'Loaded state: ' + (debugCollapsed ? 'collapsed' : 'expanded'));
        
        function setDebugState(collapsed, saveToCookie = true) {
            debugCollapsed = collapsed;
            const debugPanel = document.getElementById('debug-panel');
            const debugSceneBar = document.getElementById('debug-scene-bar');
            const debugControlsBar = document.getElementById('debug-controls-bar');
            
            if (debugCollapsed) {
                debugLog.classList.add('collapsed');
                debugPanel.classList.add('collapsed');
                debugCollapseIcon.classList.add('collapsed');
                debugSceneBar.style.display = 'none';
                debugControlsBar.style.display = 'none';
            } else {
                debugLog.classList.remove('collapsed');
                debugPanel.classList.remove('collapsed');
                debugCollapseIcon.classList.remove('collapsed');
                debugSceneBar.style.display = 'block';
                debugControlsBar.style.display = 'flex';
            }
            if (saveToCookie) {
                setCookie('thd_console_collapsed', debugCollapsed.toString(), 30); // 30 days
                addDebug('CONSOLE_SAVE', 'State saved: ' + (debugCollapsed ? 'collapsed' : 'expanded'));
            }
        }
        
        debugHeader.addEventListener('click', function(e) {
            // Don't toggle if clicking on session ID tag
            if (e.target === sessionIdTag) return;
            setDebugState(!debugCollapsed, true);
        });
        
        // Session ID tag click handler for copying
        sessionIdTag.addEventListener('click', function(e) {
            e.stopPropagation(); // Prevent header click
            const sessionId = currentSessionId;
            if (sessionId) {
                const shortId = sessionId.replace('session-', '');
                navigator.clipboard.writeText(shortId).then(() => {
                    // Visual feedback
                    sessionIdTag.classList.add('copied');
                    setTimeout(() => {
                        sessionIdTag.classList.remove('copied');
                    }, 500);
                    addDebug('COPY_SESSION', {id: shortId});
                }).catch(err => {
                    // Fallback for older browsers
                    const textArea = document.createElement('textarea');
                    textArea.value = shortId;
                    document.body.appendChild(textArea);
                    textArea.select();
                    document.execCommand('copy');
                    document.body.removeChild(textArea);
                    addDebug('COPY_SESSION_FALLBACK', {id: shortId});
                });
            }
        });
        
        // Initialize debug panel state from cookie (don't save back to cookie)
        setDebugState(debugCollapsed, false);
        
        // Start connection
        connectWebSocket();
    </script>
</body>
</html>`
	
	// Replace version placeholder and serve
	html = ReplaceVersionPlaceholder(html)
	
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}