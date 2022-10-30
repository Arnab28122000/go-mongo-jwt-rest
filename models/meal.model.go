package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Meal struct {
	ID                primitive.ObjectID `bson:"_id"`
	Name              *string            `json:"name" validate:"required,min=2,max=100"`
	Type              *string            `json:"type"` //veg/non-veg/vegan
	Calories          *string            `json:"calories" validate:"required,min=2"`
	Fat               *string            `json:"fat" validate:"required"`
	Carbohydrate      *string            `json:"required"`
	Protein           *string            `json:"protein"`
	Vitamins_Minerals *string            `json:"vitamins_minerals"`
	MealImage         *string            `json:"meal_image" validate:"required,min=2"`
	Ingredients       []string           `json:"ingredients"`
	Created_at        time.Time          `json:"created_at"`
	Updated_at        time.Time          `json:"updated_at"`
	Meal_id           string             `json:"meal_id"`
}
