# HD1 Logging Standards

## Overview

HD1 implements unified, structured logging with strict standards for message quality, audience targeting, and operational effectiveness. All logging follows the format:

```
timestamp [processid:threadid] [level] functionname.filename:line_num message
```

## Log Levels & Audience

### TRACE
- **Audience**: Developers debugging specific functionality
- **Usage**: Module-specific tracing (websocket, entities, sessions, sync)
- **Performance**: Zero overhead when disabled
- **Control**: `HD1_TRACE_MODULES=websocket,entities`
- **Example**: `logging.Trace("websocket", "client connection established")`

### DEBUG  
- **Audience**: Developers in development environment
- **Usage**: Detailed operation flow, internal state changes
- **Performance**: Zero overhead when disabled
- **Example**: `logging.Debug("entity component updated")`

### INFO
- **Audience**: Production users, SREs, system administrators
- **Usage**: Successful operations, state changes, system events
- **Performance**: Always enabled in production
- **Example**: `logging.Info("session created", data)`

### WARN
- **Audience**: Production users, SREs requiring attention
- **Usage**: Recoverable issues, degraded performance, missing resources
- **Performance**: Always enabled
- **Example**: `logging.Warn("asset loading timeout - using fallback")`

### ERROR
- **Audience**: Production users, SREs requiring immediate action
- **Usage**: Operation failures, invalid requests, system problems
- **Performance**: Always enabled
- **Example**: `logging.Error("database connection failed", data)`

### FATAL
- **Audience**: All users - system terminating
- **Usage**: Critical failures causing process termination
- **Performance**: Always enabled, calls os.Exit(1)
- **Example**: `logging.Fatal("configuration initialization failed", data)`

## Message Quality Standards

### ✅ Good Messages
- **Actionable**: Clear what happened and what to do
- **Concise**: No redundant information (filename/function already in format)
- **Audience-appropriate**: Match the log level audience
- **Specific**: Include relevant context data

```go
// Good Examples
logging.Info("session created", map[string]interface{}{
    "session_id": sessionID,
})

logging.Error("asset loading failed", map[string]interface{}{
    "asset_path": path,
    "error": err.Error(),
})

logging.Warn("high memory usage detected", map[string]interface{}{
    "usage_mb": memoryMB,
    "threshold_mb": thresholdMB,
})
```

### ❌ Bad Messages  
- **Redundant**: Including function/file info already in format
- **Non-actionable**: Vague descriptions
- **Wrong audience**: Technical details at INFO level
- **Inconsistent format**: Using SUCCESS:, [STATUS], etc.

```go
// Bad Examples - DO NOT USE
logging.Info("SUCCESS: CreateEntityHandler completed successfully") // Redundant
logging.Info("[WEBSOCKET] Connection established in hub.go:123")    // Wrong format
logging.Error("Error occurred while processing request")            // Non-actionable
logging.Debug("Server starting on port 8080")                      // Wrong level
```

## Format Standards

### Message Structure
```
Short, actionable message without redundant context
```

### Data Structure
```go
map[string]interface{}{
    "key": value,  // Relevant context only
}
```

### Avoid These Patterns
- ❌ `"SUCCESS: operation completed"`
- ❌ `"[MODULE] message"`  
- ❌ `"Error in function_name()"`
- ❌ `"Processing request in handler"`
- ❌ Including filename/line numbers in message

## Runtime Configuration

### Environment Variables
```bash
HD1_LOG_LEVEL=INFO                    # TRACE,DEBUG,INFO,WARN,ERROR,FATAL
HD1_TRACE_MODULES=websocket,entities  # Comma-separated module list
HD1_LOG_DIR=/opt/hd1/build/logs      # Log directory path
```

### Command Line Flags
```bash
--log-level=DEBUG
--trace-modules=websocket,entities,sessions
--log-dir=/custom/log/path
```

### API Control
```bash
# Change log level at runtime
curl -X POST /api/admin/logging/level -d '{"level":"DEBUG"}'

# Enable/disable trace modules
curl -X POST /api/admin/logging/trace -d '{"enable":["websocket","entities"]}'

# Get current configuration
curl /api/admin/logging/config
```

## Module-Based Tracing

