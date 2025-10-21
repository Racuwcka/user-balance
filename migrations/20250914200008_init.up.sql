CREATE TABLE balances
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT    NOT NULL UNIQUE,
    balance    INT       NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_balances_user_id ON balances (user_id);

CREATE TABLE transactions
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT    NOT NULL,
    amount     INT       NOT NULL,
    type       TEXT      NOT NULL CHECK (type IN ('deposit', 'reserve', 'revenue')), -- deposit, withdraw, transfer
    service_id BIGINT,                                                                -- связь с услугой
    order_id   BIGINT,                                                                -- связь с заказом
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_transactions_service_id ON transactions (service_id);
CREATE INDEX idx_transactions_order_id ON transactions (order_id);

CREATE TABLE reserves
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT    NOT NULL,
    amount     INT       NOT NULL,
    status     TEXT      NOT NULL CHECK (status IN ('success', 'wait', 'rejected')),
    service_id BIGINT, -- связь с услугой
    order_id   BIGINT, -- связь с заказом
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_reserves_user_id ON reserves (user_id);
CREATE INDEX idx_reserves_status ON reserves (status);