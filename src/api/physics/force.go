package physics

import (
	"encoding/json"
	"net/http"
	"holodeck1/logging"
	"holodeck1/server"
)

// ApplyForce handles POST /sessions/{sessionId}/physics/rigidbodies/{entityId}/force
func ApplyForceHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	entityId := extractEntityId(r)
	
	var request struct {
		Force    []float64 `json:"force"`
		Position []float64 `json:"position"`
		Mode     string    `json:"mode"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	logging.Info("applying force to rigid body", map[string]interface{}{
		"sessionId": sessionId,
		"entityId":  entityId,
		"force":     request.Force,
		"mode":      request.Mode,
	})
	
	// Try PlayCanvas integration first
	if IsPlayCanvasAvailable() {
		bridge := NewPlayCanvasBridge(sessionId)
		result, err := ApplyForceInPlayCanvas(bridge, entityId, map[string]interface{}{
			"force":    request.Force,
			"position": request.Position,
			"mode":     request.Mode,
		})
		
		if err == nil && result["success"] == true {
			response := map[string]interface{}{
				"success":       true,
				"entity_id":     entityId,
				"force_applied": request.Force,
				"new_velocity":  result["new_velocity"],
				"source":        "PlayCanvas Engine",
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		
		logging.Warn("PlayCanvas force application failed, using mock", map[string]interface{}{
			"error": err,
			"sessionId": sessionId,
		})
	}
	
	// Mock force application
	response := map[string]interface{}{
		"success":       true,
		"entity_id":     entityId,
		"force_applied": request.Force,
		"new_velocity":  []float64{5.2, 0, 0}, // Mock resulting velocity
		"source":        "Mock (PlayCanvas unavailable)",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}