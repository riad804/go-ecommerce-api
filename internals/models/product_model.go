package models

import (
	"context"
	"log"
	"time"

	"github.com/riad804/go_ecommerce_api/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GenderAge string

const (
	Men   GenderAge = "men"
	Women GenderAge = "women"
	Kids  GenderAge = "kids"
)

type Product struct {
	ID                primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name              string               `bson:"name,omitempty" json:"name"`
	Description       string               `bson:"description,omitempty" json:"description"`
	Price             float64              `bson:"price,omitempty" json:"price"`
	Rating            float64              `bson:"rating,omitempty" json:"rating"`
	Colors            []string             `bson:"colors,omitempty" json:"colors"`
	Image             string               `bson:"image,omitempty" json:"image"`
	Images            []string             `bson:"images,omitempty" json:"images"`
	Reviews           []primitive.ObjectID `bson:"reviews_ids,omitempty" json:"reviews_ids"`
	NumberOfReviews   int                  `bson:"number_of_reviews,omitempty" json:"number_of_reviews"`
	Sizes             []string             `bson:"sizes,omitempty" json:"sizes"`
	Category          []primitive.ObjectID `bson:"category_ids,omitempty" json:"category_ids"`
	GenderAgeCategory GenderAge            `bson:"gender_age_category,omitempty" json:"gender_age_category"`
	CountInStock      int                  `bson:"count_in_stock,omitempty" json:"count_in_stock"`
	DateAdded         time.Time            `bson:"date_added,omitempty" json:"date_added"`
}

type Review struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Comment string             `bson:"comment"`
	Rating  int                `bson:"rating"`
}

func (p *Product) Pre(ctx context.Context, db *mongo.Database) Product {
	if len(p.Reviews) > 0 {
		reviews := p.getReviews(ctx, db)

		totalRating := helpers.Reduce(reviews, 0, func(acc int, val Review) int {
			return acc + val.Rating
		})
		p.Rating = float64(totalRating) / float64(len(reviews))
		p.NumberOfReviews = len(reviews)
	}
	return *p
}

func (p *Product) getReviews(ctx context.Context, db *mongo.Database) []Review {
	var reviews []Review
	cursor, err := db.Collection("reviews").Find(ctx, bson.M{
		"_id": bson.M{"$in": p.Reviews},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var r Review
		cursor.Decode(&r)
		reviews = append(reviews, r)
	}

	return reviews
}

func EnsureProductIndexes(collection *mongo.Collection) error {
	mod := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
			{Key: "description", Value: 1},
		},
		Options: options.Index(),
	}
	_, err := collection.Indexes().CreateOne(context.Background(), mod)
	return err
}
