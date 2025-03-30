package main

import (


	"github.com/gin-gonic/gin"
	
	"url-shortener/internal/handlers"
	"url-shortener/internal/repository"

)



func main() {
    repository.ConnectDB()

    r := gin.Default()
    
    r.POST("/shorten", handlers.CreateShortURL)
    r.GET("/shorten/:shortCode", handlers.GetOriginalURL)
    r.PUT("/shorten/:shortCode", handlers.UpdateShortURL)
    r.DELETE("/shorten/:shortCode", handlers.DeleteShortURL)
    r.GET("/shorten/:shortCode/stats", handlers.GetURLStats)

    r.Run(":8080")
}
