package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"inventory-platform-data-collector/internal/ebay/auth/service"
	"inventory-platform-data-collector/internal/ebay/integration/models/inventory"
	"inventory-platform-data-collector/internal/ebay/integration/models/pagination"
	"io"
	"net/http"
	"strconv"
)

type InventoryItemsClient struct {
	*BaseClient
	authService *service.Service
}

func CreateInventoryItemsClient(config *Config, authService *service.Service) *InventoryItemsClient {
	return &InventoryItemsClient{
		BaseClient:  NewBaseClient(config),
		authService: authService,
	}
}

func (c *InventoryItemsClient) GetInventoryItemsDetails(ctx context.Context, params pagination.PaginationParams, headers http.Header) (*inventory.InventoryResponse, error) {
	userID := headers.Get("X-User-ID")
	if userID == "" {
		return nil, fmt.Errorf("missing X-User-ID header")
	}

	token, err := c.authService.GetAccessToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	url := fmt.Sprintf("%ssell/inventory/v1/inventory_item", c.config.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.setHeaders(httpReq, token)
	q := httpReq.URL.Query()
	if params.Limit != nil {
		q.Add("limit", strconv.Itoa(*params.Limit))
	}
	if params.Offset != nil {
		q.Add("offset", strconv.Itoa(*params.Offset))
	}
	httpReq.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %v", resp.StatusCode, string(body))
	}
	var result inventory.InventoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}
