package browser

import (
	"encoding/json"
	"net/http"
	"holodeck/server"
)

// ForceRefreshHandler - POST /browser/refresh - Force browser refresh and session reset
func ForceRefreshHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	h, ok := hub.(*server.Hub)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Parse request body for session ID
	var req struct {
		SessionID string `json:"session_id"`
		ClearStorage bool `json:"clear_storage,omitempty"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	// Validate session exists
	if req.SessionID != "" {
		if _, exists := h.GetStore().GetSession(req.SessionID); !exists {
			http.Error(w, "Session not found", http.StatusNotFound)
			return
		}
	}

	// Broadcast force refresh message
	refreshMessage := map[string]interface{}{
		"type": "force_refresh",
		"session_id": req.SessionID,
		"clear_storage": req.ClearStorage,
		"timestamp": "now",
	}

	if jsonData, err := json.Marshal(refreshMessage); err == nil {
		if req.SessionID != "" {
			// Session-specific refresh - only refresh browsers in this session
			h.BroadcastToSession(req.SessionID, "force_refresh", refreshMessage)
		} else {
			// Global refresh - refresh all connected browsers
			h.BroadcastMessage(jsonData)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Browser refresh command sent",
		"session_id": req.SessionID,
	})
}

// SetCanvasHandler - POST /browser/canvas - Direct WebGL control
func SetCanvasHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	h, ok := hub.(*server.Hub)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Parse request body for canvas commands
	var req struct {
		SessionID string `json:"session_id"`
		Command   string `json:"command"`
		Objects   []map[string]interface{} `json:"objects,omitempty"`
		Camera    map[string]interface{} `json:"camera,omitempty"`
		Clear     bool `json:"clear,omitempty"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	// Validate session exists
	if req.SessionID != "" {
		if _, exists := h.GetStore().GetSession(req.SessionID); !exists {
			http.Error(w, "Session not found", http.StatusNotFound)
			return
		}
	}

	// Create canvas control message
	canvasMessage := map[string]interface{}{
		"type": "canvas_control",
		"session_id": req.SessionID,
		"command": req.Command,
	}

	// Add optional parameters
	if req.Objects != nil {
		canvasMessage["objects"] = req.Objects
	}
	if req.Camera != nil {
		canvasMessage["camera"] = req.Camera
	}
	if req.Clear {
		canvasMessage["clear"] = true
	}

	// Broadcast canvas control message
	if jsonData, err := json.Marshal(canvasMessage); err == nil {
		if req.SessionID != "" {
			// Session-specific canvas control
			h.BroadcastToSession(req.SessionID, "canvas_control", canvasMessage)
		} else {
			// Global canvas control - all sessions
			h.BroadcastMessage(jsonData)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Canvas control command sent",
		"command": req.Command,
		"session_id": req.SessionID,
	})
}