package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subscription struct {
	ID                primitive.ObjectID `bson:"_id"`
	Name              *string            `json:"name" validate:"required,min=2,max=100"`
	Description       *string            `json:"description"`
	Club_id           string             `bson:"club"`
	Plan_id           string             `bson:"plan"`              //Meal_id,
	Subscription_type string             `bson:"subscription_type"` //Monthly, Yearly
	Subscription_id   string             `bson:"subscription_id"`
	Payment_method    string             `bson:"payment_method"` //cash/Online
	Price             string             `bson:"price"`
	Duration          string             `bson:"duration"` //Month
	Expiry_date       time.Time          `bson:"expiry_date"`
	Created_at        time.Time          `json:"created_at"`
	Updated_at        time.Time          `json:"updated_at"`
}
