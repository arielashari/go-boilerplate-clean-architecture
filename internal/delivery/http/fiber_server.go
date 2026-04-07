package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/configs"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/delivery/http/middleware"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/delivery/http/response"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/pkg/apperror"
	customervalidator "github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/pkg/validator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

type fiberServer struct {
	app *fiber.App
	cfg *configs.Config
}

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func NewFiberServer(cfg configs.Config) Server {

	fs := &fiberServer{
		cfg: &cfg,
	}

	fiberConfig := fiber.Config{
		AppName:         cfg.App.Name,
		StructValidator: &structValidator{validate: customervalidator.GetValidator()},
		ErrorHandler:    fs.GlobalErrorHandler,
	}

	fs.app = fiber.New(fiberConfig)

	fs.app.Use(helmet.New())
	fs.app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORS.AllowOrigins,
		AllowHeaders: cfg.CORS.AllowHeaders,
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}))
	fs.app.Use(limiter.New(limiter.Config{
		Max:        cfg.RateLimit.MaxRequests,
		Expiration: time.Duration(cfg.RateLimit.ExpirationSeconds) * time.Second,
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c fiber.Ctx) error {
			return response.Error(c, fiber.StatusTooManyRequests, "too many requests", nil)
		},
	}))
	fs.app.Use(requestid.New())
	fs.app.Use(middleware.MetricsMiddleware())
	fs.app.Use(middleware.LoggerMiddleware(&fs.cfg.App))

	return fs
}

func (fs *fiberServer) Start() {
	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("Server is starting", "port", fs.cfg.App.Port)
		err := fs.app.Listen(fmt.Sprintf(":%d", fs.cfg.App.Port), fiber.ListenConfig{
			DisableStartupMessage: true,
		})
		if err != nil {
			serverErrors <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		slog.Error("Startup failed", "error", err)
	case sig := <-shutdown:
		slog.Info("Signal received, shutting down...", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := fs.app.ShutdownWithContext(ctx); err != nil {
			slog.Error("Graceful shutdown failed", "error", err)
		} else {
			slog.Info("Server shutdown complete")
		}
	}
}

func (fs *fiberServer) GetFiberApp() *fiber.App {
	return fs.app
}

func (fs *fiberServer) GlobalErrorHandler(c fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	var valErrs validator.ValidationErrors
	if errors.As(err, &valErrs) {
		formatted := customervalidator.FormatValidationError(valErrs)
		return response.Error(c, fiber.StatusBadRequest, "Validation failed", formatted)
	}

	if appErr, ok := apperror.As(err); ok {
		statusCode := mapCodeToStatus(appErr.Code)
		fs.logAppError(c.Context(), appErr, statusCode)
		return response.Error(c, statusCode, appErr.Message, nil)
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		slog.Warn("fiber infrastructure error",
			slog.Int("status", fiberErr.Code),
			slog.String("message", fiberErr.Message),
		)
		return response.Error(c, fiberErr.Code, fiberErr.Message, nil)
	}

	slog.Error("unhandled system error", slog.String("error", err.Error()))
	return response.Error(c, fiber.StatusInternalServerError, "Internal server error", nil)
}

func mapCodeToStatus(code apperror.ErrorCode) int {
	switch code {
	case apperror.CodeValidation:
		return fiber.StatusBadRequest
	case apperror.CodeUnauthorized, apperror.CodeInvalidCreds:
		return fiber.StatusUnauthorized
	case apperror.CodeForbidden:
		return fiber.StatusForbidden
	case apperror.CodeNotFound:
		return fiber.StatusNotFound
	case apperror.CodeConflict, apperror.CodeEmailTaken:
		return fiber.StatusConflict
	case apperror.CodeTimeout:
		return fiber.StatusRequestTimeout
	case apperror.CodeUnavailable:
		return fiber.StatusServiceUnavailable
	default:
		return fiber.StatusInternalServerError
	}
}

func (fs *fiberServer) logAppError(ctx context.Context, appErr *apperror.AppError, status int) {
	level := slog.LevelWarn
	if status >= 500 {
		level = slog.LevelError
	}

	slog.Log(ctx, level, "application error",
		"code", string(appErr.Code),
		"message", appErr.Message,
		"status", status,
		"operation", appErr.Operation,
		"internal", appErr.Internal,
	)
}
