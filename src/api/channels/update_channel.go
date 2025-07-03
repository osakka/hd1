package channels

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

// UpdateChannelHandler handles PUT /channels/{channelId} - update channel configuration
func UpdateChannelHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Extract channel ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/channels/")
	channelId := strings.Split(path, "/")[0]
	
	logging.Debug("channel update requested", map[string]interface{}{
		"channel_id": channelId,
		"method":     r.Method,
		"path":       r.URL.Path,
	})

	if channelId == "" {
		logging.Warn("missing channel ID in update request", map[string]interface{}{
			"path": r.URL.Path,
		})
		http.Error(w, `{"success": false, "message": "Channel ID is required"}`, http.StatusBadRequest)
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
		logging.Error("failed to parse channel update request", map[string]interface{}{
			"error":      err.Error(),
			"channel_id": channelId,
		})
		http.Error(w, `{"success": false, "message": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate channel exists
	channelsDir := "/opt/hd1/share/channels"
	channelPath := filepath.Join(channelsDir, channelId)
	configPath := filepath.Join(channelPath, "config.yaml")
	
	if _, err := os.Stat(channelPath); os.IsNotExist(err) {
		logging.Warn("channel not found for update", map[string]interface{}{
			"channel_id": channelId,
			"path":       channelPath,
		})
		http.Error(w, fmt.Sprintf(`{"success": false, "message": "Channel '%s' not found"}`, channelId), http.StatusNotFound)
		return
	}

	// Load existing configuration
	configData, err := os.ReadFile(configPath)
	if err != nil {
		logging.Error("failed to read channel config for update", map[string]interface{}{
			"error":      err.Error(),
			"channel_id": channelId,
			"path":       configPath,
		})
		http.Error(w, `{"success": false, "message": "Failed to load channel configuration"}`, http.StatusInternalServerError)
		return
	}

	// Parse existing YAML configuration
	var config map[string]interface{}
	if err := yaml.Unmarshal(configData, &config); err != nil {
		logging.Error("failed to parse existing channel config", map[string]interface{}{
			"error":      err.Error(),
			"channel_id": channelId,
		})
		http.Error(w, `{"success": false, "message": "Invalid existing channel configuration"}`, http.StatusInternalServerError)
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
				"value":      *updates.MaxClients,
				"channel_id": channelId,
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
			"channel_id": channelId,
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
			"error":      err.Error(),
			"channel_id": channelId,
		})
		http.Error(w, `{"success": false, "message": "Failed to process configuration updates"}`, http.StatusInternalServerError)
		return
	}

	// Write updated configuration
	if err := os.WriteFile(configPath, updatedData, 0644); err != nil {
		logging.Error("failed to write updated channel config", map[string]interface{}{
			"error":      err.Error(),
			"channel_id": channelId,
			"path":       configPath,
		})
		http.Error(w, `{"success": false, "message": "Failed to save channel configuration"}`, http.StatusInternalServerError)
		return
	}

	logging.Info("channel configuration updated", map[string]interface{}{
		"channel_id":   channelId,
		"updates":      updateCount,
		"name":         config["name"],
		"environment":  config["environment"],
		"max_clients":  config["max_clients"],
		"enabled":      config["enabled"],
	})

	// Build response with updated configuration
	response := map[string]interface{}{
		"success": true,
		"message": "Channel updated successfully",
		"channel": map[string]interface{}{
			"id":          channelId,
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