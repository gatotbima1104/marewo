package entity

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
}

type CreateSchedulesRes struct {
	ScheduleIds []string `json:"schedule_ids"`
}
