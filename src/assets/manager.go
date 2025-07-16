package assets

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"holodeck1/database"
	"holodeck1/logging"
)

type Manager struct {
	db           *database.DB
	storagePath  string
	maxFileSize  int64
	allowedTypes map[string]bool
	mu           sync.RWMutex
}

type Asset struct {
	ID          uuid.UUID         `json:"id"`
	SessionID   uuid.UUID         `json:"session_id"`
	UserID      uuid.UUID         `json:"user_id"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	MimeType    string            `json:"mime_type"`
	Size        int64             `json:"size"`
	Hash        string            `json:"hash"`
	URL         string            `json:"url"`
	StoragePath string            `json:"storage_path"`
	Tags        []string          `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
	Status      string            `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type AssetVersion struct {
	ID        uuid.UUID `json:"id"`
	AssetID   uuid.UUID `json:"asset_id"`
	Version   int       `json:"version"`
	Hash      string    `json:"hash"`
	Size      int64     `json:"size"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
}

type AssetUsage struct {
	ID        uuid.UUID `json:"id"`
	AssetID   uuid.UUID `json:"asset_id"`
	SessionID uuid.UUID `json:"session_id"`
	UserID    uuid.UUID `json:"user_id"`
	Usage     string    `json:"usage"`
	Context   map[string]interface{} `json:"context"`
	CreatedAt time.Time `json:"created_at"`
}

