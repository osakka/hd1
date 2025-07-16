-- Enterprise Features schema for Phase 4: Universal Platform

-- Organizations
CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    domain VARCHAR(255) UNIQUE,
    logo_url VARCHAR(500),
    website VARCHAR(500),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    subscription_tier VARCHAR(50) NOT NULL DEFAULT 'basic',
    max_users INTEGER NOT NULL DEFAULT 100,
    max_sessions INTEGER NOT NULL DEFAULT 50,
    max_storage_gb INTEGER NOT NULL DEFAULT 10,
    features JSONB NOT NULL DEFAULT '{}',
    settings JSONB NOT NULL DEFAULT '{}',
    billing_info JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'
);

CREATE INDEX IF NOT EXISTS idx_organizations_name ON organizations(name);
CREATE INDEX IF NOT EXISTS idx_organizations_domain ON organizations(domain);
CREATE INDEX IF NOT EXISTS idx_organizations_status ON organizations(status);
CREATE INDEX IF NOT EXISTS idx_organizations_subscription_tier ON organizations(subscription_tier);

-- Organization Users (Many-to-Many relationship)
CREATE TABLE IF NOT EXISTS organization_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    permissions JSONB NOT NULL DEFAULT '[]',
    department VARCHAR(100),
    job_title VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    invited_by UUID REFERENCES users(id),
    invited_at TIMESTAMP WITH TIME ZONE,
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_active TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}',
    UNIQUE(organization_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_organization_users_organization_id ON organization_users(organization_id);
CREATE INDEX IF NOT EXISTS idx_organization_users_user_id ON organization_users(user_id);
CREATE INDEX IF NOT EXISTS idx_organization_users_role ON organization_users(role);
CREATE INDEX IF NOT EXISTS idx_organization_users_status ON organization_users(status);
CREATE INDEX IF NOT EXISTS idx_organization_users_department ON organization_users(department);

-- Roles and Permissions
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    permissions JSONB NOT NULL DEFAULT '[]',
    is_system_role BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}',
    UNIQUE(organization_id, name)
);

CREATE INDEX IF NOT EXISTS idx_roles_organization_id ON roles(organization_id);
CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);
CREATE INDEX IF NOT EXISTS idx_roles_is_system_role ON roles(is_system_role);

-- User Roles (Many-to-Many relationship)
CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    granted_by UUID REFERENCES users(id),
    granted_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}',
    UNIQUE(user_id, role_id, organization_id)
);

CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_organization_id ON user_roles(organization_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_expires_at ON user_roles(expires_at);

-- Analytics Events
CREATE TABLE IF NOT EXISTS analytics_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    session_id UUID REFERENCES sessions(id) ON DELETE SET NULL,
    event_type VARCHAR(100) NOT NULL,
    event_category VARCHAR(50) NOT NULL,
    event_action VARCHAR(100) NOT NULL,
    event_label VARCHAR(255),
    event_value NUMERIC,
    properties JSONB NOT NULL DEFAULT '{}',
    user_agent TEXT,
    ip_address INET,
    referrer VARCHAR(500),
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'
);

CREATE INDEX IF NOT EXISTS idx_analytics_events_organization_id ON analytics_events(organization_id);
CREATE INDEX IF NOT EXISTS idx_analytics_events_user_id ON analytics_events(user_id);
CREATE INDEX IF NOT EXISTS idx_analytics_events_session_id ON analytics_events(session_id);
CREATE INDEX IF NOT EXISTS idx_analytics_events_event_type ON analytics_events(event_type);
CREATE INDEX IF NOT EXISTS idx_analytics_events_event_category ON analytics_events(event_category);
CREATE INDEX IF NOT EXISTS idx_analytics_events_timestamp ON analytics_events(timestamp);
CREATE INDEX IF NOT EXISTS idx_analytics_events_created_at ON analytics_events(created_at);

-- Analytics Aggregates (for performance)
CREATE TABLE IF NOT EXISTS analytics_aggregates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    aggregate_type VARCHAR(50) NOT NULL,
    dimension VARCHAR(100) NOT NULL,
    dimension_value VARCHAR(255) NOT NULL,
    metric VARCHAR(100) NOT NULL,
    value NUMERIC NOT NULL,
    period_type VARCHAR(20) NOT NULL, -- hour, day, week, month, year
    period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}',
    UNIQUE(organization_id, aggregate_type, dimension, dimension_value, metric, period_type, period_start)
);

