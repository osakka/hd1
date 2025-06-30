const scene = document.getElementById('holodeck-scene');
const debugLog = document.getElementById('debug-log');
const debugHeader = document.getElementById('debug-header');
const debugCollapseIcon = document.getElementById('debug-collapse-icon');

// Status bar elements
const statusConnectionIndicator = document.getElementById('status-connection-indicator');
const statusConnectionText = document.getElementById('status-connection-text');
const statusLockIndicator = document.getElementById('status-lock-indicator');
const sessionIdTagStatus = document.getElementById('session-id-tag-status');
const statusBar = document.getElementById('debug-status-bar');
const statusCollapseArrow = document.getElementById('status-collapse-arrow');

let hd1Manager;
let ws;
let lastMessageTime = 0;
let reconnectAttempts = 0;
let maxReconnectAttempts = 99;
let reconnectTimeout;
let jsVersion = '${JS_VERSION}'; // Server will replace this

// Status management
function setStatus(status, message) {
    // Update status bar connection indicator
    statusConnectionIndicator.className = status;
    statusConnectionIndicator.setAttribute('data-status', message || status);
    
    switch(status) {
        case 'connecting':
            statusConnectionText.textContent = 'Connecting';
            break;
        case 'connected':
            statusConnectionText.textContent = 'Connected';
            break;
        case 'disconnected':
            statusConnectionText.textContent = 'Disconnected';
            break;
        case 'error':
            statusConnectionText.textContent = 'Error';
            break;
        case 'receiving':
            statusConnectionText.textContent = 'Receiving';
            break;
        default:
            statusConnectionText.textContent = 'Unknown';
    }
}

// Lock status management
function setLockStatus(status, message) {
    // Update status bar lock indicator
    statusLockIndicator.className = status;
    statusLockIndicator.setAttribute('data-status', message || status);
    statusLockIndicator.textContent = status === 'locked' ? 'ESC' : '';
}

// Update debug session ID when session changes
function updateDebugSession(sessionId) {
    // Update status bar session ID tag
    if (sessionId && sessionId !== 'No Session') {
        sessionIdTagStatus.textContent = sessionId;
        sessionIdTagStatus.style.display = 'inline';
    } else {
        sessionIdTagStatus.style.display = 'none';
    }
}

// Update environmental info display (read-only)
function updateEnvironmentalDisplay() {
    if (!hd1Manager || !hd1Manager.state) return;
    
    const sceneState = hd1Manager.state;
    const environment = sceneState.environment;
    
    // Get display elements
    const scaleDisplay = document.getElementById('environment-scale');
    const gravityDisplay = document.getElementById('environment-gravity');
    const temperatureDisplay = document.getElementById('environment-temperature');
    const atmosphereDisplay = document.getElementById('environment-atmosphere');
    
    // Scale info mapping
    const scaleInfo = {
        'nm': { name: 'Molecular', unit: 'nm' },
        'μm': { name: 'Microscopic', unit: 'μm' },
        'mm': { name: 'Precision', unit: 'mm' },
        'cm': { name: 'Small Objects', unit: 'cm' },
        'm': { name: 'Human Scale', unit: 'm' },
        'km': { name: 'Landscape', unit: 'km' },
        'Mm': { name: 'Continental', unit: 'Mm' },
        'Gm': { name: 'Planetary', unit: 'Gm' }
    };
    
    // Atmosphere type mapping
    const atmosphereTypes = {
        'air': 'Earth Surface',
        'vacuum': 'Space Vacuum',
        'thin_air': 'High Altitude',
        'liquid': 'Underwater'
    };
    
    const scale = scaleInfo[environment.scale] || { name: 'Unknown', unit: environment.scale };
    const atmosphereType = atmosphereTypes[environment.atmosphere.composition] || 'Custom';
    
    // Update displays
    if (scaleDisplay) {
        scaleDisplay.textContent = `${scale.name} (${scale.unit})`;
    }
    
    if (gravityDisplay) {
        gravityDisplay.textContent = `${environment.gravity} m/s²`;
    }
    
    if (temperatureDisplay) {
        const celsius = (environment.temperature - 273.15).toFixed(1);
        temperatureDisplay.textContent = `${celsius}°C`;
    }
    
    if (atmosphereDisplay) {
        const density = environment.atmosphere.density;
        atmosphereDisplay.textContent = `${atmosphereType} (${density} kg/m³)`;
    }
    
    console.log(`[HD1] Environmental display updated: ${scale.name} scale, ${environment.gravity}m/s², ${(environment.temperature-273.15).toFixed(1)}°C`);
}

