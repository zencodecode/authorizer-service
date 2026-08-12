CREATE TABLE applications (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    client_id VARCHAR(255) NOT NULL UNIQUE,
    client_secret_hash VARCHAR(255) NOT NULL,
    redirect_uris JSONB NOT NULL DEFAULT '[]',
    allowed_grant_types JSONB NOT NULL DEFAULT '["authorization_code", "refresh_token"]',
    requires_organization BOOLEAN NOT NULL DEFAULT TRUE,
    metadata JSONB,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_applications_client_id ON applications(client_id);
CREATE INDEX idx_applications_slug ON applications(slug);
CREATE INDEX idx_applications_is_active ON applications(is_active);

CREATE TRIGGER set_applications_updated_at
    BEFORE UPDATE ON applications
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();
