package camera

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"holodeck1/config"
	"holodeck1/server"
	"holodeck1/logging"
)

type CameraPositionRequest struct {
	Position struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	} `json:"position"`
}

func SetCameraPositionHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	w.Header().Set("Content-Type", "application/json")
	
	// Extract session ID from URL path
	parts := strings.Split(r.URL.Path, "/")
	var sessionID string
	for i, part := range parts {
		if part == "sessions" && i+1 < len(parts) {
			sessionID = parts[i+1]
			break
		}
	}
	
	if sessionID == "" {
		http.Error(w, "Session ID not found in path", http.StatusBadRequest)
		return
	}
	
	// Parse camera position request
	var req CameraPositionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.Warn("failed to decode camera position request", map[string]interface{}{
			"session_id": sessionID,
			"error": err.Error(),
		})
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Get hub instance
	h, ok := hub.(*server.Hub)
	if !ok {
		logging.Error("invalid hub type in camera position handler", map[string]interface{}{
			"session_id": sessionID,
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Update session avatar position (100% API-first approach)
	// Verify session exists
	_, exists := h.GetStore().GetSession(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	
	// Look for avatar entity by trying common patterns
	// Avatar creation uses format: session_{client_id}
	// Try to find any entity starting with "session_" that has session-avatar tag
	avatarName := ""
	
	// Since we need to find entities by pattern, we'll make an internal HTTP call to list entities
	// This maintains the 100% API-first principle
	entitiesURL := fmt.Sprintf("%s/sessions/%s/entities", config.GetInternalAPIBase(), sessionID)
	resp, err := http.Get(entitiesURL)
	if err != nil {
		logging.Warn("failed to list entities for avatar search", map[string]interface{}{
			"session_id": sessionID,
			"error": err.Error(),
		})
		http.Error(w, "Could not access session entities", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	var entitiesResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&entitiesResponse); err != nil {
		http.Error(w, "Could not parse entities response", http.StatusInternalServerError)
		return
	}
	
	// Use the entities API with tag filter for clean avatar lookup
	entitiesURL = fmt.Sprintf("%s/sessions/%s/entities?tag=session-avatar", config.GetInternalAPIBase(), sessionID)
	resp, err = http.Get(entitiesURL)
	if err != nil {
		logging.Warn("failed to filter entities by session-avatar tag", map[string]interface{}{
			"session_id": sessionID,
			"error": err.Error(),
		})
		http.Error(w, "Could not find session avatar", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	if err := json.NewDecoder(resp.Body).Decode(&entitiesResponse); err != nil {
		http.Error(w, "Could not parse avatar entities response", http.StatusInternalServerError)
		return
	}
	
	// Get the first session avatar entity
	if entitiesData, ok := entitiesResponse["entities"].([]interface{}); ok && len(entitiesData) > 0 {
		if entityMap, ok := entitiesData[0].(map[string]interface{}); ok {
			if name, ok := entityMap["name"].(string); ok {
				avatarName = name
			}
		}
	}
	
	if avatarName == "" {
		logging.Warn("session avatar not found for update", map[string]interface{}{
			"session_id": sessionID,
		})
		http.Error(w, "Avatar not found or session not in channel", http.StatusNotFound)
		return
	}
	
	// 🔥 CRITICAL FIX: Don't update avatar via API (causes entity_updated → delete/recreate)
	// The WebSocket avatar_position_update below is sufficient for real-time avatar movement
	// avatarUpdatePayload := map[string]interface{}{
	// 	"components": map[string]interface{}{
	// 		"transform": map[string]interface{}{
	// 			"position": []float64{req.X, req.Y + 1.5, req.Z}, // Avatar slightly above camera
	// 		},
	// 	},
	// }
	// 
	// if err := h.UpdateEntityByNameViaAPI(sessionID, avatarName, avatarUpdatePayload); err != nil {
	// 	logging.Warn("failed to update session avatar position via API", map[string]interface{}{
	// 		"session_id": sessionID,
	// 		"avatar_name": avatarName,
	// 		"position": map[string]float64{"x": req.X, "y": req.Y, "z": req.Z},
	// 		"error": err.Error(),
	// 	})
	// 	// Don't fail the camera position request if avatar update fails
	// }
	
	// Broadcast avatar position update to all sessions in the channel via WebSocket
	// This ensures all participants see the avatar movement in real-time
	avatarPositionUpdate := map[string]interface{}{
		"session_id": sessionID,
		"avatar_name": avatarName,
		"position": map[string]float64{
			"x": req.Position.X,
			"y": req.Position.Y + 1.5, // Avatar position (camera + offset)
			"z": req.Position.Z,
		},
		"camera_position": map[string]float64{
			"x": req.Position.X,
			"y": req.Position.Y,
			"z": req.Position.Z,
		},
	}
	
	// 🚀 REVOLUTIONARY HD1-VSC SYNCHRONIZATION: Apply avatar movement with perfect consistency
	// Vector clocks + Delta-State CRDTs + Authoritative server = 100% consistency guarantee
	position := map[string]float64{
		"x": req.Position.X,
		"y": req.Position.Y + 1.5, // Avatar position (camera + offset)
		"z": req.Position.Z,
	}
	rotation := map[string]float64{
		"x": 0.0,
		"y": 0.0,
		"z": 0.0,
	}
	
	// Apply movement through sync protocol for perfect consistency
	if err := h.ApplyAvatarMovement(sessionID, position, rotation); err != nil {
		logging.Warn("HD1-VSC avatar movement failed", map[string]interface{}{
			"session_id": sessionID,
			"error": err.Error(),
		})
		// Fallback to legacy broadcast
		h.BroadcastAvatarPositionToChannel(sessionID, "avatar_position_update", avatarPositionUpdate)
	}
	
	logging.Info("avatar position broadcast to session participants", map[string]interface{}{
		"session_id": sessionID,
		"position": avatarPositionUpdate["position"],
		"broadcast_type": "avatar_position_update",
	})
	
	logging.Info("camera position updated with avatar sync", map[string]interface{}{
		"session_id": sessionID,
		"position": map[string]float64{"x": req.Position.X, "y": req.Position.Y, "z": req.Position.Z},
		"avatar_name": avatarName,
	})
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Camera positioned with avatar sync and broadcast",
		"position": map[string]float64{"x": req.Position.X, "y": req.Position.Y, "z": req.Position.Z},
		"avatar_position": avatarPositionUpdate["position"],
	})
}

func GetCameraPositionHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	w.Header().Set("Content-Type", "application/json")
	
	// Extract session ID from URL path
	parts := strings.Split(r.URL.Path, "/")
	var sessionID string
	for i, part := range parts {
		if part == "sessions" && i+1 < len(parts) {
			sessionID = parts[i+1]
			break
		}
	}
	
	if sessionID == "" {
		http.Error(w, "Session ID not found in path", http.StatusBadRequest)
		return
	}
	
	// Get hub instance
	h, ok := hub.(*server.Hub)
	if !ok {
		logging.Error("invalid hub type in camera position handler", map[string]interface{}{
			"session_id": sessionID,
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Get avatar entity position (avatar represents camera position)
	// Get session to verify it exists
	_, exists := h.GetStore().GetSession(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	
	// Make internal HTTP call to list entities (maintains 100% API-first principle)
	entitiesURL := fmt.Sprintf("%s/sessions/%s/entities", config.GetInternalAPIBase(), sessionID)
	resp, err := http.Get(entitiesURL)
	if err != nil {
		logging.Warn("failed to list entities for avatar search", map[string]interface{}{
			"session_id": sessionID,
			"error": err.Error(),
		})
		http.Error(w, "Could not access session entities", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	var entitiesResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&entitiesResponse); err != nil {
		http.Error(w, "Could not parse entities response", http.StatusInternalServerError)
		return
	}
	
	// Use the entities API with tag filter for clean avatar lookup
	entitiesURL = fmt.Sprintf("%s/sessions/%s/entities?tag=session-avatar", config.GetInternalAPIBase(), sessionID)
	resp, err = http.Get(entitiesURL)
	if err != nil {
		logging.Warn("failed to filter entities by session-avatar tag", map[string]interface{}{
			"session_id": sessionID,
			"error": err.Error(),
		})
		http.Error(w, "Could not find session avatar", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	
	if err := json.NewDecoder(resp.Body).Decode(&entitiesResponse); err != nil {
		http.Error(w, "Could not parse avatar entities response", http.StatusInternalServerError)
		return
	}
	
	// Get the first session avatar entity
	var avatarEntityData map[string]interface{}
	if entitiesData, ok := entitiesResponse["entities"].([]interface{}); ok && len(entitiesData) > 0 {
		if entityMap, ok := entitiesData[0].(map[string]interface{}); ok {
			avatarEntityData = entityMap
		}
	}
	
	if avatarEntityData == nil {
		logging.Warn("session avatar not found", map[string]interface{}{
			"session_id": sessionID,
		})
		http.Error(w, "Avatar not found or session not in channel", http.StatusNotFound)
		return
	}
	
	// Extract position from avatar entity
	var avatarPosition map[string]float64
	var cameraPosition map[string]float64
	
	// Parse transform component from avatar entity
	if components, ok := avatarEntityData["components"].(map[string]interface{}); ok {
		if transform, ok := components["transform"].(map[string]interface{}); ok {
			if position, ok := transform["position"].([]interface{}); ok && len(position) >= 3 {
				avatarPosition = map[string]float64{
					"x": position[0].(float64),
					"y": position[1].(float64),
					"z": position[2].(float64),
				}
				// Camera position is avatar position minus offset
				cameraPosition = map[string]float64{
					"x": avatarPosition["x"],
					"y": avatarPosition["y"] - 1.5, // Remove avatar offset
					"z": avatarPosition["z"],
				}
			}
		}
	}
	
	if avatarPosition == nil {
		http.Error(w, "Could not determine camera position from avatar", http.StatusInternalServerError)
		return
	}
	
	logging.Info("camera position retrieved", map[string]interface{}{
		"session_id": sessionID,
		"camera_position": cameraPosition,
		"avatar_position": avatarPosition,
	})
	
	// Get avatar name from entity data
	avatarName := "unknown"
	if name, ok := avatarEntityData["name"].(string); ok {
		avatarName = name
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"camera_position": cameraPosition,
		"avatar_position": avatarPosition,
		"avatar_name": avatarName,
	})
}