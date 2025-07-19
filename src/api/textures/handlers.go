package textures

import (
	"encoding/json"
	"net/http"
	"time"

	"holodeck1/api/shared"
	"holodeck1/logging"
	"holodeck1/sync"
)

// LoadTexture handles POST /textures/load
func LoadTexture(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	textureID := generateTextureID()
	clientID := shared.GetClientID(r)

	// Create sync operation
	operation := &sync.Operation{
		ClientID:  clientID,
		Type:      "texture_load",
		Data: map[string]interface{}{
			"texture_id": textureID,
			"url":        getString(req, "url", ""),
			"wrapS":      getString(req, "wrapS", "ClampToEdgeWrapping"),
			"wrapT":      getString(req, "wrapT", "ClampToEdgeWrapping"),
			"magFilter":  getString(req, "magFilter", "LinearFilter"),
			"minFilter":  getString(req, "minFilter", "LinearMipmapLinearFilter"),
			"flipY":      getBool(req, "flipY", true),
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

	logging.Info("Texture loaded", map[string]interface{}{
		"url":        getString(req, "url", ""),
		"wrapS":      getString(req, "wrapS", "ClampToEdgeWrapping"),
		"wrapT":      getString(req, "wrapT", "ClampToEdgeWrapping"),
		"texture_id": textureID,
		"seq_num":    seqNum,
		"endpoint":   "/textures/load",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"texture_id": textureID,
		"seq_num":    seqNum,
	})
}

// CreateProceduralTexture handles POST /textures/create
func CreateProceduralTexture(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	textureID := generateTextureID()
	clientID := shared.GetClientID(r)

	// Create sync operation
	operation := &sync.Operation{
		ClientID:  clientID,
		Type:      "texture_create",
		Data: map[string]interface{}{
			"texture_id": textureID,
			"type":       getString(req, "type", "canvas"),
			"width":      getInt(req, "width", 512),
			"height":     getInt(req, "height", 512),
			"pattern":    getString(req, "pattern", "checkerboard"),
			"color1":     getString(req, "color1", "#ffffff"),
			"color2":     getString(req, "color2", "#000000"),
			"scale":      getFloat(req, "scale", 1.0),
			"rotation":   getFloat(req, "rotation", 0.0),
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

	logging.Info("Procedural texture created", map[string]interface{}{
		"type":       getString(req, "type", "canvas"),
		"width":      getInt(req, "width", 512),
		"height":     getInt(req, "height", 512),
		"pattern":    getString(req, "pattern", "checkerboard"),
		"color1":     getString(req, "color1", "#ffffff"),
		"color2":     getString(req, "color2", "#000000"),
		"texture_id": textureID,
		"seq_num":    seqNum,
		"endpoint":   "/textures/create",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"texture_id": textureID,
		"seq_num":    seqNum,
	})
}

// Helper functions
func getInt(req map[string]interface{}, key string, defaultValue int) int {
	if val, ok := req[key]; ok {
		if f, ok := val.(float64); ok {
			return int(f)
		}
	}
	return defaultValue
}

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

func generateTextureID() string {
	return "texture_" + "generated"
}