package sessions

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

// JoinSessionWorld - POST /sessions/{sessionId}/world/join
func JoinSessionWorldHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
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
	if len(pathParts) < 4 {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}
	sessionID := pathParts[2] // /api/sessions/{sessionId}/world/join
	
	// Parse request body
	var request struct {
		ClientID  string `json:"client_id"`
		WorldID   string `json:"world_id"`
		Reconnect bool   `json:"reconnect"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	if request.ClientID == "" {
		http.Error(w, "client_id is required", http.StatusBadRequest)
		return
	}
	
	// Verify session exists
	_, exists := h.GetStore().GetSession(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	
	// If world_id provided, load named world configuration
	if request.WorldID != "" {
		logging.Info("loading named world configuration", map[string]interface{}{
			"session_id": sessionID,
			"world_id": request.WorldID,
		})
		
		// CRITICAL FIX: Remove client from old world before joining new one
		if oldSession, exists := h.GetStore().GetSession(sessionID); exists && oldSession.WorldID != "" && oldSession.WorldID != request.WorldID {
			if removed, _ := h.LeaveSessionWorld(sessionID, request.ClientID); removed {
				logging.Info("client removed from old world", map[string]interface{}{
					"session_id": sessionID,
					"old_world": oldSession.WorldID,
					"new_world": request.WorldID,
					"client_id": request.ClientID,
				})
			}
		}
		
		// PHASE 1 FIX: Clear existing entities AND avatars (API-first with proper broadcasts)
		if err := h.ClearSessionEntitiesWithBroadcast(sessionID); err != nil {
			logging.Error("failed to clear session entities", map[string]interface{}{
				"session_id": sessionID,
				"error": err.Error(),
			})
			// Continue anyway - don't fail world join for clearing issues
		}
		
		// CRITICAL FIX: Clear avatar world association (don't remove avatar entirely)
		if err := h.GetSyncProtocol().ClearAvatarWorld(sessionID); err != nil {
			logging.Warn("failed to clear avatar world association", map[string]interface{}{
				"session_id": sessionID,
				"error": err.Error(),
			})
			// Continue anyway - avatar clearing failure shouldn't block world join
		}
		
		// Load world configuration and create entities
		if err := h.LoadNamedWorldIntoSession(sessionID, request.WorldID); err != nil {
			logging.Error("failed to load named world", map[string]interface{}{
				"session_id": sessionID,
				"world_id": request.WorldID,
				"error": err.Error(),
			})
			http.Error(w, "Failed to load world configuration", http.StatusInternalServerError)
			return
		}
		
		// Update session world association
		if err := h.GetStore().UpdateSessionWorld(sessionID, request.WorldID); err != nil {
			logging.Error("failed to update session world", map[string]interface{}{
				"session_id": sessionID,
				"world_id": request.WorldID,
				"error": err.Error(),
			})
		}
	}
	
	// Join the session world (create world if it doesn't exist)
	world, clientCount, graphState := h.JoinSessionWorld(sessionID, request.ClientID, request.Reconnect)
	
	// Create session avatar entity via API (100% API-first approach)
	avatarName := fmt.Sprintf("session_%s", request.ClientID)
	
	// Determine avatar type based on world defaults and client type
	avatarType := "humanoid_avatar" // System default
	
	// Use world-specific defaults
	switch request.WorldID {
	case "world_one":
		avatarType = "humanoid_avatar"
	case "world_two":  
		avatarType = "fox_avatar"
	case "world_three":
		avatarType = "humanoid_avatar"
	default:
		avatarType = "humanoid_avatar" // Fallback for unknown worlds
	}
	
	// Override for Claude clients
	if strings.Contains(request.ClientID, "claude") || strings.Contains(sessionID, "claude") {
		avatarType = "claude_avatar"
	}
	
	avatarPayload := map[string]interface{}{
		"name": avatarName,
		"tags": []string{"session-avatar", "user-representation", fmt.Sprintf("avatar-%s", avatarType)},
		"components": map[string]interface{}{
			"transform": map[string]interface{}{
				"position": []float64{0, 1.5, 0}, // Default spawn position
			},
			"model": map[string]interface{}{
				"type": "asset",
				"asset_path": "model.glb", // GLB model will be loaded via HTTP
			},
			"text": map[string]interface{}{
				"text": request.ClientID,
				"offset": []float64{0, 1, 0}, // Float above avatar
				"color": "#FFFFFF",
				"size": 0.3,
			},
		},
	}
	
	// Use HD1's API to create the avatar entity
	if err := h.CreateEntityViaAPI(sessionID, avatarPayload); err != nil {
		logging.Warn("failed to create session avatar via API", map[string]interface{}{
			"session_id": sessionID,
			"client_id": request.ClientID,
			"error": err.Error(),
		})
		// Continue - avatar creation failure shouldn't block world join
	} else {
		logging.Info("session avatar created via API", map[string]interface{}{
			"session_id": sessionID,
			"client_id": request.ClientID,
			"avatar_name": avatarName,
		})
		
		// CRITICAL FIX: Register avatar with HD1-VSC sync protocol
		avatarState := &hd1sync.AvatarState{
			SessionID:   sessionID,
			WorldID:     request.WorldID,
			Position:    hd1sync.Vector3{X: 0, Y: 1.5, Z: 0}, // Default spawn position
			Rotation:    hd1sync.Vector3{X: 0, Y: 0, Z: 0},   // Default rotation
			Animation:   "idle",
			Metadata: map[string]interface{}{
				"entity_id":   "", // Will be populated when entity is created
				"entity_name": avatarName,
				"avatar_type": avatarType, // Use correct avatar type for GLB loading
				"client_id":   request.ClientID,
				"instance_id": "0",
			},
			LastUpdate:  time.Now(),
			VectorClock: make(hd1sync.VectorClock),
		}
		
		// Register with sync protocol for perfect consistency
		h.GetSyncProtocol().RegisterAvatar(sessionID, avatarState)
		
		logging.Info("avatar registered with HD1-VSC sync protocol", map[string]interface{}{
			"session_id": sessionID,
			"client_id": request.ClientID,
			"avatar_name": avatarName,
			"world_id": request.WorldID,
			"sync_protocol": "HD1-VSC",
		})
	}
	
	logging.Info("client joined session world", map[string]interface{}{
		"session_id": sessionID,
		"client_id":  request.ClientID,
		"reconnect":  request.Reconnect,
		"client_count": clientCount,
	})
	
	// Broadcast world join event to other clients in the world
	h.BroadcastToSessionWorld(sessionID, "client_joined", map[string]interface{}{
		"client_id":    request.ClientID,
		"client_count": clientCount,
		"joined_at":    time.Now(),
	}, request.ClientID) // Exclude the joining client
	
	// Get session to return the correct world_id
	session, _ := h.GetStore().GetSession(sessionID)
	responseWorldID := world.GetID() // Default to session-based world ID
	if session.WorldID != "" {
		responseWorldID = session.WorldID // Use named world ID if available
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"world_id":   responseWorldID,
		"session_id":   sessionID,
		"client_count": clientCount,
		"graph_state":  graphState,
		"message":      "Joined session world",
	})
}

