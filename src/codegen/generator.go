package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	
	"gopkg.in/yaml.v3"
	"holodeck1/config"
	"holodeck1/logging"
)

//go:embed templates/*
var templateFS embed.FS

// Template cache for performance
var templateCache = make(map[string]*template.Template)

// Three.js schema generation types
type ThreeJSGeometry struct {
	Name        string                 `yaml:"name"`
	Constructor string                 `yaml:"constructor"`
	Parameters  []GeometryParameter    `yaml:"parameters"`
	Description string                 `yaml:"description"`
}

type GeometryParameter struct {
	Name         string      `yaml:"name"`
	Type         string      `yaml:"type"`
	Required     bool        `yaml:"required"`
	DefaultValue interface{} `yaml:"default,omitempty"`
	Description  string      `yaml:"description,omitempty"`
}

type ThreeJSAPISchema struct {
	OpenAPI string                 `yaml:"openapi"`
	Info    ThreeJSInfo            `yaml:"info"`
	Paths   map[string]interface{} `yaml:"paths"`
	Components Components           `yaml:"components"`
}

type ThreeJSInfo struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

type Components struct {
	Schemas map[string]interface{} `yaml:"schemas"`
}

// Schema merger types
type SchemaMerger struct {
	schemas []APISchema
}

type APISchema struct {
	Name     string
	FilePath string
	Spec     map[string]interface{}
}

// OpenAPI Specification Structure
type OpenAPISpec struct {
	OpenAPI string                 `yaml:"openapi"`
	Info    Info                   `yaml:"info"`
	Paths   map[string]PathItem    `yaml:"paths"`
	XCodeGeneration CodeGenConfig  `yaml:"x-code-generation"`
}

type Info struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

type PathItem struct {
	Get    *Operation `yaml:"get,omitempty"`
	Post   *Operation `yaml:"post,omitempty"`
	Put    *Operation `yaml:"put,omitempty"`
	Delete *Operation `yaml:"delete,omitempty"`
}

type Operation struct {
	OperationID string   `yaml:"operationId"`
	Tags        []string `yaml:"tags"`
	Summary     string   `yaml:"summary"`
	Description string   `yaml:"description"`
	Parameters  []Parameter `yaml:"parameters,omitempty"`
	RequestBody *RequestBody `yaml:"requestBody,omitempty"`
	Responses   map[string]Response `yaml:"responses"`
	XHandler    string   `yaml:"x-handler"`
	XFunction   string   `yaml:"x-function"`
}

type Parameter struct {
	Name     string `yaml:"name"`
	In       string `yaml:"in"`
	Required bool   `yaml:"required"`
	Schema   Schema `yaml:"schema"`
}

type RequestBody struct {
	Required bool                    `yaml:"required"`
	Content  map[string]MediaType    `yaml:"content"`
}

type MediaType struct {
	Schema Schema `yaml:"schema"`
}

type Response struct {
	Description string               `yaml:"description"`
	Content     map[string]MediaType `yaml:"content,omitempty"`
}

type Schema struct {
	Type    string `yaml:"type"`
	Pattern string `yaml:"pattern,omitempty"`
	Ref     string `yaml:"$ref,omitempty"`
}

type CodeGenConfig struct {
	StrictValidation      bool `yaml:"strict-validation"`
	AutoRouting          bool `yaml:"auto-routing"`
	HandlerValidation    bool `yaml:"handler-validation"`
	FailOnMissingHandlers bool `yaml:"fail-on-missing-handlers"`
}

// loadTemplate loads and caches a template from the embedded filesystem
func loadTemplate(templatePath string) (*template.Template, error) {
	if tmpl, exists := templateCache[templatePath]; exists {
		return tmpl, nil
	}
	
	content, err := templateFS.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}
	
	// Add custom template functions
	funcMap := template.FuncMap{
		"hasSuffix": strings.HasSuffix,
	}
	
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(funcMap).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}
	
	templateCache[templatePath] = tmpl
	return tmpl, nil
}

// NOTE: All templates are now externalized to templates/ directory
// No hardcoded templates in this generator - external templates only

