package clients

import (
	"inventory-platform-data-collector/internal/ebay/auth/service"
	"inventory-platform-data-collector/internal/ebay/integration/clients"
)

type ClientFactory struct {
	authService *service.Service
	userID      string
	config      *clients.Config
}

func NewClientFactory(authService *service.Service, userID string, environment string) *ClientFactory {
	baseURL := "https://api.sandbox.ebay.com"
	if environment == "production" {
		baseURL = "https://api.ebay.com"
	}

	return &ClientFactory{
		authService: authService,
		userID:      userID,
		config: &clients.Config{
			BaseURL:     baseURL,
			Environment: environment,
		},
	}
}

func (f *ClientFactory) NewFindingClient() *clients.FindingClient {
	return clients.NewFindingClient(f.config, f.authService, f.userID)
}

func (f *ClientFactory) NewMerchandisingClient() *clients.MerchandisingClient {
	return clients.NewMerchandisingClient(f.config, f.authService, f.userID)
}

func (f *ClientFactory) NewProductClient() *clients.ProductClient {
	return clients.NewProductClient(f.config, f.authService, f.userID)
}

func (f *ClientFactory) NewTradingClient() *clients.TradingClient {
	return clients.NewTradingClient(f.config, f.authService, f.userID)
}
