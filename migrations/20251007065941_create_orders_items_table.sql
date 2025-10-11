-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    estimate_id UUID NOT NULL,
    merchant_id UUID NOT NULL,
    merchant_item_id UUID NOT NULL,
    quantity INT NOT NULL
);

CREATE INDEX idx_orders_items_merchant_id ON orders_items (merchant_id);
CREATE INDEX idx_orders_items_merchant_item_id ON orders_items(merchant_item_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_orders_items_merchant_id;
DROP INDEX IF EXISTS idx_orders_items_merchant_item_id;

DROP TABLE IF EXISTS orders_items;
-- +goose StatementEnd
