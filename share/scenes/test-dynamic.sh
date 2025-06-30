#!/bin/bash
# Test scene to verify dynamic dropdown updates

# Set THD_ROOT and source functions
HD1_ROOT="/opt/holo-deck"
source "${HD1_ROOT}/lib/thdlib.sh" 2>/dev/null || {
    echo "ERROR: HD1 functions not available"
    exit 1
}

# Set session ID for THD functions
HD1_SESSION_ID="$SESSION_ID"

# Create a simple test cube
thd::create_object "test_cube" "cube" 0 1 0

echo "Dynamic dropdown test scene loaded"