CREATE TYPE order_status_enum AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');
CREATE TYPE transaction_type_enum AS ENUM ('+', '-');

CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     login VARCHAR(255) NOT NULL UNIQUE,
                                     password VARCHAR(255) NOT NULL,
                                     balance INT NOT NULL DEFAULT 0
);

CREATE INDEX idx_user_login ON users (login);

CREATE TABLE IF NOT EXISTS orders (
                                      code VARCHAR(255) NOT NULL PRIMARY KEY,
                                      uploaded_at timestamp with time zone default current_timestamp,
                                      status order_status_enum NOT NULL,
                                      user_id INT REFERENCES users(id)
);

CREATE INDEX idx_order_status ON orders (status);
CREATE INDEX idx_order_user_id ON orders (user_id);

CREATE TABLE bonus_transactions (
                                    order_code VARCHAR(255) NOT NULL,
                                    user_id INT REFERENCES users(id) NOT NULL,
                                    processed_at timestamp with time zone default current_timestamp,
                                    transaction_type transaction_type_enum NOT NULL,
                                    sum INT
);

CREATE INDEX idx_bonus_transactions_order_code ON bonus_transactions (order_code);
CREATE INDEX idx_bonus_transactions_user_id ON bonus_transactions (user_id);
CREATE INDEX idx_bonus_transactions_processed_at ON bonus_transactions (processed_at);
CREATE INDEX idx_bonus_transactions_transaction_type ON bonus_transactions (transaction_type);