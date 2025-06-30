#!/bin/bash

SCENE_NAME="Test WebSocket Updates"

# Set THD_ROOT and source functions
HD1_ROOT="/opt/holo-deck"
source "${HD1_ROOT}/lib/thdlib.sh" 2>/dev/null || {
    echo "ERROR: HD1 functions not available"
    exit 1
}

# Simple test scene to verify WebSocket scene update system
thd::create_object "test_websocket_cube" "cube" 0 2 -5
echo "Test WebSocket scene created successfully"