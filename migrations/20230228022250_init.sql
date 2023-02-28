-- +goose Up
-- +goose StatementBegin
CREATE TABLE account(
                        id INTEGER PRIMARY KEY,
                        first_name VARCHAR,
                        last_name VARCHAR,
                        email VARCHAR
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE account;
-- +goose StatementEnd
