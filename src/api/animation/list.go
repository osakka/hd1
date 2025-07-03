package animation

import (
	"encoding/json"
	"net/http"
	"holodeck1/logging"
	"holodeck1/server"
)

// ListAnimations handles GET /sessions/{sessionId}/animations
func ListAnimationsHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	
	logging.Info("listing animations", map[string]interface{}{
		"sessionId": sessionId,
		"operation": "list_animations",
	})
	
	// Try PlayCanvas integration first
	if IsPlayCanvasAvailable() {
		bridge := NewPlayCanvasBridge(sessionId)
		result, err := GetAnimationsFromPlayCanvas(bridge)
		
		if err == nil && result["success"] == true {
			response := map[string]interface{}{
				"success":    true,
				"animations": result["animations"],
				"source":     "PlayCanvas Engine",
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		
		logging.Warn("PlayCanvas animation integration failed, falling back to mock", map[string]interface{}{
			"error": err,
			"sessionId": sessionId,
		})
	}
	
	// Fallback to mock animations
	animations := []map[string]interface{}{
		{
			"id":           "anim-demo-rotation",
			"name":         "Demo Rotation",
			"entity_id":    "entity-demo",
			"duration":     3.0,
			"loop":         true,
			"playing":      false,
			"current_time": 0.0,
		},
	}
	
	response := map[string]interface{}{
		"success":    true,
		"animations": animations,
		"source":     "Mock (PlayCanvas unavailable)",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}