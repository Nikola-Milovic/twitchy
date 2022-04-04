package response

type RefreshResponse struct {
	JWT          string `json:"jwt"`
	RefreshToken string `json:"refresh_token"`
}
