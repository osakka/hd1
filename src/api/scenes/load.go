package scenes

import (
	"encoding/json"
	"net/http"
	"holodeck/server"
	"fmt"
	"log"
	"strings"
)

type LoadSceneRequest struct {
	SessionID string `json:"session_id"`
}

type LoadSceneResponse struct {
	Success        bool   `json:"success"`
	SceneID        string `json:"scene_id"`
	SessionID      string `json:"session_id"`
	ObjectsCreated int    `json:"objects_created"`
	Message        string `json:"message"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

// LoadSceneHandler loads a predefined scene into the specified session
func LoadSceneHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type
	h, ok := hub.(*server.Hub)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "internal_error",
			Message: "Internal server error",
		})
		return
	}

	// Extract scene ID from URL path
	sceneID := extractSceneID(r.URL.Path)
	if sceneID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "missing_scene_id",
			Message: "Scene ID is required",
		})
		return
	}

	// Parse request body
	var req LoadSceneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "invalid_request",
			Message: "Invalid JSON request body",
		})
		return
	}

	if req.SessionID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "missing_session_id",
			Message: "Session ID is required",
		})
		return
	}

	// Verify session exists
	if _, exists := h.GetStore().GetSession(req.SessionID); !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "session_not_found",
			Message: "Session does not exist",
		})
		return
	}

	// Clear existing objects in session
	clearAllObjects(req.SessionID, h)
	
	// Send clear command via WebSocket
	h.BroadcastToSession(req.SessionID, "canvas_control", map[string]interface{}{
		"command": "clear",
		"clear":   true,
	})

	var objectsCreated int
	var message string

	// Load scene based on ID
	switch sceneID {
	case "empty":
		message = "Empty grid scene loaded"
		objectsCreated = 0
		
	case "anime-ui":
		objectsCreated = loadAnimeUIScene(req.SessionID, h)
		message = fmt.Sprintf("Anime UI scene loaded with %d objects", objectsCreated)
		
	case "ultimate":
		objectsCreated = loadUltimateScene(req.SessionID, h)
		message = fmt.Sprintf("Ultimate demo scene loaded with %d objects", objectsCreated)
		
	case "basic-shapes":
		objectsCreated = loadBasicShapesScene(req.SessionID, h)
		message = fmt.Sprintf("Basic shapes scene loaded with %d objects", objectsCreated)
		
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "scene_not_found",
			Message: "Scene '" + sceneID + "' not found",
		})
		return
	}

	log.Printf("[THD] Scene '%s' loaded into session '%s' with %d objects", sceneID, req.SessionID, objectsCreated)

	response := LoadSceneResponse{
		Success:        true,
		SceneID:        sceneID,
		SessionID:      req.SessionID,
		ObjectsCreated: objectsCreated,
		Message:        message,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// extractSceneID extracts scene ID from URL path
func extractSceneID(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "scenes" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// clearAllObjects removes all objects from a session
func clearAllObjects(sessionID string, h *server.Hub) {
	objects := h.GetStore().ListObjects(sessionID)
	for _, obj := range objects {
		h.GetStore().DeleteObject(sessionID, obj.Name)
	}
}

// Scene creation functions
func loadAnimeUIScene(sessionID string, h *server.Hub) int {
	objects := []map[string]interface{}{
		// Central anime ring interface
		{
			"name": "central_ring",
			"type": "cylinder",
			"x": 0, "y": 2, "z": 0,
			"scale": 3,
			"color": map[string]float64{"r": 0.2, "g": 0.8, "b": 1.0, "a": 0.3},
			"wireframe": true,
		},
		// Floating UI cubes
		{
			"name": "ui_cube_1",
			"type": "cube",
			"x": -4, "y": 3, "z": 2,
			"scale": 0.8,
			"color": map[string]float64{"r": 1.0, "g": 0.4, "b": 0.8, "a": 0.8},
		},
		{
			"name": "ui_cube_2",
			"type": "cube",
			"x": 4, "y": 3.5, "z": -2,
			"scale": 0.6,
			"color": map[string]float64{"r": 0.4, "g": 1.0, "b": 0.8, "a": 0.8},
		},
		// Data visualization spheres
		{
			"name": "data_sphere_1",
			"type": "sphere",
			"x": -2, "y": 4, "z": -4,
			"scale": 0.5,
			"color": map[string]float64{"r": 0.8, "g": 0.2, "b": 1.0, "a": 0.9},
		},
		{
			"name": "data_sphere_2",
			"type": "sphere",
			"x": 3, "y": 2.5, "z": 4,
			"scale": 0.7,
			"color": map[string]float64{"r": 1.0, "g": 0.8, "b": 0.2, "a": 0.9},
		},
		// Info panel (replacing text)
		{
			"name": "info_panel",
			"type": "plane",
			"x": 0, "y": 5, "z": -3,
			"scale": 2,
			"color": map[string]float64{"r": 0.2, "g": 1.0, "b": 1.0, "a": 0.3},
		},
	}

	// Send objects via WebSocket
	h.BroadcastToSession(sessionID, "canvas_control", map[string]interface{}{
		"command": "create",
		"objects": objects,
	})
	
	return len(objects)
}

func loadUltimateScene(sessionID string, h *server.Hub) int {
	objects := []map[string]interface{}{
		// Sky environment
		{
			"name": "holodeck_sky",
			"type": "sky",
			"x": 0, "y": 0, "z": 0,
			"color": map[string]float64{"r": 0.1, "g": 0.2, "b": 0.4, "a": 1.0},
		},
		// Central metallic platform
		{
			"name": "central_platform",
			"type": "cylinder",
			"x": 0, "y": 0.2, "z": 0,
			"scale": 4,
			"color": map[string]float64{"r": 0.7, "g": 0.7, "b": 0.8, "a": 1.0},
			"material": map[string]interface{}{
				"metalness": 0.8,
				"roughness": 0.2,
			},
		},
		// Crystal formations
		{
			"name": "crystal_1",
			"type": "cone",
			"x": -3, "y": 2, "z": -3,
			"scale": 1.5,
			"color": map[string]float64{"r": 0.8, "g": 0.2, "b": 0.8, "a": 0.7},
			"material": map[string]interface{}{
				"transparent": true,
				"metalness": 0.1,
				"roughness": 0.1,
			},
		},
		{
			"name": "crystal_2",
			"type": "cone",
			"x": 3, "y": 1.8, "z": 3,
			"scale": 1.2,
			"color": map[string]float64{"r": 0.2, "g": 0.8, "b": 0.8, "a": 0.7},
			"material": map[string]interface{}{
				"transparent": true,
				"metalness": 0.1,
				"roughness": 0.1,
			},
		},
		// Additional metallic structures (replacing particle effects)
		{
			"name": "metallic_pillar_1",
			"type": "cylinder",
			"x": -5, "y": 2, "z": 0,
			"scale": 0.5,
			"color": map[string]float64{"r": 0.8, "g": 0.3, "b": 0.1, "a": 1.0},
			"material": map[string]interface{}{
				"metalness": 0.9,
				"roughness": 0.1,
			},
		},
		{
			"name": "energy_sphere",
			"type": "sphere",
			"x": 5, "y": 2, "z": 0,
			"scale": 0.8,
			"color": map[string]float64{"r": 1.0, "g": 1.0, "b": 0.2, "a": 0.8},
			"material": map[string]interface{}{
				"transparent": true,
				"metalness": 0.1,
				"roughness": 0.1,
			},
		},
		// Professional lighting
		{
			"name": "main_light",
			"type": "light",
			"x": 5, "y": 8, "z": 5,
			"lightType": "directional",
			"intensity": 1.2,
			"color": map[string]float64{"r": 1.0, "g": 1.0, "b": 0.9, "a": 1.0},
		},
		{
			"name": "accent_light",
			"type": "light",
			"x": -3, "y": 6, "z": -3,
			"lightType": "point",
			"intensity": 0.8,
			"color": map[string]float64{"r": 0.2, "g": 0.8, "b": 1.0, "a": 1.0},
		},
		// Status display panel
		{
			"name": "status_display",
			"type": "plane",
			"x": 0, "y": 6, "z": -5,
			"scale": 3,
			"color": map[string]float64{"r": 1.0, "g": 1.0, "b": 1.0, "a": 0.2},
		},
	}

	// Send objects via WebSocket
	h.BroadcastToSession(sessionID, "canvas_control", map[string]interface{}{
		"command": "create",
		"objects": objects,
	})
	
	return len(objects)
}

func loadBasicShapesScene(sessionID string, h *server.Hub) int {
	objects := []map[string]interface{}{
		// Basic cube
		{
			"name": "demo_cube",
			"type": "cube",
			"x": -3, "y": 1, "z": 0,
			"scale": 1,
			"color": map[string]float64{"r": 1.0, "g": 0.2, "b": 0.2, "a": 1.0},
		},
		// Basic sphere
		{
			"name": "demo_sphere",
			"type": "sphere",
			"x": 0, "y": 1, "z": 0,
			"scale": 1,
			"color": map[string]float64{"r": 0.2, "g": 1.0, "b": 0.2, "a": 1.0},
		},
		// Basic cylinder
		{
			"name": "demo_cylinder",
			"type": "cylinder",
			"x": 3, "y": 1, "z": 0,
			"scale": 1,
			"color": map[string]float64{"r": 0.2, "g": 0.2, "b": 1.0, "a": 1.0},
		},
		// Basic cone
		{
			"name": "demo_cone",
			"type": "cone",
			"x": -1.5, "y": 1, "z": -3,
			"scale": 1,
			"color": map[string]float64{"r": 1.0, "g": 1.0, "b": 0.2, "a": 1.0},
		},
		// Wireframe cube
		{
			"name": "wireframe_cube",
			"type": "cube",
			"x": 1.5, "y": 1, "z": -3,
			"scale": 1,
			"color": map[string]float64{"r": 0.8, "g": 0.8, "b": 0.8, "a": 1.0},
			"wireframe": true,
		},
		// Label panel
		{
			"name": "shapes_label",
			"type": "plane",
			"x": 0, "y": 3, "z": -1,
			"scale": 2,
			"color": map[string]float64{"r": 1.0, "g": 1.0, "b": 1.0, "a": 0.3},
		},
	}

	// Send objects via WebSocket
	h.BroadcastToSession(sessionID, "canvas_control", map[string]interface{}{
		"command": "create",
		"objects": objects,
	})
	
	return len(objects)
}