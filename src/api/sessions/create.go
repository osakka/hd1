package sessions

import (
	"encoding/json"
	"net/http"

	"holodeck1/logging"
	"holodeck1/server"
)

// CreateSessionHandler - POST /sessions
func CreateSessionHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type
	h, ok := hub.(*server.Hub)
	if !ok {
		logging.Error("hub interface cast failed", map[string]interface{}{
			"expected_type": "*server.Hub",
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Create session using SessionStore
	session := h.GetStore().CreateSession()
	logging.Info("session created successfully", map[string]interface{}{
		"session_id": session.ID,
		"created_at": session.CreatedAt,
	})
	
	// Broadcast session creation for real-time updates
	h.BroadcastUpdate("session_created", map[string]interface{}{
		"session": session,
	})
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"session_id": session.ID,
		"created_at": session.CreatedAt,
		"status":     session.Status,
		"message":    "Session created - ready to join channels",
	})
}

