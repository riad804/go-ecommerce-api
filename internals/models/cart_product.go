package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartProduct struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId            primitive.ObjectID `bson:"user_id,omitempty" json:"user_id"`
	Quantity          int                `bson:"quantity,omitempty" json:"quantity"`
	SelectedSize      string             `bson:"selected_size,omitempty" json:"selected_size"`
	SelectedColor     string             `bson:"selected_color,omitempty" json:"selected_color"`
	ProductName       string             `bson:"product_name,omitempty" json:"product_name"`
	ProductImage      string             `bson:"product_image,omitempty" json:"product_image"`
	ProductPrice      float64            `bson:"product_price,omitempty" json:"product_price"`
	ReservationExpiry time.Time          `bson:"reservation_exp,omitempty" json:"reservation_exp"` // default 30d
	Reserved          bool               `bson:"reserved,omitempty" json:"reserved"`
}
