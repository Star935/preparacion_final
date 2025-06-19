package handlers

import "go.mongodb.org/mongo-driver/mongo"

type Handler struct {
	Books *mongo.Collection
	Users *mongo.Collection
	Loans *mongo.Collection
}

func NewHandler(books, users, loans *mongo.Collection) *Handler {
	return &Handler{
		Books: books, 
		Users: users,
		Loans: loans,
	}
}