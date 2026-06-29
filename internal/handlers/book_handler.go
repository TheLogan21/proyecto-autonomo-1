package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"ebook-system/internal/models"
	"ebook-system/internal/services"
)

type BookHandler struct {
	service services.BookService
}

func NewBookHandler(service services.BookService) *BookHandler {
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

func (h *BookHandler) ListBooksJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	books, err := h.service.ListBooks(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al obtener libros"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}

// CreateBookJSON procesa la creación de un nuevo e-book mediante el servicio web.
// Cumple con la serialización JSON obligatoria de la Unidad 4 y maneja errores semánticos.
func (h *BookHandler) CreateBookJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Método no permitido"})
		return
	}

	var book models.Book
	// Deserialización de los datos binarios entrantes (JSON) hacia la estructura interna del modelo
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cuerpo JSON inválido o malformado"})
		return
	}

	// VALIDACIÓN TÉCNICA OBLIGATORIA DE ENTRADA
	if book.Title == "" || book.Author == "" || book.Price <= 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"error": "Datos del libro incompletos o precio inválido"})
		return
	}

	// Ejecución de la persistencia interactuando con la interfaz desacoplada de la lógica de negocio
	err := h.service.AddBook(r.Context(), &book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Serialización de salida con confirmación del recurso creado de manera exitosa
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}


