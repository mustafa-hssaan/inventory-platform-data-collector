package integration

import (
	"context"
	"fmt"

	authService "inventory-platform-data-collector/internal/ebay/auth/service"
)

type Config struct {
	UserID      string
	Environment string
	authService *authService.Service
}
type TokenProvider struct {
	config      *Config
	authService *authService.Service
}

func NewTokenProvider(config *Config, authService *authService.Service) *TokenProvider {
	return &TokenProvider{
		config:      config,
		authService: authService,
	}
}
func (tp *TokenProvider) GetToken(ctx context.Context) (string, error) {
	token, err := tp.authService.GetAccessToken(ctx, tp.config.UserID)
	if err != nil {
		return "", fmt.Errorf("getting access token: %w", err)
	}
	return token, nil
}
