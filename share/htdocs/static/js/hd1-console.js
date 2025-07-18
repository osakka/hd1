// HD1 Three.js Console - Minimal Implementation
const canvas = document.getElementById('holodeck-canvas');
const debugLog = document.getElementById('debug-log');
const debugHeader = document.getElementById('debug-header');
const debugCollapseIcon = document.getElementById('debug-collapse-icon');

// Status bar elements
const statusConnectionIndicator = document.getElementById('status-connection-indicator');
const statusConnectionText = document.getElementById('status-connection-text');
const statusLockIndicator = document.getElementById('status-lock-indicator');
const statusMouselookIndicator = document.getElementById('status-mouselook-indicator');

let ws;
let reconnectAttempts = 0;
let maxReconnectAttempts = 99;
let reconnectTimeout;
let hd1Id = null;
let apiClient = null;

// Status management
function setStatus(status, message) {
    // Update connection text and indicator
    switch(status) {
        case 'connecting':
            statusConnectionText.textContent = 'Connecting';
            statusConnectionIndicator.className = 'connecting';
            break;
        case 'connected':
            statusConnectionText.textContent = 'Connected';
            statusConnectionIndicator.className = 'connected';
            break;
        case 'disconnected':
            statusConnectionText.textContent = 'Disconnected';
            statusConnectionIndicator.className = 'disconnected';
            break;
        case 'error':
            statusConnectionText.textContent = 'Error';
            statusConnectionIndicator.className = 'disconnected';
            break;
        case 'receiving':
            statusConnectionText.textContent = 'Receiving';
            statusConnectionIndicator.className = 'receiving';
            break;
        default:
            statusConnectionText.textContent = 'Unknown';
            statusConnectionIndicator.className = 'disconnected';
    }
}

// Mouse look status management
function setMouselookStatus(active) {
    if (statusMouselookIndicator) {
        statusMouselookIndicator.className = active ? 'active' : 'off';
    }
}

// Make globally available
window.setMouselookStatus = setMouselookStatus;

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
    window.ws = ws; // Make globally available
    
    ws.onopen = function() {
        addDebug('WS_OPEN', 'WebSocket connected');
        setStatus('connected');
        reconnectAttempts = 0;
        
        if (reconnectTimeout) {
            clearTimeout(reconnectTimeout);
            reconnectTimeout = null;
        }
        
        // Send existing hd1_id for reconnection if we have one
        if (hd1Id) {
            const reconnectMsg = {
                type: 'client_reconnect',
                hd1_id: hd1Id
            };
            ws.send(JSON.stringify(reconnectMsg));
            addDebug('CLIENT_RECONNECT', 'Sent existing hd1_id: ' + hd1Id);
        }
    };
    
    ws.onmessage = function(event) {
        setStatus('receiving');
        
        try {
            const data = JSON.parse(event.data);
            addDebug('WS_MSG', data);
            
            // Handle client initialization from server
            if (data.type === 'client_init' && data.hd1_id) {
                hd1Id = data.hd1_id;
                window.hd1Id = hd1Id; // Make globally available
                
                // Update API client with server-provided hd1_id
                if (apiClient) {
                    apiClient.setHd1Id(hd1Id);
                }
                
                updateRebootstrapButton();
                addDebug('CLIENT_INIT', 'Server-provided hd1_id: ' + hd1Id);
                
                // Client is already initialized with proper hd1_id - no need to join separate session
                addDebug('HD1_READY', 'Client initialized with unified HD1 ID: ' + hd1Id);
                
                // Request full sync to get all existing operations
                requestFullSync();
            }
            
            // Handle successful client reconnection
            if (data.type === 'client_reconnect_success' && data.hd1_id) {
                hd1Id = data.hd1_id;
                window.hd1Id = hd1Id; // Make globally available
                
                // Update API client with reconnected hd1_id
                if (apiClient) {
                    apiClient.setHd1Id(hd1Id);
                }
                
                updateRebootstrapButton();
                addDebug('CLIENT_RECONNECT_SUCCESS', 'Reconnected with hd1_id: ' + hd1Id);
            }
            
            // Handle sync operations from server
            if (data.type === 'sync_operation' && data.operation) {
                // Forward sync operation to Three.js scene manager
                if (window.hd1ThreeJS) {
                    window.hd1ThreeJS.handleSyncOperation(data.operation);
                    addDebug('SYNC_OP', 'Applied operation: ' + data.operation.type + ' seq:' + data.operation.seq_num);
                } else {
                    addDebug('SYNC_OP_ERROR', 'Three.js scene manager not available');
                }
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

// Update rebootstrap button with current hd1_id
function updateRebootstrapButton() {
    const btn = document.getElementById('rebootstrap-btn');
    if (hd1Id) {
        btn.textContent = hd1Id;
        btn.title = `Rebootstrap: HD1 ${hd1Id} - Clear storage and reload page`;
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

// Request full sync from server
async function requestFullSync() {
    if (!apiClient) {
        addDebug('BOOTSTRAP_ERROR', 'API client not available');
        return;
    }
    
    try {
        addDebug('BOOTSTRAP_START', 'Requesting full sync...');
        const response = await apiClient.getFullSync();
        
        if (response.success && response.operations) {
            addDebug('BOOTSTRAP_SUCCESS', `Received ${response.operations.length} operations`);
            
            // Apply all operations to Three.js scene
            if (window.hd1ThreeJS) {
                for (const opWrapper of response.operations) {
                    window.hd1ThreeJS.handleSyncOperation(opWrapper.operation);
                }
                addDebug('BOOTSTRAP_APPLIED', `Applied ${response.operations.length} operations to scene`);
            } else {
                addDebug('BOOTSTRAP_ERROR', 'Three.js scene manager not available');
            }
        } else {
            addDebug('BOOTSTRAP_ERROR', 'Failed to get full sync: ' + JSON.stringify(response));
        }
    } catch (error) {
        addDebug('BOOTSTRAP_ERROR', 'Full sync request failed: ' + error.message);
    }
}

// Initialize console
function initConsole() {
    addDebug('INIT', 'HD1 Three.js Console initializing...');
    
    // Initialize API client without client ID (server will provide it via WebSocket)
    if (window.HD1ThreeJSAPIClient) {
        apiClient = new window.HD1ThreeJSAPIClient();
        window.apiClient = apiClient; // Make globally available for Three.js
        addDebug('API_CLIENT', 'API client initialized and made globally available');
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