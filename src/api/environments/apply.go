package environments

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"os"
	"os/exec"
	
	"holodeck1/logging"
	"holodeck1/server"
)

// ApplyEnvironmentRequest matches the OpenAPI request schema exactly
type ApplyEnvironmentRequest struct {
	SessionID string `json:"session_id"`
}

// ApplyEnvironmentResponse matches the OpenAPI response schema exactly
type ApplyEnvironmentResponse struct {
	Success       bool   `json:"success"`
	EnvironmentID string `json:"environment_id"`
	SessionID     string `json:"session_id"`
	Message       string `json:"message"`
}

// ErrorResponse matches the OpenAPI error schema
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ApplyEnvironmentHandler applies environment settings to specified session
func ApplyEnvironmentHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
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

	// Extract environment ID from URL path (follows existing pattern from scenes)
	environmentID := extractEnvironmentID(r.URL.Path)
	if environmentID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "missing_environment_id",
			Message: "Environment ID is required",
		})
		return
	}

	// Parse request body
	var req ApplyEnvironmentRequest
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

	if req.SessionID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "missing_session_id",
			Message: "Session ID is required",
		})
		return
	}

	// Verify session exists
	if _, exists := h.GetStore().GetSession(req.SessionID); !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "session_not_found",
			Message: "Session does not exist",
		})
		return
	}

	// Find and execute environment script
	environmentScript := getEnvironmentScript(environmentID)
	if environmentScript == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "environment_not_found",
			Message: "Environment '" + environmentID + "' not found",
		})
		return
	}

	// Execute environment script with session ID
	err := executeEnvironmentScript(environmentScript, req.SessionID, h)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "environment_application_failed",
			Message: "Failed to apply environment: " + err.Error(),
		})
		return
	}

	logging.Info("environment applied successfully", map[string]interface{}{
		"environment_id": environmentID,
		"session_id":     req.SessionID,
	})

	// Environments are now applied via channels, not session tracking
	logging.Info("environment applied via channel configuration", map[string]interface{}{
		"environment_id": environmentID,
		"session_id":     req.SessionID,
	})

	// Return success response
	response := ApplyEnvironmentResponse{
		Success:       true,
		EnvironmentID: environmentID,
		SessionID:     req.SessionID,
		Message:       "Environment '" + getEnvironmentName(environmentID) + "' applied successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// extractEnvironmentID extracts environment ID from URL path (follows scenes pattern)
func extractEnvironmentID(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "environments" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// getEnvironmentScript returns the script path for a given environment ID
func getEnvironmentScript(environmentID string) string {
	environmentsDir := "/opt/hd1/share/environments"
	scriptPath := filepath.Join(environmentsDir, environmentID+".sh")
	
	// Check if file exists
	if _, err := os.Stat(scriptPath); err == nil {
		return scriptPath
	}
	
	return "" // Environment not found
}

// executeEnvironmentScript runs the environment script and applies settings to session
func executeEnvironmentScript(scriptPath string, sessionID string, h *server.Hub) error {
	logging.Debug("executing environment script", map[string]interface{}{
		"script_path": scriptPath,
		"session_id":  sessionID,
	})
	
	// Execute the environment script with session ID
	cmd := exec.Command("/bin/bash", scriptPath, sessionID)
	cmd.Dir = "/opt/hd1"  // Set working directory
	
	// Set environment variables for the script
	cmd.Env = append(os.Environ(), 
		"HD1_SESSION_ID="+sessionID,
		"HD1_ROOT=/opt/hd1")
	
	output, err := cmd.Output()
	if err != nil {
		logging.Error("script execution failed", map[string]interface{}{
			"script_path": scriptPath,
			"session_id":  sessionID,
			"error":       err.Error(),
		})
		return err
	}
	
	logging.Debug("environment script executed successfully", map[string]interface{}{
		"script_path": scriptPath,
		"session_id":  sessionID,
		"output":      string(output),
	})
	
	// Send environment change notification via WebSocket
	h.BroadcastToSession(sessionID, "environment_changed", map[string]interface{}{
		"environment_applied": true,
		"session_id":         sessionID,
		"timestamp":          "now",
	})
	
	return nil
}

// getEnvironmentName returns human-readable name for environment ID
func getEnvironmentName(environmentID string) string {
	nameMap := map[string]string{
		"earth-surface":    "Earth Surface",
		"molecular-scale":  "Molecular Scale", 
		"space-vacuum":     "Space Vacuum",
		"underwater":       "Underwater",
	}
	
	if name, exists := nameMap[environmentID]; exists {
		return name
	}
	
	return environmentID // Fallback to ID
}