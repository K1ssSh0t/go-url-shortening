package main

import (
	"log"
	"os"
	_ "url-shortener/docs"            // Importa los docs generados por swag
	"url-shortener/internal/database" // Ruta a tu paquete de base de datos
	"url-shortener/internal/handlers" // Ruta a tu paquete de handlers

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv" // Para cargar .env
)

func setupRoutes(app *fiber.App) {
	// Documentación Swagger
	app.Get("/swagger/*", swagger.HandlerDefault) // Acceso en /swagger/index.html

	// Rutas de la API para gestionar las URLs
	apiGroup := app.Group("/shorten") // Agrupamos bajo /shorten
	apiGroup.Post("/", handlers.CreateShortURL)         // POST /shorten
	apiGroup.Get("/:shortCode", handlers.GetOriginalURLData) // GET /shorten/abc123 (devuelve datos)
	apiGroup.Put("/:shortCode", handlers.UpdateShortURL)     // PUT /shorten/abc123
	apiGroup.Delete("/:shortCode", handlers.DeleteShortURL)   // DELETE /shorten/abc123
	apiGroup.Get("/:shortCode/stats", handlers.GetURLStats) // GET /shorten/abc123/stats

	// Ruta para la redirección directa usando el código corto
	// Esta ruta está fuera del grupo /shorten para que sea más corta (ej: tudominio.com/abc123)
	app.Get("/:shortCode", handlers.RedirectShortURL)
}

func main() {
	// Cargar variables de entorno del archivo .env (opcional)
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: No se pudo cargar el archivo .env")
	}

	// Conectar a la base de datos
	database.ConnectDB()

	// Crear instancia de Fiber
	app := fiber.New()

	// Middleware
	app.Use(logger.New()) // Registrar logs de peticiones HTTP
	app.Use(cors.New())   // Habilitar CORS para permitir peticiones desde otros dominios (frontend)

	// Configurar las rutas
	setupRoutes(app)

	// Obtener el puerto de la variable de entorno o usar uno por defecto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Puerto por defecto
	}

	log.Printf("Servidor escuchando en el puerto %s", port)
	// Iniciar el servidor
	log.Fatal(app.Listen(":" + port)) // log.Fatal si Listen devuelve un error
}