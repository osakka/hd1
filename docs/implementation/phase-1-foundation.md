# Phase 1: Foundation Implementation Plan
**Duration**: 2 months  
**Goal**: Transform HD1 from single-tenant to multi-tenant universal platform  
**Endpoints**: 11 → 30

## Overview
Phase 1 establishes the foundational architecture for the universal 3D interface platform. We migrate from single-session to multi-tenant architecture and create the service registry that enables any service to integrate with HD1.

## Technical Objectives

### 1. Multi-Tenant Session Architecture
**Current State**: Single WebSocket connection with basic session association  
**Target State**: Full multi-tenant sessions with isolation, permissions, and scalability

**Implementation Steps**:
1. **Database Schema Evolution**
   - Add `sessions` table with full metadata
   - Add `participants` table for session membership
   - Add `permissions` table for role-based access
   - Add indexes for performance optimization

2. **Session Management Service**
   - Create `SessionManager` struct with CRUD operations
   - Implement session lifecycle management
   - Add participant tracking and permissions
   - Include session analytics and monitoring

3. **WebSocket Hub Enhancement**
   - Modify hub to support multiple sessions
   - Add session-based message routing
   - Implement session isolation in broadcasts
   - Add connection pooling and load balancing

### 2. Service Registry Implementation
**Current State**: No external service integration  
**Target State**: Universal service registry where any API can register and render 3D interfaces

**Implementation Steps**:
1. **Service Registry Database**
   - Add `services` table with metadata
   - Add `service_health` table for monitoring
   - Add `ui_mappings` table for 3D interface mappings
   - Add `service_permissions` table for access control

2. **Service Management API**
   - Create service registration endpoints
   - Implement service discovery and health checks
   - Add service lifecycle management
   - Include service analytics and monitoring

3. **Service Integration Framework**
   - Create service adapter pattern
   - Implement automatic UI-to-3D mapping
   - Add service communication protocols
   - Include service sandbox and security

### 3. Enhanced Authentication System
**Current State**: No authentication  
**Target State**: OAuth 2.0 with role-based access control

**Implementation Steps**:
1. **Authentication Infrastructure**
   - Add OAuth 2.0 server implementation
   - Create JWT token management
   - Add user management database tables
   - Implement role-based access control

2. **API Security**
   - Add authentication middleware
   - Implement authorization checks
   - Add rate limiting and throttling
   - Include security audit logging

### 4. Database Migration and Scaling
**Current State**: In-memory data structures  
**Target State**: PostgreSQL with Redis caching and horizontal scaling

**Implementation Steps**:
1. **Database Setup**
   - PostgreSQL primary database
   - Redis for caching and sessions
   - Database migration scripts
   - Connection pooling and optimization

2. **Data Layer**
   - Create data access layer (DAL)
   - Implement repository pattern
   - Add database transaction management
   - Include database monitoring

## Detailed Implementation

### Step 1: Database Schema Design
```sql
-- Sessions table
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    owner_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'active',
    visibility VARCHAR(20) DEFAULT 'private',
    max_participants INTEGER DEFAULT 10,
    current_participants INTEGER DEFAULT 0,
    settings JSONB,
    metadata JSONB
);

-- Participants table
CREATE TABLE participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID REFERENCES sessions(id),
    user_id UUID NOT NULL,
    role VARCHAR(20) DEFAULT 'participant',
    joined_at TIMESTAMP DEFAULT NOW(),
    last_seen TIMESTAMP DEFAULT NOW(),
    permissions JSONB,
    metadata JSONB
);

-- Services table
CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    endpoint VARCHAR(500) NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    capabilities JSONB,
    ui_mapping JSONB,
    permissions JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    metadata JSONB
);

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    last_login TIMESTAMP,
    status VARCHAR(20) DEFAULT 'active',
    profile JSONB,
    metadata JSONB
);

-- Permissions table
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID NOT NULL,
    user_id UUID NOT NULL,
    permission_type VARCHAR(50) NOT NULL,
    granted_at TIMESTAMP DEFAULT NOW(),
    granted_by UUID,
    metadata JSONB
);
```

