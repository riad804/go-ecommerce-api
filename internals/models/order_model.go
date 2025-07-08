package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	Pending        OrderStatus = "pending"
	Processed      OrderStatus = "processed"
	Shipped        OrderStatus = "shipped"
	OutForDelivery OrderStatus = "out-for-delivery"
	Delivered      OrderStatus = "delivered"
	Cancelled      OrderStatus = "cancelled"
	OnHold         OrderStatus = "on-hold"
	Expired        OrderStatus = "expired"
)

type Order struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	ShippingAddress string               `bson:"shipping_address,omitempty" json:"shipping_address"`
	City            string               `bson:"city,omitempty" json:"city"`
	PostalCode      string               `bson:"postal_code,omitempty" json:"postal_code"`
	Country         string               `bson:"country,omitempty" json:"country"`
	Phone           string               `bson:"phone,omitempty" json:"phone"`
	PaymentId       *string              `bson:"payment_id,omitempty" json:"payment_id"`
	Status          OrderStatus          `bson:"status,omitempty" json:"status"`
	StatusHistory   []OrderStatus        `bson:"status_history,omitempty" json:"status_history"`
	TotalPrice      float64              `bson:"total_price,omitempty" json:"total_price"`
	UserId          primitive.ObjectID   `bson:"user_id,omitempty" json:"user_id"`
	DateOrdered     time.Time            `bson:"date_ordered,omitempty" json:"date_ordered"`
	OrderItems      []primitive.ObjectID `bson:"order_items,omitempty" json:"order_items"`
}

type OrderItem struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductId     primitive.ObjectID `bson:"product_id,omitempty" json:"product_id"`
	ProductName   string             `bson:"product_name,omitempty" json:"product_name"`
	ProductImage  string             `bson:"product_image,omitempty" json:"product_image"`
	ProductPrice  float64            `bson:"product_price,omitempty" json:"product_price"`
	Quantity      int                `bson:"quantity,omitempty" json:"quantity"`
	SelectedSize  string             `bson:"selected_size,omitempty" json:"selected_size"`
	SelectedColor string             `bson:"selected_color,omitempty" json:"selected_color"`
}
