package enterprise

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"holodeck1/auth"
	"holodeck1/database"
	"holodeck1/enterprise"
	"holodeck1/logging"
)

type Handlers struct {
	orgManager      *enterprise.OrganizationManager
	rbacManager     *enterprise.RBACManager
	analyticsManager *enterprise.AnalyticsManager
	securityManager *enterprise.SecurityManager
}

func NewHandlers(db *database.DB) *Handlers {
	return &Handlers{
		orgManager:      enterprise.NewOrganizationManager(db),
		rbacManager:     enterprise.NewRBACManager(db),
		analyticsManager: enterprise.NewAnalyticsManager(db),
		securityManager: enterprise.NewSecurityManager(db),
	}
}

// Organization handlers

func (h *Handlers) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uuid.UUID)

	var req enterprise.OrganizationCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	org, err := h.orgManager.CreateOrganization(r.Context(), &req)
	if err != nil {
		logging.Error("failed to create organization", err, map[string]interface{}{
			"user_id": userID,
		})
		http.Error(w, "Failed to create organization", http.StatusInternalServerError)
		return
	}

	// Assign owner role to creating user
	err = h.rbacManager.InitializeDefaultRoles(r.Context(), org.ID)
	if err != nil {
		logging.Error("failed to initialize roles", err, map[string]interface{}{
			"organization_id": org.ID,
		})
	}

	// Get owner role
	roles, err := h.rbacManager.GetUserRoles(r.Context(), org.ID, userID)
	if err == nil && len(roles) == 0 {
		// Assign owner role
		var ownerRoleID uuid.UUID
		// Find owner role ID (would normally query for it)
		h.rbacManager.AssignRole(r.Context(), org.ID, &enterprise.UserRoleAssignment{
			UserID:    userID,
			RoleID:    ownerRoleID,
			GrantedBy: userID,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(org)
}

func (h *Handlers) GetOrganization(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	// Check permission
	hasPermission, err := h.rbacManager.CheckPermission(r.Context(), orgID, userID, "organizations:read")
	if err != nil || !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	org, err := h.orgManager.GetOrganization(r.Context(), orgID)
	if err != nil {
		http.Error(w, "Organization not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(org)
}

func (h *Handlers) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	// Check permission
	hasPermission, err := h.rbacManager.CheckPermission(r.Context(), orgID, userID, "organizations:write")
	if err != nil || !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req enterprise.OrganizationUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.orgManager.UpdateOrganization(r.Context(), orgID, &req)
	if err != nil {
		http.Error(w, "Failed to update organization", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	// Check permission
	hasPermission, err := h.rbacManager.CheckPermission(r.Context(), orgID, userID, "organizations:delete")
	if err != nil || !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err = h.orgManager.DeleteOrganization(r.Context(), orgID)
	if err != nil {
		http.Error(w, "Failed to delete organization", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) ListOrganizations(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uuid.UUID)

	filter := &enterprise.OrganizationFilter{
		Status: r.URL.Query().Get("status"),
		SubscriptionTier: r.URL.Query().Get("tier"),
		Domain: r.URL.Query().Get("domain"),
		Limit: 100,
		Offset: 0,
	}

	orgs, err := h.orgManager.ListOrganizations(r.Context(), filter)
	if err != nil {
		http.Error(w, "Failed to list organizations", http.StatusInternalServerError)
		return
	}

	// Filter by user access (simplified - would normally use joins)
	var userOrgs []*enterprise.Organization
	for _, org := range orgs {
		hasPermission, _ := h.rbacManager.CheckPermission(r.Context(), org.ID, userID, "organizations:read")
		if hasPermission {
			userOrgs = append(userOrgs, org)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userOrgs)
}

// RBAC handlers

func (h *Handlers) CreateRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgID, err := uuid.Parse(vars["orgId"])
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	// Check permission
	hasPermission, err := h.rbacManager.CheckPermission(r.Context(), orgID, userID, "users:manage_roles")
	if err != nil || !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req enterprise.RoleCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.OrganizationID = orgID

	role, err := h.rbacManager.CreateRole(r.Context(), &req)
	if err != nil {
		http.Error(w, "Failed to create role", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

func (h *Handlers) AssignRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgID, err := uuid.Parse(vars["orgId"])
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	// Check permission
	hasPermission, err := h.rbacManager.CheckPermission(r.Context(), orgID, userID, "users:manage_roles")
	if err != nil || !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req enterprise.UserRoleAssignment
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.GrantedBy = userID

	err = h.rbacManager.AssignRole(r.Context(), orgID, &req)
	if err != nil {
		http.Error(w, "Failed to assign role", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) GetUserPermissions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgID, err := uuid.Parse(vars["orgId"])
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	targetUserID, err := uuid.Parse(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	// Users can check their own permissions or admins can check others
	if targetUserID != userID {
		hasPermission, err := h.rbacManager.CheckPermission(r.Context(), orgID, userID, "users:read")
		if err != nil || !hasPermission {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	permissions, err := h.rbacManager.GetUserPermissions(r.Context(), orgID, targetUserID)
	if err != nil {
		http.Error(w, "Failed to get permissions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": targetUserID,
		"organization_id": orgID,
		"permissions": permissions,
	})
}

// Analytics handlers

func (h *Handlers) TrackEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgID, err := uuid.Parse(vars["orgId"])
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	var req enterprise.EventTrackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set user ID if not provided
	if req.UserID == nil {
		req.UserID = &userID
	}

	// Get client info
	req.UserAgent = r.UserAgent()
	req.IPAddress = r.RemoteAddr

	err = h.analyticsManager.TrackEvent(r.Context(), orgID, &req)
	if err != nil {
		http.Error(w, "Failed to track event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) GetAnalyticsReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgID, err := uuid.Parse(vars["orgId"])
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	// Check permission
	hasPermission, err := h.rbacManager.CheckPermission(r.Context(), orgID, userID, "analytics:read")
	if err != nil || !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse time range
	startTimeStr := r.URL.Query().Get("start_time")
	endTimeStr := r.URL.Query().Get("end_time")

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		startTime = time.Now().AddDate(0, 0, -7) // Default to last 7 days
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		endTime = time.Now()
	}

	report, err := h.analyticsManager.GenerateReport(r.Context(), orgID, startTime, endTime)
	if err != nil {
		http.Error(w, "Failed to generate report", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// Security handlers

func (h *Handlers) LogSecurityEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgID, err := uuid.Parse(vars["orgId"])
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	var req enterprise.SecurityEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set user ID if not provided
	if req.UserID == nil {
		req.UserID = &userID
	}

	// Get client info
	req.UserAgent = r.UserAgent()
	req.IPAddress = r.RemoteAddr

	err = h.securityManager.LogSecurityEvent(r.Context(), orgID, &req)
	if err != nil {
		http.Error(w, "Failed to log security event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) GetSecurityAuditLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgID, err := uuid.Parse(vars["orgId"])
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	// Check permission
	hasPermission, err := h.rbacManager.CheckPermission(r.Context(), orgID, userID, "security:read_audit_logs")
	if err != nil || !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse time range
	startTimeStr := r.URL.Query().Get("start_time")
	endTimeStr := r.URL.Query().Get("end_time")

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		startTime = time.Now().AddDate(0, 0, -7) // Default to last 7 days
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		endTime = time.Now()
	}

	filters := map[string]string{
		"event_category": r.URL.Query().Get("category"),
		"risk_level": r.URL.Query().Get("risk_level"),
		"result": r.URL.Query().Get("result"),
	}

	events, err := h.securityManager.GetSecurityEvents(r.Context(), orgID, startTime, endTime, filters)
	if err != nil {
		http.Error(w, "Failed to get security events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (h *Handlers) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgID, err := uuid.Parse(vars["orgId"])
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	// Check permission
	hasPermission, err := h.rbacManager.CheckPermission(r.Context(), orgID, userID, "security:manage_api_keys")
	if err != nil || !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req enterprise.APIKeyCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	apiKey, err := h.securityManager.CreateAPIKey(r.Context(), orgID, userID, &req)
	if err != nil {
		http.Error(w, "Failed to create API key", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiKey)
}

func (h *Handlers) CreateComplianceRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orgID, err := uuid.Parse(vars["orgId"])
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	// Check permission
	hasPermission, err := h.rbacManager.CheckPermission(r.Context(), orgID, userID, "security:manage_compliance")
	if err != nil || !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req enterprise.ComplianceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.CreatedBy = userID

	record, err := h.securityManager.CreateComplianceRecord(r.Context(), orgID, &req)
	if err != nil {
		http.Error(w, "Failed to create compliance record", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}

// Middleware

func (h *Handlers) RequireOrganizationAccess(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orgID, err := uuid.Parse(vars["orgId"])
		if err != nil {
			http.Error(w, "Invalid organization ID", http.StatusBadRequest)
			return
		}

		userID := r.Context().Value("userID").(uuid.UUID)

		// Check if user has any permission in the organization
		permissions, err := h.rbacManager.GetUserPermissions(r.Context(), orgID, userID)
		if err != nil || len(permissions) == 0 {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}