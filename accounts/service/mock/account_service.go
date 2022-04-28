package mock

import "nikolamilovic/twitchy/common/event"

type AccountServiceMock struct {
}

func (a *AccountServiceMock) CreateUser(ev event.AccountCreatedEvent) error {
	return nil
}
