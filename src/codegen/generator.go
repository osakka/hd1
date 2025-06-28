package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	
	"gopkg.in/yaml.v3"
)

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

// Generated Router Template
const routerTemplate = `// AUTO-GENERATED FROM api.yaml - DO NOT EDIT MANUALLY
// This file is the SINGLE SOURCE OF TRUTH for routing
package main

import (
	"net/http"
	"strings"
	"log"
{{range .Imports}}
	"{{.}}"{{end}}
)

// Route represents a single API route
type Route struct {
	Path       string
	Method     string
	Handler    func(http.ResponseWriter, *http.Request)
	OperationID string
}

// APIRouter manages all auto-generated routes
type APIRouter struct {
	routes []Route
	hub    *server.Hub
}

// NewAPIRouter creates router from specification
func NewAPIRouter(hub *server.Hub) *APIRouter {
	router := &APIRouter{hub: hub}
	router.generateRoutes()
	return router
}

// ServeHTTP implements http.Handler interface
func (r *APIRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Route matching with parameter extraction
	path := strings.TrimPrefix(req.URL.Path, "/api")
	for _, route := range r.routes {
		if r.matchRoute(route.Path, path) && route.Method == req.Method {
			log.Printf("API: %s %s -> %s", req.Method, req.URL.Path, route.OperationID)
			route.Handler(w, req)
			return
		}
	}

	// Route not found
	http.Error(w, "API endpoint not found", http.StatusNotFound)
}

// Match route with parameter support (e.g. /sessions/{sessionId})
func (r *APIRouter) matchRoute(pattern, path string) bool {
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return false
	}

	for i, part := range patternParts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			// Parameter placeholder - matches any value
			continue
		}
		if part != pathParts[i] {
			return false
		}
	}
	return true
}

// AUTO-GENERATED ROUTE DEFINITIONS
func (r *APIRouter) generateRoutes() {
	r.routes = []Route{
{{range .Routes}}
		{
			Path:       "{{.Path}}",
			Method:     "{{.Method}}",
			Handler:    r.{{.HandlerFunc}},
			OperationID: "{{.OperationID}}",
		},{{end}}
	}
}

{{range .HandlerStubs}}
// {{.Comment}}
func (r *APIRouter) {{.FuncName}}(w http.ResponseWriter, req *http.Request) {
	{{.Package}}.{{.FuncName}}Handler(w, req, r.hub)
}
{{end}}
`

// Handler validation and generation
func main() {
	fmt.Println("🧠 BRAIN SURGEON CODE GENERATOR - SPEC-DRIVEN DEVELOPMENT")
	fmt.Println("========================================================")

	// Load API specification
	specData, err := os.ReadFile("api.yaml")
	if err != nil {
		log.Fatal("❌ FATAL: Cannot read api.yaml - Specification is REQUIRED!")
	}

	var spec OpenAPISpec
	if err := yaml.Unmarshal(specData, &spec); err != nil {
		log.Fatal("❌ FATAL: Invalid YAML in api.yaml:", err)
	}

	fmt.Printf("✅ Loaded API Spec: %s v%s\n", spec.Info.Title, spec.Info.Version)
	fmt.Printf("📊 Found %d paths to process\n", len(spec.Paths))

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

			fmt.Printf("🔍 Processing: %s %s -> %s\n", method, path, op.OperationID)

			// Validate handler file exists
			if op.XHandler != "" {
				handlerPath := op.XHandler
				if !fileExists(handlerPath) {
					missingHandlers = append(missingHandlers, fmt.Sprintf("%s %s -> %s", method, path, handlerPath))
				}

				// Extract package for import
				if dir := filepath.Dir(handlerPath); dir != "." {
					importPath := fmt.Sprintf("holodeck/%s", strings.ReplaceAll(dir, "/", "/"))
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
		fmt.Println("\n❌ BUILD FAILURE - MISSING REQUIRED HANDLERS:")
		for _, missing := range missingHandlers {
			fmt.Printf("   🚫 %s\n", missing)
		}
		fmt.Println("\n💡 Create the missing handler files or disable strict validation")
		os.Exit(1)
	}

	// Add server import
	if !contains(imports, "holodeck/server") {
		imports = append(imports, "holodeck/server")
	}

	// Generate router code
	fmt.Printf("\n🏗️  Generating auto-router with %d routes...\n", len(routes))

	tmpl := template.Must(template.New("router").Parse(routerTemplate))
	
	routerFile, err := os.Create("auto_router.go")
	if err != nil {
		log.Fatal("❌ Cannot create auto_router.go:", err)
	}
	defer routerFile.Close()

	templateData := struct {
		Routes       []RouteInfo
		HandlerStubs []HandlerStub
		Imports      []string
	}{
		Routes:       routes,
		HandlerStubs: handlerStubs,
		Imports:      imports,
	}

	if err := tmpl.Execute(routerFile, templateData); err != nil {
		log.Fatal("❌ Template generation failed:", err)
	}

	fmt.Println("✅ SUCCESS: auto_router.go generated")
	fmt.Printf("📊 Generated %d routes\n", len(routes))
	fmt.Printf("🔗 Generated %d handler stubs\n", len(handlerStubs))
	
	if len(missingHandlers) > 0 {
		fmt.Printf("⚠️  WARNING: %d handlers missing (build will continue)\n", len(missingHandlers))
	}

	fmt.Println("\n🚀 SPEC-DRIVEN DEVELOPMENT COMPLETE!")
	fmt.Println("   • API specification drives all routing")
	fmt.Println("   • Handler files validated at build time")
	fmt.Println("   • Zero manual route configuration needed")
	fmt.Println("   • Change spec = change API automatically")
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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}