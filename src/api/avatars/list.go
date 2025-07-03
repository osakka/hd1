package avatars

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"holodeck1/logging"
	"gopkg.in/yaml.v3"
)

// AvatarType represents basic avatar information
type AvatarType struct {
	Name        string   `json:"name" yaml:"name"`
	Type        string   `json:"type" yaml:"type"`
	Description string   `json:"description" yaml:"description"`
	Contexts    []string `json:"contexts" yaml:"contexts"`
	AssetPath   string   `json:"asset_path" yaml:"path"`
}

// AvatarConfig represents the avatar system configuration
type AvatarConfig struct {
	System struct {
		Name        string `yaml:"name"`
		Version     string `yaml:"version"`
		Description string `yaml:"description"`
	} `yaml:"system"`
	AvatarTypes map[string]struct {
		Path        string   `yaml:"path"`
		Description string   `yaml:"description"`
		Contexts    []string `yaml:"contexts"`
	} `yaml:"avatar_types"`
}

// ListAvatarsHandler handles GET /avatars
func ListAvatarsHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	logger := logging.GetLogger()
	logger.Info("Listing available avatar types")

	// Read avatar config
	configPath := "/opt/hd1/share/avatars/config.yaml"
	configData, err := os.ReadFile(configPath)
	if err != nil {
		logger.Error("Failed to read avatar config", map[string]interface{}{
			"error": err,
			"path":  configPath,
		})
		http.Error(w, "Failed to read avatar configuration", http.StatusInternalServerError)
		return
	}

	var config AvatarConfig
	if err := yaml.Unmarshal(configData, &config); err != nil {
		logger.Error("Failed to parse avatar config", map[string]interface{}{
			"error": err,
		})
		http.Error(w, "Failed to parse avatar configuration", http.StatusInternalServerError)
		return
	}

	// Build avatar types list
	var avatars []AvatarType
	for avatarType, avatarInfo := range config.AvatarTypes {
		// Read the individual avatar spec to get the name
		avatarSpecPath := filepath.Join("/opt/hd1/share/avatars", avatarInfo.Path)
		specData, err := os.ReadFile(avatarSpecPath)
		if err != nil {
			logger.Warn("Could not read avatar spec", map[string]interface{}{
				"type":  avatarType,
				"path":  avatarSpecPath,
				"error": err,
			})
			continue
		}

		var spec struct {
			Metadata struct {
				Name        string `yaml:"name"`
				Description string `yaml:"description"`
			} `yaml:"metadata"`
		}
		if err := yaml.Unmarshal(specData, &spec); err != nil {
			logger.Warn("Could not parse avatar spec", map[string]interface{}{
				"type":  avatarType,
				"error": err,
			})
			continue
		}

		avatars = append(avatars, AvatarType{
			Name:        spec.Metadata.Name,
			Type:        avatarType,
			Description: avatarInfo.Description,
			Contexts:    avatarInfo.Contexts,
			AssetPath:   avatarInfo.Path,
		})
	}

	// Return response
	response := map[string]interface{}{
		"success": true,
		"avatars": avatars,
		"message": "Avatar types retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Failed to encode avatar list response", map[string]interface{}{
			"error": err,
		})
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.Info("Avatar types listed successfully", map[string]interface{}{
		"count": len(avatars),
	})
}