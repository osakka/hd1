#!/bin/bash

# =========================================================================
# GLIBSH Core: World Management - Professional 3D World Orchestration
# =========================================================================
#
# World coordinate system management providing:
# - THD coordinate system initialization (25×25×25 grid)
# - World state management and validation
# - Professional world reset and cleanup procedures
# - 3D space organization and optimization
# =========================================================================

set -euo pipefail

# =========================================================================
# THD WORLD CONSTANTS & CONFIGURATION
# =========================================================================

# THD World Specification (professional defaults)
readonly THD_WORLD_SIZE=25
readonly THD_WORLD_MIN_BOUND=-12
readonly THD_WORLD_MAX_BOUND=12
readonly THD_DEFAULT_TRANSPARENCY=0.1
readonly THD_DEFAULT_GRID_SPACING=1.0

# Default camera positioning for optimal viewing
readonly THD_DEFAULT_CAMERA_X=10
readonly THD_DEFAULT_CAMERA_Y=10
readonly THD_DEFAULT_CAMERA_Z=10

# =========================================================================
# WORLD INITIALIZATION & MANAGEMENT
# =========================================================================

# Initialize THD world with professional standards
initialize_world() {
    local session_id="$1"
    local size="${2:-${THD_WORLD_SIZE}}"
    local transparency="${3:-${THD_DEFAULT_TRANSPARENCY}}"
    local camera_x="${4:-${THD_DEFAULT_CAMERA_X}}"
    local camera_y="${5:-${THD_DEFAULT_CAMERA_Y}}"
    local camera_z="${6:-${THD_DEFAULT_CAMERA_Z}}"
    
    log "INFO" "Initializing THD world for session: ${session_id}"
    log "DEBUG" "World parameters: size=${size}, transparency=${transparency}, camera=(${camera_x},${camera_y},${camera_z})"
    
    # Validate world parameters
    validate_world_parameters "${size}" "${transparency}" "${camera_x}" "${camera_y}" "${camera_z}"
    
    # Create world initialization payload
    local world_payload
    world_payload=$(create_world_payload "${size}" "${transparency}" "${camera_x}" "${camera_y}" "${camera_z}")
    
    # Initialize world via API
    local response
    response=$(call_thd_api "sessions/${session_id}/world" "POST" "${world_payload}") || {
        error_exit "Failed to initialize world for session ${session_id}"
    }
    
    # Validate successful initialization
    if is_api_success "${response}"; then
        log "INFO" "World initialized successfully for session: ${session_id}"
        
        # Set optimal camera position
        set_optimal_camera_position "${session_id}" "${camera_x}" "${camera_y}" "${camera_z}"
        
        return 0
    else
        local error_msg
        error_msg=$(extract_error_message "${response}")
        error_exit "World initialization failed: ${error_msg}"
    fi
}

# Create world initialization payload with professional structure
create_world_payload() {
    local size="$1"
    local transparency="$2"
    local camera_x="$3"
    local camera_y="$4"
    local camera_z="$5"
    
    # Professional JSON payload construction
    jq -n \
        --arg size "${size}" \
        --arg transparency "${transparency}" \
        --arg grid_spacing "${THD_DEFAULT_GRID_SPACING}" \
        --arg camera_x "${camera_x}" \
        --arg camera_y "${camera_y}" \
        --arg camera_z "${camera_z}" \
        '{
            size: ($size | tonumber),
            transparency: ($transparency | tonumber),
            grid_spacing: ($grid_spacing | tonumber),
            camera_x: ($camera_x | tonumber),
            camera_y: ($camera_y | tonumber),
            camera_z: ($camera_z | tonumber),
            bounds: {
                min: -12,
                max: 12
            },
            coordinate_system: "thd_standard",
            grid_enabled: true,
            axis_helpers: true
        }'
}

