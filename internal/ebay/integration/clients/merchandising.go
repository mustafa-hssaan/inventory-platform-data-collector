package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"inventory-platform-data-collector/internal/ebay/auth/service"
	"inventory-platform-data-collector/internal/ebay/integration/models"
	"net/http"
)

type MerchandisingClient struct {
	*BaseClient
	authService *service.Service
}

func NewMerchandisingClient(config *Config, authService *service.Service) *MerchandisingClient {
	return &MerchandisingClient{
		BaseClient:  NewBaseClient(config),
		authService: authService,
	}
}

func (c *MerchandisingClient) GetDeals(ctx context.Context, req models.GetDealsRequest, headers http.Header) (*models.DealsResponse, error) {
	userID := headers.Get("X-User-ID")
	if userID == "" {
		return nil, fmt.Errorf("missing X-User-ID header")
	}

	token, err := c.authService.GetAccessToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	url := fmt.Sprintf("%s/buy/marketing/v1_beta/merchandised_product", c.config.BaseURL)
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.setHeaders(request, token)
	q := request.URL.Query()
	q.Add("category_ids", req.CategoryID)
	if req.Limit > 0 {
		q.Add("limit", fmt.Sprintf("%d", req.Limit))
	}
	request.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	var result models.DealsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}
