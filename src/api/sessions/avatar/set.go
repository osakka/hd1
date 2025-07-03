package avatar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
	"gopkg.in/yaml.v3"
)

// SetAvatarRequest represents the request body for setting session avatar
type SetAvatarRequest struct {
	AvatarType    string             `json:"avatar_type"`
	SpawnPosition map[string]float64 `json:"spawn_position,omitempty"`
}

// AvatarSpec represents a simplified avatar specification
type AvatarSpec struct {
	Metadata struct {
		Name string `yaml:"name"`
		Type string `yaml:"type"`
	} `yaml:"metadata"`
	Physical struct {
		DefaultTransforms struct {
			Position map[string]float64 `yaml:"position"`
			Rotation map[string]float64 `yaml:"rotation"`
			Scale    map[string]float64 `yaml:"scale"`
		} `yaml:"default_transforms"`
		CameraOffset map[string]float64 `yaml:"camera_offset"`
	} `yaml:"physical"`
	Entity struct {
		Tags         []string               `yaml:"tags"`
		NameTemplate string                 `yaml:"name_template"`
		Components   map[string]interface{} `yaml:"components"`
	} `yaml:"entity"`
	Movement struct {
		WalkSpeed float64 `yaml:"walk_speed"`
		RunSpeed  float64 `yaml:"run_speed"`
	} `yaml:"movement"`
}

// SetSessionAvatarHandler handles POST /sessions/{sessionId}/avatar
func SetSessionAvatarHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	logger := logging.GetLogger()
	sessionID := extractSessionID(r.URL.Path)

	logger.Info("Setting session avatar", map[string]interface{}{
		"session_id": sessionID,
	})

	// Parse request body
	var req SetAvatarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to parse request body", map[string]interface{}{
			"error": err,
		})
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate avatar type
	validTypes := map[string]bool{
		"default":  true,
		"business": true,
		"casual":   true,
	}

	if !validTypes[req.AvatarType] {
		logger.Warn("Invalid avatar type requested", map[string]interface{}{
			"avatar_type": req.AvatarType,
		})
		http.Error(w, "Invalid avatar type", http.StatusBadRequest)
		return
	}

	// Cast hub to proper type
	h, ok := hub.(*server.Hub)
	if !ok {
		logger.Error("Hub not available")
		http.Error(w, "Server not available", http.StatusInternalServerError)
		return
	}

	// Get session
	session, exists := h.GetStore().GetSession(sessionID)
	if !exists {
		logger.Warn("Session not found", map[string]interface{}{
			"session_id": sessionID,
		})
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Load avatar specification
	specPath := filepath.Join("/opt/hd1/share/avatars", req.AvatarType, "avatar.yaml")
	specData, err := os.ReadFile(specPath)
	if err != nil {
		logger.Error("Failed to read avatar specification", map[string]interface{}{
			"avatar_type": req.AvatarType,
			"path":        specPath,
			"error":       err,
		})
		http.Error(w, "Avatar specification not found", http.StatusNotFound)
		return
	}

	var spec AvatarSpec
	if err := yaml.Unmarshal(specData, &spec); err != nil {
		logger.Error("Failed to parse avatar specification", map[string]interface{}{
			"avatar_type": req.AvatarType,
			"error":       err,
		})
		http.Error(w, "Failed to parse avatar specification", http.StatusInternalServerError)
		return
	}

	// Remove existing avatar entity if it exists
	entities, err := h.GetStore().GetEntities(sessionID)
	if err == nil {
		for _, entity := range entities {
			for _, tag := range entity.Tags {
				if tag == "session-avatar" || tag == "avatar" {
					logger.Info("Removing existing avatar entity", map[string]interface{}{
						"entity_id": entity.ID,
					})
					// Remove from session entities
					delete(session.Entities, entity.ID)
					break
				}
			}
		}
	}

	// Generate entity name from template
	entityName := strings.ReplaceAll(spec.Entity.NameTemplate, "{session_id}", sessionID)
	entityName = strings.ReplaceAll(entityName, "{user_id}", "user") // TODO: Get actual user ID

	// Determine spawn position
	spawnPos := spec.Physical.DefaultTransforms.Position
	if req.SpawnPosition != nil {
		spawnPos = req.SpawnPosition
	}

	// Create avatar entity
	entityID := fmt.Sprintf("avatar-%s-%d", sessionID, time.Now().Unix())
	
	// Build entity components from spec
	components := make(map[string]interface{})
	
	// Transform component
	components["transform"] = map[string]interface{}{
		"position": spawnPos,
		"rotation": spec.Physical.DefaultTransforms.Rotation,
		"scale":    spec.Physical.DefaultTransforms.Scale,
	}

	// Copy other components from spec
	for componentType, componentData := range spec.Entity.Components {
		if componentType != "transform" {
			components[componentType] = componentData
		}
	}

	// Create entity with avatar tags
	tags := append(spec.Entity.Tags, "avatar-"+req.AvatarType)
	
	// Create entity directly
	entity := &server.Entity{
		ID:             entityID,
		Name:           entityName,
		PlayCanvasGUID: "", // Will be set by PlayCanvas
		Components:     components,
		Tags:           tags,
		CreatedAt:      time.Now(),
		Enabled:        true,
	}
	if err := h.GetStore().AddEntity(sessionID, entity); err != nil {
		logger.Error("Failed to create avatar entity", map[string]interface{}{
			"error": err,
		})
		http.Error(w, "Failed to create avatar entity", http.StatusInternalServerError)
		return
	}

	// Broadcast avatar creation via WebSocket
	broadcastData := map[string]interface{}{
		"type":       "avatar_created",
		"session_id": sessionID,
		"entity_id":  entityID,
		"avatar_type": req.AvatarType,
		"position":   spawnPos,
		"timestamp":  time.Now(),
	}

	h.BroadcastToSession(sessionID, "avatar_created", broadcastData)

	// Build session avatar response
	sessionAvatar := SessionAvatar{
		AvatarType: req.AvatarType,
		EntityID:   entityID,
		Created:    time.Now(),
		Position:   spawnPos,
		Rotation:   spec.Physical.DefaultTransforms.Rotation,
		Status:     "active",
		LastUpdate: time.Now(),
		SessionID:  sessionID,
	}

	response := map[string]interface{}{
		"success":   true,
		"avatar":    sessionAvatar,
		"entity_id": entityID,
		"message":   "Avatar set successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Failed to encode set avatar response", map[string]interface{}{
			"error": err,
		})
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.Info("Session avatar set successfully", map[string]interface{}{
		"session_id":  sessionID,
		"avatar_type": req.AvatarType,
		"entity_id":   entityID,
	})
}