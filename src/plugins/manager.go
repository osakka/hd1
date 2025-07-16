package plugins

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"holodeck1/database"
	"holodeck1/logging"
)

type Manager struct {
	db          *database.DB
	plugins     map[uuid.UUID]*Plugin
	registry    map[string]*PluginRegistry
	hooks       map[string][]Hook
	pluginsPath string
	mu          sync.RWMutex
}

type Plugin struct {
	ID           uuid.UUID              `json:"id"`
	Name         string                 `json:"name"`
	Version      string                 `json:"version"`
	Description  string                 `json:"description"`
	Author       string                 `json:"author"`
	Type         string                 `json:"type"`
	Category     string                 `json:"category"`
	Entry        string                 `json:"entry"`
	Dependencies []PluginDependency     `json:"dependencies"`
	Permissions  []string               `json:"permissions"`
	Configuration map[string]interface{} `json:"configuration"`
	Status       string                 `json:"status"`
	LoadedAt     *time.Time             `json:"loaded_at"`
	Metadata     map[string]interface{} `json:"metadata"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type PluginDependency struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Required     bool   `json:"required"`
	MinVersion   string `json:"min_version"`
	MaxVersion   string `json:"max_version"`
}

type PluginRegistry struct {
	Name        string                 `json:"name"`
	URL         string                 `json:"url"`
	Type        string                 `json:"type"`
	Credentials map[string]interface{} `json:"credentials"`
	Status      string                 `json:"status"`
	LastSync    time.Time              `json:"last_sync"`
}

type Hook struct {
	ID         uuid.UUID              `json:"id"`
	PluginID   uuid.UUID              `json:"plugin_id"`
	Event      string                 `json:"event"`
	Handler    string                 `json:"handler"`
	Priority   int                    `json:"priority"`
	Conditions map[string]interface{} `json:"conditions"`
	CreatedAt  time.Time              `json:"created_at"`
}

type PluginManifest struct {
	Name         string                 `json:"name"`
	Version      string                 `json:"version"`
	Description  string                 `json:"description"`
	Author       string                 `json:"author"`
	License      string                 `json:"license"`
	Type         string                 `json:"type"`
	Category     string                 `json:"category"`
	Entry        string                 `json:"entry"`
	Dependencies []PluginDependency     `json:"dependencies"`
	Permissions  []string               `json:"permissions"`
	Hooks        []HookDefinition       `json:"hooks"`
	Configuration []ConfigField         `json:"configuration"`
	Assets       []string               `json:"assets"`
	Metadata     map[string]interface{} `json:"metadata"`
}

type HookDefinition struct {
	Event      string                 `json:"event"`
	Handler    string                 `json:"handler"`
	Priority   int                    `json:"priority"`
	Conditions map[string]interface{} `json:"conditions"`
}

type ConfigField struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default"`
	Description string      `json:"description"`
	Options     []string    `json:"options,omitempty"`
}

type PluginEvent struct {
	ID        uuid.UUID              `json:"id"`
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`
	Data      map[string]interface{} `json:"data"`
	Context   map[string]interface{} `json:"context"`
	Timestamp time.Time              `json:"timestamp"`
}

type InstallRequest struct {
	Source      string                 `json:"source"`
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Registry    string                 `json:"registry"`
	Config      map[string]interface{} `json:"config"`
	AutoEnable  bool                   `json:"auto_enable"`
}

type PluginFilter struct {
	Type     string `json:"type"`
	Category string `json:"category"`
	Status   string `json:"status"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
}

func NewManager(db *database.DB, pluginsPath string) *Manager {
	// Create plugins directory if it doesn't exist
	os.MkdirAll(pluginsPath, 0755)

	return &Manager{
		db:          db,
		plugins:     make(map[uuid.UUID]*Plugin),
		registry:    make(map[string]*PluginRegistry),
		hooks:       make(map[string][]Hook),
		pluginsPath: pluginsPath,
	}
}

