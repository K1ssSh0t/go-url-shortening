package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"url-shortener/internal/database" // Asegúrate de que la ruta sea correcta
	"url-shortener/internal/models"   // Asegúrate de que la ruta sea correcta
	"url-shortener/internal/utils"    // Para IsValidURL si la usas

	"github.com/gofiber/fiber/v2"
	"github.com/teris-io/shortid" // Generador de IDs cortos
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

// CreateShortURL maneja la creación de una nueva URL corta
// @Summary Crea una URL corta
// @Description Genera un código corto para una URL original y la almacena en la base de datos
// @Tags URLs
// @Accept json
// @Produce json
// @Param data body models.CreateURLRequest true "URL a acortar"
// @Success 201 {object} models.URL
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /shorten/ [post]
func CreateShortURL(c *fiber.Ctx) error {
	req := new(models.CreateURLRequest)

	// Parsear el cuerpo de la petición
	if err := c.BodyParser(req); err != nil {
		log.Printf("Error al parsear el cuerpo de la petición: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo procesar la petición",
		})
	}

	// Validar la URL (usando nuestra función o simplemente verificando que no esté vacía)
	if req.URL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El campo 'url' es requerido",
		})
	}
	// Validación más robusta (opcional)
	if !utils.IsValidURL(req.URL) {
		 return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			 "error": "La URL proporcionada no es válida",
		 })
	}


	// Generar un código corto único
	shortCode, err := shortid.Generate()
	if err != nil {
		log.Printf("Error al generar shortid: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error interno al generar el código corto",
		})
	}
	// Podrías añadir lógica para asegurar unicidad si shortid no fuera suficiente
	// (consultar la BD si el código ya existe y regenerar si es necesario).

	now := time.Now()
	newURL := models.URL{
		// ID se generará automáticamente por MongoDB
		OriginalURL: req.URL,
		ShortCode:   shortCode,
		CreatedAt:   now,
		UpdatedAt:   now,
		AccessCount: 0,
	}

	// Insertar en la base de datos
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := database.UrlCollection.InsertOne(ctx, newURL)
	if err != nil {
		log.Printf("Error al insertar URL en la base de datos: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al guardar la URL",
		})
	}

	// Recuperar el ID insertado para devolverlo en la respuesta
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("Error al obtener el ID insertado")
        // Aún así podemos devolver el resto de la info, o un error genérico
        // Devolveremos la info sin el ID en este caso raro
         return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener ID después de guardar",
		})
	}

	newURL.ID = insertedID // Asignamos el ID generado

	log.Printf("URL corta creada: %s -> %s", newURL.ShortCode, newURL.OriginalURL)
	return c.Status(fiber.StatusCreated).JSON(newURL) // Devuelve el objeto completo
}

// GetOriginalURLData maneja la obtención de la URL original a partir del código corto
// @Summary Obtiene datos de la URL original
// @Description Devuelve los datos de la URL original asociada a un código corto (no redirige)
// @Tags URLs
// @Produce json
// @Param shortCode path string true "Código corto"
// @Success 200 {object} models.URL
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /shorten/{shortCode} [get]
func GetOriginalURLData(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")

	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El código corto es requerido en la ruta",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var urlData models.URL
	filter := bson.M{"shortCode": shortCode}

	err := database.UrlCollection.FindOne(ctx, filter).Decode(&urlData)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("Código corto no encontrado: %s", shortCode)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "URL corta no encontrada",
			})
		}
		log.Printf("Error al buscar URL en la base de datos: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al buscar la URL",
		})
	}

	// NO incrementamos el contador aquí, eso se hace en la redirección.
	// Opcionalmente, podrías decidir incrementarlo también aquí si este endpoint se usa mucho.

	log.Printf("Datos de URL recuperados para: %s", shortCode)
	return c.Status(fiber.StatusOK).JSON(urlData)
}

