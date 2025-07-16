-- LLM Avatar System tables for Phase 3: AI Integration

-- LLM Avatars
CREATE TABLE IF NOT EXISTS llm_avatars (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    model VARCHAR(100) NOT NULL,
    personality TEXT,
    knowledge TEXT,
    capabilities JSONB NOT NULL DEFAULT '[]',
    configuration JSONB NOT NULL DEFAULT '{}',
    state VARCHAR(20) NOT NULL DEFAULT 'active',
    position JSONB,
    avatar_3d JSONB,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_llm_avatars_session_id ON llm_avatars(session_id);
CREATE INDEX IF NOT EXISTS idx_llm_avatars_type ON llm_avatars(type);
CREATE INDEX IF NOT EXISTS idx_llm_avatars_provider ON llm_avatars(provider);
CREATE INDEX IF NOT EXISTS idx_llm_avatars_state ON llm_avatars(state);
CREATE INDEX IF NOT EXISTS idx_llm_avatars_created_at ON llm_avatars(created_at);

-- LLM Conversations
CREATE TABLE IF NOT EXISTS llm_conversations (
    id UUID PRIMARY KEY,
    avatar_id UUID NOT NULL REFERENCES llm_avatars(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    response TEXT NOT NULL,
    context TEXT,
    metadata JSONB NOT NULL DEFAULT '{}',
    usage JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_llm_conversations_avatar_id ON llm_conversations(avatar_id);
CREATE INDEX IF NOT EXISTS idx_llm_conversations_user_id ON llm_conversations(user_id);
CREATE INDEX IF NOT EXISTS idx_llm_conversations_created_at ON llm_conversations(created_at);

-- LLM Providers
CREATE TABLE IF NOT EXISTS llm_providers (
    id UUID PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    type VARCHAR(50) NOT NULL,
    endpoint VARCHAR(500) NOT NULL,
    api_key_hash VARCHAR(255),
    capabilities JSONB NOT NULL DEFAULT '[]',
    configuration JSONB NOT NULL DEFAULT '{}',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_llm_providers_name ON llm_providers(name);
CREATE INDEX IF NOT EXISTS idx_llm_providers_type ON llm_providers(type);
CREATE INDEX IF NOT EXISTS idx_llm_providers_status ON llm_providers(status);

-- LLM Models
CREATE TABLE IF NOT EXISTS llm_models (
    id UUID PRIMARY KEY,
    provider_id UUID NOT NULL REFERENCES llm_providers(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    capabilities JSONB NOT NULL DEFAULT '[]',
    pricing JSONB NOT NULL DEFAULT '{}',
    limits JSONB NOT NULL DEFAULT '{}',
    configuration JSONB NOT NULL DEFAULT '{}',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    UNIQUE(provider_id, name)
);

CREATE INDEX IF NOT EXISTS idx_llm_models_provider_id ON llm_models(provider_id);
CREATE INDEX IF NOT EXISTS idx_llm_models_name ON llm_models(name);
CREATE INDEX IF NOT EXISTS idx_llm_models_status ON llm_models(status);

-- LLM Avatar Interactions
CREATE TABLE IF NOT EXISTS llm_avatar_interactions (
    id UUID PRIMARY KEY,
    avatar_id UUID NOT NULL REFERENCES llm_avatars(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    interaction_type VARCHAR(50) NOT NULL,
    data JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_llm_avatar_interactions_avatar_id ON llm_avatar_interactions(avatar_id);
CREATE INDEX IF NOT EXISTS idx_llm_avatar_interactions_user_id ON llm_avatar_interactions(user_id);
CREATE INDEX IF NOT EXISTS idx_llm_avatar_interactions_type ON llm_avatar_interactions(interaction_type);
CREATE INDEX IF NOT EXISTS idx_llm_avatar_interactions_created_at ON llm_avatar_interactions(created_at);

-- LLM Usage Stats
CREATE TABLE IF NOT EXISTS llm_usage_stats (
    id UUID PRIMARY KEY,
    avatar_id UUID NOT NULL REFERENCES llm_avatars(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    model VARCHAR(100) NOT NULL,
    input_tokens INTEGER NOT NULL DEFAULT 0,
    output_tokens INTEGER NOT NULL DEFAULT 0,
    total_tokens INTEGER NOT NULL DEFAULT 0,
    processing_time FLOAT NOT NULL DEFAULT 0,
    cost FLOAT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_llm_usage_stats_avatar_id ON llm_usage_stats(avatar_id);
CREATE INDEX IF NOT EXISTS idx_llm_usage_stats_user_id ON llm_usage_stats(user_id);
CREATE INDEX IF NOT EXISTS idx_llm_usage_stats_provider ON llm_usage_stats(provider);
CREATE INDEX IF NOT EXISTS idx_llm_usage_stats_model ON llm_usage_stats(model);
CREATE INDEX IF NOT EXISTS idx_llm_usage_stats_created_at ON llm_usage_stats(created_at);

-- LLM Avatar Presets
CREATE TABLE IF NOT EXISTS llm_avatar_presets (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    model VARCHAR(100) NOT NULL,
    personality TEXT,
    knowledge TEXT,
    capabilities JSONB NOT NULL DEFAULT '[]',
    configuration JSONB NOT NULL DEFAULT '{}',
    avatar_3d JSONB,
    is_public BOOLEAN NOT NULL DEFAULT false,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_llm_avatar_presets_name ON llm_avatar_presets(name);
CREATE INDEX IF NOT EXISTS idx_llm_avatar_presets_type ON llm_avatar_presets(type);
CREATE INDEX IF NOT EXISTS idx_llm_avatar_presets_provider ON llm_avatar_presets(provider);
CREATE INDEX IF NOT EXISTS idx_llm_avatar_presets_is_public ON llm_avatar_presets(is_public);
CREATE INDEX IF NOT EXISTS idx_llm_avatar_presets_created_by ON llm_avatar_presets(created_by);

-- Triggers for updating timestamps
CREATE TRIGGER update_llm_avatars_updated_at
    BEFORE UPDATE ON llm_avatars
    FOR EACH ROW
    EXECUTE FUNCTION update_rtc_updated_at_column();

CREATE TRIGGER update_llm_providers_updated_at
    BEFORE UPDATE ON llm_providers
    FOR EACH ROW
    EXECUTE FUNCTION update_rtc_updated_at_column();

CREATE TRIGGER update_llm_models_updated_at
    BEFORE UPDATE ON llm_models
    FOR EACH ROW
    EXECUTE FUNCTION update_rtc_updated_at_column();

CREATE TRIGGER update_llm_avatar_presets_updated_at
    BEFORE UPDATE ON llm_avatar_presets
    FOR EACH ROW
    EXECUTE FUNCTION update_rtc_updated_at_column();