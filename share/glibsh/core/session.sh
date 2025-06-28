#!/bin/bash

# =========================================================================
# GLIBSH Core: Session Management - Professional Session Orchestration
# =========================================================================
#
# Session lifecycle management with surgical precision:
# - Session creation and validation
# - Session health monitoring  
# - Error recovery and cleanup
# - Professional session state management
# =========================================================================

set -euo pipefail

# =========================================================================
# SESSION MANAGEMENT CORE
# =========================================================================

# Global session state
CURRENT_SESSION_ID=""

# Get or create session with professional validation
ensure_session() {
    log "INFO" "Ensuring active THD session exists"
    
    # Try to get existing sessions first
    local sessions_response
    sessions_response=$(call_thd_api "sessions" "GET") || {
        error_exit "Failed to query existing sessions"
    }
    
    # Parse sessions (expecting JSON with sessions array)
    if echo "${sessions_response}" | jq -e '.sessions | length > 0' > /dev/null 2>&1; then
        # Use first available session
        CURRENT_SESSION_ID=$(echo "${sessions_response}" | jq -r '.sessions[0].id')
        log "INFO" "Using existing session: ${CURRENT_SESSION_ID}"
        
        # Validate session is healthy
        validate_session "${CURRENT_SESSION_ID}" || {
            log "WARN" "Existing session unhealthy, creating new session"
            create_new_session
        }
    else
        # No sessions exist, create new one
        log "INFO" "No existing sessions found, creating new session"
        create_new_session
    fi
    
    # Ensure session has initialized world
    ensure_world_initialized "${CURRENT_SESSION_ID}"
    
    log "INFO" "Session ready: ${CURRENT_SESSION_ID}"
    echo "${CURRENT_SESSION_ID}"
}

# Create new session with professional error handling
create_new_session() {
    log "INFO" "Creating new THD session"
    
    local create_response
    create_response=$(call_thd_api "create-session" "POST") || {
        error_exit "Failed to create new session"
    }
    
    # Extract session ID from response
    CURRENT_SESSION_ID=$(echo "${create_response}" | jq -r '.session_id // .id // empty')
    
    if [[ -z "${CURRENT_SESSION_ID}" ]]; then
        log "ERROR" "Session creation response: ${create_response}"
        error_exit "Failed to extract session ID from creation response"
    fi
    
    log "INFO" "Created session: ${CURRENT_SESSION_ID}"
}

# Validate session health and responsiveness
validate_session() {
    local session_id="$1"
    
    log "DEBUG" "Validating session health: ${session_id}"
    
    # Query session details
    local session_response
    session_response=$(call_thd_api "sessions/${session_id}" "GET" 2>/dev/null) || {
        log "WARN" "Session ${session_id} not responding"
        return 1
    }
    
    # Check if session is active
    local session_status
    session_status=$(echo "${session_response}" | jq -r '.status // "unknown"')
    
    if [[ "${session_status}" != "active" ]]; then
        log "WARN" "Session ${session_id} status: ${session_status}"
        return 1
    fi
    
    log "DEBUG" "Session ${session_id} validated successfully"
    return 0
}

# Ensure session has initialized world coordinate system
ensure_world_initialized() {
    local session_id="$1"
    
    log "DEBUG" "Ensuring world initialized for session: ${session_id}"
    
    # Check if world exists
    local world_response
    world_response=$(call_thd_api "sessions/${session_id}/world" "GET" 2>/dev/null) || {
        log "INFO" "World not initialized, initializing now"
        initialize_session_world "${session_id}"
        return
    }
    
    # Validate world has proper bounds
    local world_bounds
    world_bounds=$(echo "${world_response}" | jq -r '.bounds // empty')
    
    if [[ -z "${world_bounds}" ]]; then
        log "INFO" "World bounds not set, re-initializing world"
        initialize_session_world "${session_id}"
    else
        log "DEBUG" "World already initialized with bounds: ${world_bounds}"
    fi
}

