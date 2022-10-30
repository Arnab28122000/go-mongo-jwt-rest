package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID               primitive.ObjectID `bson:"_id"`
	Subscription_id  string             `bson:"subscription_id"`
	Payment_method   string             `bson:"payment_method"` //cash/Online
	Payment_due_date string             `bson:"payment_due_date"`
	Created_at       time.Time          `json:"created_at"`
	Updated_at       time.Time          `json:"updated_at"`
	Invoice_id       string             `json:"invoice_id"`
}
