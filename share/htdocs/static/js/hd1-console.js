const canvas = document.getElementById('holodeck-canvas');
const debugLog = document.getElementById('debug-log');
const debugHeader = document.getElementById('debug-header');
const debugCollapseIcon = document.getElementById('debug-collapse-icon');

// Status bar elements
const statusConnectionIndicator = document.getElementById('status-connection-indicator');
const statusConnectionText = document.getElementById('status-connection-text');
const statusLockIndicator = document.getElementById('status-lock-indicator');
const sessionIdTagStatus = document.getElementById('session-id-tag-status');

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
                    setTimeout(() => {
                        const savedScene = getCookie('hd1_scene');
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
                setTimeout(() => {
                    const savedScene = getCookie('hd1_scene');
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

// WebSocket connection with intelligent rebootstrap
function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;
    
    addDebug('WS_CONNECT', {url: wsUrl, attempt: reconnectAttempts + 1});
    setStatus('connecting');
    
    ws = new WebSocket(wsUrl);
    
    ws.onopen = function() {
        addDebug('WS_OPEN', 'WebSocket connected');
        setStatus('connected');
        reconnectAttempts = 0;
        lastMessageTime = Date.now();
        
        // Clear any pending reconnect timeout
        if (reconnectTimeout) {
            clearTimeout(reconnectTimeout);
            reconnectTimeout = null;
        }
    };
    
    ws.onmessage = function(event) {
        setStatus('receiving');
        lastMessageTime = Date.now();
        
        try {
            const data = JSON.parse(event.data);
            addDebug('WS_MSG', data);
            
            // Handle different message types
            if (data.type === 'session_update') {
                updateDebugSession(data.session_id);
            }
        } catch (error) {
            addDebug('WS_ERROR', 'Failed to parse message: ' + error.message);
        }
        
        // Reset status after brief delay
        setTimeout(() => setStatus('connected'), 200);
    };
    
    ws.onclose = function(event) {
        addDebug('WS_CLOSE', {code: event.code, reason: event.reason});
        setStatus('disconnected');
        
        // Attempt reconnection with exponential backoff
        reconnectAttempts++;
        
        if (reconnectAttempts >= maxReconnectAttempts) {
            addDebug('REBOOTSTRAP', `Max reconnect attempts (${maxReconnectAttempts}) reached - triggering rebootstrap`);
            triggerRebootstrap();
            return;
        }
        
        const delay = Math.min(1000 * Math.pow(2, reconnectAttempts), 30000);
        addDebug('WS_RECONNECT', {attempt: reconnectAttempts, delay: delay});
        
        reconnectTimeout = setTimeout(() => {
            connectWebSocket();
        }, delay);
    };
    
    ws.onerror = function(error) {
        addDebug('WS_ERROR', 'WebSocket error: ' + error.message);
        setStatus('error');
    };
}

// Rebootstrap functionality - clears storage and reloads page
function triggerRebootstrap() {
    addDebug('REBOOTSTRAP', 'Clearing storage and reloading page...');
    
    // Clear all browser storage
    localStorage.clear();
    sessionStorage.clear();
    
    // Clear cookies
    document.cookie.split(";").forEach(function(c) { 
        document.cookie = c.replace(/^ +/, "").replace(/=.*/, "=;expires=" + new Date().toUTCString() + ";path=/"); 
    });
    
    // Clear any cached data
    if ('caches' in window) {
        caches.keys().then(names => {
            names.forEach(name => {
                caches.delete(name);
            });
        });
    }
    
    // Force page reload after small delay
    setTimeout(() => {
        window.location.reload(true);
    }, 1000);
}

// Scene loading
function loadScene(sceneName) {
    addDebug('SCENE_LOAD', {scene: sceneName});
    
    if (sceneName === 'empty') {
        // Clear scene to empty grid
        const objects = document.getElementById('holodeck-objects');
        if (objects) {
            objects.innerHTML = '';
        }
    } else if (sceneName === 'basic-shapes') {
        // Load basic shapes scene
        const objects = document.getElementById('holodeck-objects');
        if (objects) {
            objects.innerHTML = `
                <a-box position="-1 0.5 -3" rotation="0 45 0" color="#4CC3D9"></a-box>
                <a-sphere position="0 1.25 -5" radius="1.25" color="#EF2D5E"></a-sphere>
                <a-cylinder position="1 0.75 -3" radius="0.5" height="1.5" color="#FFC65D"></a-cylinder>
                <a-plane position="0 0 -4" rotation="-90 0 0" width="4" height="4" color="#7BC8A4"></a-plane>
            `;
        }
    }
    
    // Save scene selection
    setCookie('hd1_scene', sceneName);
}

// Console panel collapse functionality
let debugCollapsed = false;

debugHeader.addEventListener('click', function() {
    debugCollapsed = !debugCollapsed;
    const debugContent = document.getElementById('debug-content');
    
    if (debugCollapsed) {
        debugContent.classList.add('collapsed');
        debugCollapseIcon.classList.add('collapsed');
        debugCollapseIcon.innerHTML = '&#8681;'; // Down arrow
    } else {
        debugContent.classList.remove('collapsed');
        debugCollapseIcon.classList.remove('collapsed');
        debugCollapseIcon.innerHTML = '&#8679;'; // Up arrow
    }
    
    addDebug('CONSOLE_TOGGLE', {collapsed: debugCollapsed});
});

// Session ID click to copy
sessionIdTagStatus.addEventListener('click', function() {
    if (navigator.clipboard && sessionIdTagStatus.textContent) {
        navigator.clipboard.writeText(sessionIdTagStatus.textContent).then(() => {
            addDebug('SESSION_COPY', sessionIdTagStatus.textContent);
        });
    }
});

// Scene selector
const debugSceneSelect = document.getElementById('debug-scene-select');
debugSceneSelect.addEventListener('change', function() {
    const selectedScene = this.value;
    if (selectedScene) {
        loadScene(selectedScene);
    }
});

// Rebootstrap button
const rebootstrapBtn = document.getElementById('rebootstrap-btn');
rebootstrapBtn.addEventListener('click', function() {
    if (confirm('This will clear all storage and reload the page. Continue?')) {
        triggerRebootstrap();
    }
});

// Initialize console on page load
document.addEventListener('DOMContentLoaded', function() {
    addDebug('INIT', 'HD1 Console initializing...');
    
    // Initialize WebSocket connection
    connectWebSocket();
    
    // Restore saved scene selection
    const savedScene = getCookie('hd1_scene');
    if (savedScene && debugSceneSelect) {
        debugSceneSelect.value = savedScene;
    }
    
    addDebug('READY', 'HD1 Console ready');
});

// Initialize immediately if DOM already loaded
if (document.readyState === 'loading') {
    // Wait for DOMContentLoaded
} else {
    // DOM already loaded
    addDebug('INIT', 'HD1 Console initializing...');
    connectWebSocket();
    
    const savedScene = getCookie('hd1_scene');
    if (savedScene && debugSceneSelect) {
        debugSceneSelect.value = savedScene;
    }
    
    addDebug('READY', 'HD1 Console ready');
}