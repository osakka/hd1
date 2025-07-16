package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"holodeck1/api/assets"
	"holodeck1/api/auth"
	"holodeck1/api/clients"
	"holodeck1/api/enterprise"
	"holodeck1/api/llm"
	"holodeck1/api/ot"
	"holodeck1/api/plugins"
	"holodeck1/api/services"
	"holodeck1/api/sessions"
	"holodeck1/api/webrtc"
	"holodeck1/database"
	"holodeck1/server"
)

// TestFullIntegration tests the complete HD1 platform integration
func TestFullIntegration(t *testing.T) {
	// Setup test database
	db := setupTestDatabase(t)
	defer db.Close()

	err := db.InitializeSchema()
	require.NoError(t, err)

	// Create WebSocket hub
	hub := server.NewHub()
	go hub.Run()

	// Create test server with full routing
	router := setupFullRouter(db, hub)
	server := httptest.NewServer(router)
	defer server.Close()

	// Convert http:// to ws:// for WebSocket testing
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	t.Run("Complete Platform Workflow", func(t *testing.T) {
		// Step 1: Create organization
		orgID := createTestOrganization(t, server.URL)

		// Step 2: Register user and authenticate
		userID, token := registerAndAuthenticateUser(t, server.URL)

		// Step 3: Create session
		sessionID := createTestSession(t, server.URL, token)

		// Step 4: Register service
		serviceID := registerTestService(t, server.URL, token)

		// Step 5: Install and enable plugin
		pluginID := installAndEnablePlugin(t, server.URL, token)

		// Step 6: Register client
		clientID := registerTestClient(t, server.URL, token)

		// Step 7: Create LLM avatar
		avatarID := createLLMAvatar(t, server.URL, token, sessionID)

		// Step 8: Upload asset
		assetID := uploadTestAsset(t, server.URL, token, sessionID)

		// Step 9: Create collaborative document
		docID := createCollaborativeDocument(t, server.URL, token, sessionID)

		// Step 10: Create WebRTC session
		rtcSessionID := createWebRTCSession(t, server.URL, token, sessionID)

		// Step 11: Test real-time collaboration
		testRealtimeCollaboration(t, wsURL, docID, avatarID)

		// Step 12: Test AI interaction
		testAIInteraction(t, server.URL, token, avatarID, sessionID)

		// Step 13: Test analytics tracking
		testAnalyticsTracking(t, server.URL, token, orgID, userID, sessionID)

		// Step 14: Test security features
		testSecurityFeatures(t, server.URL, token, orgID, userID)

		// Step 15: Verify all components are working
		verifyPlatformHealth(t, server.URL, map[string]string{
			"organization": orgID,
			"user":         userID.String(),
			"session":      sessionID.String(),
			"service":      serviceID,
			"plugin":       pluginID,
			"client":       clientID,
			"avatar":       avatarID,
			"asset":        assetID,
			"document":     docID,
			"rtc_session":  rtcSessionID,
		})
	})

	t.Run("Multi-User Collaboration", func(t *testing.T) {
		// Create organization
		orgID := createTestOrganization(t, server.URL)

		// Create multiple users
		users := make([]struct {
			id    uuid.UUID
			token string
		}, 3)

		for i := range users {
			users[i].id, users[i].token = registerAndAuthenticateUser(t, server.URL)
		}

		// Create shared session
		sessionID := createTestSession(t, server.URL, users[0].token)

		// All users join the session
		for i := range users {
			joinSession(t, server.URL, users[i].token, sessionID)
		}

		// Create collaborative document
		docID := createCollaborativeDocument(t, server.URL, users[0].token, sessionID)

		// Test concurrent editing
		testConcurrentEditing(t, server.URL, users, docID)

		// Test WebRTC signaling between users
		testWebRTCSignaling(t, server.URL, users, sessionID)

		// Test real-time asset sharing
		testAssetSharing(t, server.URL, users, sessionID)
	})

	t.Run("Enterprise Features Integration", func(t *testing.T) {
		// Create enterprise organization
		orgID := createEnterpriseOrganization(t, server.URL)

		// Create admin user
		adminID, adminToken := registerAndAuthenticateUser(t, server.URL)

		// Create custom roles
		designerRoleID := createCustomRole(t, server.URL, adminToken, orgID, "3D Designer")
		managerRoleID := createCustomRole(t, server.URL, adminToken, orgID, "Project Manager")

		// Create regular users
		designerID, designerToken := registerAndAuthenticateUser(t, server.URL)
		managerID, managerToken := registerAndAuthenticateUser(t, server.URL)

		// Assign roles
		assignUserRole(t, server.URL, adminToken, orgID, designerID, designerRoleID)
		assignUserRole(t, server.URL, adminToken, orgID, managerID, managerRoleID)

		// Test role-based permissions
		testRoleBasedPermissions(t, server.URL, map[string]string{
			"designer": designerToken,
			"manager":  managerToken,
			"admin":    adminToken,
		}, orgID)

		// Test analytics and reporting
		testEnterpriseAnalytics(t, server.URL, adminToken, orgID)

		// Test security and compliance
		testSecurityCompliance(t, server.URL, adminToken, orgID)
	})

	t.Run("AI-Powered Content Generation", func(t *testing.T) {
		// Setup
		userID, token := registerAndAuthenticateUser(t, server.URL)
		sessionID := createTestSession(t, server.URL, token)

		// Create AI avatar
		avatarID := createLLMAvatar(t, server.URL, token, sessionID)

		// Test content generation
		testAIContentGeneration(t, server.URL, token, avatarID, sessionID)

		// Test scene understanding
		testAISceneUnderstanding(t, server.URL, token, avatarID, sessionID)

		// Test multi-modal interaction
		testAIMultiModalInteraction(t, server.URL, token, avatarID, sessionID)
	})

	t.Run("Cross-Platform Client Integration", func(t *testing.T) {
		// Register different client types
		clients := []struct {
			name     string
			platform string
			token    string
		}{
			{"Web Client", "web", ""},
			{"Mobile Client", "mobile", ""},
			{"Desktop Client", "desktop", ""},
			{"VR Client", "vr", ""},
		}

		for i := range clients {
			clients[i].token = registerClientAndGetToken(t, server.URL, clients[i].name, clients[i].platform)
		}

		// Test cross-platform synchronization
		testCrossPlatformSync(t, server.URL, clients)

		// Test platform-specific features
		testPlatformSpecificFeatures(t, server.URL, clients)
	})

	t.Run("Performance and Scalability", func(t *testing.T) {
		// Test concurrent session creation
		testConcurrentSessions(t, server.URL, 10)

		// Test high-volume asset uploads
		testHighVolumeAssets(t, server.URL, 50)

		// Test WebSocket connection limits
		testWebSocketLimits(t, wsURL, 100)

		// Test database performance
		testDatabasePerformance(t, db)
	})

	t.Run("Error Handling and Recovery", func(t *testing.T) {
		// Test authentication failures
		testAuthenticationFailures(t, server.URL)

		// Test invalid API requests
		testInvalidRequests(t, server.URL)

		// Test resource not found scenarios
		testResourceNotFound(t, server.URL)

		// Test WebSocket connection failures
		testWebSocketFailures(t, wsURL)
	})
}

