package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"holodeck1/api/auth"
	"holodeck1/api/services"
	"holodeck1/api/sessions"
	authPkg "holodeck1/auth"
	"holodeck1/service"
	"holodeck1/session"
)

func SetupFoundationRoutes(r *mux.Router, sessionManager *session.Manager, serviceRegistry *service.Registry, authManager *authPkg.Manager) {
	// Initialize handlers
	sessionHandler := sessions.NewHandler(sessionManager)
	serviceHandler := services.NewHandler(serviceRegistry)
	authHandler := auth.NewHandler(authManager)
	
	// Initialize middleware
	authMiddleware := authPkg.NewMiddleware(authManager)
	
	// API base path
	api := r.PathPrefix("/api").Subrouter()
	
	// Authentication routes (no auth required)
	authRoutes := api.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/register", authHandler.Register).Methods("POST")
	authRoutes.HandleFunc("/login", authHandler.Login).Methods("POST")
	authRoutes.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")
	
	// Protected auth routes
	protectedAuthRoutes := api.PathPrefix("/auth").Subrouter()
	protectedAuthRoutes.Use(authMiddleware.AuthRequired)
	protectedAuthRoutes.HandleFunc("/logout", authHandler.Logout).Methods("POST")
	protectedAuthRoutes.HandleFunc("/me", authHandler.GetCurrentUser).Methods("GET")
	
	// Session routes (require authentication)
	sessionRoutes := api.PathPrefix("/sessions").Subrouter()
	sessionRoutes.Use(authMiddleware.AuthRequired)
	sessionRoutes.HandleFunc("", sessionHandler.ListSessions).Methods("GET")
	sessionRoutes.HandleFunc("", sessionHandler.CreateSession).Methods("POST")
	sessionRoutes.HandleFunc("/{sessionId}", sessionHandler.GetSession).Methods("GET")
	sessionRoutes.HandleFunc("/{sessionId}", sessionHandler.UpdateSession).Methods("PUT")
	sessionRoutes.HandleFunc("/{sessionId}", sessionHandler.DeleteSession).Methods("DELETE")
	sessionRoutes.HandleFunc("/{sessionId}/join", sessionHandler.JoinSession).Methods("POST")
	sessionRoutes.HandleFunc("/{sessionId}/leave", sessionHandler.LeaveSession).Methods("POST")
	sessionRoutes.HandleFunc("/{sessionId}/participants", sessionHandler.GetSessionParticipants).Methods("GET")
	
	// Service routes (require authentication)
	serviceRoutes := api.PathPrefix("/services").Subrouter()
	serviceRoutes.Use(authMiddleware.AuthRequired)
	serviceRoutes.HandleFunc("", serviceHandler.ListServices).Methods("GET")
	serviceRoutes.HandleFunc("", serviceHandler.RegisterService).Methods("POST")
	serviceRoutes.HandleFunc("/{serviceId}", serviceHandler.GetService).Methods("GET")
	serviceRoutes.HandleFunc("/{serviceId}", serviceHandler.UpdateService).Methods("PUT")
	serviceRoutes.HandleFunc("/{serviceId}", serviceHandler.DeleteService).Methods("DELETE")
	serviceRoutes.HandleFunc("/{serviceId}/health", serviceHandler.CheckServiceHealth).Methods("GET")
	serviceRoutes.HandleFunc("/{serviceId}/health/history", serviceHandler.GetServiceHealthHistory).Methods("GET")
	
	// Health check endpoint (no auth required)
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	}).Methods("GET")
}