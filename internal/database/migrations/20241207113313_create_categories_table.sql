-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories (
    "id" SERIAL PRIMARY KEY NOT NULL,
    "user_id" UUID,
    "name" VARCHAR(255) NOT NULL,
    "theme" VARCHAR(255) NOT NULL,
    "icon" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd