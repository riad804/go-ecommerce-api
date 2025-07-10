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
const ORDER_ITEMS = "order_items"
const CART_PRODUCTS = "cart_products"

type OrderRepository interface {
	FindAllOrders() ([]models.Order, error)
	FindOrderByUserId(userId primitive.ObjectID) ([]models.Order, error)
	DeleteOrderByUserId(userId primitive.ObjectID) error
	DeleteOrderItems(ids []primitive.ObjectID) error
	DeleteCartByUserId(userId primitive.ObjectID) error
}

type orderRepository struct {
	db *mongo.Database
}

func NewOrderRepository(db *mongo.Database) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) FindAllOrders() ([]models.Order, error) {
	collection := r.db.Collection(ORDERS)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
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

func (r *orderRepository) FindOrderByUserId(userId primitive.ObjectID) ([]models.Order, error) {
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

func (r *orderRepository) DeleteOrderByUserId(userId primitive.ObjectID) error {
	collection := r.db.Collection(ORDERS)
	filter := bson.M{"user_id": userId}
	_, err := collection.DeleteMany(context.Background(), filter)
	return err
}

func (r *orderRepository) DeleteOrderItems(ids []primitive.ObjectID) error {
	collection := r.db.Collection(ORDER_ITEMS)
	filter := bson.M{"_id": bson.M{"$in": ids}}
	_, err := collection.DeleteMany(context.Background(), filter)
	return err
}

func (r *orderRepository) DeleteCartByUserId(userId primitive.ObjectID) error {
	collection := r.db.Collection(CART_PRODUCTS)
	filter := bson.M{"user_id": userId}
	_, err := collection.DeleteMany(context.Background(), filter)
	return err
}
