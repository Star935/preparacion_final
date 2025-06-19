package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "backend/handlers"
    "backend/models"
    "github.com/labstack/echo/v4"
)

func TestCreateBookValidationTitleEmpty(t *testing.T) {
    e := echo.New()
    h := &handlers.Handler{Books: nil}

    book := models.Book{Title: "", Author: "Title empty", Isbn: "Title is empty", Availability: 1}
    body, _ := json.Marshal(book)
    req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    if err := h.CreateBook(c); err != nil {
        t.Fatal(err)
    }
    if rec.Code != http.StatusBadRequest {
        t.Errorf("Esperado 400, obtuvo %d", rec.Code)
    }
}

func TestCreateBookValidationAuthorEmpty(t *testing.T) {
    e := echo.New()
    h := &handlers.Handler{Books: nil}

    book := models.Book{Title: "Author empty", Author: "", Isbn: "Author is empty", Availability: 2}
    body, _ := json.Marshal(book)
    req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    if err := h.CreateBook(c); err != nil {
        t.Fatal(err)
    }
    if rec.Code != http.StatusBadRequest {
        t.Errorf("Esperado 400, obtuvo %d", rec.Code)
    }
}
