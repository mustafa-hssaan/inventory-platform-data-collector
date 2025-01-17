package clients

import (
	"inventory-platform-data-collector/internal/ebay/auth/service"
	"inventory-platform-data-collector/internal/ebay/integration/clients"
)

type ClientFactory struct {
	authService *service.Service
	config      *clients.Config
}

func NewClientFactory(authService *service.Service, environment string) *ClientFactory {
	baseURL := "https://api.sandbox.ebay.com"
	if environment == "production" {
		baseURL = "https://api.ebay.com"
	}

	return &ClientFactory{
		authService: authService,
		config: &clients.Config{
			BaseURL:     baseURL,
			Environment: environment,
		},
	}
}

func (f *ClientFactory) NewItemBrowseClient() *clients.ItemsBrowseClient {
	return clients.CreateItemBrowseClient(f.config, f.authService)
}
func (f *ClientFactory) NewAnalyticsClient() *clients.AnalyticsClient {
	return clients.CreateAnalyticsClient(f.config, f.authService)
}
