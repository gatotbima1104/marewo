package ports

import (
	"codebase-app/internal/module/partner/entity"
	"context"
)

type PartnerRepository interface {
	CreatePartner(ctx context.Context, req *entity.CreatePartnerReq) (*entity.CreatePartnerResp, error)
	GetPartners(ctx context.Context, req *entity.GetPartnersReq) (*entity.GetPartnersResp, error)
}

type PartnerService interface {
	CreatePartner(ctx context.Context, req *entity.CreatePartnerReq) (*entity.CreatePartnerResp, error)
	GetPartners(ctx context.Context, req *entity.GetPartnersReq) (*entity.GetPartnersResp, error)
}
