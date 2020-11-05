package handlers

import (
	"fmt"
	"hearMeMail/global"
	"hearMeMail/services"
	"log"
	"net/http"
)

type LogoutHandler struct {
	config       *global.Config
	loginService *services.LoginService
}

func LogoutHandlerBuild(config *global.Config, loginService *services.LoginService) LogoutHandler {
	return LogoutHandler{
		config:       config,
		loginService: loginService,
	}
}

func (handler *LogoutHandler) Handler(rw http.ResponseWriter, req *http.Request) {
	if req == nil {
		writeResponseHeader(rw, http.StatusBadRequest, "json body containing 'username' and 'password' is required")
		return
	}
	if req.Method != http.MethodPost {
		message := fmt.Sprintf("Request method not available for this endpoint: %s", req.Method)
		writeResponseHeader(rw, http.StatusBadRequest, message)
		return
	}
	requestBody := LogoutRequestBody{}
	err := requestBody.decodeLogoutRequestBody(req)
	if err != nil {
		message := fmt.Sprintf("Failed to decode request body: err=%+v", err)
		writeResponseHeader(rw, http.StatusBadRequest, message)
		return
	}

	err = handler.loginService.Logout(requestBody.Username, requestBody.HeaderToken)
	if err != nil {
		message := fmt.Sprintf("%+v", err)
		log.Printf(message)
		writeResponseHeader(rw, http.StatusInternalServerError, message)
		return
	}
	writeResponseHeader(rw, http.StatusOK, "")
}
