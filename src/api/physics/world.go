package physics

import (
	"encoding/json"
	"net/http"
	"holodeck1/logging"
	"holodeck1/server"
)

// GetPhysicsWorld handles GET /sessions/{sessionId}/physics/world
func GetPhysicsWorldHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	
	logging.Info("getting physics world", map[string]interface{}{
		"sessionId": sessionId,
		"operation": "get_physics_world",
	})
	
	// Try PlayCanvas integration first
	if IsPlayCanvasAvailable() {
		bridge := NewPlayCanvasBridge(sessionId)
		result, err := GetPhysicsWorldFromPlayCanvas(bridge)
		
		if err == nil && result["success"] == true {
			response := map[string]interface{}{
				"success":       true,
				"physics_world": result["physics_world"],
				"source":        "PlayCanvas Engine",
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		
		logging.Warn("PlayCanvas physics world integration failed, falling back to mock", map[string]interface{}{
			"error": err,
			"sessionId": sessionId,
		})
	}
	
	// Fallback to mock physics world
	physicsWorld := map[string]interface{}{
		"gravity":           []float64{0, -9.8, 0},
		"timestep":          0.016,
		"enabled":           true,
		"solver_iterations": 10,
		"active_bodies":     5,
	}
	
	response := map[string]interface{}{
		"success":       true,
		"physics_world": physicsWorld,
		"source":        "Mock (PlayCanvas unavailable)",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdatePhysicsWorld handles PUT /sessions/{sessionId}/physics/world
func UpdatePhysicsWorldHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	
	var request struct {
		Gravity          []float64 `json:"gravity"`
		Timestep         float64   `json:"timestep"`
		Enabled          bool      `json:"enabled"`
		SolverIterations int       `json:"solver_iterations"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	logging.Info("updating physics world", map[string]interface{}{
		"sessionId": sessionId,
		"gravity":   request.Gravity,
		"enabled":   request.Enabled,
	})
	
	// Try PlayCanvas integration first
	if IsPlayCanvasAvailable() {
		bridge := NewPlayCanvasBridge(sessionId)
		result, err := UpdatePhysicsWorldInPlayCanvas(bridge, map[string]interface{}{
			"gravity":          request.Gravity,
			"timestep":         request.Timestep,
			"enabled":          request.Enabled,
			"solverIterations": request.SolverIterations,
		})
		
		if err == nil && result["success"] == true {
			response := map[string]interface{}{
				"success":            true,
				"updated_properties": result["updated_properties"],
				"source":             "PlayCanvas Engine",
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		
		logging.Warn("PlayCanvas physics world update failed, using mock", map[string]interface{}{
			"error": err,
			"sessionId": sessionId,
		})
	}
	
	// Mock physics world update
	updatedProperties := []string{}
	if len(request.Gravity) > 0 {
		updatedProperties = append(updatedProperties, "gravity")
	}
	if request.Timestep > 0 {
		updatedProperties = append(updatedProperties, "timestep")
	}
	
	response := map[string]interface{}{
		"success":            true,
		"updated_properties": updatedProperties,
		"source":             "Mock (PlayCanvas unavailable)",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}