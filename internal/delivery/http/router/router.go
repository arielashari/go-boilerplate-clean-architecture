package router

import (
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/configs"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/delivery/http/handler"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/delivery/http/middleware"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type Router struct {
	App         *fiber.App
	UserHandler *handler.UserHandler
	AuthHandler *handler.AuthHandler
	RoleHandler *handler.RoleHandler
	FileHandler *handler.FileHandler
	cfg         *configs.Config
	authUseCase usecase.AuthUseCase
}

func NewRouter(app *fiber.App, userHandler *handler.UserHandler, authHandler *handler.AuthHandler, roleHandler *handler.RoleHandler, fileHandler *handler.FileHandler, cfg *configs.Config, authUseCase usecase.AuthUseCase) *Router {
	return &Router{
		App:         app,
		UserHandler: userHandler,
		AuthHandler: authHandler,
		RoleHandler: roleHandler,
		FileHandler: fileHandler,
		cfg:         cfg,
		authUseCase: authUseCase,
	}
}

func (r *Router) Setup() {

	r.App.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"env":    r.cfg.App.Env,
		})
	})

	r.App.Get("/ready", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ready",
		})
	})

	r.App.Get("/metrics", func(c fiber.Ctx) error {
		fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())(c.RequestCtx())
		return nil
	})

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
	r.FileHandler.RegisterRoutes(protected)
}
