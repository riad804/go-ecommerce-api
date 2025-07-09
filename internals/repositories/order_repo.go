package repositories

import (
	"context"
	"time"

	"github.com/riad804/go_ecommerce_api/internals/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const ORDERS = "orders"

type OrderRepository interface {
	// FindByID()
	// FindAll()
	FindByUserId(userId primitive.ObjectID) ([]models.Order, error)
}

type orderRepository struct {
	db *mongo.Database
}

func NewOrderRepository(db *mongo.Database) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) FindByUserId(userId primitive.ObjectID) ([]models.Order, error) {
	collection := r.db.Collection(ORDERS)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id": userId,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	for cursor.Next(ctx) {
		var order models.Order
		if err := cursor.Decode(&order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
