package entities

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
)

// ListEntitiesHandler handles GET /sessions/{sessionId}/entities
func ListEntitiesHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
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
	// Extract sessionId from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/"), "/")
	if len(pathParts) < 3 || pathParts[0] != "sessions" || pathParts[2] != "entities" {
		logging.Error("invalid URL path for list entities", map[string]interface{}{
			"path": r.URL.Path,
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
	
	// Parse query parameters
	query := r.URL.Query()
	tagFilter := query.Get("tag")
	enabledFilter := query.Get("enabled")
	limitParam := query.Get("limit")
	
	limit := 100 // Default limit
	if limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 && parsedLimit <= 1000 {
			limit = parsedLimit
		}
	}
	
	// Get entities from session storage
	storedEntities, err := h.GetStore().GetEntities(sessionID)
	if err != nil {
		logging.Error("failed to get entities", map[string]interface{}{
			"session_id": sessionID,
			"error": err.Error(),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to retrieve entities",
		})
		return
	}
	
	// Convert to response format
	entities := make([]map[string]interface{}, 0, len(storedEntities))
	for _, entity := range storedEntities {
		// Apply filters
		if tagFilter != "" {
			// Tag filtering not implemented yet
			continue
		}
		
		if enabledFilter != "" {
			if enabledFilter == "true" && !entity.Enabled {
				continue
			}
			if enabledFilter == "false" && entity.Enabled {
				continue
			}
		}
		
		entityData := map[string]interface{}{
			"entity_id":        entity.ID,
			"name":             entity.Name,
			"playcanvas_guid":  entity.PlayCanvasGUID,
			"components":       entity.Components,
			"created_at":       entity.CreatedAt,
			"enabled":          entity.Enabled,
		}
		
		entities = append(entities, entityData)
		
		// Apply limit
		if len(entities) >= limit {
			break
		}
	}
	
	logging.Info("entities listed", map[string]interface{}{
		"session_id":    sessionID,
		"tag_filter":    tagFilter,
		"enabled_filter": enabledFilter,
		"limit":         limit,
		"count":         len(entities),
	})
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"entities":  entities,
		"total":     len(entities),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}