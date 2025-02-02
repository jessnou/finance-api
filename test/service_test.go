package service

import (
	"context"
	"errors"
	"finance-api/internal/model"
	"finance-api/internal/repository"
	"finance-api/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_Deposit(t *testing.T) {
	mockRepo := new(repository.MockRepository)
	service := service.NewService(mockRepo)

	userID := int64(1)
	amount := 100.0
	ctx := context.Background()

	// Ожидаем вызов метода Deposit и имитируем успешный результат
	mockRepo.On("Deposit", ctx, userID, amount).Return(nil)

	err := service.Deposit(ctx, userID, amount)
	assert.NoError(t, err) // Проверяем, что ошибки нет

	// Проверяем, что метод был вызван с правильными аргументами
	mockRepo.AssertCalled(t, "Deposit", ctx, userID, amount)
}

func TestService_Transfer(t *testing.T) {
	mockRepo := new(repository.MockRepository)
	service := service.NewService(mockRepo)

	fromID := int64(1)
	toID := int64(2)
	amount := 50.0
	ctx := context.Background()

	mockRepo.On("Transfer", ctx, fromID, toID, amount).Return(nil)

	err := service.Transfer(ctx, fromID, toID, amount)
	assert.NoError(t, err)

	mockRepo.AssertCalled(t, "Transfer", ctx, fromID, toID, amount)
}

func TestService_GetTransactions(t *testing.T) {
	mockRepo := new(repository.MockRepository)
	service := service.NewService(mockRepo)

	userID := int64(1)
	ctx := context.Background()

	expectedTransactions := []model.Transaction{
		{ID: 1, UserID: userID, Amount: 100.0, Type: "deposit"},
		{ID: 2, UserID: userID, Amount: -50.0, Type: "transfer"},
	}

	mockRepo.On("GetTransactions", ctx, userID).Return(expectedTransactions, nil)

	transactions, err := service.GetTransactions(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedTransactions, transactions)

	mockRepo.AssertCalled(t, "GetTransactions", ctx, userID)
}

func TestService_Transfer_Error(t *testing.T) {
	mockRepo := new(repository.MockRepository)
	service := service.NewService(mockRepo)

	fromID := int64(1)
	toID := int64(2)
	amount := 50.0
	ctx := context.Background()

	mockRepo.On("Transfer", ctx, fromID, toID, amount).Return(errors.New("insufficient funds"))

	err := service.Transfer(ctx, fromID, toID, amount)
	assert.Error(t, err)
	assert.Equal(t, "insufficient funds", err.Error())

	mockRepo.AssertCalled(t, "Transfer", ctx, fromID, toID, amount)
}
