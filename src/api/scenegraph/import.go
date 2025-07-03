package scenegraph

import (
	"encoding/json"
	"net/http"
	"holodeck1/logging"
	"holodeck1/server"
)

// ImportSceneDefinition handles POST /sessions/{sessionId}/scene/import
func ImportSceneDefinitionHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	
	var request struct {
		SceneData map[string]interface{} `json:"scene_data"`
		Format    string                 `json:"format"`
		MergeMode string                 `json:"merge_mode"`
		Validate  bool                   `json:"validate"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// Set defaults
	if request.Format == "" {
		request.Format = "json"
	}
	if request.MergeMode == "" {
		request.MergeMode = "replace"
	}
	
	logging.Info("importing scene definition", map[string]interface{}{
		"sessionId":  sessionId,
		"operation":  "import_scene_definition",
		"format":     request.Format,
		"merge_mode": request.MergeMode,
		"validate":   request.Validate,
	})
	
	// Validate scene data if requested
	warnings := []string{}
	if request.Validate {
		// Mock validation process
		if _, ok := request.SceneData["scene"]; !ok {
			http.Error(w, "Invalid scene data: missing 'scene' property", http.StatusBadRequest)
			return
		}
		warnings = append(warnings, "Scene data validation completed successfully")
	}
	
	// Process import (mock implementation)
	entitiesCreated := 0
	
	if scene, ok := request.SceneData["scene"].(map[string]interface{}); ok {
		if entities, ok := scene["entities"].([]interface{}); ok {
			entitiesCreated = len(entities)
		}
	}
	
	// Add merge mode specific warnings
	switch request.MergeMode {
	case "merge":
		warnings = append(warnings, "Entities merged with existing scene")
	case "append":
		warnings = append(warnings, "Entities appended to existing scene")
	case "replace":
		warnings = append(warnings, "Existing scene replaced with imported data")
	}
	
	response := map[string]interface{}{
		"success":          true,
		"entities_created": entitiesCreated,
		"warnings":         warnings,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}