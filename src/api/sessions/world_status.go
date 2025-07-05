package sessions

import (
	"encoding/json"
	"net/http"
	"strings"

	"holodeck1/logging"
	"holodeck1/server"
)

// GetSessionWorldStatus - GET /sessions/{sessionId}/world/status
func GetSessionWorldStatusHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
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
	
	// Get world status
	worldStatus := h.GetSessionWorldStatus(sessionID)
	
	logging.Debug("session world status retrieved", map[string]interface{}{
		"session_id":  sessionID,
		"world_active": worldStatus.WorldActive,
		"client_count": len(worldStatus.ConnectedClients),
	})
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":            true,
		"session_id":         sessionID,
		"world_active":       worldStatus.WorldActive,
		"connected_clients":  worldStatus.ConnectedClients,
		"graph_summary":      worldStatus.GraphSummary,
		"health_metrics":     worldStatus.HealthMetrics,
	})
}

