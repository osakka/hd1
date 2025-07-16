package session

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"holodeck1/database"
	"holodeck1/logging"
)

type Manager struct {
	db *database.DB
}

type Session struct {
	ID                  uuid.UUID              `json:"id"`
	Name                string                 `json:"name"`
	Description         string                 `json:"description"`
	OwnerID             uuid.UUID              `json:"owner_id"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
	Status              string                 `json:"status"`
	Visibility          string                 `json:"visibility"`
	MaxParticipants     int                    `json:"max_participants"`
	CurrentParticipants int                    `json:"current_participants"`
	Settings            map[string]interface{} `json:"settings"`
	Metadata            map[string]interface{} `json:"metadata"`
}

type Participant struct {
	ID          uuid.UUID              `json:"id"`
	SessionID   uuid.UUID              `json:"session_id"`
	UserID      uuid.UUID              `json:"user_id"`
	Role        string                 `json:"role"`
	JoinedAt    time.Time              `json:"joined_at"`
	LeftAt      *time.Time             `json:"left_at"`
	LastSeen    time.Time              `json:"last_seen"`
	Permissions map[string]interface{} `json:"permissions"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type CreateSessionRequest struct {
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	OwnerID         uuid.UUID              `json:"owner_id"`
	Visibility      string                 `json:"visibility"`
	MaxParticipants int                    `json:"max_participants"`
	Settings        map[string]interface{} `json:"settings"`
	Metadata        map[string]interface{} `json:"metadata"`
}

type UpdateSessionRequest struct {
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	MaxParticipants int                    `json:"max_participants"`
	Settings        map[string]interface{} `json:"settings"`
	Metadata        map[string]interface{} `json:"metadata"`
}

type JoinSessionRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
}

