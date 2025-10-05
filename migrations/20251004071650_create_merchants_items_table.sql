-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE purchase_categories_enum AS ENUM (
    'Beverage',
    'Food',
    'Snack',
    'Condiments',
    'Additions'
);

CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID NOT NULL,
    name VARCHAR(30) NOT NULL,
    price INT NOT NULL,
    image_url TEXT NOT NULL,
    category purchase_categories_enum NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_items_category ON items (category);
CREATE INDEX idx_items_merchant_id ON items(merchant_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_items_category;
DROP INDEX IF EXISTS idx_items_merchant_id;

DROP TABLE IF EXISTS items;

DROP TYPE IF EXISTS purchase_categories_enum;

DROP EXTENSION IF EXISTS pgcrypto CASCADE;
-- +goose StatementEnd