func (m *Manager) InstallPlugin(ctx context.Context, req *InstallRequest) (*Plugin, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if plugin already exists
	existingPlugin := m.findPluginByName(req.Name)
	if existingPlugin != nil {
		return nil, fmt.Errorf("plugin %s already installed", req.Name)
	}

	// Download/copy plugin
	pluginPath, err := m.downloadPlugin(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to download plugin: %w", err)
	}

	// Load manifest
	manifest, err := m.loadManifest(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load manifest: %w", err)
	}

	// Validate dependencies
	err = m.validateDependencies(manifest.Dependencies)
	if err != nil {
		return nil, fmt.Errorf("dependency validation failed: %w", err)
	}

	// Create plugin record
	pluginID := uuid.New()
	plugin := &Plugin{
		ID:            pluginID,
		Name:          manifest.Name,
		Version:       manifest.Version,
		Description:   manifest.Description,
		Author:        manifest.Author,
		Type:          manifest.Type,
		Category:      manifest.Category,
		Entry:         manifest.Entry,
		Dependencies:  manifest.Dependencies,
		Permissions:   manifest.Permissions,
		Configuration: req.Config,
		Status:        "installed",
		Metadata:      manifest.Metadata,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Store in database
	query := `
		INSERT INTO plugins (id, name, version, description, author, type, category, entry, dependencies, permissions, configuration, status, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	
	dependenciesJSON, _ := json.Marshal(plugin.Dependencies)
	permissionsJSON, _ := json.Marshal(plugin.Permissions)
	configJSON, _ := json.Marshal(plugin.Configuration)
	metadataJSON, _ := json.Marshal(plugin.Metadata)
	
	_, err = m.db.Conn.ExecContext(ctx, query,
		plugin.ID, plugin.Name, plugin.Version, plugin.Description, plugin.Author,
		plugin.Type, plugin.Category, plugin.Entry, dependenciesJSON,
		permissionsJSON, configJSON, plugin.Status, metadataJSON,
		plugin.CreatedAt, plugin.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to store plugin: %w", err)
	}

	m.plugins[pluginID] = plugin

	// Register hooks
	err = m.registerHooks(ctx, plugin, manifest.Hooks)
	if err != nil {
		logging.Error("failed to register hooks", map[string]interface{}{
			"plugin_id": pluginID,
			"error":     err.Error(),
		})
	}

	// Auto-enable if requested
	if req.AutoEnable {
		err = m.enablePlugin(ctx, pluginID)
		if err != nil {
			logging.Error("failed to auto-enable plugin", map[string]interface{}{
				"plugin_id": pluginID,
				"error":     err.Error(),
			})
		}
	}

	logging.Info("installed plugin", map[string]interface{}{
		"plugin_id": pluginID,
		"name":      plugin.Name,
		"version":   plugin.Version,
		"type":      plugin.Type,
		"category":  plugin.Category,
	})

	return plugin, nil
}

func (m *Manager) EnablePlugin(ctx context.Context, pluginID uuid.UUID) error {
	return m.enablePlugin(ctx, pluginID)
}

func (m *Manager) enablePlugin(ctx context.Context, pluginID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	plugin, exists := m.plugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin not found")
	}

	if plugin.Status == "enabled" {
		return fmt.Errorf("plugin already enabled")
	}

	// Load plugin code (implementation would vary by plugin type)
	err := m.loadPlugin(plugin)
	if err != nil {
		return fmt.Errorf("failed to load plugin: %w", err)
	}

	// Update status
	plugin.Status = "enabled"
	loadedAt := time.Now()
	plugin.LoadedAt = &loadedAt
	plugin.UpdatedAt = time.Now()

	// Update database
	query := `UPDATE plugins SET status = $1, loaded_at = $2, updated_at = $3 WHERE id = $4`
	_, err = m.db.Conn.ExecContext(ctx, query, plugin.Status, plugin.LoadedAt, plugin.UpdatedAt, pluginID)
	if err != nil {
		return fmt.Errorf("failed to update plugin status: %w", err)
	}

	logging.Info("enabled plugin", map[string]interface{}{
		"plugin_id": pluginID,
		"name":      plugin.Name,
	})

	return nil
}

func (m *Manager) DisablePlugin(ctx context.Context, pluginID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	plugin, exists := m.plugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin not found")
	}

	if plugin.Status != "enabled" {
		return fmt.Errorf("plugin not enabled")
	}

	// Unload plugin code
	err := m.unloadPlugin(plugin)
	if err != nil {
		return fmt.Errorf("failed to unload plugin: %w", err)
	}

	// Update status
	plugin.Status = "disabled"
	plugin.LoadedAt = nil
	plugin.UpdatedAt = time.Now()

	// Update database
	query := `UPDATE plugins SET status = $1, loaded_at = NULL, updated_at = $2 WHERE id = $3`
	_, err = m.db.Conn.ExecContext(ctx, query, plugin.Status, plugin.UpdatedAt, pluginID)
	if err != nil {
		return fmt.Errorf("failed to update plugin status: %w", err)
	}

	logging.Info("disabled plugin", map[string]interface{}{
		"plugin_id": pluginID,
		"name":      plugin.Name,
	})

	return nil
}

func (m *Manager) UninstallPlugin(ctx context.Context, pluginID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	plugin, exists := m.plugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin not found")
	}

	// Disable plugin first
	if plugin.Status == "enabled" {
		err := m.unloadPlugin(plugin)
		if err != nil {
			logging.Error("failed to unload plugin during uninstall", map[string]interface{}{
				"plugin_id": pluginID,
				"error":     err.Error(),
			})
		}
	}

	// Remove hooks
	m.unregisterHooks(pluginID)

	// Delete from database
	query := `DELETE FROM plugins WHERE id = $1`
	_, err := m.db.Conn.ExecContext(ctx, query, pluginID)
	if err != nil {
		return fmt.Errorf("failed to delete plugin: %w", err)
	}

	// Remove plugin files
	pluginDir := filepath.Join(m.pluginsPath, plugin.Name)
	err = os.RemoveAll(pluginDir)
	if err != nil {
		logging.Error("failed to remove plugin files", map[string]interface{}{
			"plugin_id": pluginID,
			"path":      pluginDir,
			"error":     err.Error(),
		})
	}

	// Remove from memory
	delete(m.plugins, pluginID)

	logging.Info("uninstalled plugin", map[string]interface{}{
		"plugin_id": pluginID,
		"name":      plugin.Name,
	})

	return nil
}

func (m *Manager) GetPlugin(ctx context.Context, pluginID uuid.UUID) (*Plugin, error) {
	m.mu.RLock()
	plugin, exists := m.plugins[pluginID]
	m.mu.RUnlock()

	if exists {
		return plugin, nil
	}

	// Load from database
	query := `
		SELECT id, name, version, description, author, type, category, entry, dependencies, permissions, configuration, status, loaded_at, metadata, created_at, updated_at
		FROM plugins
		WHERE id = $1
	`
	
	var dbPlugin Plugin
	var dependenciesJSON, permissionsJSON, configJSON, metadataJSON []byte
	var loadedAt sql.NullTime
	
	err := m.db.Conn.QueryRowContext(ctx, query, pluginID).Scan(
		&dbPlugin.ID, &dbPlugin.Name, &dbPlugin.Version, &dbPlugin.Description,
		&dbPlugin.Author, &dbPlugin.Type, &dbPlugin.Category, &dbPlugin.Entry,
		&dependenciesJSON, &permissionsJSON, &configJSON, &dbPlugin.Status,
		&loadedAt, &metadataJSON, &dbPlugin.CreatedAt, &dbPlugin.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("plugin not found")
		}
		return nil, fmt.Errorf("failed to get plugin: %w", err)
	}

	// Unmarshal JSON fields
	json.Unmarshal(dependenciesJSON, &dbPlugin.Dependencies)
	json.Unmarshal(permissionsJSON, &dbPlugin.Permissions)
	json.Unmarshal(configJSON, &dbPlugin.Configuration)
	json.Unmarshal(metadataJSON, &dbPlugin.Metadata)

	if loadedAt.Valid {
		dbPlugin.LoadedAt = &loadedAt.Time
	}

	m.mu.Lock()
	m.plugins[pluginID] = &dbPlugin
	m.mu.Unlock()

	return &dbPlugin, nil
}

func (m *Manager) ListPlugins(ctx context.Context, filter *PluginFilter) ([]*Plugin, error) {
	query := `
		SELECT id, name, version, description, author, type, category, entry, dependencies, permissions, configuration, status, loaded_at, metadata, created_at, updated_at
		FROM plugins
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if filter.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, filter.Type)
		argIndex++
	}

	if filter.Category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIndex)
		args = append(args, filter.Category)
		argIndex++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
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

	rows, err := m.db.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query plugins: %w", err)
	}
	defer rows.Close()

	var plugins []*Plugin
	for rows.Next() {
		var plugin Plugin
		var dependenciesJSON, permissionsJSON, configJSON, metadataJSON []byte
		var loadedAt sql.NullTime
		
		err := rows.Scan(
			&plugin.ID, &plugin.Name, &plugin.Version, &plugin.Description,
			&plugin.Author, &plugin.Type, &plugin.Category, &plugin.Entry,
			&dependenciesJSON, &permissionsJSON, &configJSON, &plugin.Status,
			&loadedAt, &metadataJSON, &plugin.CreatedAt, &plugin.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan plugin: %w", err)
		}

		// Unmarshal JSON fields
		json.Unmarshal(dependenciesJSON, &plugin.Dependencies)
		json.Unmarshal(permissionsJSON, &plugin.Permissions)
		json.Unmarshal(configJSON, &plugin.Configuration)
		json.Unmarshal(metadataJSON, &plugin.Metadata)

		if loadedAt.Valid {
			plugin.LoadedAt = &loadedAt.Time
		}

		plugins = append(plugins, &plugin)
	}

	return plugins, nil
}

