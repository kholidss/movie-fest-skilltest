-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS movie_genre (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    movie_id VARCHAR(36) NOT NULL,
    genre_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP NULL,
    KEY movie_id(movie_id),
    KEY genre_id(genre_id),
    KEY created_at(created_at),
    KEY updated_at(updated_at),
    KEY is_deleted(is_deleted),
    KEY deleted_at(deleted_at),
    CONSTRAINT fk_movie_genre_movie_id FOREIGN KEY (movie_id) REFERENCES movies(id),
    CONSTRAINT fk_movie_genre_genre_id FOREIGN KEY (genre_id) REFERENCES genres(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS movie_genre;
-- +goose StatementEnd