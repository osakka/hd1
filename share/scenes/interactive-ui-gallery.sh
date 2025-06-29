#!/bin/bash

# =========================================================================
# THD Scene: Interactive UI Gallery - A-Frame Style Interface
# =========================================================================
#
# Recreates an A-Frame laser input UI with gallery-style poster frames
# Based on the Studio Ghibli movie poster UI example
#
# Usage: ./interactive-ui-gallery.sh [SESSION_ID]
# =========================================================================

set -euo pipefail

# Scene configuration
SCENE_NAME="Interactive UI Gallery"
SCENE_DESCRIPTION="Gallery-style interface with poster frames and info panel, based on A-Frame laser input UI"

# Get session ID from argument or use active session
SESSION_ID="${1:-${THD_SESSION:-}}"

if [[ -z "$SESSION_ID" ]]; then
    echo "ERROR: No session ID provided and THD_SESSION not set"
    exit 1
fi

# Set THD_ROOT and source functions
THD_ROOT="/opt/holo-deck"
source "${THD_ROOT}/lib/thdlib.sh" 2>/dev/null || {
    echo "ERROR: THD functions not available"
    exit 1
}

echo "Creating $SCENE_NAME scene..."

# Background environment - dark purple sphere
curl -X POST "http://localhost:8080/api/sessions/$SESSION_ID/objects" \
  -H "Content-Type: application/json" \
  -d '{"name":"background_sphere","type":"sphere","x":0,"y":1.6,"z":0,"color":{"r":0.13,"g":0.07,"b":0.2,"a":1.0},"material":{"shader":"standard","transparent":false}}' \
  -s > /dev/null

# Three poster frames at eye level
curl -X POST "http://localhost:8080/api/sessions/$SESSION_ID/objects" \
  -H "Content-Type: application/json" \
  -d '{"name":"left_poster_frame","type":"plane","x":-2,"y":1.6,"z":-2.5,"color":{"r":1.0,"g":1.0,"b":1.0,"a":1.0},"material":{"shader":"standard","metalness":0.0,"roughness":0.3}}' \
  -s > /dev/null

curl -X POST "http://localhost:8080/api/sessions/$SESSION_ID/objects" \
  -H "Content-Type: application/json" \
  -d '{"name":"center_poster_frame","type":"plane","x":0,"y":1.6,"z":-2.5,"color":{"r":1.0,"g":1.0,"b":1.0,"a":1.0},"material":{"shader":"standard","metalness":0.0,"roughness":0.3}}' \
  -s > /dev/null

curl -X POST "http://localhost:8080/api/sessions/$SESSION_ID/objects" \
  -H "Content-Type: application/json" \
  -d '{"name":"right_poster_frame","type":"plane","x":2,"y":1.6,"z":-2.5,"color":{"r":1.0,"g":1.0,"b":1.0,"a":1.0},"material":{"shader":"standard","metalness":0.0,"roughness":0.3}}' \
  -s > /dev/null

# Info panel for displaying content
curl -X POST "http://localhost:8080/api/sessions/$SESSION_ID/objects" \
  -H "Content-Type: application/json" \
  -d '{"name":"info_panel","type":"plane","x":0,"y":1.6,"z":-2,"color":{"r":0.2,"g":0.2,"b":0.2,"a":1.0},"material":{"shader":"standard","metalness":0.1,"roughness":0.8}}' \
  -s > /dev/null

# Title text
curl -X POST "http://localhost:8080/api/sessions/$SESSION_ID/objects" \
  -H "Content-Type: application/json" \
  -d '{"name":"title_text_plane","type":"plane","x":-1,"y":2.2,"z":-1.9,"color":{"r":1.0,"g":1.0,"b":1.0,"a":1.0},"text":"Interactive UI Gallery","material":{"shader":"standard","transparent":true}}' \
  -s > /dev/null

# Standard lighting setup
curl -X POST "http://localhost:8080/api/sessions/$SESSION_ID/objects" \
  -H "Content-Type: application/json" \
  -d '{"name":"ambient_light","type":"light","x":0,"y":3,"z":0,"lightType":"ambient","intensity":0.3,"color":{"r":0.9,"g":0.9,"b":1.0,"a":1.0}}' \
  -s > /dev/null

curl -X POST "http://localhost:8080/api/sessions/$SESSION_ID/objects" \
  -H "Content-Type: application/json" \
  -d '{"name":"directional_light","type":"light","x":1,"y":4,"z":1,"lightType":"directional","intensity":0.8,"color":{"r":1.0,"g":0.95,"b":0.8,"a":1.0}}' \
  -s > /dev/null

echo "THD Scene '$SCENE_NAME' loaded successfully"
echo "Objects created: 8"
echo "Session: $SESSION_ID"