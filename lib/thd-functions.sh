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

echo "🎯 THD Holodeck Functions Loaded"