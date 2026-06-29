package services

import (
	"context"
	"ebook-system/internal/models"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// userService gestiona usuarios
type userService struct {
	repo UserRepository
}

// NewUserService crea un UserService
func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

// Register registra un nuevo usuario
func (s *userService) Register(ctx context.Context, username, email, password string) error {
	if username == "" || email == "" || password == "" {
		return errors.New("todos los campos son requeridos")
	}

	existing, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("error al verificar email existente: %w", err)
	}
	if existing != nil {
		return errors.New("el email ya está registrado")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error al encriptar contraseña: %w", err)
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("error al crear el usuario en el repositorio: %w", err)
	}
	return nil
}

// Login valida las credenciales de un usuario
func (s *userService) Login(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario para login: %w", err)
	}
	if user == nil {
		return nil, errors.New("credenciales inválidas")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	return user, nil
}

// AddBalance añade saldo a un usuario
func (s *userService) AddBalance(ctx context.Context, userID int, amount float64) error {
	if amount <= 0 {
		return errors.New("la cantidad debe ser mayor a cero")
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("error al obtener usuario para añadir saldo: %w", err)
	}
	if user == nil {
		return errors.New("usuario no encontrado")
	}

	newBalance := user.Balance + amount
	err = s.repo.UpdateBalance(ctx, userID, newBalance)
	if err != nil {
		return fmt.Errorf("error al actualizar saldo del usuario en el repositorio: %w", err)
	}
	return nil
}

// GetUserByID obtiene un usuario por ID
func (s *userService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario por ID %d del repositorio: %w", id, err)
	}
	return user, nil
}


