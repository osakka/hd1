# HD1 Code Documentation Standards

**Version: v5.0.5**  
**Last Updated: 2025-07-03**

## 🎯 Documentation Philosophy

HD1 follows the principle of **"Just Perfect Documentation"** - not too much, not too little. Every piece of code should be self-documenting with strategic comments that enhance comprehension without redundancy.

## 📋 Documentation Requirements

### **1. Package Documentation**
Every package MUST have package-level documentation:

```go
// Package entities provides HTTP handlers for entity lifecycle management
// in HD1 sessions. Entities represent 3D objects with components that can
// be manipulated via the REST API and synchronized via WebSocket.
//
// Key concepts:
//   - Entity: PlayCanvas 3D object with transform, model, material components
//   - Session: Isolated game world containing entities
//   - WebSocket sync: Real-time entity updates across clients
package entities
```

### **2. Function Documentation**

#### **Exported Functions (Required)**
All exported functions MUST have documentation following Go conventions:

```go
// CreateEntityHandler handles POST /sessions/{sessionId}/entities requests.
// Creates a new PlayCanvas entity in the specified session with proper
// component structure and WebSocket broadcasting.
//
// The request body should contain entity name and optional components.
// Returns 201 Created with entity details on success, or appropriate
// error status codes for validation failures.
//
// WebSocket event: Broadcasts 'entity_created' to all session clients.
func CreateEntityHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
```

#### **Unexported Functions (When Complex)**
Complex unexported functions should have documentation:

```go
// extractSessionID parses session ID from URL path following the pattern
// /sessions/{sessionId}/... and validates the format.
// Returns empty string if parsing fails or format is invalid.
func extractSessionID(path string) string {
```

### **3. Struct and Type Documentation**

#### **Exported Types (Required)**
```go
// CreateEntityRequest represents the request body for entity creation.
// All fields except Name are optional and will use PlayCanvas defaults.
type CreateEntityRequest struct {
    Name       string                 `json:"name"`        // Required: entity display name (1-64 chars)
    Tags       []string               `json:"tags,omitempty"` // Optional: entity tags for organization
    Components map[string]interface{} `json:"components,omitempty"` // Optional: PlayCanvas components
}
```

#### **Configuration Structs**
```go
// ServerConfig holds HTTP server configuration with environment variable
// support and command-line flag overrides following the priority:
// Flags > Environment Variables > .env File > Defaults
type ServerConfig struct {
    Host string `json:"host"` // HTTP server bind address (default: localhost)
    Port string `json:"port"` // HTTP server port (default: 8080)
}
```

### **4. Global Variables and Constants**

#### **Package-level Variables**
```go
// templateCache stores compiled templates for performance optimization.
// Templates are loaded once at startup and reused across requests.
var templateCache = make(map[string]*template.Template)

// defaultSessionTimeout defines the maximum idle time before session cleanup.
const defaultSessionTimeout = 30 * time.Minute
```

### **5. Complex Business Logic**

#### **Section Comments for Large Functions**
```go
func ProcessEntityUpdate(sessionID string, entityID string, data UpdateData) error {
    // Validation phase: Check session exists and entity is valid
    session, exists := hub.sessions[sessionID]
    if !exists {
        return ErrSessionNotFound
    }

    // Component processing: Apply updates to PlayCanvas components
    for componentType, componentData := range data.Components {
        if err := validateComponent(componentType, componentData); err != nil {
            return fmt.Errorf("invalid component %s: %w", componentType, err)
        }
    }

    // WebSocket broadcast: Notify all session clients of changes
    hub.BroadcastToSession(sessionID, EntityUpdateEvent{
        EntityID: entityID,
        Data:     data,
    })

    return nil
}
```

## 🎨 Style Guidelines

### **Commenting Style**
- **Clarity over cleverness**: Use simple, direct language
- **Present tense**: "Returns user data" not "Will return user data"
- **No redundancy**: Don't repeat what the code obviously does
- **Focus on why and what**: Explain purpose and business logic, not syntax
- **Atomic Naming**: Use lowercase_snake_case for consistency and semantic clarity

### **Good Examples**
```go
// ✅ Good: Explains purpose and behavior
// validate_session_id checks if the session exists and is active.
// Returns ErrSessionNotFound if session doesn't exist or expired.

// ✅ Good: Explains business context
// broadcast_entity_update sends real-time entity changes to all clients
// in the session, excluding the client that initiated the change.

// ✅ Good: Atomic naming with clear semantics
// convert_to_daemon_process detaches the HD1 process from terminal
// and creates PID file for system service management.
```

### **Bad Examples**
```go
// ❌ Bad: States the obvious
// GetSessionID gets the session ID

// ❌ Bad: Redundant with function signature
// CreateEntity creates an entity with name and components

// ❌ Bad: Too vague
// HandleRequest handles the request

// ❌ Bad: Non-indicative naming patterns
// func fixSession()      // What kind of fix?
// func tmpHandler()      // Temporary implies incomplete
// func ensureDirectories() // Vague about what directories
```

## 🔧 Implementation Checklist

### **For Every Source File:**
- [ ] Package documentation at top of main file
- [ ] All exported functions documented
- [ ] All exported types documented
- [ ] Complex unexported functions documented
- [ ] Global variables explained
- [ ] Section comments for large functions (>20 lines)

### **Quality Criteria:**
- [ ] Documentation enhances code comprehension
- [ ] Business logic and architectural decisions explained
- [ ] Error conditions and return values documented
- [ ] Integration points and dependencies clarified
- [ ] Performance considerations noted where relevant

## 📊 Documentation Categories by Priority

### **Priority 1: Critical (Must Have)**
- API handlers (user-facing functionality)
- Server infrastructure (hub, client management)
- Configuration system
- Core business logic

### **Priority 2: High (Should Have)**
- Utility functions used across packages
- Complex algorithms and data processing
- Error handling and validation logic
- Template and code generation

### **Priority 3: Medium (Nice to Have)**
- Simple getter/setter functions
- Obvious helper functions
- Test utilities

## 🚀 Automation and Tooling

### **Documentation Linting**
```bash
# Check for missing documentation
go doc -all ./... | grep -i "undocumented"

# Lint documentation quality
golangci-lint run --enable=godoc
```

### **Documentation Generation**
```bash
# Generate package documentation
go doc -all ./api/entities
```

## 📈 Success Metrics

- **Coverage**: 100% of exported functions and types documented
- **Quality**: All critical business logic explained
- **Consistency**: Uniform style across all packages
- **Maintainability**: New developers can understand code from comments alone

---

**HD1 v5.0.5** - API-First Game Engine Platform  
**Code Documentation Standards** - Just Perfect Documentation