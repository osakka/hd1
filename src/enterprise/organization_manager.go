package enterprise

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"holodeck1/database"
	"holodeck1/logging"
)

type OrganizationManager struct {
	db *database.DB
}

type Organization struct {
	ID               uuid.UUID              `json:"id"`
	Name             string                 `json:"name"`
	DisplayName      string                 `json:"display_name"`
	Description      string                 `json:"description"`
	Domain           string                 `json:"domain"`
	LogoURL          string                 `json:"logo_url"`
	Website          string                 `json:"website"`
	Status           string                 `json:"status"`
	SubscriptionTier string                 `json:"subscription_tier"`
	MaxUsers         int                    `json:"max_users"`
	MaxSessions      int                    `json:"max_sessions"`
	MaxStorageGB     int                    `json:"max_storage_gb"`
	Features         map[string]interface{} `json:"features"`
	Settings         map[string]interface{} `json:"settings"`
	BillingInfo      map[string]interface{} `json:"billing_info"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	Metadata         map[string]interface{} `json:"metadata"`
}

type OrganizationUser struct {
	ID             uuid.UUID              `json:"id"`
	OrganizationID uuid.UUID              `json:"organization_id"`
	UserID         uuid.UUID              `json:"user_id"`
	Role           string                 `json:"role"`
	Permissions    []string               `json:"permissions"`
	Department     string                 `json:"department"`
	JobTitle       string                 `json:"job_title"`
	Status         string                 `json:"status"`
	InvitedBy      *uuid.UUID             `json:"invited_by"`
	InvitedAt      *time.Time             `json:"invited_at"`
	JoinedAt       time.Time              `json:"joined_at"`
	LastActive     *time.Time             `json:"last_active"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	Metadata       map[string]interface{} `json:"metadata"`
}

type OrganizationCreateRequest struct {
	Name             string                 `json:"name"`
	DisplayName      string                 `json:"display_name"`
	Description      string                 `json:"description"`
	Domain           string                 `json:"domain"`
	LogoURL          string                 `json:"logo_url"`
	Website          string                 `json:"website"`
	SubscriptionTier string                 `json:"subscription_tier"`
	MaxUsers         int                    `json:"max_users"`
	MaxSessions      int                    `json:"max_sessions"`
	MaxStorageGB     int                    `json:"max_storage_gb"`
	Features         map[string]interface{} `json:"features"`
	Settings         map[string]interface{} `json:"settings"`
	BillingInfo      map[string]interface{} `json:"billing_info"`
}

type OrganizationUpdateRequest struct {
	DisplayName      *string                 `json:"display_name"`
	Description      *string                 `json:"description"`
	Domain           *string                 `json:"domain"`
	LogoURL          *string                 `json:"logo_url"`
	Website          *string                 `json:"website"`
	Status           *string                 `json:"status"`
	SubscriptionTier *string                 `json:"subscription_tier"`
	MaxUsers         *int                    `json:"max_users"`
	MaxSessions      *int                    `json:"max_sessions"`
	MaxStorageGB     *int                    `json:"max_storage_gb"`
	Features         *map[string]interface{} `json:"features"`
	Settings         *map[string]interface{} `json:"settings"`
	BillingInfo      *map[string]interface{} `json:"billing_info"`
}

type OrganizationFilter struct {
	Status           string `json:"status"`
	SubscriptionTier string `json:"subscription_tier"`
	Domain           string `json:"domain"`
	Limit            int    `json:"limit"`
	Offset           int    `json:"offset"`
}

type UserInviteRequest struct {
	UserID     uuid.UUID `json:"user_id"`
	Role       string    `json:"role"`
	Department string    `json:"department"`
	JobTitle   string    `json:"job_title"`
	InvitedBy  uuid.UUID `json:"invited_by"`
}

func NewOrganizationManager(db *database.DB) *OrganizationManager {
	return &OrganizationManager{
		db: db,
	}
}

