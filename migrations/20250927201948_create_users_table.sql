-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;

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
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS pgcrypto CASCADE;
-- +goose StatementEnd
