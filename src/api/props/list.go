package props

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	
	"gopkg.in/yaml.v3"
)

// PropInfo matches the OpenAPI PropInfo schema exactly
type PropInfo struct {
	ID               string             `json:"id" yaml:"id"`
	Name             string             `json:"name" yaml:"name"`
	Description      string             `json:"description" yaml:"description"`
	Category         string             `json:"category" yaml:"category"`
	ScaleCompatible  []string           `json:"scale_compatible" yaml:"scale_compatible"`
	Mass             float64            `json:"mass" yaml:"mass"`
	Material         string             `json:"material" yaml:"material"`
	Dimensions       PropDimensions     `json:"dimensions" yaml:"dimensions"`
	PhysicsProperties *PropPhysics      `json:"physics_properties,omitempty" yaml:"physics_properties,omitempty"`
	Complexity       string             `json:"complexity,omitempty" yaml:"complexity,omitempty"`
	Components       []string           `json:"components,omitempty" yaml:"components,omitempty"`
	Tags             []string           `json:"tags,omitempty" yaml:"tags,omitempty"`
	Script           string             `json:"-" yaml:"script"` // Internal use only
}

// PropDimensions represents the physical dimensions of a prop
type PropDimensions struct {
	Width  float64 `json:"width" yaml:"width"`
	Height float64 `json:"height" yaml:"height"`
	Depth  float64 `json:"depth" yaml:"depth"`
}

// PropPhysics represents physics properties of a prop
type PropPhysics struct {
	Friction     float64 `json:"friction,omitempty" yaml:"friction,omitempty"`
	Restitution  float64 `json:"restitution,omitempty" yaml:"restitution,omitempty"`
	Density      float64 `json:"density,omitempty" yaml:"density,omitempty"`
}

// ListPropsResponse matches the OpenAPI response schema
type ListPropsResponse struct {
	Success bool       `json:"success"`
	Props   []PropInfo `json:"props"`
}

// ListPropsHandler returns all available HD1 props by scanning prop directory
func ListPropsHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Dynamically discover props from prop directory
	props := []PropInfo{}
	
	propsDir := "/opt/hd1/share/props"
	
	// Scan all prop category directories
	err := filepath.Walk(propsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors, continue scanning
		}
		
		// Look for YAML prop definition files
		if strings.HasSuffix(info.Name(), ".yaml") {
			if propInfo := parsePropFile(path); propInfo != nil {
				props = append(props, *propInfo)
			}
		}
		
		return nil
	})
	
	if err != nil {
		// Fallback: Return default props if directory doesn't exist
		props = getDefaultProps()
	}
	
	// If no props found, return defaults
	if len(props) == 0 {
		props = getDefaultProps()
	}
	
	response := ListPropsResponse{
		Success: true,
		Props:   props,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// parsePropFile parses a YAML prop definition file
func parsePropFile(filePath string) *PropInfo {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil
	}
	
	var prop PropInfo
	if err := yaml.Unmarshal(data, &prop); err != nil {
		return nil
	}
	
	// Validate required fields per OpenAPI spec
	if prop.ID == "" || prop.Name == "" || prop.Description == "" || 
	   prop.Category == "" || len(prop.ScaleCompatible) == 0 || 
	   prop.Mass <= 0 || prop.Material == "" {
		return nil
	}
	
	return &prop
}

// getDefaultProps returns hardcoded props matching OpenAPI examples
func getDefaultProps() []PropInfo {
	return []PropInfo{
		{
			ID:              "wooden-chair",
			Name:            "Wooden Chair",
			Description:     "Standard wooden chair with realistic physics",
			Category:        "furniture",
			ScaleCompatible: []string{"mm", "cm", "m"},
			Mass:            5.5,
			Material:        "wood",
			Dimensions: PropDimensions{
				Width:  0.6,
				Height: 0.8,
				Depth:  0.6,
			},
			PhysicsProperties: &PropPhysics{
				Friction:    0.7,
				Restitution: 0.3,
				Density:     600.0,
			},
			Complexity: "moderate",
			Components: []string{"seat", "backrest", "legs", "armrests"},
			Tags:       []string{"seating", "wooden", "indoor"},
		},
		{
			ID:              "wooden-table",
			Name:            "Wooden Table",
			Description:     "Standard wooden table with four legs and solid top",
			Category:        "furniture",
			ScaleCompatible: []string{"cm", "m"},
			Mass:            12.0,
			Material:        "wood",
			Dimensions: PropDimensions{
				Width:  1.2,
				Height: 0.75,
				Depth:  0.8,
			},
			PhysicsProperties: &PropPhysics{
				Friction:    0.8,
				Restitution: 0.2,
				Density:     650.0,
			},
			Complexity: "simple",
			Components: []string{"table_top", "leg_fl", "leg_fr", "leg_bl", "leg_br"},
			Tags:       []string{"table", "wooden", "indoor", "furniture", "workspace"},
		},
		{
			ID:              "hammer",
			Name:            "Hammer",
			Description:     "Standard claw hammer with steel head and wooden handle",
			Category:        "tools",
			ScaleCompatible: []string{"mm", "cm", "m"},
			Mass:            0.6,
			Material:        "metal",
			Dimensions: PropDimensions{
				Width:  0.05,
				Height: 0.32,
				Depth:  0.05,
			},
			PhysicsProperties: &PropPhysics{
				Friction:    0.9,
				Restitution: 0.7,
				Density:     7800.0,
			},
			Complexity: "simple",
			Components: []string{"head", "handle", "grip"},
			Tags:       []string{"tool", "metal", "construction", "handheld"},
		},
	}
}