-- name:database-exist
SELECT 1 FROM pg_database WHERE datname='?';

-- name:database-create
CREATE DATABASE ?;