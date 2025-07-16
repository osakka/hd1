package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"holodeck1/api/assets"
	"holodeck1/api/ot"
	"holodeck1/api/webrtc"
	"holodeck1/database"
	"holodeck1/server"
)

// TestPhase2Collaboration tests all Phase 2 components
func TestPhase2Collaboration(t *testing.T) {
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
	webrtcHandlers := webrtc.NewHandlers(db, hub)
	otHandlers := ot.NewHandlers(db, hub)
	assetHandlers := assets.NewHandlers(db)

	// Register routes
	api.HandleFunc("/webrtc/sessions", webrtcHandlers.CreateRTCSession).Methods("POST")
	api.HandleFunc("/webrtc/offer", webrtcHandlers.SendOffer).Methods("POST")
	api.HandleFunc("/webrtc/answer", webrtcHandlers.SendAnswer).Methods("POST")
	api.HandleFunc("/webrtc/ice", webrtcHandlers.SendICECandidate).Methods("POST")
	api.HandleFunc("/ot/documents", otHandlers.CreateDocument).Methods("POST")
	api.HandleFunc("/ot/documents/{id}/operations", otHandlers.ApplyOperation).Methods("POST")
	api.HandleFunc("/assets", assetHandlers.UploadAsset).Methods("POST")

	// Test WebRTC session creation
	t.Run("WebRTC Session Creation", func(t *testing.T) {
		userID := createTestUser(t, db)
		sessionID := createTestSession(t, db, userID)

		reqBody := map[string]interface{}{
			"session_id": sessionID,
			"peer_id":    userID,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/webrtc/sessions", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "id")
		assert.Contains(t, response, "ice_servers")
	})

	// Test WebRTC signaling
	t.Run("WebRTC Signaling", func(t *testing.T) {
		peerID := uuid.New()
		
		// Send offer
		offerBody := map[string]interface{}{
			"peer_id": peerID,
			"offer": map[string]interface{}{
				"type": "offer",
				"sdp":  "v=0\r\no=- 123456 2 IN IP4 127.0.0.1\r\n...",
			},
		}
		body, _ := json.Marshal(offerBody)

		req := httptest.NewRequest("POST", "/api/webrtc/offer", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		// Send answer
		answerBody := map[string]interface{}{
			"peer_id": peerID,
			"answer": map[string]interface{}{
				"type": "answer",
				"sdp":  "v=0\r\no=- 654321 2 IN IP4 127.0.0.1\r\n...",
			},
		}
		body, _ = json.Marshal(answerBody)

		req = httptest.NewRequest("POST", "/api/webrtc/answer", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		// Send ICE candidate
		iceBody := map[string]interface{}{
			"peer_id": peerID,
			"candidate": map[string]interface{}{
				"candidate":     "candidate:1 1 UDP 2122252543 192.168.1.100 50000 typ host",
				"sdpMLineIndex": 0,
				"sdpMid":        "0",
			},
		}
		body, _ = json.Marshal(iceBody)

		req = httptest.NewRequest("POST", "/api/webrtc/ice", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// Test Operational Transform
	t.Run("Operational Transform", func(t *testing.T) {
		userID := createTestUser(t, db)
		sessionID := createTestSession(t, db, userID)

		// Create document
		docBody := map[string]interface{}{
			"session_id": sessionID,
			"type":       "scene",
			"content": map[string]interface{}{
				"entities": []interface{}{},
				"lights":   []interface{}{},
			},
		}
		body, _ := json.Marshal(docBody)

		req := httptest.NewRequest("POST", "/api/ot/documents", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)

		var docResponse map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &docResponse)
		docID := docResponse["id"].(string)

		// Apply operation
		opBody := map[string]interface{}{
			"type": "insert",
			"path": "/entities/0",
			"value": map[string]interface{}{
				"id":   uuid.New().String(),
				"type": "box",
				"position": map[string]float64{
					"x": 0, "y": 1, "z": 0,
				},
			},
			"version":   1,
			"client_id": userID.String(),
		}
		body, _ = json.Marshal(opBody)

		req = httptest.NewRequest("POST", fmt.Sprintf("/api/ot/documents/%s/operations", docID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var opResponse map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &opResponse)
		assert.NoError(t, err)
		assert.Contains(t, opResponse, "transformed")
		assert.Contains(t, opResponse, "version")
	})

	// Test Asset Management
	t.Run("Asset Management", func(t *testing.T) {
		userID := createTestUser(t, db)
		sessionID := createTestSession(t, db, userID)

		// Create multipart form for file upload
		var b bytes.Buffer
		w := multipart.NewWriter(&b)

		// Add file field
		fw, err := w.CreateFormFile("file", "test-model.glb")
		require.NoError(t, err)
		
		// Write dummy GLB data
		glbHeader := []byte("glTF\x02\x00\x00\x00")
		fw.Write(glbHeader)
		fw.Write(make([]byte, 100)) // Dummy data

		// Add other fields
		w.WriteField("name", "Test 3D Model")
		w.WriteField("type", "model")
		w.WriteField("session_id", sessionID.String())
		w.WriteField("user_id", userID.String())
		w.Close()

		req := httptest.NewRequest("POST", "/api/assets", &b)
		req.Header.Set("Content-Type", w.FormDataContentType())
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "id")
		assert.Contains(t, response, "url")
		assert.Equal(t, "Test 3D Model", response["name"])
	})
}

// TestWebSocketSync tests real-time WebSocket synchronization
func TestWebSocketSync(t *testing.T) {
	db := setupTestDatabase(t)
	defer db.Close()

	hub := server.NewHub()
	go hub.Run()

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.ServeWS(hub, w, r)
	}))
	defer server.Close()

	// Convert http:// to ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	t.Run("WebSocket Connection", func(t *testing.T) {
		// Connect client 1
		ws1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		require.NoError(t, err)
		defer ws1.Close()

		// Connect client 2
		ws2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		require.NoError(t, err)
		defer ws2.Close()

		// Client 1 sends message
		message := map[string]interface{}{
			"type": "entity_create",
			"data": map[string]interface{}{
				"entity_id": uuid.New().String(),
				"type":      "sphere",
				"position":  map[string]float64{"x": 1, "y": 2, "z": 3},
			},
			"sequence": 1,
		}
		err = ws1.WriteJSON(message)
		assert.NoError(t, err)

		// Client 2 should receive the message
		var received map[string]interface{}
		err = ws2.ReadJSON(&received)
		assert.NoError(t, err)
		assert.Equal(t, "entity_create", received["type"])
		assert.Equal(t, float64(1), received["sequence"])
	})
}

