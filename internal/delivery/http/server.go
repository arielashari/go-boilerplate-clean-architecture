package http

import (
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/configs"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/delivery/http/handler"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/usecase"
)

type Server interface {
	Start()
	RegisterRoutes(userHandler *handler.UserHandler, authHandler *handler.AuthHandler, roleHandler *handler.RoleHandler, fileHandler *handler.FileHandler, cfg *configs.Config, authUseCase usecase.AuthUseCase)
}
