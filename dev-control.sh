#!/bin/bash
# VWS (Virtual World Synthesizer) - Professional Development Control System
# Makes development predictable for LLM workflows with professional artifact management

set -e

PROJECT_ROOT="/home/claude-3/3dv"
SRC_DIR="$PROJECT_ROOT/src"
BUILD_DIR="$PROJECT_ROOT/build"
BIN_DIR="$BUILD_DIR/bin"
LOG_DIR="$BUILD_DIR/logs"
RUNTIME_DIR="$BUILD_DIR/runtime"
SERVER_BIN="$BIN_DIR/vws"
CLIENT_BIN="$BIN_DIR/vws-client"
PID_FILE="$RUNTIME_DIR/vws.pid"
LOG_FILE="$LOG_DIR/development.log"

# Clear previous state
clear_state() {
    echo "🧹 CLEARING VWS DEVELOPMENT STATE"
    pkill -f "vws" 2>/dev/null || true
    rm -f "$PID_FILE" "$LOG_FILE"
    mkdir -p "$BIN_DIR" "$LOG_DIR" "$RUNTIME_DIR"
    echo "✅ State cleared"
}

# Build from specification
build_server() {
    echo "🏗️  BUILDING VWS FROM SPECIFICATION"
    cd "$SRC_DIR"
    
    # Validate spec exists
    if [ ! -f "api.yaml" ]; then
        echo "❌ FATAL: api.yaml missing!"
        exit 1
    fi
    
    # Professional build with proper logging
    mkdir -p "$LOG_DIR"
    make all > "$LOG_FILE" 2>&1
    
    if [ -f "$SERVER_BIN" ]; then
        echo "✅ VWS server built successfully -> $SERVER_BIN"
        return 0
    else
        echo "❌ Build failed - check $LOG_FILE"
        tail -10 "$LOG_FILE"
        return 1
    fi
}

# Start server with proper logging
start_server() {
    if [ -f "$PID_FILE" ]; then
        echo "⚠️  VWS server already running (PID: $(cat $PID_FILE))"
        return 0
    fi
    
    echo "🚀 STARTING VWS (Virtual World Synthesizer)"
    cd "$PROJECT_ROOT"
    
    # Start with structured logging to professional directories
    nohup "$SERVER_BIN" >> "$LOG_FILE" 2>&1 &
    SERVER_PID=$!
    echo $SERVER_PID > "$PID_FILE"
    
    # Wait for startup
    sleep 3
    
    if kill -0 $SERVER_PID 2>/dev/null; then
        echo "✅ VWS server started (PID: $SERVER_PID)"
        echo "📋 Server logs: tail -f $LOG_FILE"
        echo "🌐 VWS API: http://localhost:8080/api"
        return 0
    else
        echo "❌ VWS server failed to start"
        rm -f "$PID_FILE"
        echo "📋 Check logs: tail $LOG_FILE"
        return 1
    fi
}

# Stop server cleanly
stop_server() {
    if [ ! -f "$PID_FILE" ]; then
        echo "⚠️  No server running"
        return 0
    fi
    
    PID=$(cat "$PID_FILE")
    echo "🛑 STOPPING SERVER (PID: $PID)"
    
    kill $PID 2>/dev/null || true
    sleep 2
    
    # Force kill if needed
    kill -9 $PID 2>/dev/null || true
    rm -f "$PID_FILE"
    
    echo "✅ Server stopped"
}

