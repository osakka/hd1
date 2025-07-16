package webrtc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/google/uuid"
	"github.com/pion/webrtc/v4"
	"holodeck1/database"
	"holodeck1/logging"
)

type Manager struct {
	db         *database.DB
	sessions   map[uuid.UUID]*RTCSession
	upgrader   websocket.Upgrader
	mu         sync.RWMutex
	api        *webrtc.API
	config     webrtc.Configuration
}

type RTCSession struct {
	ID           uuid.UUID
	SessionID    uuid.UUID
	Participants map[uuid.UUID]*RTCParticipant
	CreatedAt    time.Time
	UpdatedAt    time.Time
	mu           sync.RWMutex
}

type RTCParticipant struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	SessionID  uuid.UUID
	Connection *webrtc.PeerConnection
	DataChannel *webrtc.DataChannel
	WebSocket  *websocket.Conn
	Role       string
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type RTCMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	UserID    uuid.UUID   `json:"user_id"`
	SessionID uuid.UUID   `json:"session_id"`
}

type SignalingMessage struct {
	Type        string                     `json:"type"`
	SDP         *webrtc.SessionDescription `json:"sdp,omitempty"`
	ICECandidate *webrtc.ICECandidate      `json:"ice_candidate,omitempty"`
	UserID      uuid.UUID                  `json:"user_id"`
	TargetID    uuid.UUID                  `json:"target_id,omitempty"`
}

type RTCStats struct {
	SessionID         uuid.UUID `json:"session_id"`
	ParticipantCount  int       `json:"participant_count"`
	ActiveConnections int       `json:"active_connections"`
	DataChannels      int       `json:"data_channels"`
	LastActivity      time.Time `json:"last_activity"`
}

type CreateRTCSessionRequest struct {
	SessionID uuid.UUID `json:"session_id"`
	UserID    uuid.UUID `json:"user_id"`
	Role      string    `json:"role"`
}

type JoinRTCSessionRequest struct {
	SessionID uuid.UUID `json:"session_id"`
	UserID    uuid.UUID `json:"user_id"`
	Role      string    `json:"role"`
}

func NewManager(db *database.DB) *Manager {
	// Configure WebRTC
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create WebRTC API
	api := webrtc.NewAPI()

	return &Manager{
		db:       db,
		sessions: make(map[uuid.UUID]*RTCSession),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		api:    api,
		config: config,
	}
}

func (m *Manager) CreateRTCSession(ctx context.Context, req *CreateRTCSessionRequest) (*RTCSession, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	sessionID := uuid.New()
	rtcSession := &RTCSession{
		ID:           sessionID,
		SessionID:    req.SessionID,
		Participants: make(map[uuid.UUID]*RTCParticipant),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Store in database
	query := `
		INSERT INTO rtc_sessions (id, session_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := m.db.ExecContext(ctx, query, sessionID, req.SessionID, rtcSession.CreatedAt, rtcSession.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to store RTC session: %w", err)
	}

	m.sessions[sessionID] = rtcSession

	logging.Info("created RTC session", map[string]interface{}{
		"rtc_session_id": sessionID,
		"session_id":     req.SessionID,
		"user_id":        req.UserID,
	})

	return rtcSession, nil
}

func (m *Manager) JoinRTCSession(ctx context.Context, sessionID uuid.UUID, req *JoinRTCSessionRequest) (*RTCParticipant, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	rtcSession, exists := m.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("RTC session not found")
	}

	rtcSession.mu.Lock()
	defer rtcSession.mu.Unlock()

	// Check if participant already exists
	for _, participant := range rtcSession.Participants {
		if participant.UserID == req.UserID {
			return participant, nil
		}
	}

	// Create peer connection
	peerConnection, err := m.api.NewPeerConnection(m.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %w", err)
	}

	// Create data channel
	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create data channel: %w", err)
	}

	participantID := uuid.New()
	participant := &RTCParticipant{
		ID:         participantID,
		UserID:     req.UserID,
		SessionID:  sessionID,
		Connection: peerConnection,
		DataChannel: dataChannel,
		Role:       req.Role,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Store in database
	query := `
		INSERT INTO rtc_participants (id, user_id, session_id, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = m.db.ExecContext(ctx, query, participantID, req.UserID, sessionID, req.Role, true, participant.CreatedAt, participant.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to store RTC participant: %w", err)
	}

	rtcSession.Participants[participantID] = participant

	// Set up connection handlers
	m.setupConnectionHandlers(participant)

	logging.Info("user joined RTC session", map[string]interface{}{
		"participant_id": participantID,
		"user_id":        req.UserID,
		"session_id":     sessionID,
		"role":           req.Role,
	})

	return participant, nil
}

func (m *Manager) setupConnectionHandlers(participant *RTCParticipant) {
	// Handle ICE connection state changes
	participant.Connection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		logging.Debug("ICE connection state changed", map[string]interface{}{
			"participant_id": participant.ID,
			"state":          connectionState.String(),
		})
	})

	// Handle data channel messages
	participant.DataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		m.handleDataChannelMessage(participant, msg)
	})

	// Handle connection close
	participant.Connection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		if state == webrtc.PeerConnectionStateClosed || state == webrtc.PeerConnectionStateFailed {
			m.removeParticipant(participant)
		}
	})
}

