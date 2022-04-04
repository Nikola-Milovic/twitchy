package mock

type TokenServiceMock struct {
}

func (s *TokenServiceMock) RefreshToken(refreshTokenString string) (string, string, error) {
	return "JWT", "REFRESH", nil
}

//Returns JWT, RefreshToken, error
func (s *TokenServiceMock) GenerateNewTokensForUser(userId int) (string, string, error) {
	return "JWT", "REFRESH", nil
}
