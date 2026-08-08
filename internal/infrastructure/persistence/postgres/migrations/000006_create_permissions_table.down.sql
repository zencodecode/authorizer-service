-- +migrate Down
DROP TRIGGER IF EXISTS update_permissions_timestamp ON permissions;
DROP FUNCTION IF EXISTS update_timestamp();
DROP TABLE IF EXISTS permissions;
