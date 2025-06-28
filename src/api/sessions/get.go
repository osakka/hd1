package sessions

import (
	"encoding/json"
	"net/http"
	"strings"
	"holodeck/server"
)

// GetSessionHandler - GET /sessions/{sessionId}
func GetSessionHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type
	h, ok := hub.(*server.Hub)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Extract session ID from URL
	sessionID := extractSessionID(r.URL.Path)
	if sessionID == "" {
		http.Error(w, `{"error": "Session ID required"}`, http.StatusBadRequest)
		return
	}
	
	// Get session from SessionStore
	session, exists := h.GetStore().GetSession(sessionID)
	if !exists {
		http.Error(w, `{"error": "Session not found"}`, http.StatusNotFound)
		return
	}
	
	// Get related data
	objects := h.GetStore().ListObjects(sessionID)
	world, worldExists := h.GetStore().GetWorld(sessionID)
	
	// Build comprehensive session response
	sessionData := map[string]interface{}{
		"id":                session.ID,
		"created_at":        session.CreatedAt,
		"status":           session.Status,
		"object_count":     len(objects),
		"world_initialized": worldExists,
		"objects":          objects,
	}
	
	if worldExists {
		sessionData["world"] = world
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessionData)
}

func extractSessionID(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 && parts[1] == "sessions" {
		return parts[2]
	}
	return ""
}