CREATE TABLE balances
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT         NOT NULL UNIQUE,
    amount     NUMERIC(12, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP      NOT NULL DEFAULT now()
);

CREATE INDEX idx_balances_user_id ON balances (user_id);

CREATE TABLE transactions
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT         NOT NULL,
    amount     NUMERIC(12, 2) NOT NULL,
    type       VARCHAR(20)    NOT NULL, -- deposit, withdraw, transfer
    service_id BIGINT,                  -- связь с услугой
    order_id   BIGINT,                  -- связь с заказом
    created_at TIMESTAMP      NOT NULL DEFAULT now()
);

CREATE INDEX idx_transactions_service_id ON transactions (service_id);
CREATE INDEX idx_transactions_order_id ON transactions (order_id);