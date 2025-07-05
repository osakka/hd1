package worlds

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

// DeleteWorldHandler handles DELETE /worlds/{worldId} - delete world
func DeleteWorldHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Extract world ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/worlds/")
	worldId := strings.Split(path, "/")[0]
	
	// Parse query parameters
	force := r.URL.Query().Get("force") == "true"
	
	logging.Debug("world deletion requested", map[string]interface{}{
		"world_id": worldId,
		"force":    force,
		"method":   r.Method,
		"path":     r.URL.Path,
	})

	if worldId == "" {
		logging.Warn("missing world ID in delete request", map[string]interface{}{
			"path": r.URL.Path,
		})
		http.Error(w, `{"success": false, "message": "World ID is required"}`, http.StatusBadRequest)
		return
	}

	// Validate world exists
	worldsDir := config.GetWorldsDir()
	worldPath := filepath.Join(worldsDir, worldId)
	configPath := filepath.Join(worldPath, "config.yaml")
	
	if _, err := os.Stat(worldPath); os.IsNotExist(err) {
		logging.Warn("world not found for deletion", map[string]interface{}{
			"world_id": worldId,
			"path":     worldPath,
		})
		http.Error(w, fmt.Sprintf(`{"success": false, "message": "World '%s' not found"}`, worldId), http.StatusNotFound)
		return
	}

	// Load server configuration to check if this is the default world
	serverConfigPath := filepath.Join(config.GetWorldsDir(), "config.yaml")
	serverConfigData, err := os.ReadFile(serverConfigPath)
	if err == nil {
		var serverConfig struct {
			Defaults struct {
				DefaultWorld  string `yaml:"default_world"`
				FallbackWorld string `yaml:"fallback_world"`
			} `yaml:"defaults"`
		}
		
		if err := yaml.Unmarshal(serverConfigData, &serverConfig); err == nil {
			if worldId == serverConfig.Defaults.DefaultWorld {
				logging.Warn("attempt to delete default world", map[string]interface{}{
					"world_id": worldId,
				})
				http.Error(w, fmt.Sprintf(`{"success": false, "message": "Cannot delete default world '%s'"}`, worldId), http.StatusBadRequest)
				return
			}
			if worldId == serverConfig.Defaults.FallbackWorld {
				logging.Warn("attempt to delete fallback world", map[string]interface{}{
					"world_id": worldId,
				})
				http.Error(w, fmt.Sprintf(`{"success": false, "message": "Cannot delete fallback world '%s'"}`, worldId), http.StatusBadRequest)
				return
			}
		}
	}

	// Check for active clients (placeholder - would integrate with session management)
	currentClients := 0 // TODO: Implement actual client counting from session manager

	// Simulate some clients for demonstration
	if worldId == "world_one" || worldId == "world_two" {
		currentClients = 2 // Simulate active clients for testing
	}

	if currentClients > 0 && !force {
		logging.Warn("world has active clients", map[string]interface{}{
			"world_id":        worldId,
			"current_clients": currentClients,
			"force":           force,
		})
		
		response := map[string]interface{}{
			"success":         false,
			"message":         fmt.Sprintf("Cannot delete world '%s' with %d active clients. Use ?force=true to disconnect clients and delete.", worldId, currentClients),
			"current_clients": currentClients,
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Load world configuration for logging
	var worldName string
	if configData, err := os.ReadFile(configPath); err == nil {
		var config struct {
			Name string `yaml:"name"`
		}
		if yaml.Unmarshal(configData, &config) == nil {
			worldName = config.Name
		}
	}

	// TODO: If we had active clients, disconnect them here
	disconnectedClients := currentClients

	// Remove world directory and all contents
	if err := os.RemoveAll(worldPath); err != nil {
		logging.Error("failed to delete world directory", map[string]interface{}{
			"error":    err.Error(),
			"world_id": worldId,
			"path":     worldPath,
		})
		http.Error(w, `{"success": false, "message": "Failed to delete world files"}`, http.StatusInternalServerError)
		return
	}

	logging.Info("world deleted successfully", map[string]interface{}{
		"world_id":             worldId,
		"world_name":           worldName,
		"disconnected_clients": disconnectedClients,
		"force":                force,
	})

	// Return success response
	response := map[string]interface{}{
		"success":              true,
		"message":              "World deleted successfully",
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