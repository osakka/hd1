package avatars

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
	logger := logging.GetLogger()
	avatarType := extractAvatarType(r.URL.Path)

	logger.Info("Getting avatar specification", map[string]interface{}{
		"raw_path":    r.URL.Path,
		"avatar_type": avatarType,
	})

	// Validate avatar type
	validTypes := map[string]bool{
		"default":  true,
		"business": true,
		"casual":   true,
	}

	if !validTypes[avatarType] {
		logger.Warn("Invalid avatar type requested", map[string]interface{}{
			"avatar_type": avatarType,
		})
		http.Error(w, "Invalid avatar type", http.StatusBadRequest)
		return
	}

	// Read avatar specification file
	specPath := filepath.Join("/opt/hd1/share/avatars", avatarType, "avatar.yaml")
	specData, err := os.ReadFile(specPath)
	if err != nil {
		logger.Error("Failed to read avatar specification", map[string]interface{}{
			"avatar_type": avatarType,
			"path":        specPath,
			"error":       err,
		})
		http.Error(w, "Avatar specification not found", http.StatusNotFound)
		return
	}

	// Parse YAML specification
	var spec AvatarSpecification
	if err := yaml.Unmarshal(specData, &spec); err != nil {
		logger.Error("Failed to parse avatar specification", map[string]interface{}{
			"avatar_type": avatarType,
			"error":       err,
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
		logger.Error("Failed to encode avatar specification response", map[string]interface{}{
			"error": err,
		})
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.Info("Avatar specification retrieved successfully", map[string]interface{}{
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