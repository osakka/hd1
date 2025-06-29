package recording

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"holodeck/server"
)

// StopRecordingHandler - POST /sessions/{sessionId}/recording/stop - Stop temporal recording
func StopRecordingHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	h, ok := hub.(*server.Hub)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Extract session ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 || pathParts[1] != "sessions" {
		http.Error(w, "Invalid session ID in path", http.StatusBadRequest)
		return
	}
	sessionID := pathParts[2]

	// Validate session exists
	store := h.GetStore()
	_, exists := store.GetSession(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Note: Recording metadata not persistent yet - using placeholder values
	// TODO: Implement proper recording state management
	recordingID := "rec-placeholder"
	recordingName := "Test Recording"
	startTime := time.Now().Add(-5 * time.Minute).Format(time.RFC3339)
	operations := []map[string]interface{}{}

	// Calculate duration
	startTimeParsed, _ := time.Parse(time.RFC3339, startTime)
	duration := time.Since(startTimeParsed)

	// Generate recording file
	recordingPath := filepath.Join("/opt/holo-deck/recordings", recordingID+".thd")
	if err := os.MkdirAll(filepath.Dir(recordingPath), 0755); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create recordings directory: %v", err), http.StatusInternalServerError)
		return
	}

	// Create recording file content
	recordingContent := generateRecordingFile(recordingID, recordingName, startTime, duration, operations)
	
	if err := os.WriteFile(recordingPath, []byte(recordingContent), 0644); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save recording: %v", err), http.StatusInternalServerError)
		return
	}

	// Note: Recording state cleanup not implemented yet
	// TODO: Clear recording state when Session.Metadata is available

	// Broadcast recording stop
	stopMessage := map[string]interface{}{
		"type": "recording_stopped",
		"session_id": sessionID,
		"recording_id": recordingID,
		"duration": duration.String(),
		"operations_count": len(operations),
		"timestamp": time.Now().Format(time.RFC3339),
	}
	if _, err := json.Marshal(stopMessage); err == nil {
		h.BroadcastToSession(sessionID, "recording_stopped", stopMessage)
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Recording '%s' stopped and saved", recordingName),
		"session_id": sessionID,
		"recording_id": recordingID,
		"duration": duration.String(),
		"operations_recorded": len(operations),
		"recording_path": recordingPath,
	})
}

// generateRecordingFile creates a .thd recording file with temporal sequence
func generateRecordingFile(recordingID, name, startTime string, duration time.Duration, operations []map[string]interface{}) string {
	content := fmt.Sprintf(`# THD Recording: %s
# ID: %s
# Created: %s
# Duration: %s
# Operations: %d
#
# This is a temporal recording file (.thd) containing a sequence of 
# holodeck operations with timestamps for accurate playback.

recording_id: %s
name: %s
created_at: %s
duration: %s
operations_count: %d

operations:
`, name, recordingID, startTime, duration.String(), len(operations), 
   recordingID, name, startTime, duration.String(), len(operations))

	// Add each operation with timing
	for i, op := range operations {
		content += fmt.Sprintf(`
  - step: %d
    timestamp: %s
    operation: %s
    data: %s
`, i+1, op["timestamp"], op["type"], op["data"])
	}

	return content
}