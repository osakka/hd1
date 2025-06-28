# THD Share Directory - Asset Organization & Library Systems

## Overview

The `/opt/holo-deck/share/` directory contains **all shared assets, libraries, and resources** for THD (The Holo-Deck). This follows professional separation of concerns with **surgical precision**.

---

## 📁 **Directory Structure**

### **`htdocs/`** - Web Interface Assets
```
htdocs/
├── static/js/          # JavaScript: HolodeckRenderer, gl-matrix, debug
├── debug.html          # Direct object testing interface
├── force-session.html  # Session management utilities
└── assets/             # Static web resources
```
**Purpose**: Web-based 3D visualization client served by THD daemon

### **`glibsh/`** - GL Library Shell (GLIBSH)
```
glibsh/
├── glibsh              # Single entry point script
├── core/               # Core orchestration primitives
├── objects/            # Primitive object creation
├── shapes/             # Composite shape patterns  
├── scenes/             # Complete scene compositions
└── capabilities/       # High-level automation
```
**Purpose**: **Single source of truth** for THD automation - professional shell library for 3D operations

### **`configs/`** - Configuration Templates
```
configs/
└── (configuration templates and defaults)
```
**Purpose**: Professional configuration management

### **`templates/`** - Template System
```
templates/
└── (HTML/rendering templates)
```
**Purpose**: Template-driven content generation

---

## 🎯 **Design Philosophy**

### **Single Source of Truth**
- **GLIBSH**: One entry point for all THD automation
- **htdocs**: One web interface for visualization
- **No duplication**: Every asset has exactly one canonical location

### **Professional Organization**
- **Clear separation**: Web assets vs. automation vs. configuration
- **Hierarchical structure**: From primitives → composites → complete systems
- **Specification-driven**: All capabilities align with api.yaml

### **xVC Methodology**
- **Pattern consistency**: Repeatable interaction patterns
- **Surgical precision**: Targeted, exact implementations
- **Bar-raising solutions**: Professional-grade only

---

## 🔗 **Integration Points**

### **GLIBSH → THD Integration**
```bash
share/glibsh/glibsh create-scene demo
# ↓
# Orchestrates multiple thd-client calls
# ↓  
# THD Server API (specification-driven)
# ↓
# Real-time 3D visualization in browser
```

### **Web Interface → THD Integration**
```
htdocs/static/js/renderer.js (HolodeckRenderer)
# ↓
# WebSocket connection (/ws)
# ↓
# THD Server Hub (real-time updates)
# ↓
# Session store & object management
```

---

## 🛡️ **Professional Standards**

### **File Organization**
- **Absolute paths only**: No relative references
- **Professional naming**: No spaces, clear purposes
- **Version control**: All assets under git control
- **Clean structure**: No temporary or backup files

### **Asset Management**
- **Single canonical location**: Each asset has one source of truth
- **Professional build integration**: Assets validate in build pipeline
- **Documentation**: Every directory has clear purpose documentation

### **Security & Quality**
- **Input validation**: All scripts validate parameters
- **Error handling**: Professional error reporting
- **Resource cleanup**: Proper resource management
- **Professional logging**: Structured, timestamped output

---

## 📚 **Usage Patterns**

### **For Web Development**
```bash
# Serve web interface
cd /opt/holo-deck/src && make start
# Access: http://localhost:8080/
```

### **For Automation/Scripting**
```bash
# Single entry point for all operations
/opt/holo-deck/share/glibsh/glibsh --help
/opt/holo-deck/share/glibsh/glibsh create-object cube --position 0,0,0
/opt/holo-deck/share/glibsh/glibsh create-scene gallery
```

### **For Development**
```bash
# Professional build with asset validation
cd /opt/holo-deck/src && make all
# Assets automatically referenced from share/
```

---

## 🎖️ **Excellence Standards**

THD Share directory represents **professional asset organization** with:

✅ **Clear separation of concerns**  
✅ **Single entry point systems**  
✅ **Professional naming conventions**  
✅ **Complete documentation**  
✅ **Integration with specification-driven architecture**  

**"Every asset serves the vision: professional 3D visualization through surgical precision."**

---

*Last Updated: 2025-06-28*  
*THD Version: 2.0.0*  
*Authority: Professional Asset Organization Standards*