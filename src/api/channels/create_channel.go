package channels

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"holodeck1/logging"
)

// CreateChannelHandler handles POST /channels - create new channel
func CreateChannelHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	logging.Debug("channel creation requested", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	// Parse request body
	var req struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Environment string `json:"environment"`
		MaxClients  int    `json:"max_clients"`
		Enabled     *bool  `json:"enabled"`
		Priority    int    `json:"priority"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.Error("failed to parse channel creation request", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, `{"success": false, "message": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		logging.Warn("channel creation missing required field", map[string]interface{}{
			"field": "name",
		})
		http.Error(w, `{"success": false, "message": "Channel name is required"}`, http.StatusBadRequest)
		return
	}

	if req.Environment == "" {
		logging.Warn("channel creation missing required field", map[string]interface{}{
			"field": "environment",
		})
		http.Error(w, `{"success": false, "message": "Channel environment is required"}`, http.StatusBadRequest)
		return
	}

	// Auto-generate ID if not provided
	if req.ID == "" {
		// Generate ID from name: convert to lowercase, replace spaces with underscores, remove special chars
		id := strings.ToLower(req.Name)
		id = regexp.MustCompile(`[^a-z0-9_\-]`).ReplaceAllString(id, "")
		id = regexp.MustCompile(`\s+`).ReplaceAllString(id, "_")
		id = strings.Trim(id, "_-")
		
		// Ensure it starts with "channel_"
		if !strings.HasPrefix(id, "channel_") {
			id = "channel_" + id
		}
		
		req.ID = id
		logging.Debug("auto-generated channel ID", map[string]interface{}{
			"id": req.ID,
		})
	}

	// Validate ID format
	if !regexp.MustCompile(`^channel_[a-z0-9_\-]+$`).MatchString(req.ID) {
		logging.Warn("invalid channel ID format", map[string]interface{}{
			"id": req.ID,
		})
		http.Error(w, `{"success": false, "message": "Channel ID must start with 'channel_' and contain only lowercase letters, numbers, underscores, and hyphens"}`, http.StatusBadRequest)
		return
	}

	// Set defaults
	if req.MaxClients <= 0 {
		req.MaxClients = 50
	}
	if req.Enabled == nil {
		enabled := true
		req.Enabled = &enabled
	}
	if req.Priority <= 0 {
		req.Priority = 100
	}

	// Check if channel already exists
	channelsDir := "/opt/hd1/share/channels"
	channelPath := filepath.Join(channelsDir, req.ID)
	
	if _, err := os.Stat(channelPath); err == nil {
		logging.Warn("channel ID already exists", map[string]interface{}{
			"id": req.ID,
		})
		http.Error(w, fmt.Sprintf(`{"success": false, "message": "Channel ID '%s' already exists"}`, req.ID), http.StatusConflict)
		return
	}

	// Create channel directory
	if err := os.MkdirAll(channelPath, 0755); err != nil {
		logging.Error("failed to create channel directory", map[string]interface{}{
			"error": err.Error(),
			"path":  channelPath,
		})
		http.Error(w, `{"success": false, "message": "Failed to create channel directory"}`, http.StatusInternalServerError)
		return
	}

	// Create channel configuration YAML
	configContent := fmt.Sprintf(`# HD1 Channel Configuration - %s
# Auto-generated channel configuration

name: "%s"
description: "%s"
environment: "%s"
max_clients: %d
enabled: %t
priority: %d

# Channel-specific settings
auto_apply_environment: true
allow_scene_modification: true
enable_physics: true

# Collaboration settings
allow_voice_chat: true
allow_screen_sharing: false
max_session_duration_hours: 24

# Performance settings
entity_limit: 1000
physics_update_rate: 60
network_sync_rate: 30

# Created automatically via API
created_at: "%s"
created_via: "api"
`, req.Name, req.Name, req.Description, req.Environment, req.MaxClients, *req.Enabled, req.Priority, time.Now().Format(time.RFC3339))

	configPath := filepath.Join(channelPath, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		logging.Error("failed to write channel config", map[string]interface{}{
			"error": err.Error(),
			"path":  configPath,
		})
		
		// Clean up directory on failure
		os.RemoveAll(channelPath)
		http.Error(w, `{"success": false, "message": "Failed to create channel configuration"}`, http.StatusInternalServerError)
		return
	}


	logging.Info("channel created successfully", map[string]interface{}{
		"channel_id":   req.ID,
		"name":         req.Name,
		"environment":  req.Environment,
		"max_clients":  req.MaxClients,
		"enabled":      *req.Enabled,
	})

	// Return success response
	response := map[string]interface{}{
		"success": true,
		"message": "Channel created successfully",
		"channel": map[string]interface{}{
			"id":          req.ID,
			"name":        req.Name,
			"description": req.Description,
			"environment": req.Environment,
			"max_clients": req.MaxClients,
			"enabled":     *req.Enabled,
			"priority":    req.Priority,
			"created_at":  time.Now().Format(time.RFC3339),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}