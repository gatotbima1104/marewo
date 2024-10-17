package service

import (
	"codebase-app/internal/module/partner/entity"
	"codebase-app/internal/module/partner/ports"
	"context"
)

var _ ports.PartnerService = &partnerService{}

type partnerService struct {
	repo ports.PartnerService
}

func NewPartnerService(repo ports.PartnerRepository) *partnerService {
	return &partnerService{
		repo: repo,
	}
}

func (s *partnerService) CreatePartner(ctx context.Context, req *entity.CreatePartnerReq) (*entity.CreatePartnerResp, error) {
	return s.repo.CreatePartner(ctx, req)
}

func (s *partnerService) GetPartners(ctx context.Context, req *entity.GetPartnersReq) (*entity.GetPartnersResp, error) {
	return s.repo.GetPartners(ctx, req)
}
