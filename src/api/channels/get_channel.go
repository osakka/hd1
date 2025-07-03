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

// GetChannelHandler handles GET /channels/{channelId} - get specific channel details
func GetChannelHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Extract channel ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/channels/")
	channelId := strings.Split(path, "/")[0]
	
	logging.Debug("channel details requested", map[string]interface{}{
		"channel_id": channelId,
		"method":     r.Method,
		"path":       r.URL.Path,
	})

	if channelId == "" {
		logging.Warn("missing channel ID in request", map[string]interface{}{
			"path": r.URL.Path,
		})
		http.Error(w, `{"success": false, "message": "Channel ID is required"}`, http.StatusBadRequest)
		return
	}

	// Validate channel exists
	channelsDir := "/opt/hd1/share/channels"
	configPath := filepath.Join(channelsDir, channelId+".yaml")
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logging.Warn("channel not found", map[string]interface{}{
			"channel_id": channelId,
			"path":       configPath,
		})
		http.Error(w, fmt.Sprintf(`{"success": false, "message": "Channel '%s' not found"}`, channelId), http.StatusNotFound)
		return
	}

	// Load channel configuration
	configData, err := os.ReadFile(configPath)
	if err != nil {
		logging.Error("failed to read channel config", map[string]interface{}{
			"error":      err.Error(),
			"channel_id": channelId,
			"path":       configPath,
		})
		http.Error(w, `{"success": false, "message": "Failed to load channel configuration"}`, http.StatusInternalServerError)
		return
	}

	// Parse YAML configuration using the same structs as list_channels.go
	var fullConfig struct {
		Channel struct {
			ID          string `yaml:"id"`
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
			MaxClients  int    `yaml:"max_clients"`
		} `yaml:"channel"`
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
		logging.Error("failed to parse channel config", map[string]interface{}{
			"error":      err.Error(),
			"channel_id": channelId,
		})
		http.Error(w, `{"success": false, "message": "Invalid channel configuration"}`, http.StatusInternalServerError)
		return
	}

	// Get file modification time as updated_at
	var updatedAt string
	if stat, err := os.Stat(configPath); err == nil {
		updatedAt = stat.ModTime().Format(time.RFC3339)
	}

	// Get current client count (placeholder - would integrate with session management)
	currentClients := 0 // TODO: Implement actual client counting from session manager

	logging.Debug("channel configuration loaded", map[string]interface{}{
		"channel_id":      channelId,
		"name":            fullConfig.Channel.Name,
		"current_clients": currentClients,
		"has_playcanvas":  fullConfig.PlayCanvas != nil,
	})

	// Build response with PlayCanvas configuration
	channelResponse := map[string]interface{}{
		"id":              fullConfig.Channel.ID,
		"name":            fullConfig.Channel.Name,
		"description":     fullConfig.Channel.Description,
		"max_clients":     fullConfig.Channel.MaxClients,
		"current_clients": currentClients,
		"updated_at":      updatedAt,
	}

	// Include PlayCanvas configuration if present
	if fullConfig.PlayCanvas != nil {
		channelResponse["playcanvas"] = fullConfig.PlayCanvas
	}

	response := map[string]interface{}{
		"success": true,
		"channel": channelResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}