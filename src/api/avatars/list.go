package avatars

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"holodeck1/config"
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
	// TRACE: Module-specific debugging for avatar operations
	logging.Trace("avatars", "avatar list request initiated")

	// Read avatar config from configured path
	configPath := config.GetAvatarsConfigFile()
	configData, err := os.ReadFile(configPath)
	if err != nil {
		logging.Error("avatar config read failed", map[string]interface{}{
			"error": err.Error(),
			"path":  configPath,
		})
		http.Error(w, "Failed to read avatar configuration", http.StatusInternalServerError)
		return
	}

	var avatarConfig AvatarConfig
	if err := yaml.Unmarshal(configData, &avatarConfig); err != nil {
		logging.Error("avatar config parse failed", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Failed to parse avatar configuration", http.StatusInternalServerError)
		return
	}

	// Build avatar types list
	var avatars []AvatarType
	for avatarType, avatarInfo := range avatarConfig.AvatarTypes {
		// Read the individual avatar spec to get the name
		avatarSpecPath := filepath.Join(config.GetAvatarsDir(), avatarInfo.Path)
		specData, err := os.ReadFile(avatarSpecPath)
		if err != nil {
			logging.Warn("avatar spec read failed", map[string]interface{}{
				"type":  avatarType,
				"path":  avatarSpecPath,
				"error": err.Error(),
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
			logging.Warn("avatar spec parse failed", map[string]interface{}{
				"type":  avatarType,
				"error": err.Error(),
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
		logging.Error("response encoding failed", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// INFO: Production-appropriate completion logging
	logging.Info("avatar types retrieved", map[string]interface{}{
		"count": len(avatars),
	})
}