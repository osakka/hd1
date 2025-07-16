package enterprise

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"holodeck1/database"
	"holodeck1/logging"
)

type RBACManager struct {
	db *database.DB
}

type Role struct {
	ID           uuid.UUID              `json:"id"`
	OrganizationID uuid.UUID            `json:"organization_id"`
	Name         string                 `json:"name"`
	DisplayName  string                 `json:"display_name"`
	Description  string                 `json:"description"`
	Permissions  []string               `json:"permissions"`
	IsSystemRole bool                   `json:"is_system_role"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	Metadata     map[string]interface{} `json:"metadata"`
}

type UserRole struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	RoleID         uuid.UUID  `json:"role_id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	GrantedBy      *uuid.UUID `json:"granted_by"`
	GrantedAt      time.Time  `json:"granted_at"`
	ExpiresAt      *time.Time `json:"expires_at"`
	CreatedAt      time.Time  `json:"created_at"`
	Metadata       map[string]interface{} `json:"metadata"`
}

type Permission struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
	Scope    string `json:"scope"`
}

type RoleCreateRequest struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	Name          string    `json:"name"`
	DisplayName   string    `json:"display_name"`
	Description   string    `json:"description"`
	Permissions   []string  `json:"permissions"`
}

type RoleUpdateRequest struct {
	DisplayName *string   `json:"display_name"`
	Description *string   `json:"description"`
	Permissions *[]string `json:"permissions"`
}

type UserRoleAssignment struct {
	UserID    uuid.UUID  `json:"user_id"`
	RoleID    uuid.UUID  `json:"role_id"`
	GrantedBy uuid.UUID  `json:"granted_by"`
	ExpiresAt *time.Time `json:"expires_at"`
}

var SystemPermissions = []string{
	// Organization Management
	"organizations:read",
	"organizations:write",
	"organizations:delete",
	"organizations:manage_users",
	"organizations:manage_billing",
	
	// User Management
	"users:read",
	"users:write",
	"users:delete",
	"users:manage_roles",
	
	// Session Management
	"sessions:read",
	"sessions:write",
	"sessions:delete",
	"sessions:join",
	"sessions:manage_participants",
	
	// Content Management
	"content:read",
	"content:write",
	"content:delete",
	"content:publish",
	
	// Analytics
	"analytics:read",
	"analytics:export",
	
	// Security
	"security:read_audit_logs",
	"security:manage_compliance",
	"security:manage_api_keys",
	
	// Plugins
	"plugins:read",
	"plugins:install",
	"plugins:manage",
}

var DefaultRoles = map[string][]string{
	"owner": {
		"organizations:*",
		"users:*",
		"sessions:*",
		"content:*",
		"analytics:*",
		"security:*",
		"plugins:*",
	},
	"admin": {
		"organizations:read",
		"organizations:write",
		"organizations:manage_users",
		"users:*",
		"sessions:*",
		"content:*",
		"analytics:*",
		"security:read_audit_logs",
		"plugins:*",
	},
	"manager": {
		"organizations:read",
		"users:read",
		"users:write",
		"users:manage_roles",
		"sessions:*",
		"content:*",
		"analytics:read",
		"analytics:export",
	},
	"member": {
		"organizations:read",
		"users:read",
		"sessions:read",
		"sessions:write",
		"sessions:join",
		"content:read",
		"content:write",
		"analytics:read",
	},
	"viewer": {
		"organizations:read",
		"users:read",
		"sessions:read",
		"sessions:join",
		"content:read",
		"analytics:read",
	},
}

func NewRBACManager(db *database.DB) *RBACManager {
	return &RBACManager{
		db: db,
	}
}

func (rm *RBACManager) InitializeDefaultRoles(ctx context.Context, orgID uuid.UUID) error {
	for roleName, permissions := range DefaultRoles {
		role := &Role{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Name:           roleName,
			DisplayName:    capitalizeFirst(roleName),
			Description:    fmt.Sprintf("Default %s role", roleName),
			Permissions:    permissions,
			IsSystemRole:   true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			Metadata:       make(map[string]interface{}),
		}

		query := `
			INSERT INTO roles (id, organization_id, name, display_name, description, permissions, is_system_role, created_at, updated_at, metadata)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			ON CONFLICT (organization_id, name) DO NOTHING
		`

		permissionsJSON, _ := json.Marshal(role.Permissions)
		metadataJSON, _ := json.Marshal(role.Metadata)

		_, err := rm.db.Conn.ExecContext(ctx, query,
			role.ID, role.OrganizationID, role.Name, role.DisplayName, role.Description,
			permissionsJSON, role.IsSystemRole, role.CreatedAt, role.UpdatedAt, metadataJSON)
		if err != nil {
			return fmt.Errorf("failed to initialize role %s: %w", roleName, err)
		}
	}

	logging.Info("initialized default roles", map[string]interface{}{
		"organization_id": orgID,
		"roles_count":     len(DefaultRoles),
	})

	return nil
}

