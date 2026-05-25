package repository

import (
	"context"
	"ebook-system/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepository struct {
	db *pgxpool.Pool
}

func NewBookRepository(db *pgxpool.Pool) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) Create(ctx context.Context, book *models.Book) error {
	query := `
		INSERT INTO books (title, author, published_year, isbn, price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query, book.Title, book.Author, book.PublishedYear, book.ISBN, book.Price).
		Scan(&book.ID, &book.CreatedAt)
	if err != nil {
		return fmt.Errorf("error creating book: %w", err)
	}

	return nil
}

func (r *BookRepository) GetAll(ctx context.Context) ([]models.Book, error) {
	query := `SELECT id, title, author, published_year, isbn, price, created_at FROM books ORDER BY id DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying books: %w", err)
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.PublishedYear, &b.ISBN, &b.Price, &b.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning book: %w", err)
		}
		books = append(books, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating books: %w", err)
	}

	return books, nil
}

func (r *BookRepository) GetByID(ctx context.Context, id int) (*models.Book, error) {
	query := `SELECT id, title, author, published_year, isbn, price, created_at FROM books WHERE id = $1`
	var b models.Book
	err := r.db.QueryRow(ctx, query, id).Scan(&b.ID, &b.Title, &b.Author, &b.PublishedYear, &b.ISBN, &b.Price, &b.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error querying book by id: %w", err)
	}
	return &b, nil
}
