package handlers

import (
	"errors"
	"net/http"
)

type LogoutRequestBody struct {
	Username    string
	HeaderToken string
}

func (logoutRequest *LogoutRequestBody) decodeLogoutRequestBody(req *http.Request) error {
	username := req.Header["Username"]
	token := req.Header["Token"]
	if username == nil || token == nil || len(username) < 1 || len(token) < 1 {
		return errors.New("username and token not valid")
	}
	logoutRequest.Username = username[0]
	logoutRequest.HeaderToken = token[0]
	return nil
}
