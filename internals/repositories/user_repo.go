package repositories

import (
	"context"
	"time"

	"github.com/riad804/go_ecommerce_api/internals/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const USERS = "users"

type UserRepository interface {
	Create(user models.User) (*mongo.InsertOneResult, error)
	Update(user models.User) (*mongo.UpdateResult, error)
	FindByID(id primitive.ObjectID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Create(user models.User) (*mongo.InsertOneResult, error) {
	collection := u.db.Collection(USERS)
	user.CreatedAt = time.Now()
	result, err := collection.InsertOne(context.Background(), user)
	return result, err
}

func (u *userRepository) Update(user models.User) (*mongo.UpdateResult, error) {
	collection := u.db.Collection(USERS)
	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"name":                       user.Name,
			"email":                      user.Email,
			"password":                   user.Password,
			"street":                     user.Street,
			"apartment":                  user.Apartment,
			"city":                       user.City,
			"postal_code":                user.PostalCode,
			"phone":                      user.Phone,
			"is_admin":                   user.IsAdmin,
			"reset_password_otp":         user.ResetPasswordOtp,
			"reset_password_otp_expires": user.ResetPasswordOtpExpires,
			"wishlist":                   user.Wishlist,
			"updated_at":                 time.Now(),
		},
	}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	return result, err
}

func (u *userRepository) FindByID(id primitive.ObjectID) (*models.User, error) {
	collection := u.db.Collection(USERS)
	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	return &user, err
}

func (u *userRepository) FindByEmail(email string) (*models.User, error) {
	collection := u.db.Collection(USERS)
	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	return &user, err
}
