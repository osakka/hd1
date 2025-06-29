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

# Source THD shell functions
source "$(dirname "$0")/../client/thd-shell.sh"

# Use specified session
thd::session use "$SESSION_ID"

# Empty grid scene - no objects to create
# The grid system is automatically present in every holodeck session

echo "THD Scene '$SCENE_NAME' loaded successfully"
echo "Objects created: 0"
echo "Session: $SESSION_ID"