func (m *Manager) handleDataChannelMessage(participant *RTCParticipant, msg webrtc.DataChannelMessage) {
	var rtcMessage RTCMessage
	if err := json.Unmarshal(msg.Data, &rtcMessage); err != nil {
		logging.Error("failed to unmarshal data channel message", map[string]interface{}{
			"participant_id": participant.ID,
			"error":          err.Error(),
		})
		return
	}

	rtcMessage.UserID = participant.UserID
	rtcMessage.SessionID = participant.SessionID
	rtcMessage.Timestamp = time.Now()

	// Broadcast to other participants
	m.broadcastToSession(participant.SessionID, rtcMessage, participant.ID)
}

func (m *Manager) broadcastToSession(sessionID uuid.UUID, message RTCMessage, excludeParticipant uuid.UUID) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, rtcSession := range m.sessions {
		if rtcSession.SessionID == sessionID {
			rtcSession.mu.RLock()
			for participantID, participant := range rtcSession.Participants {
				if participantID != excludeParticipant && participant.IsActive {
					go m.sendToParticipant(participant, message)
				}
			}
			rtcSession.mu.RUnlock()
			break
		}
	}
}

func (m *Manager) sendToParticipant(participant *RTCParticipant, message RTCMessage) {
	data, err := json.Marshal(message)
	if err != nil {
		logging.Error("failed to marshal message", map[string]interface{}{
			"participant_id": participant.ID,
			"error":          err.Error(),
		})
		return
	}

	if participant.DataChannel != nil && participant.DataChannel.ReadyState() == webrtc.DataChannelStateOpen {
		err = participant.DataChannel.Send(data)
		if err != nil {
			logging.Error("failed to send data channel message", map[string]interface{}{
				"participant_id": participant.ID,
				"error":          err.Error(),
			})
		}
	}
}

func (m *Manager) removeParticipant(participant *RTCParticipant) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, rtcSession := range m.sessions {
		if rtcSession.SessionID == participant.SessionID {
			rtcSession.mu.Lock()
			delete(rtcSession.Participants, participant.ID)
			rtcSession.mu.Unlock()
			break
		}
	}

	// Update database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE rtc_participants SET is_active = false, updated_at = $1 WHERE id = $2`
	_, err := m.db.ExecContext(ctx, query, time.Now(), participant.ID)
	if err != nil {
		logging.Error("failed to update participant status", map[string]interface{}{
			"participant_id": participant.ID,
			"error":          err.Error(),
		})
	}

	// Close connections
	if participant.Connection != nil {
		participant.Connection.Close()
	}
	if participant.WebSocket != nil {
		participant.WebSocket.Close()
	}

	logging.Info("participant removed from RTC session", map[string]interface{}{
		"participant_id": participant.ID,
		"user_id":        participant.UserID,
		"session_id":     participant.SessionID,
	})
}

func (m *Manager) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request, userID uuid.UUID, sessionID uuid.UUID) {
	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.Error("failed to upgrade websocket connection", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	defer conn.Close()

	// Find participant
	var participant *RTCParticipant
	m.mu.RLock()
	for _, rtcSession := range m.sessions {
		if rtcSession.SessionID == sessionID {
			rtcSession.mu.RLock()
			for _, p := range rtcSession.Participants {
				if p.UserID == userID {
					participant = p
					break
				}
			}
			rtcSession.mu.RUnlock()
			break
		}
	}
	m.mu.RUnlock()

	if participant == nil {
		logging.Error("participant not found", map[string]interface{}{
			"user_id":    userID,
			"session_id": sessionID,
		})
		return
	}

	participant.WebSocket = conn

	// Handle WebSocket messages
	for {
		var msg SignalingMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logging.Error("websocket error", map[string]interface{}{
					"error": err.Error(),
				})
			}
			break
		}

		m.handleSignalingMessage(participant, msg)
	}
}

