package objects

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"holodeck1/server"
)

// UpdateObjectRequest represents the request body for updating objects
type UpdateObjectRequest struct {
	X        *float64 `json:"x,omitempty"`
	Y        *float64 `json:"y,omitempty"`
	Z        *float64 `json:"z,omitempty"`
	Color    *Color   `json:"color,omitempty"`
	Scale    *float64 `json:"scale,omitempty"`
	Rotation *string  `json:"rotation,omitempty"`
}

func UpdateObjectHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type
	h, ok := hub.(*server.Hub)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Extract session ID and object name from path
	sessionID := extractSessionID(r.URL.Path)
	objectName := extractObjectNameFromPath(r.URL.Path)
	if sessionID == "" || objectName == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid parameters",
		})
		return
	}
	
	// Parse request body
	var req UpdateObjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid JSON",
		})
		return
	}
	
	// Get existing object
	store := h.GetStore()
	object, exists := store.GetObject(sessionID, objectName)
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Object not found",
		})
		return
	}
	
	// Prepare updates including position changes
	updates := map[string]interface{}{}
	if req.X != nil {
		updates["x"] = *req.X
	}
	if req.Y != nil {
		updates["y"] = *req.Y
	}
	if req.Z != nil {
		updates["z"] = *req.Z
	}
	if req.Scale != nil {
		updates["scale"] = *req.Scale
	}
	if req.Rotation != nil {
		updates["rotation"] = *req.Rotation
	}
	
	// Mark object as "modified" if it was previously "base"
	if object.TrackingStatus == "base" {
		updates["tracking_status"] = "modified"
	}
	updates["created_at"] = time.Now() // Update modification time
	
	// Save updated object
	store.UpdateObject(sessionID, objectName, updates)
	
	// Broadcast update
	h.BroadcastToSession(sessionID, "object_updated", map[string]interface{}{
		"object_name": objectName,
		"tracking_status": object.TrackingStatus,
		"timestamp": time.Now().Format(time.RFC3339),
	})
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true, 
		"message": "Object updated",
		"tracking_status": object.TrackingStatus,
	})
}

// extractObjectNameFromPath extracts object name from URL path
func extractObjectNameFromPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 5 && parts[1] == "api" && parts[2] == "sessions" && parts[4] == "objects" {
		return parts[5]
	}
	return ""
}

