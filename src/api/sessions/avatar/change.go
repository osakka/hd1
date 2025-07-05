package avatar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
	hd1sync "holodeck1/sync"
)

// ChangeSessionAvatarRequest represents the request to change an avatar
type ChangeSessionAvatarRequest struct {
	AvatarType       string `json:"avatar_type"`
	PreservePosition bool   `json:"preserve_position"`
	Broadcast        bool   `json:"broadcast"`
}

// ChangeSessionAvatarResponse represents the response after changing an avatar
type ChangeSessionAvatarResponse struct {
	Success        bool                   `json:"success"`
	OldAvatarType  string                 `json:"old_avatar_type"`
	NewAvatarType  string                 `json:"new_avatar_type"`
	EntityID       string                 `json:"entity_id"`
	Position       map[string]interface{} `json:"position"`
	Message        string                 `json:"message"`
}

// ChangeSessionAvatar - PUT /sessions/{sessionId}/avatar
func ChangeSessionAvatarHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type
	h, ok := hub.(*server.Hub)
	if !ok {
		logging.Error("hub interface cast failed", map[string]interface{}{
			"expected_type": "*server.Hub",
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Extract session ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}
	sessionID := pathParts[2] // /api/sessions/{sessionId}/avatar
	
	// Parse request body
	var request ChangeSessionAvatarRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// Set defaults
	if request.PreservePosition == false && r.Header.Get("Content-Type") == "application/json" {
		request.PreservePosition = true // Default to preserving position
	}
	if request.Broadcast == false && r.Header.Get("Content-Type") == "application/json" {
		request.Broadcast = true // Default to broadcasting
	}
	
	logging.Info("avatar change request", map[string]interface{}{
		"session_id":        sessionID,
		"new_avatar_type":   request.AvatarType,
		"preserve_position": request.PreservePosition,
		"broadcast":         request.Broadcast,
	})
	
	// Validate avatar type
	validTypes := []string{"humanoid_avatar", "fox_avatar", "claude_avatar", "default"}
	isValid := false
	for _, validType := range validTypes {
		if request.AvatarType == validType {
			isValid = true
			break
		}
	}
	if !isValid {
		http.Error(w, fmt.Sprintf("Invalid avatar type. Valid types: %v", validTypes), http.StatusBadRequest)
		return
	}
	
	// Check if session exists
	session, exists := h.GetStore().GetSession(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	
	// Get current avatar information
	oldAvatarType := "unknown"
	var currentPosition map[string]interface{}
	
	// For now, use default position (we'll need to implement proper entity lookup later)
	currentPosition = map[string]interface{}{
		"x": 0.0,
		"y": 1.5,
		"z": 0.0,
	}
	
	// TODO: Implement proper entity lookup when entity store methods are available
	/*
	// Find current avatar entity
	entities := h.GetStore().GetSessionEntities(sessionID)
	var oldAvatarEntity map[string]interface{}
	for _, entity := range entities {
		if entityMap, ok := entity.(map[string]interface{}); ok {
			if tags, tagExists := entityMap["tags"].([]interface{}); tagExists {
				isSessionAvatar := false
				for _, tag := range tags {
					if tagStr, ok := tag.(string); ok && tagStr == "session-avatar" {
						isSessionAvatar = true
						break
					}
				}
				if isSessionAvatar {
					oldAvatarEntity = entityMap
					// Extract old avatar type from tags
					for _, tag := range tags {
						if tagStr, ok := tag.(string); ok && strings.HasPrefix(tagStr, "avatar-") {
							oldAvatarType = strings.TrimPrefix(tagStr, "avatar-")
							break
						}
					}
					// Get current position if preserving
					if request.PreservePosition {
						if components, compExists := entityMap["components"].(map[string]interface{}); compExists {
							if transform, transExists := components["transform"].(map[string]interface{}); transExists {
								if pos, posExists := transform["position"].([]interface{}); posExists && len(pos) >= 3 {
									currentPosition = map[string]interface{}{
										"x": pos[0],
										"y": pos[1], 
										"z": pos[2],
									}
								}
							}
						}
					}
					break
				}
			}
		}
	}
	*/
	
	// TODO: Remove old avatar entity when entity store methods are available
	// For now, clear from sync protocol
	if err := h.GetSyncProtocol().ClearAvatarWorld(sessionID); err != nil {
		logging.Warn("failed to clear avatar from sync protocol", map[string]interface{}{
			"session_id": sessionID,
			"error":      err.Error(),
		})
	}
	
	// Create new avatar entity with the specified type
	avatarName := fmt.Sprintf("session_client_%s", sessionID)
	avatarPayload := map[string]interface{}{
		"name": avatarName,
		"tags": []string{"session-avatar", "user-representation", fmt.Sprintf("avatar-%s", request.AvatarType)},
		"components": map[string]interface{}{
			"transform": map[string]interface{}{
				"position": []float64{
					currentPosition["x"].(float64),
					currentPosition["y"].(float64), 
					currentPosition["z"].(float64),
				},
			},
			"model": map[string]interface{}{
				"type":       "asset",
				"asset_path": "model.glb",
			},
			"text": map[string]interface{}{
				"text":   fmt.Sprintf("client_%s", sessionID),
				"offset": []float64{0, 1, 0},
				"color":  "#FFFFFF",
				"size":   0.3,
			},
		},
	}
	
	// Create the new avatar entity
	if err := h.CreateEntityViaAPI(sessionID, avatarPayload); err != nil {
		logging.Error("failed to create new avatar entity", map[string]interface{}{
			"session_id":      sessionID,
			"new_avatar_type": request.AvatarType,
			"error":           err.Error(),
		})
		http.Error(w, "Failed to create new avatar", http.StatusInternalServerError)
		return
	}
	
	// Register new avatar with sync protocol
	avatarState := &hd1sync.AvatarState{
		SessionID: sessionID,
		WorldID: session.WorldID,
		Position: hd1sync.Vector3{
			X: currentPosition["x"].(float64),
			Y: currentPosition["y"].(float64),
			Z: currentPosition["z"].(float64),
		},
		Rotation:  hd1sync.Vector3{X: 0, Y: 0, Z: 0},
		Animation: "idle",
		Metadata: map[string]interface{}{
			"entity_id":   "", // Will be populated when entity is created
			"entity_name": avatarName,
			"avatar_type": request.AvatarType,
			"client_id":   fmt.Sprintf("client_%s", sessionID),
			"instance_id": "0",
		},
		LastUpdate:  time.Now(),
		VectorClock: make(hd1sync.VectorClock),
	}
	
	h.GetSyncProtocol().RegisterAvatar(sessionID, avatarState)
	
	// Broadcast avatar change if requested
	if request.Broadcast {
		changeMessage := map[string]interface{}{
			"session_id":      sessionID,
			"old_avatar_type": oldAvatarType,
			"new_avatar_type": request.AvatarType,
			"position":        currentPosition,
			"timestamp":       time.Now().Unix(),
		}
		
		// Broadcast to all clients in the session
		h.BroadcastToSession(sessionID, "avatar_changed", changeMessage)
		
		logging.Info("avatar change broadcasted", map[string]interface{}{
			"session_id":      sessionID,
			"world_id":        session.WorldID,
			"old_avatar_type": oldAvatarType,
			"new_avatar_type": request.AvatarType,
		})
	}
	
	// Return success response
	response := ChangeSessionAvatarResponse{
		Success:       true,
		OldAvatarType: oldAvatarType,
		NewAvatarType: request.AvatarType,
		EntityID:      fmt.Sprintf("avatar-%s-%s", sessionID, request.AvatarType),
		Position:      currentPosition,
		Message:       fmt.Sprintf("Avatar changed from %s to %s", oldAvatarType, request.AvatarType),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	
	logging.Info("avatar changed successfully", map[string]interface{}{
		"session_id":      sessionID,
		"old_avatar_type": oldAvatarType,
		"new_avatar_type": request.AvatarType,
		"position":        currentPosition,
		"broadcast":       request.Broadcast,
	})
}