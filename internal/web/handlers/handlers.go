package handlers

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gitslim/gophermart/internal/errs"
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

// handleError обрабатывает ошибки
func handleError(c *gin.Context, err error) {
	var e *errs.AppError
	if errors.As(err, &e) {
		c.JSON(e.Type.HTTPStatus, gin.H{"error": e.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func bindDTO(c *gin.Context, dto interface{}) error {
	if err := c.ShouldBindJSON(&dto); err != nil {
		return errs.NewAppError(errs.ErrBadRequest, "invalid request")
	}
	return nil
}

// getUserID возвращает ID пользователя из контекста
func getUserID(c *gin.Context) (int64, error) {
	err := errs.NewAppError(errs.ErrUnauthorized, "user not found")

	userIDRaw, exists := c.Get(userIDKey)
	if !exists {
		return 0, err
	}

	userID, ok := userIDRaw.(int64)
	if !ok {
		return 0, err
	}

	return userID, nil
}

// validateOrderLuhn проверяет номер заказа по алгоритму Луна
func validateOrderLuhn(number string) error {
	err := errs.NewAppError(errs.ErrUnprocessableEntity, "invalid order number")
	digits := make([]int, 0, len(number))
	for _, r := range number {
		if d, err := strconv.Atoi(string(r)); err == nil {
			digits = append(digits, d)
		} else {
			return err
		}
	}

	if len(digits) == 0 {
		return err
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

	if checksum%10 != 0 {
		return err
	}

	return nil
}

// Register обрабатывает регистрацию пользователя
func (h *Handler) Register(c *gin.Context) {
	var req dto.UserRequest
	err := bindDTO(c, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	user, err := h.userService.Register(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		handleError(c, err)
		return
	}

	token, err := h.auth.GenerateToken(user.ID)
	if err != nil {
		handleError(c, err)
		return
	}

	h.auth.SetAuthCookie(c, token)
	c.Set(userIDKey, user.ID)

	c.Status(http.StatusOK)
}

// Login обрабатывает аутентификацию пользователя
func (h *Handler) Login(c *gin.Context) {
	var req dto.UserRequest
	err := bindDTO(c, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	user, err := h.userService.Login(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		handleError(c, err)
		return
	}

	token, err := h.auth.GenerateToken(user.ID)
	if err != nil {
		handleError(c, err)
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
		handleError(c, err)
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		handleError(c, err)
		return
	}

	orderNumber := string(body)
	if err := validateOrderLuhn(orderNumber); err != nil {
		handleError(c, err)
		return
	}

	err = h.orderService.UploadOrder(c.Request.Context(), userID, orderNumber)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusAccepted)
}

// GetOrders возвращает список заказов пользователя
func (h *Handler) GetOrders(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		handleError(c, err)
		return
	}

	orders, err := h.orderService.GetUserOrders(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
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
		handleError(c, err)
		return
	}

	balance, err := h.balanceService.GetBalance(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	withdrawals, err := h.balanceService.GetWithdrawals(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
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
		handleError(c, err)
		return
	}

	var req dto.WithdrawRequest
	err = bindDTO(c, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	if err := validateOrderLuhn(req.Order); err != nil {
		handleError(c, err)
		return
	}

	err = h.balanceService.Withdraw(c.Request.Context(), userID, req.Order, req.Sum)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// GetWithdrawals возвращает историю списаний средств
func (h *Handler) GetWithdrawals(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		handleError(c, err)
		return
	}

	withdrawals, err := h.balanceService.GetWithdrawals(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	if len(withdrawals) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, withdrawals)
}
