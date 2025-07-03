package animation

import (
	"encoding/json"
	"net/http"
	"holodeck1/logging"
	"holodeck1/server"
)

// PlayAnimation handles POST /sessions/{sessionId}/animations/{animationId}/play
func PlayAnimationHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	animationId := extractAnimationId(r)
	
	var request struct {
		Speed     float64 `json:"speed"`
		StartTime float64 `json:"start_time"`
		BlendMode string  `json:"blend_mode"`
		FadeIn    float64 `json:"fade_in"`
	}
	
	// Parse optional request body
	json.NewDecoder(r.Body).Decode(&request)
	
	logging.Info("playing animation", map[string]interface{}{
		"sessionId":   sessionId,
		"animationId": animationId,
		"speed":       request.Speed,
		"blendMode":   request.BlendMode,
	})
	
	// Try PlayCanvas integration first
	if IsPlayCanvasAvailable() {
		bridge := NewPlayCanvasBridge(sessionId)
		result, err := PlayAnimationInPlayCanvas(bridge, animationId, map[string]interface{}{
			"speed":      request.Speed,
			"startTime":  request.StartTime,
			"blendMode":  request.BlendMode,
			"fadeIn":     request.FadeIn,
		})
		
		if err == nil && result["success"] == true {
			response := map[string]interface{}{
				"success":      true,
				"animation_id": animationId,
				"playing":      true,
				"start_time":   request.StartTime,
				"source":       "PlayCanvas Engine",
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		
		logging.Warn("PlayCanvas animation play failed, using mock", map[string]interface{}{
			"error": err,
			"sessionId": sessionId,
		})
	}
	
	// Mock animation playback
	response := map[string]interface{}{
		"success":      true,
		"animation_id": animationId,
		"playing":      true,
		"start_time":   request.StartTime,
		"source":       "Mock (PlayCanvas unavailable)",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}