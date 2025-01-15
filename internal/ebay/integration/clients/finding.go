package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"inventory-platform-data-collector/internal/ebay/auth/service"
	"inventory-platform-data-collector/internal/ebay/integration/models/browseitem"
	"io"
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

func (c *FindingClient) FindItemDetailsByID(ctx context.Context, itemRequest browseitem.GetBrowseItemRequest, params browseitem.ItemBrowseParams, headers http.Header) (*browseitem.BrowseItemResponse, error) {
	userID := headers.Get("X-User-ID")
	if userID == "" {
		return nil, fmt.Errorf("missing X-User-ID header")
	}

	token, err := c.authService.GetAccessToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}
	itemID := fmt.Sprintf("v1|%s|0", itemRequest.ItemID)

	url := fmt.Sprintf("%s/buy/browse/v1/item/%s", c.config.BaseURL, itemID)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.setHeaders(httpReq, token)
	q := httpReq.URL.Query()
	if params.FieldGroups != "" {
		q.Add("fieldGroups", params.FieldGroups)
	}
	if params.QuantityForShippingEstimate != "" {
		q.Add("quantity_for_shipping_estimate", params.QuantityForShippingEstimate)
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
	var result browseitem.BrowseItemResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}