func (rm *RBACManager) CreateRole(ctx context.Context, req *RoleCreateRequest) (*Role, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, fmt.Errorf("role name is required")
	}
	if req.DisplayName == "" {
		return nil, fmt.Errorf("role display name is required")
	}

	// Validate permissions
	for _, perm := range req.Permissions {
		if !isValidPermission(perm) {
			return nil, fmt.Errorf("invalid permission: %s", perm)
		}
	}

	role := &Role{
		ID:             uuid.New(),
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		DisplayName:    req.DisplayName,
		Description:    req.Description,
		Permissions:    req.Permissions,
		IsSystemRole:   false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Metadata:       make(map[string]interface{}),
	}

	query := `
		INSERT INTO roles (id, organization_id, name, display_name, description, permissions, is_system_role, created_at, updated_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	permissionsJSON, _ := json.Marshal(role.Permissions)
	metadataJSON, _ := json.Marshal(role.Metadata)

	_, err := rm.db.Conn.ExecContext(ctx, query,
		role.ID, role.OrganizationID, role.Name, role.DisplayName, role.Description,
		permissionsJSON, role.IsSystemRole, role.CreatedAt, role.UpdatedAt, metadataJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	logging.Info("created role", map[string]interface{}{
		"role_id":         role.ID,
		"organization_id": role.OrganizationID,
		"name":            role.Name,
		"permissions":     len(role.Permissions),
	})

	return role, nil
}

func (rm *RBACManager) GetRole(ctx context.Context, roleID uuid.UUID) (*Role, error) {
	query := `
		SELECT id, organization_id, name, display_name, description, permissions, is_system_role, created_at, updated_at, metadata
		FROM roles
		WHERE id = $1
	`

	var role Role
	var permissionsJSON, metadataJSON []byte

	err := rm.db.Conn.QueryRowContext(ctx, query, roleID).Scan(
		&role.ID, &role.OrganizationID, &role.Name, &role.DisplayName, &role.Description,
		&permissionsJSON, &role.IsSystemRole, &role.CreatedAt, &role.UpdatedAt, &metadataJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	// Unmarshal JSON fields
	json.Unmarshal(permissionsJSON, &role.Permissions)
	json.Unmarshal(metadataJSON, &role.Metadata)

	return &role, nil
}

func (rm *RBACManager) UpdateRole(ctx context.Context, roleID uuid.UUID, req *RoleUpdateRequest) error {
	// Check if role is system role
	var isSystemRole bool
	checkQuery := `SELECT is_system_role FROM roles WHERE id = $1`
	err := rm.db.Conn.QueryRowContext(ctx, checkQuery, roleID).Scan(&isSystemRole)
	if err != nil {
		return fmt.Errorf("failed to check role: %w", err)
	}

	if isSystemRole {
		return fmt.Errorf("cannot update system role")
	}

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

	if req.Permissions != nil {
		// Validate permissions
		for _, perm := range *req.Permissions {
			if !isValidPermission(perm) {
				return fmt.Errorf("invalid permission: %s", perm)
			}
		}

		permissionsJSON, _ := json.Marshal(*req.Permissions)
		setParts = append(setParts, fmt.Sprintf("permissions = $%d", argIndex))
		args = append(args, permissionsJSON)
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
	args = append(args, roleID)

	query := fmt.Sprintf("UPDATE roles SET %s WHERE id = $%d", 
		fmt.Sprintf("%s", setParts), argIndex)

	result, err := rm.db.Conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("role not found")
	}

	logging.Info("updated role", map[string]interface{}{
		"role_id":        roleID,
		"fields_updated": len(setParts) - 1, // exclude updated_at
	})

	return nil
}

func (rm *RBACManager) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	// Check if role is system role
	var isSystemRole bool
	checkQuery := `SELECT is_system_role FROM roles WHERE id = $1`
	err := rm.db.Conn.QueryRowContext(ctx, checkQuery, roleID).Scan(&isSystemRole)
	if err != nil {
		return fmt.Errorf("failed to check role: %w", err)
	}

	if isSystemRole {
		return fmt.Errorf("cannot delete system role")
	}

	query := `DELETE FROM roles WHERE id = $1`

	result, err := rm.db.Conn.ExecContext(ctx, query, roleID)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("role not found")
	}

	logging.Info("deleted role", map[string]interface{}{
		"role_id": roleID,
	})

	return nil
}

func (rm *RBACManager) AssignRole(ctx context.Context, orgID uuid.UUID, assignment *UserRoleAssignment) error {
	// Check if assignment already exists
	existingQuery := `SELECT id FROM user_roles WHERE user_id = $1 AND role_id = $2 AND organization_id = $3`
	var existingID uuid.UUID
	err := rm.db.Conn.QueryRowContext(ctx, existingQuery, assignment.UserID, assignment.RoleID, orgID).Scan(&existingID)
	if err != sql.ErrNoRows {
		return fmt.Errorf("role already assigned to user")
	}

	userRole := &UserRole{
		ID:             uuid.New(),
		UserID:         assignment.UserID,
		RoleID:         assignment.RoleID,
		OrganizationID: orgID,
		GrantedBy:      &assignment.GrantedBy,
		GrantedAt:      time.Now(),
		ExpiresAt:      assignment.ExpiresAt,
		CreatedAt:      time.Now(),
		Metadata:       make(map[string]interface{}),
	}

	query := `
		INSERT INTO user_roles (id, user_id, role_id, organization_id, granted_by, granted_at, expires_at, created_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	metadataJSON, _ := json.Marshal(userRole.Metadata)

	_, err = rm.db.Conn.ExecContext(ctx, query,
		userRole.ID, userRole.UserID, userRole.RoleID, userRole.OrganizationID,
		userRole.GrantedBy, userRole.GrantedAt, userRole.ExpiresAt,
		userRole.CreatedAt, metadataJSON)
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	logging.Info("assigned role to user", map[string]interface{}{
		"user_id":         assignment.UserID,
		"role_id":         assignment.RoleID,
		"organization_id": orgID,
		"granted_by":      assignment.GrantedBy,
		"expires_at":      assignment.ExpiresAt,
	})

	return nil
}

