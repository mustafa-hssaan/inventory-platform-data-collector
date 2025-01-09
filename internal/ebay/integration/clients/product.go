package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"inventory-platform-data-collector/internal/ebay/auth/service"
	"inventory-platform-data-collector/internal/ebay/integration/models"
	"net/http"
)

type ProductClient struct {
	*BaseClient
	authService *service.Service
	userID      string
}

func NewProductClient(config *Config, authService *service.Service, userID string) *ProductClient {
	return &ProductClient{
		BaseClient:  NewBaseClient(config),
		authService: authService,
		userID:      userID,
	}
}

func (c *ProductClient) GetProduct(ctx context.Context, req models.GetProductRequest) (*models.ProductResponse, error) {
	token, err := c.authService.GetAccessToken(ctx, c.userID)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}
	url := fmt.Sprintf("%s/buy/browse/v1/item/%s", c.config.BaseURL)

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.setHeaders(request, token)

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	var result models.ProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}
