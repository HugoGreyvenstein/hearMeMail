package services

import (
	"errors"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"hearMeMail/global"
	"hearMeMail/models"
	"hearMeMail/repositories"
	"net/http"
)

type EmailService struct {
	config             *global.Config
	client             *sendgrid.Client
	emailLogRepository *repositories.EmailLogRepository
	userLogRepository  *repositories.UserRepository
}

type SendEmailParameters struct {
	FromUsername string
	To           string
	Subject      string
	Body         string
}

func EmailServiceBuild(config *global.Config, emailLogRepository *repositories.EmailLogRepository, userLogRepository *repositories.UserRepository) *EmailService {
	service := &EmailService{
		config:             config,
		emailLogRepository: emailLogRepository,
		userLogRepository:  userLogRepository,
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
	user, err := service.userLogRepository.FindByUsername(params.FromUsername)
	if err != nil {
		return err
	}
	emailToLog := &models.EmailLog{
		Subject: params.Subject,
		Body:    params.Body,
		To:      params.To,
		Success: true,
		UserID:  user.ID,
	}
	if response.StatusCode != http.StatusAccepted {
		emailToLog.Success = false
		err := service.emailLogRepository.LogEmail(user, emailToLog)
		return errors.New(fmt.Sprintf("Failed to send email. Status code not `OK`: %+v, database-err=%+v", *response, err))
	}
	return service.emailLogRepository.LogEmail(user, emailToLog)
}
