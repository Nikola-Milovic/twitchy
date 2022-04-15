package mock

import "nikolamilovic/twitchy/accounts/model/event"

type AccountServiceMock struct {
}

func (a *AccountServiceMock) CreateUser(ev event.CreateAccountEvent) error {
	return nil
}
