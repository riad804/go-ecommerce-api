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

const ORDERS = "orders"
const ORDER_ITEMS = "order_items"
const CART_PRODUCTS = "cart_products"

type OrderRepository interface {
	FindOrdersCount() (*int64, error)
	FindAllOrders() ([]models.OrderResponse, error)
	FindOrderByUserId(userId primitive.ObjectID) ([]models.Order, error)
	DeleteOrderByUserId(userId primitive.ObjectID) error
	DeleteOrderItems(ids []primitive.ObjectID) error
	DeleteCartByUserId(userId primitive.ObjectID) error
	FindOrderById(id primitive.ObjectID) (*models.Order, error)
	UpdateOrder(order models.Order) (*models.Order, error)
	DeleteOrderById(id primitive.ObjectID) *mongo.SingleResult
}

type orderRepository struct {
	db *mongo.Database
}

func NewOrderRepository(db *mongo.Database) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) FindOrdersCount() (*int64, error) {
	collection := r.db.Collection(ORDERS)
	count, err := collection.CountDocuments(context.Background(), bson.M{})
	return &count, err
}

func (r *orderRepository) FindAllOrders() ([]models.OrderResponse, error) {
	collection := r.db.Collection(ORDERS)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pipeline := mongo.Pipeline{
		// Join user
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "user"},
		}}},
		{{Key: "$unwind", Value: "$user"}},

		// Join order_items
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "order_items"},
			{Key: "localField", Value: "order_items"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "order_items"},
		}}},

		// Unwind and join each product inside order_items
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$order_items"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "products"},
			{Key: "localField", Value: "order_items.product_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "order_items.product"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$order_items.product"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},

		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "reviews"},
			{Key: "localField", Value: "order_items.product.reviews"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "order_items.product.reviews"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$order_items.product.reviews"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "categories"},
			{Key: "localField", Value: "order_items.product.category"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "order_items.product.category"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$order_items.product.category"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},

		// Group back to reconstruct order_items array
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id"},
			{Key: "shipping_address", Value: bson.D{{Key: "$first", Value: "$shipping_address"}}},
			{Key: "city", Value: bson.D{{Key: "$first", Value: "$city"}}},
			{Key: "postal_code", Value: bson.D{{Key: "$first", Value: "$postal_code"}}},
			{Key: "country", Value: bson.D{{Key: "$first", Value: "$country"}}},
			{Key: "phone", Value: bson.D{{Key: "$first", Value: "$phone"}}},
			{Key: "payment_id", Value: bson.D{{Key: "$first", Value: "$payment_id"}}},
			{Key: "status", Value: bson.D{{Key: "$first", Value: "$status"}}},
			{Key: "total_price", Value: bson.D{{Key: "$first", Value: "$total_price"}}},
			{Key: "date_ordered", Value: bson.D{{Key: "$first", Value: "$date_ordered"}}},
			{Key: "user", Value: bson.D{{Key: "$first", Value: "$user"}}},
			{Key: "order_items", Value: bson.D{{Key: "$push", Value: "$order_items"}}},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []models.OrderResponse
	if err := cursor.All(ctx, &orders); err != nil {
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

func (r *orderRepository) FindOrderById(id primitive.ObjectID) (*models.Order, error) {
	collection := r.db.Collection(ORDERS)
	var order models.Order
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&order)
	return &order, err
}

func (r *orderRepository) UpdateOrder(order models.Order) (*models.Order, error) {
	collection := r.db.Collection(ORDERS)
	filter := bson.M{"_id": order.ID}
	data, err := helpers.StructToBsonMap(order)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse user model to bson")
	}
	update := bson.M{
		"$set": data,
	}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	return &order, err
}

func (r *orderRepository) DeleteOrderById(id primitive.ObjectID) *mongo.SingleResult {
	collection := r.db.Collection(ORDERS)
	filter := bson.M{"_id": id}
	result := collection.FindOneAndDelete(context.Background(), filter)
	return result
}