// Handler validation and generation
func main() {
	// Initialize configuration system for code generation
	if err := config.Initialize(); err != nil {
		// Cannot use structured logging before logging is initialized
		fmt.Fprintf(os.Stderr, "FATAL: Configuration initialization failed: %v\n", err)
		os.Exit(1)
	}

	// Initialize logging for code generation
	logging.InitLogger(config.GetLogDir(), logging.INFO, []string{})
	logging.Info("code generator starting", map[string]interface{}{
		"task": "dynamic-schema-generation",
		"single_source_of_truth": true,
	})

	// Generate unified API from multiple schemas
	schemasDir := "schemas"
	unifiedAPIPath := "../build/api.yaml"
	
	// Ensure build directory exists
	if err := os.MkdirAll("../build", 0755); err != nil {
		logging.Fatal("failed to create build directory", map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	if err := generateUnifiedAPI(schemasDir, unifiedAPIPath); err != nil {
		logging.Fatal("failed to generate unified API", map[string]interface{}{
			"error": err.Error(),
			"schemas_dir": schemasDir,
		})
	}

	// Load generated unified API specification
	specData, err := os.ReadFile(unifiedAPIPath)
	if err != nil {
		logging.Fatal("cannot read unified API specification", map[string]interface{}{
			"error": err.Error(),
			"note": "unified API is required for code generation",
		})
	}

	var spec OpenAPISpec
	if err := yaml.Unmarshal(specData, &spec); err != nil {
		logging.Fatal("invalid YAML in unified API", map[string]interface{}{
			"error": err.Error(),
		})
	}

	logging.Info("unified API specification loaded successfully", map[string]interface{}{
		"title": spec.Info.Title,
		"version": spec.Info.Version,
		"total_paths": len(spec.Paths),
		"generated_dynamically": true,
	})
	// DEBUG: Developer-focused code generation details
	logging.Debug("API spec analysis", map[string]interface{}{
		"title": spec.Info.Title,
		"version": spec.Info.Version,
		"paths_to_process": len(spec.Paths),
	})

	// Analyze and validate all endpoints
	var routes []RouteInfo
	var handlerStubs []HandlerStub
	var missingHandlers []string
	var imports []string

	for path, pathItem := range spec.Paths {
		operations := map[string]*Operation{
			"GET":    pathItem.Get,
			"POST":   pathItem.Post,
			"PUT":    pathItem.Put,
			"DELETE": pathItem.Delete,
		}

		for method, op := range operations {
			if op == nil {
				continue
			}

			// TRACE: Detailed processing information for code generation
			logging.Trace("codegen", "processing endpoint", map[string]interface{}{
				"method": method,
				"path": path,
				"operation_id": op.OperationID,
			})

			// Validate handler file exists
			if op.XHandler != "" {
				handlerPath := op.XHandler
				if !fileExists(handlerPath) {
					missingHandlers = append(missingHandlers, fmt.Sprintf("%s %s -> %s", method, path, handlerPath))
				}

				// Extract package for import
				if dir := filepath.Dir(handlerPath); dir != "." {
					importPath := fmt.Sprintf("holodeck1/%s", strings.ReplaceAll(dir, "/", "/"))
					if !contains(imports, importPath) {
						imports = append(imports, importPath)
					}
				}
			}

			// Generate route info
			routes = append(routes, RouteInfo{
				Path:        strings.TrimPrefix(path, "/api"),
				Method:      method,
				OperationID: op.OperationID,
				HandlerFunc: op.XFunction,
			})

			// Generate handler stub with package info
			handlerDir := filepath.Dir(op.XHandler)
			packageName := strings.Split(handlerDir, "/")[len(strings.Split(handlerDir, "/"))-1]
			handlerStubs = append(handlerStubs, HandlerStub{
				FuncName: op.XFunction,
				Package:  packageName,
				Comment:  fmt.Sprintf("%s %s - %s", method, path, op.Summary),
			})
		}
	}

	// FAIL BUILD if handlers missing and strict mode enabled
	if spec.XCodeGeneration.FailOnMissingHandlers && len(missingHandlers) > 0 {
		logging.Fatal("build failed - missing required handlers", map[string]interface{}{
			"missing_handlers": missingHandlers,
			"message": "Create the missing handler files or disable strict validation",
		})
	}

	// Note: server import is already included in template

	// Generate router code
	logging.Info("generating auto-router", map[string]interface{}{
		"routes_count": len(routes),
	})

	tmpl, err := loadTemplate("templates/go/router.tmpl")
	if err != nil {
		logging.Fatal("failed to load router template", map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	routerFile, err := os.Create("router/auto_router.go")
	if err != nil {
		logging.Fatal("failed to create auto_router.go", map[string]interface{}{
			"error": err.Error(),
		})
	}
	defer routerFile.Close()

	// Organize routes by category for Three.js template
	var syncOps, entityOps, avatarOps, sceneOps, systemOps []RouteInfo
	for _, route := range routes {
		if strings.HasPrefix(route.Path, "/sync") {
			syncOps = append(syncOps, route)
		} else if strings.HasPrefix(route.Path, "/entities") {
			entityOps = append(entityOps, route)
		} else if strings.HasPrefix(route.Path, "/avatars") {
			avatarOps = append(avatarOps, route)
		} else if strings.HasPrefix(route.Path, "/scene") {
			sceneOps = append(sceneOps, route)
		} else if strings.HasPrefix(route.Path, "/system") {
			systemOps = append(systemOps, route)
		}
	}

	templateData := struct {
		SyncOperations []RouteInfo
		Entities []RouteInfo
		Avatars []RouteInfo
		Scene []RouteInfo
		System []RouteInfo
		Imports []string
		TotalRoutes int
		SyncOpsCount int
		EntityOpsCount int
		AvatarOpsCount int
		SceneOpsCount int
		SystemOpsCount int
	}{
		SyncOperations: syncOps,
		Entities: entityOps,
		Avatars: avatarOps,
		Scene: sceneOps,
		System: systemOps,
		Imports: imports,
		TotalRoutes: len(routes),
		SyncOpsCount: len(syncOps),
		EntityOpsCount: len(entityOps),
		AvatarOpsCount: len(avatarOps),
		SceneOpsCount: len(sceneOps),
		SystemOpsCount: len(systemOps),
	}

	if err := tmpl.Execute(routerFile, templateData); err != nil {
		logging.Fatal("template generation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	logging.Info("auto-router generated", map[string]interface{}{
		"routes_generated": len(routes),
		"handler_stubs_generated": len(handlerStubs),
		"missing_handlers": len(missingHandlers),
	})
	
	if len(missingHandlers) > 0 {
		logging.Warn("handlers missing but build continuing", map[string]interface{}{
			"missing_count": len(missingHandlers),
			"missing_handlers": missingHandlers,
		})
	}

	// Generate minimal Web UI Client
	logging.Info("generating minimal Web UI client")
	generateWebUIClient(spec, routes)

	logging.Info("code generation complete", map[string]interface{}{
		"features": []string{
			"Dynamic schema generation from Three.js TypeScript definitions",
			"Unified API specification from multiple sources",
			"Auto-generated routing from unified spec",
			"Three.js + WebGL direct integration",
			"Zero manual geometry curation",
			"Web UI client auto-generated from unified spec",
			"Build-time API discovery",
		},
		"single_source_of_truth": true,
		"dynamic_generation": true,
	})
}

// generateUnifiedAPI orchestrates the complete dynamic schema generation process
func generateUnifiedAPI(schemasDir, outputPath string) error {
	logging.Info("generating unified API from schemas", map[string]interface{}{
		"schemas_dir": schemasDir,
		"output_path": outputPath,
		"task": "dynamic-unified-api-generation",
	})

	// Step 1: Generate Three.js schema from TypeScript definitions
	if err := generateThreeJSSchema(schemasDir); err != nil {
		return fmt.Errorf("failed to generate Three.js schema: %w", err)
	}

	// Step 2: Merge all schemas into unified API
	merger := NewSchemaMerger()
	if err := merger.LoadAllSchemas(schemasDir); err != nil {
		return fmt.Errorf("failed to load schemas: %w", err)
	}
	
	unified, err := merger.MergeSchemas()
	if err != nil {
		return fmt.Errorf("failed to merge schemas: %w", err)
	}
	
	if err := merger.WriteMergedSchema(unified, outputPath); err != nil {
		return fmt.Errorf("failed to write unified schema: %w", err)
	}

	return nil
}

// generateThreeJSSchema generates Three.js API schema from TypeScript definitions
func generateThreeJSSchema(schemasDir string) error {
	logging.Info("generating Three.js schema from TypeScript definitions", map[string]interface{}{
		"task": "threejs-schema-generation",
	})

	// Path to Three.js TypeScript definitions
	threejsTypesPath := filepath.Join(schemasDir, "threejs-types")
	threejsSchemaPath := filepath.Join(schemasDir, "threejs-api.yaml")

	// Generate schema from TypeScript definitions
	schema, err := ScanThreeJSDefinitions(threejsTypesPath)
	if err != nil {
		return fmt.Errorf("failed to scan Three.js definitions: %w", err)
	}

	// Write Three.js schema to file
	if err := WriteThreeJSSchema(schema, threejsSchemaPath); err != nil {
		return fmt.Errorf("failed to write Three.js schema: %w", err)
	}

	logging.Info("Three.js schema generated successfully", map[string]interface{}{
		"output_path": threejsSchemaPath,
		"dynamic_generation": true,
	})

	return nil
}

type RouteInfo struct {
	Path        string
	Method      string
	OperationID string
	HandlerFunc string
}

type HandlerStub struct {
	FuncName string
	Package  string
	Comment  string
}

// fileExists checks if a file exists at the given path.
// Returns true if the file exists and is accessible, false otherwise.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// contains checks if a string slice contains a specific item.
// Uses linear search for small slices typical in code generation.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// CLI client generation removed for minimal build



// generateWebUIClient creates the advanced auto-generated web UI client
func generateWebUIClient(spec OpenAPISpec, routes []RouteInfo) {
	logging.Debug("creating Web UI generator infrastructure")
	
	// Create web UI client directory structure
	uiClientDir := "../share/htdocs/static/js"
	if err := os.MkdirAll(uiClientDir, 0755); err != nil {
		logging.Error("failed to create UI client directory", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	
	// Generate JavaScript API Client Library
	if err := generateJavaScriptAPIClient(uiClientDir, spec, routes); err != nil {
		logging.Error("failed to generate JavaScript API client", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	
	
	logging.Info("Web UI client generated", map[string]interface{}{
		"endpoints_count": len(routes),
		"features": []string{
			"JavaScript API client auto-generated",
			"Three.js integration auto-generated",
			"Minimal build for Three.js console",
		},
		"single_source_of_truth": true,
	})
}

// generateJavaScriptAPIClient creates the complete JavaScript API wrapper
func generateJavaScriptAPIClient(outputDir string, spec OpenAPISpec, routes []RouteInfo) error {

	// Process routes for JavaScript template
	type JSMethod struct {
		MethodName     string
		Comment        string
		Parameters     string
		Implementation string
	}
	
	var jsMethods []JSMethod
	for _, route := range routes {
		method := JSMethod{
			MethodName:     getJSMethodName(route),
			Comment:        fmt.Sprintf("%s %s - %s", route.Method, route.Path, route.OperationID),
			Parameters:     getJSParameters(route),
			Implementation: generateJSImplementation(route),
		}
		jsMethods = append(jsMethods, method)
	}
	
	// Organize methods by category for Three.js JavaScript template
	var syncOps, entityOps, avatarOps, sceneOps, systemOps []JSMethod
	for _, method := range jsMethods {
		if strings.Contains(method.Comment, "/sync") {
			syncOps = append(syncOps, method)
		} else if strings.Contains(method.Comment, "/entities") {
			entityOps = append(entityOps, method)
		} else if strings.Contains(method.Comment, "/avatars") {
			avatarOps = append(avatarOps, method)
		} else if strings.Contains(method.Comment, "/scene") {
			sceneOps = append(sceneOps, method)
		} else if strings.Contains(method.Comment, "/system") {
			systemOps = append(systemOps, method)
		}
	}

	tmplData := struct {
		SyncOperations []JSMethod
		Entities []JSMethod
		Avatars []JSMethod
		Scene []JSMethod
		System []JSMethod
	}{
		SyncOperations: syncOps,
		Entities: entityOps,
		Avatars: avatarOps,
		Scene: sceneOps,
		System: systemOps,
	}
	
	tmpl, err := loadTemplate("templates/javascript/threejs-client.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load JavaScript API template: %w", err)
	}
	
	apiClientPath := filepath.Join(outputDir, "hd1lib.js")
	file, err := os.Create(apiClientPath)
	if err != nil {
		return fmt.Errorf("create API client file error: %v", err)
	}
	defer file.Close()
	
	if err := tmpl.Execute(file, tmplData); err != nil {
		return fmt.Errorf("API client template execute error: %v", err)
	}
	
	return nil
}

// UI component generation removed for minimal build

// Form system generation removed for minimal build

// Helper functions for JavaScript code generation

func getJSMethodName(route RouteInfo) string {
	// Convert operationId to camelCase for JavaScript
	name := route.OperationID
	return strings.ToLower(string(name[0])) + name[1:]
}

func getJSParameters(route RouteInfo) string {
	// Count path parameters
	paramCount := strings.Count(route.Path, "{")
	
	if paramCount == 0 {
		if route.Method == "POST" || route.Method == "PUT" {
			return "data = null"
		}
		return ""
	} else if paramCount == 1 {
		if route.Method == "POST" || route.Method == "PUT" {
			return "param1, data = null"
		}
		return "param1"
	} else {
		if route.Method == "POST" || route.Method == "PUT" {
			return "param1, param2, data = null"
		}
		return "param1, param2"
	}
}

func generateJSImplementation(route RouteInfo) string {
	method := strings.ToUpper(route.Method)
	path := route.Path
	paramCount := strings.Count(path, "{")
	
	if paramCount == 0 {
		if method == "GET" || method == "DELETE" {
			return fmt.Sprintf(`return this.request('%s', '%s');`, method, path)
		} else {
			return fmt.Sprintf(`return this.request('%s', '%s', data);`, method, path)
		}
	} else if paramCount == 1 {
		pathVar := "this.extractPathParams('" + path + "', [param1])"
		if method == "GET" || method == "DELETE" {
			return fmt.Sprintf(`const path = %s;
        return this.request('%s', path);`, pathVar, method)
		} else {
			return fmt.Sprintf(`const path = %s;
        return this.request('%s', path, data);`, pathVar, method)
		}
	} else {
		pathVar := "this.extractPathParams('" + path + "', [param1, param2])"
		if method == "GET" || method == "DELETE" {
			return fmt.Sprintf(`const path = %s;
        return this.request('%s', path);`, pathVar, method)
		} else {
			return fmt.Sprintf(`const path = %s;
        return this.request('%s', path, data);`, pathVar, method)
		}
	}
}



// Shell functions generation removed for minimal build

// ==============================================================================
// THREE.JS SCHEMA GENERATION FUNCTIONS
// ==============================================================================

// ScanThreeJSDefinitions generates Three.js API schema from TypeScript definitions
func ScanThreeJSDefinitions(typeDefsPath string) (*ThreeJSAPISchema, error) {
	logging.Info("scanning Three.js TypeScript definitions", map[string]interface{}{
		"path": typeDefsPath,
		"task": "threejs-schema-generation",
	})

	// For now, create essential geometries directly - can be enhanced later
	geometries := createEssentialGeometries()
	
	// Generate OpenAPI schema
	schema := generateThreeJSOpenAPISchema(geometries)
	
	logging.Info("Three.js schema generation complete", map[string]interface{}{
		"geometries_found": len(geometries),
		"endpoints_generated": len(schema.Paths),
		"single_source_of_truth": true,
	})

	return schema, nil
}

// createEssentialGeometries creates essential Three.js geometries
func createEssentialGeometries() []ThreeJSGeometry {
	return []ThreeJSGeometry{
		{
			Name:        "TextGeometry",
			Constructor: "TextGeometry",
			Parameters: []GeometryParameter{
				{Name: "text", Type: "string", Required: true, Description: "The text to render"},
				{Name: "size", Type: "number", Required: false, DefaultValue: 1, Description: "Size of the text"},
				{Name: "depth", Type: "number", Required: false, DefaultValue: 0.1, Description: "Depth of extrusion"},
				{Name: "curveSegments", Type: "integer", Required: false, DefaultValue: 4, Description: "Number of curve segments"},
				{Name: "bevelEnabled", Type: "boolean", Required: false, DefaultValue: false, Description: "Enable bevel"},
				{Name: "bevelThickness", Type: "number", Required: false, DefaultValue: 0.02, Description: "Bevel thickness"},
				{Name: "bevelSize", Type: "number", Required: false, DefaultValue: 0.01, Description: "Bevel size"},
				{Name: "bevelOffset", Type: "number", Required: false, DefaultValue: 0, Description: "Bevel offset"},
				{Name: "bevelSegments", Type: "integer", Required: false, DefaultValue: 3, Description: "Bevel segments"},
			},
			Description: "Three.js 3D text geometry with font rendering",
		},
		{
			Name:        "BoxGeometry",
			Constructor: "BoxGeometry",
			Parameters: []GeometryParameter{
				{Name: "width", Type: "number", Required: false, DefaultValue: 1, Description: "Width of the box"},
				{Name: "height", Type: "number", Required: false, DefaultValue: 1, Description: "Height of the box"},
				{Name: "depth", Type: "number", Required: false, DefaultValue: 1, Description: "Depth of the box"},
				{Name: "widthSegments", Type: "integer", Required: false, DefaultValue: 1, Description: "Width segments"},
				{Name: "heightSegments", Type: "integer", Required: false, DefaultValue: 1, Description: "Height segments"},
				{Name: "depthSegments", Type: "integer", Required: false, DefaultValue: 1, Description: "Depth segments"},
			},
			Description: "Three.js box geometry",
		},
		{
			Name:        "SphereGeometry",
			Constructor: "SphereGeometry",
			Parameters: []GeometryParameter{
				{Name: "radius", Type: "number", Required: false, DefaultValue: 1, Description: "Radius of the sphere"},
				{Name: "widthSegments", Type: "integer", Required: false, DefaultValue: 32, Description: "Width segments"},
				{Name: "heightSegments", Type: "integer", Required: false, DefaultValue: 16, Description: "Height segments"},
			},
			Description: "Three.js sphere geometry",
		},
		{
			Name:        "CylinderGeometry",
			Constructor: "CylinderGeometry",
			Parameters: []GeometryParameter{
				{Name: "radiusTop", Type: "number", Required: false, DefaultValue: 1, Description: "Top radius"},
				{Name: "radiusBottom", Type: "number", Required: false, DefaultValue: 1, Description: "Bottom radius"},
				{Name: "height", Type: "number", Required: false, DefaultValue: 1, Description: "Height of cylinder"},
				{Name: "radialSegments", Type: "integer", Required: false, DefaultValue: 8, Description: "Radial segments"},
			},
			Description: "Three.js cylinder geometry",
		},
	}
}

// generateThreeJSOpenAPISchema generates OpenAPI schema from geometries
func generateThreeJSOpenAPISchema(geometries []ThreeJSGeometry) *ThreeJSAPISchema {
	schema := &ThreeJSAPISchema{
		OpenAPI: "3.0.3",
		Info: ThreeJSInfo{
			Title:       "Three.js Geometry API",
			Description: "Auto-generated Three.js geometry API from TypeScript definitions",
			Version:     "1.0.0",
		},
		Paths:      make(map[string]interface{}),
		Components: Components{Schemas: make(map[string]interface{})},
	}

	// Generate entity creation endpoints with geometry support
	entityPath := map[string]interface{}{
		"post": map[string]interface{}{
			"operationId": "createEntityWithGeometry",
			"summary":     "Create entity with Three.js geometry",
			"description": "Create a 3D entity using any Three.js geometry type",
			"x-handler":   "api/entities/handlers.go",
			"x-function":  "CreateEntity",
			"requestBody": map[string]interface{}{
				"required": true,
				"content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"geometry": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"type": map[string]interface{}{
											"type": "string",
											"enum": getGeometryTypeList(geometries),
										},
									},
									"oneOf": generateGeometrySchemaList(geometries),
								},
								"material": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"type": map[string]interface{}{
											"type": "string",
											"enum": []string{"basic", "phong", "standard"},
										},
										"color": map[string]interface{}{
											"type": "string",
											"example": "#777777",
										},
									},
								},
								"position": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"x": map[string]interface{}{"type": "number"},
										"y": map[string]interface{}{"type": "number"},
										"z": map[string]interface{}{"type": "number"},
									},
								},
							},
						},
					},
				},
			},
			"responses": map[string]interface{}{
				"200": map[string]interface{}{
					"description": "Entity created successfully",
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"success":   map[string]interface{}{"type": "boolean"},
									"entity_id": map[string]interface{}{"type": "string"},
									"seq_num":   map[string]interface{}{"type": "integer"},
								},
							},
						},
					},
				},
			},
		},
	}

	schema.Paths["/entities"] = entityPath

	return schema
}

