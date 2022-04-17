package service

import (
	"context"
	"errors"
	"nikolamilovic/twitchy/auth/model"
	emitterMock "nikolamilovic/twitchy/auth/emitter/mock"
	serviceMock "nikolamilovic/twitchy/auth/service/mock"
	"testing"

	"github.com/pashagolub/pgxmock"
)

func TestRegistration(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	sut := &AuthService{
		DB:           mock,
		TokenService: &serviceMock.TokenServiceMock{},
		Emitter:      &emitterMock.AccountEmitterMock{},
	}

	// before we actually execute our api function, we need to expect required DB actions
	rows := pgxmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery("INSERT INTO users").WillReturnRows(rows)

	jwt, refresh, id, err := sut.Register("test@gmail.com", "123qwe")

	if jwt != "JWT" {
		t.Fatalf("Expected jwt to be %s got %s", "JWT", jwt)
	}

	if refresh != "REFRESH" {
		t.Fatalf("Expected refresh to be %s got %s", "REFRESH", refresh)
	}

	if err != nil {
		t.Fatalf("an error '%s' was not expected when creating auth", err)
	}

	if id != 1 {
		t.Fatalf("Expected id to be %d got %d", 1, id)
	}
}

func TestLoginCheck(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	sut := &AuthService{
		DB: mock,
	}

	hashedPassword, _ := hashPassword("password")

	// before we actually execute our api function, we need to expect required DB actions
	rows := pgxmock.NewRows([]string{"id", "password"}).AddRow(1, hashedPassword)

	mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)

	id, err := sut.checkLogin("test@gmail.com", "password")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when creating auth", err)
	}

	if id != 1 {
		t.Fatalf("Expected id to be %d got %d", 1, id)
	}
}

func TestLoginCheckWrongPassword(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	sut := &AuthService{
		DB: mock,
	}

	hashedPassword, _ := hashPassword("password")

	// before we actually execute our api function, we need to expect required DB actions
	rows := pgxmock.NewRows([]string{"id", "password"}).AddRow(1, hashedPassword)

	mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)

	id, err := sut.checkLogin("test@gmail.com", "wrongpassword")

	if id != -1 {
		t.Fatalf("Expected id to be %d got %d", -1, id)
	}

	if !errors.Is(err, model.WrongPasswordError) {
		t.Fatalf("wrong error , expected %e, got %e", model.WrongPasswordError, err)
	}
}
