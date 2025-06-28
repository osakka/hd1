#!/bin/bash

# =========================================================================
# GLIBSH Core: Validation - Professional Input Validation & Error Prevention
# =========================================================================
#
# Comprehensive validation system providing:
# - Coordinate bounds enforcement (THD: [-12, +12])
# - Parameter type checking and sanitization
# - Professional error messages with guidance
# - Input sanitization and security validation
# =========================================================================

set -euo pipefail

# =========================================================================
# THD COORDINATE SYSTEM VALIDATION
# =========================================================================

# THD coordinate system bounds: [-12, +12] on all axes
readonly THD_MIN_COORD=-12
readonly THD_MAX_COORD=12

# Validate 3D coordinate within THD bounds
validate_coordinate() {
    local coord="$1"
    local axis_name="${2:-coordinate}"
    
    # Check if coordinate is a valid number
    if ! [[ "${coord}" =~ ^-?[0-9]+(\.[0-9]+)?$ ]]; then
        error_exit "Invalid ${axis_name}: '${coord}' is not a valid number"
    fi
    
    # Check bounds
    if (( $(echo "${coord} < ${THD_MIN_COORD}" | bc -l) )); then
        error_exit "${axis_name} ${coord} below minimum bound ${THD_MIN_COORD}"
    fi
    
    if (( $(echo "${coord} > ${THD_MAX_COORD}" | bc -l) )); then
        error_exit "${axis_name} ${coord} above maximum bound ${THD_MAX_COORD}"
    fi
    
    log "DEBUG" "Coordinate validated: ${axis_name}=${coord}"
    return 0
}

