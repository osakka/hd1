package database

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
	"holodeck1/config"
	"holodeck1/logging"
)

type DB struct {
	*sql.DB
}

func NewConnection() (*DB, error) {
	dbHost := config.GetString("DB_HOST", "localhost")
	dbPort := config.GetString("DB_PORT", "5432")
	dbUser := config.GetString("DB_USER", "hd1")
	dbPassword := config.GetString("DB_PASSWORD", "hd1")
	dbName := config.GetString("DB_NAME", "hd1")
	dbSSLMode := config.GetString("DB_SSL_MODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logging.Info("database connection established", map[string]interface{}{
		"host": dbHost,
		"port": dbPort,
		"name": dbName,
	})

	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func (db *DB) InitializeSchema() error {
	schemaSQL := `
	-- Check if schema is already initialized
	SELECT COUNT(*) FROM information_schema.tables 
	WHERE table_schema = 'public' AND table_name = 'users';
	`

	var count int
	err := db.QueryRow(schemaSQL).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check schema: %w", err)
	}

	if count > 0 {
		logging.Info("database schema already initialized", nil)
		// Apply WebRTC schema for Phase 2
		err = db.applyWebRTCSchema()
		if err != nil {
			return fmt.Errorf("failed to apply WebRTC schema: %w", err)
		}
		// Apply OT schema for Phase 2
		err = db.applyOTSchema()
		if err != nil {
			return fmt.Errorf("failed to apply OT schema: %w", err)
		}
		// Apply Assets schema for Phase 2
		err = db.applyAssetsSchema()
		if err != nil {
			return fmt.Errorf("failed to apply Assets schema: %w", err)
		}
		// Apply LLM schema for Phase 3
		err = db.applyLLMSchema()
		if err != nil {
			return fmt.Errorf("failed to apply LLM schema: %w", err)
		}
		// Apply Clients schema for Phase 4
		err = db.applyClientsSchema()
		if err != nil {
			return fmt.Errorf("failed to apply Clients schema: %w", err)
		}
		// Apply Enterprise schema for Phase 4
		err = db.applyEnterpriseSchema()
		if err != nil {
			return fmt.Errorf("failed to apply Enterprise schema: %w", err)
		}
		return nil
	}

	// Read and execute schema file
	logging.Info("initializing database schema", nil)
	
	// This would normally read from schema.sql file
	// For now, we'll create tables programmatically
	err = db.createTables()
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	// Apply WebRTC schema for Phase 2
	err = db.applyWebRTCSchema()
	if err != nil {
		return fmt.Errorf("failed to apply WebRTC schema: %w", err)
	}

	// Apply OT schema for Phase 2
	err = db.applyOTSchema()
	if err != nil {
		return fmt.Errorf("failed to apply OT schema: %w", err)
	}

	// Apply Assets schema for Phase 2
	err = db.applyAssetsSchema()
	if err != nil {
		return fmt.Errorf("failed to apply Assets schema: %w", err)
	}

	// Apply LLM schema for Phase 3
	err = db.applyLLMSchema()
	if err != nil {
		return fmt.Errorf("failed to apply LLM schema: %w", err)
	}

	// Apply Clients schema for Phase 4
	err = db.applyClientsSchema()
	if err != nil {
		return fmt.Errorf("failed to apply Clients schema: %w", err)
	}

	// Apply Enterprise schema for Phase 4
	err = db.applyEnterpriseSchema()
	if err != nil {
		return fmt.Errorf("failed to apply Enterprise schema: %w", err)
	}

	logging.Info("database schema initialized successfully", nil)
	return nil
}

