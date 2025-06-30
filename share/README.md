# HD1 Share Directory - Asset Organization & Library Systems

## Overview

The `/opt/holodeck-one/share/` directory contains **all shared assets, libraries, and resources** for HD1 (Holodeck One). This follows standard separation of concerns with **precise**.

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
**Purpose**: Web-based 3D visualization client served by HD1 daemon

### **`scenes/`** - Standard Scene Collection
```
scenes/
├── empty-grid.sh       # Clean holodeck foundation
├── basic-shapes.sh     # Fundamental geometric demonstrations
├── anime-ui.sh         # Advanced UI showcase
├── complete-demo.sh    # Complete capabilities demonstration
└── my-custom-scene.sh  # Custom development template
```
**Purpose**: Standard holodeck scene collection accessible via API endpoints

### **`configs/`** - Configuration Templates
```
configs/
└── (configuration templates and defaults)
```
**Purpose**: Standard configuration management

### **`templates/`** - Template System
```
templates/
└── (HTML/rendering templates)
```
**Purpose**: Template-driven content generation

---

## 🎯 **Design Philosophy**

### **Single Source of Truth**
- **Scenes**: Standard scene collection via API integration
- **htdocs**: One web interface for visualization
- **No duplication**: Every asset has exactly one canonical location

### **Standard Organization**
- **Clear separation**: Web assets vs. scenes vs. configuration
- **API-driven scenes**: Complete scene management via /api/scenes endpoints
- **Specification-driven**: All capabilities align with api.yaml

### **xVC Methodology**
- **Pattern consistency**: Repeatable interaction patterns
- **Surgical precision**: Targeted, exact implementations
- **Quality solutions**: Standard-grade only

---

## 🔗 **Integration Points**

### **Scene → HD1 Integration**
```bash
# List available scenes
curl http://localhost:8080/api/scenes
# ↓
# Load scene via API
curl -X POST http://localhost:8080/api/scenes/complete-demo
# ↓  
# HD1 Server executes scene script
# ↓
# Real-time 3D visualization in browser
```

### **Web Interface → HD1 Integration**
```
htdocs/static/js/renderer.js (HolodeckRenderer)
# ↓
# WebSocket connection (/ws)
# ↓
# HD1 Server Hub (real-time updates)
# ↓
# Session store & object management
```

---

## 🛡️ **Standard Standards**

### **File Organization**
- **Absolute paths only**: No relative references
- **Standard naming**: No spaces, clear purposes
- **Version control**: All assets under git control
- **Clean structure**: No temporary or backup files

### **Asset Management**
- **Single canonical location**: Each asset has one source of truth
- **Standard build integration**: Assets validate in build pipeline
- **Documentation**: Every directory has clear purpose documentation

### **Security & Quality**
- **Input validation**: All scripts validate parameters
- **Error handling**: Standard error reporting
- **Resource cleanup**: Proper resource management
- **Standard logging**: Structured, timestamped output

---

## 📚 **Usage Patterns**

### **For Web Development**
```bash
# Serve web interface
cd /opt/holodeck-one/src && make start
# Access: http://localhost:8080/
```

### **For Scene Management**
```bash
# Standard scene management via API
curl http://localhost:8080/api/scenes                    # List scenes
curl -X POST http://localhost:8080/api/scenes/complete   # Load scene
# Or use web interface scene dropdown with cookie persistence
```

### **For Development**
```bash
# Standard build with asset validation
cd /opt/holodeck-one/src && make all
# Assets automatically referenced from share/
```

---

## 🎖️ **Excellence Standards**

HD1 Share directory represents **standard asset organization** with:

✅ **Clear separation of concerns**  
✅ **Single entry point systems**  
✅ **Standard naming conventions**  
✅ **Complete documentation**  
✅ **Integration with specification-driven architecture**  

**"Every asset serves the vision: standard 3D visualization through precise."**

---

*Last Updated: 2025-06-28*  
*HD1 Version: 2.0.0*  
*Authority: Standard Asset Organization Standards*