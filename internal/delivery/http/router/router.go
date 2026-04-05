package router

import (
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/configs"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/delivery/http/handler"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/delivery/http/middleware"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

type Router struct {
	App         *fiber.App
	UserHandler *handler.UserHandler
	AuthHandler *handler.AuthHandler
	RoleHandler *handler.RoleHandler
	cfg         *configs.Config
	authUseCase usecase.AuthUseCase
}

func NewRouter(app *fiber.App, userHandler *handler.UserHandler, authHandler *handler.AuthHandler, roleHandler *handler.RoleHandler, cfg *configs.Config, authUseCase usecase.AuthUseCase) *Router {
	return &Router{
		App:         app,
		UserHandler: userHandler,
		AuthHandler: authHandler,
		RoleHandler: roleHandler,
		cfg:         cfg,
		authUseCase: authUseCase,
	}
}

func (r *Router) Setup() {

	// Serve API reference at /documentation
	r.App.Get("/documentation", func(c fiber.Ctx) error {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./api/openapi.yaml",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Go Boilerplate API Reference",
			},
			HideDownloadButton: true,
			DarkMode:           true,
			Theme:              scalar.ThemeMoon,
		})

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		c.Set("Content-Type", "text/html")
		return c.SendString(htmlContent)
	})

	// API Versioning Group
	v1 := r.App.Group("/api/v1")
	r.AuthHandler.RegisterRoutes(v1, middleware.JWTMiddleware(&r.cfg.JWT, r.authUseCase))

	protected := v1.Group("/", middleware.JWTMiddleware(&r.cfg.JWT, r.authUseCase))
	r.UserHandler.RegisterRoutes(protected)
	r.RoleHandler.RegisterRoutes(protected)
}
