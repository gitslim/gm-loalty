package balance

import (
	"context"
	"fmt"
	"time"

	"github.com/gitslim/gophermart/internal/models"
	"github.com/gitslim/gophermart/internal/service"
	"github.com/gitslim/gophermart/internal/storage"
)

// BalanceServiceImpl реализует интерфейс service.BalanceService
type BalanceServiceImpl struct {
	storage storage.Storage
}

// NewBalanceService создает новый экземпляр сервиса баланса
func NewBalanceService(storage storage.Storage) service.BalanceService {
	return &BalanceServiceImpl{
		storage: storage,
	}
}

// GetBalance возвращает текущий баланс пользователя
func (s *BalanceServiceImpl) GetBalance(ctx context.Context, userID int64) (float64, error) {
	user, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return 0, fmt.Errorf("user not found")
	}

	return user.Balance, nil
}

// Withdraw списывает средства с баланса пользователя
func (s *BalanceServiceImpl) Withdraw(ctx context.Context, userID int64, orderNumber string, amount float64) error {
	// Проверяем баланс пользователя
	user, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	if user.Balance < amount {
		return fmt.Errorf("insufficient funds")
	}

	// Создаем запись о списании
	withdrawal := &models.Withdrawal{
		UserID:      userID,
		Order:       orderNumber,
		Sum:         amount,
		ProcessedAt: time.Now(),
	}

	if err := s.storage.CreateWithdrawal(ctx, withdrawal); err != nil {
		return fmt.Errorf("failed to create withdrawal: %w", err)
	}

	// Обновляем баланс пользователя
	if err := s.storage.UpdateBalance(ctx, userID, -amount); err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return nil
}

// GetWithdrawals возвращает историю списаний пользователя
func (s *BalanceServiceImpl) GetWithdrawals(ctx context.Context, userID int64) ([]*models.Withdrawal, error) {
	return s.storage.GetUserWithdrawals(ctx, userID)
}
