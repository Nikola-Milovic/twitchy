package event

const (
	AccountCreatedType    = "account_created"
	AccountCreatedAckType = "account_created_ack"
)

type AccountCreatedEventData struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type AccountCreatedAckData struct {
	ID      int64  `json:"id"`
	Service string `json:"service"`
}