# Validate world initialization parameters
validate_world_parameters() {
    local size="$1"
    local transparency="$2"
    local camera_x="$3"
    local camera_y="$4" 
    local camera_z="$5"
    
    # Validate size
    if ! [[ "${size}" =~ ^[0-9]+$ ]] || [[ ${size} -lt 5 ]] || [[ ${size} -gt 50 ]]; then
        error_exit "Invalid world size: ${size}. Must be integer between 5-50"
    fi
    
    # Validate transparency
    if ! [[ "${transparency}" =~ ^[0-9]+(\.[0-9]+)?$ ]] || (( $(echo "${transparency} < 0" | bc -l) )) || (( $(echo "${transparency} > 1" | bc -l) )); then
        error_exit "Invalid transparency: ${transparency}. Must be between 0.0-1.0"
    fi
    
    # Validate camera coordinates
    validate_coordinate "${camera_x}" "Camera X"
    validate_coordinate "${camera_y}" "Camera Y"
    validate_coordinate "${camera_z}" "Camera Z"
    
    log "DEBUG" "World parameters validated successfully"
}

# =========================================================================
# WORLD STATE MANAGEMENT
# =========================================================================

# Get current world state from session
get_world_state() {
    local session_id="$1"
    
    log "DEBUG" "Retrieving world state for session: ${session_id}"
    
    local response
    response=$(call_thd_api "sessions/${session_id}/world" "GET") || {
        log "ERROR" "Failed to retrieve world state for session ${session_id}"
        return 1
    }
    
    echo "${response}"
}

# Validate world exists and is properly initialized
validate_world_state() {
    local session_id="$1"
    
    log "DEBUG" "Validating world state for session: ${session_id}"
    
    local world_state
    world_state=$(get_world_state "${session_id}") || {
        log "WARN" "No world state found for session ${session_id}"
        return 1
    }
    
    # Check for required world properties
    local bounds size
    bounds=$(echo "${world_state}" | jq -r '.bounds // empty')
    size=$(echo "${world_state}" | jq -r '.size // empty')
    
    if [[ -z "${bounds}" ]] || [[ -z "${size}" ]]; then
        log "WARN" "World state incomplete for session ${session_id}"
        return 1
    fi
    
    log "DEBUG" "World state validated for session: ${session_id}"
    return 0
}

# Reset world to clean state
reset_world() {
    local session_id="$1"
    local preserve_objects="${2:-false}"
    
    log "INFO" "Resetting world for session: ${session_id} (preserve_objects=${preserve_objects})"
    
    if [[ "${preserve_objects}" != "true" ]]; then
        # Clear all objects first
        clear_all_objects "${session_id}"
    fi
    
    # Reinitialize world with default parameters
    initialize_world "${session_id}"
    
    log "INFO" "World reset completed for session: ${session_id}"
}

# Clear all objects from world
clear_all_objects() {
    local session_id="$1"
    
    log "INFO" "Clearing all objects from session: ${session_id}"
    
    # Get list of all objects
    local objects_response
    objects_response=$(api_list_objects "${session_id}") || {
        log "WARN" "Failed to list objects for session ${session_id}"
        return 1
    }
    
    # Extract object names and delete each one
    local object_names
    object_names=$(echo "${objects_response}" | jq -r '.objects[]?.name // empty' 2>/dev/null)
    
    if [[ -n "${object_names}" ]]; then
        while IFS= read -r object_name; do
            if [[ -n "${object_name}" ]]; then
                log "DEBUG" "Deleting object: ${object_name}"
                api_delete_object "${session_id}" "${object_name}" || {
                    log "WARN" "Failed to delete object: ${object_name}"
                }
            fi
        done <<< "${object_names}"
        
        log "INFO" "All objects cleared from session: ${session_id}"
    else
        log "DEBUG" "No objects found to clear in session: ${session_id}"
    fi
}

# =========================================================================
# CAMERA & VIEWING OPTIMIZATION
# =========================================================================

