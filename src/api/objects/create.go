package objects

import (
	"encoding/json"
	"net/http"
	"strings"
	"holodeck/server"
)

// Color represents RGBA color values
type Color struct {
	R float64 `json:"r"`
	G float64 `json:"g"`
	B float64 `json:"b"`
	A float64 `json:"a"`
}

// CreateObjectRequest represents the request body for creating objects
type CreateObjectRequest struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Z     float64 `json:"z"`
	Color *Color  `json:"color,omitempty"`
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
	
	// Use provided color or default to green
	objectColor := map[string]interface{}{
		"r": 0.2,
		"g": 0.8,
		"b": 0.2,
		"a": 1.0,
	}
	if req.Color != nil {
		objectColor = map[string]interface{}{
			"r": req.Color.R,
			"g": req.Color.G,
			"b": req.Color.B,
			"a": req.Color.A,
		}
	}

	// Broadcast object creation for real-time updates via canvas control
	h.BroadcastUpdate("canvas_control", map[string]interface{}{
		"command": "create",
		"objects": []map[string]interface{}{
			{
				"id":   object.Name,
				"name": object.Name,
				"type": object.Type,
				"transform": map[string]interface{}{
					"position": map[string]interface{}{
						"x": object.X,
						"y": object.Y,
						"z": object.Z,
					},
					"scale": map[string]interface{}{
						"x": object.Scale,
						"y": object.Scale,
						"z": object.Scale,
					},
					"rotation": map[string]interface{}{
						"x": 0,
						"y": 0,
						"z": 0,
					},
				},
				"color": objectColor,
				"material": map[string]interface{}{
					"shader":     "standard",
					"metalness":  0.1,
					"roughness":  0.7,
					"transparent": false,
				},
				"physics": map[string]interface{}{
					"enabled": false,
					"mass":    1.0,
					"type":    "static",
				},
				"lighting": map[string]interface{}{
					"castShadow":    true,
					"receiveShadow": true,
				},
				"visible":   true,
				"wireframe": false,
			},
		},
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