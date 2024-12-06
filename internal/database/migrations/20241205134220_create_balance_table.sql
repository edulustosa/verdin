-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS balances (
    "id" UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
    "user_id" UUID NOT NULL,
    "current" DECIMAL NOT NULL DEFAULT 0,
    "income" DECIMAL NOT NULL DEFAULT 0,
    "expenses" DECIMAL NOT NULL DEFAULT 0,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS balances;
-- +goose StatementEnd