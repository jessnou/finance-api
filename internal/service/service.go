package service

import (
	"context"
	"finance-api/internal/model"
	"finance-api/internal/repository"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Deposit(ctx context.Context, userID int64, amount float64) error {
	return s.repo.Deposit(ctx, userID, amount)
}

func (s *Service) Transfer(ctx context.Context, fromID, toID int64, amount float64) error {
	return s.repo.Transfer(ctx, fromID, toID, amount)
}

func (s *Service) GetTransactions(ctx context.Context, userID int64) ([]model.Transaction, error) {
	return s.repo.GetTransactions(ctx, userID)
}
