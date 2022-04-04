package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"nikolamilovic/twitchy/auth/db"
	"nikolamilovic/twitchy/auth/model"
	"time"

	"github.com/golang-jwt/jwt"
)

var secret = []byte("test secret")
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type ITokenService interface {
	RefreshToken(refreshTokenString string) (string, string, error)
	GenerateNewTokensForUser(userId int) (string, string, error)
}

type TokenService struct {
	DB db.PgxIface
}

func (s *TokenService) RefreshToken(refreshTokenString string) (string, string, error) {
	res, err := s.DB.Query(context.Background(), "SELECT user_id, token, expires FROM refresh_tokens WHERE token = $1", refreshTokenString)

	if err != nil {
		return "", "", fmt.Errorf("Refresh token: %w", err)
	}

	if res.CommandTag().RowsAffected() == 0 || !res.Next() {
		return "", "", fmt.Errorf("No RefreshToken: %w", errors.New("No rows returned"))
	}

	var refreshToken model.RefreshToken
	err = res.Scan(&refreshToken.UserId, &refreshToken.Token, &refreshToken.Expires)

	if err != nil {
		return "", "", fmt.Errorf("RefreshToken: %w", err)
	}

	verifyRefreshToken(refreshToken)

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

func (s *TokenService) saveRefreshToken(token string, userId int) error {
	expiresAt := time.Now().Add(time.Hour * 24 * 7)
	res, err := s.DB.Exec(context.Background(), "INSERT INTO refresh_tokens (user_id, token, expires) VALUES ($1, $2, $3)", userId, token, expiresAt)

	if err != nil {
		return fmt.Errorf("Insert refresh token: %w", err)
	}
	if res.RowsAffected() > 0 {
		return nil
	} else {
		return fmt.Errorf("Insert refresh token: %w", errors.New("No rows affected"))
	}
}

func generateTokens(userId int) (string, string, error) {
	claims := model.UserClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)

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

func CheckJWTToken(tokenString string) (bool, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if !token.Valid {
		return false, model.InvalidJWTError
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		err := claims.Valid()
		if err != nil {
			return false, fmt.Errorf("Claims Invalid: %w", err)
		}
	} else {
		return false, fmt.Errorf("Claims Invalid: %w", err)
	}

	return true, nil
}
