package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"inventory-platform-data-collector/internal/ebay/auth/handlers"
	"inventory-platform-data-collector/internal/ebay/auth/service"
	client_handlers "inventory-platform-data-collector/internal/ebay/integration/handlers"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisURL := os.Getenv("REDIS_URL")
	userID := os.Getenv("EBAY_USER_ID")
	environment := os.Getenv("EBAY_ENVIRONMENT")

	authService, err := service.NewService(ctx, redisURL)
	if err != nil {
		log.Fatal(err)
	}

	handler := handlers.NewHandler(authService)

	http.HandleFunc("/auth/config", handler.RegisterConfig)
	http.HandleFunc("/auth/authorize", handler.HandleAuth)
	http.HandleFunc("/auth/callback", handler.HandleCallback)

	clientHandler := client_handlers.NewClientHandler(authService, userID, environment)
	http.HandleFunc("/auth/token", handler.GetToken)

	http.HandleFunc("/api/search", clientHandler.FindingHandler)
	http.HandleFunc("/api/deals", clientHandler.MerchandisingHandler)
	http.HandleFunc("/api/product", clientHandler.ProductHandler)
	http.HandleFunc("/api/item", clientHandler.TradingHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
