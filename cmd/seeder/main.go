package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/configs"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/entity"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/repository/postgres"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/pkg/database"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/pkg/logger"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	viper.Set("APP_ENV", env)

	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	logger.InitLogger(&cfg.App)
	log := slog.Default()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pg, err := database.NewPostgresConnection(*cfg, log)

	if err != nil {
		log.Error("failed to connect to database", "error", err)
		panic(err)
	}
	defer pg.Close()

	userRepo := postgres.NewUserPostgresRepository(pg)
	roleRepo := postgres.NewRolePostgresRepository(pg)
	seedRoles(ctx, roleRepo)
	seedUsers(ctx, userRepo, roleRepo)

}

func seedRoles(ctx context.Context, roleRepo entity.RolePostgresRepository) {
	roles := []string{"Admin", "User"}

	for _, roleName := range roles {
		role, err := roleRepo.GetByName(ctx, roleName)
		if err != nil && !errors.Is(err, entity.ErrNotFound) {
			log.Printf("failed to get role %s: %v", roleName, err)
			continue
		}
		if role != nil {
			log.Printf("role %s already exists, skipping", roleName)
			continue
		}

		newRole := &entity.Role{
			ID:   uuid.New().String(),
			Name: roleName,
		}
		_, err = roleRepo.Create(ctx, newRole)
		if err != nil {
			log.Printf("failed to create role %s: %v", roleName, err)
			continue
		}
		log.Printf("role %s created successfully", roleName)
	}
}

func seedUsers(ctx context.Context, userRepo entity.UserPostgresRepository, roleRepo entity.RolePostgresRepository) {
	adminRole, err := roleRepo.GetByName(ctx, "Admin")
	if err != nil || adminRole == nil && !errors.Is(err, entity.ErrNotFound) {
		log.Fatalf("cannot seed users: Admin role not found")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("@Dmin1234"), bcrypt.DefaultCost)
	users := []entity.User{
		{
			ID:          uuid.NewString(),
			Email:       "admin@gmail.com",
			Password:    string(hashedPassword),
			FirstName:   "Admin",
			LastName:    "User",
			PhonePrefix: "+62",
			PhoneNumber: "81213141516",
			RoleID:      adminRole.ID,
		},
	}
	for _, user := range users {
		_, err := userRepo.Create(ctx, &user)
		if err != nil {
			log.Printf("failed to create user %s: %v", user.Email, err)
			continue
		}
		log.Printf("user %s created successfully", user.Email)
	}
}
