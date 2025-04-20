package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client     *mongo.Client
	UrlCollection *mongo.Collection
)

func ConnectDB() {
	mongoURI := os.Getenv("MONGO_URI") // Carga desde variable de entorno
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // Valor por defecto si no está definida
		log.Println("Advertencia: MONGO_URI no definida, usando valor por defecto:", mongoURI)
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "urlshortener" // Nombre de BD por defecto
		log.Println("Advertencia: DB_NAME no definida, usando valor por defecto:", dbName)
	}

	collectionName := "urls" // Nombre de la colección

	log.Println("Conectando a MongoDB...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Contexto con timeout
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error al conectar a MongoDB: %v", err)
	}

	// Comprobar la conexión
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Error al hacer ping a MongoDB: %v", err)
	}

	Client = client
	database := client.Database(dbName)
	UrlCollection = database.Collection(collectionName)

	log.Printf("¡Conectado a MongoDB! Base de datos: %s, Colección: %s", dbName, collectionName)
}

// Función para obtener la colección (útil si tienes más colecciones)
func GetCollection(collectionName string) *mongo.Collection {
	if Client == nil {
		log.Fatal("El cliente de MongoDB no está inicializado. Llama a ConnectDB primero.")
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "url_shortener_db"
	}
	return Client.Database(dbName).Collection(collectionName)
}