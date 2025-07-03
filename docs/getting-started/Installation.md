# HD1 Installation Guide

**HD1 v5.0.1 Installation and Setup**

## 📋 System Requirements

### **Minimum Requirements**
- **OS**: Linux (Ubuntu 20.04+ recommended)
- **Architecture**: x86_64
- **Memory**: 2GB RAM minimum, 4GB recommended
- **Storage**: 1GB available disk space
- **Network**: Internet connection for dependencies

### **Development Requirements**
- **Go**: 1.19+ required for building from source
- **Make**: Build system requirement
- **Git**: Version control access

## 🚀 Quick Installation

### **Method 1: Binary Release (Recommended)**
```bash
# Download latest release
wget https://github.com/your-org/hd1/releases/latest/hd1-linux-amd64.tar.gz

# Extract and install
tar -xzf hd1-linux-amd64.tar.gz
sudo mv hd1 /usr/local/bin/
```

### **Method 2: Build from Source**
```bash
# Clone repository
git clone https://github.com/your-org/hd1.git
cd hd1

# Build HD1
cd src && make clean && make

# Start HD1 daemon
make start
```

## ⚙️ Configuration

### **Environment Variables (Recommended)**
```bash
# Create .env file
cat > .env << 'EOF'
HD1_HOST=localhost
HD1_PORT=8080
HD1_API_BASE=http://localhost:8080/api
HD1_STATIC_DIR=/opt/hd1/share/htdocs
HD1_DAEMON=false
HD1_LOG_LEVEL=INFO
EOF
```

### **Command Line Flags**
```bash
# Start with custom configuration
./hd1 --host=0.0.0.0 --port=9090 --log-level=DEBUG
```

## 🧪 Verification

### **Test Installation**
```bash
# Check HD1 is running
curl http://localhost:8080/api/version

# Expected response
{"version": "v5.0.1", "status": "ok"}
```

### **Create Test Session**
```bash
# Create a new session
curl -X POST http://localhost:8080/api/sessions \
  -H "Content-Type: application/json" \
  -d '{"name": "test-session"}'

# Expected response with session_id
{"session_id": "session-abc123", "name": "test-session"}
```

## 🔧 Common Issues

### **Port Already in Use**
```bash
# Find process using port 8080
sudo lsof -i :8080

# Kill the process
sudo kill -9 <PID>
```

### **Permission Issues**
```bash
# Make HD1 executable
chmod +x hd1

# For system-wide installation
sudo chown root:root /usr/local/bin/hd1
```

## 📚 Next Steps

- **[Quick Start Guide](Quick-Start.md)** - Create your first 3D scene
- **[Configuration Reference](../reference/Configuration.md)** - Complete configuration options
- **[API Usage Guide](../user-guides/API-Usage.md)** - Learn the API endpoints

---

**Back to**: [Getting Started](README.md) | **Next**: [Quick Start](Quick-Start.md)