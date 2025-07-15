# Phase 2: Collaboration Implementation Plan
**Duration**: 3 months  
**Goal**: Transform HD1 into real-time collaborative platform  
**Endpoints**: 30 → 60

## Overview
Phase 2 builds upon the foundation to create a comprehensive real-time collaboration platform. We implement WebRTC for peer-to-peer connections, operational transform for conflict resolution, and spatial collaboration features including voice chat, screen sharing, and shared cursors.

## Technical Objectives

### 1. Real-time Collaboration Hub
**Current State**: Basic WebSocket broadcasting  
**Target State**: Advanced real-time collaboration with sub-100ms latency and conflict resolution

**Implementation Steps**:
1. **WebRTC Integration**
   - P2P mesh for direct client connections
   - STUN/TURN servers for NAT traversal
   - Signaling server for connection establishment
   - Bandwidth optimization and quality control

2. **Operational Transform (OT)**
   - Conflict resolution for concurrent operations
   - Operation ordering and consistency
   - Rollback and recovery mechanisms
   - State synchronization across clients

3. **Collaboration State Management**
   - Shared cursors and user presence
   - Real-time activity tracking
   - Collaborative object locking
   - Session state persistence

### 2. Spatial Voice Chat System
**Current State**: No voice communication  
**Target State**: Positional audio with spatial processing and quality control

**Implementation Steps**:
1. **Audio Processing**
   - WebRTC audio streams
   - Spatial audio positioning
   - Noise cancellation and echo suppression
   - Dynamic audio quality adjustment

2. **Voice Chat Management**
   - Voice channel creation and management
   - Participant audio controls
   - Push-to-talk and voice activation
   - Audio recording and playback

### 3. Screen Sharing as 3D Surfaces
**Current State**: No screen sharing  
**Target State**: Screen sharing rendered as 3D surfaces in virtual space

**Implementation Steps**:
1. **Screen Capture**
   - WebRTC screen sharing API
   - Multi-monitor support
   - Application window capture
   - Performance optimization

2. **3D Surface Rendering**
   - Screen content as Three.js textures
   - Spatial positioning of shared screens
   - Interactive screen surfaces
   - Multi-screen collaboration

### 4. Asset Streaming Pipeline
**Current State**: Static asset loading  
**Target State**: Progressive asset streaming with CDN integration

**Implementation Steps**:
1. **Asset Management**
   - Asset versioning and lifecycle
   - Progressive loading strategies
   - CDN integration and caching
   - Asset optimization and compression

2. **Streaming Infrastructure**
   - Chunked asset delivery
   - Bandwidth-adaptive streaming
   - Offline asset caching
   - Asset synchronization across clients

### 5. User Management System
**Current State**: No user management  
**Target State**: Complete user lifecycle with profiles and preferences

**Implementation Steps**:
1. **User Profiles**
   - User profile creation and management
   - Avatar customization
   - Preference settings
   - Social features and connections

2. **User Analytics**
   - User behavior tracking
   - Session analytics
   - Performance monitoring
   - Usage insights and reporting

## Detailed Implementation