type SessionFilter struct {
	Status     string `json:"status"`
	Visibility string `json:"visibility"`
	OwnerID    string `json:"owner_id"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}

func NewManager(db *database.DB) *Manager {
	return &Manager{db: db}
}

func (m *Manager) CreateSession(ctx context.Context, req *CreateSessionRequest) (*Session, error) {
	session := &Session{
		ID:              uuid.New(),
		Name:            req.Name,
		Description:     req.Description,
		OwnerID:         req.OwnerID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Status:          "active",
		Visibility:      req.Visibility,
		MaxParticipants: req.MaxParticipants,
		Settings:        req.Settings,
		Metadata:        req.Metadata,
	}

	if session.Visibility == "" {
		session.Visibility = "private"
	}
	if session.MaxParticipants == 0 {
		session.MaxParticipants = 10
	}
	if session.Settings == nil {
		session.Settings = make(map[string]interface{})
	}
	if session.Metadata == nil {
		session.Metadata = make(map[string]interface{})
	}

	settingsJSON, _ := json.Marshal(session.Settings)
	metadataJSON, _ := json.Marshal(session.Metadata)

	query := `
		INSERT INTO sessions (id, name, description, owner_id, status, visibility, max_participants, settings, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at, current_participants
	`

	err := m.db.QueryRowContext(ctx, query,
		session.ID,
		session.Name,
		session.Description,
		session.OwnerID,
		session.Status,
		session.Visibility,
		session.MaxParticipants,
		settingsJSON,
		metadataJSON,
	).Scan(&session.CreatedAt, &session.UpdatedAt, &session.CurrentParticipants)

	if err != nil {
		logging.Error("failed to create session", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	logging.Info("session created", map[string]interface{}{
		"session_id": session.ID,
		"name":       session.Name,
		"owner_id":   session.OwnerID,
	})

	return session, nil
}

func (m *Manager) GetSession(ctx context.Context, sessionID uuid.UUID) (*Session, error) {
	session := &Session{}
	var settingsJSON, metadataJSON []byte

	query := `
		SELECT id, name, description, owner_id, created_at, updated_at, status, 
		       visibility, max_participants, current_participants, settings, metadata
		FROM sessions WHERE id = $1
	`

	err := m.db.QueryRowContext(ctx, query, sessionID).Scan(
		&session.ID,
		&session.Name,
		&session.Description,
		&session.OwnerID,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.Status,
		&session.Visibility,
		&session.MaxParticipants,
		&session.CurrentParticipants,
		&settingsJSON,
		&metadataJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		logging.Error("failed to get session", map[string]interface{}{
			"session_id": sessionID,
			"error":      err.Error(),
		})
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	json.Unmarshal(settingsJSON, &session.Settings)
	json.Unmarshal(metadataJSON, &session.Metadata)

	return session, nil
}

func (m *Manager) ListSessions(ctx context.Context, filter *SessionFilter) ([]*Session, error) {
	query := `
		SELECT id, name, description, owner_id, created_at, updated_at, status, 
		       visibility, max_participants, current_participants, settings, metadata
		FROM sessions
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	if filter.Visibility != "" {
		query += fmt.Sprintf(" AND visibility = $%d", argIndex)
		args = append(args, filter.Visibility)
		argIndex++
	}

	if filter.OwnerID != "" {
		query += fmt.Sprintf(" AND owner_id = $%d", argIndex)
		args = append(args, filter.OwnerID)
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
		argIndex++
	}

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		logging.Error("failed to list sessions", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}
	defer rows.Close()

	sessions := []*Session{}
	for rows.Next() {
		session := &Session{}
		var settingsJSON, metadataJSON []byte

		err := rows.Scan(
			&session.ID,
			&session.Name,
			&session.Description,
			&session.OwnerID,
			&session.CreatedAt,
			&session.UpdatedAt,
			&session.Status,
			&session.Visibility,
			&session.MaxParticipants,
			&session.CurrentParticipants,
			&settingsJSON,
			&metadataJSON,
		)

		if err != nil {
			logging.Error("failed to scan session", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}

		json.Unmarshal(settingsJSON, &session.Settings)
		json.Unmarshal(metadataJSON, &session.Metadata)

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (m *Manager) UpdateSession(ctx context.Context, sessionID uuid.UUID, req *UpdateSessionRequest) (*Session, error) {
	session, err := m.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		session.Name = req.Name
	}
	if req.Description != "" {
		session.Description = req.Description
	}
	if req.MaxParticipants > 0 {
		session.MaxParticipants = req.MaxParticipants
	}
	if req.Settings != nil {
		session.Settings = req.Settings
	}
	if req.Metadata != nil {
		session.Metadata = req.Metadata
	}

	settingsJSON, _ := json.Marshal(session.Settings)
	metadataJSON, _ := json.Marshal(session.Metadata)

	query := `
		UPDATE sessions 
		SET name = $2, description = $3, max_participants = $4, settings = $5, metadata = $6
		WHERE id = $1
		RETURNING updated_at
	`

	err = m.db.QueryRowContext(ctx, query,
		sessionID,
		session.Name,
		session.Description,
		session.MaxParticipants,
		settingsJSON,
		metadataJSON,
	).Scan(&session.UpdatedAt)

	if err != nil {
		logging.Error("failed to update session", map[string]interface{}{
			"session_id": sessionID,
			"error":      err.Error(),
		})
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	logging.Info("session updated", map[string]interface{}{
		"session_id": sessionID,
		"name":       session.Name,
	})

	return session, nil
}

func (m *Manager) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	query := `DELETE FROM sessions WHERE id = $1`

	result, err := m.db.ExecContext(ctx, query, sessionID)
	if err != nil {
		logging.Error("failed to delete session", map[string]interface{}{
			"session_id": sessionID,
			"error":      err.Error(),
		})
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	logging.Info("session deleted", map[string]interface{}{
		"session_id": sessionID,
	})

	return nil
}

func (m *Manager) JoinSession(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, role string) error {
	if role == "" {
		role = "participant"
	}

	participant := &Participant{
		ID:        uuid.New(),
		SessionID: sessionID,
		UserID:    userID,
		Role:      role,
		JoinedAt:  time.Now(),
		LastSeen:  time.Now(),
	}

	query := `
		INSERT INTO participants (id, session_id, user_id, role, joined_at, last_seen)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (session_id, user_id) 
		DO UPDATE SET left_at = NULL, last_seen = NOW()
	`

	_, err := m.db.ExecContext(ctx, query,
		participant.ID,
		participant.SessionID,
		participant.UserID,
		participant.Role,
		participant.JoinedAt,
		participant.LastSeen,
	)

	if err != nil {
		logging.Error("failed to join session", map[string]interface{}{
			"session_id": sessionID,
			"user_id":    userID,
			"error":      err.Error(),
		})
		return fmt.Errorf("failed to join session: %w", err)
	}

	logging.Info("user joined session", map[string]interface{}{
		"session_id": sessionID,
		"user_id":    userID,
		"role":       role,
	})

	return nil
}

func (m *Manager) LeaveSession(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) error {
	query := `
		UPDATE participants 
		SET left_at = NOW() 
		WHERE session_id = $1 AND user_id = $2 AND left_at IS NULL
	`

	result, err := m.db.ExecContext(ctx, query, sessionID, userID)
	if err != nil {
		logging.Error("failed to leave session", map[string]interface{}{
			"session_id": sessionID,
			"user_id":    userID,
			"error":      err.Error(),
		})
		return fmt.Errorf("failed to leave session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("participant not found in session")
	}

	logging.Info("user left session", map[string]interface{}{
		"session_id": sessionID,
		"user_id":    userID,
	})

	return nil
}

func (m *Manager) GetSessionParticipants(ctx context.Context, sessionID uuid.UUID) ([]*Participant, error) {
	query := `
		SELECT id, session_id, user_id, role, joined_at, left_at, last_seen, permissions, metadata
		FROM participants 
		WHERE session_id = $1 AND left_at IS NULL
		ORDER BY joined_at ASC
	`

	rows, err := m.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		logging.Error("failed to get session participants", map[string]interface{}{
			"session_id": sessionID,
			"error":      err.Error(),
		})
		return nil, fmt.Errorf("failed to get session participants: %w", err)
	}
	defer rows.Close()

	participants := []*Participant{}
	for rows.Next() {
		participant := &Participant{}
		var permissionsJSON, metadataJSON []byte

		err := rows.Scan(
			&participant.ID,
			&participant.SessionID,
			&participant.UserID,
			&participant.Role,
			&participant.JoinedAt,
			&participant.LeftAt,
			&participant.LastSeen,
			&permissionsJSON,
			&metadataJSON,
		)

		if err != nil {
			logging.Error("failed to scan participant", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}

		json.Unmarshal(permissionsJSON, &participant.Permissions)
		json.Unmarshal(metadataJSON, &participant.Metadata)

		participants = append(participants, participant)
	}

	return participants, nil
}

func (m *Manager) UpdateLastSeen(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) error {
	query := `
		UPDATE participants 
		SET last_seen = NOW() 
		WHERE session_id = $1 AND user_id = $2 AND left_at IS NULL
	`

	_, err := m.db.ExecContext(ctx, query, sessionID, userID)
	if err != nil {
		logging.Error("failed to update last seen", map[string]interface{}{
			"session_id": sessionID,
			"user_id":    userID,
			"error":      err.Error(),
		})
		return fmt.Errorf("failed to update last seen: %w", err)
	}

	return nil
}