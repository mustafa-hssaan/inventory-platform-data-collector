package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"inventory-platform-data-collector/internal/ebay/auth/service"
	"inventory-platform-data-collector/internal/ebay/integration/models"
	"net/http"
)

type FindingClient struct {
	*BaseClient
	authService *service.Service
}

func NewFindingClient(config *Config, authService *service.Service) *FindingClient {
	return &FindingClient{
		BaseClient:  NewBaseClient(config),
		authService: authService,
	}
}
func (c *FindingClient) Search(ctx context.Context, params models.SearchParams, headers http.Header) (*models.SearchResponse, error) {
	userID := headers.Get("X-User-ID")
	if userID == "" {
		return nil, fmt.Errorf("missing X-User-ID header")
	}

	token, err := c.authService.GetAccessToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	url := fmt.Sprintf("%s/buy/browse/v1/item_summary/search", c.config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.setHeaders(req, token)
	q := req.URL.Query()
	q.Add("q", params.Q)
	if params.CategoryID != "" {
		q.Add("category_ids", params.CategoryID)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	var result models.SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}
