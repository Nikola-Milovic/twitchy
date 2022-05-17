package mock

type AuthServiceMock struct {
}

func (a *AuthServiceMock) Register(email, password, username string) (string, string, int, error) {
	return "JWT", "REFRESH", 1, nil
}

func (a *AuthServiceMock) Login(email, password string) (string, string, int, error) {
	return "JWT", "REFRESH", 1, nil
}
