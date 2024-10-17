package handler

import (
	"codebase-app/internal/adapter"
	m "codebase-app/internal/middleware"
	"codebase-app/internal/module/delivery/entity"
	"codebase-app/internal/module/delivery/ports"
	"codebase-app/internal/module/delivery/repository"
	"codebase-app/internal/module/delivery/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type deliveryHandler struct {
	service ports.DeliveryService
}

func NewDeliveryHandler() *deliveryHandler {
	var (
		repo    = repository.NewDeliveryRepository()
		service = service.NewDeliveryService(repo)
		handler = new(deliveryHandler)
	)
	handler.service = service

	return handler
}

func (h *deliveryHandler) Register(router fiber.Router) {
	router.Post("/",
		m.AuthBearer,
		m.AuthRole([]string{"admin"}),
		h.createSchedules,
	)
}

func (h *deliveryHandler) createSchedules(c *fiber.Ctx) error {
	var (
		req = new(entity.CreateSchedulesReq)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = m.GetLocals(c)
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createSchedules - invalid request")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.UserId = l.UserId

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("req", req).Msg("handler::createSchedules - invalid request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.CreateSchedules(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(resp, ""))

}
