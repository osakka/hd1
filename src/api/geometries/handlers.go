package geometries

import (
	"encoding/json"
	"net/http"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
	"holodeck1/sync"
)

// CreateBoxGeometry handles POST /geometries/box
func CreateBoxGeometry(w http.ResponseWriter, r *http.Request) {
	hub := r.Context().Value("hub").(*server.Hub)
	
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create Three.js box geometry via sync operation
	operation := &sync.Operation{
		ClientID:  getClientID(r),
		Type:      "entity_create",
		Data: map[string]interface{}{
			"geometry": map[string]interface{}{
				"type":           "box",
				"width":          getFloat(req, "width", 1.0),
				"height":         getFloat(req, "height", 1.0),
				"depth":          getFloat(req, "depth", 1.0),
				"widthSegments":  getInt(req, "widthSegments", 1),
				"heightSegments": getInt(req, "heightSegments", 1),
				"depthSegments":  getInt(req, "depthSegments", 1),
			},
			"position": req["position"],
			"rotation": req["rotation"],
			"scale":    req["scale"],
			"material": req["material"],
		},
		Timestamp: time.Now(),
	}

	hub.GetSync().SubmitOperation(operation)
	seqNum := operation.SeqNum

	logging.Info("Box geometry created", map[string]interface{}{
		"width":    getFloat(req, "width", 1.0),
		"height":   getFloat(req, "height", 1.0),
		"depth":    getFloat(req, "depth", 1.0),
		"seq_num":  seqNum,
		"endpoint": "/geometries/box",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"entity_id": generateEntityID(),
		"seq_num":   seqNum,
	})
}

// CreateSphereGeometry handles POST /geometries/sphere
func CreateSphereGeometry(w http.ResponseWriter, r *http.Request) {
	hub := r.Context().Value("hub").(*server.Hub)
	
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create Three.js sphere geometry via sync operation
	operation := &sync.Operation{
		ClientID:  getClientID(r),
		Type:      "entity_create",
		Data: map[string]interface{}{
			"geometry": map[string]interface{}{
				"type":           "sphere",
				"radius":         getFloat(req, "radius", 1.0),
				"widthSegments":  getInt(req, "widthSegments", 32),
				"heightSegments": getInt(req, "heightSegments", 16),
				"phiStart":       getFloat(req, "phiStart", 0.0),
				"phiLength":      getFloat(req, "phiLength", 6.283185307179586),
				"thetaStart":     getFloat(req, "thetaStart", 0.0),
				"thetaLength":    getFloat(req, "thetaLength", 3.141592653589793),
			},
			"position": req["position"],
			"rotation": req["rotation"],
			"scale":    req["scale"],
			"material": req["material"],
		},
		Timestamp: time.Now(),
	}

	hub.GetSync().SubmitOperation(operation)
	seqNum := operation.SeqNum

	logging.Info("Sphere geometry created", map[string]interface{}{
		"radius":  getFloat(req, "radius", 1.0),
		"seq_num": seqNum,
		"endpoint": "/geometries/sphere",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"entity_id": generateEntityID(),
		"seq_num":   seqNum,
	})
}

// CreateCylinderGeometry handles POST /geometries/cylinder
func CreateCylinderGeometry(w http.ResponseWriter, r *http.Request) {
	hub := r.Context().Value("hub").(*server.Hub)
	
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create Three.js cylinder geometry via sync operation
	operation := &sync.Operation{
		ClientID:  getClientID(r),
		Type:      "entity_create",
		Data: map[string]interface{}{
			"geometry": map[string]interface{}{
				"type":           "cylinder",
				"radiusTop":      getFloat(req, "radiusTop", 1.0),
				"radiusBottom":   getFloat(req, "radiusBottom", 1.0),
				"height":         getFloat(req, "height", 1.0),
				"radialSegments": getInt(req, "radialSegments", 32),
				"heightSegments": getInt(req, "heightSegments", 1),
				"openEnded":      getBool(req, "openEnded", false),
				"thetaStart":     getFloat(req, "thetaStart", 0.0),
				"thetaLength":    getFloat(req, "thetaLength", 6.283185307179586),
			},
			"position": req["position"],
			"rotation": req["rotation"],
			"scale":    req["scale"],
			"material": req["material"],
		},
		Timestamp: time.Now(),
	}

	hub.GetSync().SubmitOperation(operation)
	seqNum := operation.SeqNum

	logging.Info("Cylinder geometry created", map[string]interface{}{
		"radiusTop":    getFloat(req, "radiusTop", 1.0),
		"radiusBottom": getFloat(req, "radiusBottom", 1.0),
		"height":       getFloat(req, "height", 1.0),
		"seq_num":      seqNum,
		"endpoint":     "/geometries/cylinder",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"entity_id": generateEntityID(),
		"seq_num":   seqNum,
	})
}

