package objects

import (
	"encoding/json"
	"net/http"
	"strings"
	"holodeck/server"
)

// CreateObjectRequest represents the request body for creating objects
type CreateObjectRequest struct {
	Name string  `json:"name"`
	Type string  `json:"type"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Z    float64 `json:"z"`
}

// CreateObjectHandler - POST /sessions/{sessionId}/objects
func CreateObjectHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type
	h, ok := hub.(*server.Hub)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Extract session ID from path
	sessionID := extractSessionID(r.URL.Path)
	if sessionID == "" {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}
	
	// Validate session exists
	if _, exists := h.GetStore().GetSession(sessionID); !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	
	// Parse request body
	var req CreateObjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}
	
	// Validate required fields
	if req.Name == "" || req.Type == "" {
		http.Error(w, "Missing required fields: name, type", http.StatusBadRequest)
		return
	}
	
	// Create object using SessionStore with coordinate validation
	object, err := h.GetStore().CreateObject(sessionID, req.Name, req.Type, req.X, req.Y, req.Z)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Broadcast object creation for real-time updates
	h.BroadcastUpdate("object_created", map[string]interface{}{
		"session_id": sessionID,
		"object":     object,
	})
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"object":  object,
		"message": "Object created successfully",
	})
}

// extractSessionID extracts session ID from URL path
func extractSessionID(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "sessions" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}