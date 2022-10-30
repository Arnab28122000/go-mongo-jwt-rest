package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Exercise struct {
	ID             primitive.ObjectID `bson:"_id"`
	Name           *string            `json:"name" validate:"required,min=2,max=100"`
	Exercise_image string             `json:"workout_image" validate:"required,min=2"`
	Created_at     time.Time          `json:"created_at"`
	Updated_at     time.Time          `json:"updated_at"`
	Exercise_id    string             `json:"workout_id"`
}
