# ADR-007: Single Source of Truth Avatar Management

**Status**: Accepted  
**Date**: 2025-07-17  
**Authors**: HD1 Development Team  

## Context

HD1's avatar system suffered from complex lifecycle management with multiple control paths, session ambiguity, and inconsistent state tracking. The system had parallel avatar creation/destruction mechanisms that led to orphaned avatars, memory leaks, and synchronization issues across WebSocket connections.

### Problems Identified

1. **Multiple Avatar Sources**: Avatars could be created via movement operations, explicit API calls, and WebSocket connections
2. **Session Complexity**: Multiple ID types (client_id, session_id, avatar_id) created confusion and inconsistency  
3. **Cleanup Ambiguity**: No clear ownership model for avatar lifecycle management
4. **State Drift**: Avatars persisting after client disconnection violated real-time multiplayer expectations

## Decision

We implement **Single Source of Truth Avatar Management** with surgical precision lifecycle control.

### Core Principles

1. **WebSocket Connection = Avatar Lifecycle**: Avatars exist only while WebSocket connection is active
2. **Immediate Registration**: Client registration triggers immediate avatar creation and sync broadcasting
3. **Automatic Cleanup**: WebSocket disconnection immediately triggers avatar removal
4. **Session Inactivity Integration**: Avatar cleanup integrates with session timeout mechanisms
5. **Zero Orphan Policy**: No avatar exists without an active WebSocket connection

### Implementation Strategy

```go
// Avatar Registry - Single Source of Truth
type AvatarRegistry struct {
    avatars map[string]*Avatar  // hd1_id -> Avatar
    hub     *Hub               // WebSocket hub reference
}

// Lifecycle tied to WebSocket events
func (h *Hub) registerClient(client *Client) {
    // Immediate avatar creation and broadcasting
    avatar := h.avatarRegistry.CreateAvatar(client)
    h.broadcastAvatarCreate(avatar)
}

func (h *Hub) unregisterClient(client *Client) {
    // Immediate avatar cleanup and broadcasting  
    h.avatarRegistry.RemoveAvatar(client.GetAvatarID())
    h.broadcastAvatarRemove(client.GetAvatarID())
}
```

### Breaking Changes

- **Avatar Movement**: Movement operations no longer create avatars automatically
- **Explicit Creation**: Removed separate avatar creation API endpoints
- **Session Binding**: Avatars permanently bound to WebSocket connections

## Consequences

### Positive

- **Predictable Lifecycle**: Clear avatar creation/destruction semantics
- **Real-time Accuracy**: Avatar presence accurately reflects client connectivity
- **Memory Efficiency**: Zero memory leaks from orphaned avatars
- **Synchronization Integrity**: All clients see consistent avatar state
- **Performance**: Reduced overhead from lifecycle management complexity

### Negative

- **Breaking Change**: Existing avatar management code requires updates
- **WebSocket Dependency**: Avatar functionality completely dependent on WebSocket infrastructure

## Implementation Timeline

- **2025-07-17 13:00**: Surgical precision avatar movement implementation
- **2025-07-17 14:24**: Complete avatar movement with API client integration
- **2025-07-17 18:52**: Bar-raising single source of truth avatar lifecycle cleanup

## Compliance

This ADR establishes HD1's commitment to architectural purity and single source of truth principles in avatar management, eliminating complexity while ensuring real-time multiplayer integrity.