#!/bin/bash

# =========================================================================
# HD1 Scene: Ultimate Demo - Complete Holodeck Showcase
# =========================================================================
#
# Standard holodeck demonstration with advanced materials, lighting,
# and complex 3D structures showcasing HD1 capabilities
#
# Usage: ./ultimate-demo.sh [SESSION_ID]
# =========================================================================

set -euo pipefail

# Scene configuration
SCENE_NAME="Ultimate Demo"
SCENE_DESCRIPTION="Complete holodeck showcase with metallic platforms, crystal formations, and cinematic lighting"

# Get session ID from argument or use active session
SESSION_ID="${1:-${HD1_SESSION:-}}"

if [[ -z "$SESSION_ID" ]]; then
    echo "Error: Session ID required" >&2
    exit 1
fi

# Set THD_ROOT and source functions
HD1_ROOT="/opt/holo-deck"
source "${HD1_ROOT}/lib/hd1lib.sh" 2>/dev/null || {
    echo "ERROR: HD1 functions not available"
    exit 1
}

# Set session ID for THD functions
HD1_SESSION_ID="$SESSION_ID"

echo "Creating $SCENE_NAME scene..."

# Central metallic platform
thd::create_object "central_platform" "cylinder" 0 0.2 0

# Crystal formations
thd::create_object "crystal_1" "box" -3 2 -3
thd::create_object "crystal_2" "box" 3 1.8 3

# Metallic structures
thd::create_object "metallic_pillar_1" "cylinder" -5 2 0
thd::create_object "energy_sphere" "sphere" 5 2 0

# Standard lighting setup
thd::create_object "main_light" "light" 5 8 5
thd::create_object "accent_light" "light" -3 6 -3

# Status display panel
thd::create_object "status_display" "plane" 0 6 -5

echo "HD1 Scene '$SCENE_NAME' loaded successfully"
echo "Objects created: 7"
echo "Session: $SESSION_ID"