func (m *Manager) handleSignalingMessage(participant *RTCParticipant, msg SignalingMessage) {
	switch msg.Type {
	case "offer":
		m.handleOffer(participant, msg)
	case "answer":
		m.handleAnswer(participant, msg)
	case "ice-candidate":
		m.handleICECandidate(participant, msg)
	default:
		logging.Warn("unknown signaling message type", map[string]interface{}{
			"type":           msg.Type,
			"participant_id": participant.ID,
		})
	}
}

func (m *Manager) handleOffer(participant *RTCParticipant, msg SignalingMessage) {
	if msg.SDP == nil {
		return
	}

	err := participant.Connection.SetRemoteDescription(*msg.SDP)
	if err != nil {
		logging.Error("failed to set remote description", map[string]interface{}{
			"participant_id": participant.ID,
			"error":          err.Error(),
		})
		return
	}

	answer, err := participant.Connection.CreateAnswer(nil)
	if err != nil {
		logging.Error("failed to create answer", map[string]interface{}{
			"participant_id": participant.ID,
			"error":          err.Error(),
		})
		return
	}

	err = participant.Connection.SetLocalDescription(answer)
	if err != nil {
		logging.Error("failed to set local description", map[string]interface{}{
			"participant_id": participant.ID,
			"error":          err.Error(),
		})
		return
	}

	// Send answer back
	response := SignalingMessage{
		Type:   "answer",
		SDP:    &answer,
		UserID: participant.UserID,
	}

	if participant.WebSocket != nil {
		participant.WebSocket.WriteJSON(response)
	}
}

func (m *Manager) handleAnswer(participant *RTCParticipant, msg SignalingMessage) {
	if msg.SDP == nil {
		return
	}

	err := participant.Connection.SetRemoteDescription(*msg.SDP)
	if err != nil {
		logging.Error("failed to set remote description", map[string]interface{}{
			"participant_id": participant.ID,
			"error":          err.Error(),
		})
	}
}

func (m *Manager) handleICECandidate(participant *RTCParticipant, msg SignalingMessage) {
	if msg.ICECandidate == nil {
		return
	}

	candidate := webrtc.ICECandidateInit{
		Candidate: msg.ICECandidate.ToJSON().Candidate,
		SDPMid:    &msg.ICECandidate.SDPMid,
		SDPMLineIndex: &msg.ICECandidate.SDPMLineIndex,
	}
	err := participant.Connection.AddICECandidate(candidate)
	if err != nil {
		logging.Error("failed to add ICE candidate", map[string]interface{}{
			"participant_id": participant.ID,
			"error":          err.Error(),
		})
	}
}

func (m *Manager) GetRTCStats(ctx context.Context, sessionID uuid.UUID) (*RTCStats, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var rtcSession *RTCSession
	for _, session := range m.sessions {
		if session.SessionID == sessionID {
			rtcSession = session
			break
		}
	}

	if rtcSession == nil {
		return nil, fmt.Errorf("RTC session not found")
	}

	rtcSession.mu.RLock()
	defer rtcSession.mu.RUnlock()

	activeConnections := 0
	dataChannels := 0
	for _, participant := range rtcSession.Participants {
		if participant.IsActive {
			activeConnections++
			if participant.DataChannel != nil {
				dataChannels++
			}
		}
	}

	return &RTCStats{
		SessionID:         sessionID,
		ParticipantCount:  len(rtcSession.Participants),
		ActiveConnections: activeConnections,
		DataChannels:      dataChannels,
		LastActivity:      rtcSession.UpdatedAt,
	}, nil
}

func (m *Manager) LeaveRTCSession(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, rtcSession := range m.sessions {
		if rtcSession.SessionID == sessionID {
			rtcSession.mu.RLock()
			for _, participant := range rtcSession.Participants {
				if participant.UserID == userID {
					rtcSession.mu.RUnlock()
					m.removeParticipant(participant)
					return nil
				}
			}
			rtcSession.mu.RUnlock()
			break
		}
	}

	return fmt.Errorf("participant not found in RTC session")
}

func (m *Manager) CleanupInactiveSessions(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for sessionID, rtcSession := range m.sessions {
		rtcSession.mu.Lock()
		if len(rtcSession.Participants) == 0 || time.Since(rtcSession.UpdatedAt) > 24*time.Hour {
			delete(m.sessions, sessionID)
			logging.Info("cleaned up inactive RTC session", map[string]interface{}{
				"session_id": sessionID,
			})
		}
		rtcSession.mu.Unlock()
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
			m.CleanupInactiveSessions(ctx)
		}
	}
}