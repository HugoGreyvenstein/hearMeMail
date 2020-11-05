package handlers

import (
	"io"
	"log"
	"net/http"
)

func writeResponseHeader(rw http.ResponseWriter, statusCode int, body string) {
	rw.WriteHeader(statusCode)
	_, err := rw.Write([]byte(body))
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