func (rm *RBACManager) RevokeRole(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, roleID uuid.UUID) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2 AND organization_id = $3`

	result, err := rm.db.Conn.ExecContext(ctx, query, userID, roleID, orgID)
	if err != nil {
		return fmt.Errorf("failed to revoke role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("role assignment not found")
	}

	logging.Info("revoked role from user", map[string]interface{}{
		"user_id":         userID,
		"role_id":         roleID,
		"organization_id": orgID,
	})

	return nil
}

func (rm *RBACManager) GetUserRoles(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) ([]*Role, error) {
	query := `
		SELECT r.id, r.organization_id, r.name, r.display_name, r.description, r.permissions, r.is_system_role, r.created_at, r.updated_at, r.metadata
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.organization_id = $1 AND ur.user_id = $2 AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
		ORDER BY r.name
	`

	rows, err := rm.db.Conn.QueryContext(ctx, query, orgID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	defer rows.Close()

	var roles []*Role
	for rows.Next() {
		var role Role
		var permissionsJSON, metadataJSON []byte

		err := rows.Scan(
			&role.ID, &role.OrganizationID, &role.Name, &role.DisplayName, &role.Description,
			&permissionsJSON, &role.IsSystemRole, &role.CreatedAt, &role.UpdatedAt, &metadataJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}

		// Unmarshal JSON fields
		json.Unmarshal(permissionsJSON, &role.Permissions)
		json.Unmarshal(metadataJSON, &role.Metadata)

		roles = append(roles, &role)
	}

	return roles, nil
}

func (rm *RBACManager) GetUserPermissions(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) ([]string, error) {
	roles, err := rm.GetUserRoles(ctx, orgID, userID)
	if err != nil {
		return nil, err
	}

	// Collect all unique permissions
	permissionSet := make(map[string]bool)
	for _, role := range roles {
		for _, perm := range role.Permissions {
			// Handle wildcard permissions
			if strings.HasSuffix(perm, ":*") {
				resource := strings.TrimSuffix(perm, ":*")
				// Add all actions for this resource
				for _, sysPerm := range SystemPermissions {
					if strings.HasPrefix(sysPerm, resource+":") {
						permissionSet[sysPerm] = true
					}
				}
			} else {
				permissionSet[perm] = true
			}
		}
	}

	// Convert set to slice
	permissions := make([]string, 0, len(permissionSet))
	for perm := range permissionSet {
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

func (rm *RBACManager) CheckPermission(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, permission string) (bool, error) {
	permissions, err := rm.GetUserPermissions(ctx, orgID, userID)
	if err != nil {
		return false, err
	}

	for _, perm := range permissions {
		if perm == permission {
			return true, nil
		}
	}

	return false, nil
}

// Helper functions

func isValidPermission(permission string) bool {
	// Check if it's a wildcard permission
	if strings.HasSuffix(permission, ":*") {
		resource := strings.TrimSuffix(permission, ":*")
		// Check if resource exists
		for _, sysPerm := range SystemPermissions {
			if strings.HasPrefix(sysPerm, resource+":") {
				return true
			}
		}
		return false
	}

	// Check if it's in the system permissions list
	for _, sysPerm := range SystemPermissions {
		if permission == sysPerm {
			return true
		}
	}

	return false
}

func capitalizeFirst(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}