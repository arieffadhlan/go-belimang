-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE users (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(32) NOT NULL,
    username VARCHAR(32) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    is_admin BOOLEAN
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS users;
-- +goose StatementEnd
