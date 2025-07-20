# HD1 (Holodeck One) 🚀

**Turn any service into a 3D interface with simple HTTP calls**

HD1 is a **pure WebGL REST platform** that lets any application render rich 3D interfaces through simple HTTP API calls. Think "GraphQL for 3D Graphics" - a universal API that makes building 3D interfaces as easy as REST endpoints.

```bash
# Create a 3D sphere
curl -X POST http://localhost:8080/api/geometries/sphere \
  -H "Content-Type: application/json" \
  -d '{"radius": 2, "color": "#ff0000", "position": {"x": 0, "y": 1, "z": 0}}'

# Add some lighting  
curl -X POST http://localhost:8080/api/lights/directional \
  -d '{"color": "#ffffff", "intensity": 1, "position": {"x": 10, "y": 10, "z": 5}}'
```

## 🎯 What HD1 Does

**For Developers:**
- ✅ **No Three.js Knowledge Required** - Simple REST API calls create 3D objects
- ✅ **69 API Endpoints** - Complete coverage of geometries, materials, lighting, cameras, animations, textures
- ✅ **Real-Time Sync** - WebSocket updates, multiple users see changes instantly
- ✅ **Production Ready** - Multi-tenant, mobile controls, auto-generated client libraries

**For Applications:**
- 📧 **Email Apps** → Floating 3D mail objects in space
- 📅 **Calendar Apps** → Spatial time blocks and 3D scheduling  
- 🤖 **AI Services** → Interactive 3D avatars with visual interfaces
- 📊 **Analytics** → Data visualization in immersive 3D environments
- 🎮 **Any Service** → Rich 3D interfaces without WebGL complexity

## 🚀 Quick Start

```bash
# 1. Start HD1
cd src && make && make start

# 2. Open the console
open http://localhost:8080

# 3. Create your first 3D scene via API
curl -X POST http://localhost:8080/api/geometries/box \
  -H "Content-Type: application/json" \
  -d '{"width": 2, "height": 2, "depth": 2, "color": "#00ff00"}'

# 4. View it live at http://localhost:8080
```

## 🎨 What You Get Out of the Box

### **Complete Three.js API Coverage (69 Endpoints)**
- **10 Geometries**: Box, sphere, cylinder, cone, torus, plane, ring, circle, capsule, torusknot
- **4 Materials**: Basic, phong, standard PBR, physical advanced PBR  
- **5 Lighting**: Directional, point, spot, ambient, hemisphere
- **2 Cameras**: Perspective and orthographic with full control
- **2 Animations**: Keyframe animations and timeline control
- **2 Textures**: URL loading and procedural generation

### **Real-Time Collaboration**
- **WebSocket Sync**: Changes appear instantly for all connected users
- **Multi-User Support**: Unlimited concurrent sessions
- **Avatar System**: Real-time multiplayer with automatic cleanup
- **Mobile Controls**: Touch-optimized for phones/tablets

### **Production Features**
- **Auto-Generated Clients**: JavaScript library with all 69 API methods
- **Configuration Management**: Environment variables, flags, .env files
- **Logging System**: Structured JSON logging with runtime level control  
- **Health Monitoring**: API status and performance metrics

## 🏗️ Architecture

```
HTTP API → Sync Operations → WebSocket Events → Three.js Rendering
```

**Single Source of Truth**: Everything driven by `src/schemas/hd1-api.yaml`

- **Backend**: Go server with auto-generated routes
- **Frontend**: Pure Three.js r170 with zero abstraction layers  
- **Sync**: TCP-simple reliability with sequence numbers
- **Config**: Zero hardcoded values, fully configurable

## 📊 API Examples

### Create 3D Objects
```bash
# Spinning torus with PBR material
curl -X POST http://localhost:8080/api/geometries/torus \
  -d '{"radius": 3, "tube": 1, "radialSegments": 16}'

curl -X POST http://localhost:8080/api/materials/physical \
  -d '{"color": "#8844ff", "metalness": 0.8, "roughness": 0.2}'
```