### Step 1: WebRTC Integration
```go
// src/webrtc/manager.go
package webrtc

import (
    "context"
    "encoding/json"
    "time"
    
    "github.com/google/uuid"
    "github.com/pion/webrtc/v3"
    "holodeck1/logging"
)

type WebRTCManager struct {
    api          *webrtc.API
    config       webrtc.Configuration
    connections  map[uuid.UUID]*webrtc.PeerConnection
    dataChannels map[uuid.UUID]*webrtc.DataChannel
}

type SignalingMessage struct {
    Type      string                 `json:"type"`
    SessionID uuid.UUID              `json:"session_id"`
    UserID    uuid.UUID              `json:"user_id"`
    TargetID  uuid.UUID              `json:"target_id"`
    Data      map[string]interface{} `json:"data"`
}

func NewWebRTCManager() *WebRTCManager {
    config := webrtc.Configuration{
        ICEServers: []webrtc.ICEServer{
            {
                URLs: []string{"stun:stun.l.google.com:19302"},
            },
            {
                URLs:       []string{"turn:turn.example.com:3478"},
                Username:   "username",
                Credential: "password",
            },
        },
    }
    
    return &WebRTCManager{
        api:          webrtc.NewAPI(),
        config:       config,
        connections:  make(map[uuid.UUID]*webrtc.PeerConnection),
        dataChannels: make(map[uuid.UUID]*webrtc.DataChannel),
    }
}

func (w *WebRTCManager) CreatePeerConnection(ctx context.Context, userID uuid.UUID) (*webrtc.PeerConnection, error) {
    pc, err := w.api.NewPeerConnection(w.config)
    if err != nil {
        logging.Error("failed to create peer connection", map[string]interface{}{
            "user_id": userID,
            "error": err.Error(),
        })
        return nil, err
    }
    
    // Set up data channel for real-time collaboration
    dataChannel, err := pc.CreateDataChannel("collaboration", nil)
    if err != nil {
        logging.Error("failed to create data channel", map[string]interface{}{
            "user_id": userID,
            "error": err.Error(),
        })
        return nil, err
    }
    
    dataChannel.OnOpen(func() {
        logging.Info("data channel opened", map[string]interface{}{
            "user_id": userID,
        })
    })
    
    dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
        w.handleCollaborationMessage(userID, msg.Data)
    })
    
    w.connections[userID] = pc
    w.dataChannels[userID] = dataChannel
    
    return pc, nil
}

func (w *WebRTCManager) handleCollaborationMessage(userID uuid.UUID, data []byte) {
    var msg CollaborationMessage
    if err := json.Unmarshal(data, &msg); err != nil {
        logging.Error("failed to parse collaboration message", map[string]interface{}{
            "user_id": userID,
            "error": err.Error(),
        })
        return
    }
    
    switch msg.Type {
    case "cursor_update":
        w.handleCursorUpdate(userID, msg.Data)
    case "object_lock":
        w.handleObjectLock(userID, msg.Data)
    case "operation":
        w.handleOperation(userID, msg.Data)
    }
}

func (w *WebRTCManager) BroadcastToSession(sessionID uuid.UUID, message []byte) error {
    // Broadcast message to all peers in session
    for userID, dataChannel := range w.dataChannels {
        if dataChannel.ReadyState() == webrtc.DataChannelStateOpen {
            err := dataChannel.Send(message)
            if err != nil {
                logging.Error("failed to send message to peer", map[string]interface{}{
                    "user_id": userID,
                    "error": err.Error(),
                })
            }
        }
    }
    
    return nil
}
```

