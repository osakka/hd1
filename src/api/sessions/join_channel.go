package sessions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
)

// JoinSessionChannel - POST /sessions/{sessionId}/channel/join
func JoinSessionChannelHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
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
	sessionID := pathParts[2] // /api/sessions/{sessionId}/channel/join
	
	// Parse request body
	var request struct {
		ClientID  string `json:"client_id"`
		ChannelID string `json:"channel_id"`
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
	
	// If channel_id provided, load named channel configuration
	if request.ChannelID != "" {
		logging.Info("loading named channel configuration", map[string]interface{}{
			"session_id": sessionID,
			"channel_id": request.ChannelID,
		})
		
		// CRITICAL FIX: Remove client from old channel before joining new one
		if oldSession, exists := h.GetStore().GetSession(sessionID); exists && oldSession.ChannelID != "" && oldSession.ChannelID != request.ChannelID {
			if removed, _ := h.LeaveSessionChannel(sessionID, request.ClientID); removed {
				logging.Info("client removed from old channel", map[string]interface{}{
					"session_id": sessionID,
					"old_channel": oldSession.ChannelID,
					"new_channel": request.ChannelID,
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
			// Continue anyway - don't fail channel join for clearing issues
		}
		
		// CRITICAL FIX: Clear avatar from sync protocol when changing channels
		if err := h.ClearAvatarFromSyncProtocol(sessionID); err != nil {
			logging.Warn("failed to clear avatar from sync protocol", map[string]interface{}{
				"session_id": sessionID,
				"error": err.Error(),
			})
			// Continue anyway - avatar clearing failure shouldn't block channel join
		}
		
		// Load channel configuration and create entities
		if err := h.LoadNamedChannelIntoSession(sessionID, request.ChannelID); err != nil {
			logging.Error("failed to load named channel", map[string]interface{}{
				"session_id": sessionID,
				"channel_id": request.ChannelID,
				"error": err.Error(),
			})
			http.Error(w, "Failed to load channel configuration", http.StatusInternalServerError)
			return
		}
		
		// Update session channel association
		if err := h.GetStore().UpdateSessionChannel(sessionID, request.ChannelID); err != nil {
			logging.Error("failed to update session channel", map[string]interface{}{
				"session_id": sessionID,
				"channel_id": request.ChannelID,
				"error": err.Error(),
			})
		}
	}
	
	// Join the session channel (create channel if it doesn't exist)
	channel, clientCount, graphState := h.JoinSessionChannel(sessionID, request.ClientID, request.Reconnect)
	
	// Create session avatar entity via API (100% API-first approach)
	avatarName := fmt.Sprintf("session_%s", request.ClientID)
	avatarPayload := map[string]interface{}{
		"name": avatarName,
		"tags": []string{"session-avatar", "user-representation"},
		"components": map[string]interface{}{
			"transform": map[string]interface{}{
				"position": []float64{0, 1.5, 0}, // Default spawn position
			},
			"render": map[string]interface{}{
				"geometry": "capsule",
				"material": map[string]interface{}{
					"color": "#00FFFF", // Teal color for session avatar
				},
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
		// Continue - avatar creation failure shouldn't block channel join
	} else {
		logging.Info("session avatar created via API", map[string]interface{}{
			"session_id": sessionID,
			"client_id": request.ClientID,
			"avatar_name": avatarName,
		})
	}
	
	logging.Info("client joined session channel", map[string]interface{}{
		"session_id": sessionID,
		"client_id":  request.ClientID,
		"reconnect":  request.Reconnect,
		"client_count": clientCount,
	})
	
	// Broadcast channel join event to other clients in the channel
	h.BroadcastToSessionChannel(sessionID, "client_joined", map[string]interface{}{
		"client_id":    request.ClientID,
		"client_count": clientCount,
		"joined_at":    time.Now(),
	}, request.ClientID) // Exclude the joining client
	
	// Get session to return the correct channel_id
	session, _ := h.GetStore().GetSession(sessionID)
	responseChannelID := channel.GetID() // Default to session-based channel ID
	if session.ChannelID != "" {
		responseChannelID = session.ChannelID // Use named channel ID if available
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"channel_id":   responseChannelID,
		"session_id":   sessionID,
		"client_count": clientCount,
		"graph_state":  graphState,
		"message":      "Successfully joined session channel",
	})
}