// Package service defines the interfaces for interacting with the user repository,
// providing methods for managing user accounts such as retrieving, listing, creating,
// updating, deleting, and checking user existence.
package service

import (
	"context"
	"gophKeeper/server/internal/domain/users/model"
)

// RepoDBI defines the interface for database interactions related to user accounts.
// It includes methods for retrieving, listing, creating, updating, deleting, and checking
// the existence of users in the database.
type RepoDBI interface {
	Get(ctx context.Context, pars *model.GetPars) (*model.Main, bool, error)
	List(ctx context.Context, pars *model.ListPars) ([]*model.Main, int64, error)
	Create(ctx context.Context, obj *model.Edit) error
	Update(ctx context.Context, pars *model.GetPars, obj *model.Edit) error
	Delete(ctx context.Context, pars *model.GetPars) error
	Exists(ctx context.Context, pars *model.GetPars) (bool, error)
}
