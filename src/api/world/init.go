package world

import (
	"encoding/json"
	"net/http"
	"strings"
	"holodeck/server"
)

// InitializeWorldRequest represents the request body for world initialization
type InitializeWorldRequest struct {
	Size                int     `json:"size"`
	GridTransparency    float64 `json:"grid_transparency"`     // Internal grid lines
	SurfaceTransparency float64 `json:"surface_transparency"`  // Boundary surfaces
	CameraX             float64 `json:"camera_x"`
	CameraY             float64 `json:"camera_y"`
	CameraZ             float64 `json:"camera_z"`
}

// InitializeWorldHandler - POST /sessions/{sessionId}/world
func InitializeWorldHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
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
	
	// Parse request body with defaults
	var req InitializeWorldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Use defaults if no body provided
		req = InitializeWorldRequest{
			Size:                25,
			GridTransparency:    0.01,  // 99% transparent grid
			SurfaceTransparency: 0.3,   // 70% transparent surfaces
			CameraX:             10,
			CameraY:             10,
			CameraZ:             10,
		}
	}
	
	// Validate world size
	if req.Size < 5 || req.Size > 50 {
		req.Size = 25 // Default to 25
	}
	
	// Initialize world using SessionStore
	world, err := h.GetStore().InitializeWorld(sessionID, req.Size, req.GridTransparency, req.CameraX, req.CameraY, req.CameraZ)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Broadcast world initialization for real-time updates
	h.BroadcastUpdate("world_initialized", map[string]interface{}{
		"session_id":           sessionID,
		"world":               world,
		"grid_transparency":    req.GridTransparency,
		"surface_transparency": req.SurfaceTransparency,
	})
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":           true,
		"world":            world,
		"bounds":           map[string]int{"min": -12, "max": 12},
		"coordinate_system": "fixed_grid",
		"message":          "World initialized - ready for objects",
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