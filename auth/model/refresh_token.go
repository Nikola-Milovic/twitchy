package model

type RefreshToken struct {
	Token   string `json:"token"`
	UserId  int    `json:"user_id"`
	Expires uint   `json:"expires"`
	ID 		int    `json:"id"`	
}
