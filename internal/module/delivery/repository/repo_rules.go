package repository

import (
	"context"
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
