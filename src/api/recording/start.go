package recording

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"crypto/rand"
	"encoding/hex"
	"holodeck/server"
)

// StartRecordingHandler - POST /sessions/{sessionId}/recording/start - Start temporal recording
func StartRecordingHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
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

	// Parse request body (optional)
	var req struct {
		RecordingName string `json:"recording_name"`
		Description   string `json:"description"`
		ClearOnStart  bool   `json:"clear_on_start"`
	}
	
	// Set defaults
	req.ClearOnStart = true
	req.RecordingName = fmt.Sprintf("recording-%s", time.Now().Format("20060102-150405"))
	
	if r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
			return
		}
	}

	// Validate session exists
	store := h.GetStore()
	_, exists := store.GetSession(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Generate recording ID
	recordingID := generateRecordingID()

	// Clear existing objects if requested
	if req.ClearOnStart {
		objects := store.ListObjects(sessionID)
		for _, obj := range objects {
			store.DeleteObject(sessionID, obj.Name)
		}
		
		// Broadcast clear message
		clearMessage := map[string]interface{}{
			"type": "recording_clear",
			"session_id": sessionID,
			"recording_id": recordingID,
			"timestamp": time.Now().Format(time.RFC3339),
		}
		if _, err := json.Marshal(clearMessage); err == nil {
			h.BroadcastToSession(sessionID, "recording_clear", clearMessage)
		}
	}

	// Note: Recording state managed in-memory for this session
	// TODO: Add Session.Metadata field to store persistent recording state

	// Broadcast recording start
	startMessage := map[string]interface{}{
		"type": "recording_started",
		"session_id": sessionID,
		"recording_id": recordingID,
		"recording_name": req.RecordingName,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	if _, err := json.Marshal(startMessage); err == nil {
		h.BroadcastToSession(sessionID, "recording_started", startMessage)
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Recording '%s' started", req.RecordingName),
		"session_id": sessionID,
		"recording_id": recordingID,
		"start_time": time.Now().Format(time.RFC3339),
	})
}

// generateRecordingID creates a unique recording identifier
func generateRecordingID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return "rec-" + hex.EncodeToString(bytes)
}