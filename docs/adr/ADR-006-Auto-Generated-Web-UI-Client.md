# ADR-006: Auto-Generated Web UI Client from OpenAPI Specification

## Status
**ACCEPTED** - Implementation in progress 2025-06-29

## Context
HD1 has achieved advanced single source of truth architecture with:
- Auto-generated Go API routing from OpenAPI specification
- Auto-generated CLI client (`hd1-client`) from specification
- Manual JavaScript UI requiring synchronization with API changes

The **complete evolution** is to **auto-generate the complete web UI client** from the OpenAPI specification, achieving **100% single source of truth** where:
- API routing ← Generated from spec
- CLI client ← Generated from spec  
- **Web UI client ← Generated from spec** ← **THE CROWN JEWEL**

## Decision
Implement **Advanced Auto-Generated Web UI Client System** with complete OpenAPI-driven interface generation.

### Core Architecture
- **UI Generator**: Auto-generate JavaScript API client from OpenAPI spec
- **Component Library**: Generate UI components for each endpoint
- **Form Generation**: Dynamic forms from request schemas
- **Response Handling**: Auto-generated response processors
- **Real-time Integration**: WebSocket client auto-generation

### Single Source of Truth Philosophy
```
┌─────────────────────────────────────────────────┐
│                 api.yaml                        │  ← SINGLE SOURCE OF TRUTH
│           OpenAPI 3.0.3 Specification          │
└─────────────────┬───────────────────────────────┘
                  │
         ┌────────┼────────┐
         │        │        │
         ▼        ▼        ▼
    ┌────────┐ ┌─────────┐ ┌─────────────────┐
    │ Go API │ │ CLI     │ │ Web UI Client   │  ← THE CROWN JEWEL
    │ Router │ │ Client  │ │ (JavaScript)    │
    └────────┘ └─────────┘ └─────────────────┘
```

## Technical Implementation

### Web UI Generator Architecture
```go
type WebUIGenerator struct {
    Spec         *openapi3.T
    OutputDir    string
    TemplateDir  string
    Components   []UIComponent
}

type UIComponent struct {
    Name         string
    Endpoint     string
    Method       string
    RequestForm  FormSchema
    ResponseView ResponseSchema
}
```

### Generated Components
1. **API Client Library**: Complete JavaScript API wrapper
2. **Form Components**: Dynamic forms for all request schemas
3. **Response Viewers**: Formatted displays for all response types
4. **Navigation System**: Auto-generated menus from endpoint tags
5. **Real-time Updates**: WebSocket integration for live data

### Code Generation Strategy
- **Template-Based Generation**: Mustache/Handlebars templates for UI components
- **Schema-Driven Forms**: Generate forms directly from request schemas
- **Type-Safe JavaScript**: Generate TypeScript definitions from schemas
- **Component Modularity**: Each endpoint becomes a reusable UI component

## Generated UI Features

### Standard Holodeck Interface
- **Session Management**: Auto-generated session controls
- **Scene Operations**: Fork/save/load forms from specification
- **Recording Controls**: Temporal sequence UI from recording endpoints
- **Object Management**: Create/update/delete forms from object schemas
- **Real-time Updates**: WebSocket integration with A-Frame rendering

### Example Generated Components
```javascript
// Auto-generated from POST /sessions/{sessionId}/scenes/save
class SaveSceneComponent {
    constructor(apiClient) {
        this.api = apiClient;
        this.generateForm();
    }
    
    generateForm() {
        // Form generated from saveSceneFromSession schema
        this.form = new Form({
            fields: {
                scene_id: { type: 'string', required: true },
                name: { type: 'string', required: true },
                description: { type: 'string', required: true },
                overwrite: { type: 'boolean', default: false }
            }
        });
    }
}
```

## Implementation Phases

### Phase 1: Core Generator Engine
- OpenAPI specification parser enhancement
- Template engine integration
- Basic JavaScript API client generation
- Simple form generation for core endpoints

### Phase 2: Advanced UI Components
- Complex form handling (nested objects, arrays)
- Response formatting and display
- Error handling and validation
- Standard holodeck styling integration

### Phase 3: Real-time Integration
- WebSocket client auto-generation
- A-Frame integration components
- Live session state synchronization
- Standard console enhancement

### Phase 4: Complete Single Source of Truth
- Complete UI replacement with generated components
- Zero manual JavaScript code (except templates)
- 100% specification-driven interface
- Automatic UI updates when spec changes

## Benefits

### Complete Bar-Raising Achievement
✅ **100% Single Source of Truth**: Everything generated from specification
✅ **Zero API Drift**: UI automatically stays in sync with backend
✅ **Instant Feature Development**: New endpoints = automatic UI
✅ **Standard Standards**: Generated code follows consistent patterns
✅ **Type Safety**: Complete type checking from specification to UI

### Development Velocity
- **Instant UI**: New API endpoints automatically generate UI components
- **Zero Maintenance**: UI stays synchronized with API changes
- **Standard Quality**: Generated code follows enterprise standards
- **Complete Testing**: Auto-generated tests for all UI components

## Implementation Notes
- Generator integrates with existing code generation pipeline
- Templates maintain HD1's standard holodeck aesthetic
- Generated components work seamlessly with existing A-Frame integration
- Maintains backward compatibility during transition

## Related ADRs
- ADR-004: Scene Forking System Implementation
- ADR-005: Temporal Recording System Implementation
- ADR-003: Standard UI Enhancement

---
*HD1 - The complete single source of truth in holodeck engineering*