package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
)

/*
=======================================================================================================
Create Mock CDN Method
=======================================================================================================
*/

type MockCDN struct {
	mock.Mock
}

func (m *MockCDN) Put(ctx context.Context, name string, contents []byte) (any, error) {
	args := m.Called(ctx, name, contents)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockCDN) Delete(ctx context.Context, identifier string) error {
	args := m.Called(ctx, identifier)
	return args.Error(0)
}

func (m *MockCDN) Get(ctx context.Context, identifier string) ([]byte, error) {
	args := m.Called(ctx, identifier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}
