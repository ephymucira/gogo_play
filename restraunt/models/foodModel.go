package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Food struct {
	ID          primitive.ObjectID  `bson:"_id"`
	Name        *string             `json:"name" validate:"required,max=100,min=2"`
	Description *string             `json:"description"`
	Price       *float64            `json:"price" validate:"required"`
	Food_image  *string             `json:"food_image" validate:"required"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Food_id	    string              `json:"food_id"`
	Menu_id     *string             `json:"menu_id"`
}