# Set optimal camera position for world viewing
set_optimal_camera_position() {
    local session_id="$1"
    local camera_x="${2:-${THD_DEFAULT_CAMERA_X}}"
    local camera_y="${3:-${THD_DEFAULT_CAMERA_Y}}"
    local camera_z="${4:-${THD_DEFAULT_CAMERA_Z}}"
    
    log "INFO" "Setting optimal camera position: (${camera_x}, ${camera_y}, ${camera_z})"
    
    api_set_camera_position "${session_id}" "${camera_x}" "${camera_y}" "${camera_z}" || {
        log "WARN" "Failed to set camera position for session ${session_id}"
        return 1
    }
    
    log "DEBUG" "Camera position set successfully"
    return 0
}

# Calculate optimal camera position based on world content
calculate_optimal_camera() {
    local session_id="$1"
    
    log "DEBUG" "Calculating optimal camera position for session: ${session_id}"
    
    # Get all objects in the world
    local objects_response
    objects_response=$(api_list_objects "${session_id}") || {
        # If no objects, use default position
        echo "${THD_DEFAULT_CAMERA_X} ${THD_DEFAULT_CAMERA_Y} ${THD_DEFAULT_CAMERA_Z}"
        return 0
    }
    
    # Calculate bounding box of all objects
    local min_x=${THD_WORLD_MAX_BOUND} max_x=${THD_WORLD_MIN_BOUND}
    local min_y=${THD_WORLD_MAX_BOUND} max_y=${THD_WORLD_MIN_BOUND}
    local min_z=${THD_WORLD_MAX_BOUND} max_z=${THD_WORLD_MIN_BOUND}
    
    # Extract object positions and calculate bounds
    while IFS= read -r object_data; do
        if [[ -n "${object_data}" ]]; then
            local x y z
            x=$(echo "${object_data}" | jq -r '.x // .transform.position.x // 0')
            y=$(echo "${object_data}" | jq -r '.y // .transform.position.y // 0')
            z=$(echo "${object_data}" | jq -r '.z // .transform.position.z // 0')
            
            # Update bounds
            if (( $(echo "${x} < ${min_x}" | bc -l) )); then min_x="${x}"; fi
            if (( $(echo "${x} > ${max_x}" | bc -l) )); then max_x="${x}"; fi
            if (( $(echo "${y} < ${min_y}" | bc -l) )); then min_y="${y}"; fi
            if (( $(echo "${y} > ${max_y}" | bc -l) )); then max_y="${y}"; fi
            if (( $(echo "${z} < ${min_z}" | bc -l) )); then min_z="${z}"; fi
            if (( $(echo "${z} > ${max_z}" | bc -l) )); then max_z="${z}"; fi
        fi
    done <<< "$(echo "${objects_response}" | jq -c '.objects[]? // empty' 2>/dev/null)"
    
    # Calculate center point and optimal distance
    local center_x center_y center_z
    center_x=$(echo "scale=2; (${max_x} + ${min_x}) / 2" | bc -l)
    center_y=$(echo "scale=2; (${max_y} + ${min_y}) / 2" | bc -l)
    center_z=$(echo "scale=2; (${max_z} + ${min_z}) / 2" | bc -l)
    
    # Calculate optimal viewing distance (1.5x the maximum extent)
    local extent_x extent_y extent_z max_extent
    extent_x=$(echo "scale=2; ${max_x} - ${min_x}" | bc -l)
    extent_y=$(echo "scale=2; ${max_y} - ${min_y}" | bc -l) 
    extent_z=$(echo "scale=2; ${max_z} - ${min_z}" | bc -l)
    
    # Find maximum extent
    max_extent="${extent_x}"
    if (( $(echo "${extent_y} > ${max_extent}" | bc -l) )); then max_extent="${extent_y}"; fi
    if (( $(echo "${extent_z} > ${max_extent}" | bc -l) )); then max_extent="${extent_z}"; fi
    
    # Calculate camera position (45-degree isometric view)
    local distance
    distance=$(echo "scale=2; ${max_extent} * 1.5 + 5" | bc -l)
    
    local camera_x camera_y camera_z
    camera_x=$(echo "scale=2; ${center_x} + ${distance}" | bc -l)
    camera_y=$(echo "scale=2; ${center_y} + ${distance}" | bc -l)
    camera_z=$(echo "scale=2; ${center_z} + ${distance}" | bc -l)
    
    # Ensure camera position is within world bounds
    if (( $(echo "${camera_x} > ${THD_WORLD_MAX_BOUND}" | bc -l) )); then camera_x="${THD_WORLD_MAX_BOUND}"; fi
    if (( $(echo "${camera_y} > ${THD_WORLD_MAX_BOUND}" | bc -l) )); then camera_y="${THD_WORLD_MAX_BOUND}"; fi
    if (( $(echo "${camera_z} > ${THD_WORLD_MAX_BOUND}" | bc -l) )); then camera_z="${THD_WORLD_MAX_BOUND}"; fi
    
    log "DEBUG" "Calculated optimal camera: (${camera_x}, ${camera_y}, ${camera_z})"
    echo "${camera_x} ${camera_y} ${camera_z}"
}

