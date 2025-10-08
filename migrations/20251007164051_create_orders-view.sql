-- +goose Up
-- +goose StatementBegin
ALTER TABLE items ALTER COLUMN merchant_id TYPE uuid USING merchant_id::uuid;
ALTER TABLE orders ALTER COLUMN estimate_id TYPE uuid USING estimate_id::uuid;
ALTER TABLE orders_items ALTER COLUMN estimate_id TYPE uuid USING estimate_id::uuid;
ALTER TABLE orders_items ALTER COLUMN merchant_id TYPE uuid USING merchant_id::uuid;
ALTER TABLE orders_items ALTER COLUMN merchant_item_id TYPE uuid USING merchant_item_id::uuid;

CREATE OR REPLACE VIEW order_history_view AS
SELECT
    o.id AS order_id,
    e.user_id AS user_id,
    m.id AS merchant_id,
    m.name AS merchant_name,
    m.category AS merchant_category,
    m.image_url AS merchant_image_url,
    ST_Y(m.location::geometry) AS merchant_lat,
    ST_X(m.location::geometry) AS merchant_long,
    m.created_at AS merchant_created_at,
    i.id AS item_id,
    i.name AS item_name,
    i.category AS item_category,
    i.price AS item_price,
    oi.quantity AS quantity,
    i.image_url AS item_image_url,
    i.created_at AS item_created_at
FROM orders o
JOIN estimates e ON e.id = o.estimate_id
JOIN orders_items oi ON oi.estimate_id = e.id
JOIN merchants m ON m.id = oi.merchant_id
JOIN items i ON i.id = oi.merchant_item_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS order_history_view;
ALTER TABLE items ALTER COLUMN merchant_id TYPE TEXT USING merchant_id::text;
ALTER TABLE orders ALTER COLUMN estimate_id TYPE TEXT USING estimate_id::text;
ALTER TABLE orders_items ALTER COLUMN estimate_id TYPE TEXT USING estimate_id::text;
ALTER TABLE orders_items ALTER COLUMN merchant_id TYPE TEXT USING merchant_id::text;
ALTER TABLE orders_items ALTER COLUMN merchant_item_id TYPE TEXT USING merchant_item_id::text;
-- +goose StatementEnd
