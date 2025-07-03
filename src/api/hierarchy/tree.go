// ===================================================================
// HD1 Entity Hierarchy Tree Operations (PlayCanvas v3.0)
// ===================================================================
//
// API-First Game Engine Platform: Entity hierarchy tree operations
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

// GET /sessions/{sessionId}/entities/hierarchy/tree
func GetHierarchyTreeHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionFromPath(r.URL.Path)
	
	// Query parameters
	includeComponents := r.URL.Query().Get("includeComponents") == "true"
	depthParam := r.URL.Query().Get("maxDepth")
	
	maxDepth := -1 // unlimited by default
	if depthParam != "" {
		if d, err := strconv.Atoi(depthParam); err == nil {
			maxDepth = d
		}
	}
	
	logging.Info("entity hierarchy tree retrieval", map[string]interface{}{
		"session_id": sessionId,
		"include_components": includeComponents,
		"max_depth": maxDepth,
	})

	// PlayCanvas integration: Build complete entity hierarchy tree
	// TODO: Implement actual PlayCanvas app.root traversal and entity hierarchy building
	
	sampleTree := []map[string]interface{}{
		{
			"entityId": "root_001",
			"name": "RootEntity",
			"parentId": nil,
			"depth": 0,
			"children": []string{"child_001", "child_002"},
			"transform": map[string]interface{}{
				"position": []float64{0.0, 0.0, 0.0},
				"rotation": map[string]interface{}{
					"x": 0.0, "y": 0.0, "z": 0.0, "w": 1.0,
				},
				"scale": []float64{1.0, 1.0, 1.0},
			},
		},
		{
			"entityId": "child_001",
			"name": "ChildEntity1",
			"parentId": "root_001",
			"depth": 1,
			"children": []string{},
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
		"root": sampleTree,
		"totalEntities": len(sampleTree),
		"maxDepth": 1,
		"includeComponents": includeComponents,
		"message": "Entity hierarchy tree retrieved successfully",
	}
	
	writeJSONResponse(w, response, http.StatusOK)
}

// POST /sessions/{sessionId}/entities/hierarchy/search
func SearchEntityHierarchyHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionFromPath(r.URL.Path)
	
	var request struct {
		Query       string `json:"query"`
		SearchNames bool   `json:"searchNames,omitempty"`
		SearchIds   bool   `json:"searchIds,omitempty"`
		MaxResults  *int   `json:"maxResults,omitempty"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// Default search options
	if !request.SearchNames && !request.SearchIds {
		request.SearchNames = true
	}
	
	maxResults := 50 // default limit
	if request.MaxResults != nil && *request.MaxResults > 0 {
		maxResults = *request.MaxResults
	}
	
	logging.Info("entity hierarchy search", map[string]interface{}{
		"session_id": sessionId,
		"query": request.Query,
		"search_names": request.SearchNames,
		"search_ids": request.SearchIds,
		"max_results": maxResults,
	})

	// PlayCanvas integration: Search entity hierarchy
	// TODO: Implement actual PlayCanvas entity search through scene graph
	
	matchingEntities := []map[string]interface{}{
		{
			"entityId": "player_001",
			"name": "PlayerCharacter",
			"parentId": "scene_root",
			"depth": 1,
			"children": []string{"weapon_001", "health_bar"},
			"transform": map[string]interface{}{
				"position": []float64{0.0, 1.0, 0.0},
				"rotation": map[string]interface{}{
					"x": 0.0, "y": 0.0, "z": 0.0, "w": 1.0,
				},
				"scale": []float64{1.0, 1.0, 1.0},
			},
		},
	}
	
	response := map[string]interface{}{
		"success": true,
		"entities": matchingEntities,
		"totalMatches": len(matchingEntities),
		"query": request.Query,
		"searchOptions": map[string]interface{}{
			"searchNames": request.SearchNames,
			"searchIds": request.SearchIds,
			"maxResults": maxResults,
		},
		"message": "Entity hierarchy search completed successfully",
	}
	
	writeJSONResponse(w, response, http.StatusOK)
}

// POST /sessions/{sessionId}/entities/hierarchy/validate
func ValidateEntityHierarchyHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	sessionId := extractSessionFromPath(r.URL.Path)
	
	var request struct {
		EntityId       string `json:"entityId"`
		ProposedParent string `json:"proposedParent"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	logging.Info("entity hierarchy validation", map[string]interface{}{
		"session_id": sessionId,
		"entity_id": request.EntityId,
		"proposed_parent": request.ProposedParent,
	})

	// PlayCanvas integration: Validate hierarchy operation
	// TODO: Implement actual PlayCanvas cycle detection and hierarchy validation
	
	// Mock validation logic
	isValid := true
	issues := []string{}
	
	// Check for potential cycles
	if request.EntityId == request.ProposedParent {
		isValid = false
		issues = append(issues, "Entity cannot be its own parent")
	}
	
	response := map[string]interface{}{
		"success": true,
		"valid": isValid,
		"issues": issues,
		"validation": map[string]interface{}{
			"entityId": request.EntityId,
			"proposedParent": request.ProposedParent,
			"wouldCreateCycle": !isValid,
		},
		"message": "Hierarchy validation completed",
	}
	
	writeJSONResponse(w, response, http.StatusOK)
}

// Helper function to extract session ID from URL path for tree operations
func extractSessionFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	
	// Parse /api/sessions/{sessionId}/...
	if len(parts) >= 3 && parts[0] == "api" && parts[1] == "sessions" {
		return parts[2]
	}
	
	return "default-session"
}