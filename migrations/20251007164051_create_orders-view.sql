-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE VIEW order_history_view AS
SELECT
    od.id AS order_id,
    es.user_id AS user_id,
    mc.id AS merchant_id,
    mc.name AS merchant_name,
    mc.category AS merchant_category,
    mc.imageurl AS merchant_imageurl,
    ST_Y(mc.location::geometry) AS merchant_lat,
    ST_X(mc.location::geometry) AS merchant_lon,
    mc.created_at AS merchant_created_at,
    it.id AS item_id,
    it.name AS item_name,
    it.category AS item_category,
    it.imageurl AS item_imageurl,
    it.price AS item_price,
    oi.quantity AS quantity,
    it.created_at AS item_created_at
FROM orders od
JOIN estimates es ON es.id = od.estimate_id
JOIN orders_items oi ON oi.estimate_id = od.estimate_id
JOIN merchants mc ON mc.id = oi.merchant_id
JOIN items it ON it.id = oi.merchant_item_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS order_history_view;
-- +goose StatementEnd
