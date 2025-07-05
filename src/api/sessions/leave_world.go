package sessions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
)

// LeaveSessionWorld - POST /sessions/{sessionId}/world/leave
func LeaveSessionWorldHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
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
	
	// Leave the session world
	success, clientCount := h.LeaveSessionWorld(sessionID, request.ClientID)
	
	if !success {
		http.Error(w, "Client not found in world or world does not exist", http.StatusNotFound)
		return
	}
	
	// Remove session avatar entity via API (100% API-first approach)
	avatarName := fmt.Sprintf("session_%s", request.ClientID)
	if err := h.DeleteEntityByNameViaAPI(sessionID, avatarName); err != nil {
		logging.Warn("failed to delete session avatar via API", map[string]interface{}{
			"session_id": sessionID,
			"client_id": request.ClientID,
			"avatar_name": avatarName,
			"error": err.Error(),
		})
		// Continue - avatar deletion failure shouldn't block world leave
	} else {
		logging.Info("session avatar deleted via API", map[string]interface{}{
			"session_id": sessionID,
			"client_id": request.ClientID,
			"avatar_name": avatarName,
		})
	}
	
	logging.Info("client left session world", map[string]interface{}{
		"session_id":   sessionID,
		"client_id":    request.ClientID,
		"client_count": clientCount,
	})
	
	// Broadcast world leave event to remaining clients in the world
	h.BroadcastToSessionWorld(sessionID, "client_left", map[string]interface{}{
		"client_id":    request.ClientID,
		"client_count": clientCount,
		"left_at":      time.Now(),
	}, "")
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Left session world",
	})
}

