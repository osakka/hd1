package scenes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"holodeck1/server"
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
	scenePath := filepath.Join("/opt/hd1/share/scenes", req.SceneID+".sh")
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

	// HD1 v3.0: Get session's current channel to extract entities
	session, exists := store.GetSession(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	
	// Check if session is joined to a channel
	if session.ChannelID == "" {
		http.Error(w, "Session must be joined to a channel to save scene", http.StatusBadRequest)
		return
	}
	
	// In HD1 v3.0, entities are defined in channel YAML, not stored in session
	// This API creates a script that references the channel configuration

	// Generate scene script that references the channel configuration
	scriptContent := generateSceneScript(req.SceneID, req.Name, req.Description, session.ChannelID)
	
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
		"channel_id": session.ChannelID,
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
		"channel_id": session.ChannelID,
	})
}

// generateSceneScript creates a bash script that references channel configuration
func generateSceneScript(sceneID, name, description, channelID string) string {
	script := fmt.Sprintf(`#!/bin/bash

# =========================================================================
# HD1 Scene: %s - %s
# =========================================================================
#
# %s
#
# Usage: ./%s.sh [SESSION_ID]
# Auto-generated from channel %s configuration on %s
# =========================================================================

set -euo pipefail

# Scene configuration
SCENE_NAME="%s"
SCENE_DESCRIPTION="%s"
CHANNEL_ID="%s"

# Get session ID from argument or use active session
SESSION_ID="${1:-${HD1_SESSION:-}}"

if [[ -z "$SESSION_ID" ]]; then
    echo "Error: Session ID required" >&2
    exit 1
fi

# Path to auto-generated HD1 client
HD1_CLIENT="/opt/hd1/build/bin/hd1-client"

echo "Creating %s scene from channel %s..."

# Join session to channel to load entities
"$HD1_CLIENT" join-session-channel --session-id "$SESSION_ID" --channel-id "$CHANNEL_ID"

`, name, description, description, sceneID, channelID, time.Now().Format("2006-01-02 15:04:05"), name, description, channelID, strings.ToLower(name), channelID)

	// HD1 v3.0: Entities are defined in channel YAML configuration
	// The scene script simply joins the session to the channel
	
	// Add footer
	script += fmt.Sprintf(`
echo "HD1 Scene '$SCENE_NAME' loaded successfully"
echo "Channel: $CHANNEL_ID entities loaded"
echo "Session: $SESSION_ID"
`)

	return script
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