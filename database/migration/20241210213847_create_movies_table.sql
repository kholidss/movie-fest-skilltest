-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS movies (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    minutes_duration INT NOT NULL DEFAULT 0,
    view_number BIGINT NOT NULL DEFAULT 0,
    artist TEXT NOT NULL,
    watch_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP NULL,
    KEY title(title),
    KEY description(description(100)),
    KEY minutes_duration(minutes_duration),
    KEY view_number(view_number),
    KEY artist(artist(100)),
    KEY created_at(created_at),
    KEY updated_at(updated_at),
    KEY is_deleted(is_deleted),
    KEY deleted_at(deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS movies;
-- +goose StatementEnd
