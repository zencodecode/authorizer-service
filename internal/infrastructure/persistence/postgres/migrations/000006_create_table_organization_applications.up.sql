CREATE TABLE organization_applications (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    activated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(organization_id, application_id)
);

CREATE INDEX idx_org_apps_org_id ON organization_applications(organization_id);
CREATE INDEX idx_org_apps_app_id ON organization_applications(application_id);
