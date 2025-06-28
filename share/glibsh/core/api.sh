#!/bin/bash

# =========================================================================
# GLIBSH Core: API Abstraction - Professional THD Client Orchestration
# =========================================================================
#
# High-level API abstraction layer providing:
# - Intelligent thd-client wrapper functions
# - Response parsing and validation
# - Error handling and recovery
# - Professional request/response management
# =========================================================================

set -euo pipefail

# =========================================================================
# API CORE ORCHESTRATION
# =========================================================================

# Central API call function with professional error handling
call_thd_api() {
    local endpoint="$1"
    local method="${2:-GET}"
    local payload="${3:-}"
    
    log "DEBUG" "API Call: ${method} ${endpoint}"
    
    # Validate THD client exists
    if [[ ! -x "${THD_CLIENT}" ]]; then
        error_exit "THD client not found or not executable: ${THD_CLIENT}"
    fi
    
    local response
    local exit_code
    
    case "${method}" in
        "GET")
            response=$(call_thd_get "${endpoint}")
            exit_code=$?
            ;;
        "POST")
            response=$(call_thd_post "${endpoint}" "${payload}")
            exit_code=$?
            ;;
        "PUT")
            response=$(call_thd_put "${endpoint}" "${payload}")
            exit_code=$?
            ;;
        "DELETE")
            response=$(call_thd_delete "${endpoint}")
            exit_code=$?
            ;;
        *)
            error_exit "Unsupported HTTP method: ${method}"
            ;;
    esac
    
    # Validate response and exit code
    if [[ ${exit_code} -ne 0 ]]; then
        log "ERROR" "API call failed: ${method} ${endpoint} (exit code: ${exit_code})"
        return ${exit_code}
    fi
    
    # Validate JSON response format
    if ! echo "${response}" | jq . >/dev/null 2>&1; then
        log "WARN" "API response not valid JSON: ${response}"
        # Don't fail - some endpoints might return non-JSON
    fi
    
    log "DEBUG" "API Response: ${response}"
    echo "${response}"
    return 0
}

# =========================================================================
# HTTP METHOD IMPLEMENTATIONS
# =========================================================================

# GET request wrapper
call_thd_get() {
    local endpoint="$1"
    
    case "${endpoint}" in
        "sessions")
            "${THD_CLIENT}" sessions 2>/dev/null || echo '{"sessions":[]}'
            ;;
        "sessions/"*)
            # Extract session ID for specific session queries
            local session_id="${endpoint#sessions/}"
            session_id="${session_id%%/*}"  # Remove any sub-paths
            
            # Use direct curl for specific session endpoints
            curl -s "http://localhost:8080/api/${endpoint}" 2>/dev/null || echo '{"error":"not_found"}'
            ;;
        *)
            # Generic GET using curl
            curl -s "http://localhost:8080/api/${endpoint}" 2>/dev/null || echo '{"error":"endpoint_not_found"}'
            ;;
    esac
}

# POST request wrapper  
call_thd_post() {
    local endpoint="$1"
    local payload="${2:-{}}"
    
    case "${endpoint}" in
        "create-session"|"sessions")
            "${THD_CLIENT}" create-session 2>/dev/null || echo '{"error":"session_creation_failed"}'
            ;;
        "sessions/"*)
            # POST to specific session endpoints with payload
            if [[ -n "${payload}" && "${payload}" != "{}" ]]; then
                curl -s -X POST \
                    -H "Content-Type: application/json" \
                    -d "${payload}" \
                    "http://localhost:8080/api/${endpoint}" 2>/dev/null || echo '{"error":"post_failed"}'
            else
                curl -s -X POST "http://localhost:8080/api/${endpoint}" 2>/dev/null || echo '{"error":"post_failed"}'
            fi
            ;;
        *)
            # Generic POST
            if [[ -n "${payload}" && "${payload}" != "{}" ]]; then
                curl -s -X POST \
                    -H "Content-Type: application/json" \
                    -d "${payload}" \
                    "http://localhost:8080/api/${endpoint}" 2>/dev/null || echo '{"error":"post_failed"}'
            else
                curl -s -X POST "http://localhost:8080/api/${endpoint}" 2>/dev/null || echo '{"error":"post_failed"}'
            fi
            ;;
    esac
}

