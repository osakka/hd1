# HD1 (Holodeck One) - Three.js Game Engine Platform

HD1 is an API-first, specification-driven Three.js game engine platform that provides complete 3D game development capabilities through REST endpoints with real-time WebSocket synchronization.

## 🚀 Quick Start

```bash
# Build and start HD1
cd src && make clean && make && make start

# Access the console
open http://localhost:8080

# Check API status
curl http://localhost:8080/api/system/version
```

## 🏗️ Architecture

- **API-First Design**: 86 REST endpoints auto-generated from OpenAPI specification
- **Real-Time Sync**: WebSocket hub with TCP-simple sequence-based synchronization  
- **Three.js Integration**: Direct WebGL rendering with zero abstraction layers
- **Specification-Driven**: Single source of truth in `src/api.yaml`
- **Ultra-Minimal Build**: Optimized codebase with only essential components

## 📁 Project Structure

```
/opt/hd1/
├── src/           # Go server source code
├── share/         # Static assets and configuration
├── docs/          # Complete documentation
├── build/         # Build artifacts and binaries
└── CLAUDE.md      # Development context and principles
```

## 🛠️ Core Features

- **Three.js Console**: Ultra-minimal debug panel with WebSocket monitoring
- **Rebootstrap**: Intelligent recovery system clearing storage on connection failures  
- **Auto-Generation**: Complete routing, client libraries, and API documentation
- **Configuration Management**: Environment variables, flags, and .env file support
- **Production Ready**: Comprehensive logging, error handling, and performance optimization

## 📖 Documentation

- **[Architecture Overview](docs/architecture/README.md)** - System design and components
- **[User Guides](docs/guides/)** - Development and deployment guides
- **[ADR](docs/adr/)** - Architectural decision records
- **[API Reference](src/api.yaml)** - Complete OpenAPI specification

## 🔧 Development

HD1 follows specification-driven development where `src/api.yaml` is the single source of truth:

```bash
# Generate code from specification
make generate

# Run development build
make dev

# View logs
make logs
```

## 📊 Status

**Current Version**: v6.0.0  
**Build Status**: ✅ Production Ready  
**API Endpoints**: 11 active routes  
**Code Quality**: Ultra-minimal optimized build

## 📄 License

Professional holodeck platform - See documentation for details.

---

*HD1 v6.0.0: Where OpenAPI specifications become immersive Three.js game worlds.*