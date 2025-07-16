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
	"holodeck1/api/clients"
	"holodeck1/api/enterprise"
	"holodeck1/api/plugins"
	"holodeck1/database"
	"holodeck1/server"
)

// TestPhase4UniversalPlatform tests all Phase 4 Universal Platform and Enterprise components
func TestPhase4UniversalPlatform(t *testing.T) {
	db := setupTestDatabase(t)
	defer db.Close()

	err := db.InitializeSchema()
	require.NoError(t, err)

	// Create WebSocket hub
	hub := server.NewHub()
	go hub.Run()

	// Create router
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()

	// Setup handlers
	clientHandlers := clients.NewHandlers(db)
	pluginHandlers := plugins.NewHandlers(db)
	enterpriseHandlers := enterprise.NewHandlers(db)

	// Register client routes
	api.HandleFunc("/clients/register", clientHandlers.RegisterClient).Methods("POST")
	api.HandleFunc("/clients/{id}", clientHandlers.GetClient).Methods("GET")
	api.HandleFunc("/clients/{id}", clientHandlers.UpdateClient).Methods("PUT")
	api.HandleFunc("/clients/{id}/capabilities", clientHandlers.UpdateCapabilities).Methods("PUT")
	api.HandleFunc("/clients/{id}/heartbeat", clientHandlers.Heartbeat).Methods("POST")
	api.HandleFunc("/clients/{id}/message", clientHandlers.SendMessage).Methods("POST")

	// Register plugin routes
	api.HandleFunc("/plugins", pluginHandlers.InstallPlugin).Methods("POST")
	api.HandleFunc("/plugins", pluginHandlers.ListPlugins).Methods("GET")
	api.HandleFunc("/plugins/{id}", pluginHandlers.GetPlugin).Methods("GET")
	api.HandleFunc("/plugins/{id}", pluginHandlers.UpdatePlugin).Methods("PUT")
	api.HandleFunc("/plugins/{id}", pluginHandlers.DeletePlugin).Methods("DELETE")
	api.HandleFunc("/plugins/{id}/enable", pluginHandlers.EnablePlugin).Methods("POST")
	api.HandleFunc("/plugins/{id}/disable", pluginHandlers.DisablePlugin).Methods("POST")

	// Register enterprise routes
	api.HandleFunc("/organizations", enterpriseHandlers.CreateOrganization).Methods("POST")
	api.HandleFunc("/organizations", enterpriseHandlers.ListOrganizations).Methods("GET")
	api.HandleFunc("/organizations/{id}", enterpriseHandlers.GetOrganization).Methods("GET")
	api.HandleFunc("/organizations/{id}", enterpriseHandlers.UpdateOrganization).Methods("PUT")
	api.HandleFunc("/organizations/{id}", enterpriseHandlers.DeleteOrganization).Methods("DELETE")

	// Organization-specific routes
	orgApi := api.PathPrefix("/organizations/{orgId}").Subrouter()
	orgApi.HandleFunc("/roles", enterpriseHandlers.CreateRole).Methods("POST")
	orgApi.HandleFunc("/roles/assign", enterpriseHandlers.AssignRole).Methods("POST")
	orgApi.HandleFunc("/users/{userId}/permissions", enterpriseHandlers.GetUserPermissions).Methods("GET")
	orgApi.HandleFunc("/analytics/track", enterpriseHandlers.TrackEvent).Methods("POST")
	orgApi.HandleFunc("/analytics/report", enterpriseHandlers.GetAnalyticsReport).Methods("GET")
	orgApi.HandleFunc("/security/events", enterpriseHandlers.LogSecurityEvent).Methods("POST")
	orgApi.HandleFunc("/security/audit", enterpriseHandlers.GetSecurityAuditLog).Methods("GET")
	orgApi.HandleFunc("/security/api-keys", enterpriseHandlers.CreateAPIKey).Methods("POST")
	orgApi.HandleFunc("/security/compliance", enterpriseHandlers.CreateComplianceRecord).Methods("POST")

	// Test Client Registration
	t.Run("Client Registration", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":        "HD1 Web Client",
			"platform":    "web",
			"version":     "1.0.0",
			"description": "Official HD1 web client",
			"capabilities": []string{
				"3d_rendering",
				"webgl",
				"audio",
				"video",
				"webrtc",
			},
			"metadata": map[string]interface{}{
				"user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
				"resolution": "1920x1080",
				"webgl_version": "2.0",
			},
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/clients/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "id")
		assert.Contains(t, response, "token")
		assert.Equal(t, "HD1 Web Client", response["name"])
		assert.Equal(t, "active", response["status"])
	})

	// Test Client Management
	t.Run("Client Management", func(t *testing.T) {
		// Register client
		reqBody := map[string]interface{}{
			"name":     "Test Mobile Client",
			"platform": "mobile",
			"version":  "1.0.0",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/clients/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		require.Equal(t, http.StatusCreated, rr.Code)

		var clientResponse map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &clientResponse)
		clientID := clientResponse["id"].(string)

		// Update capabilities
		capabilitiesBody := map[string]interface{}{
			"capabilities": []string{
				"3d_rendering",
				"touch_input",
				"accelerometer",
				"camera",
			},
		}
		body, _ = json.Marshal(capabilitiesBody)

		req = httptest.NewRequest("PUT", fmt.Sprintf("/api/clients/%s/capabilities", clientID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		// Send heartbeat
		heartbeatBody := map[string]interface{}{
			"status": "active",
			"metrics": map[string]interface{}{
				"cpu_usage":    45.2,
				"memory_usage": 67.8,
				"fps":          60,
			},
		}
		body, _ = json.Marshal(heartbeatBody)

		req = httptest.NewRequest("POST", fmt.Sprintf("/api/clients/%s/heartbeat", clientID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		// Send message
		messageBody := map[string]interface{}{
			"type": "scene_update",
			"data": map[string]interface{}{
				"entity_id": "box1",
				"position": map[string]float64{
					"x": 1, "y": 2, "z": 3,
				},
			},
		}
		body, _ = json.Marshal(messageBody)

		req = httptest.NewRequest("POST", fmt.Sprintf("/api/clients/%s/message", clientID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// Test Plugin Management
	t.Run("Plugin Management", func(t *testing.T) {
		// Install plugin
		reqBody := map[string]interface{}{
			"name":        "Physics Plugin",
			"description": "Advanced physics simulation plugin",
			"version":     "2.1.0",
			"author":      "HD1 Team",
			"repository":  "https://github.com/hd1/physics-plugin",
			"capabilities": []string{
				"physics_simulation",
				"collision_detection",
				"rigid_body_dynamics",
			},
			"configuration": map[string]interface{}{
				"gravity": -9.81,
				"timestep": 0.016,
				"max_substeps": 4,
			},
			"hooks": []string{
				"on_entity_create",
				"on_entity_update",
				"on_frame_update",
			},
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/plugins", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "id")
		assert.Equal(t, "Physics Plugin", response["name"])
		assert.Equal(t, "installed", response["status"])

		pluginID := response["id"].(string)

		// Enable plugin
		req = httptest.NewRequest("POST", fmt.Sprintf("/api/plugins/%s/enable", pluginID), nil)
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		// Get plugin details
		req = httptest.NewRequest("GET", fmt.Sprintf("/api/plugins/%s", pluginID), nil)
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var pluginDetails map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &pluginDetails)
		assert.NoError(t, err)
		assert.Equal(t, "enabled", pluginDetails["status"])

		// Disable plugin
		req = httptest.NewRequest("POST", fmt.Sprintf("/api/plugins/%s/disable", pluginID), nil)
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		// Update plugin
		updateBody := map[string]interface{}{
			"configuration": map[string]interface{}{
				"gravity": -12.0,
				"timestep": 0.020,
			},
		}
		body, _ = json.Marshal(updateBody)

		req = httptest.NewRequest("PUT", fmt.Sprintf("/api/plugins/%s", pluginID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// Test Organization Management
	t.Run("Organization Management", func(t *testing.T) {
		// Create organization
		reqBody := map[string]interface{}{
			"name":        "Acme Corporation",
			"description": "Leading provider of 3D experiences",
			"plan":        "enterprise",
			"settings": map[string]interface{}{
				"max_users":        1000,
				"max_sessions":     100,
				"storage_limit":    "1TB",
				"api_rate_limit":   10000,
				"features": []string{
					"analytics",
					"audit_logging",
					"custom_branding",
					"sso",
				},
			},
			"billing": map[string]interface{}{
				"email": "billing@acme.com",
				"plan":  "enterprise",
				"cycle": "monthly",
			},
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/organizations", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "id")
		assert.Equal(t, "Acme Corporation", response["name"])
		assert.Equal(t, "active", response["status"])

		orgID := response["id"].(string)

		// Update organization
		updateBody := map[string]interface{}{
			"description": "Updated description",
			"settings": map[string]interface{}{
				"max_users": 2000,
			},
		}
		body, _ = json.Marshal(updateBody)

		req = httptest.NewRequest("PUT", fmt.Sprintf("/api/organizations/%s", orgID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		// Get organization
		req = httptest.NewRequest("GET", fmt.Sprintf("/api/organizations/%s", orgID), nil)
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var orgDetails map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &orgDetails)
		assert.NoError(t, err)
		assert.Equal(t, "Updated description", orgDetails["description"])
	})

	// Test RBAC System
	t.Run("RBAC System", func(t *testing.T) {
		// Create organization first
		orgBody := map[string]interface{}{
			"name": "Test Organization",
			"plan": "enterprise",
		}
		body, _ := json.Marshal(orgBody)

		req := httptest.NewRequest("POST", "/api/organizations", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		require.Equal(t, http.StatusCreated, rr.Code)

		var orgResponse map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &orgResponse)
		orgID := orgResponse["id"].(string)

		// Create custom role
		roleBody := map[string]interface{}{
			"name":        "3D Designer",
			"description": "Can create and modify 3D content",
			"permissions": []string{
				"session.create",
				"session.update",
				"asset.upload",
				"asset.download",
				"object.create",
				"object.update",
				"object.delete",
			},
			"restrictions": map[string]interface{}{
				"max_sessions": 10,
				"max_assets":   100,
				"storage_limit": "10GB",
			},
		}
		body, _ = json.Marshal(roleBody)

		req = httptest.NewRequest("POST", fmt.Sprintf("/api/organizations/%s/roles", orgID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)

		var roleResponse map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &roleResponse)
		assert.NoError(t, err)
		assert.Contains(t, roleResponse, "id")
		assert.Equal(t, "3D Designer", roleResponse["name"])

		roleID := roleResponse["id"].(string)
		userID := createTestUser(t, db)

		// Assign role to user
		assignBody := map[string]interface{}{
			"user_id": userID.String(),
			"role_id": roleID,
			"expires_at": time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
		}
		body, _ = json.Marshal(assignBody)

		req = httptest.NewRequest("POST", fmt.Sprintf("/api/organizations/%s/roles/assign", orgID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		// Get user permissions
		req = httptest.NewRequest("GET", fmt.Sprintf("/api/organizations/%s/users/%s/permissions", orgID, userID.String()), nil)
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var permissionsResponse map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &permissionsResponse)
		assert.NoError(t, err)
		assert.Contains(t, permissionsResponse, "permissions")
		
		permissions := permissionsResponse["permissions"].([]interface{})
		assert.Contains(t, permissions, "session.create")
		assert.Contains(t, permissions, "object.create")
	})

	// Test Analytics System
	t.Run("Analytics System", func(t *testing.T) {
		// Create organization first
		orgBody := map[string]interface{}{
			"name": "Analytics Test Org",
			"plan": "enterprise",
		}
		body, _ := json.Marshal(orgBody)

		req := httptest.NewRequest("POST", "/api/organizations", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		require.Equal(t, http.StatusCreated, rr.Code)

		var orgResponse map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &orgResponse)
		orgID := orgResponse["id"].(string)

		userID := createTestUser(t, db)

		// Track events
		events := []map[string]interface{}{
			{
				"event":      "session.created",
				"user_id":    userID.String(),
				"session_id": uuid.New().String(),
				"properties": map[string]interface{}{
					"session_type": "collaboration",
					"duration":     1800,
				},
				"timestamp": time.Now().Format(time.RFC3339),
			},
			{
				"event":      "object.created",
				"user_id":    userID.String(),
				"properties": map[string]interface{}{
					"object_type": "box",
					"size":        "large",
				},
				"timestamp": time.Now().Format(time.RFC3339),
			},
			{
				"event":      "asset.uploaded",
				"user_id":    userID.String(),
				"properties": map[string]interface{}{
					"asset_type": "model",
					"file_size":  1024000,
				},
				"timestamp": time.Now().Format(time.RFC3339),
			},
		}

		for _, event := range events {
			body, _ := json.Marshal(event)

			req = httptest.NewRequest("POST", fmt.Sprintf("/api/organizations/%s/analytics/track", orgID), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rr = httptest.NewRecorder()

			router.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusOK, rr.Code)
		}

		// Get analytics report
		reportBody := map[string]interface{}{
			"start_date": time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
			"end_date":   time.Now().Format(time.RFC3339),
			"metrics": []string{
				"user_activity",
				"session_stats",
				"asset_usage",
				"performance_metrics",
			},
			"filters": map[string]interface{}{
				"user_id": userID.String(),
			},
		}
		body, _ = json.Marshal(reportBody)

		req = httptest.NewRequest("GET", fmt.Sprintf("/api/organizations/%s/analytics/report", orgID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var reportResponse map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &reportResponse)
		assert.NoError(t, err)
		assert.Contains(t, reportResponse, "metrics")
		assert.Contains(t, reportResponse, "summary")
		assert.Contains(t, reportResponse, "period")
	})

	// Test Security and Compliance
	t.Run("Security and Compliance", func(t *testing.T) {
		// Create organization first
		orgBody := map[string]interface{}{
			"name": "Security Test Org",
			"plan": "enterprise",
		}
		body, _ := json.Marshal(orgBody)

		req := httptest.NewRequest("POST", "/api/organizations", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		require.Equal(t, http.StatusCreated, rr.Code)

		var orgResponse map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &orgResponse)
		orgID := orgResponse["id"].(string)

		userID := createTestUser(t, db)

		// Log security event
		securityEventBody := map[string]interface{}{
			"event_type":   "login_attempt",
			"user_id":      userID.String(),
			"ip_address":   "192.168.1.100",
			"user_agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"success":      true,
			"risk_score":   0.2,
			"details": map[string]interface{}{
				"location":      "New York, NY",
				"device_type":   "desktop",
				"authentication": "password",
			},
			"timestamp": time.Now().Format(time.RFC3339),
		}
		body, _ = json.Marshal(securityEventBody)

		req = httptest.NewRequest("POST", fmt.Sprintf("/api/organizations/%s/security/events", orgID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)

		// Create API key
		apiKeyBody := map[string]interface{}{
			"name":        "Integration API Key",
			"description": "API key for third-party integration",
			"permissions": []string{
				"session.read",
				"asset.read",
				"user.read",
			},
			"rate_limit": 1000,
			"expires_at": time.Now().Add(365 * 24 * time.Hour).Format(time.RFC3339),
		}
		body, _ = json.Marshal(apiKeyBody)

		req = httptest.NewRequest("POST", fmt.Sprintf("/api/organizations/%s/security/api-keys", orgID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)

		var apiKeyResponse map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &apiKeyResponse)
		assert.NoError(t, err)
		assert.Contains(t, apiKeyResponse, "id")
		assert.Contains(t, apiKeyResponse, "key")
		assert.Equal(t, "Integration API Key", apiKeyResponse["name"])

		// Create compliance record
		complianceBody := map[string]interface{}{
			"type":        "gdpr_data_export",
			"user_id":     userID.String(),
			"description": "User requested data export under GDPR",
			"status":      "completed",
			"data": map[string]interface{}{
				"export_size":    "15MB",
				"export_format":  "json",
				"files_included": []string{
					"user_profile.json",
					"sessions.json",
					"assets.json",
				},
			},
			"compliance_officer": "compliance@example.com",
			"timestamp":          time.Now().Format(time.RFC3339),
		}
		body, _ = json.Marshal(complianceBody)

		req = httptest.NewRequest("POST", fmt.Sprintf("/api/organizations/%s/security/compliance", orgID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)

		// Get security audit log
		req = httptest.NewRequest("GET", fmt.Sprintf("/api/organizations/%s/security/audit?start_date=%s&end_date=%s",
			orgID,
			time.Now().Add(-24*time.Hour).Format(time.RFC3339),
			time.Now().Format(time.RFC3339),
		), nil)
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var auditResponse map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &auditResponse)
		assert.NoError(t, err)
		assert.Contains(t, auditResponse, "events")
		assert.Contains(t, auditResponse, "summary")
	})
}

// TestClientManager tests the client management system
func TestClientManager(t *testing.T) {
	db := setupTestDatabase(t)
	defer db.Close()

	err := db.InitializeSchema()
	require.NoError(t, err)

	clientManager := clients.NewClientManager(db)

	t.Run("Client Lifecycle", func(t *testing.T) {
		// Register client
		client, err := clientManager.RegisterClient(&clients.ClientRegistration{
			Name:        "Test Desktop Client",
			Platform:    "desktop",
			Version:     "1.0.0",
			Description: "Test desktop client application",
			Capabilities: []string{
				"3d_rendering",
				"opengl",
				"audio",
				"file_system",
			},
			Metadata: map[string]interface{}{
				"os":           "Windows 10",
				"architecture": "x64",
				"cpu_cores":    8,
				"memory_gb":    16,
			},
		})
		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, "Test Desktop Client", client.Name)
		assert.Equal(t, "active", client.Status)
		assert.NotEmpty(t, client.Token)

		// Get client
		retrieved, err := clientManager.GetClient(client.ID)
		require.NoError(t, err)
		assert.Equal(t, client.Name, retrieved.Name)
		assert.Equal(t, client.Platform, retrieved.Platform)

		// Update client
		err = clientManager.UpdateClient(client.ID, map[string]interface{}{
			"description": "Updated description",
			"version":     "1.1.0",
		})
		require.NoError(t, err)

		// Update capabilities
		err = clientManager.UpdateCapabilities(client.ID, []string{
			"3d_rendering",
			"opengl",
			"vulkan",
			"audio",
			"file_system",
			"gamepad",
		})
		require.NoError(t, err)

		// Send heartbeat
		err = clientManager.Heartbeat(client.ID, &clients.HeartbeatData{
			Status: "active",
			Metrics: map[string]interface{}{
				"cpu_usage":     25.5,
				"memory_usage":  8.2,
				"gpu_usage":     45.0,
				"fps":           60,
				"render_time":   16.67,
			},
		})
		require.NoError(t, err)

		// List clients
		clientsList, err := clientManager.ListClients(&clients.ClientFilter{
			Platform: "desktop",
			Status:   "active",
		})
		require.NoError(t, err)
		assert.Len(t, clientsList, 1)
		assert.Equal(t, "Updated description", clientsList[0].Description)

		// Deactivate client
		err = clientManager.DeactivateClient(client.ID)
		assert.NoError(t, err)
	})

	t.Run("Client Messaging", func(t *testing.T) {
		// Register client
		client, err := clientManager.RegisterClient(&clients.ClientRegistration{
			Name:     "Messaging Test Client",
			Platform: "web",
			Version:  "1.0.0",
		})
		require.NoError(t, err)

		// Send message
		err = clientManager.SendMessage(client.ID, &clients.Message{
			Type: "scene_update",
			Data: map[string]interface{}{
				"entities": []interface{}{
					map[string]interface{}{
						"id":   "box1",
						"type": "box",
						"position": map[string]float64{
							"x": 1, "y": 2, "z": 3,
						},
					},
				},
			},
		})
		require.NoError(t, err)

		// Get messages
		messages, err := clientManager.GetMessages(client.ID, 10)
		require.NoError(t, err)
		assert.Len(t, messages, 1)
		assert.Equal(t, "scene_update", messages[0].Type)
	})

	t.Run("Client Synchronization", func(t *testing.T) {
		// Register multiple clients
		client1, err := clientManager.RegisterClient(&clients.ClientRegistration{
			Name:     "Sync Client 1",
			Platform: "web",
			Version:  "1.0.0",
		})
		require.NoError(t, err)

		client2, err := clientManager.RegisterClient(&clients.ClientRegistration{
			Name:     "Sync Client 2",
			Platform: "mobile",
			Version:  "1.0.0",
		})
		require.NoError(t, err)

		// Broadcast message to all clients
		err = clientManager.BroadcastMessage(&clients.Message{
			Type: "global_announcement",
			Data: map[string]interface{}{
				"message": "Server maintenance in 5 minutes",
				"priority": "high",
			},
		})
		require.NoError(t, err)

		// Check that both clients received the message
		messages1, err := clientManager.GetMessages(client1.ID, 10)
		require.NoError(t, err)
		assert.Len(t, messages1, 1)

		messages2, err := clientManager.GetMessages(client2.ID, 10)
		require.NoError(t, err)
		assert.Len(t, messages2, 1)
	})
}

// TestPluginManager tests the plugin management system
func TestPluginManager(t *testing.T) {
	db := setupTestDatabase(t)
	defer db.Close()

	err := db.InitializeSchema()
	require.NoError(t, err)

	pluginManager := plugins.NewPluginManager(db)

	t.Run("Plugin Lifecycle", func(t *testing.T) {
		// Install plugin
		plugin, err := pluginManager.InstallPlugin(&plugins.PluginInstallation{
			Name:        "Audio Plugin",
			Description: "3D spatial audio processing",
			Version:     "1.0.0",
			Author:      "HD1 Team",
			Repository:  "https://github.com/hd1/audio-plugin",
			Capabilities: []string{
				"spatial_audio",
				"audio_effects",
				"sound_synthesis",
			},
			Configuration: map[string]interface{}{
				"max_channels":     32,
				"sample_rate":      44100,
				"buffer_size":      512,
				"doppler_enabled":  true,
			},
			Hooks: []string{
				"on_audio_play",
				"on_audio_stop",
				"on_listener_move",
			},
		})
		require.NoError(t, err)
		assert.NotNil(t, plugin)
		assert.Equal(t, "Audio Plugin", plugin.Name)
		assert.Equal(t, "installed", plugin.Status)

		// Get plugin
		retrieved, err := pluginManager.GetPlugin(plugin.ID)
		require.NoError(t, err)
		assert.Equal(t, plugin.Name, retrieved.Name)
		assert.Equal(t, plugin.Version, retrieved.Version)

		// Update plugin
		err = pluginManager.UpdatePlugin(plugin.ID, map[string]interface{}{
			"description": "Updated audio plugin description",
			"configuration": map[string]interface{}{
				"max_channels": 64,
				"sample_rate":  48000,
			},
		})
		require.NoError(t, err)

		// Enable plugin
		err = pluginManager.EnablePlugin(plugin.ID)
		require.NoError(t, err)

		// Verify plugin is enabled
		enabled, err := pluginManager.GetPlugin(plugin.ID)
		require.NoError(t, err)
		assert.Equal(t, "enabled", enabled.Status)

		// Disable plugin
		err = pluginManager.DisablePlugin(plugin.ID)
		require.NoError(t, err)

		// List plugins
		pluginsList, err := pluginManager.ListPlugins(&plugins.PluginFilter{
			Status: "disabled",
		})
		require.NoError(t, err)
		assert.Len(t, pluginsList, 1)
		assert.Equal(t, "disabled", pluginsList[0].Status)

		// Uninstall plugin
		err = pluginManager.UninstallPlugin(plugin.ID)
		assert.NoError(t, err)
	})

	t.Run("Plugin Hooks", func(t *testing.T) {
		// Install plugin with hooks
		plugin, err := pluginManager.InstallPlugin(&plugins.PluginInstallation{
			Name:         "Hook Test Plugin",
			Description:  "Plugin for testing hook system",
			Version:      "1.0.0",
			Author:       "Test Author",
			Capabilities: []string{"testing"},
			Hooks: []string{
				"on_entity_create",
				"on_entity_update",
				"on_entity_delete",
			},
		})
		require.NoError(t, err)

		// Enable plugin
		err = pluginManager.EnablePlugin(plugin.ID)
		require.NoError(t, err)

		// Test hook execution
		hookData := map[string]interface{}{
			"entity_id": "test_entity",
			"type":      "box",
			"position": map[string]float64{
				"x": 0, "y": 0, "z": 0,
			},
		}

		err = pluginManager.ExecuteHook("on_entity_create", hookData)
		require.NoError(t, err)

		// Get hook execution history
		history, err := pluginManager.GetHookHistory(plugin.ID, 10)
		require.NoError(t, err)
		assert.Len(t, history, 1)
		assert.Equal(t, "on_entity_create", history[0].Hook)
	})
}

// TestEnterpriseManager tests the enterprise management system
func TestEnterpriseManager(t *testing.T) {
	db := setupTestDatabase(t)
	defer db.Close()

	err := db.InitializeSchema()
	require.NoError(t, err)

	enterpriseManager := enterprise.NewEnterpriseManager(db)

	t.Run("Organization Management", func(t *testing.T) {
		// Create organization
		org, err := enterpriseManager.CreateOrganization(&enterprise.OrganizationRequest{
			Name:        "Test Enterprise",
			Description: "Test enterprise organization",
			Plan:        "enterprise",
			Settings: map[string]interface{}{
				"max_users":     500,
				"max_sessions":  50,
				"storage_limit": "500GB",
			},
		})
		require.NoError(t, err)
		assert.NotNil(t, org)
		assert.Equal(t, "Test Enterprise", org.Name)
		assert.Equal(t, "active", org.Status)

		// Get organization
		retrieved, err := enterpriseManager.GetOrganization(org.ID)
		require.NoError(t, err)
		assert.Equal(t, org.Name, retrieved.Name)

		// Update organization
		err = enterpriseManager.UpdateOrganization(org.ID, map[string]interface{}{
			"description": "Updated enterprise description",
			"settings": map[string]interface{}{
				"max_users": 1000,
			},
		})
		require.NoError(t, err)

		// List organizations
		orgs, err := enterpriseManager.ListOrganizations(&enterprise.OrganizationFilter{
			Plan:   "enterprise",
			Status: "active",
		})
		require.NoError(t, err)
		assert.Len(t, orgs, 1)
		assert.Equal(t, "Updated enterprise description", orgs[0].Description)

		// Deactivate organization
		err = enterpriseManager.DeactivateOrganization(org.ID)
		assert.NoError(t, err)
	})

	t.Run("RBAC System", func(t *testing.T) {
		// Create organization
		org, err := enterpriseManager.CreateOrganization(&enterprise.OrganizationRequest{
			Name: "RBAC Test Org",
			Plan: "enterprise",
		})
		require.NoError(t, err)

		// Create role
		role, err := enterpriseManager.CreateRole(org.ID, &enterprise.RoleRequest{
			Name:        "Content Manager",
			Description: "Manages 3D content and assets",
			Permissions: []string{
				"asset.create",
				"asset.update",
				"asset.delete",
				"session.create",
				"session.update",
			},
		})
		require.NoError(t, err)
		assert.NotNil(t, role)
		assert.Equal(t, "Content Manager", role.Name)

		userID := createTestUser(t, db)

		// Assign role
		err = enterpriseManager.AssignRole(org.ID, userID, role.ID, nil)
		require.NoError(t, err)

		// Get user permissions
		permissions, err := enterpriseManager.GetUserPermissions(org.ID, userID)
		require.NoError(t, err)
		assert.Contains(t, permissions, "asset.create")
		assert.Contains(t, permissions, "session.create")

		// Check permission
		hasPermission, err := enterpriseManager.CheckPermission(org.ID, userID, "asset.create")
		require.NoError(t, err)
		assert.True(t, hasPermission)

		// Revoke role
		err = enterpriseManager.RevokeRole(org.ID, userID, role.ID)
		assert.NoError(t, err)
	})

	t.Run("Analytics System", func(t *testing.T) {
		// Create organization
		org, err := enterpriseManager.CreateOrganization(&enterprise.OrganizationRequest{
			Name: "Analytics Test Org",
			Plan: "enterprise",
		})
		require.NoError(t, err)

		userID := createTestUser(t, db)

		// Track events
		events := []enterprise.AnalyticsEvent{
			{
				Event:     "user.login",
				UserID:    userID,
				Timestamp: time.Now(),
				Properties: map[string]interface{}{
					"ip_address": "192.168.1.100",
					"user_agent": "Mozilla/5.0...",
				},
			},
			{
				Event:     "session.created",
				UserID:    userID,
				Timestamp: time.Now(),
				Properties: map[string]interface{}{
					"session_type": "collaboration",
					"duration":     3600,
				},
			},
		}

		for _, event := range events {
			err = enterpriseManager.TrackEvent(org.ID, &event)
			require.NoError(t, err)
		}

		// Get analytics report
		report, err := enterpriseManager.GetAnalyticsReport(org.ID, &enterprise.AnalyticsQuery{
			StartDate: time.Now().Add(-24 * time.Hour),
			EndDate:   time.Now(),
			Metrics:   []string{"user_activity", "session_stats"},
		})
		require.NoError(t, err)
		assert.NotNil(t, report)
		assert.Contains(t, report.Metrics, "user_activity")
		assert.Contains(t, report.Metrics, "session_stats")
	})
}