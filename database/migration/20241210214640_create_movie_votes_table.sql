-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS movie_votes (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    movie_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP NULL,
    KEY user_id(user_id),
    KEY movie_id(movie_id),
    KEY created_at(created_at),
    KEY updated_at(updated_at),
    KEY is_deleted(is_deleted),
    KEY deleted_at(deleted_at),
    CONSTRAINT fk_movie_votes_user_id FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_movie_votes_movie_id FOREIGN KEY (movie_id) REFERENCES movies(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS movie_votes;
-- +goose StatementEnd