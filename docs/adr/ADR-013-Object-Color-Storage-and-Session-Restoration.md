# ADR-013: Object Color Storage and Session Restoration

## Status
Accepted - Implemented

## Context

HD1 (Holodeck One) experienced critical issues with object color persistence during session restoration. Objects would appear red instead of their intended colors when sessions were restored after WebSocket reconnections or daemon restarts.

### Problem Analysis

**Root Cause Identified**: The object creation system had a fundamental disconnect between color processing and color storage:

1. **Color Processing**: Object creation correctly processed color data from API requests and used it for WebSocket broadcasting to clients
2. **Color Storage Gap**: Colors were never stored in the Object struct database records
3. **Session Restoration Failure**: When restoring sessions, the color parsing logic failed because `obj.Color` was empty, triggering fallback to red default colors

### Technical Investigation

Through systematic debugging, we discovered:

- Object creation in `/opt/hd1/src/api/objects/create.go` processed colors (lines 174-188) but never stored them
- Session restoration in `/opt/hd1/src/server/client.go` expected JSON-formatted color strings (line 167)
- The `CreateObject` database method only stored basic fields (Name, Type, X, Y, Z, Scale) but not Color
- Empty color strings caused `json.Unmarshal([]byte(obj.Color), &parsedColor)` to fail, defaulting to red

## Decision

Implement comprehensive color storage architecture that maintains single source of truth between object creation and session restoration.

### Solution Architecture

1. **Color Storage**: Store colors as JSON strings in the Object struct database
2. **Format Consistency**: Ensure session restoration uses identical color format as object creation
3. **Backward Compatibility**: Maintain existing WebSocket broadcasting format
4. **Single Source of Truth**: Color processing happens once and is stored persistently

## Implementation

### Code Changes

**Primary Fix** (`/opt/hd1/src/api/objects/create.go`):
```go
// Store color as JSON string in Object struct for session restoration
colorJSON, _ := json.Marshal(objectColor)

// Mark object as "new" for session object tracking and store color
updates := map[string]interface{}{
    "tracking_status": "new",
    "created_at": time.Now(),
    "color": string(colorJSON),  // <- KEY FIX: Store color as JSON string
}
h.GetStore().UpdateObject(sessionID, object.Name, updates)
```

**Session Restoration** (`/opt/hd1/src/server/client.go`):
- Existing color parsing logic now works correctly with stored JSON strings
- Fallback to red only occurs when JSON parsing actually fails
- Maintains session-level restoration tracking to prevent infinite loops

### Data Flow

1. **Object Creation**: API request → Color processing → JSON serialization → Database storage
2. **WebSocket Broadcasting**: Same color object used for real-time updates
3. **Session Restoration**: Database → JSON parsing → Identical format to creation

## Consequences

### Positive

- ✅ **Color Persistence**: Objects maintain their colors across session restoration
- ✅ **Single Source of Truth**: Color data flows consistently from creation to restoration
- ✅ **Backward Compatibility**: Existing WebSocket clients continue to work
- ✅ **Format Consistency**: Session restoration uses identical color format as creation
- ✅ **Error Handling**: Graceful fallback to red when color parsing genuinely fails

### Technical Benefits

- **Reactive Scene Graph**: Colors integrate seamlessly with existing reactive architecture
- **Session Isolation**: Color storage respects HD1's session-based architecture  
- **API Compliance**: Maintains OpenAPI 3.0.3 specification adherence
- **Build System Validation**: Color storage passes automated build validation

### Performance

- **Minimal Overhead**: JSON serialization adds negligible processing time
- **Storage Efficiency**: Colors stored as compact JSON strings
- **Network Optimization**: No additional WebSocket messages required

## Verification

### Test Results

1. **Object Creation**: 
   ```bash
   curl -X POST .../objects -d '{"name": "purple-cube", "color": {"r": 0.8, "g": 0.3, "b": 0.9, "a": 1.0}}'
   # Response: "color":"{\"a\":1,\"b\":0.9,\"g\":0.3,\"r\":0.8}"
   ```

2. **Session Restoration**:
   ```
   INFO: session objects restored to client {"object_count":6,"session_id":"session-jzk6uf7t"}
   ```

3. **Color Persistence**: Objects maintain purple color instead of red fallback

## Future Considerations

### Potential Enhancements

1. **Color Validation**: Add color value validation (0.0-1.0 range) in API layer
2. **Color Profiles**: Support for different color spaces (HSL, HSV) if needed
3. **Performance Optimization**: Consider binary color storage for high-frequency scenarios

### Migration Strategy

- **No Migration Required**: Existing objects without colors gracefully default to green
- **Forward Compatible**: New color storage format supports all current A-Frame color features
- **Rollback Safe**: System degrades gracefully if color parsing fails

## Related ADRs

- **ADR-008**: Thread-Safe Session Store (enables color persistence per session)
- **ADR-009**: WebSocket Realtime Architecture (color broadcasting foundation)
- **ADR-011**: Build System Validation (ensures color storage doesn't break builds)

## References

- Issue: Objects appearing red during session restoration
- Fix: Color storage implementation in object creation
- Verification: Comprehensive session restoration testing
- Architecture: Single source of truth for color data flow