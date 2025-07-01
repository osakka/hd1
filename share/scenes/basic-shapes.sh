#!/bin/bash

# =========================================================================
# HD1 Scene: Basic Shapes - Fundamental 3D Demonstration  
# =========================================================================
#
# Educational scene showcasing fundamental 3D shapes with various materials
# Perfect for learning HD1 object creation and material properties
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

# Path to auto-generated HD1 client
HD1_CLIENT="/opt/hd1/build/bin/hd1-client"

echo "Creating basic shapes scene..."

# Basic cube - red (API-first format)
$HD1_CLIENT create-object "$SESSION_ID" '{
    "name": "demo_cube",
    "type": "cube", 
    "position": [-3, 1, 0],
    "scale": [1, 1, 1],
    "color": [1.0, 0.2, 0.2]
}' > /dev/null

# Basic sphere - green (API-first format)
$HD1_CLIENT create-object "$SESSION_ID" '{
    "name": "demo_sphere",
    "type": "sphere",
    "position": [0, 1, 0],
    "scale": [1, 1, 1],
    "color": [0.2, 1.0, 0.2]
}' > /dev/null

# Basic cylinder - blue (API-first format)
$HD1_CLIENT create-object "$SESSION_ID" '{
    "name": "demo_cylinder", 
    "type": "cylinder",
    "position": [3, 1, 0],
    "scale": [1, 1, 1],
    "color": [0.2, 0.2, 1.0]
}' > /dev/null

# Basic cone - yellow (API-first format)
$HD1_CLIENT create-object "$SESSION_ID" '{
    "name": "demo_cone",
    "type": "cone",
    "position": [-1.5, 1, -3],
    "scale": [1, 1, 1],
    "color": [1.0, 1.0, 0.2]
}' > /dev/null

# Wireframe cube - white wireframe (API-first format)
$HD1_CLIENT create-object "$SESSION_ID" '{
    "name": "wireframe_cube",
    "type": "cube",
    "position": [1.5, 1, -3],
    "scale": [1, 1, 1],
    "color": [0.8, 0.8, 0.8]
}' > /dev/null

# Info panel (API-first format)
$HD1_CLIENT create-object "$SESSION_ID" '{
    "name": "shapes_label",
    "type": "plane", 
    "position": [0, 3, -1],
    "scale": [2, 2, 2],
    "color": [1.0, 1.0, 1.0]
}' > /dev/null

echo "HD1 Scene '$SCENE_NAME' loaded successfully"
echo "Objects created: 6"
echo "Session: $SESSION_ID"