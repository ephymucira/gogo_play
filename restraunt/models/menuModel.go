package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Menu struct {
	ID          primitive.ObjectID  `bson:"_id"`
	Name        *string             `json:"name" validate:"required,max=100,min=2"`
	Description *string             `json:"description"`
	Category    *string             `json:"category" validate:"required"`
	Start_Date  *time.Time          `json:"start_date"`
	End_Date    *time.Time          `json:"end_date"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Menu_id	    string              `json:"menu_id"`
}