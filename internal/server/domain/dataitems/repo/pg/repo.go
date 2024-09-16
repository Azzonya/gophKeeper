// Package pg provides a PostgreSQL-based implementation for managing data items,
// including operations such as retrieving, listing, creating, updating, and deleting records.
package pg

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gophKeeper/internal/server/domain/dataitems/model"
	"gophKeeper/internal/server/errs"
)

// Repo provides methods to interact with the PostgreSQL database for data item operations.
// It holds a connection pool to manage database connections.
type Repo struct {
	Con *pgxpool.Pool
}

// New creates a new instance of Repo with the given PostgreSQL connection pool.
func New(con *pgxpool.Pool) *Repo {
	return &Repo{
		con,
	}
}

// Get retrieves a single data item based on the provided query parameters.
// It returns the item if found, a boolean indicating its existence, and any error encountered.
func (r *Repo) Get(ctx context.Context, pars *model.GetPars) (*model.DataItems, bool, error) {
	if !pars.IsValid() {
		return nil, false, errs.InvalidInput
	}

	var result model.DataItems

	queryBuilder := squirrel.Select("*").From("data_items")

	if len(pars.ID) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.ID})
	}

	if len(pars.UserID) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"user_id": pars.UserID})
	}

	if len(pars.Type) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"type": pars.Type})
	}

	queryBuilder = queryBuilder.Limit(1)

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, false, err
	}

	err = r.Con.QueryRow(ctx, sql, args...).Scan(&result.ID, &result.UserID, &result.Type, &result.Data, &result.Meta, &result.URL, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &result, true, nil
}

// List retrieves multiple data items based on the provided query parameters,
// supporting filters like ID, user ID, type, and timestamps. It returns the list
// of items, the total count, and any error encountered.
func (r *Repo) List(ctx context.Context, pars *model.ListPars) ([]*model.DataItems, int64, error) {
	queryBuilder := squirrel.
		Select("id", "user_id", "type", "data", "created_at", "updated_at").
		From("data_items")

	if pars.ID != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.ID})
	}

	if pars.IDs != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.IDs})
	}

	if pars.UserID != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"user_id": pars.UserID})
	}

	if pars.UserIDs != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"user_id": pars.UserIDs})
	}

	if pars.Type != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"type": pars.Type})
	}

	if pars.CreatedBefore != nil {
		queryBuilder = queryBuilder.Where(squirrel.LtOrEq{"created_at": pars.CreatedBefore})
	}

	if pars.CreatedAfter != nil {
		queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"created_at": pars.CreatedAfter})
	}

	if pars.UpdatedBefore != nil {
		queryBuilder = queryBuilder.Where(squirrel.LtOrEq{"updated_at": pars.UpdatedBefore})
	}

	if pars.UpdatedAfter != nil {
		queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"updated_at": pars.UpdatedAfter})
	}

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.Con.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	var result []*model.DataItems
	for rows.Next() {
		var data model.DataItems
		err = rows.Scan(&data.ID, &data.UserID, &data.Type, &data.Data, &data.CreatedAt, &data.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, &data)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return result, int64(len(result)), nil
}

// Create inserts a new data item into the database based on the provided Edit object,
// returning the ID of the newly created item and any error encountered.
func (r *Repo) Create(ctx context.Context, obj *model.Edit) error {
	insert := squirrel.Insert("data_items").
		Columns("id", "user_id", "type", "data", "meta").
		Values(obj.ID, obj.UserID, obj.Type, obj.Data, obj.Meta).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := insert.ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// Update modifies an existing data item based on the provided query parameters and Edit object,
// returning any error encountered during the operation.
func (r *Repo) Update(ctx context.Context, pars *model.GetPars, obj *model.Edit) error {
	if !pars.IsValid() {
		return errs.InvalidInput
	}

	queryBuilder := squirrel.Update("data_items")

	if obj.UserID != nil {
		queryBuilder = queryBuilder.Set("user_id", obj.UserID)
	}

	if obj.Type != nil {
		queryBuilder = queryBuilder.Set("type", obj.Type)
	}

	if obj.Data != nil {
		queryBuilder = queryBuilder.Set("data", obj.Data)
	}

	if obj.Meta != nil {
		queryBuilder = queryBuilder.Set("meta", obj.Meta)
	}

	if obj.UpdatedAt != nil {
		queryBuilder = queryBuilder.Set("updated_at", obj.UpdatedAt)
	}

	if obj.URL != nil {
		queryBuilder = queryBuilder.Set("url", obj.URL)
	}

	if len(pars.ID) > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.ID})
	}

	if len(pars.UserID) > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"user_id": pars.UserID})
	}

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, sql, args...)
	return err
}

// Delete removes a data item from the database based on the provided query parameters,
// returning any error encountered during the operation.
func (r *Repo) Delete(ctx context.Context, pars *model.GetPars) error {
	if !pars.IsValid() {
		return errs.InvalidInput
	}

	queryBuilder := squirrel.Delete("data_items")

	queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.ID})

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, sql, args...)
	return err
}

func (r *Repo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.Con.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (r *Repo) CommitTx(ctx context.Context, tx pgx.Tx) error {
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (r *Repo) RollbackTx(ctx context.Context, tx pgx.Tx) error {
	if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
		return err
	}
	return nil
}

func (r *Repo) HandleTxCompletion(tx pgx.Tx, err *error) {
	if p := recover(); p != nil {
		_ = tx.Rollback(context.Background())
		panic(p)
	} else if *err != nil {
		_ = tx.Rollback(context.Background())
	} else {
		*err = tx.Commit(context.Background())
	}
}
