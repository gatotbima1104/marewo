package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/user/entity"
	"codebase-app/internal/module/user/ports"
	"codebase-app/pkg/errmsg"
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.UserRepository = &userRepo{}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository() *userRepo {
	return &userRepo{
		db: adapter.Adapters.Postgres,
	}
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*entity.UserResult, error) {
	query := `
		SELECT
			u.id, u.name,u.email, u.password, r.name as role_name, c.name as company_name, b.name as branch_name
		FROM
			users u
		JOIN
			roles r
			ON r.id = u.role_id
		JOIN
			companies c
			ON c.id = u.company_id
		LEFT JOIN
			branches b
			ON b.id = u.branch_id
		WHERE
			u.email = ?
	`

	var user entity.UserResult

	err := r.db.GetContext(ctx, &user, r.db.Rebind(query), email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Err(err).Any("email", email).Msg("repo::FindByEmail - user not found")
			return nil, errmsg.NewCustomErrors(400).SetMessage("Kredensial tidak valid")
		}

		log.Error().Err(err).Any("email", email).Msg("repo::FindByEmail - failed to get user by email")
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) FindById(ctx context.Context, id string) (*entity.UserResult, error) {
	return nil, nil
}
