CREATE TABLE oauth_access_tokens (
    id UUID PRIMARY KEY,
    token_hash VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID REFERENCES organizations(id) ON DELETE SET NULL,
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    scopes JSONB NOT NULL DEFAULT '[]',
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_oauth_access_tokens_token_hash ON oauth_access_tokens(token_hash);
CREATE INDEX idx_oauth_access_tokens_user_id ON oauth_access_tokens(user_id);
CREATE INDEX idx_oauth_access_tokens_expires_at ON oauth_access_tokens(expires_at);
CREATE INDEX idx_oauth_access_tokens_app_id ON oauth_access_tokens(application_id);
