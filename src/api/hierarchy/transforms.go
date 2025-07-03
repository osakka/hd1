// ===================================================================
// HD1 Entity Hierarchy Transform Management (PlayCanvas v3.0)
// ===================================================================
//
// API-First Game Engine Platform: Entity transform hierarchy
// Single source of truth: api.yaml drives all functionality
// Generated: Part of HD1 v3.0 PlayCanvas transformation
//
// ===================================================================
package hierarchy

import (
	"encoding/json"
	"net/http"
	"holodeck1/logging"
	"holodeck1/server"
)

// GET /sessions/{sessionId}/entities/{entityId}/hierarchy/transforms
func GetEntityTransformsHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId, entityId := extractSessionAndEntity(r.URL.Path)
	
	logging.Info("entity local transform retrieval", map[string]interface{}{
		"session_id": sessionId,
		"entity_id": entityId,
	})

	// PlayCanvas integration: Get entity local transform
	// TODO: Implement actual PlayCanvas entity.getLocalPosition(), getLocalRotation(), getLocalScale()
	
	transform := map[string]interface{}{
		"position": []float64{0.0, 0.0, 0.0},
		"rotation": map[string]interface{}{
			"x": 0.0, "y": 0.0, "z": 0.0, "w": 1.0,
		},
		"scale": []float64{1.0, 1.0, 1.0},
	}
	
	response := map[string]interface{}{
		"success": true,
		"transform": transform,
		"space": "local",
		"message": "Local transform retrieved successfully",
	}
	
	writeJSONResponse(w, response, http.StatusOK)
}

// PUT /sessions/{sessionId}/entities/{entityId}/hierarchy/transforms  
func SetEntityTransformsHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId, entityId := extractSessionAndEntity(r.URL.Path)
	
	var request struct {
		Position *[]float64 `json:"position,omitempty"`
		Rotation *struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
			Z float64 `json:"z"`
			W float64 `json:"w"`
		} `json:"rotation,omitempty"`
		Scale *[]float64 `json:"scale,omitempty"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	logging.Info("entity local transform update", map[string]interface{}{
		"session_id": sessionId,
		"entity_id": entityId,
		"has_position": request.Position != nil,
		"has_rotation": request.Rotation != nil,
		"has_scale": request.Scale != nil,
	})

	// PlayCanvas integration: Update entity local transform
	// TODO: Implement actual PlayCanvas entity.setLocalPosition(), setLocalRotation(), setLocalScale()
	
	// Build response transform from request
	transform := map[string]interface{}{}
	if request.Position != nil {
		transform["position"] = request.Position
	}
	if request.Rotation != nil {
		transform["rotation"] = request.Rotation
	}
	if request.Scale != nil {
		transform["scale"] = request.Scale
	}
	
	response := map[string]interface{}{
		"success": true,
		"transform": transform,
		"space": "local",
		"message": "Local transform updated successfully",
	}
	
	writeJSONResponse(w, response, http.StatusOK)
}

// GET /sessions/{sessionId}/entities/{entityId}/hierarchy/transforms/world
func GetEntityWorldTransformHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId, entityId := extractSessionAndEntity(r.URL.Path)
	
	logging.Info("entity world transform retrieval", map[string]interface{}{
		"session_id": sessionId,
		"entity_id": entityId,
	})

	// PlayCanvas integration: Get entity world transform
	// TODO: Implement actual PlayCanvas entity.getPosition(), getRotation(), getLocalScale() (world)
	
	transform := map[string]interface{}{
		"position": []float64{0.0, 0.0, 0.0},
		"rotation": map[string]interface{}{
			"x": 0.0, "y": 0.0, "z": 0.0, "w": 1.0,
		},
		"scale": []float64{1.0, 1.0, 1.0},
	}
	
	response := map[string]interface{}{
		"success": true,
		"transform": transform,
		"space": "world",
		"message": "World transform retrieved successfully",
	}
	
	writeJSONResponse(w, response, http.StatusOK)
}

// PUT /sessions/{sessionId}/entities/{entityId}/hierarchy/transforms/world
func SetEntityWorldTransformHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId, entityId := extractSessionAndEntity(r.URL.Path)
	
	var request struct {
		Position *[]float64 `json:"position,omitempty"`
		Rotation *struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
			Z float64 `json:"z"`
			W float64 `json:"w"`
		} `json:"rotation,omitempty"`
		Scale *[]float64 `json:"scale,omitempty"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	logging.Info("entity world transform update", map[string]interface{}{
		"session_id": sessionId,
		"entity_id": entityId,
		"has_position": request.Position != nil,
		"has_rotation": request.Rotation != nil,
		"has_scale": request.Scale != nil,
	})

	// PlayCanvas integration: Update entity world transform
	// TODO: Implement actual PlayCanvas entity.setPosition(), setRotation(), setLocalScale() (world)
	
	// Build response transform from request
	transform := map[string]interface{}{}
	if request.Position != nil {
		transform["position"] = request.Position
	}
	if request.Rotation != nil {
		transform["rotation"] = request.Rotation
	}
	if request.Scale != nil {
		transform["scale"] = request.Scale
	}
	
	response := map[string]interface{}{
		"success": true,
		"transform": transform,
		"space": "world",
		"message": "World transform updated successfully",
	}
	
	writeJSONResponse(w, response, http.StatusOK)
}