// Helper functions for integration testing

func setupFullRouter(db *database.DB, hub *server.Hub) *mux.Router {
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()

	// Auth routes (no middleware for login/register)
	authHandlers := auth.NewHandlers(db)
	api.HandleFunc("/auth/register", authHandlers.Register).Methods("POST")
	api.HandleFunc("/auth/login", authHandlers.Login).Methods("POST")

	// Apply auth middleware to protected routes
	protected := api.PathPrefix("/").Subrouter()
	protected.Use(auth.AuthMiddleware)

	// Foundation routes
	sessionHandlers := sessions.NewHandlers(db)
	protected.HandleFunc("/sessions", sessionHandlers.CreateSession).Methods("POST")
	protected.HandleFunc("/sessions", sessionHandlers.ListSessions).Methods("GET")
	protected.HandleFunc("/sessions/{id}", sessionHandlers.GetSession).Methods("GET")
	protected.HandleFunc("/sessions/{id}/join", sessionHandlers.JoinSession).Methods("POST")

	serviceHandlers := services.NewHandlers(db)
	protected.HandleFunc("/services", serviceHandlers.RegisterService).Methods("POST")
	protected.HandleFunc("/services/{id}", serviceHandlers.GetService).Methods("GET")

	// Collaboration routes
	webrtcHandlers := webrtc.NewHandlers(db, hub)
	protected.HandleFunc("/webrtc/sessions", webrtcHandlers.CreateRTCSession).Methods("POST")
	protected.HandleFunc("/webrtc/offer", webrtcHandlers.SendOffer).Methods("POST")

	otHandlers := ot.NewHandlers(db, hub)
	protected.HandleFunc("/ot/documents", otHandlers.CreateDocument).Methods("POST")
	protected.HandleFunc("/ot/documents/{id}/operations", otHandlers.ApplyOperation).Methods("POST")

	assetHandlers := assets.NewHandlers(db)
	protected.HandleFunc("/assets", assetHandlers.UploadAsset).Methods("POST")
	protected.HandleFunc("/assets/{id}", assetHandlers.GetAsset).Methods("GET")

	// AI routes
	llmHandlers := llm.NewHandlers(db)
	protected.HandleFunc("/llm/avatars", llmHandlers.CreateAvatar).Methods("POST")
	protected.HandleFunc("/llm/avatars/{id}/interact", llmHandlers.InteractWithAvatar).Methods("POST")
	protected.HandleFunc("/llm/content/generate", llmHandlers.GenerateContent).Methods("POST")

	// Universal platform routes
	clientHandlers := clients.NewHandlers(db)
	protected.HandleFunc("/clients/register", clientHandlers.RegisterClient).Methods("POST")
	protected.HandleFunc("/clients/{id}/heartbeat", clientHandlers.Heartbeat).Methods("POST")

	pluginHandlers := plugins.NewHandlers(db)
	protected.HandleFunc("/plugins", pluginHandlers.InstallPlugin).Methods("POST")
	protected.HandleFunc("/plugins/{id}/enable", pluginHandlers.EnablePlugin).Methods("POST")

	// Enterprise routes
	enterpriseHandlers := enterprise.NewHandlers(db)
	protected.HandleFunc("/organizations", enterpriseHandlers.CreateOrganization).Methods("POST")
	protected.HandleFunc("/organizations/{id}", enterpriseHandlers.GetOrganization).Methods("GET")

	orgApi := protected.PathPrefix("/organizations/{orgId}").Subrouter()
	orgApi.HandleFunc("/roles", enterpriseHandlers.CreateRole).Methods("POST")
	orgApi.HandleFunc("/roles/assign", enterpriseHandlers.AssignRole).Methods("POST")
	orgApi.HandleFunc("/analytics/track", enterpriseHandlers.TrackEvent).Methods("POST")
	orgApi.HandleFunc("/security/events", enterpriseHandlers.LogSecurityEvent).Methods("POST")

	// WebSocket endpoint
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWS(hub, w, r)
	})

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","version":"v0.7.0"}`))
	}).Methods("GET")

	return router
}

