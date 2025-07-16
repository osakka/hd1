package enterprise

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"holodeck1/database"
	"holodeck1/logging"
)

type SecurityManager struct {
	db *database.DB
}

type SecurityAuditLog struct {
	ID             uuid.UUID              `json:"id"`
	OrganizationID uuid.UUID              `json:"organization_id"`
	UserID         *uuid.UUID             `json:"user_id"`
	EventType      string                 `json:"event_type"`
	EventCategory  string                 `json:"event_category"`
	ResourceType   string                 `json:"resource_type"`
	ResourceID     *uuid.UUID             `json:"resource_id"`
	Action         string                 `json:"action"`
	Result         string                 `json:"result"`
	RiskLevel      string                 `json:"risk_level"`
	Details        map[string]interface{} `json:"details"`
	IPAddress      string                 `json:"ip_address"`
	UserAgent      string                 `json:"user_agent"`
	GeoLocation    map[string]interface{} `json:"geo_location"`
	Timestamp      time.Time              `json:"timestamp"`
	CreatedAt      time.Time              `json:"created_at"`
	Metadata       map[string]interface{} `json:"metadata"`
}

type ComplianceRecord struct {
	ID              uuid.UUID              `json:"id"`
	OrganizationID  uuid.UUID              `json:"organization_id"`
	ComplianceType  string                 `json:"compliance_type"`
	RecordType      string                 `json:"record_type"`
	SubjectType     string                 `json:"subject_type"`
	SubjectID       *uuid.UUID             `json:"subject_id"`
	Status          string                 `json:"status"`
	Priority        string                 `json:"priority"`
	Description     string                 `json:"description"`
	RequestDetails  map[string]interface{} `json:"request_details"`
	ResponseDetails map[string]interface{} `json:"response_details"`
	Evidence        map[string]interface{} `json:"evidence"`
	DueDate         *time.Time             `json:"due_date"`
	CompletedAt     *time.Time             `json:"completed_at"`
	CreatedBy       *uuid.UUID             `json:"created_by"`
	AssignedTo      *uuid.UUID             `json:"assigned_to"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	Metadata        map[string]interface{} `json:"metadata"`
}

type APIKey struct {
	ID             uuid.UUID              `json:"id"`
	OrganizationID uuid.UUID              `json:"organization_id"`
	UserID         uuid.UUID              `json:"user_id"`
	Name           string                 `json:"name"`
	Key            string                 `json:"key"`
	KeyHash        string                 `json:"-"`
	Permissions    []string               `json:"permissions"`
	RateLimit      int                    `json:"rate_limit"`
	ExpiresAt      *time.Time             `json:"expires_at"`
	LastUsedAt     *time.Time             `json:"last_used_at"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	Metadata       map[string]interface{} `json:"metadata"`
}

type SecurityEventRequest struct {
	UserID       *uuid.UUID             `json:"user_id"`
	EventType    string                 `json:"event_type"`
	EventCategory string                `json:"event_category"`
	ResourceType string                 `json:"resource_type"`
	ResourceID   *uuid.UUID             `json:"resource_id"`
	Action       string                 `json:"action"`
	Result       string                 `json:"result"`
	RiskLevel    string                 `json:"risk_level"`
	Details      map[string]interface{} `json:"details"`
	IPAddress    string                 `json:"ip_address"`
	UserAgent    string                 `json:"user_agent"`
}

type ComplianceRequest struct {
	ComplianceType  string                 `json:"compliance_type"`
	RecordType      string                 `json:"record_type"`
	SubjectType     string                 `json:"subject_type"`
	SubjectID       *uuid.UUID             `json:"subject_id"`
	Priority        string                 `json:"priority"`
	Description     string                 `json:"description"`
	RequestDetails  map[string]interface{} `json:"request_details"`
	DueDate         *time.Time             `json:"due_date"`
	CreatedBy       uuid.UUID              `json:"created_by"`
	AssignedTo      *uuid.UUID             `json:"assigned_to"`
}

