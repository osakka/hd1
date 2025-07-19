# ADR-008: Unified HD1 ID System

**Status**: Accepted  
**Date**: 2025-07-18  
**Authors**: HD1 Development Team  

## Context

HD1 suffered from architectural inconsistency with multiple parallel ID systems (`client_id`, `session_id`, `avatar_id`) that created confusion, maintenance overhead, and violated single source of truth principles. Different parts of the system used different ID types, leading to mapping complexity and potential synchronization issues.

### Problems Identified

1. **Multiple ID Types**: Three different ID systems serving overlapping purposes
2. **Mapping Complexity**: Constant translation between ID types across system layers
3. **Inconsistent APIs**: Different endpoints expected different ID types
4. **Logging Confusion**: Mixed ID types in logs made debugging difficult
5. **Client Ambiguity**: Frontend code switching between ID types unpredictably

### System Analysis

```
BEFORE (Multiple ID Systems):
WebSocket: client_id → Avatar: avatar_id → Session: session_id
         ↓             ↓                   ↓
   API Headers    JSON Responses     Database Queries
```

## Decision

We implement **Unified HD1 ID System** with complete elimination of parallel ID types.

### Core Principles

1. **Single ID Type**: `hd1_id` replaces all client_id, session_id, avatar_id references
2. **End-to-End Consistency**: Same ID flows from WebSocket → API → Database → Logging
3. **Protocol Unification**: All WebSocket messages use `hd1_id` field exclusively
4. **Template Integration**: Auto-generation produces unified ID system code
5. **Zero Legacy**: Complete elimination of old ID systems without compatibility layers

### Implementation Strategy

```
AFTER (Unified HD1 ID System):
WebSocket: hd1_id → Avatar: hd1_id → Session: hd1_id → Database: hd1_id
         ↓
   Single source of truth flows through entire stack
```

#### WebSocket Protocol Changes
```json
// Client initialization
{"type": "client_init", "hd1_id": "hd1-12345-67890"}

// Avatar operations  
{"type": "avatar_move", "data": {"hd1_id": "hd1-12345-67890", ...}}
```

#### API Header Standardization
```javascript
headers: {
    'Content-Type': 'application/json',
    'X-Client-ID': this.hd1Id  // Content is hd1_id value
}
```

#### Database Consistency
```sql
SELECT user_id::text as hd1_id FROM participants WHERE...
```

### Breaking Changes

- **WebSocket Messages**: All `client_id` fields renamed to `hd1_id`
- **Avatar JSON**: Avatar serialization uses `hd1_id` instead of `client_id`
- **Database Aliases**: Query result aliases changed from `client_id` to `hd1_id`
- **Logging Fields**: All log entries use `hd1_id` field names

## Implementation Timeline

- **2025-07-18 13:43**: Initial unified HD1 ID system implementation
- **2025-07-18 13:59**: Legacy session association system elimination
- **2025-07-18 14:16**: Complete unified ID system with operation broadcasting
- **2025-07-19 10:01**: Complete unified hd1_id system across entire stack
- **2025-07-19 11:44**: Surgical unification with template system fixes

## Template System Integration

Critical fix: Updated auto-generation templates to prevent regression to old ID systems:

```javascript
// Template generates unified code
constructor(baseURL = '/api', hd1Id = null) {
    this.hd1Id = hd1Id;
}

setHd1Id(hd1Id) {
    this.hd1Id = hd1Id;  
}
```

## Consequences

### Positive

- **Single Source of Truth**: Complete elimination of ID system ambiguity
- **Reduced Complexity**: No ID type mapping or translation required
- **Consistent Debugging**: Unified logging with single ID type
- **Maintainable Code**: Single ID pattern throughout codebase
- **Template Consistency**: Auto-generation maintains architectural purity

### Negative

- **Breaking Changes**: Major protocol changes requiring client updates
- **Migration Overhead**: All existing integrations require updates

## Compliance

This ADR represents the pinnacle of single source of truth architecture, achieving complete consistency without compromising functionality. The unified HD1 ID system establishes HD1 as a reference implementation for distributed 3D platform architecture.