### Scene Management  
```bash
# Set background color
curl -X PUT http://localhost:8080/api/scene \
  -d '{"background": "#87CEEB"}'

# Add fog for atmosphere
curl -X PUT http://localhost:8080/api/scene \
  -d '{"fog": {"color": "#ffffff", "near": 1, "far": 100}}'
```

### Camera Control
```bash
# Position the camera
curl -X POST http://localhost:8080/api/cameras/perspective \
  -d '{"fov": 75, "position": {"x": 10, "y": 5, "z": 10}, "lookAt": {"x": 0, "y": 0, "z": 0}}'
```

## 🎮 Controls

- **Desktop**: WASD movement, mouse look, ESC to exit pointer lock
- **Mobile**: Left side for movement, right side for camera look
- **Touch**: Fully optimized for mobile devices

## 📁 Project Structure

```
/opt/hd1/
├── src/                      # Go backend source
│   ├── schemas/hd1-api.yaml  # Single source of truth API spec
│   ├── api/*/handlers.go     # Generated HTTP handlers  
│   ├── router/auto_router.go # Generated routing
│   └── config/config.go      # Configuration management
├── share/htdocs/static/      # Frontend assets
│   ├── js/hd1lib.js         # Generated API client (69 methods)
│   ├── js/hd1-threejs.js    # Three.js scene manager
│   └── vendor/threejs/      # Three.js r170 library
└── build/                   # Build artifacts
    ├── bin/hd1             # Compiled server
    └── logs/               # Application logs
```

## 🔧 Configuration

HD1 uses a priority-based config system: **Flags > Environment Variables > .env File > Defaults**

```bash
# Environment variables (HD1_ prefix)
export HD1_HOST=0.0.0.0
export HD1_PORT=8080
export HD1_LOG_LEVEL=INFO

# Command line flags  
./hd1 --host=127.0.0.1 --port=9090 --log-level=DEBUG

# .env file support
echo "HD1_PORT=3000" > .env
```

## 🛠️ Development

```bash
# Build and start
cd src && make && make start

# Clean rebuild  
make clean && make build

# Stop daemon
make stop

# View logs
make logs

# Check status
make status
```

## 📚 API Documentation

- **[Complete API Reference](src/schemas/hd1-api.yaml)** - OpenAPI 3.0.3 specification
- **[Development Context](CLAUDE.md)** - Technical implementation details
- **JavaScript Client**: Auto-generated `hd1lib.js` with all 69 methods

## 🌟 Why HD1?

**Before HD1:**
```javascript
// Complex Three.js setup
const scene = new THREE.Scene();
const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
const renderer = new THREE.WebGLRenderer();
const geometry = new THREE.BoxGeometry(1, 1, 1);
const material = new THREE.MeshBasicMaterial({color: 0x00ff00});
const cube = new THREE.Mesh(geometry, material);
scene.add(cube);
// ... 50+ more lines of WebGL setup
```

**With HD1:**
```bash
# One HTTP call
curl -X POST http://localhost:8080/api/geometries/box -d '{"color": "#00ff00"}'
```

**Perfect for:**
- 🚀 **Rapid Prototyping** - 3D interfaces in minutes, not weeks
- 🌐 **Any Backend** - Python, Node.js, Java, PHP - if it can make HTTP calls, it can render 3D
- 📱 **Mobile First** - Touch controls built-in, works on any device
- 👥 **Team Development** - Designers work with HTTP APIs, not WebGL code
- 🔄 **Real-Time Apps** - Built-in WebSocket sync for collaborative experiences

## 📈 Status

- **Version**: v1.0.0 (Production Ready)
- **API Endpoints**: 69 (Complete Three.js coverage)
- **Platform**: Pure WebGL REST platform
- **Architecture**: Single source of truth, zero hardcoded values
- **Mobile**: Full touch control support
- **Real-Time**: WebSocket synchronization with sequence-based reliability

## 🚦 Getting Help

1. **API Issues**: Check `make logs` for detailed error information
2. **Configuration**: See environment variables section above  
3. **Development**: Review `CLAUDE.md` for technical context
4. **Performance**: Built-in metrics at `/api/sync/stats`

---

**HD1 v1.0.0**: Turn any service into a 3D interface with simple HTTP calls. The GraphQL for 3D Graphics. 🎯