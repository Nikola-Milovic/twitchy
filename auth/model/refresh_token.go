package model

type RefreshToken struct {
	Token   string `json:"token"`
	UserId  int    `json:"user_id"`
	Expires int64   `json:"expires"`
	ID 		int    `json:"id"`	
}
