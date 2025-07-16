package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"holodeck1/database"
	"holodeck1/logging"
)

type Registry struct {
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

type ServiceHealth struct {
	ID           uuid.UUID              `json:"id"`
	ServiceID    uuid.UUID              `json:"service_id"`
	Status       string                 `json:"status"`
	ResponseTime int                    `json:"response_time"`
	ErrorMessage string                 `json:"error_message"`
	CheckedAt    time.Time              `json:"checked_at"`
	Metadata     map[string]interface{} `json:"metadata"`
}

type RegisterServiceRequest struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Type         string                 `json:"type"`
	Endpoint     string                 `json:"endpoint"`
	Capabilities []string               `json:"capabilities"`
	UIMapping    map[string]interface{} `json:"ui_mapping"`
	Permissions  map[string]interface{} `json:"permissions"`
	Metadata     map[string]interface{} `json:"metadata"`
}

type UpdateServiceRequest struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Endpoint     string                 `json:"endpoint"`
	Status       string                 `json:"status"`
	Capabilities []string               `json:"capabilities"`
	UIMapping    map[string]interface{} `json:"ui_mapping"`
	Permissions  map[string]interface{} `json:"permissions"`
	Metadata     map[string]interface{} `json:"metadata"`
}

type ServiceFilter struct {
	Type   string `json:"type"`
	Status string `json:"status"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

func NewRegistry(db *database.DB) *Registry {
	return &Registry{db: db}
}

func (r *Registry) RegisterService(ctx context.Context, req *RegisterServiceRequest) (*Service, error) {
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

	if service.Capabilities == nil {
		service.Capabilities = []string{}
	}
	if service.UIMapping == nil {
		service.UIMapping = make(map[string]interface{})
	}
	if service.Permissions == nil {
		service.Permissions = make(map[string]interface{})
	}
	if service.Metadata == nil {
		service.Metadata = make(map[string]interface{})
	}

	capabilitiesJSON, _ := json.Marshal(service.Capabilities)
	uiMappingJSON, _ := json.Marshal(service.UIMapping)
	permissionsJSON, _ := json.Marshal(service.Permissions)
	metadataJSON, _ := json.Marshal(service.Metadata)

	query := `
		INSERT INTO services (id, name, description, type, endpoint, status, capabilities, ui_mapping, permissions, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		service.ID,
		service.Name,
		service.Description,
		service.Type,
		service.Endpoint,
		service.Status,
		capabilitiesJSON,
		uiMappingJSON,
		permissionsJSON,
		metadataJSON,
	).Scan(&service.CreatedAt, &service.UpdatedAt)

	if err != nil {
		logging.Error("failed to register service", map[string]interface{}{
			"service_name": req.Name,
			"error":        err.Error(),
		})
		return nil, fmt.Errorf("failed to register service: %w", err)
	}

	logging.Info("service registered", map[string]interface{}{
		"service_id":   service.ID,
		"service_name": service.Name,
		"service_type": service.Type,
		"endpoint":     service.Endpoint,
	})

	return service, nil
}

