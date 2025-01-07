package service

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrTokenNotFound   = errors.New("token not found")
	ErrTokenExpired    = errors.New("token expired")
	ErrRedisConnection = errors.New("redis connection error")
)
