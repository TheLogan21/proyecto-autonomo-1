package services

import (
	"context"
	"ebook-system/internal/models"
)

// UserRepository define el almacenamiento de usuarios
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	UpdateBalance(ctx context.Context, id int, newBalance float64) error
}

// BookRepository define el almacenamiento de libros
type BookRepository interface {
	Create(ctx context.Context, book *models.Book) error
	GetAll(ctx context.Context) ([]models.Book, error)
	GetByID(ctx context.Context, id int) (*models.Book, error)
}

// PurchaseRepository define el almacenamiento de compras
type PurchaseRepository interface {
	CreatePurchase(ctx context.Context, userID, bookID int) error
	HasPurchased(ctx context.Context, userID, bookID int) (bool, error)
	IncrementDownload(ctx context.Context, userID, bookID int) error
	GetUserPurchases(ctx context.Context, userID int) ([]models.UserBook, error)
}
