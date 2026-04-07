package usecase

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/configs"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/entity"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/model"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/model/mapper"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/pkg/apperror"
	_jwt "github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/pkg/jwt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Login(ctx context.Context, request *model.LoginRequest) (*model.AuthResponse, error)
	Register(ctx context.Context, request *model.RegisterRequest) (*model.UserResponse, error)
	Logout(ctx context.Context, userID string) error
	Refresh(ctx context.Context, refreshToken string) (*model.AuthResponse, error)
	ValidateSession(ctx context.Context, userID, tokenID string) (bool, error)
}

type authUseCase struct {
	authRepository entity.AuthRepository
	userRepository entity.UserRepository
	transactor     entity.Transactor
	mailer         entity.EmailSender
	cfg            *configs.Config
}

func NewAuthUseCase(authRepo entity.AuthRepository, userRepo entity.UserRepository, transactor entity.Transactor, mailer entity.EmailSender, cfg *configs.Config) AuthUseCase {
	return &authUseCase{
		authRepository: authRepo,
		userRepository: userRepo,
		transactor:     transactor,
		mailer:         mailer,
		cfg:            cfg,
	}
}

func (a *authUseCase) Login(ctx context.Context, request *model.LoginRequest) (*model.AuthResponse, error) {
	user, err := a.userRepository.GetUserForAuth(ctx, request.Email)
	if err != nil {
		return nil, entity.ErrInvalidCredentials.WithOperation("AuthUseCase.Login.GetUser")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, entity.ErrInvalidCredentials.WithInternal(err).WithOperation("AuthUseCase.Login.PasswordCheck")
	}

	td, err := _jwt.GenerateTokenPair(user.ID, &a.cfg.JWT)
	if err != nil {
		return nil, apperror.New(apperror.CodeInternal, "failed to generate tokens").
			WithInternal(err).
			WithOperation("AuthUseCase.Login.GenerateTokens")
	}

	err = a.authRepository.SetSession(ctx, user.ID, td.TokenID, time.Hour*24*7)
	if err != nil {
		return nil, apperror.New(apperror.CodeInternal, "failed to create session").
			WithInternal(err).
			WithOperation("AuthUseCase.Login.SetSession")
	}

	return &model.AuthResponse{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
	}, nil
}

func (a *authUseCase) Register(ctx context.Context, request *model.RegisterRequest) (*model.UserResponse, error) {
	existingUser, err := a.userRepository.GetByEmail(ctx, request.Email)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, entity.ErrEmailAlreadyExists.WithOperation("AuthUseCase.Register")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.New(apperror.CodeInternal, "failed to hash password").
			WithInternal(err).
			WithOperation("AuthUseCase.Register.HashPassword")
	}

	user := &entity.User{
		ID:          uuid.NewString(),
		Email:       request.Email,
		Password:    string(hashedPassword),
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		PhonePrefix: request.PhonePrefix,
		PhoneNumber: request.PhoneNumber,
		RoleID:      request.RoleID,
	}

	var createdUser *entity.User
	err = a.transactor.WithTx(ctx, func(ctx context.Context) error {
		createdUser, err = a.userRepository.Create(ctx, user)
		return err
	})
	if err != nil {
		return nil, err
	}

	go func() {
		if err := a.mailer.SendWelcomeEmail(context.Background(), createdUser.Email, createdUser.FirstName); err != nil {
			slog.Warn("failed to send welcome email", "error", err, "userID", createdUser.ID)
		}
	}()

	return mapper.UserToResponse(createdUser), nil
}

func (a *authUseCase) Logout(ctx context.Context, userID string) error {
	return a.authRepository.DeleteSession(ctx, userID)
}

func (a *authUseCase) Refresh(ctx context.Context, refreshToken string) (*model.AuthResponse, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, entity.ErrUnauthorized.WithOperation("AuthUseCase.Refresh.InvalidSigningMethod")
		}
		return []byte(a.cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, entity.ErrUnauthorized.WithInternal(err).WithOperation("AuthUseCase.Refresh.ParseToken")
	}

	if !token.Valid {
		return nil, entity.ErrUnauthorized.WithOperation("AuthUseCase.Refresh.InvalidToken")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, entity.ErrUnauthorized.WithOperation("AuthUseCase.Refresh.InvalidClaims")
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, entity.ErrUnauthorized.WithOperation("AuthUseCase.Refresh.InvalidTokenType")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, entity.ErrUnauthorized.WithOperation("AuthUseCase.Refresh.MissingUserID")
	}

	tokenID, ok := claims["jti"].(string)
	if !ok {
		return nil, entity.ErrUnauthorized.WithOperation("AuthUseCase.Refresh.MissingTokenID")
	}

	isValid, err := a.authRepository.CheckSession(ctx, userID, tokenID)
	if err != nil {
		return nil, apperror.New(apperror.CodeInternal, "failed to validate session").
			WithInternal(err).
			WithOperation("AuthUseCase.Refresh.CheckSession")
	}

	if !isValid {
		return nil, apperror.New(apperror.CodeUnauthorized, "invalid session").
			WithOperation("AuthUseCase.Refresh.SessionInvalid")
	}

	td, err := _jwt.GenerateTokenPair(userID, &a.cfg.JWT)
	if err != nil {
		return nil, apperror.New(apperror.CodeInternal, "failed to generate tokens").
			WithInternal(err).
			WithOperation("AuthUseCase.Refresh.GenerateTokens")
	}

	err = a.authRepository.SetSession(ctx, userID, td.TokenID, time.Hour*24*7)
	if err != nil {
		return nil, apperror.New(apperror.CodeInternal, "failed to update session").
			WithInternal(err).
			WithOperation("AuthUseCase.Refresh.SetSession")
	}

	return &model.AuthResponse{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
	}, nil
}

func (a *authUseCase) ValidateSession(ctx context.Context, userID, tokenID string) (bool, error) {
	return a.authRepository.CheckSession(ctx, userID, tokenID)
}
