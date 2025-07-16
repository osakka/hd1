package scene

import (
	"encoding/json"
	"net/http"
	"time"

	"holodeck1/api/shared"
	"holodeck1/logging"
	"holodeck1/sync"
)

// SceneResponse represents the current scene state
type SceneResponse struct {
	Success         bool                   `json:"success"`
	Scene           map[string]interface{} `json:"scene"`
	CurrentSequence uint64                 `json:"current_sequence"`
}

// UpdateSceneRequest represents the request to update scene properties
type UpdateSceneRequest struct {
	Background string    `json:"background,omitempty"`
	Fog        *SceneFog `json:"fog,omitempty"`
}

// SceneFog represents scene fog properties
type SceneFog struct {
	Color string  `json:"color"`
	Near  float64 `json:"near"`
	Far   float64 `json:"far"`
}

// UpdateSceneResponse represents the response after updating scene
type UpdateSceneResponse struct {
	Success bool   `json:"success"`
	SeqNum  uint64 `json:"seq_num"`
}

// GetScene handles GET /api/threejs/scene
func GetScene(w http.ResponseWriter, r *http.Request) {
	// Get hub from context
	hub := shared.GetHubFromContext(r)
	if hub == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get current sequence
	currentSeq := hub.GetSync().GetCurrentSequence()

	// Build scene state by reconstructing from operations
	// In a real implementation, you might cache this or have a scene state manager
	sceneState := map[string]interface{}{
		"entities": []interface{}{},
		"avatars":  []interface{}{},
		"background": "#87CEEB", // Default sky blue
		"fog": nil,
	}

	// Return response
	response := SceneResponse{
		Success:         true,
		Scene:           sceneState,
		CurrentSequence: currentSeq,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	logging.Debug("scene state retrieved via API", map[string]interface{}{
		"current_sequence": currentSeq,
	})
}

// UpdateScene handles PUT /api/threejs/scene
func UpdateScene(w http.ResponseWriter, r *http.Request) {
	var req UpdateSceneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get client ID
	clientID := shared.GetClientID(r)

	// Create operation data
	operationData := map[string]interface{}{}

	// Add provided updates
	if req.Background != "" {
		operationData["background"] = req.Background
	}
	if req.Fog != nil {
		operationData["fog"] = map[string]interface{}{
			"color": req.Fog.Color,
			"near":  req.Fog.Near,
			"far":   req.Fog.Far,
		}
	}

	// Only proceed if there are updates
	if len(operationData) == 0 {
		http.Error(w, "No updates provided", http.StatusBadRequest)
		return
	}

	// Create operation
	operation := &sync.Operation{
		ClientID:  clientID,
		Type:      "scene_update",
		Data:      operationData,
		Timestamp: time.Now(),
	}

	// Get hub and submit operation
	hub := shared.GetHubFromContext(r)
	if hub == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	hub.GetSync().SubmitOperation(operation)

	// Return response
	response := UpdateSceneResponse{
		Success: true,
		SeqNum:  operation.SeqNum,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	logging.Info("scene updated via API", map[string]interface{}{
		"client_id": clientID,
		"seq_num":   operation.SeqNum,
		"updates":   len(operationData),
	})
}