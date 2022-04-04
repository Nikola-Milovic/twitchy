package handler

import (
	"encoding/json"

	"io/ioutil"
	"nikolamilovic/twitchy/auth/model/response"
	"nikolamilovic/twitchy/auth/service/mock"
	"strings"

	// "net/http"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestRegistration(t *testing.T) {
	reader := strings.NewReader(`{
		"email":"test@gmail.com",
		"password":"123qwe123"
	 }`)
	req := httptest.NewRequest(http.MethodPost, "/register", reader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv := &AuthHandler{}
	srv.authService = &mock.AuthServiceMock{}
	srv.validator = validator.New()

	srv.handleRegistration()(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	// We should get a good status code
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}

	var responseData response.AuthResponse

	json.Unmarshal(data, &responseData)

	expectedResponse := response.AuthResponse{
		ID:           1,
		JWT:          "JWT",
		RefreshToken: "REFRESH",
	}

	if got, want := responseData, expectedResponse; want != got {
		t.Fatalf("expected a %v, instead got: %v", want, got)
	}

}
func TestRegistrationErrors(t *testing.T) {

	type registrationTest struct {
		description    string
		input          string
		expected       string
		expectedStatus int
	}

	for _, scenario := range []registrationTest{
		{
			description: "invalid email",
			input: `{
				"email":"invalid",
				"password":"123qwe123"
			 }`,
			expected:       "Key: 'RegistrationRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag\n",
			expectedStatus: http.StatusBadRequest,
		},
		{
			description: "missing password",
			input: `{
				"email":"valid@gmail.com",
				"password":"123"
			 }`,
			expected:       "Key: 'RegistrationRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag\n",
			expectedStatus: http.StatusBadRequest,
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			//GIVEN
			reader := strings.NewReader(scenario.input)

			req := httptest.NewRequest(http.MethodPost, "/register", reader)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			//INITIALIZING SERVER
			srv := AuthHandler{}
			srv.authService = &mock.AuthServiceMock{}
			srv.Routes()
			srv.validator = validator.New()

			srv.handleRegistration()(w, req)

			//SHOULD
			res := w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			// We should get a good status code
			if want, got := scenario.expectedStatus, w.Result().StatusCode; want != got {
				t.Fatalf("expected a %d, instead got: %d", want, got)
			}

			if want, got := scenario.expected, string(data); strings.Compare(want, got) != 0 {
				t.Fatalf("expected a %s, instead got: %s", want, got)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	reader := strings.NewReader(`{
		"email":"test@gmail.com",
		"password":"123qwe1"
	 }`)
	req := httptest.NewRequest(http.MethodPost, "/register", reader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv := &AuthHandler{}
	srv.authService = &mock.AuthServiceMock{}
	srv.validator = validator.New()

	srv.handleLogin()(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	// We should get a good status code
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}

	var responseData response.AuthResponse

	json.Unmarshal(data, &responseData)

	expectedResponse := response.AuthResponse{
		ID:           1,
		JWT:          "JWT",
		RefreshToken: "REFRESH",
	}

	if got, want := responseData, expectedResponse; want != got {
		t.Fatalf("expected a %v, instead got: %v", want, got)
	}
}

func TestRefresh(t *testing.T) {
	reader := strings.NewReader(`{
		"refresh_token":"REFRESH"
	 }`)
	req := httptest.NewRequest(http.MethodPost, "/refresh", reader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv := &AuthHandler{}
	srv.authService = &mock.AuthServiceMock{}
	srv.tokenService = &mock.TokenServiceMock{}
	srv.validator = validator.New()

	srv.handleRefresh()(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	// We should get a good status code
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}

	var responseData response.RefreshResponse

	json.Unmarshal(data, &responseData)

	expectedResponse := response.RefreshResponse{
		JWT:          "JWT",
		RefreshToken: "REFRESH",
	}

	if got, want := responseData, expectedResponse; want != got {
		t.Fatalf("expected a %v, instead got: %v", want, got)
	}
}
