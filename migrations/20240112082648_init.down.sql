-- +goose Down
DROP TABLE bonus_transactions CASCADE;

DROP TABLE orders CASCADE;

DROP TABLE users CASCADE;
