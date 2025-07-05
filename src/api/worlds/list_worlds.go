package worlds

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"holodeck1/config"
	"holodeck1/logging"
	"gopkg.in/yaml.v3"
)

// PlayCanvasScene represents PlayCanvas scene configuration
type PlayCanvasScene struct {
	AmbientLight interface{} `json:"ambientLight,omitempty" yaml:"ambientLight,omitempty"` // Can be string or []float64
	Gravity      []float64   `json:"gravity,omitempty" yaml:"gravity,omitempty"`
	Fog          *struct {
		Type  string `json:"type,omitempty" yaml:"type,omitempty"`
		Color string `json:"color,omitempty" yaml:"color,omitempty"`
		Near  int    `json:"near,omitempty" yaml:"near,omitempty"`
		Far   int    `json:"far,omitempty" yaml:"far,omitempty"`
	} `json:"fog,omitempty" yaml:"fog,omitempty"`
}

// PlayCanvasEntity represents a native PlayCanvas entity configuration
type PlayCanvasEntity struct {
	Name       string                 `json:"name" yaml:"name"`
	Components map[string]interface{} `json:"components" yaml:"components"`
}

// PlayCanvasConfig represents the full PlayCanvas configuration section
type PlayCanvasConfig struct {
	Scene    PlayCanvasScene    `json:"scene,omitempty" yaml:"scene,omitempty"`
	Entities []PlayCanvasEntity `json:"entities,omitempty" yaml:"entities,omitempty"`
}

// WorldInfo represents information about an available world
type WorldInfo struct {
	ID          string            `json:"id" yaml:"id"`
	Name        string            `json:"name" yaml:"name"`
	Description string            `json:"description" yaml:"description"`
	Environment string            `json:"environment" yaml:"type"`
	MaxClients  int               `json:"max_clients" yaml:"max_clients"`
	PlayCanvas  *PlayCanvasConfig `json:"playcanvas,omitempty" yaml:"playcanvas,omitempty"`
}

// WorldConfig represents the full world configuration
type WorldConfig struct {
	World       WorldInfo         `yaml:"world"`
	PlayCanvas  *PlayCanvasConfig `yaml:"playcanvas,omitempty"`
	Settings    struct {
		MaxClients int `yaml:"max_clients"`
	} `yaml:"settings,omitempty"`
	Environment struct {
		Type string `yaml:"type"`
	} `yaml:"environment,omitempty"`
}

// ServerConfig represents the server-side world configuration
type ServerConfig struct {
	Defaults struct {
		DefaultWorld            string `yaml:"default_world"`
		FallbackWorld           string `yaml:"fallback_world"`
		AutoJoinOnSessionCreate bool   `yaml:"auto_join_on_session_create"`
	} `yaml:"defaults"`
	Worlds map[string]struct {
		Enabled             bool   `yaml:"enabled"`
		Priority            int    `yaml:"priority"`
		Description         string `yaml:"description"`
		MaxClients          int    `yaml:"max_clients"`
		AutoApplyEnvironment bool   `yaml:"auto_apply_environment"`
	} `yaml:"worlds"`
	Server struct {
		CleanupEmptyWorldsAfterMinutes int  `yaml:"cleanup_empty_worlds_after_minutes"`
		MaxWorldSwitchesPerMinute      int  `yaml:"max_world_switches_per_minute"`
		SyncEnvironmentOnJoin          bool `yaml:"sync_environment_on_join"`
		LogWorldEvents                 bool `yaml:"log_world_events"`
	} `yaml:"server"`
}

// Enhanced response with server defaults
type ListWorldsResponseEnhanced struct {
	Success      bool        `json:"success"`
	Message      string      `json:"message"`
	Worlds       []WorldInfo `json:"worlds"`
	DefaultWorld string      `json:"default_world"`
	ServerConfig struct {
		AutoJoin              bool `json:"auto_join_on_session_create"`
		SyncEnvironmentOnJoin bool `json:"sync_environment_on_join"`
	} `json:"server_config"`
}

// ListWorldsResponse represents the response for listing available worlds
type ListWorldsResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Worlds  []WorldInfo `json:"worlds"`
}

// convertMapInterfaceToStringInterface recursively converts map[interface{}]interface{} to map[string]interface{}
func convertMapInterfaceToStringInterface(input interface{}) interface{} {
	switch v := input.(type) {
	case map[interface{}]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			if strKey, ok := key.(string); ok {
				result[strKey] = convertMapInterfaceToStringInterface(value)
			}
		}
		return result
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			result[key] = convertMapInterfaceToStringInterface(value)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = convertMapInterfaceToStringInterface(item)
		}
		return result
	default:
		return v
	}
}

