package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type EmailRequestBody struct {
	FromUsername string `json:"-"`
	Email        string `json:"email"`
	Subject      string `json:"subject"`
	Body         string `json:"body"`
}

func (emailRequest EmailRequestBody) validateEmailRequest() error {
	if emailRequest.Email == "" {
		return errors.New("destination email is missing")
	}
	return nil
}

func (emailRequest *EmailRequestBody) decodeEmailRequestBody(req *http.Request) error {
	if req.Header != nil {
		username := req.Header[headerUsername]
		if len(username) > 0 {
			emailRequest.FromUsername = username[0]
		}
	}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(emailRequest)
	defer closeRequestBody(req.Body)
	if err == io.EOF {
		return errors.New("failed to decode body. Body is empty")
	}

	if err != nil {
		return err
	}
	return err
}
