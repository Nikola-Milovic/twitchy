package event

type AccountCreatedEvent struct {
	ID    int  `json:"id"`
	Email string `json:"email"`
}

type AccountCreatedAck struct {
	ID      int64  `json:"id"`
	Service string `json:"service"`
}
