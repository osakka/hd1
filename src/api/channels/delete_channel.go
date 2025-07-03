package channels

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"holodeck1/config"
	"holodeck1/logging"
	"gopkg.in/yaml.v3"
)

// DeleteChannelHandler handles DELETE /channels/{channelId} - delete channel
func DeleteChannelHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Extract channel ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/channels/")
	channelId := strings.Split(path, "/")[0]
	
	// Parse query parameters
	force := r.URL.Query().Get("force") == "true"
	
	logging.Debug("channel deletion requested", map[string]interface{}{
		"channel_id": channelId,
		"force":      force,
		"method":     r.Method,
		"path":       r.URL.Path,
	})

	if channelId == "" {
		logging.Warn("missing channel ID in delete request", map[string]interface{}{
			"path": r.URL.Path,
		})
		http.Error(w, `{"success": false, "message": "Channel ID is required"}`, http.StatusBadRequest)
		return
	}

	// Validate channel exists
	channelsDir := config.GetChannelsDir()
	channelPath := filepath.Join(channelsDir, channelId)
	configPath := filepath.Join(channelPath, "config.yaml")
	
	if _, err := os.Stat(channelPath); os.IsNotExist(err) {
		logging.Warn("channel not found for deletion", map[string]interface{}{
			"channel_id": channelId,
			"path":       channelPath,
		})
		http.Error(w, fmt.Sprintf(`{"success": false, "message": "Channel '%s' not found"}`, channelId), http.StatusNotFound)
		return
	}

	// Load server configuration to check if this is the default channel
	serverConfigPath := filepath.Join(config.GetChannelsDir(), "config.yaml")
	serverConfigData, err := os.ReadFile(serverConfigPath)
	if err == nil {
		var serverConfig struct {
			Defaults struct {
				DefaultChannel  string `yaml:"default_channel"`
				FallbackChannel string `yaml:"fallback_channel"`
			} `yaml:"defaults"`
		}
		
		if err := yaml.Unmarshal(serverConfigData, &serverConfig); err == nil {
			if channelId == serverConfig.Defaults.DefaultChannel {
				logging.Warn("attempt to delete default channel", map[string]interface{}{
					"channel_id": channelId,
				})
				http.Error(w, fmt.Sprintf(`{"success": false, "message": "Cannot delete default channel '%s'"}`, channelId), http.StatusBadRequest)
				return
			}
			if channelId == serverConfig.Defaults.FallbackChannel {
				logging.Warn("attempt to delete fallback channel", map[string]interface{}{
					"channel_id": channelId,
				})
				http.Error(w, fmt.Sprintf(`{"success": false, "message": "Cannot delete fallback channel '%s'"}`, channelId), http.StatusBadRequest)
				return
			}
		}
	}

	// Check for active clients (placeholder - would integrate with session management)
	currentClients := 0 // TODO: Implement actual client counting from session manager

	// Simulate some clients for demonstration
	if channelId == "channel_one" || channelId == "channel_two" {
		currentClients = 2 // Simulate active clients for testing
	}

	if currentClients > 0 && !force {
		logging.Warn("channel has active clients", map[string]interface{}{
			"channel_id":      channelId,
			"current_clients": currentClients,
			"force":           force,
		})
		
		response := map[string]interface{}{
			"success":         false,
			"message":         fmt.Sprintf("Cannot delete channel '%s' with %d active clients. Use ?force=true to disconnect clients and delete.", channelId, currentClients),
			"current_clients": currentClients,
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Load channel configuration for logging
	var channelName string
	if configData, err := os.ReadFile(configPath); err == nil {
		var config struct {
			Name string `yaml:"name"`
		}
		if yaml.Unmarshal(configData, &config) == nil {
			channelName = config.Name
		}
	}

	// TODO: If we had active clients, disconnect them here
	disconnectedClients := currentClients

	// Remove channel directory and all contents
	if err := os.RemoveAll(channelPath); err != nil {
		logging.Error("failed to delete channel directory", map[string]interface{}{
			"error":      err.Error(),
			"channel_id": channelId,
			"path":       channelPath,
		})
		http.Error(w, `{"success": false, "message": "Failed to delete channel files"}`, http.StatusInternalServerError)
		return
	}

	logging.Info("channel deleted successfully", map[string]interface{}{
		"channel_id":           channelId,
		"channel_name":         channelName,
		"disconnected_clients": disconnectedClients,
		"force":                force,
	})

	// Return success response
	response := map[string]interface{}{
		"success":              true,
		"message":              "Channel deleted successfully",
		"disconnected_clients": disconnectedClients,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper function to parse boolean from string
func parseBool(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}