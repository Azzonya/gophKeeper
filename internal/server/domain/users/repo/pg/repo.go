// Package pg provides a PostgreSQL-based implementation for managing user accounts,
// including operations for retrieving, listing, creating, updating, deleting, and checking user existence.
package pg

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gophKeeper/internal/server/domain/users/model"
	"gophKeeper/internal/server/errs"
)

// Repo provides methods to interact with the PostgreSQL database for user account operations.
// It holds a connection pool to manage database interactions.
type Repo struct {
	Con *pgxpool.Pool
}

// New creates a new Repo instance with the given PostgreSQL connection pool.
func New(con *pgxpool.Pool) *Repo {
	return &Repo{
		con,
	}
}

// Get retrieves a user based on the provided query parameters. It returns the user if found,
// a boolean indicating the user's existence, and any error encountered.
func (r *Repo) Get(ctx context.Context, pars *model.GetPars) (*model.User, bool, error) {
	if !pars.IsValid() {
		return nil, false, errs.InvalidInput
	}

	var result model.User

	queryBuilder := squirrel.Select("*").From("users")

	if len(pars.UserID) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.UserID})
	}

	if len(pars.Username) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"username": pars.Username})
	}

	queryBuilder = queryBuilder.Limit(1)

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, false, err
	}

	err = r.Con.QueryRow(ctx, sql, args...).Scan(&result.UserID, &result.Username, &result.PasswordHash, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &result, true, nil
}

// List retrieves multiple users based on the provided filtering parameters,
// supporting filters like UserIDs, Username, and timestamps. It returns the list
// of users, the total count, and any error encountered.
func (r *Repo) List(ctx context.Context, pars *model.ListPars) ([]*model.User, int64, error) {
	queryBuilder := squirrel.
		Select("id", "username", "password_hash", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"true": true})

	if pars.UserID != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"user_id": pars.UserID})
	}

	if pars.UserIDs != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"user_id": pars.UserIDs})
	}

	if pars.Username != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"username": pars.Username})
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

	var result []*model.User
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.UserID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return result, int64(len(result)), nil
}

// Create inserts a new user into the database based on the provided Edit object,
// returning any error encountered during the operation.
func (r *Repo) Create(ctx context.Context, obj *model.Edit) error {
	insert := squirrel.Insert("users").
		Columns("username", "password_hash").
		Values(obj.Username, obj.PasswordHash)

	query, args, err := insert.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, query, args...)
	return err
}

// Update modifies an existing user record based on the provided query parameters and Edit object,
// returning any error encountered during the operation.
func (r *Repo) Update(ctx context.Context, pars *model.GetPars, obj *model.Edit) error {
	if !pars.IsValid() {
		return errs.InvalidInput
	}

	queryBuilder := squirrel.Update("users")

	if obj.Username != nil {
		queryBuilder = queryBuilder.Set("username", obj.Username)
	}

	if obj.PasswordHash != nil {
		queryBuilder = queryBuilder.Set("password_hash", obj.PasswordHash)
	}

	if obj.UpdatedAt != nil {
		queryBuilder = queryBuilder.Set("updated_at", obj.UpdatedAt)
	}

	queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.UserID})

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, sql, args...)
	return err
}

// Delete removes a user from the database based on the provided query parameters,
// returning any error encountered during the operation.
func (r *Repo) Delete(ctx context.Context, pars *model.GetPars) error {
	if !pars.IsValid() {
		return errs.InvalidInput
	}

	queryBuilder := squirrel.Delete("users")

	queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.UserID})

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, sql, args...)
	return err
}

// Exists checks whether a user exists in the database based on the provided query parameters.
// It returns a boolean indicating existence and any error encountered.
func (r *Repo) Exists(ctx context.Context, pars *model.GetPars) (bool, error) {
	if !pars.IsValid() {
		return false, errs.InvalidInput
	}

	existsQuery := "SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)"
	var exist bool

	err := r.Con.QueryRow(ctx, existsQuery, pars.Username).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, nil
}
