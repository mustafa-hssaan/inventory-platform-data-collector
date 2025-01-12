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
}

func NewProductClient(config *Config, authService *service.Service) *ProductClient {
	return &ProductClient{
		BaseClient:  NewBaseClient(config),
		authService: authService,
	}
}

func (c *ProductClient) GetProduct(ctx context.Context, req models.GetProductRequest, headers http.Header) (*models.ProductResponse, error) {
	userID := headers.Get("X-User-ID")
	if userID == "" {
		return nil, fmt.Errorf("missing X-User-ID header")
	}

	token, err := c.authService.GetAccessToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}
	url := fmt.Sprintf("%s/buy/browse/v1/item/%s", c.config.BaseURL, req.ProductID)

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
