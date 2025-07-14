package sync

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"holodeck1/logging"
	"holodeck1/sync"
)

// MissingOperationsResponse represents the response for missing operations
type MissingOperationsResponse struct {
	Success    bool                   `json:"success"`
	Operations []OperationWithSeqNum  `json:"operations"`
}

// OperationWithSeqNum represents an operation with its sequence number
type OperationWithSeqNum struct {
	SeqNum    uint64           `json:"seq_num"`
	Operation *sync.Operation  `json:"operation"`
}

// GetMissingOperations handles GET /api/sync/missing/{from}/{to}
func GetMissingOperations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	
	// Parse parameters
	fromStr := vars["from"]
	toStr := vars["to"]
	
	from, err := strconv.ParseUint(fromStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid 'from' parameter", http.StatusBadRequest)
		return
	}
	
	to, err := strconv.ParseUint(toStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid 'to' parameter", http.StatusBadRequest)
		return
	}
	
	if from > to {
		http.Error(w, "'from' must be <= 'to'", http.StatusBadRequest)
		return
	}
	
	if to - from > 10000 {
		http.Error(w, "Range too large (max 10000 operations)", http.StatusBadRequest)
		return
	}

	// Get hub from context
	hub := getHubFromContext(r)
	if hub == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get missing operations
	operations := hub.GetSync().GetMissingOperations(from, to)

	// Convert to response format
	var operationsWithSeq []OperationWithSeqNum
	for _, op := range operations {
		operationsWithSeq = append(operationsWithSeq, OperationWithSeqNum{
			SeqNum:    op.SeqNum,
			Operation: op,
		})
	}

	// Return response
	response := MissingOperationsResponse{
		Success:    true,
		Operations: operationsWithSeq,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	logging.Debug("missing operations retrieved via API", map[string]interface{}{
		"from":  from,
		"to":    to,
		"count": len(operations),
	})
}