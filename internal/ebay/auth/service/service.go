package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"inventory-platform-data-collector/internal/ebay/auth/models"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	configs map[string]*models.OAuthConfig
	redis   *redis.Client
	mu      sync.RWMutex
}

func NewService(ctx context.Context, redisURL string) (*Service, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &Service{
		configs: make(map[string]*models.OAuthConfig),
		redis:   client,
	}, nil
}
func (s *Service) RegisterConfig(userID string, config *models.OAuthConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.configs[userID] = config
}

func (s *Service) GetConfig(userID string) (*models.OAuthConfig, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	config, exists := s.configs[userID]
	return config, exists
}

func (s *Service) GetAuthURL(userID string, state string) (string, error) {
	config, exists := s.GetConfig(userID)
	if !exists {
		return "", fmt.Errorf("configuration not found for user: %s", userID)
	}

	baseURL := "https://auth.sandbox.ebay.com/oauth2/authorize"
	if config.Environment == "production" {
		baseURL = "https://auth.ebay.com/oauth2/authorize"
	}

	params := url.Values{}
	params.Add("client_id", config.ClientID)
	params.Add("response_type", "code")
	params.Add("redirect_uri", config.RedirectURI)
	params.Add("scope", strings.Join(config.Scopes, " "))
	params.Add("state", state)

	return fmt.Sprintf("%s?%s", baseURL, params.Encode()), nil
}

func (s *Service) ExchangeCodeForToken(ctx context.Context, userID, code string) error {
	config, exists := s.GetConfig(userID)
	if !exists {
		return fmt.Errorf("configuration not found for user: %s", userID)
	}

	baseURL := "https://api.sandbox.ebay.com/identity/v1/oauth2/token"
	if config.Environment == "production" {
		baseURL = "https://api.ebay.com/identity/v1/oauth2/token"
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", config.RedirectURI)

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(config.ClientID + ":" + config.ClientSecret))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var tokenResp models.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return err
	}
	token := &models.TokenCache{
		AccessToken:  tokenResp.AccessToken,
		TokenType:    tokenResp.TokenType,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
	}
	tokenJSON, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return s.redis.Set(ctx, fmt.Sprintf("ebay:tokens:%s", userID), tokenJSON, 0).Err()
}

func (s *Service) GetAccessToken(ctx context.Context, userID string) (string, error) {
	tokenJSON, err := s.redis.Get(ctx, fmt.Sprintf("ebay:tokens:%s", userID)).Result()
	if err != nil {
		return "", fmt.Errorf("no valid token found for user: %s", err)
	}
	var token models.TokenCache
	if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
		return "", err
	}

	if time.Now().After(token.ExpiresAt) {
		return "", fmt.Errorf("token expired for user: %s", userID)
	}

	return token.AccessToken, nil
}
