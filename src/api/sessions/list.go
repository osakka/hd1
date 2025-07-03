package sessions

import (
	"encoding/json"
	"net/http"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
)

// ListSessionsHandler - GET /sessions
func ListSessionsHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type
	h, ok := hub.(*server.Hub)
	if !ok {
		logging.Error("hub interface cast failed", map[string]interface{}{
			"expected_type": "*server.Hub",
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Get all sessions from SessionStore
	sessions := h.GetStore().ListSessions()
	logging.Debug("sessions list requested", map[string]interface{}{
		"total_sessions": len(sessions),
	})
	
	// Transform to response format with current architecture
	var sessionList []map[string]interface{}
	for _, session := range sessions {
		sessionData := map[string]interface{}{
			"id":         session.ID,
			"created_at": session.CreatedAt,
			"status":     session.Status,
			"channel_id": session.ChannelID, // Current channel joined
		}
		
		sessionList = append(sessionList, sessionData)
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"sessions":  sessionList,
		"total":     len(sessionList),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}