// TestAssetProcessing tests asset upload and processing
func TestAssetProcessing(t *testing.T) {
	db := setupTestDatabase(t)
	defer db.Close()

	err := db.InitializeSchema()
	require.NoError(t, err)

	assetManager := assets.NewAssetManager(db)
	userID := createTestUser(t, db)
	sessionID := createTestSession(t, db, userID)

	t.Run("Asset Lifecycle", func(t *testing.T) {
		// Create test file
		tempFile, err := os.CreateTemp("", "test-texture-*.png")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())

		// Write PNG header
		pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
		tempFile.Write(pngHeader)
		tempFile.Close()

		// Open file for reading
		file, err := os.Open(tempFile.Name())
		require.NoError(t, err)
		defer file.Close()

		// Upload asset
		asset, err := assetManager.UploadAsset(&assets.AssetUpload{
			Name:      "Test Texture",
			Type:      "texture",
			File:      file,
			FileSize:  8,
			MimeType:  "image/png",
			SessionID: sessionID,
			UserID:    userID,
		})
		require.NoError(t, err)
		assert.NotNil(t, asset)
		assert.Equal(t, "Test Texture", asset.Name)

		// Get asset
		retrieved, err := assetManager.GetAsset(asset.ID)
		require.NoError(t, err)
		assert.Equal(t, asset.Name, retrieved.Name)

		// Update asset metadata
		err = assetManager.UpdateAsset(asset.ID, map[string]interface{}{
			"tags": []string{"texture", "test"},
			"metadata": map[string]interface{}{
				"resolution": "1024x1024",
			},
		})
		assert.NoError(t, err)

		// List assets
		assets, err := assetManager.ListAssets(&assets.AssetFilter{
			SessionID: sessionID,
			Type:      "texture",
			Limit:     10,
		})
		require.NoError(t, err)
		assert.Len(t, assets, 1)

		// Process asset (optimization)
		err = assetManager.ProcessAsset(asset.ID, &assets.ProcessingOptions{
			Optimize:   true,
			MaxSize:    1024,
			Format:     "webp",
		})
		assert.NoError(t, err)

		// Delete asset
		err = assetManager.DeleteAsset(asset.ID)
		assert.NoError(t, err)
	})
}

