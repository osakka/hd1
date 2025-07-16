package sessions

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"holodeck1/auth"
	"holodeck1/logging"
	"holodeck1/session"
)

type Handler struct {
	sessionManager *session.Manager
}

func NewHandler(sessionManager *session.Manager) *Handler {
	return &Handler{
		sessionManager: sessionManager,
	}
}

func (h *Handler) ListSessions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Parse query parameters
	filter := &session.SessionFilter{
		Status:     r.URL.Query().Get("status"),
		Visibility: r.URL.Query().Get("visibility"),
		OwnerID:    r.URL.Query().Get("owner_id"),
	}
	
	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filter.Limit = l
		}
	}
	
	if offset := r.URL.Query().Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			filter.Offset = o
		}
	}
	
	// Set default limit if not provided
	if filter.Limit == 0 {
		filter.Limit = 50
	}
	
	sessions, err := h.sessionManager.ListSessions(ctx, filter)
	if err != nil {
		logging.Error("failed to list sessions", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success":  true,
		"sessions": sessions,
		"pagination": map[string]interface{}{
			"limit":    filter.Limit,
			"offset":   filter.Offset,
			"has_more": len(sessions) >= filter.Limit,
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var req session.CreateSessionRequest
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
	
	req.OwnerID = user.ID
	
	session, err := h.sessionManager.CreateSession(ctx, &req)
	if err != nil {
		logging.Error("failed to create session", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success": true,
		"session": session,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionId"])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}
	
	session, err := h.sessionManager.GetSession(ctx, sessionID)
	if err != nil {
		if err.Error() == "session not found" {
			http.Error(w, "Session not found", http.StatusNotFound)
			return
		}
		logging.Error("failed to get session", map[string]interface{}{
			"session_id": sessionID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success": true,
		"session": session,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) UpdateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionId"])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}
	
	var req session.UpdateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	session, err := h.sessionManager.UpdateSession(ctx, sessionID, &req)
	if err != nil {
		if err.Error() == "session not found" {
			http.Error(w, "Session not found", http.StatusNotFound)
			return
		}
		logging.Error("failed to update session", map[string]interface{}{
			"session_id": sessionID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success": true,
		"session": session,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionId"])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}
	
	err = h.sessionManager.DeleteSession(ctx, sessionID)
	if err != nil {
		if err.Error() == "session not found" {
			http.Error(w, "Session not found", http.StatusNotFound)
			return
		}
		logging.Error("failed to delete session", map[string]interface{}{
			"session_id": sessionID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) JoinSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionId"])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}
	
	var req session.JoinSessionRequest
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
	
	err = h.sessionManager.JoinSession(ctx, sessionID, req.UserID, req.Role)
	if err != nil {
		logging.Error("failed to join session", map[string]interface{}{
			"session_id": sessionID,
			"user_id":    req.UserID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success": true,
		"message": "Successfully joined session",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) LeaveSession(w http.ResponseWriter, r *http.Request) {
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
	
	err = h.sessionManager.LeaveSession(ctx, sessionID, user.ID)
	if err != nil {
		logging.Error("failed to leave session", map[string]interface{}{
			"session_id": sessionID,
			"user_id":    user.ID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success": true,
		"message": "Successfully left session",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetSessionParticipants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionId"])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}
	
	participants, err := h.sessionManager.GetSessionParticipants(ctx, sessionID)
	if err != nil {
		logging.Error("failed to get session participants", map[string]interface{}{
			"session_id": sessionID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success":      true,
		"participants": participants,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}