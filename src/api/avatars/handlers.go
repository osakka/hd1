package avatars

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"holodeck1/api/shared"
	"holodeck1/logging"
	"holodeck1/sync"
)

// MoveAvatarRequest represents the request to move an avatar
type MoveAvatarRequest struct {
	Position  shared.Vector3 `json:"position"`
	Rotation  *shared.Vector3 `json:"rotation,omitempty"`
	Animation string  `json:"animation,omitempty"`
}

// MoveAvatarResponse represents the response after moving an avatar
type MoveAvatarResponse struct {
	Success bool   `json:"success"`
	SeqNum  uint64 `json:"seq_num"`
}

// GetAvatars handles GET /api/threejs/avatars
func GetAvatars(w http.ResponseWriter, r *http.Request) {
	// Get hub from context
	hub := shared.GetHubFromContext(r)
	if hub == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get all connected avatars from registry
	avatars := hub.GetAvatarRegistry().GetAllAvatars()

	// Return list of connected avatars
	response := map[string]interface{}{
		"success": true,
		"avatars": avatars,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	logging.Debug("avatars listed via API", map[string]interface{}{
		"client_id":     shared.GetClientID(r),
		"avatar_count":  len(avatars),
	})
}

// CreateAvatar handles POST /api/threejs/avatars
func CreateAvatar(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name         string    `json:"name"`
		Position     shared.Vector3   `json:"position"`
		Capabilities []string  `json:"capabilities"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get client ID
	clientID := shared.GetClientID(r)

	// Generate avatar ID
	avatarID := fmt.Sprintf("avatar-%s", clientID)

	// Create operation data
	operationData := map[string]interface{}{
		"hd1_id":       avatarID,
		"name":         req.Name,
		"position":     req.Position,
		"capabilities": req.Capabilities,
	}

	// Create operation
	operation := &sync.Operation{
		ClientID:  clientID,
		Type:      "avatar_create",
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
	response := map[string]interface{}{
		"success":   true,
		"avatar_id": avatarID,
		"seq_num":   operation.SeqNum,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	logging.Debug("avatar created via API", map[string]interface{}{
		"avatar_id": avatarID,
		"client_id": clientID,
		"seq_num":   operation.SeqNum,
		"name":      req.Name,
	})
}

// UpdateAvatar handles PUT /api/threejs/avatars/{avatarId}
func UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	avatarID := vars["avatarId"]

	if avatarID == "" {
		http.Error(w, "Avatar ID required", http.StatusBadRequest)
		return
	}

	var req struct {
		Position  *shared.Vector3 `json:"position,omitempty"`
		Rotation  *shared.Vector3 `json:"rotation,omitempty"`
		Animation string   `json:"animation,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get hub from context
	hub := shared.GetHubFromContext(r)
	if hub == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get client ID
	clientID := shared.GetClientID(r)

	// Create operation data for sync
	operationData := map[string]interface{}{
		"hd1_id": avatarID,
	}

	// Prepare updates for registry
	updates := make(map[string]interface{})

	// Add optional properties
	if req.Position != nil {
		operationData["position"] = req.Position
		updates["position"] = map[string]interface{}{
			"x": req.Position.X,
			"y": req.Position.Y,
			"z": req.Position.Z,
		}
	}
	if req.Rotation != nil {
		operationData["rotation"] = req.Rotation
		updates["rotation"] = map[string]interface{}{
			"x": req.Rotation.X,
			"y": req.Rotation.Y,
			"z": req.Rotation.Z,
		}
	}
	if req.Animation != "" {
		operationData["animation"] = req.Animation
		updates["animation"] = req.Animation
	}

	// Update avatar in registry
	if err := hub.GetAvatarRegistry().UpdateAvatar(avatarID, updates); err != nil {
		http.Error(w, "Avatar not found", http.StatusNotFound)
		return
	}

	// Create operation
	operation := &sync.Operation{
		ClientID:  clientID,
		Type:      "avatar_update",
		Data:      operationData,
		Timestamp: time.Now(),
	}

	hub.GetSync().SubmitOperation(operation)

	// Return response
	response := map[string]interface{}{
		"success": true,
		"seq_num": operation.SeqNum,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	logging.Debug("avatar updated via API", map[string]interface{}{
		"avatar_id": avatarID,
		"client_id": clientID,
		"seq_num":   operation.SeqNum,
	})
}

// RemoveAvatar handles DELETE /api/threejs/avatars/{avatarId}
func RemoveAvatar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	avatarID := vars["avatarId"]

	if avatarID == "" {
		http.Error(w, "Avatar ID required", http.StatusBadRequest)
		return
	}

	// Get client ID
	clientID := shared.GetClientID(r)

	// Create operation data
	operationData := map[string]interface{}{
		"hd1_id": avatarID,
	}

	// Create operation
	operation := &sync.Operation{
		ClientID:  clientID,
		Type:      "avatar_remove",
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
	response := map[string]interface{}{
		"success": true,
		"seq_num": operation.SeqNum,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	logging.Debug("avatar removed via API", map[string]interface{}{
		"avatar_id": avatarID,
		"client_id": clientID,
		"seq_num":   operation.SeqNum,
	})
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
	clientID := shared.GetClientID(r)

	// Create operation data
	operationData := map[string]interface{}{
		"hd1_id":   sessionID,  // sessionID is actually the hd1_id
		"position": req.Position,
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
	hub := shared.GetHubFromContext(r)
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