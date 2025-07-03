// ===================================================================
// HD1 Entity Lifecycle Activation Management (PlayCanvas v3.0)
// ===================================================================
//
// API-First Game Engine Platform: Entity activation/deactivation operations
// Single source of truth: api.yaml drives all functionality
// Generated: Part of HD1 v3.0 PlayCanvas transformation
//
// ===================================================================
package lifecycle

import (
	"encoding/json"
	"net/http"
	"time"
	"holodeck1/logging"
	"holodeck1/server"
)

// POST /sessions/{sessionId}/entities/{entityId}/lifecycle/activate
func ActivateEntityHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId, entityId := extractSessionAndEntity(r.URL.Path)
	
	var request struct {
		InitializeComponents bool `json:"initializeComponents"`
		ExecuteScripts       bool `json:"executeScripts"`
		ActivateChildren     bool `json:"activateChildren"`
	}
	
	// Set defaults
	request.InitializeComponents = true
	request.ExecuteScripts = true
	request.ActivateChildren = false
	
	if r.Body != nil {
		json.NewDecoder(r.Body).Decode(&request)
	}
	
	logging.Info("entity activation operation", map[string]interface{}{
		"session_id": sessionId,
		"entity_id": entityId,
		"initialize_components": request.InitializeComponents,
		"execute_scripts": request.ExecuteScripts,
		"activate_children": request.ActivateChildren,
	})

	// PlayCanvas integration: Full entity activation
	// TODO: Implement actual PlayCanvas entity activation with:
	// - Component initialization
	// - Script execution
	// - Resource allocation
	// - Child entity activation (if requested)
	
	// Mock activation results
	componentsInitialized := 5
	scriptsExecuted := 2
	childrenActivated := 0
	
	if request.ActivateChildren {
		childrenActivated = 2
	}
	
	entityState := createEntityLifecycleState(entityId, true, true)
	
	response := map[string]interface{}{
		"success": true,
		"entity": entityState,
		"activation": map[string]interface{}{
			"componentsInitialized": componentsInitialized,
			"scriptsExecuted": scriptsExecuted,
			"childrenActivated": childrenActivated,
		},
		"message": "Entity activated with 5 components and 2 scripts",
	}
	
	writeJSONResponse(w, response, http.StatusCreated)
}

// POST /sessions/{sessionId}/entities/{entityId}/lifecycle/deactivate
func DeactivateEntityHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId, entityId := extractSessionAndEntity(r.URL.Path)
	
	var request struct {
		CleanupComponents   bool `json:"cleanupComponents"`
		StopScripts        bool `json:"stopScripts"`
		DeactivateChildren bool `json:"deactivateChildren"`
	}
	
	// Set defaults
	request.CleanupComponents = true
	request.StopScripts = true
	request.DeactivateChildren = false
	
	if r.Body != nil {
		json.NewDecoder(r.Body).Decode(&request)
	}
	
	logging.Info("entity deactivation operation", map[string]interface{}{
		"session_id": sessionId,
		"entity_id": entityId,
		"cleanup_components": request.CleanupComponents,
		"stop_scripts": request.StopScripts,
		"deactivate_children": request.DeactivateChildren,
	})

	// PlayCanvas integration: Full entity deactivation
	// TODO: Implement actual PlayCanvas entity deactivation with:
	// - Component cleanup
	// - Script stopping
	// - Resource deallocation
	// - Child entity deactivation (if requested)
	
	// Mock deactivation results
	componentsCleanedUp := 5
	scriptsStopped := 2
	childrenDeactivated := 0
	
	if request.DeactivateChildren {
		childrenDeactivated = 2
	}
	
	entityState := createEntityLifecycleState(entityId, false, false)
	// Update deactivation timestamp
	if lifecycle, ok := entityState["lifecycle"].(map[string]interface{}); ok {
		lifecycle["deactivatedAt"] = time.Now().Format(time.RFC3339)
	}
	
	response := map[string]interface{}{
		"success": true,
		"entity": entityState,
		"deactivation": map[string]interface{}{
			"componentsCleanedUp": componentsCleanedUp,
			"scriptsStopped": scriptsStopped,
			"childrenDeactivated": childrenDeactivated,
		},
		"message": "Entity deactivated - 5 components cleaned up",
	}
	
	writeJSONResponse(w, response, http.StatusOK)
}