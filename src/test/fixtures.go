package test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
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
)

// TestFixtures contains all test data and fixtures
type TestFixtures struct {
	Users         []TestUser
	Organizations []TestOrganization
	Sessions      []TestSession
	Services      []TestService
	Plugins       []TestPlugin
	Assets        []TestAsset
	Documents     []TestDocument
	Avatars       []TestAvatar
	Clients       []TestClient
	Roles         []TestRole
}

// TestUser represents a test user fixture
type TestUser struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	PasswordHash string    `json:"password_hash"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// TestOrganization represents a test organization fixture
type TestOrganization struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Plan        string                 `json:"plan"`
	Status      string                 `json:"status"`
	Settings    map[string]interface{} `json:"settings"`
	OwnerID     uuid.UUID              `json:"owner_id"`
	CreatedAt   time.Time              `json:"created_at"`
}

// TestSession represents a test session fixture
type TestSession struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	OwnerID     uuid.UUID              `json:"owner_id"`
	Visibility  string                 `json:"visibility"`
	Status      string                 `json:"status"`
	Settings    map[string]interface{} `json:"settings"`
	CreatedAt   time.Time              `json:"created_at"`
}

// TestService represents a test service fixture
type TestService struct {
	ID           uuid.UUID              `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Type         string                 `json:"type"`
	Endpoint     string                 `json:"endpoint"`
	Capabilities []string               `json:"capabilities"`
	Status       string                 `json:"status"`
	UIMapping    map[string]interface{} `json:"ui_mapping"`
	CreatedAt    time.Time              `json:"created_at"`
}

// TestPlugin represents a test plugin fixture
type TestPlugin struct {
	ID            uuid.UUID              `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Version       string                 `json:"version"`
	Author        string                 `json:"author"`
	Repository    string                 `json:"repository"`
	Capabilities  []string               `json:"capabilities"`
	Configuration map[string]interface{} `json:"configuration"`
	Hooks         []string               `json:"hooks"`
	Status        string                 `json:"status"`
	CreatedAt     time.Time              `json:"created_at"`
}

// TestAsset represents a test asset fixture
type TestAsset struct {
	ID        uuid.UUID              `json:"id"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	URL       string                 `json:"url"`
	Size      int64                  `json:"size"`
	MimeType  string                 `json:"mime_type"`
	SessionID uuid.UUID              `json:"session_id"`
	UserID    uuid.UUID              `json:"user_id"`
	Tags      []string               `json:"tags"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
}

// TestDocument represents a test document fixture
type TestDocument struct {
	ID        uuid.UUID              `json:"id"`
	Type      string                 `json:"type"`
	Content   map[string]interface{} `json:"content"`
	Version   int                    `json:"version"`
	SessionID uuid.UUID              `json:"session_id"`
	CreatedAt time.Time              `json:"created_at"`
}

// TestAvatar represents a test avatar fixture
type TestAvatar struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Provider    string                 `json:"provider"`
	Model       string                 `json:"model"`
	SessionID   uuid.UUID              `json:"session_id"`
	Personality map[string]interface{} `json:"personality"`
	Appearance  map[string]interface{} `json:"appearance"`
	Status      string                 `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
}

// TestClient represents a test client fixture
type TestClient struct {
	ID           uuid.UUID              `json:"id"`
	Name         string                 `json:"name"`
	Platform     string                 `json:"platform"`
	Version      string                 `json:"version"`
	Description  string                 `json:"description"`
	Capabilities []string               `json:"capabilities"`
	Status       string                 `json:"status"`
	Token        string                 `json:"token"`
	Metadata     map[string]interface{} `json:"metadata"`
	CreatedAt    time.Time              `json:"created_at"`
}

