// ===================================================================
// HD1 Entity Hierarchy Children Management (PlayCanvas v3.0)
// ===================================================================
//
// API-First Game Engine Platform: Entity children management
// Single source of truth: api.yaml drives all functionality
// Generated: Part of HD1 v3.0 PlayCanvas transformation
//
// ===================================================================
package hierarchy

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"holodeck1/logging"
	"holodeck1/server"
)

// GET /sessions/{sessionId}/entities/{entityId}/hierarchy/children
func GetEntityChildrenHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId, entityId := extractSessionAndEntity(r.URL.Path)
	
	// Query parameters
	recursive := r.URL.Query().Get("recursive") == "true"
	depthParam := r.URL.Query().Get("depth")
	
	depth := -1 // unlimited by default
	if depthParam != "" {
		if d, err := strconv.Atoi(depthParam); err == nil {
			depth = d
		}
	}
	
	logging.Info("entity children retrieval", map[string]interface{}{
		"session_id": sessionId,
		"entity_id": entityId,
		"recursive": recursive,
		"depth": depth,
	})

	// PlayCanvas integration: Get entity children from scene graph
	// TODO: Implement actual PlayCanvas entity.children traversal
	
	children := []map[string]interface{}{
		{
			"entityId": "child_001",
			"name": "ChildEntity1",
			"transform": map[string]interface{}{
				"position": []float64{1.0, 0.0, 0.0},
				"rotation": map[string]interface{}{
					"x": 0.0, "y": 0.0, "z": 0.0, "w": 1.0,
				},
				"scale": []float64{1.0, 1.0, 1.0},
			},
		},
	}
	
	response := map[string]interface{}{
		"success": true,
		"children": children,
		"totalChildren": len(children),
		"recursive": recursive,
		"depth": depth,
		"message": "Entity children retrieved successfully",
	}
	
	writeJSONResponse(w, response, http.StatusOK)
}

// POST /sessions/{sessionId}/entities/{entityId}/hierarchy/children/{childId}
func AddEntityChildHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId, entityId := extractSessionAndEntity(r.URL.Path)
	
	// Extract childId from URL
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	childId := "unknown-child"
	if len(pathParts) >= 8 {
		childId = pathParts[7] // /api/sessions/{id}/entities/{id}/hierarchy/children/{childId}
	}
	
	var request struct {
		Index *int `json:"index,omitempty"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	logging.Info("entity child addition", map[string]interface{}{
		"session_id": sessionId,
		"entity_id": entityId,
		"child_id": childId,
		"index": request.Index,
	})

	// PlayCanvas integration: Add child to entity in scene graph
	// TODO: Implement actual PlayCanvas childEntity.reparent(parentEntity) call
	
	response := map[string]interface{}{
		"success": true,
		"relationship": map[string]interface{}{
			"parentId": entityId,
			"childId": childId,
			"index": request.Index,
		},
		"message": "Child entity added successfully",
	}
	
	writeJSONResponse(w, response, http.StatusCreated)
}

// DELETE /sessions/{sessionId}/entities/{entityId}/hierarchy/children/{childId}
func RemoveEntityChildHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId, entityId := extractSessionAndEntity(r.URL.Path)
	
	// Extract childId from URL
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	childId := "unknown-child"
	if len(pathParts) >= 8 {
		childId = pathParts[7] // /api/sessions/{id}/entities/{id}/hierarchy/children/{childId}
	}
	
	logging.Info("entity child removal", map[string]interface{}{
		"session_id": sessionId,
		"entity_id": entityId,
		"child_id": childId,
	})

	// PlayCanvas integration: Remove child from entity (make root)
	// TODO: Implement actual PlayCanvas childEntity.reparent(app.root) call
	
	response := map[string]interface{}{
		"success": true,
		"removed": map[string]interface{}{
			"parentId": entityId,
			"childId": childId,
			"newParent": nil, // Child becomes root
		},
		"message": "Child entity removed from parent",
	}
	
	writeJSONResponse(w, response, http.StatusOK)
}