package repository

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"url-shortener/internal/models"
)

type URLRepository struct {
	collection *mongo.Collection
}

func NewURLRepository(client *mongo.Client, dbName string) *URLRepository {
	collection := client.Database(dbName).Collection("urls")
	
	// Create unique index on short_code
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"short_code": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create unique index")
	}

	return &URLRepository{collection: collection}
}

func (r *URLRepository) Create(ctx context.Context, url *models.URL) error {
	url.CreatedAt = time.Now()
	url.UpdatedAt = time.Now()
	url.AccessCount = 0

	// Attempt to generate a unique short code
	for attempts := 0; attempts < 5; attempts++ {
		url.GenerateShortCode()
		
		_, err := r.collection.InsertOne(ctx, url)
		if err == nil {
			return nil
		}
		
		// Check if it's a duplicate key error
		if isDuplicateKeyError(err) {
			continue
		}
		
		return err
	}

	return models.ErrDuplicateShortCode
}

func (r *URLRepository) FindByShortCode(ctx context.Context, shortCode string) (*models.URL, error) {
	var url models.URL
	err := r.collection.FindOne(ctx, bson.M{"short_code": shortCode}).Decode(&url)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &url, nil
}

func (r *URLRepository) Update(ctx context.Context, shortCode string, newURL string) (*models.URL, error) {
	filter := bson.M{"short_code": shortCode}
	update := bson.M{
		"$set": bson.M{
			"original_url": newURL,
			"updated_at":   time.Now(),
		},
	}
	
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedURL models.URL
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedURL)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	
	return &updatedURL, nil
}

func (r *URLRepository) Delete(ctx context.Context, shortCode string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"short_code": shortCode})
	if err != nil {
		return err
	}
	
	if result.DeletedCount == 0 {
		return nil // Treat as not found
	}
	
	return nil
}

func (r *URLRepository) IncrementAccessCount(ctx context.Context, shortCode string) error {
	filter := bson.M{"short_code": shortCode}
	update := bson.M{"$inc": bson.M{"access_count": 1}}
	
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Helper function to check for duplicate key errors
func isDuplicateKeyError(err error) bool {
	// Check if the error is a MongoDB duplicate key error
	mongoErr, ok := err.(mongo.WriteException)
	if !ok {
		return false
	}
	
	for _, writeErr := range mongoErr.WriteErrors {
		if writeErr.Code == 11000 {
			return true
		}
	}
	
	return false
}