// TestRole represents a test role fixture
type TestRole struct {
	ID            uuid.UUID              `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Permissions   []string               `json:"permissions"`
	Restrictions  map[string]interface{} `json:"restrictions"`
	OrganizationID uuid.UUID             `json:"organization_id"`
	CreatedAt     time.Time              `json:"created_at"`
}

// GenerateTestFixtures creates comprehensive test fixtures
func GenerateTestFixtures() *TestFixtures {
	fixtures := &TestFixtures{
		Users:         generateTestUsers(),
		Organizations: generateTestOrganizations(),
		Sessions:      generateTestSessions(),
		Services:      generateTestServices(),
		Plugins:       generateTestPlugins(),
		Assets:        generateTestAssets(),
		Documents:     generateTestDocuments(),
		Avatars:       generateTestAvatars(),
		Clients:       generateTestClients(),
		Roles:         generateTestRoles(),
	}

	// Create relationships between fixtures
	linkFixtures(fixtures)

	return fixtures
}

// generateTestUsers creates test user fixtures
func generateTestUsers() []TestUser {
	users := []TestUser{
		{
			ID:           uuid.New(),
			Email:        "admin@hd1.com",
			Username:     "admin",
			Password:     "AdminPass123!",
			PasswordHash: "$2b$10$hash...",
			FirstName:    "Admin",
			LastName:     "User",
			Role:         "admin",
			CreatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Email:        "designer@hd1.com",
			Username:     "designer",
			Password:     "DesignPass123!",
			PasswordHash: "$2b$10$hash...",
			FirstName:    "Jane",
			LastName:     "Designer",
			Role:         "designer",
			CreatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Email:        "developer@hd1.com",
			Username:     "developer",
			Password:     "DevPass123!",
			PasswordHash: "$2b$10$hash...",
			FirstName:    "John",
			LastName:     "Developer",
			Role:         "developer",
			CreatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Email:        "manager@hd1.com",
			Username:     "manager",
			Password:     "ManagerPass123!",
			PasswordHash: "$2b$10$hash...",
			FirstName:    "Sarah",
			LastName:     "Manager",
			Role:         "manager",
			CreatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Email:        "viewer@hd1.com",
			Username:     "viewer",
			Password:     "ViewerPass123!",
			PasswordHash: "$2b$10$hash...",
			FirstName:    "Bob",
			LastName:     "Viewer",
			Role:         "viewer",
			CreatedAt:    time.Now(),
		},
	}

	return users
}

// generateTestOrganizations creates test organization fixtures
func generateTestOrganizations() []TestOrganization {
	organizations := []TestOrganization{
		{
			ID:          uuid.New(),
			Name:        "Acme Corporation",
			Description: "Leading provider of 3D experiences",
			Plan:        "enterprise",
			Status:      "active",
			Settings: map[string]interface{}{
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
			CreatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Startup Inc",
			Description: "Innovative startup using 3D technology",
			Plan:        "professional",
			Status:      "active",
			Settings: map[string]interface{}{
				"max_users":      50,
				"max_sessions":   10,
				"storage_limit":  "100GB",
				"api_rate_limit": 1000,
				"features": []string{
					"basic_analytics",
					"standard_support",
				},
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Educational Institute",
			Description: "Using 3D for educational purposes",
			Plan:        "education",
			Status:      "active",
			Settings: map[string]interface{}{
				"max_users":      200,
				"max_sessions":   30,
				"storage_limit":  "500GB",
				"api_rate_limit": 5000,
				"features": []string{
					"educational_tools",
					"classroom_management",
					"student_progress",
				},
			},
			CreatedAt: time.Now(),
		},
	}

	return organizations
}

// generateTestSessions creates test session fixtures
func generateTestSessions() []TestSession {
	sessions := []TestSession{
		{
			ID:          uuid.New(),
			Name:        "Product Design Review",
			Description: "Collaborative 3D product design session",
			Visibility:  "private",
			Status:      "active",
			Settings: map[string]interface{}{
				"max_participants": 10,
				"recording_enabled": true,
				"auto_save":        true,
				"theme":            "dark",
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Architecture Walkthrough",
			Description: "Virtual building walkthrough session",
			Visibility:  "public",
			Status:      "active",
			Settings: map[string]interface{}{
				"max_participants": 50,
				"recording_enabled": false,
				"auto_save":        true,
				"theme":            "light",
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Game Development Demo",
			Description: "Showcasing game development in 3D",
			Visibility:  "public",
			Status:      "active",
			Settings: map[string]interface{}{
				"max_participants": 100,
				"recording_enabled": true,
				"auto_save":        false,
				"theme":            "gaming",
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Training Simulation",
			Description: "3D training simulation environment",
			Visibility:  "private",
			Status:      "active",
			Settings: map[string]interface{}{
				"max_participants": 20,
				"recording_enabled": true,
				"auto_save":        true,
				"theme":            "professional",
			},
			CreatedAt: time.Now(),
		},
	}

	return sessions
}

// generateTestServices creates test service fixtures
func generateTestServices() []TestService {
	services := []TestService{
		{
			ID:          uuid.New(),
			Name:        "3D Renderer Service",
			Description: "High-performance 3D rendering service",
			Type:        "renderer",
			Endpoint:    "http://renderer.local:8080",
			Capabilities: []string{
				"threejs",
				"webgl2",
				"physics",
				"particles",
				"post_processing",
			},
			Status: "active",
			UIMapping: map[string]interface{}{
				"panel": "renderer-controls",
				"menu":  "rendering-options",
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Physics Engine",
			Description: "Advanced physics simulation service",
			Type:        "physics",
			Endpoint:    "http://physics.local:8081",
			Capabilities: []string{
				"rigid_body",
				"soft_body",
				"fluid_simulation",
				"collision_detection",
			},
			Status: "active",
			UIMapping: map[string]interface{}{
				"panel": "physics-controls",
				"menu":  "physics-settings",
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Audio Engine",
			Description: "3D spatial audio processing service",
			Type:        "audio",
			Endpoint:    "http://audio.local:8082",
			Capabilities: []string{
				"spatial_audio",
				"audio_effects",
				"sound_synthesis",
				"voice_chat",
			},
			Status: "active",
			UIMapping: map[string]interface{}{
				"panel": "audio-controls",
				"menu":  "audio-settings",
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "AI Assistant Service",
			Description: "AI-powered assistance and content generation",
			Type:        "ai",
			Endpoint:    "http://ai.local:8083",
			Capabilities: []string{
				"content_generation",
				"natural_language",
				"scene_understanding",
				"recommendations",
			},
			Status: "active",
			UIMapping: map[string]interface{}{
				"panel": "ai-assistant",
				"menu":  "ai-tools",
			},
			CreatedAt: time.Now(),
		},
	}

	return services
}

// generateTestPlugins creates test plugin fixtures
func generateTestPlugins() []TestPlugin {
	plugins := []TestPlugin{
		{
			ID:          uuid.New(),
			Name:        "Physics Plugin",
			Description: "Advanced physics simulation plugin",
			Version:     "2.1.0",
			Author:      "HD1 Team",
			Repository:  "https://github.com/hd1/physics-plugin",
			Capabilities: []string{
				"physics_simulation",
				"collision_detection",
				"rigid_body_dynamics",
			},
			Configuration: map[string]interface{}{
				"gravity":       -9.81,
				"timestep":      0.016,
				"max_substeps":  4,
				"solver_iterations": 8,
			},
			Hooks: []string{
				"on_entity_create",
				"on_entity_update",
				"on_frame_update",
			},
			Status:    "enabled",
			CreatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Animation Plugin",
			Description: "Advanced animation system plugin",
			Version:     "1.5.0",
			Author:      "Animation Studio",
			Repository:  "https://github.com/animation-studio/hd1-plugin",
			Capabilities: []string{
				"keyframe_animation",
				"skeletal_animation",
				"morph_targets",
				"timeline_editor",
			},
			Configuration: map[string]interface{}{
				"fps":           60,
				"interpolation": "linear",
				"compression":   true,
			},
			Hooks: []string{
				"on_animation_start",
				"on_animation_end",
				"on_keyframe_change",
			},
			Status:    "enabled",
			CreatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Lighting Plugin",
			Description: "Professional lighting system plugin",
			Version:     "3.0.0",
			Author:      "Lighting Pro",
			Repository:  "https://github.com/lighting-pro/hd1-plugin",
			Capabilities: []string{
				"dynamic_lighting",
				"shadow_mapping",
				"global_illumination",
				"hdri_environment",
			},
			Configuration: map[string]interface{}{
				"shadow_quality": "high",
				"ambient_occlusion": true,
				"bloom_enabled":     true,
			},
			Hooks: []string{
				"on_light_create",
				"on_light_update",
				"on_material_change",
			},
			Status:    "installed",
			CreatedAt: time.Now(),
		},
	}

	return plugins
}

// generateTestAssets creates test asset fixtures
func generateTestAssets() []TestAsset {
	assets := []TestAsset{
		{
			ID:       uuid.New(),
			Name:     "Spaceship Model",
			Type:     "model",
			URL:      "/assets/spaceship.glb",
			Size:     2048000,
			MimeType: "model/gltf-binary",
			Tags:     []string{"sci-fi", "vehicle", "spaceship"},
			Metadata: map[string]interface{}{
				"vertices":  15000,
				"triangles": 12000,
				"materials": 5,
				"animations": []string{"idle", "fly", "land"},
			},
			CreatedAt: time.Now(),
		},
		{
			ID:       uuid.New(),
			Name:     "Futuristic Texture",
			Type:     "texture",
			URL:      "/assets/futuristic_texture.jpg",
			Size:     512000,
			MimeType: "image/jpeg",
			Tags:     []string{"texture", "sci-fi", "metal"},
			Metadata: map[string]interface{}{
				"resolution": "2048x2048",
				"format":     "jpeg",
				"channels":   3,
			},
			CreatedAt: time.Now(),
		},
		{
			ID:       uuid.New(),
			Name:     "Ambient Sound",
			Type:     "audio",
			URL:      "/assets/ambient_sound.wav",
			Size:     1024000,
			MimeType: "audio/wav",
			Tags:     []string{"audio", "ambient", "background"},
			Metadata: map[string]interface{}{
				"duration":    120.5,
				"sample_rate": 44100,
				"channels":    2,
				"bitrate":     320,
			},
			CreatedAt: time.Now(),
		},
		{
			ID:       uuid.New(),
			Name:     "Environment HDR",
			Type:     "environment",
			URL:      "/assets/environment.hdr",
			Size:     4096000,
			MimeType: "image/hdr",
			Tags:     []string{"environment", "lighting", "hdri"},
			Metadata: map[string]interface{}{
				"resolution": "4096x2048",
				"format":     "hdr",
				"exposure":   1.0,
			},
			CreatedAt: time.Now(),
		},
	}

	return assets
}

// generateTestDocuments creates test document fixtures
func generateTestDocuments() []TestDocument {
	documents := []TestDocument{
		{
			ID:   uuid.New(),
			Type: "scene",
			Content: map[string]interface{}{
				"entities": []interface{}{
					map[string]interface{}{
						"id":   "box1",
						"type": "box",
						"position": map[string]float64{
							"x": 0, "y": 1, "z": 0,
						},
						"rotation": map[string]float64{
							"x": 0, "y": 0, "z": 0,
						},
						"scale": map[string]float64{
							"x": 1, "y": 1, "z": 1,
						},
						"material": map[string]interface{}{
							"color": "#FF0000",
							"type":  "phong",
						},
					},
				},
				"lights": []interface{}{
					map[string]interface{}{
						"id":        "directional1",
						"type":      "directional",
						"intensity": 1.0,
						"color":     "#FFFFFF",
						"position": map[string]float64{
							"x": 5, "y": 5, "z": 5,
						},
					},
				},
				"camera": map[string]interface{}{
					"position": map[string]float64{
						"x": 0, "y": 5, "z": 10,
					},
					"rotation": map[string]float64{
						"x": -0.5, "y": 0, "z": 0,
					},
					"fov": 75,
				},
			},
			Version:   1,
			CreatedAt: time.Now(),
		},
		{
			ID:   uuid.New(),
			Type: "animation",
			Content: map[string]interface{}{
				"timeline": map[string]interface{}{
					"duration": 10.0,
					"fps":      60,
				},
				"keyframes": []interface{}{
					map[string]interface{}{
						"time":      0.0,
						"entity_id": "box1",
						"property":  "position",
						"value": map[string]float64{
							"x": 0, "y": 0, "z": 0,
						},
					},
					map[string]interface{}{
						"time":      5.0,
						"entity_id": "box1",
						"property":  "position",
						"value": map[string]float64{
							"x": 5, "y": 0, "z": 0,
						},
					},
				},
			},
			Version:   1,
			CreatedAt: time.Now(),
		},
	}

	return documents
}

// generateTestAvatars creates test avatar fixtures
func generateTestAvatars() []TestAvatar {
	avatars := []TestAvatar{
		{
			ID:          uuid.New(),
			Name:        "AI Assistant",
			Description: "Helpful AI assistant avatar",
			Provider:    "openai",
			Model:       "gpt-4",
			Personality: map[string]interface{}{
				"traits": []string{"helpful", "friendly", "knowledgeable"},
				"style":  "conversational",
				"tone":   "professional",
			},
			Appearance: map[string]interface{}{
				"model_url": "https://example.com/avatar.glb",
				"color":     "#4A90E2",
				"scale":     1.2,
			},
			Status:    "active",
			CreatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Creative Assistant",
			Description: "AI assistant specialized in creative tasks",
			Provider:    "claude",
			Model:       "claude-3-opus",
			Personality: map[string]interface{}{
				"traits": []string{"creative", "imaginative", "artistic"},
				"style":  "inspirational",
				"tone":   "enthusiastic",
			},
			Appearance: map[string]interface{}{
				"model_url": "https://example.com/creative_avatar.glb",
				"color":     "#E24A90",
				"scale":     1.0,
			},
			Status:    "active",
			CreatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Technical Expert",
			Description: "AI assistant for technical questions",
			Provider:    "gemini",
			Model:       "gemini-pro",
			Personality: map[string]interface{}{
				"traits": []string{"analytical", "precise", "technical"},
				"style":  "detailed",
				"tone":   "expert",
			},
			Appearance: map[string]interface{}{
				"model_url": "https://example.com/tech_avatar.glb",
				"color":     "#90E24A",
				"scale":     1.1,
			},
			Status:    "active",
			CreatedAt: time.Now(),
		},
	}

	return avatars
}

// generateTestClients creates test client fixtures
func generateTestClients() []TestClient {
	clients := []TestClient{
		{
			ID:          uuid.New(),
			Name:        "HD1 Web Client",
			Platform:    "web",
			Version:     "1.0.0",
			Description: "Official HD1 web client",
			Capabilities: []string{
				"3d_rendering",
				"webgl",
				"audio",
				"video",
				"webrtc",
			},
			Status: "active",
			Token:  "web_client_token_123",
			Metadata: map[string]interface{}{
				"user_agent":     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
				"resolution":     "1920x1080",
				"webgl_version":  "2.0",
				"audio_context":  true,
				"fullscreen":     true,
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "HD1 Mobile App",
			Platform:    "mobile",
			Version:     "1.0.0",
			Description: "HD1 mobile application",
			Capabilities: []string{
				"3d_rendering",
				"touch_input",
				"accelerometer",
				"camera",
				"location",
			},
			Status: "active",
			Token:  "mobile_client_token_456",
			Metadata: map[string]interface{}{
				"device_type":   "smartphone",
				"os_version":    "iOS 15.0",
				"screen_size":   "6.1 inches",
				"battery_level": 85,
				"network_type":  "5G",
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "HD1 Desktop App",
			Platform:    "desktop",
			Version:     "1.0.0",
			Description: "HD1 desktop application",
			Capabilities: []string{
				"3d_rendering",
				"opengl",
				"vulkan",
				"file_system",
				"native_ui",
			},
			Status: "active",
			Token:  "desktop_client_token_789",
			Metadata: map[string]interface{}{
				"os":           "Windows 10",
				"architecture": "x64",
				"cpu_cores":    8,
				"memory_gb":    16,
				"gpu_model":    "NVIDIA RTX 3080",
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "HD1 VR Client",
			Platform:    "vr",
			Version:     "1.0.0",
			Description: "HD1 virtual reality client",
			Capabilities: []string{
				"vr_rendering",
				"6dof_tracking",
				"hand_tracking",
				"spatial_audio",
				"haptic_feedback",
			},
			Status: "active",
			Token:  "vr_client_token_101",
			Metadata: map[string]interface{}{
				"headset_model": "Meta Quest 3",
				"tracking_type": "inside_out",
				"refresh_rate":  90,
				"resolution":    "2064x2208",
				"fov":           110,
			},
			CreatedAt: time.Now(),
		},
	}

	return clients
}

// generateTestRoles creates test role fixtures
func generateTestRoles() []TestRole {
	roles := []TestRole{
		{
			ID:          uuid.New(),
			Name:        "3D Designer",
			Description: "Can create and modify 3D content",
			Permissions: []string{
				"session.create",
				"session.update",
				"asset.upload",
				"asset.download",
				"object.create",
				"object.update",
				"object.delete",
				"material.create",
				"material.update",
			},
			Restrictions: map[string]interface{}{
				"max_sessions": 10,
				"max_assets":   100,
				"storage_limit": "10GB",
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Project Manager",
			Description: "Can manage projects and users",
			Permissions: []string{
				"session.create",
				"session.update",
				"session.delete",
				"user.invite",
				"user.manage",
				"analytics.view",
				"report.generate",
			},
			Restrictions: map[string]interface{}{
				"max_sessions": 50,
				"max_users":    20,
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Developer",
			Description: "Can develop and deploy services",
			Permissions: []string{
				"service.register",
				"service.update",
				"service.delete",
				"plugin.install",
				"plugin.configure",
				"api.access",
				"debug.access",
			},
			Restrictions: map[string]interface{}{
				"max_services": 5,
				"max_plugins":  10,
				"api_rate_limit": 10000,
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Content Viewer",
			Description: "Can view and interact with content",
			Permissions: []string{
				"session.view",
				"session.join",
				"asset.view",
				"asset.download",
				"object.view",
			},
			Restrictions: map[string]interface{}{
				"max_sessions": 5,
				"download_limit": "1GB",
			},
			CreatedAt: time.Now(),
		},
	}

	return roles
}

// linkFixtures creates relationships between different fixtures
func linkFixtures(fixtures *TestFixtures) {
	// Link sessions to users (owners)
	for i := range fixtures.Sessions {
		if i < len(fixtures.Users) {
			fixtures.Sessions[i].OwnerID = fixtures.Users[i].ID
		}
	}

	// Link organizations to users (owners)
	for i := range fixtures.Organizations {
		if i < len(fixtures.Users) {
			fixtures.Organizations[i].OwnerID = fixtures.Users[i].ID
		}
	}

	// Link assets to sessions and users
	for i := range fixtures.Assets {
		if i < len(fixtures.Sessions) {
			fixtures.Assets[i].SessionID = fixtures.Sessions[i].ID
		}
		if i < len(fixtures.Users) {
			fixtures.Assets[i].UserID = fixtures.Users[i].ID
		}
	}

	// Link documents to sessions
	for i := range fixtures.Documents {
		if i < len(fixtures.Sessions) {
			fixtures.Documents[i].SessionID = fixtures.Sessions[i].ID
		}
	}

	// Link avatars to sessions
	for i := range fixtures.Avatars {
		if i < len(fixtures.Sessions) {
			fixtures.Avatars[i].SessionID = fixtures.Sessions[i].ID
		}
	}

	// Link roles to organizations
	for i := range fixtures.Roles {
		if i < len(fixtures.Organizations) {
			fixtures.Roles[i].OrganizationID = fixtures.Organizations[i].ID
		}
	}
}

// SaveFixturesToFile saves fixtures to a JSON file
func SaveFixturesToFile(fixtures *TestFixtures, filename string) error {
	data, err := json.MarshalIndent(fixtures, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// LoadFixturesFromFile loads fixtures from a JSON file
func LoadFixturesFromFile(filename string) (*TestFixtures, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var fixtures TestFixtures
	err = json.Unmarshal(data, &fixtures)
	if err != nil {
		return nil, err
	}

	return &fixtures, nil
}

// CreateTestDataDirectory creates a directory for test data
func CreateTestDataDirectory() (string, error) {
	testDataDir := filepath.Join("test_data")
	
	if err := os.MkdirAll(testDataDir, 0755); err != nil {
		return "", err
	}

	return testDataDir, nil
}

// GenerateAndSaveFixtures generates and saves test fixtures
func GenerateAndSaveFixtures() error {
	fixtures := GenerateTestFixtures()
	
	testDataDir, err := CreateTestDataDirectory()
	if err != nil {
		return err
	}

	filename := filepath.Join(testDataDir, "test_fixtures.json")
	return SaveFixturesToFile(fixtures, filename)
}

// GetFixtureByID gets a fixture by ID from the appropriate collection
func (f *TestFixtures) GetFixtureByID(id uuid.UUID, fixtureType string) interface{} {
	switch fixtureType {
	case "user":
		for _, user := range f.Users {
			if user.ID == id {
				return user
			}
		}
	case "organization":
		for _, org := range f.Organizations {
			if org.ID == id {
				return org
			}
		}
	case "session":
		for _, session := range f.Sessions {
			if session.ID == id {
				return session
			}
		}
	case "service":
		for _, service := range f.Services {
			if service.ID == id {
				return service
			}
		}
	case "plugin":
		for _, plugin := range f.Plugins {
			if plugin.ID == id {
				return plugin
			}
		}
	case "asset":
		for _, asset := range f.Assets {
			if asset.ID == id {
				return asset
			}
		}
	case "document":
		for _, doc := range f.Documents {
			if doc.ID == id {
				return doc
			}
		}
	case "avatar":
		for _, avatar := range f.Avatars {
			if avatar.ID == id {
				return avatar
			}
		}
	case "client":
		for _, client := range f.Clients {
			if client.ID == id {
				return client
			}
		}
	case "role":
		for _, role := range f.Roles {
			if role.ID == id {
				return role
			}
		}
	}
	return nil
}

// GetRandomFixture gets a random fixture from a collection
func (f *TestFixtures) GetRandomFixture(fixtureType string) interface{} {
	switch fixtureType {
	case "user":
		if len(f.Users) > 0 {
			return f.Users[0]
		}
	case "organization":
		if len(f.Organizations) > 0 {
			return f.Organizations[0]
		}
	case "session":
		if len(f.Sessions) > 0 {
			return f.Sessions[0]
		}
	case "service":
		if len(f.Services) > 0 {
			return f.Services[0]
		}
	case "plugin":
		if len(f.Plugins) > 0 {
			return f.Plugins[0]
		}
	case "asset":
		if len(f.Assets) > 0 {
			return f.Assets[0]
		}
	case "document":
		if len(f.Documents) > 0 {
			return f.Documents[0]
		}
	case "avatar":
		if len(f.Avatars) > 0 {
			return f.Avatars[0]
		}
	case "client":
		if len(f.Clients) > 0 {
			return f.Clients[0]
		}
	case "role":
		if len(f.Roles) > 0 {
			return f.Roles[0]
		}
	}
	return nil
}