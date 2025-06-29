package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	
	"holodeck/logging"
)

// AFrameSchemaReader extracts component schemas from A-Frame vendor files
type AFrameSchemaReader struct {
	VendorPath     string
	ComponentSpecs map[string]*AFrameComponentSpec
}

// AFrameComponentSpec represents a complete A-Frame component specification
type AFrameComponentSpec struct {
	Name        string                    `json:"name"`
	Schema      map[string]*PropertySpec  `json:"schema"`
	Methods     []string                  `json:"methods"`
	Events      []string                  `json:"events"`
	Description string                    `json:"description"`
	Category    string                    `json:"category"` // geometry, material, light, etc.
}

// PropertySpec defines A-Frame component property schema
type PropertySpec struct {
	Type        string      `json:"type"`
	Default     interface{} `json:"default"`
	Min         *float64    `json:"min,omitempty"`
	Max         *float64    `json:"max,omitempty"`
	OneOf       []string    `json:"oneOf,omitempty"`
	Description string      `json:"description"`
}

// NewAFrameSchemaReader creates schema reader for downstream A-Frame dependency
func NewAFrameSchemaReader(vendorPath string) *AFrameSchemaReader {
	return &AFrameSchemaReader{
		VendorPath:     vendorPath,
		ComponentSpecs: make(map[string]*AFrameComponentSpec),
	}
}

// ExtractSchemas reads A-Frame component definitions from vendor files
func (r *AFrameSchemaReader) ExtractSchemas() error {
	logging.Info("extracting A-Frame schemas from downstream dependency", map[string]interface{}{
		"vendor_path": r.VendorPath,
		"task": "downstream-api-integration",
	})

	// Extract core A-Frame primitive schemas
	if err := r.extractCoreSchemas(); err != nil {
		return fmt.Errorf("failed to extract core A-Frame schemas: %w", err)
	}

	// Extract component library schemas
	if err := r.extractComponentSchemas(); err != nil {
		return fmt.Errorf("failed to extract component schemas: %w", err)
	}

	// Extract physics system schemas
	if err := r.extractPhysicsSchemas(); err != nil {
		return fmt.Errorf("failed to extract physics schemas: %w", err)
	}

	logging.Info("A-Frame schema extraction completed", map[string]interface{}{
		"total_components": len(r.ComponentSpecs),
		"categories": r.getCategoryStats(),
	})

	return nil
}