// TestOperationalTransform tests OT conflict resolution
func TestOperationalTransform(t *testing.T) {
	db := setupTestDatabase(t)
	defer db.Close()

	err := db.InitializeSchema()
	require.NoError(t, err)

	otManager := ot.NewOTManager(db)
	sessionID := uuid.New()

	t.Run("Concurrent Operations", func(t *testing.T) {
		// Create document
		doc, err := otManager.CreateDocument(&ot.CreateDocumentRequest{
			SessionID: sessionID,
			Type:      "scene",
			Content: map[string]interface{}{
				"entities": []interface{}{
					map[string]interface{}{
						"id":   "entity1",
						"type": "box",
						"position": map[string]float64{
							"x": 0, "y": 0, "z": 0,
						},
					},
				},
			},
		})
		require.NoError(t, err)

		// Client A operation
		opA := &ot.Operation{
			Type:     "update",
			Path:     "/entities/0/position/x",
			Value:    float64(5),
			Version:  1,
			ClientID: "clientA",
		}

		// Client B operation (concurrent)
		opB := &ot.Operation{
			Type:     "update",
			Path:     "/entities/0/position/y",
			Value:    float64(10),
			Version:  1,
			ClientID: "clientB",
		}

		// Apply operations
		resultA, err := otManager.ApplyOperation(doc.ID, opA)
		require.NoError(t, err)
		assert.Equal(t, 2, resultA.Version)

		resultB, err := otManager.ApplyOperation(doc.ID, opB)
		require.NoError(t, err)
		assert.Equal(t, 3, resultB.Version)

		// Get final document state
		snapshot, err := otManager.GetSnapshot(doc.ID)
		require.NoError(t, err)

		entities := snapshot.Content["entities"].([]interface{})
		entity := entities[0].(map[string]interface{})
		position := entity["position"].(map[string]interface{})

		assert.Equal(t, float64(5), position["x"])
		assert.Equal(t, float64(10), position["y"])
		assert.Equal(t, float64(0), position["z"])
	})

	t.Run("Conflict Resolution", func(t *testing.T) {
		// Create document with array
		doc, err := otManager.CreateDocument(&ot.CreateDocumentRequest{
			SessionID: sessionID,
			Type:      "scene",
			Content: map[string]interface{}{
				"entities": []interface{}{},
			},
		})
		require.NoError(t, err)

		// Two clients insert at same position
		opA := &ot.Operation{
			Type:  "insert",
			Path:  "/entities/0",
			Value: map[string]interface{}{"id": "A", "type": "box"},
			Version:  1,
			ClientID: "clientA",
		}

		opB := &ot.Operation{
			Type:  "insert",
			Path:  "/entities/0",
			Value: map[string]interface{}{"id": "B", "type": "sphere"},
			Version:  1,
			ClientID: "clientB",
		}

		// Apply both operations
		_, err = otManager.ApplyOperation(doc.ID, opA)
		require.NoError(t, err)

		resultB, err := otManager.ApplyOperation(doc.ID, opB)
		require.NoError(t, err)

		// Check transformed operation
		assert.True(t, resultB.Transformed)

		// Verify both entities exist
		snapshot, err := otManager.GetSnapshot(doc.ID)
		require.NoError(t, err)

		entities := snapshot.Content["entities"].([]interface{})
		assert.Len(t, entities, 2)
	})
}

// Helper function to create test session
func createTestSession(t *testing.T, db *database.DB, userID uuid.UUID) uuid.UUID {
	sessionID := uuid.New()
	_, err := db.Exec(`
		INSERT INTO sessions (id, name, owner_id, created_at)
		VALUES ($1, $2, $3, $4)
	`, sessionID, "Test Session", userID, time.Now())
	require.NoError(t, err)
	return sessionID
}