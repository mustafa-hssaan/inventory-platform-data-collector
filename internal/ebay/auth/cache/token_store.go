package cache

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenCache struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	ExpiresAt    time.Time
}
type TokenStore struct {
	redisClient *redis.Client
	gcm         cipher.AEAD
	namespace   string
	mu          sync.RWMutex
}

func NewTokenStore(ctx context.Context, redisURL, encryptionKey, namespace string) (*TokenStore, error) {
	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	redisDb := redis.NewClient(opt)

	return &TokenStore{
		redisClient: redisDb,
		gcm:         gcm,
		namespace:   namespace,
	}, nil
}
func (tokenStore *TokenStore) encrypt(data []byte) ([]byte, error) {
	nonce := make([]byte, tokenStore.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return tokenStore.gcm.Seal(nonce, nonce, data, nil), nil
}

func (tokenStore *TokenStore) decrypt(data []byte) ([]byte, error) {
	nonceSize := tokenStore.gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("data too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return tokenStore.gcm.Open(nil, nonce, ciphertext, nil)
}

func (tokenStore *TokenStore) Set(ctx context.Context, userID string, token *TokenCache) error {
	tokenStore.mu.Lock()
	defer tokenStore.mu.Unlock()

	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	encrypted, err := tokenStore.encrypt(data)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s:%s", tokenStore.namespace, userID)
	encodedData := base64.StdEncoding.EncodeToString(encrypted)

	expiration := time.Until(token.ExpiresAt)
	return tokenStore.redisClient.Set(ctx, key, encodedData, expiration).Err()
}
func (tokenStore *TokenStore) Get(ctx context.Context, userID string) (*TokenCache, bool) {
	tokenStore.mu.RLock()
	defer tokenStore.mu.RUnlock()

	key := fmt.Sprintf("%s:%s", tokenStore.namespace, userID)

	encodedData, err := tokenStore.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, false
	}

	encrypted, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, false
	}

	data, err := tokenStore.decrypt(encrypted)
	if err != nil {
		return nil, false
	}

	var token TokenCache
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, false
	}

	if token.ExpiresAt.Before(time.Now()) {
		tokenStore.redisClient.Del(ctx, key)
		return nil, false
	}

	return &token, true
}
func (tokenStore *TokenStore) Delete(ctx context.Context, userID string) error {
	key := fmt.Sprintf("%s:%s", tokenStore.namespace, userID)
	return tokenStore.redisClient.Del(ctx, key).Err()
}
