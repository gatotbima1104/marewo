package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var _ ports.ProductRepository = &productRepo{}

type productRepo struct {
	db *sqlx.DB
}

func NewProductRepository() *productRepo {
	return &productRepo{
		db: adapter.Adapters.Postgres,
	}
}

func (r *productRepo) CreateProduct(ctx context.Context, req *entity.CreateProductReq) (*entity.CreateProductResp, error) {
	resp := new(entity.CreateProductResp)
	Id := ulid.Make().String()
	query := `
		INSERT INTO products (id, parent_id, company_id, branch_id, name, price, stock)
		VALUES (?, ?,
			(SELECT company_id FROM users WHERE id = ?),
			(SELECT branch_id FROM users WHERE id = ?),
			?, ?, ?
		)
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query),
		Id, req.ParentId,
		req.UserId, req.UserId,
		req.Name, req.Price, req.Stock,
	)
	if err != nil {
		log.Error().Err(err).Any("req", req).Msg("repo::CreateProduct - failed to create product")
		return nil, err
	}

	resp.Id = Id
	resp.ParentId = req.ParentId
	resp.Name = req.Name
	resp.Price = req.Price
	resp.Stock = req.Stock

	return resp, nil
}

func (r *productRepo) GetProduct(ctx context.Context, req *entity.GetProductReq) (*entity.GetProductResp, error) {
	return nil, nil
}

func (r *productRepo) GetProducts(ctx context.Context, req *entity.GetProductsReq) (*entity.GetProductsResp, error) {
	return nil, nil
}
