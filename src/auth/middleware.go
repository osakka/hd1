package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"holodeck1/logging"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

type Middleware struct {
	authManager *Manager
}

func NewMiddleware(authManager *Manager) *Middleware {
	return &Middleware{
		authManager: authManager,
	}
}

func (m *Middleware) AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := m.extractToken(r)
		if token == "" {
			m.respondWithError(w, http.StatusUnauthorized, "Authorization token required")
			return
		}

		user, err := m.authManager.ValidateToken(r.Context(), token)
		if err != nil {
			logging.Warn("token validation failed", map[string]interface{}{
				"error": err.Error(),
				"ip":    r.RemoteAddr,
			})
			m.respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) AuthOptional(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := m.extractToken(r)
		if token != "" {
			user, err := m.authManager.ValidateToken(r.Context(), token)
			if err == nil {
				// Add user to context if token is valid
				ctx := context.WithValue(r.Context(), UserContextKey, user)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) extractToken(r *http.Request) string {
	// Check Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Expected format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	// Check query parameter
	token := r.URL.Query().Get("access_token")
	if token != "" {
		return token
	}

	// Check cookie
	cookie, err := r.Cookie("hd1_token")
	if err == nil {
		return cookie.Value
	}

	return ""
}

func (m *Middleware) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := map[string]interface{}{
		"success": false,
		"error":   message,
	}
	
	json.NewEncoder(w).Encode(response)
}

func GetUserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(UserContextKey).(*User)
	return user, ok
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	user, ok := GetUserFromContext(ctx)
	if !ok {
		return uuid.Nil, false
	}
	return user.ID, true
}