package handler

import (
	"codebase-app/internal/adapter"
	m "codebase-app/internal/middleware"
	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	"codebase-app/internal/module/product/repository"
	"codebase-app/internal/module/product/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type productHandler struct {
	service ports.ProductService
}

func NewProduct() *productHandler {
	var (
		repo    = repository.NewProductRepository()
		service = service.NewProductService(repo)
		handler = new(productHandler)
	)
	handler.service = service

	return handler
}

func (h *productHandler) Register(router fiber.Router) {
	router.Post("/",
		m.AuthBearer,
		m.AuthRole([]string{"admin"}),
		h.CreateProduct,
	)

	router.Get("/",
		m.AuthBearer,
		m.AuthRole([]string{"admin"}),
		h.GetProducts,
	)

	router.Get("/:id", h.GetProduct)
}

func (h *productHandler) CreateProduct(c *fiber.Ctx) error {
	var (
		req = new(entity.CreateProductReq)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = m.GetLocals(c)
	)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.UserId = l.UserId

	if err := v.Validate(req); err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.CreateProduct(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(res, ""))
}

func (h *productHandler) GetProduct(c *fiber.Ctx) error {
	var (
		req = new(entity.GetProductReq)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	req.Id = c.Params("id")

	if err := v.Validate(req); err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.GetProduct(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(res, ""))
}

func (h *productHandler) GetProducts(c *fiber.Ctx) error {
	var (
		req = new(entity.GetProductsReq)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = m.GetLocals(c)
	)

	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.UserId = l.UserId
	req.SetDefault()

	if err := v.Validate(req); err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.GetProducts(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(res, ""))
}
