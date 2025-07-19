package materials

import (
	"encoding/json"
	"net/http"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
	"holodeck1/sync"
)

// CreateBasicMaterial handles POST /materials/basic
func CreateBasicMaterial(w http.ResponseWriter, r *http.Request) {
	hub := r.Context().Value("hub").(*server.Hub)
	
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create Three.js basic material via sync operation
	operation := &sync.Operation{
		ClientID:  getClientID(r),
		Type:      "material_create",
		Data: map[string]interface{}{
			"material": map[string]interface{}{
				"type":      "basic",
				"color":     getString(req, "color", "#ffffff"),
				"wireframe": getBool(req, "wireframe", false),
			},
		},
		Timestamp: time.Now(),
	}

	hub.GetSync().SubmitOperation(operation)
	seqNum := operation.SeqNum

	logging.Info("Basic material created", map[string]interface{}{
		"color":      getString(req, "color", "#ffffff"),
		"wireframe":  getBool(req, "wireframe", false),
		"seq_num":    seqNum,
		"endpoint":   "/materials/basic",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"material_id": generateMaterialID(),
		"seq_num":     seqNum,
	})
}

// CreatePhongMaterial handles POST /materials/phong
func CreatePhongMaterial(w http.ResponseWriter, r *http.Request) {
	hub := r.Context().Value("hub").(*server.Hub)
	
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create Three.js phong material via sync operation
	operation := &sync.Operation{
		ClientID:  getClientID(r),
		Type:      "material_create",
		Data: map[string]interface{}{
			"material": map[string]interface{}{
				"type":       "phong",
				"color":      getString(req, "color", "#ffffff"),
				"emissive":   getString(req, "emissive", "#000000"),
				"specular":   getString(req, "specular", "#111111"),
				"shininess":  getFloat(req, "shininess", 30.0),
				"wireframe":  getBool(req, "wireframe", false),
			},
		},
		Timestamp: time.Now(),
	}

	hub.GetSync().SubmitOperation(operation)
	seqNum := operation.SeqNum

	logging.Info("Phong material created", map[string]interface{}{
		"color":      getString(req, "color", "#ffffff"),
		"emissive":   getString(req, "emissive", "#000000"),
		"specular":   getString(req, "specular", "#111111"),
		"shininess":  getFloat(req, "shininess", 30.0),
		"wireframe":  getBool(req, "wireframe", false),
		"seq_num":    seqNum,
		"endpoint":   "/materials/phong",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"material_id": generateMaterialID(),
		"seq_num":     seqNum,
	})
}

// CreateStandardMaterial handles POST /materials/standard
func CreateStandardMaterial(w http.ResponseWriter, r *http.Request) {
	hub := r.Context().Value("hub").(*server.Hub)
	
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create Three.js standard material via sync operation
	operation := &sync.Operation{
		ClientID:  getClientID(r),
		Type:      "material_create",
		Data: map[string]interface{}{
			"material": map[string]interface{}{
				"type":       "standard",
				"color":      getString(req, "color", "#ffffff"),
				"emissive":   getString(req, "emissive", "#000000"),
				"metalness":  getFloat(req, "metalness", 0.0),
				"roughness":  getFloat(req, "roughness", 1.0),
				"wireframe":  getBool(req, "wireframe", false),
			},
		},
		Timestamp: time.Now(),
	}

	hub.GetSync().SubmitOperation(operation)
	seqNum := operation.SeqNum

	logging.Info("Standard material created", map[string]interface{}{
		"color":      getString(req, "color", "#ffffff"),
		"emissive":   getString(req, "emissive", "#000000"),
		"metalness":  getFloat(req, "metalness", 0.0),
		"roughness":  getFloat(req, "roughness", 1.0),
		"wireframe":  getBool(req, "wireframe", false),
		"seq_num":    seqNum,
		"endpoint":   "/materials/standard",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"material_id": generateMaterialID(),
		"seq_num":     seqNum,
	})
}

// CreatePhysicalMaterial handles POST /materials/physical
func CreatePhysicalMaterial(w http.ResponseWriter, r *http.Request) {
	hub := r.Context().Value("hub").(*server.Hub)
	
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create Three.js physical material via sync operation
	operation := &sync.Operation{
		ClientID:  getClientID(r),
		Type:      "material_create",
		Data: map[string]interface{}{
			"material": map[string]interface{}{
				"type":               "physical",
				"color":              getString(req, "color", "#ffffff"),
				"emissive":           getString(req, "emissive", "#000000"),
				"metalness":          getFloat(req, "metalness", 0.0),
				"roughness":          getFloat(req, "roughness", 1.0),
				"clearcoat":          getFloat(req, "clearcoat", 0.0),
				"clearcoatRoughness": getFloat(req, "clearcoatRoughness", 0.0),
				"transmission":       getFloat(req, "transmission", 0.0),
				"thickness":          getFloat(req, "thickness", 0.0),
				"ior":                getFloat(req, "ior", 1.5),
			},
		},
		Timestamp: time.Now(),
	}

	hub.GetSync().SubmitOperation(operation)
	seqNum := operation.SeqNum

	logging.Info("Physical material created", map[string]interface{}{
		"color":             getString(req, "color", "#ffffff"),
		"emissive":          getString(req, "emissive", "#000000"),
		"metalness":         getFloat(req, "metalness", 0.0),
		"roughness":         getFloat(req, "roughness", 1.0),
		"clearcoat":         getFloat(req, "clearcoat", 0.0),
		"clearcoatRoughness": getFloat(req, "clearcoatRoughness", 0.0),
		"transmission":      getFloat(req, "transmission", 0.0),
		"thickness":         getFloat(req, "thickness", 0.0),
		"ior":               getFloat(req, "ior", 1.5),
		"seq_num":           seqNum,
		"endpoint":          "/materials/physical",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"material_id": generateMaterialID(),
		"seq_num":     seqNum,
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

func generateMaterialID() string {
	return "material_" + "generated"
}

func getClientID(r *http.Request) string {
	// Try to get client ID from various sources
	if clientID := r.Header.Get("X-HD1-ID"); clientID != "" {
		return clientID
	}
	
	// Generate a client ID if none provided
	return "api-client-" + time.Now().Format("20060102150405")
}