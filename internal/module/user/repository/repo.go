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

func (r *userRepo) GetCouriers(ctx context.Context, req *entity.GetCouriersReq) (*entity.GetCouriersRes, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.CourierItem
	}

	var (
		data = make([]dao, 0)
		res  entity.GetCouriersRes
	)
	res.Items = make([]entity.CourierItem, 0)

	query := `
		WITH role_courier AS (
			SELECT
				id
			FROM
				roles
			WHERE
				name = 'courier'
		),
		user_company AS (
			SELECT
				company_id
			FROM
				users
			WHERE
				id = ?
		)
		SELECT
			COUNT(*) OVER() AS total_data,
			u.id,
			u.branch_id,
			u.name,
			b.name as branch_name,
			u.email,
			u.created_at
		FROM
			users u
		LEFT JOIN
			branches b
			ON b.id = u.branch_id
		WHERE
			u.role_id = (SELECT id FROM role_courier)
			AND u.company_id = (SELECT company_id FROM user_company)
		LIMIT ? OFFSET ?
	`

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), req.UserId, req.Paginate, (req.Page-1)*req.Paginate)
	if err != nil {
		log.Error().Err(err).Any("req", req).Msg("repo::GetCouriers - failed to get couriers")
		return nil, err
	}

	for _, d := range data {
		res.Items = append(res.Items, d.CourierItem)
	}

	if len(res.Items) > 0 {
		res.Meta.TotalData = data[0].TotalData
	}

	res.Meta.CountTotalPage(req.Page, req.Paginate, res.Meta.TotalData)

	return &res, nil
}
