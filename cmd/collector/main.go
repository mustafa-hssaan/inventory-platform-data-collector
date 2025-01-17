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
	environment := os.Getenv("EBAY_ENVIRONMENT")

	authService, err := service.NewService(ctx, redisURL)
	if err != nil {
		log.Fatal(err)
	}

	handler := handlers.NewHandler(authService)
	http.HandleFunc("/auth/callback", handler.HandleCallback)
	http.HandleFunc("/auth/register-and-start", handler.RegisterConfigAndStart)
	http.HandleFunc("/auth/complete", handler.CompleteAuth)

	clientHandler := client_handlers.NewClientHandler(authService, environment)

	http.HandleFunc("/api/browseItemsDetails", clientHandler.FindingHandler)
	http.HandleFunc("/api/trafficReport", clientHandler.GetTrafficReport)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
