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
    <title>Holodeck</title>
    <style>
        body { margin: 0; padding: 0; background: #000; overflow: hidden; font-family: monospace; }
        canvas { display: block; }
        
        #status-led {
            position: absolute;
            bottom: 20px;
            right: 20px;
            width: 16px;
            height: 16px;
            border-radius: 50%;
            background: #666;
            box-shadow: 0 0 8px rgba(102, 102, 102, 0.5);
            transition: all 0.3s ease;
            cursor: pointer;
            z-index: 100;
        }
        
        #status-led.connecting {
            background: #ff9500;
            box-shadow: 0 0 16px rgba(255, 149, 0, 0.8);
            animation: pulse 1.5s infinite;
        }
        
        #status-led.connected {
            background: #00ff00;
            box-shadow: 0 0 16px rgba(0, 255, 0, 0.8);
        }
        
        #status-led.receiving {
            background: #00ffff;
            box-shadow: 0 0 20px rgba(0, 255, 255, 1);
            animation: flicker 0.2s;
        }
        
        #status-led.error {
            background: #ff0000;
            box-shadow: 0 0 16px rgba(255, 0, 0, 0.8);
            animation: pulse 0.8s infinite;
        }
        
        @keyframes pulse {
            0%, 100% { opacity: 1; transform: scale(1); }
            50% { opacity: 0.7; transform: scale(1.1); }
        }
        
        @keyframes flicker {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.3; }
        }
        
        #status-tooltip {
            position: absolute;
            bottom: 45px;
            right: 0;
            background: rgba(0, 0, 0, 0.9);
            color: white;
            padding: 8px 12px;
            border-radius: 6px;
            font-size: 12px;
            white-space: nowrap;
            opacity: 0;
            visibility: hidden;
            transition: all 0.3s ease;
            border: 1px solid rgba(255, 255, 255, 0.2);
        }
        
        #status-led:hover #status-tooltip {
            opacity: 1;
            visibility: visible;
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
        
        #session-info {
            position: absolute;
            bottom: 20px;
            left: 20px;
            color: rgba(255, 255, 255, 0.7);
            font-size: 12px;
            font-family: 'Courier New', monospace;
            background: rgba(0, 0, 0, 0.3);
            padding: 8px 12px;
            border-radius: 4px;
            border: 1px solid rgba(255, 255, 255, 0.1);
            z-index: 100;
        }
        
        #session-id {
            font-weight: bold;
        }
    </style>
