package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	State   string `json:"state" bson:"state"`
	City    string `json:"city" bson:"city"`
	Country string `json:"country"`
	Pincode int    `json:"pincode" bson:"pincode"`
}

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	First_name    *string            `json:"first_name" validate:"required,min=2,max=100"`
	Last_name     *string            `json:"last_name" validate:"required,min=2,max=100"`
	Password      *string            `json:"password" validate:"required,min=6,max=100"`
	Email         *string            `json:"email" validate:"email"`
	Phone         *string            `json:"phone" validate:"required,min=10,max=10"`
	UserImage     *string            `json:"user_image" validate:"required,min=2"`
	BloodGroup    *string            `json:"blood_group" validate:"required,min=2"`
	BirthDate     time.Time          `json:"birth_date" validate:"required"`
	Age           int64              `json:"age"`
	Address       Address            `json:"address"`
	Token         *string            `json:"token"`
	Club          []string           `json:"gym"`
	Workout       []string           `json:"workout"`
	Diet          []string           `json:"diet"`
	Subscription  []string           `json:"subscription"`
	Refresh_token *string            `json:"refresh_token"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"user_id"`
}
