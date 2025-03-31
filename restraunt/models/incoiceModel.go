package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID               primitive.ObjectID  `bson:"_id" json:"id"`
	Invoice_id       *string            `json:"invoice_id"`
	Order_id         *string            `json:"order_id"`
	Payment_method   *string            `json:"payment_method" validate:"eq=CARD|eq=CASH"`
	Payment_status   *string            `json:"payment_status" validate:"eq=PAID|eq=PENDING"`
	Payment_due_date *time.Time         `json:"payment_due_date"`
	Total_amount     *float64           `json:"total_amount"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}