func (db *DB) createTables() error {
	// Execute schema creation
	schemas := []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`,
		`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`,
		createUsersTable(),
		createSessionsTable(),
		createParticipantsTable(),
		createServicesTable(),
		createServiceHealthTable(),
		createPermissionsTable(),
		createOAuthTokensTable(),
		createAuditLogTable(),
		createIndexes(),
		createTriggers(),
	}

	for _, schema := range schemas {
		_, err := db.Exec(schema)
		if err != nil {
			return fmt.Errorf("failed to execute schema: %w", err)
		}
	}

	return nil
}

func createUsersTable() string {
	return `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		email VARCHAR(255) UNIQUE NOT NULL,
		username VARCHAR(100) UNIQUE NOT NULL,
		password_hash VARCHAR(255),
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW(),
		last_login TIMESTAMP,
		status VARCHAR(20) DEFAULT 'active',
		profile JSONB DEFAULT '{}',
		metadata JSONB DEFAULT '{}'
	)`
}

func createSessionsTable() string {
	return `
	CREATE TABLE IF NOT EXISTS sessions (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name VARCHAR(255) NOT NULL,
		description TEXT,
		owner_id UUID NOT NULL REFERENCES users(id),
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW(),
		status VARCHAR(20) DEFAULT 'active',
		visibility VARCHAR(20) DEFAULT 'private',
		max_participants INTEGER DEFAULT 10,
		current_participants INTEGER DEFAULT 0,
		settings JSONB DEFAULT '{}',
		metadata JSONB DEFAULT '{}'
	)`
}

func createParticipantsTable() string {
	return `
	CREATE TABLE IF NOT EXISTS participants (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
		user_id UUID NOT NULL REFERENCES users(id),
		role VARCHAR(20) DEFAULT 'participant',
		joined_at TIMESTAMP DEFAULT NOW(),
		left_at TIMESTAMP,
		last_seen TIMESTAMP DEFAULT NOW(),
		permissions JSONB DEFAULT '{}',
		metadata JSONB DEFAULT '{}',
		UNIQUE(session_id, user_id)
	)`
}

func createServicesTable() string {
	return `
	CREATE TABLE IF NOT EXISTS services (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name VARCHAR(255) NOT NULL,
		description TEXT,
		type VARCHAR(50) NOT NULL,
		endpoint VARCHAR(500) NOT NULL,
		status VARCHAR(20) DEFAULT 'active',
		capabilities JSONB DEFAULT '[]',
		ui_mapping JSONB DEFAULT '{}',
		permissions JSONB DEFAULT '{}',
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW(),
		metadata JSONB DEFAULT '{}',
		UNIQUE(name)
	)`
}

func createServiceHealthTable() string {
	return `
	CREATE TABLE IF NOT EXISTS service_health (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
		status VARCHAR(20) NOT NULL,
		response_time INTEGER,
		error_message TEXT,
		checked_at TIMESTAMP DEFAULT NOW(),
		metadata JSONB DEFAULT '{}'
	)`
}

func createPermissionsTable() string {
	return `
	CREATE TABLE IF NOT EXISTS permissions (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		resource_type VARCHAR(50) NOT NULL,
		resource_id UUID NOT NULL,
		user_id UUID NOT NULL REFERENCES users(id),
		permission_type VARCHAR(50) NOT NULL,
		granted_at TIMESTAMP DEFAULT NOW(),
		granted_by UUID REFERENCES users(id),
		metadata JSONB DEFAULT '{}'
	)`
}

func createOAuthTokensTable() string {
	return `
	CREATE TABLE IF NOT EXISTS oauth_tokens (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID NOT NULL REFERENCES users(id),
		access_token VARCHAR(1000) NOT NULL,
		refresh_token VARCHAR(1000),
		token_type VARCHAR(20) DEFAULT 'Bearer',
		expires_at TIMESTAMP,
		scope VARCHAR(500),
		created_at TIMESTAMP DEFAULT NOW(),
		metadata JSONB DEFAULT '{}'
	)`
}

func createAuditLogTable() string {
	return `
	CREATE TABLE IF NOT EXISTS audit_log (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		user_id UUID REFERENCES users(id),
		session_id UUID REFERENCES sessions(id),
		action VARCHAR(100) NOT NULL,
		resource_type VARCHAR(50) NOT NULL,
		resource_id UUID,
		old_values JSONB,
		new_values JSONB,
		ip_address INET,
		user_agent TEXT,
		created_at TIMESTAMP DEFAULT NOW(),
		metadata JSONB DEFAULT '{}'
	)`
}

func createIndexes() string {
	return `
	CREATE INDEX IF NOT EXISTS idx_sessions_owner_id ON sessions(owner_id);
	CREATE INDEX IF NOT EXISTS idx_sessions_status ON sessions(status);
	CREATE INDEX IF NOT EXISTS idx_participants_session_id ON participants(session_id);
	CREATE INDEX IF NOT EXISTS idx_participants_user_id ON participants(user_id);
	CREATE INDEX IF NOT EXISTS idx_services_name ON services(name);
	CREATE INDEX IF NOT EXISTS idx_services_type ON services(type);
	CREATE INDEX IF NOT EXISTS idx_permissions_user_id ON permissions(user_id);
	CREATE INDEX IF NOT EXISTS idx_oauth_tokens_user_id ON oauth_tokens(user_id);
	CREATE INDEX IF NOT EXISTS idx_audit_log_user_id ON audit_log(user_id);
	`
}

func createTriggers() string {
	return `
	CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = NOW();
		RETURN NEW;
	END;
	$$ language 'plpgsql';

	DROP TRIGGER IF EXISTS update_users_updated_at ON users;
	CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
		FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

	DROP TRIGGER IF EXISTS update_sessions_updated_at ON sessions;
	CREATE TRIGGER update_sessions_updated_at BEFORE UPDATE ON sessions
		FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

	DROP TRIGGER IF EXISTS update_services_updated_at ON services;
	CREATE TRIGGER update_services_updated_at BEFORE UPDATE ON services
		FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
	`
}

func (db *DB) applyWebRTCSchema() error {
	// Check if WebRTC tables already exist
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM information_schema.tables 
		WHERE table_schema = 'public' AND table_name = 'rtc_sessions'
	`).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check WebRTC schema: %w", err)
	}

	if count > 0 {
		logging.Info("WebRTC schema already exists", nil)
		return nil
	}

	// Read WebRTC schema file
	schemaPath := filepath.Join("src", "database", "webrtc_schema.sql")
	schemaData, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read WebRTC schema file: %w", err)
	}

	// Execute WebRTC schema
	_, err = db.Exec(string(schemaData))
	if err != nil {
		return fmt.Errorf("failed to execute WebRTC schema: %w", err)
	}

	logging.Info("WebRTC schema applied successfully", nil)
	return nil
}