func (om *OrganizationManager) CreateOrganization(ctx context.Context, req *OrganizationCreateRequest) (*Organization, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, fmt.Errorf("organization name is required")
	}
	if req.DisplayName == "" {
		return nil, fmt.Errorf("organization display name is required")
	}

	// Set defaults
	if req.SubscriptionTier == "" {
		req.SubscriptionTier = "basic"
	}
	if req.MaxUsers == 0 {
		req.MaxUsers = 100
	}
	if req.MaxSessions == 0 {
		req.MaxSessions = 50
	}
	if req.MaxStorageGB == 0 {
		req.MaxStorageGB = 10
	}
	if req.Features == nil {
		req.Features = make(map[string]interface{})
	}
	if req.Settings == nil {
		req.Settings = make(map[string]interface{})
	}
	if req.BillingInfo == nil {
		req.BillingInfo = make(map[string]interface{})
	}

	org := &Organization{
		ID:               uuid.New(),
		Name:             req.Name,
		DisplayName:      req.DisplayName,
		Description:      req.Description,
		Domain:           req.Domain,
		LogoURL:          req.LogoURL,
		Website:          req.Website,
		Status:           "active",
		SubscriptionTier: req.SubscriptionTier,
		MaxUsers:         req.MaxUsers,
		MaxSessions:      req.MaxSessions,
		MaxStorageGB:     req.MaxStorageGB,
		Features:         req.Features,
		Settings:         req.Settings,
		BillingInfo:      req.BillingInfo,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Metadata:         make(map[string]interface{}),
	}

	// Store in database
	query := `
		INSERT INTO organizations (id, name, display_name, description, domain, logo_url, website, status, subscription_tier, max_users, max_sessions, max_storage_gb, features, settings, billing_info, created_at, updated_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	featuresJSON, _ := json.Marshal(org.Features)
	settingsJSON, _ := json.Marshal(org.Settings)
	billingJSON, _ := json.Marshal(org.BillingInfo)
	metadataJSON, _ := json.Marshal(org.Metadata)

	_, err := om.db.Conn.ExecContext(ctx, query,
		org.ID, org.Name, org.DisplayName, org.Description, org.Domain,
		org.LogoURL, org.Website, org.Status, org.SubscriptionTier,
		org.MaxUsers, org.MaxSessions, org.MaxStorageGB,
		featuresJSON, settingsJSON, billingJSON,
		org.CreatedAt, org.UpdatedAt, metadataJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	logging.Info("created organization", map[string]interface{}{
		"organization_id": org.ID,
		"name":            org.Name,
		"display_name":    org.DisplayName,
		"domain":          org.Domain,
		"tier":            org.SubscriptionTier,
	})

	return org, nil
}

func (om *OrganizationManager) GetOrganization(ctx context.Context, orgID uuid.UUID) (*Organization, error) {
	query := `
		SELECT id, name, display_name, description, domain, logo_url, website, status, subscription_tier, max_users, max_sessions, max_storage_gb, features, settings, billing_info, created_at, updated_at, metadata
		FROM organizations
		WHERE id = $1
	`

	var org Organization
	var featuresJSON, settingsJSON, billingJSON, metadataJSON []byte

	err := om.db.Conn.QueryRowContext(ctx, query, orgID).Scan(
		&org.ID, &org.Name, &org.DisplayName, &org.Description, &org.Domain,
		&org.LogoURL, &org.Website, &org.Status, &org.SubscriptionTier,
		&org.MaxUsers, &org.MaxSessions, &org.MaxStorageGB,
		&featuresJSON, &settingsJSON, &billingJSON,
		&org.CreatedAt, &org.UpdatedAt, &metadataJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	// Unmarshal JSON fields
	json.Unmarshal(featuresJSON, &org.Features)
	json.Unmarshal(settingsJSON, &org.Settings)
	json.Unmarshal(billingJSON, &org.BillingInfo)
	json.Unmarshal(metadataJSON, &org.Metadata)

	return &org, nil
}

func (om *OrganizationManager) UpdateOrganization(ctx context.Context, orgID uuid.UUID, req *OrganizationUpdateRequest) error {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.DisplayName != nil {
		setParts = append(setParts, fmt.Sprintf("display_name = $%d", argIndex))
		args = append(args, *req.DisplayName)
		argIndex++
	}

	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}

	if req.Domain != nil {
		setParts = append(setParts, fmt.Sprintf("domain = $%d", argIndex))
		args = append(args, *req.Domain)
		argIndex++
	}

	if req.LogoURL != nil {
		setParts = append(setParts, fmt.Sprintf("logo_url = $%d", argIndex))
		args = append(args, *req.LogoURL)
		argIndex++
	}

	if req.Website != nil {
		setParts = append(setParts, fmt.Sprintf("website = $%d", argIndex))
		args = append(args, *req.Website)
		argIndex++
	}

	if req.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *req.Status)
		argIndex++
	}

	if req.SubscriptionTier != nil {
		setParts = append(setParts, fmt.Sprintf("subscription_tier = $%d", argIndex))
		args = append(args, *req.SubscriptionTier)
		argIndex++
	}

	if req.MaxUsers != nil {
		setParts = append(setParts, fmt.Sprintf("max_users = $%d", argIndex))
		args = append(args, *req.MaxUsers)
		argIndex++
	}

	if req.MaxSessions != nil {
		setParts = append(setParts, fmt.Sprintf("max_sessions = $%d", argIndex))
		args = append(args, *req.MaxSessions)
		argIndex++
	}

	if req.MaxStorageGB != nil {
		setParts = append(setParts, fmt.Sprintf("max_storage_gb = $%d", argIndex))
		args = append(args, *req.MaxStorageGB)
		argIndex++
	}

	if req.Features != nil {
		featuresJSON, _ := json.Marshal(*req.Features)
		setParts = append(setParts, fmt.Sprintf("features = $%d", argIndex))
		args = append(args, featuresJSON)
		argIndex++
	}

	if req.Settings != nil {
		settingsJSON, _ := json.Marshal(*req.Settings)
		setParts = append(setParts, fmt.Sprintf("settings = $%d", argIndex))
		args = append(args, settingsJSON)
		argIndex++
	}

	if req.BillingInfo != nil {
		billingJSON, _ := json.Marshal(*req.BillingInfo)
		setParts = append(setParts, fmt.Sprintf("billing_info = $%d", argIndex))
		args = append(args, billingJSON)
		argIndex++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Add updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add WHERE clause
	args = append(args, orgID)

	query := fmt.Sprintf("UPDATE organizations SET %s WHERE id = $%d", 
		fmt.Sprintf("%s", setParts), argIndex)

	result, err := om.db.Conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update organization: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("organization not found")
	}

	logging.Info("updated organization", map[string]interface{}{
		"organization_id": orgID,
		"fields_updated":  len(setParts) - 1, // exclude updated_at
	})

	return nil
}

func (om *OrganizationManager) DeleteOrganization(ctx context.Context, orgID uuid.UUID) error {
	query := `DELETE FROM organizations WHERE id = $1`

	result, err := om.db.Conn.ExecContext(ctx, query, orgID)
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("organization not found")
	}

	logging.Info("deleted organization", map[string]interface{}{
		"organization_id": orgID,
	})

	return nil
}

func (om *OrganizationManager) ListOrganizations(ctx context.Context, filter *OrganizationFilter) ([]*Organization, error) {
	query := `
		SELECT id, name, display_name, description, domain, logo_url, website, status, subscription_tier, max_users, max_sessions, max_storage_gb, features, settings, billing_info, created_at, updated_at, metadata
		FROM organizations
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	if filter.SubscriptionTier != "" {
		query += fmt.Sprintf(" AND subscription_tier = $%d", argIndex)
		args = append(args, filter.SubscriptionTier)
		argIndex++
	}

	if filter.Domain != "" {
		query += fmt.Sprintf(" AND domain = $%d", argIndex)
		args = append(args, filter.Domain)
		argIndex++
	}

	query += " ORDER BY name ASC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	rows, err := om.db.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}
	defer rows.Close()

	var organizations []*Organization
	for rows.Next() {
		var org Organization
		var featuresJSON, settingsJSON, billingJSON, metadataJSON []byte

		err := rows.Scan(
			&org.ID, &org.Name, &org.DisplayName, &org.Description, &org.Domain,
			&org.LogoURL, &org.Website, &org.Status, &org.SubscriptionTier,
			&org.MaxUsers, &org.MaxSessions, &org.MaxStorageGB,
			&featuresJSON, &settingsJSON, &billingJSON,
			&org.CreatedAt, &org.UpdatedAt, &metadataJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan organization: %w", err)
		}

		// Unmarshal JSON fields
		json.Unmarshal(featuresJSON, &org.Features)
		json.Unmarshal(settingsJSON, &org.Settings)
		json.Unmarshal(billingJSON, &org.BillingInfo)
		json.Unmarshal(metadataJSON, &org.Metadata)

		organizations = append(organizations, &org)
	}

	return organizations, nil
}

