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

// UpdateWorldHandler handles PUT /worlds/{worldId} - update world configuration
func UpdateWorldHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Extract world ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/worlds/")
	worldId := strings.Split(path, "/")[0]
	
	logging.Debug("world update requested", map[string]interface{}{
		"world_id": worldId,
		"method":     r.Method,
		"path":       r.URL.Path,
	})

	if worldId == "" {
		logging.Warn("missing world ID in update request", map[string]interface{}{
			"path": r.URL.Path,
		})
		http.Error(w, `{"success": false, "message": "World ID is required"}`, http.StatusBadRequest)
		return
	}

	// Parse request body
	var updates struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Environment *string `json:"environment"`
		MaxClients  *int    `json:"max_clients"`
		Enabled     *bool   `json:"enabled"`
		Priority    *int    `json:"priority"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		logging.Error("failed to parse world update request", map[string]interface{}{
			"error":    err.Error(),
			"world_id": worldId,
		})
		http.Error(w, `{"success": false, "message": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate world exists
	worldsDir := "/opt/hd1/share/worlds"
	worldPath := filepath.Join(worldsDir, worldId)
	configPath := filepath.Join(worldPath, "config.yaml")
	
	if _, err := os.Stat(worldPath); os.IsNotExist(err) {
		logging.Warn("world not found for update", map[string]interface{}{
			"world_id": worldId,
			"path":     worldPath,
		})
		http.Error(w, fmt.Sprintf(`{"success": false, "message": "World '%s' not found"}`, worldId), http.StatusNotFound)
		return
	}

	// Load existing configuration
	configData, err := os.ReadFile(configPath)
	if err != nil {
		logging.Error("failed to read world config for update", map[string]interface{}{
			"error":    err.Error(),
			"world_id": worldId,
			"path":       configPath,
		})
		http.Error(w, `{"success": false, "message": "Failed to load world configuration"}`, http.StatusInternalServerError)
		return
	}

	// Parse existing YAML configuration
	var config map[string]interface{}
	if err := yaml.Unmarshal(configData, &config); err != nil {
		logging.Error("failed to parse existing world config", map[string]interface{}{
			"error":    err.Error(),
			"world_id": worldId,
		})
		http.Error(w, `{"success": false, "message": "Invalid existing world configuration"}`, http.StatusInternalServerError)
		return
	}

	// Apply updates to configuration
	updateCount := 0
	if updates.Name != nil {
		config["name"] = *updates.Name
		updateCount++
	}
	if updates.Description != nil {
		config["description"] = *updates.Description
		updateCount++
	}
	if updates.Environment != nil {
		config["environment"] = *updates.Environment
		updateCount++
	}
	if updates.MaxClients != nil {
		if *updates.MaxClients <= 0 {
			logging.Warn("invalid max_clients value", map[string]interface{}{
				"value":    *updates.MaxClients,
				"world_id": worldId,
			})
			http.Error(w, `{"success": false, "message": "max_clients must be greater than 0"}`, http.StatusBadRequest)
			return
		}
		config["max_clients"] = *updates.MaxClients
		updateCount++
	}
	if updates.Enabled != nil {
		config["enabled"] = *updates.Enabled
		updateCount++
	}
	if updates.Priority != nil {
		config["priority"] = *updates.Priority
		updateCount++
	}

	if updateCount == 0 {
		logging.Warn("no valid updates provided", map[string]interface{}{
			"world_id": worldId,
		})
		http.Error(w, `{"success": false, "message": "No valid updates provided"}`, http.StatusBadRequest)
		return
	}

	// Set updated timestamp
	config["updated_at"] = time.Now().Format(time.RFC3339)

	// Marshal back to YAML
	updatedData, err := yaml.Marshal(config)
	if err != nil {
		logging.Error("failed to marshal updated config", map[string]interface{}{
			"error":    err.Error(),
			"world_id": worldId,
		})
		http.Error(w, `{"success": false, "message": "Failed to process configuration updates"}`, http.StatusInternalServerError)
		return
	}

	// Write updated configuration
	if err := os.WriteFile(configPath, updatedData, 0644); err != nil {
		logging.Error("failed to write updated world config", map[string]interface{}{
			"error":    err.Error(),
			"world_id": worldId,
			"path":       configPath,
		})
		http.Error(w, `{"success": false, "message": "Failed to save world configuration"}`, http.StatusInternalServerError)
		return
	}

	logging.Info("world configuration updated", map[string]interface{}{
		"world_id": worldId,
		"updates":      updateCount,
		"name":         config["name"],
		"environment":  config["environment"],
		"max_clients":  config["max_clients"],
		"enabled":      config["enabled"],
	})

	// Build response with updated configuration
	response := map[string]interface{}{
		"success": true,
		"message": "World updated successfully",
		"world": map[string]interface{}{
			"id": worldId,
			"name":        config["name"],
			"description": config["description"],
			"environment": config["environment"],
			"max_clients": config["max_clients"],
			"enabled":     config["enabled"],
			"priority":    config["priority"],
			"updated_at":  config["updated_at"],
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}