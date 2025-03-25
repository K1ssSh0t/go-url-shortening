package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"url-shortener/internal/handlers"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"
)

func main() {
	// MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get MongoDB connection string from environment or use default
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	// Create MongoDB client
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Create repository, service, and handlers
	repo := repository.NewURLRepository(client, "urlshortener")
	urlService := service.NewURLService(repo)
	urlHandler := handlers.NewURLHandler(urlService)

	// Setup Gin router
	router := gin.Default()

	// URL shortener routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/shorten", urlHandler.CreateShortURL)
		v1.GET("/shorten/:shortCode", urlHandler.RedirectURL)
		v1.PUT("/shorten/:shortCode", urlHandler.UpdateShortURL)
		v1.DELETE("/shorten/:shortCode", urlHandler.DeleteShortURL)
		v1.GET("/shorten/:shortCode/stats", urlHandler.GetURLStats)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on :%s", port)
	log.Fatal(router.Run(":" + port))
}