type UploadRequest struct {
	SessionID uuid.UUID `json:"session_id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Tags      []string  `json:"tags"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type AssetFilter struct {
	SessionID uuid.UUID `json:"session_id"`
	UserID    uuid.UUID `json:"user_id"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Tags      []string  `json:"tags"`
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
}

func NewManager(db *database.DB, storagePath string) *Manager {
	// Create storage directory if it doesn't exist
	os.MkdirAll(storagePath, 0755)

	return &Manager{
		db:          db,
		storagePath: storagePath,
		maxFileSize: 100 * 1024 * 1024, // 100MB
		allowedTypes: map[string]bool{
			"image/jpeg":    true,
			"image/png":     true,
			"image/gif":     true,
			"image/webp":    true,
			"model/gltf":    true,
			"model/glb":     true,
			"audio/mpeg":    true,
			"audio/wav":     true,
			"audio/ogg":     true,
			"video/mp4":     true,
			"video/webm":    true,
			"text/plain":    true,
			"application/json": true,
			"application/javascript": true,
			"text/css":      true,
			"text/html":     true,
		},
	}
}

func (m *Manager) UploadAsset(ctx context.Context, req *UploadRequest, file multipart.File, header *multipart.FileHeader) (*Asset, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Validate file size
	if header.Size > m.maxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size")
	}

	// Validate file type
	mimeType := header.Header.Get("Content-Type")
	if !m.allowedTypes[mimeType] {
		return nil, fmt.Errorf("file type not allowed: %s", mimeType)
	}

	// Generate asset ID
	assetID := uuid.New()
	
	// Create hash
	hash, err := m.calculateFileHash(file)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate file hash: %w", err)
	}

	// Reset file pointer
	file.Seek(0, 0)

	// Create storage path
	assetDir := filepath.Join(m.storagePath, assetID.String())
	err = os.MkdirAll(assetDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create asset directory: %w", err)
	}

	// Save file
	fileName := fmt.Sprintf("%s_%s", hash, header.Filename)
	filePath := filepath.Join(assetDir, fileName)
	
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Create asset record
	asset := &Asset{
		ID:          assetID,
		SessionID:   req.SessionID,
		UserID:      req.UserID,
		Name:        req.Name,
		Type:        req.Type,
		MimeType:    mimeType,
		Size:        header.Size,
		Hash:        hash,
		URL:         fmt.Sprintf("/api/assets/%s", assetID.String()),
		StoragePath: filePath,
		Tags:        req.Tags,
		Metadata:    req.Metadata,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Store in database
	query := `
		INSERT INTO assets (id, session_id, user_id, name, type, mime_type, size, hash, url, storage_path, tags, metadata, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	
	tagsJSON, _ := json.Marshal(asset.Tags)
	metadataJSON, _ := json.Marshal(asset.Metadata)
	
	_, err = m.db.ExecContext(ctx, query,
		asset.ID, asset.SessionID, asset.UserID, asset.Name, asset.Type,
		asset.MimeType, asset.Size, asset.Hash, asset.URL, asset.StoragePath,
		tagsJSON, metadataJSON, asset.Status, asset.CreatedAt, asset.UpdatedAt)
	if err != nil {
		// Clean up file if database insert fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to store asset in database: %w", err)
	}

	// Create first version
	version := &AssetVersion{
		ID:        uuid.New(),
		AssetID:   assetID,
		Version:   1,
		Hash:      hash,
		Size:      header.Size,
		Path:      filePath,
		CreatedAt: time.Now(),
	}

	query = `
		INSERT INTO asset_versions (id, asset_id, version, hash, size, path, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = m.db.ExecContext(ctx, query, version.ID, version.AssetID, version.Version, version.Hash, version.Size, version.Path, version.CreatedAt)
	if err != nil {
		logging.Error("failed to create asset version", map[string]interface{}{
			"asset_id": assetID,
			"error":    err.Error(),
		})
	}

	logging.Info("asset uploaded successfully", map[string]interface{}{
		"asset_id":   assetID,
		"session_id": req.SessionID,
		"user_id":    req.UserID,
		"name":       req.Name,
		"type":       req.Type,
		"size":       header.Size,
	})

	return asset, nil
}

func (m *Manager) GetAsset(ctx context.Context, assetID uuid.UUID) (*Asset, error) {
	query := `
		SELECT id, session_id, user_id, name, type, mime_type, size, hash, url, storage_path, tags, metadata, status, created_at, updated_at
		FROM assets
		WHERE id = $1
	`
	
	var asset Asset
	var tagsJSON, metadataJSON []byte
	
	err := m.db.QueryRowContext(ctx, query, assetID).Scan(
		&asset.ID, &asset.SessionID, &asset.UserID, &asset.Name, &asset.Type,
		&asset.MimeType, &asset.Size, &asset.Hash, &asset.URL, &asset.StoragePath,
		&tagsJSON, &metadataJSON, &asset.Status, &asset.CreatedAt, &asset.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("asset not found")
		}
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}

	// Unmarshal JSON fields
	json.Unmarshal(tagsJSON, &asset.Tags)
	json.Unmarshal(metadataJSON, &asset.Metadata)

	return &asset, nil
}

func (m *Manager) GetAssetsBySession(ctx context.Context, filter *AssetFilter) ([]*Asset, error) {
	query := `
		SELECT id, session_id, user_id, name, type, mime_type, size, hash, url, storage_path, tags, metadata, status, created_at, updated_at
		FROM assets
		WHERE session_id = $1
	`
	args := []interface{}{filter.SessionID}
	argIndex := 2

	if filter.UserID != uuid.Nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, filter.UserID)
		argIndex++
	}

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
	}

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query assets: %w", err)
	}
	defer rows.Close()

	var assets []*Asset
	for rows.Next() {
		var asset Asset
		var tagsJSON, metadataJSON []byte
		
		err := rows.Scan(
			&asset.ID, &asset.SessionID, &asset.UserID, &asset.Name, &asset.Type,
			&asset.MimeType, &asset.Size, &asset.Hash, &asset.URL, &asset.StoragePath,
			&tagsJSON, &metadataJSON, &asset.Status, &asset.CreatedAt, &asset.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset: %w", err)
		}

		// Unmarshal JSON fields
		json.Unmarshal(tagsJSON, &asset.Tags)
		json.Unmarshal(metadataJSON, &asset.Metadata)

		assets = append(assets, &asset)
	}

	return assets, nil
}

func (m *Manager) DeleteAsset(ctx context.Context, assetID uuid.UUID, userID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Get asset to check ownership and get storage path
	asset, err := m.GetAsset(ctx, assetID)
	if err != nil {
		return err
	}

	// Check ownership
	if asset.UserID != userID {
		return fmt.Errorf("unauthorized to delete asset")
	}

	// Delete from database
	query := `DELETE FROM assets WHERE id = $1`
	result, err := m.db.ExecContext(ctx, query, assetID)
	if err != nil {
		return fmt.Errorf("failed to delete asset from database: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("asset not found")
	}

	// Delete file from storage
	assetDir := filepath.Dir(asset.StoragePath)
	err = os.RemoveAll(assetDir)
	if err != nil {
		logging.Error("failed to delete asset files", map[string]interface{}{
			"asset_id": assetID,
			"path":     assetDir,
			"error":    err.Error(),
		})
	}

	logging.Info("asset deleted successfully", map[string]interface{}{
		"asset_id": assetID,
		"user_id":  userID,
	})

	return nil
}

func (m *Manager) GetAssetContent(ctx context.Context, assetID uuid.UUID) ([]byte, string, error) {
	asset, err := m.GetAsset(ctx, assetID)
	if err != nil {
		return nil, "", err
	}

	// Check if file exists
	if _, err := os.Stat(asset.StoragePath); os.IsNotExist(err) {
		return nil, "", fmt.Errorf("asset file not found")
	}

	// Read file
	content, err := os.ReadFile(asset.StoragePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read asset file: %w", err)
	}

	return content, asset.MimeType, nil
}

func (m *Manager) TrackAssetUsage(ctx context.Context, assetID uuid.UUID, sessionID uuid.UUID, userID uuid.UUID, usage string, context map[string]interface{}) error {
	usageRecord := &AssetUsage{
		ID:        uuid.New(),
		AssetID:   assetID,
		SessionID: sessionID,
		UserID:    userID,
		Usage:     usage,
		Context:   context,
		CreatedAt: time.Now(),
	}

	query := `
		INSERT INTO asset_usage (id, asset_id, session_id, user_id, usage, context, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	contextJSON, _ := json.Marshal(usageRecord.Context)
	
	_, err := m.db.ExecContext(ctx, query,
		usageRecord.ID, usageRecord.AssetID, usageRecord.SessionID,
		usageRecord.UserID, usageRecord.Usage, contextJSON, usageRecord.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to track asset usage: %w", err)
	}

	return nil
}

func (m *Manager) GetAssetUsage(ctx context.Context, assetID uuid.UUID) ([]*AssetUsage, error) {
	query := `
		SELECT id, asset_id, session_id, user_id, usage, context, created_at
		FROM asset_usage
		WHERE asset_id = $1
		ORDER BY created_at DESC
	`
	
	rows, err := m.db.QueryContext(ctx, query, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to query asset usage: %w", err)
	}
	defer rows.Close()

	var usages []*AssetUsage
	for rows.Next() {
		var usage AssetUsage
		var contextJSON []byte
		
		err := rows.Scan(
			&usage.ID, &usage.AssetID, &usage.SessionID,
			&usage.UserID, &usage.Usage, &contextJSON, &usage.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset usage: %w", err)
		}

		// Unmarshal context
		json.Unmarshal(contextJSON, &usage.Context)

		usages = append(usages, &usage)
	}

	return usages, nil
}

func (m *Manager) calculateFileHash(file multipart.File) (string, error) {
	hasher := md5.New()
	_, err := io.Copy(hasher, file)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (m *Manager) CleanupOrphanedAssets(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Find assets that haven't been used in 30 days
	query := `
		SELECT id, storage_path FROM assets
		WHERE updated_at < $1 AND status = 'active'
	`
	
	cutoff := time.Now().AddDate(0, 0, -30)
	rows, err := m.db.QueryContext(ctx, query, cutoff)
	if err != nil {
		logging.Error("failed to query orphaned assets", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()

	var orphanedAssets []struct {
		ID          uuid.UUID
		StoragePath string
	}

	for rows.Next() {
		var asset struct {
			ID          uuid.UUID
			StoragePath string
		}
		err := rows.Scan(&asset.ID, &asset.StoragePath)
		if err != nil {
			logging.Error("failed to scan orphaned asset", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}
		orphanedAssets = append(orphanedAssets, asset)
	}

	// Clean up orphaned assets
	for _, asset := range orphanedAssets {
		// Mark as deleted in database
		query := `UPDATE assets SET status = 'deleted', updated_at = $1 WHERE id = $2`
		_, err := m.db.ExecContext(ctx, query, time.Now(), asset.ID)
		if err != nil {
			logging.Error("failed to mark asset as deleted", map[string]interface{}{
				"asset_id": asset.ID,
				"error":    err.Error(),
			})
			continue
		}

		// Delete file
		assetDir := filepath.Dir(asset.StoragePath)
		err = os.RemoveAll(assetDir)
		if err != nil {
			logging.Error("failed to delete asset directory", map[string]interface{}{
				"asset_id": asset.ID,
				"path":     assetDir,
				"error":    err.Error(),
			})
		}

		logging.Info("cleaned up orphaned asset", map[string]interface{}{
			"asset_id": asset.ID,
		})
	}
}

func (m *Manager) StartCleanupWorker(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.CleanupOrphanedAssets(ctx)
		}
	}
}