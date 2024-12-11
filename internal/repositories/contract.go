package repositories

import (
	"context"
	"database/sql"

	"github.com/kholidss/movie-fest-skilltest/internal/entity"
)

type UserRepository interface {
	Store(ctx context.Context, payload any, opts ...Option) error
	FindOne(ctx context.Context, param any, selectColumn []string) (*entity.User, error)
	Finds(ctx context.Context, param any, selectColumns []string) ([]entity.User, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}
