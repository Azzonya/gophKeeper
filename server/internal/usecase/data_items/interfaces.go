// Package data_items defines the service interface for managing data items,
// including operations for listing, creating, retrieving, updating, and deleting data items.
package data_items

import (
	"context"
	"gophKeeper/server/internal/domain/data_items/model"
)

// DataItemsServiceI defines the interface for the data items service,
// providing methods to manage data items, including listing, creating,
// retrieving, updating, and deleting operations.
type DataItemsServiceI interface {
	List(ctx context.Context, pars *model.ListPars) ([]*model.Main, int64, error)
	Create(ctx context.Context, obj *model.Edit) error
	Get(ctx context.Context, pars *model.GetPars) (*model.Main, bool, error)
	Update(ctx context.Context, pars *model.GetPars, obj *model.Edit) error
	Delete(ctx context.Context, pars *model.GetPars) error
}
