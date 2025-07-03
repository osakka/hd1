package scenegraph

import (
	"encoding/json"
	"net/http"
	"holodeck1/logging"
	"holodeck1/server"
	"time"
	"fmt"
)

// ListSessionScenes handles GET /sessions/{sessionId}/scenes
func ListSessionScenesHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	
	logging.Info("listing session scenes", map[string]interface{}{
		"sessionId": sessionId,
		"operation": "list_session_scenes",
	})
	
	// Mock session scenes data
	scenes := []map[string]interface{}{
		{
			"sceneId":      "scene_main",
			"sceneName":    "Main Game Scene",
			"isActive":     true,
			"entities":     15,
			"template":     "basic",
			"createdAt":    "2025-07-01T17:15:30.123Z",
			"lastModified": "2025-07-01T17:20:45.456Z",
		},
		{
			"sceneId":      "scene_ui",
			"sceneName":    "UI Overlay Scene",
			"isActive":     false,
			"entities":     5,
			"template":     "empty",
			"createdAt":    "2025-07-01T17:16:15.789Z",
			"lastModified": "2025-07-01T17:18:20.123Z",
		},
	}
	
	response := map[string]interface{}{
		"success":      true,
		"scenes":       scenes,
		"active_scene": "scene_main",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateSessionScene handles POST /sessions/{sessionId}/scenes
func CreateSessionSceneHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	
	var request struct {
		SceneName string `json:"scene_name"`
		Template  string `json:"template"`
		SetActive bool   `json:"set_active"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// Set defaults
	if request.Template == "" {
		request.Template = "empty"
	}
	
	logging.Info("creating session scene", map[string]interface{}{
		"sessionId":  sessionId,
		"operation":  "create_session_scene",
		"scene_name": request.SceneName,
		"template":   request.Template,
		"set_active": request.SetActive,
	})
	
	// Create scene (mock implementation)
	sceneId := fmt.Sprintf("scene_%d", time.Now().Unix())
	
	response := map[string]interface{}{
		"success":    true,
		"scene_id":   sceneId,
		"scene_name": request.SceneName,
		"is_active":  request.SetActive,
	}
	
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ActivateSessionScene handles POST /sessions/{sessionId}/scenes/{sceneId}/activate
func ActivateSessionSceneHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionId(r)
	sceneId := extractSceneId(r)
	
	logging.Info("activating session scene", map[string]interface{}{
		"sessionId": sessionId,
		"sceneId":   sceneId,
		"operation": "activate_session_scene",
	})
	
	// Activate scene (mock implementation)
	previousScene := "scene_main"
	
	response := map[string]interface{}{
		"success":        true,
		"active_scene":   sceneId,
		"previous_scene": previousScene,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}