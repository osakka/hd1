package sessions

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
)

// SyncSessionState - POST /sessions/{sessionId}/room/sync
func SyncSessionStateHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type
	h, ok := hub.(*server.Hub)
	if !ok {
		logging.Error("hub interface cast failed", map[string]interface{}{
			"expected_type": "*server.Hub",
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Extract session ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}
	sessionID := pathParts[2]
	
	// Parse request body (optional)
	var request struct {
		ForceFullSync bool `json:"force_full_sync"`
	}
	
	// Try to decode request body, but don't fail if it's empty
	json.NewDecoder(r.Body).Decode(&request)
	
	// Verify session exists
	_, exists := h.GetStore().GetSession(sessionID)
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	
	// Force synchronization of session state
	syncedClients, err := h.SyncSessionState(sessionID, request.ForceFullSync)
	if err != nil {
		logging.Error("session state sync failed", map[string]interface{}{
			"session_id": sessionID,
			"error":      err.Error(),
		})
		http.Error(w, "Failed to sync session state: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	syncTimestamp := time.Now()
	
	logging.Info("session state synchronized", map[string]interface{}{
		"session_id":      sessionID,
		"synced_clients":  syncedClients,
		"force_full_sync": request.ForceFullSync,
		"sync_timestamp":  syncTimestamp,
	})
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":        true,
		"synced_clients": syncedClients,
		"sync_timestamp": syncTimestamp,
		"message":        "Session state synchronized successfully",
	})
}