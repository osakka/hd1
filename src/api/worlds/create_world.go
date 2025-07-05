// Package worlds provides HTTP handlers for world management in HD1.
// Worlds represent YAML-based scene configurations that sessions can join
// for collaborative 3D environments with shared physics and real-time sync.
//
// Key concepts:
//   - World: YAML configuration defining 3D scene properties
//   - Environment: Base scene template (lighting, physics, entities)
//   - Session-World relationship: Sessions join worlds for scene state
//   - File-based storage: Worlds stored as YAML files in filesystem
package worlds

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

// CreateWorldRequest represents the request body for world creation.
// All fields except Name and Environment are optional with sensible defaults.
type CreateWorldRequest struct {
	ID          string `json:"id"`          // Optional: world identifier (auto-generated from name if empty)
	Name        string `json:"name"`        // Required: human-readable world name
	Description string `json:"description"` // Optional: world description
	Environment string `json:"environment"` // Required: base environment template
	MaxClients  int    `json:"max_clients"` // Optional: concurrent user limit (default: 50)
	Enabled     *bool  `json:"enabled"`     // Optional: world active state (default: true)
	Priority    int    `json:"priority"`    // Optional: world priority for ordering (default: 100)
}

// CreateWorldHandler handles POST /worlds requests.
// Creates a new world with YAML configuration file and directory structure.
// Worlds define scene environments that sessions can join for collaborative
// 3D experiences with shared physics, entities, and real-time synchronization.
//
// The request body should contain world name and environment template.
// Returns 201 Created with world details on success, or appropriate
// error status codes for validation failures.
//
// File system: Creates /opt/hd1/share/worlds/{world_id}/ directory
// with config.yaml containing world configuration.
//
// URL path: /api/worlds
// Method: POST
// Content-Type: application/json
func CreateWorldHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	logging.Debug("world creation requested", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	// Request parsing: Decode JSON request body to struct
	// Validates JSON format and structure
	var req CreateWorldRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.Error("failed to parse world creation request", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, `{"success": false, "message": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Field validation: Check required fields are present
	// Name and Environment are required for world creation
	if req.Name == "" {
		logging.Warn("world creation missing required field", map[string]interface{}{
			"field": "name",
		})
		http.Error(w, `{"success": false, "message": "World name is required"}`, http.StatusBadRequest)
		return
	}

	if req.Environment == "" {
		logging.Warn("world creation missing required field", map[string]interface{}{
			"field": "environment",
		})
		http.Error(w, `{"success": false, "message": "World environment is required"}`, http.StatusBadRequest)
		return
	}

	// ID generation: Auto-generate world ID from name if not provided
	// Ensures consistent naming convention: world_{sanitized_name}
	if req.ID == "" {
		// Generate ID from name: convert to lowercase, replace spaces with underscores, remove special chars
		id := strings.ToLower(req.Name)
		id = regexp.MustCompile(`[^a-z0-9_\-]`).ReplaceAllString(id, "")
		id = regexp.MustCompile(`\s+`).ReplaceAllString(id, "_")
		id = strings.Trim(id, "_-")
		
		// Ensure it starts with "world_"
		if !strings.HasPrefix(id, "world_") {
			id = "world_" + id
		}
		
		req.ID = id
		logging.Debug("auto-generated world ID", map[string]interface{}{
			"id": req.ID,
		})
	}

	// ID format validation: Ensure world ID follows naming convention
	// Required format: world_{alphanumeric_with_underscores_hyphens}
	if !regexp.MustCompile(`^world_[a-z0-9_\-]+$`).MatchString(req.ID) {
		logging.Warn("invalid world ID format", map[string]interface{}{
			"id": req.ID,
		})
		http.Error(w, `{"success": false, "message": "World ID must start with 'world_' and contain only lowercase letters, numbers, underscores, and hyphens"}`, http.StatusBadRequest)
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

	// Conflict detection: Verify world ID is unique
	// Prevents overwriting existing world configurations
	worldsDir := "/opt/hd1/share/worlds"
	worldPath := filepath.Join(worldsDir, req.ID)
	
	if _, err := os.Stat(worldPath); err == nil {
		logging.Warn("world ID already exists", map[string]interface{}{
			"id": req.ID,
		})
		http.Error(w, fmt.Sprintf(`{"success": false, "message": "World ID '%s' already exists"}`, req.ID), http.StatusConflict)
		return
	}

	// Directory creation: Create world directory structure
	// Each world gets isolated directory for configuration and assets
	if err := os.MkdirAll(worldPath, 0755); err != nil {
		logging.Error("failed to create world directory", map[string]interface{}{
			"error": err.Error(),
			"path":  worldPath,
		})
		http.Error(w, `{"success": false, "message": "Failed to create world directory"}`, http.StatusInternalServerError)
		return
	}

	// Configuration generation: Create YAML config file with world settings
	// Includes world metadata, collaboration settings, and performance tuning
	configContent := fmt.Sprintf(`# HD1 World Configuration - %s
# Auto-generated world configuration

name: "%s"
description: "%s"
environment: "%s"
max_clients: %d
enabled: %t
priority: %d

# World-specific settings
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

	configPath := filepath.Join(worldPath, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		logging.Error("failed to write world config", map[string]interface{}{
			"error": err.Error(),
			"path":  configPath,
		})
		
		// Cleanup on failure: Remove partially created directory
		// Maintains consistency - either complete success or clean failure
		os.RemoveAll(worldPath)
		http.Error(w, `{"success": false, "message": "Failed to create world configuration"}`, http.StatusInternalServerError)
		return
	}


	logging.Info("world created", map[string]interface{}{
		"world_id":     req.ID,
		"name":         req.Name,
		"environment":  req.Environment,
		"max_clients":  req.MaxClients,
		"enabled":      *req.Enabled,
	})

	// Success response: Return created world details
	// 201 Created with complete world configuration for client reference
	response := map[string]interface{}{
		"success": true,
		"message": "World created successfully",
		"world": map[string]interface{}{
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