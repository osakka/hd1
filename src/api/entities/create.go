// Package entities provides HTTP handlers for entity lifecycle management
// in HD1 sessions. Entities represent 3D objects with components that can
// be manipulated via the REST API and synchronized via WebSocket.
//
// Key concepts:
//   - Entity: PlayCanvas 3D object with transform, model, material components
//   - Session: Isolated game world containing entities
//   - WebSocket sync: Real-time entity updates across clients
//   - PlayCanvas GUID: Unique identifier for PlayCanvas engine integration
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

// mustMarshal converts any value to JSON bytes, ignoring errors.
// Used for internal WebSocket message marshaling where errors are not expected.
// Panics should not occur as we control the input data structures.
func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

// Transform represents PlayCanvas entity transform component data.
// All vectors use PlayCanvas coordinate system: X=right, Y=up, Z=forward.
type Transform struct {
	Position []float64 `json:"position,omitempty"` // [x, y, z] world position
	Rotation []float64 `json:"rotation,omitempty"` // [x, y, z] Euler angles in degrees
	Scale    []float64 `json:"scale,omitempty"`    // [x, y, z] scale factors
}

// CreateEntityRequest represents the request body for entity creation.
// Supports both legacy position fields and modern PlayCanvas component structure.
// All fields except Name are optional and will use PlayCanvas defaults.
type CreateEntityRequest struct {
	Name       string                 `json:"name"`                   // Required: entity display name (1-64 chars)
	Tags       []string               `json:"tags,omitempty"`         // Optional: entity tags for organization
	Enabled    *bool                  `json:"enabled,omitempty"`      // Optional: entity active state (default: true)
	Position   []float64              `json:"position,omitempty"`     // Legacy: [x,y,z] position (use components.transform instead)
	Rotation   []float64              `json:"rotation,omitempty"`     // Legacy: [x,y,z] rotation (use components.transform instead)
	Scale      []float64              `json:"scale,omitempty"`        // Legacy: [x,y,z] scale (use components.transform instead)
	Components map[string]interface{} `json:"components,omitempty"`   // PlayCanvas components (transform, model, material, etc.)
}

// CreateEntityHandler handles POST /sessions/{sessionId}/entities requests.
// Creates a new PlayCanvas entity in the specified session with proper
// component structure and WebSocket broadcasting.
//
// The request body should contain entity name and optional components.
// Supports both legacy position fields and modern PlayCanvas component structure.
// Returns 201 Created with entity details on success, or appropriate
// error status codes for validation failures.
//
// WebSocket event: Broadcasts 'entity_created' to all session clients
// with full entity data for immediate client rendering.
//
// URL path: /api/sessions/{sessionId}/entities
// Method: POST
// Content-Type: application/json
func CreateEntityHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Hub interface validation: Ensure we have a valid server hub
	// This cast should never fail in normal operation
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
	// URL parsing: Extract session ID from REST path pattern
	// Expected format: /api/sessions/{sessionId}/entities
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/"), "/")
	if len(pathParts) < 3 || pathParts[0] != "sessions" || pathParts[2] != "entities" {
		logging.Error("invalid URL path format", map[string]interface{}{
			"path":         r.URL.Path,
			"method":       r.Method,
			"expected":     "/api/sessions/{sessionId}/entities", 
			"parts_count":  len(pathParts),
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
	
	// Session validation: Verify session exists and is active
	// Prevents entity creation in non-existent sessions
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
	
	// Request parsing: Decode JSON request body to struct
	// Validates JSON format and structure
	var req CreateEntityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.Error("request body decode failed", map[string]interface{}{
			"session_id": sessionID,
			"path":       r.URL.Path,
			"method":     r.Method,
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
	
	// Field validation: Check required fields are present
	// Name is the only required field for entity creation
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
	
	// Position validation: Legacy position field must be 3D vector
	// API-first design: strict validation of coordinate arrays
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
	
	// Default values: Apply sensible defaults for optional fields
	// Enabled defaults to true for immediate visibility
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	
	// ID generation: Create unique identifiers for entity tracking
	// PlayCanvas GUID enables integration with PlayCanvas engine
	entityID := fmt.Sprintf("entity-%s-%d", sessionID[8:], time.Now().UnixNano()%1000000)
	playCanvasGUID := fmt.Sprintf("pc-guid-%d", time.Now().UnixNano())
	
	// Component processing: Build PlayCanvas component structure
	// Supports both modern components and legacy position fields
	components := make(map[string]interface{})
	
	// Component assignment: Use provided components or empty map
	if req.Components != nil {
		components = req.Components
	}
	
	// Legacy compatibility: Convert position field to transform component
	// Maintains backward compatibility with v4.x API clients
	if req.Position != nil && len(req.Position) == 3 {
		if _, hasTransform := components["transform"]; !hasTransform {
			components["transform"] = map[string]interface{}{
				"position": req.Position,
			}
		}
	}
	
	// Entity construction: Build complete entity with all metadata
	entity := &server.Entity{
		ID:             entityID,
		Name:           req.Name,
		PlayCanvasGUID: playCanvasGUID,
		Components:     components,
		Tags:           req.Tags,
		CreatedAt:      time.Now(),
		Enabled:        enabled,
	}
	
	// Entity persistence: Store entity in session state for retrieval
	// Thread-safe session store handles concurrent access
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
		"playcanvas_guid":  playCanvasGUID,
	})
	
	// WebSocket notification: Broadcast entity creation to all session clients
	// Includes full entity data for immediate PlayCanvas rendering
	h.BroadcastUpdate("entity_created", map[string]interface{}{
		"session_id":       sessionID,
		"entity_id":        entityID,
		"name":             req.Name,
		"enabled":          enabled,
		"tags":             req.Tags,
		"playcanvas_guid":  playCanvasGUID,
		"components":       req.Components, // CRITICAL: Include components for client rendering
	})
	
	// Success response: Return created entity details
	// 201 Created with entity metadata for client reference
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":          true,
		"entity_id":        entityID,
		"name":             req.Name,
		"playcanvas_guid":  playCanvasGUID,
		"created_at":       time.Now().Format(time.RFC3339),
	})
}