# PUT request wrapper
call_thd_put() {
    local endpoint="$1"
    local payload="${2:-{}}"
    
    if [[ -n "${payload}" && "${payload}" != "{}" ]]; then
        curl -s -X PUT \
            -H "Content-Type: application/json" \
            -d "${payload}" \
            "http://localhost:8080/api/${endpoint}" 2>/dev/null || echo '{"error":"put_failed"}'
    else
        curl -s -X PUT "http://localhost:8080/api/${endpoint}" 2>/dev/null || echo '{"error":"put_failed"}'
    fi
}

# DELETE request wrapper
call_thd_delete() {
    local endpoint="$1"
    
    curl -s -X DELETE "http://localhost:8080/api/${endpoint}" 2>/dev/null || echo '{"error":"delete_failed"}'
}

# =========================================================================
# OBJECT MANAGEMENT API WRAPPERS
# =========================================================================

# Create object in session
api_create_object() {
    local session_id="$1"
    local object_data="$2"
    
    log "INFO" "Creating object in session ${session_id}"
    log "DEBUG" "Object data: ${object_data}"
    
    local response
    response=$(call_thd_api "sessions/${session_id}/objects" "POST" "${object_data}")
    
    # Validate object creation response
    if echo "${response}" | jq -e '.success == true or .id' >/dev/null 2>&1; then
        log "INFO" "Object created successfully"
        echo "${response}"
        return 0
    else
        log "ERROR" "Object creation failed: ${response}"
        return 1
    fi
}

# Get object from session
api_get_object() {
    local session_id="$1"
    local object_name="$2"
    
    log "DEBUG" "Getting object ${object_name} from session ${session_id}"
    
    call_thd_api "sessions/${session_id}/objects/${object_name}" "GET"
}

# Update object in session
api_update_object() {
    local session_id="$1"
    local object_name="$2"
    local update_data="$3"
    
    log "INFO" "Updating object ${object_name} in session ${session_id}"
    
    call_thd_api "sessions/${session_id}/objects/${object_name}" "PUT" "${update_data}"
}

# Delete object from session
api_delete_object() {
    local session_id="$1" 
    local object_name="$2"
    
    log "INFO" "Deleting object ${object_name} from session ${session_id}"
    
    call_thd_api "sessions/${session_id}/objects/${object_name}" "DELETE"
}

# List all objects in session
api_list_objects() {
    local session_id="$1"
    
    log "DEBUG" "Listing objects in session ${session_id}"
    
    call_thd_api "sessions/${session_id}/objects" "GET"
}

# =========================================================================
# CAMERA CONTROL API WRAPPERS  
# =========================================================================

# Set camera position
api_set_camera_position() {
    local session_id="$1"
    local x="$2"
    local y="$3" 
    local z="$4"
    
    local camera_data="{\"x\": ${x}, \"y\": ${y}, \"z\": ${z}}"
    
    log "INFO" "Setting camera position: (${x}, ${y}, ${z})"
    
    call_thd_api "sessions/${session_id}/camera/position" "PUT" "${camera_data}"
}

# Start camera orbit animation
api_start_camera_orbit() {
    local session_id="$1"
    local radius="${2:-10}"
    local speed="${3:-1.0}"
    
    local orbit_data="{\"radius\": ${radius}, \"speed\": ${speed}}"
    
    log "INFO" "Starting camera orbit: radius=${radius}, speed=${speed}"
    
    call_thd_api "sessions/${session_id}/camera/orbit" "POST" "${orbit_data}"
}

