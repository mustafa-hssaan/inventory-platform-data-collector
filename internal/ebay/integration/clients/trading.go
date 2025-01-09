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
	userID      string
}

func NewTradingClient(config *Config, authService *service.Service, userID string) *TradingClient {
	return &TradingClient{
		BaseClient:  NewBaseClient(config),
		authService: authService,
		userID:      userID,
	}
}

func (c *TradingClient) GetItem(ctx context.Context, req models.GetItemRequest) (*models.GetItemResponse, error) {
	token, err := c.authService.GetAccessToken(ctx, c.userID)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}
	url := fmt.Sprintf("%s/ws/api.dll", c.config.BaseURL)

	request, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	c.setHeaders(request, token)

	// request.Header.Set("X-EBAY-API-CALL-NAME", "GetItem")
	// request.Header.Set("X-EBAY-API-APP-NAME", c.config.AppID)
	// request.Header.Set("X-EBAY-API-VERSION", "967")

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
