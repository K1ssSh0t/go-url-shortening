package models

import "time"

type ShortURL struct {
    ID        string    `json:"id" bson:"_id,omitempty"`
    URL       string    `json:"url" bson:"url"`
    ShortCode string    `json:"shortCode" bson:"short_code"`
    CreatedAt time.Time `json:"createdAt" bson:"created_at"`
    UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
    AccessCount int     `json:"accessCount" bson:"access_count"`
}
