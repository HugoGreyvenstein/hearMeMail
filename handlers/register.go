package handlers

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"hearMeMail/global"
	"hearMeMail/models"
	"hearMeMail/repositories"
	"log"
	"net/http"
)

type RegisterHandler struct {
	config         *global.Config
	userRepository *repositories.UserRepository
}

func RegisterHandlerBuild(config *global.Config, userRepository *repositories.UserRepository) RegisterHandler {
	return RegisterHandler{
		config:         config,
		userRepository: userRepository,
	}
}

type RegisterResponse struct {
	Username string `json:"username"`
}

func (handler *RegisterHandler) Handler(rw http.ResponseWriter, req *http.Request) {
	if req == nil {
		writeResponseHeader(rw, http.StatusBadRequest, "json body containing 'username' and 'password' is required")
		return
	}
	if req.Method != http.MethodPost {
		message := fmt.Sprintf("Request method not available for this endpoint: %s", req.Method)
		writeResponseHeader(rw, http.StatusBadRequest, message)
		return
	}
	requestBody := RegisterRequestBody{}
	err := requestBody.decodeRegisterRequestBody(req)
	if err != nil {
		message := fmt.Sprintf("Failed to decode request body: err=%+v", err)
		writeResponseHeader(rw, http.StatusBadRequest, message)
		return
	}
	err = requestBody.validateRegisterRequest()
	if err != nil {
		message := fmt.Sprintf("Missing arguments: err=%+v", err)
		writeResponseHeader(rw, http.StatusBadRequest, message)
		return
	}
	existingUser, err := handler.userRepository.FindByUsername(requestBody.Username)
	if existingUser != nil {
		message := fmt.Sprintf("User with username already exists: err=%+v", err)
		writeResponseHeader(rw, http.StatusUnauthorized, message)
		return
	}
	if err != nil && err != repositories.ErrNotFound {
		message := fmt.Sprintf("Database error: err=%+v", err)
		writeResponseHeader(rw, http.StatusInternalServerError, message)
		return
	}
	newUser := &models.User{}
	newUser.Username = requestBody.Username
	newUser.Password, _ = bcrypt.GenerateFromPassword([]byte(requestBody.Password), handler.config.Bcrypt.Cost)
	user, err := handler.userRepository.Insert(newUser)
	if err != nil {
		message := fmt.Sprintf("Failed to insert new user: %+v", err)
		log.Printf(message)
		writeResponseHeader(rw, http.StatusInternalServerError, message)
		return
	}
	body, err := json.Marshal(&RegisterResponse{
		Username: user.Username,
	})
	if err != nil {
		message := fmt.Sprintf("%+v", err)
		log.Printf(message)
		writeResponseHeader(rw, http.StatusInternalServerError, message)
		return
	}
	_, err = rw.Write(body)
	if err != nil {
		message := fmt.Sprintf("Unable to write body to response: err=%+v", err)
		log.Printf(message)
		writeResponseHeader(rw, http.StatusInternalServerError, message)
		return
	}
}
