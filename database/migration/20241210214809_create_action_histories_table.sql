-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS action_histories (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    identifier_id VARCHAR(100) NOT NULL,
    identifier_name VARCHAR(100) NOT NULL,
    user_agent VARCHAR(255) NOT NULL DEFAULT '-',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP NULL,
    KEY identifier_id(identifier_id),
    KEY identifier_name(identifier_name),
    KEY created_at(created_at),
    KEY updated_at(updated_at),
    KEY is_deleted(is_deleted),
    KEY deleted_at(deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS action_histories;
-- +goose StatementEnd