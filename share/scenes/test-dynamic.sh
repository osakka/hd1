#!/bin/bash
# Test scene to verify dynamic dropdown updates

# Set THD_ROOT and source functions
THD_ROOT="/opt/holo-deck"
source "${THD_ROOT}/lib/thdlib.sh" 2>/dev/null || {
    echo "ERROR: THD functions not available"
    exit 1
}

# Set session ID for THD functions
THD_SESSION_ID="$SESSION_ID"

# Create a simple test cube
thd::create_object "test_cube" "cube" 0 1 0

echo "Dynamic dropdown test scene loaded"