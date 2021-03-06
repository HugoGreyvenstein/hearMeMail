package handlers

import (
	"encoding/json"
	"fmt"
	"hearMeMail/global"
	"hearMeMail/repositories"
	"hearMeMail/services"
	"log"
	"net/http"
)

type LoginHandler struct {
	config       *global.Config
	loginService *services.LoginService
}

func LoginHandlerBuild(config *global.Config, loginService *services.LoginService) LoginHandler {
	return LoginHandler{
		config:       config,
		loginService: loginService,
	}
}

type LoginResponse struct {
	UserId      string `json:"userId"`
	HeaderToken string `json:"headerToken"`
}

func (handler *LoginHandler) Handler(rw http.ResponseWriter, req *http.Request) {
	if req == nil {
		writeResponseHeader(rw, http.StatusBadRequest, "json body containing 'username' and 'password' is required")
		return
	}
	if req.Method != http.MethodPost {
		message := fmt.Sprintf("Request method not available for this endpoint: %s", req.Method)
		writeResponseHeader(rw, http.StatusBadRequest, message)
		return
	}
	requestBody := LoginRequestBody{}
	err := requestBody.decodeLoginRequestBody(req)
	if err != nil {
		message := fmt.Sprintf("Failed to decode request body: err=%+v", err)
		writeResponseHeader(rw, http.StatusBadRequest, message)
		return
	}

	err = requestBody.validateLoginRequest()
	if err != nil {
		message := fmt.Sprintf("Missing arguments: err=%+v", err)
		writeResponseHeader(rw, http.StatusBadRequest, message)
		return
	}

	result, err := handler.loginService.Login(requestBody.Username, requestBody.Password)
	if err == services.LoginCredentialsIncorrect {
		writeResponseHeader(rw, http.StatusUnauthorized, "Username or password incorrect")
		return
	}
	if err != nil {
		message := fmt.Sprintf("%+v", err)
		log.Printf(message)
		writeResponseHeader(rw, http.StatusInternalServerError, message)
		return
	}

	body, err := json.Marshal(&LoginResponse{
		UserId:      result.UserId,
		HeaderToken: result.HeaderToken,
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

func (handler *LoginHandler) TokenChecker(handlerFunc func(rw http.ResponseWriter, req *http.Request)) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		username := req.Header[headerUsername]
		token := req.Header[headerToken]
		if username == nil || token == nil || len(username) < 1 || len(token) < 1 {
			writeResponseHeader(rw, http.StatusUnauthorized, "Login token not valid")
			return
		}
		isValid, err := handler.loginService.TokenValid(username[0], token[0])
		if err != nil && err != repositories.ErrNotFound {
			log.Printf("Could not determine if token was valid: username=%s, token=%s, err=%+v", username, token, err)
			writeResponseHeader(rw, http.StatusInternalServerError, "Could not determine if token was valid")
			return
		}
		if !isValid {
			writeResponseHeader(rw, http.StatusUnauthorized, "Login token not valid")
			return
		}
		handlerFunc(rw, req)
	}
}
