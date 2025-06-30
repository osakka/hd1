#!/bin/bash

# =========================================================================
# HD1 Scene: Empty Grid - Clean Holodeck Environment
# =========================================================================
#
# Standard minimal scene with just the coordinate grid system
# Perfect for custom object creation and experimentation
#
# Usage: ./empty-grid.sh [SESSION_ID]
# =========================================================================

set -euo pipefail

# Scene configuration
SCENE_NAME="Empty Grid"
SCENE_DESCRIPTION="Clean holodeck with just the coordinate grid system"

# Get session ID from argument or use active session
SESSION_ID="${1:-${HD1_SESSION:-}}"

if [[ -z "$SESSION_ID" ]]; then
    echo "Error: Session ID required" >&2
    exit 1
fi

# Set HD1_ROOT and source functions
HD1_ROOT="/opt/hd1"
source "${HD1_ROOT}/lib/hd1lib.sh" 2>/dev/null || {
    echo "ERROR: HD1 functions not available"
    exit 1
}

# Set session ID for HD1 functions
HD1_SESSION_ID="$SESSION_ID"

echo "Creating $SCENE_NAME scene..."

# Empty grid scene - clear any existing objects
hd1::clear

# The grid system is automatically present in every holodeck session

echo "HD1 Scene '$SCENE_NAME' loaded successfully"
echo "Objects created: 0"
echo "Session: $SESSION_ID"