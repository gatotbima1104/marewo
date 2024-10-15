package service

import (
	"codebase-app/internal/module/user/entity"
	"codebase-app/internal/module/user/ports"
	"codebase-app/pkg"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/jwthandler"
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

var _ ports.UserService = &userService{}

type userService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *userService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Login(ctx context.Context, req *entity.LoginReq) (*entity.LoginRes, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// compare password
	match := pkg.ComparePassword(user.Password, req.Password)
	if !match {
		log.Warn().Any("req", req).Msg("service::Login - invalid credential")
		return nil, errmsg.NewCustomErrors(400).SetMessage("Kredensial tidak valid")
	}

	var (
		expiredAt = time.Now().Add(time.Hour * 12)
		payload   = jwthandler.CostumClaimsPayload{
			UserId:          user.Id,
			Role:            user.RoleName,
			TokenExpiration: expiredAt,
		}
	)

	token, err := jwthandler.GenerateTokenString(payload)
	if err != nil {
		log.Error().Err(err).Any("req", req).Msg("service::Login - failed to generate token")
		return nil, errmsg.NewCustomErrors(500).SetMessage("Internal Server Error")
	}

	resp := new(entity.LoginRes)
	resp.Id = user.Id
	resp.Name = user.Name
	resp.RoleName = user.RoleName
	resp.CompanyName = user.CompanyName
	resp.BranchName = user.BranchName
	resp.Token = token

	return resp, nil
}
