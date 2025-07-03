// Package channels provides HTTP handlers for channel management in HD1.
// Channels represent YAML-based scene configurations that sessions can join
// for collaborative 3D environments with shared physics and real-time sync.
//
// Key concepts:
//   - Channel: YAML configuration defining 3D scene properties
//   - Environment: Base scene template (lighting, physics, entities)
//   - Session-Channel relationship: Sessions join channels for scene state
//   - File-based storage: Channels stored as YAML files in filesystem
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

// CreateChannelRequest represents the request body for channel creation.
// All fields except Name and Environment are optional with sensible defaults.
type CreateChannelRequest struct {
	ID          string `json:"id"`          // Optional: channel identifier (auto-generated from name if empty)
	Name        string `json:"name"`        // Required: human-readable channel name
	Description string `json:"description"` // Optional: channel description
	Environment string `json:"environment"` // Required: base environment template
	MaxClients  int    `json:"max_clients"` // Optional: concurrent user limit (default: 50)
	Enabled     *bool  `json:"enabled"`     // Optional: channel active state (default: true)
	Priority    int    `json:"priority"`    // Optional: channel priority for ordering (default: 100)
}

// CreateChannelHandler handles POST /channels requests.
// Creates a new channel with YAML configuration file and directory structure.
// Channels define scene environments that sessions can join for collaborative
// 3D experiences with shared physics, entities, and real-time synchronization.
//
// The request body should contain channel name and environment template.
// Returns 201 Created with channel details on success, or appropriate
// error status codes for validation failures.
//
// File system: Creates /opt/hd1/share/channels/{channel_id}/ directory
// with config.yaml containing channel configuration.
//
// URL path: /api/channels
// Method: POST
// Content-Type: application/json
func CreateChannelHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	logging.Debug("channel creation requested", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	// Request parsing: Decode JSON request body to struct
	// Validates JSON format and structure
	var req CreateChannelRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.Error("failed to parse channel creation request", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, `{"success": false, "message": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Field validation: Check required fields are present
	// Name and Environment are required for channel creation
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

	// ID generation: Auto-generate channel ID from name if not provided
	// Ensures consistent naming convention: channel_{sanitized_name}
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

	// ID format validation: Ensure channel ID follows naming convention
	// Required format: channel_{alphanumeric_with_underscores_hyphens}
	if !regexp.MustCompile(`^channel_[a-z0-9_\-]+$`).MatchString(req.ID) {
		logging.Warn("invalid channel ID format", map[string]interface{}{
			"id": req.ID,
		})
		http.Error(w, `{"success": false, "message": "Channel ID must start with 'channel_' and contain only lowercase letters, numbers, underscores, and hyphens"}`, http.StatusBadRequest)
		return
	}

	// Default values: Apply sensible defaults for optional fields
	// MaxClients=50, Enabled=true, Priority=100
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

	// Conflict detection: Verify channel ID is unique
	// Prevents overwriting existing channel configurations
	channelsDir := "/opt/hd1/share/channels"
	channelPath := filepath.Join(channelsDir, req.ID)
	
	if _, err := os.Stat(channelPath); err == nil {
		logging.Warn("channel ID already exists", map[string]interface{}{
			"id": req.ID,
		})
		http.Error(w, fmt.Sprintf(`{"success": false, "message": "Channel ID '%s' already exists"}`, req.ID), http.StatusConflict)
		return
	}

	// Directory creation: Create channel directory structure
	// Each channel gets isolated directory for configuration and assets
	if err := os.MkdirAll(channelPath, 0755); err != nil {
		logging.Error("failed to create channel directory", map[string]interface{}{
			"error": err.Error(),
			"path":  channelPath,
		})
		http.Error(w, `{"success": false, "message": "Failed to create channel directory"}`, http.StatusInternalServerError)
		return
	}

	// Configuration generation: Create YAML config file with channel settings
	// Includes channel metadata, collaboration settings, and performance tuning
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
		
		// Cleanup on failure: Remove partially created directory
		// Maintains consistency - either complete success or clean failure
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

	// Success response: Return created channel details
	// 201 Created with complete channel configuration for client reference
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