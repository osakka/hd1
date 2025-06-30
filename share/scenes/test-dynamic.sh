#!/bin/bash
# Test scene to verify dynamic dropdown updates

# Set HD1_ROOT and source functions
HD1_ROOT="/opt/hd1"
source "${HD1_ROOT}/lib/hd1lib.sh" 2>/dev/null || {
    echo "ERROR: HD1 functions not available"
    exit 1
}

# Set session ID for HD1 functions
HD1_SESSION_ID="$SESSION_ID"

# Create a simple test cube
hd1::create_object "test_cube" "cube" 0 1 0

echo "Dynamic dropdown test scene loaded"