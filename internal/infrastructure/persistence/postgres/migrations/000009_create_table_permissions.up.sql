CREATE TABLE permissions (
    id UUID PRIMARY KEY,
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    slug VARCHAR(255) NOT NULL,
    resource VARCHAR(255) NOT NULL,
    action VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(application_id, slug)
);

CREATE INDEX idx_permissions_app_id ON permissions(application_id);
CREATE INDEX idx_permissions_slug ON permissions(slug);
