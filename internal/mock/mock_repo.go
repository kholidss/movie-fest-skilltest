package mock

import (
	"context"
	"database/sql"

	"github.com/kholidss/movie-fest-skilltest/internal/entity"
	"github.com/kholidss/movie-fest-skilltest/internal/repositories"
	"github.com/stretchr/testify/mock"
)

/*
=======================================================================================================
Create Mock Repo User
=======================================================================================================
*/
type MockRepoUser struct {
	mock.Mock
}

func (m *MockRepoUser) Store(ctx context.Context, payload any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, opts)
	return args.Error(0)
}

func (m *MockRepoUser) Update(ctx context.Context, payload any, where any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, opts)
	return args.Error(0)
}

func (m *MockRepoUser) FindOne(ctx context.Context, param any, selectColumn []string) (*entity.User, error) {
	args := m.Called(ctx, param, selectColumn)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockRepoUser) Finds(ctx context.Context, param any, selectColumns []string) ([]entity.User, error) {
	args := m.Called(ctx, param, selectColumns)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]entity.User), args.Error(1)
}

func (m *MockRepoUser) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	args := m.Called(ctx, opts)
	return nil, args.Error(1)
}

/*
=======================================================================================================
Create Mock Repo Movie
=======================================================================================================
*/
type MockRepoMovie struct {
	mock.Mock
}

func (m *MockRepoMovie) Store(ctx context.Context, payload any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, opts)
	return args.Error(0)
}

func (m *MockRepoMovie) Update(ctx context.Context, payload any, where any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, where, opts)
	return args.Error(0)
}

func (m *MockRepoMovie) FindOne(ctx context.Context, param any, selectColumns []string) (*entity.Movie, error) {
	args := m.Called(ctx, param, selectColumns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Movie), args.Error(1)
}

func (m *MockRepoMovie) FindOneWithForUpdate(ctx context.Context, param any, opts ...repositories.Option) (*entity.Movie, error) {
	args := m.Called(ctx, param, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Movie), args.Error(1)
}

func (m *MockRepoMovie) Finds(ctx context.Context, param any, selectColumns []string) ([]entity.Movie, error) {
	args := m.Called(ctx, param, selectColumns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Movie), args.Error(1)
}

func (m *MockRepoMovie) List(ctx context.Context, meta entity.MetaPagination, selectColumns []string) ([]entity.Movie, int, error) {
	args := m.Called(ctx, meta, selectColumns)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]entity.Movie), args.Int(1), args.Error(2)
}

func (m *MockRepoMovie) ListMostView(ctx context.Context, meta entity.MetaPagination, selectColumns []string) ([]entity.Movie, int, error) {
	args := m.Called(ctx, meta, selectColumns)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]entity.Movie), args.Int(1), args.Error(2)
}

func (m *MockRepoMovie) ListWithLike(ctx context.Context, meta entity.MetaPagination, param any, selectColumns []string) ([]entity.Movie, int, error) {
	args := m.Called(ctx, meta, param, selectColumns)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]entity.Movie), args.Int(1), args.Error(2)
}

func (m *MockRepoMovie) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	args := m.Called(ctx, opts)
	return nil, args.Error(1)
}

/*
=======================================================================================================
Create Mock Repo MovieGenre
=======================================================================================================
*/
type MockRepoMovieGenre struct {
	mock.Mock
}

func (m *MockRepoMovieGenre) Store(ctx context.Context, payload any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, opts)
	return args.Error(0)
}

func (m *MockRepoMovieGenre) Update(ctx context.Context, payload any, where any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, where, opts)
	return args.Error(0)
}

func (m *MockRepoMovieGenre) FindOne(ctx context.Context, param any, selectColumns []string) (*entity.MovieGenre, error) {
	args := m.Called(ctx, param, selectColumns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MovieGenre), args.Error(1)
}

func (m *MockRepoMovieGenre) Finds(ctx context.Context, param any, selectColumns []string) ([]entity.MovieGenre, error) {
	args := m.Called(ctx, param, selectColumns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.MovieGenre), args.Error(1)
}

func (m *MockRepoMovieGenre) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(*sql.Tx), args.Error(1)
}

/*
=======================================================================================================
Create Mock Repo MovieVote
=======================================================================================================
*/
type MockRepoMovieVote struct {
	mock.Mock
}

