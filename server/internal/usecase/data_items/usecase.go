// Package data_items implements the use case logic for managing data items,
// coordinating operations like retrieval, creation, updating, and deletion
// through the service layer.
package data_items

import (
	"context"
	"gophKeeper/server/internal/domain/data_items/model"
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

// GetData retrieves a data item based on the provided query parameters.
func (u *Usecase) GetData(ctx context.Context, obj *model.GetPars) (*model.Main, bool, error) {
	return u.dataItemsService.Get(ctx, obj)
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
