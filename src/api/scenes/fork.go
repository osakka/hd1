package scenes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"holodeck1/server"
)

// ForkSceneHandler - POST /scenes/{sceneId}/fork - Fork scene to session for editing
func ForkSceneHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	h, ok := hub.(*server.Hub)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Extract scene ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 || pathParts[1] != "scenes" {
		http.Error(w, "Invalid scene ID in path", http.StatusBadRequest)
		return
	}
	sceneID := pathParts[2]

	// Parse request body
	var req struct {
		SessionID     string `json:"session_id"`
		ClearExisting bool   `json:"clear_existing"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	// Validate session exists
	store := h.GetStore()
	_, exists := store.GetSession(req.SessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Clear existing objects if requested
	if req.ClearExisting {
		// Scene forking now works with channel entities, not session objects
		// Clear handled via PlayCanvas/channels
		
		// Broadcast clear message
		clearMessage := map[string]interface{}{
			"type": "objects_cleared",
			"session_id": req.SessionID,
			"timestamp": time.Now().Format(time.RFC3339),
		}
		if _, err := json.Marshal(clearMessage); err == nil {
			h.BroadcastToSession(req.SessionID, "objects_cleared", clearMessage)
		}
	}

	// Execute scene script to load objects into session
	scriptPath := filepath.Join("/opt/hd1/share/scenes", sceneID+".sh")
	objectCount, message := executeSceneScript(scriptPath, req.SessionID)
	
	if objectCount == 0 && strings.Contains(message, "failed") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": message,
		})
		return
	}

	// Mark all entities in session as "base" from forked scene
	// Scene forking now works with channel entities, not session objects
	// Entity tracking handled via PlayCanvas/channels

	// Broadcast fork completion
	forkMessage := map[string]interface{}{
		"type": "scene_forked",
		"scene_id": sceneID,
		"session_id": req.SessionID,
		"objects_forked": objectCount,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	if _, err := json.Marshal(forkMessage); err == nil {
		h.BroadcastToSession(req.SessionID, "scene_forked", forkMessage)
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Scene '%s' forked to session for editing", sceneID),
		"scene_id": sceneID,
		"session_id": req.SessionID,
		"objects_forked": objectCount,
		"fork_timestamp": time.Now().Format(time.RFC3339),
	})
}