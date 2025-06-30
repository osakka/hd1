#!/bin/bash

# =========================================================================
# HD1 Scene: Anime UI Demo - Interactive Holographic Interface
# =========================================================================
#
# Interactive anime-style holodeck with floating UI elements, dynamic lighting,
# and data visualization components
#
# Usage: ./anime-ui.sh [SESSION_ID]
# =========================================================================

set -euo pipefail

# Scene configuration
SCENE_NAME="Anime UI Demo"
SCENE_DESCRIPTION="Interactive anime-style holodeck with floating UI elements, blue lighting, and data visualization cubes"

# Get session ID from argument or use active session
SESSION_ID="${1:-${HD1_SESSION:-}}"

if [[ -z "$SESSION_ID" ]]; then
    echo "Error: Session ID required" >&2
    exit 1
fi

# Path to auto-generated THD client
THD_CLIENT="/opt/holo-deck/build/bin/thd-client"

echo "Creating anime UI demo scene..."

# Central anime ring interface - wireframe cyan
$THD_CLIENT create-object "$SESSION_ID" '{
    "name": "central_ring",
    "type": "cylinder",
    "x": 0, "y": 2, "z": 0,
    "scale": 3,
    "color": {"r": 0.2, "g": 0.8, "b": 1.0, "a": 0.3},
    "wireframe": true
}' > /dev/null

# Floating UI cubes with anime colors
$THD_CLIENT create-object "$SESSION_ID" '{
    "name": "ui_cube_1",
    "type": "cube",
    "x": -4, "y": 3, "z": 2,
    "scale": 0.8,
    "color": {"r": 1.0, "g": 0.4, "b": 0.8, "a": 0.8}
}' > /dev/null

$THD_CLIENT create-object "$SESSION_ID" '{
    "name": "ui_cube_2", 
    "type": "cube",
    "x": 4, "y": 3.5, "z": -2,
    "scale": 0.6,
    "color": {"r": 0.4, "g": 1.0, "b": 0.8, "a": 0.8}
}' > /dev/null

# Data visualization spheres
$THD_CLIENT create-object "$SESSION_ID" '{
    "name": "data_sphere_1",
    "type": "sphere", 
    "x": -2, "y": 4, "z": -4,
    "scale": 0.5,
    "color": {"r": 0.8, "g": 0.2, "b": 1.0, "a": 0.9}
}' > /dev/null

$THD_CLIENT create-object "$SESSION_ID" '{
    "name": "data_sphere_2",
    "type": "sphere",
    "x": 3, "y": 2.5, "z": 4,
    "scale": 0.7,
    "color": {"r": 1.0, "g": 0.8, "b": 0.2, "a": 0.9}
}' > /dev/null

# Info panel for anime interface
$THD_CLIENT create-object "$SESSION_ID" '{
    "name": "info_panel",
    "type": "plane",
    "x": 0, "y": 5, "z": -3,
    "scale": 2,
    "color": {"r": 0.2, "g": 1.0, "b": 1.0, "a": 0.3}
}' > /dev/null

echo "HD1 Scene '$SCENE_NAME' loaded successfully"
echo "Objects created: 6"
echo "Session: $SESSION_ID"