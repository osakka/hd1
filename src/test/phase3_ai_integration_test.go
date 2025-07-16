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
	"holodeck1/api/llm"
	"holodeck1/database"
)

// TestPhase3AIIntegration tests all Phase 3 AI components
func TestPhase3AIIntegration(t *testing.T) {
	db := setupTestDatabase(t)
	defer db.Close()

	err := db.InitializeSchema()
	require.NoError(t, err)

	// Create router
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()

	// Setup handlers
	llmHandlers := llm.NewHandlers(db)

	// Register routes
	api.HandleFunc("/llm/avatars", llmHandlers.CreateAvatar).Methods("POST")
	api.HandleFunc("/llm/avatars", llmHandlers.ListAvatars).Methods("GET")
	api.HandleFunc("/llm/avatars/{id}", llmHandlers.GetAvatar).Methods("GET")
	api.HandleFunc("/llm/avatars/{id}", llmHandlers.UpdateAvatar).Methods("PUT")
	api.HandleFunc("/llm/avatars/{id}", llmHandlers.DeleteAvatar).Methods("DELETE")
	api.HandleFunc("/llm/avatars/{id}/interact", llmHandlers.InteractWithAvatar).Methods("POST")
	api.HandleFunc("/llm/content/generate", llmHandlers.GenerateContent).Methods("POST")

	// Test LLM Avatar Creation
	t.Run("LLM Avatar Creation", func(t *testing.T) {
		userID := createTestUser(t, db)
		sessionID := createTestSession(t, db, userID)

		reqBody := map[string]interface{}{
			"name":        "AI Assistant",
			"description": "Helpful AI assistant avatar",
			"provider":    "openai",
			"model":       "gpt-4",
			"session_id":  sessionID,
			"personality": map[string]interface{}{
				"traits": []string{"helpful", "friendly", "knowledgeable"},
				"style":  "conversational",
			},
			"appearance": map[string]interface{}{
				"model_url": "https://example.com/avatar.glb",
				"color":     "#4A90E2",
				"scale":     1.2,
			},
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/llm/avatars", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "id")
		assert.Equal(t, "AI Assistant", response["name"])
		assert.Equal(t, "active", response["status"])
	})

	// Test Avatar Interaction
	t.Run("Avatar Interaction", func(t *testing.T) {
		userID := createTestUser(t, db)
		sessionID := createTestSession(t, db, userID)

		// Create avatar first
		avatarBody := map[string]interface{}{
			"name":       "Test Avatar",
			"provider":   "openai",
			"model":      "gpt-3.5-turbo",
			"session_id": sessionID,
		}
		body, _ := json.Marshal(avatarBody)

		req := httptest.NewRequest("POST", "/api/llm/avatars", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		require.Equal(t, http.StatusCreated, rr.Code)

		var avatarResponse map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &avatarResponse)
		avatarID := avatarResponse["id"].(string)

		// Test interaction
		interactionBody := map[string]interface{}{
			"message": "Hello, how are you?",
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
			"user_id": userID.String(),
		}
		body, _ = json.Marshal(interactionBody)

		req = httptest.NewRequest("POST", fmt.Sprintf("/api/llm/avatars/%s/interact", avatarID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "message")
		assert.Contains(t, response, "tokens_used")
		assert.Contains(t, response, "cost")
	})

	// Test Content Generation
	t.Run("Content Generation", func(t *testing.T) {
		userID := createTestUser(t, db)
		sessionID := createTestSession(t, db, userID)

		reqBody := map[string]interface{}{
			"type":       "3d_scene",
			"provider":   "openai",
			"model":      "gpt-4",
			"session_id": sessionID,
			"prompt":     "Create a futuristic city scene with flying cars",
			"parameters": map[string]interface{}{
				"style":       "cyberpunk",
				"complexity":  "medium",
				"object_count": 10,
			},
			"template": map[string]interface{}{
				"scene_type": "urban",
				"lighting":   "neon",
				"weather":    "rain",
			},
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/llm/content/generate", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "content")
		assert.Contains(t, response, "generation_id")
		assert.Contains(t, response, "tokens_used")
		assert.Contains(t, response, "cost")

		// Verify generated content structure
		content := response["content"].(map[string]interface{})
		assert.Contains(t, content, "entities")
		assert.Contains(t, content, "lighting")
		assert.Contains(t, content, "environment")
	})

	// Test Multi-Provider Support
	t.Run("Multi-Provider Support", func(t *testing.T) {
		userID := createTestUser(t, db)
		sessionID := createTestSession(t, db, userID)

		providers := []map[string]interface{}{
			{
				"provider": "openai",
				"model":    "gpt-4",
				"name":     "OpenAI Avatar",
			},
			{
				"provider": "claude",
				"model":    "claude-3-opus",
				"name":     "Claude Avatar",
			},
			{
				"provider": "gemini",
				"model":    "gemini-pro",
				"name":     "Gemini Avatar",
			},
		}

		for _, providerConfig := range providers {
			reqBody := map[string]interface{}{
				"name":       providerConfig["name"],
				"provider":   providerConfig["provider"],
				"model":      providerConfig["model"],
				"session_id": sessionID,
			}
			body, _ := json.Marshal(reqBody)

			req := httptest.NewRequest("POST", "/api/llm/avatars", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusCreated, rr.Code)

			var response map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, providerConfig["name"], response["name"])
			assert.Equal(t, providerConfig["provider"], response["provider"])
		}
	})
}

// TestLLMManager tests the LLM manager directly
func TestLLMManager(t *testing.T) {
	db := setupTestDatabase(t)
	defer db.Close()

	err := db.InitializeSchema()
	require.NoError(t, err)

	llmManager := llm.NewLLMManager(db)
	userID := createTestUser(t, db)
	sessionID := createTestSession(t, db, userID)

	t.Run("Avatar Lifecycle", func(t *testing.T) {
		// Create avatar
		avatar, err := llmManager.CreateAvatar(&llm.CreateAvatarRequest{
			Name:        "Test Avatar",
			Description: "Test avatar for lifecycle testing",
			Provider:    "openai",
			Model:       "gpt-3.5-turbo",
			SessionID:   sessionID,
			Personality: map[string]interface{}{
				"traits": []string{"helpful", "creative"},
				"style":  "formal",
			},
			Appearance: map[string]interface{}{
				"model_url": "https://example.com/avatar.glb",
				"color":     "#FF6B6B",
			},
		})
		require.NoError(t, err)
		assert.NotNil(t, avatar)
		assert.Equal(t, "Test Avatar", avatar.Name)
		assert.Equal(t, "active", avatar.Status)

		// Get avatar
		retrieved, err := llmManager.GetAvatar(avatar.ID)
		require.NoError(t, err)
		assert.Equal(t, avatar.Name, retrieved.Name)
		assert.Equal(t, avatar.Provider, retrieved.Provider)

		// Update avatar
		err = llmManager.UpdateAvatar(avatar.ID, map[string]interface{}{
			"description": "Updated description",
			"status":      "maintenance",
		})
		require.NoError(t, err)

		// List avatars
		avatars, err := llmManager.ListAvatars(&llm.AvatarFilter{
			SessionID: sessionID,
			Provider:  "openai",
			Status:    "maintenance",
		})
		require.NoError(t, err)
		assert.Len(t, avatars, 1)
		assert.Equal(t, "Updated description", avatars[0].Description)

		// Delete avatar
		err = llmManager.DeleteAvatar(avatar.ID)
		assert.NoError(t, err)

		// Verify deletion
		_, err = llmManager.GetAvatar(avatar.ID)
		assert.Error(t, err)
	})

	t.Run("Avatar Interaction", func(t *testing.T) {
		// Create avatar
		avatar, err := llmManager.CreateAvatar(&llm.CreateAvatarRequest{
			Name:      "Interactive Avatar",
			Provider:  "openai",
			Model:     "gpt-3.5-turbo",
			SessionID: sessionID,
		})
		require.NoError(t, err)

		// Test interaction
		response, err := llmManager.InteractWithAvatar(avatar.ID, &llm.InteractionRequest{
			Message: "What objects are in the scene?",
			Context: &llm.SceneContext{
				Objects: []map[string]interface{}{
					{
						"id":   "sphere1",
						"type": "sphere",
						"position": map[string]float64{
							"x": 1, "y": 2, "z": 3,
						},
						"material": map[string]interface{}{
							"color": "#00FF00",
						},
					},
				},
				Camera: map[string]interface{}{
					"position": map[string]float64{
						"x": 0, "y": 5, "z": 10,
					},
					"rotation": map[string]float64{
						"x": -0.5, "y": 0, "z": 0,
					},
				},
			},
			UserID: userID,
		})
		require.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotEmpty(t, response.Message)
		assert.Greater(t, response.TokensUsed, 0)
		assert.Greater(t, response.Cost, 0.0)

		// Test conversation history
		history, err := llmManager.GetConversationHistory(avatar.ID, userID, 10)
		require.NoError(t, err)
		assert.Len(t, history, 1)
		assert.Equal(t, "What objects are in the scene?", history[0].UserMessage)
	})

	t.Run("Content Generation", func(t *testing.T) {
		// Test scene generation
		content, err := llmManager.GenerateContent(&llm.ContentGenerationRequest{
			Type:      "3d_scene",
			Provider:  "openai",
			Model:     "gpt-4",
			SessionID: sessionID,
			Prompt:    "Create a peaceful forest scene",
			Parameters: map[string]interface{}{
				"style":        "realistic",
				"object_count": 5,
				"lighting":     "natural",
			},
			Template: map[string]interface{}{
				"scene_type": "outdoor",
				"biome":      "forest",
				"time_of_day": "morning",
			},
		})
		require.NoError(t, err)
		assert.NotNil(t, content)
		assert.NotEmpty(t, content.GenerationID)
		assert.Contains(t, content.Content, "entities")
		assert.Greater(t, content.TokensUsed, 0)

		// Test object generation
		objectContent, err := llmManager.GenerateContent(&llm.ContentGenerationRequest{
			Type:      "3d_object",
			Provider:  "openai",
			Model:     "gpt-3.5-turbo",
			SessionID: sessionID,
			Prompt:    "Create a medieval castle",
			Parameters: map[string]interface{}{
				"style":      "gothic",
				"complexity": "high",
				"size":       "large",
			},
		})
		require.NoError(t, err)
		assert.NotNil(t, objectContent)
		assert.Contains(t, objectContent.Content, "geometry")
		assert.Contains(t, objectContent.Content, "materials")
	})

	t.Run("Usage Tracking", func(t *testing.T) {
		// Create avatar
		avatar, err := llmManager.CreateAvatar(&llm.CreateAvatarRequest{
			Name:      "Usage Test Avatar",
			Provider:  "openai",
			Model:     "gpt-3.5-turbo",
			SessionID: sessionID,
		})
		require.NoError(t, err)

		// Perform multiple interactions
		for i := 0; i < 3; i++ {
			_, err := llmManager.InteractWithAvatar(avatar.ID, &llm.InteractionRequest{
				Message: fmt.Sprintf("Test message %d", i+1),
				UserID:  userID,
			})
			require.NoError(t, err)
		}

		// Get usage statistics
		usage, err := llmManager.GetUsageStats(&llm.UsageFilter{
			SessionID: sessionID,
			Provider:  "openai",
			StartTime: time.Now().Add(-1 * time.Hour),
			EndTime:   time.Now(),
		})
		require.NoError(t, err)
		assert.NotNil(t, usage)
		assert.Greater(t, usage.TotalTokens, 0)
		assert.Greater(t, usage.TotalCost, 0.0)
		assert.Equal(t, 3, usage.InteractionCount)

		// Get session usage
		sessionUsage, err := llmManager.GetSessionUsage(sessionID)
		require.NoError(t, err)
		assert.NotNil(t, sessionUsage)
		assert.Greater(t, sessionUsage.TotalTokens, 0)
		assert.Greater(t, sessionUsage.TotalCost, 0.0)
	})
}

// TestLLMProviders tests different LLM provider implementations
func TestLLMProviders(t *testing.T) {
	db := setupTestDatabase(t)
	defer db.Close()

	err := db.InitializeSchema()
	require.NoError(t, err)

	t.Run("Provider Configuration", func(t *testing.T) {
		providers := []struct {
			name   string
			config map[string]interface{}
		}{
			{
				name: "openai",
				config: map[string]interface{}{
					"api_key":     "test-key",
					"model":       "gpt-4",
					"temperature": 0.7,
					"max_tokens":  1000,
				},
			},
			{
				name: "claude",
				config: map[string]interface{}{
					"api_key":     "test-key",
					"model":       "claude-3-opus",
					"temperature": 0.8,
					"max_tokens":  2000,
				},
			},
			{
				name: "gemini",
				config: map[string]interface{}{
					"api_key":     "test-key",
					"model":       "gemini-pro",
					"temperature": 0.6,
					"max_tokens":  1500,
				},
			},
		}

		for _, provider := range providers {
			providerInstance, err := llm.NewProvider(provider.name, provider.config)
			require.NoError(t, err)
			assert.NotNil(t, providerInstance)
			assert.Equal(t, provider.name, providerInstance.Name())
		}
	})

	t.Run("Provider Capabilities", func(t *testing.T) {
		provider, err := llm.NewProvider("openai", map[string]interface{}{
			"api_key": "test-key",
			"model":   "gpt-4",
		})
		require.NoError(t, err)

		capabilities := provider.GetCapabilities()
		assert.Contains(t, capabilities, "text_generation")
		assert.Contains(t, capabilities, "conversation")
		assert.Contains(t, capabilities, "code_generation")
		assert.Contains(t, capabilities, "content_creation")
	})

	t.Run("Rate Limiting", func(t *testing.T) {
		provider, err := llm.NewProvider("openai", map[string]interface{}{
			"api_key":    "test-key",
			"model":      "gpt-3.5-turbo",
			"rate_limit": 10, // 10 requests per minute
		})
		require.NoError(t, err)

		// Test rate limiting
		for i := 0; i < 5; i++ {
			allowed := provider.CheckRateLimit()
			assert.True(t, allowed, "Request %d should be allowed", i)
		}
	})
}

// TestAIContentGeneration tests AI-powered content generation
func TestAIContentGeneration(t *testing.T) {
	db := setupTestDatabase(t)
	defer db.Close()

	err := db.InitializeSchema()
	require.NoError(t, err)

	contentGenerator := llm.NewContentGenerator(db)
	userID := createTestUser(t, db)
	sessionID := createTestSession(t, db, userID)

	t.Run("Scene Generation", func(t *testing.T) {
		scene, err := contentGenerator.GenerateScene(&llm.SceneGenerationRequest{
			Theme:     "space station",
			Style:     "futuristic",
			Complexity: "high",
			Objects:   []string{"control_panels", "windows", "machinery"},
			Lighting:  "artificial",
			Provider:  "openai",
			Model:     "gpt-4",
			SessionID: sessionID,
		})
		require.NoError(t, err)
		assert.NotNil(t, scene)
		assert.Contains(t, scene.Content, "entities")
		assert.Contains(t, scene.Content, "lighting")
		assert.Contains(t, scene.Content, "environment")
		assert.Greater(t, len(scene.Content["entities"].([]interface{})), 0)
	})

	t.Run("Object Generation", func(t *testing.T) {
		object, err := contentGenerator.GenerateObject(&llm.ObjectGenerationRequest{
			Type:        "vehicle",
			Description: "futuristic hover car",
			Style:       "sleek",
			Materials:   []string{"metal", "glass", "neon"},
			Provider:    "openai",
			Model:       "gpt-3.5-turbo",
			SessionID:   sessionID,
		})
		require.NoError(t, err)
		assert.NotNil(t, object)
		assert.Contains(t, object.Content, "geometry")
		assert.Contains(t, object.Content, "materials")
		assert.Contains(t, object.Content, "properties")
	})

	t.Run("Animation Generation", func(t *testing.T) {
		animation, err := contentGenerator.GenerateAnimation(&llm.AnimationGenerationRequest{
			ObjectID:    "test-object",
			Type:        "rotation",
			Duration:    5.0,
			Easing:      "ease-in-out",
			Loop:        true,
			Description: "slowly rotating object",
			Provider:    "openai",
			Model:       "gpt-3.5-turbo",
			SessionID:   sessionID,
		})
		require.NoError(t, err)
		assert.NotNil(t, animation)
		assert.Contains(t, animation.Content, "keyframes")
		assert.Contains(t, animation.Content, "duration")
		assert.Contains(t, animation.Content, "easing")
	})

	t.Run("Template-Based Generation", func(t *testing.T) {
		templates := []struct {
			name     string
			template map[string]interface{}
		}{
			{
				name: "office_environment",
				template: map[string]interface{}{
					"type":        "indoor",
					"objects":     []string{"desk", "chair", "computer", "plants"},
					"lighting":    "office",
					"color_scheme": "professional",
				},
			},
			{
				name: "fantasy_forest",
				template: map[string]interface{}{
					"type":        "outdoor",
					"objects":     []string{"trees", "rocks", "mushrooms", "flowers"},
					"lighting":    "natural",
					"color_scheme": "earth_tones",
				},
			},
		}

		for _, template := range templates {
			content, err := contentGenerator.GenerateFromTemplate(&llm.TemplateGenerationRequest{
				Template:  template.template,
				Variables: map[string]interface{}{
					"size":       "medium",
					"complexity": "normal",
				},
				Provider:  "openai",
				Model:     "gpt-4",
				SessionID: sessionID,
			})
			require.NoError(t, err)
			assert.NotNil(t, content)
			assert.Contains(t, content.Content, "entities")
		}
	})
}