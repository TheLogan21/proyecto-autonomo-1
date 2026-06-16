package services

import (
	"context"
	"ebook-system/internal/models"
	"errors"
	"fmt"
)

// PurchaseService gestiona la compra y descarga de libros
type PurchaseService struct {
	purchaseRepo PurchaseRepository
	userRepo     UserRepository
	bookRepo     BookRepository
}

// NewPurchaseService crea un PurchaseService
func NewPurchaseService(pr PurchaseRepository, ur UserRepository, br BookRepository) *PurchaseService {
	return &PurchaseService{
		purchaseRepo: pr,
		userRepo:     ur,
		bookRepo:     br,
	}
}

// BuyBook procesa la compra de un libro
func (s *PurchaseService) BuyBook(ctx context.Context, userID, bookID int) error {
	hasPurchased, err := s.purchaseRepo.HasPurchased(ctx, userID, bookID)
	if err != nil {
		return fmt.Errorf("error al verificar libro duplicado: %w", err)
	}
	if hasPurchased {
		return errors.New("ya posees este libro")
	}

	book, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("error al obtener libro para compra: %w", err)
	}
	if book == nil {
		return errors.New("libro no encontrado")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("error al obtener usuario para compra: %w", err)
	}
	if user == nil {
		return errors.New("usuario no encontrado")
	}

	if user.Balance < book.Price {
		return errors.New("saldo insuficiente")
	}

	newBalance := user.Balance - book.Price
	err = s.userRepo.UpdateBalance(ctx, userID, newBalance)
	if err != nil {
		return fmt.Errorf("error al actualizar saldo: %w", err)
	}

	err = s.purchaseRepo.CreatePurchase(ctx, userID, bookID)
	if err != nil {
		// Restaura el saldo anterior en caso de fallo
		_ = s.userRepo.UpdateBalance(ctx, userID, user.Balance)
		return fmt.Errorf("error al registrar compra: %w", err)
	}

	return nil
}

// GetUserPurchases obtiene los libros comprados por un usuario
func (s *PurchaseService) GetUserPurchases(ctx context.Context, userID int) ([]models.UserBook, error) {
	purchases, err := s.purchaseRepo.GetUserPurchases(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener compras del usuario desde el repositorio: %w", err)
	}
	return purchases, nil
}

// DownloadBook descarga un libro comprado
func (s *PurchaseService) DownloadBook(ctx context.Context, userID, bookID int) ([]byte, error) {
	hasPurchased, err := s.purchaseRepo.HasPurchased(ctx, userID, bookID)
	if err != nil {
		return nil, fmt.Errorf("error al verificar compra para descarga: %w", err)
	}
	if !hasPurchased {
		return nil, errors.New("no has comprado este libro")
	}

	err = s.purchaseRepo.IncrementDownload(ctx, userID, bookID)
	if err != nil {
		return nil, fmt.Errorf("error al incrementar descargas: %w", err)
	}

	content := []byte("esto es una prueba no esperemos mucho")
	return content, nil
}

