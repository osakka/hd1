package environments

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// EnvironmentInfo matches the OpenAPI EnvironmentInfo schema exactly
type EnvironmentInfo struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	ScaleUnit   string            `json:"scale_unit"`
	Gravity     float64           `json:"gravity"`
	Atmosphere  string            `json:"atmosphere"`
	Boundaries  *Boundaries       `json:"boundaries,omitempty"`
	Complexity  string            `json:"complexity,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
}

// Boundaries matches the OpenAPI boundaries schema
type Boundaries struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// ListEnvironmentsResponse matches the OpenAPI response schema
type ListEnvironmentsResponse struct {
	Success      bool              `json:"success"`
	Environments []EnvironmentInfo `json:"environments"`
}

// ListEnvironmentsHandler returns all available HD1 environments by scanning script directory
func ListEnvironmentsHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Dynamically discover environments from script directory
	environments := []EnvironmentInfo{}
	
	environmentsDir := "/opt/hd1/share/environments"
	
	// Scan environment scripts directory
	if files, err := os.ReadDir(environmentsDir); err == nil {
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".sh") {
				envID := strings.TrimSuffix(file.Name(), ".sh")
				
				// Parse environment script for metadata
				if envInfo := parseEnvironmentScript(filepath.Join(environmentsDir, file.Name())); envInfo != nil {
					envInfo.ID = envID
					environments = append(environments, *envInfo)
				}
			}
		}
	}
	
	// Fallback: Return default environments if directory doesn't exist
	if len(environments) == 0 {
		environments = getDefaultEnvironments()
	}
	
	response := ListEnvironmentsResponse{
		Success:      true,
		Environments: environments,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// parseEnvironmentScript extracts metadata from environment script comments
func parseEnvironmentScript(scriptPath string) *EnvironmentInfo {
	file, err := os.Open(scriptPath)
	if err != nil {
		return nil
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	env := &EnvironmentInfo{
		Boundaries: &Boundaries{Min: -12, Max: 12},
		Complexity: "simple",
		Tags:       []string{},
	}
	
	// Parse script header comments for metadata
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Stop parsing when we hit non-comment lines
		if !strings.HasPrefix(line, "#") && line != "" {
			break
		}
		
		// Extract metadata from comments
		if strings.Contains(line, "Name:") {
			env.Name = extractValue(line, "Name:")
		} else if strings.Contains(line, "Description:") {
			env.Description = extractValue(line, "Description:")
		} else if strings.Contains(line, "Scale:") {
			env.ScaleUnit = extractValue(line, "Scale:")
		} else if strings.Contains(line, "Gravity:") {
			if gravityStr := extractValue(line, "Gravity:"); gravityStr != "" {
				// Parse gravity value - simplified parsing
				switch gravityStr {
				case "9.8", "earth":
					env.Gravity = 9.8
				case "0.0", "zero", "space":
					env.Gravity = 0.0
				case "1.6", "moon":
					env.Gravity = 1.6
				default:
					env.Gravity = 9.8
				}
			}
		} else if strings.Contains(line, "Atmosphere:") {
			env.Atmosphere = extractValue(line, "Atmosphere:")
		} else if strings.Contains(line, "Complexity:") {
			env.Complexity = extractValue(line, "Complexity:")
		} else if strings.Contains(line, "Tags:") {
			tagsStr := extractValue(line, "Tags:")
			if tagsStr != "" {
				env.Tags = strings.Split(tagsStr, ",")
				for i := range env.Tags {
					env.Tags[i] = strings.TrimSpace(env.Tags[i])
				}
			}
		}
	}
	
	// Validate required fields per OpenAPI spec
	if env.Name == "" || env.Description == "" || env.ScaleUnit == "" || env.Atmosphere == "" {
		return nil
	}
	
	return env
}

// extractValue extracts value after "Key:" in comment line
func extractValue(line, key string) string {
	if idx := strings.Index(line, key); idx != -1 {
		value := strings.TrimSpace(line[idx+len(key):])
		return strings.Trim(value, "\"'")
	}
	return ""
}

// getDefaultEnvironments returns hardcoded environments matching OpenAPI examples
func getDefaultEnvironments() []EnvironmentInfo {
	return []EnvironmentInfo{
		{
			ID:          "earth-surface",
			Name:        "Earth Surface",
			Description: "Standard terrestrial environment",
			ScaleUnit:   "m",
			Gravity:     9.8,
			Atmosphere:  "air",
			Boundaries:  &Boundaries{Min: -12, Max: 12},
			Complexity:  "simple",
			Tags:        []string{"terrestrial", "standard"},
		},
		{
			ID:          "molecular-scale",
			Name:        "Molecular Scale",
			Description: "Nanometer-scale molecular environment",
			ScaleUnit:   "nm",
			Gravity:     0.0,
			Atmosphere:  "vacuum",
			Boundaries:  &Boundaries{Min: -12, Max: 12},
			Complexity:  "moderate",
			Tags:        []string{"molecular", "scientific"},
		},
		{
			ID:          "space-vacuum",
			Name:        "Space Vacuum",
			Description: "Zero gravity space environment",
			ScaleUnit:   "km",
			Gravity:     0.0,
			Atmosphere:  "vacuum",
			Boundaries:  &Boundaries{Min: -12, Max: 12},
			Complexity:  "simple",
			Tags:        []string{"space", "vacuum"},
		},
		{
			ID:          "underwater",
			Name:        "Underwater",
			Description: "Submerged liquid environment",
			ScaleUnit:   "m",
			Gravity:     8.8,
			Atmosphere:  "liquid",
			Boundaries:  &Boundaries{Min: -12, Max: 12},
			Complexity:  "moderate",
			Tags:        []string{"underwater", "liquid"},
		},
	}
}