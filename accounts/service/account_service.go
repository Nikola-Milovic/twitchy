package service

import (
	"context"
	"fmt"
	"nikolamilovic/twitchy/accounts/db"
	"nikolamilovic/twitchy/accounts/model/event"
)

type IAccountService interface {
	CreateUser(ev event.CreateAccountEvent) error
}

type AccountService struct {
	DB db.PgxIface
}

func (s *AccountService) CreateUser(ev event.CreateAccountEvent) error {
	rows, err := s.DB.Query(context.Background(), "INSERT INTO users (id, email) VALUES ($1,$2)", ev.UserId, ev.Email)

	if err != nil {
		return fmt.Errorf("CreateUser %w", err)
	}

	defer rows.Close()

	return nil
}
