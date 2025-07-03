// Package sessions provides HTTP handlers for session lifecycle management
// in HD1. Sessions represent isolated 3D game worlds that can contain
// entities, physics simulation, and real-time collaborative features.
//
// Key concepts:
//   - Session: Isolated game world with unique ID and state
//   - Channel: Scene configuration that sessions can join
//   - Real-time sync: WebSocket updates for session state changes
//   - Session store: Thread-safe storage for session persistence
package sessions

import (
	"encoding/json"
	"net/http"

	"holodeck1/logging"
	"holodeck1/server"
)

// CreateSessionHandler handles POST /sessions requests.
// Creates a new isolated game session with unique ID and default state.
// Sessions provide the foundation for 3D world management and can later
// join channels for scene configuration and collaborative features.
//
// Returns 201 Created with session details on success.
// Broadcasts 'session_created' event via WebSocket for real-time monitoring.
//
// URL path: /api/sessions
// Method: POST
// Content-Type: application/json (no body required)
func CreateSessionHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Hub interface validation: Ensure we have a valid server hub
	// This cast should never fail in normal operation
	h, ok := hub.(*server.Hub)
	if !ok {
		logging.Error("hub interface cast failed", map[string]interface{}{
			"expected_type": "*server.Hub",
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Session creation: Use thread-safe session store to create new session
	// Generates unique ID and initializes default state
	session := h.GetStore().CreateSession()
	logging.Info("session created successfully", map[string]interface{}{
		"session_id": session.ID,
		"created_at": session.CreatedAt,
	})
	
	// WebSocket notification: Broadcast session creation to monitoring clients
	// Enables real-time session tracking and administrative monitoring
	h.BroadcastUpdate("session_created", map[string]interface{}{
		"session": session,
	})
	
	// Success response: Return session details
	// 201 Created with session metadata and next steps guidance
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"session_id": session.ID,
		"created_at": session.CreatedAt,
		"status":     session.Status,
		"message":    "Session created - ready to join channels",
	})
}

