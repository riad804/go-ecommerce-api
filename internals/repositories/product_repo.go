package repositories

import (
	"context"
	"fmt"

	"github.com/riad804/go_ecommerce_api/helpers"
	"github.com/riad804/go_ecommerce_api/internals/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const PRODUCTS = "products"
const CATEGORIES = "categories"

type ProductRepository interface {
	CategoryFindOne(id primitive.ObjectID) (*models.Category, error)
	CategorySave(category models.Category) (*mongo.InsertOneResult, error)
	CategoryDeleteById(id primitive.ObjectID) error
	CategoryUpdate(cat models.Category) (*models.Category, error)

	// products
	CountAllProducts() (int64, error)
	ProductSave(product models.Product) (*mongo.InsertOneResult, error)
	ProductUpdate(product models.Product) (*models.Product, error)
	ProductFindOne(id primitive.ObjectID) (*models.Product, error)
}

type productRepository struct {
	db *mongo.Database
}

func NewProductRepository(db *mongo.Database) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) CategoryFindOne(id primitive.ObjectID) (*models.Category, error) {
	collection := r.db.Collection(CATEGORIES)
	var cat models.Category
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&cat)
	return &cat, err
}

func (r *productRepository) CategorySave(category models.Category) (*mongo.InsertOneResult, error) {
	collection := r.db.Collection(CATEGORIES)
	result, err := collection.InsertOne(context.Background(), category)
	return result, err
}

func (r *productRepository) CategoryDeleteById(id primitive.ObjectID) error {
	collection := r.db.Collection(CATEGORIES)
	var category models.Category
	collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(category)
	category.MarkedForDelete = true
	data, err := helpers.StructToBsonMap(category)
	if err != nil {
		return fmt.Errorf("couldn't parse user model to bson")
	}
	update := bson.M{
		"$set": data,
	}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	return err
}

func (r *productRepository) CategoryUpdate(cat models.Category) (*models.Category, error) {
	collection := r.db.Collection(CATEGORIES)
	data, err := helpers.StructToBsonMap(cat)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse user model to bson")
	}
	update := bson.M{
		"$set": data,
	}
	var updatedCategory models.Category
	result := collection.FindOneAndUpdate(context.Background(), bson.M{"_id": cat.ID}, update).Decode(&updatedCategory)
	if result != nil {
		return nil, result
	}
	return &updatedCategory, nil
}

func (r *productRepository) CountAllProducts() (int64, error) {
	collection := r.db.Collection(PRODUCTS)
	count, err := collection.CountDocuments(context.Background(), bson.M{})
	return count, err
}

func (r *productRepository) ProductSave(product models.Product) (*mongo.InsertOneResult, error) {
	product = product.Pre(context.Background(), r.db)
	collection := r.db.Collection(PRODUCTS)
	result, err := collection.InsertOne(context.Background(), product)
	return result, err
}

func (r *productRepository) ProductFindOne(id primitive.ObjectID) (*models.Product, error) {
	collection := r.db.Collection(PRODUCTS)
	var product models.Product
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&product)
	return &product, err
}

func (r *productRepository) ProductUpdate(product models.Product) (*models.Product, error) {
	collection := r.db.Collection(PRODUCTS)
	data, err := helpers.StructToBsonMap(product)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse product model to bson")
	}
	update := bson.M{
		"$set": data,
	}
	var updatedProduct models.Product
	result := collection.FindOneAndUpdate(context.Background(), bson.M{"_id": product.ID}, update).Decode(&updatedProduct)
	if result != nil {
		return nil, result
	}
	return &updatedProduct, nil
}
