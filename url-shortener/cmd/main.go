package main

import (
	"context"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"url-shortener/internal/handlers"

	_ "url-shortener/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title URL Shortener API
// @version 1.0
// @description This is a URL shortening service API
// @host localhost:8080
// @BasePath /
func main() {
	// Leer la variable de entorno MONGO_URI
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // Valor por defecto si la variable no está definida
	}

	// Configurar las opciones del cliente
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Conectar a MongoDB
	_, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

    r := gin.Default()

    // Swagger documentation endpoint
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    r.POST("/shorten", handlers.CreateShortURL)
    r.GET("/shorten/:shortCode", handlers.GetOriginalURL)
    r.PUT("/shorten/:shortCode", handlers.UpdateShortURL)
    r.DELETE("/shorten/:shortCode", handlers.DeleteShortURL)
    r.GET("/shorten/:shortCode/stats", handlers.GetURLStats)

    r.Run(":8080")
}
