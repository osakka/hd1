package channels

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

// ChannelInfo represents information about an available channel
type ChannelInfo struct {
	ID          string            `json:"id" yaml:"id"`
	Name        string            `json:"name" yaml:"name"`
	Description string            `json:"description" yaml:"description"`
	Environment string            `json:"environment" yaml:"type"`
	MaxClients  int               `json:"max_clients" yaml:"max_clients"`
	PlayCanvas  *PlayCanvasConfig `json:"playcanvas,omitempty" yaml:"playcanvas,omitempty"`
}

// ChannelConfig represents the full channel configuration
type ChannelConfig struct {
	Channel     ChannelInfo       `yaml:"channel"`
	PlayCanvas  *PlayCanvasConfig `yaml:"playcanvas,omitempty"`
	Settings    struct {
		MaxClients int `yaml:"max_clients"`
	} `yaml:"settings,omitempty"`
	Environment struct {
		Type string `yaml:"type"`
	} `yaml:"environment,omitempty"`
}

// ServerConfig represents the server-side channel configuration
type ServerConfig struct {
	Defaults struct {
		DefaultChannel           string `yaml:"default_channel"`
		FallbackChannel          string `yaml:"fallback_channel"`
		AutoJoinOnSessionCreate  bool   `yaml:"auto_join_on_session_create"`
	} `yaml:"defaults"`
	Channels map[string]struct {
		Enabled             bool   `yaml:"enabled"`
		Priority            int    `yaml:"priority"`
		Description         string `yaml:"description"`
		MaxClients          int    `yaml:"max_clients"`
		AutoApplyEnvironment bool   `yaml:"auto_apply_environment"`
	} `yaml:"channels"`
	Server struct {
		CleanupEmptyChannelsAfterMinutes int  `yaml:"cleanup_empty_channels_after_minutes"`
		MaxChannelSwitchesPerMinute      int  `yaml:"max_channel_switches_per_minute"`
		SyncEnvironmentOnJoin            bool `yaml:"sync_environment_on_join"`
		LogChannelEvents                 bool `yaml:"log_channel_events"`
	} `yaml:"server"`
}

// Enhanced response with server defaults
type ListChannelsResponseEnhanced struct {
	Success        bool          `json:"success"`
	Message        string        `json:"message"`
	Channels       []ChannelInfo `json:"channels"`
	DefaultChannel string        `json:"default_channel"`
	ServerConfig   struct {
		AutoJoin              bool `json:"auto_join_on_session_create"`
		SyncEnvironmentOnJoin bool `json:"sync_environment_on_join"`
	} `json:"server_config"`
}

// ListChannelsResponse represents the response for listing available channels
type ListChannelsResponse struct {
	Success  bool          `json:"success"`
	Message  string        `json:"message"`
	Channels []ChannelInfo `json:"channels"`
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

// ListChannelsHandler handles requests to list available channels
func ListChannelsHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Get channels directory from configuration
	channelsDir := config.GetChannelsDir()
	
	logging.Info("listing available channels", map[string]interface{}{
		"channels_dir": channelsDir,
	})
	
	// Read server configuration
	serverConfig, err := loadServerConfig(channelsDir)
	if err != nil {
		logging.Error("failed to load server config", map[string]interface{}{
			"error": err.Error(),
		})
		// Continue with defaults if config fails to load
		serverConfig = &ServerConfig{}
		serverConfig.Defaults.DefaultChannel = "channel_one"
		serverConfig.Defaults.AutoJoinOnSessionCreate = true
		serverConfig.Server.SyncEnvironmentOnJoin = true
	}

	// Read channel configuration files
	channels, err := discoverChannels(channelsDir, serverConfig)
	if err != nil {
		logging.Error("failed to discover channels", map[string]interface{}{
			"error": err.Error(),
			"channels_dir": channelsDir,
		})
		http.Error(w, "Failed to discover channels", http.StatusInternalServerError)
		return
	}

	// Build enhanced response
	response := ListChannelsResponseEnhanced{
		Success:        true,
		Message:        "Available channels retrieved successfully",
		Channels:       channels,
		DefaultChannel: serverConfig.Defaults.DefaultChannel,
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
		logging.Error("failed to encode channels response", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

	logging.Info("channels list response sent", map[string]interface{}{
		"channel_count": len(channels),
	})
}

// loadServerConfig loads the server configuration from config.yaml
func loadServerConfig(channelsDir string) (*ServerConfig, error) {
	configPath := filepath.Join(channelsDir, "config.yaml")
	
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	
	var config ServerConfig
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, err
	}
	
	logging.Info("server channel configuration loaded", map[string]interface{}{
		"default_channel": config.Defaults.DefaultChannel,
		"auto_join": config.Defaults.AutoJoinOnSessionCreate,
		"channel_count": len(config.Channels),
	})
	
	return &config, nil
}

// discoverChannels scans the channels directory for available channel configurations
func discoverChannels(channelsDir string, serverConfig *ServerConfig) ([]ChannelInfo, error) {
	var channels []ChannelInfo

	// Read directory contents
	files, err := ioutil.ReadDir(channelsDir)
	if err != nil {
		return nil, err
	}

	// Process each YAML configuration file
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}

		configPath := filepath.Join(channelsDir, file.Name())
		
		// Read and parse channel configuration
		configData, err := ioutil.ReadFile(configPath)
		if err != nil {
			logging.Error("failed to read channel config", map[string]interface{}{
				"file": configPath,
				"error": err.Error(),
			})
			continue
		}

		// Parse YAML configuration with interface{} types first
		var rawConfig map[string]interface{}
		if err := yaml.Unmarshal(configData, &rawConfig); err != nil {
			logging.Error("failed to parse channel config", map[string]interface{}{
				"file": configPath,
				"error": err.Error(),
			})
			continue
		}

		// Convert to proper JSON-serializable format
		convertedConfig := convertMapInterfaceToStringInterface(rawConfig).(map[string]interface{})

		// Re-marshal and unmarshal to get proper struct types
		convertedData, err := json.Marshal(convertedConfig)
		if err != nil {
			logging.Error("failed to convert channel config", map[string]interface{}{
				"file": configPath,
				"error": err.Error(),
			})
			continue
		}

		var config ChannelConfig
		if err := json.Unmarshal(convertedData, &config); err != nil {
			logging.Error("failed to unmarshal converted channel config", map[string]interface{}{
				"file": configPath,
				"error": err.Error(),
			})
			continue
		}

		// Create channel info
		channelInfo := ChannelInfo{
			ID:          config.Channel.ID,
			Name:        config.Channel.Name,
			Description: config.Channel.Description,
			Environment: config.Environment.Type,
			MaxClients:  config.Settings.MaxClients,
			PlayCanvas:  config.PlayCanvas,
		}

		// Apply server configuration overrides if available
		if serverConfig != nil && serverConfig.Channels != nil {
			if channelOverride, exists := serverConfig.Channels[channelInfo.ID]; exists {
				// Skip disabled channels
				if !channelOverride.Enabled {
					logging.Debug("skipping disabled channel", map[string]interface{}{
						"channel_id": channelInfo.ID,
					})
					continue
				}
				
				// Apply overrides
				if channelOverride.Description != "" {
					channelInfo.Description = channelOverride.Description
				}
				if channelOverride.MaxClients > 0 {
					channelInfo.MaxClients = channelOverride.MaxClients
				}
			}
		}

		channels = append(channels, channelInfo)
		
		logging.Debug("discovered channel", map[string]interface{}{
			"channel_id": channelInfo.ID,
			"name": channelInfo.Name,
			"environment": channelInfo.Environment,
		})
	}

	return channels, nil
}