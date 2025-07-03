package components

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
)

// GetComponentHandler handles GET /sessions/{sessionId}/entities/{entityId}/components/{componentType}
func GetComponentHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
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
	
	// Extract sessionId, entityId, and componentType from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/"), "/")
	if len(pathParts) < 6 || pathParts[0] != "sessions" || pathParts[2] != "entities" || pathParts[4] != "components" {
		logging.Error("invalid URL path for get component", map[string]interface{}{
			"path": r.URL.Path,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid URL path",
			"message": "Expected /api/sessions/{sessionId}/entities/{entityId}/components/{componentType}",
		})
		return
	}
	
	sessionID := pathParts[1]
	entityID := pathParts[3]
	componentType := pathParts[5]
	
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
	
	// Validate component type
	validTypes := []string{"model", "camera", "light", "rigidbody", "script", "sound", "animation", "collision", "element", "particlesystem", "render", "sprite"}
	isValidType := false
	for _, validType := range validTypes {
		if componentType == validType {
			isValidType = true
			break
		}
	}
	
	if !isValidType {
		logging.Error("invalid component type", map[string]interface{}{
			"component_type": componentType,
			"valid_types":    validTypes,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid component type",
			"message": "Component type must be one of: " + strings.Join(validTypes, ", "),
		})
		return
	}
	
	// TODO: Implement PlayCanvas component retrieval
	// For now, return mock component details
	componentID := "comp-" + componentType + "-" + entityID + "-1"
	
	// Mock component details based on type
	var properties map[string]interface{}
	var schema map[string]interface{}
	var dependencies []string
	
	switch componentType {
	case "model":
		properties = map[string]interface{}{
			"type":           "box",
			"castShadows":    true,
			"receiveShadows": true,
			"material":       "material-asset-id",
		}
		schema = map[string]interface{}{
			"type": map[string]interface{}{
				"type":    "string",
				"enum":    []string{"asset", "box", "sphere", "cylinder", "cone", "capsule", "plane"},
				"default": "box",
			},
			"castShadows": map[string]interface{}{
				"type":    "boolean",
				"default": true,
			},
		}
		dependencies = []string{"transform"}
	case "camera":
		properties = map[string]interface{}{
			"fov":        45.0,
			"nearClip":   0.1,
			"farClip":    1000.0,
			"projection": "perspective",
		}
		schema = map[string]interface{}{
			"fov": map[string]interface{}{
				"type":    "number",
				"minimum": 1.0,
				"maximum": 179.0,
				"default": 45.0,
			},
		}
		dependencies = []string{"transform"}
	case "light":
		properties = map[string]interface{}{
			"type":      "directional",
			"color":     "#ffffff",
			"intensity": 1.0,
		}
		schema = map[string]interface{}{
			"type": map[string]interface{}{
				"type":    "string",
				"enum":    []string{"directional", "point", "spot", "ambient"},
				"default": "directional",
			},
		}
		dependencies = []string{"transform"}
	default:
		properties = map[string]interface{}{}
		schema = map[string]interface{}{}
		dependencies = []string{}
	}
	
	logging.Info("component details retrieved", map[string]interface{}{
		"session_id":     sessionID,
		"entity_id":      entityID,
		"component_id":   componentID,
		"component_type": componentType,
	})
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"component_id":     componentID,
		"type":            componentType,
		"enabled":         true,
		"playcanvas_id":   "pc-comp-guid-123",
		"properties_count": len(properties),
		"last_updated":    time.Now().Format(time.RFC3339),
		"properties":      properties,
		"schema":          schema,
		"events":          []map[string]interface{}{},
		"dependencies":    dependencies,
	})
}