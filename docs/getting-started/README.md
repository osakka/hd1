# Getting Started with HD1

**HD1 v5.0.1 - API-First Game Engine Platform**

This guide will help you quickly get started with HD1, from installation to creating your first entities and scenes.

## 🚀 **Quick Start**

### **Prerequisites**
- Go 1.21+ (for building from source)
- Modern web browser (Chrome, Firefox, Safari, Edge)
- Basic understanding of REST APIs and 3D concepts

### **Installation**

#### **Option 1: Build from Source**
```bash
# Clone the repository
git clone https://git.uk.home.arpa/itdlabs/holo-deck.git
cd holo-deck

# Build HD1
cd src
make clean && make

# Start the server
make start
```

#### **Option 2: Docker (Coming Soon)**
```bash
# Docker support planned for v5.1.0
docker run -p 8080:8080 hd1/holodeck-one:latest
```

### **Verify Installation**
1. **Check server status**:
   ```bash
   curl http://localhost:8080/api/system/version
   ```

2. **Open web interface**:
   Navigate to http://localhost:8080 in your browser

3. **Test API connectivity**:
   ```bash
   # Create a session
   curl -X POST http://localhost:8080/api/sessions
   ```

## 🎯 **First Steps**

### **1. Understanding HD1 Architecture**

HD1 follows an **API-first** architecture where everything is controlled via REST endpoints:

```
HTTP APIs → Game Commands → Server State → WebSocket Events → PlayCanvas Rendering
```

**Key Concepts:**
- **Sessions**: Isolated game instances
- **Channels**: Scene configurations and collaborative spaces
- **Entities**: Game objects with components
- **Components**: Behavior and properties (transform, model, physics, etc.)

### **2. Create Your First Session**

```bash
# Create a new session
curl -X POST http://localhost:8080/api/sessions \
  -H "Content-Type: application/json" \
  -d '{}'

# Response: {"session_id": "session-abc123", "status": "active"}
```

### **3. Join a Channel**

Channels define the 3D scene configuration:

```bash
# Join channel_one (contains a red box on a floor)
curl -X POST http://localhost:8080/api/sessions/session-abc123/channel/join \
  -H "Content-Type: application/json" \
  -d '{"channel_id": "channel_one"}'
```

### **4. Create Your First Entity**

```bash
# Create a blue cube entity
curl -X POST http://localhost:8080/api/sessions/session-abc123/entities \
  -H "Content-Type: application/json" \
  -d '{
    "name": "blue_cube",
    "components": {
      "transform": {
        "position": [2, 1, 0],
        "rotation": [0, 45, 0],
        "scale": [1, 1, 1]
      },
      "model": {
        "type": "box"
      },
      "material": {
        "diffuse": "#0000ff",
        "metalness": 0.0,
        "roughness": 0.7
      }
    }
  }'
```

### **5. View Your Scene**

Open http://localhost:8080 in your browser and you'll see:
- A gray floor
- A red box (from channel configuration)
- Your new blue cube at position [2, 1, 0]

## 📚 **Next Steps**

### **Learn Core Concepts**
- **[User Guide](../user-guide/README.md)** - Complete user documentation
- **[API Reference](../reference/api-specification.md)** - All 59 endpoints documented
- **[Architecture Overview](../architecture/overview.md)** - System architecture

### **Build Something**
- **[Examples](examples/)** - Code samples and tutorials
- **[Developer Guide](../developer-guide/README.md)** - Development documentation
- **[API Development](../developer-guide/api-development.md)** - Building with HD1 APIs

### **Deploy to Production**
- **[Operations Guide](../operations/README.md)** - Production deployment
- **[Security](../operations/security.md)** - Security configuration
- **[Monitoring](../operations/monitoring.md)** - System monitoring

## 🛠️ **Development Workflow**

### **Basic Development Commands**
```bash
# Build and start
cd src && make clean && make && make start

# Stop server
make stop

# Regenerate auto-generated code
make generate

# View logs
tail -f build/logs/hd1.log
```

### **Testing Your Changes**
```bash
# Test with CLI client
./build/bin/hd1-client --help

# Test API endpoints
curl http://localhost:8080/api/sessions

# View real-time logs
curl http://localhost:8080/api/admin/logging/level \
  -X POST -d '{"level": "DEBUG"}'
```

## 🔍 **Troubleshooting**

### **Common Issues**

**Port already in use:**
```bash
# Find process using port 8080
lsof -i :8080

# Kill existing process
pkill hd1
```

**Build failures:**
```bash
# Clean and rebuild
cd src && make clean && make
```

**API not responding:**
```bash
# Check server status
curl http://localhost:8080/api/system/version

# Check logs
tail -f build/logs/hd1.log
```

### **Getting Help**
- **[User Troubleshooting](../user-guide/troubleshooting.md)** - User issues and solutions
- **[Operations Troubleshooting](../operations/troubleshooting.md)** - Production issues
- **[Issue Tracker](https://github.com/hd1/issues)** - Report bugs and request features

## 📖 **Additional Resources**

- **[Complete API Documentation](../reference/api-specification.md)** - All 59 endpoints
- **[Channel Configuration](../user-guide/channels.md)** - YAML-based scene management
- **[Entity-Component System](../user-guide/entities-components.md)** - Game object architecture
- **[Architectural Decisions](../decisions/README.md)** - Design rationale and history

---

**Next**: [User Guide](../user-guide/README.md) | **Back to**: [Documentation Home](../README.md)