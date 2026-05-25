package services

import (
	"context"
	"ebook-system/internal/models"
	"ebook-system/internal/repository"
	"errors"
)

type BookService struct {
	repo *repository.BookRepository
}

func NewBookService(repo *repository.BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) AddBook(ctx context.Context, book *models.Book) error {
	if book.Title == "" || book.Author == "" {
		return errors.New("título y autor son obligatorios")
	}
	if book.PublishedYear <= 0 {
		return errors.New("año de publicación inválido")
	}

	return s.repo.Create(ctx, book)
}

func (s *BookService) ListBooks(ctx context.Context) ([]models.Book, error) {
	return s.repo.GetAll(ctx)
}
