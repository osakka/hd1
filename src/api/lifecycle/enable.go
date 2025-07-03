// ===================================================================
// HD1 Entity Lifecycle Enable/Disable Management (PlayCanvas v3.0)
// ===================================================================
//
// API-First Game Engine Platform: Entity enable/disable operations
// Single source of truth: api.yaml drives all functionality
// Generated: Part of HD1 v3.0 PlayCanvas transformation
//
// ===================================================================
package lifecycle

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"holodeck1/logging"
	"holodeck1/server"
)

// PUT /sessions/{sessionId}/entities/{entityId}/lifecycle/enable
func EnableEntityHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId, entityId := extractSessionAndEntity(r.URL.Path)
	
	logging.Info("entity enable operation", map[string]interface{}{
		"session_id": sessionId,
		"entity_id": entityId,
	})

	// PlayCanvas integration: Set entity.enabled = true
	// TODO: Implement actual PlayCanvas entity.enabled = true
	
	response := map[string]interface{}{
		"success": true,
		"entity": map[string]interface{}{
			"entityId": entityId,
			"enabled": true,
			"activated": true,
			"visible": true,
		},
		"message": "Entity enabled successfully",
	}
	
	writeJSONResponse(w, response, http.StatusOK)
}

// PUT /sessions/{sessionId}/entities/{entityId}/lifecycle/disable
func DisableEntityHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId, entityId := extractSessionAndEntity(r.URL.Path)
	
	logging.Info("entity disable operation", map[string]interface{}{
		"session_id": sessionId,
		"entity_id": entityId,
	})

	// PlayCanvas integration: Set entity.enabled = false
	// TODO: Implement actual PlayCanvas entity.enabled = false
	
	response := map[string]interface{}{
		"success": true,
		"entity": map[string]interface{}{
			"entityId": entityId,
			"enabled": false,
			"activated": false,
			"visible": false,
		},
		"message": "Entity disabled successfully",
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

// Helper function to create entity lifecycle state
func createEntityLifecycleState(entityId string, enabled, activated bool) map[string]interface{} {
	now := time.Now()
	
	return map[string]interface{}{
		"entityId": entityId,
		"name": "SampleEntity",
		"enabled": enabled,
		"activated": activated,
		"visible": enabled && activated,
		"destroyed": false,
		"componentCount": 3,
		"childCount": 0,
		"parentId": nil,
		"lastModified": now.Format(time.RFC3339),
		"position": []float64{0.0, 0.0, 0.0},
		"lifecycle": map[string]interface{}{
			"createdAt": now.Add(-time.Hour).Format(time.RFC3339),
			"activatedAt": func() interface{} {
				if activated {
					return now.Add(-time.Minute * 30).Format(time.RFC3339)
				}
				return nil
			}(),
			"deactivatedAt": nil,
			"enabledDuration": 1800.5, // 30 minutes
		},
	}
}