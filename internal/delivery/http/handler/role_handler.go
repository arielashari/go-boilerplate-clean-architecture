package handler

import (
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/delivery/http/response"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/model"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

type RoleHandler struct {
	useCase usecase.RoleUseCase
}

func NewRoleHandler(useCase usecase.RoleUseCase) *RoleHandler {
	return &RoleHandler{useCase: useCase}
}

func (h *RoleHandler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/roles")
	group.Post("", h.Create)
	group.Get("/", h.List)
	group.Get("/:id", h.GetByID)
	group.Patch("/:id", h.Update)
	group.Delete("/:id", h.Delete)
}

func (h *RoleHandler) Create(c fiber.Ctx) error {
	var req model.CreateRoleRequest

	if err := c.Bind().Body(&req); err != nil {
		return err
	}
	resp, err := h.useCase.Create(c.Context(), &req)
	if err != nil {
		return err
	}

	return response.Send(c, fiber.StatusCreated, "Role created successfully", resp)
}

func (h *RoleHandler) Update(c fiber.Ctx) error {
	var req model.UpdateRoleRequest

	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	resp, err := h.useCase.Update(c.Context(), &req)
	if err != nil {
		return err
	}

	return response.Send(c, fiber.StatusOK, "Role updated successfully", resp)
}

func (h *RoleHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	resp, err := h.useCase.GetByID(c.Context(), id)
	if err != nil {
		return err
	}

	return response.Send(c, fiber.StatusOK, "Role retrieved successfully", resp)
}

func (h *RoleHandler) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	err := h.useCase.Delete(c.Context(), id)
	if err != nil {
		return err
	}

	return response.Send[any](c, fiber.StatusOK, "Role deleted successfully", nil)
}

func (h *RoleHandler) List(c fiber.Ctx) error {
	resp, err := h.useCase.List(c.Context(), 1, 10)
	if err != nil {
		return err
	}

	return response.Send(c, fiber.StatusOK, "Roles retrieved successfully", resp)
}
