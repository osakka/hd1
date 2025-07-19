package lights

import (
	"encoding/json"
	"net/http"
	"time"

	"holodeck1/api/shared"
	"holodeck1/logging"
	"holodeck1/sync"
)

// CreateDirectionalLight handles POST /lights/directional
func CreateDirectionalLight(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	operation := &sync.Operation{
		ClientID:  getClientID(r),
		Type:      "scene_update",
		Data: map[string]interface{}{
			"operation": "add_light",
			"light": map[string]interface{}{
				"type":       "directional",
				"color":      getString(req, "color", "#ffffff"),
				"intensity":  getFloat(req, "intensity", 1.0),
				"position":   req["position"],
				"target":     req["target"],
				"castShadow": getBool(req, "castShadow", false),
			},
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
	lightID := generateLightID()

	logging.Info("Directional light created", map[string]interface{}{
		"color":     getString(req, "color", "#ffffff"),
		"intensity": getFloat(req, "intensity", 1.0),
		"light_id":  lightID,
		"seq_num":   seqNum,
		"endpoint":  "/lights/directional",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"light_id": lightID,
		"seq_num":  seqNum,
	})
}

// CreatePointLight handles POST /lights/point
func CreatePointLight(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	operation := &sync.Operation{
		ClientID:  getClientID(r),
		Type:      "scene_update",
		Data: map[string]interface{}{
			"operation": "add_light",
			"light": map[string]interface{}{
				"type":       "point",
				"color":      getString(req, "color", "#ffffff"),
				"intensity":  getFloat(req, "intensity", 1.0),
				"distance":   getFloat(req, "distance", 0.0),
				"decay":      getFloat(req, "decay", 2.0),
				"position":   req["position"],
				"castShadow": getBool(req, "castShadow", false),
			},
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
	lightID := generateLightID()

	logging.Info("Point light created", map[string]interface{}{
		"color":     getString(req, "color", "#ffffff"),
		"intensity": getFloat(req, "intensity", 1.0),
		"distance":  getFloat(req, "distance", 0.0),
		"light_id":  lightID,
		"seq_num":   seqNum,
		"endpoint":  "/lights/point",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"light_id": lightID,
		"seq_num":  seqNum,
	})
}

// CreateSpotLight handles POST /lights/spot
func CreateSpotLight(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	operation := &sync.Operation{
		ClientID:  getClientID(r),
		Type:      "scene_update",
		Data: map[string]interface{}{
			"operation": "add_light",
			"light": map[string]interface{}{
				"type":       "spot",
				"color":      getString(req, "color", "#ffffff"),
				"intensity":  getFloat(req, "intensity", 1.0),
				"distance":   getFloat(req, "distance", 0.0),
				"angle":      getFloat(req, "angle", 1.0471975511965976),
				"penumbra":   getFloat(req, "penumbra", 0.0),
				"decay":      getFloat(req, "decay", 2.0),
				"position":   req["position"],
				"target":     req["target"],
				"castShadow": getBool(req, "castShadow", false),
			},
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
	lightID := generateLightID()

	logging.Info("Spot light created", map[string]interface{}{
		"color":     getString(req, "color", "#ffffff"),
		"intensity": getFloat(req, "intensity", 1.0),
		"angle":     getFloat(req, "angle", 1.0471975511965976),
		"light_id":  lightID,
		"seq_num":   seqNum,
		"endpoint":  "/lights/spot",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"light_id": lightID,
		"seq_num":  seqNum,
	})
}

// CreateAmbientLight handles POST /lights/ambient
func CreateAmbientLight(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	operation := &sync.Operation{
		ClientID:  getClientID(r),
		Type:      "scene_update",
		Data: map[string]interface{}{
			"operation": "add_light",
			"light": map[string]interface{}{
				"type":      "ambient",
				"color":     getString(req, "color", "#ffffff"),
				"intensity": getFloat(req, "intensity", 1.0),
			},
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
	lightID := generateLightID()

	logging.Info("Ambient light created", map[string]interface{}{
		"color":     getString(req, "color", "#ffffff"),
		"intensity": getFloat(req, "intensity", 1.0),
		"light_id":  lightID,
		"seq_num":   seqNum,
		"endpoint":  "/lights/ambient",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"light_id": lightID,
		"seq_num":  seqNum,
	})
}

// CreateHemisphereLight handles POST /lights/hemisphere
func CreateHemisphereLight(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	operation := &sync.Operation{
		ClientID:  getClientID(r),
		Type:      "scene_update",
		Data: map[string]interface{}{
			"operation": "add_light",
			"light": map[string]interface{}{
				"type":        "hemisphere",
				"skyColor":    getString(req, "skyColor", "#ffffff"),
				"groundColor": getString(req, "groundColor", "#444444"),
				"intensity":   getFloat(req, "intensity", 1.0),
				"position":    req["position"],
			},
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
	lightID := generateLightID()

	logging.Info("Hemisphere light created", map[string]interface{}{
		"skyColor":    getString(req, "skyColor", "#ffffff"),
		"groundColor": getString(req, "groundColor", "#444444"),
		"intensity":   getFloat(req, "intensity", 1.0),
		"light_id":    lightID,
		"seq_num":     seqNum,
		"endpoint":    "/lights/hemisphere",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"light_id": lightID,
		"seq_num":  seqNum,
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

func generateLightID() string {
	return "light_" + "generated"
}

func getClientID(r *http.Request) string {
	// Try to get client ID from various sources
	if clientID := r.Header.Get("X-HD1-ID"); clientID != "" {
		return clientID
	}
	
	// Fallback to generating a client ID based on request info
	return "api-client-" + time.Now().Format("20060102150405")
}