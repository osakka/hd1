package sessions

import (
	"encoding/json"
	"net/http"
	"holodeck/server"
)

// CreateSessionHandler - POST /sessions
func CreateSessionHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type
	h, ok := hub.(*server.Hub)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Create session using SessionStore
	session := h.GetStore().CreateSession()
	
	// Automatically initialize world with default THD settings (25x25x25 grid, [-12,+12] bounds)
	world, err := h.GetStore().InitializeWorld(session.ID, 25, 0.1, 10, 10, 10)
	if err != nil {
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
		"message":          "Session created with world ready - THD holo-deck activated",
	})
}

