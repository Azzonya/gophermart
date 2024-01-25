CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     login VARCHAR(255) NOT NULL UNIQUE,
                                     password VARCHAR(255) NOT NULL DEFAULT ''
);

CREATE INDEX idx_user_login ON users (login);

CREATE TABLE IF NOT EXISTS orders (
                                      code VARCHAR(255) NOT NULL PRIMARY KEY,
                                      uploaded_at timestamp with time zone default current_timestamp,
                                      status VARCHAR(255) NOT NULL DEFAULT '',
                                      user_id INT REFERENCES users(id)
);

CREATE INDEX idx_order_status ON orders (status);
CREATE INDEX idx_order_user_id ON orders (user_id);

CREATE TABLE bonus_transactions (
                                    order_code VARCHAR(255),
                                    user_id INT REFERENCES users(id),
                                    processed_at timestamp with time zone default current_timestamp,
                                    transaction_type VARCHAR(255),
                                    sum INT
);

CREATE INDEX idx_bonus_transactions_order_code ON bonus_transactions (order_code);
CREATE INDEX idx_bonus_transactions_user_id ON bonus_transactions (user_id);
CREATE INDEX idx_bonus_transactions_processed_at ON bonus_transactions (processed_at);
CREATE INDEX idx_bonus_transactions_transaction_type ON bonus_transactions (transaction_type);