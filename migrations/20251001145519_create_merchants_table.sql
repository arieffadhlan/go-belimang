-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE IF NOT EXISTS merchant_category_enum AS ENUM (
    'SmallRestaurant',
    'MediumRestaurant',
    'LargeRestaurant',
    'MerchandiseRestaurant',
    'BoothKiosk',
    'ConvenienceStore'
);

CREATE TYPE IF NOT EXISTS product_category_enum AS ENUM (
    'Beverage',
    'Food',
    'Snack',
    'Condiments',
    'Additions'
);

CREATE TABLE IF NOT EXISTS merchants (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    name VARCHAR(30) NOT NULL CHECK (char_length(name) >= 2),
    merchant_category merchant_category_enum NOT NULL,
    image_url TEXT NOT NULL,
    location GEOGRAPHY(Point, 4326) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id TEXT NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    name VARCHAR(30) NOT NULL CHECK (char_length(name) >= 2),
    product_category product_category_enum NOT NULL,
    price NUMERIC(12,2) NOT NULL CHECK (price >= 1),
    image_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_merchants_name ON merchants (LOWER(name));
CREATE INDEX idx_merchants_category ON merchants (merchant_category);
CREATE INDEX idx_items_name ON items (LOWER(name));
CREATE INDEX idx_items_category ON items (product_category);
CREATE INDEX idx_merchants_location ON merchants USING GIST (location);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP INDEX IF EXISTS idx_merchants_location;
DROP INDEX IF EXISTS idx_items_category;
DROP INDEX IF EXISTS idx_items_name;
DROP INDEX IF EXISTS idx_merchants_category;
DROP INDEX IF EXISTS idx_merchants_name;

DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS merchants;

DROP TYPE IF EXISTS product_category_enum;
DROP TYPE IF EXISTS merchant_category_enum;

DROP EXTENSION IF EXISTS postgis;
DROP EXTENSION IF EXISTS pgcrypto;
-- +goose StatementEnd
