package scenes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"holodeck/server"
)

// SaveSceneFromSessionHandler - POST /sessions/{sessionId}/scenes/save - Save session state as scene
func SaveSceneFromSessionHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	h, ok := hub.(*server.Hub)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Extract session ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 || pathParts[1] != "sessions" {
		http.Error(w, "Invalid session ID in path", http.StatusBadRequest)
		return
	}
	sessionID := pathParts[2]

	// Parse request body
	var req struct {
		SceneID     string `json:"scene_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Overwrite   bool   `json:"overwrite"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	// Validate session exists
	store := h.GetStore()
	_, exists := store.GetSession(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Validate scene ID format (alphanumeric, hyphens, underscores)
	if !isValidSceneID(req.SceneID) {
		http.Error(w, "Invalid scene_id format. Use alphanumeric, hyphens, and underscores only", http.StatusBadRequest)
		return
	}

	// Check if scene already exists
	scenePath := filepath.Join("/opt/holo-deck/share/glibsh/scenes", req.SceneID+".sh")
	if _, err := os.Stat(scenePath); err == nil && !req.Overwrite {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": fmt.Sprintf("Scene '%s' already exists. Use overwrite=true to replace", req.SceneID),
			"scene_path": scenePath,
		})
		return
	}

	// Get current session objects (PHOTO - current state snapshot)
	objects := store.ListObjects(sessionID)
	if len(objects) == 0 {
		http.Error(w, "No objects in session to save", http.StatusBadRequest)
		return
	}

	// Generate scene script from current session state
	scriptContent := generateSceneScript(req.SceneID, req.Name, req.Description, convertObjectPointers(objects))
	
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(scenePath), 0755); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create scenes directory: %v", err), http.StatusInternalServerError)
		return
	}

	// Write scene script file
	if err := os.WriteFile(scenePath, []byte(scriptContent), 0755); err != nil {
		http.Error(w, fmt.Sprintf("Failed to write scene script: %v", err), http.StatusInternalServerError)
		return
	}

	// Broadcast scene save notification
	saveMessage := map[string]interface{}{
		"type": "scene_saved",
		"scene_id": req.SceneID,
		"session_id": sessionID,
		"objects_exported": len(objects),
		"timestamp": time.Now().Format(time.RFC3339),
	}
	if _, err := json.Marshal(saveMessage); err == nil {
		h.BroadcastToSession(sessionID, "scene_saved", saveMessage)
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Scene '%s' saved successfully", req.Name),
		"scene_id": req.SceneID,
		"scene_path": scenePath,
		"objects_exported": len(objects),
	})
}

// generateSceneScript creates a bash script from session objects
func generateSceneScript(sceneID, name, description string, objects []server.Object) string {
	script := fmt.Sprintf(`#!/bin/bash

# =========================================================================
# THD Scene: %s - %s
# =========================================================================
#
# %s
#
# Usage: ./%s.sh [SESSION_ID]
# Auto-generated from session state on %s
# =========================================================================

set -euo pipefail

# Scene configuration
SCENE_NAME="%s"
SCENE_DESCRIPTION="%s"

# Get session ID from argument or use active session
SESSION_ID="${1:-${THD_SESSION:-}}"

if [[ -z "$SESSION_ID" ]]; then
    echo "Error: Session ID required" >&2
    exit 1
fi

# Path to auto-generated THD client
THD_CLIENT="/opt/holo-deck/build/bin/thd-client"

echo "Creating %s scene..."

`, name, description, description, sceneID, time.Now().Format("2006-01-02 15:04:05"), name, description, strings.ToLower(name))

	// Add object creation commands
	for _, obj := range objects {
		// Convert object to JSON for thd-client
		objectJSON := fmt.Sprintf(`{
    "name": "%s",
    "type": "%s",
    "x": %v, "y": %v, "z": %v,
    "scale": %v`, obj.Name, obj.Type, obj.X, obj.Y, obj.Z, obj.Scale)

		// Add color if present (Color is stored as string in current Object struct)
		if obj.Color != "" {
			objectJSON += fmt.Sprintf(`,
    "color": "%s"`, obj.Color)
		}

		// Note: Wireframe property not in current Object struct
		// TODO: Add wireframe support to Object struct

		objectJSON += "\n}"

		script += fmt.Sprintf(`
# %s - %s
$THD_CLIENT create-object "$SESSION_ID" '%s' > /dev/null
`, strings.Title(strings.ReplaceAll(obj.Name, "_", " ")), obj.Type, objectJSON)
	}

	// Add footer
	script += fmt.Sprintf(`
echo "THD Scene '$SCENE_NAME' loaded successfully"
echo "Objects created: %d"
echo "Session: $SESSION_ID"
`, len(objects))

	return script
}

// convertObjectPointers converts []*Object to []Object for generateSceneScript
func convertObjectPointers(objects []*server.Object) []server.Object {
	result := make([]server.Object, len(objects))
	for i, obj := range objects {
		result[i] = *obj
	}
	return result
}

// isValidSceneID validates scene ID format
func isValidSceneID(sceneID string) bool {
	if len(sceneID) == 0 || len(sceneID) > 50 {
		return false
	}
	
	for _, r := range sceneID {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || 
			 (r >= '0' && r <= '9') || r == '-' || r == '_') {
			return false
		}
	}
	return true
}