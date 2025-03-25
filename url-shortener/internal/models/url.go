package models

import (
	"errors"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URL struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OriginalURL string             `bson:"original_url" json:"url" binding:"required,url"`
	ShortCode   string             `bson:"short_code" json:"shortCode"`
	AccessCount int64              `bson:"access_count" json:"accessCount"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

// GenerateShortCode creates a unique short code for the URL
func (u *URL) GenerateShortCode() {
	// We'll implement a custom short code generation method
	u.ShortCode = generateUniqueShortCode()
}

// Validate checks if the URL is valid
func (u *URL) Validate() error {
	if u.OriginalURL == "" {
		return ErrInvalidURL
	}
	return nil
}

// Error types
var (
	ErrInvalidURL         = errors.New("invalid URL")
	ErrDuplicateShortCode = errors.New("short code already exists")
)

// Helper function to generate unique short code
func generateUniqueShortCode() string {
	// Generate a random 6-character alphanumeric code
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
