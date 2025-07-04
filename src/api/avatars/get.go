package avatars

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"holodeck1/config"
	"holodeck1/logging"
	"gopkg.in/yaml.v3"
)

// AvatarSpecification represents the complete avatar specification
type AvatarSpecification struct {
	Metadata struct {
		Name        string `json:"name" yaml:"name"`
		Description string `json:"description" yaml:"description"`
		Version     string `json:"version" yaml:"version"`
		Type        string `json:"type" yaml:"type"`
		Created     string `json:"created" yaml:"created"`
	} `json:"metadata" yaml:"metadata"`
	Assets struct {
		Model      string   `json:"model" yaml:"model"`
		Textures   []string `json:"textures" yaml:"textures"`
		Animations []string `json:"animations" yaml:"animations"`
	} `json:"assets" yaml:"assets"`
	Physical struct {
		DefaultTransforms struct {
			Position map[string]float64 `json:"position" yaml:"position"`
			Rotation map[string]float64 `json:"rotation" yaml:"rotation"`
			Scale    map[string]float64 `json:"scale" yaml:"scale"`
		} `json:"default_transforms" yaml:"default_transforms"`
		CameraOffset         map[string]float64 `json:"camera_offset" yaml:"camera_offset"`
		CameraFollowDistance float64            `json:"camera_follow_distance" yaml:"camera_follow_distance"`
		Physics              struct {
			Mass            float64 `json:"mass" yaml:"mass"`
			Height          float64 `json:"height" yaml:"height"`
			CollisionRadius float64 `json:"collision_radius" yaml:"collision_radius"`
		} `json:"physics" yaml:"physics"`
	} `json:"physical" yaml:"physical"`
	Entity struct {
		Tags         []string               `json:"tags" yaml:"tags"`
		NameTemplate string                 `json:"name_template" yaml:"name_template"`
		Components   map[string]interface{} `json:"components" yaml:"components"`
	} `json:"entity" yaml:"entity"`
	Movement struct {
		WalkSpeed    float64 `json:"walk_speed" yaml:"walk_speed"`
		RunSpeed     float64 `json:"run_speed" yaml:"run_speed"`
		TurnSpeed    float64 `json:"turn_speed" yaml:"turn_speed"`
		Acceleration float64 `json:"acceleration" yaml:"acceleration"`
		Deceleration float64 `json:"deceleration" yaml:"deceleration"`
	} `json:"movement" yaml:"movement"`
	Animations map[string]interface{} `json:"animations" yaml:"animations"`
	Network    struct {
		SyncFrequency     int     `json:"sync_frequency" yaml:"sync_frequency"`
		PositionThreshold float64 `json:"position_threshold" yaml:"position_threshold"`
		RotationThreshold float64 `json:"rotation_threshold" yaml:"rotation_threshold"`
		Interpolation     struct {
			Enabled bool    `json:"enabled" yaml:"enabled"`
			Method  string  `json:"method" yaml:"method"`
			Time    float64 `json:"time" yaml:"time"`
		} `json:"interpolation" yaml:"interpolation"`
	} `json:"network" yaml:"network"`
}

// GetAvatarSpecificationHandler handles GET /avatars/{avatarType}
func GetAvatarSpecificationHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	avatarType := extractAvatarType(r.URL.Path)

	// TRACE: Module-specific debugging for avatar operations
	logging.Trace("avatars", "avatar specification request", map[string]interface{}{
		"avatar_type": avatarType,
	})

	// Dynamically validate avatar type from configuration
	validTypes, err := getValidAvatarTypes()
	if err != nil {
		logging.Error("failed to load avatar configuration", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Avatar configuration error", http.StatusInternalServerError)
		return
	}

	if !validTypes[avatarType] {
		validTypesList := make([]string, 0, len(validTypes))
		for t := range validTypes {
			validTypesList = append(validTypesList, t)
		}
		logging.Warn("invalid avatar type", map[string]interface{}{
			"avatar_type": avatarType,
			"valid_types": validTypesList,
		})
		http.Error(w, "Invalid avatar type", http.StatusBadRequest)
		return
	}

	// Read avatar specification file from configured path
	specPath := filepath.Join(config.GetAvatarsDir(), avatarType, "avatar.yaml")
	specData, err := os.ReadFile(specPath)
	if err != nil {
		logging.Error("avatar specification read failed", map[string]interface{}{
			"avatar_type": avatarType,
			"path":        specPath,
			"error":       err.Error(),
		})
		http.Error(w, "Avatar specification not found", http.StatusNotFound)
		return
	}

	// Parse YAML specification
	var spec AvatarSpecification
	if err := yaml.Unmarshal(specData, &spec); err != nil {
		logging.Error("avatar specification parse failed", map[string]interface{}{
			"avatar_type": avatarType,
			"error":       err.Error(),
		})
		http.Error(w, "Failed to parse avatar specification", http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]interface{}{
		"success": true,
		"avatar":  spec,
		"message": "Avatar specification retrieved successfully",
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
	logging.Info("avatar specification retrieved", map[string]interface{}{
		"avatar_type": avatarType,
	})
}

// extractAvatarType extracts avatar type from URL path
func extractAvatarType(path string) string {
	// Handle full path like /api/avatars/{avatarType}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	
	// Find the index of "avatars" and get the next part
	for i, part := range parts {
		if part == "avatars" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	
	return ""
}


// getValidAvatarTypes dynamically loads valid avatar types from configuration
func getValidAvatarTypes() (map[string]bool, error) {
	configPath := filepath.Join(config.GetAvatarsDir(), "config.yaml")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Use AvatarConfig struct from list.go to avoid duplication
	var avatarConfig struct {
		AvatarTypes map[string]struct {
			Path        string   `yaml:"path"`
			Description string   `yaml:"description"`
			Contexts    []string `yaml:"contexts"`
		} `yaml:"avatar_types"`
	}
	
	if err := yaml.Unmarshal(configData, &avatarConfig); err != nil {
		return nil, err
	}

	validTypes := make(map[string]bool)
	for avatarType := range avatarConfig.AvatarTypes {
		validTypes[avatarType] = true
	}

	return validTypes, nil
}