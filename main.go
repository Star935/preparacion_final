package main

import (
	"context"
	"log"

	"backend/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Instancia de Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	// Define la base de datos y la coleccion
	db := client.Database("practica_parcial_final")

	h := handlers.NewHandler(db.Collection("books"), db.Collection("users"), db.Collection("loans"))

	// Rutas para la gestion de inventarios
	e.GET("/books", h.GetBooks)
	e.GET("/books/:id", h.GetBookById)
	e.POST("/books", h.CreateBook)
	e.PUT("/books/:id", h.UpdateBook)
	e.DELETE("/books/:id", h.DeleteBook)

	// Rutas para la gestion de inventarios
	e.GET("/users", h.GetUsers)
	e.GET("/users/:id", h.GetUserById)
	e.POST("/users", h.CreateUser)
	e.DELETE("/users/:id", h.DeleteUser)

	// Rutas para la gestion de inventarios
	e.GET("/loans", h.GetLoans)
	e.POST("/loans", h.CreateLoan)
	e.PUT("/return-loan/:id", h.ReturnLoan)

	e.Logger.Fatal(e.Start(":8080"))
	// Analisis estatico
	// github.com/securego/gosec/v2/cmd/gosec@latest
	// gosec ./...

	// Analisis de vulnerabilidades conocidas
	// go install golang.org/x/vuln/cmd/govulncheck@latest
	// govulncheck ./...
}