CREATE INDEX IF NOT EXISTS idx_analytics_aggregates_organization_id ON analytics_aggregates(organization_id);
CREATE INDEX IF NOT EXISTS idx_analytics_aggregates_aggregate_type ON analytics_aggregates(aggregate_type);
CREATE INDEX IF NOT EXISTS idx_analytics_aggregates_dimension ON analytics_aggregates(dimension);
CREATE INDEX IF NOT EXISTS idx_analytics_aggregates_metric ON analytics_aggregates(metric);
CREATE INDEX IF NOT EXISTS idx_analytics_aggregates_period_type ON analytics_aggregates(period_type);
CREATE INDEX IF NOT EXISTS idx_analytics_aggregates_period_start ON analytics_aggregates(period_start);

-- Security Audit Log
CREATE TABLE IF NOT EXISTS security_audit_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    event_type VARCHAR(100) NOT NULL,
    event_category VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50),
    resource_id UUID,
    action VARCHAR(100) NOT NULL,
    result VARCHAR(20) NOT NULL, -- success, failure, warning
    risk_level VARCHAR(20) NOT NULL DEFAULT 'low', -- low, medium, high, critical
    details JSONB NOT NULL DEFAULT '{}',
    ip_address INET,
    user_agent TEXT,
    geo_location JSONB,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'
);

CREATE INDEX IF NOT EXISTS idx_security_audit_log_organization_id ON security_audit_log(organization_id);
CREATE INDEX IF NOT EXISTS idx_security_audit_log_user_id ON security_audit_log(user_id);
CREATE INDEX IF NOT EXISTS idx_security_audit_log_event_type ON security_audit_log(event_type);
CREATE INDEX IF NOT EXISTS idx_security_audit_log_event_category ON security_audit_log(event_category);
CREATE INDEX IF NOT EXISTS idx_security_audit_log_result ON security_audit_log(result);
CREATE INDEX IF NOT EXISTS idx_security_audit_log_risk_level ON security_audit_log(risk_level);
CREATE INDEX IF NOT EXISTS idx_security_audit_log_timestamp ON security_audit_log(timestamp);

-- Compliance Records
CREATE TABLE IF NOT EXISTS compliance_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    compliance_type VARCHAR(50) NOT NULL, -- gdpr, hipaa, sox, etc.
    record_type VARCHAR(50) NOT NULL, -- data_processing, access_request, deletion_request
    subject_type VARCHAR(50) NOT NULL, -- user, data, system
    subject_id UUID,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    priority VARCHAR(20) NOT NULL DEFAULT 'medium',
    description TEXT,
    request_details JSONB NOT NULL DEFAULT '{}',
    response_details JSONB NOT NULL DEFAULT '{}',
    evidence JSONB NOT NULL DEFAULT '{}',
    due_date TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    assigned_to UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    metadata JSONB NOT NULL DEFAULT '{}'
);

CREATE INDEX IF NOT EXISTS idx_compliance_records_organization_id ON compliance_records(organization_id);
CREATE INDEX IF NOT EXISTS idx_compliance_records_compliance_type ON compliance_records(compliance_type);
CREATE INDEX IF NOT EXISTS idx_compliance_records_record_type ON compliance_records(record_type);
CREATE INDEX IF NOT EXISTS idx_compliance_records_status ON compliance_records(status);
CREATE INDEX IF NOT EXISTS idx_compliance_records_priority ON compliance_records(priority);
CREATE INDEX IF NOT EXISTS idx_compliance_records_due_date ON compliance_records(due_date);
CREATE INDEX IF NOT EXISTS idx_compliance_records_created_by ON compliance_records(created_by);
CREATE INDEX IF NOT EXISTS idx_compliance_records_assigned_to ON compliance_records(assigned_to);

-- Update organizations table to reference root sessions
ALTER TABLE sessions ADD COLUMN IF NOT EXISTS organization_id UUID REFERENCES organizations(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_sessions_organization_id ON sessions(organization_id);

-- Update users table for enterprise features
ALTER TABLE users ADD COLUMN IF NOT EXISTS organization_id UUID REFERENCES organizations(id) ON DELETE SET NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_system_user BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_password_change TIMESTAMP WITH TIME ZONE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_expires_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS security_settings JSONB NOT NULL DEFAULT '{}';

CREATE INDEX IF NOT EXISTS idx_users_organization_id ON users(organization_id);
CREATE INDEX IF NOT EXISTS idx_users_is_system_user ON users(is_system_user);
CREATE INDEX IF NOT EXISTS idx_users_password_expires_at ON users(password_expires_at);

-- Triggers for updating timestamps
CREATE TRIGGER update_organizations_updated_at
    BEFORE UPDATE ON organizations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_organization_users_updated_at
    BEFORE UPDATE ON organization_users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_roles_updated_at
    BEFORE UPDATE ON roles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_analytics_aggregates_updated_at
    BEFORE UPDATE ON analytics_aggregates
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_compliance_records_updated_at
    BEFORE UPDATE ON compliance_records
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();