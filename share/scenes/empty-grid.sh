#!/bin/bash

# =========================================================================
# THD Scene: Empty Grid - Clean Holodeck Environment
# =========================================================================
#
# Professional minimal scene with just the coordinate grid system
# Perfect for custom object creation and experimentation
#
# Usage: ./empty-grid.sh [SESSION_ID]
# =========================================================================

set -euo pipefail

# Scene configuration
SCENE_NAME="Empty Grid"
SCENE_DESCRIPTION="Clean holodeck with just the coordinate grid system"

# Get session ID from argument or use active session
SESSION_ID="${1:-${THD_SESSION:-}}"

if [[ -z "$SESSION_ID" ]]; then
    echo "Error: Session ID required" >&2
    exit 1
fi

# Set THD_ROOT and source functions
THD_ROOT="/opt/holo-deck"
source "${THD_ROOT}/lib/thdlib.sh" 2>/dev/null || {
    echo "ERROR: THD functions not available"
    exit 1
}

# Set session ID for THD functions
THD_SESSION_ID="$SESSION_ID"

echo "Creating $SCENE_NAME scene..."

# Empty grid scene - clear any existing objects
thd::clear

# The grid system is automatically present in every holodeck session

echo "THD Scene '$SCENE_NAME' loaded successfully"
echo "Objects created: 0"
echo "Session: $SESSION_ID"