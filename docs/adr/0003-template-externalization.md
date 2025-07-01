# ADR-003: Template Externalization for HD1 Code Generator

## Status
Accepted - Implemented 2025-07-01

## Context
The HD1 code generator contained 8 massive hardcoded template strings embedded directly in Go source code (generator.go), totaling over 2,000 lines of template content. This approach created several maintenance challenges:

- **Poor maintainability**: Templates mixed with Go logic made both harder to read and modify
- **Developer experience**: No syntax highlighting or IDE support for templates
- **Version control**: Large diffs when modifying templates obscured actual logic changes  
- **Code organization**: Single 2,000+ line file violated separation of concerns
- **Testing**: Difficult to test templates independently from generator logic
- **Collaboration**: Frontend developers couldn't easily modify JS/HTML templates

### Original Architecture Problems
```go
// Before: Hardcoded templates in Go code
const routerTemplate = `// 122 lines of Go template code...`
const clientTemplate = `// 98 lines of Go template code...` 
const apiClientTemplate = `// 91 lines of JavaScript template...`
// ... 5 more massive template constants
```

## Decision
Externalize all hardcoded templates to dedicated template files organized by language, using Go's embed filesystem for distribution as a single binary.

### New Template Architecture
```
src/codegen/templates/
├── go/
│   ├── router.tmpl       # Auto-router generation (122 lines)
│   └── client.tmpl       # Go CLI client (98 lines)
├── javascript/
│   ├── api-client.tmpl   # JS API wrapper (91 lines) 
│   ├── ui-components.tmpl # UI components (117 lines)
│   ├── form-system.tmpl  # Dynamic forms (119 lines)
│   └── aframe-bridge.tmpl # A-Frame integration (199 lines)
└── shell/
    ├── core-functions.tmpl    # Core API functions (177 lines)
    └── aframe-functions.tmpl  # Enhanced A-Frame (214 lines)
```

### Implementation Strategy
- **Go embed filesystem**: `//go:embed templates/*` for single binary distribution
- **Template caching**: Performance optimization with `map[string]*template.Template`
- **Surgical refactoring**: Zero regression policy with identical output validation
- **Error handling**: Comprehensive error messages for missing/invalid templates

## Consequences

### Positive
- **Maintainability**: Templates now have proper syntax highlighting and IDE support
- **Separation of concerns**: Logic and templates cleanly separated
- **Developer experience**: Frontend developers can directly edit JS/HTML templates
- **Version control**: Clean diffs showing only actual changes
- **Testing**: Templates can be tested independently
- **Organization**: Language-based directory structure for intuitive navigation
- **Performance**: Template caching reduces file I/O overhead
- **Distribution**: Single binary maintains deployment simplicity

### Negative  
- **File count**: 8 additional template files to maintain
- **Complexity**: Slightly more complex build process with embed directives
- **Learning curve**: Developers need to understand new template location conventions

### Risk Mitigation
- **Zero regression testing**: Validated all generated files remain identical
- **Comprehensive error handling**: Clear messages for missing/invalid templates
- **Documentation**: Complete ADR, README, and workflow documentation
- **Backup preservation**: Original generator.go preserved for rollback if needed

## Implementation Notes

### Template Loading System
```go
//go:embed templates/*
var templateFS embed.FS

var templateCache = make(map[string]*template.Template)

func loadTemplate(templatePath string) (*template.Template, error) {
    if tmpl, exists := templateCache[templatePath]; exists {
        return tmpl, nil
    }
    
    content, err := templateFS.ReadFile(templatePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read template %s: %w", templatePath, err)
    }
    
    tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(content))
    if err != nil {
        return nil, fmt.Errorf("failed to parse template %s: %w", templatePath, err)
    }
    
    templateCache[templatePath] = tmpl
    return tmpl, nil
}
```

### Migration Results
- **Before**: 1 file, 2,000+ lines of mixed Go code and templates
- **After**: 9 files with clean separation, identical functionality
- **Build time**: No performance regression (caching optimization)
- **Generated output**: 100% identical to pre-refactoring (verified)

## Alternatives Considered

1. **External template directory**: Rejected due to deployment complexity
2. **Database-stored templates**: Rejected as overkill for build-time generation  
3. **Code generation from templates**: Rejected as circular dependency
4. **Template inheritance system**: Deferred to future iteration for simplicity

## References
- Go embed documentation: https://pkg.go.dev/embed
- Template best practices: Go template/text package documentation
- HD1 specification-driven development principles
- Single source of truth architecture requirements

---
**Decision made by**: Claude Code Assistant  
**Stakeholders consulted**: HD1 Development Team (implied via user directive)  
**Implementation date**: 2025-07-01  
**Review date**: Next major template system changes