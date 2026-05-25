package handlers

import (
	"fmt"
	"html/template"
	"net/http"
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

func (h *BookHandler) render(w http.ResponseWriter, page string, data interface{}) {
	t, err := template.ParseFiles("ui/html/base.tmpl", fmt.Sprintf("ui/html/pages/%s.tmpl", page))
	if err != nil {
		http.Error(w, "Error parsing templates", http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func (h *BookHandler) Home(w http.ResponseWriter, r *http.Request) {
	h.render(w, "home", nil)
}

func (h *BookHandler) Catalog(w http.ResponseWriter, r *http.Request) {
	books, err := h.service.ListBooks(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener libros", http.StatusInternalServerError)
		return
	}
	h.render(w, "catalog", books)
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	year, _ := strconv.Atoi(r.FormValue("published_year"))
	book := &models.Book{
		Title:         r.FormValue("title"),
		Author:        r.FormValue("author"),
		PublishedYear: year,
		ISBN:          r.FormValue("isbn"),
	}

	err = h.service.AddBook(r.Context(), book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/catalog", http.StatusSeeOther)
}
