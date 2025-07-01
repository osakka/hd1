#!/bin/bash
# ===================================================================
# HD1 Lighting Library: Point Light Component
# ===================================================================
# Name: Point Light
# Description: Creates realistic point light sources with proper A-Frame integration
# Type: lighting_utility
# API Integration: Uses HD1 API calls internally (tied API pattern)
# ===================================================================

# Try multiple paths for HD1 core library
if [[ -f "${HD1_ROOT}/lib/hd1lib.sh" ]]; then
    source "${HD1_ROOT}/lib/hd1lib.sh"
elif [[ -f "/opt/hd1/lib/hd1lib.sh" ]]; then
    source "/opt/hd1/lib/hd1lib.sh"
elif [[ -f "../lib/hd1lib.sh" ]]; then
    source "../lib/hd1lib.sh"
else
    echo "ERROR: HD1 core library not found"
    exit 1
fi

# Point light creation function - called by API and calls API
# This demonstrates the bidirectional script bridge pattern
create_point_light() {
    local name="$1"
    local x="$2" 
    local y="$3"
    local z="$4"
    local color="${5:-#ffffff}"
    local intensity="${6:-1.0}"
    local distance="${7:-0}"
    local decay="${8:-2}"
    local session_id="${9:-${HD1_SESSION_ID}}"
    
    if [[ -z "$name" || -z "$x" || -z "$y" || -z "$z" ]]; then
        echo "ERROR: create_point_light requires name, x, y, z parameters"
        return 1
    fi
    
    if [[ -z "$session_id" ]]; then
        echo "ERROR: Session ID required for lighting operations"
        return 1
    fi
    
    echo "LIGHTING: Creating point light '$name' at ($x,$y,$z)"
    
    # VALIDATION: Tied API Architecture Pattern
    # This script is called BY the API and CALLS the API internally
    
    # Call HD1 API to create the light object (script → API)
    local light_data=$(cat <<EOF
{
    "name": "$name",
    "type": "light",
    "x": $x,
    "y": $y,
    "z": $z,
    "visible": true,
    "light": {
        "type": "point",
        "color": "$color",
        "intensity": $intensity,
        "distance": $distance,
        "decay": $decay
    }
}
EOF
)
    
    # Execute API call using HD1 API client
    hd1::api_call "POST" "/sessions/$session_id/objects" "$light_data"
    
    if [[ $? -eq 0 ]]; then
        echo "LIGHT: Point light '$name' created successfully"
        echo "PROPERTIES: Color=$color, Intensity=$intensity, Distance=$distance, Decay=$decay"
        echo "API_VALIDATION: ✅ Script called API successfully (tied architecture)"
        return 0
    else
        echo "ERROR: Failed to create point light '$name'"
        echo "API_VALIDATION: ❌ Script → API call failed"
        return 1
    fi
}

# Color temperature utility for realistic lighting
kelvin_to_rgb() {
    local kelvin="$1"
    
    case $kelvin in
        1000|"candle") echo "#ff8c00" ;;
        2700|"warm_white") echo "#ffe5b4" ;;
        3000|"soft_white") echo "#fff2d6" ;;
        4000|"cool_white") echo "#fff8f0" ;;
        5500|"daylight") echo "#ffffff" ;;
        6500|"bright_daylight") echo "#f0f8ff" ;;
        10000|"sky") echo "#87ceeb" ;;
        *) echo "#ffffff" ;;
    esac
}

# Preset lighting configurations
create_warm_bulb() {
    local name="$1"
    local x="$2" y="$3" z="$4"
    local intensity="${5:-1.0}"
    
    local warm_color=$(kelvin_to_rgb "warm_white")
    create_point_light "$name" "$x" "$y" "$z" "$warm_color" "$intensity" "5.0" "2"
}

create_daylight() {
    local name="$1"
    local x="$2" y="$3" z="$4"
    local intensity="${5:-0.8}"
    
    local daylight_color=$(kelvin_to_rgb "daylight")
    create_point_light "$name" "$x" "$y" "$z" "$daylight_color" "$intensity" "10.0" "1"
}

echo "HD1: Point Light Library Loaded"
echo "FUNCTIONS: create_point_light, kelvin_to_rgb, create_warm_bulb, create_daylight"
echo "ARCHITECTURE: Tied API pattern - script calls API internally"
echo "VALIDATION: Ready for API trigger and API consumption"