// extractCoreSchemas extracts A-Frame primitive schemas (geometry, material, light)
func (r *AFrameSchemaReader) extractCoreSchemas() error {
	// A-Frame core primitive schemas (based on A-Frame 1.4.0 documentation)
	
	// Geometry primitives
	r.ComponentSpecs["geometry-box"] = &AFrameComponentSpec{
		Name:        "geometry-box",
		Category:    "geometry",
		Description: "Box geometry primitive",
		Schema: map[string]*PropertySpec{
			"width":  {Type: "number", Default: 1.0, Min: &[]float64{0}[0], Description: "Width of the box"},
			"height": {Type: "number", Default: 1.0, Min: &[]float64{0}[0], Description: "Height of the box"},
			"depth":  {Type: "number", Default: 1.0, Min: &[]float64{0}[0], Description: "Depth of the box"},
		},
	}

	r.ComponentSpecs["geometry-sphere"] = &AFrameComponentSpec{
		Name:        "geometry-sphere",
		Category:    "geometry",
		Description: "Sphere geometry primitive",
		Schema: map[string]*PropertySpec{
			"radius":         {Type: "number", Default: 1.0, Min: &[]float64{0}[0], Description: "Radius of the sphere"},
			"segmentsWidth":  {Type: "int", Default: 18, Min: &[]float64{3}[0], Description: "Number of horizontal segments"},
			"segmentsHeight": {Type: "int", Default: 36, Min: &[]float64{2}[0], Description: "Number of vertical segments"},
		},
	}

	r.ComponentSpecs["geometry-cylinder"] = &AFrameComponentSpec{
		Name:        "geometry-cylinder",
		Category:    "geometry", 
		Description: "Cylinder geometry primitive",
		Schema: map[string]*PropertySpec{
			"radius":       {Type: "number", Default: 1.0, Min: &[]float64{0}[0], Description: "Radius of the cylinder"},
			"height":       {Type: "number", Default: 1.0, Min: &[]float64{0}[0], Description: "Height of the cylinder"},
			"segmentsRadial": {Type: "int", Default: 36, Min: &[]float64{3}[0], Description: "Number of radial segments"},
		},
	}

	r.ComponentSpecs["geometry-cone"] = &AFrameComponentSpec{
		Name:        "geometry-cone",
		Category:    "geometry",
		Description: "Cone geometry primitive", 
		Schema: map[string]*PropertySpec{
			"radius":       {Type: "number", Default: 1.0, Min: &[]float64{0}[0], Description: "Radius of the cone base"},
			"height":       {Type: "number", Default: 1.0, Min: &[]float64{0}[0], Description: "Height of the cone"},
			"segmentsRadial": {Type: "int", Default: 36, Min: &[]float64{3}[0], Description: "Number of radial segments"},
		},
	}

	r.ComponentSpecs["geometry-plane"] = &AFrameComponentSpec{
		Name:        "geometry-plane",
		Category:    "geometry",
		Description: "Plane geometry primitive",
		Schema: map[string]*PropertySpec{
			"width":  {Type: "number", Default: 1.0, Min: &[]float64{0}[0], Description: "Width of the plane"},
			"height": {Type: "number", Default: 1.0, Min: &[]float64{0}[0], Description: "Height of the plane"},
		},
	}

	// Material properties
	r.ComponentSpecs["material-standard"] = &AFrameComponentSpec{
		Name:        "material-standard",
		Category:    "material",
		Description: "Standard PBR material",
		Schema: map[string]*PropertySpec{
			"color":       {Type: "color", Default: "#ffffff", Description: "Base color of the material"},
			"metalness":   {Type: "number", Default: 0.0, Min: &[]float64{0}[0], Max: &[]float64{1}[0], Description: "Metallic property"},
			"roughness":   {Type: "number", Default: 0.5, Min: &[]float64{0}[0], Max: &[]float64{1}[0], Description: "Surface roughness"},
			"transparent": {Type: "boolean", Default: false, Description: "Enable transparency"},
			"opacity":     {Type: "number", Default: 1.0, Min: &[]float64{0}[0], Max: &[]float64{1}[0], Description: "Material opacity"},
			"emissive":    {Type: "color", Default: "#000000", Description: "Emissive color"},
			"emissiveIntensity": {Type: "number", Default: 1.0, Min: &[]float64{0}[0], Description: "Emissive intensity"},
		},
	}

	// Light types
	r.ComponentSpecs["light-directional"] = &AFrameComponentSpec{
		Name:        "light-directional",
		Category:    "light",
		Description: "Directional light source",
		Schema: map[string]*PropertySpec{
			"color":     {Type: "color", Default: "#ffffff", Description: "Light color"},
			"intensity": {Type: "number", Default: 1.0, Min: &[]float64{0}[0], Description: "Light intensity"},
			"castShadow": {Type: "boolean", Default: false, Description: "Enable shadow casting"},
		},
	}

	r.ComponentSpecs["light-point"] = &AFrameComponentSpec{
		Name:        "light-point",
		Category:    "light",
		Description: "Point light source",
		Schema: map[string]*PropertySpec{
			"color":     {Type: "color", Default: "#ffffff", Description: "Light color"},
			"intensity": {Type: "number", Default: 1.0, Min: &[]float64{0}[0], Description: "Light intensity"},
			"distance":  {Type: "number", Default: 0.0, Min: &[]float64{0}[0], Description: "Light attenuation distance"},
			"decay":     {Type: "number", Default: 1.0, Min: &[]float64{0}[0], Description: "Light decay rate"},
		},
	}

	r.ComponentSpecs["light-ambient"] = &AFrameComponentSpec{
		Name:        "light-ambient",
		Category:    "light",
		Description: "Ambient light source",
		Schema: map[string]*PropertySpec{
			"color":     {Type: "color", Default: "#ffffff", Description: "Ambient light color"},
			"intensity": {Type: "number", Default: 1.0, Min: &[]float64{0}[0], Description: "Ambient light intensity"},
		},
	}

	logging.Info("core A-Frame schemas extracted", map[string]interface{}{
		"geometry_primitives": 5,
		"material_types": 1,
		"light_types": 3,
	})

	return nil
}

