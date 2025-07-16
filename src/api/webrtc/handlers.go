package webrtc

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"holodeck1/auth"
	"holodeck1/logging"
	"holodeck1/webrtc"
)

type Handler struct {
	webrtcManager *webrtc.Manager
}

func NewHandler(webrtcManager *webrtc.Manager) *Handler {
	return &Handler{
		webrtcManager: webrtcManager,
	}
}

func (h *Handler) CreateRTCSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req webrtc.CreateRTCSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user from context
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	req.UserID = user.ID

	// Validate required fields
	if req.SessionID == uuid.Nil {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	rtcSession, err := h.webrtcManager.CreateRTCSession(ctx, &req)
	if err != nil {
		logging.Error("failed to create RTC session", map[string]interface{}{
			"session_id": req.SessionID,
			"user_id":    req.UserID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":     true,
		"rtc_session": rtcSession,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) JoinRTCSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionId"])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	var req webrtc.JoinRTCSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user from context
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	req.UserID = user.ID
	req.SessionID = sessionID

	participant, err := h.webrtcManager.JoinRTCSession(ctx, sessionID, &req)
	if err != nil {
		logging.Error("failed to join RTC session", map[string]interface{}{
			"session_id": sessionID,
			"user_id":    req.UserID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":     true,
		"participant": participant,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) LeaveRTCSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionId"])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// Get user from context
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	err = h.webrtcManager.LeaveRTCSession(ctx, sessionID, user.ID)
	if err != nil {
		logging.Error("failed to leave RTC session", map[string]interface{}{
			"session_id": sessionID,
			"user_id":    user.ID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Successfully left RTC session",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetRTCStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionId"])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	stats, err := h.webrtcManager.GetRTCStats(ctx, sessionID)
	if err != nil {
		logging.Error("failed to get RTC stats", map[string]interface{}{
			"session_id": sessionID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"stats":   stats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionId"])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// Get user from context
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	h.webrtcManager.HandleWebSocketConnection(w, r, user.ID, sessionID)
}