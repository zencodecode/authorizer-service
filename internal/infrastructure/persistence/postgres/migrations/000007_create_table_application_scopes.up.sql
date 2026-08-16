CREATE TABLE application_scopes (
    id UUID PRIMARY KEY,
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    scope VARCHAR(255) NOT NULL,
    description TEXT,
    UNIQUE(application_id, scope)
);

CREATE INDEX idx_app_scopes_app_id ON application_scopes(application_id);
