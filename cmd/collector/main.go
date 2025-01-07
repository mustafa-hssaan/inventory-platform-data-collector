package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"inventory-platform-data-collector/internal/ebay/auth/handlers"
	"inventory-platform-data-collector/internal/ebay/auth/service"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisURL := os.Getenv("REDIS_URL")

	authService, err := service.NewService(ctx, redisURL)
	if err != nil {
		log.Fatal(err)
	}

	handler := handlers.NewHandler(authService)

	http.HandleFunc("/auth/config", handler.RegisterConfig)
	http.HandleFunc("/auth/authorize", handler.HandleAuth)
	http.HandleFunc("/auth/callback", handler.HandleCallback)
	http.HandleFunc("/auth/token", handler.GetToken)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
