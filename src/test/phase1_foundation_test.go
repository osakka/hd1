package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"holodeck1/api/auth"
	"holodeck1/api/services"
	"holodeck1/api/sessions"
	"holodeck1/database"
)

// TestPhase1Foundation tests all Phase 1 components
func TestPhase1Foundation(t *testing.T) {
	// Setup test database
	db := setupTestDatabase(t)
	defer db.Close()

	// Initialize schema
	err := db.InitializeSchema()
	require.NoError(t, err, "Failed to initialize database schema")

	// Create router
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()

	// Setup handlers
	authHandlers := auth.NewHandlers(db)
	sessionHandlers := sessions.NewHandlers(db)
	serviceHandlers := services.NewHandlers(db)

	// Register routes
	api.HandleFunc("/auth/register", authHandlers.Register).Methods("POST")
	api.HandleFunc("/auth/login", authHandlers.Login).Methods("POST")
	api.HandleFunc("/sessions", sessionHandlers.CreateSession).Methods("POST")
	api.HandleFunc("/services", serviceHandlers.RegisterService).Methods("POST")

	// Test user registration
	t.Run("User Registration", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    "test@example.com",
			"username": "testuser",
			"password": "TestPass123!",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "user")
		assert.Contains(t, response, "token")
	})

	// Test user login
	t.Run("User Login", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    "test@example.com",
			"password": "TestPass123!",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "token")
		assert.Contains(t, response, "refresh_token")
	})

	// Test session creation (requires auth)
	t.Run("Session Creation", func(t *testing.T) {
		// First login to get token
		loginBody := map[string]string{
			"email":    "test@example.com",
			"password": "TestPass123!",
		}
		body, _ := json.Marshal(loginBody)

		loginReq := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRR := httptest.NewRecorder()
		router.ServeHTTP(loginRR, loginReq)

		var loginResp map[string]interface{}
		json.Unmarshal(loginRR.Body.Bytes(), &loginResp)
		token := loginResp["token"].(string)

		// Create session with auth
		sessionBody := map[string]interface{}{
			"name":        "Test Session",
			"description": "A test 3D session",
			"visibility":  "public",
		}
		body, _ = json.Marshal(sessionBody)

		// Add auth middleware for this test
		api.Use(auth.AuthMiddleware)

		req := httptest.NewRequest("POST", "/api/sessions", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "id")
		assert.Equal(t, "Test Session", response["name"])
	})

	// Test service registration
	t.Run("Service Registration", func(t *testing.T) {
		serviceBody := map[string]interface{}{
			"name":        "Test Service",
			"description": "A test 3D rendering service",
			"type":        "renderer",
			"endpoint":    "http://localhost:3000",
			"capabilities": []string{"threejs", "webgl", "physics"},
		}
		body, _ := json.Marshal(serviceBody)

		req := httptest.NewRequest("POST", "/api/services", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "id")
		assert.Equal(t, "Test Service", response["name"])
		assert.Equal(t, "active", response["status"])
	})
}

