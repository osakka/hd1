# HD1 Code Generation Templates

This directory contains all external templates for the HD1 code generator, organized by target language. This architecture replaces the previous hardcoded template strings embedded in `generator.go` for better maintainability and developer experience.

## Directory Structure

```
templates/
├── go/                     # Go language templates
│   ├── router.tmpl        # HTTP router and handler generation
│   └── client.tmpl        # CLI client application
├── javascript/            # JavaScript templates
│   ├── api-client.tmpl    # Browser API client library
│   ├── ui-components.tmpl # Dynamic UI components
│   ├── form-system.tmpl   # Form generation system
│   └── aframe-bridge.tmpl # A-Frame WebXR integration
└── shell/                 # Shell script templates
    ├── core-functions.tmpl     # Core API wrapper functions
    └── aframe-functions.tmpl   # Enhanced A-Frame functions
```

## Template Usage

### Template Loading System
Templates are loaded using Go's embed filesystem with caching:

```go
//go:embed templates/*
var templateFS embed.FS

// Templates are cached for performance
tmpl, err := loadTemplate("templates/go/router.tmpl")
```

### Template Syntax
All templates use Go's `text/template` syntax:

```go
{{range .Routes}}
// {{.Comment}}
func {{.FunctionName}}() {
    {{.Implementation}}
}
{{end}}
```

## Development Workflow

### Modifying Templates

1. **Edit template files directly** - No Go code changes needed
2. **Test generation**: `cd src && make clean && make`
3. **Validate output**: Check generated files in expected locations
4. **Commit changes**: Include both template and any generated file updates

### Adding New Templates

1. **Create template file** in appropriate language directory
2. **Add loadTemplate() call** in `generator.go`
3. **Update this README** with new template documentation
4. **Test build process** to ensure no regressions

### Template Guidelines

- **Preserve formatting**: Maintain consistent indentation and spacing
- **Comment comprehensively**: Explain template variables and logic
- **Test thoroughly**: Validate generated code compiles and runs
- **Follow conventions**: Match existing template patterns and structure

## Template Reference

### Go Templates

#### `go/router.tmpl`
- **Purpose**: Generates HTTP router with auto-mapped handlers
- **Output**: `src/auto_router.go`
- **Variables**: `.Routes`, `.HandlerStubs`, `.Imports`
- **Size**: ~122 lines

#### `go/client.tmpl`
- **Purpose**: Generates CLI client with command mapping
- **Output**: `src/client/main.go`
- **Variables**: `.Routes` (with `.CommandName`, `.FunctionName`, etc.)
- **Size**: ~98 lines

### JavaScript Templates

#### `javascript/api-client.tmpl`
- **Purpose**: Browser-compatible API client library
- **Output**: `share/htdocs/static/js/hd1lib.js`
- **Variables**: Template-driven but mostly static
- **Size**: ~91 lines

#### `javascript/ui-components.tmpl`
- **Purpose**: Dynamic UI components for web interface
- **Output**: `share/htdocs/static/js/hd1-ui-components.js`
- **Variables**: Route-based component generation
- **Size**: ~117 lines

#### `javascript/form-system.tmpl`
- **Purpose**: Dynamic form generation from API schemas
- **Output**: `share/htdocs/static/js/hd1-form-system.js`
- **Variables**: `.FormSchemas` with validation rules
- **Size**: ~119 lines

#### `javascript/aframe-bridge.tmpl`
- **Purpose**: WebXR/A-Frame integration with validation
- **Output**: `lib/downstream/aframelib.js`
- **Variables**: Static template with A-Frame schema validation
- **Size**: ~199 lines

### Shell Templates

#### `shell/core-functions.tmpl`
- **Purpose**: Core shell function library from API spec
- **Output**: `lib/hd1lib.sh`
- **Variables**: Auto-generated from OpenAPI specification
- **Size**: ~177 lines

#### `shell/aframe-functions.tmpl`
- **Purpose**: Enhanced shell functions with A-Frame integration
- **Output**: `lib/downstream/aframelib.sh`
- **Variables**: Static template with enhanced validation
- **Size**: ~214 lines

## Performance Considerations

- **Template caching**: Templates are cached in memory after first load
- **Single binary**: All templates embedded in final HD1 binary
- **Build time**: No significant impact on build performance
- **Runtime**: Template loading happens only during code generation

## Troubleshooting

### Common Issues

**Template not found**
```
Error: failed to read template templates/new-template.tmpl: file does not exist
```
- Ensure template file exists in correct directory
- Check embed directive includes new template path

**Template parse error**
```
Error: failed to parse template: unexpected "}" in command
```
- Validate Go template syntax
- Check balanced braces and correct variable references

**Generated code doesn't compile**
```
Error: undefined: SomeFunction
```
- Verify template generates syntactically correct target language
- Check template variable data matches expected structure

### Validation Commands

```bash
# Full build test
cd src && make clean && make

# Template syntax validation (for Go templates)
go run -c 'template.Must(template.ParseFiles("templates/go/router.tmpl"))'

# Generated output verification
diff expected_output.go generated_output.go
```

## Related Documentation

- **ADR-003**: Template Externalization Architecture Decision
- **CLAUDE.md**: Development context and template architecture overview
- **src/codegen/generator.go**: Template loading implementation
- **Makefile**: Build process integration

## Maintenance

This template system is designed for:
- **Long-term maintainability**: Clear separation of concerns
- **Developer productivity**: Direct template editing with IDE support
- **Zero regression**: Identical output through surgical refactoring
- **Single source of truth**: API specification drives all generation

For questions or issues, refer to the ADR-003 documentation or examine the `loadTemplate()` implementation in `generator.go`.