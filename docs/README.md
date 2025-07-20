# HD1 Documentation

Documentation for HD1 (Holodeck One) - **experimental WebGL REST platform** exploring HTTP-to-3D concepts.

⚠️ **EXPERIMENTAL PROJECT** - This is research/prototype code, not production software!

## 📚 Documentation Structure

### 🔍 Navigation & Reference
- **[VERSION.md](VERSION.md)** - Single source of truth for version information
- **[INDEX.md](INDEX.md)** - Complete cross-reference index
- **[API Endpoints](api/ENDPOINTS.md)** - Current API reference (~40 endpoints, many incomplete)

### 📋 Experimental Development
- **[Development Context](../CLAUDE.md)** - Current technical implementation details
- **[Changelog](../CHANGELOG.md)** - Actual changes and current limitations

### 🏗️ Architecture
- **[System Overview](architecture/overview.md)** - High-level system design
- **[API Specification](../src/schemas/hd1-api.yaml)** - Current experimental API specification

### 📖 User Guides
- **[Quick Start](guides/quick-start.md)** - Get up and running
- **[Development Guide](guides/development.md)** - Building and extending HD1
- **[Configuration Guide](guides/configuration.md)** - Environment and settings
- **[Troubleshooting](guides/troubleshooting.md)** - Common issues and solutions

### Architectural Decision Records (ADR)
⚠️ **Note**: Many ADRs contain aspirational content that doesn't match current experimental reality
- **[ADR-001](adr/001-api-first-architecture.md)** - API-First Architecture
- **[ADR-002](adr/002-specification-driven-development.md)** - Specification-Driven Development
- **[ADR-003](adr/003-threejs-minimal-console.md)** - Three.js Minimal Console

### Experimental Status
- **Current State**: ~40 API endpoints, many incomplete or non-functional
- **Working Features**: Basic Three.js operations, WebSocket sync, configuration
- **Known Issues**: Limited error handling, performance not optimized, mobile UX rough

## 🔍 Quick Navigation

### Current Experimental Platform (v1.0.0)
**For Developers:**
- Start with [Quick Start Guide](guides/quick-start.md)
- Review [Development Guide](guides/development.md)
- Understand [System Overview](architecture/overview.md)

**For Architects:**
- Read [System Overview](architecture/overview.md)
- Review all [ADRs](adr/) for design decisions
- Study [API Specification](../src/api.yaml)

**For Operations:**
- Configure via [Configuration Guide](guides/configuration.md)
- Use [Troubleshooting Guide](guides/troubleshooting.md)

### Experimental Development
**For Researchers:**
- Review [Current Limitations](../CHANGELOG.md#what-s-broken-missing)
- Study [Working Features](../CHANGELOG.md#what-actually-works)
- Understand HTTP-to-3D exploration approach

**For Contributors:**
- Check [Development Context](../CLAUDE.md) for technical details
- Review [XVC Development Dependency](../README.md#development)
- Understand experimental nature and expect bugs

## 📖 Document Conventions

- **Strategic docs**: Business plans, transformation strategy, and vision
- **Architecture docs**: Technical system design and components
- **Guides**: Step-by-step practical instructions
- **ADRs**: Historical architectural decisions with context and rationale
- **Implementation plans**: Detailed phase-by-phase technical specifications

## 🔄 Experimental Evolution

### Current State (v1.0.0)
- **~40 API Endpoints**: Basic Three.js operations, many incomplete
- **WebSocket Sync**: Simple real-time updates for basic operations
- **Configuration**: Environment variables, command line flags working
- **Experimental Quality**: Expect bugs, crashes, and missing features

### Development Reality
- **Research Project**: Exploring HTTP-to-3D concepts
- **XVC Dependency**: Developed with experimental versioning tool
- **Limited Scope**: Basic 3D object creation and scene management
- **No Timeline**: Experimental development with no specific roadmap

## 🔄 Keeping Documentation Accurate

This documentation should reflect the experimental reality of HD1 v1.0.0. When making changes:

1. Update guides to match actual working features
2. Document what actually works vs what's broken
3. Remove aspirational content that doesn't match reality
4. Keep experimental warnings prominent
5. Maintain honest assessment in changelog

---

*Documentation for HD1 v1.0.0 - Experimental WebGL REST platform exploring HTTP-to-3D concepts*  
*This is research/prototype code - expect bugs, limitations, and missing features!*