</head>
<body>
    <canvas id="canvas"></canvas>
    <div id="text-overlay"></div>
    <div id="status-led" class="connecting">
        <div id="status-tooltip">connecting</div>
    </div>
    <div id="session-info">
        <div id="session-id">No Session</div>
    </div>
    
    <script src="/static/js/gl-matrix.js"></script>
    <script src="/static/js/renderer.js"></script>
    <script>
        const canvas = document.getElementById('canvas');
        const statusLed = document.getElementById('status-led');
        const tooltip = document.getElementById('status-tooltip');
        
        let renderer;
        let ws;
        let lastMessageTime = 0;
        let reconnectAttempts = 0;
        let maxReconnectAttempts = 99;
        let reconnectTimeout;
        let jsVersion = '${JS_VERSION}'; // Server will replace this
        
        // Status management
        function setStatus(status, message) {
            statusLed.className = status;
            tooltip.textContent = message;
        }
        
        // Persistent THD session management
        let currentSessionId = localStorage.getItem('thd_session_id');
        async function ensureSession() {
            try {
                // Check if we have a persistent session
                if (currentSessionId) {
                    // Verify session still exists
                    const checkResponse = await fetch('/api/sessions/' + currentSessionId);
                    if (checkResponse.ok) {
                        const sessionData = await checkResponse.json();
                        document.getElementById('session-id').textContent = currentSessionId;
                        
                        // Initialize world grid if it exists
                        if (renderer && renderer.initializeWorld && sessionData.world) {
                            renderer.initializeWorld(sessionData.world);
                        }
                        
                        // Load existing objects in the session
                        try {
                            const objectsResponse = await fetch('/api/sessions/' + currentSessionId + '/objects');
                            if (objectsResponse.ok) {
                                const objectsData = await objectsResponse.json();
                                if (objectsData.objects && objectsData.objects.length > 0) {
                                    // Convert server objects to renderer format
                                    const rendererObjects = objectsData.objects.map(obj => ({
                                        id: obj.name,
                                        name: obj.name,
                                        type: obj.type,
                                        transform: {
                                            position: { x: obj.x || 0, y: obj.y || 0, z: obj.z || 0 },
                                            scale: { x: obj.scale || 1, y: obj.scale || 1, z: obj.scale || 1 },
                                            rotation: { x: 0, y: 0, z: 0 }
                                        },
                                        color: { r: 1, g: 1, b: 1, a: 1 },
                                        visible: true
                                    }));
                                    
                                    // Send converted objects to renderer
                                    renderer.processMessage({
                                        type: 'create',
                                        objects: rendererObjects
                                    });
                                    console.log('Loaded', rendererObjects.length, 'existing objects');
                                }
                            }
                        } catch (e) {
                            console.log('No existing objects to load');
                        }
                        
                        console.log('THD Session restored:', currentSessionId);
                        setStatus('connected', 'THD session: ' + currentSessionId.slice(-8));
                        return;
                    } else {
                        // Session expired, clear it
                        localStorage.removeItem('thd_session_id');
                        currentSessionId = null;
                    }
                }
                
                // Create new session only if needed
                const response = await fetch('/api/sessions', { method: 'POST' });
                const sessionData = await response.json();
                
                if (sessionData.success) {
                    currentSessionId = sessionData.session_id;
                    localStorage.setItem('thd_session_id', currentSessionId);
                    document.getElementById('session-id').textContent = currentSessionId;
                    
                    // Initialize world grid in renderer
                    if (renderer && renderer.initializeWorld) {
                        renderer.initializeWorld(sessionData.world);
                    }
                    
                    // Load any existing objects (in case session already had objects)
                    try {
                        const objectsResponse = await fetch('/api/sessions/' + currentSessionId + '/objects');
                        if (objectsResponse.ok) {
                            const objectsData = await objectsResponse.json();
                            if (objectsData.objects && objectsData.objects.length > 0) {
                                // Convert server objects to renderer format
                                const rendererObjects = objectsData.objects.map(obj => ({
                                    id: obj.name,
                                    name: obj.name,
                                    type: obj.type,
                                    transform: {
                                        position: { x: obj.x || 0, y: obj.y || 0, z: obj.z || 0 },
                                        scale: { x: obj.scale || 1, y: obj.scale || 1, z: obj.scale || 1 },
                                        rotation: { x: 0, y: 0, z: 0 }
                                    },
                                    color: { r: 1, g: 1, b: 1, a: 1 },
                                    visible: true
                                }));
                                
                                renderer.processMessage({
                                    type: 'create',
                                    objects: rendererObjects
                                });
                                console.log('Loaded', rendererObjects.length, 'existing objects');
                            }
                        }
                    } catch (e) {
                        console.log('No existing objects to load');
                    }
                    
                    console.log('THD Session created:', currentSessionId);
                    setStatus('connected', 'THD session: ' + currentSessionId.slice(-8));
                } else {
                    console.error('Failed to create session:', sessionData);
                    document.getElementById('session-id').textContent = 'Session Failed';
                }
            } catch (error) {
                console.error('Error managing session:', error);
                document.getElementById('session-id').textContent = 'Connection Error';
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
                ws.send(JSON.stringify({
                    type: 'version_check',
                    js_version: jsVersion
                }));
                
                // Send client capabilities
                setTimeout(sendClientInfo, 500);
                
                // Ensure session exists and initialize world with grid (delay to avoid race conditions)
                setTimeout(ensureSession, 2000);
            };
            
            ws.onmessage = function(event) {
                try {
                    const message = JSON.parse(event.data);
                    
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
                        console.log('[THD] Canvas control command:', message.command, message.objects);
                        if (message.clear) {
                            renderer.processMessage({type: 'clear'});
                        }
                        if (message.objects) {
                            // Convert server objects to renderer format
                            const rendererObjects = message.objects.map(obj => {
                                const converted = {
                                    id: obj.id || obj.name,
                                    name: obj.name || obj.id,
                                    type: obj.type,
                                    transform: obj.transform || {
                                        position: { x: obj.x || 0, y: obj.y || 0, z: obj.z || 0 },
                                        scale: { x: obj.scale || 1, y: obj.scale || 1, z: obj.scale || 1 },
                                        rotation: { x: 0, y: 0, z: 0 }
                                    },
                                    color: obj.color || { r: 1, g: 1, b: 1, a: 1 },
                                    wireframe: obj.wireframe || false,
                                    visible: obj.visible !== undefined ? obj.visible : true
                                };
                                console.log('[THD] Converted object:', converted);
                                return converted;
                            });
                            renderer.processMessage({
                                type: message.command,
                                objects: rendererObjects
                            });
                            console.log('[THD] Sent to renderer:', message.command, rendererObjects);
                        }
                        if (message.camera) {
                            renderer.processMessage({
                                type: 'camera',
                                ...message.camera
                            });
                        }
                        setStatus('receiving', 'THD canvas control');
                        setTimeout(() => {
                            setStatus('connected', 'THD active • objects: ' + renderer.objects.size);
                        }, 500);
                        return;
                    }
                    
                    // Regular 3D messages
                    renderer.processMessage(message);
                    
                    setStatus('receiving', 'receiving data');
                    clearTimeout(window.receivingTimeout);
                    window.receivingTimeout = setTimeout(() => {
                        setStatus('connected', 'connected • objects: ' + renderer.objects.size);
                    }, 200);
                    
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
                ws.send(JSON.stringify({
                    type: 'client_log',
                    level: level,
                    message: message,
                    data: data,
                    timestamp: new Date().toISOString(),
                    url: window.location.href,
                    userAgent: navigator.userAgent
                }));
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
            renderer = new HolodeckRenderer(canvas);
            sendLog('info', 'HolodeckRenderer initialized successfully');
        } catch (e) {
            setStatus('error', 'webgl not supported');
            sendLog('error', 'WebGL initialization failed', e.message);
        }
        
        // Resize canvas to window
        function resize() {
            canvas.width = window.innerWidth;
            canvas.height = window.innerHeight;
        }
        window.addEventListener('resize', resize);
        resize();
        
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
                        width: canvas.width,
                        height: canvas.height
                    },
                    capabilities: {
                        webgl: !!renderer,
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
                ws.send(JSON.stringify(clientInfo));
            }
        }
        
        // Send client info on connect and resize
        window.addEventListener('resize', function() {
            resize();
            setTimeout(sendClientInfo, 100); // Delay to get accurate dimensions
        });
        
        // Add click interaction tracking
        canvas.addEventListener('click', function(e) {
            if (ws && ws.readyState === WebSocket.OPEN) {
                const rect = canvas.getBoundingClientRect();
                const interaction = {
                    type: 'interaction',
                    event: 'click',
                    position: {
                        x: e.clientX - rect.left,
                        y: e.clientY - rect.top,
                        normalized: {
                            x: (e.clientX - rect.left) / canvas.width,
                            y: (e.clientY - rect.top) / canvas.height
                        }
                    },
                    timestamp: Date.now()
                };
                ws.send(JSON.stringify(interaction));
            }
        });
        
        // Wait for renderer to be ready, then add mouse controls
        function setupMouseControls() {
            if (!renderer) {
                setTimeout(setupMouseControls, 100);
                return;
            }
            
            let mouseDown = false;
            let lastMouseX = 0;
            let lastMouseY = 0;
            let cameraDistance = 8;
            let cameraAngleX = 0.3;  // Start slightly angled
            let cameraAngleY = 0.3;
            
            function updateCameraFromAngles() {
                const x = Math.cos(cameraAngleY) * Math.cos(cameraAngleX) * cameraDistance;
                const y = Math.sin(cameraAngleX) * cameraDistance;
                const z = Math.sin(cameraAngleY) * Math.cos(cameraAngleX) * cameraDistance;
                
                renderer.camera.position = [x, y, z];
                renderer.camera.target = [0, 0, 0];
            }
            
            // Mouse down
            canvas.onmousedown = function(e) {
                mouseDown = true;
                lastMouseX = e.clientX;
                lastMouseY = e.clientY;
                setStatus('connected', 'MOUSE ACTIVE - drag to orbit, wheel to zoom');
                return false;
            };
            
            // Mouse up
            document.onmouseup = function(e) {
                mouseDown = false;
                setStatus('connected', 'connected - click and drag to control camera');
            };
            
            // Mouse move
            document.onmousemove = function(e) {
                if (mouseDown) {
                    const deltaX = e.clientX - lastMouseX;
                    const deltaY = e.clientY - lastMouseY;
                    
                    cameraAngleY += deltaX * 0.005;
                    cameraAngleX -= deltaY * 0.005;
                    
                    // Clamp vertical angle
                    cameraAngleX = Math.max(-1.4, Math.min(1.4, cameraAngleX));
                    
                    updateCameraFromAngles();
                    
                    lastMouseX = e.clientX;
                    lastMouseY = e.clientY;
                }
            };
            
            // Mouse wheel
            canvas.onwheel = function(e) {
                cameraDistance += e.deltaY * 0.005;
                cameraDistance = Math.max(1, Math.min(15, cameraDistance));
                updateCameraFromAngles();
                e.preventDefault();
                return false;
            };
            
            // Set initial camera position
            updateCameraFromAngles();
            setStatus('connected', 'Mouse controls ready - click and drag to orbit');
        }
        
        // Start mouse controls after a delay
        setTimeout(setupMouseControls, 1000);
        
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