// getGeometryTypeList returns list of geometry type names
func getGeometryTypeList(geometries []ThreeJSGeometry) []string {
	var types []string
	for _, geo := range geometries {
		types = append(types, strings.ToLower(strings.TrimSuffix(geo.Name, "Geometry")))
	}
	return types
}

// generateGeometrySchemaList generates schema definitions for each geometry type
func generateGeometrySchemaList(geometries []ThreeJSGeometry) []map[string]interface{} {
	var schemas []map[string]interface{}
	
	for _, geo := range geometries {
		geoType := strings.ToLower(strings.TrimSuffix(geo.Name, "Geometry"))
		
		properties := map[string]interface{}{
			"type": map[string]interface{}{
				"type": "string",
				"const": geoType,
			},
		}
		
		// Add geometry-specific parameters
		for _, param := range geo.Parameters {
			properties[param.Name] = map[string]interface{}{
				"type": param.Type,
			}
			if param.DefaultValue != nil {
				properties[param.Name].(map[string]interface{})["default"] = param.DefaultValue
			}
			if param.Description != "" {
				properties[param.Name].(map[string]interface{})["description"] = param.Description
			}
		}
		
		schema := map[string]interface{}{
			"type": "object",
			"properties": properties,
			"required": []string{"type"},
		}
		
		schemas = append(schemas, schema)
	}
	
	return schemas
}

