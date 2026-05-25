package services

import (
	"context"
	"ebook-system/internal/models"
	"ebook-system/internal/repository"
	"errors"
	"fmt"
)

type PurchaseService struct {
	purchaseRepo *repository.PurchaseRepository
	userRepo     *repository.UserRepository
	bookRepo     *repository.BookRepository
}

func NewPurchaseService(pr *repository.PurchaseRepository, ur *repository.UserRepository, br *repository.BookRepository) *PurchaseService {
	return &PurchaseService{
		purchaseRepo: pr,
		userRepo:     ur,
		bookRepo:     br,
	}
}

func (s *PurchaseService) BuyBook(ctx context.Context, userID, bookID int) error {
	hasPurchased, err := s.purchaseRepo.HasPurchased(ctx, userID, bookID)
	if err != nil {
		return err
	}
	if hasPurchased {
		return errors.New("ya posees este libro")
	}

	book, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return err
	}
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
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
		return fmt.Errorf("error al actualizar saldo: %v", err)
	}

	err = s.purchaseRepo.CreatePurchase(ctx, userID, bookID)
	if err != nil {
		s.userRepo.UpdateBalance(ctx, userID, user.Balance)
		return fmt.Errorf("error al registrar compra: %v", err)
	}

	return nil
}

func (s *PurchaseService) GetUserPurchases(ctx context.Context, userID int) ([]models.UserBook, error) {
	return s.purchaseRepo.GetUserPurchases(ctx, userID)
}

func (s *PurchaseService) DownloadBook(ctx context.Context, userID, bookID int) ([]byte, error) {
	hasPurchased, err := s.purchaseRepo.HasPurchased(ctx, userID, bookID)
	if err != nil {
		return nil, err
	}
	if !hasPurchased {
		return nil, errors.New("no has comprado este libro")
	}

	err = s.purchaseRepo.IncrementDownload(ctx, userID, bookID)
	if err != nil {
		return nil, err
	}

	content := []byte("esto es una prueba no esperemos mucho")
	return content, nil
}
