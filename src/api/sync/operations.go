package sync

import (
	"encoding/json"
	"net/http"
	"time"

	"holodeck1/logging"
	"holodeck1/server"
	"holodeck1/sync"
)

// SubmitOperationRequest represents the request to submit an operation
type SubmitOperationRequest struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// SubmitOperationResponse represents the response after submitting an operation
type SubmitOperationResponse struct {
	Success bool   `json:"success"`
	SeqNum  uint64 `json:"seq_num"`
	Message string `json:"message"`
}

// SubmitOperation handles POST /api/sync/operations
func SubmitOperation(w http.ResponseWriter, r *http.Request) {
	var req SubmitOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate operation type
	validTypes := map[string]bool{
		"avatar_move":   true,
		"entity_create": true,
		"entity_update": true,
		"entity_delete": true,
		"scene_update":  true,
	}

	if !validTypes[req.Type] {
		http.Error(w, "Invalid operation type", http.StatusBadRequest)
		return
	}

	// Get client ID from request (could be from session, header, etc.)
	clientID := getClientID(r)

	// Create operation
	operation := &sync.Operation{
		ClientID:  clientID,
		Type:      req.Type,
		Data:      req.Data,
		Timestamp: time.Now(),
	}

	// Get hub from context (needs to be injected by router)
	hub := getHubFromContext(r)
	if hub == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Submit operation to sync system
	hub.GetSync().SubmitOperation(operation)

	// Return response
	response := SubmitOperationResponse{
		Success: true,
		SeqNum:  operation.SeqNum,
		Message: "Operation submitted",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	logging.Info("operation submitted via API", map[string]interface{}{
		"client_id": clientID,
		"type":      req.Type,
		"seq_num":   operation.SeqNum,
	})
}

// Helper functions
func getClientID(r *http.Request) string {
	// Try to get client ID from various sources
	if clientID := r.Header.Get("X-Client-ID"); clientID != "" {
		return clientID
	}
	
	// Generate a client ID if none provided
	return "api-client-" + time.Now().Format("20060102150405")
}

func getHubFromContext(r *http.Request) *server.Hub {
	// This will be injected by the router
	if hub := r.Context().Value("hub"); hub != nil {
		if h, ok := hub.(*server.Hub); ok {
			return h
		}
	}
	return nil
}