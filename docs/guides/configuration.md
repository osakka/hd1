# HD1 Configuration Guide

Complete configuration reference for HD1 (Holodeck One).

## Configuration Priority

HD1 uses a hierarchical configuration system with the following priority order:

1. **Command-line flags** (highest priority)
2. **Environment variables** 
3. **.env file**
4. **Default values** (lowest priority)

## Environment Variables

All HD1 environment variables use the `HD1_` prefix.

### Server Configuration
```bash
# Network binding
HD1_HOST=0.0.0.0                          # Server bind host (default: 0.0.0.0)
HD1_PORT=8080                             # Server bind port (default: 8080)

# External URLs
HD1_API_BASE=http://0.0.0.0:8080/api     # External API base URL
HD1_INTERNAL_API_BASE=http://localhost:8080/api  # Internal API communications

# Process management  
HD1_VERSION=v6.0.0                       # HD1 version identifier
HD1_DAEMON=true                          # Run in daemon mode (default: false)
HD1_PID_FILE=/opt/hd1/build/hd1.pid     # PID file location
```

### Directory Paths
```bash
# Core directories
HD1_ROOT_DIR=/opt/hd1                     # HD1 root directory
HD1_STATIC_DIR=/opt/hd1/share/htdocs/static  # Static files directory
HD1_LOG_DIR=/opt/hd1/build/logs          # Log directory

# Configuration directories
HD1_WORLDS_DIR=/opt/hd1/share/worlds     # Worlds configuration
HD1_AVATARS_DIR=/opt/hd1/share/avatars   # Avatars configuration
HD1_RECORDINGS_DIR=/opt/hd1/share/recordings  # Recording storage

# Build directories
HD1_BUILD_DIR=/opt/hd1/build             # Build artifacts
HD1_BIN_DIR=/opt/hd1/build/bin           # Binary output
HD1_RUNTIME_DIR=/opt/hd1/build/runtime   # Runtime data
```

### Logging Configuration
```bash
# Log levels: TRACE, DEBUG, INFO, WARN, ERROR, FATAL
HD1_LOG_LEVEL=INFO                        # Base logging level
HD1_TRACE_MODULES=websocket,entities     # Comma-separated trace modules
HD1_LOG_FILE=/opt/hd1/build/logs/hd1.log # Log file path (optional)
```

### WebSocket Configuration
```bash
# Timeouts and limits
HD1_WEBSOCKET_WRITE_TIMEOUT=10s          # Write timeout
HD1_WEBSOCKET_PONG_TIMEOUT=60s           # Pong wait timeout
HD1_WEBSOCKET_PING_PERIOD=54s            # Ping interval
HD1_WEBSOCKET_MAX_MESSAGE_SIZE=512       # Maximum message size (bytes)

# Buffer sizes
HD1_WEBSOCKET_READ_BUFFER_SIZE=4096      # Read buffer size
HD1_WEBSOCKET_WRITE_BUFFER_SIZE=4096     # Write buffer size
HD1_WEBSOCKET_CLIENT_BUFFER=256          # Client send buffer size
```

### World System Configuration
```bash
# World management
HD1_WORLDS_DEFAULT_WORLD=world_one       # Default world name
HD1_WORLDS_PROTECTED_LIST=world_one,world_two  # Protected worlds (comma-separated)
```

## Command-Line Flags

All environment variables have corresponding command-line flags:

```bash
# Server configuration
./hd1 --host=127.0.0.1                   # Bind to specific host
./hd1 --port=8081                        # Use different port
./hd1 --daemon                           # Run as daemon
./hd1 --pid-file=/custom/path/hd1.pid    # Custom PID file location

# Logging configuration
./hd1 --log-level=DEBUG                  # Set log level
./hd1 --log-file=/var/log/hd1.log       # Enable file logging
./hd1 --trace-modules=websocket,api      # Enable module tracing

# Directory configuration
./hd1 --root-dir=/custom/hd1             # Custom root directory
./hd1 --static-dir=/custom/static        # Custom static files directory

# Advanced options
./hd1 --internal-api-base=http://internal:8080/api  # Internal API URL
./hd1 --protected-worlds=secure,admin    # Specify protected worlds
./hd1 --version=v1.0.0                  # Override version string
```