// TestDatabaseOperations tests core database functionality
func TestDatabaseOperations(t *testing.T) {
	db := setupTestDatabase(t)
	defer db.Close()

	err := db.InitializeSchema()
	require.NoError(t, err)

	t.Run("Health Check", func(t *testing.T) {
		err := db.HealthCheck()
		assert.NoError(t, err)
	})

	t.Run("Transaction Support", func(t *testing.T) {
		tx, err := db.Begin()
		require.NoError(t, err)

		// Insert test data
		userID := uuid.New()
		_, err = tx.Exec(`
			INSERT INTO users (id, email, username, password_hash, created_at)
			VALUES ($1, $2, $3, $4, $5)
		`, userID, "tx@test.com", "txuser", "hash", time.Now())
		require.NoError(t, err)

		// Rollback
		err = tx.Rollback()
		assert.NoError(t, err)

		// Verify data not committed
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", userID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}

// TestAuthenticationFlow tests the complete auth flow
func TestAuthenticationFlow(t *testing.T) {
	db := setupTestDatabase(t)
	defer db.Close()

	err := db.InitializeSchema()
	require.NoError(t, err)

	authManager := auth.NewAuthManager(db)

	t.Run("Complete Auth Flow", func(t *testing.T) {
		// Register user
		user, err := authManager.RegisterUser("authtest@example.com", "authuser", "SecurePass123!")
		require.NoError(t, err)
		assert.NotNil(t, user)

		// Login
		token, refreshToken, err := authManager.Login("authtest@example.com", "SecurePass123!")
		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotEmpty(t, refreshToken)

		// Validate token
		claims, err := authManager.ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, user.ID.String(), claims.UserID)

		// Refresh token
		newToken, err := authManager.RefreshToken(refreshToken)
		require.NoError(t, err)
		assert.NotEmpty(t, newToken)
		assert.NotEqual(t, token, newToken)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		_, _, err := authManager.Login("authtest@example.com", "WrongPassword")
		assert.Error(t, err)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		_, err := authManager.ValidateToken("invalid.token.here")
		assert.Error(t, err)
	})
}

// TestSessionManagement tests session operations
func TestSessionManagement(t *testing.T) {
	db := setupTestDatabase(t)
	defer db.Close()

	err := db.InitializeSchema()
	require.NoError(t, err)

	sessionManager := sessions.NewSessionManager(db)
	userID := createTestUser(t, db)

	t.Run("Session Lifecycle", func(t *testing.T) {
		// Create session
		session, err := sessionManager.CreateSession(&sessions.CreateSessionRequest{
			Name:        "Test 3D World",
			Description: "A test 3D environment",
			OwnerID:     userID,
			Visibility:  "public",
		})
		require.NoError(t, err)
		assert.NotNil(t, session)

		// Get session
		retrieved, err := sessionManager.GetSession(session.ID)
		require.NoError(t, err)
		assert.Equal(t, session.Name, retrieved.Name)

		// Update session
		err = sessionManager.UpdateSession(session.ID, map[string]interface{}{
			"description": "Updated description",
			"visibility":  "private",
		})
		require.NoError(t, err)

		// List sessions
		sessions, err := sessionManager.ListSessions(&sessions.SessionFilter{
			OwnerID: userID,
			Status:  "active",
			Limit:   10,
		})
		require.NoError(t, err)
		assert.Len(t, sessions, 1)

		// Join session
		participant, err := sessionManager.JoinSession(session.ID, userID)
		require.NoError(t, err)
		assert.NotNil(t, participant)

		// Leave session
		err = sessionManager.LeaveSession(session.ID, userID)
		assert.NoError(t, err)

		// Delete session
		err = sessionManager.DeleteSession(session.ID)
		assert.NoError(t, err)
	})
}

// TestServiceRegistry tests service management
func TestServiceRegistry(t *testing.T) {
	db := setupTestDatabase(t)
	defer db.Close()

	err := db.InitializeSchema()
	require.NoError(t, err)

	serviceManager := services.NewServiceManager(db)

	t.Run("Service Operations", func(t *testing.T) {
		// Register service
		service, err := serviceManager.RegisterService(&services.ServiceRegistration{
			Name:        "3D Renderer",
			Description: "High-performance 3D rendering service",
			Type:        "renderer",
			Endpoint:    "http://renderer.local:8080",
			Capabilities: []string{"threejs", "webgl2", "physics", "particles"},
			UIMapping: map[string]interface{}{
				"panel": "renderer-controls",
				"menu":  "rendering-options",
			},
		})
		require.NoError(t, err)
		assert.NotNil(t, service)

		// Get service
		retrieved, err := serviceManager.GetService(service.ID)
		require.NoError(t, err)
		assert.Equal(t, service.Name, retrieved.Name)

		// Update service
		err = serviceManager.UpdateService(service.ID, map[string]interface{}{
			"status":   "maintenance",
			"endpoint": "http://renderer.local:8081",
		})
		require.NoError(t, err)

		// Health check
		health, err := serviceManager.CheckHealth(service.ID)
		require.NoError(t, err)
		assert.NotNil(t, health)

		// List services
		services, err := serviceManager.ListServices(&services.ServiceFilter{
			Type:   "renderer",
			Status: "maintenance",
		})
		require.NoError(t, err)
		assert.Len(t, services, 1)

		// Delete service
		err = serviceManager.DeleteService(service.ID)
		assert.NoError(t, err)
	})
}

// Helper functions

func setupTestDatabase(t *testing.T) *database.DB {
	// Use test configuration
	db, err := database.NewConnection()
	require.NoError(t, err, "Failed to connect to test database")
	return db
}

func createTestUser(t *testing.T, db *database.DB) uuid.UUID {
	userID := uuid.New()
	_, err := db.Exec(`
		INSERT INTO users (id, email, username, password_hash, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, fmt.Sprintf("user_%s@test.com", userID), fmt.Sprintf("user_%s", userID), "hash", time.Now())
	require.NoError(t, err)
	return userID
}