#!/bin/bash

# =========================================================================
# THD Scene: My Custom Scene - A customized version of basic shapes scene
# =========================================================================
#
# A customized version of basic shapes scene
#
# Usage: ./my-custom-scene.sh [SESSION_ID]
# Auto-generated from session state on 2025-06-29 18:13:56
# =========================================================================

set -euo pipefail

# Scene configuration
SCENE_NAME="My Custom Scene"
SCENE_DESCRIPTION="A customized version of basic shapes scene"

# Get session ID from argument or use active session
SESSION_ID="${1:-${THD_SESSION:-}}"

if [[ -z "$SESSION_ID" ]]; then
    echo "Error: Session ID required" >&2
    exit 1
fi

# Path to auto-generated THD client
THD_CLIENT="/opt/holo-deck/build/bin/thd-client"

echo "Creating my custom scene scene..."


# Demo Cube - cube
$THD_CLIENT create-object "$SESSION_ID" '{
    "name": "demo_cube",
    "type": "cube",
    "x": -3, "y": 1, "z": 0,
    "scale": 1
}' > /dev/null

# Demo Sphere - sphere
$THD_CLIENT create-object "$SESSION_ID" '{
    "name": "demo_sphere",
    "type": "sphere",
    "x": 0, "y": 1, "z": 0,
    "scale": 1
}' > /dev/null

# Demo Cylinder - cylinder
$THD_CLIENT create-object "$SESSION_ID" '{
    "name": "demo_cylinder",
    "type": "cylinder",
    "x": 3, "y": 1, "z": 0,
    "scale": 1
}' > /dev/null

# Demo Cone - cone
$THD_CLIENT create-object "$SESSION_ID" '{
    "name": "demo_cone",
    "type": "cone",
    "x": -1.5, "y": 1, "z": -3,
    "scale": 1
}' > /dev/null

# Wireframe Cube - cube
$THD_CLIENT create-object "$SESSION_ID" '{
    "name": "wireframe_cube",
    "type": "cube",
    "x": 1.5, "y": 1, "z": -3,
    "scale": 1
}' > /dev/null

# Shapes Label - plane
$THD_CLIENT create-object "$SESSION_ID" '{
    "name": "shapes_label",
    "type": "plane",
    "x": 0, "y": 3, "z": -1,
    "scale": 1
}' > /dev/null

echo "THD Scene '$SCENE_NAME' loaded successfully"
echo "Objects created: 6"
echo "Session: $SESSION_ID"
