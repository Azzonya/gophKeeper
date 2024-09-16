// Package service implements the business logic for managing user accounts,
// including operations like password hashing, user validation, and coordinating
// interactions with the database repository.
package service

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gophKeeper/internal/server/domain/users/model"
)

// Service provides methods to manage user accounts, handling password operations,
// user validation, and CRUD operations through the repository interface.
type Service struct {
	repoDB RepoDBI
}

// New creates a new Service instance with the given database repository.
func New(repoDB RepoDBI) *Service {
	return &Service{
		repoDB: repoDB,
	}
}

// RepoDBI defines the interface for database interactions related to user accounts.
// It includes methods for retrieving, listing, creating, updating, deleting, and checking
// the existence of users in the database.
type RepoDBI interface {
	Get(ctx context.Context, pars *model.GetPars) (*model.User, bool, error)
	List(ctx context.Context, pars *model.ListPars) ([]*model.User, int64, error)
	Create(ctx context.Context, obj *model.Edit) error
	Update(ctx context.Context, pars *model.GetPars, obj *model.Edit) error
	Delete(ctx context.Context, pars *model.GetPars) error
	Exists(ctx context.Context, pars *model.GetPars) (bool, error)
}

// IsValidPassword compares a hashed password with a plain password to verify a match.
func (s *Service) IsValidPassword(password string, plainPassword string) bool {
	// Сравниваем хэшированный пароль из базы данных с переданным паролем
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(plainPassword))
	return err == nil
}

// HashPassword hashes a plain password using bcrypt before storing it in the database.
func (s *Service) HashPassword(password string) (string, error) {
	// Хэшируем пароль перед сохранением в базу данных
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("can not hash password - %w", err)
	}

	return string(hashedPassword), nil
}

// IsLoginTaken checks if a username is already taken by querying the database.
func (s *Service) IsLoginTaken(ctx context.Context, username string) (bool, error) {
	return s.Exists(ctx, &model.GetPars{Username: username})
}

// List retrieves a list of users based on the provided filtering parameters,
// delegating the operation to the database repository.
func (s *Service) List(ctx context.Context, pars *model.ListPars) ([]*model.User, int64, error) {
	return s.repoDB.List(ctx, pars)
}

// Create stores a new user account in the database.
func (s *Service) Create(ctx context.Context, obj *model.Edit) error {
	return s.repoDB.Create(ctx, obj)
}

// Get retrieves a user account from the database based on the provided query parameters.
func (s *Service) Get(ctx context.Context, pars *model.GetPars) (*model.User, bool, error) {
	return s.repoDB.Get(ctx, pars)
}

// Update modifies an existing user account in the database.
func (s *Service) Update(ctx context.Context, pars *model.GetPars, obj *model.Edit) error {
	return s.repoDB.Update(ctx, pars, obj)
}

// Delete removes a user account from the database.
func (s *Service) Delete(ctx context.Context, pars *model.GetPars) error {
	return s.repoDB.Delete(ctx, pars)
}

// Exists checks whether a user account exists in the database based on the provided query parameters.
func (s *Service) Exists(ctx context.Context, pars *model.GetPars) (bool, error) {
	return s.repoDB.Exists(ctx, pars)
}
