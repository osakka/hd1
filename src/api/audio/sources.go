package audio

import (
	"encoding/json"
	"net/http"
	"holodeck1/logging"
	"holodeck1/server"
)

// ListAudioSources handles GET /sessions/{sessionId}/audio/sources
func ListAudioSourcesHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	
	logging.Info("listing audio sources", map[string]interface{}{
		"sessionId": sessionId,
		"operation": "list_audio_sources",
	})
	
	// Try PlayCanvas integration first
	if IsPlayCanvasAvailable() {
		bridge := NewPlayCanvasBridge(sessionId)
		result, err := GetAudioSourcesFromPlayCanvas(bridge)
		
		if err == nil && result["success"] == true {
			response := map[string]interface{}{
				"success":       true,
				"audio_sources": result["audio_sources"],
				"source":        "PlayCanvas Engine",
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		
		logging.Warn("PlayCanvas audio sources integration failed, falling back to mock", map[string]interface{}{
			"error": err,
			"sessionId": sessionId,
		})
	}
	
	// Fallback to mock audio sources
	audioSources := []map[string]interface{}{
		{
			"id":           "audio-demo-engine",
			"entity_id":    "entity-demo",
			"name":         "Demo Engine Sound",
			"playing":      false,
			"volume":       0.8,
			"pitch":        1.0,
			"loop":         true,
			"positional":   true,
			"max_distance": 50.0,
		},
	}
	
	response := map[string]interface{}{
		"success":       true,
		"audio_sources": audioSources,
		"source":        "Mock (PlayCanvas unavailable)",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateAudioSource handles POST /sessions/{sessionId}/audio/sources
func CreateAudioSourceHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	
	var request struct {
		Name        string  `json:"name"`
		EntityId    string  `json:"entity_id"`
		AudioClip   string  `json:"audio_clip"`
		Volume      float64 `json:"volume"`
		Pitch       float64 `json:"pitch"`
		Loop        bool    `json:"loop"`
		Positional  bool    `json:"positional"`
		MaxDistance float64 `json:"max_distance"`
		AutoPlay    bool    `json:"auto_play"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	logging.Info("creating audio source", map[string]interface{}{
		"sessionId":  sessionId,
		"name":       request.Name,
		"entityId":   request.EntityId,
		"audioClip":  request.AudioClip,
		"positional": request.Positional,
	})
	
	// Try PlayCanvas integration first
	if IsPlayCanvasAvailable() {
		bridge := NewPlayCanvasBridge(sessionId)
		result, err := CreateAudioSourceInPlayCanvas(bridge, map[string]interface{}{
			"name":        request.Name,
			"entityId":    request.EntityId,
			"audioClip":   request.AudioClip,
			"volume":      request.Volume,
			"pitch":       request.Pitch,
			"loop":        request.Loop,
			"positional":  request.Positional,
			"maxDistance": request.MaxDistance,
			"autoPlay":    request.AutoPlay,
		})
		
		if err == nil && result["success"] == true {
			response := map[string]interface{}{
				"success":   true,
				"audio_id":  result["audio_id"],
				"name":      request.Name,
				"entity_id": request.EntityId,
				"source":    "PlayCanvas Engine",
			}
			
			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		
		logging.Warn("PlayCanvas audio source creation failed, using mock", map[string]interface{}{
			"error": err,
			"sessionId": sessionId,
		})
	}
	
	// Mock audio source creation
	response := map[string]interface{}{
		"success":   true,
		"audio_id":  "audio-" + sessionId + "-" + request.Name,
		"name":      request.Name,
		"entity_id": request.EntityId,
		"source":    "Mock (PlayCanvas unavailable)",
	}
	
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}