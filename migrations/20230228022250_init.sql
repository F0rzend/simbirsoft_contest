-- +goose Up
-- +goose StatementBegin
CREATE TABLE account(
    id integer PRIMARY KEY,
    first_name varchar,
    last_name varchar,
    email varchar
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE account;
-- +goose StatementEnd
