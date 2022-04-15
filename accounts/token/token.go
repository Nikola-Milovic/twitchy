package token

import (
	"fmt"
	"nikolamilovic/twitchy/accounts/model"
	"time"

	"github.com/golang-jwt/jwt"
)

//TODO
var secret = "test secret"

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

func GenerateTokens(userId int) (string, error) {
	claims := model.UserClaims{
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
