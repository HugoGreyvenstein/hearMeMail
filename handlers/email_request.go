package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type EmailRequest struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (emailRequest EmailRequest) validateEmailRequest() error {
	if emailRequest.Email == "" {
		return errors.New("destination email is missing")
	}
	return nil
}

func (emailRequest *EmailRequest) decodeEmailRequestBody(req *http.Request) error {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(emailRequest)
	defer closeRequestBody(req.Body)
	if err == io.EOF {
		return errors.New("Failed to decode body. Body is empty")
	}

	if err != nil {
		return err
	}
	return err
}
