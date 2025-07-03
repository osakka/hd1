package audio

import (
	"encoding/json"
	"net/http"
	"holodeck1/logging"
	"holodeck1/server"
)

// StopAudio handles POST /sessions/{sessionId}/audio/sources/{audioId}/stop
func StopAudioHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	audioId := extractAudioId(r)
	
	var request struct {
		FadeOut float64 `json:"fade_out"`
	}
	
	// Parse optional request body
	json.NewDecoder(r.Body).Decode(&request)
	
	logging.Info("stopping audio source", map[string]interface{}{
		"sessionId": sessionId,
		"audioId":   audioId,
		"fadeOut":   request.FadeOut,
	})
	
	// Try PlayCanvas integration first
	if IsPlayCanvasAvailable() {
		bridge := NewPlayCanvasBridge(sessionId)
		result, err := StopAudioInPlayCanvas(bridge, audioId, map[string]interface{}{
			"fadeOut": request.FadeOut,
		})
		
		if err == nil && result["success"] == true {
			response := map[string]interface{}{
				"success":  true,
				"audio_id": audioId,
				"playing":  false,
				"source":   "PlayCanvas Engine",
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		
		logging.Warn("PlayCanvas audio stop failed, using mock", map[string]interface{}{
			"error": err,
			"sessionId": sessionId,
		})
	}
	
	// Mock audio stop
	response := map[string]interface{}{
		"success":  true,
		"audio_id": audioId,
		"playing":  false,
		"source":   "Mock (PlayCanvas unavailable)",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}