package repository

import (
	"context"
	"ebook-system/internal/models"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PurchaseRepository struct {
	db *pgxpool.Pool
}

func NewPurchaseRepository(db *pgxpool.Pool) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

func (r *PurchaseRepository) CreatePurchase(ctx context.Context, userID, bookID int) error {
	query := `INSERT INTO user_books (user_id, book_id) VALUES ($1, $2)`
	_, err := r.db.Exec(ctx, query, userID, bookID)
	return err
}

func (r *PurchaseRepository) HasPurchased(ctx context.Context, userID, bookID int) (bool, error) {
	query := `SELECT 1 FROM user_books WHERE user_id = $1 AND book_id = $2`
	var exists int
	err := r.db.QueryRow(ctx, query, userID, bookID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *PurchaseRepository) IncrementDownload(ctx context.Context, userID, bookID int) error {
	query := `UPDATE user_books SET download_count = download_count + 1 WHERE user_id = $1 AND book_id = $2`
	_, err := r.db.Exec(ctx, query, userID, bookID)
	return err
}

func (r *PurchaseRepository) GetUserPurchases(ctx context.Context, userID int) ([]models.UserBook, error) {
	query := `
		SELECT ub.user_id, ub.book_id, ub.download_count, ub.purchased_at,
		       b.id, b.title, b.author, b.published_year, b.isbn, b.price, b.created_at
		FROM user_books ub
		JOIN books b ON ub.book_id = b.id
		WHERE ub.user_id = $1
		ORDER BY ub.purchased_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying purchases: %w", err)
	}
	defer rows.Close()

	var purchases []models.UserBook
	for rows.Next() {
		var ub models.UserBook
		var b models.Book
		if err := rows.Scan(
			&ub.UserID, &ub.BookID, &ub.DownloadCount, &ub.PurchasedAt,
			&b.ID, &b.Title, &b.Author, &b.PublishedYear, &b.ISBN, &b.Price, &b.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning purchase: %w", err)
		}
		ub.Book = &b
		purchases = append(purchases, ub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating purchases: %w", err)
	}

	return purchases, nil
}
