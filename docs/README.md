# HD1 Documentation Center

**Three-Layer Game Engine Architecture Documentation for HD1 (Holodeck One)**

Welcome to the comprehensive documentation center for HD1 - the world's first three-layer game engine holodeck platform combining Environment + Props + Scene architecture with specification-driven development and industry-standard VR/AR capabilities.

## 📚 Documentation Structure

### 📚 [User Guides](guides/)
**Getting Started & User Documentation**
- [Getting Started Guide](guides/getting-started.md) - Quick start for new users
- Installation, configuration, and basic usage
- Step-by-step tutorials and examples

### 🔧 [API Documentation](api/)
**Complete API Reference & Integration**
- [API Reference](api/README.md) - Complete endpoint documentation
- [Enhanced Integration](api/enhanced-integration.md) - Advanced API patterns
- Authentication, examples, and best practices

### 🏗️ [Architecture](architecture/)
**System Design & Technical Foundation**
- [Architecture Overview](architecture/overview.md) - Three-layer system design
- [System Architecture](architecture/system-architecture.md) - Detailed system design
- [Coordinate System](architecture/coordinate-system.md) - 3D space implementation

### 👩‍💻 [Development](development/)
**Developer Resources & Standards**
- [Contributing Guide](development/contributing.md) - How to contribute to HD1
- [Three-Layer Implementation](development/three-layer-implementation.md) - Implementation plan
- [Session Handoff](development/session-handoff.md) - Development context

### 📜 [Architecture Decision Records (ADRs)](adr/)
**Architectural Evolution Timeline**
- [ADR Timeline and Index](adr/README.md) - Complete chronological decision history
- Advanced milestones and standard standards implementation
- Impact assessment and architectural principles
- **Latest**: [ADR-014](adr/ADR-014-Three-Layer-Architecture-Environment-Props-System.md) - Three-Layer Architecture Environment + Props System

### 📋 [Guidelines](guidelines/)
**Development Standards & Best Practices**
- [Development Standards](guidelines/development-standards.md) - Core engineering principles
- Standard coding standards and quality assurance
- API design, testing, and documentation guidelines

### 📈 [Trackers](trackers/)
**Work Tracking & Progress Monitoring**
- Feature development tracking
- Performance monitoring and optimization goals
- Quality assurance milestones

## 🎯 Quick Start Guide

### For Developers
1. **Read the Whitepaper**: Start with [Advanced Holodeck Architecture](whitepaper/advanced-holodeck-architecture.md)
2. **Review Guidelines**: Check [Development Standards](guidelines/development-standards.md)
3. **Explore ADRs**: Understand architectural evolution in [ADR Timeline](adr/README.md)
4. **API Reference**: Use [API Documentation](design/api/) for implementation

### For Architects
1. **System Architecture**: Review [Architecture Documentation](design/architecture/)
2. **Decision History**: Study [ADR Timeline](adr/README.md) for context
3. **Technical Foundation**: Deep dive into [Advanced Architecture Whitepaper](whitepaper/advanced-holodeck-architecture.md)

### For Project Managers
1. **Progress Tracking**: Monitor [Trackers](trackers/) for current status
2. **Quality Standards**: Review [Guidelines](guidelines/) for quality metrics
3. **Impact Assessment**: Check [ADR Status Overview](adr/README.md#-adr-status-overview)

## 🚀 Advanced Features

### Three-Layer Game Engine Architecture
- **Environment System**: 4 physics contexts with realistic material adaptation
- **Props System**: 6 categories of reusable objects with YAML-based definitions
- **Scene Orchestration**: Smart composition (future Phase 3)
- **Game Engine Parity**: Matches Unity/Unreal Engine architectural patterns

### Single Source of Truth Architecture
- **OpenAPI 3.0.3 Driven**: All routing and clients auto-generated (31 endpoints)
- **Zero Manual Synchronization**: Perfect specification-to-implementation consistency
- **Three-Layer APIs**: Environment and Props APIs auto-generated from specification

### Physics Cohesion Engine
- **Environment-Aware Props**: Props automatically adapt physics based on session environment
- **Real-time Adaptation**: Physics recalculated instantly on environment changes
- **Material Accuracy**: Realistic properties (wood: 600 kg/m³, metal: 7800 kg/m³)

### Standard VR/AR Platform
- **100% Local A-Frame**: Complete 2.5MB WebXR ecosystem
- **Holodeck Containment**: Escape-proof [-12, +12] boundaries
- **Enterprise Quality**: Thread-safe session management

## 📈 Version Information

**Current Version**: v4.0.0 - Three-Layer Architecture Revolution

**Key Milestones**:
- v4.0.0: Three-layer game engine architecture (Environment + Props + Scenes)
- v3.6.0: Simple scene updates & precise code audit
- v3.5.0: Advanced API bridge system
- v3.4.0: Auto-generated web UI client
- v3.3.0: Standard UI excellence & scene management
- v3.2.0: Scene forking and temporal recording
- v3.1.0: A-Frame WebXR integration
- v3.0.0: Specification-driven development foundation

## 🔗 External Resources

- **Git Repository**: `https://git.uk.home.arpa/itdlabs/holo-deck.git`
- **Project Module**: `holodeck` (Go module)
- **Binary**: `hd1` (standard daemon)
- **Runtime Location**: `/opt/hd1/`

---

**HD1: Where three-layer game engine architecture meets immersive holodeck technology.**

*This documentation represents the complete technical foundation for the revolutionary HD1 platform - the world's first three-layer game engine holodeck with Environment + Props + Scene architecture and specification-driven development.*