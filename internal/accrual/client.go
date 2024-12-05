package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gitslim/gophermart/internal/conf"
	"github.com/gitslim/gophermart/internal/errs"
)

// Response представляет ответ от системы расчета начислений
type Response struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

// Client представляет клиент для взаимодействия с системой расчета начислений
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient создает новый экземпляр клиента системы начислений
func NewClient(config *conf.Config) *Client {
	return &Client{
		baseURL: config.AccrualSystemAddress,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetOrderAccrual получает информацию о начислении баллов за заказ
func (c *Client) GetOrderAccrual(ctx context.Context, orderNumber string) (response *Response, StatusCode int, err error) {
	url := fmt.Sprintf("%s/api/orders/%s", c.baseURL, orderNumber)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, errs.NewAppError(errs.ErrInternal, "failed to create request")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, errs.NewAppError(errs.ErrInternal, "failed to send request")
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var response Response
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, 0, errs.NewAppError(errs.ErrInternal, "failed to decode response")
		}
		return &response, resp.StatusCode, nil
	default:
		return nil, resp.StatusCode, nil
	}
}