### Step 2: Operational Transform Implementation
```go
// src/collaboration/transform.go
package collaboration

import (
    "context"
    "encoding/json"
    "sync"
    "time"
    
    "github.com/google/uuid"
    "holodeck1/logging"
)

type OperationalTransform struct {
    mu                sync.RWMutex
    operationHistory  []Operation
    currentState      map[string]interface{}
    pendingOperations map[uuid.UUID][]Operation
}

type Operation struct {
    ID        uuid.UUID              `json:"id"`
    Type      string                 `json:"type"`
    UserID    uuid.UUID              `json:"user_id"`
    SessionID uuid.UUID              `json:"session_id"`
    Timestamp time.Time              `json:"timestamp"`
    Data      map[string]interface{} `json:"data"`
    Vector    []int                  `json:"vector"`  // Vector clock for ordering
}

type TransformResult struct {
    Operation     Operation `json:"operation"`
    Transformed   bool      `json:"transformed"`
    ConflictsWith []uuid.UUID `json:"conflicts_with"`
}

func NewOperationalTransform() *OperationalTransform {
    return &OperationalTransform{
        operationHistory:  make([]Operation, 0),
        currentState:      make(map[string]interface{}),
        pendingOperations: make(map[uuid.UUID][]Operation),
    }
}

func (ot *OperationalTransform) ApplyOperation(ctx context.Context, op Operation) (*TransformResult, error) {
    ot.mu.Lock()
    defer ot.mu.Unlock()
    
    // Check for conflicts with pending operations
    conflicts := ot.detectConflicts(op)
    
    // Transform operation if conflicts exist
    transformed := false
    if len(conflicts) > 0 {
        transformedOp, err := ot.transformOperation(op, conflicts)
        if err != nil {
            logging.Error("failed to transform operation", map[string]interface{}{
                "operation_id": op.ID,
                "error": err.Error(),
            })
            return nil, err
        }
        op = transformedOp
        transformed = true
    }
    
    // Apply operation to current state
    err := ot.applyToState(op)
    if err != nil {
        logging.Error("failed to apply operation to state", map[string]interface{}{
            "operation_id": op.ID,
            "error": err.Error(),
        })
        return nil, err
    }
    
    // Add to operation history
    ot.operationHistory = append(ot.operationHistory, op)
    
    result := &TransformResult{
        Operation:     op,
        Transformed:   transformed,
        ConflictsWith: conflicts,
    }
    
    return result, nil
}

func (ot *OperationalTransform) detectConflicts(op Operation) []uuid.UUID {
    conflicts := make([]uuid.UUID, 0)
    
    // Check for operations affecting the same object
    for _, historyOp := range ot.operationHistory {
        if ot.operationsConflict(op, historyOp) {
            conflicts = append(conflicts, historyOp.ID)
        }
    }
    
    return conflicts
}

func (ot *OperationalTransform) operationsConflict(op1, op2 Operation) bool {
    // Check if operations affect the same object
    obj1ID, ok1 := op1.Data["object_id"]
    obj2ID, ok2 := op2.Data["object_id"]
    
    if ok1 && ok2 && obj1ID == obj2ID {
        // Same object, check for conflicting operations
        return ot.isConflictingOperation(op1.Type, op2.Type)
    }
    
    return false
}

func (ot *OperationalTransform) isConflictingOperation(type1, type2 string) bool {
    conflictMap := map[string][]string{
        "entity_update": {"entity_update", "entity_delete"},
        "entity_delete": {"entity_update", "entity_delete"},
        "avatar_move":   {"avatar_move"},
    }
    
    conflicts, exists := conflictMap[type1]
    if !exists {
        return false
    }
    
    for _, conflictType := range conflicts {
        if type2 == conflictType {
            return true
        }
    }
    
    return false
}

func (ot *OperationalTransform) transformOperation(op Operation, conflicts []uuid.UUID) (Operation, error) {
    // Transform operation based on conflicts
    switch op.Type {
    case "entity_update":
        return ot.transformEntityUpdate(op, conflicts)
    case "entity_delete":
        return ot.transformEntityDelete(op, conflicts)
    case "avatar_move":
        return ot.transformAvatarMove(op, conflicts)
    default:
        return op, nil
    }
}

func (ot *OperationalTransform) transformEntityUpdate(op Operation, conflicts []uuid.UUID) (Operation, error) {
    // For entity updates, merge changes or apply last-writer-wins
    transformedOp := op
    
    // Apply transformation logic based on conflict resolution strategy
    for _, conflictID := range conflicts {
        conflictOp := ot.findOperation(conflictID)
        if conflictOp != nil {
            // Merge updates or apply precedence rules
            transformedOp = ot.mergeEntityUpdates(transformedOp, *conflictOp)
        }
    }
    
    return transformedOp, nil
}

func (ot *OperationalTransform) applyToState(op Operation) error {
    switch op.Type {
    case "entity_create":
        return ot.applyEntityCreate(op)
    case "entity_update":
        return ot.applyEntityUpdate(op)
    case "entity_delete":
        return ot.applyEntityDelete(op)
    case "avatar_move":
        return ot.applyAvatarMove(op)
    default:
        return nil
    }
}

func (ot *OperationalTransform) applyEntityUpdate(op Operation) error {
    entityID, ok := op.Data["entity_id"].(string)
    if !ok {
        return fmt.Errorf("invalid entity_id in operation")
    }
    
    // Update entity in current state
    if entity, exists := ot.currentState[entityID]; exists {
        entityMap := entity.(map[string]interface{})
        
        // Apply updates
        if position, ok := op.Data["position"]; ok {
            entityMap["position"] = position
        }
        if rotation, ok := op.Data["rotation"]; ok {
            entityMap["rotation"] = rotation
        }
        if scale, ok := op.Data["scale"]; ok {
            entityMap["scale"] = scale
        }
        
        ot.currentState[entityID] = entityMap
    }
    
    return nil
}

func (ot *OperationalTransform) GetCurrentState() map[string]interface{} {
    ot.mu.RLock()
    defer ot.mu.RUnlock()
    
    // Return deep copy of current state
    stateCopy := make(map[string]interface{})
    for k, v := range ot.currentState {
        stateCopy[k] = v
    }
    
    return stateCopy
}
```

