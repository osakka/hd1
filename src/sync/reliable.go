// Package sync implements HD1's TCP-simple reliable synchronization protocol
// Based on sequence numbers - the same principle that makes TCP bulletproof
package sync

import (
	"sync"
	"time"
	
	"holodeck1/logging"
)

// Operation represents a single synchronized operation with sequence number
type Operation struct {
	SeqNum    uint64                 `json:"seq_num"`    // Global sequence number
	ClientID  string                 `json:"client_id"`  // Who sent it
	Type      string                 `json:"type"`       // "avatar_move", "entity_create", etc.
	Data      map[string]interface{} `json:"data"`       // The actual change
	Timestamp time.Time              `json:"timestamp"`  // When it happened
}

// ReliableSync implements TCP-simple synchronization using sequence numbers
type ReliableSync struct {
	// Core state
	nextSeqNum     uint64
	operations     map[uint64]*Operation
	mutex          sync.RWMutex
	
	// Per-client tracking
	clientLastSeen map[string]uint64
	clients        map[string]chan *Operation
	
	// Cleanup
	maxOperations  int
	cleanupCounter uint64
}

// NewReliableSync creates a new TCP-simple sync system
func NewReliableSync() *ReliableSync {
	return &ReliableSync{
		nextSeqNum:     1,
		operations:     make(map[uint64]*Operation),
		clientLastSeen: make(map[string]uint64),
		clients:        make(map[string]chan *Operation),
		maxOperations:  100000, // Keep last 100k operations
		cleanupCounter: 0,
	}
}

// SubmitOperation adds an operation to the global sequence
func (rs *ReliableSync) SubmitOperation(op *Operation) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	
	// Assign sequence number
	op.SeqNum = rs.nextSeqNum
	op.Timestamp = time.Now()
	rs.nextSeqNum++
	
	// Store operation
	rs.operations[op.SeqNum] = op
	
	logging.Debug("operation submitted", map[string]interface{}{
		"seq_num":   op.SeqNum,
		"client_id": op.ClientID,
		"type":      op.Type,
	})
	
	// Broadcast to all clients
	rs.broadcastOperation(op)
	
	// Periodic cleanup
	rs.cleanupCounter++
	if rs.cleanupCounter%1000 == 0 {
		rs.cleanup()
	}
}

// RegisterClient registers a new client for synchronization
func (rs *ReliableSync) RegisterClient(clientID string) chan *Operation {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	
	// Create client channel
	clientChan := make(chan *Operation, 1000)
	rs.clients[clientID] = clientChan
	rs.clientLastSeen[clientID] = 0
	
	logging.Info("client registered", map[string]interface{}{
		"client_id": clientID,
	})
	
	return clientChan
}

// UnregisterClient removes a client
func (rs *ReliableSync) UnregisterClient(clientID string) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	
	if clientChan, exists := rs.clients[clientID]; exists {
		close(clientChan)
		delete(rs.clients, clientID)
		delete(rs.clientLastSeen, clientID)
		
		logging.Info("client unregistered", map[string]interface{}{
			"client_id": clientID,
		})
	}
}

// GetCurrentSequence returns the current sequence number
func (rs *ReliableSync) GetCurrentSequence() uint64 {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	return rs.nextSeqNum - 1
}

// GetMissingOperations returns operations from 'from' to 'to' (inclusive)
func (rs *ReliableSync) GetMissingOperations(from, to uint64) []*Operation {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	
	var missing []*Operation
	for seq := from; seq <= to && seq < rs.nextSeqNum; seq++ {
		if op, exists := rs.operations[seq]; exists {
			missing = append(missing, op)
		}
	}
	
	return missing
}

// GetAllOperations returns all operations for new client sync
func (rs *ReliableSync) GetAllOperations() []*Operation {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	
	var allOps []*Operation
	for seq := uint64(1); seq < rs.nextSeqNum; seq++ {
		if op, exists := rs.operations[seq]; exists {
			allOps = append(allOps, op)
		}
	}
	
	return allOps
}

// UpdateClientLastSeen updates the last seen sequence for a client
func (rs *ReliableSync) UpdateClientLastSeen(clientID string, seqNum uint64) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	
	if lastSeen, exists := rs.clientLastSeen[clientID]; exists {
		if seqNum > lastSeen {
			rs.clientLastSeen[clientID] = seqNum
		}
	}
}

// GetClientLastSeen returns the last seen sequence for a client
func (rs *ReliableSync) GetClientLastSeen(clientID string) uint64 {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	
	return rs.clientLastSeen[clientID]
}

// GetCurrentSequence - REMOVED: Duplicate method, already exists above

// broadcastOperation sends operation to all connected clients
func (rs *ReliableSync) broadcastOperation(op *Operation) {
	for clientID, clientChan := range rs.clients {
		select {
		case clientChan <- op:
			// Successfully sent
		default:
			// Client channel full - skip this client
			logging.Warn("client channel full", map[string]interface{}{
				"client_id": clientID,
				"seq_num":   op.SeqNum,
			})
		}
	}
}

// cleanup removes old operations to prevent memory growth
func (rs *ReliableSync) cleanup() {
	if len(rs.operations) <= rs.maxOperations {
		return
	}
	
	// Find minimum last seen sequence across all clients
	minLastSeen := rs.nextSeqNum
	for _, lastSeen := range rs.clientLastSeen {
		if lastSeen < minLastSeen {
			minLastSeen = lastSeen
		}
	}
	
	// Keep operations after (minLastSeen - 1000) to provide buffer
	keepAfter := minLastSeen - 1000
	if keepAfter < 1 {
		keepAfter = 1
	}
	
	// Remove old operations
	removed := 0
	for seq := range rs.operations {
		if seq < keepAfter {
			delete(rs.operations, seq)
			removed++
		}
	}
	
	logging.Info("operations cleaned up", map[string]interface{}{
		"removed":    removed,
		"remaining":  len(rs.operations),
		"keep_after": keepAfter,
	})
}

// GetStats returns synchronization statistics
func (rs *ReliableSync) GetStats() map[string]interface{} {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	
	return map[string]interface{}{
		"next_sequence":    rs.nextSeqNum,
		"stored_operations": len(rs.operations),
		"connected_clients": len(rs.clients),
		"max_operations":   rs.maxOperations,
	}
}

// GetPendingOperations returns operations that need to be broadcast
func (rs *ReliableSync) GetPendingOperations() []*Operation {
	// For this simple implementation, we don't queue pending operations
	// Operations are broadcast immediately when submitted
	return []*Operation{}
}

// GetOperationsInRange returns operations within a sequence range
func (rs *ReliableSync) GetOperationsInRange(from, to uint64) []*Operation {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	
	var operations []*Operation
	for seq := from; seq <= to; seq++ {
		if op, exists := rs.operations[seq]; exists {
			operations = append(operations, op)
		}
	}
	return operations
}