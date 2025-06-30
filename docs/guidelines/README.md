# HD1 Development Guidelines

**Standard Standards for HD1 (Holodeck One) Development**

This directory contains comprehensive development guidelines and standards for maintaining the standard quality and architectural integrity of HD1.

## 📋 Guidelines Overview

### [Development Standards](development-standards.md)
Core development principles, coding standards, and standard engineering practices for HD1 development.

### [API Design Guidelines](api-design-guidelines.md)
Standards for OpenAPI specification design, endpoint naming, and API evolution.

### [Testing Standards](testing-standards.md)
Comprehensive testing practices, validation procedures, and quality assurance protocols.

### [Documentation Standards](documentation-standards.md)
Standard documentation practices, ADR processes, and technical writing guidelines.

## 🎯 Core Principles

### Single Source of Truth
- OpenAPI 3.0.3 specification drives all development
- Auto-generated routing and clients eliminate manual synchronization
- One authoritative source for all system behavior

### Standard Engineering Standards
- Quality solutions only - every decision elevates system quality
- Zero regressions through comprehensive validation
- High-quality logging, monitoring, and error handling

### Advanced Architecture
- Specification-driven development eliminating manual configuration
- Upstream/downstream API bridge maintaining perfect integration
- Standard VR/AR holodeck platform with WebXR standards

## 🔧 Development Workflow

### 1. Specification First
```bash
# Edit OpenAPI specification
vim src/api.yaml

# Regenerate all routing and clients
make generate

# Validate complete system
make validate
```

### 2. Standard Standards
```bash
# Build with validation
make all

# Standard daemon control
make start
make status
make stop
```

### 3. Quality Assurance
```bash
# Run comprehensive tests
make test

# Validate API specification
make validate

# Check standard standards
make lint
```

## 📚 Quick Reference

### Code Quality
- **No Emojis**: Standard output only
- **Absolute Paths**: All configurations use absolute paths
- **Long Flags**: No short flags to eliminate confusion
- **Standard Logging**: Timestamped, structured, actionable

### API Standards
- **OpenAPI 3.0.3**: Single source of truth specification
- **Auto-Generated**: All routing and clients from specification
- **Standard Validation**: Comprehensive error handling
- **Backward Compatibility**: Zero regressions maintained

### VR/AR Integration
- **A-Frame WebXR**: Complete local 2.5MB ecosystem
- **Standard Boundaries**: [-12, +12] holodeck containment
- **60fps Monitoring**: Real-time position tracking
- **Session Isolation**: Thread-safe multi-user support

---

*These guidelines ensure HD1 maintains its position as the advanced standard VR/AR holodeck platform with high-quality engineering standards.*