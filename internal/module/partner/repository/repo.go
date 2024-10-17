package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/partner/entity"
	"codebase-app/internal/module/partner/ports"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var _ ports.PartnerRepository = &partnerRepo{}

type partnerRepo struct {
	db *sqlx.DB
}

func NewPartnerRepository() *partnerRepo {
	return &partnerRepo{
		db: adapter.Adapters.Postgres,
	}
}

func (r *partnerRepo) CreatePartner(ctx context.Context, req *entity.CreatePartnerReq) (*entity.CreatePartnerResp, error) {
	query := `
		INSERT INTO partners (id, company_id , name, address, phone_country_code, phone, latitude, longitude)
		VALUES (?, (SELECT company_id FROM users WHERE id = ?), ?, ?, ?, ?, ?, ?)
	`

	id := ulid.Make().String()

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), id, req.UserId,
		req.Name, req.Address, req.PhoneCountryCode, req.Phone, req.Lat, req.Lng,
	)
	if err != nil {
		log.Error().Err(err).Any("req", req).Msg("repo::CreatePartner - failed to create partner")
		return nil, err
	}

	return &entity.CreatePartnerResp{
		Id: id,
	}, nil
}

func (r *partnerRepo) GetPartners(ctx context.Context, req *entity.GetPartnersReq) (*entity.GetPartnersResp, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.PartnerItem
	}

	var (
		data = make([]dao, 0)
		res  = new(entity.GetPartnersResp)
	)
	res.Items = make([]entity.PartnerItem, 0)

	query := `
		WITH user_company AS (
			SELECT company_id FROM users WHERE id = ?
		)
		SELECT
			COUNT(*) OVER() AS total_data,
			id,
			name,
			address,
			latitude,
			longitude,
			phone_country_code,
			phone
		FROM partners
		WHERE company_id = (SELECT company_id FROM user_company)
		LIMIT ? OFFSET ?
	`

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query),
		req.UserId,
		req.Paginate, (req.Page-1)*req.Paginate,
	)
	if err != nil {
		log.Error().Err(err).Any("req", req).Msg("repo::GetPartners - failed to get partners")
		return nil, err
	}

	if len(data) > 0 {
		res.Meta.TotalData = data[0].TotalData
	}

	for _, d := range data {
		res.Items = append(res.Items, d.PartnerItem)
	}

	res.Meta.CountTotalPage(req.Page, req.Paginate, res.Meta.TotalData)

	return res, nil
}
