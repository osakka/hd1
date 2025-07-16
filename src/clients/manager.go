package clients

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"holodeck1/database"
	"holodeck1/logging"
)

type Manager struct {
	db      *database.DB
	clients map[uuid.UUID]*Client
	mu      sync.RWMutex
}

type Client struct {
	ID           uuid.UUID              `json:"id"`
	SessionID    uuid.UUID              `json:"session_id"`
	UserID       uuid.UUID              `json:"user_id"`
	Type         string                 `json:"type"`
	Platform     string                 `json:"platform"`
	Version      string                 `json:"version"`
	Capabilities []string               `json:"capabilities"`
	Configuration map[string]interface{} `json:"configuration"`
	Status       string                 `json:"status"`
	LastSeen     time.Time              `json:"last_seen"`
	Metadata     map[string]interface{} `json:"metadata"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type ClientCapability struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type ClientRegistration struct {
	SessionID     uuid.UUID                     `json:"session_id"`
	UserID        uuid.UUID                     `json:"user_id"`
	Type          string                        `json:"type"`
	Platform      string                        `json:"platform"`
	Version       string                        `json:"version"`
	Capabilities  []string                      `json:"capabilities"`
	Configuration map[string]interface{}        `json:"configuration"`
	Metadata      map[string]interface{}        `json:"metadata"`
}

type ClientUpdate struct {
	Status        string                 `json:"status"`
	Configuration map[string]interface{} `json:"configuration"`
	Metadata      map[string]interface{} `json:"metadata"`
}

type ClientMessage struct {
	ID        uuid.UUID              `json:"id"`
	ClientID  uuid.UUID              `json:"client_id"`
	Type      string                 `json:"type"`
	Action    string                 `json:"action"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

type ClientFilter struct {
	SessionID uuid.UUID `json:"session_id"`
	UserID    uuid.UUID `json:"user_id"`
	Type      string    `json:"type"`
	Platform  string    `json:"platform"`
	Status    string    `json:"status"`
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
}

type PlatformAdapter interface {
	GetPlatformInfo() *PlatformInfo
	ValidateConfiguration(config map[string]interface{}) error
	FormatMessage(message *ClientMessage) (interface{}, error)
	ParseMessage(data interface{}) (*ClientMessage, error)
	GetRequiredCapabilities() []string
}

type PlatformInfo struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Capabilities []string `json:"capabilities"`
	Requirements []string `json:"requirements"`
	Metadata     map[string]interface{} `json:"metadata"`
}

func NewManager(db *database.DB) *Manager {
	return &Manager{
		db:      db,
		clients: make(map[uuid.UUID]*Client),
	}
}

