package service

import (
	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	"context"
)

var _ ports.ProductService = &productService{}

type productService struct {
	repo ports.ProductRepository
}

func NewProductService(repo ports.ProductRepository) *productService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req *entity.CreateProductReq) (*entity.CreateProductResp, error) {
	if req.ParentId != nil {
		err := s.repo.IsProductValid(ctx, *req.ParentId, req.UserId)
		if err != nil {
			return nil, err
		}
	}

	return s.repo.CreateProduct(ctx, req)
}

func (s *productService) GetProduct(ctx context.Context, req *entity.GetProductReq) (*entity.GetProductResp, error) {
	return s.repo.GetProduct(ctx, req)
}

func (s *productService) GetProducts(ctx context.Context, req *entity.GetProductsReq) (*entity.GetProductsResp, error) {
	return s.repo.GetProducts(ctx, req)
}
