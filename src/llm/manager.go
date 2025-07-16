package llm

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
	db       *database.DB
	avatars  map[uuid.UUID]*LLMAvatar
	providers map[string]LLMProvider
	mu       sync.RWMutex
}

type LLMAvatar struct {
	ID          uuid.UUID              `json:"id"`
	SessionID   uuid.UUID              `json:"session_id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Provider    string                 `json:"provider"`
	Model       string                 `json:"model"`
	Personality string                 `json:"personality"`
	Knowledge   string                 `json:"knowledge"`
	Capabilities []string              `json:"capabilities"`
	Configuration map[string]interface{} `json:"configuration"`
	State       string                 `json:"state"`
	Position    *Position3D            `json:"position"`
	Avatar3D    *Avatar3DConfig        `json:"avatar_3d"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type Position3D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Avatar3DConfig struct {
	ModelURL    string            `json:"model_url"`
	Scale       float64           `json:"scale"`
	Animations  []string          `json:"animations"`
	Materials   map[string]string `json:"materials"`
	Lighting    map[string]interface{} `json:"lighting"`
}

type LLMProvider interface {
	GenerateResponse(ctx context.Context, prompt string, config map[string]interface{}) (*LLMResponse, error)
	GetCapabilities() []string
	ValidateConfig(config map[string]interface{}) error
}

type LLMResponse struct {
	Content     string                 `json:"content"`
	Type        string                 `json:"type"`
	Metadata    map[string]interface{} `json:"metadata"`
	Usage       *UsageStats            `json:"usage"`
	Timestamp   time.Time              `json:"timestamp"`
}

type UsageStats struct {
	InputTokens     int     `json:"input_tokens"`
	OutputTokens    int     `json:"output_tokens"`
	TotalTokens     int     `json:"total_tokens"`
	ProcessingTime  float64 `json:"processing_time"`
	Cost            float64 `json:"cost"`
}

type CreateAvatarRequest struct {
	SessionID     uuid.UUID              `json:"session_id"`
	Name          string                 `json:"name"`
	Type          string                 `json:"type"`
	Provider      string                 `json:"provider"`
	Model         string                 `json:"model"`
	Personality   string                 `json:"personality"`
	Knowledge     string                 `json:"knowledge"`
	Capabilities  []string               `json:"capabilities"`
	Configuration map[string]interface{} `json:"configuration"`
	Position      *Position3D            `json:"position"`
	Avatar3D      *Avatar3DConfig        `json:"avatar_3d"`
}

type ChatRequest struct {
	AvatarID uuid.UUID `json:"avatar_id"`
	Message  string    `json:"message"`
	Context  string    `json:"context"`
	UserID   uuid.UUID `json:"user_id"`
}

type ChatResponse struct {
	ID         uuid.UUID  `json:"id"`
	AvatarID   uuid.UUID  `json:"avatar_id"`
	UserID     uuid.UUID  `json:"user_id"`
	Message    string     `json:"message"`
	Response   string     `json:"response"`
	Context    string     `json:"context"`
	Metadata   map[string]interface{} `json:"metadata"`
	Usage      *UsageStats `json:"usage"`
	Timestamp  time.Time   `json:"timestamp"`
}

type AvatarFilter struct {
	SessionID uuid.UUID `json:"session_id"`
	Type      string    `json:"type"`
	Provider  string    `json:"provider"`
	State     string    `json:"state"`
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
}

func NewManager(db *database.DB) *Manager {
	return &Manager{
		db:        db,
		avatars:   make(map[uuid.UUID]*LLMAvatar),
		providers: make(map[string]LLMProvider),
	}
}

func (m *Manager) RegisterProvider(name string, provider LLMProvider) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.providers[name] = provider
	logging.Info("registered LLM provider", map[string]interface{}{
		"provider": name,
		"capabilities": provider.GetCapabilities(),
	})
}

