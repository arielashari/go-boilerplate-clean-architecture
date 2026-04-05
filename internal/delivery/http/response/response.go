package response

import (
	"net/http"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/model"
	"github.com/gofiber/fiber/v3"
)

func Send[T any](c fiber.Ctx, code int, message string, data T) error {
	return c.Status(code).JSON(model.WebResponse[T]{
		Code:    code,
		Status:  http.StatusText(code),
		Message: message,
		Data:    data,
	})
}

func SendPaging[T any](c fiber.Ctx, code int, message string, data T, paging *model.PageMetadata) error {
	return c.Status(code).JSON(model.WebResponse[T]{
		Code:    code,
		Status:  http.StatusText(code),
		Message: message,
		Data:    data,
		Paging:  paging,
	})
}

func Error(c fiber.Ctx, code int, message string, errors any) error {
	return c.Status(code).JSON(model.WebResponse[any]{
		Code:    code,
		Status:  http.StatusText(code),
		Message: message,
		Errors:  errors,
	})
}
