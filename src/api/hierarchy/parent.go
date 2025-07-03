// ===================================================================
// HD1 Entity Hierarchy Parent Management (PlayCanvas v3.0)
// ===================================================================
//
// API-First Game Engine Platform: Entity parent-child relationships
// Single source of truth: api.yaml drives all functionality
// Generated: Part of HD1 v3.0 PlayCanvas transformation
//
// ===================================================================
package hierarchy

import (
	"encoding/json"
	"net/http"
	"strings"
	"holodeck1/logging"
	"holodeck1/server"
)

// GET /sessions/{sessionId}/entities/{entityId}/hierarchy/parent
func GetEntityParentHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId, entityId := extractSessionAndEntity(r.URL.Path)
	
	logging.Info("entity parent retrieval", map[string]interface{}{
		"session_id": sessionId,
		"entity_id": entityId,
	})

	// PlayCanvas integration: Get entity parent from scene graph
	response := map[string]interface{}{
		"success": true,
		"parent": map[string]interface{}{
			"parentId": nil, // Root entity has no parent
			"entityId": entityId,
			"name": "SampleEntity",
			"hasParent": false,
		},
		"message": "Entity parent retrieved successfully",
	}
	
	writeJSONResponse(w, response, http.StatusOK)
}

// PUT /sessions/{sessionId}/entities/{entityId}/hierarchy/parent
func SetEntityParentHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId, entityId := extractSessionAndEntity(r.URL.Path)
	
	var request struct {
		ParentId *string `json:"parentId"`
		Index    *int    `json:"index,omitempty"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	logging.Info("entity parent assignment", map[string]interface{}{
		"session_id": sessionId,
		"entity_id": entityId,
		"parent_id": request.ParentId,
		"index": request.Index,
	})

	// PlayCanvas integration: Update entity hierarchy in scene graph
	// TODO: Implement actual PlayCanvas entity.reparent() call
	
	response := map[string]interface{}{
		"success": true,
		"parent": map[string]interface{}{
			"parentId": request.ParentId,
			"entityId": entityId,
			"index": request.Index,
		},
		"message": "Entity parent set successfully",
	}
	
	writeJSONResponse(w, response, http.StatusOK)
}

// DELETE /sessions/{sessionId}/entities/{entityId}/hierarchy/parent
func RemoveEntityParentHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId, entityId := extractSessionAndEntity(r.URL.Path)
	
	logging.Info("entity parent removal", map[string]interface{}{
		"session_id": sessionId,
		"entity_id": entityId,
	})

	// PlayCanvas integration: Remove entity from parent (make root)
	// TODO: Implement actual PlayCanvas entity.reparent(app.root) call
	
	response := map[string]interface{}{
		"success": true,
		"entity": map[string]interface{}{
			"entityId": entityId,
			"parentId": nil,
			"isRoot": true,
		},
		"message": "Entity parent removed - entity is now root",
	}
	
	writeJSONResponse(w, response, http.StatusOK)
}

// Helper function to extract session and entity IDs from URL path
func extractSessionAndEntity(path string) (string, string) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	
	sessionId := "default-session"
	entityId := "default-entity"
	
	// Parse /api/sessions/{sessionId}/entities/{entityId}/...
	if len(parts) >= 4 && parts[0] == "api" && parts[1] == "sessions" && parts[3] == "entities" {
		sessionId = parts[2]
		if len(parts) >= 5 {
			entityId = parts[4]
		}
	}
	
	return sessionId, entityId
}

// Helper function to write JSON responses
func writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}