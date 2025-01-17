package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"inventory-platform-data-collector/internal/ebay/auth/service"
	"inventory-platform-data-collector/internal/ebay/integration/models/analytics"
	"io"
	"net/http"
)

type AnalyticsClient struct {
	*BaseClient
	authService *service.Service
}

func CreateAnalyticsClient(config *Config, authService *service.Service) *AnalyticsClient {
	return &AnalyticsClient{
		BaseClient:  NewBaseClient(config),
		authService: authService,
	}
}
func (c *AnalyticsClient) GetTrafficReport(ctx context.Context, params analytics.TrafficReportParam, headers http.Header) (*analytics.TrafficReportResponse, error) {
	userID := headers.Get("X-User-ID")
	if userID == "" {
		return nil, fmt.Errorf("missing X-User-ID header")
	}

	token, err := c.authService.GetAccessToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}
	url := fmt.Sprintf("%s/sell/analytics/v1/traffic_report", c.config.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.setHeaders(httpReq, token)
	q := httpReq.URL.Query()
	if params.Dimension != "" {
		q.Add("dimension", params.Dimension)
	}
	if params.Filter != "" {
		q.Add("filter", params.Filter)
	}
	if params.Metric != "" {
		q.Add("metric", params.Metric)
	}
	if params.Sort != "" {
		q.Add("sort", params.Sort)
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
	var result analytics.TrafficReportResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}
