package sessions

import (
	"encoding/json"
	"net/http"
	"strings"
	"holodeck/server"
)

// DeleteSessionHandler - DELETE /sessions/{sessionId}
func DeleteSessionHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type
	h, ok := hub.(*server.Hub)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	sessionID := extractSessionIDFromPath(r.URL.Path)
	if sessionID == "" {
		http.Error(w, `{"error": "Session ID required"}`, http.StatusBadRequest)
		return
	}
	
	// Check if session exists before deleting
	if _, exists := h.GetStore().GetSession(sessionID); !exists {
		http.Error(w, `{"error": "Session not found"}`, http.StatusNotFound)
		return
	}
	
	// Delete session from SessionStore (cascades to objects and world)
	deleted := h.GetStore().DeleteSession(sessionID)
	if !deleted {
		http.Error(w, `{"error": "Failed to delete session"}`, http.StatusInternalServerError)
		return
	}
	
	// Broadcast session deletion for real-time updates
	h.BroadcastUpdate("session_deleted", map[string]interface{}{
		"session_id": sessionID,
	})
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Session terminated and all data removed",
		"session_id": sessionID,
	})
}

func extractSessionIDFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 && parts[1] == "sessions" {
		return parts[2]
	}
	return ""
}