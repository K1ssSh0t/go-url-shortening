package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URL struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"` // _id es el estándar en MongoDB
	OriginalURL string             `bson:"originalUrl" json:"url"`
	ShortCode   string             `bson:"shortCode" json:"shortCode"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	AccessCount int64              `bson:"accessCount" json:"accessCount"` // Usamos int64 para contadores potencialmente grandes
}

// Request para crear una URL
type CreateURLRequest struct {
	URL string `json:"url" validate:"required,url"` // Añadimos validación básica
}

// Request para actualizar una URL
type UpdateURLRequest struct {
	URL string `json:"url" validate:"required,url"`
}