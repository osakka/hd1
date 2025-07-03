// ===================================================================
// HD1 Entity Lifecycle Bulk Operations (PlayCanvas v3.0)
// ===================================================================
//
// API-First Game Engine Platform: Bulk entity lifecycle operations
// Single source of truth: api.yaml drives all functionality
// Generated: Part of HD1 v3.0 PlayCanvas transformation
//
// ===================================================================
package lifecycle

import (
	"encoding/json"
	"net/http"
	"strings"
	"holodeck1/logging"
	"holodeck1/server"
)

// POST /sessions/{sessionId}/entities/lifecycle/bulk
func BulkEntityLifecycleOperationHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionFromPath(r.URL.Path)
	
	var request struct {
		Operation string   `json:"operation"`
		Entities  []string `json:"entities"`
		Options   struct {
			IncludeChildren bool `json:"includeChildren"`
			ForceOperation  bool `json:"forceOperation"`
			StopOnError     bool `json:"stopOnError"`
		} `json:"options"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// Validate operation
	validOperations := []string{"enable", "disable", "activate", "deactivate", "destroy"}
	isValidOperation := false
	for _, op := range validOperations {
		if request.Operation == op {
			isValidOperation = true
			break
		}
	}
	
	if !isValidOperation {
		http.Error(w, "Invalid operation. Must be one of: enable, disable, activate, deactivate, destroy", http.StatusBadRequest)
		return
	}
	
	if len(request.Entities) == 0 {
		http.Error(w, "At least one entity ID is required", http.StatusBadRequest)
		return
	}
	
	logging.Info("bulk entity lifecycle operation", map[string]interface{}{
		"session_id": sessionId,
		"operation": request.Operation,
		"entity_count": len(request.Entities),
		"include_children": request.Options.IncludeChildren,
		"force_operation": request.Options.ForceOperation,
		"stop_on_error": request.Options.StopOnError,
	})

	// PlayCanvas integration: Perform bulk lifecycle operations
	// TODO: Implement actual PlayCanvas bulk operations:
	// - Process entities in batches for performance
	// - Handle dependencies and conflicts
	// - Provide detailed success/failure reporting
	// - Support atomic operations where possible
	
	results := []map[string]interface{}{}
	successful := 0
	failed := 0
	
	for _, entityId := range request.Entities {
		result := map[string]interface{}{
			"entityId": entityId,
			"success": true,
			"error": nil,
		}
		
		// Mock operation execution
		success := performBulkOperation(request.Operation, entityId, request.Options)
		
		if success {
			successful++
		} else {
			failed++
			result["success"] = false
			result["error"] = "Mock operation failed for demonstration"
			
			// Stop on error if requested
			if request.Options.StopOnError {
				results = append(results, result)
				break
			}
		}
		
		results = append(results, result)
	}
	
	response := map[string]interface{}{
		"success": true,
		"operation": request.Operation,
		"results": results,
		"summary": map[string]interface{}{
			"totalEntities": len(request.Entities),
			"successful": successful,
			"failed": failed,
		},
	}
	
	writeJSONResponse(w, response, http.StatusOK)
}

// Helper function to perform individual bulk operations
func performBulkOperation(operation, entityId string, options struct {
	IncludeChildren bool `json:"includeChildren"`
	ForceOperation  bool `json:"forceOperation"`
	StopOnError     bool `json:"stopOnError"`
}) bool {
	// Mock implementation - in real PlayCanvas integration, this would:
	// - Call the appropriate entity method based on operation
	// - Handle children if includeChildren is true
	// - Apply force flags for destruction
	// - Return actual success/failure status
	
	// For demo: fail if entity ID contains "fail"
	return !strings.Contains(entityId, "fail")
}

// Helper function to extract session ID from URL path for bulk operations
func extractSessionFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	
	// Parse /api/sessions/{sessionId}/...
	if len(parts) >= 3 && parts[0] == "api" && parts[1] == "sessions" {
		return parts[2]
	}
	
	return "default-session"
}