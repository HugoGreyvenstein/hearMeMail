package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

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

func isMethodAvailable(requestMethod string, availableMethods ...string) error {
	for _, method := range availableMethods {
		if requestMethod == method {
			return nil
		}
	}
	message := fmt.Sprintf("method not available: requested=%s, available=%s", requestMethod, availableMethods)
	return errors.New(message)
}