func (r *Registry) GetService(ctx context.Context, serviceID uuid.UUID) (*Service, error) {
	service := &Service{}
	var capabilitiesJSON, uiMappingJSON, permissionsJSON, metadataJSON []byte

	query := `
		SELECT id, name, description, type, endpoint, status, capabilities, ui_mapping, permissions, created_at, updated_at, metadata
		FROM services WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, serviceID).Scan(
		&service.ID,
		&service.Name,
		&service.Description,
		&service.Type,
		&service.Endpoint,
		&service.Status,
		&capabilitiesJSON,
		&uiMappingJSON,
		&permissionsJSON,
		&service.CreatedAt,
		&service.UpdatedAt,
		&metadataJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("service not found")
		}
		logging.Error("failed to get service", map[string]interface{}{
			"service_id": serviceID,
			"error":      err.Error(),
		})
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	json.Unmarshal(capabilitiesJSON, &service.Capabilities)
	json.Unmarshal(uiMappingJSON, &service.UIMapping)
	json.Unmarshal(permissionsJSON, &service.Permissions)
	json.Unmarshal(metadataJSON, &service.Metadata)

	return service, nil
}

func (r *Registry) GetServiceByName(ctx context.Context, name string) (*Service, error) {
	service := &Service{}
	var capabilitiesJSON, uiMappingJSON, permissionsJSON, metadataJSON []byte

	query := `
		SELECT id, name, description, type, endpoint, status, capabilities, ui_mapping, permissions, created_at, updated_at, metadata
		FROM services WHERE name = $1
	`

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&service.ID,
		&service.Name,
		&service.Description,
		&service.Type,
		&service.Endpoint,
		&service.Status,
		&capabilitiesJSON,
		&uiMappingJSON,
		&permissionsJSON,
		&service.CreatedAt,
		&service.UpdatedAt,
		&metadataJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("service not found")
		}
		logging.Error("failed to get service by name", map[string]interface{}{
			"service_name": name,
			"error":        err.Error(),
		})
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	json.Unmarshal(capabilitiesJSON, &service.Capabilities)
	json.Unmarshal(uiMappingJSON, &service.UIMapping)
	json.Unmarshal(permissionsJSON, &service.Permissions)
	json.Unmarshal(metadataJSON, &service.Metadata)

	return service, nil
}

func (r *Registry) ListServices(ctx context.Context, filter *ServiceFilter) ([]*Service, error) {
	query := `
		SELECT id, name, description, type, endpoint, status, capabilities, ui_mapping, permissions, created_at, updated_at, metadata
		FROM services
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if filter.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, filter.Type)
		argIndex++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
		argIndex++
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		logging.Error("failed to list services", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to list services: %w", err)
	}
	defer rows.Close()

	services := []*Service{}
	for rows.Next() {
		service := &Service{}
		var capabilitiesJSON, uiMappingJSON, permissionsJSON, metadataJSON []byte

		err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.Description,
			&service.Type,
			&service.Endpoint,
			&service.Status,
			&capabilitiesJSON,
			&uiMappingJSON,
			&permissionsJSON,
			&service.CreatedAt,
			&service.UpdatedAt,
			&metadataJSON,
		)

		if err != nil {
			logging.Error("failed to scan service", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}

		json.Unmarshal(capabilitiesJSON, &service.Capabilities)
		json.Unmarshal(uiMappingJSON, &service.UIMapping)
		json.Unmarshal(permissionsJSON, &service.Permissions)
		json.Unmarshal(metadataJSON, &service.Metadata)

		services = append(services, service)
	}

	return services, nil
}

func (r *Registry) UpdateService(ctx context.Context, serviceID uuid.UUID, req *UpdateServiceRequest) (*Service, error) {
	service, err := r.GetService(ctx, serviceID)
	if err != nil {
		return nil, err
	}

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
	if req.Permissions != nil {
		service.Permissions = req.Permissions
	}
	if req.Metadata != nil {
		service.Metadata = req.Metadata
	}

	capabilitiesJSON, _ := json.Marshal(service.Capabilities)
	uiMappingJSON, _ := json.Marshal(service.UIMapping)
	permissionsJSON, _ := json.Marshal(service.Permissions)
	metadataJSON, _ := json.Marshal(service.Metadata)

	query := `
		UPDATE services 
		SET name = $2, description = $3, endpoint = $4, status = $5, capabilities = $6, ui_mapping = $7, permissions = $8, metadata = $9
		WHERE id = $1
		RETURNING updated_at
	`

	err = r.db.QueryRowContext(ctx, query,
		serviceID,
		service.Name,
		service.Description,
		service.Endpoint,
		service.Status,
		capabilitiesJSON,
		uiMappingJSON,
		permissionsJSON,
		metadataJSON,
	).Scan(&service.UpdatedAt)

	if err != nil {
		logging.Error("failed to update service", map[string]interface{}{
			"service_id": serviceID,
			"error":      err.Error(),
		})
		return nil, fmt.Errorf("failed to update service: %w", err)
	}

	logging.Info("service updated", map[string]interface{}{
		"service_id":   serviceID,
		"service_name": service.Name,
	})

	return service, nil
}

