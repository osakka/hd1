package ot

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"holodeck1/database"
	"holodeck1/logging"
)

type Manager struct {
	db        *database.DB
	documents map[uuid.UUID]*Document
	mu        sync.RWMutex
}

type Document struct {
	ID          uuid.UUID         `json:"id"`
	SessionID   uuid.UUID         `json:"session_id"`
	Content     string            `json:"content"`
	Version     int64             `json:"version"`
	Operations  []Operation       `json:"operations"`
	Clients     map[uuid.UUID]*Client `json:"clients"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	mu          sync.RWMutex
}

type Client struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Version   int64     `json:"version"`
	State     string    `json:"state"`
	Buffer    []Operation `json:"buffer"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Operation struct {
	ID          uuid.UUID   `json:"id"`
	Type        string      `json:"type"`
	Position    int         `json:"position"`
	Content     string      `json:"content,omitempty"`
	Length      int         `json:"length,omitempty"`
	Version     int64       `json:"version"`
	ClientID    uuid.UUID   `json:"client_id"`
	UserID      uuid.UUID   `json:"user_id"`
	Timestamp   time.Time   `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type OperationResult struct {
	Success       bool        `json:"success"`
	Document      *Document   `json:"document,omitempty"`
	Operation     *Operation  `json:"operation,omitempty"`
	Transformed   []Operation `json:"transformed,omitempty"`
	Error         string      `json:"error,omitempty"`
}

type CreateDocumentRequest struct {
	SessionID uuid.UUID `json:"session_id"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
}

type JoinDocumentRequest struct {
	DocumentID uuid.UUID `json:"document_id"`
	UserID     uuid.UUID `json:"user_id"`
}

type ApplyOperationRequest struct {
	DocumentID uuid.UUID `json:"document_id"`
	ClientID   uuid.UUID `json:"client_id"`
	Operation  Operation `json:"operation"`
}

func NewManager(db *database.DB) *Manager {
	return &Manager{
		db:        db,
		documents: make(map[uuid.UUID]*Document),
	}
}

func (m *Manager) CreateDocument(ctx context.Context, req *CreateDocumentRequest) (*Document, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	documentID := uuid.New()
	doc := &Document{
		ID:        documentID,
		SessionID: req.SessionID,
		Content:   req.Content,
		Version:   1,
		Operations: []Operation{},
		Clients:   make(map[uuid.UUID]*Client),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Store in database
	query := `
		INSERT INTO ot_documents (id, session_id, content, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := m.db.ExecContext(ctx, query, documentID, req.SessionID, req.Content, doc.Version, doc.CreatedAt, doc.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to store document: %w", err)
	}

	m.documents[documentID] = doc

	logging.Info("created OT document", map[string]interface{}{
		"document_id": documentID,
		"session_id":  req.SessionID,
		"content_len": len(req.Content),
	})

	return doc, nil
}

func (m *Manager) JoinDocument(ctx context.Context, req *JoinDocumentRequest) (*Client, error) {
	m.mu.RLock()
	doc, exists := m.documents[req.DocumentID]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("document not found")
	}

	doc.mu.Lock()
	defer doc.mu.Unlock()

	// Check if client already exists
	for _, client := range doc.Clients {
		if client.UserID == req.UserID {
			return client, nil
		}
	}

	clientID := uuid.New()
	client := &Client{
		ID:        clientID,
		UserID:    req.UserID,
		Version:   doc.Version,
		State:     "synchronized",
		Buffer:    []Operation{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Store in database
	query := `
		INSERT INTO ot_clients (id, document_id, user_id, version, state, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := m.db.ExecContext(ctx, query, clientID, req.DocumentID, req.UserID, client.Version, client.State, client.CreatedAt, client.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to store client: %w", err)
	}

	doc.Clients[clientID] = client

	logging.Info("client joined OT document", map[string]interface{}{
		"hd1_id":   clientID,
		"document_id": req.DocumentID,
		"user_id":     req.UserID,
		"version":     client.Version,
	})

	return client, nil
}

func (m *Manager) ApplyOperation(ctx context.Context, req *ApplyOperationRequest) (*OperationResult, error) {
	m.mu.RLock()
	doc, exists := m.documents[req.DocumentID]
	m.mu.RUnlock()

	if !exists {
		return &OperationResult{
			Success: false,
			Error:   "document not found",
		}, nil
	}

	doc.mu.Lock()
	defer doc.mu.Unlock()

	client, exists := doc.Clients[req.ClientID]
	if !exists {
		return &OperationResult{
			Success: false,
			Error:   "client not found",
		}, nil
	}

	// Assign operation ID and timestamp
	req.Operation.ID = uuid.New()
	req.Operation.Timestamp = time.Now()
	req.Operation.ClientID = req.ClientID

	// Transform operation based on current document state
	transformedOps, err := m.transformOperation(doc, client, req.Operation)
	if err != nil {
		return &OperationResult{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// Apply transformed operations to document
	for _, op := range transformedOps {
		err = m.applyOperationToDocument(doc, op)
		if err != nil {
			return &OperationResult{
				Success: false,
				Error:   err.Error(),
			}, nil
		}
	}

	// Update document version
	doc.Version++
	doc.UpdatedAt = time.Now()

	// Update client state
	client.Version = doc.Version
	client.UpdatedAt = time.Now()

	// Store operation in database
	operationData, _ := json.Marshal(req.Operation)
	query := `
		INSERT INTO ot_operations (id, document_id, client_id, user_id, type, position, content, length, version, data, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err = m.db.ExecContext(ctx, query, req.Operation.ID, req.DocumentID, req.ClientID, req.Operation.UserID, req.Operation.Type, req.Operation.Position, req.Operation.Content, req.Operation.Length, req.Operation.Version, operationData, req.Operation.Timestamp)
	if err != nil {
		logging.Error("failed to store operation", map[string]interface{}{
			"operation_id": req.Operation.ID,
			"error":        err.Error(),
		})
	}

	// Update document in database
	query = `UPDATE ot_documents SET content = $1, version = $2, updated_at = $3 WHERE id = $4`
	_, err = m.db.ExecContext(ctx, query, doc.Content, doc.Version, doc.UpdatedAt, req.DocumentID)
	if err != nil {
		logging.Error("failed to update document", map[string]interface{}{
			"document_id": req.DocumentID,
			"error":       err.Error(),
		})
	}

	// Add to document operations history
	doc.Operations = append(doc.Operations, req.Operation)

	result := &OperationResult{
		Success:     true,
		Document:    doc,
		Operation:   &req.Operation,
		Transformed: transformedOps,
	}

	logging.Debug("applied OT operation", map[string]interface{}{
		"document_id":  req.DocumentID,
		"hd1_id":    req.ClientID,
		"operation_id": req.Operation.ID,
		"type":         req.Operation.Type,
		"version":      doc.Version,
	})

	return result, nil
}

func (m *Manager) transformOperation(doc *Document, client *Client, op Operation) ([]Operation, error) {
	// Simple operational transform implementation
	// In a production system, this would be more sophisticated

	transformed := []Operation{op}
	
	// If client is behind document version, need to transform against intervening operations
	if client.Version < doc.Version {
		// Get operations that happened since client's version
		for _, docOp := range doc.Operations {
			if docOp.Version > client.Version {
				// Transform the operation against this document operation
				transformed = m.transformAgainstOperation(transformed, docOp)
			}
		}
	}

	return transformed, nil
}

func (m *Manager) transformAgainstOperation(ops []Operation, docOp Operation) []Operation {
	result := []Operation{}
	
	for _, op := range ops {
		transformedOp := op
		
		// Transform based on operation types
		switch docOp.Type {
		case "insert":
			if op.Type == "insert" {
				if op.Position >= docOp.Position {
					transformedOp.Position += len(docOp.Content)
				}
			} else if op.Type == "delete" {
				if op.Position >= docOp.Position {
					transformedOp.Position += len(docOp.Content)
				}
			}
		case "delete":
			if op.Type == "insert" {
				if op.Position > docOp.Position {
					transformedOp.Position -= docOp.Length
				}
			} else if op.Type == "delete" {
				if op.Position > docOp.Position {
					transformedOp.Position -= docOp.Length
				} else if op.Position == docOp.Position {
					// Same position, skip this operation
					continue
				}
			}
		}
		
		result = append(result, transformedOp)
	}
	
	return result
}

func (m *Manager) applyOperationToDocument(doc *Document, op Operation) error {
	switch op.Type {
	case "insert":
		if op.Position < 0 || op.Position > len(doc.Content) {
			return fmt.Errorf("invalid insert position: %d", op.Position)
		}
		doc.Content = doc.Content[:op.Position] + op.Content + doc.Content[op.Position:]
	case "delete":
		if op.Position < 0 || op.Position+op.Length > len(doc.Content) {
			return fmt.Errorf("invalid delete range: %d-%d", op.Position, op.Position+op.Length)
		}
		doc.Content = doc.Content[:op.Position] + doc.Content[op.Position+op.Length:]
	case "retain":
		// Retain operations don't change content
	default:
		return fmt.Errorf("unknown operation type: %s", op.Type)
	}
	
	return nil
}

func (m *Manager) GetDocument(ctx context.Context, documentID uuid.UUID) (*Document, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	doc, exists := m.documents[documentID]
	if !exists {
		return nil, fmt.Errorf("document not found")
	}

	return doc, nil
}

func (m *Manager) GetDocumentsBySession(ctx context.Context, sessionID uuid.UUID) ([]*Document, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var documents []*Document
	for _, doc := range m.documents {
		if doc.SessionID == sessionID {
			documents = append(documents, doc)
		}
	}

	return documents, nil
}

func (m *Manager) LeaveDocument(ctx context.Context, documentID uuid.UUID, userID uuid.UUID) error {
	m.mu.RLock()
	doc, exists := m.documents[documentID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("document not found")
	}

	doc.mu.Lock()
	defer doc.mu.Unlock()

	// Find and remove client
	var clientID uuid.UUID
	for id, client := range doc.Clients {
		if client.UserID == userID {
			clientID = id
			break
		}
	}

	if clientID == uuid.Nil {
		return fmt.Errorf("client not found")
	}

	delete(doc.Clients, clientID)

	// Update database
	query := `DELETE FROM ot_clients WHERE id = $1`
	_, err := m.db.ExecContext(ctx, query, clientID)
	if err != nil {
		logging.Error("failed to remove client from database", map[string]interface{}{
			"hd1_id": clientID,
			"error":     err.Error(),
		})
	}

	logging.Info("client left OT document", map[string]interface{}{
		"hd1_id":   clientID,
		"document_id": documentID,
		"user_id":     userID,
	})

	return nil
}

func (m *Manager) GetOperationHistory(ctx context.Context, documentID uuid.UUID, since time.Time) ([]Operation, error) {
	query := `
		SELECT id, type, position, content, length, version, client_id, user_id, created_at, data
		FROM ot_operations
		WHERE document_id = $1 AND created_at > $2
		ORDER BY created_at ASC
	`
	
	rows, err := m.db.QueryContext(ctx, query, documentID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to query operations: %w", err)
	}
	defer rows.Close()

	var operations []Operation
	for rows.Next() {
		var op Operation
		var data []byte
		
		err := rows.Scan(&op.ID, &op.Type, &op.Position, &op.Content, &op.Length, &op.Version, &op.ClientID, &op.UserID, &op.Timestamp, &data)
		if err != nil {
			return nil, fmt.Errorf("failed to scan operation: %w", err)
		}
		
		// Unmarshal metadata
		if len(data) > 0 {
			json.Unmarshal(data, &op.Metadata)
		}
		
		operations = append(operations, op)
	}

	return operations, nil
}

func (m *Manager) CleanupInactiveDocuments(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for docID, doc := range m.documents {
		doc.mu.Lock()
		if len(doc.Clients) == 0 && time.Since(doc.UpdatedAt) > 24*time.Hour {
			delete(m.documents, docID)
			logging.Info("cleaned up inactive OT document", map[string]interface{}{
				"document_id": docID,
			})
		}
		doc.mu.Unlock()
	}
}

func (m *Manager) StartCleanupWorker(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.CleanupInactiveDocuments(ctx)
		}
	}
}