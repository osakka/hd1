package physics

import (
	"encoding/json"
	"net/http"
	"holodeck1/logging"
	"holodeck1/server"
)

// ListRigidBodies handles GET /sessions/{sessionId}/physics/rigidbodies
func ListRigidBodiesHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	
	logging.Info("listing rigid bodies", map[string]interface{}{
		"sessionId": sessionId,
		"operation": "list_rigid_bodies",
	})
	
	// Try PlayCanvas integration first
	if IsPlayCanvasAvailable() {
		bridge := NewPlayCanvasBridge(sessionId)
		result, err := GetRigidBodiesFromPlayCanvas(bridge)
		
		if err == nil && result["success"] == true {
			response := map[string]interface{}{
				"success":      true,
				"rigid_bodies": result["rigid_bodies"],
				"source":       "PlayCanvas Engine",
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		
		logging.Warn("PlayCanvas rigid bodies integration failed, falling back to mock", map[string]interface{}{
			"error": err,
			"sessionId": sessionId,
		})
	}
	
	// Fallback to mock rigid bodies
	rigidBodies := []map[string]interface{}{
		{
			"entity_id":        "entity-demo-cube",
			"name":             "Demo Cube",
			"mass":             10.0,
			"velocity":         []float64{0, 0, 0},
			"angular_velocity": []float64{0, 0, 0},
			"body_type":        "dynamic",
			"sleeping":         false,
		},
	}
	
	response := map[string]interface{}{
		"success":      true,
		"rigid_bodies": rigidBodies,
		"source":       "Mock (PlayCanvas unavailable)",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}