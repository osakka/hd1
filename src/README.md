# HD1 Source Code

Ultra-minimal Go server implementation for the HD1 Three.js game engine platform.

## 🏗️ Structure

```
src/
├── api.yaml           # OpenAPI specification (single source of truth)
├── main.go            # Server entry point and HTTP handlers
├── auto_router.go     # Auto-generated routing (DO NOT EDIT)
├── Makefile           # Build system and development commands
├── api/               # API endpoint handlers
├── codegen/           # Code generation from specification
├── config/            # Configuration management system
├── logging/           # Structured logging with module tracing
├── server/            # WebSocket hub and client management
└── threejs/           # Three.js bridge and utilities
```

## ⚙️ Build Commands

```bash
# Development workflow
make clean && make && make start

# Individual steps
make generate    # Generate code from api.yaml
make build      # Build server binary
make start      # Start daemon
make stop       # Stop daemon
make logs       # View logs
```

## 🔧 Key Files

- **`api.yaml`** - Complete API specification, drives all code generation
- **`main.go`** - Server initialization, routing setup, static file serving
- **`auto_router.go`** - Auto-generated from api.yaml (never edit manually)
- **`codegen/generator.go`** - Specification-driven code generator

## 📊 Code Statistics

- **4,989 lines** of Go code total
- **538 lines** in code generator (after optimization)
- **27 core source files** in entire project
- **11 active API routes** auto-generated from specification

## 🛡️ Code Quality

- **Zero manual routing** - All routes generated from specification
- **Ultra-minimal build** - Only essential code included
- **Specification-driven** - Single source of truth in api.yaml
- **Development-focused** - Comprehensive error handling and logging

---

*Source code for HD1 v0.7.0 - Specification-driven Three.js platform*