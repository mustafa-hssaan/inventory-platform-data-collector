package models

import "time"

type OAuthConfig struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURI  string   `json:"redirect_uri"`
	Scopes       []string `json:"scopes"`
	Environment  string   `json:"environment"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type TokenCache struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type RegisterConfigRequest struct {
	UserID string      `json:"user_id"`
	Config OAuthConfig `json:"config"`
}
