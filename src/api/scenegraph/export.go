package scenegraph

import (
	"encoding/json"
	"net/http"
	"holodeck1/logging"
	"holodeck1/server"
	"time"
)

// ExportSceneDefinition handles GET /sessions/{sessionId}/scene/export
func ExportSceneDefinitionHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	
	// Parse query parameters
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}
	
	includeAssets := r.URL.Query().Get("include_assets") == "true"
	
	logging.Info("exporting scene definition", map[string]interface{}{
		"sessionId":      sessionId,
		"operation":      "export_scene_definition",
		"format":         format,
		"include_assets": includeAssets,
	})
	
	// Mock scene export data
	exportData := map[string]interface{}{
		"version": "3.0.0",
		"scene": map[string]interface{}{
			"name": "Exported Scene",
			"settings": map[string]interface{}{
				"lighting": map[string]interface{}{
					"ambientColor": "#404040",
					"skybox":       "urban-environment",
				},
				"physics": map[string]interface{}{
					"gravity": []float64{0, -9.8, 0},
					"enabled": true,
				},
			},
			"entities": []map[string]interface{}{
				{
					"id":   "entity_camera",
					"name": "Main Camera",
					"components": map[string]interface{}{
						"camera": map[string]interface{}{
							"clearColor": []float64{0.2, 0.2, 0.3, 1.0},
							"fov":        60,
						},
					},
					"transform": map[string]interface{}{
						"position": []float64{0, 5, 10},
						"rotation": map[string]float64{"x": 0, "y": 0, "z": 0, "w": 1},
						"scale":    []float64{1, 1, 1},
					},
				},
				{
					"id":   "entity_light",
					"name": "Directional Light",
					"components": map[string]interface{}{
						"light": map[string]interface{}{
							"type":      "directional",
							"color":     []float64{1, 1, 1},
							"intensity": 1.0,
						},
					},
					"transform": map[string]interface{}{
						"position": []float64{10, 10, 5},
						"rotation": map[string]float64{"x": -0.3, "y": 0.2, "z": 0, "w": 0.9},
						"scale":    []float64{1, 1, 1},
					},
				},
			},
		},
	}
	
	// Add assets if requested
	if includeAssets {
		exportData["assets"] = []map[string]interface{}{
			{
				"id":   "model-cube",
				"type": "model",
				"url":  "/assets/models/cube.glb",
			},
			{
				"id":   "texture-metal",
				"type": "texture",
				"url":  "/assets/textures/metal.jpg",
			},
		}
	}
	
	metadata := map[string]interface{}{
		"format":     format,
		"version":    "3.0.0",
		"timestamp":  time.Now().Format(time.RFC3339),
		"entities":   2,
		"assets":     2,
		"size":       524288,
	}
	
	response := map[string]interface{}{
		"success":     true,
		"export_data": exportData,
		"metadata":    metadata,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}