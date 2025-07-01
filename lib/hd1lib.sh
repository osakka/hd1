#!/bin/bash
# ===================================================================
# HD1 Core Shell Function Library - AUTO-GENERATED
# ===================================================================
#
# GENERATED FROM: api.yaml specification
# SINGLE SOURCE OF TRUTH: All functions auto-generated from API spec
# PURPOSE: Standard shell wrapper for HD1 API endpoints
# 
# DO NOT EDIT MANUALLY - Regenerate with: make generate
# ===================================================================

# Configuration - respects environment variables and defaults
HD1_API_BASE="${HD1_API_BASE:-http://localhost:8080/api}"
HD1_SESSION_ID="${HD1_SESSION_ID:-${SESSION_ID:-session-19cdcfgj}}"

# Standard HTTP client with error handling
hd1::api_call() {
    local method="$1"
    local endpoint="$2"
    local payload="$3"
    local content_type="${4:-application/json}"
    
    local response
    if [[ -n "$payload" ]]; then
        response=$(curl -s -X "$method" "$HD1_API_BASE$endpoint" \
                        -H "Content-Type: $content_type" \
                        -d "$payload")
    else
        response=$(curl -s -X "$method" "$HD1_API_BASE$endpoint")
    fi
    
    # Standard JSON response parsing
    if echo "$response" | jq . >/dev/null 2>&1; then
        echo "$response" | jq -r '.message // .success // "Success"'
    else
        echo "ERROR: $response"
        return 1
    fi
}


# Auto-generated from DELETE /sessions/{sessionId}/objects/{objectName}
hd1::delete_object() {
    local name="$1"
    
    if [[ -z "$name" ]]; then
        echo "Usage: hd1::delete_object <name>"
        return 1
    fi
    
    hd1::api_call "DELETE" "/sessions/$HD1_SESSION_ID/objects/$name"
    echo "DELETE: Object $name"
}

# Auto-generated from GET /sessions/{sessionId}/objects/{objectName}
hd1::get_object() {
    local name="$1"
    
    if [[ -z "$name" ]]; then
        echo "Usage: hd1::get_object <name>"
        return 1
    fi
    
    hd1::api_call "GET" "/sessions/$HD1_SESSION_ID/objects/$name"
}

# Auto-generated from PUT /sessions/{sessionId}/objects/{objectName}
hd1::update_object() {
    local name="$1"
    local property="$2"
    local value="$3"
    
    if [[ -z "$name" || -z "$property" || -z "$value" ]]; then
        echo "Usage: hd1::update_object <name> <property> <value>"
        return 1
    fi
    
    local payload=$(cat <<EOF
{
    "$property": "$value"
}
EOF
)
    
    hd1::api_call "PUT" "/sessions/$HD1_SESSION_ID/objects/$name" "$payload"
    echo "UPDATE: Object $name property $property"
}

# Auto-generated from PUT /sessions/{sessionId}/camera/position
hd1::camera() {
    local x="$1" y="$2" z="$3"
    
    if [[ -z "$x" || -z "$y" || -z "$z" ]]; then
        echo "Usage: hd1::camera <x> <y> <z>"
        return 1
    fi
    
    local payload=$(cat <<EOF
{
    "x": $x,
    "y": $y,
    "z": $z
}
EOF
)
    
    hd1::api_call "PUT" "/sessions/$HD1_SESSION_ID/camera/position" "$payload"
    echo "CAMERA: Positioned at ($x,$y,$z)"
}

# Auto-generated from POST /browser/canvas
hd1::canvas_control() {
    local command="$1"
    shift
    local objects="$@"
    
    if [[ -z "$command" ]]; then
        echo "Usage: hd1::canvas_control <command> [objects...]"
        return 1
    fi
    
    local payload=$(cat <<EOF
{
    "command": "$command",
    "objects": [$objects]
}
EOF
)
    
    hd1::api_call "POST" "/browser/canvas" "$payload"
}

# Clear holodeck (uses canvas control)
hd1::clear() {
    echo "CLEAR: Clearing holodeck..."
    hd1::canvas_control "clear"
}

# Auto-generated from GET /sessions
hd1::list_sessions() {
    hd1::api_call "GET" "/sessions"
}

# Auto-generated from POST /sessions
hd1::create_session() {
    hd1::api_call "POST" "/sessions"
}

# Auto-generated from GET /sessions/{sessionId}
hd1::get_session() {
    local session_id="${1:-$HD1_SESSION_ID}"
    hd1::api_call "GET" "/sessions/$session_id"
}

# Auto-generated from GET /sessions/{sessionId}/objects
hd1::list_objects() {
    hd1::api_call "GET" "/sessions/$HD1_SESSION_ID/objects"
}

# Auto-generated from POST /sessions/{sessionId}/objects
hd1::create_object() {
    local name="$1"
    local type="$2" 
    local x="$3"
    local y="$4"
    local z="$5"
    
    if [[ -z "$name" || -z "$type" || -z "$x" || -z "$y" || -z "$z" ]]; then
        echo "Usage: hd1::create_object <name> <type> <x> <y> <z>"
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
    
    hd1::api_call "POST" "/sessions/$HD1_SESSION_ID/objects" "$payload"
    echo "OBJECT: $name at ($x,$y,$z)"
}


echo "HD1: Core Functions Loaded - AUTO-GENERATED FROM API SPEC"
echo "SPEC: Generated from api.yaml specification"
echo "SYNC: Single source of truth - Zero manual synchronization"
echo "FUNCS: create_object, update_object, camera, canvas_control, clear, list_objects"
echo "SESSION: create_session, get_session, init_world"
echo "STATUS: Bar-raising achieved"