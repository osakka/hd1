# HD1 Share Directory - v5.0.0 Native PC Architecture

## Overview

The `/opt/hd1/share/` directory contains **production assets and configuration** for HD1 v5.0.0 API-first game engine platform. This represents 100% native PC architecture with PlayCanvas integration and single source of truth design.

---

## 📁 **Directory Structure (v5.0.0)**

### **`channels/`** - Channel-Based Scene Configuration
```
channels/
├── channel_one.yaml     # Primary collaborative environment
├── channel_two.yaml     # Advanced physics simulation  
├── channel_three.yaml   # Underwater exploration
└── config.yaml          # Channel server configuration
```
**Purpose**: YAML-based scene definitions for PlayCanvas game engine  
**API Integration**: `/api/channels` endpoints read these configurations  
**Architecture**: Single source of truth for collaborative environments

### **`props/`** - Reusable Entity Definitions
```
props/
├── decorative/          # Visual elements (flower-pot.yaml)
├── electronic/          # Electronic devices (lightbulb.yaml)
├── furniture/           # Furniture items (wooden-chair.yaml, wooden-table.yaml)
├── lighting/            # Lighting elements (hd1-test-lighting.yaml)
├── structural/          # Building elements (concrete-pillar.yaml)
└── tools/              # Tool objects (hammer.yaml)
```
**Purpose**: YAML prop definitions with PlayCanvas component specifications  
**API Integration**: `/api/props` and `/api/sessions/{id}/props/{propId}` endpoints  
**Architecture**: Reusable entity-component-system definitions

### **`htdocs/`** - Web Interface Assets
```
htdocs/
├── index.html           # Main PlayCanvas web interface
├── static/js/           # Auto-generated JavaScript API clients
├── static/css/          # Console styling
├── assets/              # Static web resources (audio, models, textures)
└── debug.html           # Development testing interface
```
**Purpose**: PlayCanvas web client served by HD1 daemon  
**API Integration**: Consumes 85 REST endpoints via auto-generated clients  
**Architecture**: Professional 3D rendering with real-time WebSocket sync

---

## 🎯 **V5.0.0 Architecture Principles**

### **Single Source of Truth**
- **Channel Configuration**: YAML files drive PlayCanvas scene setup
- **Prop Definitions**: YAML specifications define reusable entities
- **API Generation**: All clients auto-generated from api.yaml specification
- **No Duplication**: Every asset has exactly one canonical definition

### **PlayCanvas Native**
- **Entity-Component-System**: Native PlayCanvas architecture
- **Professional 3D**: WebGL game engine rendering
- **Real-Time Sync**: <10ms WebSocket synchronization
- **85 API Endpoints**: Complete game engine control via REST

### **Native PC Optimization**
- **No Legacy Systems**: A-Frame completely removed
- **No Shell Scripts**: YAML configuration only
- **Clean Dependencies**: Minimal, purpose-built asset structure
- **Production Ready**: Optimized for deployment and performance

---

## 🔗 **Integration Points**

### **Channel System Flow**
```
1. YAML Configuration (channels/*.yaml)
   ↓
2. API Endpoints (/api/channels)
   ↓  
3. PlayCanvas Scene Loading
   ↓
4. Real-Time WebSocket Sync
   ↓
5. Multi-User Collaboration
```

### **Props System Flow**
```
1. YAML Definitions (props/*/*.yaml)
   ↓
2. API Endpoints (/api/props, /api/sessions/{id}/props/{propId})
   ↓
3. PlayCanvas Entity Instantiation
   ↓
4. Component System Integration
   ↓
5. Physics & Rendering
```

### **Web Interface Flow**
```
1. Auto-Generated Clients (htdocs/static/js/)
   ↓
2. 85 REST API Endpoints
   ↓
3. PlayCanvas Game Engine
   ↓
4. Professional 3D Rendering
```

---

## 🛠️ **Development Workflow**

### **Channel Development**
```bash
# Edit channel configuration
vim /opt/hd1/share/channels/channel_one.yaml

# Restart HD1 to reload channels
cd /opt/hd1/src && make restart

# Test via API
curl http://localhost:8080/api/channels
```

### **Props Development**
```bash
# Add new prop definition
vim /opt/hd1/share/props/furniture/new-table.yaml

# List available props
curl http://localhost:8080/api/props

# Instantiate in session
curl -X POST http://localhost:8080/api/sessions/{sessionId}/props/new-table
```

### **Web Interface Testing**
```bash
# Start development server
cd /opt/hd1/src && make start

# Access PlayCanvas interface
open http://localhost:8080

# View real-time updates in browser console
```

---

## 📊 **V5.0.0 Architecture Benefits**

### **Performance Optimized**
- **<50ms API Response**: Professional game engine performance
- **<10ms WebSocket Latency**: Real-time collaboration
- **Minimal Assets**: Clean, purpose-built directory structure
- **Native PlayCanvas**: Optimized 3D rendering pipeline

### **Developer Experience**
- **YAML Configuration**: Human-readable scene definitions
- **Auto-Generated Clients**: Zero manual API synchronization
- **Single Source of Truth**: api.yaml drives all functionality
- **Clean Architecture**: No legacy systems or deprecated code

### **Production Ready**
- **85 REST Endpoints**: Complete game engine API surface
- **Entity-Component-System**: Professional game development patterns
- **Real-Time Collaboration**: Multi-user environment support
- **Comprehensive Documentation**: Complete system understanding

---

## 📚 **File Organization Standards**

### **Required Files Only**
- ✅ **channels/**: YAML scene configurations for PlayCanvas
- ✅ **props/**: YAML entity definitions for reusable components
- ✅ **htdocs/**: Web interface with auto-generated API clients

### **Removed Legacy**
- ❌ **environments/**: Legacy shell scripts (removed)
- ❌ **scenes/**: Legacy shell scripts (removed)
- ❌ **lighting/**: Obsolete directory (removed)
- ❌ **configs/**: Unused directory (removed)
- ❌ **templates/**: Unused directory (removed)
- ❌ **All .sh files**: Legacy A-Frame scripts (removed)

---

## 🎖️ **Quality Standards**

HD1 v5.0.0 Share directory represents **production-grade asset organization**:

✅ **100% Native PC Architecture**: Optimized for local development and deployment  
✅ **Single Source of Truth**: YAML + api.yaml drive all functionality  
✅ **PlayCanvas Integration**: Professional 3D game engine architecture  
✅ **Clean Dependencies**: No legacy systems or deprecated code  
✅ **API-First Design**: REST endpoints control all game engine features  

**"Professional game development through specification-driven architecture."**

---

*Last Updated: 2025-07-03*  
*HD1 Version: 5.0.0*  
*Architecture: API-First Game Engine Platform*