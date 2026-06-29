package services

import (
	"context"
	"ebook-system/internal/models"
	"errors"
	"fmt"
)

// bookService gestiona las operaciones de libros
type bookService struct {
	repo BookRepository
}

// NewBookService crea un BookService
func NewBookService(repo BookRepository) BookService {
	return &bookService{repo: repo}
}

// AddBook registra un nuevo libro
func (s *bookService) AddBook(ctx context.Context, book *models.Book) error {
	if book.Title == "" || book.Author == "" {
		return errors.New("título y autor son obligatorios")
	}
	if book.PublishedYear <= 0 {
		return errors.New("año de publicación inválido")
	}
	if book.Price < 0 {
		return errors.New("el precio no puede ser negativo")
	}

	err := s.repo.Create(ctx, book)
	if err != nil {
		return fmt.Errorf("error al agregar libro en el repositorio: %w", err)
	}
	return nil
}

// ListBooks obtiene la lista de libros
func (s *bookService) ListBooks(ctx context.Context) ([]models.Book, error) {
	books, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al listar libros del repositorio: %w", err)
	}
	return books, nil
}

// GetBookByID obtiene un libro por ID
func (s *bookService) GetBookByID(ctx context.Context, id int) (*models.Book, error) {
	book, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener libro por ID %d del repositorio: %w", id, err)
	}
	return book, nil
}


