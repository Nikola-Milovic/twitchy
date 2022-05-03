package service

import (
	"context"
	"errors"
	"nikolamilovic/twitchy/common/token"
	tok "nikolamilovic/twitchy/common/token"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pashagolub/pgxmock"
)

func TestExpiredToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test secret")
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

	isValid, err := tok.CheckJWTToken(tokenString, string(secret))
	if isValid {
		t.Fatalf("Expected JWT to be invalid, got valid")
	}

	if !errors.Is(err, tok.InvalidJWTError) {
		t.Fatalf("Expected error to be %e when checking JWT, got %e", tok.InvalidJWTError, err)
	}
}

func TestGenarateNewTokens(t *testing.T) {
	t.Setenv("JWT_SECRET", "test secret")
	jwt, refresh, err := generateTokens(1)
	if err != nil {
		t.Fatalf("Expected error to be nil, got %v", err.Error())
	}

	if len(refresh) != 128 {
		t.Fatalf("Expected refresh token to be 128 characters long, got %d", len(refresh))
	}

	isValid, err := tok.CheckJWTToken(jwt, string(secret))

	if err != nil {
		t.Fatalf("Expected error to be nil when checking JWT, got %v", err.Error())
	}

	if !isValid {
		t.Fatalf("Expected JWT to be valid, got invalid")
	}

	isValid, err = tok.CheckJWTToken(jwt+"a", string(secret))

	if err == nil {
		t.Fatalf("Expected error to be not nil when checking JWT, got nil")
	}

	if isValid {
		t.Fatalf("Expected JWT to be invalid, got valid")
	}
}

func TestRefreshToken(t *testing.T) {
	//Setup
	t.Setenv("JWT_SECRET", "test secret")

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	_, jwt, err := generateTokens(1)

	if err != nil {
		t.Fatalf("Expected error to be nil, got %v", err.Error())
	}

	rows := pgxmock.NewRows([]string{"user_id", "token", "expires"}).AddRow(1, jwt, time.Now().Add(time.Minute*5).Unix()).AddRow(1, jwt, time.Now().Add(time.Minute*5).Unix())

	mock.ExpectQuery("SELECT user_id, token, expires FROM refresh_tokens").WithArgs("correct_token").
		WillReturnRows(rows)
	mock.ExpectExec("INSERT INTO refresh_tokens").WithArgs(1, pgxmock.AnyArg(), pgxmock.AnyArg())
		mock.ExpectQuery("SELECT user_id, token, expires FROM refresh_tokens").WithArgs("incorrect_token").
		WillReturnRows(pgxmock.NewRows([]string{"user_id", "token", "expires"}))

	s := &TokenService{
		DB: mock,
	}
	correctJwt, correctRefresh, err := s.RefreshToken("correct_token")
	if err != nil {
		t.Fatalf("Expected error to be nil, got %v", err.Error())
	}
	valid, err := token.CheckJWTToken(correctJwt, "test secret")

	if err != nil {
		t.Fatalf("Expected error to be nil, got %v", err.Error())
	}

	if !valid {
		t.Fatalf("Expected JWT to be valid, got invalid")
	}

	if len(correctRefresh) != 128 {
		t.Errorf("TokenService.RefreshToken() refresh token length not 128, got %d", len(correctRefresh))
	}

	_, _, err = s.RefreshToken("incorrect_token")

	if err == nil {
		t.Fatalf("Expected error to be not nil, got nil")
	}

}
