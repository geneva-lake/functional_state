package service

import (
	"context"
	"database/sql"

	"github.com/geneva-lake/functional_state/general"
	"github.com/google/uuid"
)

//   - -------------------------------------------------------------------------------------------------------------------
//     Repository struct
//   - -------------------------------------------------------------------------------------------------------------------
type Repository struct {
	pgsql *general.Pgsql
}

func NewRepository(connection *general.Pgsql) *Repository {
	return &Repository{
		pgsql: connection,
	}
}

//   - -------------------------------------------------------------------------------------------------------------------
//     Get user by id from storage with stored procedure
//   - -------------------------------------------------------------------------------------------------------------------
func (r *Repository) UserGetByID(ctx context.Context, id uuid.UUID) (*StoredUser, error) {
	u := new(StoredUser)
	err := r.pgsql.QueryRowContext(ctx, "select users.user__get_by_id($1::uuid)", id).Scan(u)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}
