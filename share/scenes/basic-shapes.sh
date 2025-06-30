#!/bin/bash

# =========================================================================
# HD1 Scene: Basic Shapes - Fundamental 3D Demonstration  
# =========================================================================
#
# Educational scene showcasing fundamental 3D shapes with various materials
# Perfect for learning THD object creation and material properties
#
# Usage: ./basic-shapes.sh [SESSION_ID]
# =========================================================================

set -euo pipefail

# Scene configuration
SCENE_NAME="Basic Shapes"
SCENE_DESCRIPTION="Fundamental 3D shapes demonstration with various materials and colors"

# Get session ID from argument or use active session
SESSION_ID="${1:-${HD1_SESSION:-}}"

if [[ -z "$SESSION_ID" ]]; then
    echo "Error: Session ID required" >&2
    exit 1
fi

# Path to auto-generated THD client
THD_CLIENT="/opt/holo-deck/build/bin/thd-client"

echo "Creating basic shapes scene..."

# Basic cube - red
$THD_CLIENT create-object "$SESSION_ID" '{
    "name": "demo_cube",
    "type": "cube", 
    "x": -3, "y": 1, "z": 0,
    "scale": 1,
    "color": {"r": 1.0, "g": 0.2, "b": 0.2, "a": 1.0}
}' > /dev/null

# Basic sphere - green
$THD_CLIENT create-object "$SESSION_ID" '{
    "name": "demo_sphere",
    "type": "sphere",
    "x": 0, "y": 1, "z": 0, 
    "scale": 1,
    "color": {"r": 0.2, "g": 1.0, "b": 0.2, "a": 1.0}
}' > /dev/null

# Basic cylinder - blue
$THD_CLIENT create-object "$SESSION_ID" '{
    "name": "demo_cylinder", 
    "type": "cylinder",
    "x": 3, "y": 1, "z": 0,
    "scale": 1,
    "color": {"r": 0.2, "g": 0.2, "b": 1.0, "a": 1.0}
}' > /dev/null

# Basic cone - yellow
$THD_CLIENT create-object "$SESSION_ID" '{
    "name": "demo_cone",
    "type": "cone",
    "x": -1.5, "y": 1, "z": -3,
    "scale": 1,
    "color": {"r": 1.0, "g": 1.0, "b": 0.2, "a": 1.0}
}' > /dev/null

# Wireframe cube - white wireframe
$THD_CLIENT create-object "$SESSION_ID" '{
    "name": "wireframe_cube",
    "type": "cube",
    "x": 1.5, "y": 1, "z": -3,
    "scale": 1,
    "color": {"r": 0.8, "g": 0.8, "b": 0.8, "a": 1.0},
    "wireframe": true
}' > /dev/null

# Info panel
$THD_CLIENT create-object "$SESSION_ID" '{
    "name": "shapes_label",
    "type": "plane", 
    "x": 0, "y": 3, "z": -1,
    "scale": 2,
    "color": {"r": 1.0, "g": 1.0, "b": 1.0, "a": 0.3}
}' > /dev/null

echo "HD1 Scene '$SCENE_NAME' loaded successfully"
echo "Objects created: 6"
echo "Session: $SESSION_ID"