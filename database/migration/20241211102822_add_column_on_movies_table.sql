-- +goose Up
-- +goose StatementBegin
ALTER TABLE movies
    ADD COLUMN genre_ids TEXT AFTER title,
    ADD COLUMN created_by VARCHAR(36) AFTER watch_url,
    ADD KEY genre_ids(genre_ids(255)),
    ADD KEY created_by(created_by),
    ADD CONSTRAINT fk_movies_created_by FOREIGN KEY (created_by) REFERENCES users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE movies
    DROP FOREIGN KEY fk_movies_created_by,
    DROP COLUMN genre_ids,
    DROP COLUMN created_by;
-- +goose StatementEnd
