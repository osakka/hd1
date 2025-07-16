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
		"task": "upstream-downstream-integration",
		"single_source_of_truth": true,
	})

	// Load API specification
	specData, err := os.ReadFile("api.yaml")
	if err != nil {
		logging.Fatal("cannot read api.yaml specification", map[string]interface{}{
			"error": err.Error(),
			"note": "specification is required for code generation",
		})
	}

	var spec OpenAPISpec
	if err := yaml.Unmarshal(specData, &spec); err != nil {
		logging.Fatal("invalid YAML in api.yaml", map[string]interface{}{
			"error": err.Error(),
		})
	}

	logging.Info("API specification loaded successfully", map[string]interface{}{
		"title": spec.Info.Title,
		"version": spec.Info.Version,
		"total_paths": len(spec.Paths),
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
	
	routerFile, err := os.Create("auto_router.go")
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
			"API specification drives all routing",
			"Auto-generated from API spec (SINGLE SOURCE)",
			"Three.js + WebGL direct integration",
			"Zero manual route configuration needed",
			"Web UI client auto-generated from spec",
			"Minimal build optimized for Three.js console",
		},
		"single_source_of_truth": true,
	})
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