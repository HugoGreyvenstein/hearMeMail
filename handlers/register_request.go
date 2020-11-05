package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type RegisterRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (registerRequest RegisterRequestBody) validateRegisterRequest() error {
	if registerRequest.Username == "" {
		return errors.New("username is missing")
	}
	if registerRequest.Password == "" {
		return errors.New("password is missing")
	}
	return nil
}

func (registerRequest *RegisterRequestBody) decodeRegisterRequestBody(req *http.Request) error {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(registerRequest)
	defer closeRequestBody(req.Body)
	if err == io.EOF {
		return errors.New("failed to decode body. Body is empty")
	}
	if err != nil {
		return err
	}
	return err
}