# =========================================================================
# WORLD INFORMATION & MONITORING
# =========================================================================

# Display comprehensive world information
show_world_info() {
    local session_id="$1"
    
    log "INFO" "Retrieving world information for session: ${session_id}"
    
    # Get world state
    local world_state
    world_state=$(get_world_state "${session_id}") || {
        echo "❌ World not initialized or not accessible"
        return 1
    }
    
    # Get object count
    local objects_response object_count
    objects_response=$(api_list_objects "${session_id}" 2>/dev/null) || objects_response='{"objects":[]}'
    object_count=$(echo "${objects_response}" | jq '.objects | length' 2>/dev/null || echo "0")
    
    # Display professional world summary
    echo "THD World Information"
    echo "===================="
    echo "Session ID: ${session_id}"
    echo "World Size: $(echo "${world_state}" | jq -r '.size // "Unknown"')"
    echo "Bounds: [${THD_WORLD_MIN_BOUND}, ${THD_WORLD_MAX_BOUND}] on all axes"
    echo "Grid Spacing: $(echo "${world_state}" | jq -r '.grid_spacing // "1.0"')"
    echo "Transparency: $(echo "${world_state}" | jq -r '.transparency // "0.1"')"
    echo "Objects: ${object_count}"
    echo "Coordinate System: THD Standard (25×25×25)"
    echo "Status: $(if validate_world_state "${session_id}" >/dev/null 2>&1; then echo "✅ Healthy"; else echo "⚠️  Needs attention"; fi)"
}

# Monitor world performance and health
monitor_world_health() {
    local session_id="$1"
    local duration="${2:-60}"
    
    log "INFO" "Monitoring world health for ${duration} seconds"
    
    local start_time end_time
    start_time=$(date +%s)
    end_time=$((start_time + duration))
    
    while [[ $(date +%s) -lt ${end_time} ]]; do
        if validate_world_state "${session_id}" >/dev/null 2>&1; then
            local object_count
            object_count=$(api_list_objects "${session_id}" | jq '.objects | length' 2>/dev/null || echo "0")
            
            local elapsed=$(($(date +%s) - start_time))
            echo "[${elapsed}s] World Health: ✅ OK | Objects: ${object_count}"
        else
            echo "[${elapsed}s] World Health: ❌ UNHEALTHY"
        fi
        
        sleep 5
    done
    
    log "INFO" "World health monitoring completed"
}

# =========================================================================
# WORLD TEMPLATES & PRESETS
# =========================================================================

# Initialize world with preset configuration
initialize_world_preset() {
    local session_id="$1"
    local preset="${2:-standard}"
    
    log "INFO" "Initializing world with preset: ${preset}"
    
    case "${preset}" in
        "minimal")
            initialize_world "${session_id}" 15 0.05 8 8 8
            ;;
        "standard")
            initialize_world "${session_id}" 25 0.1 10 10 10
            ;;
        "expanded")
            initialize_world "${session_id}" 35 0.15 12 12 12
            ;;
        "performance")
            initialize_world "${session_id}" 20 0.2 6 6 6
            ;;
        *)
            error_exit "Unknown world preset: ${preset}. Available: minimal, standard, expanded, performance"
            ;;
    esac
    
    log "INFO" "World preset '${preset}' initialized successfully"
}