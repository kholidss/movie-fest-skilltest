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
