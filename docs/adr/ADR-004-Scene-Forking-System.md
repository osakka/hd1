# ADR-004: Scene Forking System Implementation

## Status
**ACCEPTED** - Implemented 2025-06-29

## Context
THD needed a revolutionary approach to scene management that enables users to:
- Load existing scenes for modification without affecting originals
- Save modified sessions as new scenes  
- Track object lifecycle and provenance
- Implement professional "photo vs video" paradigm for content creation

## Decision
Implement a comprehensive **Scene Forking System** with the following components:

### Core Architecture
- **Fork API**: `POST /scenes/{sceneId}/fork` - Load scene into session for editing
- **Save API**: `POST /sessions/{sessionId}/scenes/save` - Export session state as new scene
- **Object Tracking**: Three-state system (base/modified/new)
- **Script Generation**: Dynamic creation of executable scene files

### Photo vs Video Paradigm
- **📷 Photo Mode**: Current state snapshots saved as scenes
- **🎥 Video Mode**: Temporal recording for complete session playback
- **Professional UI**: Console controls for both modes

## Technical Implementation

### Scene Fork Process
1. Clear existing session objects (optional)
2. Execute source scene script in target session
3. Mark all loaded objects as "base" with source scene reference
4. Enable modification tracking for all subsequent changes

### Object Lifecycle Tracking
```go
type Object struct {
    // ... existing fields
    TrackingStatus string    `json:"tracking_status,omitempty"` // "base", "modified", "new"
    SourceScene    string    `json:"source_scene,omitempty"`    // Original scene ID
    CreatedAt      time.Time `json:"created_at,omitempty"`      // Timestamp
}
```

### Scene Script Generation
- Parse session objects into executable bash scripts
- Include metadata headers for dynamic discovery
- Use auto-generated thd-client for object creation
- Maintain 100% API compatibility

## Consequences

### Positive
✅ **Professional Content Creation**: Industry-standard photo/video workflow
✅ **Non-destructive Editing**: Original scenes never modified
✅ **Complete Audit Trail**: Full object provenance tracking
✅ **API-Driven Architecture**: Single source of truth maintained
✅ **Dynamic Discovery**: No hardcoded scene metadata

### Negative
⚠️ **Increased Complexity**: More API endpoints and object metadata
⚠️ **Storage Requirements**: Additional tracking fields per object

## Implementation Notes
- All scene operations maintain session isolation
- Fork operations support selective object clearing
- Generated scripts include professional headers and metadata
- Object tracking integrates seamlessly with existing APIs

## Related ADRs
- ADR-001: A-Frame WebXR Integration
- ADR-003: Professional UI Enhancement

---
*THD - Where immersive holodeck technology meets professional engineering*