package ot

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"holodeck1/auth"
	"holodeck1/logging"
	"holodeck1/ot"
)

type Handler struct {
	otManager *ot.Manager
}

func NewHandler(otManager *ot.Manager) *Handler {
	return &Handler{
		otManager: otManager,
	}
}

func (h *Handler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req ot.CreateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.SessionID == uuid.Nil {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	document, err := h.otManager.CreateDocument(ctx, &req)
	if err != nil {
		logging.Error("failed to create OT document", map[string]interface{}{
			"session_id": req.SessionID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":  true,
		"document": document,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) JoinDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	documentID, err := uuid.Parse(vars["documentId"])
	if err != nil {
		http.Error(w, "Invalid document ID", http.StatusBadRequest)
		return
	}

	var req ot.JoinDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user from context
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	req.DocumentID = documentID
	req.UserID = user.ID

	client, err := h.otManager.JoinDocument(ctx, &req)
	if err != nil {
		logging.Error("failed to join OT document", map[string]interface{}{
			"document_id": documentID,
			"user_id":     req.UserID,
			"error":       err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"client":  client,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) ApplyOperation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	documentID, err := uuid.Parse(vars["documentId"])
	if err != nil {
		http.Error(w, "Invalid document ID", http.StatusBadRequest)
		return
	}

	var req ot.ApplyOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user from context
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	req.DocumentID = documentID
	req.Operation.UserID = user.ID

	// Validate operation
	if req.Operation.Type == "" {
		http.Error(w, "Operation type is required", http.StatusBadRequest)
		return
	}

	result, err := h.otManager.ApplyOperation(ctx, &req)
	if err != nil {
		logging.Error("failed to apply OT operation", map[string]interface{}{
			"document_id": documentID,
			"client_id":   req.ClientID,
			"user_id":     user.ID,
			"error":       err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !result.Success {
		http.Error(w, result.Error, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) GetDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	documentID, err := uuid.Parse(vars["documentId"])
	if err != nil {
		http.Error(w, "Invalid document ID", http.StatusBadRequest)
		return
	}

	document, err := h.otManager.GetDocument(ctx, documentID)
	if err != nil {
		if err.Error() == "document not found" {
			http.Error(w, "Document not found", http.StatusNotFound)
			return
		}
		logging.Error("failed to get OT document", map[string]interface{}{
			"document_id": documentID,
			"error":       err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":  true,
		"document": document,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetDocumentsBySession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionId"])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	documents, err := h.otManager.GetDocumentsBySession(ctx, sessionID)
	if err != nil {
		logging.Error("failed to get documents by session", map[string]interface{}{
			"session_id": sessionID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":   true,
		"documents": documents,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) LeaveDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	documentID, err := uuid.Parse(vars["documentId"])
	if err != nil {
		http.Error(w, "Invalid document ID", http.StatusBadRequest)
		return
	}

	// Get user from context
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	err = h.otManager.LeaveDocument(ctx, documentID, user.ID)
	if err != nil {
		logging.Error("failed to leave OT document", map[string]interface{}{
			"document_id": documentID,
			"user_id":     user.ID,
			"error":       err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Successfully left document",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetOperationHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	documentID, err := uuid.Parse(vars["documentId"])
	if err != nil {
		http.Error(w, "Invalid document ID", http.StatusBadRequest)
		return
	}

	// Parse query parameters
	since := time.Time{}
	if sinceStr := r.URL.Query().Get("since"); sinceStr != "" {
		since, err = time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			http.Error(w, "Invalid since parameter", http.StatusBadRequest)
			return
		}
	}

	operations, err := h.otManager.GetOperationHistory(ctx, documentID, since)
	if err != nil {
		logging.Error("failed to get operation history", map[string]interface{}{
			"document_id": documentID,
			"since":       since,
			"error":       err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":    true,
		"operations": operations,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}