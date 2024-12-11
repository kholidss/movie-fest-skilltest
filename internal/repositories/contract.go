package repositories

import (
	"context"
	"database/sql"

	"github.com/kholidss/movie-fest-skilltest/internal/entity"
)

type UserRepository interface {
	Store(ctx context.Context, payload any, opts ...Option) error
	Update(ctx context.Context, payload any, where any, opts ...Option) error
	FindOne(ctx context.Context, param any, selectColumns []string) (*entity.User, error)
	Finds(ctx context.Context, param any, selectColumns []string) ([]entity.User, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type MovieRepository interface {
	Store(ctx context.Context, payload any, opts ...Option) error
	Update(ctx context.Context, payload any, where any, opts ...Option) error
	FindOne(ctx context.Context, param any, selectColumns []string) (*entity.Movie, error)
	Finds(ctx context.Context, param any, selectColumns []string) ([]entity.Movie, error)
	ListMostView(ctx context.Context, meta entity.MetaPagination, selectColumns []string) ([]entity.Movie, int, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type MovieGenreRepository interface {
	Store(ctx context.Context, payload any, opts ...Option) error
	Update(ctx context.Context, payload any, where any, opts ...Option) error
	FindOne(ctx context.Context, param any, selectColumns []string) (*entity.MovieGenre, error)
	Finds(ctx context.Context, param any, selectColumns []string) ([]entity.MovieGenre, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type MovieVoteRepository interface {
	Store(ctx context.Context, payload any, opts ...Option) error
	Update(ctx context.Context, payload any, where any, opts ...Option) error
	FindOne(ctx context.Context, param any, selectColumns []string) (*entity.MovieVote, error)
	Finds(ctx context.Context, param any, selectColumns []string) ([]entity.MovieVote, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type GenreRepository interface {
	Store(ctx context.Context, payload any, opts ...Option) error
	Update(ctx context.Context, payload any, where any, opts ...Option) error
	FindOne(ctx context.Context, param any, selectColumns []string) (*entity.Genre, error)
	Finds(ctx context.Context, param any, selectColumns []string) ([]entity.Genre, error)
	ListMostView(ctx context.Context, meta entity.MetaPagination, selectColumns []string) ([]entity.Genre, int, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type BucketRepository interface {
	Store(ctx context.Context, payload any, opts ...Option) error
	Update(ctx context.Context, payload any, where any, opts ...Option) error
	FindOne(ctx context.Context, param any, selectColumns []string) (*entity.Bucket, error)
	Finds(ctx context.Context, param any, selectColumns []string) ([]entity.Bucket, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type ActionHistoryRepository interface {
	Store(ctx context.Context, payload any, opts ...Option) error
	Update(ctx context.Context, payload any, where any, opts ...Option) error
	FindOne(ctx context.Context, param any, selectColumns []string) (*entity.ActionHistory, error)
	Finds(ctx context.Context, param any, selectColumns []string) ([]entity.ActionHistory, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}
