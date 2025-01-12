package clients

import (
	"context"
	"encoding/xml"
	"fmt"
	"inventory-platform-data-collector/internal/ebay/auth/service"
	"inventory-platform-data-collector/internal/ebay/integration/models"
	"net/http"
)

type TradingClient struct {
	*BaseClient
	authService *service.Service
}

func NewTradingClient(config *Config, authService *service.Service) *TradingClient {
	return &TradingClient{
		BaseClient:  NewBaseClient(config),
		authService: authService,
	}
}

func (c *TradingClient) GetItem(ctx context.Context, req models.GetItemRequest, headers http.Header) (*models.GetItemResponse, error) {
	userID := headers.Get("X-User-ID")
	if userID == "" {
		return nil, fmt.Errorf("missing X-User-ID header")
	}

	token, err := c.authService.GetAccessToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}
	url := fmt.Sprintf("%s/ws/api.dll", c.config.BaseURL)

	request, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	c.setHeaders(request, token)

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	var result models.GetItemResponse
	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}
