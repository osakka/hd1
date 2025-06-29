package recording

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"holodeck/server"
)

// GetRecordingStatusHandler - GET /sessions/{sessionId}/recording/status - Get recording status
func GetRecordingStatusHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
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

	// Build response
	response := map[string]interface{}{
		"session_id": sessionID,
		"recording_active": false,
		"current_recording": nil,
		"available_recordings": []map[string]interface{}{},
	}

	// Note: Recording state not persistent yet - return placeholder
	// TODO: Check session recording metadata when available

	// Get available recordings
	recordingsDir := "/opt/holo-deck/recordings"
	if entries, err := os.ReadDir(recordingsDir); err == nil {
		recordings := []map[string]interface{}{}
		
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".thd") {
				recordingID := strings.TrimSuffix(entry.Name(), ".thd")
				info, _ := entry.Info()
				
				recordings = append(recordings, map[string]interface{}{
					"recording_id": recordingID,
					"name": recordingID, // Could parse from file if needed
					"created_at": info.ModTime().Format("2006-01-02T15:04:05Z07:00"),
					"size": info.Size(),
				})
			}
		}
		
		response["available_recordings"] = recordings
	}

	// Return status
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}