package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Loan struct {
	ID    		primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  		string 			   `json:"name" bson:"name"`
	Description string 			   `json:"description" bson:"description"`
	UserId      string 			   `json:"user_id" bson:"user_id"`
	BookId      string 			   `json:"book_id" bson:"book_id"`
	IsReturned  bool      		   `json:"is_returned" bson:"is_returned"`
}