package ports

import (
	"codebase-app/internal/module/delivery/entity"
	"context"
)

type DeliveryRepository interface {
	CreateSchedules(ctx context.Context, req *entity.CreateSchedulesReq) (*entity.CreateSchedulesRes, error)
	CreateScheduleTemplates(ctx context.Context, req *entity.CreateScheduleTemplatesReq) (*entity.CreateScheduleTemplatesRes, error)
	ApplyScheduleTemplates(ctx context.Context, req *entity.ApplyScheduleTemplatesReq) error

	IsValidPartner(ctx context.Context, partnerId string, userId string) error
	IsValidCourier(ctx context.Context, courierId string, userId string) error
	IsValidDeliveryTemplate(ctx context.Context, templateId string, userId string) error
	BindCompanyAndBranch(ctx context.Context, userId string, companyId, branchId *string) error
}

type DeliveryService interface {
	CreateSchedules(ctx context.Context, req *entity.CreateSchedulesReq) (*entity.CreateSchedulesRes, error)
	CreateScheduleTemplates(ctx context.Context, req *entity.CreateScheduleTemplatesReq) (*entity.CreateScheduleTemplatesRes, error)

	ApplyScheduleTemplates(ctx context.Context, req *entity.ApplyScheduleTemplatesReq) error
}