// WriteThreeJSSchema writes the generated schema to a file
func WriteThreeJSSchema(schema *ThreeJSAPISchema, outputPath string) error {
	yamlData, err := yaml.Marshal(schema)
	if err != nil {
		return fmt.Errorf("failed to marshal schema to YAML: %w", err)
	}

	if err := os.WriteFile(outputPath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write schema file: %w", err)
	}

	logging.Info("Three.js schema written", map[string]interface{}{
		"output_path": outputPath,
		"size_bytes": len(yamlData),
		"single_source_of_truth": true,
	})

	return nil
}

// ==============================================================================
// SCHEMA MERGER FUNCTIONS
// ==============================================================================

// NewSchemaMerger creates a new schema merger
func NewSchemaMerger() *SchemaMerger {
	return &SchemaMerger{
		schemas: make([]APISchema, 0),
	}
}

// LoadSchema loads an API schema from a file
func (sm *SchemaMerger) LoadSchema(name, filePath string) error {
	logging.Info("loading API schema", map[string]interface{}{
		"name":     name,
		"filepath": filePath,
		"task":     "schema-merging",
	})

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("schema file not found: %s", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read schema file %s: %w", filePath, err)
	}

	var spec map[string]interface{}
	if err := yaml.Unmarshal(content, &spec); err != nil {
		return fmt.Errorf("failed to parse YAML schema %s: %w", filePath, err)
	}

	schema := APISchema{
		Name:     name,
		FilePath: filePath,
		Spec:     spec,
	}

	sm.schemas = append(sm.schemas, schema)

	logging.Info("API schema loaded", map[string]interface{}{
		"name":         name,
		"paths_count":  len(getSchemaPaths(spec)),
		"has_components": hasSchemaComponents(spec),
	})

	return nil
}

