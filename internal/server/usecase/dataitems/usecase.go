// Package dataitems implements the use case logic for managing data items,
// coordinating operations like retrieval, creation, updating, and deletion
// through the service layer.
package dataitems

import (
	"context"
	"gophKeeper/internal/server/domain/dataitems/model"
)

// Usecase provides the business logic for managing data items,
// leveraging a data items service interface to perform operations.
type Usecase struct {
	dataItemsService DataItemsServiceI
}

// New creates a new Usecase instance with the provided data items service.
func New(dataItemsService DataItemsServiceI) *Usecase {
	return &Usecase{
		dataItemsService: dataItemsService,
	}
}

// DataItemsServiceI defines the interface for the data items service,
// providing methods to manage data items, including listing, creating,
// retrieving, updating, and deleting operations.
type DataItemsServiceI interface {
	List(ctx context.Context, pars *model.ListPars) ([]*model.DataItems, int64, error)
	Create(ctx context.Context, obj *model.Edit) error
	Get(ctx context.Context, pars *model.GetPars) (*model.DataItems, bool, error)
	Update(ctx context.Context, pars *model.GetPars, obj *model.Edit) error
	Delete(ctx context.Context, pars *model.GetPars) error
}

// GetData retrieves a data item based on the provided query parameters.
func (u *Usecase) GetData(ctx context.Context, obj *model.GetPars) (*model.DataItems, bool, error) {
	return u.dataItemsService.Get(ctx, obj)
}

// ListAll retrieves the all data item based on the provided user id
func (u *Usecase) ListAll(ctx context.Context, obj *model.ListPars) ([]*model.DataItems, int64, error) {
	return u.dataItemsService.List(ctx, obj)
}

// CreateData creates a new data item using the provided model.Edit object.
func (u *Usecase) CreateData(ctx context.Context, obj *model.Edit) error {
	return u.dataItemsService.Create(ctx, obj)
}

// EditData updates an existing data item identified by the provided model.Edit object.
func (u *Usecase) EditData(ctx context.Context, obj *model.Edit) error {
	return u.dataItemsService.Update(ctx, &model.GetPars{
		ID: obj.ID,
	}, obj)
}

// DeleteData deletes a data item based on the provided query parameters.
func (u *Usecase) DeleteData(ctx context.Context, obj *model.GetPars) error {
	return u.dataItemsService.Delete(ctx, obj)
}
