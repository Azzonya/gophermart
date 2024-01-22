CREATE TABLE IF NOT EXISTS user (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL DEFAULT '',
    balance INT NOT NULL DEFAULT 0
);

CREATE INDEX idx_user_login ON "user" (login);

CREATE TABLE IF NOT EXISTS order (
    code VARCHAR(255) NOT NULL PRIMARY KEY,
    uploaded_at timestamp default current_timestamp,
    status VARCHAR(255) NOT NULL DEFAULT '',
    user_id INT REFERENCES user(id),
    accrual INT NOT NULL DEFAULT 0
);

CREATE INDEX idx_order_user_id ON "order" (user_id);

CREATE TABLE IF NOT EXISTS withdrawal (
    id SERIAL PRIMARY KEY,
    order_code INT REFERENCES order(code),
    amount INT NOT NULL DEFAULT 0
);

CREATE INDEX idx_withdrawal_order_code ON withdrawal (order_code);

