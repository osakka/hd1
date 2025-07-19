package cameras

import (
	"encoding/json"
	"net/http"
	"time"

	"holodeck1/api/shared"
	"holodeck1/logging"
	"holodeck1/sync"
)

// SetPerspectiveCamera handles POST /cameras/perspective
func SetPerspectiveCamera(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	operation := &sync.Operation{
		ClientID: getClientID(r),
		Type:     "scene_update",
		Data: map[string]interface{}{
			"operation": "set_camera",
			"camera": map[string]interface{}{
				"type":     "perspective",
				"fov":      getFloat(req, "fov", 75.0),
				"aspect":   getFloat(req, "aspect", 1.0),
				"near":     getFloat(req, "near", 0.1),
				"far":      getFloat(req, "far", 1000.0),
				"position": req["position"],
				"rotation": req["rotation"],
				"lookAt":   req["lookAt"],
			},
		},
		Timestamp: time.Now(),
	}

	hub := shared.GetHubFromContext(r)
	if hub == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	hub.GetSync().SubmitOperation(operation)
	seqNum := operation.SeqNum

	logging.Info("Perspective camera configured", map[string]interface{}{
		"fov":      getFloat(req, "fov", 75.0),
		"aspect":   getFloat(req, "aspect", 1.0),
		"near":     getFloat(req, "near", 0.1),
		"far":      getFloat(req, "far", 1000.0),
		"seq_num":  seqNum,
		"endpoint": "/cameras/perspective",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"camera_type": "perspective",
		"seq_num":     seqNum,
	})
}

// SetOrthographicCamera handles POST /cameras/orthographic
func SetOrthographicCamera(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	operation := &sync.Operation{
		ClientID: getClientID(r),
		Type:     "scene_update",
		Data: map[string]interface{}{
			"operation": "set_camera",
			"camera": map[string]interface{}{
				"type":     "orthographic",
				"left":     getFloat(req, "left", -1.0),
				"right":    getFloat(req, "right", 1.0),
				"top":      getFloat(req, "top", 1.0),
				"bottom":   getFloat(req, "bottom", -1.0),
				"near":     getFloat(req, "near", 0.1),
				"far":      getFloat(req, "far", 1000.0),
				"position": req["position"],
				"rotation": req["rotation"],
			},
		},
		Timestamp: time.Now(),
	}

	hub := shared.GetHubFromContext(r)
	if hub == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	hub.GetSync().SubmitOperation(operation)
	seqNum := operation.SeqNum

	logging.Info("Orthographic camera configured", map[string]interface{}{
		"left":     getFloat(req, "left", -1.0),
		"right":    getFloat(req, "right", 1.0),
		"top":      getFloat(req, "top", 1.0),
		"bottom":   getFloat(req, "bottom", -1.0),
		"near":     getFloat(req, "near", 0.1),
		"far":      getFloat(req, "far", 1000.0),
		"seq_num":  seqNum,
		"endpoint": "/cameras/orthographic",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"camera_type": "orthographic",
		"seq_num":     seqNum,
	})
}

// Helper functions
func getClientID(r *http.Request) string {
	// Extract client ID from request headers or context
	if clientID := r.Header.Get("X-Client-ID"); clientID != "" {
		return clientID
	}
	// Fallback to remote address if no client ID header
	return r.RemoteAddr
}

func getFloat(req map[string]interface{}, key string, defaultValue float64) float64 {
	if val, ok := req[key]; ok {
		if f, ok := val.(float64); ok {
			return f
		}
	}
	return defaultValue
}