func (m *Manager) RegisterClient(ctx context.Context, req *ClientRegistration) (*Client, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	clientID := uuid.New()
	client := &Client{
		ID:            clientID,
		SessionID:     req.SessionID,
		UserID:        req.UserID,
		Type:          req.Type,
		Platform:      req.Platform,
		Version:       req.Version,
		Capabilities:  req.Capabilities,
		Configuration: req.Configuration,
		Status:        "active",
		LastSeen:      time.Now(),
		Metadata:      req.Metadata,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Store in database
	query := `
		INSERT INTO clients (id, session_id, user_id, type, platform, version, capabilities, configuration, status, last_seen, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	
	capabilitiesJSON, _ := json.Marshal(client.Capabilities)
	configJSON, _ := json.Marshal(client.Configuration)
	metadataJSON, _ := json.Marshal(client.Metadata)
	
	_, err := m.db.Conn.ExecContext(ctx, query,
		client.ID, client.SessionID, client.UserID, client.Type, client.Platform,
		client.Version, capabilitiesJSON, configJSON, client.Status,
		client.LastSeen, metadataJSON, client.CreatedAt, client.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to register client: %w", err)
	}

	m.clients[clientID] = client

	logging.Info("registered client", map[string]interface{}{
		"client_id":  clientID,
		"session_id": req.SessionID,
		"user_id":    req.UserID,
		"type":       req.Type,
		"platform":   req.Platform,
		"version":    req.Version,
	})

	return client, nil
}

func (m *Manager) GetClient(ctx context.Context, clientID uuid.UUID) (*Client, error) {
	m.mu.RLock()
	client, exists := m.clients[clientID]
	m.mu.RUnlock()

	if exists {
		return client, nil
	}

	// Load from database
	query := `
		SELECT id, session_id, user_id, type, platform, version, capabilities, configuration, status, last_seen, metadata, created_at, updated_at
		FROM clients
		WHERE id = $1
	`
	
	var dbClient Client
	var capabilitiesJSON, configJSON, metadataJSON []byte
	
	err := m.db.Conn.QueryRowContext(ctx, query, clientID).Scan(
		&dbClient.ID, &dbClient.SessionID, &dbClient.UserID, &dbClient.Type,
		&dbClient.Platform, &dbClient.Version, &capabilitiesJSON, &configJSON,
		&dbClient.Status, &dbClient.LastSeen, &metadataJSON, &dbClient.CreatedAt, &dbClient.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("client not found")
		}
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	// Unmarshal JSON fields
	json.Unmarshal(capabilitiesJSON, &dbClient.Capabilities)
	json.Unmarshal(configJSON, &dbClient.Configuration)
	json.Unmarshal(metadataJSON, &dbClient.Metadata)

	m.mu.Lock()
	m.clients[clientID] = &dbClient
	m.mu.Unlock()

	return &dbClient, nil
}

func (m *Manager) GetClientsBySession(ctx context.Context, filter *ClientFilter) ([]*Client, error) {
	query := `
		SELECT id, session_id, user_id, type, platform, version, capabilities, configuration, status, last_seen, metadata, created_at, updated_at
		FROM clients
		WHERE session_id = $1
	`
	args := []interface{}{filter.SessionID}
	argIndex := 2

	if filter.UserID != uuid.Nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, filter.UserID)
		argIndex++
	}

	if filter.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, filter.Type)
		argIndex++
	}

	if filter.Platform != "" {
		query += fmt.Sprintf(" AND platform = $%d", argIndex)
		args = append(args, filter.Platform)
		argIndex++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	query += " ORDER BY last_seen DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	rows, err := m.db.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query clients: %w", err)
	}
	defer rows.Close()

	var clients []*Client
	for rows.Next() {
		var client Client
		var capabilitiesJSON, configJSON, metadataJSON []byte
		
		err := rows.Scan(
			&client.ID, &client.SessionID, &client.UserID, &client.Type,
			&client.Platform, &client.Version, &capabilitiesJSON, &configJSON,
			&client.Status, &client.LastSeen, &metadataJSON, &client.CreatedAt, &client.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan client: %w", err)
		}

		// Unmarshal JSON fields
		json.Unmarshal(capabilitiesJSON, &client.Capabilities)
		json.Unmarshal(configJSON, &client.Configuration)
		json.Unmarshal(metadataJSON, &client.Metadata)

		clients = append(clients, &client)
	}

	return clients, nil
}

func (m *Manager) UpdateClient(ctx context.Context, clientID uuid.UUID, update *ClientUpdate) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, exists := m.clients[clientID]
	if !exists {
		return fmt.Errorf("client not found")
	}

	// Update client
	if update.Status != "" {
		client.Status = update.Status
	}
	if update.Configuration != nil {
		client.Configuration = update.Configuration
	}
	if update.Metadata != nil {
		client.Metadata = update.Metadata
	}
	client.LastSeen = time.Now()
	client.UpdatedAt = time.Now()

	// Update database
	configJSON, _ := json.Marshal(client.Configuration)
	metadataJSON, _ := json.Marshal(client.Metadata)
	
	query := `
		UPDATE clients SET status = $1, configuration = $2, metadata = $3, last_seen = $4, updated_at = $5
		WHERE id = $6
	`
	_, err := m.db.Conn.ExecContext(ctx, query,
		client.Status, configJSON, metadataJSON, client.LastSeen, client.UpdatedAt, clientID)
	if err != nil {
		return fmt.Errorf("failed to update client: %w", err)
	}

	return nil
}

func (m *Manager) UnregisterClient(ctx context.Context, clientID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Delete from database
	query := `DELETE FROM clients WHERE id = $1`
	result, err := m.db.Conn.ExecContext(ctx, query, clientID)
	if err != nil {
		return fmt.Errorf("failed to unregister client: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("client not found")
	}

	// Remove from memory
	delete(m.clients, clientID)

	logging.Info("unregistered client", map[string]interface{}{
		"client_id": clientID,
	})

	return nil
}

func (m *Manager) BroadcastToSession(ctx context.Context, sessionID uuid.UUID, message *ClientMessage) error {
	// Get all clients in session
	filter := &ClientFilter{
		SessionID: sessionID,
		Status:    "active",
	}
	
	clients, err := m.GetClientsBySession(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to get clients: %w", err)
	}

	// Store message in database
	messageID := uuid.New()
	query := `
		INSERT INTO client_messages (id, session_id, client_id, type, action, data, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	dataJSON, _ := json.Marshal(message.Data)
	_, err = m.db.Conn.ExecContext(ctx, query,
		messageID, sessionID, message.ClientID, message.Type, message.Action,
		dataJSON, time.Now())
	if err != nil {
		logging.Error("failed to store client message", map[string]interface{}{
			"session_id": sessionID,
			"client_id":  message.ClientID,
			"error":      err.Error(),
		})
	}

	// Broadcast to all clients (implementation would vary by platform)
	for _, client := range clients {
		if client.ID != message.ClientID {
			go m.sendMessageToClient(client, message)
		}
	}

	return nil
}

func (m *Manager) sendMessageToClient(client *Client, message *ClientMessage) {
	// Platform-specific message sending would be implemented here
	// This is a placeholder for the actual implementation
	logging.Debug("sending message to client", map[string]interface{}{
		"client_id":    client.ID,
		"client_type":  client.Type,
		"platform":     client.Platform,
		"message_type": message.Type,
		"action":       message.Action,
	})
}

func (m *Manager) GetClientStats(ctx context.Context, sessionID uuid.UUID) (*ClientStats, error) {
	query := `
		SELECT 
			platform,
			COUNT(*) as count,
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active_count,
			AVG(EXTRACT(EPOCH FROM (NOW() - last_seen))) as avg_last_seen_seconds
		FROM clients
		WHERE session_id = $1
		GROUP BY platform
	`
	
	rows, err := m.db.Conn.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query client stats: %w", err)
	}
	defer rows.Close()

	stats := &ClientStats{
		SessionID:     sessionID,
		TotalClients:  0,
		ActiveClients: 0,
		Platforms:     make(map[string]*PlatformStats),
	}

	for rows.Next() {
		var platform string
		var count, activeCount int
		var avgLastSeenSeconds float64
		
		err := rows.Scan(&platform, &count, &activeCount, &avgLastSeenSeconds)
		if err != nil {
			return nil, fmt.Errorf("failed to scan client stats: %w", err)
		}

		stats.TotalClients += count
		stats.ActiveClients += activeCount
		stats.Platforms[platform] = &PlatformStats{
			Platform:      platform,
			TotalClients:  count,
			ActiveClients: activeCount,
			AvgLastSeen:   time.Duration(avgLastSeenSeconds) * time.Second,
		}
	}

	return stats, nil
}

func (m *Manager) GetClientCapabilities(ctx context.Context, clientID uuid.UUID) ([]string, error) {
	client, err := m.GetClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	return client.Capabilities, nil
}

func (m *Manager) UpdateClientCapabilities(ctx context.Context, clientID uuid.UUID, capabilities []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, exists := m.clients[clientID]
	if !exists {
		return fmt.Errorf("client not found")
	}

	client.Capabilities = capabilities
	client.UpdatedAt = time.Now()

	// Update database
	capabilitiesJSON, _ := json.Marshal(capabilities)
	query := `UPDATE clients SET capabilities = $1, updated_at = $2 WHERE id = $3`
	_, err := m.db.Conn.ExecContext(ctx, query, capabilitiesJSON, client.UpdatedAt, clientID)
	if err != nil {
		return fmt.Errorf("failed to update client capabilities: %w", err)
	}

	return nil
}

func (m *Manager) HeartbeatClient(ctx context.Context, clientID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, exists := m.clients[clientID]
	if !exists {
		return fmt.Errorf("client not found")
	}

	client.LastSeen = time.Now()
	client.UpdatedAt = time.Now()

	// Update database
	query := `UPDATE clients SET last_seen = $1, updated_at = $2 WHERE id = $3`
	_, err := m.db.Conn.ExecContext(ctx, query, client.LastSeen, client.UpdatedAt, clientID)
	if err != nil {
		return fmt.Errorf("failed to update client heartbeat: %w", err)
	}

	return nil
}

func (m *Manager) CleanupInactiveClients(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Find clients that haven't been seen in 5 minutes
	cutoff := time.Now().Add(-5 * time.Minute)
	query := `
		SELECT id FROM clients
		WHERE last_seen < $1 AND status = 'active'
	`
	
	rows, err := m.db.Conn.QueryContext(ctx, query, cutoff)
	if err != nil {
		logging.Error("failed to query inactive clients", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()

	var inactiveClients []uuid.UUID
	for rows.Next() {
		var clientID uuid.UUID
		err := rows.Scan(&clientID)
		if err != nil {
			logging.Error("failed to scan inactive client", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}
		inactiveClients = append(inactiveClients, clientID)
	}

	// Mark inactive clients
	for _, clientID := range inactiveClients {
		query := `UPDATE clients SET status = 'inactive', updated_at = $1 WHERE id = $2`
		_, err := m.db.Conn.ExecContext(ctx, query, time.Now(), clientID)
		if err != nil {
			logging.Error("failed to mark client as inactive", map[string]interface{}{
				"client_id": clientID,
				"error":     err.Error(),
			})
			continue
		}

		delete(m.clients, clientID)
		logging.Info("marked client as inactive", map[string]interface{}{
			"client_id": clientID,
		})
	}
}

func (m *Manager) StartCleanupWorker(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.CleanupInactiveClients(ctx)
		}
	}
}

type ClientStats struct {
	SessionID     uuid.UUID                   `json:"session_id"`
	TotalClients  int                         `json:"total_clients"`
	ActiveClients int                         `json:"active_clients"`
	Platforms     map[string]*PlatformStats   `json:"platforms"`
}

type PlatformStats struct {
	Platform      string        `json:"platform"`
	TotalClients  int           `json:"total_clients"`
	ActiveClients int           `json:"active_clients"`
	AvgLastSeen   time.Duration `json:"avg_last_seen"`
}