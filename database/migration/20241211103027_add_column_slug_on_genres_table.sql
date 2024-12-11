-- +goose Up
-- +goose StatementBegin
ALTER TABLE genres
    ADD COLUMN slug VARCHAR(255) NOT NULL AFTER name,
    ADD KEY slug(slug);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE genres
    DROP COLUMN slug;
-- +goose StatementEnd
