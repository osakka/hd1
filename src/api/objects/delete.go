package objects

import (
	"encoding/json"
	"net/http"
	"strings"
	"holodeck1/server"
)

// DeleteObjectHandler - DELETE /sessions/{sessionId}/objects/{objectName}
func DeleteObjectHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type
	h, ok := hub.(*server.Hub)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Internal server error",
			"message": "Hub interface error",
		})
		return
	}
	
	// Extract session ID and object name from path
	sessionID := extractSessionID(r.URL.Path)
	objectName := extractObjectName(r.URL.Path)
	
	if sessionID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid session ID",
			"message": "Session ID is required",
		})
		return
	}
	
	if objectName == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid object name",
			"message": "Object name is required",
		})
		return
	}
	
	// Validate session exists
	if _, exists := h.GetStore().GetSession(sessionID); !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Session not found",
			"message": "Session does not exist or has expired",
		})
		return
	}
	
	// Check if object exists before deletion
	if _, exists := h.GetStore().GetObject(sessionID, objectName); !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Object not found",
			"message": "Object does not exist in session",
		})
		return
	}
	
	// Delete object from session store
	if !h.GetStore().DeleteObject(sessionID, objectName) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Delete failed",
			"message": "Failed to delete object from session",
		})
		return
	}
	
	// Broadcast deletion to session clients for real-time update
	deleteMessage := map[string]interface{}{
		"command":     "delete",
		"object_name": objectName,
		"session_id":  sessionID,
	}
	h.BroadcastToSession(sessionID, "canvas_control", deleteMessage)
	
	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Object deleted successfully",
		"object":  objectName,
		"session": sessionID,
	})
}

// Helper function to extract object name from URL path
func extractObjectName(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 5 && parts[len(parts)-1] != "" {
		return parts[len(parts)-1]
	}
	return ""
}