### Step 2: Session Management Implementation
```go
// src/session/manager.go
package session

import (
    "context"
    "database/sql"
    "time"
    
    "github.com/google/uuid"
    "holodeck1/database"
    "holodeck1/logging"
)

type SessionManager struct {
    db *database.DB
}

type Session struct {
    ID                 uuid.UUID              `json:"id"`
    Name               string                 `json:"name"`
    Description        string                 `json:"description"`
    OwnerID            uuid.UUID              `json:"owner_id"`
    CreatedAt          time.Time              `json:"created_at"`
    UpdatedAt          time.Time              `json:"updated_at"`
    Status             string                 `json:"status"`
    Visibility         string                 `json:"visibility"`
    MaxParticipants    int                    `json:"max_participants"`
    CurrentParticipants int                   `json:"current_participants"`
    Settings           map[string]interface{} `json:"settings"`
    Metadata           map[string]interface{} `json:"metadata"`
}

func NewSessionManager(db *database.DB) *SessionManager {
    return &SessionManager{db: db}
}

func (sm *SessionManager) CreateSession(ctx context.Context, req *CreateSessionRequest) (*Session, error) {
    session := &Session{
        ID:              uuid.New(),
        Name:            req.Name,
        Description:     req.Description,
        OwnerID:         req.OwnerID,
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
        Status:          "active",
        Visibility:      req.Visibility,
        MaxParticipants: req.MaxParticipants,
        Settings:        req.Settings,
        Metadata:        req.Metadata,
    }
    
    err := sm.db.CreateSession(ctx, session)
    if err != nil {
        logging.Error("failed to create session", map[string]interface{}{
            "error": err.Error(),
        })
        return nil, err
    }
    
    return session, nil
}

func (sm *SessionManager) GetSession(ctx context.Context, sessionID uuid.UUID) (*Session, error) {
    session, err := sm.db.GetSession(ctx, sessionID)
    if err != nil {
        logging.Error("failed to get session", map[string]interface{}{
            "session_id": sessionID,
            "error": err.Error(),
        })
        return nil, err
    }
    
    return session, nil
}

func (sm *SessionManager) ListSessions(ctx context.Context, filter *SessionFilter) ([]*Session, error) {
    sessions, err := sm.db.ListSessions(ctx, filter)
    if err != nil {
        logging.Error("failed to list sessions", map[string]interface{}{
            "error": err.Error(),
        })
        return nil, err
    }
    
    return sessions, nil
}

func (sm *SessionManager) UpdateSession(ctx context.Context, sessionID uuid.UUID, req *UpdateSessionRequest) (*Session, error) {
    session, err := sm.db.GetSession(ctx, sessionID)
    if err != nil {
        return nil, err
    }
    
    // Update fields
    if req.Name != "" {
        session.Name = req.Name
    }
    if req.Description != "" {
        session.Description = req.Description
    }
    if req.MaxParticipants > 0 {
        session.MaxParticipants = req.MaxParticipants
    }
    if req.Settings != nil {
        session.Settings = req.Settings
    }
    
    session.UpdatedAt = time.Now()
    
    err = sm.db.UpdateSession(ctx, session)
    if err != nil {
        logging.Error("failed to update session", map[string]interface{}{
            "session_id": sessionID,
            "error": err.Error(),
        })
        return nil, err
    }
    
    return session, nil
}

func (sm *SessionManager) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
    err := sm.db.DeleteSession(ctx, sessionID)
    if err != nil {
        logging.Error("failed to delete session", map[string]interface{}{
            "session_id": sessionID,
            "error": err.Error(),
        })
        return err
    }
    
    return nil
}

func (sm *SessionManager) JoinSession(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, role string) error {
    participant := &Participant{
        ID:        uuid.New(),
        SessionID: sessionID,
        UserID:    userID,
        Role:      role,
        JoinedAt:  time.Now(),
        LastSeen:  time.Now(),
    }
    
    err := sm.db.CreateParticipant(ctx, participant)
    if err != nil {
        logging.Error("failed to join session", map[string]interface{}{
            "session_id": sessionID,
            "user_id": userID,
            "error": err.Error(),
        })
        return err
    }
    
    // Update participant count
    err = sm.db.IncrementParticipantCount(ctx, sessionID)
    if err != nil {
        logging.Error("failed to update participant count", map[string]interface{}{
            "session_id": sessionID,
            "error": err.Error(),
        })
        return err
    }
    
    return nil
}
```

