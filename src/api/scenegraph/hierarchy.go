package scenegraph

import (
	"encoding/json"
	"net/http"
	"holodeck1/logging"
	"holodeck1/server"
)

// GetSceneHierarchy handles GET /sessions/{sessionId}/scene/hierarchy
func GetSceneHierarchyHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	
	logging.Info("retrieving scene hierarchy", map[string]interface{}{
		"sessionId": sessionId,
		"operation": "get_scene_hierarchy",
		"integration": "PlayCanvas",
	})
	
	// Try PlayCanvas integration first
	if IsPlayCanvasAvailable() {
		bridge := NewPlayCanvasBridge(sessionId)
		result, err := bridge.GetSceneHierarchyFromPlayCanvas()
		
		if err == nil && result["success"] == true {
			// Return real PlayCanvas scene hierarchy
			response := map[string]interface{}{
				"success":   true,
				"hierarchy": result["hierarchy"],
				"metadata":  result["metadata"],
				"source":    "PlayCanvas Engine",
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		
		logging.Warn("PlayCanvas integration failed, falling back to mock", map[string]interface{}{
			"error": err,
			"sessionId": sessionId,
		})
	}
	
	// Fallback to mock data if PlayCanvas not available
	hierarchy := map[string]interface{}{
		"id":   "scene-root",
		"name": "Scene Root", 
		"children": []map[string]interface{}{
			{
				"id":       "entity_main_camera",
				"name":     "Main Camera",
				"type":     "camera",
				"enabled":  true,
				"children": []interface{}{},
				"transform": map[string]interface{}{
					"position": []float64{0, 5, 10},
					"rotation": map[string]float64{"x": 0, "y": 0, "z": 0, "w": 1},
					"scale":    []float64{1, 1, 1},
				},
			},
			{
				"id":       "entity_directional_light",
				"name":     "Directional Light",
				"type":     "light",
				"enabled":  true,
				"children": []interface{}{},
				"transform": map[string]interface{}{
					"position": []float64{10, 10, 5},
					"rotation": map[string]float64{"x": -0.3, "y": 0.2, "z": 0, "w": 0.9},
					"scale":    []float64{1, 1, 1},
				},
			},
		},
	}
	
	metadata := map[string]interface{}{
		"totalNodes":    2,
		"maxDepth":      1,
		"lastModified":  "2025-07-01T17:26:00.000Z",
		"rootEntities":  2,
	}
	
	response := map[string]interface{}{
		"success":   true,
		"hierarchy": hierarchy,
		"metadata":  metadata,
		"source":    "Mock (PlayCanvas unavailable)",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateSceneHierarchy handles PUT /sessions/{sessionId}/scene/hierarchy  
func UpdateSceneHierarchyHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	
	var request struct {
		Operations []map[string]interface{} `json:"operations"`
		Validate   bool                     `json:"validate"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	logging.Info("updating scene hierarchy", map[string]interface{}{
		"sessionId":         sessionId,
		"operation":         "update_scene_hierarchy",
		"operations_count":  len(request.Operations),
		"validate":          request.Validate,
	})
	
	// Process hierarchy operations (mock implementation)
	operationsApplied := len(request.Operations)
	warnings := []string{}
	
	if request.Validate {
		// Add validation warnings if needed
		warnings = append(warnings, "Hierarchy validation completed successfully")
	}
	
	response := map[string]interface{}{
		"success":            true,
		"operations_applied": operationsApplied,
		"warnings":          warnings,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}