# Initialize world coordinate system for session
initialize_session_world() {
    local session_id="$1"
    
    log "INFO" "Initializing world coordinate system for session: ${session_id}"
    
    # Initialize with THD standard parameters (25x25x25 grid, [-12,+12] bounds)
    local world_payload='{
        "size": 25,
        "transparency": 0.1,
        "grid_spacing": 1.0,
        "camera_x": 10,
        "camera_y": 10,
        "camera_z": 10
    }'
    
    local init_response
    init_response=$(call_thd_api "sessions/${session_id}/world" "POST" "${world_payload}") || {
        error_exit "Failed to initialize world for session ${session_id}"
    }
    
    log "INFO" "World initialized successfully for session: ${session_id}"
}

# Get current session ID (must be called after ensure_session)
get_current_session() {
    if [[ -z "${CURRENT_SESSION_ID}" ]]; then
        error_exit "No active session. Call ensure_session() first."
    fi
    echo "${CURRENT_SESSION_ID}"
}

# Clean session state (for cleanup/reset scenarios)
reset_session_state() {
    log "INFO" "Resetting GLIBSH session state"
    CURRENT_SESSION_ID=""
}

# Delete session with proper cleanup
delete_session() {
    local session_id="$1"
    
    log "INFO" "Deleting session: ${session_id}"
    
    call_thd_api "sessions/${session_id}" "DELETE" || {
        log "WARN" "Failed to delete session ${session_id} (may already be deleted)"
    }
    
    # Reset state if this was our current session
    if [[ "${session_id}" == "${CURRENT_SESSION_ID}" ]]; then
        reset_session_state
    fi
    
    log "INFO" "Session deletion completed: ${session_id}"
}

# List all active sessions
list_sessions() {
    log "INFO" "Listing all active THD sessions"
    
    local sessions_response
    sessions_response=$(call_thd_api "sessions" "GET") || {
        error_exit "Failed to list sessions"
    }
    
    echo "${sessions_response}" | jq -r '.sessions[] | "ID: \(.id) | Status: \(.status) | Created: \(.created_at)"'
}

# Professional session information display
show_session_info() {
    local session_id="${1:-${CURRENT_SESSION_ID}}"
    
    if [[ -z "${session_id}" ]]; then
        echo "No active session"
        return 1
    fi
    
    log "INFO" "Retrieving session information: ${session_id}"
    
    local session_response
    session_response=$(call_thd_api "sessions/${session_id}" "GET") || {
        error_exit "Failed to get session information"
    }
    
    echo "Session Information:"
    echo "==================="
    echo "${session_response}" | jq -r '
        "ID: \(.id)",
        "Status: \(.status)", 
        "Created: \(.created_at)",
        "Objects: \(.object_count // 0)",
        "World: \(if .world then "Initialized" else "Not initialized" end)"
    '
}

# =========================================================================
# SESSION VALIDATION & ERROR RECOVERY
# =========================================================================

# Comprehensive session health check
health_check_session() {
    local session_id="${1:-${CURRENT_SESSION_ID}}"
    
    if [[ -z "${session_id}" ]]; then
        echo "❌ No session to check"
        return 1
    fi
    
    log "INFO" "Performing comprehensive health check: ${session_id}"
    
    # Test session responsiveness
    if ! validate_session "${session_id}"; then
        echo "❌ Session ${session_id} failed health check"
        return 1
    fi
    
    # Test world initialization
    local world_response
    world_response=$(call_thd_api "sessions/${session_id}/world" "GET" 2>/dev/null) || {
        echo "⚠️  Session ${session_id} world not initialized"
        return 1
    }
    
    # Test object listing
    local objects_response
    objects_response=$(call_thd_api "sessions/${session_id}/objects" "GET" 2>/dev/null) || {
        echo "⚠️  Session ${session_id} objects not accessible"
        return 1
    }
    
    echo "✅ Session ${session_id} passed comprehensive health check"
    return 0
}

# Recovery from session failures
recover_session() {
    log "WARN" "Attempting session recovery"
    
    # Reset current state
    reset_session_state
    
    # Try to create fresh session
    create_new_session || {
        error_exit "Session recovery failed - unable to create new session"
    }
    
    log "INFO" "Session recovery successful: ${CURRENT_SESSION_ID}"
}

# =========================================================================
# PROFESSIONAL LOGGING & MONITORING
# =========================================================================

# Log session operation with context
log_session_operation() {
    local operation="$1"
    local session_id="$2"
    local result="$3"
    
    log "INFO" "SESSION_OP: ${operation} | Session: ${session_id} | Result: ${result}"
}