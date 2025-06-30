#!/bin/bash
# ===================================================================
# HD1 Environment Script: Earth Surface
# ===================================================================
# Name: Earth Surface
# Description: Standard terrestrial environment with human-scale physics
# Scale: m
# Gravity: 9.8
# Atmosphere: air
# Complexity: simple
# Tags: terrestrial, standard, human-scale
# ===================================================================

source "${HD1_ROOT}/lib/hd1lib.sh" 2>/dev/null || {
    echo "ERROR: HD1 core library not found"
    exit 1
}

HD1_SESSION_ID="${1:-${HD1_SESSION_ID:-}}"

if [[ -z "$HD1_SESSION_ID" ]]; then
    echo "ERROR: Session ID required"
    exit 1
fi

echo "ENVIRONMENT: Applying Earth Surface environment to session $HD1_SESSION_ID"

# Environment physics setup
echo "PHYSICS: Standard Earth gravity (9.81 m/s²)"
echo "SCALE: Meter units - human scale interaction"
echo "ATMOSPHERE: Standard air pressure and composition"
echo "BOUNDARIES: 100m x 100m x 50m operating volume"

# Set up basic lighting for Earth surface
hd1::create_object "sun-light" "light" 10 10 5
hd1::create_object "ambient-light" "light" 0 0 0

# Create ground plane
hd1::create_object "ground" "plane" 0 0 0

# Set camera to human height
hd1::camera 0 1.7 3

echo "ENVIRONMENT: Earth Surface applied successfully"
echo "Environment applied: earth-surface"