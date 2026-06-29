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

// BookService define la lógica de negocio para libros
type BookService interface {
	AddBook(ctx context.Context, book *models.Book) error
	ListBooks(ctx context.Context) ([]models.Book, error)
	GetBookByID(ctx context.Context, id int) (*models.Book, error)
}

// UserService define la lógica de negocio para usuarios
type UserService interface {
	Register(ctx context.Context, username, email, password string) error
	Login(ctx context.Context, email, password string) (*models.User, error)
	AddBalance(ctx context.Context, userID int, amount float64) error
	GetUserByID(ctx context.Context, id int) (*models.User, error)
}

// PurchaseService define la lógica de negocio para transacciones de libros
type PurchaseService interface {
	BuyBook(ctx context.Context, userID, bookID int) error
	GetUserPurchases(ctx context.Context, userID int) ([]models.UserBook, error)
	DownloadBook(ctx context.Context, userID, bookID int) ([]byte, error)
}

