package event

type CreateAccountEvent struct {
	UserId int64 `json:"user_id"`
	Email string `json:"email"`
}
