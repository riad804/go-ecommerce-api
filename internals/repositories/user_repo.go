package repositories

import "go.mongodb.org/mongo-driver/mongo"

type UserRepository interface{}

type userRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		db: db,
	}
}
