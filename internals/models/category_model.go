package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name            string             `bson:"name,omitempty" json:"name"`
	Color           string             `bson:"color,omitempty" json:"color"`
	Image           string             `bson:"image,omitempty" json:"image"`
	MarkedForDelete bool               `bson:"marked_for_delete,omitempty" json:"marked_for_delete"`
}

type CategoryCreateRequest struct {
	Name  string `form:"name" validate:"required"`
	Color string `form:"color" validate:"required"`
}
