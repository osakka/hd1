package scenes

import (
	"encoding/json"
	"net/http"
)

type SceneInfo struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ObjectCount int      `json:"object_count"`
	Complexity  string   `json:"complexity"`
	Tags        []string `json:"tags"`
}

type ListScenesResponse struct {
	Success bool        `json:"success"`
	Scenes  []SceneInfo `json:"scenes"`
}

// ListScenesHandler returns all available THD scenes for holodeck environments
func ListScenesHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Predefined scenes available in THD
	scenes := []SceneInfo{
		{
			ID:          "empty",
			Name:        "Empty Grid",
			Description: "Clean holodeck with just the coordinate grid system",
			ObjectCount: 0,
			Complexity:  "simple",
			Tags:        []string{"basic", "grid", "clean"},
		},
		{
			ID:          "anime-ui",
			Name:        "Anime UI Demo",
			Description: "Interactive anime-style holodeck with floating UI elements, blue lighting, and data visualization cubes",
			ObjectCount: 25,
			Complexity:  "moderate",
			Tags:        []string{"ui", "interactive", "demo", "anime", "lighting"},
		},
		{
			ID:          "ultimate",
			Name:        "Ultimate Demo",
			Description: "Complete holodeck showcase with 200+ objects: metallic platforms, crystal formations, particle effects, and cinematic lighting",
			ObjectCount: 200,
			Complexity:  "complex",
			Tags:        []string{"showcase", "complex", "particles", "lighting", "materials"},
		},
		{
			ID:          "basic-shapes",
			Name:        "Basic Shapes",
			Description: "Fundamental 3D shapes demonstration with various materials and colors",
			ObjectCount: 10,
			Complexity:  "simple",
			Tags:        []string{"basic", "shapes", "materials", "educational"},
		},
	}

	response := ListScenesResponse{
		Success: true,
		Scenes:  scenes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}