package threejs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"holodeck1/logging"
	"holodeck1/sync"
)

// MoveAvatarRequest represents the request to move an avatar
type MoveAvatarRequest struct {
	Position  Vector3 `json:"position"`
	Rotation  *Vector3 `json:"rotation,omitempty"`
	Animation string  `json:"animation,omitempty"`
}

// MoveAvatarResponse represents the response after moving an avatar
type MoveAvatarResponse struct {
	Success bool   `json:"success"`
	SeqNum  uint64 `json:"seq_num"`
}

// MoveAvatar handles POST /api/threejs/avatars/{sessionId}/move
func MoveAvatar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	if sessionID == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	var req MoveAvatarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate animation if provided
	if req.Animation != "" {
		validAnimations := map[string]bool{
			"idle": true,
			"walk": true,
			"run":  true,
		}
		if !validAnimations[req.Animation] {
			http.Error(w, "Invalid animation", http.StatusBadRequest)
			return
		}
	}

	// Get client ID
	clientID := getClientID(r)

	// Create operation data
	operationData := map[string]interface{}{
		"session_id": sessionID,
		"position":   req.Position,
	}

	// Add optional properties
	if req.Rotation != nil {
		operationData["rotation"] = req.Rotation
	}
	if req.Animation != "" {
		operationData["animation"] = req.Animation
	}

	// Create operation
	operation := &sync.Operation{
		ClientID:  clientID,
		Type:      "avatar_move",
		Data:      operationData,
		Timestamp: time.Now(),
	}

	// Get hub and submit operation
	hub := getHubFromContext(r)
	if hub == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	hub.GetSync().SubmitOperation(operation)

	// Return response
	response := MoveAvatarResponse{
		Success: true,
		SeqNum:  operation.SeqNum,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	logging.Debug("avatar moved via API", map[string]interface{}{
		"session_id": sessionID,
		"client_id":  clientID,
		"seq_num":    operation.SeqNum,
		"position":   fmt.Sprintf("%.2f,%.2f,%.2f", req.Position.X, req.Position.Y, req.Position.Z),
	})
}