package token

import (
	"fmt"

	"errors"

	"github.com/golang-jwt/jwt"
)

var InvalidJWTError = errors.New("Invalid JWT Token")

func CheckJWTToken(tokenString string, secret []byte) (bool, error) {
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
		return false, InvalidJWTError
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
