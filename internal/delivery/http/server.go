package http

import (
	"github.com/gofiber/fiber/v3"
)

type Server interface {
	Start()
	GetFiberApp() *fiber.App
}
