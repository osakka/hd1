// ===================================================================
// HD1 Entity Lifecycle Destruction Management (PlayCanvas v3.0)
// ===================================================================
//
// API-First Game Engine Platform: Entity destruction operations
// Single source of truth: api.yaml drives all functionality
// Generated: Part of HD1 v3.0 PlayCanvas transformation
//
// ===================================================================
package lifecycle

import (
	"encoding/json"
	"net/http"
	"holodeck1/logging"
	"holodeck1/server"
)

// DELETE /sessions/{sessionId}/entities/{entityId}/lifecycle/destroy
func DestroyEntityHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId, entityId := extractSessionAndEntity(r.URL.Path)
	
	var request struct {
		DestroyChildren bool `json:"destroyChildren"`
		ForceDestroy    bool `json:"forceDestroy"`
	}
	
	// Set defaults
	request.DestroyChildren = false
	request.ForceDestroy = false
	
	if r.Body != nil {
		json.NewDecoder(r.Body).Decode(&request)
	}
	
	logging.Info("entity destruction operation", map[string]interface{}{
		"session_id": sessionId,
		"entity_id": entityId,
		"destroy_children": request.DestroyChildren,
		"force_destroy": request.ForceDestroy,
	})

	// PlayCanvas integration: Entity destruction with dependency checking
	// TODO: Implement actual PlayCanvas entity destruction:
	// - Check for dependencies (children, references)
	// - Clean up all components
	// - Remove from scene graph
	// - Free memory and resources
	// - Handle child destruction if requested
	
	// Mock dependency checking
	hasActiveDependencies := !request.DestroyChildren && hasActiveChildren(entityId)
	
	if hasActiveDependencies && !request.ForceDestroy {
		// Return conflict - cannot destroy entity with dependencies
		dependencies := []string{"child_entity_123", "child_entity_456"}
		
		response := map[string]interface{}{
			"error": "Cannot destroy entity with active children",
			"dependencies": dependencies,
		}
		
		writeJSONResponse(w, response, http.StatusConflict)
		return
	}
	
	// Proceed with destruction
	logging.Info("entity destroyed successfully", map[string]interface{}{
		"session_id": sessionId,
		"entity_id": entityId,
		"children_destroyed": request.DestroyChildren,
		"forced": request.ForceDestroy,
	})
	
	// Return 204 No Content for successful destruction
	w.WriteHeader(http.StatusNoContent)
}

// Helper function to check if entity has active children
func hasActiveChildren(entityId string) bool {
	// Mock implementation - in real PlayCanvas integration, 
	// this would check entity.children.length > 0
	
	// For demo purposes, assume entities ending in "parent" have children
	return len(entityId) > 6 && entityId[len(entityId)-6:] == "parent"
}