// Hardware monitoring
function initializeHardwareMonitoring() {
    const fpsDisplay = document.getElementById('hardware-fps');
    const gpuDisplay = document.getElementById('hardware-gpu');
    const memoryDisplay = document.getElementById('hardware-memory');
    const rendererDisplay = document.getElementById('hardware-renderer');
    
    let frameCount = 0;
    let lastTime = performance.now();
    let fps = 0;
    
    // FPS monitoring
    function updateFPS() {
        frameCount++;
        const currentTime = performance.now();
        
        if (currentTime - lastTime >= 1000) {
            fps = Math.round((frameCount * 1000) / (currentTime - lastTime));
            frameCount = 0;
            lastTime = currentTime;
            
            if (fpsDisplay) {
                fpsDisplay.textContent = fps;
            }
        }
        
        requestAnimationFrame(updateFPS);
    }
    
    // GPU and WebGL info
    function updateGPUInfo() {
        try {
            // Create temporary canvas to avoid context conflicts
            const tempCanvas = document.createElement('canvas');
            const gl = tempCanvas.getContext('webgl') || tempCanvas.getContext('experimental-webgl');
            
            if (gl) {
                const debugInfo = gl.getExtension('WEBGL_debug_renderer_info');
                let gpuInfo = 'WebGL Ready';
                
                if (debugInfo) {
                    const renderer = gl.getParameter(debugInfo.UNMASKED_RENDERER_WEBGL);
                    gpuInfo = renderer.split('(')[0].trim();
                    if (gpuInfo.length > 15) {
                        gpuInfo = gpuInfo.substring(0, 15) + '...';
                    }
                }
                
                if (gpuDisplay) {
                    gpuDisplay.textContent = gpuInfo;
                }
                
                // Clean up
                const loseContext = gl.getExtension('WEBGL_lose_context');
                if (loseContext) loseContext.loseContext();
            } else {
                if (gpuDisplay) {
                    gpuDisplay.textContent = 'WebGL Not Available';
                }
            }
        } catch (error) {
            if (gpuDisplay) {
                gpuDisplay.textContent = 'WebGL Detection Error';
            }
        }
    }
    
    // Memory monitoring  
    function updateMemoryInfo() {
        if (performance.memory) {
            const used = Math.round(performance.memory.usedJSHeapSize / 1048576); // MB
            const total = Math.round(performance.memory.totalJSHeapSize / 1048576); // MB
            
            if (memoryDisplay) {
                memoryDisplay.textContent = `${used}/${total} MB`;
            }
        }
    }
    
    // Initialize monitoring
    updateFPS();
    updateGPUInfo();
    
    // Update GPU info once A-Frame is ready
    setTimeout(updateGPUInfo, 2000);
    
    // Update memory info periodically
    setInterval(updateMemoryInfo, 2000);
    
    console.log('[HD1] Hardware monitoring initialized');
}