# Validate 3D position string (format: "x,y,z")
validate_position() {
    local position="$1"
    
    # Split position into coordinates
    local IFS=','
    read -ra coords <<< "${position}"
    
    if [[ ${#coords[@]} -ne 3 ]]; then
        error_exit "Position must be in format 'x,y,z', got: '${position}'"
    fi
    
    # Validate each coordinate
    validate_coordinate "${coords[0]}" "X coordinate"
    validate_coordinate "${coords[1]}" "Y coordinate" 
    validate_coordinate "${coords[2]}" "Z coordinate"
    
    log "DEBUG" "Position validated: ${position}"
    return 0
}

# Parse and validate position string, return as separate variables
parse_position() {
    local position="$1"
    
    validate_position "${position}"
    
    local IFS=','
    read -ra coords <<< "${position}"
    
    # Export coordinates for use by caller
    PARSED_X="${coords[0]}"
    PARSED_Y="${coords[1]}"
    PARSED_Z="${coords[2]}"
    
    log "DEBUG" "Position parsed: X=${PARSED_X}, Y=${PARSED_Y}, Z=${PARSED_Z}"
}

# =========================================================================
# OBJECT PROPERTY VALIDATION
# =========================================================================

# Validate object type
validate_object_type() {
    local object_type="$1"
    
    local valid_types=("cube" "sphere" "mesh" "plane" "cylinder" "cone")
    
    for valid_type in "${valid_types[@]}"; do
        if [[ "${object_type}" == "${valid_type}" ]]; then
            log "DEBUG" "Object type validated: ${object_type}"
            return 0
        fi
    done
    
    error_exit "Invalid object type: '${object_type}'. Valid types: ${valid_types[*]}"
}

# Validate color specification (name or hex)
validate_color() {
    local color="$1"
    
    # Predefined color names
    local valid_colors=("red" "green" "blue" "yellow" "cyan" "magenta" "white" "black" "orange" "purple" "pink" "brown" "gray" "grey")
    
    # Check if it's a named color
    for valid_color in "${valid_colors[@]}"; do
        if [[ "${color,,}" == "${valid_color}" ]]; then
            log "DEBUG" "Color validated (name): ${color}"
            return 0
        fi
    done
    
    # Check if it's a hex color (format: #RRGGBB or #RGB)
    if [[ "${color}" =~ ^#[0-9A-Fa-f]{6}$ ]] || [[ "${color}" =~ ^#[0-9A-Fa-f]{3}$ ]]; then
        log "DEBUG" "Color validated (hex): ${color}"
        return 0
    fi
    
    error_exit "Invalid color: '${color}'. Use color name (${valid_colors[*]}) or hex format (#RRGGBB)"
}

# Validate scale factor
validate_scale() {
    local scale="$1"
    
    # Check if scale is a positive number
    if ! [[ "${scale}" =~ ^[0-9]+(\.[0-9]+)?$ ]]; then
        error_exit "Invalid scale: '${scale}' must be a positive number"
    fi
    
    # Check reasonable bounds (0.1 to 50.0)
    if (( $(echo "${scale} < 0.1" | bc -l) )); then
        error_exit "Scale ${scale} too small (minimum: 0.1)"
    fi
    
    if (( $(echo "${scale} > 50.0" | bc -l) )); then
        error_exit "Scale ${scale} too large (maximum: 50.0)"
    fi
    
    log "DEBUG" "Scale validated: ${scale}"
    return 0
}

# Validate object name (alphanumeric, hyphens, underscores)
validate_object_name() {
    local name="$1"
    
    # Check format
    if ! [[ "${name}" =~ ^[a-zA-Z0-9_-]+$ ]]; then
        error_exit "Invalid object name: '${name}'. Use alphanumeric characters, hyphens, and underscores only"
    fi
    
    # Check length
    if [[ ${#name} -lt 1 ]] || [[ ${#name} -gt 50 ]]; then
        error_exit "Object name length must be 1-50 characters, got ${#name}: '${name}'"
    fi
    
    log "DEBUG" "Object name validated: ${name}"
    return 0
}

# =========================================================================
# SHAPE PARAMETER VALIDATION
# =========================================================================

# Validate grid size specification (format: "NxM" or "N")
validate_grid_size() {
    local size_spec="$1"
    
    local width height
    
    if [[ "${size_spec}" =~ ^([0-9]+)x([0-9]+)$ ]]; then
        # Format: "NxM"
        width="${BASH_REMATCH[1]}"
        height="${BASH_REMATCH[2]}"
    elif [[ "${size_spec}" =~ ^[0-9]+$ ]]; then
        # Format: "N" (square grid)
        width="${size_spec}"
        height="${size_spec}"
    else
        error_exit "Invalid grid size: '${size_spec}'. Use format 'NxM' or 'N'"
    fi
    
    # Validate dimensions
    if [[ ${width} -lt 1 ]] || [[ ${width} -gt 25 ]]; then
        error_exit "Grid width ${width} out of range (1-25)"
    fi
    
    if [[ ${height} -lt 1 ]] || [[ ${height} -gt 25 ]]; then
        error_exit "Grid height ${height} out of range (1-25)"
    fi
    
    # Export for use by caller
    PARSED_GRID_WIDTH="${width}"
    PARSED_GRID_HEIGHT="${height}"
    
    log "DEBUG" "Grid size validated: ${width}x${height}"
    return 0
}

# Validate spacing parameter
validate_spacing() {
    local spacing="$1"
    
    # Check if spacing is a positive number
    if ! [[ "${spacing}" =~ ^[0-9]+(\.[0-9]+)?$ ]]; then
        error_exit "Invalid spacing: '${spacing}' must be a positive number"
    fi
    
    # Check reasonable bounds (0.1 to 10.0)
    if (( $(echo "${spacing} < 0.1" | bc -l) )); then
        error_exit "Spacing ${spacing} too small (minimum: 0.1)"
    fi
    
    if (( $(echo "${spacing} > 10.0" | bc -l) )); then
        error_exit "Spacing ${spacing} too large (maximum: 10.0)"
    fi
    
    log "DEBUG" "Spacing validated: ${spacing}"
    return 0
}

# Validate count parameter (for objects, iterations, etc.)
validate_count() {
    local count="$1"
    local max_count="${2:-1000}"
    
    # Check if count is a positive integer
    if ! [[ "${count}" =~ ^[0-9]+$ ]]; then
        error_exit "Invalid count: '${count}' must be a positive integer"
    fi
    
    # Check bounds
    if [[ ${count} -lt 1 ]]; then
        error_exit "Count must be at least 1, got: ${count}"
    fi
    
    if [[ ${count} -gt ${max_count} ]]; then
        error_exit "Count ${count} exceeds maximum ${max_count}"
    fi
    
    log "DEBUG" "Count validated: ${count} (max: ${max_count})"
    return 0
}

# =========================================================================
# ANIMATION & CAMERA VALIDATION
# =========================================================================

# Validate animation pattern type
validate_animation_pattern() {
    local pattern="$1"
    
    local valid_patterns=("rotate" "orbit" "bounce" "scale" "fade" "spiral")
    
    for valid_pattern in "${valid_patterns[@]}"; do
        if [[ "${pattern}" == "${valid_pattern}" ]]; then
            log "DEBUG" "Animation pattern validated: ${pattern}"
            return 0
        fi
    done
    
    error_exit "Invalid animation pattern: '${pattern}'. Valid patterns: ${valid_patterns[*]}"
}

# Validate camera preset
validate_camera_preset() {
    local preset="$1"
    
    local valid_presets=("overview" "close" "orbital" "cinematic" "top" "side" "front" "isometric")
    
    for valid_preset in "${valid_presets[@]}"; do
        if [[ "${preset}" == "${valid_preset}" ]]; then
            log "DEBUG" "Camera preset validated: ${preset}"
            return 0
        fi
    done
    
    error_exit "Invalid camera preset: '${preset}'. Valid presets: ${valid_presets[*]}"
}

# Validate duration (in seconds)
validate_duration() {
    local duration="$1"
    
    # Check if duration is a positive number
    if ! [[ "${duration}" =~ ^[0-9]+(\.[0-9]+)?$ ]]; then
        error_exit "Invalid duration: '${duration}' must be a positive number (seconds)"
    fi
    
    # Check reasonable bounds (0.1 to 3600 seconds)
    if (( $(echo "${duration} < 0.1" | bc -l) )); then
        error_exit "Duration ${duration} too short (minimum: 0.1 seconds)"
    fi
    
    if (( $(echo "${duration} > 3600" | bc -l) )); then
        error_exit "Duration ${duration} too long (maximum: 3600 seconds)"
    fi
    
    log "DEBUG" "Duration validated: ${duration} seconds"
    return 0
}

# =========================================================================
# FILE & PATH VALIDATION
# =========================================================================

# Validate file path exists and is readable
validate_file_path() {
    local file_path="$1"
    local description="${2:-file}"
    
    if [[ ! -f "${file_path}" ]]; then
        error_exit "${description} not found: ${file_path}"
    fi
    
    if [[ ! -r "${file_path}" ]]; then
        error_exit "${description} not readable: ${file_path}"
    fi
    
    log "DEBUG" "File path validated: ${file_path}"
    return 0
}

# Validate export format
validate_export_format() {
    local format="$1"
    
    local valid_formats=("json" "yaml" "csv" "xml")
    
    for valid_format in "${valid_formats[@]}"; do
        if [[ "${format}" == "${valid_format}" ]]; then
            log "DEBUG" "Export format validated: ${format}"
            return 0
        fi
    done
    
    error_exit "Invalid export format: '${format}'. Valid formats: ${valid_formats[*]}"
}

# =========================================================================
# SECURITY & SANITIZATION
# =========================================================================

# Sanitize string for safe usage (remove dangerous characters)
sanitize_string() {
    local input="$1"
    
    # Remove potentially dangerous characters
    local sanitized
    sanitized=$(echo "${input}" | tr -d '`$(){}[]|&;<>*?~^')
    
    log "DEBUG" "String sanitized: '${input}' -> '${sanitized}'"
    echo "${sanitized}"
}

# Validate session ID format
validate_session_id() {
    local session_id="$1"
    
    # THD session IDs should be alphanumeric with hyphens
    if ! [[ "${session_id}" =~ ^[a-zA-Z0-9-]+$ ]]; then
        error_exit "Invalid session ID format: '${session_id}'"
    fi
    
    # Check reasonable length
    if [[ ${#session_id} -lt 8 ]] || [[ ${#session_id} -gt 64 ]]; then
        error_exit "Session ID length out of range (8-64 chars): '${session_id}'"
    fi
    
    log "DEBUG" "Session ID validated: ${session_id}"
    return 0
}

# =========================================================================
# BATCH VALIDATION FUNCTIONS
# =========================================================================

# Validate complete object specification
validate_object_spec() {
    local position="$1"
    local object_type="$2"
    local color="${3:-white}"
    local scale="${4:-1.0}"
    local name="${5:-auto}"
    
    log "DEBUG" "Validating complete object specification"
    
    validate_position "${position}"
    validate_object_type "${object_type}"
    validate_color "${color}"
    validate_scale "${scale}"
    
    if [[ "${name}" != "auto" ]]; then
        validate_object_name "${name}"
    fi
    
    log "DEBUG" "Object specification validated successfully"
    return 0
}

# Professional validation summary
show_validation_summary() {
    cat << 'EOF'
GLIBSH Validation Standards
==========================

Coordinates: [-12, +12] on all axes (THD coordinate system)
Position:    "x,y,z" format with valid coordinates
Colors:      Named colors or hex format (#RRGGBB)
Scale:       0.1 to 50.0
Names:       Alphanumeric, hyphens, underscores (1-50 chars)
Count:       1 to 1000 (configurable maximum)
Duration:    0.1 to 3600 seconds
Grid Size:   1x1 to 25x25
Spacing:     0.1 to 10.0

All inputs are validated for security and THD compatibility.
EOF
}