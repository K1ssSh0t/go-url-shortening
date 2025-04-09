package handlers

import (
	"context"
	"time"
	"url-shortener/internal/models"
	"url-shortener/internal/repository"

	"math/rand"

	"github.com/gofiber/fiber/v2"
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
func CreateShortURL(c *fiber.Ctx) error {
    var input struct {
        URL string `json:"url"`
    }

    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
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
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not save URL"})
    }

    return c.Status(fiber.StatusCreated).JSON(newURL)
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
func GetOriginalURL(c *fiber.Ctx) error {
    shortCode := c.Params("shortCode")
    var result models.ShortURL

    err := repository.DB.FindOne(context.TODO(), bson.M{"short_code": shortCode}).Decode(&result)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Short URL not found"})
    }

    repository.DB.UpdateOne(context.TODO(), bson.M{"short_code": shortCode}, bson.M{"$inc": bson.M{"access_count": 1}})
    return c.Status(fiber.StatusOK).JSON(result)
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
func UpdateShortURL(c *fiber.Ctx) error {
    shortCode := c.Params("shortCode")
    var input struct {
        URL string `json:"url"`
    }

    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
    }

    update := bson.M{"$set": bson.M{"url": input.URL, "updated_at": time.Now()}}
    res, err := repository.DB.UpdateOne(context.TODO(), bson.M{"short_code": shortCode}, update)
    if err != nil || res.ModifiedCount == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Short URL not found"})
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "URL updated"})
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
func DeleteShortURL(c *fiber.Ctx) error {
    shortCode := c.Params("shortCode")
    res, err := repository.DB.DeleteOne(context.TODO(), bson.M{"short_code": shortCode})
    if err != nil || res.DeletedCount == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Short URL not found"})
    }

    return c.Status(fiber.StatusNoContent).Send(nil)
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
func GetURLStats(c *fiber.Ctx) error {
    shortCode := c.Params("shortCode")
    var result models.ShortURL

    err := repository.DB.FindOne(context.TODO(), bson.M{"short_code": shortCode}).Decode(&result)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Short URL not found"})
    }

    return c.Status(fiber.StatusOK).JSON(result)
}
