CREATE TABLE oauth_authorization_codes (
    id UUID PRIMARY KEY,
    code_hash VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID REFERENCES organizations(id) ON DELETE SET NULL,
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    redirect_uri VARCHAR(2048) NOT NULL,
    code_challenge VARCHAR(255) NOT NULL,
    code_challenge_method VARCHAR(10) NOT NULL DEFAULT 'S256',
    scopes JSONB NOT NULL DEFAULT '[]',
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_oauth_auth_codes_code_hash ON oauth_authorization_codes(code_hash);
CREATE INDEX idx_oauth_auth_codes_user_id ON oauth_authorization_codes(user_id);
CREATE INDEX idx_oauth_auth_codes_expires_at ON oauth_authorization_codes(expires_at);
