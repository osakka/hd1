package avatars

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"holodeck1/config"
	"holodeck1/logging"
)

// GetAvatarAssetHandler serves the GLB asset file for a specific avatar type via HTTP
// This endpoint provides proper HTTP delivery of 3D model assets for PlayCanvas
func GetAvatarAssetHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	avatarType := extractAvatarType(r.URL.Path)

	if avatarType == "" {
		http.Error(w, "Avatar type is required", http.StatusBadRequest)
		return
	}

	logging.Info("avatar asset request", map[string]interface{}{
		"avatar_type": avatarType,
		"method":      "HTTP",
	})

	// Validate avatar type
	validTypes, err := getValidAvatarTypes()
	if err != nil {
		logging.Error("failed to load avatar configuration", map[string]interface{}{
			"avatar_type": avatarType,
			"error":       err.Error(),
		})
		http.Error(w, "Failed to load avatar configuration", http.StatusInternalServerError)
		return
	}

	if !validTypes[avatarType] {
		logging.Warn("invalid avatar type requested", map[string]interface{}{
			"avatar_type":   avatarType,
			"valid_types":   getMapKeys(validTypes),
		})
		http.Error(w, "Invalid avatar type", http.StatusNotFound)
		return
	}

	// Construct path to GLB asset file
	avatarsDir := config.GetAvatarsDir()
	glbPath := filepath.Join(avatarsDir, avatarType, "model.glb")

	// Read GLB file
	glbData, err := ioutil.ReadFile(glbPath)
	if err != nil {
		logging.Error("failed to read GLB asset file", map[string]interface{}{
			"avatar_type": avatarType,
			"glb_path":    glbPath,
			"error":       err.Error(),
		})
		http.Error(w, "Avatar asset not found", http.StatusNotFound)
		return
	}

	// Set proper headers for GLB binary data
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(glbData)))
	w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour
	
	// Allow CORS for PlayCanvas asset loading
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Write GLB binary data
	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write(glbData)
	if writeErr != nil {
		logging.Error("failed to write GLB asset response", map[string]interface{}{
			"avatar_type": avatarType,
			"error":       writeErr.Error(),
		})
		return
	}

	logging.Info("GLB asset served via HTTP", map[string]interface{}{
		"avatar_type": avatarType,
		"size_bytes": len(glbData),
		"path":       glbPath,
	})
}

// Helper function to get map keys
func getMapKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}