// ===================================================================
// HD1 Entity Lifecycle Status Management (PlayCanvas v3.0)
// ===================================================================
//
// API-First Game Engine Platform: Entity lifecycle status operations
// Single source of truth: api.yaml drives all functionality
// Generated: Part of HD1 v3.0 PlayCanvas transformation
//
// ===================================================================
package lifecycle

import (
	"net/http"
	"holodeck1/logging"
	"holodeck1/server"
)

// GET /sessions/{sessionId}/entities/{entityId}/lifecycle/status
func GetEntityLifecycleStatusHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId, entityId := extractSessionAndEntity(r.URL.Path)
	
	// Query parameters
	includeComponents := r.URL.Query().Get("includeComponents") == "true"
	includeChildren := r.URL.Query().Get("includeChildren") == "true"
	
	logging.Info("entity lifecycle status retrieval", map[string]interface{}{
		"session_id": sessionId,
		"entity_id": entityId,
		"include_components": includeComponents,
		"include_children": includeChildren,
	})

	// PlayCanvas integration: Get entity lifecycle status
	// TODO: Implement actual PlayCanvas entity status retrieval:
	// - entity.enabled state
	// - component status for each attached component
	// - children status if requested
	// - activation timestamps and duration
	
	// Create base entity status
	entityState := createEntityLifecycleState(entityId, true, true)
	
	response := map[string]interface{}{
		"success": true,
		"entity": entityState,
	}
	
	// Add component information if requested
	if includeComponents {
		components := []map[string]interface{}{
			{
				"type": "model",
				"enabled": true,
				"initialized": true,
			},
			{
				"type": "animation",
				"enabled": true,
				"initialized": true,
			},
			{
				"type": "rigidbody",
				"enabled": false,
				"initialized": false,
			},
		}
		response["components"] = components
	}
	
	// Add children information if requested
	if includeChildren {
		children := []map[string]interface{}{
			createEntityLifecycleState("child_entity_001", true, true),
			createEntityLifecycleState("child_entity_002", false, false),
		}
		response["children"] = children
	}
	
	writeJSONResponse(w, response, http.StatusOK)
}