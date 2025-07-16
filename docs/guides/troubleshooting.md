# HD1 Troubleshooting Guide

Common issues and solutions for HD1 (Holodeck One).

## Build Issues

### Make Generate Fails
```bash
# Error: failed to parse api.yaml
# Solution: Validate OpenAPI specification syntax
make validate

# Error: missing handler files
# Solution: Create missing handler files or disable strict validation
vim src/api.yaml  # Set fail-on-missing-handlers: false
```

### Go Build Errors
```bash
# Error: package not found
go mod tidy
go mod download

# Error: undefined references after cleanup
make clean && make generate && make build

# Error: import cycle
# Solution: Review import dependencies in affected files
```

### Code Generation Issues
```bash
# Error: template execution failed
# Solution: Check template syntax in src/codegen/templates/
ls -la src/codegen/templates/

# Error: permission denied writing generated files
chmod 755 src/
```

## Server Issues

### Server Won't Start
```bash
# Check if port is already in use
lsof -i :8080
# Solution: Stop existing process or use different port
HD1_PORT=8081 make start

# Check configuration
env | grep HD1_
# Verify all required directories exist
```

### Server Crashes on Startup
```bash
# Check logs for detailed error
make logs

# Common causes:
# 1. Missing static directory
ls -la /opt/hd1/share/htdocs/static/
# 2. Permission issues
chmod -R 755 /opt/hd1/build/

# 3. Invalid configuration
HD1_LOG_LEVEL=DEBUG make start
```

### API Endpoints Not Responding
```bash
# Test basic connectivity
curl -v http://localhost:8080/api/system/version

# Check if router is generated correctly
grep -n "HandleFunc" src/auto_router.go

# Verify handler exists
ls -la src/api/system/version.go
```

## WebSocket Issues

### Connection Failures
```javascript
// Check browser console for errors
console.log('WebSocket state:', ws.readyState);

// Common states:
// 0 = CONNECTING
// 1 = OPEN  
// 2 = CLOSING
// 3 = CLOSED
```

```bash
# Test WebSocket connectivity
wscat -c ws://localhost:8080/ws

# Check server logs for WebSocket errors
grep -i websocket /opt/hd1/build/logs/hd1.log
```

### Frequent Disconnections
```bash
# Check WebSocket timeout settings
env | grep HD1_WEBSOCKET_

# Increase timeout values for unstable connections
export HD1_WEBSOCKET_PONG_TIMEOUT=120s
export HD1_WEBSOCKET_PING_PERIOD=60s
make start
```

### Messages Not Broadcasting
```bash
# Check session association
# Send session_associate message first:
{"type": "session_associate", "session_id": "test_session"}

# Verify hub registration in server logs
grep "client joined session" /opt/hd1/build/logs/hd1.log
```

## Console Issues

### Console Not Loading
```bash
# Check static file serving
curl -v http://localhost:8080/static/js/hd1-console.js

# Verify static directory exists and has files
ls -la /opt/hd1/share/htdocs/static/js/

# Check for JavaScript errors in browser console
# Open DevTools (F12) and check Console tab
```

### Rebootstrap Not Working
```javascript
// Manual rebootstrap trigger
localStorage.clear();
sessionStorage.clear();
window.location.reload(true);

// Check if cookies are being cleared properly
document.cookie.split(";").forEach(cookie => {
    console.log('Cookie:', cookie);
});
```

### Version Mismatch Issues
```bash
# Check server JS version
curl http://localhost:8080/api/system/version | jq .js_version

# Force cache refresh
# Hard refresh browser (Ctrl+Shift+R or Cmd+Shift+R)

# Clear browser cache completely
# Chrome: Settings > Privacy > Clear browsing data
```

## Performance Issues

### High Memory Usage
```bash
# Check memory usage
ps aux | grep hd1

# Enable memory profiling
go tool pprof http://localhost:8080/debug/pprof/heap

# Check for memory leaks in WebSocket connections
grep "unregister" /opt/hd1/build/logs/hd1.log
```

### Slow Response Times
```bash
# Check CPU usage
top -p $(pgrep hd1)

# Profile CPU usage
go tool pprof http://localhost:8080/debug/pprof/profile

# Check for blocked goroutines
go tool pprof http://localhost:8080/debug/pprof/goroutine
```

### WebSocket Message Backlog
```bash
# Check WebSocket buffer settings
env | grep HD1_WEBSOCKET_.*BUFFER

# Increase buffer sizes for high-throughput scenarios
export HD1_WEBSOCKET_CLIENT_BUFFER=1024
export HD1_WEBSOCKET_READ_BUFFER_SIZE=8192
export HD1_WEBSOCKET_WRITE_BUFFER_SIZE=8192
```

## Development Issues

### API Changes Not Reflected
```bash
# Ensure code generation after API changes
make generate

# Check if auto_router.go was updated
stat -c %Y src/auto_router.go

# Restart server to load new routing
make stop && make start
```

### Debug Logging Not Working
```bash
# Verify log level setting
echo $HD1_LOG_LEVEL

# Enable debug logging
HD1_LOG_LEVEL=DEBUG make start

# Enable module-specific tracing
HD1_TRACE_MODULES=websocket,entities make start
```

### Build Cache Issues
```bash
# Clean all build artifacts
make clean

# Clear Go module cache
go clean -modcache

# Regenerate everything from scratch
make clean && make generate && make build
```

