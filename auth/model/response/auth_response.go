package response

type AuthResponse struct {
	ID           int    `json:"id"`
	JWT          string `json:"jwt"`
	RefreshToken string `json:"refresh_token"`
}