// CreateConeGeometry handles POST /geometries/cone
func CreateConeGeometry(w http.ResponseWriter, r *http.Request) {
	hub := r.Context().Value("hub").(*server.Hub)
	
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create Three.js cone geometry via sync operation
	operation := &sync.Operation{
		ClientID:  getClientID(r),
		Type:      "entity_create",
		Data: map[string]interface{}{
			"geometry": map[string]interface{}{
				"type":           "cone",
				"radius":         getFloat(req, "radius", 1.0),
				"height":         getFloat(req, "height", 1.0),
				"radialSegments": getInt(req, "radialSegments", 32),
				"heightSegments": getInt(req, "heightSegments", 1),
				"openEnded":      getBool(req, "openEnded", false),
				"thetaStart":     getFloat(req, "thetaStart", 0.0),
				"thetaLength":    getFloat(req, "thetaLength", 6.283185307179586),
			},
			"position": req["position"],
			"rotation": req["rotation"],
			"scale":    req["scale"],
			"material": req["material"],
		},
		Timestamp: time.Now(),
	}

	hub.GetSync().SubmitOperation(operation)
	seqNum := operation.SeqNum

	logging.Info("Cone geometry created", map[string]interface{}{
		"radius":  getFloat(req, "radius", 1.0),
		"height":  getFloat(req, "height", 1.0),
		"seq_num": seqNum,
		"endpoint": "/geometries/cone",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"entity_id": generateEntityID(),
		"seq_num":   seqNum,
	})
}

// CreateTorusGeometry handles POST /geometries/torus
func CreateTorusGeometry(w http.ResponseWriter, r *http.Request) {
	hub := r.Context().Value("hub").(*server.Hub)
	
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create Three.js torus geometry via sync operation
	operation := &sync.Operation{
		ClientID:  getClientID(r),
		Type:      "entity_create",
		Data: map[string]interface{}{
			"geometry": map[string]interface{}{
				"type":            "torus",
				"radius":          getFloat(req, "radius", 1.0),
				"tube":            getFloat(req, "tube", 0.4),
				"radialSegments":  getInt(req, "radialSegments", 12),
				"tubularSegments": getInt(req, "tubularSegments", 48),
				"arc":             getFloat(req, "arc", 6.283185307179586),
			},
			"position": req["position"],
			"rotation": req["rotation"],
			"scale":    req["scale"],
			"material": req["material"],
		},
		Timestamp: time.Now(),
	}

	hub.GetSync().SubmitOperation(operation)
	seqNum := operation.SeqNum

	logging.Info("Torus geometry created", map[string]interface{}{
		"radius":  getFloat(req, "radius", 1.0),
		"tube":    getFloat(req, "tube", 0.4),
		"seq_num": seqNum,
		"endpoint": "/geometries/torus",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"entity_id": generateEntityID(),
		"seq_num":   seqNum,
	})
}

// CreateTorusKnotGeometry handles POST /geometries/torusknot
func CreateTorusKnotGeometry(w http.ResponseWriter, r *http.Request) {
	hub := r.Context().Value("hub").(*server.Hub)
	
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create Three.js torus knot geometry via sync operation
	operation := &sync.Operation{
		ClientID:  getClientID(r),
		Type:      "entity_create",
		Data: map[string]interface{}{
			"geometry": map[string]interface{}{
				"type":            "torusknot",
				"radius":          getFloat(req, "radius", 1.0),
				"tube":            getFloat(req, "tube", 0.4),
				"tubularSegments": getInt(req, "tubularSegments", 64),
				"radialSegments":  getInt(req, "radialSegments", 8),
				"p":               getInt(req, "p", 2),
				"q":               getInt(req, "q", 3),
			},
			"position": req["position"],
			"rotation": req["rotation"],
			"scale":    req["scale"],
			"material": req["material"],
		},
		Timestamp: time.Now(),
	}

	hub.GetSync().SubmitOperation(operation)
	seqNum := operation.SeqNum

	logging.Info("Torus knot geometry created", map[string]interface{}{
		"radius":  getFloat(req, "radius", 1.0),
		"tube":    getFloat(req, "tube", 0.4),
		"p":       getInt(req, "p", 2),
		"q":       getInt(req, "q", 3),
		"seq_num": seqNum,
		"endpoint": "/geometries/torusknot",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"entity_id": generateEntityID(),
		"seq_num":   seqNum,
	})
}

// CreatePlaneGeometry handles POST /geometries/plane
func CreatePlaneGeometry(w http.ResponseWriter, r *http.Request) {
	hub := r.Context().Value("hub").(*server.Hub)
	
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create Three.js plane geometry via sync operation
	operation := &sync.Operation{
		ClientID:  getClientID(r),
		Type:      "entity_create",
		Data: map[string]interface{}{
			"geometry": map[string]interface{}{
				"type":           "plane",
				"width":          getFloat(req, "width", 1.0),
				"height":         getFloat(req, "height", 1.0),
				"widthSegments":  getInt(req, "widthSegments", 1),
				"heightSegments": getInt(req, "heightSegments", 1),
			},
			"position": req["position"],
			"rotation": req["rotation"],
			"scale":    req["scale"],
			"material": req["material"],
		},
		Timestamp: time.Now(),
	}

	hub.GetSync().SubmitOperation(operation)
	seqNum := operation.SeqNum

	logging.Info("Plane geometry created", map[string]interface{}{
		"width":   getFloat(req, "width", 1.0),
		"height":  getFloat(req, "height", 1.0),
		"seq_num": seqNum,
		"endpoint": "/geometries/plane",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"entity_id": generateEntityID(),
		"seq_num":   seqNum,
	})
}

