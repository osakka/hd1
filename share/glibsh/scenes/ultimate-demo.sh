#!/bin/bash

# =========================================================================
# THD Scene: Ultimate Demo - Complete Holodeck Showcase
# =========================================================================
#
# Professional holodeck demonstration with advanced materials, lighting,
# and complex 3D structures showcasing THD capabilities
#
# Usage: ./ultimate-demo.sh [SESSION_ID]
# =========================================================================

set -euo pipefail

# Scene configuration
SCENE_NAME="Ultimate Demo"
SCENE_DESCRIPTION="Complete holodeck showcase with metallic platforms, crystal formations, and cinematic lighting"

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

# Create ultimate demo scene
echo "Creating ultimate demo scene..."

# Sky environment
thd::sky holodeck_sky color "dark blue"

# Central metallic platform
thd::cylinder central_platform at 0,0.2,0 size 4 color silver metallic

# Crystal formations - transparent cones
thd::cone crystal_1 at -3,2,-3 size 1.5 color magenta transparent
thd::cone crystal_2 at 3,1.8,3 size 1.2 color cyan transparent

# Metallic structures
thd::cylinder metallic_pillar_1 at -5,2,0 size 0.5 color orange metallic
thd::sphere energy_sphere at 5,2,0 size 0.8 color yellow transparent

# Professional lighting setup
thd::light main_light at 5,8,5 intensity 1.2 color white
thd::light accent_light at -3,6,-3 intensity 0.8 color cyan

# Status display panel
thd::plane status_display at 0,6,-5 size 3 color white transparent

echo "THD Scene '$SCENE_NAME' loaded successfully" 
echo "Objects created: 9"
echo "Session: $SESSION_ID"