// Objects tracking display - counts A-Frame entities directly
function updateObjectsDisplay() {
    // Count objects directly from A-Frame scene instead of broken reactive state
    const scene = document.getElementById('holodeck-scene');
    const objectsContainer = document.getElementById('holodeck-objects');
    
    if (!scene || !objectsContainer) return;
    
    // Count all child entities in the objects container
    const allEntities = objectsContainer.querySelectorAll('a-entity, a-box, a-sphere, a-cylinder, a-plane, a-cone, a-light');
    const totalObjects = allEntities.length;
    
    let visibleObjects = 0;
    let lightObjects = 0;
    let primitiveObjects = 0;
    
    allEntities.forEach(entity => {
        // Check visibility
        const visible = entity.getAttribute('visible');
        if (visible !== false && visible !== 'false') visibleObjects++;
        
        // Check type by tag name or geometry component
        const tagName = entity.tagName.toLowerCase();
        const geometry = entity.getAttribute('geometry');
        
        if (tagName === 'a-light' || entity.hasAttribute('light')) {
            lightObjects++;
        } else if (['a-box', 'a-sphere', 'a-cylinder', 'a-plane', 'a-cone'].includes(tagName) || 
                   (geometry && ['box', 'sphere', 'cylinder', 'plane', 'cone'].includes(geometry.primitive))) {
            primitiveObjects++;
        }
    });
    
    // Get current scene name
    const sceneSelect = document.getElementById('debug-scene-select');
    const currentScene = sceneSelect ? sceneSelect.value : 'None';
    const sceneName = currentScene === '' ? 'None' : currentScene;
    
    // Update displays
    const totalDisplay = document.getElementById('objects-total');
    const visibleDisplay = document.getElementById('objects-visible');
    const lightsDisplay = document.getElementById('objects-lights');
    const primitivesDisplay = document.getElementById('objects-primitives');
    const lastUpdateDisplay = document.getElementById('objects-last-update');
    const sceneDisplay = document.getElementById('objects-scene');
    
    if (totalDisplay) totalDisplay.textContent = totalObjects;
    if (visibleDisplay) visibleDisplay.textContent = visibleObjects;
    if (lightsDisplay) lightsDisplay.textContent = lightObjects;
    if (primitivesDisplay) primitivesDisplay.textContent = primitiveObjects;
    if (lastUpdateDisplay) lastUpdateDisplay.textContent = new Date().toLocaleTimeString();
    if (sceneDisplay) sceneDisplay.textContent = sceneName;
    
    // Only log when objects change significantly
    if (totalObjects > 0) {
        console.log(`[HD1] Objects: ${totalObjects} total, ${visibleObjects} visible, ${lightObjects} lights, ${primitiveObjects} primitives`);
    }
}

// Debounced objects update using MutationObserver
let objectsUpdateTimeout;
function scheduleObjectsUpdate() {
    clearTimeout(objectsUpdateTimeout);
    objectsUpdateTimeout = setTimeout(() => {
        updateObjectsDisplay();
    }, 100); // Small delay to batch multiple DOM changes
}

