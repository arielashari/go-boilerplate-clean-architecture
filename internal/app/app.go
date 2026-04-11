package app

import (
	"log/slog"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/configs"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/delivery/http"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/delivery/http/handler"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/repository/postgres"
	_redis "github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/repository/redis"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/usecase"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/pkg/database"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/pkg/logger"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/pkg/mailer"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/pkg/storage"
)

type App struct {
	Config        *configs.Config
	Server        http.Server
	DatabaseInfra *database.DatabaseInfrastructure
}

func NewApp(cfg *configs.Config) *App {
	logger.InitLogger(&cfg.App)
	log := slog.Default()

	pg, err := database.NewPostgresConnection(*cfg, log)

	if err := database.RunMigrations(pg); err != nil {
		log.Error("migration failed", "error", err)
		panic(err)
	}

	if err != nil {
		log.Error("failed to connect to database", "error", err)
		panic(err)
	}

	mailer, err := mailer.NewSMTPMailer(&cfg.SMTP)
	if err != nil {
		log.Error("failed to connect to mailer service", "error", err)
		panic(err)
	}

	redis, err := database.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Error("failed to connect to redis", "error", err)
		panic(err)
	}

	dbInfra := &database.DatabaseInfrastructure{
		Postgres: pg,
		Redis:    redis,
	}

	transactor := postgres.NewTransactor(pg)
	userRepo := postgres.NewUserPostgresRepository(pg)
	roleRepo := postgres.NewRolePostgresRepository(pg)
	authRepo := _redis.NewAuthRedisRepository(redis)

	userUseCase := usecase.NewUserUseCase(userRepo)
	roleUseCase := usecase.NewRoleUseCase(roleRepo)

	authUseCase := usecase.NewAuthUseCase(authRepo, userRepo, transactor, mailer, cfg)

	authHandler := handler.NewAuthHandler(authUseCase)
	userHandler := handler.NewUserHandler(userUseCase)
	roleHandler := handler.NewRoleHandler(roleUseCase)

	s3Storage, err := storage.NewS3Storage(&cfg.S3)
	if err != nil {
		log.Error("failed to initialize S3 client", "error", err)
		panic(err)
	}

	fileUploadUseCase := usecase.NewFileUploadUseCase(s3Storage)
	fileHandler := handler.NewFileHandler(fileUploadUseCase)

	httpserver := http.NewFiberServer(*cfg)
	httpserver.RegisterRoutes(userHandler, authHandler, roleHandler, fileHandler, cfg, authUseCase)

	return &App{
		Config:        cfg,
		DatabaseInfra: dbInfra,
		Server:        httpserver,
	}
}
