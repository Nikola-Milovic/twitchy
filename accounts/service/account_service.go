package service

import (
	"context"
	"fmt"
	db "nikolamilovic/twitchy/common/db"
	event "nikolamilovic/twitchy/common/event"
)

type IAccountService interface {
	CreateUser(ev event.AccountCreatedEventData) error
}

type AccountService struct {
	DB db.PgxIface
}

func NewAccountService(db db.PgxIface) IAccountService {
	return &AccountService{
		DB: db,
	}
}

func (s *AccountService) CreateUser(ev event.AccountCreatedEventData) error {
	rows, err := s.DB.Query(context.Background(), "INSERT INTO users (id, email) VALUES ($1,$2)", ev.ID, ev.Email)

	if err != nil {
		return fmt.Errorf("CreateUser %w", err)
	}

	defer rows.Close()

	return nil
}
