#!/bin/bash
# ===================================================================
# HD1 Environment Script: Underwater
# ===================================================================
# Name: Underwater
# Description: Submerged liquid environment with fluid dynamics
# Scale: m
# Gravity: 8.8
# Atmosphere: liquid
# Complexity: moderate
# Tags: underwater, liquid, fluid-dynamics
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

echo "ENVIRONMENT: Applying Underwater environment to session $HD1_SESSION_ID"

# Environment physics setup
echo "PHYSICS: Reduced gravity (2.0 m/s²) with buoyancy"
echo "SCALE: Meter units - underwater exploration"
echo "ATMOSPHERE: Liquid water medium"
echo "BOUNDARIES: 50m x 50m x 30m operating volume"

# Set up underwater lighting
hd1::create_object "surface-light" "light" 0 25 0
hd1::create_object "deep-ambient" "light" 0 0 0

# Create water surface and seafloor
hd1::create_object "water-surface" "plane" 0 25 0
hd1::create_object "seafloor" "plane" 0 -5 0

# Set camera for underwater view
hd1::camera 0 10 15

echo "ENVIRONMENT: Underwater applied successfully"
echo "Environment applied: underwater"