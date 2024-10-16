package ports

import (
	"codebase-app/internal/module/product/entity"
	"context"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, req *entity.CreateProductReq) (*entity.CreateProductResp, error)
	GetProduct(ctx context.Context, req *entity.GetProductReq) (*entity.GetProductResp, error)
	GetProducts(ctx context.Context, req *entity.GetProductsReq) (*entity.GetProductsResp, error)

	IsProductValid(ctx context.Context, productId, userId string) error
}

type ProductService interface {
	CreateProduct(ctx context.Context, req *entity.CreateProductReq) (*entity.CreateProductResp, error)
	GetProduct(ctx context.Context, req *entity.GetProductReq) (*entity.GetProductResp, error)
	GetProducts(ctx context.Context, req *entity.GetProductsReq) (*entity.GetProductsResp, error)
}
