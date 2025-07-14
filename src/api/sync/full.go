package sync

import (
	"encoding/json"
	"net/http"

	"holodeck1/logging"
)

// FullSyncResponse represents the response for full synchronization
type FullSyncResponse struct {
	Success         bool                  `json:"success"`
	Operations      []OperationWithSeqNum `json:"operations"`
	CurrentSequence uint64                `json:"current_sequence"`
}

// GetFullSync handles GET /api/sync/full
func GetFullSync(w http.ResponseWriter, r *http.Request) {
	// Get hub from context
	hub := getHubFromContext(r)
	if hub == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get all operations
	operations := hub.GetSync().GetAllOperations()
	currentSeq := hub.GetSync().GetCurrentSequence()

	// Convert to response format
	var operationsWithSeq []OperationWithSeqNum
	for _, op := range operations {
		operationsWithSeq = append(operationsWithSeq, OperationWithSeqNum{
			SeqNum:    op.SeqNum,
			Operation: op,
		})
	}

	// Return response
	response := FullSyncResponse{
		Success:         true,
		Operations:      operationsWithSeq,
		CurrentSequence: currentSeq,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	logging.Info("full sync retrieved via API", map[string]interface{}{
		"count":            len(operations),
		"current_sequence": currentSeq,
	})
}