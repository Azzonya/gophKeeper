// Package service implements the business logic for managing data items,
// coordinating between the database and S3 storage repositories.
package service

import (
	"context"
	"fmt"
	"gophKeeper/server/internal/domain/dataItems/model"
)

// Service provides methods to manage data items, handling both database operations
// and S3 file storage interactions based on the type of data being processed.
type Service struct {
	repoDB RepoDBI
	repoS3 RepoS3
}

// New creates a new Service instance with the given database and S3 repositories.
func New(repoDB RepoDBI, repoS3 RepoS3) *Service {
	return &Service{
		repoDB: repoDB,
		repoS3: repoS3,
	}
}

// List retrieves data items based on the provided filtering parameters.
// It delegates the operation to the database repository.
func (s *Service) List(ctx context.Context, pars *model.ListPars) ([]*model.Main, int64, error) {
	return s.repoDB.List(ctx, pars)
}

// Create stores a new data item in the database and, if the item is of binary type,
// uploads the binary data to S3 and updates the database with the file's URL.
func (s *Service) Create(ctx context.Context, obj *model.Edit) error {
	tx, err := s.repoDB.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer s.repoDB.HandleTxCompletion(tx, &err)

	err = s.repoDB.Create(ctx, obj)
	if err != nil {
		return err
	}

	if *obj.Type == model.BinaryDataType {
		var url string
		url, err = s.repoS3.UploadFile(ctx, obj.ID, *obj.Data)
		if err != nil {
			return err
		}

		err = s.Update(ctx, &model.GetPars{
			ID: obj.ID,
		}, &model.Edit{
			URL: &url,
		})
		if err != nil {
			_ = s.repoS3.DeleteFile(ctx, &model.GetPars{
				ID: obj.ID,
			})
			return err
		}
	}

	return nil
}

// Get retrieves a data item from the database and, if it is of binary type,
// fetches the associated file from S3 and returns it as part of the response.
func (s *Service) Get(ctx context.Context, pars *model.GetPars) (*model.Main, bool, error) {
	obj, found, err := s.repoDB.Get(ctx, pars)
	if err != nil {
		return nil, false, err
	}
	if !found {
		return nil, false, nil
	}

	if obj.Type == model.BinaryDataType {
		file, found, err := s.repoS3.GetFile(ctx, &model.GetPars{ID: obj.ID})
		if err != nil {
			return nil, false, err
		}
		if !found {
			return nil, false, nil
		}

		obj.Data = file
	}

	return obj, found, nil
}

// Update modifies an existing data item in the database. If the item is of binary type
// and contains updated data, it uploads the new data to S3 and updates the item's URL.
func (s *Service) Update(ctx context.Context, pars *model.GetPars, obj *model.Edit) error {
	existingObj, found, err := s.repoDB.Get(ctx, pars)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("record not found")
	}

	if existingObj.Type == model.BinaryDataType && obj.Data != nil {
		url, err := s.repoS3.UploadFile(ctx, pars.ID, *obj.Data)
		if err != nil {
			return err
		}
		obj.URL = &url
	}

	return s.repoDB.Update(ctx, pars, obj)
}

// Delete removes a data item from the database. If the item is of binary type,
// it also deletes the associated file from S3.
func (s *Service) Delete(ctx context.Context, pars *model.GetPars) error {
	existingObj, found, err := s.repoDB.Get(ctx, pars)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("record not found")
	}

	if existingObj.Type == model.BinaryDataType {
		err = s.repoS3.DeleteFile(ctx, &model.GetPars{ID: existingObj.ID})
		if err != nil {
			return err
		}
	}

	return s.repoDB.Delete(ctx, pars)
}
