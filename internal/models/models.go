package models

import (
	"time"
)

// User представляет пользователя системы
type User struct {
	ID           int64     `json:"-" db:"id"`
	Login        string    `json:"login" db:"login"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Balance      float64   `json:"balance" db:"balance"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// Order представляет заказ в системе
type Order struct {
	ID          int64     `json:"-" db:"id"`
	Number      string    `json:"number" db:"number"`
	UserID      int64     `json:"-" db:"user_id"`
	Status      string    `json:"status" db:"status"`
	Accrual     float64   `json:"accrual,omitempty" db:"accrual"`
	UploadedAt  time.Time `json:"uploaded_at" db:"uploaded_at"`
	ProcessedAt time.Time `json:"processed_at,omitempty" db:"processed_at"`
}

// Withdrawal представляет операцию списания баллов
type Withdrawal struct {
	ID          int64     `json:"-" db:"id"`
	UserID      int64     `json:"-" db:"user_id"`
	Order       string    `json:"order" db:"order_number"`
	Sum         float64   `json:"sum" db:"sum"`
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
}

// OrderStatus определяет возможные статусы заказа
const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)
