package components

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
)

// BulkOperation represents a single operation in a bulk request
type BulkOperation struct {
	Action        string                 `json:"action"`
	ComponentType string                 `json:"component_type"`
	Properties    map[string]interface{} `json:"properties,omitempty"`
}

// BulkComponentRequest represents the request body for bulk component operations
type BulkComponentRequest struct {
	Operations []BulkOperation `json:"operations"`
	Atomic     *bool           `json:"atomic,omitempty"`
}

// BulkOperationResult represents the result of a single bulk operation
type BulkOperationResult struct {
	Action        string `json:"action"`
	ComponentType string `json:"component_type"`
	Success       bool   `json:"success"`
	Error         string `json:"error,omitempty"`
}

// BulkComponentOperationHandler handles POST /sessions/{sessionId}/entities/{entityId}/components/bulk
func BulkComponentOperationHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
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
	
	// Extract sessionId and entityId from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/"), "/")
	if len(pathParts) < 6 || pathParts[0] != "sessions" || pathParts[2] != "entities" || pathParts[4] != "components" || pathParts[5] != "bulk" {
		logging.Error("invalid URL path for bulk component operations", map[string]interface{}{
			"path": r.URL.Path,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid URL path",
			"message": "Expected /api/sessions/{sessionId}/entities/{entityId}/components/bulk",
		})
		return
	}
	
	sessionID := pathParts[1]
	entityID := pathParts[3]
	
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
	
	// Parse request body
	var req BulkComponentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.Error("failed to parse request body", map[string]interface{}{
			"error": err.Error(),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
			"message": "Failed to parse JSON request",
		})
		return
	}
	
	// Set default atomic behavior
	atomic := true
	if req.Atomic != nil {
		atomic = *req.Atomic
	}
	
	// Validate operations
	validTypes := []string{"model", "camera", "light", "rigidbody", "script", "sound", "animation", "collision", "element", "particlesystem", "render", "sprite"}
	validActions := []string{"add", "remove", "update", "enable", "disable"}
	
	var results []BulkOperationResult
	completed := 0
	failed := 0
	
	// TODO: Implement atomic transaction logic for PlayCanvas operations
	// For now, process each operation independently
	
	for _, operation := range req.Operations {
		result := BulkOperationResult{
			Action:        operation.Action,
			ComponentType: operation.ComponentType,
			Success:       false,
		}
		
		// Validate action
		isValidAction := false
		for _, validAction := range validActions {
			if operation.Action == validAction {
				isValidAction = true
				break
			}
		}
		
		if !isValidAction {
			result.Error = "Invalid action: " + operation.Action + ". Must be one of: " + strings.Join(validActions, ", ")
			results = append(results, result)
			failed++
			if atomic {
				// In atomic mode, fail all operations if any fails
				break
			}
			continue
		}
		
		// Validate component type
		isValidType := false
		for _, validType := range validTypes {
			if operation.ComponentType == validType {
				isValidType = true
				break
			}
		}
		
		if !isValidType {
			result.Error = "Invalid component type: " + operation.ComponentType + ". Must be one of: " + strings.Join(validTypes, ", ")
			results = append(results, result)
			failed++
			if atomic {
				// In atomic mode, fail all operations if any fails
				break
			}
			continue
		}
		
		// TODO: Implement actual PlayCanvas component operations
		// For now, mock successful operations
		switch operation.Action {
		case "add":
			// Mock component addition
			result.Success = true
			completed++
		case "remove":
			// Mock component removal
			result.Success = true
			completed++
		case "update":
			// Mock component update
			result.Success = true
			completed++
		case "enable", "disable":
			// Mock component enable/disable
			result.Success = true
			completed++
		}
		
		results = append(results, result)
	}
	
	// If atomic mode and any operation failed, mark all as failed
	if atomic && failed > 0 {
		for i := range results {
			if results[i].Success {
				results[i].Success = false
				results[i].Error = "Transaction rolled back due to atomic failure"
			}
		}
		completed = 0
		failed = len(results)
	}
	
	logging.Info("bulk component operations completed", map[string]interface{}{
		"session_id":  sessionID,
		"entity_id":   entityID,
		"total_ops":   len(req.Operations),
		"completed":   completed,
		"failed":      failed,
		"atomic":      atomic,
	})
	
	// Broadcast bulk operations event via WebSocket
	if h != nil && completed > 0 {
		h.BroadcastToSession(sessionID, "bulk_component_operations", map[string]interface{}{
			"entity_id": entityID,
			"completed": completed,
			"failed":    failed,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   failed == 0,
		"results":   results,
		"completed": completed,
		"failed":    failed,
	})
}