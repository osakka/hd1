package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"holodeck1/auth"
	"holodeck1/logging"
)

type Handler struct {
	authManager *auth.Manager
}

func NewHandler(authManager *auth.Manager) *Handler {
	return &Handler{
		authManager: authManager,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var req auth.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate required fields
	if req.Email == "" || req.Username == "" || req.Password == "" {
		http.Error(w, "Email, username, and password are required", http.StatusBadRequest)
		return
	}
	
	user, err := h.authManager.Register(ctx, &req)
	if err != nil {
		if err.Error() == "user already exists" {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		logging.Error("failed to register user", map[string]interface{}{
			"email": req.Email,
			"error": err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success": true,
		"user":    user,
		"message": "User registered successfully",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate required fields
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}
	
	user, token, err := h.authManager.Login(ctx, &req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		logging.Error("failed to login user", map[string]interface{}{
			"email": req.Email,
			"error": err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success":       true,
		"user":          user,
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"token_type":    token.TokenType,
		"expires_at":    token.ExpiresAt,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var req auth.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate required fields
	if req.RefreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}
	
	token, err := h.authManager.RefreshToken(ctx, &req)
	if err != nil {
		if err.Error() == "invalid refresh token" || err.Error() == "refresh token expired" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		logging.Error("failed to refresh token", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success":       true,
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"token_type":    token.TokenType,
		"expires_at":    token.ExpiresAt,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Extract access token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header required", http.StatusBadRequest)
		return
	}
	
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Invalid authorization header format", http.StatusBadRequest)
		return
	}
	
	accessToken := parts[1]
	
	err := h.authManager.Logout(ctx, accessToken)
	if err != nil {
		if err.Error() == "token not found" {
			http.Error(w, "Token not found", http.StatusNotFound)
			return
		}
		logging.Error("failed to logout user", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get user from context (set by auth middleware)
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}
	
	response := map[string]interface{}{
		"success": true,
		"user":    user,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}