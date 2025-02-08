CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE outbox_events (
    id UUID DEFAULT uuid_generate_v4(),
    aggregate_id TEXT NOT NULL,
    aggregate_type TEXT NOT NULL,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (id)
)

