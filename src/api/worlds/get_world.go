package worlds

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"holodeck1/logging"
	"gopkg.in/yaml.v3"
)

// GetWorldHandler handles GET /worlds/{worldId} - get specific world details
func GetWorldHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Extract world ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/worlds/")
	worldId := strings.Split(path, "/")[0]
	
	logging.Debug("world details requested", map[string]interface{}{
		"world_id": worldId,
		"method":     r.Method,
		"path":       r.URL.Path,
	})

	if worldId == "" {
		logging.Warn("missing world ID in request", map[string]interface{}{
			"path": r.URL.Path,
		})
		http.Error(w, `{"success": false, "message": "World ID is required"}`, http.StatusBadRequest)
		return
	}

	// Validate world exists
	worldsDir := "/opt/hd1/share/worlds"
	configPath := filepath.Join(worldsDir, worldId, "config.yaml")
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logging.Warn("world not found", map[string]interface{}{
			"world_id": worldId,
			"path":       configPath,
		})
		http.Error(w, fmt.Sprintf(`{"success": false, "message": "World '%s' not found"}`, worldId), http.StatusNotFound)
		return
	}

	// Load world configuration
	configData, err := os.ReadFile(configPath)
	if err != nil {
		logging.Error("failed to read world config", map[string]interface{}{
			"error":    err.Error(),
			"world_id": worldId,
			"path":       configPath,
		})
		http.Error(w, `{"success": false, "message": "Failed to load world configuration"}`, http.StatusInternalServerError)
		return
	}

	// Parse YAML configuration using the same structs as list_worlds.go
	var fullConfig struct {
		ID          string `yaml:"id,omitempty"`
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		Environment string `yaml:"environment"`
		MaxClients  int    `yaml:"max_clients"`
		Enabled     *bool  `yaml:"enabled"`
		Priority    int    `yaml:"priority"`
		PlayCanvas *struct {
			Scene struct {
				AmbientLight interface{} `yaml:"ambientLight,omitempty"` // Can be string or []float64
				Gravity      []float64 `yaml:"gravity,omitempty"`
				Fog          *struct {
					Type  string `yaml:"type,omitempty"`
					Color string `yaml:"color,omitempty"`
					Near  int    `yaml:"near,omitempty"`
					Far   int    `yaml:"far,omitempty"`
				} `yaml:"fog,omitempty"`
			} `yaml:"scene,omitempty"`
			Entities []struct {
				Name       string                 `yaml:"name"`
				Components map[string]interface{} `yaml:"components"`
			} `yaml:"entities,omitempty"`
		} `yaml:"playcanvas,omitempty"`
	}

	if err := yaml.Unmarshal(configData, &fullConfig); err != nil {
		logging.Error("failed to parse world config", map[string]interface{}{
			"error":    err.Error(),
			"world_id": worldId,
		})
		http.Error(w, `{"success": false, "message": "Invalid world configuration"}`, http.StatusInternalServerError)
		return
	}

	// Get file modification time as updated_at
	var updatedAt string
	if stat, err := os.Stat(configPath); err == nil {
		updatedAt = stat.ModTime().Format(time.RFC3339)
	}

	// Get current client count (placeholder - would integrate with session management)
	currentClients := 0 // TODO: Implement actual client counting from session manager

	logging.Debug("world configuration loaded", map[string]interface{}{
		"world_id": worldId,
		"name":     fullConfig.Name,
		"current_clients": currentClients,
		"has_playcanvas":  fullConfig.PlayCanvas != nil,
	})

	// Build response with PlayCanvas configuration
	worldResponse := map[string]interface{}{
		"id":          worldId,
		"name":        fullConfig.Name,
		"description": fullConfig.Description,
		"environment": fullConfig.Environment,
		"max_clients": fullConfig.MaxClients,
		"enabled":     fullConfig.Enabled,
		"priority":    fullConfig.Priority,
		"current_clients": currentClients,
		"updated_at":      updatedAt,
	}

	// Include PlayCanvas configuration if present
	if fullConfig.PlayCanvas != nil {
		worldResponse["playcanvas"] = fullConfig.PlayCanvas
	}

	response := map[string]interface{}{
		"success": true,
		"world":   worldResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}