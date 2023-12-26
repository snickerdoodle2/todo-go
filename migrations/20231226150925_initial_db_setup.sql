-- Add migration script here
CREATE TABLE IF NOT EXISTS todos
(
    id         UUID PRIMARY KEY,
    content    TEXT                  NOT NULL UNIQUE,
    finished   BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMPTZ           NOT NULL,
    updated_at TIMESTAMPTZ           NOT NULL
)