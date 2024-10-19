package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"
)

func (r *deliveryRepo) IsValidPartner(ctx context.Context, partnerId string, userId string) error {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM partners
			WHERE id = ?
			AND company_id = (SELECT company_id FROM users WHERE id = ?)
		)
	`

	var isValid bool
	if err := r.db.GetContext(ctx, &isValid, r.db.Rebind(query), partnerId, userId); err != nil {
		log.Error().Err(err).Str("partner_id", partnerId).Str("user_id", userId).Msg("repo::IsValidPartner - failed to check partner")
		return errors.New("bukan partner yang valid")
	}

	if !isValid {
		return errors.New("bukan partner yang valid")
	}

	return nil
}

func (r *deliveryRepo) IsValidCourier(ctx context.Context, courierId string, userId string) error {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE
				id = ?
				AND company_id = (SELECT company_id FROM users WHERE id = ?)
				AND role_id = (SELECT id FROM roles WHERE name = 'courier')
		)
	`

	var isValid bool
	if err := r.db.GetContext(ctx, &isValid, r.db.Rebind(query), courierId, userId); err != nil {
		log.Error().Err(err).Str("courier_id", courierId).Str("user_id", userId).Msg("repo::IsValidCourier - failed to check courier")
		return errors.New("bukan kurir yang valid")
	}

	if !isValid {
		return errors.New("bukan kurir yang valid")
	}

	return nil
}

func (r *deliveryRepo) IsValidDeliveryTemplate(ctx context.Context, templateId string, userId string) error {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM delivery_templates
			WHERE id = ?
			AND company_id = (SELECT company_id FROM users WHERE id = ?)
		)
	`

	var isValid bool
	if err := r.db.GetContext(ctx, &isValid, r.db.Rebind(query), templateId, userId); err != nil {
		log.Error().Err(err).Str("template_id", templateId).Str("user_id", userId).Msg("repo::IsValidDeliveryTemplate - failed to check delivery template")
		return errors.New("bukan template pengiriman yang valid")
	}

	if !isValid {
		return errors.New("bukan template pengiriman yang valid")
	}

	return nil
}

func (r *deliveryRepo) BindCompanyAndBranch(ctx context.Context, userId string, companyId *string, branchId *string) error {
	if branchId != nil {
		query := `SELECT company_id, branch_id FROM users WHERE id = ?`
		row := r.db.QueryRowxContext(ctx, r.db.Rebind(query), userId)
		if err := row.Scan(companyId, branchId); err != nil {
			if err == sql.ErrNoRows {
				log.Warn().Err(err).Str("user_id", userId).Msg("repo::BindCompanyAndBranch - user not found")
				return errors.New("pengguna tidak ditemukan")
			}
			log.Error().Err(err).Str("user_id", userId).Msg("repo::BindCompanyAndBranch - failed to get company and branch")
			return err
		}
	} else {
		query := `SELECT company_id FROM users WHERE id = ?`
		row := r.db.QueryRowxContext(ctx, r.db.Rebind(query), userId)
		if err := row.Scan(companyId); err != nil {
			if err == sql.ErrNoRows {
				log.Warn().Err(err).Str("user_id", userId).Msg("repo::BindCompanyAndBranch - user not found")
				return errors.New("pengguna tidak ditemukan")
			}
			log.Error().Err(err).Str("user_id", userId).Msg("repo::BindCompanyAndBranch - failed to get company")
			return err
		}
	}

	return nil
}
