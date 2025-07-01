package objects

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
)

// Color represents RGBA color values
type Color struct {
	R float64 `json:"r"`
	G float64 `json:"g"`
	B float64 `json:"b"`
	A float64 `json:"a"`
}

// Material represents A-Frame material properties
type Material struct {
	Shader      string  `json:"shader,omitempty"`
	Metalness   float64 `json:"metalness,omitempty"`
	Roughness   float64 `json:"roughness,omitempty"`
	Transparent bool    `json:"transparent,omitempty"`
	Emissive    bool    `json:"emissive,omitempty"`
}

// Physics represents physics simulation properties
type Physics struct {
	Enabled bool    `json:"enabled,omitempty"`
	Mass    float64 `json:"mass,omitempty"`
	Type    string  `json:"type,omitempty"` // static, dynamic, kinematic
}

// Lighting represents shadow and lighting properties
type Lighting struct {
	CastShadow    bool `json:"castShadow,omitempty"`
	ReceiveShadow bool `json:"receiveShadow,omitempty"`
}

// Light represents light source properties
type Light struct {
	LightType string  `json:"lightType,omitempty"` // ambient, directional, point, spot
	Intensity float64 `json:"intensity,omitempty"`
}

// Particle represents particle system properties
type Particle struct {
	ParticleType string `json:"particleType,omitempty"` // fire, smoke, sparkle
	Count        int    `json:"count,omitempty"`
}

// CreateObjectRequest represents the complete request body for creating A-Frame objects
type CreateObjectRequest struct {
	Name         string     `json:"name"`
	Type         string     `json:"type"`
	X            float64    `json:"x"`
	Y            float64    `json:"y"`
	Z            float64    `json:"z"`
	Color        *Color     `json:"color,omitempty"`
	Material     *Material  `json:"material,omitempty"`
	Physics      *Physics   `json:"physics,omitempty"`
	Lighting     *Lighting  `json:"lighting,omitempty"`
	Light        *Light     `json:"light,omitempty"`
	Particle     *Particle  `json:"particle,omitempty"`
	Text         string     `json:"text,omitempty"`
	LightType    string     `json:"lightType,omitempty"`
	Intensity    float64    `json:"intensity,omitempty"`
	ParticleType string     `json:"particleType,omitempty"`
	Count        int        `json:"count,omitempty"`
}