## Network Issues

### CORS Errors
```javascript
// Check browser console for CORS errors
// HD1 automatically sets CORS headers:
// Access-Control-Allow-Origin: *
// Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
// Access-Control-Allow-Headers: Content-Type, X-Client-ID
```

```bash
# Test CORS with curl
curl -X OPTIONS -H "Access-Control-Request-Method: POST" \
     -H "Origin: http://example.com" \
     http://localhost:8080/api/threejs/entities
```

### Firewall/Proxy Issues
```bash
# Test direct connection
telnet localhost 8080

# Check if running behind proxy
env | grep -i proxy

# Test WebSocket through proxy
# Some proxies don't support WebSocket upgrades
```

### DNS Resolution Issues
```bash
# Test with IP address instead of hostname
curl http://127.0.0.1:8080/api/system/version

# Check hosts file
cat /etc/hosts | grep localhost
```

## File System Issues

### Permission Denied Errors
```bash
# Check directory permissions
ls -la /opt/hd1/
ls -la /opt/hd1/build/
ls -la /opt/hd1/share/

# Fix permissions
sudo chown -R $USER:$USER /opt/hd1/
chmod -R 755 /opt/hd1/
chmod -R 644 /opt/hd1/src/*.go
```

### Disk Space Issues
```bash
# Check available disk space
df -h /opt/hd1/

# Check log file sizes
du -sh /opt/hd1/build/logs/

# Clean old logs
find /opt/hd1/build/logs/ -name "*.log" -mtime +7 -delete
```

### Missing Files
```bash
# Verify all required files exist
ls -la /opt/hd1/src/api.yaml
ls -la /opt/hd1/share/htdocs/index.html
ls -la /opt/hd1/share/htdocs/static/js/hd1-console.js

# Restore from git if files are missing
git status
git checkout HEAD -- missing_file.ext
```

## Debugging Techniques

### Enable Verbose Logging
```bash
# Maximum logging verbosity
HD1_LOG_LEVEL=TRACE HD1_TRACE_MODULES=websocket,entities,api,server make start

# Monitor logs in real-time
tail -f /opt/hd1/build/logs/hd1.log

# Filter specific log types
grep "ERROR" /opt/hd1/build/logs/hd1.log
grep "websocket" /opt/hd1/build/logs/hd1.log
```

### Browser Developer Tools
```javascript
// Monitor WebSocket traffic
// 1. Open DevTools (F12)
// 2. Go to Network tab
// 3. Filter by "WS" (WebSocket)
// 4. Click on WebSocket connection to see messages

// Debug console JavaScript
// 1. Open Console tab
// 2. Look for errors or warnings
// 3. Check HD1Console object state
console.log(window.hd1Console);
```

### Server Profiling
```bash
# Start server with profiling enabled
HD1_LOG_LEVEL=DEBUG make start

# CPU profiling (30 seconds)
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30

# Memory profiling
go tool pprof http://localhost:8080/debug/pprof/heap

# Goroutine analysis
go tool pprof http://localhost:8080/debug/pprof/goroutine
```

## Recovery Procedures

### Complete System Reset
```bash
# 1. Stop all HD1 processes
make stop
pkill -f hd1

# 2. Clean all build artifacts
make clean
rm -rf /opt/hd1/build/*

# 3. Reset configuration to defaults
unset $(env | grep HD1_ | cut -d= -f1)

# 4. Rebuild from scratch
make clean && make && make start
```

### Emergency Recovery
```bash
# If system is completely broken:
# 1. Reset to known good state
git reset --hard HEAD
git clean -fd

# 2. Restore default configuration
cp .env.example .env  # If example exists

# 3. Rebuild everything
make clean && make generate && make build

# 4. Test basic functionality
curl http://localhost:8080/api/system/version
```

### Data Recovery
```bash
# If configuration or data is corrupted:
# 1. Backup current state
cp -r /opt/hd1/build /tmp/hd1-backup-$(date +%Y%m%d)

# 2. Restore from backup
# (Restore from your backup system)

# 3. Verify integrity
make validate
```

## Getting Help

### Log Analysis
```bash
# Collect diagnostic information
echo "=== System Information ===" > debug.log
uname -a >> debug.log
go version >> debug.log
echo "=== HD1 Configuration ===" >> debug.log
env | grep HD1_ >> debug.log
echo "=== Recent Logs ===" >> debug.log
tail -100 /opt/hd1/build/logs/hd1.log >> debug.log
```

### Creating Bug Reports
Include the following information:
1. **HD1 Version**: `curl http://localhost:8080/api/system/version`
2. **Environment**: OS, Go version, browser (if applicable)
3. **Configuration**: Relevant environment variables
4. **Steps to Reproduce**: Exact commands/actions that trigger the issue
5. **Expected vs Actual**: What should happen vs what actually happens
6. **Logs**: Relevant log entries with timestamps
7. **Error Messages**: Complete error messages and stack traces

### Performance Issues
For performance problems, include:
1. **System Resources**: CPU, memory, disk usage
2. **Load Characteristics**: Number of connections, message frequency
3. **Profiling Data**: CPU and memory profiles if available
4. **Timeline**: When the issue started and any recent changes

---

*Troubleshooting Guide for HD1 v0.7.0 - Three.js Game Engine Platform*