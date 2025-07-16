package assets

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"holodeck1/assets"
	"holodeck1/auth"
	"holodeck1/logging"
)

type Handler struct {
	assetManager *assets.Manager
}

func NewHandler(assetManager *assets.Manager) *Handler {
	return &Handler{
		assetManager: assetManager,
	}
}

func (h *Handler) UploadAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse multipart form
	err := r.ParseMultipartForm(100 << 20) // 100MB
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get user from context
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse session ID
	sessionIDStr := r.FormValue("session_id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// Create upload request
	req := &assets.UploadRequest{
		SessionID: sessionID,
		UserID:    user.ID,
		Name:      r.FormValue("name"),
		Type:      r.FormValue("type"),
		Tags:      []string{},
		Metadata:  make(map[string]interface{}),
	}

	// Parse tags if provided
	if tagsStr := r.FormValue("tags"); tagsStr != "" {
		err = json.Unmarshal([]byte(tagsStr), &req.Tags)
		if err != nil {
			http.Error(w, "Invalid tags format", http.StatusBadRequest)
			return
		}
	}

	// Parse metadata if provided
	if metadataStr := r.FormValue("metadata"); metadataStr != "" {
		err = json.Unmarshal([]byte(metadataStr), &req.Metadata)
		if err != nil {
			http.Error(w, "Invalid metadata format", http.StatusBadRequest)
			return
		}
	}

	// Upload asset
	asset, err := h.assetManager.UploadAsset(ctx, req, file, header)
	if err != nil {
		logging.Error("failed to upload asset", map[string]interface{}{
			"session_id": sessionID,
			"user_id":    user.ID,
			"name":       req.Name,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"asset":   asset,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	assetID, err := uuid.Parse(vars["assetId"])
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	asset, err := h.assetManager.GetAsset(ctx, assetID)
	if err != nil {
		if err.Error() == "asset not found" {
			http.Error(w, "Asset not found", http.StatusNotFound)
			return
		}
		logging.Error("failed to get asset", map[string]interface{}{
			"asset_id": assetID,
			"error":    err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"asset":   asset,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetAssetsBySession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	sessionID, err := uuid.Parse(vars["sessionId"])
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// Parse query parameters
	filter := &assets.AssetFilter{
		SessionID: sessionID,
		Type:      r.URL.Query().Get("type"),
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

	assets, err := h.assetManager.GetAssetsBySession(ctx, filter)
	if err != nil {
		logging.Error("failed to get assets by session", map[string]interface{}{
			"session_id": sessionID,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"assets":  assets,
		"pagination": map[string]interface{}{
			"limit":    filter.Limit,
			"offset":   filter.Offset,
			"has_more": len(assets) >= filter.Limit,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) DeleteAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	assetID, err := uuid.Parse(vars["assetId"])
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	// Get user from context
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	err = h.assetManager.DeleteAsset(ctx, assetID, user.ID)
	if err != nil {
		if err.Error() == "asset not found" {
			http.Error(w, "Asset not found", http.StatusNotFound)
			return
		}
		if err.Error() == "unauthorized to delete asset" {
			http.Error(w, "Unauthorized to delete asset", http.StatusForbidden)
			return
		}
		logging.Error("failed to delete asset", map[string]interface{}{
			"asset_id": assetID,
			"user_id":  user.ID,
			"error":    err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetAssetContent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	assetID, err := uuid.Parse(vars["assetId"])
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	content, mimeType, err := h.assetManager.GetAssetContent(ctx, assetID)
	if err != nil {
		if err.Error() == "asset not found" {
			http.Error(w, "Asset not found", http.StatusNotFound)
			return
		}
		if err.Error() == "asset file not found" {
			http.Error(w, "Asset file not found", http.StatusNotFound)
			return
		}
		logging.Error("failed to get asset content", map[string]interface{}{
			"asset_id": assetID,
			"error":    err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.Write(content)
}

func (h *Handler) TrackAssetUsage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	assetID, err := uuid.Parse(vars["assetId"])
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	var req struct {
		SessionID uuid.UUID              `json:"session_id"`
		Usage     string                 `json:"usage"`
		Context   map[string]interface{} `json:"context"`
	}

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

	err = h.assetManager.TrackAssetUsage(ctx, assetID, req.SessionID, user.ID, req.Usage, req.Context)
	if err != nil {
		logging.Error("failed to track asset usage", map[string]interface{}{
			"asset_id":   assetID,
			"session_id": req.SessionID,
			"user_id":    user.ID,
			"usage":      req.Usage,
			"error":      err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Asset usage tracked successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetAssetUsage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	assetID, err := uuid.Parse(vars["assetId"])
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	usage, err := h.assetManager.GetAssetUsage(ctx, assetID)
	if err != nil {
		logging.Error("failed to get asset usage", map[string]interface{}{
			"asset_id": assetID,
			"error":    err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"usage":   usage,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}