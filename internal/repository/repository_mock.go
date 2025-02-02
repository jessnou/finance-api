package repository

import (
	"context"
	"finance-api/internal/model"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Deposit(ctx context.Context, userID int64, amount float64) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

func (m *MockRepository) Transfer(ctx context.Context, fromID, toID int64, amount float64) error {
	args := m.Called(ctx, fromID, toID, amount)
	return args.Error(0)
}

func (m *MockRepository) GetTransactions(ctx context.Context, userID int64) ([]model.Transaction, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Transaction), args.Error(1)
}
