#!/bin/bash
# ===================================================================
# HD1 Channel One - Default Collaborative Environment
# ===================================================================
#
# Purpose: Main collaboration channel for general holodeck activities
# Environment: Earth Surface physics with standard lighting
# Capacity: Unlimited clients
# Default Objects: Basic scene setup with ground plane and lighting
#
# ===================================================================

# Load HD1 core library
source "${HD1_ROOT}/lib/hd1lib.sh" 2>/dev/null || {
    echo "ERROR: HD1 core library not found"
    exit 1
}

# Channel configuration
CHANNEL_ID="channel_one"
CHANNEL_NAME="Channel One - Main Collaboration"
DEFAULT_ENVIRONMENT="earth-surface"
AUTO_CREATE_SCENE=true

# Channel initialization
echo "CHANNEL: Initializing ${CHANNEL_NAME}"
echo "ID: ${CHANNEL_ID}"
echo "ENV: ${DEFAULT_ENVIRONMENT}"

# Set default session for this channel
HD1_SESSION_ID="${HD1_SESSION_ID:-${CHANNEL_ID}}"

# Initialize channel environment
hd1::init_channel() {
    echo "INIT: Setting up ${CHANNEL_NAME}"
    
    # Apply default environment
    echo "ENV: Applying ${DEFAULT_ENVIRONMENT} physics"
    # hd1::apply_environment "${DEFAULT_ENVIRONMENT}"
    
    # Create default scene objects if enabled
    if [[ "${AUTO_CREATE_SCENE}" == "true" ]]; then
        echo "SCENE: Creating default objects"
        
        # Ground plane
        hd1::create_object "ground-plane" "plane" 0 0 0
        echo "OBJECT: Ground plane created"
        
        # Main lighting
        hd1::create_object "main-light" "light" 0 5 5
        echo "LIGHT: Main illumination created"
        
        # Welcome beacon
        hd1::create_object "welcome-beacon" "sphere" 0 2 0
        hd1::update_object "welcome-beacon" "color" "#00ff00"
        echo "BEACON: Welcome marker created"
    fi
    
    echo "READY: ${CHANNEL_NAME} initialized successfully"
}

# Channel cleanup
hd1::cleanup_channel() {
    echo "CLEANUP: Shutting down ${CHANNEL_NAME}"
    # Add any cleanup logic here
}

# Export channel functions
export CHANNEL_ID CHANNEL_NAME DEFAULT_ENVIRONMENT
export -f hd1::init_channel hd1::cleanup_channel

# Auto-initialize if sourced directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    hd1::init_channel
fi

echo "CHANNEL: ${CHANNEL_NAME} script loaded"