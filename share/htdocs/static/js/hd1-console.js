// HD1 Three.js Console - Minimal Implementation
const canvas = document.getElementById('holodeck-canvas');
const debugLog = document.getElementById('debug-log');
const debugHeader = document.getElementById('debug-header');
const debugCollapseIcon = document.getElementById('debug-collapse-icon');

// Status bar elements
const statusConnectionIndicator = document.getElementById('status-connection-indicator');
const statusConnectionText = document.getElementById('status-connection-text');
const statusLockIndicator = document.getElementById('status-lock-indicator');

let ws;
let reconnectAttempts = 0;
let maxReconnectAttempts = 99;
let reconnectTimeout;
let clientId = null;
let sessionId = null;
let avatarId = null;

// Status management
function setStatus(status, message) {
    // Update connection text only (removed connection indicator)
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

// WebSocket connection with rebootstrap
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
        
        if (reconnectTimeout) {
            clearTimeout(reconnectTimeout);
            reconnectTimeout = null;
        }
    };
    
    ws.onmessage = function(event) {
        setStatus('receiving');
        
        try {
            const data = JSON.parse(event.data);
            addDebug('WS_MSG', data);
            
            // Capture client/session/avatar IDs from server messages
            if (data.client_id) {
                clientId = data.client_id;
                updateRebootstrapButton();
            }
            if (data.session_id) {
                sessionId = data.session_id;
                updateRebootstrapButton();
            }
            if (data.avatar_id) {
                avatarId = data.avatar_id;
                updateRebootstrapButton();
            }
        } catch (error) {
            addDebug('WS_ERROR', 'Failed to parse message: ' + error.message);
        }
        
        setTimeout(() => setStatus('connected'), 200);
    };
    
    ws.onclose = function(event) {
        addDebug('WS_CLOSE', {code: event.code, reason: event.reason});
        setStatus('disconnected');
        
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

// Rebootstrap functionality
function triggerRebootstrap() {
    addDebug('REBOOTSTRAP', 'Clearing storage and reloading page...');
    
    localStorage.clear();
    sessionStorage.clear();
    
    document.cookie.split(";").forEach(function(c) { 
        document.cookie = c.replace(/^ +/, "").replace(/=.*/, "=;expires=" + new Date().toUTCString() + ";path=/"); 
    });
    
    if ('caches' in window) {
        caches.keys().then(names => {
            names.forEach(name => {
                caches.delete(name);
            });
        });
    }
    
    setTimeout(() => {
        window.location.reload(true);
    }, 1000);
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

// Update rebootstrap button with current session/avatar/client ID
function updateRebootstrapButton() {
    const btn = document.getElementById('rebootstrap-btn');
    if (avatarId) {
        btn.textContent = avatarId;
        btn.title = `Rebootstrap: Avatar ${avatarId} - Clear storage and reload page`;
    } else if (sessionId) {
        btn.textContent = sessionId;
        btn.title = `Rebootstrap: Session ${sessionId} - Clear storage and reload page`;
    } else if (clientId) {
        btn.textContent = clientId;
        btn.title = `Rebootstrap: Client ${clientId} - Clear storage and reload page`;
    } else {
        btn.textContent = 'REBOOTSTRAP';
        btn.title = 'Rebootstrap: Clear storage and reload page';
    }
}

// Rebootstrap button
const rebootstrapBtn = document.getElementById('rebootstrap-btn');
rebootstrapBtn.addEventListener('click', function() {
    if (confirm('This will clear all storage and reload the page. Continue?')) {
        triggerRebootstrap();
    }
});

// Initialize console
function initConsole() {
    addDebug('INIT', 'HD1 Three.js Console initializing...');
    
    // Initialize API client and get client ID
    if (window.HD1ThreeJSAPIClient) {
        const apiClient = new window.HD1ThreeJSAPIClient();
        clientId = apiClient.clientId;
        updateRebootstrapButton();
        addDebug('CLIENT_ID', clientId);
    }
    
    connectWebSocket();
    addDebug('READY', 'HD1 Three.js Console ready');
}

// Start console when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initConsole);
} else {
    initConsole();
}