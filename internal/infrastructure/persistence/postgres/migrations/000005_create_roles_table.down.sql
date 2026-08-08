-- +migrate Down
DROP TRIGGER IF EXISTS update_roles_timestamp ON roles;
DROP TABLE IF EXISTS roles;