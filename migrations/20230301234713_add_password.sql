-- +goose Up
-- +goose StatementBegin
ALTER TABLE account
    ALTER COLUMN id SET DATA TYPE bigint,
    ADD COLUMN password varchar;

CREATE SEQUENCE account_id
    AS bigint
    INCREMENT BY 1
    MINVALUE 0 NO MAXVALUE
    START WITH 1
    NO CYCLE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE account
    ALTER COLUMN id SET DATA TYPE integer,
    DROP COLUMN password;
-- +goose StatementEnd
