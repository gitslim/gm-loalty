package dto

// UserRequest представляет запрос для регистрации/входа пользователя
type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// BalanceResponse представляет ответ с информацией о балансе
type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

// WithdrawRequest представляет запрос на списание средств
type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
