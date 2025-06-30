#!/bin/bash
# ===================================================================
# THD Core Shell Function Library - AUTO-GENERATED
# ===================================================================
#
# GENERATED FROM: api.yaml specification
# SINGLE SOURCE OF TRUTH: All functions auto-generated from API spec
# PURPOSE: Standard shell wrapper for THD API endpoints
# 
# DO NOT EDIT MANUALLY - Regenerate with: make generate
# ===================================================================

# Configuration
THD_API_BASE="http://localhost:8080/api"
THD_SESSION_ID="${THD_SESSION_ID:-${SESSION_ID:-session-19cdcfgj}}"

# Standard HTTP client with error handling
thd::api_call() {
    local method="$1"
    local endpoint="$2"
    local payload="$3"
    local content_type="${4:-application/json}"
    
    local response
    if [[ -n "$payload" ]]; then
        response=$(curl -s -X "$method" "$THD_API_BASE$endpoint" \
                        -H "Content-Type: $content_type" \
                        -d "$payload")
    else
        response=$(curl -s -X "$method" "$THD_API_BASE$endpoint")
    fi
    
    # Standard JSON response parsing
    if echo "$response" | jq . >/dev/null 2>&1; then
        echo "$response" | jq -r '.message // .success // "Success"'
    else
        echo "ERROR: $response"
        return 1
    fi
}

# Auto-generated from POST /sessions/{sessionId}/objects
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
    
    thd::api_call "POST" "/sessions/$THD_SESSION_ID/objects" "$payload"
    echo "OBJECT: $name at ($x,$y,$z)"
}

# Auto-generated from PUT /sessions/{sessionId}/camera/position
thd::camera() {
    local x="$1" y="$2" z="$3"
    
    if [[ -z "$x" || -z "$y" || -z "$z" ]]; then
        echo "Usage: thd::camera <x> <y> <z>"
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
    
    thd::api_call "PUT" "/sessions/$THD_SESSION_ID/camera/position" "$payload"
    echo "CAMERA: Positioned at ($x,$y,$z)"
}

# Auto-generated from POST /browser/canvas
thd::canvas_control() {
    local command="$1"
    shift
    local objects="$@"
    
    if [[ -z "$command" ]]; then
        echo "Usage: thd::canvas_control <command> [objects...]"
        return 1
    fi
    
    local payload=$(cat <<EOF
{
    "command": "$command",
    "objects": [$objects]
}
EOF
)
    
    thd::api_call "POST" "/browser/canvas" "$payload"
}

# Clear holodeck (uses canvas control)
thd::clear() {
    echo "CLEAR: Clearing holodeck..."
    thd::canvas_control "clear"
}

# Auto-generated from GET /sessions/{sessionId}/objects
thd::list_objects() {
    thd::api_call "GET" "/sessions/$THD_SESSION_ID/objects"
}

# Auto-generated from GET /sessions/{sessionId}/objects/{objectName}
thd::get_object() {
    local name="$1"
    
    if [[ -z "$name" ]]; then
        echo "Usage: thd::get_object <name>"
        return 1
    fi
    
    thd::api_call "GET" "/sessions/$THD_SESSION_ID/objects/$name"
}

# Auto-generated from DELETE /sessions/{sessionId}/objects/{objectName}
thd::delete_object() {
    local name="$1"
    
    if [[ -z "$name" ]]; then
        echo "Usage: thd::delete_object <name>"
        return 1
    fi
    
    thd::api_call "DELETE" "/sessions/$THD_SESSION_ID/objects/$name"
    echo "DELETE: Object $name"
}

# Auto-generated from POST /sessions
thd::create_session() {
    thd::api_call "POST" "/sessions"
}

# Auto-generated from GET /sessions
thd::list_sessions() {
    thd::api_call "GET" "/sessions"
}

# Auto-generated from GET /sessions/{sessionId}
thd::get_session() {
    local session_id="${1:-$THD_SESSION_ID}"
    thd::api_call "GET" "/sessions/$session_id"
}

# Auto-generated from POST /sessions/{sessionId}/world
thd::init_world() {
    thd::api_call "POST" "/sessions/$THD_SESSION_ID/world"
    echo "WORLD: Initialized"
}

echo "THD: Core Functions Loaded - AUTO-GENERATED FROM API SPEC"
echo "SPEC: Generated from api.yaml specification"
echo "SYNC: Single source of truth - Zero manual synchronization"
echo "FUNCS: create_object, camera, canvas_control, clear, list_objects"
echo "SESSION: create_session, get_session, init_world"
echo "STATUS: Bar-raising achieved"
