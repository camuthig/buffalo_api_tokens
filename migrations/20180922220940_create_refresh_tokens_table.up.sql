CREATE TABLE refresh_tokens (
    id CHAR(80) PRIMARY KEY,
    user_id UUID REFERENCES users (id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);