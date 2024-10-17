package ports

import (
	"codebase-app/internal/module/delivery/entity"
	"context"
)

type DeliveryRepository interface {
	CreateSchedules(ctx context.Context, req *entity.CreateSchedulesReq) (*entity.CreateSchedulesRes, error)

	IsValidPartner(ctx context.Context, partnerId string, userId string) error
	IsValidCourier(ctx context.Context, courierId string, userId string) error
}

type DeliveryService interface {
	CreateSchedules(ctx context.Context, req *entity.CreateSchedulesReq) (*entity.CreateSchedulesRes, error)
}
