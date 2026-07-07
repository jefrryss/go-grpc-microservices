-- +goose Up

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    part_ids UUID[] NOT NULL,
    total_price DOUBLE PRECISION NOT NULL
        CHECK (total_price >= 0),
    transaction_id UUID,
    payment_method TEXT NOT NULL DEFAULT 'UNKNOWN'
        CHECK (
            payment_method IN (
                'UNKNOWN',
                'CARD',
                'SBP',
                'CREDIT_CARD',
                'INVESTOR_MONEY'
            )
        ),
    status TEXT NOT NULL DEFAULT 'PENDING_PAYMENT'
        CHECK (
            status IN (
                'PENDING_PAYMENT',
                'PAID',
                'CANCELLED',
                'COMPLETED'
            )
        ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id
    ON orders (user_id);

CREATE INDEX IF NOT EXISTS idx_orders_status
    ON orders (status);

-- +goose Down

DROP TABLE IF EXISTS orders;