// ListWorldsHandler handles requests to list available worlds
func ListWorldsHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Get worlds directory from configuration
	worldsDir := config.GetWorldsDir()
	
	logging.Info("listing available worlds", map[string]interface{}{
		"worlds_dir": worldsDir,
	})
	
	// Read server configuration
	serverConfig, err := loadServerConfig(worldsDir)
	if err != nil {
		logging.Error("failed to load server config", map[string]interface{}{
			"error": err.Error(),
		})
		// Continue with defaults if config fails to load
		serverConfig = &ServerConfig{}
		serverConfig.Defaults.DefaultWorld = "world_one"
		serverConfig.Defaults.AutoJoinOnSessionCreate = true
		serverConfig.Server.SyncEnvironmentOnJoin = true
	}

	// Read world configuration files
	worlds, err := discoverWorlds(worldsDir, serverConfig)
	if err != nil {
		logging.Error("failed to discover worlds", map[string]interface{}{
			"error":      err.Error(),
			"worlds_dir": worldsDir,
		})
		http.Error(w, "Failed to discover worlds", http.StatusInternalServerError)
		return
	}

	// Build enhanced response
	response := ListWorldsResponseEnhanced{
		Success:      true,
		Message:      "Available worlds retrieved successfully",
		Worlds:       worlds,
		DefaultWorld: serverConfig.Defaults.DefaultWorld,
		ServerConfig: struct {
			AutoJoin              bool `json:"auto_join_on_session_create"`
			SyncEnvironmentOnJoin bool `json:"sync_environment_on_join"`
		}{
			AutoJoin:              serverConfig.Defaults.AutoJoinOnSessionCreate,
			SyncEnvironmentOnJoin: serverConfig.Server.SyncEnvironmentOnJoin,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logging.Error("failed to encode worlds response", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

	logging.Info("worlds list response sent", map[string]interface{}{
		"world_count": len(worlds),
	})
}

// loadServerConfig loads the server configuration from config.yaml
func loadServerConfig(worldsDir string) (*ServerConfig, error) {
	configPath := filepath.Join(worldsDir, "config.yaml")
	
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	
	var config ServerConfig
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, err
	}
	
	logging.Info("server world configuration loaded", map[string]interface{}{
		"default_world": config.Defaults.DefaultWorld,
		"auto_join":     config.Defaults.AutoJoinOnSessionCreate,
		"world_count":   len(config.Worlds),
	})
	
	return &config, nil
}

// discoverWorlds scans the worlds directory for available world configurations
func discoverWorlds(worldsDir string, serverConfig *ServerConfig) ([]WorldInfo, error) {
	var worlds []WorldInfo

	// Read directory contents
	files, err := ioutil.ReadDir(worldsDir)
	if err != nil {
		return nil, err
	}

	// Process each YAML configuration file
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}

		configPath := filepath.Join(worldsDir, file.Name())
		
		// Read and parse world configuration
		configData, err := ioutil.ReadFile(configPath)
		if err != nil {
			logging.Error("failed to read world config", map[string]interface{}{
				"file":  configPath,
				"error": err.Error(),
			})
			continue
		}

		// Parse YAML configuration with interface{} types first
		var rawConfig map[string]interface{}
		if err := yaml.Unmarshal(configData, &rawConfig); err != nil {
			logging.Error("failed to parse world config", map[string]interface{}{
				"file":  configPath,
				"error": err.Error(),
			})
			continue
		}

		// Convert to proper JSON-serializable format
		convertedConfig := convertMapInterfaceToStringInterface(rawConfig).(map[string]interface{})

		// Re-marshal and unmarshal to get proper struct types
		convertedData, err := json.Marshal(convertedConfig)
		if err != nil {
			logging.Error("failed to convert world config", map[string]interface{}{
				"file":  configPath,
				"error": err.Error(),
			})
			continue
		}

		var config WorldConfig
		if err := json.Unmarshal(convertedData, &config); err != nil {
			logging.Error("failed to unmarshal converted world config", map[string]interface{}{
				"file":  configPath,
				"error": err.Error(),
			})
			continue
		}

		// Create world info
		worldInfo := WorldInfo{
			ID:          config.World.ID,
			Name:        config.World.Name,
			Description: config.World.Description,
			Environment: config.Environment.Type,
			MaxClients:  config.Settings.MaxClients,
			PlayCanvas:  config.PlayCanvas,
		}

		// Apply server configuration overrides if available
		if serverConfig != nil && serverConfig.Worlds != nil {
			if worldOverride, exists := serverConfig.Worlds[worldInfo.ID]; exists {
				// Skip disabled worlds
				if !worldOverride.Enabled {
					logging.Debug("skipping disabled world", map[string]interface{}{
						"world_id": worldInfo.ID,
					})
					continue
				}
				
				// Apply overrides
				if worldOverride.Description != "" {
					worldInfo.Description = worldOverride.Description
				}
				if worldOverride.MaxClients > 0 {
					worldInfo.MaxClients = worldOverride.MaxClients
				}
			}
		}

		worlds = append(worlds, worldInfo)
		
		logging.Debug("discovered world", map[string]interface{}{
			"world_id":    worldInfo.ID,
			"name":        worldInfo.Name,
			"environment": worldInfo.Environment,
		})
	}

	return worlds, nil
}

