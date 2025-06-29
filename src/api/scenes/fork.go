package scenes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"holodeck/server"
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
		objects := store.ListObjects(req.SessionID)
		for _, obj := range objects {
			store.DeleteObject(req.SessionID, obj.Name)
		}
		
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
	scriptPath := filepath.Join("/opt/holo-deck/share/glibsh/scenes", sceneID+".sh")
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

	// Mark all objects in session as "base" from forked scene
	objects := store.ListObjects(req.SessionID)
	for _, obj := range objects {
		updates := map[string]interface{}{
			"tracking_status": "base",
			"source_scene": sceneID,
			"created_at": time.Now(),
		}
		store.UpdateObject(req.SessionID, obj.Name, updates)
	}

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