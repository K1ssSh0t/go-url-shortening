package handlers

import (
    "context"
    "net/http"
    "time"
    "url-shortener/internal/repository"
    "url-shortener/internal/models"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "math/rand"
)

func generateShortCode() string {
    letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
    code := make([]rune, 6)
    for i := range code {
        code[i] = letters[rand.Intn(len(letters))]
    }
    return string(code)
}

func CreateShortURL(c *gin.Context) {
    var input struct {
        URL string `json:"url" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    shortCode := generateShortCode()
    newURL := models.ShortURL{
        ID:          primitive.NewObjectID().Hex(),
        URL:         input.URL,
        ShortCode:   shortCode,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
        AccessCount: 0,
    }

    _, err := repository.DB.InsertOne(context.TODO(), newURL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save URL"})
        return
    }

    c.JSON(http.StatusCreated, newURL)
}

func GetOriginalURL(c *gin.Context) {
    shortCode := c.Param("shortCode")
    var result models.ShortURL

    err := repository.DB.FindOne(context.TODO(), bson.M{"short_code": shortCode}).Decode(&result)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
        return
    }

    repository.DB.UpdateOne(context.TODO(), bson.M{"short_code": shortCode}, bson.M{"$inc": bson.M{"access_count": 1}})
    c.JSON(http.StatusOK, result)
}

func UpdateShortURL(c *gin.Context) {
    shortCode := c.Param("shortCode")
    var input struct {
        URL string `json:"url" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    update := bson.M{"$set": bson.M{"url": input.URL, "updated_at": time.Now()}}
    res, err := repository.DB.UpdateOne(context.TODO(), bson.M{"short_code": shortCode}, update)
    if err != nil || res.ModifiedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "URL updated"})
}

func DeleteShortURL(c *gin.Context) {
    shortCode := c.Param("shortCode")
    res, err := repository.DB.DeleteOne(context.TODO(), bson.M{"short_code": shortCode})
    if err != nil || res.DeletedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
        return
    }

    c.JSON(http.StatusNoContent, nil)
}

func GetURLStats(c *gin.Context) {
    shortCode := c.Param("shortCode")
    var result models.ShortURL

    err := repository.DB.FindOne(context.TODO(), bson.M{"short_code": shortCode}).Decode(&result)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
        return
    }

    c.JSON(http.StatusOK, result)
}