### Step 3: Service Registry Implementation
```go
// src/service/registry.go
package service

import (
    "context"
    "time"
    
    "github.com/google/uuid"
    "holodeck1/database"
    "holodeck1/logging"
)

type ServiceRegistry struct {
    db *database.DB
}

type Service struct {
    ID           uuid.UUID              `json:"id"`
    Name         string                 `json:"name"`
    Description  string                 `json:"description"`
    Type         string                 `json:"type"`
    Endpoint     string                 `json:"endpoint"`
    Status       string                 `json:"status"`
    Capabilities []string               `json:"capabilities"`
    UIMapping    map[string]interface{} `json:"ui_mapping"`
    Permissions  map[string]interface{} `json:"permissions"`
    CreatedAt    time.Time              `json:"created_at"`
    UpdatedAt    time.Time              `json:"updated_at"`
    Metadata     map[string]interface{} `json:"metadata"`
}

func NewServiceRegistry(db *database.DB) *ServiceRegistry {
    return &ServiceRegistry{db: db}
}

func (sr *ServiceRegistry) RegisterService(ctx context.Context, req *RegisterServiceRequest) (*Service, error) {
    service := &Service{
        ID:           uuid.New(),
        Name:         req.Name,
        Description:  req.Description,
        Type:         req.Type,
        Endpoint:     req.Endpoint,
        Status:       "active",
        Capabilities: req.Capabilities,
        UIMapping:    req.UIMapping,
        Permissions:  req.Permissions,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
        Metadata:     req.Metadata,
    }
    
    err := sr.db.CreateService(ctx, service)
    if err != nil {
        logging.Error("failed to register service", map[string]interface{}{
            "service_name": req.Name,
            "error": err.Error(),
        })
        return nil, err
    }
    
    return service, nil
}

func (sr *ServiceRegistry) GetService(ctx context.Context, serviceID uuid.UUID) (*Service, error) {
    service, err := sr.db.GetService(ctx, serviceID)
    if err != nil {
        logging.Error("failed to get service", map[string]interface{}{
            "service_id": serviceID,
            "error": err.Error(),
        })
        return nil, err
    }
    
    return service, nil
}

func (sr *ServiceRegistry) ListServices(ctx context.Context, filter *ServiceFilter) ([]*Service, error) {
    services, err := sr.db.ListServices(ctx, filter)
    if err != nil {
        logging.Error("failed to list services", map[string]interface{}{
            "error": err.Error(),
        })
        return nil, err
    }
    
    return services, nil
}

func (sr *ServiceRegistry) UpdateService(ctx context.Context, serviceID uuid.UUID, req *UpdateServiceRequest) (*Service, error) {
    service, err := sr.db.GetService(ctx, serviceID)
    if err != nil {
        return nil, err
    }
    
    // Update fields
    if req.Name != "" {
        service.Name = req.Name
    }
    if req.Description != "" {
        service.Description = req.Description
    }
    if req.Endpoint != "" {
        service.Endpoint = req.Endpoint
    }
    if req.Status != "" {
        service.Status = req.Status
    }
    if req.Capabilities != nil {
        service.Capabilities = req.Capabilities
    }
    if req.UIMapping != nil {
        service.UIMapping = req.UIMapping
    }
    
    service.UpdatedAt = time.Now()
    
    err = sr.db.UpdateService(ctx, service)
    if err != nil {
        logging.Error("failed to update service", map[string]interface{}{
            "service_id": serviceID,
            "error": err.Error(),
        })
        return nil, err
    }
    
    return service, nil
}

func (sr *ServiceRegistry) DeleteService(ctx context.Context, serviceID uuid.UUID) error {
    err := sr.db.DeleteService(ctx, serviceID)
    if err != nil {
        logging.Error("failed to delete service", map[string]interface{}{
            "service_id": serviceID,
            "error": err.Error(),
        })
        return err
    }
    
    return nil
}
```

