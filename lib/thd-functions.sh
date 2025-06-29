#!/bin/bash
# THD HOLODECK SHELL FUNCTIONS
# Production holodeck scenario building toolkit

# Configuration
THD_API_BASE="http://localhost:8080/api"
THD_SESSION_ID="${THD_SESSION_ID:-${SESSION_ID:-session-19cdcfgj}}"

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
    
    local response=$(curl -s -X POST "$THD_API_BASE/sessions/$THD_SESSION_ID/objects" \
         -H "Content-Type: application/json" \
         -d "$payload")
    
    # Parse JSON response safely
    if echo "$response" | jq . >/dev/null 2>&1; then
        echo "$response" | jq -r '.message // "Created"'
    else
        echo "ERROR: $response"
    fi
    
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
    
    local response=$(curl -s -X POST "$THD_API_BASE/sessions/$THD_SESSION_ID/objects" \
         -H "Content-Type: application/json" \
         -d "$payload")
    
    # Parse JSON response safely
    if echo "$response" | jq . >/dev/null 2>&1; then
        echo "$response" | jq -r '.message // "Light Created"'
    else
        echo "ERROR: $response"
    fi
    
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
    
    local response=$(curl -s -X POST "$THD_API_BASE/sessions/$THD_SESSION_ID/objects" \
         -H "Content-Type: application/json" \
         -d "$payload")
    
    # Parse JSON response safely
    if echo "$response" | jq . >/dev/null 2>&1; then
        echo "$response" | jq -r '.message // "Physics Object Created"'
    else
        echo "ERROR: $response"
    fi
    
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
    
    local response=$(curl -s -X POST "$THD_API_BASE/sessions/$THD_SESSION_ID/objects" \
         -H "Content-Type: application/json" \
         -d "$payload")
    
    # Parse JSON response safely
    if echo "$response" | jq . >/dev/null 2>&1; then
        echo "$response" | jq -r '.message // "Material Object Created"'
    else
        echo "ERROR: $response"
    fi
    
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
    
    local response=$(curl -s -X POST "$THD_API_BASE/sessions/$THD_SESSION_ID/objects" \
         -H "Content-Type: application/json" \
         -d "$payload")
    
    # Parse JSON response safely
    if echo "$response" | jq . >/dev/null 2>&1; then
        echo "$response" | jq -r '.message // "Sky Created"'
    else
        echo "ERROR: $response"
    fi
    
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
    
    local response=$(curl -s -X POST "$THD_API_BASE/sessions/$THD_SESSION_ID/objects" \
         -H "Content-Type: application/json" \
         -d "$payload")
    
    # Parse JSON response safely
    if echo "$response" | jq . >/dev/null 2>&1; then
        echo "$response" | jq -r '.message // "Enhanced Object Created"'
    else
        echo "ERROR: $response"
    fi
    
    echo "✨ $name (enhanced) at ($x,$y,$z)"
}

# Create 3D text
thd::create_text() {
    local name="$1"
    local text="$2"
    local x="$3" y="$4" z="$5"
    local color_r="${6:-1.0}"
    local color_g="${7:-1.0}"
    local color_b="${8:-1.0}"
    
    local payload=$(cat <<EOF
{
    "name": "$name",
    "type": "text",
    "text": "$text",
    "x": $x,
    "y": $y,
    "z": $z,
    "color": {
        "r": $color_r,
        "g": $color_g,
        "b": $color_b,
        "a": 1.0
    }
}
EOF
)
    
    local response=$(curl -s -X POST "$THD_API_BASE/sessions/$THD_SESSION_ID/objects" \
         -H "Content-Type: application/json" \
         -d "$payload")
    
    # Parse JSON response safely
    if echo "$response" | jq . >/dev/null 2>&1; then
        echo "$response" | jq -r '.message // "Text Created"'
    else
        echo "ERROR: $response"
    fi
    
    echo "📝 $name text ('$text') at ($x,$y,$z)"
}

# Create particle effects
thd::create_particles() {
    local name="$1"
    local type="$2"  # fire, smoke, sparkle
    local x="$3" y="$4" z="$5"
    local count="${6:-500}"
    
    local payload=$(cat <<EOF
{
    "name": "$name",
    "type": "particle",
    "particleType": "$type",
    "x": $x,
    "y": $y,
    "z": $z,
    "count": $count
}
EOF
)
    
    local response=$(curl -s -X POST "$THD_API_BASE/sessions/$THD_SESSION_ID/objects" \
         -H "Content-Type: application/json" \
         -d "$payload")
    
    # Parse JSON response safely
    if echo "$response" | jq . >/dev/null 2>&1; then
        echo "$response" | jq -r '.message // "Particles Created"'
    else
        echo "ERROR: $response"
    fi
    
    echo "✨ $name particles ($type) at ($x,$y,$z)"
}

# Ultimate holodeck object creation
thd::create_ultimate() {
    local name="$1"
    local type="$2"
    local x="$3" y="$4" z="$5"
    local color_r="${6:-0.5}"
    local color_g="${7:-0.8}"
    local color_b="${8:-1.0}"
    local metalness="${9:-0.5}"
    local roughness="${10:-0.3}"
    local emissive="${11:-false}"
    
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
        "roughness": $roughness,
        "emissive": $emissive
    },
    "physics": {
        "enabled": true,
        "mass": 1.0,
        "type": "dynamic"
    },
    "lighting": {
        "castShadow": true,
        "receiveShadow": true
    }
}
EOF
)
    
    local response=$(curl -s -X POST "$THD_API_BASE/sessions/$THD_SESSION_ID/objects" \
         -H "Content-Type: application/json" \
         -d "$payload")
    
    # Parse JSON response safely
    if echo "$response" | jq . >/dev/null 2>&1; then
        echo "$response" | jq -r '.message // "Ultimate Object Created"'
    else
        echo "ERROR: $response"
    fi
    
    echo "🌟 $name (ultimate holodeck object) at ($x,$y,$z)"
}

echo "🎯 THD Holodeck Functions Loaded (ULTIMATE A-Frame)"
echo "💡 Functions: thd::create_light, thd::create_physics, thd::create_material, thd::create_sky, thd::create_enhanced"
echo "🌟 Ultimate: thd::create_text, thd::create_particles, thd::create_ultimate"