package service

import (
	"context"
	"nikolamilovic/twitchy/common/event"
	"testing"

	"github.com/pashagolub/pgxmock"
)

func TestUserCreation(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	sut := &AccountService{
		DB: mock,
	}

	// before we actually execute our api function, we need to expect required DB actions
	rows := pgxmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery("INSERT INTO users").WithArgs(1, "email@gmail.com", "username").WillReturnRows(rows)

	err = sut.CreateUser(event.AccountCreatedEventData{
		ID:       1,
		Email:    "email@gmail.com",
		Username: "username",
	})

	if err != nil {
		t.Fatalf("an error '%s' was not expected when creating user", err)
	}
}
