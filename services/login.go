package services

import (
	"errors"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"golang.org/x/crypto/bcrypt"
	"hearMeMail/global"
	"hearMeMail/repositories"
	"log"
	"strconv"
	"time"
)

type LoginService struct {
	config         *global.Config
	client         *sendgrid.Client
	userRepository *repositories.UserRepository
}

func LoginServiceBuild(config *global.Config, userRepository *repositories.UserRepository) *LoginService {
	return &LoginService{
		config:         config,
		userRepository: userRepository,
	}
}

type LoginResult struct {
	UserId       string    `json:"userId"`
	HeaderToken  string    `json:"headerToken"`
	HeaderExpiry time.Time `json:"headerExpiry"`
}

var LoginCredentialsIncorrect = errors.New("username or password incorrect")

func (service *LoginService) Login(username string, password string) (*LoginResult, error) {
	user, err := service.userRepository.FindByUsername(username)
	if err == repositories.ErrNotFound {
		return nil, LoginCredentialsIncorrect
	}
	if err != nil {
		errMessage := repositories.ErrDatabaseError.Error() + fmt.Sprintf(": %+v", err)
		return nil, errors.New(errMessage)
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, LoginCredentialsIncorrect
	}
	if err != nil {
		log.Printf("Password hash error: err=%+v", err)
		return nil, err
	}
	unixNano := time.Now().UnixNano()
	headerTokenSource := username + strconv.FormatInt(unixNano, 10)
	headerToken, err := bcrypt.GenerateFromPassword([]byte(headerTokenSource),
		service.config.Bcrypt.Cost)

	user, err = service.userRepository.UpdateHeaderToken(user, headerToken, time.Now().Add(24*time.Hour))
	if err != nil {
		return nil, err
	}
	return &LoginResult{
		UserId:       username,
		HeaderToken:  string(headerToken),
		HeaderExpiry: *user.HeaderExpiry,
	}, nil
}

func (service *LoginService) Logout(username string) error {
	err := service.userRepository.DeleteHeaderToken(username)
	if err != nil {
		return err
	}
	return nil
}

func (service *LoginService) TokenValid(username string, token string) (bool, error) {
	user, err := service.userRepository.FindByUsername(username)
	if err == repositories.ErrNotFound {
		log.Printf("User not found: %s", username)
		return false, err
	}
	if err != nil {
		return false, err
	}
	hasTokenExpired := token == string(user.HeaderToken) && time.Now().Before(*user.HeaderExpiry)
	return hasTokenExpired, nil
}
