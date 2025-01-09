package clients

import (
	"net/http"
	"time"
)

type Config struct {
	BaseURL     string
	Environment string
}

type BaseClient struct {
	httpClient *http.Client
	config     *Config
}

func NewBaseClient(config *Config) *BaseClient {
	return &BaseClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		config: config,
	}
}

func (c *BaseClient) setHeaders(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-EBAY-C-MARKETPLACE-ID", "EBAY_US")
}
