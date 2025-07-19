package animations

import (
	"encoding/json"
	"net/http"
	"time"

	"holodeck1/api/shared"
	"holodeck1/logging"
	"holodeck1/sync"
)

// CreateKeyframeAnimation handles POST /animations/keyframe
func CreateKeyframeAnimation(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	animationID := generateAnimationID()
	clientID := shared.GetClientID(r)

	// Create sync operation
	operation := &sync.Operation{
		ClientID:  clientID,
		Type:      "animation_create",
		Data: map[string]interface{}{
			"animation_id": animationID,
			"target":       getString(req, "target", ""),
			"property":     getString(req, "property", ""),
			"duration":     getFloat(req, "duration", 1.0),
			"loop":         getBool(req, "loop", false),
			"easing":       getString(req, "easing", "linear"),
			"keyframes":    req["keyframes"],
		},
		Timestamp: time.Now(),
	}

	// Get hub and submit operation
	hub := shared.GetHubFromContext(r)
	if hub == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	hub.GetSync().SubmitOperation(operation)
	seqNum := operation.SeqNum

	logging.Info("Keyframe animation created", map[string]interface{}{
		"target":        getString(req, "target", ""),
		"property":      getString(req, "property", ""),
		"duration":      getFloat(req, "duration", 1.0),
		"loop":          getBool(req, "loop", false),
		"easing":        getString(req, "easing", "linear"),
		"animation_id":  animationID,
		"seq_num":       seqNum,
		"endpoint":      "/animations/keyframe",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"animation_id": animationID,
		"seq_num":      seqNum,
	})
}

// ControlTimeline handles POST /animations/timeline
func ControlTimeline(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	action := getString(req, "action", "play")
	speed := getFloat(req, "speed", 1.0)
	timeValue := getFloat(req, "time", 0.0)
	clientID := shared.GetClientID(r)

	// Create sync operation
	operation := &sync.Operation{
		ClientID:  clientID,
		Type:      "animation_control",
		Data: map[string]interface{}{
			"action": action,
			"speed":  speed,
			"time":   timeValue,
			"target": getString(req, "target", ""),
		},
		Timestamp: time.Now(),
	}

	// Get hub and submit operation
	hub := shared.GetHubFromContext(r)
	if hub == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	hub.GetSync().SubmitOperation(operation)
	seqNum := operation.SeqNum

	logging.Info("Animation timeline controlled", map[string]interface{}{
		"action":   action,
		"speed":    speed,
		"time":     timeValue,
		"seq_num":  seqNum,
		"endpoint": "/animations/timeline",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"action":  action,
		"time":    timeValue,
		"seq_num": seqNum,
	})
}

// Helper functions
func getFloat(req map[string]interface{}, key string, defaultValue float64) float64 {
	if val, ok := req[key]; ok {
		if f, ok := val.(float64); ok {
			return f
		}
	}
	return defaultValue
}

func getBool(req map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := req[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return defaultValue
}

func getString(req map[string]interface{}, key string, defaultValue string) string {
	if val, ok := req[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return defaultValue
}

func generateAnimationID() string {
	return "animation_" + "generated"
}