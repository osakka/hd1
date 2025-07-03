package animation

import (
	"encoding/json"
	"net/http"
	"holodeck1/logging"
	"holodeck1/server"
)

// CreateAnimation handles POST /sessions/{sessionId}/animations
func CreateAnimationHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	
	var request struct {
		Name      string `json:"name"`
		EntityId  string `json:"entity_id"`
		Duration  float64 `json:"duration"`
		Loop      bool   `json:"loop"`
		Keyframes []map[string]interface{} `json:"keyframes"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	logging.Info("creating animation", map[string]interface{}{
		"sessionId": sessionId,
		"name":      request.Name,
		"entityId":  request.EntityId,
		"duration":  request.Duration,
	})
	
	// Try PlayCanvas integration first
	if IsPlayCanvasAvailable() {
		bridge := NewPlayCanvasBridge(sessionId)
		result, err := CreateAnimationInPlayCanvas(bridge, map[string]interface{}{
			"name":      request.Name,
			"entityId":  request.EntityId,
			"duration":  request.Duration,
			"loop":      request.Loop,
			"keyframes": request.Keyframes,
		})
		
		if err == nil && result["success"] == true {
			response := map[string]interface{}{
				"success":      true,
				"animation_id": result["animation_id"],
				"name":         request.Name,
				"duration":     request.Duration,
				"source":       "PlayCanvas Engine",
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		
		logging.Warn("PlayCanvas animation creation failed, using mock", map[string]interface{}{
			"error": err,
			"sessionId": sessionId,
		})
	}
	
	// Mock animation creation
	response := map[string]interface{}{
		"success":      true,
		"animation_id": "anim-" + sessionId + "-" + request.Name,
		"name":         request.Name,
		"duration":     request.Duration,
		"source":       "Mock (PlayCanvas unavailable)",
	}
	
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}