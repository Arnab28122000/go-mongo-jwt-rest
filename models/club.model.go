package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Club struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       *string            `json:"name" validate:"required,min=2,max=100"`
	Email      *string            `json:"email" validate:"email"`
	Phone      *string            `json:"phone" validate:"required,min=10,max=10"`
	Address    Address            `json:"address"`
	ClubImage  *string            `json:"club_image" validate:"required,min=2"`
	Admins     []string           `json:"admins"`
	SubAdmins  []string           `json:"subadmins"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
	Club_id    string             `json:"club_id"`
}
