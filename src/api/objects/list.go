package objects

import (
	"encoding/json"
	"net/http"
	"strings"
	"holodeck/server"
)

// ListObjectsHandler - GET /sessions/{sessionId}/objects
func ListObjectsHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type
	h, ok := hub.(*server.Hub)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Extract session ID from path
	sessionID := extractSessionIDFromPath(r.URL.Path)
	if sessionID == "" {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}
	
	// Validate session exists
	if _, exists := h.GetStore().GetSession(sessionID); !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	
	// Get objects from store
	objects := h.GetStore().ListObjects(sessionID)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"objects": objects,
		"total":   len(objects),
	})
}

// extractSessionIDFromPath extracts session ID from URL path
func extractSessionIDFromPath(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "sessions" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}
