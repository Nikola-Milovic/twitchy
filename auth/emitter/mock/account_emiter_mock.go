package mock

import "nikolamilovic/twitchy/auth/model"

type AccountEmitterMock struct {}

func (e *AccountEmitterMock) Emit(event model.AccountCreatedEvent) error {
	return nil
}