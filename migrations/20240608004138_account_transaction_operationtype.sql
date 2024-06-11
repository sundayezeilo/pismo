-- +goose Up
-- +goose StatementBegin
CREATE TYPE op_desc_enum AS ENUM ('purchase', 'installment purchase', 'withdrawal', 'payment');
CREATE TYPE op_type_enum AS ENUM ('debit', 'credit');

CREATE TABLE IF NOT EXISTS operation_types (
    id SERIAL PRIMARY KEY,
    "description" op_desc_enum UNIQUE NOT NULL,
    op_type op_type_enum NOT NULL,
    active_support BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO operation_types ("description", op_type)
VALUES
    ('purchase', 'debit'),
    ('installment purchase', 'debit'),
    ('withdrawal', 'debit'),
    ('payment', 'credit');

CREATE TABLE IF NOT EXISTS accounts (
    id BIGSERIAL PRIMARY KEY,
    document_number VARCHAR(11) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    account_id INT REFERENCES accounts(id),
    operation_type_id INT REFERENCES operation_types(id),
    amount DECIMAL(10, 2) NOT NULL,
    event_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transactions_account_id ON transactions (account_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS operation_types;

DROP TYPE IF EXISTS op_desc_enum;
DROP TYPE IF EXISTS op_type_enum;
-- +goose StatementEnd
