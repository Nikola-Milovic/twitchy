package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"nikolamilovic/twitchy/auth/model"
	db "nikolamilovic/twitchy/common/db"
	tok "nikolamilovic/twitchy/common/token"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type ITokenService interface {
	RefreshToken(refreshTokenString string) (string, string, error)
	GenerateNewTokensForUser(userId int) (string, string, error)
}

type TokenService struct {
	DB db.PgxIface
}

func (s *TokenService) RefreshToken(refreshTokenString string) (string, string, error) {

	refreshToken, err := s.fetchRefreshToken(refreshTokenString)
	if err != nil {
		return "", "", fmt.Errorf("RefreshToken: %w", err)
	}

	isValid := verifyRefreshToken(refreshToken)

	if !isValid {
		return "", "", fmt.Errorf("RefreshToken: %w", errors.New("Refresh token is not valid"))
	}

	jwt, refresh, err := generateTokens(refreshToken.UserId)

	if err != nil {
		return "", "", fmt.Errorf("RefreshToken: %w", err)
	}

	err = s.saveRefreshToken(refresh, refreshToken.UserId)

	if err != nil {
		return "", "", fmt.Errorf("RefreshToken: %w", err)
	}

	return jwt, refresh, nil
}

//Returns JWT, RefreshToken, error
func (s *TokenService) GenerateNewTokensForUser(userId int) (string, string, error) {
	jwt, refresh, err := generateTokens(userId)

	if err != nil {
		return "", "", err
	}

	err = s.saveRefreshToken(refresh, userId)

	if err != nil {
		return "", "", fmt.Errorf("GenerateNewTokensForUser: %w", err)
	}

	return jwt, refresh, nil
}

//get Refresh token
func (s *TokenService) fetchRefreshToken(refreshTokenString string) (model.RefreshToken, error) {
	res, err := s.DB.Query(context.Background(), "SELECT user_id, token, expires FROM refresh_tokens WHERE token = $1", refreshTokenString)

	if err != nil {
		return model.RefreshToken{}, fmt.Errorf("fetchRefreshToken: %w", err)
	}

	if !res.Next() {
		return model.RefreshToken{}, fmt.Errorf("fetchRefreshToken: %w", errors.New("No rows returned"))
	}

	var refreshToken model.RefreshToken
	err = res.Scan(&refreshToken.UserId, &refreshToken.Token, &refreshToken.Expires)

	if err != nil {
		return model.RefreshToken{}, fmt.Errorf("fetchRefreshToken: %w", err)
	}

	return refreshToken, nil
}

func (s *TokenService) saveRefreshToken(token string, userId int) error {
	expiresAt := time.Now().Add(time.Hour * 24 * 7).Unix()
	//INSERT if refresh token for given user doesnt exist already, otherwise update
	res, err := s.DB.Exec(context.Background(), "INSERT INTO refresh_tokens (user_id, token, expires) VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET token = $2, expires = $3", userId, token, expiresAt)

	if err != nil {
		return fmt.Errorf("saveRefreshToken: %w", err)
	}
	if res.RowsAffected() > 0 {
		return nil
	} else {
		return fmt.Errorf("saveRefreshToken: %w", errors.New("No rows affected"))
	}
}

func generateTokens(userId int) (string, string, error) {
	claims := tok.UserClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			Issuer:    "twitchy", //TODO
			IssuedAt:  time.Now().Unix(),
			Subject:   fmt.Sprintf("%d", userId),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))

	b := make([]rune, 128)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	refreshToken := string(b)

	if err != nil {
		return "", "", fmt.Errorf("GenerateTokens: %w", err)
	}

	return tokenString, refreshToken, nil
}

func verifyRefreshToken(token model.RefreshToken) bool {
	expires := time.Unix(int64(token.Expires), 0)
	if expires.Before(time.Now()) {
		return false
	}

	return true
}