## .env File Configuration

Create a `.env` file in the project root for local development:

```bash
# .env file example
HD1_HOST=127.0.0.1
HD1_PORT=9090
HD1_LOG_LEVEL=DEBUG
HD1_DAEMON=false
HD1_TRACE_MODULES=websocket,entities

# Development paths
HD1_ROOT_DIR=/home/developer/hd1
HD1_LOG_DIR=/tmp/hd1-logs

# Development WebSocket settings
HD1_WEBSOCKET_WRITE_TIMEOUT=30s
HD1_WEBSOCKET_PONG_TIMEOUT=120s
```

## Configuration Examples

### Development Configuration
```bash
# Local development with debug logging
export HD1_HOST=127.0.0.1
export HD1_PORT=8080
export HD1_LOG_LEVEL=DEBUG
export HD1_TRACE_MODULES=websocket,entities,api
export HD1_DAEMON=false

# Start development server
make start
```

### Production Configuration
```bash
# Production deployment
export HD1_HOST=0.0.0.0
export HD1_PORT=8080
export HD1_LOG_LEVEL=INFO
export HD1_DAEMON=true
export HD1_PID_FILE=/var/run/hd1.pid
export HD1_LOG_FILE=/var/log/hd1.log

# Production paths
export HD1_ROOT_DIR=/opt/hd1
export HD1_STATIC_DIR=/opt/hd1/share/htdocs/static
export HD1_LOG_DIR=/var/log

# Start production server
./hd1 --daemon --log-file=/var/log/hd1.log
```

### Docker Configuration
```bash
# Docker environment
HD1_HOST=0.0.0.0
HD1_PORT=8080
HD1_LOG_LEVEL=INFO
HD1_ROOT_DIR=/app
HD1_STATIC_DIR=/app/share/htdocs/static
HD1_LOG_DIR=/app/logs
```

### Load Balancer Configuration
```bash
# Behind load balancer
export HD1_HOST=127.0.0.1                # Bind to localhost only
export HD1_PORT=8081                     # Non-standard port
export HD1_API_BASE=https://api.domain.com/api  # External API URL
export HD1_INTERNAL_API_BASE=http://127.0.0.1:8081/api  # Internal URL
```

## Configuration Validation

### Checking Current Configuration
```go
// View configuration at startup
go run main.go --help

// Check specific values
curl http://localhost:8080/api/system/version | jq .
```

### Environment Variable Validation
```bash
# Check all HD1 environment variables
env | grep HD1_

# Validate specific settings
echo "Host: $HD1_HOST"
echo "Port: $HD1_PORT"
echo "Log Level: $HD1_LOG_LEVEL"
```

### Configuration Testing
```bash
# Test configuration changes
HD1_LOG_LEVEL=DEBUG HD1_PORT=9090 make start

# Verify configuration is applied
curl http://localhost:9090/api/system/version
```

## Advanced Configuration

### Custom Directory Layout
```bash
# Custom directory structure
export HD1_ROOT_DIR=/custom/hd1
export HD1_STATIC_DIR=/custom/assets
export HD1_LOG_DIR=/custom/logs
export HD1_BUILD_DIR=/custom/build
export HD1_RUNTIME_DIR=/custom/runtime

# Ensure directories exist
mkdir -p $HD1_STATIC_DIR $HD1_LOG_DIR $HD1_BUILD_DIR $HD1_RUNTIME_DIR
```

### Multi-Instance Configuration
```bash
# Instance 1
export HD1_HOST=127.0.0.1
export HD1_PORT=8080
export HD1_PID_FILE=/var/run/hd1-1.pid
export HD1_LOG_FILE=/var/log/hd1-1.log

# Instance 2  
export HD1_HOST=127.0.0.1
export HD1_PORT=8081
export HD1_PID_FILE=/var/run/hd1-2.pid
export HD1_LOG_FILE=/var/log/hd1-2.log
```

