package service

import (
	"codebase-app/internal/module/delivery/entity"
	"codebase-app/internal/module/delivery/ports"
	"codebase-app/pkg/errmsg"
	"context"
	"strconv"

	"github.com/rs/zerolog/log"
)

var _ ports.DeliveryService = &deliveryService{}

type deliveryService struct {
	repo ports.DeliveryRepository
}

func NewDeliveryService(repo ports.DeliveryRepository) *deliveryService {
	return &deliveryService{
		repo: repo,
	}
}

func (s *deliveryService) CreateSchedules(ctx context.Context, req *entity.CreateSchedulesReq) (*entity.CreateSchedulesRes, error) {
	errs := errmsg.NewCustomErrors(400)

	for partnersIdx, partnerSchedule := range req.Schedules {
		pIdx := strconv.Itoa(partnersIdx)

		if err := s.repo.IsValidPartner(ctx, partnerSchedule.PartnerId, req.UserId); err != nil {
			errs.Add("partner_schedules["+pIdx+"].partner_id", err.Error())
		}

		for schedulesIdx, schedule := range partnerSchedule.DeliverySchedule {
			if err := s.repo.IsValidCourier(ctx, schedule.CourierId, req.UserId); err != nil {
				cIdx := strconv.Itoa(schedulesIdx)
				errs.Add("partner_schedules["+pIdx+"].schedules["+cIdx+"].courier_id", err.Error())
			}
		}
	}

	if errs.HasErrors() {
		log.Warn().Err(errs).Any("req", req).Msg("service::CreateSchedules - invalid request")
		return nil, errs
	}

	return s.repo.CreateSchedules(ctx, req)
}
