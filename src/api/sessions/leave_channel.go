package sessions

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
)

// LeaveSessionChannel - POST /sessions/{sessionId}/channel/leave
func LeaveSessionChannelHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
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
		ClientID string `json:"client_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	if request.ClientID == "" {
		http.Error(w, "client_id is required", http.StatusBadRequest)
		return
	}
	
	// Leave the session channel
	success, clientCount := h.LeaveSessionChannel(sessionID, request.ClientID)
	
	if !success {
		http.Error(w, "Client not found in channel or channel does not exist", http.StatusNotFound)
		return
	}
	
	logging.Info("client left session channel", map[string]interface{}{
		"session_id":   sessionID,
		"client_id":    request.ClientID,
		"client_count": clientCount,
	})
	
	// Broadcast channel leave event to remaining clients in the channel
	h.BroadcastToSessionChannel(sessionID, "client_left", map[string]interface{}{
		"client_id":    request.ClientID,
		"client_count": clientCount,
		"left_at":      time.Now(),
	}, "")
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Successfully left session channel",
	})
}