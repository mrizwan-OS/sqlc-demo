-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN status VARCHAR(20) DEFAULT 'active';
ALTER TABLE users ADD COLUMN last_login TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN status;
ALTER TABLE users DROP COLUMN last_login;
-- +goose StatementEnd