func (m *Manager) EmitEvent(ctx context.Context, event *PluginEvent) error {
	m.mu.RLock()
	hooks, exists := m.hooks[event.Type]
	m.mu.RUnlock()

	if !exists {
		return nil
	}

	// Execute hooks in priority order
	for _, hook := range hooks {
		go m.executeHook(ctx, hook, event)
	}

	return nil
}

func (m *Manager) executeHook(ctx context.Context, hook Hook, event *PluginEvent) {
	// Check conditions
	if !m.checkHookConditions(hook.Conditions, event) {
		return
	}

	// Get plugin
	plugin, exists := m.plugins[hook.PluginID]
	if !exists || plugin.Status != "enabled" {
		return
	}

	// Execute hook handler (implementation would vary by plugin type)
	err := m.callHookHandler(plugin, hook.Handler, event)
	if err != nil {
		logging.Error("failed to execute hook", map[string]interface{}{
			"hook_id":    hook.ID,
			"plugin_id":  hook.PluginID,
			"event_type": event.Type,
			"handler":    hook.Handler,
			"error":      err.Error(),
		})
	}
}

func (m *Manager) checkHookConditions(conditions map[string]interface{}, event *PluginEvent) bool {
	// Simple condition checking implementation
	for key, expectedValue := range conditions {
		if actualValue, exists := event.Data[key]; exists {
			if actualValue != expectedValue {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func (m *Manager) findPluginByName(name string) *Plugin {
	for _, plugin := range m.plugins {
		if plugin.Name == name {
			return plugin
		}
	}
	return nil
}

func (m *Manager) downloadPlugin(ctx context.Context, req *InstallRequest) (string, error) {
	// Placeholder implementation for plugin download
	// This would handle different sources (file, URL, registry, etc.)
	pluginDir := filepath.Join(m.pluginsPath, req.Name)
	err := os.MkdirAll(pluginDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create plugin directory: %w", err)
	}
	return pluginDir, nil
}

func (m *Manager) loadManifest(pluginPath string) (*PluginManifest, error) {
	// Placeholder implementation for manifest loading
	// This would read and parse the plugin's manifest file
	manifest := &PluginManifest{
		Name:        "example-plugin",
		Version:     "1.0.0",
		Description: "Example plugin",
		Author:      "HD1",
		Type:        "service",
		Category:    "utility",
		Entry:       "main.js",
		Dependencies: []PluginDependency{},
		Permissions:  []string{"read", "write"},
		Hooks:        []HookDefinition{},
		Configuration: []ConfigField{},
		Assets:       []string{},
		Metadata:     make(map[string]interface{}),
	}
	return manifest, nil
}

func (m *Manager) validateDependencies(dependencies []PluginDependency) error {
	// Placeholder implementation for dependency validation
	return nil
}

func (m *Manager) loadPlugin(plugin *Plugin) error {
	// Placeholder implementation for plugin loading
	// This would vary by plugin type (JavaScript, WebAssembly, etc.)
	return nil
}

func (m *Manager) unloadPlugin(plugin *Plugin) error {
	// Placeholder implementation for plugin unloading
	return nil
}

func (m *Manager) registerHooks(ctx context.Context, plugin *Plugin, hookDefs []HookDefinition) error {
	for _, hookDef := range hookDefs {
		hookID := uuid.New()
		hook := Hook{
			ID:         hookID,
			PluginID:   plugin.ID,
			Event:      hookDef.Event,
			Handler:    hookDef.Handler,
			Priority:   hookDef.Priority,
			Conditions: hookDef.Conditions,
			CreatedAt:  time.Now(),
		}

		// Store in database
		query := `
			INSERT INTO plugin_hooks (id, plugin_id, event, handler, priority, conditions, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		conditionsJSON, _ := json.Marshal(hook.Conditions)
		_, err := m.db.Conn.ExecContext(ctx, query,
			hook.ID, hook.PluginID, hook.Event, hook.Handler,
			hook.Priority, conditionsJSON, hook.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to store hook: %w", err)
		}

		// Add to memory
		m.hooks[hook.Event] = append(m.hooks[hook.Event], hook)
	}

	return nil
}

func (m *Manager) unregisterHooks(pluginID uuid.UUID) {
	// Remove from memory
	for event, hooks := range m.hooks {
		var filteredHooks []Hook
		for _, hook := range hooks {
			if hook.PluginID != pluginID {
				filteredHooks = append(filteredHooks, hook)
			}
		}
		m.hooks[event] = filteredHooks
	}
}

func (m *Manager) callHookHandler(plugin *Plugin, handler string, event *PluginEvent) error {
	// Placeholder implementation for calling hook handlers
	// This would vary by plugin type
	logging.Debug("calling hook handler", map[string]interface{}{
		"plugin":  plugin.Name,
		"handler": handler,
		"event":   event.Type,
	})
	return nil
}