-- +goose Up
-- +goose StatementBegin
ALTER TABLE movies
    ADD COLUMN vote_number BIGINT DEFAULT 0 NOT NULL AFTER view_number,
    ADD KEY vote_number(vote_number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE movies
DROP COLUMN vote_number;
-- +goose StatementEnd
