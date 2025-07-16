-- Asset Management tables for Phase 2: Collaboration

-- Assets
CREATE TABLE IF NOT EXISTS assets (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    size BIGINT NOT NULL,
    hash VARCHAR(64) NOT NULL,
    url VARCHAR(500) NOT NULL,
    storage_path VARCHAR(500) NOT NULL,
    tags JSONB NOT NULL DEFAULT '[]',
    metadata JSONB NOT NULL DEFAULT '{}',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_assets_session_id ON assets(session_id);
CREATE INDEX IF NOT EXISTS idx_assets_user_id ON assets(user_id);
CREATE INDEX IF NOT EXISTS idx_assets_type ON assets(type);
CREATE INDEX IF NOT EXISTS idx_assets_status ON assets(status);
CREATE INDEX IF NOT EXISTS idx_assets_hash ON assets(hash);
CREATE INDEX IF NOT EXISTS idx_assets_created_at ON assets(created_at);
CREATE INDEX IF NOT EXISTS idx_assets_tags ON assets USING gin(tags);

-- Asset Versions
CREATE TABLE IF NOT EXISTS asset_versions (
    id UUID PRIMARY KEY,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    hash VARCHAR(64) NOT NULL,
    size BIGINT NOT NULL,
    path VARCHAR(500) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    UNIQUE(asset_id, version)
);

CREATE INDEX IF NOT EXISTS idx_asset_versions_asset_id ON asset_versions(asset_id);
CREATE INDEX IF NOT EXISTS idx_asset_versions_version ON asset_versions(version);
CREATE INDEX IF NOT EXISTS idx_asset_versions_hash ON asset_versions(hash);
CREATE INDEX IF NOT EXISTS idx_asset_versions_created_at ON asset_versions(created_at);

-- Asset Usage Tracking
CREATE TABLE IF NOT EXISTS asset_usage (
    id UUID PRIMARY KEY,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    usage VARCHAR(50) NOT NULL,
    context JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_asset_usage_asset_id ON asset_usage(asset_id);
CREATE INDEX IF NOT EXISTS idx_asset_usage_session_id ON asset_usage(session_id);
CREATE INDEX IF NOT EXISTS idx_asset_usage_user_id ON asset_usage(user_id);
CREATE INDEX IF NOT EXISTS idx_asset_usage_usage ON asset_usage(usage);
CREATE INDEX IF NOT EXISTS idx_asset_usage_created_at ON asset_usage(created_at);

-- Asset Dependencies (for tracking asset relationships)
CREATE TABLE IF NOT EXISTS asset_dependencies (
    id UUID PRIMARY KEY,
    parent_asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    child_asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    dependency_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    UNIQUE(parent_asset_id, child_asset_id, dependency_type)
);

CREATE INDEX IF NOT EXISTS idx_asset_dependencies_parent_asset_id ON asset_dependencies(parent_asset_id);
CREATE INDEX IF NOT EXISTS idx_asset_dependencies_child_asset_id ON asset_dependencies(child_asset_id);
CREATE INDEX IF NOT EXISTS idx_asset_dependencies_dependency_type ON asset_dependencies(dependency_type);

-- Asset Tags (for better tag management)
CREATE TABLE IF NOT EXISTS asset_tags (
    id UUID PRIMARY KEY,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    tag VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    UNIQUE(asset_id, tag)
);

CREATE INDEX IF NOT EXISTS idx_asset_tags_asset_id ON asset_tags(asset_id);
CREATE INDEX IF NOT EXISTS idx_asset_tags_tag ON asset_tags(tag);

-- Asset Collections (for grouping assets)
CREATE TABLE IF NOT EXISTS asset_collections (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_asset_collections_session_id ON asset_collections(session_id);
CREATE INDEX IF NOT EXISTS idx_asset_collections_user_id ON asset_collections(user_id);
CREATE INDEX IF NOT EXISTS idx_asset_collections_name ON asset_collections(name);

-- Asset Collection Items
CREATE TABLE IF NOT EXISTS asset_collection_items (
    id UUID PRIMARY KEY,
    collection_id UUID NOT NULL REFERENCES asset_collections(id) ON DELETE CASCADE,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    order_index INTEGER NOT NULL DEFAULT 0,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    UNIQUE(collection_id, asset_id)
);

CREATE INDEX IF NOT EXISTS idx_asset_collection_items_collection_id ON asset_collection_items(collection_id);
CREATE INDEX IF NOT EXISTS idx_asset_collection_items_asset_id ON asset_collection_items(asset_id);
CREATE INDEX IF NOT EXISTS idx_asset_collection_items_order_index ON asset_collection_items(order_index);

-- Triggers for updating timestamps
CREATE TRIGGER update_assets_updated_at
    BEFORE UPDATE ON assets
    FOR EACH ROW
    EXECUTE FUNCTION update_rtc_updated_at_column();

CREATE TRIGGER update_asset_collections_updated_at
    BEFORE UPDATE ON asset_collections
    FOR EACH ROW
    EXECUTE FUNCTION update_rtc_updated_at_column();