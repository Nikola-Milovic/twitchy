package handler

import (
	"bytes"
	"io/ioutil"
	"nikolamilovic/twitchy/accounts/service/mock"
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
	req := httptest.NewRequest(http.MethodPost, "/test", reader)
	req.Header.Set("Content-Type", "application/json")

	srv := &AuthHandler{}
	srv.accountService = &mock.AccountServiceMock{}
	srv.validator = validator.New()
	srv.Routes()

	resp, err := srv.Router.Test(req)

	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	// We should get a good status code
	if want, got := http.StatusOK, resp.StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}

	expectedResponse := []byte("hello there")

	if got, want := data, expectedResponse; !bytes.Equal(got, want) {
		t.Fatalf("expected a %v, instead got: %v", want, got)
	}
}