func (m *MockRepoMovieVote) Store(ctx context.Context, payload any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, opts)
	return args.Error(0)
}

func (m *MockRepoMovieVote) Update(ctx context.Context, payload any, where any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, where, opts)
	return args.Error(0)
}

func (m *MockRepoMovieVote) FindOne(ctx context.Context, param any, selectColumns []string) (*entity.MovieVote, error) {
	args := m.Called(ctx, param, selectColumns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MovieVote), args.Error(1)
}

func (m *MockRepoMovieVote) FindOneWithForUpdate(ctx context.Context, param any, opts ...repositories.Option) (*entity.MovieVote, error) {
	args := m.Called(ctx, param, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MovieVote), args.Error(1)
}

func (m *MockRepoMovieVote) Finds(ctx context.Context, param any, selectColumns []string) ([]entity.MovieVote, error) {
	args := m.Called(ctx, param, selectColumns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.MovieVote), args.Error(1)
}

func (m *MockRepoMovieVote) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(*sql.Tx), args.Error(1)
}

/*
=======================================================================================================
Create Mock Repo Genre
=======================================================================================================
*/
type MockRepoGenre struct {
	mock.Mock
}

func (m *MockRepoGenre) Store(ctx context.Context, payload any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, opts)
	return args.Error(0)
}

func (m *MockRepoGenre) Update(ctx context.Context, payload any, where any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, where, opts)
	return args.Error(0)
}

func (m *MockRepoGenre) FindOne(ctx context.Context, param any, selectColumns []string) (*entity.Genre, error) {
	args := m.Called(ctx, param, selectColumns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Genre), args.Error(1)
}

func (m *MockRepoGenre) FindOneWithForUpdate(ctx context.Context, param any, opts ...repositories.Option) (*entity.Genre, error) {
	args := m.Called(ctx, param, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Genre), args.Error(1)
}

func (m *MockRepoGenre) Finds(ctx context.Context, param any, selectColumns []string) ([]entity.Genre, error) {
	args := m.Called(ctx, param, selectColumns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Genre), args.Error(1)
}

func (m *MockRepoGenre) ListMostView(ctx context.Context, meta entity.MetaPagination, selectColumns []string) ([]entity.Genre, int, error) {
	args := m.Called(ctx, meta, selectColumns)
	return args.Get(0).([]entity.Genre), args.Int(1), args.Error(2)
}

func (m *MockRepoGenre) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(*sql.Tx), args.Error(1)
}

/*
=======================================================================================================
Create Mock Repo Bucket
=======================================================================================================
*/

type MockRepoBucket struct {
	mock.Mock
}

func (m *MockRepoBucket) Store(ctx context.Context, payload any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, opts)
	return args.Error(0)
}

func (m *MockRepoBucket) Update(ctx context.Context, payload any, where any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, where, opts)
	return args.Error(0)
}

func (m *MockRepoBucket) FindOne(ctx context.Context, param any, selectColumns []string) (*entity.Bucket, error) {
	args := m.Called(ctx, param, selectColumns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Bucket), args.Error(1)
}

func (m *MockRepoBucket) Finds(ctx context.Context, param any, selectColumns []string) ([]entity.Bucket, error) {
	args := m.Called(ctx, param, selectColumns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Bucket), args.Error(1)
}

func (m *MockRepoBucket) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(*sql.Tx), args.Error(1)
}

/*
=======================================================================================================
Create Mock Repo Action History
=======================================================================================================
*/

type MockRepoActionHistory struct {
	mock.Mock
}

func (m *MockRepoActionHistory) Store(ctx context.Context, payload any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, opts)
	return args.Error(0)
}

func (m *MockRepoActionHistory) Update(ctx context.Context, payload any, where any, opts ...repositories.Option) error {
	args := m.Called(ctx, payload, where, opts)
	return args.Error(0)
}

func (m *MockRepoActionHistory) FindOne(ctx context.Context, param any, selectColumns []string) (*entity.ActionHistory, error) {
	args := m.Called(ctx, param, selectColumns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ActionHistory), args.Error(1)
}

func (m *MockRepoActionHistory) Finds(ctx context.Context, param any, selectColumns []string) ([]entity.ActionHistory, error) {
	args := m.Called(ctx, param, selectColumns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.ActionHistory), args.Error(1)
}

func (m *MockRepoActionHistory) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(*sql.Tx), args.Error(1)
}
