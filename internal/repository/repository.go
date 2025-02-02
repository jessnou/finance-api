package repository

import (
	"context"
	"finance-api/internal/model"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Deposit(ctx context.Context, userID int64, amount float64) error
	Transfer(ctx context.Context, fromID, toID int64, amount float64) error
	GetTransactions(ctx context.Context, userID int64) ([]model.Transaction, error)
}

type postgresRepository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) Repository {
	return &postgresRepository{db: db}
}

// Deposit adds the specified amount to a user's balance and records the transaction.
func (r *postgresRepository) Deposit(ctx context.Context, userID int64, amount float64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "UPDATE users SET balance = balance + $1 WHERE id = $2", amount, userID)
	if err != nil {
		log.Printf("Error updating balance for user %d: %v", userID, err)
		return fmt.Errorf("failed to update user balance: %w", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO transactions (user_id, amount, type) VALUES ($1, $2, 'deposit')", userID, amount)
	if err != nil {
		log.Printf("Error inserting transaction for user %d: %v", userID, err)
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *postgresRepository) Transfer(ctx context.Context, fromID, toID int64, amount float64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var balance float64
	err = tx.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1", fromID).Scan(&balance)
	if err != nil {
		log.Printf("Error getting balance for user %d: %v", fromID, err)
		return fmt.Errorf("failed to get balance for user %d: %w", fromID, err)
	}

	if balance < amount {
		log.Printf("Insufficient funds for user %d", fromID)
		return fmt.Errorf("insufficient funds for user %d", fromID)
	}

	_, err = tx.Exec(ctx, `
		UPDATE users 
		SET balance = balance + CASE 
			WHEN id = $2 THEN $1 * -1
			WHEN id = $3 THEN $1
			ELSE 0 
		END
		WHERE id IN ($2, $3);
	`, amount, fromID, toID)
	if err != nil {
		log.Printf("Error updating balances: %v", err)
		return fmt.Errorf("failed to update balances: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO transactions (user_id, amount, type) 
		VALUES 
			($1, $2, 'transfer'), 
			($3, $4, 'transfer')
	`, fromID, -amount, toID, amount)

	if err := tx.Commit(ctx); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *postgresRepository) GetTransactions(ctx context.Context, userID int64) ([]model.Transaction, error) {
	rows, err := r.db.Query(ctx, "SELECT id, user_id, amount, type, created_at FROM transactions WHERE user_id = $1 ORDER BY created_at DESC LIMIT 10", userID)
	if err != nil {
		log.Printf("Error querying transactions for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.Amount, &t.Type, &t.CreatedAt); err != nil {
			log.Printf("Error scanning transaction: %v", err)
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over transactions: %v", err)
		return nil, fmt.Errorf("failed to iterate over transactions: %w", err)
	}

	return transactions, nil
}
