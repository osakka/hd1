package sessions

import (
	"encoding/json"
	"net/http"

	"holodeck1/logging"
	"holodeck1/server"
)

// CreateSessionHandler - POST /sessions
func CreateSessionHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type
	h, ok := hub.(*server.Hub)
	if !ok {
		logging.Error("hub interface cast failed", map[string]interface{}{
			"expected_type": "*server.Hub",
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Create session using SessionStore
	session := h.GetStore().CreateSession()
	logging.Info("session created successfully", map[string]interface{}{
		"session_id": session.ID,
		"created_at": session.CreatedAt,
	})
	
	// Automatically initialize world with holodeck coordinate system (floor=Y:0, human eye level=Y:1.7)
	world, err := h.GetStore().InitializeWorld(session.ID, 25, 0.01, 0, 1.7, 0)
	if err != nil {
		logging.Error("world initialization failed", map[string]interface{}{
			"session_id": session.ID,
			"error": err.Error(),
		})
		http.Error(w, "Failed to initialize world: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Broadcast session creation with world initialization for real-time updates
	h.BroadcastUpdate("session_created", map[string]interface{}{
		"session": session,
		"world":   world,
	})
	
	// Broadcast world initialization separately for clients expecting this event
	h.BroadcastUpdate("world_initialized", map[string]interface{}{
		"session_id": session.ID,
		"world":      world,
	})
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":           true,
		"session_id":        session.ID,
		"created_at":        session.CreatedAt,
		"status":            session.Status,
		"world":            world,
		"bounds":           map[string]int{"min": -12, "max": 12},
		"coordinate_system": "fixed_grid",
		"message":          "Session created with world ready - HD1 holodeck activated",
	})
}

