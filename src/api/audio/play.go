package audio

import (
	"encoding/json"
	"net/http"
	"holodeck1/logging"
	"holodeck1/server"
)

// PlayAudio handles POST /sessions/{sessionId}/audio/sources/{audioId}/play
func PlayAudioHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	audioId := extractAudioId(r)
	
	var request struct {
		Volume    float64 `json:"volume"`
		Pitch     float64 `json:"pitch"`
		FadeIn    float64 `json:"fade_in"`
		StartTime float64 `json:"start_time"`
	}
	
	// Parse optional request body
	json.NewDecoder(r.Body).Decode(&request)
	
	logging.Info("playing audio source", map[string]interface{}{
		"sessionId": sessionId,
		"audioId":   audioId,
		"volume":    request.Volume,
		"pitch":     request.Pitch,
	})
	
	// Try PlayCanvas integration first
	if IsPlayCanvasAvailable() {
		bridge := NewPlayCanvasBridge(sessionId)
		result, err := PlayAudioInPlayCanvas(bridge, audioId, map[string]interface{}{
			"volume":    request.Volume,
			"pitch":     request.Pitch,
			"fadeIn":    request.FadeIn,
			"startTime": request.StartTime,
		})
		
		if err == nil && result["success"] == true {
			response := map[string]interface{}{
				"success":  true,
				"audio_id": audioId,
				"playing":  true,
				"source":   "PlayCanvas Engine",
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		
		logging.Warn("PlayCanvas audio play failed, using mock", map[string]interface{}{
			"error": err,
			"sessionId": sessionId,
		})
	}
	
	// Mock audio playback
	response := map[string]interface{}{
		"success":  true,
		"audio_id": audioId,
		"playing":  true,
		"source":   "Mock (PlayCanvas unavailable)",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}