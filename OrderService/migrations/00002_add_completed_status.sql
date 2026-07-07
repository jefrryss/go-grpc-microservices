-- +goose Up

ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_status_check;
ALTER TABLE orders ADD CONSTRAINT orders_status_check
    CHECK (status IN ('PENDING_PAYMENT', 'PAID', 'CANCELLED', 'COMPLETED'));

-- +goose Down

ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_status_check;
ALTER TABLE orders ADD CONSTRAINT orders_status_check
    CHECK (status IN ('PENDING_PAYMENT', 'PAID', 'CANCELLED'));