// extractComponentSchemas extracts schemas from A-Frame component library files
func (r *AFrameSchemaReader) extractComponentSchemas() error {
	componentFiles := []string{
		"aframe-environment-component.min.js",
		"aframe-physics-system.min.js", 
		"aframe-animation-component.min.js",
		"aframe-particle-system.js",
	}

	for _, file := range componentFiles {
		filePath := filepath.Join(r.VendorPath, file)
		if err := r.extractComponentFromFile(filePath); err != nil {
			logging.Warn("failed to extract schema from component file", map[string]interface{}{
				"file": file,
				"error": err.Error(),
			})
			continue
		}
	}

	return nil
}

// extractComponentFromFile extracts component schema from JavaScript file
func (r *AFrameSchemaReader) extractComponentFromFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("component file not found: %s", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read component file: %w", err)
	}

	// Extract component registrations and schemas using regex patterns
	return r.parseComponentDefinitions(string(content), filepath.Base(filePath))
}

// parseComponentDefinitions extracts component schemas from JavaScript content
func (r *AFrameSchemaReader) parseComponentDefinitions(content, filename string) error {
	// Pattern to match AFRAME.registerComponent calls
	componentPattern := regexp.MustCompile(`AFRAME\.registerComponent\s*\(\s*['"]([^'"]+)['"]\s*,\s*{([^}]+(?:{[^}]*}[^}]*)*)}`)
	schemaPattern := regexp.MustCompile(`schema\s*:\s*{([^}]+(?:{[^}]*}[^}]*)*)}`)

	matches := componentPattern.FindAllStringSubmatch(content, -1)
	
	for _, match := range matches {
		if len(match) >= 3 {
			componentName := match[1]
			componentBody := match[2]
			
			// Extract schema from component body
			schemaMatches := schemaPattern.FindStringSubmatch(componentBody)
			if len(schemaMatches) >= 2 {
				schema := r.parseSchemaDefinition(schemaMatches[1])
				
				r.ComponentSpecs[componentName] = &AFrameComponentSpec{
					Name:        componentName,
					Category:    r.inferCategory(componentName),
					Description: fmt.Sprintf("Component from %s", filename),
					Schema:      schema,
				}
				
				logging.Debug("extracted component schema", map[string]interface{}{
					"component": componentName,
					"file": filename,
					"properties": len(schema),
				})
			}
		}
	}

	return nil
}

// parseSchemaDefinition parses A-Frame schema definition syntax
func (r *AFrameSchemaReader) parseSchemaDefinition(schemaContent string) map[string]*PropertySpec {
	schema := make(map[string]*PropertySpec)
	
	// Pattern to match property definitions
	propertyPattern := regexp.MustCompile(`(\w+)\s*:\s*{([^}]+)}`)
	matches := propertyPattern.FindAllStringSubmatch(schemaContent, -1)
	
	for _, match := range matches {
		if len(match) >= 3 {
			propertyName := match[1]
			propertyDef := match[2]
			
			prop := &PropertySpec{}
			
			// Extract type
			if typeMatch := regexp.MustCompile(`type\s*:\s*['"]([^'"]+)['"]`).FindStringSubmatch(propertyDef); len(typeMatch) >= 2 {
				prop.Type = typeMatch[1]
			}
			
			// Extract default value
			if defaultMatch := regexp.MustCompile(`default\s*:\s*([^,}]+)`).FindStringSubmatch(propertyDef); len(defaultMatch) >= 2 {
				prop.Default = strings.TrimSpace(defaultMatch[1])
			}
			
			schema[propertyName] = prop
		}
	}
	
	return schema
}

// extractPhysicsSchemas extracts physics system component schemas
func (r *AFrameSchemaReader) extractPhysicsSchemas() error {
	// Physics body types from aframe-physics-system
	r.ComponentSpecs["physics-dynamic"] = &AFrameComponentSpec{
		Name:        "physics-dynamic",
		Category:    "physics",
		Description: "Dynamic physics body",
		Schema: map[string]*PropertySpec{
			"mass":        {Type: "number", Default: 1.0, Min: &[]float64{0}[0], Description: "Mass of the body"},
			"linearDamping":  {Type: "number", Default: 0.01, Min: &[]float64{0}[0], Max: &[]float64{1}[0], Description: "Linear damping"},
			"angularDamping": {Type: "number", Default: 0.01, Min: &[]float64{0}[0], Max: &[]float64{1}[0], Description: "Angular damping"},
		},
	}

	r.ComponentSpecs["physics-static"] = &AFrameComponentSpec{
		Name:        "physics-static",
		Category:    "physics",
		Description: "Static physics body",
		Schema: map[string]*PropertySpec{
			"type": {Type: "string", Default: "static", OneOf: []string{"static"}, Description: "Physics body type"},
		},
	}

	r.ComponentSpecs["physics-kinematic"] = &AFrameComponentSpec{
		Name:        "physics-kinematic",
		Category:    "physics", 
		Description: "Kinematic physics body",
		Schema: map[string]*PropertySpec{
			"type": {Type: "string", Default: "kinematic", OneOf: []string{"kinematic"}, Description: "Physics body type"},
		},
	}

	return nil
}

