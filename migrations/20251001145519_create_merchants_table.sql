-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TYPE merchant_categories_enum AS ENUM (
    'SmallRestaurant',
    'LargeRestaurant',
    'BoothKiosk',
    'MediumRestaurant',
    'ConvenienceStore',
    'MerchandiseRestaurant'
);

CREATE TABLE IF NOT EXISTS merchants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    imageurl TEXT NOT NULL,
    category merchant_categories_enum NOT NULL,
    location GEOGRAPHY('POINT', 4326) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_merchants_category ON merchants (category);
CREATE INDEX idx_merchants_location ON merchants USING GIST(location);
CREATE INDEX idx_merchants_created_at ON merchants (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_merchants_location;
DROP INDEX IF EXISTS idx_merchants_category;
DROP INDEX IF EXISTS idx_merchants_created_at;

DROP TABLE IF EXISTS merchants;

DROP TYPE IF EXISTS merchant_categories_enum;

DROP EXTENSION IF EXISTS postgis CASCADE;
-- +goose StatementEnd