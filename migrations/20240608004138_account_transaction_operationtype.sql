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
    balance NUMERIC(10, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_accounts_document_number ON accounts (document_number);

ALTER TABLE accounts
ADD CONSTRAINT check_balance_non_negative
CHECK (balance >= 0);

CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT REFERENCES accounts(id),
    operation_type_id INT REFERENCES operation_types(id),
    amount DECIMAL(10, 2) NOT NULL,
    event_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    balance_before NUMERIC(10, 2) NOT NULL,
    balance_after NUMERIC(10, 2) NOT NULL
);

ALTER TABLE transactions 
ADD CONSTRAINT check_balance_before_after_non_negative
CHECK (balance_after >= 0 AND balance_before >= 0);

CREATE INDEX idx_transactions_account_id ON transactions (account_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_transactions_account_id;

ALTER TABLE transactions
DROP CONSTRAINT check_balance_before_after_non_negative;

DROP TABLE transactions;

ALTER TABLE accounts
DROP CONSTRAINT check_balance_non_negative;

DROP INDEX idx_accounts_document_number;

DROP TABLE accounts;

DROP TABLE operation_types;

DROP TYPE op_desc_enum;
DROP TYPE op_type_enum;
-- +goose StatementEnd
