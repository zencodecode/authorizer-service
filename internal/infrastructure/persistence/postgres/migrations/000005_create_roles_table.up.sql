-- +migrate Up
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY,
    application_id UUID,
    code TEXT NOT NULL,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    scope public.scope,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (application_id, code),

    CONSTRAINT fk_roles_application
        FOREIGN KEY (application_id) REFERENCES applications (id) ON DELETE CASCADE
);

CREATE TRIGGER update_roles_timestamp
BEFORE UPDATE ON roles
FOR EACH ROW
EXECUTE PROCEDURE update_timestamp();
