package middleware

import (
	"log/slog"
	"time"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/configs"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func LoggerMiddleware(cfg *configs.AppConfig) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		var reqBody []byte
		if cfg.Env == "dev" && (c.Method() == "POST" || c.Method() == "PUT" || c.Method() == "PATCH" || c.Method() == "DELETE") {
			reqBody = c.Body()
		}

		chainErr := c.Next()

		if chainErr != nil {
			_ = c.App().Config().ErrorHandler(c, chainErr)
		}

		stop := time.Since(start)
		status := c.Response().StatusCode()
		reqID := requestid.FromContext(c)

		attrs := []slog.Attr{
			slog.String("request_id", reqID),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", status),
			slog.Duration("latency", stop),
			slog.String("ip", c.IP()),
			slog.String("ua", c.Get("User-Agent")),
		}

		if cfg.Env == "dev" {
			if len(reqBody) > 0 {
				attrs = append(attrs, slog.String("body", string(reqBody)))
			}
			if chainErr != nil {
				attrs = append(attrs, slog.String("err_detail", chainErr.Error()))
			}
		}

		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}

		slog.LogAttrs(c.Context(), level, "http_request", attrs...)

		return nil
	}
}