func (om *OrganizationManager) InviteUser(ctx context.Context, orgID uuid.UUID, req *UserInviteRequest) (*OrganizationUser, error) {
	// Check if user is already in organization
	existingQuery := `SELECT id FROM organization_users WHERE organization_id = $1 AND user_id = $2`
	var existingID uuid.UUID
	err := om.db.Conn.QueryRowContext(ctx, existingQuery, orgID, req.UserID).Scan(&existingID)
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("user already in organization")
	}

	// Create organization user record
	orgUser := &OrganizationUser{
		ID:             uuid.New(),
		OrganizationID: orgID,
		UserID:         req.UserID,
		Role:           req.Role,
		Department:     req.Department,
		JobTitle:       req.JobTitle,
		Status:         "invited",
		InvitedBy:      &req.InvitedBy,
		InvitedAt:      &time.Time{},
		JoinedAt:       time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Metadata:       make(map[string]interface{}),
	}

	*orgUser.InvitedAt = time.Now()

	query := `
		INSERT INTO organization_users (id, organization_id, user_id, role, department, job_title, status, invited_by, invited_at, joined_at, created_at, updated_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	metadataJSON, _ := json.Marshal(orgUser.Metadata)

	_, err = om.db.Conn.ExecContext(ctx, query,
		orgUser.ID, orgUser.OrganizationID, orgUser.UserID, orgUser.Role,
		orgUser.Department, orgUser.JobTitle, orgUser.Status,
		orgUser.InvitedBy, orgUser.InvitedAt, orgUser.JoinedAt,
		orgUser.CreatedAt, orgUser.UpdatedAt, metadataJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to invite user: %w", err)
	}

	logging.Info("invited user to organization", map[string]interface{}{
		"organization_id": orgID,
		"user_id":         req.UserID,
		"role":            req.Role,
		"invited_by":      req.InvitedBy,
	})

	return orgUser, nil
}

func (om *OrganizationManager) GetOrganizationUsers(ctx context.Context, orgID uuid.UUID) ([]*OrganizationUser, error) {
	query := `
		SELECT id, organization_id, user_id, role, permissions, department, job_title, status, invited_by, invited_at, joined_at, last_active, created_at, updated_at, metadata
		FROM organization_users
		WHERE organization_id = $1
		ORDER BY joined_at DESC
	`

	rows, err := om.db.Conn.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization users: %w", err)
	}
	defer rows.Close()

	var users []*OrganizationUser
	for rows.Next() {
		var user OrganizationUser
		var permissionsJSON, metadataJSON []byte
		var invitedBy sql.NullString
		var invitedAt, lastActive sql.NullTime

		err := rows.Scan(
			&user.ID, &user.OrganizationID, &user.UserID, &user.Role,
			&permissionsJSON, &user.Department, &user.JobTitle, &user.Status,
			&invitedBy, &invitedAt, &user.JoinedAt, &lastActive,
			&user.CreatedAt, &user.UpdatedAt, &metadataJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan organization user: %w", err)
		}

		// Handle nullable fields
		if invitedBy.Valid {
			invitedByUUID, _ := uuid.Parse(invitedBy.String)
			user.InvitedBy = &invitedByUUID
		}
		if invitedAt.Valid {
			user.InvitedAt = &invitedAt.Time
		}
		if lastActive.Valid {
			user.LastActive = &lastActive.Time
		}

		// Unmarshal JSON fields
		json.Unmarshal(permissionsJSON, &user.Permissions)
		json.Unmarshal(metadataJSON, &user.Metadata)

		users = append(users, &user)
	}

	return users, nil
}