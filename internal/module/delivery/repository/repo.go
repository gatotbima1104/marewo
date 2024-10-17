package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/delivery/entity"
	"codebase-app/internal/module/delivery/ports"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var _ ports.DeliveryRepository = &deliveryRepo{}

type deliveryRepo struct {
	db *sqlx.DB
}

func NewDeliveryRepository() *deliveryRepo {
	return &deliveryRepo{
		db: adapter.Adapters.Postgres,
	}
}

func (r *deliveryRepo) CreateSchedules(ctx context.Context, req *entity.CreateSchedulesReq) (*entity.CreateSchedulesRes, error) {
	var res entity.CreateSchedulesRes
	res.ScheduleIds = make([]string, 0)

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Any("req", req).Msg("repo::CreateSchedules - failed to begin transaction")
		return nil, err
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Error().Err(err).Msg("repo::CreateSchedules - failed to rollback transaction")
			}
		} else {
			if err := tx.Commit(); err != nil {
				log.Error().Err(err).Msg("repo::CreateSchedules - failed to commit transaction")
			}
		}
	}()

	query := `
		INSERT INTO delivery_schedules
			(id, company_id, partner_id, courier_id, status, delivery_at)
		VALUES
			(?, ?, ?, ?, ?, ?)
	`

	queryProduct := `
		INSERT INTO delivery_items
			(id, delivery_schedule_id, product_id, product_name, price, quantity, total)
		VALUES
			(?, ?, ?, ?, ?, ?, ?)
	`

	queryCompany := `SELECT company_id FROM users WHERE id = ?`
	if err := tx.GetContext(ctx, &req.CompanyId, r.db.Rebind(queryCompany), req.UserId); err != nil {
		log.Error().Err(err).Any("req", req).Msg("repo::CreateSchedules - failed to get company_id")
		return nil, err
	}

	if err := r.getProductDetails(ctx, req); err != nil {
		return nil, err
	}

	for _, schedule := range req.Schedules {
		for _, deliverySchedule := range schedule.DeliverySchedule {
			id := ulid.Make().String()
			deliveryAt := time.UnixMilli(deliverySchedule.DateTime)

			_, err := tx.ExecContext(
				ctx,
				tx.Rebind(query),
				id,
				req.CompanyId,
				schedule.PartnerId,
				deliverySchedule.CourierId,
				"scheduled",
				deliveryAt,
			)
			if err != nil {
				log.Error().Err(err).Any("req", req).Msg("repo::CreateSchedules - failed to insert delivery_schedules")
				return nil, err
			}

			for _, product := range deliverySchedule.Products {
				_, err := tx.ExecContext(
					ctx,
					tx.Rebind(queryProduct),
					ulid.Make().String(),
					id,
					product.ProductId,
					product.Name,
					product.Price,
					product.Quantity,
					product.Price*float64(product.Quantity),
				)
				if err != nil {
					log.Error().Err(err).Any("req", req).Msg("repo::CreateSchedules - failed to insert delivery_items")
					return nil, err
				}
			}

			queryUpdate := `UPDATE delivery_schedules SET total = (SELECT SUM(total) FROM delivery_items WHERE delivery_schedule_id = ?) WHERE id = ?`
			_, err = tx.ExecContext(ctx, tx.Rebind(queryUpdate), id, id)
			if err != nil {
				log.Error().Err(err).Any("req", req).Msg("repo::CreateSchedules - failed to update total")
				return nil, err
			}

			res.ScheduleIds = append(res.ScheduleIds, id)
		}
	}

	return &res, nil
}

func (r *deliveryRepo) getProductDetails(ctx context.Context, req *entity.CreateSchedulesReq) error {
	query := `
		SELECT
			id AS product_id,
			name AS product_name,
			price AS product_price
		FROM products
		WHERE id IN (?)
	`
	var productIds []string

	for _, schedule := range req.Schedules {
		for _, deliverySchedule := range schedule.DeliverySchedule {
			for _, product := range deliverySchedule.Products {
				productIds = append(productIds, product.ProductId)
			}
		}
	}

	// elimate productIds that duplicate
	productIds = removeDuplicateProductIds(productIds)

	query, args, err := sqlx.In(query, productIds)
	if err != nil {
		log.Error().Err(err).Any("req", req).Msg("repo::getProductDetails - failed to bind query")
		return err
	}

	rows, err := r.db.QueryxContext(ctx, r.db.Rebind(query), args...)
	if err != nil {
		log.Error().Err(err).Any("req", req).Msg("repo::getProductDetails - failed to query products")
		return err
	}

	// make a map of product details
	productDetails := make(map[string]entity.ProductItem)
	for rows.Next() {
		var product entity.ProductItem
		if err := rows.StructScan(&product); err != nil {
			log.Error().Err(err).Any("req", req).Msg("repo::getProductDetails - failed to scan product")
			return err
		}
		productDetails[product.ProductId] = product
	}

	// assign product details to each product
	for _, schedule := range req.Schedules {
		for _, deliverySchedule := range schedule.DeliverySchedule {
			for i, product := range deliverySchedule.Products {
				if productDetail, ok := productDetails[product.ProductId]; ok {
					deliverySchedule.Products[i].Name = productDetail.Name
					deliverySchedule.Products[i].Price = productDetail.Price
				}
			}
		}
	}

	return nil
}

func removeDuplicateProductIds(ids []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range ids {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