### Step 4: API Endpoint Implementation
```go
// src/api/sessions/handlers.go
package sessions

import (
    "encoding/json"
    "net/http"
    
    "github.com/gorilla/mux"
    "github.com/google/uuid"
    "holodeck1/session"
    "holodeck1/logging"
)

type SessionHandler struct {
    sessionManager *session.SessionManager
}

func NewSessionHandler(sm *session.SessionManager) *SessionHandler {
    return &SessionHandler{sessionManager: sm}
}

func (h *SessionHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // Parse query parameters
    filter := &session.SessionFilter{
        Status:     r.URL.Query().Get("status"),
        Visibility: r.URL.Query().Get("visibility"),
        Limit:      50, // Default limit
        Offset:     0,  // Default offset
    }
    
    sessions, err := h.sessionManager.ListSessions(ctx, filter)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    response := map[string]interface{}{
        "success":  true,
        "sessions": sessions,
        "pagination": map[string]interface{}{
            "total":    len(sessions),
            "limit":    filter.Limit,
            "offset":   filter.Offset,
            "has_more": len(sessions) >= filter.Limit,
        },
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    var req session.CreateSessionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // TODO: Get user ID from authentication context
    req.OwnerID = uuid.New() // Placeholder
    
    session, err := h.sessionManager.CreateSession(ctx, &req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    response := map[string]interface{}{
        "success": true,
        "session": session,
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}

func (h *SessionHandler) GetSession(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    vars := mux.Vars(r)
    sessionID, err := uuid.Parse(vars["sessionId"])
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }
    
    session, err := h.sessionManager.GetSession(ctx, sessionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    response := map[string]interface{}{
        "success": true,
        "session": session,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *SessionHandler) UpdateSession(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    vars := mux.Vars(r)
    sessionID, err := uuid.Parse(vars["sessionId"])
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }
    
    var req session.UpdateSessionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    session, err := h.sessionManager.UpdateSession(ctx, sessionID, &req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    response := map[string]interface{}{
        "success": true,
        "session": session,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *SessionHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    vars := mux.Vars(r)
    sessionID, err := uuid.Parse(vars["sessionId"])
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }
    
    err = h.sessionManager.DeleteSession(ctx, sessionID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}

func (h *SessionHandler) JoinSession(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    vars := mux.Vars(r)
    sessionID, err := uuid.Parse(vars["sessionId"])
    if err != nil {
        http.Error(w, "Invalid session ID", http.StatusBadRequest)
        return
    }
    
    var req session.JoinSessionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // TODO: Get user ID from authentication context
    userID := uuid.New() // Placeholder
    
    err = h.sessionManager.JoinSession(ctx, sessionID, userID, req.Role)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    response := map[string]interface{}{
        "success": true,
        "message": "Successfully joined session",
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

## Implementation Timeline

### Week 1-2: Database Foundation
- [ ] Set up PostgreSQL and Redis
- [ ] Create database schema and migrations
- [ ] Implement database connection pooling
- [ ] Create data access layer

### Week 3-4: Session Management
- [ ] Implement SessionManager with full CRUD
- [ ] Add participant management
- [ ] Create session API endpoints
- [ ] Add session-based WebSocket routing

### Week 5-6: Service Registry
- [ ] Implement ServiceRegistry with full CRUD
- [ ] Add service health monitoring
- [ ] Create service API endpoints
- [ ] Add service integration framework

### Week 7-8: Authentication & Security
- [ ] Implement OAuth 2.0 server
- [ ] Add JWT token management
- [ ] Create authentication middleware
- [ ] Add rate limiting and security

## Success Criteria

### Technical Metrics
- [ ] All 19 new endpoints implemented and tested
- [ ] Database performance < 10ms for CRUD operations
- [ ] WebSocket routing supports 1000+ concurrent sessions
- [ ] Service registry handles 100+ registered services
- [ ] Authentication system supports OAuth 2.0 flows

### Quality Metrics
- [ ] 100% API test coverage
- [ ] Zero regressions in existing functionality
- [ ] Complete API documentation
- [ ] Load testing passes for 10,000 concurrent users
- [ ] Security audit passes

### Business Metrics
- [ ] Foundation supports unlimited sessions
- [ ] Service registry enables external integration
- [ ] Authentication enables user management
- [ ] Platform ready for Phase 2 collaboration features

## Risk Mitigation

### Technical Risks
- **Database Performance**: Implement connection pooling and indexing
- **WebSocket Scaling**: Use Redis pub/sub for message routing
- **Authentication Complexity**: Use proven OAuth 2.0 libraries
- **Service Integration**: Implement circuit breakers and timeouts

### Implementation Risks
- **Schema Changes**: Use database migrations with rollback capability
- **API Backwards Compatibility**: Maintain v6 endpoints during transition
- **Testing Coverage**: Implement comprehensive test suite
- **Performance Regression**: Continuous performance monitoring

## Deliverables

### Code Deliverables
- [ ] Multi-tenant session management system
- [ ] Service registry with health monitoring
- [ ] OAuth 2.0 authentication system
- [ ] Enhanced WebSocket hub with session routing
- [ ] Complete API endpoints with documentation

### Documentation Deliverables
- [ ] Updated API specification with new endpoints
- [ ] Database schema documentation
- [ ] Service integration guide
- [ ] Authentication and security guide
- [ ] Performance and scaling guide

### Testing Deliverables
- [ ] Unit tests for all new components
- [ ] Integration tests for API endpoints
- [ ] Load tests for concurrent users
- [ ] Security tests for authentication
- [ ] End-to-end tests for complete flows

This completes the detailed Phase 1 implementation plan. The foundation will transform HD1 from a simple Three.js engine to a scalable, multi-tenant platform ready for universal service integration.