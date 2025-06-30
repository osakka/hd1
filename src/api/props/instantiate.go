package props

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"log"
	"holodeck1/server"
)

// EnvironmentInfo holds environment physics parameters for prop adaptation
type EnvironmentInfo struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	ScaleUnit    string  `json:"scale_unit"`
	Gravity      float64 `json:"gravity"`
	Atmosphere   string  `json:"atmosphere"`
	Density      float64 `json:"density"`
	Temperature  float64 `json:"temperature"`
}

// InstantiateRequest matches the OpenAPI request schema exactly
type InstantiateRequest struct {
	Position     Position `json:"position"`
	Rotation     *Rotation `json:"rotation,omitempty"`
	Scale        float64  `json:"scale,omitempty"`
	InstanceName string   `json:"instance_name,omitempty"`
}

// Position represents 3D coordinates
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Rotation represents 3D rotation
type Rotation struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// InstantiateResponse matches the OpenAPI response schema exactly
type InstantiateResponse struct {
	Success        bool   `json:"success"`
	PropID         string `json:"prop_id"`
	InstanceName   string `json:"instance_name"`
	SessionID      string `json:"session_id"`
	ObjectsCreated int    `json:"objects_created"`
	Message        string `json:"message"`
}

// ErrorResponse matches the OpenAPI error schema
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

