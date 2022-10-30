package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Workout struct {
	ID            primitive.ObjectID `bson:"_id"`
	Name          *string            `json:"name" validate:"required,min=2,max=100"`
	Workout_image string             `json:"workout_image" validate:"required,min=2"`
	Exercise      []string           `json:"exercise"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	Workout_id    string             `json:"workout_id"`
}
