#!/bin/bash
# THD HOLODECK SHELL FUNCTIONS
# Production holodeck scenario building toolkit

# Configuration
THD_API_BASE="http://localhost:8080/api"
THD_SESSION_ID="${SESSION_ID:-session-19cdcfgj}"

# Core object creation function
thd::create_object() {
    local name="$1"
    local type="$2" 
    local x="$3"
    local y="$4"
    local z="$5"
    
    if [[ -z "$name" || -z "$type" || -z "$x" || -z "$y" || -z "$z" ]]; then
        echo "Usage: thd::create_object <name> <type> <x> <y> <z>"
        return 1
    fi
    
    local payload=$(cat <<EOF
{
    "name": "$name",
    "type": "$type", 
    "x": $x,
    "y": $y,
    "z": $z
}
EOF
)
    
    curl -s -X POST "$THD_API_BASE/sessions/$THD_SESSION_ID/objects" \
         -H "Content-Type: application/json" \
         -d "$payload" | jq -r '.message // "Created"'
    
    echo "📦 $name at ($x,$y,$z)"
}

# Bulk canvas control 
thd::canvas_control() {
    local command="$1"
    shift
    local objects="$@"
    
    curl -s -X POST "$THD_API_BASE/browser/canvas" \
         -H "Content-Type: application/json" \
         -d "{\"command\": \"$command\", \"objects\": [$objects]}"
}

# Clear holodeck
thd::clear() {
    echo "🧹 Clearing holodeck..."
    thd::canvas_control "clear"
}

# Camera positioning
thd::camera() {
    local x="$1" y="$2" z="$3"
    curl -s -X PUT "$THD_API_BASE/sessions/$THD_SESSION_ID/camera/position" \
         -H "Content-Type: application/json" \
         -d "{\"x\": $x, \"y\": $y, \"z\": $z}"
    echo "📷 Camera positioned at ($x,$y,$z)"
}

# Advanced A-Frame THD Functions

# Create light sources
thd::create_light() {
    local name="$1"
    local type="$2"  # ambient, directional, point, spot
    local x="$3" y="$4" z="$5"
    local intensity="${6:-1.0}"
    local color="${7:-#ffffff}"
    
    local payload=$(cat <<EOF
{
    "name": "$name",
    "type": "light",
    "lightType": "$type",
    "x": $x,
    "y": $y,
    "z": $z,
    "intensity": $intensity,
    "color": {"r": 1.0, "g": 1.0, "b": 1.0, "a": 1.0}
}
EOF
)
    
    curl -s -X POST "$THD_API_BASE/sessions/$THD_SESSION_ID/objects" \
         -H "Content-Type: application/json" \
         -d "$payload" | jq -r '.message // "Light Created"'
    
    echo "💡 $name light ($type) at ($x,$y,$z)"
}

# Create physics-enabled objects
thd::create_physics() {
    local name="$1"
    local type="$2"
    local x="$3" y="$4" z="$5"
    local mass="${6:-1.0}"
    local physics_type="${7:-dynamic}" # static, dynamic, kinematic
    
    local payload=$(cat <<EOF
{
    "name": "$name",
    "type": "$type",
    "x": $x,
    "y": $y,
    "z": $z,
    "physics": {
        "enabled": true,
        "mass": $mass,
        "type": "$physics_type"
    }
}
EOF
)
    
    curl -s -X POST "$THD_API_BASE/sessions/$THD_SESSION_ID/objects" \
         -H "Content-Type: application/json" \
         -d "$payload" | jq -r '.message // "Physics Object Created"'
    
    echo "⚡ $name (physics: $physics_type, mass: $mass) at ($x,$y,$z)"
}

# Create textured/material objects
thd::create_material() {
    local name="$1"
    local type="$2"
    local x="$3" y="$4" z="$5"
    local shader="${6:-standard}"
    local metalness="${7:-0.1}"
    local roughness="${8:-0.7}"
    
    local payload=$(cat <<EOF
{
    "name": "$name",
    "type": "$type",
    "x": $x,
    "y": $y,
    "z": $z,
    "material": {
        "shader": "$shader",
        "metalness": $metalness,
        "roughness": $roughness
    }
}
EOF
)
    
    curl -s -X POST "$THD_API_BASE/sessions/$THD_SESSION_ID/objects" \
         -H "Content-Type: application/json" \
         -d "$payload" | jq -r '.message // "Material Object Created"'
    
    echo "🎨 $name ($shader material) at ($x,$y,$z)"
}

# Create sky/environment
thd::create_sky() {
    local name="$1"
    local color="${2:-#87CEEB}" # Default sky blue
    
    local payload=$(cat <<EOF
{
    "name": "$name",
    "type": "sky",
    "x": 0,
    "y": 0,
    "z": 0,
    "color": {"r": 0.5, "g": 0.8, "b": 0.9, "a": 1.0}
}
EOF
)
    
    curl -s -X POST "$THD_API_BASE/sessions/$THD_SESSION_ID/objects" \
         -H "Content-Type: application/json" \
         -d "$payload" | jq -r '.message // "Sky Created"'
    
    echo "🌤️  $name sky environment created"
}

# Enhanced object creation with materials
thd::create_enhanced() {
    local name="$1"
    local type="$2"
    local x="$3" y="$4" z="$5"
    local color_r="${6:-0.2}"
    local color_g="${7:-0.8}" 
    local color_b="${8:-0.2}"
    local metalness="${9:-0.1}"
    local roughness="${10:-0.7}"
    
    local payload=$(cat <<EOF
{
    "name": "$name",
    "type": "$type",
    "x": $x,
    "y": $y,
    "z": $z,
    "color": {
        "r": $color_r,
        "g": $color_g,
        "b": $color_b,
        "a": 1.0
    },
    "material": {
        "shader": "standard",
        "metalness": $metalness,
        "roughness": $roughness
    },
    "lighting": {
        "castShadow": true,
        "receiveShadow": true
    }
}
EOF
)
    
    curl -s -X POST "$THD_API_BASE/sessions/$THD_SESSION_ID/objects" \
         -H "Content-Type: application/json" \
         -d "$payload" | jq -r '.message // "Enhanced Object Created"'
    
    echo "✨ $name (enhanced) at ($x,$y,$z)"
}

echo "🎯 THD Holodeck Functions Loaded (A-Frame Enhanced)"
echo "💡 New functions: thd::create_light, thd::create_physics, thd::create_material, thd::create_sky, thd::create_enhanced"