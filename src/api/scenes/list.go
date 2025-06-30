package scenes

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

// ListScenesHandler returns all available HD1 scenes by scanning script directory
func ListScenesHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Dynamically discover scenes from script directory
	scenes := []SceneInfo{}
	
	scenesDir := "/opt/holo-deck/share/scenes"
	if entries, err := os.ReadDir(scenesDir); err == nil {
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".sh") {
				sceneID := strings.TrimSuffix(entry.Name(), ".sh")
				scriptPath := filepath.Join(scenesDir, entry.Name())
				
				if sceneInfo := parseSceneScript(sceneID, scriptPath); sceneInfo != nil {
					scenes = append(scenes, *sceneInfo)
				}
			}
		}
	}
	
	// No hardcoded scenes - all scenes come from files

	response := ListScenesResponse{
		Success: true,
		Scenes:  scenes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// parseSceneScript extracts metadata from scene script files
func parseSceneScript(sceneID, scriptPath string) *SceneInfo {
	file, err := os.Open(scriptPath)
	if err != nil {
		return nil
	}
	defer file.Close()

	var name, description string
	var objectCount int
	var tags []string
	complexity := "simple" // default
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Parse scene name from comment header
		if strings.Contains(line, "THD Scene:") {
			parts := strings.Split(line, "THD Scene:")
			if len(parts) > 1 {
				namePart := strings.TrimSpace(parts[1])
				// Extract name before the dash
				if dashIndex := strings.Index(namePart, " - "); dashIndex > 0 {
					name = strings.TrimSpace(namePart[:dashIndex])
				} else {
					name = namePart
				}
			}
		}
		
		// Parse SCENE_NAME variable
		if strings.HasPrefix(line, "SCENE_NAME=") {
			name = extractQuotedValue(line, "SCENE_NAME=")
		}
		
		// Parse SCENE_DESCRIPTION variable
		if strings.HasPrefix(line, "SCENE_DESCRIPTION=") {
			description = extractQuotedValue(line, "SCENE_DESCRIPTION=")
		}
		
		// Count object creation commands
		if (strings.Contains(line, "hd1::create_object") || 
		    strings.Contains(line, "$HD1_CLIENT create-object")) && 
		   !strings.HasPrefix(strings.TrimSpace(line), "#") {
			objectCount++
		}
		
		// Continue parsing to count all object creation commands
	}
	
	// Determine complexity based on object count
	if objectCount == 0 {
		complexity = "simple"
	} else if objectCount <= 10 {
		complexity = "simple"
	} else if objectCount <= 50 {
		complexity = "moderate"
	} else {
		complexity = "complex"
	}
	
	// Generate tags based on scene characteristics
	tags = generateSceneTags(sceneID, name, description)
	
	// Use defaults if not found
	if name == "" {
		name = strings.Title(strings.ReplaceAll(sceneID, "-", " "))
	}
	if description == "" {
		description = "THD scene: " + name
	}
	
	return &SceneInfo{
		ID:          sceneID,
		Name:        name,
		Description: description,
		ObjectCount: objectCount,
		Complexity:  complexity,
		Tags:        tags,
	}
}

// extractQuotedValue extracts value from shell variable assignment
func extractQuotedValue(line, prefix string) string {
	if !strings.HasPrefix(line, prefix) {
		return ""
	}
	
	value := strings.TrimPrefix(line, prefix)
	value = strings.Trim(value, "\"'")
	return value
}

// generateSceneTags creates appropriate tags based on scene characteristics
func generateSceneTags(sceneID, name, description string) []string {
	tags := []string{}
	
	// Add tags based on scene ID
	if strings.Contains(sceneID, "basic") {
		tags = append(tags, "basic", "educational")
	}
	if strings.Contains(sceneID, "anime") {
		tags = append(tags, "anime", "ui", "interactive")
	}
	if strings.Contains(sceneID, "ultimate") {
		tags = append(tags, "showcase", "complex", "demo")
	}
	if strings.Contains(sceneID, "custom") {
		tags = append(tags, "custom", "user-created")
	}
	
	// Add tags based on description content
	lowerDesc := strings.ToLower(description)
	if strings.Contains(lowerDesc, "shape") {
		tags = append(tags, "shapes")
	}
	if strings.Contains(lowerDesc, "light") {
		tags = append(tags, "lighting")
	}
	if strings.Contains(lowerDesc, "material") {
		tags = append(tags, "materials")
	}
	if strings.Contains(lowerDesc, "particle") {
		tags = append(tags, "particles")
	}
	if strings.Contains(lowerDesc, "ui") || strings.Contains(lowerDesc, "interface") {
		tags = append(tags, "ui")
	}
	
	// Ensure we have at least one tag
	if len(tags) == 0 {
		tags = append(tags, "scene")
	}
	
	return tags
}