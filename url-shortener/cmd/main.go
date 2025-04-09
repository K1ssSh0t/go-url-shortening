package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "url-shortener/docs"
	"url-shortener/internal/handlers"
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
		mongoURI = "mongodb://localhost:27017"
	}

	// Configurar las opciones del cliente
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Conectar a MongoDB
	_, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	// Swagger documentation endpoint
	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Post("/shorten", handlers.CreateShortURL)
	app.Get("/shorten/:shortCode", handlers.GetOriginalURL)
	app.Put("/shorten/:shortCode", handlers.UpdateShortURL)
	app.Delete("/shorten/:shortCode", handlers.DeleteShortURL)
	app.Get("/shorten/:shortCode/stats", handlers.GetURLStats)

	log.Fatal(app.Listen(":8080"))
}
