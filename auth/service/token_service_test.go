package service

import (
	"errors"
	tok "nikolamilovic/twitchy/common/token"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

func TestExpiredToken(t *testing.T) {
	claims := tok.UserClaims{
		UserId: 1,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Add(time.Minute * -5).Unix()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)

	if err != nil {
		t.Fatalf("Expected error to be nil, got %v", err.Error())
	}

	isValid, err := CheckJWTToken(tokenString)
	if isValid {
		t.Fatalf("Expected JWT to be invalid, got valid")
	}

	if !errors.Is(err, tok.InvalidJWTError) {
		t.Fatalf("Expected error to be %e when checking JWT, got %e", tok.InvalidJWTError, err)
	}
}

func TestGenarateNewTokens(t *testing.T) {
	jwt, refresh, err := generateTokens(1)
	if err != nil {
		t.Fatalf("Expected error to be nil, got %v", err.Error())
	}

	if len(refresh) != 128 {
		t.Fatalf("Expected refresh token to be 128 characters long, got %d", len(refresh))
	}

	isValid, err := CheckJWTToken(jwt)

	if err != nil {
		t.Fatalf("Expected error to be nil when checking JWT, got %v", err.Error())
	}

	if !isValid {
		t.Fatalf("Expected JWT to be valid, got invalid")
	}

	isValid, err = CheckJWTToken(jwt + "a")

	if err == nil {
		t.Fatalf("Expected error to be not nil when checking JWT, got nil")
	}

	if isValid {
		t.Fatalf("Expected JWT to be invalid, got valid")
	}
}

//	mock, err := pgxmock.NewConn()
// if err != nil {
// 	t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// }
// defer mock.Close(context.Background())

// sut := &AuthService{
// 	DB: mock,
// }

// hashedPassword, _ := hashPassword("password")

// // before we actually execute our api function, we need to expect required DB actions
// rows := pgxmock.NewRows([]string{"id", "password"}).AddRow(1, hashedPassword)

// mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)

// id, err := sut.CheckLogin("test@gmail.com", "wrongpassword")

// if id != -1 {
// 	t.Fatalf("Expected id to be %d got %d", -1, id)
// }

// if !errors.Is(err, model.WrongPasswordError) {
// 	t.Fatalf("wrong error , expected %e, got %e", model.WrongPasswordError, err)
// }
