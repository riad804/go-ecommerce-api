package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID                      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name                    string             `bson:"name" json:"name" validate:"required"`
	Email                   string             `bson:"email" json:"email" validate:"required,email"`
	Password                string             `bson:"password" json:"-" validate:"required"`
	Street                  *string            `bson:"street,omitempty" json:"street,omitempty"`
	Apartment               *string            `bson:"apartment,omitempty" json:"apartment,omitempty"`
	City                    *string            `bson:"city,omitempty" json:"city,omitempty"`
	PostalCode              *string            `bson:"postal_code,omitempty" json:"postal_code,omitempty"`
	Phone                   string             `bson:"phone" json:"phone" validate:"required"`
	IsAdmin                 bool               `bson:"is_admin" json:"is_admin"`
	ResetPasswordOtp        *int64             `bson:"reset_password_otp,omitempty" json:"reset_password_otp"`
	ResetPasswordOtpExpires *time.Time         `bson:"reset_password_otp_expires,omitempty" json:"reset_password_otp_expires"`
	Wishlist                []Wishlist         `bson:"wishlist,omitempty" json:"wishlist,omitempty"`
	CreatedAt               time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt               time.Time          `bson:"updated_at" json:"updated_at"`
}

type Wishlist struct {
	ProductID    primitive.ObjectID `bson:"product_id" json:"product_id" validate:"required"`
	ProductName  string             `bson:"product_name" json:"product_name" validate:"required"`
	ProductImage string             `bson:"product_image" json:"product_image" validate:"required"`
	ProductPrice float64            `bson:"product_price" json:"product_price" validate:"required"`
}

func EnsureUserIndexes(collection *mongo.Collection) error {
	mod := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.Background(), mod)
	return err
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,password"`
	Phone    string `json:"phone" validate:"required,e164"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,password"`
}

type ForgotPassRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyOtpRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,password"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required,e164"`
}
