package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/handlers"
	"backend/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

// setupTestDB conecta a MongoDB local y prepara la colección de pruebas.
func setupTestDB(t *testing.T) (*mongo.Collection, func()) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}
	db := client.Database("testdb")
	coll := db.Collection("books")

	cleanup := func() {
		// Limpia la base de datos de pruebas y cierra la conexión
		if err := db.Drop(context.Background()); err != nil {
			t.Fatalf("Error dropping test database: %v", err)
		}
		client.Disconnect(context.Background())
	}
	return coll, cleanup
}

// TestCreateBook verifica que CreateBook inserta correctamente un libro en
func TestCreateBook(t *testing.T) {
	coll, cleanup := setupTestDB(t)
	defer cleanup()

	// Prepara el handler con la colección de pruebas
	h := &handlers.Handler{Books: coll}
	e := echo.New()

	// Crea el cuerpo JSON para el POST
	newBook := models.Book{
		Title:        "Author",
		Author:       "Test Author",
		Isbn:         "ISBN789",
		Availability: 3,
	}

	body, _ := json.Marshal(newBook)
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Ejecuta CreateBook
	if err := h.CreateBook(c); err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	// Verifica el código HTTP
	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d but got %d", http.StatusCreated, rec.Code)
	}

	// Verifica que el documento fue insertado
	count, err := coll.CountDocuments(context.Background(), bson.M{"isbn": newBook.Isbn})
	if err != nil {
		t.Fatalf("Error counting documents: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 document, found %d", count)
	}
}
