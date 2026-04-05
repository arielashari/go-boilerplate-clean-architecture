package handler

import (
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/delivery/http/response"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/model"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UserHandler struct {
	useCase usecase.UserUseCase
}

func (h *UserHandler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/users")
	group.Post("", h.Create)
	group.Get("/", h.List)
	group.Get("/:id", h.GetByID)
	group.Patch("/:id", h.Update)
	group.Delete("/:id", h.Delete)
}

func NewUserHandler(useCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		useCase: useCase,
	}
}

func (h *UserHandler) Create(c fiber.Ctx) error {
	var req model.CreateUserRequest

	if err := c.Bind().Body(&req); err != nil {
		return err
	}
	resp, err := h.useCase.Create(c.Context(), &req)
	if err != nil {
		return err
	}

	return response.Send(c, fiber.StatusCreated, "User created successfully", resp)
}

func (h *UserHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")

	_, err := uuid.Parse(id)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user ID format", nil)
	}

	resp, err := h.useCase.GetByID(c.Context(), id)
	if err != nil {
		return err
	}

	return response.Send(c, fiber.StatusOK, "User retrieved successfully", resp)
}

func (h *UserHandler) List(c fiber.Ctx) error {
	request := new(model.GetAllUsersRequest)
	if err := c.Bind().Query(request); err != nil {
		return err
	}

	resp, err := h.useCase.List(c.Context(), request)
	if err != nil {
		return err
	}

	return response.SendPaging(c, fiber.StatusOK, "Users retrieved successfully", resp.Items, &resp.Metadata)
}

func (h *UserHandler) Update(c fiber.Ctx) error {
	var req model.UpdateUserRequest

	id := c.Params("id")

	_, err := uuid.Parse(id)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user ID format", nil)
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	req.ID = id

	resp, err := h.useCase.Update(c.Context(), &req)
	if err != nil {
		return err
	}

	return response.Send(c, fiber.StatusOK, "User updated successfully", resp)
}

func (h *UserHandler) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	_, err := uuid.Parse(id)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user ID format", nil)
	}
	if err := h.useCase.Delete(c.Context(), id); err != nil {
		return err
	}
	return response.Send[any](c, fiber.StatusOK, "User deleted successfully", nil)
}
