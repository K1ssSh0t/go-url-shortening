package handlers

import (
	"context"
	"net/http"
	"time"
	"url-shortener/internal/models"
	"url-shortener/internal/repository"

	"math/rand"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func generateShortCode() string {
    letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
    code := make([]rune, 6)
    for i := range code {
        code[i] = letters[rand.Intn(len(letters))]
    }
    return string(code)
}

// @Summary      Crear una URL corta
// @Description  Crea una nueva URL corta a partir de una URL larga
// @Tags         urls
// @Accept       json
// @Produce      json
// @Param        url  body      models.ShortURL  true  "URL a acortar"
// @Success      201  {object}  models.ShortURL
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /shorten [post]
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

// @Summary      Obtener URL original
// @Description  Obtiene la URL original a partir de un código corto
// @Tags         urls
// @Accept       json
// @Produce      json
// @Param        shortCode  path      string  true  "Código corto de la URL"
// @Success      200        {object}  models.ShortURL
// @Failure      404        {object}  map[string]string
// @Router       /shorten/{shortCode} [get]
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

// @Summary      Actualizar URL corta
// @Description  Actualiza la URL original asociada a un código corto
// @Tags         urls
// @Accept       json
// @Produce      json
// @Param        shortCode  path      string  true  "Código corto de la URL"
// @Param        url        body      models.ShortURL  true  "Nueva URL"
// @Success      200        {object}  map[string]string
// @Failure      400        {object}  map[string]string
// @Failure      404        {object}  map[string]string
// @Router       /shorten/{shortCode} [put]
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

// @Summary      Eliminar URL corta
// @Description  Elimina una URL corta y su asociación
// @Tags         urls
// @Accept       json
// @Produce      json
// @Param        shortCode  path      string  true  "Código corto de la URL"
// @Success      204        {object}  nil
// @Failure      404        {object}  map[string]string
// @Router       /shorten/{shortCode} [delete]
func DeleteShortURL(c *gin.Context) {
    shortCode := c.Param("shortCode")
    res, err := repository.DB.DeleteOne(context.TODO(), bson.M{"short_code": shortCode})
    if err != nil || res.DeletedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
        return
    }

    c.JSON(http.StatusNoContent, nil)
}

// @Summary      Obtener estadísticas
// @Description  Obtiene las estadísticas de uso de una URL corta
// @Tags         urls
// @Accept       json
// @Produce      json
// @Param        shortCode  path      string  true  "Código corto de la URL"
// @Success      200        {object}  models.ShortURL
// @Failure      404        {object}  map[string]string
// @Router       /shorten/{shortCode}/stats [get]
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
