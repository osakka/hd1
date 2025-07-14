package sync

import (
	"encoding/json"
	"net/http"
)

// SyncStatsResponse represents the response for sync statistics
type SyncStatsResponse struct {
	Success bool                   `json:"success"`
	Stats   map[string]interface{} `json:"stats"`
}

// GetSyncStats handles GET /api/sync/stats
func GetSyncStats(w http.ResponseWriter, r *http.Request) {
	// Get hub from context
	hub := getHubFromContext(r)
	if hub == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get statistics from hub
	stats := hub.GetStats()

	// Return response
	response := SyncStatsResponse{
		Success: true,
		Stats:   stats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}