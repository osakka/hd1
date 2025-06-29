package sessions

import (
	"encoding/json"
	"net/http"
	"time"

	"holodeck/logging"
	"holodeck/server"
)

// ListSessionsHandler - GET /sessions
func ListSessionsHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type
	h, ok := hub.(*server.Hub)
	if !ok {
		logging.Error("failed to cast hub interface", map[string]interface{}{
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
	
	// Transform to response format with object counts
	var sessionList []map[string]interface{}
	for _, session := range sessions {
		objects := h.GetStore().ListObjects(session.ID)
		world, worldExists := h.GetStore().GetWorld(session.ID)
		
		sessionData := map[string]interface{}{
			"id":                session.ID,
			"created_at":        session.CreatedAt,
			"status":           session.Status,
			"object_count":     len(objects),
			"world_initialized": worldExists,
		}
		
		if worldExists {
			sessionData["world"] = world
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