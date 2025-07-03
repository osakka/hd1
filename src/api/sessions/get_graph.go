package sessions

import (
	"encoding/json"
	"net/http"
	"strings"

	"holodeck1/logging"
	"holodeck1/server"
)

// GetSessionGraph - GET /sessions/{sessionId}/room/graph
func GetSessionGraphHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type
	h, ok := hub.(*server.Hub)
	if !ok {
		logging.Error("hub interface cast failed", map[string]interface{}{
			"expected_type": "*server.Hub",
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Extract session ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}
	sessionID := pathParts[2]
	
	// Verify session exists
	_, exists := h.GetStore().GetSession(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	
	// Get session graph state
	graphState, clientCount, lastUpdated := h.GetSessionGraphState(sessionID)
	
	// Get all entities in the session from channel
	entities := []interface{}{} // Entities managed via PlayCanvas/channels
	
	logging.Debug("session graph state retrieved", map[string]interface{}{
		"session_id":   sessionID,
		"entity_count": len(entities),
		"client_count": clientCount,
	})
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"session_id": sessionID,
		"graph_state": map[string]interface{}{
			"entities":     entities,
			"properties":   graphState,
			"last_updated": lastUpdated,
		},
		"client_count": clientCount,
	})
}