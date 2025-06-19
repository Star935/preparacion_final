package handlers

import (
	"context"
	"net/http"
	"strings"

	"backend/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Recupera todos los inventarios junto con sus objetos anidados
func (h *Handler) GetBooks(c echo.Context) error {
	// Valida la conexion a la coleccion
	if h.Books == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound, 
			"message" : "Sin conexion a la colección",
			"data"	  : nil, 
		})
	}

	// Recupera todos los inventarios
	cur, err := h.Books.Find(context.Background(), bson.M{})
	// Valuda si recupera los inventarios
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status"  : http.StatusInternalServerError, 
			"message" : err.Error(),
			"data"	  : nil, 
		})
	}

	// Lista de libros
	var books []models.Book

	// Almacena en la lista de inventarios todos los inventarios recuperados y valida si la operacion es exitosa
	if err := cur.All(context.Background(), &books); err != nil{
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status"  : http.StatusInternalServerError, 
			"message" : err.Error(),
			"data"	  : nil, 
		})
	}

	if len(books) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound, 
			"message" : "No se encontraron libros",
			"data"	  : nil, 
		})
	}

	// Retorna estado de respuesta ok y todos los inventarios recuperados
	return c.JSON(http.StatusFound, echo.Map{
		"status"  : http.StatusFound,
		"message" : "Lista de libros encontrada",
		"data"    : books,
	})
}

// Recupera un inventario mediante su id
func (h *Handler) GetBookById(c echo.Context) error {
	// Valida la conexion a la coleccion
	if h.Books == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound, 
			"message" : "Sin conexion a la colección",
			"data"	  : nil, 
		})
	}

	// Recupera el parametro de consulta el id
	idParam := c.Param("id")

	// Convierte el id en ObjectID
	id, err := primitive.ObjectIDFromHex(idParam)
	// Valida si la operacion fue exitosa
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "Id invalido",
			"data"	  : nil, 
		})
	}

	// Instancia de Book
	var book models.Book

	// Recupera el inventario mediante su id y lo decodifica en el espacio de memoria de la instancia de inventario
	err = h.Books.FindOne(context.Background(), bson.M{"_id": id}).Decode(&book)
	// Valuda si no existe el documento
	if err == mongo.ErrNoDocuments {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound, 
			"message" : "Libro no encontrado",
			"data"	  : nil, 
		})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status"  : http.StatusInternalServerError, 
			"message" : err.Error(),
			"data"	  : nil, 
		})
	}

	// Retorna estado de respuesta ok y todos los inventarios recuperados
	return c.JSON(http.StatusFound, echo.Map{
		"status"  : http.StatusFound,
		"message" : "Libro encontrado",
		"data"    : book,
	})
}

// Crea un nuevo inventario
func (h *Handler) CreateBook(c echo.Context) error {

	var book models.Book

	if err := c.Bind(&book); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "Input invalido",
			"data"	  : nil, 
		})
	}

	if strings.TrimSpace(book.Title) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El titulo es obligatorio",
			"data"	  : nil, 
		})
	}

	if strings.TrimSpace(book.Author) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El autor es obligatorio",
			"data"	  : nil, 
		})
	}

	if strings.TrimSpace(book.Isbn) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El isbn es obligatorio",
			"data"	  : nil, 
		})
	}

	if book.Availability <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "La disponibilidad es obligatoria y valida",
			"data"	  : nil, 
		})
	}

	_, err := h.Books.InsertOne(context.Background(), book)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status"  : http.StatusInternalServerError, 
			"message" : err.Error(),
			"data"	  : nil, 
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"status":  http.StatusCreated,
        "message": "Libro creado exitosamente",
        "data":    book,
    })
}

// Actualiza un libro existente
func (h *Handler) UpdateBook(c echo.Context) error {
	// Valida la conexion a la coleccion
	if h.Books == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound, 
			"message" : "Sin conexion a la colección",
			"data"	  : nil, 
		})
	}

	// Recupera el parametro de consulta el id
	idParam := c.Param("id")

	// Convierte el id en ObjectID
	id, err := primitive.ObjectIDFromHex(idParam)
	// Valida si la operacion fue exitosa
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "Id invalido",
			"data"	  : nil, 
		})
	}

	var book models.Book

	if err := c.Bind(&book); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "Input invalido",
			"data"	  : nil, 
		})
	}

	if strings.TrimSpace(book.Title) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El titulo es obligatorio",
			"data"	  : nil, 
		})
	}

	if strings.TrimSpace(book.Author) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El autor es obligatorio",
			"data"	  : nil, 
		})
	}

	if strings.TrimSpace(book.Isbn) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El isbn es obligatorio",
			"data"	  : nil, 
		})
	}

	if book.Availability <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "La disponibilidad es obligatoria y valida",
			"data"	  : nil, 
		})
	}

	// Prepara filtro y documento de actualización
    filter := bson.M{"_id": id}
    update := bson.M{
        "$set": bson.M{
            "title":        book.Title,
            "author":       book.Author,
            "isbn":         book.Isbn,
            "availability": book.Availability,
        },
    }

	// Actualiza el documento
    res, err := h.Books.UpdateOne(context.Background(), filter, update)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  http.StatusInternalServerError,
            "message": err.Error(),
            "data":    nil,
        })
    }

    if res.MatchedCount == 0 {
        return c.JSON(http.StatusNotFound, echo.Map{
            "status":  http.StatusNotFound,
			"message": "Libro no encontrado",
            "data":    nil,
        })
    }

	return c.JSON(http.StatusCreated, echo.Map{
		"status":  http.StatusCreated,
        "message": "Libro actualizado exitosamente",
        "data":    book,
    })
}

// Elimina un inventario
func (h *Handler) DeleteBook(c echo.Context) error {
	// Valida la conexion a la coleccion
	if h.Books == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound, 
			"message" : "Sin conexion a la colección",
			"data"	  : nil, 
		})
	}

	// Recupera el parametro de consulta el id
	idParam := c.Param("id")

	// Convierte el id en ObjectID
	id, err := primitive.ObjectIDFromHex(idParam)
	// Valida si la operacion fue exitosa
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "Id invalido",
			"data"	  : nil,
		})
	}

	// Realiza operacion de eliminado mediante el id recuperado del parametro de consulta
	res, err := h.Books.DeleteOne(context.Background(), bson.M{"_id": id})
	// Valida si la operacion fue exitosa
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status"  : http.StatusInternalServerError, 
			"message" : err.Error(),
			"data"	  : nil, 
		})
	}
	
	// Valida si se elimino algun documento
	if res.DeletedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound, 
			"message" : "Libro no encontrado",
			"data"	  : nil, 
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status"  : http.StatusOK,
        "message" : "Libro eliminado exitosamente",
        "data"    : nil,
    })
}