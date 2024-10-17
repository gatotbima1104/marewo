package handler

import (
	"codebase-app/internal/adapter"
	m "codebase-app/internal/middleware"
	"codebase-app/internal/module/user/entity"
	"codebase-app/internal/module/user/ports"
	"codebase-app/internal/module/user/repository"
	"codebase-app/internal/module/user/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type userHandler struct {
	service ports.UserService
}

func NewUserHandler() *userHandler {
	var (
		repo    = repository.NewUserRepository()
		service = service.NewUserService(repo)
		handler = new(userHandler)
	)
	handler.service = service

	return handler
}

func (h *userHandler) Register(router fiber.Router) {
	router.Post("/login", h.login)
	router.Get("/couriers",
		m.AuthBearer,
		h.getCouriers)
}

func (h *userHandler) login(c *fiber.Ctx) error {
	var (
		req = new(entity.LoginReq)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::Login - invalid request")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("handler::Login - invalid request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.Login(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(res, ""))
}

func (h *userHandler) getCouriers(c *fiber.Ctx) error {
	var (
		req = new(entity.GetCouriersReq)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = m.GetLocals(c)
	)

	if err := c.QueryParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetCouriers - invalid request")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefault()
	req.UserId = l.UserId

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("handler::GetCouriers - invalid request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.GetCouriers(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(res, ""))
}
