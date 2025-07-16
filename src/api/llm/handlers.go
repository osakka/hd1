package llm

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"holodeck1/auth"
	"holodeck1/llm"
	"holodeck1/logging"
)

type Handler struct {
	llmManager *llm.Manager
}

func NewHandler(llmManager *llm.Manager) *Handler {
	return &Handler{
		llmManager: llmManager,
	}
}

func (h *Handler) CreateAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req llm.CreateAvatarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.SessionID == uuid.Nil {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	if req.Provider == "" {
		http.Error(w, "Provider is required", http.StatusBadRequest)
		return
	}
	if req.Model == "" {
		http.Error(w, "Model is required", http.StatusBadRequest)
		return
	}

	avatar, err := h.llmManager.CreateAvatar(ctx, &req)
	if err != nil {
		logging.Error("failed to create LLM avatar", map[string]interface{}{
			"session_id": req.SessionID,
			"name":       req.Name,
			"provider":   req.Provider,
			"model":      req.Model,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"avatar":  avatar,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	avatarID, err := uuid.Parse(vars["avatarId"])
	if err != nil {
		http.Error(w, "Invalid avatar ID", http.StatusBadRequest)
		return
	}

	avatar, err := h.llmManager.GetAvatar(ctx, avatarID)
	if err != nil {
		if err.Error() == "avatar not found" {
			http.Error(w, "Avatar not found", http.StatusNotFound)
			return
		}
		logging.Error("failed to get LLM avatar", map[string]interface{}{
			"avatar_id": avatarID,
			"error":     err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"avatar":  avatar,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetAvatarsBySession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionId"])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// Parse query parameters
	filter := &llm.AvatarFilter{
		SessionID: sessionID,
		Type:      r.URL.Query().Get("type"),
		Provider:  r.URL.Query().Get("provider"),
		State:     r.URL.Query().Get("state"),
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

	avatars, err := h.llmManager.GetAvatarsBySession(ctx, filter)
	if err != nil {
		logging.Error("failed to get avatars by session", map[string]interface{}{
			"session_id": sessionID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"avatars": avatars,
		"pagination": map[string]interface{}{
			"limit":    filter.Limit,
			"offset":   filter.Offset,
			"has_more": len(avatars) >= filter.Limit,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Chat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	avatarID, err := uuid.Parse(vars["avatarId"])
	if err != nil {
		http.Error(w, "Invalid avatar ID", http.StatusBadRequest)
		return
	}

	var req llm.ChatRequest
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

	req.AvatarID = avatarID
	req.UserID = user.ID

	// Validate required fields
	if req.Message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	chatResponse, err := h.llmManager.Chat(ctx, &req)
	if err != nil {
		logging.Error("failed to generate chat response", map[string]interface{}{
			"avatar_id": avatarID,
			"user_id":   user.ID,
			"error":     err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"chat":    chatResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) UpdateAvatarPosition(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	avatarID, err := uuid.Parse(vars["avatarId"])
	if err != nil {
		http.Error(w, "Invalid avatar ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Position *llm.Position3D `json:"position"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Position == nil {
		http.Error(w, "Position is required", http.StatusBadRequest)
		return
	}

	err = h.llmManager.UpdateAvatarPosition(ctx, avatarID, req.Position)
	if err != nil {
		logging.Error("failed to update avatar position", map[string]interface{}{
			"avatar_id": avatarID,
			"position":  req.Position,
			"error":     err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Avatar position updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) UpdateAvatarState(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	avatarID, err := uuid.Parse(vars["avatarId"])
	if err != nil {
		http.Error(w, "Invalid avatar ID", http.StatusBadRequest)
		return
	}

	var req struct {
		State string `json:"state"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.State == "" {
		http.Error(w, "State is required", http.StatusBadRequest)
		return
	}

	err = h.llmManager.UpdateAvatarState(ctx, avatarID, req.State)
	if err != nil {
		logging.Error("failed to update avatar state", map[string]interface{}{
			"avatar_id": avatarID,
			"state":     req.State,
			"error":     err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Avatar state updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) DeleteAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	avatarID, err := uuid.Parse(vars["avatarId"])
	if err != nil {
		http.Error(w, "Invalid avatar ID", http.StatusBadRequest)
		return
	}

	err = h.llmManager.DeleteAvatar(ctx, avatarID)
	if err != nil {
		if err.Error() == "avatar not found" {
			http.Error(w, "Avatar not found", http.StatusNotFound)
			return
		}
		logging.Error("failed to delete avatar", map[string]interface{}{
			"avatar_id": avatarID,
			"error":     err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetConversationHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	avatarID, err := uuid.Parse(vars["avatarId"])
	if err != nil {
		http.Error(w, "Invalid avatar ID", http.StatusBadRequest)
		return
	}

	// Parse limit parameter
	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	conversations, err := h.llmManager.GetConversationHistory(ctx, avatarID, limit)
	if err != nil {
		logging.Error("failed to get conversation history", map[string]interface{}{
			"avatar_id": avatarID,
			"limit":     limit,
			"error":     err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":       true,
		"conversations": conversations,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetProviders(w http.ResponseWriter, r *http.Request) {
	providers := h.llmManager.GetProviders()

	response := map[string]interface{}{
		"success":   true,
		"providers": providers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}