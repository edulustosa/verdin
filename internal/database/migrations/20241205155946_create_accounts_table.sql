-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS accounts (
    "id" UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
    "user_id" UUID NOT NULL,
    "title" VARCHAR(255) NOT NULL,
    "balance" DECIMAL(10, 2) NOT NULL DEFAULT 0,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF NOT EXISTS accounts;
-- +goose StatementEnd