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
func (h *Handler) GetUsers(c echo.Context) error {
	// Valida la conexion a la coleccion
	if h.Users == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound, 
			"message" : "Sin conexion a la colección",
			"data"	  : nil,
		})
	}

	// Recupera todos los inventarios
	cur, err := h.Users.Find(context.Background(), bson.M{})
	// Valuda si recupera los inventarios
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status" : http.StatusInternalServerError,
			"message": err.Error(),
			"data"   : nil,
		})
	}

	// Lista de usuarios
	var users []models.User

	// Almacena en la lista de inventarios todos los inventarios recuperados y valida si la operacion es exitosa
	if err := cur.All(context.Background(), &users); err != nil{
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status"  : http.StatusInternalServerError,
			"message" : err.Error(),
			"data"    : nil,
		})
	}

	if len(users) == 0 {
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
		"data"    : users,
	})
}

// Recupera un inventario mediante su id
func (h *Handler) GetUserById(c echo.Context) error {
	// Valida la conexion a la coleccion
	if h.Users == nil {
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
			"data"    : nil,
		})
	}

	// Instancia de User
	var user models.User

	// Recupera el inventario mediante su id y lo decodifica en el espacio de memoria de la instancia de inventario
	err = h.Users.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	// Valuda si no existe el documento
	if err == mongo.ErrNoDocuments {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound,
			"message" : "Usuario no encontrado",
			"data"    : nil,
		})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status"  : http.StatusInternalServerError,
			"message" : err.Error(),
			"data"    : nil,
		})
	}

	// Retorna estado de respuesta ok y todos los inventarios recuperados
	return c.JSON(http.StatusFound, echo.Map{
		"status"  : http.StatusFound,
		"message" : "Usuario encontrado",
		"data"    : user,
	})
}

// Crea un nuevo inventario
func (h *Handler) CreateUser(c echo.Context) error {
	// Valida la conexion a la coleccion
	if h.Users == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound,
			"message" : "Sin conexion a la colección",
			"data"    : nil,
		})
	}

	var user models.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "Input invalido",
			"data"    : nil,
		})
	}

	if strings.TrimSpace(user.Name) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El nombre es obligatorio",
			"data"    : nil,
		})
	}

	if strings.TrimSpace(user.Email) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El correo electronico es obligatorio",
			"data"    : nil,
		})
	}

	_, err := h.Users.InsertOne(context.Background(), user)
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
		"data"	  : user,
	})
}

// Actualiza un usuario existente
func (h *Handler) UpdateUser(c echo.Context) error {
	// Valida la conexion a la coleccion
	if h.Users == nil {
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

	// Valida la conexion a la coleccion
	if h.Users == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound, 
			"message" : "Sin conexion a la colección",
			"data"    : nil,
		})
	}

	var user models.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "Input invalido",
			"data"	  : nil,
		})
	}

	if strings.TrimSpace(user.Name) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El nombre es obligatorio",
			"data"    : nil,
		})
	}

	if strings.TrimSpace(user.Email) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El correo electronico es obligatorio",
			"data"    : nil,
		})
	}

	// Prepara filtro y documento de actualización
    filter := bson.M{"_id": id}
    update := bson.M{
        "$set": bson.M{
            "name"  : user.Name,
            "email" : user.Email,
        },
    }

	// Actualiza el documento
    res, err := h.Users.UpdateOne(context.Background(), filter, update)
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
            "message" : "Usuario no encontrado",
            "data"    : nil,
        })
    }

	return c.JSON(http.StatusCreated, echo.Map{
        "status"  : http.StatusCreated,
        "message" : "Usuario actualizado exitosamente",
        "data"    : user,
    })
}

// Elimina un inventario
func (h *Handler) DeleteUser(c echo.Context) error {
	// Valida la conexion a la coleccion
	if h.Users == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status": http.StatusNotFound, 
			"message": "Sin conexion a la colección",
			"data": nil,
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
			"data" 	  : nil,
		})
	}

	// Realiza operacion de eliminado mediante el id recuperado del parametro de consulta
	res, err := h.Users.DeleteOne(context.Background(), bson.M{"_id": id})
	// Valida si la operacion fue exitosa
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status"  : http.StatusInternalServerError, 
			"message" : err.Error(),
			"data" 	  : nil,
		})
	}
	
	// Valida si se elimino algun documento
	if res.DeletedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound, 
			"message" : "Usuario no encontrado",
			"data" 	  : nil,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
        "status"  : http.StatusOK,
        "message" : "Usuario eliminado exitosamente",
        "data"    : nil,
    })
}