/**
 * HD1 TCP-Simple Client Synchronization
 * Based on sequence numbers - the same principle that makes TCP bulletproof
 */

class HD1Sync {
    constructor(websocket, threeJS) {
        this.ws = websocket;
        this.threeJS = threeJS;
        
        // TCP-simple state
        this.lastSeenSeq = 0;
        this.operationBuffer = [];
        this.clientId = this.generateClientId();
        
        // Connection state
        this.connected = false;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 10;
        this.reconnectDelay = 1000;
        
        // Statistics
        this.stats = {
            operationsReceived: 0,
            operationsApplied: 0,
            operationsSent: 0,
            missingRequests: 0,
            reconnections: 0
        };
        
        this.setupWebSocket();
        console.log('[HD1-Sync] TCP-simple sync initialized, client:', this.clientId);
    }
    
    generateClientId() {
        return 'client-' + Math.random().toString(36).substr(2, 9) + '-' + Date.now();
    }
    
    setupWebSocket() {
        this.ws.onopen = () => {
            console.log('[HD1-Sync] WebSocket connected');
            this.connected = true;
            this.reconnectAttempts = 0;
            
            // Request full synchronization
            this.requestFullSync();
        };
        
        this.ws.onmessage = (event) => {
            this.handleMessage(JSON.parse(event.data));
        };
        
        this.ws.onclose = () => {
            console.log('[HD1-Sync] WebSocket disconnected');
            this.connected = false;
            this.attemptReconnect();
        };
        
        this.ws.onerror = (error) => {
            console.error('[HD1-Sync] WebSocket error:', error);
        };
    }
    
    handleMessage(data) {
        switch (data.type) {
            case 'operation':
                this.handleOperation(data);
                break;
            case 'missing_operations':
                this.handleMissingOperations(data);
                break;
            case 'full_sync':
                this.handleFullSync(data);
                break;
            case 'sync_complete':
                this.handleSyncComplete(data);
                break;
            default:
                console.warn('[HD1-Sync] Unknown message type:', data.type);
        }
    }
    
    handleOperation(data) {
        this.stats.operationsReceived++;
        
        const operation = data.operation;
        const seqNum = data.seq_num;
        
        // Is this the next expected sequence?
        if (seqNum === this.lastSeenSeq + 1) {
            // Perfect! Apply immediately
            this.applyOperation(operation);
            this.lastSeenSeq = seqNum;
            
            // Process any buffered operations that are now ready
            this.processBuffer();
        } else if (seqNum > this.lastSeenSeq + 1) {
            // Future operation - buffer it
            this.bufferOperation(data);
            
            // Request missing operations
            this.requestMissing(this.lastSeenSeq + 1, seqNum - 1);
        } else {
            // Old operation - ignore (already processed)
            console.debug('[HD1-Sync] Ignoring old operation:', seqNum, 'last seen:', this.lastSeenSeq);
        }
    }
    
    handleMissingOperations(data) {
        console.log('[HD1-Sync] Received missing operations:', data.operations.length);
        
        // Sort operations by sequence number
        const operations = data.operations.sort((a, b) => a.seq_num - b.seq_num);
        
        // Apply each operation in order
        for (const opData of operations) {
            if (opData.seq_num === this.lastSeenSeq + 1) {
                this.applyOperation(opData.operation);
                this.lastSeenSeq = opData.seq_num;
            } else if (opData.seq_num > this.lastSeenSeq + 1) {
                // Still missing some operations
                this.bufferOperation(opData);
            }
        }
        
        // Process buffer after applying missing operations
        this.processBuffer();
    }
    
    handleFullSync(data) {
        console.log('[HD1-Sync] Full sync received:', data.operations.length, 'operations');
        
        // Clear current state
        this.lastSeenSeq = 0;
        this.operationBuffer = [];
        
        // Apply all operations in sequence
        const operations = data.operations.sort((a, b) => a.seq_num - b.seq_num);
        
        for (const opData of operations) {
            if (opData.seq_num === this.lastSeenSeq + 1) {
                this.applyOperation(opData.operation);
                this.lastSeenSeq = opData.seq_num;
            } else {
                console.error('[HD1-Sync] Sequence gap in full sync at:', opData.seq_num);
            }
        }
        
        console.log('[HD1-Sync] Full sync complete, last sequence:', this.lastSeenSeq);
    }
    
    handleSyncComplete(data) {
        console.log('[HD1-Sync] Sync complete, current sequence:', data.current_sequence);
    }
    
