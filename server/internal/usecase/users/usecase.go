// Package users implements the business logic for user management,
// coordinating operations like registration, login, and authentication
// through the service layer.
package users

import (
	"context"
	"gophKeeper/server/internal/domain/users/model"
	"gophKeeper/server/internal/errs"
)

// Usecase provides the business logic for managing users and handling
// authentication, using the user and authentication services to perform operations.
type Usecase struct {
	usersService UsersServiceI
	authService  AuthServiceI
}

// New creates a new Usecase instance with the provided user and authentication services.
func New(usersService UsersServiceI, authService AuthServiceI) *Usecase {
	return &Usecase{
		usersService: usersService,
		authService:  authService,
	}
}

// Register registers a new user by checking if the username is available,
// hashing the password, and creating the user in the database.
func (u *Usecase) Register(ctx context.Context, username string, password string) error {
	if username == "" || password == "" {
		return errs.InvalidInput
	}

	taken, err := u.usersService.IsLoginTaken(ctx, username)
	if err != nil {
		return err
	}
	if taken {
		return errs.UsernameAlreadyExists
	}

	passwordHash, err := u.usersService.HashPassword(password)
	if err != nil {
		return err
	}

	err = u.usersService.Create(ctx, &model.Edit{
		Username:     &username,
		PasswordHash: &passwordHash,
	})
	if err != nil {
		return err
	}

	createdUser, _, err := u.usersService.Get(ctx, &model.GetPars{
		Username: username,
	})
	if err != nil {
		return err
	}

	if createdUser == nil {
		return errs.UserNotFound
	}

	return nil
}

// Login handles user login by validating the username and password,
// and generating a JWT token if the credentials are correct.
func (u *Usecase) Login(ctx context.Context, username string, password string) (*string, error) {
	if username == "" || password == "" {
		return nil, errs.InvalidInput
	}

	user, found, err := u.usersService.Get(ctx, &model.GetPars{
		Username: username,
	})
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errs.UserNotFound
	}

	isValidPassword := u.usersService.IsValidPassword(user.PasswordHash, password)
	if !isValidPassword {
		return nil, errs.InvalidPassword
	}

	token, err := u.authService.CreateToken(user)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// GetUserIDFromContext extracts the user ID from the context,
// using the authentication service.
func (u *Usecase) GetUserIDFromContext(ctx context.Context) (string, error) {
	return u.authService.GetUserIDFromContext(ctx)
}
