package world

import (
	"encoding/json"
	"net/http"
)

// GetWorldSpecHandler - GET /sessions/{sessionId}/world
func GetWorldSpecHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"bounds": map[string][]int{
			"x": {-12, 12},
			"y": {-12, 12}, 
			"z": {-12, 12},
		},
		"grid_size": 25,
		"total_points": 15625,
		"coordinate_system": "fixed_grid",
	})
}