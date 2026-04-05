package handler

import (
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/delivery/http/response"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/model"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	useCase usecase.AuthUseCase
}

func NewAuthHandler(useCase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		useCase: useCase,
	}
}

func (h *AuthHandler) RegisterRoutes(router fiber.Router, jwtMiddleware fiber.Handler) {
	group := router.Group("/auth")
	group.Post("/login", h.Login)
	group.Post("/register", h.Register)
	group.Post("/logout", jwtMiddleware, h.Logout)
	group.Post("/refresh", h.Refresh)
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	resp, err := h.useCase.Login(c.Context(), &req)
	if err != nil {
		return err
	}

	return response.Send(c, fiber.StatusOK, "Login successful", resp)
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	resp, err := h.useCase.Register(c.Context(), &req)
	if err != nil {
		return err
	}

	return response.Send(c, fiber.StatusCreated, "User registered successfully", resp)
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	if err := h.useCase.Logout(c.Context(), c.Locals("user_id").(string)); err != nil {
		return err
	}

	return response.Send[any](c, fiber.StatusOK, "User logged out successfully", nil)
}

func (h *AuthHandler) Refresh(c fiber.Ctx) error {
	var req model.RefreshRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	resp, err := h.useCase.Refresh(c.Context(), req.RefreshToken)
	if err != nil {
		return err
	}

	return response.Send(c, fiber.StatusOK, "Token refreshed successfully", resp)
}
