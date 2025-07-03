package animation

import (
	"encoding/json"
	"net/http"
	"holodeck1/logging"
	"holodeck1/server"
)

// StopAnimation handles POST /sessions/{sessionId}/animations/{animationId}/stop
func StopAnimationHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	animationId := extractAnimationId(r)
	
	var request struct {
		FadeOut float64 `json:"fade_out"`
	}
	
	// Parse optional request body
	json.NewDecoder(r.Body).Decode(&request)
	
	logging.Info("stopping animation", map[string]interface{}{
		"sessionId":   sessionId,
		"animationId": animationId,
		"fadeOut":     request.FadeOut,
	})
	
	// Try PlayCanvas integration first
	if IsPlayCanvasAvailable() {
		bridge := NewPlayCanvasBridge(sessionId)
		result, err := StopAnimationInPlayCanvas(bridge, animationId, map[string]interface{}{
			"fadeOut": request.FadeOut,
		})
		
		if err == nil && result["success"] == true {
			response := map[string]interface{}{
				"success":      true,
				"animation_id": animationId,
				"playing":      false,
				"source":       "PlayCanvas Engine",
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		
		logging.Warn("PlayCanvas animation stop failed, using mock", map[string]interface{}{
			"error": err,
			"sessionId": sessionId,
		})
	}
	
	// Mock animation stop
	response := map[string]interface{}{
		"success":      true,
		"animation_id": animationId,
		"playing":      false,
		"source":       "Mock (PlayCanvas unavailable)",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}