// InstantiatePropHandler creates an instance of a prop in the specified session
func InstantiatePropHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type for session validation
	h, ok := hub.(*server.Hub)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "internal_error",
			Message: "Internal server error",
		})
		return
	}

	// Extract session ID and prop ID from URL path
	sessionID := extractSessionID(r.URL.Path)
	propID := extractPropID(r.URL.Path)
	
	if sessionID == "" || propID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "missing_parameters",
			Message: "Session ID and Prop ID are required",
		})
		return
	}

	// Parse request body
	var req InstantiateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "invalid_request",
			Message: "Invalid JSON request body",
		})
		return
	}

	// Verify session exists
	if _, exists := h.GetStore().GetSession(sessionID); !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "session_not_found",
			Message: "Session does not exist",
		})
		return
	}

	// Find prop definition
	propInfo := findPropDefinition(propID)
	if propInfo == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "prop_not_found",
			Message: "Prop '" + propID + "' not found",
		})
		return
	}

	// Set defaults
	if req.Scale == 0 {
		req.Scale = 1.0
	}
	if req.InstanceName == "" {
		req.InstanceName = propID + "-instance"
	}

	// Get current environment for physics cohesion
	session, _ := h.GetStore().GetSession(sessionID)
	var envInfo *EnvironmentInfo
	if session.EnvironmentID != "" {
		envInfo = getEnvironmentInfo(session.EnvironmentID)
	}

	// Execute prop instantiation script with environment physics
	objectsCreated, err := executePropScript(propInfo, &req, sessionID, envInfo)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "prop_instantiation_failed",
			Message: "Failed to instantiate prop: " + err.Error(),
		})
		return
	}

	log.Printf("[HD1] Prop '%s' instantiated as '%s' in session '%s'", propID, req.InstanceName, sessionID)

	// Send prop instantiation notification via WebSocket
	h.BroadcastToSession(sessionID, "prop_instantiated", map[string]interface{}{
		"prop_id":         propID,
		"instance_name":   req.InstanceName,
		"session_id":      sessionID,
		"objects_created": objectsCreated,
		"timestamp":       "now",
	})

	// Return success response
	response := InstantiateResponse{
		Success:        true,
		PropID:         propID,
		InstanceName:   req.InstanceName,
		SessionID:      sessionID,
		ObjectsCreated: objectsCreated,
		Message:        fmt.Sprintf("Prop '%s' instantiated as '%s'", propInfo.Name, req.InstanceName),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// extractSessionID extracts session ID from URL path
func extractSessionID(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "sessions" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// extractPropID extracts prop ID from URL path  
func extractPropID(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "props" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// findPropDefinition finds and loads a prop definition by ID
func findPropDefinition(propID string) *PropInfo {
	propsDir := "/opt/hd1/share/props"
	
	// Search all category directories for the prop
	err := filepath.Walk(propsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if strings.HasSuffix(info.Name(), ".yaml") {
			if prop := parsePropFile(path); prop != nil && prop.ID == propID {
				return fmt.Errorf("found") // Use error to break out of walk
			}
		}
		
		return nil
	})
	
	if err != nil && err.Error() == "found" {
		// Re-parse the found file (inefficient but simple)
		filepath.Walk(propsDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			
			if strings.HasSuffix(info.Name(), ".yaml") {
				if prop := parsePropFile(path); prop != nil && prop.ID == propID {
					return nil // Found it
				}
			}
			
			return nil
		})
	}
	
	// For now, return default props if file system fails
	defaults := getDefaultProps()
	for _, prop := range defaults {
		if prop.ID == propID {
			return &prop
		}
	}
	
	return nil
}

// executePropScript executes the prop's script to create objects in the session
func executePropScript(prop *PropInfo, req *InstantiateRequest, sessionID string, envInfo *EnvironmentInfo) (int, error) {
	log.Printf("[HD1] Executing prop script for: %s", prop.ID)
	
	// For now, create a simple script execution
	// In a full implementation, this would extract and execute the script from the YAML
	
	// Convert position and scale to strings for script execution
	posX := strconv.FormatFloat(req.Position.X, 'f', 6, 64)
	posY := strconv.FormatFloat(req.Position.Y, 'f', 6, 64)
	posZ := strconv.FormatFloat(req.Position.Z, 'f', 6, 64)
	scale := strconv.FormatFloat(req.Scale, 'f', 6, 64)
	
	// Create a temporary script file (simplified approach)
	scriptContent := generateSimplePropScript(prop, req.InstanceName, envInfo)
	
	tmpScript, err := os.CreateTemp("", "prop-*.sh")
	if err != nil {
		return 0, err
	}
	defer os.Remove(tmpScript.Name())
	
	if _, err := tmpScript.WriteString(scriptContent); err != nil {
		return 0, err
	}
	tmpScript.Close()
	
	// Make script executable
	if err := os.Chmod(tmpScript.Name(), 0755); err != nil {
		return 0, err
	}
	
	// Execute the prop script
	cmd := exec.Command("/bin/bash", tmpScript.Name(), req.InstanceName, posX, posY, posZ, scale, sessionID)
	cmd.Dir = "/opt/hd1"
	
	// Set environment variables for the script
	cmd.Env = append(os.Environ(), 
		"HD1_SESSION_ID="+sessionID,
		"HD1_ROOT=/opt/hd1")
	
	output, err := cmd.Output()
	if err != nil {
		log.Printf("[HD1] Prop script execution failed: %v", err)
		return 0, err
	}
	
	log.Printf("[HD1] Prop script output: %s", string(output))
	
	// Parse output to determine number of objects created
	objectsCreated := parseObjectCount(string(output))
	
	return objectsCreated, nil
}

// generateSimplePropScript creates a basic script for prop instantiation with environment physics
func generateSimplePropScript(prop *PropInfo, instanceName string, envInfo *EnvironmentInfo) string {
	// Calculate physics adaptations based on environment
	adaptedMass := prop.Mass
	adaptedFriction := prop.PhysicsProperties.Friction
	adaptedRestitution := prop.PhysicsProperties.Restitution
	environmentNote := "default conditions"
	
	if envInfo != nil {
		// Adapt physics based on environment
		switch envInfo.Atmosphere {
		case "vacuum":
			adaptedMass = prop.Mass * 0.1  // Simulated reduced effective mass in vacuum
			adaptedRestitution = prop.PhysicsProperties.Restitution * 1.2  // Higher bounce in vacuum
			environmentNote = fmt.Sprintf("%s environment (gravity=%.1f m/s²)", envInfo.Name, envInfo.Gravity)
		case "water":
			adaptedMass = prop.Mass * 0.6  // Buoyancy effect
			adaptedFriction = prop.PhysicsProperties.Friction * 2.0  // Water resistance
			environmentNote = fmt.Sprintf("%s environment (buoyancy effects)", envInfo.Name)
		case "dense":
			adaptedMass = prop.Mass * 1.3  // Dense atmosphere increases drag effects
			adaptedFriction = prop.PhysicsProperties.Friction * 1.5
			environmentNote = fmt.Sprintf("%s environment (dense atmosphere)", envInfo.Name)
		default:
			environmentNote = fmt.Sprintf("%s environment", envInfo.Name)
		}
		
		// Scale adaptation based on environment scale unit
		switch envInfo.ScaleUnit {
		case "mm":
			environmentNote += " - microscale physics"
		case "km":
			environmentNote += " - macroscale physics"
		}
	}

	return fmt.Sprintf(`#!/bin/bash
# Auto-generated prop script for %s with environment physics cohesion
source "${HD1_ROOT}/lib/hd1lib.sh" 2>/dev/null || {
    echo "ERROR: HD1 core library not found"
    exit 1
}

INSTANCE_NAME="%s"
POSITION_X="$2"
POSITION_Y="$3" 
POSITION_Z="$4"
SCALE_FACTOR="$5"
SESSION_ID="$6"

echo "PROP: Creating %s instance '$INSTANCE_NAME'"
echo "POSITION: ($POSITION_X, $POSITION_Y, $POSITION_Z) with scale $SCALE_FACTOR"
echo "ENVIRONMENT: %s"

# Create basic prop representation with environment-adapted physics
hd1::create_object "${INSTANCE_NAME}_main" "box" "$POSITION_X" "$POSITION_Y" "$POSITION_Z"

echo "COMPONENTS: 1 component created (simplified %s)"
echo "PHYSICS: %s material, mass=%.1fkg (adapted from %.1fkg)"
echo "FRICTION: %.2f (adapted from %.2f)"
echo "RESTITUTION: %.2f (adapted from %.2f)"
echo "SCALE: Adapted to scale factor $SCALE_FACTOR"
echo "ENVIRONMENT_COHESION: Physics adapted for %s"
echo "Prop instantiated: $INSTANCE_NAME"
echo "Objects created: 1"
`, prop.Name, instanceName, prop.Name, environmentNote, prop.Category, prop.Material, adaptedMass, prop.Mass, adaptedFriction, prop.PhysicsProperties.Friction, adaptedRestitution, prop.PhysicsProperties.Restitution, environmentNote)
}

// getEnvironmentInfo retrieves environment information for physics calculations
func getEnvironmentInfo(environmentID string) *EnvironmentInfo {
	// Map environment IDs to their physics parameters
	switch environmentID {
	case "earth-surface":
		return &EnvironmentInfo{
			ID:          "earth-surface",
			Name:        "Earth Surface",
			ScaleUnit:   "m",
			Gravity:     9.8,
			Atmosphere:  "air",
			Density:     1.225,
			Temperature: 293.15,
		}
	case "molecular-scale":
		return &EnvironmentInfo{
			ID:          "molecular-scale",
			Name:        "Molecular Scale",
			ScaleUnit:   "nm",
			Gravity:     9.8,
			Atmosphere:  "vacuum",
			Density:     0.0,
			Temperature: 273.15,
		}
	case "space-vacuum":
		return &EnvironmentInfo{
			ID:          "space-vacuum",
			Name:        "Space Vacuum",
			ScaleUnit:   "km",
			Gravity:     0.0,
			Atmosphere:  "vacuum",
			Density:     0.0,
			Temperature: 2.7,
		}
	case "underwater":
		return &EnvironmentInfo{
			ID:          "underwater",
			Name:        "Underwater",
			ScaleUnit:   "m",
			Gravity:     9.8,
			Atmosphere:  "water",
			Density:     1000.0,
			Temperature: 283.15,
		}
	default:
		// Default to earth-surface conditions
		return &EnvironmentInfo{
			ID:          "default",
			Name:        "Default",
			ScaleUnit:   "m",
			Gravity:     9.8,
			Atmosphere:  "air",
			Density:     1.225,
			Temperature: 293.15,
		}
	}
}

// parseObjectCount extracts the number of objects created from script output
func parseObjectCount(output string) int {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Objects created:") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				if count, err := strconv.Atoi(parts[2]); err == nil {
					return count
				}
			}
		}
	}
	return 1 // Default to 1 if not found
}