### WebSocket Tuning
```bash
# High-performance WebSocket settings
export HD1_WEBSOCKET_READ_BUFFER_SIZE=8192
export HD1_WEBSOCKET_WRITE_BUFFER_SIZE=8192
export HD1_WEBSOCKET_CLIENT_BUFFER=1024
export HD1_WEBSOCKET_WRITE_TIMEOUT=5s
export HD1_WEBSOCKET_PING_PERIOD=30s
```

## Configuration Management

### Environment-Specific Configs
```bash
# configs/development.env
HD1_LOG_LEVEL=DEBUG
HD1_TRACE_MODULES=websocket,entities
HD1_DAEMON=false

# configs/staging.env  
HD1_LOG_LEVEL=INFO
HD1_DAEMON=true
HD1_LOG_FILE=/var/log/hd1-staging.log

# configs/production.env
HD1_LOG_LEVEL=WARN
HD1_DAEMON=true
HD1_LOG_FILE=/var/log/hd1.log

# Load environment-specific config
source configs/production.env && ./hd1
```

### Configuration Validation Script
```bash
#!/bin/bash
# validate-config.sh

# Required variables
required_vars=(
    "HD1_HOST"
    "HD1_PORT" 
    "HD1_ROOT_DIR"
    "HD1_STATIC_DIR"
)

# Check required variables
for var in "${required_vars[@]}"; do
    if [[ -z "${!var}" ]]; then
        echo "ERROR: $var is not set"
        exit 1
    fi
done

# Validate directories exist
if [[ ! -d "$HD1_ROOT_DIR" ]]; then
    echo "ERROR: HD1_ROOT_DIR does not exist: $HD1_ROOT_DIR"
    exit 1
fi

if [[ ! -d "$HD1_STATIC_DIR" ]]; then
    echo "ERROR: HD1_STATIC_DIR does not exist: $HD1_STATIC_DIR"
    exit 1
fi

echo "Configuration validation passed"
```

## Troubleshooting Configuration

### Common Configuration Issues

#### Port Already in Use
```bash
# Check what's using the port
lsof -i :8080

# Use different port
HD1_PORT=8081 make start
```

#### Permission Denied
```bash
# Check directory permissions
ls -la $HD1_ROOT_DIR
ls -la $HD1_LOG_DIR

# Fix permissions
sudo chown -R $USER:$USER $HD1_ROOT_DIR
chmod -R 755 $HD1_ROOT_DIR
```

#### Invalid Configuration Values
```bash
# Check log for configuration errors
tail -f $HD1_LOG_DIR/hd1.log

# Validate log level
echo $HD1_LOG_LEVEL  # Should be: TRACE, DEBUG, INFO, WARN, ERROR, FATAL

# Reset to defaults
unset HD1_LOG_LEVEL
make start
```

### Configuration Debugging
```bash
# Enable configuration debugging
HD1_LOG_LEVEL=DEBUG make start 2>&1 | grep -i config

# Test minimal configuration
env -i HD1_HOST=127.0.0.1 HD1_PORT=8080 ./hd1

# Validate specific components
curl -f http://localhost:8080/api/system/version || echo "API not accessible"
```

## Security Considerations

### Production Security
```bash
# Restrict network binding
HD1_HOST=127.0.0.1  # Only localhost access

# Secure file permissions
chmod 600 .env       # Restrict .env file access
chmod 700 $HD1_LOG_DIR  # Restrict log directory

# Use non-root user
export HD1_USER=hd1
export HD1_GROUP=hd1
```

### Sensitive Configuration
```bash
# Avoid logging sensitive values
HD1_LOG_LEVEL=INFO  # Don't use DEBUG in production

# Use secure file locations
HD1_PID_FILE=/var/run/hd1/hd1.pid  # Secure PID file location
HD1_LOG_FILE=/var/log/hd1/hd1.log  # Secure log file location
```

---

*Configuration Guide for HD1 v6.0.0 - Three.js Game Engine Platform*