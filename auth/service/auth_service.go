package service

import (
	"context"
	"fmt"
	"nikolamilovic/twitchy/auth/client"
	"nikolamilovic/twitchy/auth/model"
	"nikolamilovic/twitchy/auth/model/event"
	db "nikolamilovic/twitchy/common/db"

	"golang.org/x/crypto/bcrypt"
)

// NOTE: we could and probably should decouple the DB logic from service and have a repository layer as well, but adding
// another layer on top of this simple example is just a bit overkill for the time being.

type IAuthService interface {
	Register(email, password string) (string, string, int, error)
	Login(email, password string) (string, string, int, error)
}

type AuthService struct {
	DB                  db.PgxIface
	TokenService        ITokenService
	AccountRabbitClient client.IAccountClient
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//Return JWT, refresh token and the user ID
func (s *AuthService) Register(email, password string) (string, string, int, error) {
	id, err := s.createUser(email, password)

	if err != nil {
		return "", "", -1, fmt.Errorf("Register create user %w", err)
	}

	//TODO add events to a workbox in DB to be eventually emitted
	// currently this poses an issue if the event emittion fails but we successfuly created a user
	// so the user and event should be saved to the DB in a transaction

	err = s.AccountRabbitClient.PublishAccountCreatedEvent(event.AccountCreatedEvent{Id: id})

	if err != nil {
		fmt.Printf("Error emitting account created event %s", err.Error())
		return "", " ", -1, err
	}

	jwt, refresh, err := s.TokenService.GenerateNewTokensForUser(id)

	if err != nil {
		return "", "", -1, fmt.Errorf("Register generate tokens %w", err)
	}

	//Emit user created event
	return jwt, refresh, id, nil
}

//Check checkLogin first, then if it's ok, generate tokens and return JWT, refresh token and the user ID
func (a *AuthService) Login(email, password string) (string, string, int, error) {
	id, err := a.checkLogin(email, password)

	if err != nil {
		return "", "", -1, fmt.Errorf("Login check %w", err)
	}

	jwt, refresh, err := a.TokenService.GenerateNewTokensForUser(id)

	if err != nil {
		return "", "", -1, fmt.Errorf("Login generate new tokens %w", err)
	}

	return jwt, refresh, id, nil
}

func (a *AuthService) checkLogin(email, password string) (int, error) {
	rows, err := a.DB.Query(context.Background(), "SELECT id, password FROM users WHERE email=$1", email)

	if err != nil {
		fmt.Printf("Error getting user %s", err.Error())
		return -1, fmt.Errorf("CheckLogin: %w", err)
	}

	defer rows.Close()

	var id = -1
	var hashedPassword string
	if rows.Next() {
		err = rows.Scan(&id, &hashedPassword)

		if err != nil {
			fmt.Printf("Error scanning login row %s", err.Error())
			return -1, fmt.Errorf("CheckLogin: %w", err)
		}
	}

	if checkPasswordHash(password, hashedPassword) {
		return id, nil
	} else {
		return -1, fmt.Errorf("CheckLogin: %w", model.WrongPasswordError)
	}
}

func (s *AuthService) createUser(email, password string) (int, error) {
	fmt.Printf("Creating user with %s email and %s password\n", email, password)

	hashedPassword, err := hashPassword(password)

	if err != nil {
		fmt.Printf("Error hashing password %s \n", err.Error())
		return -1, err
	}

	rows, err := s.DB.Query(context.Background(), "INSERT INTO users (email, password) VALUES ($1,$2) RETURNING id", email, hashedPassword)

	if err != nil {
		fmt.Printf("Error creating user %s", err.Error())
		return -1, err
	}

	defer rows.Close()

	var id = -1
	if rows.Next() {
		rows.Scan(&id)
	}

	return id, nil
}
