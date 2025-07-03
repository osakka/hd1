package sessions

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
)

// UpdateSessionGraph - PUT /sessions/{sessionId}/room/graph
func UpdateSessionGraphHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
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
	
	// Parse request body
	var request struct {
		GraphUpdates map[string]interface{} `json:"graph_updates"`
		ClientID     string                 `json:"client_id"`
		Atomic       bool                   `json:"atomic"`
	}
	
	// Set default values
	request.Atomic = true
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	if request.ClientID == "" {
		http.Error(w, "client_id is required", http.StatusBadRequest)
		return
	}
	
	if request.GraphUpdates == nil {
		http.Error(w, "graph_updates is required", http.StatusBadRequest)
		return
	}
	
	// Verify session exists
	_, exists := h.GetStore().GetSession(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	
	// Update session graph state
	broadcastCount, err := h.UpdateSessionGraphState(sessionID, request.ClientID, request.GraphUpdates, request.Atomic)
	if err != nil {
		logging.Error("session graph update failed", map[string]interface{}{
			"session_id": sessionID,
			"client_id":  request.ClientID,
			"error":      err.Error(),
		})
		http.Error(w, "Failed to update graph state: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	updatedAt := time.Now()
	
	logging.Info("session graph state updated", map[string]interface{}{
		"session_id":      sessionID,
		"client_id":       request.ClientID,
		"broadcast_count": broadcastCount,
		"atomic":          request.Atomic,
	})
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":         true,
		"updated_at":      updatedAt,
		"broadcast_count": broadcastCount,
		"message":         "Graph state updated successfully",
	})
}