func (db *DB) applyOTSchema() error {
	// Check if OT tables already exist
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM information_schema.tables 
		WHERE table_schema = 'public' AND table_name = 'ot_documents'
	`).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check OT schema: %w", err)
	}

	if count > 0 {
		logging.Info("OT schema already exists", nil)
		return nil
	}

	// Read OT schema file
	schemaPath := filepath.Join("src", "database", "ot_schema.sql")
	schemaData, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read OT schema file: %w", err)
	}

	// Execute OT schema
	_, err = db.Exec(string(schemaData))
	if err != nil {
		return fmt.Errorf("failed to execute OT schema: %w", err)
	}

	logging.Info("OT schema applied successfully", nil)
	return nil
}

func (db *DB) applyAssetsSchema() error {
	// Check if Assets tables already exist
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM information_schema.tables 
		WHERE table_schema = 'public' AND table_name = 'assets'
	`).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check Assets schema: %w", err)
	}

	if count > 0 {
		logging.Info("Assets schema already exists", nil)
		return nil
	}

	// Read Assets schema file
	schemaPath := filepath.Join("src", "database", "assets_schema.sql")
	schemaData, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read Assets schema file: %w", err)
	}

	// Execute Assets schema
	_, err = db.Exec(string(schemaData))
	if err != nil {
		return fmt.Errorf("failed to execute Assets schema: %w", err)
	}

	logging.Info("Assets schema applied successfully", nil)
	return nil
}

func (db *DB) applyLLMSchema() error {
	// Check if LLM tables already exist
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM information_schema.tables 
		WHERE table_schema = 'public' AND table_name = 'llm_avatars'
	`).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check LLM schema: %w", err)
	}

	if count > 0 {
		logging.Info("LLM schema already exists", nil)
		return nil
	}

	// Read LLM schema file
	schemaPath := filepath.Join("src", "database", "llm_schema.sql")
	schemaData, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read LLM schema file: %w", err)
	}

	// Execute LLM schema
	_, err = db.Exec(string(schemaData))
	if err != nil {
		return fmt.Errorf("failed to execute LLM schema: %w", err)
	}

	logging.Info("LLM schema applied successfully", nil)
	return nil
}

func (db *DB) applyClientsSchema() error {
	// Check if Clients tables already exist
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM information_schema.tables 
		WHERE table_schema = 'public' AND table_name = 'clients'
	`).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check Clients schema: %w", err)
	}

	if count > 0 {
		logging.Info("Clients schema already exists", nil)
		return nil
	}

	// Read Clients schema file
	schemaPath := filepath.Join("src", "database", "clients_schema.sql")
	schemaData, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read Clients schema file: %w", err)
	}

	// Execute Clients schema
	_, err = db.Exec(string(schemaData))
	if err != nil {
		return fmt.Errorf("failed to execute Clients schema: %w", err)
	}

	logging.Info("Clients schema applied successfully", nil)
	return nil
}

func (db *DB) applyEnterpriseSchema() error {
	// Check if Enterprise tables already exist
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM information_schema.tables 
		WHERE table_schema = 'public' AND table_name = 'organizations'
	`).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check Enterprise schema: %w", err)
	}

	if count > 0 {
		logging.Info("Enterprise schema already exists", nil)
		return nil
	}

	// Read Enterprise schema file
	schemaPath := filepath.Join("src", "database", "enterprise_schema.sql")
	schemaData, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read Enterprise schema file: %w", err)
	}

	// Execute Enterprise schema
	_, err = db.Exec(string(schemaData))
	if err != nil {
		return fmt.Errorf("failed to execute Enterprise schema: %w", err)
	}

	logging.Info("Enterprise schema applied successfully", nil)
	return nil
}

func (db *DB) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}