### Step 3: Collaboration API Endpoints
```go
// src/api/collaboration/handlers.go
package collaboration

import (
    "encoding/json"
    "net/http"
    
    "github.com/gorilla/mux"
    "github.com/google/uuid"
    "holodeck1/collaboration"
    "holodeck1/logging"
)

type CollaborationHandler struct {
    transform *collaboration.OperationalTransform
    webrtc    *webrtc.Manager
}

func NewCollaborationHandler(transform *collaboration.OperationalTransform, webrtc *webrtc.Manager) *CollaborationHandler {
    return &CollaborationHandler{
        transform: transform,
        webrtc:    webrtc,
    }
}

func (h *CollaborationHandler) GetSharedCursors(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    vars := mux.Vars(r)
    sessionID, err := uuid.Parse(vars["sessionId"])
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }
    
    cursors, err := h.getCursorsForSession(ctx, sessionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    response := map[string]interface{}{
        "success": true,
        "cursors": cursors,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *CollaborationHandler) UpdateSharedCursor(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    vars := mux.Vars(r)
    sessionID, err := uuid.Parse(vars["sessionId"])
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }
    
    var req UpdateCursorRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // Create cursor update operation
    op := collaboration.Operation{
        ID:        uuid.New(),
        Type:      "cursor_update",
        UserID:    req.UserID,
        SessionID: sessionID,
        Timestamp: time.Now(),
        Data: map[string]interface{}{
            "position": req.Position,
            "rotation": req.Rotation,
            "color":    req.Color,
        },
    }
    
    // Apply operation through transform
    result, err := h.transform.ApplyOperation(ctx, op)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Broadcast to session via WebRTC
    message, _ := json.Marshal(result.Operation)
    h.webrtc.BroadcastToSession(sessionID, message)
    
    response := map[string]interface{}{
        "success": true,
        "message": "Cursor updated successfully",
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *CollaborationHandler) JoinVoiceChat(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    vars := mux.Vars(r)
    sessionID, err := uuid.Parse(vars["sessionId"])
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }
    
    var req JoinVoiceChatRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // Create WebRTC peer connection for voice
    pc, err := h.webrtc.CreatePeerConnection(ctx, req.UserID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Set up audio track handling
    pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
        h.handleAudioTrack(sessionID, req.UserID, track)
    })
    
    response := map[string]interface{}{
        "success": true,
        "message": "Voice chat joined successfully",
        "peer_connection_id": pc.ID(),
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *CollaborationHandler) StartScreenShare(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    vars := mux.Vars(r)
    sessionID, err := uuid.Parse(vars["sessionId"])
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }
    
    var req StartScreenShareRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // Create 3D surface for screen sharing
    surfaceID := uuid.New()
    surface := map[string]interface{}{
        "id":       surfaceID,
        "type":     "screen_share",
        "user_id":  req.UserID,
        "position": req.Position,
        "scale":    req.Scale,
        "rotation": req.Rotation,
    }
    
    // Create operation for screen share surface
    op := collaboration.Operation{
        ID:        uuid.New(),
        Type:      "entity_create",
        UserID:    req.UserID,
        SessionID: sessionID,
        Timestamp: time.Now(),
        Data: map[string]interface{}{
            "entity_type": "screen_share_surface",
            "surface":     surface,
        },
    }
    
    // Apply operation
    result, err := h.transform.ApplyOperation(ctx, op)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Broadcast to session
    message, _ := json.Marshal(result.Operation)
    h.webrtc.BroadcastToSession(sessionID, message)
    
    response := map[string]interface{}{
        "success":    true,
        "surface_id": surfaceID,
        "message":    "Screen sharing started successfully",
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

### Step 4: Asset Management System
```go
// src/assets/manager.go
package assets

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "io"
    "time"
    
    "github.com/google/uuid"
    "holodeck1/database"
    "holodeck1/logging"
)

type AssetManager struct {
    db  *database.DB
    cdn *CDNManager
}

