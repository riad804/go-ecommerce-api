package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/riad804/go_ecommerce_api/helpers"
	"github.com/riad804/go_ecommerce_api/internals/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const USERS = "users"

type UserRepository interface {
	Create(user models.User) (*mongo.InsertOneResult, error)
	Update(user models.User) (*models.User, error)
	FindByID(id primitive.ObjectID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindAll() ([]models.User, error)
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

func (u *userRepository) Update(user models.User) (*models.User, error) {
	collection := u.db.Collection(USERS)
	filter := bson.M{"_id": user.ID}
	user.UpdatedAt = time.Now()
	data, err := helpers.StructToBsonMap(user)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse user model to bson")
	}
	update := bson.M{
		"$set": data,
	}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	return &user, err
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

func (u *userRepository) FindAll() ([]models.User, error) {
	collection := u.db.Collection(USERS)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		user.Password = ""
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// "$set": bson.M{
// 	"name":                       user.Name,
// 	"email":                      user.Email,
// 	"password":                   user.Password,
// 	"street":                     user.Street,
// 	"apartment":                  user.Apartment,
// 	"city":                       user.City,
// 	"postal_code":                user.PostalCode,
// 	"phone":                      user.Phone,
// 	"is_admin":                   user.IsAdmin,
// 	"reset_password_otp":         user.ResetPasswordOtp,
// 	"reset_password_otp_expires": user.ResetPasswordOtpExpires,
// 	"wishlist":                   user.Wishlist,
// 	"updated_at":                 time.Now(),
// },