// RedirectShortURL maneja la redirección al encontrar una URL corta
// @Summary Redirige a la URL original
// @Description Redirige al usuario a la URL original usando el código corto
// @Tags URLs
// @Param shortCode path string true "Código corto"
// @Success 301
// @Failure 400 {object} map[string]string
// @Failure 404 {string} string "URL corta no encontrada"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /{shortCode} [get]
func RedirectShortURL(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")

	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El código corto es requerido",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var urlData models.URL
	filter := bson.M{"shortCode": shortCode}

	// Usamos FindOneAndUpdate para encontrar e incrementar el contador atómicamente
	update := bson.M{
		"$inc": bson.M{"accessCount": 1},         // Incrementa accessCount en 1
		"$set": bson.M{"updatedAt": time.Now()}, // Actualiza la fecha de último acceso/update
	}
	options := options.FindOneAndUpdate().SetReturnDocument(options.After) // Devuelve el documento después de actualizar

	err := database.UrlCollection.FindOneAndUpdate(ctx, filter, update, options).Decode(&urlData)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("Código corto no encontrado para redirección: %s", shortCode)
			return c.Status(fiber.StatusNotFound).SendString("URL corta no encontrada") // O una página HTML de error
		}
		log.Printf("Error al buscar y actualizar URL para redirección: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error interno del servidor")
	}

	log.Printf("Redirigiendo: %s -> %s (Accesos: %d)", urlData.ShortCode, urlData.OriginalURL, urlData.AccessCount)
	// Realiza la redirección permanente (301) o temporal (302/307)
	return c.Redirect(urlData.OriginalURL, http.StatusMovedPermanently) // O fiber.StatusTemporaryRedirect
}

// UpdateShortURL maneja la actualización de la URL original asociada a un código corto
// @Summary Actualiza la URL original
// @Description Actualiza la URL original asociada a un código corto
// @Tags URLs
// @Accept json
// @Produce json
// @Param shortCode path string true "Código corto"
// @Param data body models.UpdateURLRequest true "Nueva URL"
// @Success 200 {object} models.URL
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /shorten/{shortCode} [put]
func UpdateShortURL(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El código corto es requerido en la ruta",
		})
	}

	req := new(models.UpdateURLRequest)
	if err := c.BodyParser(req); err != nil {
		log.Printf("Error al parsear cuerpo en actualización: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo procesar la petición",
		})
	}

	// Validar la nueva URL
	if req.URL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El campo 'url' es requerido",
		})
	}
	if !utils.IsValidURL(req.URL) {
		 return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			 "error": "La nueva URL proporcionada no es válida",
		 })
	}


	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"shortCode": shortCode}
	update := bson.M{
		"$set": bson.M{
			"originalUrl": req.URL,
			"updatedAt":   time.Now(),
		},
	}
	options := options.FindOneAndUpdate().SetReturnDocument(options.After) // Devolver el documento actualizado

	var updatedURL models.URL
	err := database.UrlCollection.FindOneAndUpdate(ctx, filter, update, options).Decode(&updatedURL)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("Código corto no encontrado para actualizar: %s", shortCode)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "URL corta no encontrada",
			})
		}
		log.Printf("Error al actualizar URL en la base de datos: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al actualizar la URL",
		})
	}

	log.Printf("URL actualizada: %s ahora apunta a %s", updatedURL.ShortCode, updatedURL.OriginalURL)
	return c.Status(fiber.StatusOK).JSON(updatedURL)
}

// DeleteShortURL maneja la eliminación de una URL corta
// @Summary Elimina una URL corta
// @Description Elimina la URL corta asociada a un código
// @Tags URLs
// @Param shortCode path string true "Código corto"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /shorten/{shortCode} [delete]
func DeleteShortURL(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El código corto es requerido en la ruta",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"shortCode": shortCode}
	result, err := database.UrlCollection.DeleteOne(ctx, filter)

	if err != nil {
		log.Printf("Error al eliminar URL de la base de datos: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al eliminar la URL",
		})
	}

	if result.DeletedCount == 0 {
		log.Printf("Código corto no encontrado para eliminar: %s", shortCode)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "URL corta no encontrada",
		})
	}

	log.Printf("URL eliminada: %s", shortCode)
	return c.SendStatus(fiber.StatusNoContent) // 204 No Content para eliminación exitosa
}

// GetURLStats maneja la obtención de estadísticas para una URL corta
// @Summary Obtiene estadísticas de una URL corta
// @Description Devuelve el número de accesos y datos de la URL corta
// @Tags URLs
// @Produce json
// @Param shortCode path string true "Código corto"
// @Success 200 {object} models.URL
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /shorten/{shortCode}/stats [get]
func GetURLStats(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El código corto es requerido en la ruta",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var urlData models.URL // Reutilizamos el modelo ya que contiene AccessCount
	filter := bson.M{"shortCode": shortCode}

	err := database.UrlCollection.FindOne(ctx, filter).Decode(&urlData)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("Código corto no encontrado para estadísticas: %s", shortCode)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "URL corta no encontrada",
			})
		}
		log.Printf("Error al buscar URL para estadísticas: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al buscar la URL",
		})
	}

	// El campo AccessCount ya está en urlData
	log.Printf("Estadísticas recuperadas para: %s (Accesos: %d)", shortCode, urlData.AccessCount)
	return c.Status(fiber.StatusOK).JSON(urlData) // Devolvemos el objeto completo que incluye las estadísticas
}