func (r *Registry) DeleteService(ctx context.Context, serviceID uuid.UUID) error {
	query := `DELETE FROM services WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, serviceID)
	if err != nil {
		logging.Error("failed to delete service", map[string]interface{}{
			"service_id": serviceID,
			"error":      err.Error(),
		})
		return fmt.Errorf("failed to delete service: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("service not found")
	}

	logging.Info("service deleted", map[string]interface{}{
		"service_id": serviceID,
	})

	return nil
}

func (r *Registry) CheckServiceHealth(ctx context.Context, serviceID uuid.UUID) (*ServiceHealth, error) {
	service, err := r.GetService(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()
	
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Make health check request
	resp, err := client.Get(service.Endpoint + "/health")
	responseTime := int(time.Since(startTime).Milliseconds())

	health := &ServiceHealth{
		ID:           uuid.New(),
		ServiceID:    serviceID,
		ResponseTime: responseTime,
		CheckedAt:    time.Now(),
		Metadata:     make(map[string]interface{}),
	}

	if err != nil {
		health.Status = "error"
		health.ErrorMessage = err.Error()
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			health.Status = "healthy"
		} else {
			health.Status = "unhealthy"
			health.ErrorMessage = fmt.Sprintf("HTTP %d", resp.StatusCode)
		}
	}

	// Store health check result
	metadataJSON, _ := json.Marshal(health.Metadata)

	query := `
		INSERT INTO service_health (id, service_id, status, response_time, error_message, checked_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = r.db.ExecContext(ctx, query,
		health.ID,
		health.ServiceID,
		health.Status,
		health.ResponseTime,
		health.ErrorMessage,
		health.CheckedAt,
		metadataJSON,
	)

	if err != nil {
		logging.Error("failed to store health check result", map[string]interface{}{
			"service_id": serviceID,
			"error":      err.Error(),
		})
		// Don't return error here, health check result is still valid
	}

	logging.Debug("service health check completed", map[string]interface{}{
		"service_id":    serviceID,
		"status":        health.Status,
		"response_time": health.ResponseTime,
	})

	return health, nil
}

func (r *Registry) GetServiceHealth(ctx context.Context, serviceID uuid.UUID) ([]*ServiceHealth, error) {
	query := `
		SELECT id, service_id, status, response_time, error_message, checked_at, metadata
		FROM service_health 
		WHERE service_id = $1
		ORDER BY checked_at DESC
		LIMIT 100
	`

	rows, err := r.db.QueryContext(ctx, query, serviceID)
	if err != nil {
		logging.Error("failed to get service health", map[string]interface{}{
			"service_id": serviceID,
			"error":      err.Error(),
		})
		return nil, fmt.Errorf("failed to get service health: %w", err)
	}
	defer rows.Close()

	healthChecks := []*ServiceHealth{}
	for rows.Next() {
		health := &ServiceHealth{}
		var metadataJSON []byte

		err := rows.Scan(
			&health.ID,
			&health.ServiceID,
			&health.Status,
			&health.ResponseTime,
			&health.ErrorMessage,
			&health.CheckedAt,
			&metadataJSON,
		)

		if err != nil {
			logging.Error("failed to scan service health", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}

		json.Unmarshal(metadataJSON, &health.Metadata)
		healthChecks = append(healthChecks, health)
	}

	return healthChecks, nil
}

func (r *Registry) StartHealthMonitoring(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.performHealthChecks(ctx)
		}
	}
}

func (r *Registry) performHealthChecks(ctx context.Context) {
	services, err := r.ListServices(ctx, &ServiceFilter{Status: "active"})
	if err != nil {
		logging.Error("failed to list services for health check", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	for _, service := range services {
		go func(serviceID uuid.UUID) {
			_, err := r.CheckServiceHealth(ctx, serviceID)
			if err != nil {
				logging.Error("health check failed", map[string]interface{}{
					"service_id": serviceID,
					"error":      err.Error(),
				})
			}
		}(service.ID)
	}
}