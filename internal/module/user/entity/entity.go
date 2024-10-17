package entity

import (
	"codebase-app/pkg/types"
	"time"
)

type LoginReq struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required"`
}

type LoginResp struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	RoleName    string  `json:"role_name"`
	CompanyName string  `json:"company_name"`
	BranchName  *string `json:"branch_name"`
	Token       string  `json:"token"`
}

type UserResult struct {
	Id          string  `db:"id"`
	Name        string  `db:"name"`
	Email       string  `db:"email"`
	RoleName    string  `db:"role_name"`
	CompanyName string  `db:"company_name"`
	BranchName  *string `db:"branch_name"`
	Password    string  `db:"password"`
}

type GetCouriersReq struct {
	UserId string

	Page     int `json:"page" validate:"required"`
	Paginate int `json:"paginate" validate:"required"`
}

func (req *GetCouriersReq) SetDefault() {
	if req.Page < 1 {
		req.Page = 1
	}

	if req.Paginate < 1 {
		req.Paginate = 10
	}
}

type GetCouriersRes struct {
	Items []CourierItem `json:"items"`
	Meta  types.Meta    `json:"meta"`
}

type CourierItem struct {
	Id         string    `json:"id" db:"id"`
	BranchId   *string   `json:"branch_id" db:"branch_id"`
	Name       string    `json:"name" db:"name"`
	BranchName *string   `json:"branch_name" db:"branch_name"`
	Email      string    `json:"email" db:"email"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
