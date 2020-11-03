package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	log.Print("Setting up handlers")
	http.HandleFunc("/email", emailHandler)
	log.Print("Handlers successfully set up")

	log.Print("Starting email server")
	err := http.ListenAndServe(":8080", nil)
	log.Printf("Error occurred while running server: err=%+v", err)
}

type EmailRequest struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func emailHandler(rw http.ResponseWriter, req *http.Request) {
	if req == nil {
		writeResponseHeader(rw, http.StatusBadRequest, "json body containing 'subject', 'email', 'body' is required")
		return
	}
	if req.Method != http.MethodPost {
		message := fmt.Sprintf("Request method not available for this endpoint: %s", req.Method)
		writeResponseHeader(rw, http.StatusBadRequest, message)
		return
	}
	emailRequest, err := decodeEmailRequestBody(req)
	if err != nil {
		message := fmt.Sprintf("Failed to decode request body: err=%+v", err)
		writeResponseHeader(rw, http.StatusBadRequest, message)
		return
	}
	err = emailRequest.validateEmailRequest()
	if err != nil {
		message := fmt.Sprintf("Request not valid: err=%+v", err)
		writeResponseHeader(rw, http.StatusBadRequest, message)
		return
	}

	// TODO Send email

	// TODO Change status code when email functionality implemented
	writeResponseHeader(rw, http.StatusNotImplemented, "Email sending functionality not implemented")
}

func (request EmailRequest) validateEmailRequest() error {
	if request.Email == "" {
		return errors.New("destination email is missing")
	}
	return nil
}

func decodeEmailRequestBody(req *http.Request) (*EmailRequest, error) {
	emailRequest := new(EmailRequest)
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(emailRequest)
	defer closeRequestBody(req.Body)
	if err == io.EOF {
		return nil, errors.New("Failed to decode body. Body is empty")
	}

	if err != nil {
		message := fmt.Sprintf("Failed to decode body: err=%+v", err)
		return nil, errors.New(message)
	}
	return emailRequest, err
}

func writeResponseHeader(rw http.ResponseWriter, statusCode int, message string) {
	log.Printf(message)
	rw.WriteHeader(statusCode)
	_, err := rw.Write([]byte(message))
	if err != nil {
		log.Printf("Failed to write to request writer: err=%+v", err)
	}
}

func closeRequestBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		log.Printf("Failed to close request body: err=%+v", err)
	}
}
