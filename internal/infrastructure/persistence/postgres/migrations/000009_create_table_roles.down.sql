DROP TRIGGER IF EXISTS set_roles_updated_at ON roles;
DROP INDEX IF EXISTS idx_roles_null_org_app_slug;
DROP TABLE IF EXISTS roles;
