package clients

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"holodeck1/auth"
	"holodeck1/clients"
	"holodeck1/logging"
)

type Handler struct {
	clientManager *clients.Manager
}

func NewHandler(clientManager *clients.Manager) *Handler {
	return &Handler{
		clientManager: clientManager,
	}
}

func (h *Handler) RegisterClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req clients.ClientRegistration
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

	req.UserID = user.ID

	// Validate required fields
	if req.SessionID == uuid.Nil {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}
	if req.Type == "" {
		http.Error(w, "Client type is required", http.StatusBadRequest)
		return
	}
	if req.Platform == "" {
		http.Error(w, "Platform is required", http.StatusBadRequest)
		return
	}
	if req.Version == "" {
		http.Error(w, "Version is required", http.StatusBadRequest)
		return
	}

	client, err := h.clientManager.RegisterClient(ctx, &req)
	if err != nil {
		logging.Error("failed to register client", map[string]interface{}{
			"session_id": req.SessionID,
			"user_id":    req.UserID,
			"type":       req.Type,
			"platform":   req.Platform,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"client":  client,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	clientID, err := uuid.Parse(vars["clientId"])
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	client, err := h.clientManager.GetClient(ctx, clientID)
	if err != nil {
		if err.Error() == "client not found" {
			http.Error(w, "Client not found", http.StatusNotFound)
			return
		}
		logging.Error("failed to get client", map[string]interface{}{
			"client_id": clientID,
			"error":     err.Error(),
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

func (h *Handler) GetClientsBySession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionId"])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// Parse query parameters
	filter := &clients.ClientFilter{
		SessionID: sessionID,
		Type:      r.URL.Query().Get("type"),
		Platform:  r.URL.Query().Get("platform"),
		Status:    r.URL.Query().Get("status"),
	}

	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err == nil {
			filter.UserID = userID
		}
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filter.Limit = l
		}
	}

	if offset := r.URL.Query().Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			filter.Offset = o
		}
	}

	// Set default limit if not provided
	if filter.Limit == 0 {
		filter.Limit = 50
	}

	clients, err := h.clientManager.GetClientsBySession(ctx, filter)
	if err != nil {
		logging.Error("failed to get clients by session", map[string]interface{}{
			"session_id": sessionID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"clients": clients,
		"pagination": map[string]interface{}{
			"limit":    filter.Limit,
			"offset":   filter.Offset,
			"has_more": len(clients) >= filter.Limit,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	clientID, err := uuid.Parse(vars["clientId"])
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	var req clients.ClientUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.clientManager.UpdateClient(ctx, clientID, &req)
	if err != nil {
		if err.Error() == "client not found" {
			http.Error(w, "Client not found", http.StatusNotFound)
			return
		}
		logging.Error("failed to update client", map[string]interface{}{
			"client_id": clientID,
			"error":     err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Client updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) UnregisterClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	clientID, err := uuid.Parse(vars["clientId"])
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	err = h.clientManager.UnregisterClient(ctx, clientID)
	if err != nil {
		if err.Error() == "client not found" {
			http.Error(w, "Client not found", http.StatusNotFound)
			return
		}
		logging.Error("failed to unregister client", map[string]interface{}{
			"client_id": clientID,
			"error":     err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) BroadcastMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionId"])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	var req clients.ClientMessage
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

	// Set message ID and timestamp
	req.ID = uuid.New()
	req.Timestamp = time.Now()

	// Validate required fields
	if req.Type == "" {
		http.Error(w, "Message type is required", http.StatusBadRequest)
		return
	}
	if req.Action == "" {
		http.Error(w, "Message action is required", http.StatusBadRequest)
		return
	}

	err = h.clientManager.BroadcastToSession(ctx, sessionID, &req)
	if err != nil {
		logging.Error("failed to broadcast message", map[string]interface{}{
			"session_id": sessionID,
			"user_id":    user.ID,
			"type":       req.Type,
			"action":     req.Action,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Message broadcasted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetClientStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionId"])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	stats, err := h.clientManager.GetClientStats(ctx, sessionID)
	if err != nil {
		logging.Error("failed to get client stats", map[string]interface{}{
			"session_id": sessionID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"stats":   stats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetClientCapabilities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	clientID, err := uuid.Parse(vars["clientId"])
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	capabilities, err := h.clientManager.GetClientCapabilities(ctx, clientID)
	if err != nil {
		logging.Error("failed to get client capabilities", map[string]interface{}{
			"client_id": clientID,
			"error":     err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":      true,
		"capabilities": capabilities,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) UpdateClientCapabilities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	clientID, err := uuid.Parse(vars["clientId"])
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Capabilities []string `json:"capabilities"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Capabilities) == 0 {
		http.Error(w, "Capabilities list cannot be empty", http.StatusBadRequest)
		return
	}

	err = h.clientManager.UpdateClientCapabilities(ctx, clientID, req.Capabilities)
	if err != nil {
		logging.Error("failed to update client capabilities", map[string]interface{}{
			"client_id": clientID,
			"error":     err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Client capabilities updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HeartbeatClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	clientID, err := uuid.Parse(vars["clientId"])
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	err = h.clientManager.HeartbeatClient(ctx, clientID)
	if err != nil {
		if err.Error() == "client not found" {
			http.Error(w, "Client not found", http.StatusNotFound)
			return
		}
		logging.Error("failed to heartbeat client", map[string]interface{}{
			"client_id": clientID,
			"error":     err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Heartbeat received",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}