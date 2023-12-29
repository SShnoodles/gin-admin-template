package util

import (
	"errors"
	"gin-admin-template/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

func ValidateUserToken(tokenString string, signSey string) (*jwt.RegisteredClaims, error) {
	if tokenString == "" {
		return nil, errors.New("Token is empty!")
	}
	tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")
	registeredClaims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &registeredClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(signSey), nil
	})
	if err != nil && !token.Valid {
		return nil, err
	}
	return &registeredClaims, nil
}

func GenerateToken(subject int64) (string, error) {
	registeredClaims := jwt.RegisteredClaims{
		ID:        uuid.New().String(),
		Subject:   strconv.FormatInt(subject, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, 7)),
		Issuer:    "gin",
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims).SignedString([]byte(config.AppConfig.Jwt.Secret))
}