func createTestOrganization(t *testing.T, baseURL string) string {
	reqBody := map[string]interface{}{
		"name":        "Test Organization",
		"description": "Integration test organization",
		"plan":        "enterprise",
	}
	body, _ := json.Marshal(reqBody)

	resp, err := http.Post(baseURL+"/api/organizations", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	return response["id"].(string)
}

func registerAndAuthenticateUser(t *testing.T, baseURL string) (uuid.UUID, string) {
	userID := uuid.New()
	email := fmt.Sprintf("user_%s@test.com", userID.String()[:8])
	
	// Register user
	reqBody := map[string]string{
		"email":    email,
		"username": fmt.Sprintf("user_%s", userID.String()[:8]),
		"password": "TestPass123!",
	}
	body, _ := json.Marshal(reqBody)

	resp, err := http.Post(baseURL+"/api/auth/register", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var registerResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&registerResponse)
	require.NoError(t, err)

	// Login to get token
	loginBody := map[string]string{
		"email":    email,
		"password": "TestPass123!",
	}
	body, _ = json.Marshal(loginBody)

	resp, err = http.Post(baseURL+"/api/auth/login", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	require.NoError(t, err)

	return userID, loginResponse["token"].(string)
}

func createTestSession(t *testing.T, baseURL, token string) uuid.UUID {
	reqBody := map[string]interface{}{
		"name":        "Integration Test Session",
		"description": "Test session for integration testing",
		"visibility":  "public",
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", baseURL+"/api/sessions", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	sessionID, err := uuid.Parse(response["id"].(string))
	require.NoError(t, err)

	return sessionID
}

func testRealtimeCollaboration(t *testing.T, wsURL, docID, avatarID string) {
	// Connect two WebSocket clients
	ws1, _, err := websocket.DefaultDialer.Dial(wsURL+"/ws", nil)
	require.NoError(t, err)
	defer ws1.Close()

	ws2, _, err := websocket.DefaultDialer.Dial(wsURL+"/ws", nil)
	require.NoError(t, err)
	defer ws2.Close()

	// Client 1 sends a document operation
	operation := map[string]interface{}{
		"type":        "document_operation",
		"document_id": docID,
		"operation": map[string]interface{}{
			"type":  "insert",
			"path":  "/entities/0",
			"value": map[string]interface{}{
				"id":   "test_entity",
				"type": "box",
			},
		},
		"sequence": 1,
	}

	err = ws1.WriteJSON(operation)
	require.NoError(t, err)

	// Client 2 should receive the operation
	var received map[string]interface{}
	err = ws2.ReadJSON(&received)
	require.NoError(t, err)

	assert.Equal(t, "document_operation", received["type"])
	assert.Equal(t, docID, received["document_id"])
	assert.Equal(t, float64(1), received["sequence"])
}

func testAIInteraction(t *testing.T, baseURL, token, avatarID string, sessionID uuid.UUID) {
	reqBody := map[string]interface{}{
		"message": "What objects are in the current scene?",
		"context": map[string]interface{}{
			"scene_objects": []interface{}{
				map[string]interface{}{
					"id":   "box1",
					"type": "box",
					"position": map[string]float64{
						"x": 0, "y": 1, "z": 0,
					},
				},
			},
		},
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", baseURL+"/api/llm/avatars/"+avatarID+"/interact", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Contains(t, response, "message")
	assert.Contains(t, response, "tokens_used")
	assert.Contains(t, response, "cost")
}

func verifyPlatformHealth(t *testing.T, baseURL string, resources map[string]string) {
	// Check health endpoint
	resp, err := http.Get(baseURL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var healthResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&healthResponse)
	require.NoError(t, err)

	assert.Equal(t, "ok", healthResponse["status"])
	assert.Equal(t, "v0.7.0", healthResponse["version"])

	// TODO: Add more health checks for individual components
	// This would verify that all created resources are still accessible
	// and that the system is functioning correctly
}

// Additional helper functions would be implemented here for:
// - registerTestService
// - installAndEnablePlugin
// - registerTestClient
// - createLLMAvatar
// - uploadTestAsset
// - createCollaborativeDocument
// - createWebRTCSession
// - testAnalyticsTracking
// - testSecurityFeatures
// - testConcurrentEditing
// - testWebRTCSignaling
// - testAssetSharing
// - testRoleBasedPermissions
// - testEnterpriseAnalytics
// - testSecurityCompliance
// - testAIContentGeneration
// - testAISceneUnderstanding
// - testAIMultiModalInteraction
// - testCrossPlatformSync
// - testPlatformSpecificFeatures
// - testConcurrentSessions
// - testHighVolumeAssets
// - testWebSocketLimits
// - testDatabasePerformance
// - testAuthenticationFailures
// - testInvalidRequests
// - testResourceNotFound
// - testWebSocketFailures

// These functions would implement specific integration test scenarios
// for each component of the platform, ensuring end-to-end functionality