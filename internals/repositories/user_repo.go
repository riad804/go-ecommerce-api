package repositories

import (
	"context"

	"github.com/riad804/go_ecommerce_api/internals/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const USERS = "users"

type UserRepository interface {
	Create(user models.User) (*mongo.InsertOneResult, error)
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
	result, err := collection.InsertOne(context.Background(), user)
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
