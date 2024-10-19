package repository

import (
	"codebase-app/internal/module/delivery/entity"
	"codebase-app/pkg/errmsg"
	"context"
	"database/sql"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func (r *deliveryRepo) CreateScheduleTemplates(ctx context.Context, req *entity.CreateScheduleTemplatesReq) (*entity.CreateScheduleTemplatesRes, error) {
	var res entity.CreateScheduleTemplatesRes

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Any("req", req).Msg("repo::CreateScheduleTemplates - failed to begin transaction")
		return nil, err
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Error().Err(err).Msg("repo::CreateScheduleTemplates - failed to rollback transaction")
			}
		} else {
			if err := tx.Commit(); err != nil {
				log.Error().Err(err).Msg("repo::CreateScheduleTemplates - failed to commit transaction")
			}
		}
	}()

	query := `
		INSERT INTO delivery_templates (
			id,
			company_id,
			branch_id,
			partner_id,
			courier_id,
			time_start,
			time_end
		) VALUES ( ?, ?, ?, ?, ?, ?, ? )
	`

	queryItems := `
		INSERT INTO delivery_template_items (
			id,
			delivery_template_id,
			product_id,
			quantity,
			notes
		) VALUES ( ?, ?, ?, ?, ? )
	`

	queryCompany := `SELECT company_id, branch_id FROM users WHERE id = ?`
	row := r.db.QueryRowxContext(ctx, r.db.Rebind(queryCompany), req.UserId)
	if err := row.Scan(&req.CompanyId, &req.BranchId); err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Err(err).Any("req", req).Msg("repository::CreateScheduleTemplates - user not found")
			return nil, errmsg.NewCustomErrors(404).SetMessage("Pengguna tidak ditemukan")
		}
		log.Error().Err(err).Any("req", req).Msg("repository::CreateScheduleTemplates - failed to get company id")
		return nil, err
	}

	for _, partnerSchedule := range req.Schedules {
		partnerId := partnerSchedule.PartnerId
		for _, schedule := range partnerSchedule.Schedules {
			courierId := schedule.CourierId
			timeStart := schedule.TimeStart
			timeEnd := schedule.TimeEnd

			deliveryTemplateId := ulid.Make().String()
			res.PartnerScheduleIds = append(res.PartnerScheduleIds, deliveryTemplateId)
			_, err = tx.ExecContext(ctx, r.db.Rebind(query),
				deliveryTemplateId,
				req.CompanyId,
				req.BranchId,
				partnerId,
				courierId,
				timeStart,
				timeEnd,
			)
			if err != nil {
				log.Error().Err(err).Any("req", req).Msg("repository::CreateScheduleTemplates - failed to insert delivery template")
				return nil, err
			}

			for _, product := range schedule.Products {
				_, err = tx.ExecContext(ctx, r.db.Rebind(queryItems),
					ulid.Make().String(),
					deliveryTemplateId,
					product.ProductId,
					product.Quantity,
					product.Notes,
				)
				if err != nil {
					log.Error().Err(err).Any("req", req).Msg("repository::CreateScheduleTemplates - failed to insert delivery template items")
					return nil, err
				}
			}
		}
	}

	return &res, nil
}
