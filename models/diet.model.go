package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Diet struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       *string            `json:"name" validate:"required,min=2,max=100"`
	Type       *string            `json:"type"`
	DietImage  *string            `json:"diet_image" validate:"required,min=2"`
	Calories   *string            `json:"calories"`
	Meals      []string           `json:"meals"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
	Diet_id    string             `json:"club_id"`
}
