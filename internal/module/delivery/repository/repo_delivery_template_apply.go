package repository

import (
	"codebase-app/internal/module/delivery/entity"
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func (r *deliveryRepo) ApplyScheduleTemplates(ctx context.Context, req *entity.ApplyScheduleTemplatesReq) error {
	query := `
		INSERT INTO delivery_schedules
			(id, company_id, partner_id, courier_id, status, delivery_at)
		SELECT
			?, ?, dt.partner_id, dt.courier_id, 'scheduled',
			(
				TO_TIMESTAMP(
					CONCAT(?::text, ' ', TO_CHAR(dt.time_start, 'HH24:MI:SS')),
					'YYYY-MM-DD HH24:MI:SS'
				) AT TIME ZONE 'UTC' -- First, convert to UTC
			) AT TIME ZONE ? -- Then convert from UTC to the desired timezone
		FROM
			delivery_templates dt
		WHERE
			id = ?
	`

	queryItems := `
		INSERT INTO delivery_items
			(id, delivery_schedule_id, product_id, product_name, price, quantity, total)
		VALUES
			(?, ?, ?, ?, ?, ?, ?)
	`

	type deliveryTemplateItemDAO struct {
		ProductId    string  `db:"product_id"`
		Quantity     int     `db:"quantity"`
		ProductPrice float64 `db:"product_price"`
		ProductName  string  `db:"product_name"`
	}

	var deliveryTemplateItems []deliveryTemplateItemDAO

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Any("req", req).Msg("repo::ApplyScheduleTemplates - failed to begin transaction")
		return err
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Error().Err(err).Msg("repo::ApplyScheduleTemplates - failed to rollback transaction")
			}
		} else {
			if err := tx.Commit(); err != nil {
				log.Error().Err(err).Msg("repo::ApplyScheduleTemplates - failed to commit transaction")
			}
		}
	}()

	queryTemplateItems := `
		SELECT
			dti.product_id, dti.quantity, p.price AS product_price, p.name AS product_name
		FROM
			delivery_template_items dti
		JOIN
			products p ON p.id = dti.product_id
		WHERE
			dti.delivery_template_id = ?
	`

	err = tx.SelectContext(ctx, &deliveryTemplateItems, r.db.Rebind(queryTemplateItems), req.Id)
	if err != nil {
		log.Error().Err(err).Any("req", req).Msg("repo::ApplyScheduleTemplates - failed to get delivery template items")
		return err
	}

	for _, date := range req.Dates {
		scheduleId := ulid.Make().String()
		_, err = tx.ExecContext(ctx, r.db.Rebind(query),
			scheduleId,
			req.CompanyId,
			date,
			req.Timezone,
			req.Id,
		)
		if err != nil {
			log.Error().Err(err).Any("req", req).Msg("repo::ApplyScheduleTemplates - failed to insert delivery schedule")
			return err
		}

		for _, item := range deliveryTemplateItems {
			_, err = tx.ExecContext(ctx, r.db.Rebind(queryItems),
				ulid.Make().String(),
				scheduleId,
				item.ProductId,
				item.ProductName,
				item.ProductPrice,
				item.Quantity,
				item.ProductPrice*float64(item.Quantity),
			)
			if err != nil {
				log.Error().Err(err).Any("req", req).Msg("repo::ApplyScheduleTemplates - failed to insert delivery item")
				return err
			}
		}

		// update delivery schedule total
		queryUpdateTotal := `
			UPDATE delivery_schedules
			SET total = (
				SELECT SUM(total)
				FROM delivery_items
				WHERE delivery_schedule_id = ?
			)
			WHERE id = ?
		`

		_, err = tx.ExecContext(ctx, r.db.Rebind(queryUpdateTotal), scheduleId, scheduleId)
		if err != nil {
			log.Error().Err(err).Any("req", req).Msg("repo::ApplyScheduleTemplates - failed to update delivery schedule total")
			return err
		}
	}

	return nil
}
