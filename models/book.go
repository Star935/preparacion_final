package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title        string             `json:"title" bson:"title"`
	Author       string             `json:"author" bson:"author"`
	Isbn         string             `json:"isbn" bson:"isbn"`
	Availability int                `json:"availability" bson:"availability"`
}