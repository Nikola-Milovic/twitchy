package mock

import "nikolamilovic/twitchy/common/event"

type AccountServiceMock struct {
}

func (a *AccountServiceMock) CreateUser(ev event.AccountCreatedEventData) error {
	return nil
}
