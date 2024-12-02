package postgres

import (
	"context"

	"github.com/gitslim/gophermart/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	CreateWithdrawalQuery string
	GetUserWithdrawals    string
)

func init() {
	queries := map[string]*string{
		"create_withdrawal.sql":    &CreateWithdrawalQuery,
		"get_user_withdrawals.sql": &GetUserWithdrawals,
	}

	loadQueries(queries)
}

// PgWithdrawalStorage представляет хранилище операций списания
type PgWithdrawalStorage struct {
	db *sqlx.DB
}

// NewPgWithdrawalStorage создает новый экземпляр хранилища PostgreSQL
func NewPgWithdrawalStorage(db *sqlx.DB) *PgWithdrawalStorage {
	return &PgWithdrawalStorage{
		db: db,
	}
}

// CreateWithdrawal создает новую операцию списания
func (s *PgWithdrawalStorage) CreateWithdrawal(ctx context.Context, withdrawal *models.Withdrawal) error {
	_, err := s.db.ExecContext(ctx, CreateWithdrawalQuery,
		withdrawal.UserID,
		withdrawal.Order,
		withdrawal.Sum,
		withdrawal.ProcessedAt,
	)
	return err
}

// GetUserWithdrawals возвращает все операции списания пользователя
func (s *PgWithdrawalStorage) GetUserWithdrawals(ctx context.Context, userID int64) ([]*models.Withdrawal, error) {
	var withdrawals []*models.Withdrawal
	err := s.db.SelectContext(ctx, &withdrawals, GetUserWithdrawals, userID)
	return withdrawals, err
}
