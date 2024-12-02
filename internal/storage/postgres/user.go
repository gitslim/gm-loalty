package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gitslim/gophermart/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	CreateUserQuery     string
	GetUserByLoginQuery string
	GetUserByIDQuery    string
	UpdateBalanceQuery  string
)

func init() {
	queries := map[string]*string{
		"create_user.sql":       &CreateUserQuery,
		"get_user_by_login.sql": &GetUserByLoginQuery,
		"get_user_by_id.sql":    &GetUserByIDQuery,
		"update_balance.sql":    &UpdateBalanceQuery,
	}
	loadQueries(queries)
}

// PgUserStorage представляет хранилище пользователей PostgreSQL
type PgUserStorage struct {
	db *sqlx.DB
}

// NewPgUserStorage создает новый экземпляр хранилища PostgreSQL
func NewPgUserStorage(db *sqlx.DB) *PgUserStorage {
	return &PgUserStorage{
		db: db,
	}
}

// CreateUser создает нового пользователя
func (s *PgUserStorage) CreateUser(ctx context.Context, user *models.User) error {
	err := s.db.GetContext(ctx, &user.ID, CreateUserQuery,
		user.Login,
		user.PasswordHash,
		user.Balance,
		user.CreatedAt,
	)
	return err
}

// GetUserByLogin возвращает пользователя по логину
func (s *PgUserStorage) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	err := s.db.GetContext(ctx, &user, GetUserByLoginQuery, login)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &user, err
}

// GetUserByID возвращает пользователя по ID
func (s *PgUserStorage) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	err := s.db.GetContext(ctx, &user, GetUserByIDQuery, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &user, err
}

// UpdateBalance обновляет баланс пользователя
func (s *PgUserStorage) UpdateBalance(ctx context.Context, userID int64, delta float64) error {
	_, err := s.db.ExecContext(ctx, UpdateBalanceQuery, userID, delta)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return nil
}