// CreateRingGeometry handles POST /geometries/ring
func CreateRingGeometry(w http.ResponseWriter, r *http.Request) {
	hub := r.Context().Value("hub").(*server.Hub)
	
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create Three.js ring geometry via sync operation
	operation := &sync.Operation{
		ClientID:  getClientID(r),
		Type:      "entity_create",
		Data: map[string]interface{}{
			"geometry": map[string]interface{}{
				"type":         "ring",
				"innerRadius":  getFloat(req, "innerRadius", 0.5),
				"outerRadius":  getFloat(req, "outerRadius", 1.0),
				"thetaSegments": getInt(req, "thetaSegments", 32),
				"phiSegments":  getInt(req, "phiSegments", 1),
				"thetaStart":   getFloat(req, "thetaStart", 0.0),
				"thetaLength":  getFloat(req, "thetaLength", 6.283185307179586),
			},
			"position": req["position"],
			"rotation": req["rotation"],
			"scale":    req["scale"],
			"material": req["material"],
		},
		Timestamp: time.Now(),
	}

	hub.GetSync().SubmitOperation(operation)
	seqNum := operation.SeqNum

	logging.Info("Ring geometry created", map[string]interface{}{
		"innerRadius": getFloat(req, "innerRadius", 0.5),
		"outerRadius": getFloat(req, "outerRadius", 1.0),
		"seq_num":     seqNum,
		"endpoint":    "/geometries/ring",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"entity_id": generateEntityID(),
		"seq_num":   seqNum,
	})
}

// CreateCircleGeometry handles POST /geometries/circle
func CreateCircleGeometry(w http.ResponseWriter, r *http.Request) {
	hub := r.Context().Value("hub").(*server.Hub)
	
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create Three.js circle geometry via sync operation
	operation := &sync.Operation{
		ClientID:  getClientID(r),
		Type:      "entity_create",
		Data: map[string]interface{}{
			"geometry": map[string]interface{}{
				"type":        "circle",
				"radius":      getFloat(req, "radius", 1.0),
				"segments":    getInt(req, "segments", 32),
				"thetaStart":  getFloat(req, "thetaStart", 0.0),
				"thetaLength": getFloat(req, "thetaLength", 6.283185307179586),
			},
			"position": req["position"],
			"rotation": req["rotation"],
			"scale":    req["scale"],
			"material": req["material"],
		},
		Timestamp: time.Now(),
	}

	hub.GetSync().SubmitOperation(operation)
	seqNum := operation.SeqNum

	logging.Info("Circle geometry created", map[string]interface{}{
		"radius":  getFloat(req, "radius", 1.0),
		"segments": getInt(req, "segments", 32),
		"seq_num": seqNum,
		"endpoint": "/geometries/circle",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"entity_id": generateEntityID(),
		"seq_num":   seqNum,
	})
}

// CreateCapsuleGeometry handles POST /geometries/capsule
func CreateCapsuleGeometry(w http.ResponseWriter, r *http.Request) {
	hub := r.Context().Value("hub").(*server.Hub)
	
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create Three.js capsule geometry via sync operation
	operation := &sync.Operation{
		ClientID:  getClientID(r),
		Type:      "entity_create",
		Data: map[string]interface{}{
			"geometry": map[string]interface{}{
				"type":           "capsule",
				"radius":         getFloat(req, "radius", 1.0),
				"length":         getFloat(req, "length", 1.0),
				"capSegments":    getInt(req, "capSegments", 4),
				"radialSegments": getInt(req, "radialSegments", 8),
			},
			"position": req["position"],
			"rotation": req["rotation"],
			"scale":    req["scale"],
			"material": req["material"],
		},
		Timestamp: time.Now(),
	}

	hub.GetSync().SubmitOperation(operation)
	seqNum := operation.SeqNum

	logging.Info("Capsule geometry created", map[string]interface{}{
		"radius":  getFloat(req, "radius", 1.0),
		"length":  getFloat(req, "length", 1.0),
		"seq_num": seqNum,
		"endpoint": "/geometries/capsule",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"entity_id": generateEntityID(),
		"seq_num":   seqNum,
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

func getInt(req map[string]interface{}, key string, defaultValue int) int {
	if val, ok := req[key]; ok {
		if i, ok := val.(float64); ok {
			return int(i)
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

func generateEntityID() string {
	// Generate a simple entity ID for now - in production this would be a proper UUID
	return "entity_" + "generated"
}

func getClientID(r *http.Request) string {
	// Try to get client ID from various sources
	if clientID := r.Header.Get("X-HD1-ID"); clientID != "" {
		return clientID
	}
	
	// Generate a client ID if none provided
	return "api-client-" + time.Now().Format("20060102150405")
}