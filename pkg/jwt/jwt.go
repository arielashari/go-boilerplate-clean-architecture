package jwt

import (
	"time"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/configs"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenID      string
	AtExpires    int64
	RtExpires    int64
}

func GenerateTokenPair(userID string, cfg *configs.JWTConfig) (*TokenDetails, error) {
	td := &TokenDetails{
		TokenID:   uuid.New().String(),
		AtExpires: time.Now().Add(time.Minute * time.Duration(cfg.AccessExpireMinutes)).Unix(),
		RtExpires: time.Now().Add(time.Minute * time.Duration(cfg.RefreshExpireMinutes)).Unix(),
	}

	atClaims := jwt.MapClaims{
		"sub":  userID,
		"jti":  td.TokenID,
		"exp":  td.AtExpires,
		"type": "access",
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	var err error
	td.AccessToken, err = at.SignedString([]byte(cfg.Secret))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{
		"sub":  userID,
		"jti":  td.TokenID,
		"exp":  td.RtExpires,
		"type": "refresh",
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(cfg.Secret))

	return td, err
}