# =========================================================================
# BROWSER CONTROL API WRAPPERS
# =========================================================================

# Force browser refresh
api_force_browser_refresh() {
    log "INFO" "Forcing browser refresh"
    
    call_thd_api "browser/refresh" "POST"
}

# Set canvas control parameters
api_set_canvas() {
    local canvas_data="$1"
    
    log "INFO" "Setting canvas control parameters"
    
    call_thd_api "browser/canvas" "POST" "${canvas_data}"
}

# =========================================================================
# RESPONSE PARSING UTILITIES
# =========================================================================

# Extract field from JSON response
extract_json_field() {
    local json_response="$1"
    local field_path="$2"
    local default_value="${3:-}"
    
    echo "${json_response}" | jq -r "${field_path} // \"${default_value}\""
}

# Check if API response indicates success
is_api_success() {
    local response="$1"
    
    # Check various success indicators
    if echo "${response}" | jq -e '.success == true' >/dev/null 2>&1; then
        return 0
    fi
    
    if echo "${response}" | jq -e '.id' >/dev/null 2>&1; then
        return 0
    fi
    
    if echo "${response}" | jq -e '.session_id' >/dev/null 2>&1; then
        return 0
    fi
    
    # Check for error indicators
    if echo "${response}" | jq -e '.error' >/dev/null 2>&1; then
        return 1
    fi
    
    # Default to success if no clear indicators
    return 0
}

# Extract error message from API response
extract_error_message() {
    local response="$1"
    
    local error_msg
    error_msg=$(echo "${response}" | jq -r '.error // .message // "Unknown API error"')
    
    echo "${error_msg}"
}

# =========================================================================
# API HEALTH & MONITORING
# =========================================================================

# Check if THD server is responding
check_thd_server_health() {
    log "DEBUG" "Checking THD server health"
    
    local health_response
    health_response=$(curl -s --connect-timeout 5 "http://localhost:8080/api/sessions" 2>/dev/null) || {
        log "ERROR" "THD server not responding"
        return 1
    }
    
    if echo "${health_response}" | jq . >/dev/null 2>&1; then
        log "DEBUG" "THD server responding with valid JSON"
        return 0
    else
        log "WARN" "THD server responding but invalid JSON: ${health_response}"
        return 1
    fi
}

# Wait for THD server to become available
wait_for_thd_server() {
    local timeout="${1:-30}"
    local interval="${2:-2}"
    
    log "INFO" "Waiting for THD server to become available (timeout: ${timeout}s)"
    
    local elapsed=0
    while [[ ${elapsed} -lt ${timeout} ]]; do
        if check_thd_server_health; then
            log "INFO" "THD server is available"
            return 0
        fi
        
        log "DEBUG" "THD server not ready, waiting ${interval}s... (${elapsed}/${timeout}s)"
        sleep "${interval}"
        elapsed=$((elapsed + interval))
    done
    
    log "ERROR" "THD server did not become available within ${timeout}s"
    return 1
}

# Professional API call with retry logic
call_thd_api_with_retry() {
    local endpoint="$1"
    local method="${2:-GET}"
    local payload="${3:-}"
    local max_retries="${4:-3}"
    local retry_delay="${5:-2}"
    
    local attempt=1
    local response
    
    while [[ ${attempt} -le ${max_retries} ]]; do
        log "DEBUG" "API attempt ${attempt}/${max_retries}: ${method} ${endpoint}"
        
        if response=$(call_thd_api "${endpoint}" "${method}" "${payload}"); then
            echo "${response}"
            return 0
        fi
        
        if [[ ${attempt} -eq ${max_retries} ]]; then
            log "ERROR" "API call failed after ${max_retries} attempts: ${method} ${endpoint}"
            return 1
        fi
        
        log "WARN" "API attempt ${attempt} failed, retrying in ${retry_delay}s..."
        sleep "${retry_delay}"
        attempt=$((attempt + 1))
    done
}