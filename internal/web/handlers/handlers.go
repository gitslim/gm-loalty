package handlers

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/gophermart/internal/logging"
	"github.com/gitslim/gophermart/internal/service"
	"github.com/gitslim/gophermart/internal/web/dto"
	"github.com/gitslim/gophermart/internal/web/middleware"
)

const (
	userIDKey = "userID"
)

// Handler содержит обработчики HTTP запросов
type Handler struct {
	userService    service.UserService
	orderService   service.OrderService
	balanceService service.BalanceService
	log            logging.Logger
	auth           *middleware.AuthMiddleware
}

// NewHandler создает новый экземпляр Handler
func NewHandler(log logging.Logger, userService service.UserService, orderService service.OrderService, balanceService service.BalanceService, auth *middleware.AuthMiddleware) *Handler {
	return &Handler{
		userService:    userService,
		orderService:   orderService,
		balanceService: balanceService,
		log:            log,
		auth:           auth,
	}
}

// getUserID возвращает ID пользователя из контекста
func getUserID(c *gin.Context) (int64, error) {
	userIDRaw, exists := c.Get(userIDKey)
	if !exists {
		return 0, errors.New("user ID not found in context")
	}

	userID, ok := userIDRaw.(int64)
	if !ok {
		return 0, errors.New("invalid user ID type in context")
	}

	return userID, nil
}

// Register обрабатывает регистрацию пользователя
func (h *Handler) Register(c *gin.Context) {
	var req dto.UserRequest
	h.log.Debugf("Register request: %+v", req)
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Debugf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.userService.Register(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "login already taken"})
			return
		}
		h.log.Debugf("Failed to register user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	token, err := h.auth.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.auth.SetAuthCookie(c, token)
	c.Set(userIDKey, user.ID)

	c.Status(http.StatusOK)
}

// Login обрабатывает аутентификацию пользователя
func (h *Handler) Login(c *gin.Context) {
	var req dto.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.userService.Login(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	h.log.Debugf("User %d logged in", user.ID)

	token, err := h.auth.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.auth.SetAuthCookie(c, token)
	c.Set(userIDKey, user.ID)

	c.Status(http.StatusOK)
}

// UploadOrder обрабатывает загрузку номера заказа
func (h *Handler) UploadOrder(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	orderNumber := string(body)
	if !isValidLuhn(orderNumber) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid order number"})
		return
	}

	err = h.orderService.UploadOrder(c.Request.Context(), userID, orderNumber)
	if err != nil {
		switch err.Error() {
		case "order already uploaded by this user":
			c.Status(http.StatusOK)
		case "order already uploaded by another user":
			c.JSON(http.StatusConflict, gin.H{"error": "order registered by another user"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.Status(http.StatusAccepted)
}

// GetOrders возвращает список заказов пользователя
func (h *Handler) GetOrders(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orders, err := h.orderService.GetUserOrders(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if len(orders) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetBalance возвращает текущий баланс пользователя
func (h *Handler) GetBalance(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	balance, err := h.balanceService.GetBalance(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	withdrawals, err := h.balanceService.GetWithdrawals(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var withdrawn float64
	for _, w := range withdrawals {
		withdrawn += w.Sum
	}

	response := dto.BalanceResponse{
		Current:   balance,
		Withdrawn: withdrawn,
	}

	c.JSON(http.StatusOK, response)
}

// Withdraw обрабатывает запрос на списание средств
func (h *Handler) Withdraw(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if !isValidLuhn(req.Order) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid order number"})
		return
	}

	err = h.balanceService.Withdraw(c.Request.Context(), userID, req.Order, req.Sum)
	if err != nil {
		if err.Error() == "insufficient funds" {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": "insufficient funds"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusOK)
}

// GetWithdrawals возвращает историю списаний средств
func (h *Handler) GetWithdrawals(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	withdrawals, err := h.balanceService.GetWithdrawals(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if len(withdrawals) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, withdrawals)
}

// isValidLuhn проверяет номер заказа по алгоритму Луна
func isValidLuhn(number string) bool {
	digits := make([]int, 0, len(number))
	for _, r := range number {
		if d, err := strconv.Atoi(string(r)); err == nil {
			digits = append(digits, d)
		} else {
			return false
		}
	}

	if len(digits) == 0 {
		return false
	}

	checksum := 0
	for i := len(digits) - 2; i >= 0; i -= 2 {
		d := digits[i] * 2
		if d > 9 {
			d -= 9
		}
		digits[i] = d
	}

	for _, d := range digits {
		checksum += d
	}

	return checksum%10 == 0
}
