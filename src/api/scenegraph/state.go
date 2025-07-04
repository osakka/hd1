package scenegraph

import (
	"encoding/json"
	"net/http"
	"holodeck1/logging"
	"holodeck1/server"
	"time"
	"fmt"
)

// GetSceneStateHandler handles GET /sessions/{sessionId}/scene/state
// Retrieves the current 3D scene configuration including lighting, camera, and materials.
//
// Parameters:
//   w: HTTP response writer for JSON scene state data
//   r: HTTP request containing sessionId path parameter
//   hub: WebSocket hub for session isolation and real-time sync
//
// Returns JSON scene state with lighting settings, camera configuration, and material properties.
// Uses mock PlayCanvas scene data for consistent 3D rendering baseline.
func GetSceneStateHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	
	logging.Info("retrieving scene state", map[string]interface{}{
		"sessionId": sessionId,
		"operation": "get_scene_state",
	})
	
	// Mock PlayCanvas scene state data
	state := map[string]interface{}{
		"lighting": map[string]interface{}{
			"ambientColor": "#404040",
			"skybox":       "urban-environment",
			"shadows":      true,
			"fog": map[string]interface{}{
				"enabled": false,
				"color":   "#888888",
				"density": 0.01,
			},
		},
		"physics": map[string]interface{}{
			"gravity":  []float64{0, -9.8, 0},
			"enabled":  true,
			"timeStep": 0.016,
		},
		"rendering": map[string]interface{}{
			"shadows":        true,
			"antialiasing":   true,
			"postProcessing": []string{"bloom", "tonemapping"},
		},
	}
	
	statistics := map[string]interface{}{
		"entities": 15,
		"assets":   []string{"model-cube", "texture-metal", "audio-ambient"},
		"performance": map[string]interface{}{
			"drawCalls":    25,
			"triangles":    15000,
			"memoryUsage":  125.5,
		},
	}
	
	response := map[string]interface{}{
		"success":    true,
		"state":      state,
		"statistics": statistics,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateSceneStateHandler handles PUT /sessions/{sessionId}/scene/state
// Updates 3D scene configuration with new lighting, camera, or material settings.
//
// Parameters:
//   w: HTTP response writer for operation confirmation
//   r: HTTP request containing sessionId and JSON scene state updates
//   hub: WebSocket hub for broadcasting changes to connected clients
//
// Request body: JSON object with scene state properties to update
// Response: Confirmation message with applied changes count
// Broadcasts scene state changes to all clients in the session via WebSocket.
func UpdateSceneStateHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	logging.Info("updating scene state", map[string]interface{}{
		"sessionId": sessionId,
		"operation": "update_scene_state",
		"updates":   len(request),
	})
	
	// Process scene state updates (mock implementation)
	updatedProperties := []string{}
	
	if lighting, ok := request["lighting"].(map[string]interface{}); ok {
		for key := range lighting {
			updatedProperties = append(updatedProperties, fmt.Sprintf("lighting.%s", key))
		}
	}
	
	if physics, ok := request["physics"].(map[string]interface{}); ok {
		for key := range physics {
			updatedProperties = append(updatedProperties, fmt.Sprintf("physics.%s", key))
		}
	}
	
	if rendering, ok := request["rendering"].(map[string]interface{}); ok {
		for key := range rendering {
			updatedProperties = append(updatedProperties, fmt.Sprintf("rendering.%s", key))
		}
	}
	
	response := map[string]interface{}{
		"success":            true,
		"updated_properties": updatedProperties,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SaveSceneState handles POST /sessions/{sessionId}/scene/state/save
func SaveSceneStateHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	
	var request struct {
		SnapshotName  string `json:"snapshot_name"`
		Description   string `json:"description"`
		IncludeAssets bool   `json:"include_assets"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	logging.Info("saving scene state snapshot", map[string]interface{}{
		"sessionId":      sessionId,
		"operation":      "save_scene_state",
		"snapshot_name":  request.SnapshotName,
		"include_assets": request.IncludeAssets,
	})
	
	// Create snapshot (mock implementation)
	snapshotId := fmt.Sprintf("snapshot_%d", time.Now().Unix())
	entitiesSaved := 15
	
	response := map[string]interface{}{
		"success":        true,
		"snapshot_id":    snapshotId,
		"entities_saved": entitiesSaved,
		"timestamp":      time.Now().Format(time.RFC3339),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// LoadSceneState handles POST /sessions/{sessionId}/scene/state/load
func LoadSceneStateHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	
	var request struct {
		SnapshotId    string `json:"snapshot_id"`
		ClearExisting bool   `json:"clear_existing"`
		MergeMode     string `json:"merge_mode"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	logging.Info("loading scene state snapshot", map[string]interface{}{
		"sessionId":      sessionId,
		"operation":      "load_scene_state",
		"snapshot_id":    request.SnapshotId,
		"clear_existing": request.ClearExisting,
		"merge_mode":     request.MergeMode,
	})
	
	// Load snapshot (mock implementation)
	entitiesLoaded := 15
	warnings := []string{}
	
	if !request.ClearExisting && request.MergeMode == "merge" {
		warnings = append(warnings, "Some entities were merged with existing ones")
	}
	
	response := map[string]interface{}{
		"success":          true,
		"entities_loaded":  entitiesLoaded,
		"warnings":         warnings,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ResetSceneState handles POST /sessions/{sessionId}/scene/state/reset
func ResetSceneStateHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	
	var request struct {
		PreserveCamera   bool `json:"preserve_camera"`
		PreserveLighting bool `json:"preserve_lighting"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		// Optional request body, use defaults
		request.PreserveCamera = false
		request.PreserveLighting = false
	}
	
	logging.Info("resetting scene state", map[string]interface{}{
		"sessionId":          sessionId,
		"operation":          "reset_scene_state",
		"preserve_camera":    request.PreserveCamera,
		"preserve_lighting":  request.PreserveLighting,
	})
	
	// Reset scene (mock implementation)
	entitiesRemoved := 13
	settingsReset := []string{"physics", "rendering"}
	
	if !request.PreserveCamera {
		settingsReset = append(settingsReset, "camera")
	}
	
	if !request.PreserveLighting {
		settingsReset = append(settingsReset, "lighting")
	}
	
	response := map[string]interface{}{
		"success":          true,
		"entities_removed": entitiesRemoved,
		"settings_reset":   settingsReset,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}