// GetComponentsByCategory returns components grouped by category
func (r *AFrameSchemaReader) GetComponentsByCategory() map[string][]*AFrameComponentSpec {
	categories := make(map[string][]*AFrameComponentSpec)
	
	for _, spec := range r.ComponentSpecs {
		categories[spec.Category] = append(categories[spec.Category], spec)
	}
	
	return categories
}

// GetAvailableCapabilities returns summary of A-Frame capabilities
func (r *AFrameSchemaReader) GetAvailableCapabilities() map[string]interface{} {
	categories := r.GetComponentsByCategory()
	
	capabilities := map[string]interface{}{
		"total_components": len(r.ComponentSpecs),
		"categories": make(map[string]int),
		"geometry_types": []string{},
		"material_properties": []string{},
		"light_types": []string{},
		"physics_bodies": []string{},
	}
	
	for category, specs := range categories {
		capabilities["categories"].(map[string]int)[category] = len(specs)
		
		for _, spec := range specs {
			switch category {
			case "geometry":
				capabilities["geometry_types"] = append(capabilities["geometry_types"].([]string), strings.TrimPrefix(spec.Name, "geometry-"))
			case "material":
				for prop := range spec.Schema {
					capabilities["material_properties"] = append(capabilities["material_properties"].([]string), prop)
				}
			case "light":
				capabilities["light_types"] = append(capabilities["light_types"].([]string), strings.TrimPrefix(spec.Name, "light-"))
			case "physics":
				capabilities["physics_bodies"] = append(capabilities["physics_bodies"].([]string), strings.TrimPrefix(spec.Name, "physics-"))
			}
		}
	}
	
	return capabilities
}

// inferCategory infers component category from name
func (r *AFrameSchemaReader) inferCategory(name string) string {
	lowerName := strings.ToLower(name)
	
	if strings.Contains(lowerName, "geometry") || strings.Contains(lowerName, "box") || strings.Contains(lowerName, "sphere") {
		return "geometry"
	}
	if strings.Contains(lowerName, "material") || strings.Contains(lowerName, "shader") {
		return "material"
	}
	if strings.Contains(lowerName, "light") || strings.Contains(lowerName, "shadow") {
		return "light"
	}
	if strings.Contains(lowerName, "physics") || strings.Contains(lowerName, "body") {
		return "physics"
	}
	if strings.Contains(lowerName, "animation") || strings.Contains(lowerName, "tween") {
		return "animation"
	}
	if strings.Contains(lowerName, "particle") || strings.Contains(lowerName, "effect") {
		return "effects"
	}
	if strings.Contains(lowerName, "environment") || strings.Contains(lowerName, "sky") {
		return "environment"
	}
	
	return "utility"
}

// getCategoryStats returns statistics about extracted schemas
func (r *AFrameSchemaReader) getCategoryStats() map[string]int {
	stats := make(map[string]int)
	for _, spec := range r.ComponentSpecs {
		stats[spec.Category]++
	}
	return stats
}

// SaveSchemas exports extracted schemas to JSON file
func (r *AFrameSchemaReader) SaveSchemas(outputPath string) error {
	data := map[string]interface{}{
		"metadata": map[string]interface{}{
			"extracted_at": "auto-generated",
			"total_components": len(r.ComponentSpecs),
			"categories": r.getCategoryStats(),
		},
		"capabilities": r.GetAvailableCapabilities(),
		"components": r.ComponentSpecs,
	}
	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal schemas to JSON: %w", err)
	}
	
	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write schemas to file: %w", err)
	}
	
	logging.Info("A-Frame schemas saved", map[string]interface{}{
		"output_path": outputPath,
		"file_size_kb": len(jsonData) / 1024,
	})
	
	return nil
}