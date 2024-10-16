package entity

import (
	"codebase-app/pkg/types"

	"github.com/LukaGiorgadze/gonull"
)

type CreateProductReq struct {
	UserId string

	ParentId *string `json:"parent_id" validate:"omitempty,ulid,exist=products.id"`
	Name     string  `json:"name" validate:"required"`
	Price    float64 `json:"price"`
	Stock    int64   `json:"stock"`
}

type CreateProductResp struct {
	Id       string  `json:"id" db:"id"`
	ParentId *string `json:"parent_id" db:"parent_id"`
	Name     string  `json:"name" db:"name"`
	Price    float64 `json:"price" db:"price"`
	Stock    int64   `json:"stock" db:"stock"`
}

type GetProductReq struct {
	UserId string
	Id     string `json:"id" validate:"required,ulid"`
}

type GetProductResp struct {
	Id          string  `json:"id" db:"id"`
	ParentId    *string `json:"parent_id" db:"parent_id"`
	CompanyId   string  `json:"company_id" db:"company_id"`
	CompanyName string  `json:"company_name" db:"company_name"`
	Name        string  `json:"name" db:"name"`
	Price       float64 `json:"price" dbc:"price"`
	Stock       int64   `json:"stock" db:"stock"`
}

type GetProductsReq struct {
	UserId   string
	Page     int `json:"page" query:"page" validate:"required"`
	Paginate int `json:"paginate" query:"paginate" validate:"required"`
}

func (r *GetProductsReq) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}

type GetProductsResp struct {
	Items []ProductItem `json:"data"`
	Meta  types.Meta    `json:"meta"`
}

type ProductItem struct {
	Id       string  `json:"id" db:"id"`
	ParentId *string `json:"parent_id" db:"parent_id"`
	Name     string  `json:"name" db:"name"`
	Price    float64 `json:"price" db:"price"`
	Stock    int64   `json:"stock" db:"stock"`
}

type UpdateProductReq struct {
	Id       string                   `json:"id" validate:"required,ulid"`
	ParentId gonull.Nullable[string]  `json:"parent_id"`
	Name     gonull.Nullable[string]  `json:"name"`
	Price    gonull.Nullable[float64] `json:"price"`
}