type APIKeyCreateRequest struct {
	Name        string     `json:"name"`
	Permissions []string   `json:"permissions"`
	RateLimit   int        `json:"rate_limit"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

var SecurityEventCategories = []string{
	"authentication",
	"authorization",
	"data_access",
	"configuration",
	"network",
	"system",
}

var SecurityResultTypes = []string{
	"success",
	"failure",
	"warning",
}

var RiskLevels = []string{
	"low",
	"medium",
	"high",
	"critical",
}

var ComplianceTypes = []string{
	"gdpr",
	"hipaa",
	"sox",
	"pci_dss",
	"ccpa",
	"iso27001",
}

var ComplianceRecordTypes = []string{
	"data_processing",
	"access_request",
	"deletion_request",
	"audit",
	"breach_notification",
	"consent_management",
}

func NewSecurityManager(db *database.DB) *SecurityManager {
	return &SecurityManager{
		db: db,
	}
}

func (sm *SecurityManager) LogSecurityEvent(ctx context.Context, orgID uuid.UUID, req *SecurityEventRequest) error {
	// Validate required fields
	if req.EventType == "" {
		return fmt.Errorf("event type is required")
	}
	if req.EventCategory == "" {
		return fmt.Errorf("event category is required")
	}
	if req.Action == "" {
		return fmt.Errorf("action is required")
	}
	if req.Result == "" {
		return fmt.Errorf("result is required")
	}

	// Validate event category
	validCategory := false
	for _, cat := range SecurityEventCategories {
		if req.EventCategory == cat {
			validCategory = true
			break
		}
	}
	if !validCategory {
		return fmt.Errorf("invalid event category: %s", req.EventCategory)
	}

	// Validate result
	validResult := false
	for _, res := range SecurityResultTypes {
		if req.Result == res {
			validResult = true
			break
		}
	}
	if !validResult {
		return fmt.Errorf("invalid result: %s", req.Result)
	}

	// Set default risk level if not provided
	if req.RiskLevel == "" {
		req.RiskLevel = "low"
	}

	// Validate risk level
	validRisk := false
	for _, risk := range RiskLevels {
		if req.RiskLevel == risk {
			validRisk = true
			break
		}
	}
	if !validRisk {
		return fmt.Errorf("invalid risk level: %s", req.RiskLevel)
	}

	// Get geo location from IP address (simplified)
	geoLocation := sm.getGeoLocation(req.IPAddress)

	log := &SecurityAuditLog{
		ID:             uuid.New(),
		OrganizationID: orgID,
		UserID:         req.UserID,
		EventType:      req.EventType,
		EventCategory:  req.EventCategory,
		ResourceType:   req.ResourceType,
		ResourceID:     req.ResourceID,
		Action:         req.Action,
		Result:         req.Result,
		RiskLevel:      req.RiskLevel,
		Details:        req.Details,
		IPAddress:      req.IPAddress,
		UserAgent:      req.UserAgent,
		GeoLocation:    geoLocation,
		Timestamp:      time.Now(),
		CreatedAt:      time.Now(),
		Metadata:       make(map[string]interface{}),
	}

	if log.Details == nil {
		log.Details = make(map[string]interface{})
	}

	query := `
		INSERT INTO security_audit_log (id, organization_id, user_id, event_type, event_category, resource_type, resource_id, action, result, risk_level, details, ip_address, user_agent, geo_location, timestamp, created_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	detailsJSON, _ := json.Marshal(log.Details)
	geoJSON, _ := json.Marshal(log.GeoLocation)
	metadataJSON, _ := json.Marshal(log.Metadata)

	_, err := sm.db.Conn.ExecContext(ctx, query,
		log.ID, log.OrganizationID, log.UserID, log.EventType, log.EventCategory,
		log.ResourceType, log.ResourceID, log.Action, log.Result, log.RiskLevel,
		detailsJSON, log.IPAddress, log.UserAgent, geoJSON,
		log.Timestamp, log.CreatedAt, metadataJSON)
	if err != nil {
		return fmt.Errorf("failed to log security event: %w", err)
	}

	// Alert on high-risk events
	if req.RiskLevel == "high" || req.RiskLevel == "critical" {
		sm.alertHighRiskEvent(ctx, log)
	}

	logging.Info("logged security event", map[string]interface{}{
		"event_id":        log.ID,
		"organization_id": orgID,
		"event_type":      log.EventType,
		"event_category":  log.EventCategory,
		"action":          log.Action,
		"result":          log.Result,
		"risk_level":      log.RiskLevel,
		"user_id":         log.UserID,
	})

	return nil
}

func (sm *SecurityManager) CreateComplianceRecord(ctx context.Context, orgID uuid.UUID, req *ComplianceRequest) (*ComplianceRecord, error) {
	// Validate required fields
	if req.ComplianceType == "" {
		return nil, fmt.Errorf("compliance type is required")
	}
	if req.RecordType == "" {
		return nil, fmt.Errorf("record type is required")
	}
	if req.SubjectType == "" {
		return nil, fmt.Errorf("subject type is required")
	}
	if req.Description == "" {
		return nil, fmt.Errorf("description is required")
	}

	// Validate compliance type
	validCompliance := false
	for _, comp := range ComplianceTypes {
		if req.ComplianceType == comp {
			validCompliance = true
			break
		}
	}
	if !validCompliance {
		return nil, fmt.Errorf("invalid compliance type: %s", req.ComplianceType)
	}

	// Validate record type
	validRecord := false
	for _, rec := range ComplianceRecordTypes {
		if req.RecordType == rec {
			validRecord = true
			break
		}
	}
	if !validRecord {
		return nil, fmt.Errorf("invalid record type: %s", req.RecordType)
	}

	// Set default priority if not provided
	if req.Priority == "" {
		req.Priority = "medium"
	}

	record := &ComplianceRecord{
		ID:              uuid.New(),
		OrganizationID:  orgID,
		ComplianceType:  req.ComplianceType,
		RecordType:      req.RecordType,
		SubjectType:     req.SubjectType,
		SubjectID:       req.SubjectID,
		Status:          "pending",
		Priority:        req.Priority,
		Description:     req.Description,
		RequestDetails:  req.RequestDetails,
		ResponseDetails: make(map[string]interface{}),
		Evidence:        make(map[string]interface{}),
		DueDate:         req.DueDate,
		CreatedBy:       &req.CreatedBy,
		AssignedTo:      req.AssignedTo,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Metadata:        make(map[string]interface{}),
	}

	if record.RequestDetails == nil {
		record.RequestDetails = make(map[string]interface{})
	}

	query := `
		INSERT INTO compliance_records (id, organization_id, compliance_type, record_type, subject_type, subject_id, status, priority, description, request_details, response_details, evidence, due_date, completed_at, created_by, assigned_to, created_at, updated_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`

	requestJSON, _ := json.Marshal(record.RequestDetails)
	responseJSON, _ := json.Marshal(record.ResponseDetails)
	evidenceJSON, _ := json.Marshal(record.Evidence)
	metadataJSON, _ := json.Marshal(record.Metadata)

	_, err := sm.db.Conn.ExecContext(ctx, query,
		record.ID, record.OrganizationID, record.ComplianceType, record.RecordType,
		record.SubjectType, record.SubjectID, record.Status, record.Priority,
		record.Description, requestJSON, responseJSON, evidenceJSON,
		record.DueDate, record.CompletedAt, record.CreatedBy, record.AssignedTo,
		record.CreatedAt, record.UpdatedAt, metadataJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to create compliance record: %w", err)
	}

	logging.Info("created compliance record", map[string]interface{}{
		"record_id":       record.ID,
		"organization_id": orgID,
		"compliance_type": record.ComplianceType,
		"record_type":     record.RecordType,
		"subject_type":    record.SubjectType,
		"priority":        record.Priority,
		"created_by":      record.CreatedBy,
	})

	return record, nil
}

func (sm *SecurityManager) UpdateComplianceRecord(ctx context.Context, recordID uuid.UUID, status string, responseDetails map[string]interface{}, evidence map[string]interface{}) error {
	// Validate status
	validStatuses := []string{"pending", "in_progress", "completed", "rejected", "cancelled"}
	validStatus := false
	for _, s := range validStatuses {
		if status == s {
			validStatus = true
			break
		}
	}
	if !validStatus {
		return fmt.Errorf("invalid status: %s", status)
	}

	var completedAt *time.Time
	if status == "completed" || status == "rejected" || status == "cancelled" {
		now := time.Now()
		completedAt = &now
	}

	query := `
		UPDATE compliance_records 
		SET status = $1, response_details = $2, evidence = $3, completed_at = $4, updated_at = $5
		WHERE id = $6
	`

	responseJSON, _ := json.Marshal(responseDetails)
	evidenceJSON, _ := json.Marshal(evidence)

	result, err := sm.db.Conn.ExecContext(ctx, query,
		status, responseJSON, evidenceJSON, completedAt, time.Now(), recordID)
	if err != nil {
		return fmt.Errorf("failed to update compliance record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("compliance record not found")
	}

	logging.Info("updated compliance record", map[string]interface{}{
		"record_id":    recordID,
		"status":       status,
		"completed_at": completedAt,
	})

	return nil
}

func (sm *SecurityManager) CreateAPIKey(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, req *APIKeyCreateRequest) (*APIKey, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, fmt.Errorf("API key name is required")
	}

	// Generate secure API key
	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	keyString := base64.URLEncoding.EncodeToString(keyBytes)
	keyHash := sm.hashAPIKey(keyString)

	// Set default rate limit if not provided
	if req.RateLimit == 0 {
		req.RateLimit = 1000 // requests per hour
	}

	apiKey := &APIKey{
		ID:             uuid.New(),
		OrganizationID: orgID,
		UserID:         userID,
		Name:           req.Name,
		Key:            keyString,
		KeyHash:        keyHash,
		Permissions:    req.Permissions,
		RateLimit:      req.RateLimit,
		ExpiresAt:      req.ExpiresAt,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Metadata:       make(map[string]interface{}),
	}

	if apiKey.Permissions == nil {
		apiKey.Permissions = []string{}
	}

	// Store in database (store hash, not the actual key)
	query := `
		INSERT INTO api_keys (id, organization_id, user_id, name, key_hash, permissions, rate_limit, expires_at, created_at, updated_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	permissionsJSON, _ := json.Marshal(apiKey.Permissions)
	metadataJSON, _ := json.Marshal(apiKey.Metadata)

	_, err := sm.db.Conn.ExecContext(ctx, query,
		apiKey.ID, apiKey.OrganizationID, apiKey.UserID, apiKey.Name,
		apiKey.KeyHash, permissionsJSON, apiKey.RateLimit, apiKey.ExpiresAt,
		apiKey.CreatedAt, apiKey.UpdatedAt, metadataJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	// Log API key creation
	sm.LogSecurityEvent(ctx, orgID, &SecurityEventRequest{
		UserID:        &userID,
		EventType:     "api_key_created",
		EventCategory: "authentication",
		ResourceType:  "api_key",
		ResourceID:    &apiKey.ID,
		Action:        "create",
		Result:        "success",
		RiskLevel:     "low",
		Details: map[string]interface{}{
			"key_name":    apiKey.Name,
			"permissions": apiKey.Permissions,
			"rate_limit":  apiKey.RateLimit,
		},
	})

	logging.Info("created API key", map[string]interface{}{
		"api_key_id":      apiKey.ID,
		"organization_id": orgID,
		"user_id":         userID,
		"name":            apiKey.Name,
		"permissions":     len(apiKey.Permissions),
	})

	return apiKey, nil
}

func (sm *SecurityManager) RevokeAPIKey(ctx context.Context, orgID uuid.UUID, keyID uuid.UUID) error {
	query := `DELETE FROM api_keys WHERE id = $1 AND organization_id = $2`

	result, err := sm.db.Conn.ExecContext(ctx, query, keyID, orgID)
	if err != nil {
		return fmt.Errorf("failed to revoke API key: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("API key not found")
	}

	// Log API key revocation
	sm.LogSecurityEvent(ctx, orgID, &SecurityEventRequest{
		EventType:     "api_key_revoked",
		EventCategory: "authentication",
		ResourceType:  "api_key",
		ResourceID:    &keyID,
		Action:        "revoke",
		Result:        "success",
		RiskLevel:     "medium",
	})

	logging.Info("revoked API key", map[string]interface{}{
		"api_key_id":      keyID,
		"organization_id": orgID,
	})

	return nil
}

func (sm *SecurityManager) GetSecurityEvents(ctx context.Context, orgID uuid.UUID, startTime, endTime time.Time, filters map[string]string) ([]*SecurityAuditLog, error) {
	query := `
		SELECT id, organization_id, user_id, event_type, event_category, resource_type, resource_id, action, result, risk_level, details, ip_address, user_agent, geo_location, timestamp, created_at, metadata
		FROM security_audit_log
		WHERE organization_id = $1 AND timestamp >= $2 AND timestamp <= $3
	`
	args := []interface{}{orgID, startTime, endTime}
	argIndex := 4

	// Add optional filters
	if category, ok := filters["event_category"]; ok && category != "" {
		query += fmt.Sprintf(" AND event_category = $%d", argIndex)
		args = append(args, category)
		argIndex++
	}

	if riskLevel, ok := filters["risk_level"]; ok && riskLevel != "" {
		query += fmt.Sprintf(" AND risk_level = $%d", argIndex)
		args = append(args, riskLevel)
		argIndex++
	}

	if result, ok := filters["result"]; ok && result != "" {
		query += fmt.Sprintf(" AND result = $%d", argIndex)
		args = append(args, result)
		argIndex++
	}

	query += " ORDER BY timestamp DESC LIMIT 1000"

	rows, err := sm.db.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get security events: %w", err)
	}
	defer rows.Close()

	var events []*SecurityAuditLog
	for rows.Next() {
		var event SecurityAuditLog
		var userID, resourceID sql.NullString
		var detailsJSON, geoJSON, metadataJSON []byte

		err := rows.Scan(
			&event.ID, &event.OrganizationID, &userID, &event.EventType,
			&event.EventCategory, &event.ResourceType, &resourceID,
			&event.Action, &event.Result, &event.RiskLevel,
			&detailsJSON, &event.IPAddress, &event.UserAgent, &geoJSON,
			&event.Timestamp, &event.CreatedAt, &metadataJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan security event: %w", err)
		}

		// Handle nullable fields
		if userID.Valid {
			uid, _ := uuid.Parse(userID.String)
			event.UserID = &uid
		}
		if resourceID.Valid {
			rid, _ := uuid.Parse(resourceID.String)
			event.ResourceID = &rid
		}

		// Unmarshal JSON fields
		json.Unmarshal(detailsJSON, &event.Details)
		json.Unmarshal(geoJSON, &event.GeoLocation)
		json.Unmarshal(metadataJSON, &event.Metadata)

		events = append(events, &event)
	}

	return events, nil
}

// Helper functions

func (sm *SecurityManager) getGeoLocation(ipAddress string) map[string]interface{} {
	// This is a simplified implementation
	// In production, you would use a real IP geolocation service
	
	geo := make(map[string]interface{})
	
	// Check if it's a private IP
	ip := net.ParseIP(ipAddress)
	if ip != nil && ip.IsPrivate() {
		geo["country"] = "Private Network"
		geo["city"] = "Local"
		geo["lat"] = 0.0
		geo["lon"] = 0.0
		return geo
	}
	
	// Default to unknown
	geo["country"] = "Unknown"
	geo["city"] = "Unknown"
	geo["lat"] = 0.0
	geo["lon"] = 0.0
	
	return geo
}

func (sm *SecurityManager) hashAPIKey(key string) string {
	// In production, use a proper hashing algorithm like bcrypt
	// This is simplified for demonstration
	return fmt.Sprintf("sha256:%x", strings.ToUpper(key))
}

func (sm *SecurityManager) alertHighRiskEvent(ctx context.Context, event *SecurityAuditLog) {
	// In production, this would send alerts via email, Slack, PagerDuty, etc.
	logging.Warn("HIGH RISK SECURITY EVENT", map[string]interface{}{
		"event_id":        event.ID,
		"organization_id": event.OrganizationID,
		"event_type":      event.EventType,
		"action":          event.Action,
		"risk_level":      event.RiskLevel,
		"user_id":         event.UserID,
		"ip_address":      event.IPAddress,
		"details":         event.Details,
	})
}