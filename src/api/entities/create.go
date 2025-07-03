package entities

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
)

// Helper function for JSON marshaling
func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

// Transform represents entity transform data
type Transform struct {
	Position []float64 `json:"position,omitempty"`
	Rotation []float64 `json:"rotation,omitempty"`
	Scale    []float64 `json:"scale,omitempty"`
}

// CreateEntityRequest represents the request body for entity creation
type CreateEntityRequest struct {
	Name       string                 `json:"name"`
	Tags       []string               `json:"tags,omitempty"`
	Enabled    *bool                  `json:"enabled,omitempty"`
	Position   []float64              `json:"position,omitempty"`
	Rotation   []float64              `json:"rotation,omitempty"`
	Scale      []float64              `json:"scale,omitempty"`
	Components map[string]interface{} `json:"components,omitempty"`
}

// CreateEntityHandler handles POST /sessions/{sessionId}/entities
func CreateEntityHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type
	h, ok := hub.(*server.Hub)
	if !ok {
		logging.Error("hub interface cast failed", map[string]interface{}{
			"expected_type": "*server.Hub",
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Internal server error",
			"message": "Hub interface cast failed",
		})
		return
	}
	// Extract sessionId from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/"), "/")
	if len(pathParts) < 3 || pathParts[0] != "sessions" || pathParts[2] != "entities" {
		logging.Error("invalid URL path for create entity", map[string]interface{}{
			"path": r.URL.Path,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid URL path",
			"message": "Expected /api/sessions/{sessionId}/entities",
		})
		return
	}
	
	sessionID := pathParts[1]
	
	// Validate session exists
	if _, exists := h.GetStore().GetSession(sessionID); !exists {
		logging.Warn("session not found", map[string]interface{}{
			"session_id": sessionID,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Session not found",
			"message": "Session does not exist or has expired",
		})
		return
	}
	
	// Parse request body
	var req CreateEntityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.Error("invalid JSON in request", map[string]interface{}{
			"session_id": sessionID,
			"error":      err.Error(),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid JSON",
			"message": "Request body must be valid JSON",
		})
		return
	}
	
	// Validate required fields
	if req.Name == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Missing required fields",
			"message": "name is required",
		})
		return
	}
	
	// Validate position array if provided (API-first: must be [x,y,z])
	if req.Position != nil && len(req.Position) != 3 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid position",
			"message": "Position must be [x, y, z] array",
		})
		return
	}
	
	// Set defaults
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	
	// Generate entity ID and PlayCanvas GUID
	entityID := fmt.Sprintf("entity-%s-%d", sessionID[8:], time.Now().UnixNano()%1000000)
	playcanvasGUID := fmt.Sprintf("pc-guid-%d", time.Now().UnixNano())
	
	// Create entity with components (supporting full PlayCanvas entity structure)
	components := make(map[string]interface{})
	
	// Use components from request if provided
	if req.Components != nil {
		components = req.Components
	}
	
	// Legacy position support - convert to transform component if no transform provided
	if req.Position != nil && len(req.Position) == 3 {
		if _, hasTransform := components["transform"]; !hasTransform {
			components["transform"] = map[string]interface{}{
				"position": req.Position,
			}
		}
	}
	
	// Create entity
	entity := &server.Entity{
		ID:             entityID,
		Name:           req.Name,
		PlayCanvasGUID: playcanvasGUID,
		Components:     components,
		CreatedAt:      time.Now(),
		Enabled:        enabled,
	}
	
	// Store entity in session
	if err := h.GetStore().AddEntity(sessionID, entity); err != nil {
		logging.Error("failed to store entity", map[string]interface{}{
			"session_id": sessionID,
			"entity_id": entityID,
			"error": err.Error(),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to store entity",
		})
		return
	}
	
	logging.Info("entity created", map[string]interface{}{
		"session_id":       sessionID,
		"entity_id":        entityID,
		"name":             req.Name,
		"enabled":          enabled,
		"tags":             req.Tags,
		"playcanvas_guid":  playcanvasGUID,
	})
	
	// Broadcast entity creation via WebSocket with FULL entity data for client rendering
	h.BroadcastUpdate("entity_created", map[string]interface{}{
		"session_id":       sessionID,
		"entity_id":        entityID,
		"name":             req.Name,
		"enabled":          enabled,
		"tags":             req.Tags,
		"playcanvas_guid":  playcanvasGUID,
		"components":       req.Components, // CRITICAL: Include components for client rendering
	})
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":          true,
		"entity_id":        entityID,
		"name":             req.Name,
		"playcanvas_guid":  playcanvasGUID,
		"created_at":       time.Now().Format(time.RFC3339),
	})
}