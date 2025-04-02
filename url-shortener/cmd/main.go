package main

import (
	"github.com/gin-gonic/gin"

	"url-shortener/internal/handlers"
	"url-shortener/internal/repository"

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
    repository.ConnectDB()

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
