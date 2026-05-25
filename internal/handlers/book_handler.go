package handlers

import (
	"net/http"
	"net/url"
	"strconv"

	"ebook-system/internal/models"
	"ebook-system/internal/services"
)

type BookHandler struct {
	service *services.BookService
}

func NewBookHandler(service *services.BookService) *BookHandler {
	return &BookHandler{service: service}
}

func (h *BookHandler) Home(w http.ResponseWriter, r *http.Request) {
	render(w, r, "home", nil)
}

func (h *BookHandler) Catalog(w http.ResponseWriter, r *http.Request) {
	books, err := h.service.ListBooks(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener libros", http.StatusInternalServerError)
		return
	}
	render(w, r, "catalog", books)
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/catalog?error="+url.QueryEscape("Método no permitido"), http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/catalog?error="+url.QueryEscape("Datos inválidos"), http.StatusSeeOther)
		return
	}

	year, _ := strconv.Atoi(r.FormValue("published_year"))
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)

	book := &models.Book{
		Title:         r.FormValue("title"),
		Author:        r.FormValue("author"),
		PublishedYear: year,
		ISBN:          r.FormValue("isbn"),
		Price:         price,
	}

	err = h.service.AddBook(r.Context(), book)
	if err != nil {
		http.Redirect(w, r, "/catalog?error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/catalog", http.StatusSeeOther)
}