    bufferOperation(data) {
        // Add to buffer, maintaining sort order
        const insertIndex = this.operationBuffer.findIndex(op => op.seq_num > data.seq_num);
        if (insertIndex === -1) {
            this.operationBuffer.push(data);
        } else {
            this.operationBuffer.splice(insertIndex, 0, data);
        }
    }
    
    processBuffer() {
        let applied = 0;
        
        // Process buffered operations that are now ready
        while (this.operationBuffer.length > 0) {
            const next = this.operationBuffer[0];
            
            if (next.seq_num === this.lastSeenSeq + 1) {
                // This operation is ready to apply
                this.operationBuffer.shift();
                this.applyOperation(next.operation);
                this.lastSeenSeq = next.seq_num;
                applied++;
            } else {
                // Not ready yet
                break;
            }
        }
        
        if (applied > 0) {
            console.log('[HD1-Sync] Applied', applied, 'buffered operations');
        }
    }
    
    applyOperation(operation) {
        try {
            this.threeJS.applyOperation(operation);
            this.stats.operationsApplied++;
            
            console.debug('[HD1-Sync] Applied operation:', operation.type, 'seq:', this.lastSeenSeq);
        } catch (error) {
            console.error('[HD1-Sync] Error applying operation:', error, operation);
        }
    }
    
    requestMissing(from, to) {
        if (!this.connected) return;
        
        this.stats.missingRequests++;
        
        const message = {
            type: 'request_missing',
            client_id: this.clientId,
            from: from,
            to: to
        };
        
        this.ws.send(JSON.stringify(message));
        console.log('[HD1-Sync] Requested missing operations:', from, 'to', to);
    }
    
    requestFullSync() {
        if (!this.connected) return;
        
        const message = {
            type: 'request_full_sync',
            client_id: this.clientId
        };
        
        this.ws.send(JSON.stringify(message));
        console.log('[HD1-Sync] Requested full sync');
    }
    
    sendOperation(operation) {
        if (!this.connected) {
            console.warn('[HD1-Sync] Cannot send operation - not connected');
            return;
        }
        
        const message = {
            type: 'operation',
            client_id: this.clientId,
            operation: operation
        };
        
        this.ws.send(JSON.stringify(message));
        this.stats.operationsSent++;
        
        console.debug('[HD1-Sync] Sent operation:', operation.type);
    }
    
    attemptReconnect() {
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            console.error('[HD1-Sync] Max reconnect attempts reached');
            return;
        }
        
        this.reconnectAttempts++;
        this.stats.reconnections++;
        
        const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);
        console.log(`[HD1-Sync] Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`);
        
        setTimeout(() => {
            if (!this.connected) {
                console.log('[HD1-Sync] Attempting to reconnect...');
                this.ws = new WebSocket(this.ws.url);
                this.setupWebSocket();
            }
        }, delay);
    }
    
    // Public API for sending operations
    moveAvatar(sessionId, position, rotation) {
        this.sendOperation({
            type: 'avatar_move',
            data: {
                session_id: sessionId,
                position: position,
                rotation: rotation
            }
        });
    }
    
    createEntity(id, geometry, material, transform) {
        this.sendOperation({
            type: 'entity_create',
            data: {
                id: id,
                geometry: geometry,
                material: material,
                position: transform.position,
                rotation: transform.rotation,
                scale: transform.scale
            }
        });
    }
    
    updateEntity(id, updates) {
        this.sendOperation({
            type: 'entity_update',
            data: {
                id: id,
                ...updates
            }
        });
    }
    
    deleteEntity(id) {
        this.sendOperation({
            type: 'entity_delete',
            data: {
                id: id
            }
        });
    }
    
    updateScene(updates) {
        this.sendOperation({
            type: 'scene_update',
            data: updates
        });
    }
    
    // Diagnostics
    getStats() {
        return {
            ...this.stats,
            lastSeenSeq: this.lastSeenSeq,
            bufferedOperations: this.operationBuffer.length,
            connected: this.connected,
            clientId: this.clientId
        };
    }
    
    getBufferedOperations() {
        return this.operationBuffer.map(op => ({
            seq_num: op.seq_num,
            type: op.operation.type
        }));
    }
    
    // Cleanup
    disconnect() {
        if (this.ws) {
            this.ws.close();
        }
        this.connected = false;
        console.log('[HD1-Sync] Disconnected');
    }
}

// Export for module systems
if (typeof module !== 'undefined' && module.exports) {
    module.exports = HD1Sync;
}