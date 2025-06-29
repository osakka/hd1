package scenes

import (
	"encoding/json"
	"net/http"
	"holodeck/server"
	"fmt"
	"log"
	"strings"
	"os/exec"
	"strconv"
	"regexp"
	"os"
	"path/filepath"
)

type LoadSceneRequest struct {
	SessionID string `json:"session_id"`
}

type LoadSceneResponse struct {
	Success        bool   `json:"success"`
	SceneID        string `json:"scene_id"`
	SessionID      string `json:"session_id"`
	ObjectsCreated int    `json:"objects_created"`
	Message        string `json:"message"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

// LoadSceneHandler loads a predefined scene into the specified session
func LoadSceneHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type
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

	// Extract scene ID from URL path
	sceneID := extractSceneID(r.URL.Path)
	if sceneID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "missing_scene_id",
			Message: "Scene ID is required",
		})
		return
	}

	// Parse request body
	var req LoadSceneRequest
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

	// Clear existing objects in session
	clearAllObjects(req.SessionID, h)
	
	// Send clear command via WebSocket
	h.BroadcastToSession(req.SessionID, "canvas_control", map[string]interface{}{
		"command": "clear",
		"clear":   true,
	})

	var objectsCreated int
	var message string

	// Execute scene script based on ID
	sceneScript := getSceneScript(sceneID)
	if sceneScript == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "scene_not_found",
			Message: "Scene '" + sceneID + "' not found",
		})
		return
	}

	// Execute scene script with session ID
	objectsCreated, message = executeSceneScript(sceneScript, req.SessionID)

	log.Printf("[THD] Scene '%s' loaded into session '%s' with %d objects", sceneID, req.SessionID, objectsCreated)

	response := LoadSceneResponse{
		Success:        true,
		SceneID:        sceneID,
		SessionID:      req.SessionID,
		ObjectsCreated: objectsCreated,
		Message:        message,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// extractSceneID extracts scene ID from URL path
func extractSceneID(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "scenes" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// clearAllObjects removes all objects from a session
func clearAllObjects(sessionID string, h *server.Hub) {
	objects := h.GetStore().ListObjects(sessionID)
	for _, obj := range objects {
		h.GetStore().DeleteObject(sessionID, obj.Name)
	}
}

// getSceneScript returns the script path for a given scene ID
func getSceneScript(sceneID string) string {
	// Dynamic scene discovery - single source of truth
	scenesDir := "/opt/holo-deck/share/scenes"
	scriptPath := filepath.Join(scenesDir, sceneID+".sh")
	
	// Check if file exists
	if _, err := os.Stat(scriptPath); err == nil {
		return scriptPath
	}
	
	return "" // Scene not found
}

// executeSceneScript runs the scene script and parses the output
func executeSceneScript(scriptPath string, sessionID string) (int, string) {
	log.Printf("[THD] Executing scene script: %s with session %s", scriptPath, sessionID)
	
	// Execute the scene script with session ID
	cmd := exec.Command("/bin/bash", scriptPath, sessionID)
	cmd.Dir = "/opt/holo-deck"  // Set working directory
	output, err := cmd.Output()
	
	if err != nil {
		log.Printf("[THD] Scene script execution failed: %v", err)
		return 0, fmt.Sprintf("Scene script execution failed: %v", err)
	}
	
	outputStr := string(output)
	log.Printf("[THD] Scene script output: %s", outputStr)
	
	// Parse object count from script output
	objectCount := parseObjectCount(outputStr)
	
	// Extract success message
	message := parseSceneMessage(outputStr)
	if message == "" {
		message = "Scene loaded successfully"
	}
	
	return objectCount, message
}

// parseObjectCount extracts object count from script output
func parseObjectCount(output string) int {
	re := regexp.MustCompile(`Objects created: (\d+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		if count, err := strconv.Atoi(matches[1]); err == nil {
			return count
		}
	}
	return 0
}

// parseSceneMessage extracts success message from script output
func parseSceneMessage(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "THD Scene") && strings.Contains(line, "loaded successfully") {
			return line
		}
	}
	return ""
}