### Available Modules
- `websocket` - WebSocket connection management
- `entities` - Entity lifecycle operations  
- `sessions` - Session management
- `sync` - HD1-VSC synchronization protocol
- `avatars` - Avatar system operations
- `physics` - Physics engine integration

### Usage
```go
// Enable trace for specific functionality
logging.Trace("websocket", "connection established", map[string]interface{}{
    "client_id": clientID,
    "session_id": sessionID,
})

// Check if tracing is enabled before expensive operations
if logging.IsTraceEnabled("entities") {
    // Expensive debug data collection
    logging.Trace("entities", "component state", expensiveData)
}
```

## Thread Safety

All logging operations are thread-safe with minimal lock contention:
- **Single RWMutex**: One lock for level checking
- **Zero-overhead disabled levels**: ~50ns CPU when disabled
- **Atomic operations**: No race conditions
- **Buffer pooling**: Memory-efficient with object reuse

## Performance Characteristics

### Disabled Log Levels
- **CPU Overhead**: ~50ns per call
- **Memory Overhead**: Zero allocations
- **Lock Contention**: Single RWMutex read lock

### Enabled Log Levels  
- **Structured JSON**: File logging with full context
- **Human-readable**: Console output with timestamp
- **Log rotation**: Automatic at 10MB with 3 rotations
- **Memory pooling**: Reused buffers for efficiency

## Implementation Guidelines

### 1. Choose Appropriate Level
```go
// ✅ Correct level usage
logging.Info("user authenticated")              // Production event
logging.Warn("retry attempt failed")           // Recoverable issue  
logging.Error("payment processing failed")     // Requires action
logging.Debug("cache miss occurred")           // Development detail
logging.Trace("websocket", "frame received")   // Module debugging
```

### 2. Provide Context Data
```go
// ✅ Include relevant context
logging.Error("database query failed", map[string]interface{}{
    "query": sql,
    "duration_ms": duration.Milliseconds(),
    "error": err.Error(),
})
```

### 3. Avoid Redundancy
```go
// ❌ Bad - redundant information
logging.Info("CreateSession function created session successfully in sessions.go")

// ✅ Good - concise and actionable
logging.Info("session created", map[string]interface{}{
    "session_id": sessionID,
})
```

### 4. Use Conditional Tracing
```go
// ✅ Efficient tracing with conditional checks
if logging.IsTraceEnabled("sync") {
    expensiveData := collectSyncMetrics()
    logging.Trace("sync", "delta operation applied", expensiveData)
}
```

## Examples by Use Case

### API Request Processing
```go
// Request start (TRACE for debugging)
logging.Trace("entities", "create entity request received", map[string]interface{}{
    "session_id": sessionID,
    "entity_name": request.Name,
})

// Successful completion (INFO for production)
logging.Info("entity created", map[string]interface{}{
    "entity_id": entity.ID,
    "session_id": sessionID,
})

// Validation error (WARN for client errors)
logging.Warn("invalid entity name", map[string]interface{}{
    "name": request.Name,
    "validation_error": "name too long",
})

// System error (ERROR for server issues)
logging.Error("entity creation failed", map[string]interface{}{
    "session_id": sessionID,
    "error": err.Error(),
})
```

### WebSocket Operations
```go
// Connection events (INFO for monitoring)
logging.Info("client connected", map[string]interface{}{
    "client_count": len(hub.clients),
})

// Protocol details (TRACE for debugging)
logging.Trace("websocket", "message broadcast", map[string]interface{}{
    "message_type": messageType,
    "recipient_count": clientCount,
})

// Performance issues (WARN for attention)
logging.Warn("broadcast queue full", map[string]interface{}{
    "queue_size": len(broadcast),
    "max_size": maxSize,
})
```

### System Operations
```go
// Startup/shutdown (INFO for operations)
logging.Info("server starting", map[string]interface{}{
    "bind_address": bindAddr,
    "version": config.GetVersion(),
})

// Configuration changes (INFO for operations)
logging.Info("log level changed", map[string]interface{}{
    "old_level": oldLevel,
    "new_level": newLevel,
})

// Resource limits (WARN for capacity planning)
logging.Warn("memory usage high", map[string]interface{}{
    "usage_percent": usagePercent,
    "threshold_percent": 85,
})
```

---

**Single Source of Truth**: This document defines the complete logging standards for HD1. All log messages must conform to these guidelines for consistency, operational effectiveness, and proper audience targeting.