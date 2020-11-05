package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type LoginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (loginRequest LoginRequestBody) validateLoginRequest() error {
	if loginRequest.Username == "" {
		return errors.New("username is missing")
	}
	if loginRequest.Password == "" {
		return errors.New("password is missing")
	}
	return nil
}

func (loginRequest *LoginRequestBody) decodeLoginRequestBody(req *http.Request) error {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(loginRequest)
	defer closeRequestBody(req.Body)
	if err == io.EOF {
		return errors.New("failed to decode body. Body is empty")
	}

	if err != nil {
		return err
	}
	return err
}
