package handlers

import (
	"fmt"
	"hearMeMail/global"
	"hearMeMail/services"
	"log"
	"net/http"
)

type EmailHandler struct {
	config       *global.Config
	emailService *services.EmailService
}

func EmailHandlerBuild(config *global.Config, emailService *services.EmailService) EmailHandler {
	return EmailHandler{config: config, emailService: emailService}
}

func (handler *EmailHandler) Handler(rw http.ResponseWriter, req *http.Request) {
	if req == nil {
		writeResponseHeader(rw, http.StatusBadRequest, "json body containing 'subject', 'email', 'body' is required")
		return
	}
	if req.Method != http.MethodPost {
		message := fmt.Sprintf("Request method not available for this endpoint: %s", req.Method)
		writeResponseHeader(rw, http.StatusBadRequest, message)
		return
	}
	emailRequest := EmailRequestBody{}
	err := emailRequest.decodeEmailRequestBody(req)
	if err != nil {
		message := fmt.Sprintf("Failed to decode request body: err=%+v", err)
		writeResponseHeader(rw, http.StatusBadRequest, message)
		return
	}
	err = emailRequest.validateEmailRequest()
	if err != nil {
		message := fmt.Sprintf("Request body not valid: err=%+v", err)
		writeResponseHeader(rw, http.StatusBadRequest, message)
		return
	}

	err = handler.emailService.SendEmail(services.SendEmailParameters{
		To:      emailRequest.Email,
		Subject: emailRequest.Subject,
		Body:    emailRequest.Body,
	})
	if err != nil {
		message := fmt.Sprintf("Failed to send email: err=%+v", err)
		log.Printf(message)
		writeResponseHeader(rw, http.StatusInternalServerError, message)
		return
	}
	writeResponseHeader(rw, http.StatusOK, "")
}
