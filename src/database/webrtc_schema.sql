-- WebRTC tables for Phase 2: Collaboration

-- RTC Sessions
CREATE TABLE IF NOT EXISTS rtc_sessions (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_rtc_sessions_session_id ON rtc_sessions(session_id);
CREATE INDEX IF NOT EXISTS idx_rtc_sessions_created_at ON rtc_sessions(created_at);

-- RTC Participants
CREATE TABLE IF NOT EXISTS rtc_participants (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id UUID NOT NULL REFERENCES rtc_sessions(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'participant',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_rtc_participants_user_id ON rtc_participants(user_id);
CREATE INDEX IF NOT EXISTS idx_rtc_participants_session_id ON rtc_participants(session_id);
CREATE INDEX IF NOT EXISTS idx_rtc_participants_is_active ON rtc_participants(is_active);

-- RTC Connection Stats
CREATE TABLE IF NOT EXISTS rtc_connection_stats (
    id UUID PRIMARY KEY,
    participant_id UUID NOT NULL REFERENCES rtc_participants(id) ON DELETE CASCADE,
    connection_state VARCHAR(50) NOT NULL,
    ice_state VARCHAR(50) NOT NULL,
    data_channel_state VARCHAR(50) NOT NULL,
    bytes_sent BIGINT NOT NULL DEFAULT 0,
    bytes_received BIGINT NOT NULL DEFAULT 0,
    packets_sent BIGINT NOT NULL DEFAULT 0,
    packets_received BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_rtc_connection_stats_participant_id ON rtc_connection_stats(participant_id);
CREATE INDEX IF NOT EXISTS idx_rtc_connection_stats_created_at ON rtc_connection_stats(created_at);

-- RTC Data Messages
CREATE TABLE IF NOT EXISTS rtc_data_messages (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES rtc_sessions(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message_type VARCHAR(50) NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_rtc_data_messages_session_id ON rtc_data_messages(session_id);
CREATE INDEX IF NOT EXISTS idx_rtc_data_messages_sender_id ON rtc_data_messages(sender_id);
CREATE INDEX IF NOT EXISTS idx_rtc_data_messages_created_at ON rtc_data_messages(created_at);

-- RTC Session Events
CREATE TABLE IF NOT EXISTS rtc_session_events (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES rtc_sessions(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    event_type VARCHAR(50) NOT NULL,
    event_data JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_rtc_session_events_session_id ON rtc_session_events(session_id);
CREATE INDEX IF NOT EXISTS idx_rtc_session_events_user_id ON rtc_session_events(user_id);
CREATE INDEX IF NOT EXISTS idx_rtc_session_events_event_type ON rtc_session_events(event_type);
CREATE INDEX IF NOT EXISTS idx_rtc_session_events_created_at ON rtc_session_events(created_at);

-- Triggers for updating timestamps
CREATE OR REPLACE FUNCTION update_rtc_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language plpgsql;

CREATE TRIGGER update_rtc_sessions_updated_at
    BEFORE UPDATE ON rtc_sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_rtc_updated_at_column();

CREATE TRIGGER update_rtc_participants_updated_at
    BEFORE UPDATE ON rtc_participants
    FOR EACH ROW
    EXECUTE FUNCTION update_rtc_updated_at_column();