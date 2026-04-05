package http

import (
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/delivery/http/router"
	"github.com/gofiber/fiber/v3"
)

type Server interface {
	Start()
	GetFiberApp() *fiber.App
	RegisterRoutes(r *router.Router)
}
