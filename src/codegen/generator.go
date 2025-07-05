package main

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
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

	// Add server import
	if !contains(imports, "holodeck1/server") {
		imports = append(imports, "holodeck1/server")
	}

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

	// Generate HD1 Client from same spec
	logging.Info("generating HD1 client", map[string]interface{}{
		"source": "api.yaml",
	})
	generateHD1Client(spec, routes)

	// Generate Web UI Client
	logging.Info("generating Web UI client")
	generateWebUIClient(spec, routes)

	// Generate core shell functions from API specification
	logging.Info("generating core shell functions from API spec")
	if err := generateCoreShellFunctions(&spec, routes); err != nil {
		logging.Error("core shell function generation failed", map[string]interface{}{
			"error": err.Error(),
		})
		logging.Warn("core shell function generation failed", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		logging.Info("core shell functions generated", map[string]interface{}{
			"output_path": "/opt/hd1/lib/hd1lib.sh",
			"single_source_of_truth": true,
			"source": "api.yaml",
		})
	}

	// Advanced enhanced generation with A-Frame integration
	logging.Info("generating A-Frame integration")
	if err := generateEnhancedIntegration(spec, routes); err != nil {
		logging.Error("enhanced integration generation failed", map[string]interface{}{
			"error": err.Error(),
		})
		logging.Warn("enhanced generation failed", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		logging.Info("A-Frame integration generated", map[string]interface{}{
			"shell_output": "/opt/hd1/lib/downstream/playcanvaslib.sh",
			"javascript_output": "/opt/hd1/lib/downstream/playcanvaslib.js",
		})
	}

	logging.Info("code generation complete", map[string]interface{}{
		"features": []string{
			"API specification drives all routing",
			"Auto-generated from API spec (SINGLE SOURCE)",
			"PlayCanvas schemas drive function bridge",
			"Shell + JavaScript + CLI identical signatures",
			"PlayCanvas + WebGL seamless integration",
			"Zero manual route configuration needed",
			"Web UI client auto-generated from spec",
			"Change spec = change API + UI + shell functions automatically",
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

// generateHD1Client creates complete auto-generated Go HD1 client
func generateHD1Client(spec OpenAPISpec, routes []RouteInfo) {
	clientGoPath := "client/main.go"
	clientBinPath := "../build/bin/hd1-client"
	
	// Ensure directory exists
	if err := os.MkdirAll("client", 0755); err != nil {
		logging.Error("failed to create client directory", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	
	if err := os.MkdirAll(filepath.Dir(clientBinPath), 0755); err != nil {
		logging.Error("failed to create bin directory", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	
	// Generate Go client source
	if err := generateGoClient(clientGoPath, spec, routes); err != nil {
		logging.Error("failed to generate Go client", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	
	// Build Go client binary
	if err := buildGoClient(clientGoPath, clientBinPath); err != nil {
		logging.Error("failed to build client binary", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	
	logging.Info("hd1-client generated", map[string]interface{}{
		"commands_count": len(routes),
	})
}

// generateGoClient creates the Go source code for HD1 client
func generateGoClient(clientPath string, spec OpenAPISpec, routes []RouteInfo) error {

	// Process routes for template
	type ClientRoute struct {
		CommandName    string
		FunctionName   string
		Method         string
		Path           string
		Implementation string
	}
	
	var clientRoutes []ClientRoute
	for _, route := range routes {
		clientRoute := ClientRoute{
			CommandName:    getCommandName(route),
			FunctionName:   getFunctionName(route),
			Method:         strings.ToUpper(route.Method),
			Path:           route.Path,
			Implementation: generateGoImplementation(route),
		}
		clientRoutes = append(clientRoutes, clientRoute)
	}
	
	tmplData := struct {
		Routes []ClientRoute
	}{
		Routes: clientRoutes,
	}
	
	tmpl, err := loadTemplate("templates/go/client.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load client template: %w", err)
	}
	
	file, err := os.Create(clientPath)
	if err != nil {
		return fmt.Errorf("create file error: %v", err)
	}
	defer file.Close()
	
	if err := tmpl.Execute(file, tmplData); err != nil {
		return fmt.Errorf("template execute error: %v", err)
	}
	
	return nil
}

// buildGoClient compiles the Go client source into binary
func buildGoClient(sourcePath, binaryPath string) error {
	cmd := exec.Command("go", "build", "-o", binaryPath, sourcePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

// getCommandName converts route info to CLI command name.
// Transforms camelCase operationID to kebab-case for CLI consistency.
// Example: "createSession" becomes "create-session"
func getCommandName(route RouteInfo) string {
	// Convert operationId to kebab-case by inserting hyphens before uppercase letters
	name := route.OperationID
	result := ""
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result += "-"
		}
		result += strings.ToLower(string(r))
	}
	return result
}

// getFunctionName converts route info to Go function name
func getFunctionName(route RouteInfo) string {
	// Convert operationId to camelCase
	name := route.OperationID
	return strings.ToLower(string(name[0])) + name[1:]
}

// generateGoImplementation creates Go code for each route with dynamic parameter handling
func generateGoImplementation(route RouteInfo) string {
	method := strings.ToUpper(route.Method)
	path := route.Path
	
	// Extract all path parameters dynamically
	params := extractPathParameters(path)
	paramCount := len(params)
	
	if paramCount == 0 {
		// No path parameters
		if method == "GET" || method == "DELETE" {
			return fmt.Sprintf(`makeRequest("%s", "%s", nil)`, method, path)
		} else {
			// POST/PUT with no path params - optional JSON body
			return fmt.Sprintf(`var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[2]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %%v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("%s", "%s", body)`, method, path)
		}
	} else {
		// Dynamic path parameter substitution
		pathTemplate := buildPathTemplate(path, params)
		minArgsRequired := 1 + paramCount + 1 // program + command + params
		bodyArgIndex := minArgsRequired // body is after all required args
		
		argCheckStr := fmt.Sprintf("len(os.Args) < %d", minArgsRequired)
		errorMsg := "Missing required parameter"
		if paramCount > 1 {
			errorMsg = "Missing required parameters"
		}
		
		if method == "GET" || method == "DELETE" {
			return fmt.Sprintf(`if %s {
		fmt.Println("Error: %s")
		os.Exit(1)
	}
	makeRequest("%s", %s, nil)`, argCheckStr, errorMsg, method, pathTemplate)
		} else {
			return fmt.Sprintf(`if %s {
		fmt.Println("Error: %s")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > %d {
		if err := json.Unmarshal([]byte(os.Args[%d]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %%v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("%s", %s, body)`, argCheckStr, errorMsg, minArgsRequired-1, bodyArgIndex, method, pathTemplate)
		}
	}
}

// extractPathParameters extracts all path parameters from a URL path
func extractPathParameters(path string) []string {
	var params []string
	start := 0
	for {
		openIndex := strings.Index(path[start:], "{")
		if openIndex == -1 {
			break
		}
		openIndex += start
		closeIndex := strings.Index(path[openIndex:], "}")
		if closeIndex == -1 {
			break
		}
		closeIndex += openIndex
		param := path[openIndex+1 : closeIndex]
		params = append(params, param)
		start = closeIndex + 1
	}
	return params
}

// buildPathTemplate creates Go string concatenation for dynamic path building
func buildPathTemplate(path string, params []string) string {
	template := `"` + path + `"`
	
	// Replace each parameter with os.Args substitution
	for i, param := range params {
		argIndex := 2 + i // os.Args[0] = program, os.Args[1] = command, os.Args[2+] = params
		placeholder := "{" + param + "}"
		replacement := `" + os.Args[` + fmt.Sprintf("%d", argIndex) + `] + "`
		template = strings.Replace(template, placeholder, replacement, 1)
	}
	
	return template
}

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
	
	// Generate UI Component Library
	if err := generateUIComponents(uiClientDir, spec, routes); err != nil {
		logging.Error("failed to generate UI components", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	
	// Generate Dynamic Form System
	if err := generateFormSystem(uiClientDir, spec, routes); err != nil {
		logging.Error("failed to generate form system", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	
	logging.Info("Web UI client generated", map[string]interface{}{
		"endpoints_count": len(routes),
		"features": []string{
			"JavaScript API client auto-generated",
			"UI components auto-generated",
			"Dynamic forms auto-generated",
			"Zero manual UI synchronization needed",
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
	
	tmplData := struct {
		Methods []JSMethod
	}{
		Methods: jsMethods,
	}
	
	tmpl, err := loadTemplate("templates/javascript/api-client.tmpl")
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

// generateUIComponents creates auto-generated UI components for each endpoint
func generateUIComponents(outputDir string, spec OpenAPISpec, routes []RouteInfo) error {

	// Process routes for UI components
	type UIComponent struct {
		Name         string
		ClassName    string
		Comment      string
		Endpoint     string
		Method       string
		HTML         string
		ExecuteParams string
		APIMethod    string
		APIParams    string
	}
	
	var components []UIComponent
	for _, route := range routes {
		component := UIComponent{
			Name:         getComponentName(route),
			ClassName:    getComponentClassName(route),
			Comment:      fmt.Sprintf("Component for %s %s", route.Method, route.Path),
			Endpoint:     route.Path,
			Method:       route.Method,
			HTML:         generateComponentHTML(route),
			ExecuteParams: getExecuteParams(route),
			APIMethod:    getJSMethodName(route),
			APIParams:    getAPIParams(route),
		}
		components = append(components, component)
	}
	
	tmplData := struct {
		Components []UIComponent
	}{
		Components: components,
	}
	
	tmpl, err := loadTemplate("templates/javascript/ui-components.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load UI components template: %w", err)
	}
	
	componentsPath := filepath.Join(outputDir, "hd1-ui-components.js")
	file, err := os.Create(componentsPath)
	if err != nil {
		return fmt.Errorf("create UI components file error: %v", err)
	}
	defer file.Close()
	
	if err := tmpl.Execute(file, tmplData); err != nil {
		return fmt.Errorf("UI components template execute error: %v", err)
	}
	
	return nil
}

// generateFormSystem creates dynamic form generation system
func generateFormSystem(outputDir string, spec OpenAPISpec, routes []RouteInfo) error {

	// Process routes for form schemas
	type FormSchema struct {
		Name   string
		Schema string
	}
	
	var formSchemas []FormSchema
	for _, route := range routes {
		if route.Method == "POST" || route.Method == "PUT" {
			schema := FormSchema{
				Name:   getFormSchemaName(route),
				Schema: generateFormSchemaJSON(route),
			}
			formSchemas = append(formSchemas, schema)
		}
	}
	
	tmplData := struct {
		FormSchemas []FormSchema
	}{
		FormSchemas: formSchemas,
	}
	
	tmpl, err := loadTemplate("templates/javascript/form-system.tmpl")
	if err != nil {
		return fmt.Errorf("form system template load error: %v", err)
	}
	
	formsPath := filepath.Join(outputDir, "hd1-form-system.js")
	file, err := os.Create(formsPath)
	if err != nil {
		return fmt.Errorf("create form system file error: %v", err)
	}
	defer file.Close()
	
	if err := tmpl.Execute(file, tmplData); err != nil {
		return fmt.Errorf("form system template execute error: %v", err)
	}
	
	return nil
}

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

func getComponentName(route RouteInfo) string {
	return strings.ToLower(route.OperationID)
}

func getComponentClassName(route RouteInfo) string {
	return strings.Title(route.OperationID)
}

func generateComponentHTML(route RouteInfo) string {
	title := strings.Title(strings.ReplaceAll(route.OperationID, "_", " "))
	formId := route.OperationID + "-form"
	resultId := route.OperationID + "-result"
	
	html := fmt.Sprintf(`<div class="hd1-component"><h4>%s</h4><form id="%s">`, title, formId)
	
	// Add form fields based on method and path parameters
	paramCount := strings.Count(route.Path, "{")
	if paramCount > 0 {
		html += `<div class="form-field"><label for="param1">Parameter 1:</label><input type="text" name="param1" required></div>`
	}
	if paramCount > 1 {
		html += `<div class="form-field"><label for="param2">Parameter 2:</label><input type="text" name="param2" required></div>`
	}
	
	if route.Method == "POST" || route.Method == "PUT" {
		html += `<div class="form-field"><label for="data">Request Body (JSON):</label><textarea name="data" placeholder="{}"></textarea></div>`
	}
	
	html += fmt.Sprintf(`<button type="submit">Execute</button></form><div id="%s" class="result-area"></div></div>`, resultId)
	
	// Escape for JavaScript string
	html = strings.ReplaceAll(html, "'", "\\'")
	html = strings.ReplaceAll(html, "\n", "")
	
	return html
}

func getExecuteParams(route RouteInfo) string {
	return "formData"
}

func getAPIParams(route RouteInfo) string {
	paramCount := strings.Count(route.Path, "{")
	
	if paramCount == 0 {
		if route.Method == "POST" || route.Method == "PUT" {
			return "formData"
		}
		return ""
	} else if paramCount == 1 {
		if route.Method == "POST" || route.Method == "PUT" {
			return "formData.param1, formData"
		}
		return "formData.param1"
	} else {
		if route.Method == "POST" || route.Method == "PUT" {
			return "formData.param1, formData.param2, formData"
		}
		return "formData.param1, formData.param2"
	}
}

func getFormSchemaName(route RouteInfo) string {
	return route.OperationID + "Form"
}

func generateFormSchemaJSON(route RouteInfo) string {
	title := strings.Title(strings.ReplaceAll(route.OperationID, "_", " "))
	
	schema := fmt.Sprintf(`{
        "title": "%s",
        "submitText": "Execute %s",
        "fields": {`, title, route.Method)
	
	// Add path parameter fields
	paramCount := strings.Count(route.Path, "{")
	fields := []string{}
	
	if paramCount > 0 {
		fields = append(fields, `"param1": {"type": "string", "title": "Parameter 1", "required": true}`)
	}
	if paramCount > 1 {
		fields = append(fields, `"param2": {"type": "string", "title": "Parameter 2", "required": true}`)
	}
	
	if route.Method == "POST" || route.Method == "PUT" {
		fields = append(fields, `"data": {"type": "string", "title": "Request Body (JSON)", "placeholder": "{}"}`)
	}
	
	schema += strings.Join(fields, ", ")
	schema += "}}"
	
	return schema
}

// generateEnhancedIntegration creates A-Frame bridge functions
func generateEnhancedIntegration(spec OpenAPISpec, routes []RouteInfo) error {
	logging.Info("generating enhanced A-Frame integration", map[string]interface{}{
		"task": "upstream-downstream-api-bridge",
	})
	
	// Create enhanced shell functions
	if err := generateEnhancedShellFunctions(spec, routes); err != nil {
		return fmt.Errorf("failed to generate enhanced shell functions: %w", err)
	}
	
	// Create JavaScript function bridge
	if err := generateJavaScriptBridge(spec, routes); err != nil {
		return fmt.Errorf("failed to generate JavaScript bridge: %w", err)  
	}
	
	return nil
}

// generateEnhancedShellFunctions creates shell functions with A-Frame integration
func generateEnhancedShellFunctions(spec OpenAPISpec, routes []RouteInfo) error {
	outputDir := "/opt/hd1/lib"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// A-Frame capabilities mapping (used for documentation)
	_ = map[string]interface{}{
		"geometry_types": []string{"box", "sphere", "cylinder", "cone", "plane"},
		"light_types": []string{"directional", "point", "ambient", "spot"},
		"material_properties": []string{"color", "metalness", "roughness", "transparency", "emissive"},
		"physics_bodies": []string{"dynamic", "static", "kinematic"},
	}
	
	tmpl, err := loadTemplate("templates/shell/playcanvas-functions.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load PlayCanvas shell template: %w", err)
	}

	outputPath := filepath.Join(outputDir, "downstream/playcanvaslib.sh")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create enhanced shell functions file: %w", err)
	}
	defer outputFile.Close()

	// Pass configuration data to template
	templateData := struct {
		Config *config.HD1Config
	}{
		Config: config.Config,
	}
	
	if err := tmpl.Execute(outputFile, templateData); err != nil {
		return fmt.Errorf("failed to execute PlayCanvas shell template: %w", err)
	}
	
	// Set executable permissions
	if err := os.Chmod(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to write enhanced shell functions: %w", err)
	}
	
	logging.Info("enhanced shell functions generated", map[string]interface{}{
		"output_path": outputPath,
		"playcanvas_integration": true,
	})
	
	return nil
}

// generateJavaScriptBridge creates JavaScript functions with identical signatures
func generateJavaScriptBridge(spec OpenAPISpec, routes []RouteInfo) error {
	outputDir := "/opt/hd1/lib"
	
	tmpl, err := loadTemplate("templates/javascript/playcanvas-bridge.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load PlayCanvas bridge template: %w", err)
	}

	outputPath := filepath.Join(outputDir, "downstream/playcanvaslib.js")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create JavaScript bridge file: %w", err)
	}
	defer outputFile.Close()

	// Pass configuration data to template
	templateData := struct {
		Config *config.HD1Config
	}{
		Config: config.Config,
	}
	
	if err := tmpl.Execute(outputFile, templateData); err != nil {
		return fmt.Errorf("failed to execute PlayCanvas bridge template: %w", err)
	}
	
	logging.Info("JavaScript function bridge generated", map[string]interface{}{
		"output_path": outputPath,
		"identical_signatures": true,
	})
	
	return nil
}

// generateCoreShellFunctions creates hd1-functions.sh from API specification
func generateCoreShellFunctions(spec *OpenAPISpec, routes []RouteInfo) error {
	outputDir := "/opt/hd1/lib"
	
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Generate core shell function library
	tmpl, err := loadTemplate("templates/shell/core-functions.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load core shell template: %w", err)
	}
	outputPath := filepath.Join(outputDir, "hd1lib.sh")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create core shell functions file: %w", err)
	}
	defer outputFile.Close()

	// Pass routes data and configuration defaults to template for dynamic generation
	templateData := struct {
		Routes           []RouteInfo
		DefaultAPIBase   string
		DefaultSessionID string
	}{
		Routes:           routes,
		DefaultAPIBase:   config.GetAPIBase(),
		DefaultSessionID: config.GetSessionDefaultID(),
	}
	
	if err := tmpl.Execute(outputFile, templateData); err != nil {
		return fmt.Errorf("failed to execute core shell template: %w", err)
	}
	
	// Set executable permissions
	if err := os.Chmod(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to write core shell functions: %w", err)
	}
	
	logging.Info("core shell functions generated", map[string]interface{}{
		"output_path": outputPath,
		"source": "api.yaml",
		"single_source_of_truth": true,
	})
	
	return nil
}