package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	"codebase-app/pkg/errmsg"
	"context"
	"database/sql"

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
	resp := new(entity.GetProductResp)
	query := `
		SELECT p.id, p.parent_id, p.name, p.price, p.stock, c.id as company_id, c.name as company_name
		FROM
			products p
		LEFT JOIN
			companies c ON c.id = p.company_id
		WHERE p.id = ?
		AND p.deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, resp, r.db.Rebind(query), req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Err(err).Any("req", req).Msg("repo::GetProduct - product not found")
			return nil, errmsg.NewCustomErrors(404).SetMessage("Produk tidak ditemukan")
		}
		log.Error().Err(err).Str("id", req.Id).Msg("repo::GetProduct - failed to get product")
		return nil, err
	}

	return resp, nil
}

func (r *productRepo) GetProducts(ctx context.Context, req *entity.GetProductsReq) (*entity.GetProductsResp, error) {
	return nil, nil
}

func (r *productRepo) IsProductValid(ctx context.Context, productId, userId string) error {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM products
			WHERE id = ?
			AND company_id = (SELECT company_id FROM users WHERE id = ?)
		)
	`

	var isValid bool

	err := r.db.GetContext(ctx, &isValid, r.db.Rebind(query), productId, userId)
	if err != nil {
		log.Error().Err(err).Str("productId", productId).Str("userId", userId).Msg("repo::IsProductValid - failed to check product")
		return err
	}

	if !isValid {
		log.Warn().Str("productId", productId).Str("userId", userId).Msg("repo::IsProductValid - product not valid")
		return errmsg.NewCustomErrors(404).SetMessage("Produk tidak ditemukan")
	}

	return nil
}
