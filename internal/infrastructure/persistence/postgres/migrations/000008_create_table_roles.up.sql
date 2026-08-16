CREATE TABLE roles (
    id UUID PRIMARY KEY,
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT,
    is_system BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(organization_id, application_id, slug)
);

-- Partial unique index for roles without organization (NULL org_id)
CREATE UNIQUE INDEX idx_roles_null_org_app_slug
    ON roles(application_id, slug)
    WHERE organization_id IS NULL;

CREATE INDEX idx_roles_org_id ON roles(organization_id);
CREATE INDEX idx_roles_app_id ON roles(application_id);
CREATE INDEX idx_roles_slug ON roles(slug);

CREATE TRIGGER set_roles_updated_at
    BEFORE UPDATE ON roles
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();
