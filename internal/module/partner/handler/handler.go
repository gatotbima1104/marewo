package handler

import (
	"codebase-app/internal/adapter"
	m "codebase-app/internal/middleware"
	"codebase-app/internal/module/partner/entity"
	"codebase-app/internal/module/partner/ports"
	"codebase-app/internal/module/partner/repository"
	"codebase-app/internal/module/partner/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type partnerHandler struct {
	service ports.PartnerService
}

func NewPartnerHandler() *partnerHandler {
	var (
		repo    = repository.NewPartnerRepository()
		service = service.NewPartnerService(repo)
		handler = new(partnerHandler)
	)
	handler.service = service

	return handler
}

func (h *partnerHandler) Register(router fiber.Router) {
	router.Post("/",
		m.AuthBearer,
		m.AuthRole([]string{"admin"}),
		h.createPartner,
	)

	router.Get("/",
		m.AuthBearer,
		h.getPartners,
	)
}

func (h *partnerHandler) createPartner(c *fiber.Ctx) error {
	var (
		req = new(entity.CreatePartnerReq)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = m.GetLocals(c)
	)

	if err := c.BodyParser(req); err != nil {
		log.Error().Err(err).Msg("handler::createPartner - invalid request")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.UserId = l.UserId
	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Error().Err(err).Any("req", req).Msg("handler::createPartner - invalid request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.CreatePartner(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(resp, ""))
}

func (h *partnerHandler) getPartners(c *fiber.Ctx) error {
	var (
		req = new(entity.GetPartnersReq)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = m.GetLocals(c)
	)

	if err := c.QueryParser(req); err != nil {
		log.Error().Err(err).Msg("handler::getPartners - invalid request")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.UserId = l.UserId
	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Error().Err(err).Any("req", req).Msg("handler::getPartners - invalid request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.GetPartners(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(resp, ""))
}
