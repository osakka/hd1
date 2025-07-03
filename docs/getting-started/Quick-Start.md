# HD1 Quick Start Guide

**Get up and running with HD1 v5.0.1 in 5 minutes**

## 🚀 Step 1: Start HD1

```bash
# Navigate to HD1 source directory
cd /opt/hd1/src

# Build and start HD1
make clean && make && make start
```

**Expected Output:**
```
HD1 v5.0.1 starting...
Server listening on http://localhost:8080
WebSocket hub initialized
Ready for connections
```

## 🌐 Step 2: Verify API Access

```bash
# Test API connectivity
curl http://localhost:8080/api/version

# Expected response
{"version": "v5.0.1", "status": "ok"}
```

## 🎮 Step 3: Create Your First Session

```bash
# Create a new 3D session
curl -X POST http://localhost:8080/api/sessions \
  -H "Content-Type: application/json" \
  -d '{"name": "my-first-scene"}'

# Save the session_id from response
export SESSION_ID="session-abc123"
```

## 📦 Step 4: Add 3D Objects

### **Create a Red Cube**
```bash
curl -X POST "http://localhost:8080/api/sessions/$SESSION_ID/entities" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "red-cube",
    "components": {
      "model": {"type": "box"},
      "transform": {"position": [0, 1, 0]},
      "material": {"diffuse": "#ff0000"}
    }
  }'
```

### **Create a Blue Sphere**
```bash
curl -X POST "http://localhost:8080/api/sessions/$SESSION_ID/entities" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "blue-sphere", 
    "components": {
      "model": {"type": "sphere"},
      "transform": {"position": [2, 1, 0]},
      "material": {"diffuse": "#0000ff"}
    }
  }'
```

## 🎯 Step 5: Position Camera

```bash
# Set camera to view your objects
curl -X PUT "http://localhost:8080/api/sessions/$SESSION_ID/camera/position" \
  -H "Content-Type: application/json" \
  -d '{
    "position": {"x": 5, "y": 3, "z": 5},
    "lookAt": {"x": 0, "y": 1, "z": 0}
  }'
```

## 🌍 Step 6: Join Channel and View Scene

### **Join a Channel**
```bash
# Join channel for real-time collaboration
curl -X POST "http://localhost:8080/api/sessions/$SESSION_ID/channel/join" \
  -H "Content-Type: application/json" \
  -d '{"channel_id": "channel_one"}'
```

### **Access Web Interface**
Open your browser and navigate to:
```
http://localhost:8080
```

You should see your red cube and blue sphere in a 3D environment!

## 🔄 Step 7: Real-Time Updates

### **WebSocket Connection (Optional)**
```javascript
// Connect to real-time updates
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Real-time update:', data);
};
```

### **Update Object Position**
```bash
# Move the red cube
curl -X PUT "http://localhost:8080/api/sessions/$SESSION_ID/entities/$ENTITY_ID/components/transform" \
  -H "Content-Type: application/json" \
  -d '{"position": [1, 2, 1]}'
```

## 🎉 Congratulations!

You've successfully:
- ✅ Started HD1 API-first game engine
- ✅ Created a 3D session
- ✅ Added 3D objects with materials
- ✅ Positioned the camera
- ✅ Joined a collaborative channel
- ✅ Viewed your scene in the browser

## 📚 What's Next?

### **Explore More Features**
- **[API Usage Guide](../user-guides/API-Usage.md)** - Complete API examples
- **[WebSocket Events](../user-guides/WebSocket-Events.md)** - Real-time synchronization
- **[Configuration Guide](../reference/Configuration.md)** - Advanced setup options

### **Development**
- **[Developer Guide](../developer-guide/README.md)** - Contribute to HD1
- **[Architecture Overview](../architecture/README.md)** - Understand the system
- **[Build System](../developer-guide/Build-System.md)** - Development workflow

## 🔧 Troubleshooting

### **HD1 Won't Start**
```bash
# Check if port 8080 is in use
sudo lsof -i :8080

# Use different port
make start HD1_PORT=9090
```

### **API Returns 404**
- Verify HD1 is running: `curl http://localhost:8080/api/version`
- Check logs: `tail -f /opt/hd1/logs/hd1.log`

### **No 3D Objects Visible**
- Verify session was created successfully
- Check entity creation response for errors
- Ensure camera is positioned correctly

---

**Back to**: [Getting Started](README.md) | **Next**: [API Usage Guide](../user-guides/API-Usage.md)