# Test API endpoints with clear results
test_api() {
    echo "🧪 TESTING VWS API ENDPOINTS"
    
    if [ ! -f "$PID_FILE" ]; then
        echo "❌ VWS server not running - start first"
        return 1
    fi
    
    # Test core VWS endpoints
    echo "Testing GET /api/sessions..."
    if curl -s -f http://localhost:8080/api/sessions > /dev/null; then
        echo "✅ GET /api/sessions - OK"
    else
        echo "❌ GET /api/sessions - FAILED"
        return 1
    fi
    
    echo "Testing POST /api/sessions..."
    if curl -s -f -X POST http://localhost:8080/api/sessions > /dev/null; then
        echo "✅ POST /api/sessions - OK"
    else
        echo "❌ POST /api/sessions - FAILED"
        return 1
    fi
    
    # Test session workflow
    echo "Testing complete session workflow..."
    SESSION_ID=$(curl -s -X POST http://localhost:8080/api/sessions | jq -r '.session_id' 2>/dev/null || echo "")
    if [ -n "$SESSION_ID" ] && [ "$SESSION_ID" != "null" ]; then
        echo "✅ Session creation workflow - OK ($SESSION_ID)"
    else
        echo "❌ Session creation workflow - FAILED"
        return 1
    fi
    
    echo "✅ All VWS API tests passed"
}

# Show development status
status() {
    echo "📊 VWS DEVELOPMENT STATUS"
    echo "========================="
    
    # Check VWS server binary
    if [ -f "$SERVER_BIN" ]; then
        echo "✅ VWS Server: EXISTS ($SERVER_BIN)"
    else
        echo "❌ VWS Server: MISSING"
    fi
    
    # Check VWS client
    if [ -f "$CLIENT_BIN" ]; then
        echo "✅ VWS Client: EXISTS ($CLIENT_BIN)"
    else
        echo "❌ VWS Client: MISSING"
    fi
    
    # Check if running
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if kill -0 $PID 2>/dev/null; then
            echo "✅ VWS Status: RUNNING (PID: $PID)"
        else
            echo "❌ Server status: DEAD (stale PID file)"
            rm -f "$PID_FILE"
        fi
    else
        echo "⚠️  VWS Status: STOPPED"
    fi
    
    # Check port
    if netstat -tln 2>/dev/null | grep -q ":8080"; then
        echo "✅ Port 8080: LISTENING"
    else
        echo "❌ Port 8080: NOT LISTENING"
    fi
    
    # Check API specification
    if [ -f "$SRC_DIR/api.yaml" ]; then
        echo "✅ API Specification: EXISTS"
    else
        echo "❌ API Specification: MISSING"
    fi
    
    # Show build artifacts
    echo ""
    echo "📁 Build Artifacts:"
    echo "   Binaries: $BIN_DIR/"
    echo "   Logs: $LOG_DIR/"
    echo "   Runtime: $RUNTIME_DIR/"
    
    # Show recent logs
    if [ -f "$LOG_FILE" ]; then
        echo ""
        echo "📋 Recent development logs (last 5 lines):"
        tail -5 "$LOG_FILE"
    fi
}

# Full development cycle
dev_cycle() {
    echo "🔄 VWS FULL DEVELOPMENT CYCLE"
    clear_state
    build_server && start_server && test_api
}

# Create VWS client if not exists
create_client() {
    echo "📱 CREATING VWS CLIENT"
    cd "$SRC_DIR"
    make client
    echo "✅ VWS client available: $CLIENT_BIN"
}

# Help
help() {
    echo "VWS (Virtual World Synthesizer) - Development Control System"
    echo "=========================================================="
    echo ""
    echo "Core Commands:"
    echo "  ./dev-control.sh status     - Show VWS development status"
    echo "  ./dev-control.sh build      - Build VWS server from specification"
    echo "  ./dev-control.sh start      - Start VWS server"
    echo "  ./dev-control.sh stop       - Stop VWS server"
    echo "  ./dev-control.sh restart    - Restart VWS server"
    echo "  ./dev-control.sh test       - Test VWS API endpoints"
    echo ""
    echo "Workflow Commands:"
    echo "  ./dev-control.sh cycle      - Full build/start/test cycle"
    echo "  ./dev-control.sh clear      - Clear all development state"
    echo "  ./dev-control.sh client     - Create VWS API client"
    echo ""
    echo "Debugging Commands:"
    echo "  ./dev-control.sh logs       - Show recent development logs"
    echo "  ./dev-control.sh help       - Show this help"
    echo ""
    echo "Professional Artifact Locations:"
    echo "  Binaries: $BIN_DIR/"
    echo "  Logs: $LOG_DIR/"
    echo "  Runtime: $RUNTIME_DIR/"
    echo ""
    echo "VWS API: http://localhost:8080/api"
}

# Execute command
case "${1:-help}" in
    "status") status ;;
    "build") build_server ;;
    "start") start_server ;;
    "stop") stop_server ;;
    "restart") stop_server && start_server ;;
    "test") test_api ;;
    "cycle") dev_cycle ;;
    "clear") clear_state ;;
    "client") create_client ;;
    "logs") tail -20 "$LOG_FILE" 2>/dev/null || echo "📋 No development logs yet" ;;
    "help") help ;;
    *) help ;;
esac