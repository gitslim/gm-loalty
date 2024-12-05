package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gitslim/gophermart/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	CreateOrderQuery      string
	GetOrderByNumberQuery string
	GetUserOrdersQuery    string
	UpdateOrderStatus     string
	GetOrdersByStatuses   string
)

func init() {
	queries := map[string]*string{
		"create_order.sql":           &CreateOrderQuery,
		"get_order_by_number.sql":    &GetOrderByNumberQuery,
		"get_user_orders.sql":        &GetUserOrdersQuery,
		"update_order_status.sql":    &UpdateOrderStatus,
		"get_orders_by_statuses.sql": &GetOrdersByStatuses,
	}

	loadQueries(queries)
}

// PgOrderStorage представляет хранилище заказов в PostgreSQL
type PgOrderStorage struct {
	db *sqlx.DB
}

// NewPgOrderStorage создает новый экземпляр хранилища PostgreSQL
func NewPgOrderStorage(db *sqlx.DB) *PgOrderStorage {
	return &PgOrderStorage{
		db: db,
	}
}

// CreateOrder создает новый заказ
func (s *PgOrderStorage) CreateOrder(ctx context.Context, order *models.Order) error {
	_, err := s.db.ExecContext(ctx, CreateOrderQuery,
		order.Number,
		order.UserID,
		order.Status,
		order.Accrual,
		order.UploadedAt,
		order.ProcessedAt,
	)

	return err
}

// GetOrderByNumber возвращает заказ по номеру
func (s *PgOrderStorage) GetOrderByNumber(ctx context.Context, number string) (*models.Order, error) {
	var order models.Order
	err := s.db.GetContext(ctx, &order, GetOrderByNumberQuery, number)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &order, err
}

// GetUserOrders возвращает все заказы пользователя
func (s *PgOrderStorage) GetUserOrders(ctx context.Context, userID int64) ([]*models.Order, error) {
	var orders []*models.Order
	err := s.db.SelectContext(ctx, &orders, GetUserOrdersQuery, userID)
	return orders, err
}

// UpdateOrderStatus обновляет статус заказа
func (s *PgOrderStorage) UpdateOrderStatus(ctx context.Context, orderID int64, status string, accrual float64) error {
	_, err := s.db.ExecContext(ctx, UpdateOrderStatus, orderID, status, accrual)
	return err
}

// GetOrdersByStatuses возвращает заказы с указанными статусами
func (s *PgOrderStorage) GetOrdersByStatuses(ctx context.Context, statuses []string) ([]*models.Order, error) {
	var orders []*models.Order
	err := s.db.SelectContext(ctx, &orders, GetOrdersByStatuses, statuses)
	return orders, err
}
