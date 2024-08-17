// Package users defines the service interfaces for managing user accounts,
// handling operations like authentication, password management, and CRUD operations.
package users

import (
	"context"
	"gophKeeper/server/internal/domain/users/model"
)

// UsersServiceI defines the interface for user management operations,
// including password validation, user creation, updating, deletion, and
// checking if a username is already taken.
type UsersServiceI interface {
	IsValidPassword(password string, plainPassword string) bool
	HashPassword(password string) (string, error)
	IsLoginTaken(ctx context.Context, username string) (bool, error)
	List(ctx context.Context, pars *model.ListPars) ([]*model.Main, int64, error)
	Create(ctx context.Context, obj *model.Edit) error
	Get(ctx context.Context, pars *model.GetPars) (*model.Main, bool, error)
	Update(ctx context.Context, pars *model.GetPars, obj *model.Edit) error
	Delete(ctx context.Context, pars *model.GetPars) error
	Exists(ctx context.Context, pars *model.GetPars) (bool, error)
}

// AuthServiceI defines the interface for authentication operations,
// including extracting the user ID from context and creating JWT tokens.
type AuthServiceI interface {
	GetUserIDFromContext(ctx context.Context) (string, error)
	CreateToken(u *model.Main) (string, error)
}