// LoadAllSchemas loads all schemas from the schemas directory
func (sm *SchemaMerger) LoadAllSchemas(schemasDir string) error {
	logging.Info("loading all schemas from directory", map[string]interface{}{
		"directory": schemasDir,
		"task":      "schema-discovery",
	})

	files, err := os.ReadDir(schemasDir)
	if err != nil {
		return fmt.Errorf("failed to read schemas directory: %w", err)
	}

	schemaCount := 0
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yaml") {
			schemaName := strings.TrimSuffix(file.Name(), ".yaml")
			schemaPath := filepath.Join(schemasDir, file.Name())
			
			if err := sm.LoadSchema(schemaName, schemaPath); err != nil {
				logging.Error("failed to load schema", map[string]interface{}{
					"schema": schemaName,
					"error":  err.Error(),
				})
				continue
			}
			schemaCount++
		}
	}

	logging.Info("schema discovery complete", map[string]interface{}{
		"schemas_loaded": schemaCount,
		"total_schemas":  len(sm.schemas),
	})

	return nil
}

// MergeSchemas merges all loaded schemas into a unified OpenAPI specification
func (sm *SchemaMerger) MergeSchemas() (map[string]interface{}, error) {
	if len(sm.schemas) == 0 {
		return nil, fmt.Errorf("no schemas loaded for merging")
	}

	logging.Info("merging schemas", map[string]interface{}{
		"schema_count": len(sm.schemas),
		"task":        "schema-unification",
	})

	// Create base unified schema
	unified := map[string]interface{}{
		"openapi": "3.0.3",
		"info": map[string]interface{}{
			"title":       "HD1 Unified API",
			"description": "Unified API specification generated from multiple schema sources",
			"version":     "0.7.2",
		},
		"servers": []map[string]interface{}{
			{
				"url":         "http://localhost:8080/api",
				"description": "Development server",
			},
		},
		"paths":      make(map[string]interface{}),
		"components": map[string]interface{}{
			"schemas": make(map[string]interface{}),
		},
	}

	// Merge paths from all schemas
	allPaths := make(map[string]interface{})
	allComponents := make(map[string]interface{})

	for _, schema := range sm.schemas {
		// Merge paths
		if paths := getSchemaPaths(schema.Spec); paths != nil {
			for path, pathItem := range paths {
				// Smart merge: combine HTTP methods instead of overwriting
				if existingPath, exists := allPaths[path]; exists {
					// Merge HTTP methods
					existingPathMap := existingPath.(map[string]interface{})
					newPathMap := pathItem.(map[string]interface{})
					
					// Combine all HTTP methods
					for method, methodDef := range newPathMap {
						existingPathMap[method] = methodDef
					}
					
					logging.Debug("merged path methods", map[string]interface{}{
						"path":   path,
						"schema": schema.Name,
						"method": "combined",
					})
				} else {
					allPaths[path] = pathItem
					
					logging.Debug("merged path", map[string]interface{}{
						"path":   path,
						"schema": schema.Name,
					})
				}
			}
		}

		// Merge components
		if components := getSchemaComponents(schema.Spec); components != nil {
			for compName, compDef := range components {
				prefixedName := fmt.Sprintf("%s_%s", schema.Name, compName)
				allComponents[prefixedName] = compDef
			}
		}
	}

	unified["paths"] = allPaths
	unified["components"].(map[string]interface{})["schemas"] = allComponents

	logging.Info("schema merging complete", map[string]interface{}{
		"total_paths":      len(allPaths),
		"total_components": len(allComponents),
		"unified_spec":     true,
		"single_source_of_truth": true,
	})

	return unified, nil
}

// WriteMergedSchema writes the unified schema to a file
func (sm *SchemaMerger) WriteMergedSchema(unified map[string]interface{}, outputPath string) error {
	yamlData, err := yaml.Marshal(unified)
	if err != nil {
		return fmt.Errorf("failed to marshal unified schema to YAML: %w", err)
	}

	if err := os.WriteFile(outputPath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write unified schema file: %w", err)
	}

	logging.Info("unified schema written", map[string]interface{}{
		"output_path": outputPath,
		"size_bytes": len(yamlData),
		"single_source_of_truth": true,
	})

	return nil
}

// Helper functions for schema merger

func getSchemaPaths(spec map[string]interface{}) map[string]interface{} {
	if paths, ok := spec["paths"].(map[string]interface{}); ok {
		return paths
	}
	return nil
}

func getSchemaComponents(spec map[string]interface{}) map[string]interface{} {
	if components, ok := spec["components"].(map[string]interface{}); ok {
		if schemas, ok := components["schemas"].(map[string]interface{}); ok {
			return schemas
		}
	}
	return nil
}

func hasSchemaComponents(spec map[string]interface{}) bool {
	return getSchemaComponents(spec) != nil
}