// CreateObjectHandler - POST /sessions/{sessionId}/objects
func CreateObjectHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type
	h, ok := hub.(*server.Hub)
	if !ok {
		logging.Error("hub interface cast failed", map[string]interface{}{
			"expected_type": "*server.Hub",
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Extract session ID from path
	sessionID := extractSessionID(r.URL.Path)
	if sessionID == "" {
		logging.Warn("missing session ID", map[string]interface{}{
			"path": r.URL.Path,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid session ID",
			"message": "Session ID is required",
		})
		return
	}
	
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
	var req CreateObjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.Error("invalid JSON in request", map[string]interface{}{
			"session_id": sessionID,
			"error": err.Error(),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid JSON",
			"message": "Request body must be valid JSON",
		})
		return
	}
	
	// Validate required fields
	if req.Name == "" || req.Type == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Missing required fields",
			"message": "name and type are required",
		})
		return
	}
	
	// Create object using SessionStore with coordinate validation
	object, err := h.GetStore().CreateObject(sessionID, req.Name, req.Type, req.X, req.Y, req.Z)
	if err != nil {
		logging.Error("object creation failed", map[string]interface{}{
			"session_id": sessionID,
			"object_name": req.Name,
			"object_type": req.Type,
			"position": map[string]float64{"x": req.X, "y": req.Y, "z": req.Z},
			"error": err.Error(),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Object creation failed",
			"message": err.Error(),
		})
		return
	}
	
	// Use provided color or default to green
	objectColor := map[string]interface{}{
		"r": 0.2,
		"g": 0.8,
		"b": 0.2,
		"a": 1.0,
	}
	if req.Color != nil {
		objectColor = map[string]interface{}{
			"r": req.Color.R,
			"g": req.Color.G,
			"b": req.Color.B,
			"a": req.Color.A,
		}
	}
	
	// Store color as JSON string in Object struct for session restoration
	colorJSON, _ := json.Marshal(objectColor)
	
	// Mark object as "new" for session object tracking and store color
	updates := map[string]interface{}{
		"tracking_status": "new",
		"created_at": time.Now(),
		"color": string(colorJSON),
	}
	h.GetStore().UpdateObject(sessionID, object.Name, updates)

	// Broadcast object creation for real-time updates via canvas control (session-specific)
	h.BroadcastToSession(sessionID, "canvas_control", map[string]interface{}{
		"command": "create",
		"objects": []map[string]interface{}{
			{
				"id":   object.Name,
				"name": object.Name,
				"type": object.Type,
				"transform": map[string]interface{}{
					"position": map[string]interface{}{
						"x": object.X,
						"y": object.Y,
						"z": object.Z,
					},
					"scale": map[string]interface{}{
						"x": object.Scale,
						"y": object.Scale,
						"z": object.Scale,
					},
					"rotation": map[string]interface{}{
						"x": 0,
						"y": 0,
						"z": 0,
					},
				},
				"color": objectColor,
				"material": func() map[string]interface{} {
					material := map[string]interface{}{
						"shader":     "standard",
						"metalness":  0.1,
						"roughness":  0.7,
						"transparent": false,
					}
					if req.Material != nil {
						if req.Material.Shader != "" {
							material["shader"] = req.Material.Shader
						}
						if req.Material.Metalness > 0 {
							material["metalness"] = req.Material.Metalness
						}
						if req.Material.Roughness > 0 {
							material["roughness"] = req.Material.Roughness
						}
						material["transparent"] = req.Material.Transparent
						material["emissive"] = req.Material.Emissive
					}
					return material
				}(),
				"physics": func() map[string]interface{} {
					physics := map[string]interface{}{
						"enabled": false,
						"mass":    1.0,
						"type":    "static",
					}
					if req.Physics != nil {
						physics["enabled"] = req.Physics.Enabled
						if req.Physics.Mass > 0 {
							physics["mass"] = req.Physics.Mass
						}
						if req.Physics.Type != "" {
							physics["type"] = req.Physics.Type
						}
					}
					return physics
				}(),
				"lighting": func() map[string]interface{} {
					lighting := map[string]interface{}{
						"castShadow":    true,
						"receiveShadow": true,
					}
					if req.Lighting != nil {
						lighting["castShadow"] = req.Lighting.CastShadow
						lighting["receiveShadow"] = req.Lighting.ReceiveShadow
					}
					return lighting
				}(),
				"visible":   true,
				"wireframe": false,
				// A-Frame specific properties
				"text": req.Text,
				"lightType": func() string {
					if req.LightType != "" {
						return req.LightType
					}
					if req.Light != nil && req.Light.LightType != "" {
						return req.Light.LightType
					}
					return ""
				}(),
				"intensity": func() float64 {
					if req.Intensity > 0 {
						return req.Intensity
					}
					if req.Light != nil && req.Light.Intensity > 0 {
						return req.Light.Intensity
					}
					return 1.0
				}(),
				"particleType": func() string {
					if req.ParticleType != "" {
						return req.ParticleType
					}
					if req.Particle != nil && req.Particle.ParticleType != "" {
						return req.Particle.ParticleType
					}
					return ""
				}(),
				"count": func() int {
					if req.Count > 0 {
						return req.Count
					}
					if req.Particle != nil && req.Particle.Count > 0 {
						return req.Particle.Count
					}
					return 0
				}(),
			},
		},
	})
	
	logging.Info("object created", map[string]interface{}{
		"session_id": sessionID,
		"object_name": object.Name,
		"object_type": object.Type,
		"position": map[string]float64{"x": object.X, "y": object.Y, "z": object.Z},
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"object":  object,
		"message": "Object created successfully",
	})
}

// extractSessionID extracts session ID from URL path
func extractSessionID(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "sessions" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}