func (m *Manager) CreateAvatar(ctx context.Context, req *CreateAvatarRequest) (*LLMAvatar, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Validate provider
	provider, exists := m.providers[req.Provider]
	if !exists {
		return nil, fmt.Errorf("provider not found: %s", req.Provider)
	}

	// Validate configuration
	err := provider.ValidateConfig(req.Configuration)
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	avatarID := uuid.New()
	avatar := &LLMAvatar{
		ID:            avatarID,
		SessionID:     req.SessionID,
		Name:          req.Name,
		Type:          req.Type,
		Provider:      req.Provider,
		Model:         req.Model,
		Personality:   req.Personality,
		Knowledge:     req.Knowledge,
		Capabilities:  req.Capabilities,
		Configuration: req.Configuration,
		State:         "active",
		Position:      req.Position,
		Avatar3D:      req.Avatar3D,
		Metadata:      make(map[string]interface{}),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Store in database
	query := `
		INSERT INTO llm_avatars (id, session_id, name, type, provider, model, personality, knowledge, capabilities, configuration, state, position, avatar_3d, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	
	capabilitiesJSON, _ := json.Marshal(avatar.Capabilities)
	configJSON, _ := json.Marshal(avatar.Configuration)
	positionJSON, _ := json.Marshal(avatar.Position)
	avatar3DJSON, _ := json.Marshal(avatar.Avatar3D)
	metadataJSON, _ := json.Marshal(avatar.Metadata)
	
	_, err = m.db.Conn.ExecContext(ctx, query,
		avatar.ID, avatar.SessionID, avatar.Name, avatar.Type, avatar.Provider,
		avatar.Model, avatar.Personality, avatar.Knowledge, capabilitiesJSON,
		configJSON, avatar.State, positionJSON, avatar3DJSON, metadataJSON,
		avatar.CreatedAt, avatar.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to store avatar: %w", err)
	}

	m.avatars[avatarID] = avatar

	logging.Info("created LLM avatar", map[string]interface{}{
		"avatar_id":  avatarID,
		"session_id": req.SessionID,
		"name":       req.Name,
		"type":       req.Type,
		"provider":   req.Provider,
		"model":      req.Model,
	})

	return avatar, nil
}

func (m *Manager) GetAvatar(ctx context.Context, avatarID uuid.UUID) (*LLMAvatar, error) {
	m.mu.RLock()
	avatar, exists := m.avatars[avatarID]
	m.mu.RUnlock()

	if exists {
		return avatar, nil
	}

	// Load from database
	query := `
		SELECT id, session_id, name, type, provider, model, personality, knowledge, capabilities, configuration, state, position, avatar_3d, metadata, created_at, updated_at
		FROM llm_avatars
		WHERE id = $1
	`
	
	var dbAvatar LLMAvatar
	var capabilitiesJSON, configJSON, positionJSON, avatar3DJSON, metadataJSON []byte
	
	err := m.db.Conn.QueryRowContext(ctx, query, avatarID).Scan(
		&dbAvatar.ID, &dbAvatar.SessionID, &dbAvatar.Name, &dbAvatar.Type,
		&dbAvatar.Provider, &dbAvatar.Model, &dbAvatar.Personality,
		&dbAvatar.Knowledge, &capabilitiesJSON, &configJSON, &dbAvatar.State,
		&positionJSON, &avatar3DJSON, &metadataJSON, &dbAvatar.CreatedAt, &dbAvatar.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("avatar not found")
		}
		return nil, fmt.Errorf("failed to get avatar: %w", err)
	}

	// Unmarshal JSON fields
	json.Unmarshal(capabilitiesJSON, &dbAvatar.Capabilities)
	json.Unmarshal(configJSON, &dbAvatar.Configuration)
	json.Unmarshal(positionJSON, &dbAvatar.Position)
	json.Unmarshal(avatar3DJSON, &dbAvatar.Avatar3D)
	json.Unmarshal(metadataJSON, &dbAvatar.Metadata)

	m.mu.Lock()
	m.avatars[avatarID] = &dbAvatar
	m.mu.Unlock()

	return &dbAvatar, nil
}

func (m *Manager) GetAvatarsBySession(ctx context.Context, filter *AvatarFilter) ([]*LLMAvatar, error) {
	query := `
		SELECT id, session_id, name, type, provider, model, personality, knowledge, capabilities, configuration, state, position, avatar_3d, metadata, created_at, updated_at
		FROM llm_avatars
		WHERE session_id = $1
	`
	args := []interface{}{filter.SessionID}
	argIndex := 2

	if filter.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, filter.Type)
		argIndex++
	}

	if filter.Provider != "" {
		query += fmt.Sprintf(" AND provider = $%d", argIndex)
		args = append(args, filter.Provider)
		argIndex++
	}

	if filter.State != "" {
		query += fmt.Sprintf(" AND state = $%d", argIndex)
		args = append(args, filter.State)
		argIndex++
	}

	query += " ORDER BY created_at DESC"

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
		return nil, fmt.Errorf("failed to query avatars: %w", err)
	}
	defer rows.Close()

	var avatars []*LLMAvatar
	for rows.Next() {
		var avatar LLMAvatar
		var capabilitiesJSON, configJSON, positionJSON, avatar3DJSON, metadataJSON []byte
		
		err := rows.Scan(
			&avatar.ID, &avatar.SessionID, &avatar.Name, &avatar.Type,
			&avatar.Provider, &avatar.Model, &avatar.Personality,
			&avatar.Knowledge, &capabilitiesJSON, &configJSON, &avatar.State,
			&positionJSON, &avatar3DJSON, &metadataJSON, &avatar.CreatedAt, &avatar.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan avatar: %w", err)
		}

		// Unmarshal JSON fields
		json.Unmarshal(capabilitiesJSON, &avatar.Capabilities)
		json.Unmarshal(configJSON, &avatar.Configuration)
		json.Unmarshal(positionJSON, &avatar.Position)
		json.Unmarshal(avatar3DJSON, &avatar.Avatar3D)
		json.Unmarshal(metadataJSON, &avatar.Metadata)

		avatars = append(avatars, &avatar)
	}

	return avatars, nil
}

func (m *Manager) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	avatar, err := m.GetAvatar(ctx, req.AvatarID)
	if err != nil {
		return nil, err
	}

	// Get provider
	m.mu.RLock()
	provider, exists := m.providers[avatar.Provider]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("provider not found: %s", avatar.Provider)
	}

	// Build prompt with personality and knowledge
	prompt := m.buildPrompt(avatar, req.Message, req.Context)

	// Generate response
	llmResponse, err := provider.GenerateResponse(ctx, prompt, avatar.Configuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	// Create chat response
	chatResponse := &ChatResponse{
		ID:        uuid.New(),
		AvatarID:  req.AvatarID,
		UserID:    req.UserID,
		Message:   req.Message,
		Response:  llmResponse.Content,
		Context:   req.Context,
		Metadata:  llmResponse.Metadata,
		Usage:     llmResponse.Usage,
		Timestamp: time.Now(),
	}

	// Store in database
	query := `
		INSERT INTO llm_conversations (id, avatar_id, user_id, message, response, context, metadata, usage, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	metadataJSON, _ := json.Marshal(chatResponse.Metadata)
	usageJSON, _ := json.Marshal(chatResponse.Usage)
	
	_, err = m.db.Conn.ExecContext(ctx, query,
		chatResponse.ID, chatResponse.AvatarID, chatResponse.UserID,
		chatResponse.Message, chatResponse.Response, chatResponse.Context,
		metadataJSON, usageJSON, chatResponse.Timestamp)
	if err != nil {
		logging.Error("failed to store chat response", map[string]interface{}{
			"avatar_id": req.AvatarID,
			"user_id":   req.UserID,
			"error":     err.Error(),
		})
	}

	logging.Debug("generated chat response", map[string]interface{}{
		"avatar_id":      req.AvatarID,
		"user_id":        req.UserID,
		"message_length": len(req.Message),
		"response_length": len(llmResponse.Content),
		"processing_time": llmResponse.Usage.ProcessingTime,
	})

	return chatResponse, nil
}

func (m *Manager) buildPrompt(avatar *LLMAvatar, message string, context string) string {
	prompt := ""
	
	if avatar.Personality != "" {
		prompt += fmt.Sprintf("Personality: %s\n", avatar.Personality)
	}
	
	if avatar.Knowledge != "" {
		prompt += fmt.Sprintf("Knowledge: %s\n", avatar.Knowledge)
	}
	
	if context != "" {
		prompt += fmt.Sprintf("Context: %s\n", context)
	}
	
	prompt += fmt.Sprintf("Human: %s\nAssistant:", message)
	
	return prompt
}

func (m *Manager) UpdateAvatarPosition(ctx context.Context, avatarID uuid.UUID, position *Position3D) error {
	avatar, err := m.GetAvatar(ctx, avatarID)
	if err != nil {
		return err
	}

	avatar.Position = position
	avatar.UpdatedAt = time.Now()

	// Update in database
	positionJSON, _ := json.Marshal(position)
	query := `UPDATE llm_avatars SET position = $1, updated_at = $2 WHERE id = $3`
	_, err = m.db.Conn.ExecContext(ctx, query, positionJSON, avatar.UpdatedAt, avatarID)
	if err != nil {
		return fmt.Errorf("failed to update avatar position: %w", err)
	}

	return nil
}

func (m *Manager) UpdateAvatarState(ctx context.Context, avatarID uuid.UUID, state string) error {
	avatar, err := m.GetAvatar(ctx, avatarID)
	if err != nil {
		return err
	}

	avatar.State = state
	avatar.UpdatedAt = time.Now()

	// Update in database
	query := `UPDATE llm_avatars SET state = $1, updated_at = $2 WHERE id = $3`
	_, err = m.db.Conn.ExecContext(ctx, query, state, avatar.UpdatedAt, avatarID)
	if err != nil {
		return fmt.Errorf("failed to update avatar state: %w", err)
	}

	return nil
}

func (m *Manager) DeleteAvatar(ctx context.Context, avatarID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Delete from database
	query := `DELETE FROM llm_avatars WHERE id = $1`
	result, err := m.db.Conn.ExecContext(ctx, query, avatarID)
	if err != nil {
		return fmt.Errorf("failed to delete avatar: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("avatar not found")
	}

	// Remove from memory
	delete(m.avatars, avatarID)

	logging.Info("deleted LLM avatar", map[string]interface{}{
		"avatar_id": avatarID,
	})

	return nil
}

func (m *Manager) GetConversationHistory(ctx context.Context, avatarID uuid.UUID, limit int) ([]*ChatResponse, error) {
	query := `
		SELECT id, avatar_id, user_id, message, response, context, metadata, usage, created_at
		FROM llm_conversations
		WHERE avatar_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	
	rows, err := m.db.Conn.QueryContext(ctx, query, avatarID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query conversation history: %w", err)
	}
	defer rows.Close()

	var conversations []*ChatResponse
	for rows.Next() {
		var chat ChatResponse
		var metadataJSON, usageJSON []byte
		
		err := rows.Scan(
			&chat.ID, &chat.AvatarID, &chat.UserID, &chat.Message,
			&chat.Response, &chat.Context, &metadataJSON, &usageJSON, &chat.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}

		// Unmarshal JSON fields
		json.Unmarshal(metadataJSON, &chat.Metadata)
		json.Unmarshal(usageJSON, &chat.Usage)

		conversations = append(conversations, &chat)
	}

	return conversations, nil
}

func (m *Manager) GetProviders() map[string][]string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	providers := make(map[string][]string)
	for name, provider := range m.providers {
		providers[name] = provider.GetCapabilities()
	}

	return providers
}

func (m *Manager) CleanupInactiveAvatars(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Find avatars that haven't been used in 24 hours
	query := `
		SELECT id FROM llm_avatars
		WHERE updated_at < $1 AND state = 'active'
	`
	
	cutoff := time.Now().Add(-24 * time.Hour)
	rows, err := m.db.Conn.QueryContext(ctx, query, cutoff)
	if err != nil {
		logging.Error("failed to query inactive avatars", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()

	var inactiveAvatars []uuid.UUID
	for rows.Next() {
		var avatarID uuid.UUID
		err := rows.Scan(&avatarID)
		if err != nil {
			logging.Error("failed to scan inactive avatar", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}
		inactiveAvatars = append(inactiveAvatars, avatarID)
	}

	// Mark inactive avatars
	for _, avatarID := range inactiveAvatars {
		query := `UPDATE llm_avatars SET state = 'inactive', updated_at = $1 WHERE id = $2`
		_, err := m.db.Conn.ExecContext(ctx, query, time.Now(), avatarID)
		if err != nil {
			logging.Error("failed to mark avatar as inactive", map[string]interface{}{
				"avatar_id": avatarID,
				"error":     err.Error(),
			})
			continue
		}

		delete(m.avatars, avatarID)
		logging.Info("marked avatar as inactive", map[string]interface{}{
			"avatar_id": avatarID,
		})
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
			m.CleanupInactiveAvatars(ctx)
		}
	}
}