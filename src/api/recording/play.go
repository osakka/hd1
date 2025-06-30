package recording

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"gopkg.in/yaml.v3"
	"holodeck1/server"
)

// PlayRecordingHandler - POST /sessions/{sessionId}/recording/play - Play back recorded sequence
func PlayRecordingHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
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

	// Parse request body
	var req struct {
		RecordingID     string  `json:"recording_id"`
		Speed           float64 `json:"speed"`
		ClearBeforePlay bool    `json:"clear_before_play"`
	}
	
	// Set defaults
	req.Speed = 1.0
	req.ClearBeforePlay = true
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	// Validate session exists
	store := h.GetStore()
	_, exists := store.GetSession(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Load recording file
	recordingPath := filepath.Join("/opt/hd1/recordings", req.RecordingID+".hd1")
	recordingData, err := os.ReadFile(recordingPath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": fmt.Sprintf("Recording '%s' not found", req.RecordingID),
		})
		return
	}

	// Parse recording file
	var recording RecordingFile
	if err := yaml.Unmarshal(recordingData, &recording); err != nil {
		http.Error(w, fmt.Sprintf("Invalid recording file format: %v", err), http.StatusInternalServerError)
		return
	}

	// Clear existing objects if requested
	if req.ClearBeforePlay {
		objects := store.ListObjects(sessionID)
		for _, obj := range objects {
			store.DeleteObject(sessionID, obj.Name)
		}
		
		// Broadcast clear message
		clearMessage := map[string]interface{}{
			"type": "playback_clear",
			"session_id": sessionID,
			"recording_id": req.RecordingID,
			"timestamp": time.Now().Format(time.RFC3339),
		}
		if _, err := json.Marshal(clearMessage); err == nil {
			h.BroadcastToSession(sessionID, "playback_clear", clearMessage)
		}
	}

	// Start async playback
	go func() {
		playbackRecording(h, sessionID, req.RecordingID, recording, req.Speed)
	}()

	// Return immediate response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Playback of recording '%s' started", recording.Name),
		"session_id": sessionID,
		"recording_id": req.RecordingID,
		"playback_duration": recording.Duration,
		"operations_count": recording.OperationsCount,
	})
}

// RecordingFile represents the structure of a .hd1 recording file
type RecordingFile struct {
	RecordingID      string                   `yaml:"recording_id"`
	Name             string                   `yaml:"name"`
	CreatedAt        string                   `yaml:"created_at"`
	Duration         string                   `yaml:"duration"`
	OperationsCount  int                      `yaml:"operations_count"`
	Operations       []RecordingOperation     `yaml:"operations"`
}

type RecordingOperation struct {
	Step      int    `yaml:"step"`
	Timestamp string `yaml:"timestamp"`
	Operation string `yaml:"operation"`
	Data      string `yaml:"data"`
}

// playbackRecording executes the temporal sequence
func playbackRecording(hub *server.Hub, sessionID, recordingID string, recording RecordingFile, speed float64) {
	if len(recording.Operations) == 0 {
		return
	}

	// Parse start time for relative timing
	baseTime, err := time.Parse(time.RFC3339, recording.Operations[0].Timestamp)
	if err != nil {
		return
	}

	// Broadcast playback start
	startMessage := map[string]interface{}{
		"type": "playback_started",
		"session_id": sessionID,
		"recording_id": recordingID,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	if _, err := json.Marshal(startMessage); err == nil {
		hub.BroadcastToSession(sessionID, "playback_started", startMessage)
	}

	// Execute operations with timing
	for _, op := range recording.Operations {
		opTime, err := time.Parse(time.RFC3339, op.Timestamp)
		if err != nil {
			continue
		}

		// Calculate delay from start (adjusted by speed)
		delay := time.Duration(float64(opTime.Sub(baseTime)) / speed)
		if delay > 0 {
			time.Sleep(delay)
		}

		// Execute operation (would need to integrate with object creation/updates)
		// This is a placeholder - would execute actual HD1 operations
		playbackMessage := map[string]interface{}{
			"type": "playback_operation",
			"session_id": sessionID,
			"recording_id": recordingID,
			"step": op.Step,
			"operation": op.Operation,
			"data": op.Data,
			"timestamp": time.Now().Format(time.RFC3339),
		}
		if _, err := json.Marshal(playbackMessage); err == nil {
			hub.BroadcastToSession(sessionID, "playback_operation", playbackMessage)
		}
	}

	// Broadcast playback complete
	completeMessage := map[string]interface{}{
		"type": "playback_complete",
		"session_id": sessionID,
		"recording_id": recordingID,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	if _, err := json.Marshal(completeMessage); err == nil {
		hub.BroadcastToSession(sessionID, "playback_complete", completeMessage)
	}
}