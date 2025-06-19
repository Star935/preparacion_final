package handlers

import (
	"context"
	"net/http"
	"strings"

	"backend/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Recupera todos los inventarios junto con sus objetos anidados
func (h *Handler) GetLoans(c echo.Context) error {
	// Valida la conexion a la coleccion
	if h.Loans == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound, 
			"message" : "Sin conexion a la colección",
			"data"	  : nil,
		})
	}

	// Recupera todos los inventarios
	cur, err := h.Loans.Find(context.Background(), bson.M{})
	// Valuda si recupera los inventarios
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status" : http.StatusInternalServerError,
			"message": err.Error(),
			"data"   : nil,
		})
	}

	// Lista de usuarios
	var loans []models.Loan

	// Almacena en la lista de inventarios todos los inventarios recuperados y valida si la operacion es exitosa
	if err := cur.All(context.Background(), &loans); err != nil{
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status"  : http.StatusInternalServerError,
			"message" : err.Error(),
			"data"    : nil,
		})
	}

	if len(loans) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound,
			"message" : "No se encontraron usuarios",
			"data"	  : nil,
		})
	}

	// Retorna estado de respuesta ok y todos los inventarios recuperados
	return c.JSON(http.StatusFound, echo.Map{
		"status"  : http.StatusFound,
		"message" : "Lista de usuarios encontrada",
		"data"    : loans,
	})
}

// Crea un nuevo inventario
func (h *Handler) CreateLoan(c echo.Context) error {
	// Valida la conexion a la coleccion
	if h.Loans == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound,
			"message" : "Sin conexion a la colección",
			"data"    : nil,
		})
	}

	var loan models.Loan

	if err := c.Bind(&loan); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "Input invalido",
			"data"    : nil,
		})
	}

	if strings.TrimSpace(loan.Name) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El nombre es obligatorio",
			"data"    : nil,
		})
	}

	if strings.TrimSpace(loan.Description) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "La descripción es obligatoria",
			"data"    : nil,
		})
	}

	_, err := h.Loans.InsertOne(context.Background(), loan)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status"  : http.StatusInternalServerError, 
			"message" : err.Error(),
			"data"	  : nil,
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"status"  : http.StatusCreated, 
		"message" : "Usuario creado exitosamente",
		"data"	  : loan,
	})
}

// Actualiza un prestamo existente
func (h *Handler) ReturnLoan(c echo.Context) error {
	// Valida la conexion a la coleccion
	if h.Loans == nil {
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

	// Prepara filtro y documento de actualización
    filter := bson.M{"_id": id}
    update := bson.M{
        "$set": bson.M{
        	"is_returned" : true,
        },
    }

	// Actualiza el documento
    res, err := h.Loans.UpdateOne(context.Background(), filter, update)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{
            "status"  : http.StatusInternalServerError,
            "message" : err.Error(),
            "data"    : nil,
        })
    }

    if res.MatchedCount == 0 {
        return c.JSON(http.StatusNotFound, echo.Map{
            "status"  : http.StatusNotFound,
            "message" : "Prestamo no encontrado",
            "data"    : nil,
        })
    }

	return c.JSON(http.StatusCreated, echo.Map{
        "status"  : http.StatusCreated,
        "message" : "Prestamo devuelto exitosamente!",
        "data"    : nil,
    })
}