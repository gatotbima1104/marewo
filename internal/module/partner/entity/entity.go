package entity

import "codebase-app/pkg/types"

type CreatePartnerReq struct {
	UserId string `json:"-"`

	Name             string  `json:"name" validate:"required"`
	Address          string  `json:"address" validate:"required"`
	PhoneCountryCode string  `json:"phone_country_code" validate:"required"`
	Phone            string  `json:"phone" validate:"required"`
	Lat              float64 `json:"latitude" validate:"required"`
	Lng              float64 `json:"longitude" validate:"required"`
}

func (r *CreatePartnerReq) SetDefault() {
	if r.PhoneCountryCode == "" {
		r.PhoneCountryCode = "62"
	}
}

type CreatePartnerResp struct {
	Id string `json:"id"`
}

type GetPartnersReq struct {
	UserId   string `json:"-"`
	Page     int    `query:"page" validate:"required"`
	Paginate int    `query:"paginate" validate:"required"`
}

func (r *GetPartnersReq) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}

type GetPartnersResp struct {
	Items []PartnerItem `json:"items"`
	Meta  types.Meta    `json:"meta"`
}

type PartnerItem struct {
	Id               string  `json:"id" db:"id"`
	Name             string  `json:"name" db:"name"`
	Address          string  `json:"address" db:"address"`
	PhoneCountryCode string  `json:"phone_country_code" db:"phone_country_code"`
	Phone            string  `json:"phone" db:"phone"`
	Lat              float64 `json:"latitude" db:"latitude"`
	Lng              float64 `json:"longitude" db:"longitude"`
}