// Initialize MutationObserver to watch for object changes
function initializeObjectsObserver() {
    const objectsContainer = document.getElementById('holodeck-objects');
    if (!objectsContainer) return;
    
    const observer = new MutationObserver(() => {
        scheduleObjectsUpdate();
    });
    
    observer.observe(objectsContainer, {
        childList: true,
        subtree: true,
        attributes: true,
        attributeFilter: ['visible', 'geometry', 'light']
    });
    
    console.log('[HD1] Objects MutationObserver initialized');
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

// Persistent HD1 session management
let currentSessionId = localStorage.getItem('hd1_session_id');
let sessionInitialized = localStorage.getItem('hd1_session_initialized') === 'true'; // Persistent flag to prevent multiple scene loads

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
                    if (hd1Manager && hd1Manager.initializeWorld) {
                        hd1Manager.initializeWorld(sessionData.world);
                    } else {
                        // Store world data for later initialization
                        window.pendingWorldData = sessionData.world;
                    }
                }
                
                console.log('HD1 Session restored:', currentSessionId);
                setStatus('connected', 'HD1 session: ' + currentSessionId.slice(-8));
                
                // Associate WebSocket client with this session
                associateSession(currentSessionId);
                
                // Auto-load saved scene after session is restored (ONCE only)
                if (!sessionInitialized) {
                    sessionInitialized = true;
                    localStorage.setItem('hd1_session_initialized', 'true');
                    
                    // Load saved scene immediately - no delays
                    const savedScene = getCookie('hd1_scene');
                    console.log(`[HD1] Bootstrap check: savedScene="${savedScene}", debugSceneSelect=${!!debugSceneSelect}, sessionInitialized=${sessionInitialized}`);
                    
                    if (savedScene && savedScene !== '' && debugSceneSelect) {
                        debugSceneSelect.value = savedScene;
                        addDebug('AUTO_SCENE', {scene: savedScene, trigger: 'session_restore'});
                        loadScene(savedScene);
                        console.log(`[HD1] Bootstrap: Auto-loading scene "${savedScene}"`);
                    } else {
                        console.log(`[HD1] Bootstrap: No scene to auto-load - savedScene="${savedScene}"`);
                    }
                }
                return;
            } else {
                // Session expired, clear it and reset initialization flag
                localStorage.removeItem('hd1_session_id');
                localStorage.removeItem('hd1_session_initialized');
                sessionInitialized = false;
                currentSessionId = null;
            }
        }
        
        // Create new session only if needed
        const response = await fetch('/api/sessions', { method: 'POST' });
        const sessionData = await response.json();
        
        if (sessionData.success) {
            currentSessionId = sessionData.session_id;
            localStorage.setItem('hd1_session_id', currentSessionId);
            updateDebugSession(currentSessionId);
            
            // Initialize world grid in HD1 manager (defer if not ready)
            if (sessionData.world) {
                if (hd1Manager && hd1Manager.initializeWorld) {
                    hd1Manager.initializeWorld(sessionData.world);
                } else {
                    // Store world data for later initialization
                    window.pendingWorldData = sessionData.world;
                }
            }
            
            console.log('HD1 Session created:', currentSessionId);
            setStatus('connected', 'HD1 session: ' + currentSessionId.slice(-8));
            
            // Associate WebSocket client with this session
            associateSession(currentSessionId);
            
            // Auto-load saved scene after session is established (ONCE only)
            if (!sessionInitialized) {
                sessionInitialized = true;
                localStorage.setItem('hd1_session_initialized', 'true');
                
                // Load saved scene immediately - no delays
                const savedScene = getCookie('hd1_scene');
                console.log(`[HD1] Bootstrap check (new session): savedScene="${savedScene}", debugSceneSelect=${!!debugSceneSelect}`);
                
                if (savedScene && savedScene !== '' && debugSceneSelect) {
                    debugSceneSelect.value = savedScene;
                    addDebug('AUTO_SCENE', {scene: savedScene, trigger: 'session_create'});
                    loadScene(savedScene);
                    console.log(`[HD1] Bootstrap: Auto-loading scene "${savedScene}" for new session`);
                } else {
                    console.log(`[HD1] Bootstrap: No scene to auto-load for new session - savedScene="${savedScene}"`);
                }
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
                console.log('[HD1] Scene list changed, refreshing dropdown');
                addDebug('SCENE_LIST_CHANGED', 'Refreshing scene dropdown');
                await refreshSceneDropdown();
                return;
            }
            
            // Handle browser control messages
            if (message.type === 'force_refresh') {
                console.log('[HD1] Force refresh command received');
                if (message.clear_storage) {
                    localStorage.clear();
                }
                if (message.session_id) {
                    localStorage.setItem('hd1_session_id', message.session_id);
                }
                setStatus('connecting', 'HD1 forced refresh...');
                window.location.reload(true);
                return;
            }
            
            // Handle direct canvas control
            if (message.type === 'canvas_control') {
                const controlData = message.data || message;
                console.log('[HD1] Canvas control command:', controlData.command, controlData.objects);
                addDebug('CANVAS_CTRL', {cmd: controlData.command, objs: controlData.objects?.length || 0});
                if (controlData.clear) {
                    hd1Manager.processMessage({type: 'clear'});
                }
                if (controlData.command === 'delete' && controlData.object_name) {
                    hd1Manager.processMessage({
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
                        console.log('[HD1] Converted object with text:', converted.text, converted);
                        return converted;
                    });
                    if (!hd1Manager) {
                        console.error('[HD1] HD1 MANAGER NOT FOUND - This is the root cause!');
                        addDebug('ERROR', 'HD1 Manager not initialized');
                        return;
                    }
                    if (!hd1Manager.processMessage) {
                        console.error('[HD1] HD1_MANAGER.processMessage NOT FOUND - Missing method!');
                        addDebug('ERROR', 'HD1 Manager.processMessage missing');
                        return;
                    }
                    console.log('[HD1] Calling hd1Manager.processMessage with:', controlData.command, rendererObjects);
                    addDebug('RENDER_CALL', {cmd: controlData.command, count: rendererObjects.length});
                    try {
                        hd1Manager.processMessage({
                            type: controlData.command,
                            objects: rendererObjects
                        });
                        console.log('[HD1] HD1 Manager.processMessage SUCCESS');
                        addDebug('RENDER_OK', 'Objects sent to HD1 Manager');
                        
                        // Update environmental display and schedule objects update
                        if (hd1Manager.state) {
                            if (typeof updateEnvironmentalDisplay === 'function') {
                                updateEnvironmentalDisplay();
                            }
                            scheduleObjectsUpdate();
                        }
                    } catch(e) {
                        console.error('[HD1] HD1 Manager.processMessage FAILED:', e);
                        addDebug('RENDER_FAIL', {error: e.message});
                    }
                }
                if (message.camera) {
                    hd1Manager.processMessage({
                        type: 'camera',
                        ...message.camera
                    });
                }
                setStatus('receiving', 'HD1 canvas control');
                setTimeout(() => {
                    setStatus('connected', 'HD1 active • objects: ' + (hd1Manager?.getObjectCount() || 0));
                }, 500);
                return;
            }
            
            // Regular 3D messages
            if (hd1Manager) {
                hd1Manager.processMessage(message);
                
                setStatus('receiving', 'receiving data');
                clearTimeout(window.receivingTimeout);
                window.receivingTimeout = setTimeout(() => {
                    setStatus('connected', 'connected • objects: ' + (hd1Manager?.getObjectCount() || 0));
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
    // Initialize REACTIVE A-Frame manager when scene is ready
    scene.addEventListener('loaded', function() {
        hd1Manager = new HD1ReactiveManager(scene);
        sendLog('info', 'HD1 Reactive Manager initialized successfully');
        console.log('[HD1] A-Frame scene loaded and REACTIVE manager ready');
        
        // Initialize any pending world data
        if (window.pendingWorldData) {
            hd1Manager.initializeWorld(window.pendingWorldData);
            window.pendingWorldData = null;
            console.log('[HD1] Applied pending world data');
        }
        
        // Initialize displays immediately when manager is ready
        initializeDisplays();
    });
    
    // Fallback: initialize even if scene doesn't fire loaded event
    setTimeout(function() {
        if (!hd1Manager) {
            hd1Manager = new HD1ReactiveManager(scene);
            console.log('[HD1] REACTIVE A-Frame manager initialized via fallback');
            
            // Initialize any pending world data
            if (window.pendingWorldData) {
                hd1Manager.initializeWorld(window.pendingWorldData);
                window.pendingWorldData = null;
                console.log('[HD1] Applied pending world data via fallback');
            }
            
            // Initialize displays via fallback
            initializeDisplays();
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


// Associate WebSocket client with HD1 session for isolation
function associateSession(sessionId) {
    if (ws && ws.readyState === WebSocket.OPEN) {
        const associateMsg = {
            type: 'session_associate',
            session_id: sessionId
        };
        addDebug('WS_SEND', {type: 'session_associate', session: sessionId});
        ws.send(JSON.stringify(associateMsg));
        console.log('[HD1] WebSocket associated with session:', sessionId);
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
                webgl: !!hd1Manager,
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

// Initialize mouse controls when scene is ready
setupMouseControls();

// Event-driven initialization when HD1 manager becomes available
function initializeDisplays() {
    if (hd1Manager && hd1Manager.state) {
        updateEnvironmentalDisplay();
        scheduleObjectsUpdate();
        console.log('[HD1] Displays initialized via event-driven method');
    }
}

// Debug panel collapsible functionality (initialized after cookie functions)
let debugCollapsed = false;

// Scene selection management
const debugSceneSelect = document.getElementById('debug-scene-select');

// Load saved scene from cookie
function loadSavedScene() {
    const savedScene = getCookie('hd1_scene');
    if (savedScene) {
        debugSceneSelect.value = savedScene;
    }
}

// Save scene to cookie
function saveScene(sceneId) {
    setCookie('hd1_scene', sceneId, 30); // 30 days
    console.log(`[HD1] Scene saved to cookie: "${sceneId}"`);
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
            
            // Update displays immediately after scene load response
            if (typeof updateEnvironmentalDisplay === 'function') {
                updateEnvironmentalDisplay();
            }
            // Use MutationObserver to update when DOM actually changes
            scheduleObjectsUpdate();
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
            const savedScene = getCookie('hd1_scene');
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
debugCollapsed = getCookie('hd1_console_collapsed') === 'true';
addDebug('CONSOLE_COOKIE', 'Loaded state: ' + (debugCollapsed ? 'collapsed' : 'expanded'));

function setDebugState(collapsed, saveToCookie = true) {
    debugCollapsed = collapsed;
    const debugContent = document.getElementById('debug-content');
    
    if (debugCollapsed) {
        debugCollapseIcon.classList.add('collapsed');
        if (statusCollapseArrow) statusCollapseArrow.classList.add('collapsed');
        if (debugContent) debugContent.classList.add('collapsed');
    } else {
        debugCollapseIcon.classList.remove('collapsed');
        if (statusCollapseArrow) statusCollapseArrow.classList.remove('collapsed');
        if (debugContent) debugContent.classList.remove('collapsed');
    }
    if (saveToCookie) {
        setCookie('hd1_console_collapsed', debugCollapsed.toString(), 30); // 30 days
        addDebug('CONSOLE_SAVE', 'State saved: ' + (debugCollapsed ? 'collapsed' : 'expanded'));
    }
}

debugHeader.addEventListener('click', function(e) {
    setDebugState(!debugCollapsed, true);
});

// Status bar session ID tag click handler for copying
sessionIdTagStatus.addEventListener('click', function(e) {
    e.stopPropagation(); // Prevent any propagation
    const sessionId = currentSessionId;
    if (sessionId) {
        navigator.clipboard.writeText(sessionId).then(() => {
            // Visual feedback
            sessionIdTagStatus.classList.add('copied');
            setTimeout(() => {
                sessionIdTagStatus.classList.remove('copied');
            }, 500);
            addDebug('COPY_SESSION', {id: sessionId});
        }).catch(err => {
            // Fallback for older browsers
            const textArea = document.createElement('textarea');
            textArea.value = sessionId;
            document.body.appendChild(textArea);
            textArea.select();
            document.execCommand('copy');
            document.body.removeChild(textArea);
            addDebug('COPY_SESSION_FALLBACK', {id: sessionId});
        });
    }
});

// Status bar click handler for toggle
statusBar.addEventListener('click', function(e) {
    // Don't toggle if clicking on session ID tag (it has its own handler)
    if (e.target !== sessionIdTagStatus) {
        setDebugState(!debugCollapsed, true);
    }
});

// Load version information dynamically
async function loadVersionInfo() {
    try {
        const response = await fetch('/api/version');
        if (response.ok) {
            const versionData = await response.json();
            
            // Update API version
            const apiVersionElement = document.getElementById('api-version');
            if (apiVersionElement) {
                apiVersionElement.textContent = 'v' + versionData.api_version;
            }
            
            // Update JS version (first 8 characters)
            const jsVersionElement = document.getElementById('js-version');
            if (jsVersionElement) {
                jsVersionElement.textContent = versionData.js_version.substring(0, 8);
            }
            
            addDebug('VERSION_LOADED', {
                api: versionData.api_version, 
                js: versionData.js_version.substring(0, 8),
                title: versionData.title
            });
        } else {
            addDebug('VERSION_ERROR', {status: response.status});
        }
    } catch (error) {
        addDebug('VERSION_FAIL', {error: error.message});
        // Keep fallback values if version endpoint fails
    }
}

// Initialize debug panel state from cookie (don't save back to cookie)
setDebugState(debugCollapsed, false);

// Load version information
loadVersionInfo();

// Initialize hardware monitoring and objects observer
initializeHardwareMonitoring();
initializeObjectsObserver();

// Start connection
connectWebSocket();
