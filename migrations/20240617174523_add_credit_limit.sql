-- +goose Up
-- +goose StatementBegin
ALTER TABLE accounts
ADD credit_limit NUMERIC(10, 2) NOT NULL DEFAULT 0;

ALTER TABLE accounts
ADD CONSTRAINT check_credit_limit_non_negative
CHECK (credit_limit >= 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE accounts
DROP CONSTRAINT check_credit_limit_non_negative;

ALTER TABLE accounts
DROP COLUMN credit_limit;
-- +goose StatementEnd
