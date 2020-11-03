package services

import (
	"errors"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"hearMeMail/global"
	"log"
	"net/http"
)

type EmailService struct {
	config *global.Config
	client *sendgrid.Client
}

type SendEmailParameters struct {
	To      string
	Subject string
	Body    string
}

func EmailServiceBuild(config *global.Config) *EmailService {
	service := &EmailService{
		config: config,
	}
	service.client = sendgrid.NewSendClient(config.Email.ApiKey)
	return service
}

func (service EmailService) SendEmail(params SendEmailParameters) error {
	from := mail.NewEmail(service.config.Email.Name, service.config.Email.From)
	to := mail.NewEmail("", params.To)
	message := mail.NewSingleEmail(from, params.Subject, to, params.Body, "")
	response, err := service.client.Send(message)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusAccepted {
		return errors.New(fmt.Sprintf("Failed to send email. Status code not `OK`: %+v", *response))
	}
	log.Printf("response: %+v", response)
	return nil
}
