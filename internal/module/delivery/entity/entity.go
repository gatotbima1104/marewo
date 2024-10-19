package entity

import (
	"codebase-app/pkg/errmsg"
	"strconv"
)

type CreateSchedulesReq struct {
	UserId    string
	CompanyId string
	Schedules []PartnerSchedule `json:"partner_schedules" validate:"required,dive"`
}

type PartnerSchedule struct {
	PartnerId        string             `json:"partner_id" validate:"ulid"`
	DeliverySchedule []DeliverySchedule `json:"schedules" validate:"required,dive"`
}

type DeliverySchedule struct {
	CourierId string        `json:"courier_id" validate:"ulid"`
	DateTime  int64         `json:"date_time" validate:"required"`
	Products  []ProductItem `json:"products" validate:"required,dive"`
}

type ProductItem struct {
	ProductId string  `json:"product_id" db:"product_id" validate:"ulid"`
	Quantity  int     `json:"quantity" validate:"required"`
	Name      string  `db:"product_name"`
	Price     float64 `db:"product_price"`
	Notes     *string `json:"notes"`
}

type CreateSchedulesRes struct {
	ScheduleIds []string `json:"schedule_ids"`
}

// schedule template
type CreateScheduleTemplatesReq struct {
	UserId    string
	CompanyId string
	BranchId  *string
	Schedules []PartnerScheduleTemp `json:"partner_schedules" validate:"required,dive"`
}

type PartnerScheduleTemp struct {
	PartnerId string                 `json:"partner_id" validate:"ulid"`
	Schedules []DeliveryScheduleTemp `json:"schedules" validate:"required,dive"`
}

type DeliveryScheduleTemp struct {
	CourierId string        `json:"courier_id" validate:"ulid"`
	TimeStart string        `json:"time_start" validate:"required,datetime=15:04:05"`
	TimeEnd   *string       `json:"time_end" validate:"omitempty,datetime=15:04:05"`
	Products  []ProductItem `json:"products" validate:"required,dive"`
}

func (r *CreateScheduleTemplatesReq) Validate() error {
	errs := errmsg.NewCustomErrors(400)

	for PIndex, partnerSchedule := range r.Schedules {
		pIdx := strconv.Itoa(PIndex)
		for CIndex, schedule := range partnerSchedule.Schedules {
			cIdx := strconv.Itoa(CIndex)
			if schedule.TimeEnd != nil {
				if schedule.TimeStart >= *schedule.TimeEnd {
					errs.Add(`partner_schedules.[`+pIdx+`].schedules.[`+cIdx+`].time_end`, `time_end harus lebih besar dari time_start`)
				}
			}
		}
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

type CreateScheduleTemplatesRes struct {
	PartnerScheduleIds []string `json:"partner_schedule_ids"`
}

type ApplyScheduleTemplatesReq struct {
	UserId    string
	CompanyId string
	BranchId  *string
	Id        string   `json:"id" validate:"ulid"`
	Dates     []string `json:"dates" validate:"required,unique_in_slice,dive,datetime=2006-01-02"`
	Timezone  string   `json:"timezone" validate:"timezone"`
}

func (r *ApplyScheduleTemplatesReq) SetDefault() {
	if r.Timezone == "" {
		r.Timezone = "Asia/Makassar"
	}
}
