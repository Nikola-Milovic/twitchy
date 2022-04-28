package test_util

import (
	"fmt"
	"nikolamilovic/twitchy/common/token"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateTokens(userId int, secret string) (string, error) {
	claims := token.UserClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			Issuer:    "twitchy", //TODO
			IssuedAt:  time.Now().Unix(),
			Subject:   fmt.Sprintf("%d", userId),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", fmt.Errorf("GenerateTokens: %w", err)
	}

	return tokenString, nil
}