type Asset struct {
    ID          uuid.UUID `json:"id"`
    Name        string    `json:"name"`
    Type        string    `json:"type"`
    Size        int64     `json:"size"`
    Hash        string    `json:"hash"`
    URL         string    `json:"url"`
    CDNUrl      string    `json:"cdn_url"`
    Version     int       `json:"version"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    Metadata    map[string]interface{} `json:"metadata"`
}

func NewAssetManager(db *database.DB, cdn *CDNManager) *AssetManager {
    return &AssetManager{
        db:  db,
        cdn: cdn,
    }
}

func (am *AssetManager) UploadAsset(ctx context.Context, name, assetType string, data io.Reader) (*Asset, error) {
    // Calculate hash
    hasher := sha256.New()
    var buffer bytes.Buffer
    
    tee := io.TeeReader(data, &buffer)
    if _, err := io.Copy(hasher, tee); err != nil {
        return nil, err
    }
    
    hash := hex.EncodeToString(hasher.Sum(nil))
    size := int64(buffer.Len())
    
    // Check if asset already exists
    existingAsset, err := am.db.GetAssetByHash(ctx, hash)
    if err == nil {
        return existingAsset, nil
    }
    
    // Create new asset
    asset := &Asset{
        ID:        uuid.New(),
        Name:      name,
        Type:      assetType,
        Size:      size,
        Hash:      hash,
        Version:   1,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    // Upload to CDN
    cdnUrl, err := am.cdn.Upload(ctx, asset.ID.String(), &buffer)
    if err != nil {
        logging.Error("failed to upload asset to CDN", map[string]interface{}{
            "asset_id": asset.ID,
            "error": err.Error(),
        })
        return nil, err
    }
    
    asset.CDNUrl = cdnUrl
    asset.URL = cdnUrl
    
    // Save to database
    err = am.db.CreateAsset(ctx, asset)
    if err != nil {
        logging.Error("failed to save asset to database", map[string]interface{}{
            "asset_id": asset.ID,
            "error": err.Error(),
        })
        return nil, err
    }
    
    return asset, nil
}

func (am *AssetManager) GetAsset(ctx context.Context, assetID uuid.UUID) (*Asset, error) {
    asset, err := am.db.GetAsset(ctx, assetID)
    if err != nil {
        logging.Error("failed to get asset", map[string]interface{}{
            "asset_id": assetID,
            "error": err.Error(),
        })
        return nil, err
    }
    
    return asset, nil
}

func (am *AssetManager) StreamAsset(ctx context.Context, assetID uuid.UUID) (io.ReadCloser, error) {
    asset, err := am.GetAsset(ctx, assetID)
    if err != nil {
        return nil, err
    }
    
    // Stream from CDN
    stream, err := am.cdn.Stream(ctx, asset.CDNUrl)
    if err != nil {
        logging.Error("failed to stream asset from CDN", map[string]interface{}{
            "asset_id": assetID,
            "error": err.Error(),
        })
        return nil, err
    }
    
    return stream, nil
}
```

## Implementation Timeline

### Month 1: WebRTC Foundation
- [ ] Week 1: WebRTC infrastructure setup
- [ ] Week 2: Peer connection management
- [ ] Week 3: Data channels and messaging
- [ ] Week 4: Connection quality and optimization

### Month 2: Collaboration Core
- [ ] Week 1: Operational transform implementation
- [ ] Week 2: Conflict resolution algorithms
- [ ] Week 3: State synchronization
- [ ] Week 4: Collaboration API endpoints

### Month 3: Advanced Features
- [ ] Week 1: Spatial voice chat
- [ ] Week 2: Screen sharing as 3D surfaces
- [ ] Week 3: Asset streaming pipeline
- [ ] Week 4: User management and analytics

## Success Criteria

### Technical Metrics
- [ ] Sub-100ms latency for real-time operations
- [ ] WebRTC connections support 50+ participants per session
- [ ] Operational transform handles 1000+ operations per second
- [ ] Asset streaming supports 10GB+ files
- [ ] Voice chat with spatial audio positioning

### Quality Metrics
- [ ] 100% API test coverage for all 30 new endpoints
- [ ] Zero data loss during conflict resolution
- [ ] Voice quality metrics within acceptable ranges
- [ ] Screen sharing with 60fps performance
- [ ] Asset streaming with progressive loading

### Business Metrics
- [ ] Real-time collaboration enables multi-user sessions
- [ ] Voice chat increases session engagement by 40%
- [ ] Screen sharing supports presentation workflows
- [ ] Asset streaming reduces loading times by 70%
- [ ] User management supports 10,000+ concurrent users

## Risk Mitigation

### Technical Risks
- **WebRTC Complexity**: Use proven libraries and extensive testing
- **NAT Traversal**: Implement STUN/TURN servers for reliability
- **Operational Transform**: Use established algorithms and conflict resolution
- **Audio Quality**: Implement adaptive bitrate and noise cancellation

### Performance Risks
- **Bandwidth Usage**: Implement quality adaptation and compression
- **Memory Consumption**: Use object pooling and garbage collection optimization
- **CPU Usage**: Optimize audio processing and 3D rendering
- **Network Latency**: Use edge servers and CDN optimization

## Deliverables

### Code Deliverables
- [ ] WebRTC peer-to-peer infrastructure
- [ ] Operational transform conflict resolution
- [ ] Spatial voice chat system
- [ ] Screen sharing as 3D surfaces
- [ ] Asset streaming pipeline
- [ ] User management system
- [ ] 30 new API endpoints with documentation

### Documentation Deliverables
- [ ] WebRTC integration guide
- [ ] Collaboration development guide
- [ ] Voice chat configuration guide
- [ ] Asset management documentation
- [ ] Performance optimization guide

### Testing Deliverables
- [ ] WebRTC connection tests
- [ ] Conflict resolution tests
- [ ] Voice chat quality tests
- [ ] Screen sharing performance tests
- [ ] Asset streaming load tests
- [ ] End-to-end collaboration tests

This completes the detailed Phase 2 implementation plan. The collaboration features will transform HD1 into a comprehensive real-time collaborative platform with advanced features for multi-user interaction.