// Package service defines the interfaces for interacting with different repositories,
// including a database repository and an S3-compatible file storage repository.
package service

import (
	"context"
	"github.com/jackc/pgx/v5"
	"gophKeeper/server/internal/domain/dataitems/model"
)

// RepoDBI outlines the methods for interacting with the database repository,
// including operations to get, list, create, update, and delete data items.
type RepoDBI interface {
	Get(ctx context.Context, pars *model.GetPars) (*model.Main, bool, error)
	List(ctx context.Context, pars *model.ListPars) ([]*model.Main, int64, error)
	Create(ctx context.Context, obj *model.Edit) error
	Update(ctx context.Context, pars *model.GetPars, obj *model.Edit) error
	Delete(ctx context.Context, pars *model.GetPars) error
	BeginTx(ctx context.Context) (pgx.Tx, error)
	CommitTx(ctx context.Context, tx pgx.Tx) error
	RollbackTx(ctx context.Context, tx pgx.Tx) error
	HandleTxCompletion(tx pgx.Tx, err *error)
}

// RepoS3 defines the methods for interacting with an S3-compatible storage,
// including operations to get, upload, and delete files.
type RepoS3 interface {
	GetFile(ctx context.Context, pars *model.GetPars) ([]byte, bool, error)
	UploadFile(ctx context.Context, id string, data []byte) (string, error)
	DeleteFile(ctx context.Context, pars *model.GetPars) error
}
